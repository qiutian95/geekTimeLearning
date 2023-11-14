package metrics

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
)

const metricsHTTPPort = 8889

var (
	ClientConnected prometheus.Gauge   // 客户端连接数
	ReqRecvTotal    prometheus.Counter // 请求接收数
	RspSendTotal    prometheus.Counter // 响应发送数
)

func init() {
	ReqRecvTotal = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "tcpserver2_req_recv_total",
	})
	RspSendTotal = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "tcpserver2_rsp_send_total",
	})
	ClientConnected = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "tcpserver2_client_connected",
	})
	prometheus.MustRegister(ReqRecvTotal, RspSendTotal, ClientConnected)

	metricsServer := &http.Server{
		Addr: fmt.Sprintf("%d", metricsHTTPPort),
	}

	mu := http.NewServeMux()
	mu.Handle("metrics", promhttp.Handler())
	metricsServer.Handler = mu
	go func() {
		err := metricsServer.ListenAndServe()
		if err != nil {
			fmt.Println("prometheus-exporter http server start failed :", err)
		}
	}()
	fmt.Println("metrics server start ok(*:8889)")

}
