package yarxyarx

import (
    "context"
    "time"
)

// based loosely on https://medium.com/@dlagoza/playing-with-multiple-contexts-in-go-9f72cbcff56e
type mergedContext struct {
	mainCtx context.Context
	xrayCtx context.Context
	err error
}

// X-Ray doesn't use these three methods
func (c *mergedContext) Done() <-chan struct{} {
	return c.mainCtx.Done()
}

func (c *mergedContext) Err() error {
    return c.mainCtx.Err()
}

func (c *mergedContext) Deadline() (deadline time.Time, ok bool) {
	return c.mainCtx.Deadline()
}

func (c *mergedContext) Value(key interface{}) interface{} {
    // if the regular context doesn't have this value, it's probably in the X-Ray context
    v := c.mainCtx.Value(key)
    if v == nil {
        v = c.xrayCtx.Value(key)
    }
    return v
}

func (c *mergedContext) run() {
	<-c.mainCtx.Done()
}

func mergeContexts(mainCtx, otherCtx context.Context) context.Context {
	c := &mergedContext{mainCtx: mainCtx, xrayCtx: otherCtx }
	go c.run()
	return c
}

