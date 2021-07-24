package main

import (
	"bufio"
	"encoding/base64"
	"fmt"
	"net/http"
	"os"
	"strings"
)

func Home(w http.ResponseWriter, r *http.Request) {
	imgFile, err := os.Open("/home/agam/IdeaProjects/golang-youtube-api/assets/images/17f4d5f6-e4cf-11eb-ab61-9c5a443fe580/user/Firda & Agam4.png") // a QR code image
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	defer imgFile.Close()

	// create a new buffer base on file size
	fInfo, _ := imgFile.Stat()
	var size int64 = fInfo.Size()
	buf := make([]byte, size)

	// read file content into buffer
	fReader := bufio.NewReader(imgFile)
	fReader.Read(buf)

	// if you create a new image instead of loading from file, encode the image to buffer instead with png.Encode()

	// png.Encode(&buf, image)

	// convert the buffer bytes to base64 string - use buf.Bytes() for new image
	imgBase64Str := base64.StdEncoding.EncodeToString(buf)
	fmt.Println(len(imgBase64Str))
	fmt.Println(byte(len(imgBase64Str)))

	// Embed into an html without PNG file
	img2html := "<html><body><img src=\"data:image/png;base64," + imgBase64Str + "\" /></body></html>"

	w.Write([]byte(fmt.Sprintf(img2html)))

}

//var vowelArray, consonantArray []string

func main() {
	// http.Handler
	//mux := http.NewServeMux()
	//mux.HandleFunc("/", Home)
	//log.Fatalln(http.ListenAndServe(":8080", mux))
	//var countVowel, countConsonant int
	var vowelArray, consonantArray []string
	var words = "keri ganteng banget dah ah"
	fmt.Println("words :", words)
	words = strings.Replace(words, " ", "", -1)
	fmt.Println("after replace :", words+",", "length words :", len(words))
	for i := 0; i <= len(words); i++ {
		if i != 0 {
			fmt.Print(words[i-1:i])
			vowelArray, consonantArray = isVowelConsonant(words[i-1:i], words, vowelArray, consonantArray)
		} else {
			vowelArray, consonantArray = isVowelConsonant(words[i:i+1], words, vowelArray, consonantArray)
		}
	}
	fmt.Println("")
	fmt.Println("len vowel :", len(vowelArray), ",", vowelArray)
	fmt.Println("len consonant :", len(consonantArray), ",", consonantArray)
	fmt.Println(strings.Join(vowelArray, ""), strings.Join(consonantArray, ""))
}

func isVowelConsonant(alphabet, words string, vowelArray, consonantArray []string) ([]string, []string) {
	if alphabet == "a" || alphabet == "i" || alphabet == "e" || alphabet == "u" || alphabet == "o" {
		count := strings.Count(words, alphabet)
		if !strings.Contains(strings.Join(vowelArray, ""), alphabet) {
			for i := 0; i < count; i++ {
				vowelArray = append(vowelArray, alphabet)
			}
		}
	} else {
		count := strings.Count(words, alphabet)
		if !strings.Contains(strings.Join(consonantArray, ""), alphabet) {
			for i := 0; i < count; i++ {
				consonantArray = append(consonantArray, alphabet)
			}
		}
	}
	return vowelArray, consonantArray
}

//vowelArray = append(vowelArray, strings.Join(vowel, ""))
//consonantArray = append(consonantArray, strings.Join(consonant, ""))
