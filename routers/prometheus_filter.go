package routers

import (
	"fmt"
	"strings"
	"time"

	"github.com/beego/beego/context"
	"github.com/casdoor/casdoor/object"
	"github.com/casdoor/casdoor/util"
)

func recordSystemInfo(systemInfo *util.SystemInfo) {
	for i, value := range systemInfo.CpuUsage {
		object.CpuUsage.WithLabelValues(fmt.Sprintf("%d", i)).Set(value)
	}
	object.MemoryUsage.WithLabelValues("memoryUsed").Set(float64(systemInfo.MemoryUsed))
	object.MemoryUsage.WithLabelValues("memoryTotal").Set(float64(systemInfo.MemoryTotal))
}

func PrometheusBeforeFilter(ctx *context.Context) {
	method := ctx.Input.Method()
	endpoint := ctx.Input.URL()
	if strings.HasPrefix(endpoint, "/api/metrics") {
		systemInfo, err := util.GetSystemInfo()
		if err == nil {
			recordSystemInfo(systemInfo)
		}
		return
	}

	if strings.HasPrefix(endpoint, "/api") {
		ctx.Input.SetData("startTime", time.Now())
		object.TotalThroughput.Inc()
		object.APIThroughput.WithLabelValues(endpoint, method).Inc()
	}
}
