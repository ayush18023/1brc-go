package method2

//   64 *1024 	128*1024	256*1024
// 8 3m54s 		3m22s 		3m23s
// 6 3m54s 		3m21s 		3m49s

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

type entry struct {
	name string
	temp float64
}

type values struct {
	min   float64
	max   float64
	sum   float64
	count int
}

var Wg sync.WaitGroup = sync.WaitGroup{}

// func Process(toProcess []string, wg *sync.WaitGroup, entryCh chan []entry) {
// 	entriesPack := make([]entry, len(toProcess))
// 	for _, line := range toProcess {
// 		split := strings.Split(line, ";")
// 		f, _ := strconv.ParseFloat(split[1], 64)
// 		entriesPack = append(entriesPack, entry{
// 			name: split[0],
// 			temp: f,
// 			wg:   wg,
// 		})
// 	}
// 	entryCh <- entriesPack
// }

func M4() {
	// currentValue := runtime.GOMAXPROCS(0)

	// fmt.Printf("Current GOMAXPROCS value: %d\n", currentValue)
	// fmt.Printf("Current GOMAXPROCS value: %d\n", currentValue)

	// // Set GOMAXPROCS to utilize all available CPU cores
	// maxCores := runtime.NumCPU()
	// fmt.Printf("maxCores value: %d\n", maxCores)
	// newValue := maxCores - 2

	// runtime.GOMAXPROCS(newValue)
	// fmt.Printf("Updated GOMAXPROCS value: %d\n", runtime.GOMAXPROCS(0))
	Start := time.Now()
	file, err := os.Open("../1brc/measurements.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	var mapper map[string]*values = make(map[string]*values)
	mu := sync.Mutex{}

	entriesCh := make(chan []entry)
	go func() {
		for {
			select {
			case entries, ok := <-entriesCh:
				if ok {
					mu.Lock()
					for _, entry := range entries {
						// fmt.Println(entry.name)
						pack := mapper[entry.name]
						if pack == nil {
							mapper[entry.name] = &values{
								min:   entry.temp,
								max:   entry.temp,
								sum:   entry.temp,
								count: 1,
							}
						} else {
							pack.count += 1
							pack.sum += entry.temp
							if pack.max < entry.temp {
								pack.max = entry.temp
							}
							if pack.min > entry.temp {
								pack.min = entry.temp
							}
						}
						if entry.name != "" {
							Wg.Done()
						}
					}
					mu.Unlock()
				}
			}
		}
	}() //Combiner
	reader := bufio.NewReader(file)
	lines := make([]string, 0)
	linesChunkLen := 256 * 1024
	for {
		line, _, err := reader.ReadLine()
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Fatalf("a real error happened here: %v\n", err)
		}
		lines = append(lines, string(line))
		if len(lines) == linesChunkLen {
			toProcess := lines
			Wg.Add(len(lines))
			go func() {
				entriesPack := make([]entry, len(toProcess))
				for _, line := range toProcess {
					split := strings.Split(line, ";")
					f, _ := strconv.ParseFloat(split[1], 64)
					entriesPack = append(entriesPack, entry{
						name: split[0],
						temp: f,
					})
				}
				entriesCh <- entriesPack
			}()
			lines = make([]string, 0, linesChunkLen)
		}
	}
	Wg.Wait()
	close(entriesCh)
	for k, v := range mapper {
		fmt.Printf("%s=%f/%f/%f\n", k, v.min, (v.sum / float64(v.count)), v.max)
	}
	fmt.Println("it took total of ", time.Since(Start))
}
