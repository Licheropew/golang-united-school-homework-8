package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

type Arguments map[string]string

// type User struct {
// 	ID    string `json: "id"`
// 	Email string `json: "email"`
// 	Age   int    `json: "age"`
// }

var (
	operationFlag string
	fileNameFlag  string
	userData      string
	userId        string
	// parsedUser    User
)

func init() {
	flag.StringVar(&operationFlag, "operation", "", "operation to perform")
	flag.Func("item", "User items", func(s string) error {
		//err := json.Unmarshal([]byte(s), &parsedUser)
		userData = s
		return nil
	})
	flag.StringVar(&fileNameFlag, "fileName", "", "name of the file")
	flag.Int(userId, 0, "id to search")
}

func parseArgs() Arguments {
	flag.Parse()
	return Arguments{
		"operation": operationFlag,
		"filename":  fileNameFlag,
		"user":      userData,
	}
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
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()
	data := make([]byte, 2048)
	_, err = f.Read(data)
	if err != nil {
		return err
	}
	_, err = writer.Write(data)
	if err != nil {
		return err
	}
	return nil
}

func Add() error {
	return nil
}

func Perform(args Arguments, writer io.Writer) error {
	if args["operation"] == "" {
		return fmt.Errorf("-operation flag has to be specified")
	}
	if args["filename"] == "" {
		return fmt.Errorf("-fileName flag has to be specified")
	}
	switch {
	case args["operation"] == "list":
		return List(args["filename"], writer)
	case args["operation"] == "add":
		return Add()
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
