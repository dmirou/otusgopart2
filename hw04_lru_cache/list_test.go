package hw04_lru_cache //nolint:golint,stylecheck

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/require"
)

func TestList(t *testing.T) {
	t.Run("empty list", func(t *testing.T) {
		l := NewList()

		require.Equal(t, 0, l.Len())
		require.Nil(t, l.Front())
		require.Nil(t, l.Back())
	})

	t.Run("complex", func(t *testing.T) {
		l := NewList()

		l.PushFront(10) // [10]
		l.PushBack(20)  // [10, 20]
		l.PushBack(30)  // [10, 20, 30]
		require.Equal(t, 3, l.Len())

		middle := l.Front().Next // 20
		l.Remove(middle)         // [10, 30]
		require.Equal(t, 2, l.Len())

		for i, v := range [...]int{40, 50, 60, 70, 80} {
			if i%2 == 0 {
				l.PushFront(v)
			} else {
				l.PushBack(v)
			}
		} // [80, 60, 40, 10, 30, 50, 70]

		require.Equal(t, 7, l.Len())
		require.Equal(t, 80, l.Front().Value)
		require.Equal(t, 70, l.Back().Value)

		l.MoveToFront(l.Front()) // [80, 60, 40, 10, 30, 50, 70]
		l.MoveToFront(l.Back())  // [70, 80, 60, 40, 10, 30, 50]

		elems := make([]int, 0, l.Len())
		for i := l.Front(); i != nil; i = i.Next {
			elems = append(elems, i.Value.(int))
		}
		require.Equal(t, []int{70, 80, 60, 40, 10, 30, 50}, elems)
	})
}

// TestItemValue checked that the Value method correctly returns an assigned item value.
func TestItemValue(t *testing.T) {
	item := listItem{}

	values := []interface{}{
		nil,
		10,
		"str",
	}

	for _, v := range values {
		item.Value = v
		if v != item.Value {
			t.Errorf("expected value: %v, got: %v", v, item.Value)
		}
	}
}

// TestItemNext checked that the Next method correctly returns an assigned next item.
func TestItemNext(t *testing.T) {
	item := listItem{}

	nexts := []*listItem{
		nil,
		{},
		{Value: 2},
	}

	for _, next := range nexts {
		item.Next = next
		if next != item.Next {
			t.Errorf("expected next: %v, got: %v", next, item.Next)
		}
	}
}

// TestItemPrev checked that the Prev method correctly returns an assigned previous item.
func TestItemPrev(t *testing.T) {
	item := listItem{}

	prevs := []*listItem{
		nil,
		{},
		{Value: 2},
	}

	for _, prev := range prevs {
		item.Prev = prev
		if prev != item.Prev {
			t.Errorf("expected prev: %v, got: %v", prev, item.Prev)
		}
	}
}

// TestListPushFront checks that values are added to the list via PushFront method.
func TestListPushFront(t *testing.T) {
	list := NewList()
	values := []int{3, 4, 1, 2, 8}

	for _, value := range values {
		list.PushFront(value)

		if list.Front().Value != value {
			t.Errorf("expected front value: %v, got: %v", value, list.Front().Value)
		}
	}

	if list.Len() != len(values) {
		t.Errorf("expected list len: %v, got: %v", len(values), list.Len())
	}
}

// TestListPushBack checks that values are added to the list via PushBack method.
func TestListPushBack(t *testing.T) {
	list := NewList()
	values := []int{3, 4, 1, 2, 8}

	for _, value := range values {
		list.PushBack(value)

		if list.Back().Value != value {
			t.Errorf("expected back value: %v, got: %v", value, list.Back().Value)
		}
	}

	if list.Len() != len(values) {
		t.Errorf("expected list len: %v, got: %v", len(values), list.Len())
	}
}

// RemoveTestData describes input data for testing list.Remove method.
type RemoveTestData struct {
	Source        []int
	IndexToRemove int
	Result        []int
}

// TestRemove checks that a list item is removed from the list.
// nolint: funlen
func TestRemove(t *testing.T) {
	tds := []RemoveTestData{
		{
			Source:        []int{4},
			IndexToRemove: 0,
			Result:        []int{},
		},
		{
			Source:        []int{4, 2, 8, 4, 1},
			IndexToRemove: 0,
			Result:        []int{2, 8, 4, 1},
		},
		{
			Source:        []int{4, 1, 2, 10, 12, 4},
			IndexToRemove: 5,
			Result:        []int{4, 1, 2, 10, 12},
		},
		{
			Source:        []int{4, 2, 8, 1},
			IndexToRemove: 1,
			Result:        []int{4, 8, 1},
		},
		{
			Source:        []int{4, 1, 2, 10},
			IndexToRemove: 2,
			Result:        []int{4, 1, 10},
		},
	}
	for _, td := range tds {
		list := NewList()

		for _, value := range td.Source {
			list.PushBack(value)
		}

		var toRemove *listItem

		var current = list.Front()

		for i := 0; i < list.Len(); i++ {
			if i == td.IndexToRemove {
				toRemove = current
				break
			}

			current = current.Next
		}

		list.Remove(toRemove)

		var length = len(td.Result)
		if list.Len() != length {
			t.Errorf("expected length: %v, got: %v", length, list.Len())
		}

		if length == 0 {
			continue
		}

		var (
			values = make([]int, length)
			i      = 0
		)

		for cur := list.Front(); cur != nil; cur = cur.Next {
			values[i] = cur.Value.(int)
			i++
		}

		if !cmp.Equal(values, td.Result) {
			t.Errorf("expected values: %v, got: %v", td.Result, values)
		}

		values = make([]int, length)
		i = length - 1

		for curItem := list.Back(); curItem != nil; curItem = curItem.Prev {
			values[i] = curItem.Value.(int)
			i--
		}

		if !cmp.Equal(values, td.Result) {
			t.Errorf("expected values: %v, got: %v", td.Result, values)
		}
	}
}

// TestRemoveFromAnotherList checks that the list can't remove an item from a different list.
func TestRemoveFromAnotherList(t *testing.T) {
	first := NewList()
	second := NewList()
	values := []int{3, 4, 1, 2, 8}

	for _, value := range values {
		first.PushBack(value)
		second.PushBack(value)
	}

	first.Remove(second.Front())

	if first.Len() != 5 {
		t.Errorf("expected length: %d, got: %d", 5, first.Len())
	}
}

// TestMoveToFront checks that a list item is correctly moved to the front.
// nolint: funlen
func TestMoveToFront(t *testing.T) {
	list := NewList()
	values := []int{3, 4, 1, 2, 8}

	for _, value := range values {
		list.PushBack(value)
	}

	first := list.Front()
	list.MoveToFront(first)

	if first != list.Front() {
		t.Errorf("expected first item: %v, got: %v", first, list.Front())
	}

	if len(values) != list.Len() {
		t.Errorf("expected length: %v, got: %v", len(values), list.Len())
	}

	eight := list.Back()
	list.MoveToFront(eight)

	expected := []int{8, 3, 4, 1, 2}
	cur := list.Front()

	for _, v := range expected {
		if v != cur.Value {
			t.Errorf("expected value: %v, got: %v", v, cur.Value)
		}

		cur = cur.Next
	}

	cur = list.Back()

	for i := len(expected) - 1; i >= 0; i-- {
		if cur == nil {
			t.Fatalf("can not find item for value: %v", expected[i])
		}

		if expected[i] != cur.Value {
			t.Errorf("expected value: %v, got: %v", expected[i], cur.Value)
		}

		cur = cur.Prev
	}

	if len(expected) != list.Len() {
		t.Errorf("expected length: %v, got: %v", len(expected), list.Len())
	}

	four := list.Front().Next.Next
	list.MoveToFront(four)

	expected = []int{4, 8, 3, 1, 2}
	cur = list.Front()

	for _, v := range expected {
		if v != cur.Value {
			t.Errorf("expected value: %v, got: %v", v, cur.Value)
		}

		cur = cur.Next
	}

	cur = list.Back()

	for i := len(expected) - 1; i >= 0; i-- {
		if cur == nil {
			t.Fatalf("can not find item for value: %v", expected[i])
		}

		if expected[i] != cur.Value {
			t.Errorf("expected value: %v, got: %v", expected[i], cur.Value)
		}

		cur = cur.Prev
	}

	if len(expected) != list.Len() {
		t.Errorf("expected length: %v, got: %v", len(expected), list.Len())
	}
}
