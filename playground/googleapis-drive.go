package main

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"golang-youtube-api/security/google/serviceaccount"
	"golang-youtube-api/utils"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
	"log"
	"os"
)

func main() {
	//router := gin.Default()
	//router.GET("/list", GoogleDriveTest)
	//router.Run(":5000")
	ctx := context.Background()
	client := serviceaccount.NewServiceAccount("credentials-web-go-blog-service-account.json")
	service, err := drive.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		fmt.Println()
	}
	fileList, err := service.Files.List().Q(fmt.Sprintf("name = '89c61407-fac6-11eb-8b25-9c5a443fe580' and parents = '%s'", utils.DriveImagesId)).Do()
	if err != nil {
		fmt.Println(err)
	}
	for i, file := range fileList.Files {
		fmt.Println(i, file.Id, file.Name)
	}
}

func GoogleDriveTest(c *gin.Context) {
	ctx := context.Background()
	// Step 2: Get the Google Drive service
	client := serviceaccount.NewServiceAccount("credentials-web-go-blog-service-account.json")
	service, err := drive.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		log.Fatalf("Unable to retrieve drive Client %v", err)
	}
	newPermissionsService := drive.NewPermissionsService(service)
	permission, err := newPermissionsService.Create("1GuVZDjFx1oM6WLTuSWIyQlttEH4uVSV8",
		&drive.Permission{Role: "writer", Type: "user", EmailAddress: "agamwork28@gmail.com", Deleted: true}).Do()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(permission.Id, permission.EmailAddress)
	multipartForm, err := c.MultipartForm()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(c.Request.Header.Get("Content-Type"))
	fileHeaders := multipartForm.File["images"]
	for _, fileCok := range fileHeaders {
		open, err := fileCok.Open()
		if err != nil {
			fmt.Println(err)
		}
		folderId := "1GuVZDjFx1oM6WLTuSWIyQlttEH4uVSV8"
		//folderId = "1H957ISIhjC2MZckzRuxi1P83qBggeD6y"
		//folderId = "1H957ISIhjC2MZckzRuxi1P83qBggeD6y"
		//CreateFolder(service, f)
		ListFolder(service, folderId)
		ListFile(service, folderId)
		fmt.Println(fileCok.Filename)
		createFile, err := serviceaccount.CreateFile(service, fileCok.Filename, "application/octet-stream", open, folderId)
		if err != nil {
			panic(fmt.Sprintf("Could not create file: %v\n", err))
		}

		fmt.Printf("File '%s' successfully uploaded", createFile.Name)
		fmt.Printf("\nFile Id: '%s' \n", createFile.Id)
	}

}

func DeleteFolder(service *drive.Service, parents string) {
}

func ListFolder(service *drive.Service, parents string) {
	fmt.Println("Get List Folder and Sub Folder")
	query := fmt.Sprintf("mimeType = 'application/vnd.google-apps.folder'")
	//query := fmt.Sprintf("parents = '%s' and mimeType = 'application/vnd.google-apps.folder'", parents)
	fileList, err := service.Files.List().Q(query).Do()
	if err != nil {
		log.Println(err)
	}
	for i, file := range fileList.Files {
		fmt.Println(i, file.Name, file.Id, file.Parents)
		if file.Name == "6de89c64-f92e-11eb-851b-9c5a443fe580" || file.Name == "92af8c00-f92f-11eb-aca4-9c5a443fe580" {
			service.Files.Delete(file.Id).Do()
		}
	}
	service.Files.Delete("13c37Sa88Jubu_HWJzOIS_SekPVIZsOZl").Do()
}

func ListFile(service *drive.Service, parents string) {
	fmt.Println("Get List File")
	query := fmt.Sprintf("mimeType != 'application/vnd.google-apps.folder'")
	//query := fmt.Sprintf("parents = '%s' and mimeType != 'application/vnd.google-apps.folder'", parents)
	fileList, err := service.Files.List().Q(query).Do()
	if err != nil {
		log.Println(err)
	}
	for i, file := range fileList.Files {
		//fmt.Println(file.Owners)
		fmt.Println(i, file.Name, file.Id, file.Parents)
		if file.OwnedByMe == false {
			service.Files.Delete(file.Id).Do()
		}
	}
}

func CreateFolder(service *drive.Service, file *os.File) {
	fmt.Println("Create Folder")
	// Step 3: Create directory
	dir, err := serviceaccount.CreateFolder(service, "Firman Agam", "1GuVZDjFx1oM6WLTuSWIyQlttEH4uVSV8")
	fmt.Println(dir.HTTPStatusCode, dir.Name, dir.DriveId, dir.Id)
	if err != nil {
		panic(fmt.Sprintf("Could not create dir: %v\n", err))
	}

	//give your folder id here in which you want to upload or create new directory
	folderId := dir.Id
	// Step 4: create the file and upload
	createFile, err := serviceaccount.CreateFile(service, file.Name(), "application/octet-stream", file, folderId)
	if err != nil {
		panic(fmt.Sprintf("Could not create file: %v\n", err))
	}

	fmt.Printf("File '%s' successfully uploaded", createFile.Name)
	fmt.Printf("\nFile Id: '%s' \n", createFile.Id)
}
