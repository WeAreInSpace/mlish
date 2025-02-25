package mlish

import "sync"

func NewModel[T any]() *Model[T] {
	mutex := new(sync.Mutex)
	waitGroup := new(sync.WaitGroup)
	return &Model[T]{
		data: []T{},
		mx:   mutex,
		wg:   waitGroup,
	}
}

type Model[T any] struct {
	data []T
	mx   *sync.Mutex
	wg   *sync.WaitGroup
}

func (m *Model[T]) Add(data ...T) {
	m.mx.Lock()
	m.data = append(m.data, data...)
	m.mx.Unlock()
}

type ForParams[T any] struct {
	data  T
	index int
	model *Model[T]
}

func (fp *ForParams[T]) GetData() T {
	return fp.data
}

func (fp *ForParams[T]) GetIndex() int {
	return fp.index
}

func (fp *ForParams[T]) AppendData(newItem T) {
	fp.model.data = append(fp.model.data, newItem)
}

func (m *Model[T]) ForEach(
	cb func(
		item *ForParams[T],
	) T,
) {
	m.mx.Lock()
	var newData []T
	for idx, data := range m.data {
		newData = append(
			newData,
			cb(&ForParams[T]{
				data:  data,
				index: idx,
			}),
		)
	}
	m.data = newData
	m.mx.Unlock()
}

func (m *Model[T]) For(
	cb func(
		item *ForParams[T],
	),
) {
	forParams := &ForParams[T]{}
	for idx, data := range m.data {
		m.mx.Lock()
		forParams.index = idx
		forParams.data = data
		cb(forParams)
		m.mx.Unlock()
	}
}

func (m *Model[T]) Filter(cb func(item *ForParams[T]) T) (filteredModel *Model[T]) {
	filteredModel = NewModel[T]()
	m.For(
		func(item *ForParams[T]) {
			filteredModel.Add(cb(item))
		},
	)
	return
}

func And[oldType any, newType any](
	model *Model[oldType],
	migrateCb func(item *ForParams[oldType]) newType,
) *Model[newType] {
	var newModel = NewModel[newType]()
	model.For(
		func(item *ForParams[oldType]) {
			var newItem = migrateCb(item)
			newModel.Add(newItem)
		},
	)
	return newModel
}
