package gubrak

type IChainableLastIndexOfResult interface {
	ResultAndError() (int, error)
	Result() int
	Error() error
	IsError() bool
}

type resultLastIndexOf struct {
	chainable *Chainable
	IChainableLastIndexOfResult
}

func (g *resultLastIndexOf) ResultAndError() (int, error) {
	return g.Result(), g.Error()
}

func (g *resultLastIndexOf) Result() int {
	v, _ := g.chainable.data.(int)
	return v
}

func (g *resultLastIndexOf) Error() error {
	return g.chainable.lastErrorCaught
}

func (g *resultLastIndexOf) IsError() bool {
	return g.Error() != nil
}
