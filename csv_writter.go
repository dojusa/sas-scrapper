package main

import (
	"encoding/csv"
	"fmt"
	"os"
)

type CSVRow interface {
	AsCSVRow() []string
}

func WriteCSV[T CSVRow](file string, head []string, rows []T) error {
	if len(rows) == 0 {
		fmt.Println("Skipping csv write as no rows where found")
		return nil
	}

	csvFile, err := os.Create(file)
	if err != nil {
		return err
	}
	defer csvFile.Close()

	writer := csv.NewWriter(csvFile)
	defer writer.Flush()

	writer.Write(head)
	for _, s := range rows {
		writer.Write(s.AsCSVRow())
	}
	return nil
}
