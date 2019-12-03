package yarxyarx

import (
    "sync"
    "context"
)

var handlerContext context.Context
var contextLock sync.RWMutex

func UseXrayContext(context context.Context) {
    contextLock.Lock()
    defer contextLock.Unlock()
    handlerContext = context
}

func CurrentXrayContext() (context.Context) {
    contextLock.RLock()
    defer contextLock.RUnlock()
    return handlerContext
}

func WithXrayContext(otherContext context.Context) (context.Context) {
    xrayContext := CurrentXrayContext()
    return mergeContexts(otherContext, xrayContext)
}

