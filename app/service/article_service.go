package service

import (
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"os"
	"path/filepath"
	"time"

	"gocms_v3/app/db"
	"gocms_v3/app/model"
	"gocms_v3/app/response"
)

// ArticleService 文章服务
type ArticleService struct{}

// NewArticleService 创建文章服务
func NewArticleService() *ArticleService {
	return &ArticleService{}
}

// CreateArticleRequest 创建文章请求
type CreateArticleRequest struct {
	Title      string `json:"title" binding:"required"`
	Content    string `json:"content"`
	Summary    string `json:"summary"`
	Cover      string `json:"cover"`
	CategoryID uint   `json:"category_id"`
	Status     int    `json:"status"`
	AuthorID   uint   `json:"-"`
}

// CreateArticle 创建文章
func (s *ArticleService) CreateArticle(req CreateArticleRequest) error {
	log.Printf("[SERVICE] 创建文章请求 - 作者ID: %d, 标题: %s", req.AuthorID, req.Title)

	article := &model.Article{
		Title:      req.Title,
		Content:    req.Content,
		Summary:    req.Summary,
		CoverImage: req.Cover,
		CategoryID: req.CategoryID,
		Status:     0, // 默认草稿状态，需要审核
		AuthorID:   req.AuthorID,
	}

	if err := db.GetDB().Create(article).Error; err != nil {
		log.Printf("[SERVICE-ERROR] 创建文章失败 - 作者ID: %d, 标题: %s, 错误: %v", req.AuthorID, req.Title, err)
		return err
	}

	// 自动创建审核记录（待审核状态）
	auditService := NewContentAuditService()
	_, err := auditService.SubmitAudit(struct {
		ContentType   string
		ContentID     uint
		ContentTitle  string
		SubmitterID   uint
		SubmitterName string
		NextStatus    int8
	}{
		ContentType:   "article",
		ContentID:     article.ID,
		ContentTitle:  article.Title,
		SubmitterID:   req.AuthorID,
		SubmitterName: fmt.Sprintf("用户%d", req.AuthorID),
		NextStatus:    0, // 待审核
	})
	if err != nil {
		log.Printf("[SERVICE-WARN] 创建文章审核记录失败 - 文章ID: %d, 错误: %v", article.ID, err)
	} else {
		log.Printf("[SERVICE] 文章审核记录已创建 - 编号: AUD%d0001, 文章ID: %d", time.Now().UnixNano(), article.ID)
	}

	log.Printf("[SERVICE] 创建文章成功 - 文章ID: %d, 作者ID: %d, 标题: %s", article.ID, req.AuthorID, req.Title)
	return nil
}

// UpdateArticleRequest 更新文章请求
type UpdateArticleRequest struct {
	Title   *string
	Content *string
	Status  int
}

// UpdateArticle 更新文章
func (s *ArticleService) UpdateArticle(id uint, req UpdateArticleRequest) error {
	var article model.Article
	if err := db.GetDB().First(&article, id).Error; err != nil {
		log.Printf("[SERVICE-WARN] 更新文章失败 - 文章ID: %d, 原因: 文章不存在", id)
		return fmt.Errorf("文章不存在")
	}

	log.Printf("[SERVICE] 更新文章请求 - 文章ID: %d, 标题: %s", id, article.Title)

	updates := map[string]interface{}{}
	if req.Title != nil {
		updates["title"] = *req.Title
	}
	if req.Content != nil {
		updates["content"] = *req.Content
	}
	if req.Status != 0 {
		updates["status"] = req.Status
	}

	if err := db.GetDB().Model(&article).Updates(updates).Error; err != nil {
		log.Printf("[SERVICE-ERROR] 更新文章失败 - 文章ID: %d, 错误: %v", id, err)
		return err
	}

	log.Printf("[SERVICE] 更新文章成功 - 文章ID: %d", id)
	return nil
}

// DeleteArticle 删除文章（软删除）
func (s *ArticleService) DeleteArticle(id uint) error {
	log.Printf("[SERVICE] 删除文章请求 - 文章ID: %d", id)

	if err := db.GetDB().Delete(&model.Article{}, id).Error; err != nil {
		log.Printf("[SERVICE-ERROR] 删除文章失败 - 文章ID: %d, 错误: %v", id, err)
		return err
	}

	log.Printf("[SERVICE] 删除文章成功 - 文章ID: %d", id)
	return nil
}

// GetArticleList 获取文章列表（分页）
func (s *ArticleService) GetArticleList(page, pageSize int) (*response.PageData, error) {
	log.Printf("[SERVICE] 获取文章列表请求 - 页码: %d, 每页数量: %d", page, pageSize)

	var articles []model.Article
	var total int64

	offset := (page - 1) * pageSize
	if err := db.GetDB().Model(&model.Article{}).Where("deleted_at IS NULL").Count(&total).Error; err != nil {
		log.Printf("[SERVICE-ERROR] 获取文章列表失败 - 错误: %v", err)
		return nil, err
	}

	if err := db.GetDB().Preload("Tags").Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&articles).Error; err != nil {
		log.Printf("[SERVICE-ERROR] 获取文章列表失败 - 错误: %v", err)
		return nil, err
	}

	log.Printf("[SERVICE] 获取文章列表成功 - 总数: %d, 当前页: %d", total, page)
	return &response.PageData{
		List:     articles,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	}, nil
}

// GetArticleByID 获取文章详情
func (s *ArticleService) GetArticleByID(id uint) (*model.Article, error) {
	log.Printf("[SERVICE] 获取文章详情请求 - 文章ID: %d", id)

	var article model.Article
	if err := db.GetDB().Preload("Tags").First(&article, id).Error; err != nil {
		log.Printf("[SERVICE-WARN] 获取文章详情失败 - 文章ID: %d, 原因: 文章不存在", id)
		return nil, fmt.Errorf("文章不存在")
	}

	log.Printf("[SERVICE] 获取文章详情成功 - 文章ID: %d", id)
	return &article, nil
}

// GetCategoryList 获取文章分类列表（使用 ArticleCategory model）
func GetCategoryList() ([]model.ArticleCategory, error) {
	var categories []model.ArticleCategory
	if err := db.GetDB().Where("status = 1").Order("sort_order ASC, id ASC").Find(&categories).Error; err != nil {
		log.Printf("[SERVICE-ERROR] 查询分类失败 - 错误: %v", err)
		return nil, fmt.Errorf("查询分类失败: %w", err)
	}
	return categories, nil
}

// GetTagList 获取标签列表
func GetTagList() ([]model.Tag, error) {
	var tags []model.Tag
	if err := db.GetDB().Where("status = 1").Find(&tags).Error; err != nil {
		log.Printf("[SERVICE-ERROR] 查询标签失败 - 错误: %v", err)
		return nil, fmt.Errorf("查询标签失败: %w", err)
	}
	return tags, nil
}

// UploadFile 上传文件（实际实现）
func UploadFile(file *multipart.FileHeader, userID uint, tenantID uint) (string, error) {
	// 打开上传的文件
	src, err := file.Open()
	if err != nil {
		log.Printf("[SERVICE-ERROR] 打开上传文件失败: %v", err)
		return "", fmt.Errorf("打开文件失败: %w", err)
	}
	defer src.Close()

	// 创建存储目录
	uploadDir := "./storage/uploads"
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		log.Printf("[SERVICE-ERROR] 创建上传目录失败: %v", err)
		return "", fmt.Errorf("创建目录失败: %w", err)
	}

	// 生成文件名（使用时间戳避免重名）
	timestamp := time.Now().Format("20060102_150405")
	ext := filepath.Ext(file.Filename)
	filename := fmt.Sprintf("%s_%d%s", timestamp, userID, ext)

	// 创建目标文件
	dstPath := filepath.Join(uploadDir, filename)
	dst, err := os.Create(dstPath)
	if err != nil {
		log.Printf("[SERVICE-ERROR] 创建目标文件失败: %v", err)
		return "", fmt.Errorf("创建文件失败: %w", err)
	}
	defer dst.Close()

	// 复制文件内容
	if _, err := io.Copy(dst, src); err != nil {
		log.Printf("[SERVICE-ERROR] 保存文件失败: %v", err)
		return "", fmt.Errorf("保存文件失败: %w", err)
	}

	// 创建文件上传记录到数据库
	record := &model.FileUpload{
		UploadID: fmt.Sprintf("upload_%d_%s", userID, time.Now().Format("20060102150405")),
		Filename: file.Filename,
		Size:     file.Size,
		Status:   1, // 已完成
		FilePath: dstPath,
	}
	if err := db.GetDB().Create(record).Error; err != nil {
		log.Printf("[SERVICE-WARN] 创建文件上传记录失败，但文件已保存: %v", err)
		// 即使记录失败，也返回成功的路径（文件已实际保存）
		return dstPath, nil
	}

	log.Printf("[SERVICE] 文件上传成功 - 路径: %s, 大小: %d, ID: %d", dstPath, file.Size, record.ID)
	return dstPath, nil
}
