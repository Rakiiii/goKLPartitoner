package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
	"strings"
	"time"

	boolmatrixlib "github.com/Rakiiii/goBoolMatrix"
	graphlib "github.com/Rakiiii/goGraph"
	klpartitionlib "github.com/Rakiiii/goKLPartition"
)

var res = "res"

func main() {
	//flag for first time saving
	ftimeFlag := true

	//condition flag to start next itteration
	itterationFlag := true

	//saving type of flag
	typeFlag := ""

	//er init
	var er error

	//flag for reload testing
	rFlag := false

	//path to graph
	graphPath := os.Args[2]

	//parse graph
	var parser = new(graphlib.Parser)
	var g, err = parser.ParseUnweightedUndirectedGraphFromFile(graphPath)
	if err != nil {
		fmt.Println(err)
		return
	}

	//itteration stop conditions
	maxItteration := 0
	maxTime, _ := time.ParseDuration("0ms")

	//var for result of algorithm
	result := klpartitionlib.Result{Matrix: nil, Value: -1}

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
	case "-r":
		typeFlag = "-i"
		maxItteration, er = strconv.Atoi(os.Args[3])
		if er != nil {
			log.Panicln(er)
		}
		r, err := parseResult(g.AmountOfVertex())
		if err != nil {
			if os.IsNotExist(err) {
				rFlag = true
			} else {
				log.Panicln(err)
				return
			}
		} else {
			result.Matrix = r.Matrix
			result.Value = r.Value
		}
	case "-mi":
		typeFlag = "-i"
		maxItteration = -1
		ftimeFlag = false
	default:
		log.Panicln("Wrong input flag")
	}

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

		fmt.Println("Value:", result.Value)

		//getting time after algoritm stops
		endTime = time.Now()
		//fmt.Println(endTime.Nanosecond())

		//init time of fitst itteration if needs
		if firstItTime.Nanoseconds() == 0 && ftimeFlag {
			firstItTime = endTime.Sub(startTime)
			//fmt.Println("first duration:", firstItTime.Nanoseconds())
		}
		if err != nil {
			log.Println(err)
			return
		}

		//checked for local max founding
		if prevValue <= result.Value {
			firstItTime = endTime.Sub(startTime)
			writeTime(firstItTime)
			saveResult(&result, "result_"+os.Args[2])
			return
		}

		if endTime.Sub(startTime).Milliseconds() > maxTime.Milliseconds() && typeFlag == "-t" {
			itterationFlag = false
		}
		if itterationNumber >= maxItteration && typeFlag == "-i" && maxItteration != -1 {
			itterationFlag = false
		}

		if rFlag {
			rFlag = false
			saveResult(&result, res)
		}

		itterationNumber++
	}

	//save result and time
	writeTime(firstItTime)
	saveResult(&result, "result_"+os.Args[2])

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

func saveResult(result *klpartitionlib.Result, path string) {
	f, err := os.Create(path)
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

func parseResult(length int) (klpartitionlib.Result, error) {
	file, err := os.Open(res)
	if os.IsNotExist(err) {
		return klpartitionlib.Result{Value: -1, Matrix: nil}, err
	}
	defer file.Close()
	res := klpartitionlib.Result{Matrix: nil, Value: -1}
	var bm boolmatrixlib.BoolMatrix
	bm.Init(2, length)
	scanner := bufio.NewScanner(file)

	scanner.Scan()
	subVal, err := strconv.Atoi(scanner.Text())
	if err != nil {
		return klpartitionlib.Result{Value: -1, Matrix: nil}, err
	} else {
		res.Value = int64(subVal)
	}

	it := 0
	for scanner.Scan() {
		parts := strings.Fields(scanner.Text())
		fr, err := strconv.Atoi(parts[0])
		if err != nil {
			return klpartitionlib.Result{Value: -1, Matrix: nil}, err
		}
		if fr == 0 {
			bm.SetBool(it, 0, false)
		} else {
			bm.SetBool(it, 0, true)
		}
		sr, err := strconv.Atoi(parts[1])
		if err != nil {
			return klpartitionlib.Result{Value: -1, Matrix: nil}, err
		}
		if sr == 0 {
			bm.SetBool(it, 1, false)
		} else {
			bm.SetBool(it, 1, true)
		}
		it++
	}
	res.Matrix = &bm
	return res, nil
}
