package gubrak

type IChainableTwoReturnValueResult interface {
	ResultAndError() (any, any, error)
	ResultTruthy() any
	ResultFalsey() any
	Error() error
	IsError() bool
}

type resultTwoReturnValue struct {
	chainable *Chainable
	IChainableTwoReturnValueResult
}

type resultPartition = resultTwoReturnValue

func (g *resultTwoReturnValue) ResultAndError() (any, any, error) {
	return g.ResultTruthy(), g.ResultFalsey(), g.Error()
}

func (g *resultTwoReturnValue) ResultTruthy() any {
	if v, _ := g.chainable.data.([]any); len(v) > 0 {
		return v[0]
	}

	return nil
}

func (g *resultTwoReturnValue) ResultFalsey() any {
	if v, _ := g.chainable.data.([]any); len(v) > 1 {
		return v[1]
	}

	return nil
}

func (g *resultTwoReturnValue) Error() error {
	return g.chainable.lastErrorCaught
}

func (g *resultTwoReturnValue) IsError() bool {
	return g.Error() != nil
}

func (g *resultTwoReturnValue) LastSuccessOperation() Operation {
	return g.chainable.lastSuccessOperation
}

func (g *resultTwoReturnValue) LastErrorOperation() Operation {
	return g.chainable.lastErrorOperation
}

func (g *resultTwoReturnValue) LastOperation() Operation {
	return g.chainable.lastOperation
}
