package timewheel

import (
	"container/list"
)

func insertInOrder(l *list.List, task *Task) *list.List {
	element := l.Back()
	for element != nil && element.Value.(*Task).circle >= task.circle {
		element = element.Prev()
	}
	if element == nil {
		l.PushFront(task)
	} else {
		l.InsertAfter(task, element)
	}
	return l
}

// const interval = 1 * time.Second
// const soltNum = 60

// func multiColock(expired time.Duration) int64 {
// 	waitTick := expired / interval // 需要等待的tick数
// 	circle := waitTick / soltNum
// 	return int64(circle)
// }
