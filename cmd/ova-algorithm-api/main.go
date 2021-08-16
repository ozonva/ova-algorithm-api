package main

import (
	"fmt"
	"io"
	"os"
	"time"
)

func main() {
	fmt.Println("ova-algorithm-api")

	updateConfig := func(path string) error {
		file, err := os.Open(path)
		if err != nil {
			return fmt.Errorf("cannot open file \"%v\"", path)
		}
		defer file.Close()

		data := make([]byte, 256)

		for {
			count, err := file.Read(data)

			if count > 0 {
				fmt.Printf("%s", data[:count])
			}

			if err != nil {
				if err == io.EOF {
					break
				} else {
					return fmt.Errorf("error occured while config \"%v\": %w", path, err)
				}
			}
		}

		return nil
	}

	endTime := time.Now().Add(1 * time.Minute)

	for {
		err := updateConfig("configs/config.json")

		if err != nil {
			fmt.Printf("update Config failed: %v\n", err.Error())
		}

		if time.Now().After(endTime) {
			break
		}
		time.Sleep(1 * time.Second)
	}
}
