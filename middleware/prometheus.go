package middleware

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	dto "github.com/prometheus/client_model/go"
	"github.com/valyala/fasthttp/fasthttpadaptor"
)

func PrometheusHandler(c *fiber.Ctx) error {
	reg := prometheus.NewRegistry()
	newCelsiusHistogram(reg)
	newKelvinGauge(reg)
	newRequestCounter(reg, c)

	gatherers := prometheus.Gatherers{
		reg,
	}

	p := fasthttpadaptor.NewFastHTTPHandler(
		promhttp.HandlerFor(
			gatherers,
			promhttp.HandlerOpts{
				// Opt into OpenMetrics to support exemplars.
				DisableCompression: true,
			},
		),
	)
	p(c.Context())
	return nil
}

func newCelsiusHistogram(reg *prometheus.Registry) {
	temps := prometheus.NewHistogram(prometheus.HistogramOpts{
		Name:    "pond_temperature_celsius",
		Help:    "The temperature of the frog pond.", // Sorry, we can't measure how badly it smells.
		Buckets: prometheus.LinearBuckets(20, 5, 5),  // 5 buckets, each 5 centigrade wide.
	})
	temps.Observe(30)
	metric := &dto.Metric{}
	temps.Write(metric)
	reg.MustRegister(temps)
}

func newKelvinGauge(reg *prometheus.Registry) {
	temp := prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "temperature_kelvin",
			Help: "Temperature in Kelvin.",
		},
		[]string{"location"},
	)
	reg.MustRegister(temp)
	temp.WithLabelValues("outside").Set(273.14)
	temp.WithLabelValues("inside").Set(298.44)
}

func newRequestCounter(reg *prometheus.Registry, c *fiber.Ctx) {
	reqCounter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "promhttp_request_total",
			Help: "Total number of request encountered by the promhttp metric handler.",
		},
		[]string{"status", "method", "path"},
	)
	sc := c.Context().Response.StatusCode()
	me := c.Context().Method()
	pa := c.Context().Path()
	reg.MustRegister(reqCounter)
	reqCounter.WithLabelValues(fmt.Sprintf("%d", sc), string(me), string(pa)).Add(1)
}
