package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"basic_kv_cli/internal/controller"
)

func main() {
	store := controller.NewStore()
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("Simple KV Store CLI")
	fmt.Println("Commands: PUT <key> <value> <ttl_sec> | GET <key> | DELETE <key> | EXIT")

	for {
		fmt.Print("> ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)
		args := strings.Split(input, " ")

		if len(args) == 0 {
			continue
		}

		switch strings.ToUpper(args[0]) {
		case "PUT":
			if len(args) < 4 {
				fmt.Println("Usage: PUT <key> <value> <ttl_sec>")
				continue
			}
			ttl, err := strconv.Atoi(args[3])
			if err != nil {
				fmt.Println("Invalid TTL")
				continue
			}
			store.Put(args[1], args[2], ttl)
			fmt.Println("Key stored")

		case "GET":
			if len(args) < 2 {
				fmt.Println("Usage: GET <key>")
				continue
			}
			val, err := store.Get(args[1])
			if err != nil {
				fmt.Println("Error:", err)
			} else {
				fmt.Println("Value:", val)
			}

		case "DELETE":
			if len(args) < 2 {
				fmt.Println("Usage: DELETE <key>")
				continue
			}
			if err := store.Delete(args[1]); err != nil {
				fmt.Println("Error:", err)
			} else {
				fmt.Println("Key deleted")
			}

		case "EXIT":
			fmt.Println("Exiting CLI.")
			return

		default:
			fmt.Println("Unknown command:", args[0])
		}
	}
}
