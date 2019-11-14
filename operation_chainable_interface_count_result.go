package gubrak

type IChainableCountResult interface {
	ResultAndError() (int, error)
	Result() int
	Error() error
	IsError() bool
}

type resultCount struct {
	chainable *chainable
	IChainableCountResult
}

func (g *resultCount) ResultAndError() (int, error) {
	return g.Result(), g.Error()
}

func (g *resultCount) Result() int {
	v, _ := g.chainable.data.(int)
	return v
}

func (g *resultCount) Error() error {
	return g.chainable.lastErrorCaught
}

func (g *resultCount) IsError() bool {
	return g.Error() != nil
}
