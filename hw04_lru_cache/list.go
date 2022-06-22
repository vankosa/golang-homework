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
	Length    int
	FrontItem *ListItem
	BackItem  *ListItem
}

func NewList() List {
	return new(list)
}

func (l *list) Len() int {
	return l.Length
}

func (l *list) Front() *ListItem {
	return l.FrontItem
}

func (l *list) Back() *ListItem {
	return l.BackItem
}

func (l *list) PushFront(v interface{}) *ListItem {
	newElem := ListItem{
		Value: v,
		Next:  l.FrontItem,
		Prev:  nil,
	}

	if l.FrontItem != nil {
		l.FrontItem.Prev = &newElem
	}

	l.FrontItem = &newElem
	if l.Length == 0 {
		l.BackItem = &newElem
	}

	l.Length++
	return &newElem
}

func (l *list) PushBack(v interface{}) *ListItem {
	newElem := ListItem{
		Value: v,
		Next:  nil,
		Prev:  l.BackItem,
	}

	if l.BackItem != nil {
		l.BackItem.Next = &newElem
	}

	l.BackItem = &newElem
	if l.Length == 0 {
		l.FrontItem = &newElem
	}

	l.Length++
	return &newElem
}

func (l *list) Remove(i *ListItem) {
	if i.Next != nil {
		i.Next.Prev = i.Prev
	} else {
		l.BackItem = i.Prev
	}

	if i.Prev != nil {
		i.Prev.Next = i.Next
	} else {
		l.FrontItem = i.Next
	}

	l.Length--
}

func (l *list) MoveToFront(i *ListItem) {
	// clear cache
	if i == nil {
		l.Length = 0
		l.BackItem = nil
		l.FrontItem = nil
	}

	if i == l.FrontItem {
		return
	}

	// old relatives
	if i.Next != nil {
		i.Next.Prev = i.Prev
	} else {
		l.BackItem = i.Prev
	}

	if i.Prev != nil {
		i.Prev.Next = i.Next
	}

	// new relatives
	if l.FrontItem != nil {
		l.FrontItem.Prev = i
	}

	i.Prev = nil
	i.Next = l.FrontItem

	l.FrontItem = i
}
