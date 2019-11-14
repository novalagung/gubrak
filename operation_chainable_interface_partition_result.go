package gubrak

type IChainablePartitionResult interface {
	ResultAndError() (interface{}, interface{}, error)
	ResultTruthy() interface{}
	ResultFalsey() interface{}
	Error() error
	IsError() bool
}

type resultPartition struct {
	chainable *chainable
	IChainablePartitionResult
}

func (g *resultPartition) ResultAndError() (interface{}, interface{}, error) {
	return g.ResultTruthy(), g.ResultFalsey(), g.Error()
}

func (g *resultPartition) ResultTruthy() interface{} {
	if v, _ := g.chainable.data.([]interface{}); len(v) > 0 {
		return v[0]
	}

	return nil
}

func (g *resultPartition) ResultFalsey() interface{} {
	if v, _ := g.chainable.data.([]interface{}); len(v) > 1 {
		return v[1]
	}

	return nil
}

func (g *resultPartition) Error() error {
	return g.chainable.lastErrorCaught
}

func (g *resultPartition) IsError() bool {
	return g.Error() != nil
}
