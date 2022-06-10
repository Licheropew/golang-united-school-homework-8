package main

import (
	"flag"
	"fmt"
	"io"
	"os"
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
	return Arguments{
		"operation": operationFlag,
		"fileName":  fileNameFlag,
		"user":      userData,
	}
}

func Add() error {
	return nil
}

func Perform(args Arguments, writer io.Writer) error {
	if args["operation"] == "" {
		return fmt.Errorf("-operation flag has to be specified")
	}
	if args["fileName"] == "" {
		return fmt.Errorf("-fileName flag has to be specified")
	}
	switch {
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
