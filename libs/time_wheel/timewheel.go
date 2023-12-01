package timewheel

import (
	"container/list"
	"errors"
	"log"
	"sync"
	"time"
)

type Task struct {
	taskId int64 // 任务ID
	circle int64 // 当前的circle
	job    func()
}

type TimeWheel struct {
	lock sync.Mutex

	curPos    int64 //当前位置
	curCircle int64 // 当前是第几圈
	soltNum   int64
	solts     []*list.List
	ticker    *time.Ticker

	interval time.Duration
	stop     chan struct{}

	taskMap map[int64]int64 // ["taskId":"pos"] ,任务以及对应的位置；
}

func NewTimeWheel(interval time.Duration, soltNum int64) *TimeWheel {
	solts := make([]*list.List, soltNum)
	for i := 0; i < int(soltNum); i++ {
		solts[i] = list.New()
	}
	ticker := time.NewTicker(interval)
	return &TimeWheel{
		lock:      sync.Mutex{},
		curPos:    0,
		curCircle: 0,
		solts:     solts,
		soltNum:   soltNum,
		interval:  interval,
		ticker:    ticker,
		stop:      make(chan struct{}),
		taskMap:   make(map[int64]int64),
	}
}

func (t *TimeWheel) Start() {
	for {
		select {
		case <-t.ticker.C:
			t.tick()
		case <-t.stop:
			log.Println("timeWheel stoped")
			return
		}
	}
}

func (t *TimeWheel) tick() {
	t.curPos = (t.curPos + 1) % t.soltNum
	if t.curPos == 0 {
		t.curCircle++
	}

	for e := t.solts[t.curPos].Front(); e != nil; {
		x := e.Value.(*Task)
		taskCircle := x.circle
		if taskCircle > t.curCircle {
			break
		}

		n := e.Next()
		t.solts[t.curPos].Remove(e)
		e = n

		go x.job()
	}

}

func (t *TimeWheel) AddTask(expired time.Duration, taskId int64, job func()) error {
	t.lock.Lock()
	defer t.lock.Unlock()

	if _, ok := t.taskMap[taskId]; ok {
		return errors.New("任务已存在")
	}

	circle, pos := t.multiColock(expired)
	task := &Task{
		taskId: taskId,
		circle: circle,
		job:    job,
	}
	insertInOrder(t.solts[pos], task)
	t.taskMap[taskId] = pos
	return nil
}

func (t *TimeWheel) StopTask(taskId int64) {
	t.lock.Lock()
	defer t.lock.Unlock()
	if pos, ok := t.taskMap[taskId]; ok {
		for e := t.solts[pos].Front(); e != nil; e = e.Next() {
			x := e.Value.(*Task)
			if x.taskId == taskId {
				t.solts[pos].Remove(e)
				delete(t.taskMap, taskId)
				return
			}
		}
	}
}

func (t *TimeWheel) Stop() {
	t.stop <- struct{}{}
}

func (t *TimeWheel) multiColock(expired time.Duration) (int64, int64) {
	waitTick := expired / t.interval // 需要等待的tick数
	pos := (int64(waitTick) + t.curPos) % t.soltNum
	circle := (int64(waitTick) / t.soltNum) + t.curCircle
	return circle, pos
}
