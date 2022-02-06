package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

func readData(path string) {
	// read file
	data, err := ioutil.ReadFile(path + ".json")
	if err != nil {
		fmt.Print(err)
	}

	// unmarshall it
	err = json.Unmarshal(data, &_data)
	if err != nil {
		fmt.Println("error:", err)
	}
}

func main() {
	readData("keys")
	StartServer()
}
