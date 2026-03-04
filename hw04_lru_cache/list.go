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
	length int
	front  *ListItem
	back   *ListItem
}

func (l *list) Len() int {
	return l.length
}

func (l *list) Front() *ListItem {
	return l.front
}

func (l *list) Back() *ListItem {
	return l.back
}

func (l *list) PushFront(v interface{}) *ListItem {
	el := &ListItem{Value: v}
	if l.length == 0 {
		l.front = el
		l.back = el
	} else {
		l.front.Prev = el
		el.Next = l.front
		l.front = el
	}
	l.length++
	return el
}

func (l *list) PushBack(v interface{}) *ListItem {
	el := &ListItem{Value: v}
	if l.length == 0 {
		l.front = el
		l.back = el
	} else {
		l.back.Next = el
		el.Prev = l.back
		l.back = el
	}
	l.length++
	return el
}

func (l *list) Remove(i *ListItem) {
	switch {
	case l.front == l.back:
		l.front = nil
		l.back = nil
	case i == l.back:
		i.Prev.Next = nil
		l.back = i.Prev
		i.Prev = nil
	case i == l.front:
		i.Next.Prev = nil
		l.front = i.Next
		i.Next = nil
	default:
		i.Prev.Next = i.Next
		i.Next.Prev = i.Prev
		i.Prev = nil
		i.Next = nil
	}
	l.length--
}

func (l *list) MoveToFront(i *ListItem) {
	if l.back != l.front {
		l.Remove(i)
		l.front.Prev = i
		i.Next = l.front
		l.front = i
		l.length++
	}
}

func NewList() List {
	return new(list)
}
