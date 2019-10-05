package gubrak

import "fmt"

type Chainable struct {
	data                interface{}
	enforcePanicOnError bool
	firstError          error
}

func New() *Chainable {
	g := new(Chainable)
	g.enforcePanicOnError = false
	return g
}

func (g *Chainable) From(data interface{}) *Chainable {
	g.data = data
	return g
}

func (g *Chainable) Filter(callback interface{}) *Chainable {
	res, err := Filter(g.data, callback)
	return g.setResultAndError("Filter()", res, err)
}

func (g *Chainable) Map(callback interface{}) *Chainable {
	res, err := Map(g.data, callback)
	return g.setResultAndError("Map()", res, err)
}

func (g *Chainable) ToInterface() interface{} {
	return g.data
}

func (g *Chainable) Error() error {
	return g.firstError
}

func (g *Chainable) IsError() bool {
	return g.firstError != nil
}

func (g *Chainable) setResultAndError(op string, data interface{}, err error) *Chainable {
	g.data = data

	if g.firstError == nil && err != nil {
		g.firstError = fmt.Errorf("error on %s method. %s", op, err.Error())
	}

	return g
}
