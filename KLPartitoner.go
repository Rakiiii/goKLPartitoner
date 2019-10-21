package main

import (
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
	"time"

	graphlib "github.com/Rakiiii/goGraph"
	klpartitionlib "github.com/Rakiiii/goKLPartition"
)

func main() {
	itterationFlag := true
	typeFlag := ""
	var er error
	maxItteration := 0
	maxTime, _ := time.ParseDuration("0ms")

	if er != nil {
		log.Panicln(er)
		return
	}

	switch os.Args[1] {
	case "-i":
		typeFlag = os.Args[1]
		maxItteration, er = strconv.Atoi(os.Args[3])
		if er != nil {
			log.Panicln(er)
		}
	case "-t":
		typeFlag = os.Args[1]
		maxTime, er = time.ParseDuration(os.Args[3] + "ms")
		if er != nil {
			log.Println(er)
			return
		}
	default:
		log.Panicln("Wrong input flag")
	}
	graphPath := os.Args[2]

	var parser = new(graphlib.Parser)
	var g, err = parser.ParseUnweightedUndirectedGraphFromFile(graphPath)
	if err != nil {
		fmt.Println(err)
		return
	}
	//var for result of algorithm
	result := klpartitionlib.Result{Matrix: nil, Value: -1}

	//counte time
	var prevValue int64 = math.MaxInt64
	//getting time of start
	startTime := time.Now()
	//fmt.Println(startTime.Nanosecond())
	endTime := time.Now()
	//init var for saving tmie of first itteration
	firstItTime, _ := time.ParseDuration("0ms")

	itterationNumber := 0

	for itterationFlag {
		//partition of graph
		result, err = klpartitionlib.KLPartitionigAlgorithm(g, result.Matrix)

		//getting time after algoritm stops
		endTime = time.Now()
		//fmt.Println(endTime.Nanosecond())

		//init time of fitst itteration if needs
		if firstItTime.Nanoseconds() == 0 {
			firstItTime = endTime.Sub(startTime)
			//fmt.Println("first duration:", firstItTime.Nanoseconds())
		}
		if err != nil {
			log.Println(err)
			return
		}

		//checked for local max founding
		if prevValue <= result.Value {
			writeTime(firstItTime)
			saveResult(&result)
			return
		}

		if endTime.Sub(startTime).Milliseconds() > maxTime.Milliseconds() && typeFlag == "-t" {
			itterationFlag = false
		}
		if itterationNumber >= maxItteration && typeFlag == "-i" {
			itterationFlag = false
		}

		itterationNumber++
	}

	//save result and time
	writeTime(firstItTime)
	saveResult(&result)

}

func writeTime(time time.Duration) {
	timeFile, err := os.Create("time")
	defer timeFile.Close()
	if err != nil {
		fmt.Println(err)
		return
	} else {
		timeFile.WriteString(strconv.FormatInt(time.Milliseconds(), 10) + "ms")
	}
}

func saveResult(result *klpartitionlib.Result) {
	f, err := os.Create("result_" + os.Args[2])
	if err != nil {
		fmt.Println(err)
		fmt.Println(result.Value)
		for i := 0; i < result.Matrix.Heigh(); i++ {
			for j := 0; j < result.Matrix.Width(); j++ {
				fmt.Print(result.Matrix.GetBool(i, j))
			}
			fmt.Println()
		}
		return
	}
	defer f.Close()

	f.WriteString(strconv.FormatInt(result.Value, 10) + "\n")
	for i := 0; i < result.Matrix.Heigh(); i++ {
		subStr := ""
		for j := 0; j < result.Matrix.Width(); j++ {
			if result.Matrix.GetBool(i, j) {
				subStr = subStr + string("1 ")
			} else {
				subStr = subStr + string("0 ")
			}
		}
		subStr = subStr + "\n"
		f.WriteString(subStr)

	}

}
