package main

import (
	"fmt"
	"io"
	"os"
	"time"
)

type Logger struct {
	w io.Writer
}

func NewLogger(w io.Writer) *Logger {
	return &Logger{w: w}
}

func (l *Logger) Log(message string) {
	l.w.Write([]byte(message))
}

func OpenLogFile() (*os.File, error) {
	name := fmt.Sprintf("log-%d.txt", time.Now().Unix())
	return os.Create(name)
}

func cleanup(l *Logger) {
	// do some cleanup
	l.Log("cleanup completed")
}
func writeLogs(done <-chan interface{}, l *Logger) {
	defer func(log *Logger) {
		fmt.Println("cleanup")
		cleanup(log)
	}(l)
	ticker := time.NewTicker(1 * time.Second)
	for {
		select {
		case <-done:
			return
		case t := <-ticker.C:
			l.Log(t.String() + "\n")
		}
	}
}

func main() {
	done := make(chan interface{})
	// create a new log file
	file, _ := OpenLogFile()
	defer file.Close()
	// init a new logger to log to the file
	var log *Logger = NewLogger(file)

	// goroutine to write logs every second
	go writeLogs(done, log)
	// waits for 10 seconds and stops the program
	time.Sleep(10 * time.Second)
	close(done)
	// simulate ongoing work
	time.Sleep(500 * time.Millisecond)
}
