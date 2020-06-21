package main

import (
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"net/http"
	"os"
	"strconv"
	"example/metrics"
)

func main(){
	http.HandleFunc("/abc", index)
	http.Handle("/metrics", promhttp.Handler())
	metrics.Register()
	err := http.ListenAndServe(":5565", nil) // 设置监听的端口
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func index(w http.ResponseWriter, r *http.Request) {
	timer:=metrics.NewAdmissionLatency()
	metrics.RequestIncrease()
	num:=os.Getenv("Num")
	if num==""{
		SumOfIndex(10)
		_,err:=w.Write([]byte("there is no env Num. Computation successed\n"))
		if err!=nil{
			log.Println("err:"+err.Error()+" No\n")
		}
	}else{
		numInt,_:=strconv.Atoi(num)
		SumOfIndex(numInt)
		_,err:=w.Write([]byte("there is env Num. Computation successed\n"))
		if err!=nil{
			log.Println("err:"+err.Error()+" Yes\n")
		}
	}
	timer.Observe()
}

func SumOfIndex(n int)int{
	return n*(n+1)/2
}
