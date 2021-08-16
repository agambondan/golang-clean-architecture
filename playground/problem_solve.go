package main

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

func main() {
	//var arr = []int32{1, 16, 32, 523, 432, 421, 321, 13, 5, 7, 2}
	//var arr = []int32{256741038, 623958417, 467905213, 714532089, 938071625}
	//var highest, countHighest int
	////var arr = []int32{1, 2, 12, 10, 15, 5, 15, 13, 15}
	//var arr = []int32{44, 53, 31, 27, 77, 60, 66, 77, 26, 36}
	//lengthArr := len(arr)
	////var arrInt [lengthArr]int
	//var arrInt = make([]int, lengthArr)
	//for i := 0; i < len(arr); i++ {
	//	arrInt[i] = int(arr[i])
	//}
	//fmt.Println(arrInt)
	//sort.Ints(arrInt)
	//for i := 0; i < len(arrInt); i++ {
	//	if arrInt[i] == arrInt[lengthArr-1] {
	//		highest = arrInt[i]
	//		countHighest += 1
	//	}
	//}
	//fmt.Println(highest, countHighest)
	//timeString := "07:05:45PM"
	//fmt.Println(timeConversion(timeString))
	multipleFive := 67 % 5
	divFive := 97 / 5
	fmt.Println(multipleFive, divFive)
	var grades = []int32{73, 67, 38, 33}
	fmt.Println(gradingStudents(grades))
	apples := []int32{-2, 2, 1}
	oranges := []int32{5, -6}
	countApplesAndOranges(7, 11, 5, 15, apples, oranges)
	fmt.Println(kangaroo(21, 6, 47, 3))
}

func kangaroo(x1 int32, v1 int32, x2 int32, v2 int32) string {
	// Write your code here
	var output string
	if x1 == x2 && v1 == v2 {
		output = "YES"
	} else if x1 == x2 && v1 != v2 {
		output = "NO"
	} else {
		if v1 > v2 && ((x2-x1)%(v1-v2)) == 0 {
			output = "YES"
		} else {
			output = "NO"
		}
	}
	return output
}

func countApplesAndOranges(s int32, t int32, a int32, b int32, apples []int32, oranges []int32) {
	// Write your code here
	//s: integer, starting point of Sam's house location.
	//t: integer, ending location of Sam's house location.
	//a: integer, location of the Apple tree.
	//b: integer, location of the Orange tree.
	//apples: integer array, distances at which each apple falls from the tree.
	//oranges: integer array, distances at which each orange falls from the tree.
	var countApples []int
	var countOranges []int
	for i := 0; i < len(apples); i++ {
		temp := apples[i] + a
		if temp >= s && temp <= t {
			countApples = append(countApples, int(temp))
		}
	}
	for i := 0; i < len(oranges); i++ {
		temp := oranges[i] + b
		if temp >= s && temp <= t {
			countOranges = append(countOranges, int(temp))
		}
	}
	fmt.Println(len(countApples))
	fmt.Println(len(countOranges))
}

func gradingStudents(grades []int32) []int32 {
	// Write your code here
	var result = make([]int32, len(grades))
	fmt.Println(grades)
	for i := 0; i < len(grades); i++ {
		var mod, div int32
		mod = grades[i] % 5
		div = grades[i] / 5
		fmt.Println(mod, div, "LOOP ke", i)
		if mod < 3 || grades[i] < 38 {
			result[i] = grades[i]
		} else if mod >= 3 {
			if mod%2 == 0 {
				result[i] = grades[i] + 1
			} else {
				result[i] = grades[i] + 2
			}
		}
	}
	return result
}

func timeConversion(s string) string {
	// Write your code here
	var timeParseString string
	l := len(s)
	time := s[:l-2]
	timeFormat := s[l-2 : l]
	timeSplit := strings.Split(time, ":")
	hh, mm, ss := timeSplit[0], timeSplit[1], timeSplit[2]
	hhInt, _ := strconv.Atoi(hh)
	if hhInt < 12 && timeFormat == "PM" {
		hhInt = hhInt + 12
		fmt.Println("CUK")
	} else if hhInt == 12 && timeFormat == "AM" {
		hhInt = 0
		fmt.Println("CEK")
	}
	if hhInt < 10 {
		timeParseString = fmt.Sprintf("0%d:%s:%s", hhInt, mm, ss)
	} else {
		timeParseString = fmt.Sprintf("%d:%s:%s", hhInt, mm, ss)
	}
	return timeParseString
}

func miniMaxSum(arr []int32) {
	var minVal, maxVal int64
	for i := 0; i < len(arr); i++ {
		minf := sum(i, arr)
		fmt.Println("min f")
		maxf := sum(i, arr)
		fmt.Println("max f")
		if i == 0 {
			minVal = minf
			maxVal = maxf
			fmt.Println(minVal, maxVal, "MIN VAL MAX VAL")
		}
		if minf < minVal {
			minVal = minf
			fmt.Println(minVal, "MIN VAL")
		}
		if maxf > maxVal {
			maxVal = maxf
			fmt.Println(maxVal, "MIN MAX")
		}
	}
	fmt.Println(minVal, maxVal)
}

func sum(x int, arr []int32) int64 {
	var sum int64
	for i := 0; i < len(arr); i++ {
		if x != i {
			sum = int64(arr[i]) + sum
		}
	}
	return sum
}

func diagonalDifference() int32 {
	var row, col int
	fmt.Print("Enter number of rows: ")
	fmt.Scanln(&row)
	fmt.Print("Enter number of cols: ")
	fmt.Scanln(&col)
	var matrix [3][3]int32
	fmt.Println("========== Matrix1 =============")
	fmt.Println()
	for i := 0; i < row; i++ {
		for j := 0; j < col; j++ {
			fmt.Printf("Enter the element for Matrix1 %d %d :", i+1, j+1)
			fmt.Scanln(&matrix[i][j])
		}
	}
	var diag1, diag2, sum int32
	for i := 0; i < len(matrix); i++ {
		diag1 += matrix[i][i]
		diag2 += matrix[i][(len(matrix)-1)-i]
	}
	sumAbs := math.Abs(float64(diag1 - diag2))
	fmt.Printf("%f", sumAbs)
	sum = int32(sumAbs)
	return sum
}

func plusMinus(arr []int32) {
	// Write your code here
	var plus, minus, zero float64
	for i := 0; i < len(arr); i++ {
		if arr[i] == 0 {
			zero = zero + 1
		} else if arr[i] > 0 {
			plus = plus + 1
		} else {
			minus = minus + 1
		}
		fmt.Println(zero, plus, minus)
	}
	fmt.Println(arr)
	l := len(arr)
	fmt.Printf("%f \n", plus/float64(l))
	fmt.Printf("%f \n", minus/float64(l))
	fmt.Printf("%f \n", zero/float64(l))
}
