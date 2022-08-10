package consensus

type Node struct {
	Value *Block
	Prev  *Node
	Next  *Node
}

type Queue struct {
	head *Node
	tail *Node
	len  int
}

func NewQueue() Queue {
	return Queue{
		head: nil,
		tail: nil,
		len:  0,
	}
}

func (l *Queue) Len() int {
	return l.len
}

func (l *Queue) Head() *Node {
	return l.head
}

func (l *Queue) Tail() *Node {
	return l.tail
}

func (l *Queue) Add(value *Block) {
	if l.head == nil {
		l.head = &Node{
			Value: value,
			Prev:  nil,
			Next:  nil,
		}
		l.tail = l.head
	} else {
		l.tail.Next = &Node{
			Value: value,
			Prev:  l.tail,
			Next:  nil,
		}
		l.tail = l.tail.Next
	}
	l.len++
}

func (l *Queue) Pop() *Node {
	if l.len == 0 {
		return nil
	}
	result := l.head
	if l.len == 1 {
		l.head = nil
		l.tail = nil
	} else {
		l.head = l.head.Next
		l.head.Prev = nil
	}
	l.len--
	return result
}
