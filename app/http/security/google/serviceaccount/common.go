package serviceaccount

import (
	"context"
	"go-blog-api/app/lib"
	"golang.org/x/oauth2/google"
	"golang.org/x/oauth2/jwt"
	"google.golang.org/api/drive/v3"
	"io"
	"io/ioutil"
	"log"
	"net/http"
)

// NewServiceAccount Use Service account
func NewServiceAccount(secretFile string) *http.Client {
	b, err := ioutil.ReadFile(secretFile)
	if err != nil {
		log.Fatal("error while reading the credential file", err)
	}
	var s = struct {
		Email      string `json:"client_email"`
		PrivateKey string `json:"private_key"`
	}{}
	lib.Merge(b, &s)
	config := &jwt.Config{
		Email:      s.Email,
		PrivateKey: []byte(s.PrivateKey),
		Scopes: []string{
			drive.DriveScope,
		},
		TokenURL: google.JWTTokenURL,
	}
	client := config.Client(context.Background())
	return client
}

func CreateFolder(service *drive.Service, name string, parentId string) (*drive.File, error) {
	d := &drive.File{
		Name:     name,
		MimeType: "application/vnd.google-apps.folder",
		Parents:  []string{parentId},
	}

	file, err := service.Files.Create(d).Fields("id").Do()
	if err != nil {
		return nil, err
	}

	return file, nil
}

func CreateFile(service *drive.Service, name string, mimeType string, content io.Reader, parentId string) (*drive.File, error) {
	f := &drive.File{
		MimeType: mimeType,
		Name:     name,
		Parents:  []string{parentId},
	}
	//file, err := service.Files.Create(f).Do()
	file, err := service.Files.Create(f).Media(content).Do()
	if err != nil {
		return nil, err
	}

	return file, nil
}

func GetFiles(service *drive.Service, query string) (*drive.FileList, error) {
	fileList, err := service.Files.List().Q(query).Do()
	if err != nil {
		log.Println("could not get list folder", err.Error())
		return fileList, err
	}
	return fileList, err
}
