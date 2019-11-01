package gubrak

type Operation string

const OperationNone = ""

type Chainable struct {
	data                 interface{}
	lastOperation        Operation
	lastSuccessOperation Operation
	lastErrorOperation   Operation
	lastErrorCaught      error
}

type IChainable interface {
	IChainableChunk
	IChainableCompact
	IChainableConcat
	IChainableCount
	IChainableDifference
	IChainableDrop
	IChainableEach
	IChainableFill
	IChainableFilter
	IChainableFind
	IChainableFirst
	IChainableFromPairs
	IChainableGroupBy
	IChainableIncludes
	IChainableIndexOf
	IChainableInitial
	IChainableIntersection
	IChainableJoin
	IChainableKeyBy

	ResultAndError() (interface{}, error)
	Result() interface{}
	Error() error
	IsError() bool
}

func From(data interface{}) IChainable {
	g := new(Chainable)
	g.data = data
	g.lastSuccessOperation = OperationNone
	g.lastErrorOperation = OperationNone
	g.lastOperation = OperationNone
	g.lastErrorCaught = nil
	return g
}

func (g *Chainable) markError(data interface{}, err error) *Chainable {
	g.data = data
	g.lastErrorCaught = err
	g.lastErrorOperation = g.lastOperation
	return g
}

func (g *Chainable) markResult(data interface{}) *Chainable {
	g.data = data
	g.lastSuccessOperation = g.lastOperation
	return g
}

func (g *Chainable) shouldReturn() bool {
	return false
}

func (g *Chainable) ResultAndError() (interface{}, error) {
	return g.Result(), g.Error()
}

func (g *Chainable) Result() interface{} {
	return g.data
}

func (g *Chainable) Error() error {
	return g.lastErrorCaught
}

func (g *Chainable) IsError() bool {
	return g.Error() != nil
}
