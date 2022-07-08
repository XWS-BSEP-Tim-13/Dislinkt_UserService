package tracer

import (
	"context"
	"fmt"
	"github.com/grpc-ecosystem/go-grpc-middleware/util/metautils"
	opentracing "github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-client-go/config"
	"io"
	"net/http"
)

// Init returns an instance of Jaeger Tracer.
func Init() (opentracing.Tracer, io.Closer) {
	cfg, err := config.FromEnv()
	if err != nil {
		fmt.Println(err.Error())
		return nil, nil
	}

	cfg.ServiceName = "user_service"
	cfg.Sampler.Type = jaeger.SamplerTypeConst
	cfg.Sampler.Param = 1
	cfg.Reporter.LogSpans = true

	tracer, closer, err := cfg.NewTracer(config.Logger(jaeger.StdLogger))
	if err != nil {
		panic(fmt.Sprintf("ERROR: cannot init Jaeger: %v\n", err))
	}
	return tracer, closer
}

// Inject injects the outbound HTTP request with the given span's context to ensure
// correct propagation of span context throughout the trace.
func Inject(span opentracing.Span, request *http.Request) error {
	return span.Tracer().Inject(
		span.Context(),
		opentracing.HTTPHeaders,
		opentracing.HTTPHeadersCarrier(request.Header))
}

// Extract extracts the inbound HTTP request to obtain the parent span's context to ensure
// correct propagation of span context throughout the trace.
func Extract(tracer opentracing.Tracer, r *http.Request) (opentracing.SpanContext, error) {
	return tracer.Extract(
		opentracing.HTTPHeaders,
		opentracing.HTTPHeadersCarrier(r.Header))
}

// StartSpanFromRequest extracts the parent span context from the inbound HTTP request
// and starts a new child span if there is a parent span.
func StartSpanFromRequest(spanName string, tracer opentracing.Tracer, r *http.Request) opentracing.Span {
	spanCtx, _ := Extract(tracer, r)
	return tracer.StartSpan(spanName, ext.RPCServerOption(spanCtx))
}

func StartSpanFromContext(ctx context.Context, spanName string) opentracing.Span {
	span, _ := opentracing.StartSpanFromContext(ctx, spanName)
	return span
}

func ContextWithSpan(ctx context.Context, span opentracing.Span) context.Context {
	return opentracing.ContextWithSpan(ctx, span)
}

// ExtractSpanContextFromMetadata gRPC related span extraction
func ExtractSpanContextFromMetadata(tracer opentracing.Tracer, ctx context.Context) opentracing.SpanContext {
	md := metautils.ExtractIncoming(ctx)
	parentSpanContext, _ := tracer.Extract(opentracing.HTTPHeaders, metadataTextMap(md))
	return parentSpanContext
}

func StartSpanFromContextMetadata(ctx context.Context, name string) opentracing.Span {
	spanContext := ExtractSpanContextFromMetadata(opentracing.GlobalTracer(), ctx)
	span := opentracing.StartSpan(
		name,
		opentracing.ChildOf(spanContext),
	)
	return span
}
