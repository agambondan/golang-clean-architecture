package utils

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"image"
	"image/jpeg"
	"image/png"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func ToBase64(b []byte) string {
	return base64.StdEncoding.EncodeToString(b)
}

func GetImageToBase64(uuid uuid.UUID, folderName, filename string) string {
	var fileImage string
	filePath := fmt.Sprintf("/home/agam/IdeaProjects/golang-youtube-api/assets/images/%s/%s/%s", uuid.String(), folderName, filename)
	split := strings.Split(filePath, ".")
	imgFile, err := os.Open(filePath) // a QR code image
	if err != nil {
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

func WriteImage(userId string, folderName, filename string) *bytes.Buffer {
	buffer := new(bytes.Buffer)
	var filePath string
	if userId != "" {
		filePath = fmt.Sprintf("/home/agam/IdeaProjects/golang-youtube-api/assets/images/%s/%s/%s", userId, folderName, filename)
	} else {
		filePath = fmt.Sprintf("/home/agam/IdeaProjects/golang-youtube-api/assets/images/%s/%s", folderName, filename)
	}
	splitFileName := strings.Split(filePath, ".")
	open, err := os.Open(filePath)
	if err != nil {
		log.Println(err)
		return buffer
	}
	defer open.Close()
	m, s, err := image.Decode(open)
	if err != nil {
		log.Println(s, err)
		return buffer
	}
	switch strings.ToLower(splitFileName[1]) {
	case "png":
		if err := png.Encode(buffer, m); err != nil {
			log.Println("unable to encode image.")
			return buffer
		}
	case "jpeg", "jpg", "jpe":
		if err := jpeg.Encode(buffer, m, nil); err != nil {
			log.Println("unable to encode image.")
			return buffer
		}
	default:
	}
	return buffer
}

func CreateUploadPhotoMachine(c *gin.Context, userId string, pathFolder string) ([]string, error) {
	// create folder and upload foto
	var err error
	var filenames []string
	header := c.Request.Header
	if header.Get("Content-Type")[:19] == "multipart/form-data" && userId != "" {
		formUser, err := c.MultipartForm()
		if err != nil {
			return filenames, err
		}
		files := formUser.File["images"]
		for i, file := range files {
			if i == 0 {
				if file.Size != 0 {
					basename := filepath.Base(file.Filename)
					regex := After(basename, ".")
					lowerRegex := strings.ToLower(regex)
					if lowerRegex[:2] == "pn" || lowerRegex[:2] == "jp" {
						dir := filepath.Join("./assets/images/", userId, pathFolder)
						if dir != "" {
							err = os.Mkdir("./assets/images/"+userId+pathFolder, os.ModePerm)
							if err != nil {
								_ = os.Mkdir("./assets/images/"+userId, os.ModePerm)
								_ = os.Mkdir("./assets/images/"+userId+pathFolder, os.ModePerm)
							}
						}
					}
					filename := filepath.Join("./assets/images/", userId, pathFolder, basename)
					err = c.SaveUploadedFile(file, filename)
					if err != nil {
						return filenames, err
					}
					filenames = append(filenames, file.Filename)
				} else {
					dir := filepath.Join("./assets/images/", userId, pathFolder)
					if dir != "" {
						err = os.Mkdir("./assets/images/"+userId+pathFolder, os.ModePerm)
						if err != nil {
							_ = os.Mkdir("./assets/images/"+userId, os.ModePerm)
							_ = os.Mkdir("./assets/images/"+userId+pathFolder, os.ModePerm)
						}
					}
				}
			}
		}
	} else {
		formUser, err := c.MultipartForm()
		if err != nil {
			return filenames, err
		}
		files := formUser.File["images"]
		for i, file := range files {
			if i == 0 {
				if file.Size != 0 {
					basename := filepath.Base(file.Filename)
					regex := After(basename, ".")
					lowerRegex := strings.ToLower(regex)
					if lowerRegex[:2] == "pn" || lowerRegex[:2] == "jp" {
						dir := filepath.Join("./assets/images/", pathFolder)
						if dir != "" {
							err = os.Mkdir("./assets/images/"+pathFolder, os.ModePerm)
							if err != nil {
								_ = os.Mkdir("./assets/images/"+pathFolder, os.ModePerm)
							}
						}
					}
					filename := filepath.Join("./assets/images/", pathFolder, basename)
					err = c.SaveUploadedFile(file, filename)
					if err != nil {
						return filenames, err
					}
					filenames = append(filenames, file.Filename)
				} else {
					dir := filepath.Join("./assets/images/", pathFolder)
					if dir != "" {
						_ = os.Mkdir("./assets/images/"+pathFolder, os.ModePerm)
					}
				}
			}
		}
	}
	return filenames, err
}

