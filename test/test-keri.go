package main

import (
	"fmt"
	"log"
)

//var (
//	vowelArray     []string
//	consonantArray []string
//)

//func main() {
//	var words = "Firda Dwi Gameswanti"
//	words = strings.ToLower(words)
//	var newWords = strings.Replace(words, " ", "", -1)
//	for i := 0; i <= len(newWords); i++ {
//		if i != 0 {
//			IsVowelConsonant(newWords[i-1:i], newWords)
//		} else {
//			IsVowelConsonant(newWords[i:i], newWords)
//		}
//	}
//	fmt.Println("words :", words)
//	fmt.Println("after replace :", newWords)
//	fmt.Println(len(vowelArray), len(consonantArray))
//	fmt.Println("output vowel :", strings.Join(vowelArray, ""))
//	fmt.Println("output consonant :", strings.Join(consonantArray, ""))
//}
//
//func IsVowelConsonant(alphabet, words string) {
//	if alphabet == "a" || alphabet == "i" || alphabet == "e" || alphabet == "u" || alphabet == "o" {
//		count := strings.Count(words, alphabet)
//		if !strings.Contains(strings.Join(vowelArray, ""), alphabet) {
//			for i := 0; i < count; i++ {
//				vowelArray = append(vowelArray, alphabet)
//			}
//		}
//	} else {
//		count := strings.Count(words, alphabet)
//		if !strings.Contains(strings.Join(consonantArray, ""), alphabet) {
//			for i := 0; i < count; i++ {
//				consonantArray = append(consonantArray, alphabet)
//			}
//		}
//	}
//}

func main() {
	var families int
	//var familyBus []string
	fmt.Print("Input the number of families : ")
	_, err := fmt.Scanln(&families)
	if err != nil {
		fmt.Println(err)
	}
	members := make([]int, families)
	inputFamilies := make([]interface{}, families)
	fmt.Print("Input the number of members in the family (separated by space) : ")
	for i := 0; i < families; i++ {
		inputFamilies[i] = &members[i]
	}
	_, err = fmt.Scanln(inputFamilies...)
	if err != nil {
		log.Println("Input must be equal with count family")
		return
	}
	var sliceBus = make(map[int]map[int]bool)
	for i := 0; i < len(members); i++ {
		tampungBus := map[int]bool{}
		tampungBus[members[i]] = false
		sliceBus[i] = tampungBus
	}
	fmt.Println(sliceBus)
	for i := 0; i < len(sliceBus); i++ {
		m := sliceBus[i]
		for k, v := range m {
			fmt.Println(k,v)
		}
	}

	//var bus = make(map[int]int)
	//var sliceBus = make(map[int]map[int]bool)
	//for i := 0; i < len(members); i++ {
	//	if members[i] < 4 {
	//		countMinimumBus += 1
	//	}
	//	var tampungBus = make(map[int]bool)
	//	tampungBus[members[i]] = false
	//	sliceBus[i] = tampungBus
	//	bus[i] = members[i]
	//	totalPeople = totalPeople + members[i]
	//}
	//fmt.Println(sliceBus)
	//for i := 0; i < len(bus); i++ {
	//	countPeople = countPeople + bus[i]
	//}
	//countMinimumBus = totalPeople / 4
	//if totalPeople%4 != 0 {
	//	countMinimumBus += 1
	//}
	//fmt.Println(countMinimumBus)
	//fmt.Println(countPeople)
	//fmt.Println("Total People :", totalPeople)
	//countMinimumBus = 0
	////the rules is :
	////1. 1 bus cuma bisa 4 penumpang
	////2. 1 bus cuma bisa di isi 2 keluarga
	//for i := 0; i < len(bus); i++ {
	//	for j := 0; j < len(bus); j++ {
	//		if (bus[i]+bus[j]%4 == 0 || bus[i]+bus[j] > 4) && !strings.Contains(strings.Join(familyBus, ""), strconv.Itoa(bus[i])) {
	//			countMinimumBus += 1
	//			familyBus = append(familyBus, strconv.Itoa(bus[i]))
	//			fmt.Println(familyBus, "CUKS")
	//		}
	//	}
	//}
	//
	//fmt.Println(bus, "map bus")
	//fmt.Println("Minimum", countMinimumBus, "Bus")
}

//x, err := ScanWithSlice(families)
//if err != nil {
//	fmt.Println(err, "Input must be equal with count of family")
//	return
//}
//fmt.Printf("%v\n", x)

// ScanWithSlice source https://stackoverflow.com/questions/15413469/how-to-make-fmt-scanln-read-into-a-slice-of-integers
func ScanWithSlice(n int) ([]int, error) {
	x := make([]int, n)
	y := make([]interface{}, n)
	for i := range x {
		y[i] = &x[i]
	}
	n, err := fmt.Scanln(y...)
	x = x[:n]
	return x, err
}
