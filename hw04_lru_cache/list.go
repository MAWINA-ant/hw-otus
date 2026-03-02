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
	lenght int
	front  *ListItem
	back   *ListItem
}

func (l *list) Len() int {
	return l.lenght
}

func (l *list) Front() *ListItem {
	return l.front
}

func (l *list) Back() *ListItem {
	return l.back
}

func (l *list) PushFront(v interface{}) *ListItem {
	el := &ListItem{Value: v}
	if l.lenght == 0 {
		l.front = el
		l.back = el
	} else {
		l.front.Prev = el
		el.Next = l.front
		l.front = el
	}
	l.lenght++
	return el
}

func (l *list) PushBack(v interface{}) *ListItem {
	el := &ListItem{Value: v}
	if l.lenght == 0 {
		l.front = el
		l.back = el
	} else {
		l.back.Next = el
		el.Prev = l.back
		l.back = el
	}
	l.lenght++
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
	case i == l.front:
		i.Next.Prev = nil
		l.front = i.Next
	default:
		i.Prev.Next = i.Next
		i.Next.Prev = i.Prev
	}
	l.lenght--
}

func (l *list) MoveToFront(i *ListItem) {
	l.Remove(i)
	l.PushFront(i.Value)
}

func NewList() List {
	return new(list)
}
