package interceptor

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"cms-api/internal/infra/telemetry"
)

func UnaryTracing(tracer *telemetry.Tracer) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		ctx = extractTraceContext(ctx)

		ctx, span := tracer.Start(ctx, info.FullMethod,
			trace.WithSpanKind(trace.SpanKindServer),
			trace.WithAttributes(
				semconv.RPCSystemGRPC,
				semconv.RPCMethod(info.FullMethod),
				attribute.String("rpc.grpc.request_type", "unary"),
			),
		)
		defer span.End()

		resp, err := handler(ctx, req)

		if err != nil {
			st, _ := status.FromError(err)
			span.SetAttributes(semconv.RPCGRPCStatusCodeKey.Int(int(st.Code())))
			span.RecordError(err)
			span.SetStatus(codes.Error, st.Message())
		} else {
			span.SetStatus(codes.Ok, "")
		}

		return resp, err
	}
}

func StreamTracing(tracer *telemetry.Tracer) grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		ctx := extractTraceContext(ss.Context())

		ctx, span := tracer.Start(ctx, info.FullMethod,
			trace.WithSpanKind(trace.SpanKindServer),
			trace.WithAttributes(
				semconv.RPCSystemGRPC,
				semconv.RPCMethod(info.FullMethod),
				attribute.String("rpc.grpc.request_type", "stream"),
				attribute.Bool("rpc.grpc.client_stream", info.IsClientStream),
				attribute.Bool("rpc.grpc.server_stream", info.IsServerStream),
			),
		)
		defer span.End()

		wrappedStream := &tracedServerStream{
			ServerStream: ss,
			ctx:          ctx,
		}

		err := handler(srv, wrappedStream)

		if err != nil {
			st, _ := status.FromError(err)
			span.SetAttributes(semconv.RPCGRPCStatusCodeKey.Int(int(st.Code())))
			span.RecordError(err)
			span.SetStatus(codes.Error, st.Message())
		} else {
			span.SetStatus(codes.Ok, "")
		}

		return err
	}
}

type tracedServerStream struct {
	grpc.ServerStream
	ctx context.Context
}

func (s *tracedServerStream) Context() context.Context {
	return s.ctx
}

type metadataCarrier struct {
	md metadata.MD
}

func (c *metadataCarrier) Get(key string) string {
	vals := c.md.Get(key)
	if len(vals) > 0 {
		return vals[0]
	}
	return ""
}

func (c *metadataCarrier) Set(key, value string) {
	c.md.Set(key, value)
}

func (c *metadataCarrier) Keys() []string {
	keys := make([]string, 0, len(c.md))
	for k := range c.md {
		keys = append(keys, k)
	}
	return keys
}

func extractTraceContext(ctx context.Context) context.Context {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return ctx
	}

	return otel.GetTextMapPropagator().Extract(ctx, &metadataCarrier{md: md})
}
