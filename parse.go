package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

func Parse() {
	if DBVerify() == false {
		return
	}

	DBConnect()
	DBMigrate()

	var pingResults []PingResults
	files, err := ioutil.ReadDir("output/")
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		if filepath.Ext(file.Name()) == ".json" {
			fullpath := "output/" + file.Name()
			file, err := ioutil.ReadFile(fullpath)
			if err != nil {
				fmt.Println(err)
				return
			}

			err = json.Unmarshal([]byte(file), &pingResults)
			if err != nil {
				fmt.Println(err)
				return
			}

			for _, pingResult := range pingResults {
				err = DBInsertPingResults(pingResult)
				if err != nil {
					fmt.Println(err)
					return
				}

				if _, err := os.Stat(fullpath); err == nil {
					// Remove the file once processed.
					os.Remove(fullpath)
				}
			}
		}
	}
}
