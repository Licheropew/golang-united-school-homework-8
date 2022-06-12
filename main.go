package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

const (
	permissionF = 0644
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
	userIdFlag    string
	itemFlag      string
	tempUser      User
	users         []User
)

func init() {
	flag.StringVar(&operationFlag, "operation", "", "Must take operation to perform")
	flag.StringVar(&itemFlag, "item", "", "Must take user information")
	flag.StringVar(&fileNameFlag, "fileName", "", "Must take name of the file")
	flag.StringVar(&userIdFlag, "id", "", "Must take id to search")
}

func parseArgs() Arguments {
	flag.Parse()
	return Arguments{
		"id":        userIdFlag,
		"operation": operationFlag,
		"fileName":  fileNameFlag,
		"item":      itemFlag,
	}
}

func checkFlags(item, operation string) error {
	if item == "" {
		if operation == "add" {
			return fmt.Errorf("-item flag has to be specified")
		} else if operation == "remove" || operation == "findById" {
			return fmt.Errorf("-id flag has to be specified")
		}
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
	f, err := os.OpenFile(path, os.O_RDONLY|os.O_CREATE, permissionF)
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
	if n != 0 {
		_, err = writer.Write(data[:n])
		if err != nil {
			return err
		}
	}
	return nil
}

func Add(args Arguments, writer io.Writer) error {
	checkItemFlag := checkFlags(args["item"], args["operation"])
	if checkItemFlag != nil {
		return checkItemFlag
	}
	path, err := getFilePath(args["fileName"])
	if err != nil {
		return err
	}
	f, err := os.OpenFile(path, os.O_RDONLY|os.O_CREATE, permissionF)
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
	json.Unmarshal([]byte(args["item"]), &tempUser)

	addFlag := 0

	for _, user := range users {
		if user.Id == tempUser.Id {
			addFlag = 1
		}
	}
	if addFlag == 1 {
		fmt.Fprintf(writer, "Item with id %s already exists", tempUser.Id)
	} else {
		users = append(users, tempUser)
		res, err := json.Marshal(users)
		if err != nil {
			return err
		}
		outfile, _ := os.OpenFile(path, os.O_RDONLY|os.O_TRUNC, permissionF)
		defer outfile.Close()
		_, err = outfile.Write(res)
		if err != nil {
			return err
		}
	}

	return nil
}

func RemoveById(args Arguments, writer io.Writer) error {
	checkIdFlag := checkFlags(args["id"], args["operation"])
	if checkIdFlag != nil {
		return checkIdFlag
	}
	path, err := getFilePath(args["fileName"])
	if err != nil {
		return err
	}
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, permissionF)
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

	removeFlag := 0

	for i, user := range users {
		if user.Id == args["id"] {
			users = append(users[:i], users[i+1:]...)
			removeFlag = 1
		}
	}
	if removeFlag == 0 {
		fmt.Fprintf(writer, "Item with id %s not found", args["id"])
	} else {
		res, err := json.Marshal(users)
		if err != nil {
			return err
		}
		outfile, _ := os.OpenFile(path, os.O_WRONLY|os.O_TRUNC, permissionF)
		defer outfile.Close()
		_, err = outfile.Write(res)
		if err != nil {
			return err
		}
	}

	return nil
}

func FindById(args Arguments, writer io.Writer) error {
	checkIdFlag := checkFlags(args["id"], args["operation"])
	if checkIdFlag != nil {
		return checkIdFlag
	}
	path, err := getFilePath(args["fileName"])
	if err != nil {
		return err
	}
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, permissionF)
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
	findIdFlag := 0
	userNum := 0

	for i, user := range users {
		if user.Id == args["id"] {
			findIdFlag = 1
			userNum = i
		}
	}
	if findIdFlag == 1 {
		res, err := json.Marshal(users[userNum])
		if err != nil {
			return err
		}
		writer.Write(res)
	} else {
		writer.Write([]byte{})
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
	case "remove":
		return RemoveById(args, writer)
	case "findById":
		return FindById(args, writer)
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
