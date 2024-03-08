package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"runtime/pprof"
	"strconv"
	"strings"
	"sync"
	"time"
)

type info struct {
	cumulativeTemp   float64
	numberOfReadings int
	min              float64
	max              float64
}

var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")
var routines = 5

func main() {
	flag.Parse()
	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}
	start := time.Now()
	// open file
	file, err := os.Open("measurements.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	// create hashmap of result data
	// result := make(sync.Map[string]*info)
	result := sync.Map{}
	// read line

	byteSize := 1048 * 15
	b1 := make([]byte, byteSize)
	channelList := make(map[int]chan string)
	for i := range routines {
		channelList[i] = make(chan string, 10)
		go func() {
			handleData(channelList[i], &result)
		}()
	}

	// f1 := make(chan string)

	counter := 0
	b2 := make([]byte, 64)
	for {
		// read a chunk
		n, err := file.Read(b1)
		if err != nil && err != io.EOF {
			panic(err)
		}
		if n == 0 {
			break
		}
		position := 0
		var output string
		for {
			if _, err := file.Read(b2[position : position+1]); err != nil {
				break
			}
			if b2[position] == byte('\n') {
				output = output + string(b1)
				output = output + string(b2[0:position])
				position = 0
				break
			} else {
				position++
			}
		}
		channelList[counter] <- string(output)
		counter++
		counter = counter % routines
	}
	// for scanner.Scan() {
	// 	// parse line
	// 	line := scanner.Text()
	// 	// put in object
	// 	channelList[counter] <- line
	// 	// err = handleData(line, result)
	// 	if err != nil {
	// 		fmt.Println(err.Error())
	// 	}
	// 	counter++
	// 	counter = counter % 5
	// }
	// for key, value := range result {
	// 	fmt.Printf("City: %v, avg temp: %v, min: %v, max: %v\n", key, value.cumulativeTemp/float64(value.numberOfReadings), value.min, value.max)
	// }
	end := time.Now()
	fmt.Printf("Total time taken: %v", end.Sub(start))
}

func handleData(lineChan chan string, result *sync.Map) error {
	for {
		select {
		case <-lineChan:
			data := <-lineChan
			lines := strings.Split(data, "\n")
			for _, line := range lines {
				lineData := strings.Split(line, ";")
				if len(lineData) <= 1 {
					fmt.Printf("here: %v\n", lineData)
					continue
				}
				temp, err := strconv.ParseFloat(lineData[1], 32)
				if err != nil {
					fmt.Printf("error: %v", lineData[1])
					return err
				}
				_, exists := result.Load(lineData[0])
				if !exists {
					result.Store(lineData[0], &info{})
				}
				obj, _ := result.Load(lineData[0])
				i := obj.(*info)
				i.cumulativeTemp += temp
				i.numberOfReadings++
				if temp > i.max {
					i.max = temp
				}
				if temp < i.min {
					i.min = temp
				}
			}

		default:
			continue
		}
	}
}
