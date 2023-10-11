package timewheel

import (
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestTimeWheel_multiColock(t *testing.T) {
	timeWheel := NewTimeWheel(time.Second, 10)

	circle, pos := timeWheel.multiColock(10 * time.Second)
	assert.Equal(t, circle, int64(1))
	assert.Equal(t, pos, int64(0))

	for i := 0; i < 3; i++ {
		timeWheel.tick()
	}

	circle, pos = timeWheel.multiColock(10 * time.Second)
	assert.Equal(t, circle, int64(1))
	assert.Equal(t, pos, int64(3))

}

var i = 0

func TestTimeWheel_Tick(t *testing.T) {
	timeWheel := NewTimeWheel(time.Second, 10)
	var lock sync.Locker = &sync.Mutex{}

	add1 := func() {
		lock.Lock()
		defer lock.Unlock()
		i = i + 1
	}

	timeWheel.AddTask(10*time.Second, 1, add1)
	timeWheel.AddTask(10*time.Second, 2, add1)
	timeWheel.StopTask(2)
	timeWheel.AddTask(10*time.Second, 2, add1)

	for i := 0; i < 9; i++ {
		timeWheel.tick()
	}

	assert.Equal(t, i, 0)
	timeWheel.tick()
	time.Sleep(10 * time.Millisecond)
	assert.Equal(t, i, 2)
	timeWheel.tick()
	assert.Equal(t, i, 2)
	timeWheel.AddTask(10*time.Second, 3, add1)
	for i := 0; i < 10; i++ {
		timeWheel.tick()
	}
	time.Sleep(10 * time.Millisecond)
	assert.Equal(t, i, 3)
}
