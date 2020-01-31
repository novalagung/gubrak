package gubrak

type IChainableNoReturnValueResult interface {
	Error() error
	IsError() bool
}

type resultNoReturnValue struct {
	chainable *Chainable
	IChainableNoReturnValueResult
}

type resultEach = resultNoReturnValue

func (g *resultNoReturnValue) Error() error {
	return g.chainable.lastErrorCaught
}

func (g *resultNoReturnValue) IsError() bool {
	return g.Error() != nil
}

func (g *resultNoReturnValue) LastSuccessOperation() Operation {
	return g.chainable.lastSuccessOperation
}

func (g *resultNoReturnValue) LastErrorOperation() Operation {
	return g.chainable.lastErrorOperation
}

func (g *resultNoReturnValue) LastOperation() Operation {
	return g.chainable.lastOperation
}
