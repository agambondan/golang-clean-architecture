package main

import (
	"fmt"
	"sort"
)

func main() {
	var inputArray []int
	for i := 1; i <= 20; i++ {
		inputArray = append(inputArray, i)
	}
	arrayGanjil, arrayGenap := soalNomorSatu(inputArray)
	fmt.Println("Input Array :", inputArray)
	fmt.Println("Array Ganjil :", arrayGanjil)
	fmt.Println("Array Genap :", arrayGenap)
	var arrayA = []int{1, 2, 3, 5, 7, 11, 13, 15}
	var arrayB = []int{15, 12, 5, 7, 9, 11, 13, 10}
	fmt.Println("Array A :", arrayA)
	fmt.Println("Array B :", arrayB)
	irisan := soalNomorDua(arrayA, arrayB)
	fmt.Println("Output :", irisan)
}

func soalNomorSatu(inputArray []int) ([]int, []int) {
	var arrayGanjil, arrayGenap []int
	for i := 0; i < len(inputArray); i++ {
		if inputArray[i]%2 == 0 {
			arrayGenap = append(arrayGenap, inputArray[i])
		} else {
			arrayGanjil = append(arrayGanjil, inputArray[i])
		}
	}
	return arrayGanjil, arrayGenap
}

func soalNomorDua(arrayA, arrayB []int) []int {
	var irisan []int
	var checkArrayA, checkArrayB bool
	for i := 0; i < len(arrayA); i++ {
		for j := 0; j < len(arrayB); j++ {
			if arrayA[i] == arrayB[j] {
				checkArrayA = true
			}
			if arrayA[j] == arrayB[i] {
				checkArrayB = true
			}
		}
		if checkArrayA == false {
			irisan = append(irisan, arrayA[i])
		}
		if checkArrayB == false {
			irisan = append(irisan, arrayB[i])
		}
		checkArrayA, checkArrayB = false, false
	}
	sort.Ints(irisan)
	return irisan
}
