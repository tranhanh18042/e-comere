package metrics

import (
	"database/sql"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type httpMetrics struct {
	ReconnectCnt *prometheus.CounterVec
	ErrCnt       *prometheus.CounterVec
	ResponseCode *prometheus.CounterVec
	ResponseDur  *prometheus.HistogramVec
}

type apiMetrics struct {
	ErrCnt *prometheus.CounterVec
	ReqCnt *prometheus.CounterVec
	ReqDur *prometheus.HistogramVec
}

type dbMetrics struct {
	ReconnectCnt *prometheus.CounterVec
	ErrCnt       *prometheus.CounterVec
	RowsReturned *prometheus.GaugeVec
	QueryDur     *prometheus.HistogramVec
	Stats        *dbStater
}

type dbStater struct {
	MaxOpenConnections *prometheus.GaugeVec
	// Pool Status
	OpenConnections *prometheus.GaugeVec
	InUse           *prometheus.GaugeVec
	Idle            *prometheus.GaugeVec

	// Counters
	WaitCount         *prometheus.CounterVec
	WaitDuration      *prometheus.CounterVec
	MaxIdleClosed     *prometheus.CounterVec
	MaxIdleTimeClosed *prometheus.CounterVec
	MaxLifetimeClosed *prometheus.CounterVec
}

// StatDB do the statistic for given `db` status,
// mostly based on info from `sql.DBStats`.
func StatDB(env, dbName string, db *sql.DB) {
	var stats sql.DBStats
	labels := prometheus.Labels{
		"env": env,
		"db":  dbName,
	}
	var lastWaitCount, lastMaxIdleClosed, lastMaxIdleTimeClosed, lastMaxLifetimeClosed int64
	var lastWaitDuration time.Duration
	tk := time.NewTicker(time.Second)
	for {
		<-tk.C
		stats = db.Stats()
		DB.Stats.MaxOpenConnections.With(labels).Set(float64(stats.MaxOpenConnections))

		DB.Stats.OpenConnections.With(labels).Set(float64(stats.OpenConnections))
		DB.Stats.InUse.With(labels).Set(float64(stats.InUse))
		DB.Stats.Idle.With(labels).Set(float64(stats.Idle))

		DB.Stats.WaitCount.With(labels).Add(float64(stats.WaitCount - lastWaitCount))
		DB.Stats.WaitDuration.With(labels).Add(float64((stats.WaitDuration - lastWaitDuration) / time.Millisecond))
		DB.Stats.MaxIdleClosed.With(labels).Add(float64(stats.MaxIdleClosed - lastMaxIdleClosed))
		DB.Stats.MaxIdleTimeClosed.With(labels).Add(float64(stats.MaxIdleTimeClosed - lastMaxIdleTimeClosed))
		DB.Stats.MaxLifetimeClosed.With(labels).Add(float64(stats.MaxLifetimeClosed - lastMaxLifetimeClosed))

		lastWaitCount = stats.WaitCount
		lastWaitDuration = stats.WaitDuration
		lastMaxIdleClosed = stats.MaxIdleClosed
		lastMaxIdleTimeClosed = stats.MaxIdleTimeClosed
		lastMaxLifetimeClosed = stats.MaxLifetimeClosed
	}
}

var DB *dbMetrics
var API *apiMetrics
var HTTPClient *httpMetrics

var initOnce sync.Once

func init() {
	initMetricCollectors("ecom")
}

func initMetricCollectors(metricPrefix string) {
	initOnce.Do(func() {
		DB = &dbMetrics{
			ErrCnt: promauto.NewCounterVec(prometheus.CounterOpts{
				Subsystem: metricPrefix,
				Name:      "db_error_count",
				Help:      "How many errors occurred by error type",
			}, []string{"env", "type", "db", "target"}),
			ReconnectCnt: promauto.NewCounterVec(prometheus.CounterOpts{
				Subsystem: metricPrefix,
				Name:      "db_reconnect_count",
				Help:      "How many reconnect during query occurred",
			}, []string{"env", "type", "db", "target"}),
			RowsReturned: promauto.NewGaugeVec(prometheus.GaugeOpts{
				Subsystem: metricPrefix,
				Name:      "db_num_rows_returned",
				Help:      "How many rows returned by query",
			}, []string{"env", "target", "db"}),
			QueryDur: promauto.NewHistogramVec(prometheus.HistogramOpts{
				Subsystem: metricPrefix,
				Name:      "db_query_dur_ms",
				Help:      "Query duration in ms",
				Buckets:   []float64{20, 40, 80, 120, 150, 200, 250, 300, 500, 750, 1000, 1500, 2000}, // in milliseconds
			}, []string{"env", "db", "target"}),
			Stats: &dbStater{
				MaxOpenConnections: promauto.NewGaugeVec(prometheus.GaugeOpts{
					Subsystem: metricPrefix,
					Name:      "db_max_open_connections",
					Help:      "Maximum number of open connections to the database.",
				}, []string{"env", "db"}),
				OpenConnections: promauto.NewGaugeVec(prometheus.GaugeOpts{
					Subsystem: metricPrefix,
					Name:      "db_open_connections",
					Help:      "The number of established connections both in use and idle.",
				}, []string{"env", "db"}),
				InUse: promauto.NewGaugeVec(prometheus.GaugeOpts{
					Subsystem: metricPrefix,
					Name:      "db_conn_in_use",
					Help:      "The number of connections currently in use.",
				}, []string{"env", "db"}),
				Idle: promauto.NewGaugeVec(prometheus.GaugeOpts{
					Subsystem: metricPrefix,
					Name:      "db_conn_in_idle",
					Help:      "The number of idle connections.",
				}, []string{"env", "db"}),
				WaitCount: promauto.NewCounterVec(prometheus.CounterOpts{
					Subsystem: metricPrefix,
					Name:      "db_total_conn_waited_for",
					Help:      "The total number of connections waited for.",
				}, []string{"env", "db"}),
				WaitDuration: promauto.NewCounterVec(prometheus.CounterOpts{
					Subsystem: metricPrefix,
					Name:      "db_total_duration_waited_for_ms",
					Help:      "The total time blocked waiting for a new connection in ms.",
				}, []string{"env", "db"}),
				MaxIdleClosed: promauto.NewCounterVec(prometheus.CounterOpts{
					Subsystem: metricPrefix,
					Name:      "db_total_closed_conn_by_max_idle",
					Help:      "The total number of connections closed due to SetMaxIdleConns.",
				}, []string{"env", "db"}),
				MaxIdleTimeClosed: promauto.NewCounterVec(prometheus.CounterOpts{
					Subsystem: metricPrefix,
					Name:      "db_total_closed_conn_by_max_idle_time",
					Help:      "The total number of connections closed due to SetConnMaxIdleTime.",
				}, []string{"env", "db"}),
				MaxLifetimeClosed: promauto.NewCounterVec(prometheus.CounterOpts{
					Subsystem: metricPrefix,
					Name:      "db_total_closed_conn_by_max_life_time",
					Help:      "The total number of connections closed due to SetConnMaxLifetime.",
				}, []string{"env", "db"}),
			},
		}
		API = &apiMetrics{
			ErrCnt: promauto.NewCounterVec(prometheus.CounterOpts{
				Subsystem: metricPrefix,
				Name:      "api_error_count",
				Help:      "How many errors occurred by error type",
			}, []string{"svc", "env", "type", "path"}),
			ReqCnt: promauto.NewCounterVec(prometheus.CounterOpts{
				Subsystem: metricPrefix,
				Name:      "api_request_count",
				Help:      "How many requests received by API",
			}, []string{"svc", "method", "path", "env", "status"}),
			ReqDur: promauto.NewHistogramVec(prometheus.HistogramOpts{
				Subsystem: metricPrefix,
				Name:      "api_request_duration_ms",
				Help:      "The HTTP request latencies in milliseconds.",
				Buckets:   []float64{20, 40, 80, 120, 150, 200, 250, 300, 500, 750, 1000, 1500, 2000}, // in milliseconds
			}, []string{"method", "path", "env"}),
		}
		HTTPClient = &httpMetrics{
			ReconnectCnt: promauto.NewCounterVec(prometheus.CounterOpts{
				Subsystem: metricPrefix,
				Name:      "http_client_reconnect_count",
				Help:      "How many re-dial during establishing tcp connection when request to external services",
			}, []string{"env", "network", "address"}),
			ErrCnt: promauto.NewCounterVec(prometheus.CounterOpts{
				Subsystem: metricPrefix,
				Name:      "http_client_err_count",
				Help:      "How many errors occurred by error type when request to external services",
			}, []string{"env", "type", "service", "path"}),
			ResponseCode: promauto.NewCounterVec(prometheus.CounterOpts{
				Subsystem: metricPrefix,
				Name:      "http_client_response_code_count",
				Help:      "How many response by HTTP code per endpoint when request to external services",
			}, []string{"env", "service", "path", "code"}),
			ResponseDur: promauto.NewHistogramVec(prometheus.HistogramOpts{
				Subsystem: metricPrefix,
				Name:      "http_client_response_dur_ms",
				Help:      "Request duration in ms when request to external services",
				Buckets:   []float64{20, 40, 80, 120, 150, 200, 250, 300, 500, 750, 1000, 1500, 2000, 3000, 5000, 10000}, // in milliseconds
			}, []string{"env", "service", "path"}),
		}
	})
}
