package helper

import (
	"database/sql"
	"reflect"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/tranhanh18042/e-comere/services/pkg/metrics"
)

type SvcDB struct {
	*sqlx.DB
	Name string
}

func NewDBConn(dbName, connStr string) (*SvcDB, error) {
	conn, err := sqlx.Open("mysql", connStr)
	if err != nil {
		return nil, err
	}

	return &SvcDB{
		DB: conn,
		Name: dbName,
	}, nil
}

type DBMetricsLabels struct {
	DBName string
	Target string
}

func DBGetWithMetrics(labels DBMetricsLabels, db *SvcDB, dest interface{}, query string, args ...interface{}) (err error) {
	start := time.Now().UnixMilli()
	defer func(){
		if err != nil {
			if err == sql.ErrNoRows {
				metricQueryErrCnt(labels, MetricDBErrNoRow)
			} else {
				metricQueryErrCnt(labels, MetricDBErrInternal)
			}
		}

		if err == nil { // having record
			metricQueryRowsReturned(labels, 1)
		} else if err == sql.ErrNoRows { // no record found
			metricQueryRowsReturned(labels, 0)
		}

		metricQueryDur(labels, start)
	}()

	return db.Get(dest, query, args...)
}

func DBSelectWithMetrics(labels DBMetricsLabels, db *sqlx.DB, dest interface{}, query string, args ...interface{}) (err error) {
	start := time.Now().UnixMilli()
	defer func(){
		if err != nil {
			if err == sql.ErrNoRows {
				metricQueryErrCnt(labels, MetricDBErrNoRow)
			} else {
				metricQueryErrCnt(labels, MetricDBErrInternal)
			}
		}

		if err == nil { // having record
			v := reflect.ValueOf(dest)
			if v.Kind() == reflect.Pointer {
				v = v.Elem()
			}
			metricQueryRowsReturned(labels, v.Len())
		} else if err == sql.ErrNoRows { // no record found
			metricQueryRowsReturned(labels, 0)
		}

		metricQueryDur(labels, start)
	}()

	return db.Select(dest, query, args...)
}

func DBExecWithMetrics(labels DBMetricsLabels, db *SvcDB, query string, args ...interface{}) (result sql.Result, err error) {
	start := time.Now().UnixMilli()
	defer func(){
		if err != nil {
			metricQueryErrCnt(labels, MetricDBErrInternal)
		}
		metricQueryDur(labels, start)
	}()

	return db.Exec(query, args...)
}

func metricQueryDur(labels DBMetricsLabels, start int64) {
	end := time.Now().UnixMilli()
	metrics.DB.QueryDur.With(prometheus.Labels{
		"env": "local",
		"db": labels.DBName,
		"target": labels.Target,
	}).Observe(float64(end - start))
}

func metricQueryErrCnt(labels DBMetricsLabels, target string) {
	metrics.DB.ErrCnt.With(prometheus.Labels{
		"env": "local",
		"type": target,
		"db": labels.DBName,
		"target": labels.Target,
	}).Inc()
}

func metricQueryRowsReturned(labels DBMetricsLabels, numRows int) {
	metrics.DB.RowsReturned.With(prometheus.Labels{
		"env": "local",
		"db": labels.DBName,
		"target": labels.Target,
	}).Set(float64(numRows))
}
