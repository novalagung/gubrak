package gubrak

type IChainableJoinResult interface {
	ResultAndError() (string, error)
	Result() string
	Error() error
	IsError() bool
}

type resultJoin struct {
	chainable *chainable
	IChainableJoinResult
}

func (g *resultJoin) ResultAndError() (string, error) {
	return g.Result(), g.Error()
}

func (g *resultJoin) Result() string {
	v, _ := g.chainable.data.(string)
	return v
}

func (g *resultJoin) Error() error {
	return g.chainable.lastErrorCaught
}

func (g *resultJoin) IsError() bool {
	return g.Error() != nil
}
