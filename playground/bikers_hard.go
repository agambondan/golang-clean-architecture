package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

/*
 * Complete the 'bikeRacers' function below.
 *
 * The function is expected to return a LONG_INTEGER.
 * The function accepts following parameters:
 *  1. 2D_INTEGER_ARRAY bikers
 *  2. 2D_INTEGER_ARRAY bikes
 */

func bikeRacers(bikers [][]int32, bikes [][]int32) int64 {
	// Write your code here

	return 400000
}

func main() {
	reader := bufio.NewReaderSize(os.Stdin, 16*1024*1024)

	//stdout, err := os.Create(os.Getenv("OUTPUT_PATH"))
	stdout, err := os.Create("text.txt")
	checkError(err)

	defer stdout.Close()

	writer := bufio.NewWriterSize(stdout, 16*1024*1024)

	firstMultipleInput := strings.Split(strings.TrimSpace(readLine(reader)), " ")

	nTemp, err := strconv.ParseInt(firstMultipleInput[0], 10, 64)
	checkError(err)
	n := int32(nTemp)

	mTemp, err := strconv.ParseInt(firstMultipleInput[1], 10, 64)
	checkError(err)
	m := int32(mTemp)

	kTemp, err := strconv.ParseInt(firstMultipleInput[2], 10, 64)
	checkError(err)
	k := int32(kTemp)

	fmt.Println(m, k)
	var bikers [][]int32
	for i := 0; i < int(n); i++ {
		bikersRowTemp := strings.Split(strings.TrimRight(readLine(reader), " \t\r\n"), " ")

		var bikersRow []int32
		for _, bikersRowItem := range bikersRowTemp {
			bikersItemTemp, err := strconv.ParseInt(bikersRowItem, 10, 64)
			checkError(err)
			bikersItem := int32(bikersItemTemp)
			bikersRow = append(bikersRow, bikersItem)
		}

		if len(bikersRow) != 2 {
			panic("Bad input")
		}

		bikers = append(bikers, bikersRow)
	}

	var bikes [][]int32
	for i := 0; i < int(n); i++ {
		bikesRowTemp := strings.Split(strings.TrimRight(readLine(reader), " \t\r\n"), " ")

		var bikesRow []int32
		for _, bikesRowItem := range bikesRowTemp {
			bikesItemTemp, err := strconv.ParseInt(bikesRowItem, 10, 64)
			checkError(err)
			bikesItem := int32(bikesItemTemp)
			bikesRow = append(bikesRow, bikesItem)
		}

		if len(bikesRow) != 2 {
			panic("Bad input")
		}

		bikes = append(bikes, bikesRow)
	}

	result := bikeRacers(bikers, bikes)

	fmt.Fprintf(writer, "%d\n", result)

	writer.Flush()
}

func readLine(reader *bufio.Reader) string {
	str, _, err := reader.ReadLine()
	if err == io.EOF {
		return ""
	}

	return strings.TrimRight(string(str), "\r\n")
}

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}
