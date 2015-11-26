package main

import (
	"encoding/csv"
	"encoding/json"
	"log"
	"os"
)

func main() {
	var records [][]string
	if err := json.NewDecoder(os.Stdin).Decode(&records); err != nil {
		log.Fatalln(err)
	}
	w := csv.NewWriter(os.Stdout)
	for _, record := range records {
		if err := w.Write(record); err != nil {
			log.Fatalln(err)
		}
	}
	w.Flush()
}
