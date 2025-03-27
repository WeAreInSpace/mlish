package mlish

import (
	"fmt"
	"io"
	"os"
	"regexp"
	"slices"
	"sync"
)

func NewModel[T any]() *Model[T] {
	mutex := new(sync.Mutex)
	waitGroup := new(sync.WaitGroup)
	return &Model[T]{
		data: []*T{},
		mx:   mutex,
		wg:   waitGroup,
	}
}

type Model[T any] struct {
	data []*T
	mx   *sync.Mutex
	wg   *sync.WaitGroup
}

func (m *Model[T]) Just() []*T {
	return m.data
}

func (m *Model[T]) JustCopy() *Model[T] {
	newModel := NewModel[T]()
	m.For(
		func(item *ForParams[T]) {
			data := item.Data()
			newModel.Add(&data)
		},
	)
	return newModel
}

func (m *Model[T]) Add(data ...*T) {
	m.mx.Lock()
	m.data = append(m.data, data...)
	m.mx.Unlock()
}

type ForParams[T any] struct {
	data  *T
	index int
	model *Model[T]
}

func (fp *ForParams[T]) Data() T {
	return *fp.data
}

/*
<!> Should not use to set data.

use `return` or `[].Add()` or `[].Append()`
*/
func (fp *ForParams[T]) DataAddr() *T {
	return fp.data
}

func (fp *ForParams[T]) Index() int {
	return fp.index
}

func (fp *ForParams[T]) Append(newItem *T) {
	fp.model.Add(newItem)
}

func (m *Model[T]) ForEach(
	cb func(
		item *ForParams[T],
	) *T,
) {
	var newData []*T
	m.For(func(item *ForParams[T]) {
		newData = append(
			newData,
			cb(&ForParams[T]{
				data:  item.data,
				index: item.index,
			}),
		)
	})
	m.data = newData
}

func (m *Model[T]) JustForEach(
	cb func(
		item *ForParams[T],
	) *T,
) []*T {
	m.ForEach(
		cb,
	)
	return m.data
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

func (m *Model[T]) JustFor(
	cb func(
		item *ForParams[T],
	),
) []*T {
	m.For(
		cb,
	)
	return m.data
}

func (m *Model[T]) Filter(cb func(item *ForParams[T]) *T) (filteredModel *Model[T]) {
	filteredModel = NewModel[T]()
	m.For(
		func(item *ForParams[T]) {
			filteredItem := cb(item)
			if filteredItem != nil {
				filteredModel.Add(filteredItem)
			}
		},
	)
	return
}

func (m *Model[T]) FilterByRegex(regex string, cb func(item *ForParams[T]) string) (filteredModel *Model[T]) {
	filteredModel = m.Filter(
		func(item *ForParams[T]) *T {
			key := cb(item)
			matched, err := regexp.MatchString(regex, key)
			if err != nil {
				if Settings.DebugMode {
					fmt.Sprintln(Settings.Out, err)
					os.Exit(0)
				} else {
					return nil
				}
			}

			if matched {
				return item.DataAddr()
			}

			return nil
		},
	)

	return
}

func (m *Model[T]) Push(w io.Writer, cb func(item *ForParams[T]) []byte) {
	m.For(
		func(item *ForParams[T]) {
			data := cb(item)
			_, err := w.Write(data)
			if err != nil {
				if Settings.DebugMode {
					fmt.Sprintln(Settings.Out, err)
					os.Exit(0)
				} else {
					return
				}
			}
		},
	)
}

func Migrate[oldType any, newType any](
	model *Model[oldType],
	migrateCb func(item *ForParams[oldType]) *newType,
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

func Remove[T any](model *Model[T]) {
	model.data = slices.Delete(model.data, 0, len(model.data))
}
