package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
)

func main() {
	f, err := os.Open("problems.csv")
	defer f.Close()
	if err != nil {
		log.Fatal(err)
	}

	csvReader := csv.NewReader(bufio.NewReader(f))
	stdReader := bufio.NewReader(os.Stdin)
	for {
		record, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(record[0])
		stdReader.ReadLine()
	}
}
