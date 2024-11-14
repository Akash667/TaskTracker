package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"time"
)

var fileDB *os.File
var TaskDB []Task
var taskCounter int

type Task struct {
	TaskID     int       `json:"taskID"`
	TaskDesc   string    `json:"taskDescription"`
	TaskStatus string    `json:"taskStatus"`
	CreatedAt  time.Time `json:"createdAt"`
}

func init() {
	var err error
	fileDB, err = os.OpenFile("tasks.json", os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		fmt.Println("Error: unable to open file", err)
		os.Exit(1)
	}
	taskCounter = 1
	fileStat, err := fileDB.Stat()
	if err != nil {
		fmt.Println("Error in reading file stats")
		os.Exit(1)
	}

	fileSize := fileStat.Size()

	if fileSize == 0 {
		fileDB.Write([]byte("[]"))
		return
	}

	fileBuffer := make([]byte, fileSize)

	fileDB.Read(fileBuffer)

	err = json.Unmarshal(fileBuffer, &TaskDB)

	for _, task := range TaskDB {
		if task.TaskID >= taskCounter {
			taskCounter = task.TaskID + 1
		}
	}
	if err != nil {
		fmt.Println("Error while unmarshalling", err)
	}

	// fmt.Printf("Task is %+v", TaskDB)

}

func main() {
	args := os.Args

	argsWithoutFile := args[1:]

	if len(argsWithoutFile) < 1 {
		fmt.Println("Error: nothing passed")
		os.Exit(1)
	}
	action := argsWithoutFile[0]

	switch action {
	case "add":
		addTask(argsWithoutFile[1:])
	case "delete":
		deleteTask(argsWithoutFile[1:])
	case ""
	}

}

func addTask(newTask []string) {
	if fileDB == nil {
		fmt.Println("File pointer is nil")
	} else {
		fmt.Println("File opened successfully")
	}

	taskDescription := newTask[0]
	timeNow := time.Now()

	createTask := Task{
		TaskID:     taskCounter,
		TaskDesc:   taskDescription,
		TaskStatus: "todo",
		CreatedAt:  timeNow,
	}
	TaskDB = append(TaskDB, createTask)

	truncateAndWrite(&TaskDB)

}

// func updateTask()     {}
func deleteTask(inputs []string) {

	taskIndex, err := strconv.Atoi(inputs[0])
	if err != nil {
		fmt.Println("Error: invalid type of index passed")
		os.Exit(1)
	}
	var taskToDeleteIndex int
	for i, r := range TaskDB {
		if r.TaskID == taskIndex {
			taskToDeleteIndex = i
			break
		}
	}

	TaskDB = append(TaskDB[:taskToDeleteIndex], TaskDB[taskToDeleteIndex+1:]...)
	truncateAndWrite(&TaskDB)
}

func truncateAndWrite(writeData *[]Task) {

	err := fileDB.Truncate(0)
	if err != nil {
		fmt.Println("Truncate error: ", err)
	}
	fileDB.Seek(0, 0)
	fileData, err := json.Marshal(writeData)
	if err != nil {
		fmt.Println("erorr marshalling data")
		os.Exit(1)
	}

	n, err := fileDB.WriteAt(fileData, 0)
	if err != nil {
		fmt.Println("Failed while writing to file, bytes written:", n, err)
		fmt.Println(err)
	}
}

func markDoneTask(inputs []string) {
	taskIndex, err := strconv.Atoi(inputs[0])
	if err != nil {
		fmt.Println("Error: invalid type of index passed")
		os.Exit(1)
	}

	for i, r := range TaskDB {
		if r.TaskID == taskIndex {
			r.TaskStatus = "done"
		}
	}

	truncateAndWrite(&TaskDB)
}

// func markInProgTask() {}
// func listTasks()      {}
