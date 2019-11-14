package gubrak

type IChainableIncludesResult interface {
	ResultAndError() (bool, error)
	Result() bool
	Error() error
	IsError() bool
}

type resultIncludes struct {
	chainable *Chainable
	IChainableIncludesResult
}

func (g *resultIncludes) ResultAndError() (bool, error) {
	return g.Result(), g.Error()
}

func (g *resultIncludes) Result() bool {
	v, _ := g.chainable.data.(bool)
	return v
}

func (g *resultIncludes) Error() error {
	return g.chainable.lastErrorCaught
}

func (g *resultIncludes) IsError() bool {
	return g.Error() != nil
}
