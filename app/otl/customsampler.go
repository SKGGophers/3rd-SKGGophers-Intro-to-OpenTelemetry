package otl

import (
	strace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
)

type CustomSampler struct {
	ExcludedRoutes []string
	Desc           string
}

func (cs CustomSampler) ShouldSample(p strace.SamplingParameters) strace.SamplingResult {
	pctx := trace.SpanContextFromContext(p.ParentContext)

	for _, er := range cs.ExcludedRoutes {
		if er == p.Name {
			return strace.SamplingResult{
				Decision:   strace.Drop,
				Tracestate: pctx.TraceState(),
			}
		}
	}

	return strace.SamplingResult{
		Decision:   strace.RecordAndSample,
		Tracestate: pctx.TraceState(),
	}
}

func (cs CustomSampler) Description() string {
	return cs.Desc
}
