package hw04lrucache

type List interface {
	Len() int
	Front() *ListItem
	Back() *ListItem
	PushFront(v interface{}) *ListItem
	PushBack(v interface{}) *ListItem
	Remove(i *ListItem)
	MoveToFront(i *ListItem)
}

type ListItem struct {
	Value interface{}
	Key   Key
	Next  *ListItem
	Prev  *ListItem
}

type list struct {
	len   int
	front *ListItem
	back  *ListItem
}

func NewList() List {
	return new(list)
}

func (l *list) Len() int {
	return l.len
}

func (l *list) Front() *ListItem {
	return l.front
}

func (l *list) Back() *ListItem {
	return l.back
}

func (l *list) PushFront(v interface{}) *ListItem {
	newFrontItem := &ListItem{
		Value: v,
		Next:  l.front,
		Prev:  nil,
	}

	if l.len == 0 {
		l.back = newFrontItem
	} else {
		l.front.Prev = newFrontItem
	}

	l.front = newFrontItem
	l.len++
	return newFrontItem
}

func (l *list) PushBack(v interface{}) *ListItem {
	newBackItem := &ListItem{
		Value: v,
		Next:  nil,
		Prev:  l.back,
	}

	if l.len == 0 {
		l.front = newBackItem
	} else {
		l.back.Next = newBackItem
	}

	l.back = newBackItem
	l.len++
	return newBackItem
}

func (l *list) MoveToFront(i *ListItem) {
	switch {
	case l.len == 0, i == nil, i == l.front:
		return
	case i == l.back:
		l.back = l.back.Prev
		i.Prev.Next = nil
	default:
		i.Prev.Next = i.Next
		i.Next.Prev = i.Prev
	}

	i.Next = l.front
	i.Prev = nil
	l.front.Prev = i
	l.front = i
}

func (l *list) Remove(i *ListItem) {
	switch {
	case l.len == 0, i == nil:
		return
	case l.len == 1:
		l.front = nil
		l.back = nil
	case i == l.front:
		l.front = l.front.Next
		i.Next.Prev = nil
		i.Next = nil
	case i == l.back:
		l.back = l.back.Prev
		i.Prev.Next = nil
		i.Prev = nil
	default:
		i.Prev.Next = i.Next
		i.Next.Prev = i.Prev
		i.Next = nil
		i.Prev = nil
	}

	l.len--
}
