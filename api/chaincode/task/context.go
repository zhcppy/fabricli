/*
@Time 2019-09-17 18:27
@Author ZH

*/
package task

type Context interface {
	SetVar(name, value string)
	GetVar(name string) (string, bool)
}

func NewContext() Context {
	return &defaultContext{
		vars: make(map[string]string),
	}
}

type defaultContext struct {
	vars map[string]string
}

func (c *defaultContext) SetVar(k, v string) {
	c.vars[k] = v
}

func (c *defaultContext) GetVar(k string) (string, bool) {
	value, ok := c.vars[k]
	return value, ok
}
