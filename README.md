# example

## file tree

```bash
.
├── README.md   // repository readme
├── img         //图片
│   ├── 21.png
│   ├── 22.png
│   ├── 24.png
│   ├── 25.png
│   ├── 7.png
│   ├── q1.jpg
│   └── q2.png
├── report.md   //实验report
├── report.pdf  //实验report
└── src			// 源代码
    ├── deploy
    │   ├── deployment.yaml
    │   ├── prometheus.config.yml //prometheus抓取目标配置文件
    │   ├── prometheus.deploy.yml //prometheus部署所使用Deployment
    │   ├── prometheus.rbac.yml   //prometheus权限配置文件
    │   └── service.yaml
    ├── go.mod
    ├── go.sum
    ├── metrics
    │   └── metrics.go
    ├── metrics_version
    │   ├── Dockerfile  
    │   ├── example
    │   └── main.go
    └── without_metrics
        ├── Dockerfile
        └── main.go
```

## 代码逻辑

```go
import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/mem"
	"time"
)

...

var (
	//服务器带宽可利用率
	bandwidth_rate := prometheus.NewGauge(prometheus.GaugeOpts{
		Name:      "bandwidth_rate",
		Help:      "rate of bandwidth lefted",
	})
	// CPU 利用率
	cpu_rate = prometheus.NewGauge(
		prometheus.GaugeOpts{
		Name:      "cpu_rate",
		Help:      "rate of cpu percent used.",
	})
	// disk 利用率
	dis_rate = prometheus.NewGauge(
		prometheus.GaugeOpts{
		Name:      "dis_rate",
		Help:      "rate of disk memory used.",
	})	
	// memory 利用率
	mem_rate = prometheus.NewGauge(
		prometheus.GaugeOpts{
		Name:      "mem_rate",
		Help:      "rate of system memory used.",
	})
)

...

func Register() {
	prometheus.MustRegister(requestCount)
	prometheus.MustRegister(requestLatency)
	prometheus.MustRegister(bandwidth_rate)
	prometheus.MustRegister(cpu_rate)
	prometheus.MustRegister(dis_rate)
	prometheus.MustRegister(mem_rate)
}

...

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

	mem_,_ :=mem.VirtualMemory()
	mem_rate.Set(mem_.UsedPercent)
	
	cpuper, _:= cpu.Percent(time.Second, false)
	cpu_rate.Set(cpuper[0])

	parts, _ := disk.Partitions(true)
	diskInfo, _ := disk.Usage(parts[0].Mountpoint)
	dis_rate.Set(diskInfo.UsedPercent)

}
```

Exporter在原有的latency & requestCount指标基础上增加了指标bandwidth_rate,cpu_rate,dis_rate,mem_rate，metrics类型为Gauge。

bandwidth_rate 表示服务器带宽可利用率。服务器的使用明显存在高峰期&低峰期。时间为自变量，使用抛物线进行模拟，高峰期在20点到21点间，默认这段时间利用率100%.

cpu_rate,dis_rate,mem_rate对应得是硬件信息。属于*宿主机的监控数据*。

