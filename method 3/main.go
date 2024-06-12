package method4

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
	name  string
	min   float64
	max   float64
	sum   float64
	count int
}

var Wg sync.WaitGroup = sync.WaitGroup{}

//	func Process(toProcess []string, wg *sync.WaitGroup, entryCh chan []entry) {
//		entriesPack := make([]entry, len(toProcess))
//		for _, line := range toProcess {
//			split := strings.Split(line, ";")
//			f, _ := strconv.ParseFloat(split[1], 64)
//			entriesPack = append(entriesPack, entry{
//				name: split[0],
//				temp: f,
//				wg:   wg,
//			})
//		}
//		entryCh <- entriesPack
//	}
const n int = 3

func hash(val string) int {
	var hashval int = 0
	for i := 0; i < n; i++ {
		hashval += int(val[i]) * (n - i)
	}
	return hashval
}

func M3() {
	// currentValue := runtime.GOMAXPROCS(0)

	// fmt.Printf("Current GOMAXPROCS value: %d\n", currentValue)
	// fmt.Printf("Current GOMAXPROCS value: %d\n", currentValue)

	// // Set GOMAXPROCS to utilize all available CPU cores
	// maxCores := runtime.NumCPU()
	// fmt.Printf("maxCores value: %d\n", maxCores)
	// newValue := maxCores - 2

	// runtime.GOMAXPROCS(newValue)
	// fmt.Printf("Updated GOMAXPROCS value: %d\n", runtime.GOMAXPROCS(0))
	fmt.Println(hash("Mumbai"))
	Start := time.Now()
	file, err := os.Open("../1brc/measurements.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	// var mapper map[string]*values = make(map[string]*values)
	var hashmap []*values = make([]*values, 2048)
	// var boolmap []bool = make([]bool, 0, 2048)
	mu := sync.Mutex{}

	entriesCh := make(chan []entry)
	go func() {
		for {
			select {
			case entries, ok := <-entriesCh:
				if ok {
					mu.Lock()
					for _, entry := range entries {
						if entry.name != "" {
							ind := hash(entry.name)
							if hashmap[ind] == nil {
								hashmap[ind] = &values{
									name:  entry.name,
									min:   entry.temp,
									max:   entry.temp,
									sum:   entry.temp,
									count: 1,
								}
							} else {
								hashmap[ind].count += 1
								hashmap[ind].sum += entry.temp
								if hashmap[ind].max < entry.temp {
									hashmap[ind].max = entry.temp
								}
								if hashmap[ind].min > entry.temp {
									hashmap[ind].min = entry.temp
								}
							}
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
	for _, val := range hashmap {
		if val != nil {
			fmt.Printf("%s=%f/%f/%f\n", val.name, val.min, (val.sum / float64(val.count)), val.max)
		}
	}
	fmt.Println("it took total of ", time.Since(Start))
}
