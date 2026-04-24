package main

import (
	"fmt"
	"os"
)


func saveContactRequest(name, msg string) error {
	return appendToFile("contacts.txt", fmt.Sprintf("Name: %s | Message: %s\n", name, msg)	)
}

func appendToFile(filename, content string) error {
	f, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.WriteString(content)
	return err
}