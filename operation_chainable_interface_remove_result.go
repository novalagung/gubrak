package gubrak

type IChainableRemoveResult interface {
	ResultAndError() (interface{}, interface{}, error)
	ResultAfterRemoved() interface{}
	ResultRemovedValues() interface{}
	Error() error
	IsError() bool
}

type resultRemove struct {
	chainable *Chainable
	IChainableRemoveResult
}

func (g *resultRemove) ResultAndError() (interface{}, interface{}, error) {
	return g.ResultAfterRemoved(), g.ResultRemovedValues(), g.Error()
}

func (g *resultRemove) ResultAfterRemoved() interface{} {
	if v, _ := g.chainable.data.([]interface{}); len(v) > 0 {
		return v[0]
	}

	return nil
}

func (g *resultRemove) ResultRemovedValues() interface{} {
	if v, _ := g.chainable.data.([]interface{}); len(v) > 1 {
		return v[1]
	}

	return nil
}

func (g *resultRemove) Error() error {
	return g.chainable.lastErrorCaught
}

func (g *resultRemove) IsError() bool {
	return g.Error() != nil
}
