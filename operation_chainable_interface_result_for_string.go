package gubrak

type IChainableStringResult interface {
	ResultAndError() (string, error)
	Result() string
	Error() error
	IsError() bool
}

type resultString struct {
	chainable *Chainable
	IChainableStringResult
}

type resultJoin = resultString

func (g *resultString) ResultAndError() (string, error) {
	return g.Result(), g.Error()
}

func (g *resultString) Result() string {
	v, _ := g.chainable.data.(string)
	return v
}

func (g *resultString) Error() error {
	return g.chainable.lastErrorCaught
}

func (g *resultString) IsError() bool {
	return g.Error() != nil
}

func (g *resultString) LastSuccessOperation() Operation {
	return g.chainable.lastSuccessOperation
}

func (g *resultString) LastErrorOperation() Operation {
	return g.chainable.lastErrorOperation
}

func (g *resultString) LastOperation() Operation {
	return g.chainable.lastOperation
}
