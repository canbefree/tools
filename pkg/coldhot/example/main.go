package main

import (
	"fmt"
	"os"
	"time"

	"bufio"
)

type TextFile interface {
	Read() ([]byte, error)
	Write(data []byte) error
}

type TxtFile struct {
	filename string
	file     *os.File
}

func (t *TxtFile) Open() (err error) {
	t.file, err = os.OpenFile(t.filename, os.O_RDWR|os.O_CREATE, 0755)
	return err
}

func (t *TxtFile) Close() error {
	return t.file.Close()
}

func (t *TxtFile) Read() (data []byte, err error) {
	reader := bufio.NewReader(t.file)
	data, err = reader.ReadBytes('\n')
	return data, err
}

func (t *TxtFile) Write(data []byte) error {
	writer := bufio.NewWriter(t.file)
	_, err := writer.Write(data)
	if err != nil {
		return err
	}
	return writer.Flush()
}

type Queue struct {
	index int
	data  chan int
}

func (q *Queue) Push(val int) {
	q.data <- val
	q.index++
}

func main() {
	q := Queue{data: make(chan int, 10)} // create a queue with channel

	go func() {
		for i := 0; i < 10; i++ {
			q.Push(i)
			fmt.Println("Pushed", i, "to queue")
			time.Sleep(time.Millisecond * 100) // simulate async behavior
		}
		close(q.data)
	}()

	// Simulate async behavior of reading from the queue
	for val := range q.data {
		fmt.Println("Popped", val, "from queue")
		time.Sleep(time.Millisecond * 100) // simulate processing time
	}

}
