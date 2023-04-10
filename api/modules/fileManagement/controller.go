package fileManagement

import (
	"bcraft/api/errs"
	structures "bcraft/api/structures/responses"
	"github.com/gin-gonic/gin"
	"net/http"
)

type FileManagementController struct{}

// @Summary Upload File
// @Tags File Management
// @Description Upload File and get file path
// @ID upload
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "file"
// @Success 200 {object} structures.FileUploadResponse
// @Failure 401
// @Failure 500
// @Security Bearer
// @Router /file/upload [post]
func (f *FileManagementController) Upload(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		errs.NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	path, err := SaveFile(c, file)
	if err != nil {
		errs.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, structures.FileUploadResponse{FilePath: path})
}
