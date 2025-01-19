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
func writeLogs(l *Logger) {
	defer func(log *Logger) {
		fmt.Println("cleanup")
		cleanup(log)
	}(l)
	ticker := time.NewTicker(1 * time.Second)
	for t := range ticker.C {
		l.Log(t.String() + "\n")
	}
}

func main() {
	// create a new log file
	file, _ := OpenLogFile()
	defer file.Close()
	// init a new logger to log to the file
	var log *Logger = NewLogger(file)

	// goroutine to write logs every second
	go writeLogs(log)
	// waits for 10 seconds and stops the program
	time.Sleep(10 * time.Second)
}
