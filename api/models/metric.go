package models

import "github.com/prometheus/client_golang/prometheus"

var (
	//HTTPReqDuration metric:http_request_duration_seconds
	HTTPReqDuration *prometheus.HistogramVec
	//HTTPReqTotal metric:http_request_total
	HTTPReqTotal *prometheus.CounterVec
	//
	//CpuTemp *prometheus.GaugeOpts
)

func init() {
	// 监控接口请求耗时
	//HTTPReqDuration metric:http_request_duration_seconds
	// HistogramVec 是一组Histogram
	// 这里的"method"、"path" 都是label
	HTTPReqDuration = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "http_request_duration_seconds",
		Help:    "The HTTP request latencies in seconds.",
		Buckets: nil,
	}, []string{"method", "path"})


	// 监控接口请求次数
	//HTTPReqTotal metric:http_request_total
	// HistogramVec 是一组Histogram
	// 这里的"method"、"path"、"status" 都是label
	HTTPReqTotal = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "http_requests_total",
		Help: "Total number of HTTP requests made.",
	}, []string{"method", "path", "status"})

	// 添加prometheus性能监控指标
	prometheus.MustRegister(HTTPReqTotal)
	prometheus.MustRegister(HTTPReqDuration)

	//prometheus.MustRegister(CpuTemp)
	//prometheus.MustRegister(HdFailures)
}