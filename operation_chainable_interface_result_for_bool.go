package gubrak

type IChainableBoolResult interface {
	ResultAndError() (bool, error)
	Result() bool
	Error() error
	IsError() bool
	LastSuccessOperation() Operation
	LastErrorOperation() Operation
	LastOperation() Operation
}

type resultBool struct {
	chainable *Chainable
	IChainableBoolResult
}

type resultContains = resultBool

func (g *resultBool) ResultAndError() (bool, error) {
	return g.Result(), g.Error()
}

func (g *resultBool) Result() bool {
	v, _ := g.chainable.data.(bool)
	return v
}

func (g *resultBool) Error() error {
	return g.chainable.lastErrorCaught
}

func (g *resultBool) IsError() bool {
	return g.Error() != nil
}

func (g *resultBool) LastSuccessOperation() Operation {
	return g.chainable.lastSuccessOperation
}

func (g *resultBool) LastErrorOperation() Operation {
	return g.chainable.lastErrorOperation
}

func (g *resultBool) LastOperation() Operation {
	return g.chainable.lastOperation
}
