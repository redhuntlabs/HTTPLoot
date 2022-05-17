package main

import (
	"bufio"
	"log"
	"os"
	"sync"
)

var (
	ProcChan  = make(chan *ProcJob, MAX_WORKERS)
	writer, _ = NewCSVWriter(OUTCSV)
)

type ProcJob struct {
	Host string
}

func ProcessHosts(args []string) {
	if len(INPFILE) > 0 {
		file, err := os.Open(INPFILE)
		if err != nil {
			log.Println(err)
			return
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		scanner.Split(bufio.ScanLines)
		for scanner.Scan() {
			ProcChan <- &ProcJob{
				Host: scanner.Text(),
			}
		}

		if err := scanner.Err(); err != nil {
			log.Fatalln("Error reading from target file:", err.Error())
		}
	} else {
		for _, target := range args {
			ProcChan <- &ProcJob{
				Host: target,
			}
		}
	}
	close(ProcChan)
}

func execWorker(wg *sync.WaitGroup) {
	for job := range ProcChan {
		job.RunChecks()
	}
	wg.Done()
}

func InitDispatcher(workerno int) {
	wg := new(sync.WaitGroup)
	for i := 0; i < workerno; i++ {
		wg.Add(1)
		go execWorker(wg)
	}
	wg.Wait()
}

func (p *ProcJob) RunChecks() {
	p.FingerPrint()
	p.ExecuteCrawler()
}
