package http

import (
	"context"
	"encoding/json"
	"math"
	"net/http"
)

const (
	_abortIndex int8 = math.MaxInt8 / 2
)

// Context is the most important part. It allows us to pass variables between
// middleware, manage the flow, validate the JSON of a request and render a
// JSON response for example.
type Context struct {
	context.Context
	Request *http.Request
	Writer  http.ResponseWriter

	// keys is a key/value pair exclusively for the context of each request.
	keys map[string]interface{}

	// flow control
	index    int8
	handlers []Hander

	s      *Service
	method string
}

// Next should be used only inside middleware.
// It executes the pending handlers in the chain inside the calling handler.
// See example in godoc.
func (c *Context) Next() {
	s := int8(len(c.handlers))
	for ; c.index < s; c.index++ {
		c.handlers[c.index](c)
	}
}

// Abort prevents pending handlers from being called. Note that this will not stop the current handler.
// Let's say you have an authorization middleware that validates that the current request is authorized.
// If the authorization fails (ex: the password does not match), call Abort to ensure the remaining handlers
// for this request are not called.
func (c *Context) Abort() {
	c.index = _abortIndex
}

// Set is used to store a new key/value pair exclusively for this context.
// It also lazy initializes  c.Keys if it was not used previously.
func (c *Context) Set(key string, value interface{}) {
	if c.keys == nil {
		c.keys = make(map[string]interface{})
	}
	c.keys[key] = value
}

// Get returns the value for the given key, ie: (value, true).
// If the value does not exists it returns (nil, false)
func (c *Context) Get(key string) (value interface{}, exists bool) {
	value, exists = c.keys[key]
	return
}

// JSON serializes the given struct as JSON into the response body.
// It also sets the Content-Type as "application/json".
func (c *Context) JSON(d interface{}) {
	body, err := json.Marshal(d)
	if err != nil {
		c.Writer.Write([]byte(err.Error()))
		return
	}

	h := c.Writer.Header()
	h.Set("Content-Type", "application/json; charset=utf-8")

	c.Writer.Write(body)
}
