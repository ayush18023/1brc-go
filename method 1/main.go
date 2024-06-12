package method3

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

func M1() {
	Start := time.Now()
	file, err := os.Open("../1brc/measurements.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	var hashmap []*values = make([]*values, 2048)

	reader := bufio.NewReader(file)
	for {
		line, _, err := reader.ReadLine()
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Fatalf("a real error happened here: %v\n", err)
		}
		split := strings.Split(string(line), ";")
		f, _ := strconv.ParseFloat(split[1], 64)
		ind := hash(split[0])
		if hashmap[ind] == nil {
			hashmap[ind] = &values{
				name:  split[0],
				min:   f,
				max:   f,
				sum:   f,
				count: 1,
			}
		} else {
			hashmap[ind].count += 1
			hashmap[ind].sum += f
			if hashmap[ind].max < f {
				hashmap[ind].max = f
			}
			if hashmap[ind].min > f {
				hashmap[ind].min = f
			}
		}
	}
	for _, val := range hashmap {
		if val != nil {
			fmt.Printf("%s=%f/%f/%f\n", val.name, val.min, (val.sum / float64(val.count)), val.max)
		}
	}
	fmt.Println("it took total of ", time.Since(Start))
}
