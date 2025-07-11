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
	Next  *ListItem
	Prev  *ListItem
}

type list struct {
	List
	first *ListItem
	back  *ListItem
}

//func NewList(list *list) List {
//	return new (list.n)
//}

/*
Len - получает длину списка.
*/
//func (list *list) Len() int {
//	length := list.List.Len()
//
//	return length
//}

func (list *list) Front() *ListItem {
	return list.first
}

func (list *list) Back() *ListItem {
	return list.back
}

func (list *list) PushFront(v interface{}) {
	var newItem = newItem(v)

	if list.first == nil {
		list.first = newItem
		list.back = newItem
	} else {
		list.insertBefore(list.first, newItem)
	}
}

func newItem(value interface{}) *ListItem {
	var item = new(ListItem)
	item.Value = value

	return item
}

func (list *list) insertBefore(item *ListItem, newItem *ListItem) {
	newItem.Next = item

	if item.Prev == nil {
		list.first = newItem
	} else {
		newItem.Prev = item.Prev
		item.Prev.Next = newItem
	}

	item.Prev = newItem
}
