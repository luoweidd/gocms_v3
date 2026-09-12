package controller

import (
	"crypto/md5"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"gocms_v3/app/db"
	"gocms_v3/app/model"
	"gocms_v3/app/response"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// UploadController 文件上传控制器
type UploadController struct{}

func NewUploadController() *UploadController {
	return &UploadController{}
}

// InitiateChunkUpload 初始化分片上传 API: POST /upload/initiate
func (ctrl *UploadController) InitiateChunkUpload(c *gin.Context) {
	// 支持三种请求格式：
	// 1. multipart/form-data with file 字段
	// 2. multipart/form-data with filename/size 字段
	// 3. application/json with filename/size 字段

	var filename string
	var size int64

	// 先尝试从 Content-Type 判断请求类型
	ct := c.GetHeader("Content-Type")

	// 尝试从 form data 获取文件（场景1）
	file, headerErr := c.FormFile("file")
	if headerErr == nil && file != nil {
		filename = file.Filename
		size = file.Size
	} else {
		// 尝试解析表单数据（场景2）
		if c.Request.MultipartReader == nil {
			if err := c.Request.ParseMultipartForm(32 << 20); err != nil && err.Error() != "request body is available" {
				response.Error(c, 400, "解析表单失败: "+err.Error())
				return
			}
		}
		formFilename := c.Request.Form.Get("filename")
		formSizeStr := c.Request.Form.Get("size")

		if formFilename != "" && formSizeStr != "" {
			// 从 form data 获取 filename 和 size
			filename = formFilename
			fmt.Sscanf(formSizeStr, "%d", &size)
		} else if ct == "application/json" {
			// 场景3：从 JSON 获取（兼容直接传 filename/size 的场景）
			var jsonReq struct {
				Filename string `json:"filename"`
				Size     int64  `json:"size"`
			}
			if err := c.ShouldBindJSON(&jsonReq); err != nil {
				response.Error(c, 400, "请求参数错误: "+err.Error())
				return
			}
			if jsonReq.Filename == "" || jsonReq.Size == 0 {
				response.Error(c, 400, "filename 和 size 不能为空")
				return
			}
			filename = jsonReq.Filename
			size = jsonReq.Size
		} else {
			response.Error(c, 400, "缺少必要参数: filename 和 size")
			return
		}
	}

	if filename == "" || size == 0 {
		response.Error(c, 400, "文件名和文件大小不能为空")
		return
	}

	// 获取 chunk_num
	chunkNum := 0
	if chunkNumStr := c.Request.Form.Get("chunk_num"); chunkNumStr != "" {
		fmt.Sscanf(chunkNumStr, "%d", &chunkNum)
	} else if ct == "application/json" {
		var jsonChunk struct {
			ChunkNum int `json:"chunk_num"`
		}
		if err := c.ShouldBindJSON(&jsonChunk); err == nil {
			chunkNum = jsonChunk.ChunkNum
		}
	}

	// 生成唯一的uploadID
	hash := md5.Sum([]byte(filename + fmt.Sprintf("%d", size)))
	uploadID := fmt.Sprintf("%x_%d_%d", hash, chunkNum, time.Now().UnixNano())

	record := model.FileUpload{
		UploadID: uploadID,
		Filename: filename,
		Size:     size,
		ChunkNum: 0,
		Status:   0, // 初始化中
	}
	db.GetDB().Create(&record)

	response.Success(c, gin.H{
		"upload_id": uploadID,
		"chunk_num": chunkNum,
	})
}

// UploadChunk 上传分片 API: POST /upload/chunk
func (ctrl *UploadController) UploadChunk(c *gin.Context) {
	// 先解析 multipart 表单
	if err := c.Request.ParseMultipartForm(32 << 20); err != nil {
		response.Error(c, 400, "解析表单失败: "+err.Error())
		return
	}

	// 获取上传记录
	var uploadRecord model.FileUpload
	uploadID := c.Request.Form.Get("upload_id")
	if uploadID != "" {
		if err := db.GetDB().Where("upload_id = ?", uploadID).First(&uploadRecord).Error; err == nil {
			_ = uploadRecord.Filename
		}
	}

	file, err := c.FormFile("chunk")
	if err != nil {
		response.Error(c, 400, "获取文件失败: "+err.Error())
		return
	}

	// 从 Form 中获取参数
	uploadID = c.Request.Form.Get("upload_id")
	chunkNumStr := c.Request.Form.Get("chunk_num")

	if uploadID == "" {
		response.Error(c, 400, "缺少必要参数: upload_id")
		return
	}

	chunkNum := 0
	if chunkNumStr != "" {
		fmt.Sscanf(chunkNumStr, "%d", &chunkNum)
	}

	dir := "./storage/uploads/chunks/" + uploadID
	os.MkdirAll(dir, 0755)
	chunkPath := filepath.Join(dir, fmt.Sprintf("%d", chunkNum))
	src, err := file.Open()
	if err != nil {
		response.Error(c, 500, "打开文件失败")
		return
	}
	defer src.Close()
	dst, err := os.Create(chunkPath)
	if err != nil {
		response.Error(c, 500, "创建文件失败: "+err.Error())
		return
	}
	defer dst.Close()
	if _, err := io.Copy(dst, src); err != nil {
		response.Error(c, 500, "保存分片失败: "+err.Error())
		return
	}

	db.GetDB().Model(&model.FileUpload{}).Where("upload_id = ?", uploadID).Update("chunk_num", gorm.Expr("chunk_num + 1"))

	response.Success(c, gin.H{
		"upload_id": uploadID,
		"chunk_num": chunkNum,
		"saved":     true,
	})
}

// CompleteUpload 完成上传并合并 API: POST /upload/complete
func (ctrl *UploadController) CompleteUpload(c *gin.Context) {
	var req struct {
		UploadID  string `json:"upload_id" binding:"required"`
		TotalSize int64  `json:"total_size" binding:"required"`
		MD5       string `json:"md5"`
		Filename  string `json:"filename"`
		FileType  string `json:"file_type"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 400, "参数错误: "+err.Error())
		return
	}

	// 获取上传记录
	var uploadRecord model.FileUpload
	if err := db.GetDB().Where("upload_id = ?", req.UploadID).First(&uploadRecord).Error; err != nil {
		response.Error(c, 404, "上传记录不存在")
		return
	}

	dir := "./storage/uploads/chunks/" + req.UploadID
	outputDir := "./storage/uploads"
	os.MkdirAll(outputDir, 0755)

	// 使用原始文件名作为输出文件名
	filename := req.Filename
	if filename == "" {
		filename = uploadRecord.Filename
	}
	if filename == "" {
		filename = "uploaded_file" + time.Now().Format("_20060102_150405")
	}

	outputPath := filepath.Join(outputDir, filename)

	dst, err := os.Create(outputPath)
	if err != nil {
		response.Error(c, 500, "创建输出文件失败: "+err.Error())
		return
	}
	defer dst.Close()

	chunkCount := 0
	for i := 0; ; i++ {
		chunkPath := filepath.Join(dir, fmt.Sprintf("%d", i))
		src, err := os.Open(chunkPath)
		if err != nil {
			break
		}
		io.Copy(dst, src)
		src.Close()
		os.Remove(chunkPath)
		chunkCount++
	}

	// 计算实际文件大小
	fileInfo, _ := os.Stat(outputPath)
	actualSize := fileInfo.Size()

	// 确定文件类型/后缀
	fileExt := filepath.Ext(filename)
	fileType := req.FileType
	if fileType == "" {
		switch fileExt {
		case ".jpg", ".jpeg", ".png", ".gif", ".webp", ".svg":
			fileType = "image"
		case ".mp4", ".webm", ".avi", ".mov", ".mkv":
			fileType = "video"
		case ".pdf":
			fileType = "pdf"
		default:
			fileType = "file"
		}
	}

	// 更新上传记录
	db.GetDB().Model(&uploadRecord).Updates(map[string]interface{}{
		"status":      1,
		"file_path":   outputPath,
		"chunk_num":   chunkCount,
		"actual_size": actualSize,
	})

	// 获取用户信息（如果有）
	userID := uint(0)
	if uid, exists := c.Get("user_id"); exists {
		if uID, ok := uid.(uint); ok {
			userID = uID
		}
	}
	// 自动创建 MediaAsset 记录
	if userID > 0 {
		fileCategory := getFileCategory(fileType)
		mediaAsset := model.MediaAsset{
			UserID:       uint64(userID),
			UploadID:     req.UploadID,
			Name:         filename,
			OriginalName: filename,
			FilePath:     "/storage/uploads/" + filename,
			FileType:     fileCategory,
			MimeType:     getMimeType(filename),
			FileSize:     actualSize,
			Status:       1,
		}
		db.GetDB().Create(&mediaAsset)

		response.SuccessWithMsg(c, "文件合并成功", gin.H{
			"path":        outputPath,
			"url":         "/storage/uploads/" + filename,
			"filename":    filename,
			"size":        actualSize,
			"chunk_count": chunkCount,
			"asset_id":    mediaAsset.ID,
			"file_type":   fileType,
		})
	} else {
		response.SuccessWithMsg(c, "文件合并成功", gin.H{
			"path":        outputPath,
			"url":         "/storage/uploads/" + filename,
			"filename":    filename,
			"size":        actualSize,
			"chunk_count": chunkCount,
			"file_type":   fileType,
		})
	}
}

// getFileCategory 根据文件类型获取分类
func getFileCategory(fileType string) string {
	switch fileType {
	case "image":
		return "image"
	case "video":
		return "video"
	case "audio":
		return "audio"
	case "pdf", "doc", "docx", "xls", "xlsx":
		return "document"
	default:
		return "file"
	}
}

// getMimeType 根据文件名获取 MIME 类型
func getMimeType(filename string) string {
	ext := filepath.Ext(filename)
	mimeTypes := map[string]string{
		".jpg":  "image/jpeg",
		".jpeg": "image/jpeg",
		".png":  "image/png",
		".gif":  "image/gif",
		".webp": "image/webp",
		".svg":  "image/svg+xml",
		".mp4":  "video/mp4",
		".webm": "video/webm",
		".avi":  "video/x-msvideo",
		".mov":  "video/quicktime",
		".mkv":  "video/x-matroska",
		".pdf":  "application/pdf",
		".doc":  "application/msword",
		".docx": "application/vnd.openxmlformats-officedocument.wordprocessingml.document",
		".xls":  "application/vnd.ms-excel",
		".xlsx": "application/vnd/openxmlformats-officedocument.spreadsheetml.sheet",
	}
	if mime, ok := mimeTypes[ext]; ok {
		return mime
	}
	return "application/octet-stream"
}

// RegisterFileRequest 注册已有文件请求
type RegisterFileRequest struct {
	Name     string `json:"name" binding:"required"`
	FilePath string `json:"file_path" binding:"required"`
	FileType string `json:"file_type"`
	FileSize int64  `json:"file_size"`
	MimeType string `json:"mime_type"`
}

// RegisterFile 注册已有文件到媒体库 API: POST /upload/register
func (ctrl *UploadController) RegisterFile(c *gin.Context) {
	var req RegisterFileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 400, "参数错误: "+err.Error())
		return
	}

	// 从 JWT 中获取用户ID
	userID := uint(0)
	if uid, exists := c.Get("user_id"); exists {
		if uID, ok := uid.(uint); ok {
			userID = uID
		}
	}

	if userID == 0 {
		response.Error(c, 401, "未授权")
		return
	}

	// 如果未指定 file_type，根据文件名自动判断
	fileType := req.FileType
	if fileType == "" {
		ext := filepath.Ext(req.Name)
		switch ext {
		case ".jpg", ".jpeg", ".png", ".gif", ".webp", ".svg":
			fileType = "image"
		case ".mp4", ".webm", ".avi", ".mov", ".mkv":
			fileType = "video"
		case ".pdf":
			fileType = "pdf"
		default:
			fileType = "file"
		}
	}

	// 如果未指定 mime_type，根据文件名自动判断
	mimeType := req.MimeType
	if mimeType == "" {
		mimeType = getMimeType(req.Name)
	}

	// 检查是否已存在
	var existing model.MediaAsset
	result := db.GetDB().First(&existing, "name = ? OR file_path = ?", req.Name, req.FilePath)
	if result.Error == nil {
		response.Success(c, gin.H{
			"id":        existing.ID,
			"name":      existing.Name,
			"file_path": existing.FilePath,
			"message":   "文件已存在",
		})
		return
	}

	// 创建 MediaAsset 记录
	mediaAsset := model.MediaAsset{
		UserID:       uint64(userID),
		UploadID:     fmt.Sprintf("register_%d", time.Now().UnixNano()),
		Name:         req.Name,
		OriginalName: req.Name,
		FilePath:     req.FilePath,
		FileType:     fileType,
		MimeType:     mimeType,
		FileSize:     req.FileSize,
		Status:       1,
	}

	if err := db.GetDB().Create(&mediaAsset).Error; err != nil {
		response.Error(c, 500, "创建记录失败: "+err.Error())
		return
	}

	response.SuccessWithMsg(c, "注册成功", gin.H{
		"id":        mediaAsset.ID,
		"name":      mediaAsset.Name,
		"file_path": mediaAsset.FilePath,
		"url":       mediaAsset.FilePath,
	})
}
