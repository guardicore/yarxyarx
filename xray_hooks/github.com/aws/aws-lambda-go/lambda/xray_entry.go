package lambda

import (
    "context"
    "github.com/aws/aws-lambda-go/lambda"
    "yarxyarx"
)

type xrayContextHandler struct {
    Handler
}

func (lambdaHandler xrayContextHandler) Invoke(ctx context.Context, payload []byte) ([]byte, error) {
    yarxyarx.UseXrayContext(ctx)
    return lambdaHandler.Handler.Invoke(ctx, payload)
}

func Start(handler interface{}) {
	wrappedHandler := NewHandler(handler)
	StartHandler(wrappedHandler)
}

func StartHandler(handler Handler) {
    xrayHandler := xrayContextHandler { Handler: handler }
    lambda.StartHandler(xrayHandler)
}
