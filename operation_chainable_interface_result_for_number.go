package gubrak

type IChainableNumberResult interface {
	ResultAndError() (int, error)
	Result() int
	Error() error
	IsError() bool
	LastSuccessOperation() Operation
	LastErrorOperation() Operation
	LastOperation() Operation
}

type resultNumber struct {
	chainable *Chainable
	IChainableBoolResult
}

type resultCount = resultNumber
type resultLastIndexOf = resultNumber
type resultIndexOf = resultNumber

func (g *resultNumber) ResultAndError() (int, error) {
	return g.Result(), g.Error()
}

func (g *resultNumber) Result() int {
	v, _ := g.chainable.data.(int)
	return v
}

func (g *resultNumber) Error() error {
	return g.chainable.lastErrorCaught
}

func (g *resultNumber) IsError() bool {
	return g.Error() != nil
}

func (g *resultNumber) LastSuccessOperation() Operation {
	return g.chainable.lastSuccessOperation
}

func (g *resultNumber) LastErrorOperation() Operation {
	return g.chainable.lastErrorOperation
}

func (g *resultNumber) LastOperation() Operation {
	return g.chainable.lastOperation
}
