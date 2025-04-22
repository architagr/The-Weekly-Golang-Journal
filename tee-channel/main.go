package main

import (
	"context"
	"fmt"
	"time"
)

type userDTO struct {
	id    int
	name  string
	email string
}

func orDone(done <-chan struct{}, c <-chan userDTO) <-chan userDTO {
	valStream := make(chan userDTO)
	go func() {
		defer close(valStream)
		for {
			select {
			case <-done:
				return
			case v, ok := <-c:
				if ok == false {
					return
				}
				select {
				case valStream <- v:
				case <-done:
				}
			}
		}
	}()
	return valStream
}

func tee(
	done <-chan struct{},
	in <-chan userDTO,
) (<-chan userDTO, <-chan userDTO) {
	out1 := make(chan userDTO, 1_000)
	out2 := make(chan userDTO, 1_000)
	go func() {
		defer close(out1)
		defer close(out2)
		for val := range orDone(done, in) {
			var out1, out2 = out1, out2
			for i := 0; i < 2; i++ {
				select {
				case <-done:
				case out1 <- val:
					out1 = nil
				case out2 <- val:
					out2 = nil
				}
			}
		}
	}()
	return out1, out2
}

func sendWelcomeEmail(done <-chan struct{}, ch <-chan userDTO) {
	// we can have fan-out if needed in future

	for u := range orDone(done, ch) {
		fmt.Println("Sending email for the user", u)
	}
	fmt.Println("send welcome email completed")
}

func sendNotification(done <-chan struct{}, ch <-chan userDTO, slackChannel string) {
	// we can have fan-out if needed in future

	for u := range orDone(done, ch) {
		fmt.Println("Sending notification to marketing team on slack ", slackChannel, " for the user", u)
	}

	fmt.Println("send notification completed")
}

func singnup(usr userDTO, notificationChannel chan<- userDTO) {
	// save in the DB
	notificationChannel <- usr
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	notificationChannel := make(chan userDTO, 1_000)

	// Split the input to two output handlers
	ch1, ch2 := tee(ctx.Done(), notificationChannel)

	// Start handlers
	go sendNotification(ctx.Done(), ch1, "marketing")
	go sendWelcomeEmail(ctx.Done(), ch2)

	// Simulate signups
	singnup(userDTO{id: 1, name: "Archit", email: "a@example.com"}, notificationChannel)
	singnup(userDTO{id: 2, name: "Lilly", email: "lilly@example.com"}, notificationChannel)

	time.Sleep(2 * time.Second) // Let handlers finish
}
