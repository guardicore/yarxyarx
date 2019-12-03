package ctxhttp

import (
    "yarxyarx"
    "context"
    "net/http"
    "golang.org/x/net/context/ctxhttp"
    "github.com/aws/aws-xray-sdk-go/xray"
    "io"
    "net/url"
)

func Do(ctx context.Context, client *http.Client, req *http.Request) (*http.Response, error) {
    wrappedCtx := yarxyarx.WithXrayContext(ctx)
    return ctxhttp.Do(wrappedCtx, xray.Client(client), req)
}

func Get(ctx context.Context, client *http.Client, url string) (*http.Response, error) {
    wrappedCtx := yarxyarx.WithXrayContext(ctx)
    return ctxhttp.Get(wrappedCtx, xray.Client(client), url)
}

func Head(ctx context.Context, client *http.Client, url string) (*http.Response, error) {
    wrappedCtx := yarxyarx.WithXrayContext(ctx)
    return ctxhttp.Head(wrappedCtx, xray.Client(client), url)
}

func Post(ctx context.Context, client *http.Client, url string, bodyType string, body io.Reader) (*http.Response, error) {
    wrappedCtx := yarxyarx.WithXrayContext(ctx)
    return ctxhttp.Post(wrappedCtx, xray.Client(client), url, bodyType, body)
}

func PostForm(ctx context.Context, client *http.Client, url string, data url.Values) (*http.Response, error) {
    wrappedCtx := yarxyarx.WithXrayContext(ctx)
    return ctxhttp.PostForm(wrappedCtx, xray.Client(client), url, data)
}

