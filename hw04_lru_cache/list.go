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
	len   int
	front *ListItem
	back  *ListItem
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
	item := &ListItem{Value: v}
	if l.Len() == 0 {
		l.front = item
		l.back = item
	} else {
		item.Next = l.front
		l.front.Prev = item
		l.front = item
	}
	l.len++
	return item
}

func (l *list) PushBack(v interface{}) *ListItem {
	item := &ListItem{Value: v}
	if l.Len() == 0 {
		l.back = item
		l.front = item
	} else {
		item.Prev = l.back
		l.back.Next = item
		l.back = item
	}
	l.len++
	return item
}

func (l *list) Remove(i *ListItem) {
	if l.Len() == 0 {
		return
	}
	switch {
	case l.front == i && l.back == i:
		l.front = nil
		l.back = nil
	case i == l.front:
		i.Next.Prev = nil
		l.front = i.Next
		i.Next = nil
	case i == l.back:
		i.Prev.Next = nil
		l.back = i.Prev
		i.Prev = nil
	default:
		if i.Prev != nil {
			i.Prev.Next = i.Next
		}
		if i.Next != nil {
			i.Next.Prev = i.Prev
		}
		i.Next = nil
		i.Prev = nil
	}
	l.len--
}

func (l *list) MoveToFront(i *ListItem) {
	if l.Len() == 0 {
		return
	}
	if i == l.front {
		return
	}
	if i == l.back {
		i.Prev.Next = nil
		l.back = i.Prev
	} else {
		i.Prev.Next = i.Next
	}
	if i.Next != nil {
		i.Next.Prev = i.Prev
	}
	i.Prev = nil
	i.Next = l.front
	l.front.Prev = i
	l.front = i
}

func NewList() List {
	return new(list)
}
