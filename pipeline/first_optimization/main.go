package main

import (
	"fmt"
	"sync"
	"time"
)

type taskTypeEnum string

const (
	taskTypeEpic     taskTypeEnum = "Epic"
	taskTypeStory    taskTypeEnum = "Story"
	taskTypeSubTask  taskTypeEnum = "SubTask"
	taskTypeTechTask taskTypeEnum = "TechTask"
)

type PutTask struct {
	taskId  int
	toType  taskTypeEnum
	summary string
}

type f func(task PutTask) bool

func startProcessing(done <-chan bool, arr []PutTask, worker f) (success, failed <-chan PutTask) {
	validatedDataStream := make(chan PutTask, len(arr))
	failedDataStream := make(chan PutTask, len(arr))
	go func() {
		pushWg := &sync.WaitGroup{}
		for _, task := range arr {
			pushWg.Add(1)
			go func(wg *sync.WaitGroup, t PutTask) {
				defer wg.Done()
				// here we are simulating the processing time
				ticker := time.After(100 * time.Millisecond)
				for {
					select {
					case <-done:
						return
					case <-ticker:
						// here we are simulating the worker function,
						// that is passed as input to the startProcessing
						// function for each task
						if worker(t) {
							validatedDataStream <- t
							return
						}
						failedDataStream <- t
						return
					}
				}

			}(pushWg, task)
		}
		pushWg.Wait()
		close(validatedDataStream)
		close(failedDataStream)
	}()
	return validatedDataStream, failedDataStream
}

func readData(done <-chan bool, successDataStream, failedDataStream <-chan PutTask) (successData, failedData []PutTask) {
	successData, failedData = make([]PutTask, 0), make([]PutTask, 0)
	for {
		// if both the channels are closed, then return
		if successDataStream == nil && failedDataStream == nil {
			return
		}
		select {
		case task, ok := <-successDataStream:
			if !ok {
				successDataStream = nil
				continue
			}
			successData = append(successData, task)
		case task, ok := <-failedDataStream:
			if !ok {
				failedDataStream = nil
				continue
			}
			failedData = append(failedData, task)
		case <-done:
			return
		}
	}
}

func validateTasks(done <-chan bool, arr []PutTask) (success, faulty []PutTask) {
	validatedDataStream, failedDataStream := startProcessing(done, arr, func(t PutTask) bool {
		// creating a arbitary fault condition
		if t.taskId == 1 {
			fmt.Println("Validation failed for task id: ", t.taskId)
			return false
		}
		fmt.Println("Validation done for task id: ", t.taskId)
		return true
	})
	return readData(done, validatedDataStream, failedDataStream)
}

func createAudit(done <-chan bool, arr []PutTask) (success, faulty []PutTask) {
	validatedDataStream, failedDataStream := startProcessing(done, arr, func(t PutTask) bool {
		// creating a arbitary fault condition
		if t.taskId == 2 {
			fmt.Println("Audit failed for task id: ", t.taskId)
			return false
		}
		fmt.Println("Audit done for task id: ", t.taskId)
		return true
	})
	return readData(done, validatedDataStream, failedDataStream)
}

func updateTasks(done <-chan bool, arr []PutTask) (success, faulty []PutTask) {
	validatedDataStream, failedDataStream := startProcessing(done, arr, func(t PutTask) bool {
		// creating a arbitary fault condition
		if t.taskId == 3 {
			fmt.Println("Update failed for task id: ", t.taskId)
			return false
		}
		fmt.Println("Update done for task id: ", t.taskId)
		return true
	})
	return readData(done, validatedDataStream, failedDataStream)
}

func NotifyForTasks(done <-chan bool, arr []PutTask) (success []PutTask, faulty []PutTask) {
	validatedDataStream, failedDataStream := startProcessing(done, arr, func(t PutTask) bool {
		// creating a arbitary fault condition
		if t.taskId == 4 {
			fmt.Println("Notify failed for task id: ", t.taskId)
			return false
		}
		fmt.Println("Notify done for task id: ", t.taskId)
		return true
	})
	return readData(done, validatedDataStream, failedDataStream)
}

func main() {
	// initilize server and handler
	// start server
	// handle request
	// get tasks from request
	startTime := time.Now()
	listOfTasks := []PutTask{
		{taskId: 1, toType: taskTypeEpic},
		{taskId: 2, toType: taskTypeStory},
		{taskId: 3, toType: taskTypeSubTask},
		{taskId: 4, toType: taskTypeTechTask},
		{taskId: 5, summary: "This is a updated task summary for task 5"},
		{taskId: 6, summary: "The updated summary for task 6"},
	}
	done := make(chan bool)
	defer close(done)
	fmt.Println("------------ Validation Started ------------")
	validatedTasks, validationFailed := validateTasks(done, listOfTasks)
	fmt.Println("validation success", validatedTasks)
	fmt.Println("validation failed", validationFailed)
	fmt.Printf("------------ Validation Ended (%v) ------------ \n\n", time.Since(startTime))
	fmt.Println("------------ Audit Started ------------")
	auditTasks, auditFailed := createAudit(done, validatedTasks)
	fmt.Println("Audit success", auditTasks)
	fmt.Println("Audit failed", auditFailed)
	fmt.Printf("------------ Audit Ended (%v) ------------ \n\n", time.Since(startTime))
	fmt.Println("------------ Update Started ------------")
	updatedTasks, updateFailed := updateTasks(done, auditTasks)
	fmt.Println("Update success", updatedTasks)
	fmt.Println("Update failed", updateFailed)
	fmt.Printf("------------ Update Ended (%v) ------------ \n\n", time.Since(startTime))
	fmt.Println("------------ Notify Started ------------")
	notified, notificationFailed := NotifyForTasks(done, updatedTasks)
	fmt.Println("Notify success", notified)
	fmt.Println("Notify failed", notificationFailed)
	fmt.Printf("------------ Notify Ended (%v) ------------ \n\n", time.Since(startTime))
	fmt.Println("All work done ", time.Since(startTime))
}
