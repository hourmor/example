package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"time"
)

var (
	requestCount = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name:      "request_total",
			Help:      "Number of request processed by this service.",
		}, []string{},
	)

	requestLatency = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:      "request_latency_seconds",
			Help:      "Time spent in this service.",
			Buckets:   []float64{0.01, 0.02, 0.05, 0.1, 0.2, 0.5, 1.0, 2.0, 5.0, 10.0, 20.0, 30.0, 60.0, 120.0, 300.0},
		}, []string{},
	)
	//服务器带宽可利用率
	bandwidth_rate := prometheus.NewGauge(prometheus.GaugeOpts{
		Name:      "bandwidth_rate",
		Help:      "rate of bandwidth lefted",
	})
)

// AdmissionLatency measures latency / execution time of Admission Control execution
// usual usage pattern is: timer := NewAdmissionLatency() ; compute ; timer.Observe()
type RequestLatency struct {
	histo *prometheus.HistogramVec
	start time.Time
}

func Register() {
	prometheus.MustRegister(requestCount)
	prometheus.MustRegister(requestLatency)
	prometheus.MustRegister(bandwidth_rate)
}


// NewAdmissionLatency provides a timer for admission latency; call Observe() on it to measure
func NewAdmissionLatency() *RequestLatency {
	return &RequestLatency{
		histo: requestLatency,
		start: time.Now(),
	}
}

// Observe measures the execution time from when the AdmissionLatency was created
func (t *RequestLatency) Observe() {
	(*t.histo).WithLabelValues().Observe(time.Now().Sub(t.start).Seconds())
}


// RequestIncrease increases the counter of request handled by this service
func RequestIncrease() {
	requestCount.WithLabelValues().Add(1)
	// 获取时间
	h:=time.Now().Hour()
	hour:=float64(h)
	//范围映射为0~23 -> [0,200]
	bandwithleft=(hour - 20) * (hour - 20) / 2
	//在[20:00,21:00)这段时间可利用率为0%
	bandwidth_rate.Set(bandwithleft/200)
}
