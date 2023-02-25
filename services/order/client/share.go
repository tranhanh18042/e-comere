package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/tranhanh18042/e-comere/services/helper"
	"github.com/tranhanh18042/e-comere/services/pkg/logger"
	"github.com/tranhanh18042/e-comere/services/pkg/metrics"
)

type retryDialer struct {
	*net.Dialer
	options retryDialerOptions
}

type retryDialerOptions struct {
	retryCount int
}

func newRetryDialer(dialer *net.Dialer, options retryDialerOptions) *retryDialer {
	if options.retryCount == 0 {
		options.retryCount = 2
	}

	return &retryDialer{dialer, options}
}

func (d *retryDialer) DialContext(ctx context.Context, network, address string) (net.Conn, error) {
	conn, err := d.Dialer.DialContext(ctx, network, address)
	for i := 0; err != nil && shouldRetryDialErr(err) && i < d.options.retryCount && ctx.Err() == nil; i++ {
		metrics.HTTPClient.ReconnectCnt.With(prometheus.Labels{
			"env":     "local",
			"network": network,
			"address": address,
		}).Inc()
		logger.Info(ctx, "Retrying during establishing tcp connection", i, "local", network, address)
		conn, err = d.Dialer.DialContext(ctx, network, address)
	}

	return conn, err
}

func shouldRetryDialErr(err error) bool {
	return err != nil &&
		(strings.HasSuffix(err.Error(), "i/o timeout") ||
			strings.HasSuffix(err.Error(), "connect: connection refused"))
}

var (
	httpClient = &http.Client{
		Transport: &http.Transport{
			MaxIdleConnsPerHost:   10_000,
			MaxConnsPerHost:       10_000,
			IdleConnTimeout:       30 * time.Second,
			ResponseHeaderTimeout: 3 * time.Second,

			DialContext: newRetryDialer(
				&net.Dialer{Timeout: time.Second},
				retryDialerOptions{
					retryCount: 2,
				},
			).DialContext,
		},
		Timeout: 5 * time.Second,
	}
)

type responseBody struct {
	Payload any `json:"payload"`
}

func post(targetSvc, path string, client *http.Client, url string, reqBody any, resPayload any) (err error) {
	var res *http.Response
	start := time.Now().UnixMilli()
	var errType string

	defer func ()  {
		if err != nil {
			metrics.HTTPClient.ErrCnt.With(prometheus.Labels{
				"svc": "order",
				"env": "local",
				"type": errType,
				"target_svc": targetSvc,
				"path": path,
			}).Inc()
		}
		if res != nil {
			metrics.HTTPClient.ResponseCode.With(prometheus.Labels{
				"svc": "order",
				"env": "local",
				"target_svc": targetSvc,
				"path": path,
				"code": strconv.Itoa(res.StatusCode),
			}).Inc()

			end := time.Now().UnixMilli()
			metrics.HTTPClient.ResponseDur.With(prometheus.Labels{
				"svc": "order",
				"env": "local",
				"target_svc": targetSvc,
				"path": path,
			}).Observe(float64(end - start))
		}
	}()

	reqBodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		err = fmt.Errorf("cannot marshal request body, err: %v", err)
		errType = helper.MetricMarshalReqError
		return
	}

	request, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(reqBodyBytes))
	if err != nil {
		err = fmt.Errorf("cannot init post request, err: %v", err)
		errType = helper.MetricInitReqError
		return
	}

	request.Header.Set("Content-Type", "application/json; charset=UTF-8")
	res, err = client.Do(request)
	if err != nil {
		err = fmt.Errorf("cannot post request, err: %v", err)
		errType = helper.MetricPostReqError
		return
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		err = fmt.Errorf("non-200 status returned, got: %d", res.StatusCode)
		errType = helper.MetricNon200Error
		return
	}

	resBodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		err = fmt.Errorf("cannot read response body, err: %v", err)
		errType = helper.MetricReadBodyError
		return
	}
	logger.Debug(context.Background(), "post response body", string(resBodyBytes))
	resBody := &responseBody{
		Payload: resPayload,
	}
	err = json.Unmarshal(resBodyBytes, resBody)
	if err != nil {
		err = fmt.Errorf("cannot unmarshal response body to struct, err: %v", err)
		errType = helper.MetricUnmarshalBodyError
		return
	}

	return
}

func get(targetSvc, path string, client *http.Client, url string, resPayload any) (err error) {
	var res *http.Response
	start := time.Now().UnixMilli()
	var errType string

	defer func ()  {
		if err != nil {
			metrics.HTTPClient.ErrCnt.With(prometheus.Labels{
				"svc": "order",
				"env": "local",
				"type": errType,
				"target_svc": targetSvc,
				"path": path,
			}).Inc()
		}
		if res != nil {
			metrics.HTTPClient.ResponseCode.With(prometheus.Labels{
				"svc": "order",
				"env": "local",
				"target_svc": targetSvc,
				"path": path,
				"code": strconv.Itoa(res.StatusCode),
			}).Inc()

			end := time.Now().UnixMilli()
			metrics.HTTPClient.ResponseDur.With(prometheus.Labels{
				"svc": "order",
				"env": "local",
				"target_svc": targetSvc,
				"path": path,
			}).Observe(float64(end - start))
		}
	}()

	res, err =  client.Get(url)

	if err != nil {
		err = fmt.Errorf("cannot get request, err: %v", err)
		errType = helper.MetricGetReqError
		return
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		err = fmt.Errorf("non-200 status returned, got: %d", res.StatusCode)
		errType = helper.MetricNon200Error
		return
	}

	resBodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		err = fmt.Errorf("cannot read response body, err: %v", err)
		errType = helper.MetricReadBodyError
		return
	}
	resBody := &responseBody{
		Payload: resPayload,
	}
	err = json.Unmarshal(resBodyBytes, resBody)
	if err != nil {
		err = fmt.Errorf("cannot unmarshal response body to struct, err: %v", err)
		errType = helper.MetricMarshalReqError
		return
	}

	return
}

type getParams struct {
	queryParams map[string]string
	pathParams  map[string]string
}

// formatURL returns a path with params on path and query.
// path param must have same name with the key on params.pathParams
// Eg: input:
//
//	 `path` = /abc/:def/info
//	 `params` = {
//		  queryParams = {"t": "123"},
//		  pathParams = {"def": "456", "m": "789"}
//	 }
//
// then the output should be: "/abc/456/info?t=123"
// (since "m" is not existing as path param on `path`, so it will be omitted)
func (p *getParams) formatURL(path string) string {
	if p == nil {
		return path
	}

	if len(p.pathParams) > 0 {
		for k, v := range p.pathParams {
			path = strings.ReplaceAll(path, ":"+k, v)
		}
	}

	q := url.Values{}
	if len(p.queryParams) > 0 {
		for k, v := range p.queryParams {
			q.Add(k, v)
		}
	}

	if query := q.Encode(); len(query) > 0 {
		path += "?" + query
	}

	return path
}
