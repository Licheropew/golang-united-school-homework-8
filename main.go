package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

type Arguments map[string]string

type User struct {
	Id    string `json:"id"`
	Email string `json:"email"`
	Age   int    `json:"age"`
}

var (
	operationFlag string
	fileNameFlag  string
	userId        string
	itemFlag      string
)

func init() {
	flag.StringVar(&operationFlag, "operation", "", "operation to perform")
	flag.StringVar(&itemFlag, "item", "", "User items")
	flag.StringVar(&fileNameFlag, "fileName", "", "name of the file")
	flag.Int(userId, 0, "id to search")
}

func parseArgs() Arguments {
	flag.Parse()
	return Arguments{
		"id":        userId,
		"operation": operationFlag,
		"fileName":  fileNameFlag,
		"item":      itemFlag,
	}
}

func checkItem(item string) error {
	if item == "" {
		return fmt.Errorf("-item flag has to be specified")
	}
	return nil
}

func getFilePath(fileName string) (string, error) {
	path, err := os.Getwd()
	if err != nil {
		return "", err
	}
	fPath := filepath.Join(path, fileName)
	return fPath, nil
}

func List(fileName string, writer io.Writer) error {
	// get file from current dir on Windows
	path, err := getFilePath(fileName)
	if err != nil {
		return err
	}
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	data := make([]byte, 2048)
	n, err := f.Read(data)
	if err != nil {
		if err != io.EOF {
			return err
		}
	}
	_, err = writer.Write(data[:n])
	if err != nil {
		return err
	}
	return nil
}

func Add(args Arguments, writer io.Writer) error {
	checkItemFlag := checkItem(args["item"])
	var u User
	var users []User
	if checkItemFlag != nil {
		return checkItemFlag
	}
	path, err := getFilePath(args["fileName"])
	if err != nil {
		return err
	}
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	data := make([]byte, 2048)
	n, err := f.Read(data)

	if err != nil {
		if err != io.EOF {
			return err
		}
	}

	json.Unmarshal([]byte(data[:n]), &users)
	json.Unmarshal([]byte(args["item"]), &u)

	for _, user := range users {
		if user.Id == u.Id {
			return fmt.Errorf("Item with id %v already exists", u.Id)
		}
	}
	if u.Id != "" {
		users = append(users, u)
		res, err := json.Marshal(users)
		if err != nil {
			return err
		}
		_, err = f.WriteAt(res, 0)
		if err != nil {
			return nil
		}
	}

	return nil
}

func Perform(args Arguments, writer io.Writer) error {
	if args["operation"] == "" {
		return fmt.Errorf("-operation flag has to be specified")
	}
	if args["fileName"] == "" {
		return fmt.Errorf("-fileName flag has to be specified")
	}
	switch args["operation"] {
	case "list":
		return List(args["fileName"], writer)
	case "add":
		return Add(args, writer)
	default:
		return fmt.Errorf("Operation %s not allowed!", args["operation"])
	}
}

func main() {
	err := Perform(parseArgs(), os.Stdout)
	if err != nil {
		panic(err)
	}
}
