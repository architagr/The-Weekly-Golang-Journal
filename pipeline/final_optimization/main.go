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

func convertToStream(done <-chan bool, arr []PutTask) <-chan PutTask {
	stream := make(chan PutTask, len(arr))
	go func() {
		for _, task := range arr {
			select {
			case <-done:
				return
			case stream <- task:
			}
		}
		close(stream)
	}()
	return stream
}

type f func(task PutTask) bool

func genericStageTaskProcessor(done <-chan bool, wg *sync.WaitGroup, t PutTask, outputStream, failedDataStream chan<- PutTask, stageTaskProcesser f) {
	defer wg.Done()
	duration := time.Duration(100) * time.Millisecond
	if t.taskId%5 == 0 {
		duration = time.Duration(500) * time.Millisecond
	}
	// here we are simulating the processing time
	ticker := time.After(duration)
	for {
		select {
		case <-done:
			return
		case <-ticker:
			// here we are simulating the worker function,
			// that is passed as input to the startProcessing
			// function for each task
			if stageTaskProcesser(t) {
				outputStream <- t
				return
			}
			failedDataStream <- t
			return
		}
	}
}

func pipelineGenerator(done <-chan bool, readDataStream <-chan PutTask, failedDataStream chan<- PutTask, stageTaskProcesser f) <-chan PutTask {
	stream := make(chan PutTask)
	go func() {
		pushWg := &sync.WaitGroup{}
		defer close(stream)
		for {
			select {
			case task, ok := <-readDataStream:
				if !ok {
					pushWg.Wait()
					return
				}
				pushWg.Add(1)
				go genericStageTaskProcessor(done, pushWg, task, stream, failedDataStream, stageTaskProcesser)
			case <-done:
				return
			}
		}
	}()
	return stream
}

func validateTasks(task PutTask) bool {
	if task.taskId == 1 {
		fmt.Println("Validation failed for task id: ", task.taskId)
		return false
	}
	fmt.Println("Validation done for task id: ", task.taskId)
	return true
}

func createAudit(task PutTask) bool {
	if task.taskId == 2 {
		fmt.Println("Audit failed for task id: ", task.taskId)
		return false
	}
	fmt.Println("Audit done for task id: ", task.taskId)
	return true
}

func updateTasks(task PutTask) bool {
	if task.taskId == 3 {
		fmt.Println("Update failed for task id: ", task.taskId)
		return false
	}
	fmt.Println("Update done for task id: ", task.taskId)
	return true
}

func NotifyForTasks(task PutTask) bool {
	if task.taskId == 4 {
		fmt.Println("Notify failed for task id: ", task.taskId)
		return false
	}
	fmt.Println("Notify done for task id: ", task.taskId)
	return true
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
		{taskId: 7, summary: "The updated summary for task 7"},
		{taskId: 8, summary: "The updated summary for task 8"},
		{taskId: 9, summary: "The updated summary for task 9"},
		{taskId: 10, summary: "The updated summary for task 10"},
		{taskId: 11, summary: "The updated summary for task 11"},
		{taskId: 12, summary: "The updated summary for task 12"},
	}
	done := make(chan bool)
	defer close(done)
	failedTaskDataStream := make(chan PutTask)

	// Convert slice to channel
	stream := convertToStream(done, listOfTasks)

	// Stage 1: Validation
	validatedTaskStream := pipelineGenerator(done, stream, failedTaskDataStream, validateTasks)

	// Stage 2: Audit
	auditTaskStream := pipelineGenerator(done, validatedTaskStream, failedTaskDataStream, createAudit)

	// Stage 3: Database Update
	updatedTaskStream := pipelineGenerator(done, auditTaskStream, failedTaskDataStream, updateTasks)

	// Stage 4: Notify
	notifiedTaskStream := pipelineGenerator(done, updatedTaskStream, failedTaskDataStream, NotifyForTasks)

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		// Consume the final output from the notifiedTaskStream and failedTaskDataStream
		// to complete the pipeline and return the result to the caller
		completed := []PutTask{}
		failedData := []PutTask{}
		defer func() {
			fmt.Println("-------")
			fmt.Println("all success", completed)
			fmt.Println("all failed", failedData)
		}()
		for {
			select {
			case n, ok := <-notifiedTaskStream:
				if !ok {
					return
				}
				completed = append(completed, n)
			case n := <-failedTaskDataStream:
				failedData = append(failedData, n)
			}
		}
	}()
	wg.Wait()

	fmt.Println("---- All work done ", time.Since(startTime), "------")
}
