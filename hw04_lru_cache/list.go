package hw04_lru_cache //nolint:golint,stylecheck

type List interface {
	Len() int
	Front() *listItem
	Back() *listItem
	PushFront(value interface{}) *listItem
	PushBack(value interface{}) *listItem
	Remove(item *listItem)
	MoveToFront(item *listItem)
}

type listItem struct {
	Value interface{}
	list  *list
	Next  *listItem
	Prev  *listItem
}

// newItem creates a new item with the value.
func newItem(value interface{}) *listItem {
	newItem := new(listItem)
	newItem.Value = value

	return newItem
}

type list struct {
	front  *listItem
	back   *listItem
	length int
}

func NewList() List {
	return &list{}
}

// Len returns a count of elements in the list.
func (l list) Len() int {
	return l.length
}

// Front returns a front item of the list.
func (l list) Front() *listItem {
	return l.front
}

// Back returns a back item of the list.
func (l list) Back() *listItem {
	return l.back
}

// PushFront adds a value to the beginning of the list.
func (l *list) PushFront(value interface{}) *listItem {
	item := newItem(value)
	item.list = l

	if l.front == nil {
		l.front = item
		l.back = item
		l.length++

		return item
	}

	l.front.Prev = item
	item.Next = l.front
	l.front = item
	l.length++

	return item
}

// PushBack adds a value to the end of the list.
func (l *list) PushBack(value interface{}) *listItem {
	item := newItem(value)
	item.list = l

	if l.front == nil {
		l.front = item
		l.back = item
		l.length++

		return item
	}

	l.back.Next = item
	item.Prev = l.back
	l.back = item
	l.length++

	return item
}

// Remove removes an item from the list.
// If the item doesn't belong to the list, nothing will happen.
func (l *list) Remove(item *listItem) {
	if item.list != l {
		return
	}

	if item.Prev == nil && item.Next == nil {
		item.list = nil
		l.front = nil
		l.back = nil
		l.length = 0

		return
	}

	item.list = nil

	if l.Front() == item {
		l.front = item.Next
	}

	if l.Back() == item {
		l.back = item.Prev
	}

	if item.Prev != nil {
		item.Prev.Next = item.Next
	}

	if item.Next != nil {
		item.Next.Prev = item.Prev
	}

	l.length--
}

// Move an item to the beginning of the list.
// If the item doesn't belong to the list, nothing will happen.
func (l *list) MoveToFront(item *listItem) {
	if item.list != l {
		return
	}

	if item == l.front {
		return
	}

	item.Prev.Next = item.Next

	if item.Next != nil {
		item.Next.Prev = item.Prev
	}

	if item == l.back && item.Next != nil {
		l.back = item.Next
	} else if item == l.back {
		l.back = item.Prev
	}

	item.Prev = nil
	item.Next = l.front
	l.front.Prev = item
	l.front = item
}
