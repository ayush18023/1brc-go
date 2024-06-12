package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

type Values struct {
	min   float64
	max   float64
	sum   float64
	count int
}

func M2() {
	start := time.Now()
	file, err := os.Open("../1brc/measurements.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	buf := make([]byte, 0, 1024*1024)
	scanner.Buffer(buf, 10*1024*1024)
	// i := 0
	//min sum max count
	var mapper map[string]*Values = make(map[string]*Values)
	for scanner.Scan() {
		go fmt.Println(scanner.Text())
		// i += 1
		split := strings.Split(scanner.Text(), ";")
		pack := mapper[split[0]]
		f, err := strconv.ParseFloat(split[1], 32)
		if err != nil {
			panic(err)
		}
		if pack == nil {
			mapper[split[0]] = &Values{
				min:   f,
				max:   f,
				sum:   f,
				count: 1,
			}
		} else {
			pack.count += 1
			pack.sum += f
			if pack.max < f {
				pack.max = f
			}
			if pack.min > f {
				pack.min = f
			}
		}
	}
	for k, v := range mapper {
		fmt.Printf("%s=%f/%f/%f\n", k, v.min, (v.sum / float64(v.count)), v.max)
	}
	fmt.Println("it took total of ", time.Since(start))
}
