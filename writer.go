package main

import (
	"context"
	"encoding/csv"
	"log"
	"os"
	"sync"
)

type CSVWriter struct {
	mutex     *sync.Mutex
	csvWriter *csv.Writer
}

func NewCSVWriter(fileName string) (*CSVWriter, error) {
	csvFile, err := os.Create(fileName)
	if err != nil {
		return nil, err
	}
	w := csv.NewWriter(csvFile)
	return &CSVWriter{
		csvWriter: w,
		mutex:     &sync.Mutex{},
	}, nil
}

func (w *CSVWriter) Write(row []string) {
	w.mutex.Lock()
	w.csvWriter.Write(row)
	w.mutex.Unlock()
}

func (w *CSVWriter) Flush() {
	w.mutex.Lock()
	w.csvWriter.Flush()
	w.mutex.Unlock()
}

func handleInterrupt(c chan os.Signal, cancel *context.CancelFunc) {
	<-c
	(*cancel)()
	writer.Flush()
	log.Fatalln("Scan cancelled by user.")
}
