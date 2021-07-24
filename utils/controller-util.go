package utils

import (
	"bufio"
	"encoding/base64"
	"fmt"
	"github.com/google/uuid"
	"os"
	"strings"
)

func ToBase64(b []byte) string {
	return base64.StdEncoding.EncodeToString(b)
}

func GetImage(uuid uuid.UUID, folderName, filename string) string {
	var fileImage string
	filePath := fmt.Sprintf("/home/agam/IdeaProjects/golang-youtube-api/assets/images/%s/%s/%s", uuid.String(), folderName, filename)
	split := strings.Split(filePath, ".")
	imgFile, err := os.Open(filePath) // a QR code image
	if err != nil {
		fmt.Println(err)
		return ""
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
	imgBase64Str := ToBase64(buf)
	switch strings.ToLower(split[1]) {
	case "png":
		fileImage = fmt.Sprintf("data:image/png;base64,%s", imgBase64Str)
	case "jpeg", "jpg", "jpe":
		fileImage = fmt.Sprintf("data:image/png;base64,%s", imgBase64Str)
	default:
		fileImage = imgBase64Str
	}
	return fileImage
}
