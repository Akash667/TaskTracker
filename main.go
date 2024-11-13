package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

var fileDB os.File
var TaskDB *[]Task

type Task struct {
	TaskID     int       `json:"taskID"`
	TaskDesc   string    `json:"taskDescription"`
	TaskStatus string    `json:"taskStatus"`
	CreatedAt  time.Time `json:"createdAt"`
}

func init() {

	fileDB, err := os.OpenFile("tasks.json", os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		fmt.Println("Error: unable to open file", err)
		os.Exit(1)
	}

	fileStat, err := fileDB.Stat()
	if err != nil {
		fmt.Println("Error in reading file stats")
		os.Exit(1)
	}

	fileSize := fileStat.Size()

	fileBuffer := make([]byte, fileSize)

	fileDB.Read(fileBuffer)

	err = json.Unmarshal(fileBuffer, TaskDB)

	if err != nil {
		fmt.Println("Error while unmarshalling", err)
	}

	fmt.Printf("%+v", TaskDB)

}

func main() {
	args := os.Args

	argsWithoutFile := args[1:]

	action := argsWithoutFile[0]
	fmt.Println(action)
	switch action {
	case "add":

	}

}
