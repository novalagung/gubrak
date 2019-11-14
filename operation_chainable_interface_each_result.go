package gubrak

type IChainableEachResult interface {
	Error() error
	IsError() bool
}

type resultEach struct {
	chainable *Chainable
	IChainableEachResult
}

func (g *resultEach) Error() error {
	return g.chainable.lastErrorCaught
}

func (g *resultEach) IsError() bool {
	return g.Error() != nil
}
