package gubrak

type Operation string

const (
	OperationChunk            = "Chunk()"
	OperationCompact          = "Compact()"
	OperationConcatMany       = "ConcatMany()"
	OperationConcat           = "Concat()"
	OperationCountBy          = "CountBy()"
	OperationCount            = "Count()"
	OperationDifferenceMany   = "DifferenceMany()"
	OperationDifference       = "Difference()"
	OperationDrop             = "Drop()"
	OperationDropRight        = "DropRight()"
	OperationEach             = "Each()"
	OperationEachRight        = "EachRight()"
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
	OperationIncludes         = "Includes()"
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
	OperationPull             = "Pull()"
	OperationPullMany         = "PullMany()"
	OperationPullAt           = "PullAt()"
	OperationPullAtMany       = "PullAtMany()"
	OperationReduce           = "Reduce()"
	OperationReject           = "Reject()"
	OperationRemove           = "Remove()"
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

const OperationNone = ""

type Chainable struct {
	data                 interface{}
	lastOperation        Operation
	lastSuccessOperation Operation
	lastErrorOperation   Operation
	lastErrorCaught      error
}

type IChainable interface {
	IChainableOperation

	ResultAndError() (interface{}, error)
	Result() interface{}
	Error() error
	IsError() bool
}

type IChainableOperation interface {
	Chunk(int) IChainable
	Compact() IChainable
	ConcatMany(...interface{}) IChainable
	Concat(interface{}) IChainable
	CountBy(interface{}) IChainableCountResult
	Count() IChainableCountResult
	DifferenceMany(...interface{}) IChainable
	Difference(interface{}) IChainable
	Drop(int) IChainable
	DropRight(int) IChainable
	Each(interface{}) IChainableEachResult
	EachRight(interface{}) IChainableEachResult
	ForEach(interface{}) IChainableEachResult
	ForEachRight(interface{}) IChainableEachResult
	Fill(interface{}, ...int) IChainable
	Filter(interface{}) IChainable
	Find(interface{}, ...int) IChainable
	FindIndex(interface{}, ...int) IChainable
	FindLast(interface{}, ...int) IChainable
	FindLastIndex(interface{}, ...int) IChainable
	First() IChainable
	FromPairs() IChainable
	GroupBy(interface{}) IChainable
	Includes(interface{}, ...int) IChainableIncludesResult
	IndexOf(interface{}, ...int) IChainableIndexOfResult
	Initial(interface{}) IChainable
	Intersection(interface{}) IChainable
	IntersectionMany(data ...interface{}) IChainable
	Join(string) IChainableJoinResult
	KeyBy(interface{}) IChainable
	Last(interface{}) IChainable
	LastIndexOf(interface{}, ...int) IChainableLastIndexOfResult
	Map(interface{}) IChainable
	Nth(int) IChainable
	OrderBy(interface{}, ...bool) IChainable
	Partition(interface{}) IChainablePartitionResult
	Pull(interface{}) IChainable
	PullMany(...interface{}) IChainable
	PullAt(int) IChainable
	PullAtMany(...int) IChainable
	Reduce(interface{}, interface{}) IChainable
	Reject(interface{}) IChainable
	Remove(interface{}) IChainableRemoveResult
	Reverse() IChainable
	Sample() IChainable
	SampleSize(int) IChainable
	Shuffle() IChainable
	Size() IChainable
	Tail() IChainable
	Take(int) IChainable
	TakeRight(int) IChainable
	Uniq(interface{}) IChainable
	UnionMany(...interface{}) IChainable
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
