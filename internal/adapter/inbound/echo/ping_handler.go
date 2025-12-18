package echo

import (
	"fmt"
	"net/http"
	"simple-golang/internal/port/inbound"
	"simple-golang/util"
	"time"

	"github.com/labstack/echo/v4"
)

type pingHandler struct{}

func NewPingHandler() inbound.PingHandlerInterface {
	return &pingHandler{}
}

func (h *pingHandler) Ping(c echo.Context) error {
	idle0, total0 := util.GetCPUSample()
	time.Sleep(1 * time.Second)
	idle1, total1 := util.GetCPUSample()

	idleTicks := float64(idle1 - idle0)
	totalTicks := float64(total1 - total0)
	cpuUsage := 100 * (totalTicks - idleTicks) / totalTicks

	total, free, buffers, cached := util.GetMemorySample()
	coreCount := util.GetCoreSample()

	return c.JSON(http.StatusOK, map[string]any{
		"message": "pong",
		"core": []map[string]any{
			{"core": fmt.Sprintf("%d Core", coreCount)},
		},
		"cpu": []map[string]any{
			{
				"usage": fmt.Sprintf("%f %%", cpuUsage),
				"busy":  fmt.Sprintf("%f %%", totalTicks-idleTicks),
				"total": fmt.Sprintf("%f %%", totalTicks),
			},
		},
		"memory": []map[string]any{
			{
				"usage":  fmt.Sprintf("%f %%", 100*(1-float64(free)/float64(total))),
				"total":  fmt.Sprintf("%f MB", float64(total)/1024),
				"free":   fmt.Sprintf("%f MB", float64(free)/1024),
				"buffer": fmt.Sprintf("%f MB", float64(buffers)/1024),
				"cached": fmt.Sprintf("%f MB", float64(cached)/1024),
			},
		},
	})
}
