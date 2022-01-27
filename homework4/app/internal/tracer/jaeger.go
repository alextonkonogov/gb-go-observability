package tracer

import (
	"fmt"
	"io"

	"go.uber.org/zap"

	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go/config"
)

type zapWrapper struct {
	logger *zap.Logger
}

// Error logs a message at error priority
func (w *zapWrapper) Error(msg string) {
	w.logger.Error(msg)
}

// Infof logs a message at info priority
func (w *zapWrapper) Infof(msg string, args ...interface{}) {
	w.logger.Sugar().Infof(msg, args...)
}

func InitJaeger(service string, logger *zap.Logger) (opentracing.Tracer, io.Closer) {
	cfg := &config.Configuration{
		ServiceName: service,
		Sampler: &config.SamplerConfig{
			Type:  "const",
			Param: 1,
		},
		Reporter: &config.ReporterConfig{
			LogSpans: true,
		},
	}

	tracer, closer, err := cfg.NewTracer(config.Logger(&zapWrapper{logger: logger}))
	if err != nil {
		panic(fmt.Sprintf("ERROR: cannot init Jaeger: %v\n", err))
	}

	return tracer, closer
}
