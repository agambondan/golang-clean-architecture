package google

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"golang-youtube-api/security/google/serviceaccount"
	"golang-youtube-api/utils"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
)

func UploadImageFileToAssets(c *gin.Context, folder, userId, parentId string) (*drive.File, error) {
	var driveFile *drive.File
	folderId := parentId
	multipartForm, err := c.MultipartForm()
	if err != nil {
		return driveFile, err
	}
	client := serviceaccount.NewServiceAccount("credentials-web-go-blog-service-account.json")
	service, err := drive.NewService(c, option.WithHTTPClient(client))
	if err != nil {
		return driveFile, err
	}
	if userId != "" {
		q := fmt.Sprintf("name = '%s' and parents = '%s'", userId, parentId)
		fileList, err := service.Files.List().Q(q).Do()
		if err != nil {
			return driveFile, err
		}
		if len(fileList.Files) == 0 {
			driveFile, err = serviceaccount.CreateFolder(service, userId, parentId)
			if err != nil {
				return driveFile, err
			}
			_, err = GivePermission(service, driveFile.Id, utils.DrivePermissionEmailAddress, utils.DrivePermissionRole, utils.DrivePermissionType)
			if err != nil {
				return driveFile, err
			}
			folderId = driveFile.Id
		} else {
			folderId = fileList.Files[0].Id
			q = fmt.Sprintf("name = 'post' and parents = '%s'", folderId)
			fileList, err = service.Files.List().Q(q).Do()
			if err != nil {
				return driveFile, err
			}
		}
		if len(fileList.Files) == 0 {
			driveFile, err = serviceaccount.CreateFolder(service, folder, folderId)
			if err != nil {
				return driveFile, err
			}
			_, err = GivePermission(service, driveFile.Id, utils.DrivePermissionEmailAddress, utils.DrivePermissionRole, utils.DrivePermissionType)
			if err != nil {
				return driveFile, err
			}
			folderId = driveFile.Id
		} else {
			folderId = fileList.Files[0].Id
		}
	}
	fileHeaders := multipartForm.File["images"]
	for i, fileCok := range fileHeaders {
		if i == 0 {
			open, err := fileCok.Open()
			if err != nil {
				return driveFile, err
			}
			driveFile, err = serviceaccount.CreateFile(service, fileCok.Filename, "application/octet-stream", open, folderId)
			if err != nil {
				return driveFile, err
			}
			driveFile.WebViewLink = fmt.Sprintf("https://drive.google.com/uc?id=%s", driveFile.Id)
			fileWebLink, err := service.Files.Get(driveFile.Id).Fields("thumbnailLink").Do()
			driveFile.ThumbnailLink = fileWebLink.ThumbnailLink
			if err != nil {
				return driveFile, err
			}
			_, err = GivePermission(service, driveFile.Id, utils.DrivePermissionEmailAddress, utils.DrivePermissionRole, utils.DrivePermissionType)
			if err != nil {
				return driveFile, err
			}
		}
	}

	return driveFile, err
}

func GivePermission(service *drive.Service, fileId, emailAddress, role, typePermission string) (*drive.Permission, error) {
	newPermissionsService := drive.NewPermissionsService(service)
	permission, err := newPermissionsService.Create(fileId,
		&drive.Permission{Role: role, Type: typePermission, EmailAddress: emailAddress, Deleted: true}).Do()
	if err != nil {
		return permission, err
	}
	return permission, err
}

//fileWebLink, err := service.Files.Get(driveFile.Id).Fields("webViewLink", "webContentLink", "thumbnailLink", "iconLink").Do()
//driveFile.IconLink = fileWebLink.IconLink
//driveFile.WebContentLink = fileWebLink.WebContentLink
//driveFile.WebViewLink = fileWebLink.WebViewLink
//driveFile.ThumbnailLink = fileWebLink.ThumbnailLink
