package helper

import (
	"context"
	"errors"
	"sync"
)

const EventExample1 = "event1"

var EventDict = map[string]struct{}{
	EventExample1: {},
	// ...
}

var (
	ErrEventNotDefined = errors.New("event not defined")
)

var DefaultEventMange = &EventManage{
	g:       sync.WaitGroup{},
	lock:    sync.Mutex{},
	listens: make(map[string][]interface{}),
}

var EventManageInstance EventManageIFace = DefaultEventMange

type EventManageIFace interface {
	// 分发任务
	Dispatch(ctx context.Context, et string, params ...interface{}) []error

	// 添加监听
	AddListen(event EventIFace, l ListerIFace) error

	// 添加监听
	AddListenCallback(eventTag string, fn callback) error
}

type EventManage struct {
	g       sync.WaitGroup
	lock    sync.Mutex
	listens map[string][]interface{}
	log     Log
}

type EventIFace interface {
	GetEventName() string
}

type ListerIFace interface {
	Handle(ctx context.Context, params ...interface{}) error
	GetConsumberName() []byte
}

type callback func(context.Context, ...interface{}) error

// TODO 防呆设计 (AddListen 多次调用)
func (m *EventManage) AddListen(event EventIFace, l ListerIFace) error {
	m.lock.Lock()
	defer m.lock.Unlock()
	eventName := event.GetEventName()

	if _, ok := EventDict[eventName]; !ok {
		return ErrEventNotDefined
	}

	if _, ok := m.listens[eventName]; !ok {
		m.listens[eventName] = make([]interface{}, 0)
	}

	m.listens[eventName] = append(m.listens[eventName], l)
	return nil
}

// TODO 防呆设计 (AddListenCallback 多次调用)
func (m *EventManage) AddListenCallback(eventTag string, fn callback) error {
	m.lock.Lock()
	defer m.lock.Unlock()
	if _, ok := m.listens[eventTag]; !ok {
		m.listens[eventTag] = make([]interface{}, 0)
	}
	m.listens[eventTag] = append(m.listens[eventTag], fn)
	return nil
}

func (m *EventManage) Dispatch(ctx context.Context, et string, params ...interface{}) []error {
	errs := make([]error, len(m.listens))
	var i = 0
	for eventTag, v := range m.listens {
		if eventTag != et {
			continue
		}
		for _, vv := range v {
			switch vt := vv.(type) {
			case callback:
				m.g.Add(1)
				go func(i int) {
					defer func() {
						if r := recover(); r != nil {
							m.log.ErrorfWithContext(ctx, "inner_event, err:%v", r)
						}
					}()
					defer m.g.Done()
					if err := vt(ctx, params...); err != nil {
						errs[i] = err
					}
				}(i)
			case ListerIFace:
				m.g.Add(1)
				go func(i int) {
					defer func() {
						if r := recover(); r != nil {
							m.log.ErrorfWithContext(ctx, "inner_event, err:%v", r)
						}
					}()
					defer m.g.Done()
					if err := vt.Handle(ctx, params...); err != nil {
						errs[i] = err
					}
				}(i)
			}
		}
		i++
	}
	m.g.Wait()
	var n = 0
	for _, e := range errs {
		if e != nil {
			errs[n] = e
			n++
		}
	}
	return errs[:n]
}
