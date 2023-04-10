package fileManagement

import (
	"bcraft/api/utils"
	"errors"
	"github.com/gin-gonic/gin"
	"mime/multipart"
	"strings"
)

const (
	fileNameLength = 16
	ContentDirPath = "api/content/"
)

func SaveFile(c *gin.Context, file *multipart.FileHeader) (string, error) {
	fileNameParts := strings.Split(file.Filename, ".")
	if len(fileNameParts) != 2 {
		return "", errors.New("invalid file name. Extension not found")
	}
	extension := fileNameParts[1]
	filename := utils.RandomString(fileNameLength) + "." + extension
	err := c.SaveUploadedFile(file, ContentDirPath+filename)
	if err != nil {
		return "", err
	}
	return "/content/" + filename, nil
}
