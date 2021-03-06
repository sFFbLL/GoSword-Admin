package public

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"strings"

	"project/pkg/tools"
	"project/utils"
	"project/utils/app"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type FileResponse struct {
	Size     int64  `json:"size"`      //文件大小
	Path     string `json:"path"`      // 文件相对地址
	FullPath string `json:"full_path"` // 文件完整地址
	Name     string `json:"name"`      // 文件名
	Type     string `json:"type"`      // 文件类型
}

// UploadFile 文件上传（任意类型文件）
// @Summary 文件上传（任意类型文件）
// @Description Author：JiaKunLi 2021/01/27
// @Tags 文件：文件管理 File Controller
// @Accept multipart/form-data
// @Produce application/json
// @Param file formData file true "file"
// @Security ApiKeyAuth
// @Success 200 {object} models._ResponseFile
// @Router /api/file/uploadFile [post]
func UploadFile(c *gin.Context) {
	urlPerfix := fmt.Sprintf("http://%s/", c.Request.Host)
	var fileResponse FileResponse
	fileResponse, done := singleFile(c, fileResponse, urlPerfix, false)
	if done {
		return
	}
	app.ResponseSuccess(c, fileResponse)
}

// UploadFileImage 文件上传（图片）
// @Summary 文件上传（图片）
// @Description Author：JiaKunLi 2021/01/27
// @Tags 文件：文件管理 File Controller
// @Accept multipart/form-data
// @Produce application/json
// @Param file formData file true "file"
// @Security ApiKeyAuth
// @Success 200 {object} models._ResponseFile
// @Router /api/file/uploadImage [post]
func UploadFileImage(c *gin.Context) {
	urlPerfix := fmt.Sprintf("http://%s/", c.Request.Host)
	var fileResponse FileResponse
	fileResponse, done := singleFile(c, fileResponse, urlPerfix, true)
	if done {
		return
	}
	app.ResponseSuccess(c, fileResponse)
}

//func UploadFile(c *gin.Context) {
//	tag, _ := c.GetPostForm("type")
//	urlPerfix := fmt.Sprintf("http://%s/", c.Request.Host)
//	var fileResponse FileResponse
//	if tag == "" {
//		app.ResponseErrorWithMsg(c, 200, "缺少标识")
//		return
//	} else {
//		switch tag {
//		case "1": // 单图
//			fileResponse, done := singleFile(c, fileResponse, urlPerfix)
//			if done {
//				return
//			}
//			app.ResponseSuccess(c, fileResponse)
//			return
//		case "2": // 多图
//			multipartFile := multipleFile(c, urlPerfix)
//			app.ResponseSuccess(c, multipartFile)
//			return
//		case "3": // base64
//			fileResponse = baseImg(c, fileResponse, urlPerfix)
//			app.ResponseSuccess(c, fileResponse)
//		}
//	}
//}

func baseImg(c *gin.Context, fileResponse FileResponse, urlPerfix string) FileResponse {
	files, _ := c.GetPostForm("file")
	file2list := strings.Split(files, ",")
	ddd, _ := base64.StdEncoding.DecodeString(file2list[1])
	guid := uuid.New().String()
	fileName := guid + ".jpg"
	base64File := "static/uploadfile/" + fileName
	_ = ioutil.WriteFile(base64File, ddd, 0666)
	typeStr := strings.Replace(strings.Replace(file2list[0], "data:", "", -1), ";base64", "", -1)
	fileResponse = FileResponse{
		Size:     utils.GetFileSize(base64File),
		Path:     base64File,
		FullPath: urlPerfix + base64File,
		Name:     "",
		Type:     typeStr,
	}
	return fileResponse
}

func multipleFile(c *gin.Context, urlPerfix string) []FileResponse {
	files := c.Request.MultipartForm.File["file"]
	var multipartFile []FileResponse
	for _, f := range files {
		guid := uuid.New().String()
		fileName := guid + utils.GetExt(f.Filename)
		multipartFileName := "static/uploadfile/" + fileName
		e := c.SaveUploadedFile(f, multipartFileName)
		fileType, _ := tools.GetType(multipartFileName)
		if e == nil {
			fileResponse := FileResponse{
				Size:     utils.GetFileSize(multipartFileName),
				Path:     multipartFileName,
				FullPath: urlPerfix + multipartFileName,
				Name:     f.Filename,
				Type:     fileType,
			}
			multipartFile = append(multipartFile, fileResponse)
		}
	}
	return multipartFile
}

func singleFile(c *gin.Context, fileResponse FileResponse, urlPerfix string, image bool) (FileResponse, bool) {
	files, err := c.FormFile("file")

	if err != nil {
		app.ResponseError(c, app.CodeImageIsNotNull)
		return FileResponse{}, true
	}

	uploadPath := "static/uploadfile/"
	if image {
		if utils.GetFileType(tools.GetExt(files.Filename)) != "image" {
			app.ResponseError(c, app.CodeFileImageFail)
			return FileResponse{}, true
		}
		uploadPath = "static/image/"
	}

	// 上传文件至指定目录
	guid := uuid.New().String()
	fileName := guid + tools.GetExt(files.Filename)
	singleFile := uploadPath + fileName
	err = c.SaveUploadedFile(files, singleFile)
	if err != nil {
		app.ResponseError(c, app.CodeFileUploadFail)
		return FileResponse{}, true
	}
	fileType, _ := tools.GetType(singleFile)
	fileResponse = FileResponse{
		Size:     utils.GetFileSize(singleFile),
		Path:     singleFile,
		FullPath: urlPerfix + singleFile,
		Name:     files.Filename,
		Type:     fileType,
	}
	return fileResponse, false
}
