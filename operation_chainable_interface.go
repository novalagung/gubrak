package gubrak

// Operation represent the type of chainable operation
type Operation string

const (
	OperationNone             = ""
	OperationChunk            = "Chunk()"
	OperationCompact          = "Compact()"
	OperationConcatMany       = "ConcatMany()"
	OperationConcat           = "Concat()"
	OperationContains         = "Contains()"
	OperationCountBy          = "CountBy()"
	OperationCount            = "Count()"
	OperationDifferenceMany   = "DifferenceMany()"
	OperationDifference       = "Difference()"
	OperationDrop             = "Drop()"
	OperationDropRight        = "DropRight()"
	OperationEach             = "Each()"
	OperationEachRight        = "EachRight()"
	OperationExclude          = "Exclude()"
	OperationExcludeMany      = "ExcludeMany()"
	OperationExcludeAt        = "ExcludeAt()"
	OperationExcludeAtMany    = "ExcludeAtMany()"
	OperationForEach          = "ForEach()"
	OperationForEachRight     = "ForEachRight()"
	OperationFill             = "Fill()"
	OperationFilter           = "Filter()"
	OperationFind             = "Find()"
	OperationFindIndex        = "FindIndex()"
	OperationFindLast         = "FindLast()"
	OperationFindLastIndex    = "FindLastIndex()"
	OperationFirst            = "First()"
	OperationHead             = "Head()"
	OperationFromPairs        = "FromPairs()"
	OperationGroupBy          = "GroupBy()"
	OperationIndexOf          = "IndexOf()"
	OperationInitial          = "Initial()"
	OperationIntersection     = "Intersection()"
	OperationIntersectionMany = "IntersectionMany()"
	OperationJoin             = "Join()"
	OperationKeyBy            = "KeyBy()"
	OperationLast             = "Last()"
	OperationLastIndexOf      = "LastIndexOf()"
	OperationMap              = "Map()"
	OperationNth              = "Nth()"
	OperationOrderBy          = "OrderBy()"
	OperationPartition        = "Partition()"
	OperationReduce           = "Reduce()"
	OperationReject           = "Reject()"
	OperationReverse          = "Reverse()"
	OperationSample           = "Sample()"
	OperationSampleSize       = "SampleSize()"
	OperationShuffle          = "Shuffle()"
	OperationSize             = "Size()"
	OperationTail             = "Tail()"
	OperationTake             = "Take()"
	OperationTakeRight        = "TakeRight()"
	OperationUniq             = "Uniq()"
	OperationUnionMany        = "UnionMany()"
)

// IChainable is the base interface for chainable functions
// It is contain the `IChainableOperation` interface (embedded), and result-related methods
type IChainable interface {
	IChainableOperation

	ResultAndError() (interface{}, error)
	Result() interface{}
	Error() error
	IsError() bool
	LastSuccessOperation() Operation
	LastErrorOperation() Operation
	LastOperation() Operation
}

// IChainableOperation is interface for chainable functions declaration
type IChainableOperation interface {
	Chunk(int) IChainable
	Compact() IChainable
	ConcatMany(...interface{}) IChainable
	Concat(interface{}) IChainable
	CountBy(interface{}) IChainableNumberResult
	Count() IChainableNumberResult
	DifferenceMany(...interface{}) IChainable
	Difference(interface{}) IChainable
	Drop(int) IChainable
	DropRight(int) IChainable
	Each(interface{}) IChainableNoReturnValueResult
	EachRight(interface{}) IChainableNoReturnValueResult
	Exclude(interface{}) IChainable
	ExcludeMany(...interface{}) IChainable
	ExcludeAt(int) IChainable
	ExcludeAtMany(...int) IChainable
	Fill(interface{}, ...int) IChainable
	Filter(interface{}) IChainable
	Find(interface{}, ...int) IChainable
	FindIndex(interface{}, ...int) IChainable
	FindLast(interface{}, ...int) IChainable
	FindLastIndex(interface{}, ...int) IChainable
	First() IChainable
	FromPairs() IChainable
	GroupBy(interface{}) IChainable
	Contains(interface{}, ...int) IChainableBoolResult
	IndexOf(interface{}, ...int) IChainableNumberResult
	Initial() IChainable
	Intersection(interface{}) IChainable
	IntersectionMany(data ...interface{}) IChainable
	Join(string) IChainableStringResult
	KeyBy(interface{}) IChainable
	Last() IChainable
	LastIndexOf(interface{}, ...int) IChainableNumberResult
	Map(interface{}) IChainable
	Nth(int) IChainable
	OrderBy(interface{}, ...bool) IChainable
	Partition(interface{}) IChainableTwoReturnValueResult
	Reduce(interface{}, interface{}) IChainable
	Reject(interface{}) IChainable
	Reverse() IChainable
	Sample() IChainable
	SampleSize(int) IChainable
	Shuffle() IChainable
	Size() IChainable
	Tail() IChainable
	Take(int) IChainable
	TakeRight(int) IChainable
	Uniq() IChainable
	UnionMany(...interface{}) IChainable
}

// Chainable is base type of gubrak chainable operations
type Chainable struct {
	data                 interface{}
	lastOperation        Operation
	lastSuccessOperation Operation
	lastErrorOperation   Operation
	lastErrorCaught      error
}

// From is the initial function to use gubrak chainable operation.
// This function requires one argument, the data that are going to be used in operations
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

// ResultAndError returns the result after operation, and error object
func (g *Chainable) ResultAndError() (interface{}, error) {
	return g.Result(), g.Error()
}

// Result returns the result after operation
func (g *Chainable) Result() interface{} {
	return g.data
}

// Error returns the error object
func (g *Chainable) Error() error {
	return g.lastErrorCaught
}

// IsError `true` on error, otherwise `false`
func (g *Chainable) IsError() bool {
	return g.Error() != nil
}

// LastSuccessOperation return last success operation
func (g *Chainable) LastSuccessOperation() Operation {
	return g.lastSuccessOperation
}

// LastErrorOperation return last error operation
func (g *Chainable) LastErrorOperation() Operation {
	return g.lastErrorOperation
}

// LastOperation return last operation
func (g *Chainable) LastOperation() Operation {
	return g.lastOperation
}
