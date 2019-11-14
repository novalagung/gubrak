package gubrak

type IChainableContainsResult interface {
	ResultAndError() (bool, error)
	Result() bool
	Error() error
	IsError() bool
}

type resultContains struct {
	chainable *Chainable
	IChainableContainsResult
}

func (g *resultContains) ResultAndError() (bool, error) {
	return g.Result(), g.Error()
}

func (g *resultContains) Result() bool {
	v, _ := g.chainable.data.(bool)
	return v
}

func (g *resultContains) Error() error {
	return g.chainable.lastErrorCaught
}

func (g *resultContains) IsError() bool {
	return g.Error() != nil
}
