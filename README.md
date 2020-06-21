# example

```go
var (
	//服务器带宽可利用率
	bandwidth_rate := prometheus.NewGauge(prometheus.GaugeOpts{
		Name:      "bandwidth_rate",
		Help:      "rate of bandwidth lefted",
	})
)
    
func Register() {
	prometheus.MustRegister(requestCount)
	prometheus.MustRegister(requestLatency)
	prometheus.MustRegister(bandwidth_rate)
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

```

Exporter在原有的latency & requestCount指标基础上增加了1个指标bandwidth_rate，metrics类型为Gauge。

bandwidth_rate 表示服务器带宽可利用率。服务器的使用明显存在高峰期&低峰期。时间为自变量，使用抛物线进行模拟，高峰期在20点到21点间，默认这段时间利用率100%.
