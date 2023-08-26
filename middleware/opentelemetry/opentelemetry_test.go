//go:build e2e

package opentelemetry

import (
	"github.com/liquanhui-99/lr"
	"go.opentelemetry.io/otel"
	"testing"
	"time"
)

func TestOpenTelemetry(t *testing.T) {
	tracer := otel.GetTracerProvider().Tracer(instrumentationName)
	builder := MiddlewareBuilder{
		Tracer: tracer,
	}
	s := lr.NewHTTPServer("tcp", ":8081", lr.Use(builder.Build()))

	s.GET("/user/profile", func(ctx *lr.Context) {
		c, span := tracer.Start(ctx.Req.Context(), "first_layer")
		defer span.End()

		c, second := tracer.Start(c, "second_layer")
		time.Sleep(time.Second)

		c, third1 := tracer.Start(c, "third_layer_1")
		time.Sleep(100 * time.Millisecond)
		third1.End()

		c, third2 := tracer.Start(c, "third_layer_2")
		time.Sleep(300 * time.Millisecond)
		third2.End()
		second.End()
	})

	if err := s.Server(); err != nil {
		panic(err)
	}
}
