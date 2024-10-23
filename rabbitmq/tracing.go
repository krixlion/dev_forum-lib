package rabbitmq

import (
	"context"

	"github.com/krixlion/dev_forum-lib/tracing"
	amqp "github.com/rabbitmq/amqp091-go"
)

// extractAMQPHeadersFromCtx injects the trace data from the context into the header map.
func extractAMQPHeadersFromCtx(ctx context.Context) amqp.Table {
	m := tracing.ExtractMetadataFromContext(ctx)

	headers := make(map[string]interface{}, len(m))
	for k, v := range m {
		headers[k] = v
	}

	return headers
}

// injectAMQPHeadersIntoCtx extracts the trace data from the header and puts it
// into the returned context. Any extracted non-string values are discarded.
func injectAMQPHeadersIntoCtx(ctx context.Context, headers amqp.Table) context.Context {
	m := make(map[string]string, len(headers))
	for k, v := range headers {
		str, ok := v.(string)
		if !ok {
			continue
		}
		m[k] = str
	}

	return tracing.InjectMetadataIntoContext(ctx, m)
}
