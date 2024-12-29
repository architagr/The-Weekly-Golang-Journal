package main

import (
	"channel-tutorial/ranging"
	singlemessage "channel-tutorial/single_message"
	"fmt"
)

func main() {
	fmt.Println("start example for single message producing and consuming")
	singlemessage.Run()
	fmt.Println("start example for multiple messages produced and consumed")
	ranging.Run()
}
