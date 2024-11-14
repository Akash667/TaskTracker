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
	case "mark-done":
		markDoneTask(argsWithoutFile[1:])
	case "update":
		updateTask(argsWithoutFile[1:])
	case "mark-in-progress":
		markInProgTask(argsWithoutFile[1:])
	case "list":
		listTasks(argsWithoutFile[1:])
	default:
		fmt.Println("Invalid command")
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

func updateTask(inputs []string) {
	taskIndex, err := strconv.Atoi(inputs[0])
	if err != nil || len(inputs) < 2 {
		fmt.Println("Error: invalid type of index passed or no description passed")
		os.Exit(1)
	}
	newTaskDescription := inputs[1]

	for index, r := range TaskDB {
		if r.TaskID == taskIndex {
			e := &TaskDB[index]
			e.TaskDesc = newTaskDescription
		}
	}

	truncateAndWrite(&TaskDB)

}
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

	for index, r := range TaskDB {
		if r.TaskID == taskIndex {
			e := &TaskDB[index]
			e.TaskStatus = "done"
		}
	}

	truncateAndWrite(&TaskDB)
}

func markInProgTask(inputs []string) {
	taskIndex, err := strconv.Atoi(inputs[0])
	if err != nil {
		fmt.Println("Error: invalid type of index passed")
		os.Exit(1)
	}

	for index, r := range TaskDB {
		if r.TaskID == taskIndex {
			e := &TaskDB[index]
			e.TaskStatus = "in-progress"
		}
	}

	truncateAndWrite(&TaskDB)
}

func listTasks(inputs []string) {
	fmt.Println("Task ID    Task Description    Task Status")
	if len(inputs) == 0 {
		for _, r := range TaskDB {
			fmt.Printf("%d \t%s \t%s\n", r.TaskID, r.TaskDesc, r.TaskStatus)
		}
		return
	}
	listType := inputs[0]

	for _, r := range TaskDB {
		if r.TaskStatus == listType {
			fmt.Printf("%d \t%s \t%s\n", r.TaskID, r.TaskDesc, r.TaskStatus)
		}
	}

}
