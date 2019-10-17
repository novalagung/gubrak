package gubrak

type Chainable struct {
	data                 interface{}
	lastSuccessOperation Operation
	lastErrorOperation   Operation
	lastOperation        Operation
	lastErrorCaught      error
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

type IChainableCount interface {
	ResultAndError() (int, error)
	Result() int
	Error() error
	IsError() bool
}

type privateChainableCount struct {
	chainable *Chainable
	IChainableCount
}

func (g *privateChainableCount) ResultAndError() (int, error) {
	return g.Result(), g.Error()
}

func (g *privateChainableCount) Result() int {
	v, _ := g.chainable.data.(int)
	return v
}

func (g *privateChainableCount) Error() error {
	return g.chainable.lastErrorCaught
}

func (g *privateChainableCount) IsError() bool {
	return g.Error() != nil
}

type IChainable interface {
	Chunk(int) IChainable
	Compact() IChainable
	ConcatMany(...interface{}) IChainable
	Concat(interface{}) IChainable
	Count(...interface{}) IChainableCount

	ResultAndError() (interface{}, error)
	Result() interface{}
	Error() error
	IsError() bool
}

// func (g *Chainable) Difference(compareData ...interface{}) *Chainable {
// 	if g.IsError() || g.shouldReturn(OperationDifference) {
// 		return g
// 	}

// 	data, err := Difference(g.data, compareData...)
// 	if err != nil {
// 		g.lastErrorCaught = err
// 		return g
// 	}

// 	g.lastSuccessOperation = OperationDifference
// 	g.data = data
// 	g.dataObject = nil
// 	g.dataNumber = math.NaN()
// 	g.dataString = ""
// 	g.dataBool = false
// 	return g
// }

// func (g *Chainable) Drop(size int) *Chainable {
// 	if g.IsError() || g.shouldReturn(OperationDrop) {
// 		return g
// 	}

// 	data, err := Drop(g.data, size)
// 	if err != nil {
// 		g.lastErrorCaught = err
// 		return g
// 	}

// 	g.lastSuccessOperation = OperationDrop
// 	g.data = data
// 	g.dataObject = nil
// 	g.dataNumber = math.NaN()
// 	g.dataString = ""
// 	g.dataBool = false
// 	return g
// }

// func (g *Chainable) DropRight(size int) *Chainable {
// 	if g.IsError() || g.shouldReturn(OperationDropRight) {
// 		return g
// 	}

// 	data, err := DropRight(g.data, size)
// 	if err != nil {
// 		g.lastErrorCaught = err
// 		return g
// 	}

// 	g.lastSuccessOperation = OperationDropRight
// 	g.data = data
// 	g.dataObject = nil
// 	g.dataNumber = math.NaN()
// 	g.dataString = ""
// 	g.dataBool = false
// 	return g
// }

// func (g *Chainable) Each(callback interface{}) *Chainable {
// 	if g.IsError() || g.shouldReturn(OperationEach) {
// 		return g
// 	}

// 	err := Each(g.data, callback)
// 	if err != nil {
// 		g.lastErrorCaught = err
// 		return g
// 	}

// 	g.lastSuccessOperation = OperationEach
// 	g.data = nil
// 	g.dataObject = nil
// 	g.dataNumber = math.NaN()
// 	g.dataString = ""
// 	g.dataBool = false
// 	return g
// }

// func (g *Chainable) EachRight(callback interface{}) *Chainable {
// 	if g.IsError() || g.shouldReturn(OperationEachRight) {
// 		return g
// 	}

// 	err := EachRight(g.data, callback)
// 	if err != nil {
// 		g.lastErrorCaught = err
// 		return g
// 	}

// 	g.lastSuccessOperation = OperationEachRight
// 	g.data = nil
// 	g.dataObject = nil
// 	g.dataNumber = math.NaN()
// 	g.dataString = ""
// 	g.dataBool = false
// 	return g
// }

// func (g *Chainable) Fill(value interface{}, args ...int) *Chainable {
// 	if g.IsError() || g.shouldReturn(OperationFill) {
// 		return g
// 	}

// 	data, err := Fill(g.data, value, args...)
// 	if err != nil {
// 		g.lastErrorCaught = err
// 		return g
// 	}

// 	g.lastSuccessOperation = OperationFill
// 	g.data = data
// 	g.dataObject = nil
// 	g.dataNumber = math.NaN()
// 	g.dataString = ""
// 	g.dataBool = false
// 	return g
// }

// func (g *Chainable) Filter(callback interface{}) *Chainable {
// 	if g.IsError() || g.shouldReturn(OperationFilter) {
// 		return g
// 	}

// 	data, err := Filter(g.data, callback)
// 	if err != nil {
// 		g.lastErrorCaught = err
// 		return g
// 	}

// 	g.lastSuccessOperation = OperationFilter
// 	g.data = data
// 	g.dataObject = nil
// 	g.dataNumber = math.NaN()
// 	g.dataString = ""
// 	g.dataBool = false
// 	return g
// }

// func (g *Chainable) Find(callback interface{}, args ...int) *Chainable {
// 	if g.IsError() || g.shouldReturn(OperationFind) {
// 		return g
// 	}

// 	data, err := Find(g.data, callback, args...)
// 	if err != nil {
// 		g.lastErrorCaught = err
// 		return g
// 	}

// 	g.lastSuccessOperation = OperationFind
// 	g.data = nil
// 	g.dataObject = data
// 	g.dataNumber = math.NaN()
// 	g.dataString = ""
// 	g.dataBool = false
// 	return g
// }

// func (g *Chainable) FindIndex(predicate interface{}, args ...int) *Chainable {
// 	if g.IsError() || g.shouldReturn(OperationFindIndex) {
// 		return g
// 	}

// 	index, err := FindIndex(g.data, predicate, args...)
// 	if err != nil {
// 		g.lastErrorCaught = err
// 		return g
// 	}

// 	g.lastSuccessOperation = OperationFindIndex
// 	g.data = nil
// 	g.dataObject = nil
// 	g.dataNumber = float64(index)
// 	g.dataString = ""
// 	g.dataBool = false
// 	return g
// }

// func (g *Chainable) FindLast(callback interface{}, args ...int) *Chainable {
// 	if g.IsError() || g.shouldReturn(OperationFindLast) {
// 		return g
// 	}

// 	data, err := FindLast(g.data, callback, args...)
// 	if err != nil {
// 		g.lastErrorCaught = err
// 		return g
// 	}

// 	g.lastSuccessOperation = OperationFindLast
// 	g.data = nil
// 	g.dataObject = data
// 	g.dataNumber = math.NaN()
// 	g.dataString = ""
// 	g.dataBool = false
// 	return g
// }

// func (g *Chainable) FindLastIndex(predicate interface{}, args ...int) *Chainable {
// 	if g.IsError() || g.shouldReturn(OperationFindLastIndex) {
// 		return g
// 	}

// 	index, err := FindLastIndex(g.data, predicate, args...)
// 	if err != nil {
// 		g.lastErrorCaught = err
// 		return g
// 	}

// 	g.lastSuccessOperation = OperationFindLastIndex
// 	g.data = nil
// 	g.dataObject = nil
// 	g.dataNumber = float64(index)
// 	g.dataString = ""
// 	g.dataBool = false
// 	return g
// }

// func (g *Chainable) First() *Chainable {
// 	if g.IsError() || g.shouldReturn(OperationFirst) {
// 		return g
// 	}

// 	data, err := First(g.data)
// 	if err != nil {
// 		g.lastErrorCaught = err
// 		return g
// 	}

// 	g.lastSuccessOperation = OperationFirst
// 	g.data = nil
// 	g.dataObject = data
// 	g.dataNumber = math.NaN()
// 	g.dataString = ""
// 	g.dataBool = false
// 	return g
// }

// func (g *Chainable) ForEach(callback interface{}) *Chainable {
// 	if g.IsError() || g.shouldReturn(OperationForEach) {
// 		return g
// 	}

// 	err := ForEach(g.data, callback)
// 	if err != nil {
// 		g.lastErrorCaught = err
// 		return g
// 	}

// 	g.lastSuccessOperation = OperationForEach
// 	g.data = nil
// 	g.dataObject = nil
// 	g.dataNumber = math.NaN()
// 	g.dataString = ""
// 	g.dataBool = false
// 	return g
// }

// func (g *Chainable) ForEachRight(callback interface{}) *Chainable {
// 	if g.IsError() || g.shouldReturn(OperationForEachRight) {
// 		return g
// 	}

// 	err := ForEachRight(g.data, callback)
// 	if err != nil {
// 		g.lastErrorCaught = err
// 		return g
// 	}

// 	g.lastSuccessOperation = OperationForEachRight
// 	g.data = nil
// 	g.dataObject = nil
// 	g.dataNumber = math.NaN()
// 	g.dataString = ""
// 	g.dataBool = false
// 	return g
// }

// func (g *Chainable) FromPairs(callback interface{}) *Chainable {
// 	if g.IsError() || g.shouldReturn(OperationFromPairs) {
// 		return g
// 	}

// 	data, err := FromPairs(g.data)
// 	if err != nil {
// 		g.lastErrorCaught = err
// 		return g
// 	}

// 	g.lastSuccessOperation = OperationFromPairs
// 	g.data = nil
// 	g.dataObject = data
// 	g.dataNumber = math.NaN()
// 	g.dataString = ""
// 	g.dataBool = false
// 	return g
// }

// func (g *Chainable) GroupBy(callback interface{}) *Chainable {
// 	if g.IsError() || g.shouldReturn(OperationGroupBy) {
// 		return g
// 	}

// 	data, err := GroupBy(g.data, callback)
// 	if err != nil {
// 		g.lastErrorCaught = err
// 		return g
// 	}

// 	g.lastSuccessOperation = OperationGroupBy
// 	g.data = nil
// 	g.dataObject = data
// 	g.dataNumber = math.NaN()
// 	g.dataString = ""
// 	g.dataBool = false
// 	return g
// }

// func (g *Chainable) Head() *Chainable {
// 	if g.IsError() || g.shouldReturn(OperationHead) {
// 		return g
// 	}

// 	data, err := Head(g.data)
// 	if err != nil {
// 		g.lastErrorCaught = err
// 		return g
// 	}

// 	g.lastSuccessOperation = OperationHead
// 	g.data = nil
// 	g.dataObject = data
// 	g.dataNumber = math.NaN()
// 	g.dataString = ""
// 	g.dataBool = false
// 	return g
// }

// func (g *Chainable) Includes(search interface{}, args ...int) *Chainable {
// 	if g.IsError() || g.shouldReturn(OperationIncludes) {
// 		return g
// 	}

// 	yesNo, err := Includes(g.data, search, args...)
// 	if err != nil {
// 		g.lastErrorCaught = err
// 		return g
// 	}

// 	g.lastSuccessOperation = OperationIncludes
// 	g.data = nil
// 	g.dataObject = nil
// 	g.dataNumber = math.NaN()
// 	g.dataString = ""
// 	g.dataBool = yesNo
// 	return g
// }

// func (g *Chainable) IndexOf(search interface{}, args ...int) *Chainable {
// 	if g.IsError() || g.shouldReturn(OperationIndexOf) {
// 		return g
// 	}

// 	index, err := IndexOf(g.data, search, args...)
// 	if err != nil {
// 		g.lastErrorCaught = err
// 		return g
// 	}

// 	g.lastSuccessOperation = OperationIndexOf
// 	g.data = nil
// 	g.dataObject = nil
// 	g.dataNumber = float64(index)
// 	g.dataString = ""
// 	g.dataBool = false
// 	return g
// }

// func (g *Chainable) Initial() *Chainable {
// 	if g.IsError() || g.shouldReturn(OperationInitial) {
// 		return g
// 	}

// 	data, err := Initial(g.data)
// 	if err != nil {
// 		g.lastErrorCaught = err
// 		return g
// 	}

// 	g.lastSuccessOperation = OperationInitial
// 	g.data = data
// 	g.dataObject = nil
// 	g.dataNumber = math.NaN()
// 	g.dataString = ""
// 	g.dataBool = false
// 	return g
// }

// func (g *Chainable) Intersection(dataIntersects ...interface{}) *Chainable {
// 	if g.IsError() || g.shouldReturn(OperationIntersection) {
// 		return g
// 	}

// 	data, err := Intersection(g.data, dataIntersects...)
// 	if err != nil {
// 		g.lastErrorCaught = err
// 		return g
// 	}

// 	g.lastSuccessOperation = OperationIntersection
// 	g.data = data
// 	g.dataObject = nil
// 	g.dataNumber = math.NaN()
// 	g.dataString = ""
// 	g.dataBool = false
// 	return g
// }

// func (g *Chainable) Join(separator string) *Chainable {
// 	if g.IsError() || g.shouldReturn(OperationJoin) {
// 		return g
// 	}

// 	data, err := Join(g.data, separator)
// 	if err != nil {
// 		g.lastErrorCaught = err
// 		return g
// 	}

// 	g.lastSuccessOperation = OperationJoin
// 	g.data = nil
// 	g.dataObject = nil
// 	g.dataNumber = math.NaN()
// 	g.dataString = data
// 	g.dataBool = false
// 	return g
// }

// func (g *Chainable) KeyBy(callback interface{}) *Chainable {
// 	if g.IsError() || g.shouldReturn(OperationKeyBy) {
// 		return g
// 	}

// 	data, err := KeyBy(g.data, callback)
// 	if err != nil {
// 		g.lastErrorCaught = err
// 		return g
// 	}

// 	g.lastSuccessOperation = OperationKeyBy
// 	g.data = nil
// 	g.dataObject = data
// 	g.dataNumber = math.NaN()
// 	g.dataString = ""
// 	g.dataBool = false
// 	return g
// }

// func (g *Chainable) Last() *Chainable {
// 	if g.IsError() || g.shouldReturn(OperationLast) {
// 		return g
// 	}

// 	data, err := Last(g.data)
// 	if err != nil {
// 		g.lastErrorCaught = err
// 		return g
// 	}

// 	g.lastSuccessOperation = OperationLast
// 	g.data = nil
// 	g.dataObject = data
// 	g.dataNumber = math.NaN()
// 	g.dataString = ""
// 	g.dataBool = false
// 	return g
// }

// func (g *Chainable) LastIndexOf(search interface{}, args ...int) *Chainable {
// 	if g.IsError() || g.shouldReturn(OperationLastIndexOf) {
// 		return g
// 	}

// 	index, err := LastIndexOf(g.data, search, args...)
// 	if err != nil {
// 		g.lastErrorCaught = err
// 		return g
// 	}

// 	g.lastSuccessOperation = OperationLastIndexOf
// 	g.data = nil
// 	g.dataObject = nil
// 	g.dataNumber = float64(index)
// 	g.dataString = ""
// 	g.dataBool = false
// 	return g
// }

// func (g *Chainable) Map(callback interface{}) *Chainable {
// 	if g.IsError() || g.shouldReturn(OperationMap) {
// 		return g
// 	}

// 	data, err := Map(g.data, callback)
// 	if err != nil {
// 		g.lastErrorCaught = err
// 		return g
// 	}

// 	g.lastSuccessOperation = OperationMap
// 	g.data = data
// 	g.dataObject = nil
// 	g.dataNumber = math.NaN()
// 	g.dataString = ""
// 	g.dataBool = false
// 	return g
// }

// func (g *Chainable) Nth(i int) *Chainable {
// 	if g.IsError() || g.shouldReturn(OperationNth) {
// 		return g
// 	}

// 	data, err := Nth(g.data, i)
// 	if err != nil {
// 		g.lastErrorCaught = err
// 		return g
// 	}

// 	g.lastSuccessOperation = OperationNth
// 	g.data = nil
// 	g.dataObject = data
// 	g.dataNumber = math.NaN()
// 	g.dataString = ""
// 	g.dataBool = false
// 	return g
// }

// func (g *Chainable) OrderBy(callback interface{}, args ...bool) *Chainable {
// 	if g.IsError() || g.shouldReturn(OperationOrderBy) {
// 		return g
// 	}

// 	data, err := OrderBy(g.data, callback, args...)
// 	if err != nil {
// 		g.lastErrorCaught = err
// 		return g
// 	}

// 	g.lastSuccessOperation = OperationOrderBy
// 	g.data = data
// 	g.dataObject = nil
// 	g.dataNumber = math.NaN()
// 	g.dataString = ""
// 	g.dataBool = false
// 	return g
// }

// // TODO
// // Partition(data, callback interface{}) (interface{}, interface{}, error) {
// func (g *Chainable) Partition(callback interface{}) *Chainable {
// 	return g
// }

// func (g *Chainable) Pull(items ...interface{}) *Chainable {
// 	if g.IsError() || g.shouldReturn(OperationPull) {
// 		return g
// 	}

// 	data, err := Pull(g.data, items...)
// 	if err != nil {
// 		g.lastErrorCaught = err
// 		return g
// 	}

// 	g.lastSuccessOperation = OperationPull
// 	g.data = data
// 	g.dataObject = nil
// 	g.dataNumber = math.NaN()
// 	g.dataString = ""
// 	g.dataBool = false
// 	return g
// }

// func (g *Chainable) PullAll(items interface{}) *Chainable {
// 	if g.IsError() || g.shouldReturn(OperationPullAll) {
// 		return g
// 	}

// 	data, err := PullAll(g.data, items)
// 	if err != nil {
// 		g.lastErrorCaught = err
// 		return g
// 	}

// 	g.lastSuccessOperation = OperationPullAll
// 	g.data = data
// 	g.dataObject = nil
// 	g.dataNumber = math.NaN()
// 	g.dataString = ""
// 	g.dataBool = false
// 	return g
// }

// func (g *Chainable) PullAt(indexes ...int) *Chainable {
// 	if g.IsError() || g.shouldReturn(OperationPullAt) {
// 		return g
// 	}

// 	data, err := PullAt(g.data, indexes...)
// 	if err != nil {
// 		g.lastErrorCaught = err
// 		return g
// 	}

// 	g.lastSuccessOperation = OperationPullAt
// 	g.data = data
// 	g.dataObject = nil
// 	g.dataNumber = math.NaN()
// 	g.dataString = ""
// 	g.dataBool = false
// 	return g
// }

// func (g *Chainable) Reduce(callback, initial interface{}) *Chainable {
// 	if g.IsError() || g.shouldReturn(OperationReduce) {
// 		return g
// 	}

// 	data, err := Reduce(g.data, callback, initial)
// 	if err != nil {
// 		g.lastErrorCaught = err
// 		return g
// 	}

// 	g.lastSuccessOperation = OperationReduce
// 	g.data = nil
// 	g.dataObject = data
// 	g.dataNumber = math.NaN()
// 	g.dataString = ""
// 	g.dataBool = false
// 	return g
// }

// func (g *Chainable) Reject(callback interface{}) *Chainable {
// 	if g.IsError() || g.shouldReturn(OperationReject) {
// 		return g
// 	}

// 	data, err := Reject(g.data, callback)
// 	if err != nil {
// 		g.lastErrorCaught = err
// 		return g
// 	}

// 	g.lastSuccessOperation = OperationReject
// 	g.data = data
// 	g.dataObject = nil
// 	g.dataNumber = math.NaN()
// 	g.dataString = ""
// 	g.dataBool = false
// 	return g
// }

// func (g *Chainable) Remove(predicate interface{}) *Chainable {
// 	if g.IsError() || g.shouldReturn(OperationRemove) {
// 		return g
// 	}

// 	data, err := Remove(g.data, predicate)
// 	if err != nil {
// 		g.lastErrorCaught = err
// 		return g
// 	}

// 	g.lastSuccessOperation = OperationRemove
// 	g.data = data
// 	g.dataObject = nil
// 	g.dataNumber = math.NaN()
// 	g.dataString = ""
// 	g.dataBool = false
// 	return g
// }

// Remove(data interface{}, predicate interface{}) (interface{}, interface{}, error) {
// Reverse(data interface{}) (interface{}, error) {
// Sample(data interface{}) (interface{}, error) {
// SampleSize(data interface{}, take int) (interface{}, error) {
// Shuffle(data interface{}) (interface{}, error) {
// Size(data interface{}) (int, error) {
// SortBy(data, callback interface{}, args ...bool) (interface{}, error) {
// Tail(data interface{}) (interface{}, error) {
// Take(data interface{}, size int) (interface{}, error) {
// TakeRight(data interface{}, size int) (interface{}, error) {
// Union(data interface{}, slices ...interface{}) (interface{}, error) {
// Uniq(data interface{}) (interface{}, error) {
// Without(data interface{}, items ...interface{}) (interface{}, error) {

// ==================

// func (g *Chainable) Filter(callback interface{}) *Chainable {
// 	if !g.IsError() {
// 		data, err := Filter(g.data, callback)
// 		g.setResultData("Filter()", data, err)
// 	}

// 	return g
// }

// func (g *Chainable) Map(callback interface{}) *Chainable {
// 	if !g.IsError() {
// 		data, err := Map(g.data, callback)
// 		g.setResultData("Map()", data, err)
// 	}

// 	return g
// }

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

// func (g *Chainable) setResultData(op string, data interface{}, err error) {
// 	g.data = data

// 	if err != nil {
// 		g.lastErrorCaught = fmt.Errorf("error on %s method. %s", op, err.Error())
// 	}
// }

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
