package main

import (
	"log"
	"time"

	timewheel "github.com/SuperJourney/tools/libs/time_wheel"
)

func main() {
	timeWheel := timewheel.NewTimeWheel(time.Second, 10)
	timeWheel.AddTask(10*time.Second, 1, func() {
		log.Println("after 10 seconds, task 1 is done")
	})
	log.Println("time wheel start")
	addTaskAfter10Seconds(timeWheel)
	stopTimeWheelAfter30Seconds(timeWheel)
	timeWheel.Start()

}

func addTaskAfter10Seconds(timeWheel *timewheel.TimeWheel) {
	c := time.After(7 * time.Second)
	go func() {
		for range c {
			timeWheel.AddTask(10*time.Second, 2, func() {
				log.Println("after 17 seconds, task 2 is done")
			})
		}
	}()
}

func stopTimeWheelAfter30Seconds(timeWheel *timewheel.TimeWheel) {
	c := time.After(30 * time.Second)
	go func() {
		for range c {
			log.Println("stop")
			timeWheel.Stop()
			return
		}
	}()
}
