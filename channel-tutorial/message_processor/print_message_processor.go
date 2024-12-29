package messageprocessor

import (
	"fmt"
	"io"
)

type PrintMessage struct {
	w io.Writer
}

func (processor *PrintMessage) Push(data string) {
	fmt.Fprintln(processor.w, "PrintMessage:", data)
}

func InitPrintMessage(w io.Writer) *PrintMessage {
	return &PrintMessage{
		w: w,
	}
}
