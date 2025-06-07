package main

import (
	"os"
)

func main() {
	file, err := os.CreateTemp("", "test-")
	if err != nil {
		panic(err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			panic(err)
		}
	}(file)
}
