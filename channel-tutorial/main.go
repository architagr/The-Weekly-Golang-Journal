package main

import (
	_ "channel-tutorial/ranging"
	_ "channel-tutorial/single_message"
	"channel-tutorial/unidirectional"
	_ "channel-tutorial/unidirectional"
)

func main() {
	// fmt.Println("start example for single message producing and consuming")
	// singlemessage.Run()
	// fmt.Println("start example for multiple messages produced and consumed")
	// ranging.Run()
	unidirectional.Run()
}
