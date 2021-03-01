package middleware

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/valyala/fasthttp/fasthttpadaptor"
)

var (
	p = fasthttpadaptor.NewFastHTTPHandler(
		promhttp.HandlerFor(
			prometheus.DefaultGatherer,
			promhttp.HandlerOpts{
				// Opt into OpenMetrics to support exemplars.
				DisableCompression: true,
			},
		),
	)

	errCnt = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "promhttp_metric_handler_errors_total",
			Help: "Total number of internal errors encountered by the promhttp metric handler.",
		},
		[]string{"status", "method", "path"},
	)
)

func init() {
	prometheus.MustRegister(errCnt)
}

func PrometheusHandler(c *fiber.Ctx) error {
	p(c.Context())
	return nil
}

func Prometheus(c *fiber.Ctx) {
	sc := c.Context().Response.StatusCode()
	me := c.Context().Method()
	pa := c.Context().Path()
	errCnt.WithLabelValues(fmt.Sprintf("%d", sc), string(me), string(pa)).Add(1)
	c.Next()
}
