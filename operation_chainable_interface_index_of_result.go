package gubrak

type IChainableIndexOfResult interface {
	ResultAndError() (int, error)
	Result() int
	Error() error
	IsError() bool
}

type resultIndexOf struct {
	chainable *chainable
	IChainableIndexOfResult
}

func (g *resultIndexOf) ResultAndError() (int, error) {
	return g.Result(), g.Error()
}

func (g *resultIndexOf) Result() int {
	v, _ := g.chainable.data.(int)
	return v
}

func (g *resultIndexOf) Error() error {
	return g.chainable.lastErrorCaught
}

func (g *resultIndexOf) IsError() bool {
	return g.Error() != nil
}
