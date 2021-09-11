package trace

import (
	"context"
	"errors"
	"log"
	"os"

	"github.com/opentracing/opentracing-go"
	zipkinot "github.com/openzipkin-contrib/zipkin-go-opentracing"
	"github.com/openzipkin/zipkin-go"
	"github.com/openzipkin/zipkin-go/model"
	reporterhttp "github.com/openzipkin/zipkin-go/reporter/http"
	"github.com/pot-code/gobit/pkg/config"
	"github.com/pot-code/gobit/pkg/util"
	"go.elastic.co/apm"
	"go.elastic.co/apm/module/apmot"
	"go.elastic.co/apm/transport"
)

func NewEcsTracerProvider(tc *TraceConfig, b *config.BaseConfig, lm *util.LifecycleManager) *apm.Tracer {
	os.Setenv("ELASTIC_APM_SERVER_URL", tc.URL)

	_, err := transport.InitDefault()
	util.HandleFatalError("failed to init apm default", err)

	at := apm.DefaultTracer
	at.Service.Environment = b.Env
	at.Service.Name = b.AppID

	aot := apmot.New()
	opentracing.SetGlobalTracer(aot)

	lm.AddLivenessProbe(func(ctx context.Context) error {
		if !at.Active() {
			return errors.New("inactive")
		}
		return nil
	})
	lm.OnExit(func(ctx context.Context) {
		at.Close()
		log.Println("shutdown tracer")
	})

	return at
}

func NewZipkinTracerProvider(tc *TraceConfig, b *config.BaseConfig, lm *util.LifecycleManager) *zipkin.Tracer {
	reporter := reporterhttp.NewReporter(tc.URL)
	zt, err := zipkin.NewTracer(reporter, zipkin.WithLocalEndpoint(&model.Endpoint{
		ServiceName: b.AppID,
		Port:        8080,
	}))
	util.HandleFatalError("failed to init zipkin tracer", err)

	lm.OnExit(func(ctx context.Context) {
		reporter.Close()
	})

	opentracing.SetGlobalTracer(zipkinot.Wrap(zt))
	return zt
}
