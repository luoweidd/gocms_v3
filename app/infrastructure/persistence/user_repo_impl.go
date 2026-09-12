package persistence

import (
	"context"
	"fmt"
	"time"

	"gocms_v3/app/domain/entity"
	"gocms_v3/app/domain/repository"

	"gorm.io/gorm"
)

// userRepositoryImpl UserRepository的MySQL实现
type userRepositoryImpl struct {
	db *gorm.DB
}

// NewUserRepository 创建用户仓储实现
func NewUserRepository(db *gorm.DB) repository.UserRepository {
	return &userRepositoryImpl{db: db}
}

// Create 创建用户
func (r *userRepositoryImpl) Create(ctx interface{}, user *entity.User) error {
	return r.db.WithContext(ctx.(context.Context)).Create(user).Error
}

// FindByID 根据ID查找用户
func (r *userRepositoryImpl) FindByID(ctx interface{}, id uint) (*entity.User, error) {
	var user entity.User
	err := r.db.WithContext(ctx.(context.Context)).First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// FindByUsername 根据用户名查找用户
func (r *userRepositoryImpl) FindByUsername(ctx interface{}, username string) (*entity.User, error) {
	var user entity.User
	err := r.db.WithContext(ctx.(context.Context)).Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// FindByEmail 根据邮箱查找用户
func (r *userRepositoryImpl) FindByEmail(ctx interface{}, email string) (*entity.User, error) {
	var user entity.User
	err := r.db.WithContext(ctx.(context.Context)).Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// FindByPhone 根据手机号查找用户
func (r *userRepositoryImpl) FindByPhone(ctx interface{}, phone string) (*entity.User, error) {
	var user entity.User
	err := r.db.WithContext(ctx.(context.Context)).Where("phone = ?", phone).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// FindByProvider 根据Provider和OpenID查找用户
func (r *userRepositoryImpl) FindByProvider(ctx interface{}, provider, openID string) (*entity.User, error) {
	var user entity.User
	query := r.db.WithContext(ctx.(context.Context))
	switch provider {
	case "wechat":
		query = query.Where("wechat_openid = ?", openID)
	case "qq":
		query = query.Where("qq_openid = ?", openID)
	case "github":
		query = query.Where("github_openid = ?", openID)
	default:
		return nil, fmt.Errorf("不支持的提供商: %s", provider)
	}
	err := query.First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// Update 更新用户
func (r *userRepositoryImpl) Update(ctx interface{}, user *entity.User) error {
	return r.db.WithContext(ctx.(context.Context)).Save(user).Error
}

// Delete 删除用户
func (r *userRepositoryImpl) Delete(ctx interface{}, id uint) error {
	return r.db.WithContext(ctx.(context.Context)).Delete(&entity.User{}, id).Error
}

// List 列出用户（分页）
func (r *userRepositoryImpl) List(ctx interface{}, page, pageSize int) ([]entity.User, int64, error) {
	var users []entity.User
	var total int64

	query := r.db.WithContext(ctx.(context.Context)).Model(&entity.User{})
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Find(&users).Error; err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

// FindByTenant 根据租户查找用户
func (r *userRepositoryImpl) FindByTenant(ctx interface{}, tenantID uint, page, pageSize int) ([]entity.User, int64, error) {
	var users []entity.User
	var total int64

	query := r.db.WithContext(ctx.(context.Context)).Model(&entity.User{}).Where("tenant_id = ?", tenantID)
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Find(&users).Error; err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

// UsernameExists 检查用户名是否存在
func (r *userRepositoryImpl) UsernameExists(ctx interface{}, username string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx.(context.Context)).Model(&entity.User{}).Where("username = ?", username).Count(&count).Error
	return count > 0, err
}

// articleRepositoryImpl ArticleRepository的MySQL实现
type articleRepositoryImpl struct {
	db *gorm.DB
}

// NewArticleRepository 创建文章仓储实现
func NewArticleRepository(db *gorm.DB) repository.ArticleRepository {
	return &articleRepositoryImpl{db: db}
}

// Create 创建文章
func (r *articleRepositoryImpl) Create(ctx interface{}, article *entity.Article) error {
	return r.db.WithContext(ctx.(context.Context)).Create(article).Error
}

// FindByID 根据ID查找文章
func (r *articleRepositoryImpl) FindByID(ctx interface{}, id uint) (*entity.Article, error) {
	var article entity.Article
	err := r.db.WithContext(ctx.(context.Context)).First(&article, id).Error
	if err != nil {
		return nil, err
	}
	return &article, nil
}

// FindBySlug 根据Slug查找文章
func (r *articleRepositoryImpl) FindBySlug(ctx interface{}, slug string) (*entity.Article, error) {
	var article entity.Article
	err := r.db.WithContext(ctx.(context.Context)).Where("slug = ?", slug).First(&article).Error
	if err != nil {
		return nil, err
	}
	return &article, nil
}

// FindByTenant 根据租户查找文章
func (r *articleRepositoryImpl) FindByTenant(ctx interface{}, tenantID uint, page, pageSize int) ([]entity.Article, int64, error) {
	var articles []entity.Article
	var total int64

	query := r.db.WithContext(ctx.(context.Context)).Model(&entity.Article{}).Where("tenant_id = ?", tenantID)
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Find(&articles).Error; err != nil {
		return nil, 0, err
	}

	return articles, total, nil
}

// Update 更新文章
func (r *articleRepositoryImpl) Update(ctx interface{}, article *entity.Article) error {
	return r.db.WithContext(ctx.(context.Context)).Save(article).Error
}

// Delete 删除文章
func (r *articleRepositoryImpl) Delete(ctx interface{}, id uint) error {
	return r.db.WithContext(ctx.(context.Context)).Delete(&entity.Article{}, id).Error
}

// ListPublished 列出已发布的文章（分页）
func (r *articleRepositoryImpl) ListPublished(ctx interface{}, page, pageSize int) ([]entity.Article, int64, error) {
	var articles []entity.Article
	var total int64

	query := r.db.WithContext(ctx.(context.Context)).Model(&entity.Article{}).Where("status = ?", 1)
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := query.Where("status = ?", 1).Offset(offset).Limit(pageSize).Find(&articles).Error; err != nil {
		return nil, 0, err
	}

	return articles, total, nil
}

// ListByCategory 根据分类列出文章（分页）
func (r *articleRepositoryImpl) ListByCategory(ctx interface{}, categoryID, page, pageSize int) ([]entity.Article, int64, error) {
	var articles []entity.Article
	var total int64

	query := r.db.WithContext(ctx.(context.Context)).Model(&entity.Article{}).Where("category_id = ?", categoryID)
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := query.Where("category_id = ?", categoryID).Offset(offset).Limit(pageSize).Find(&articles).Error; err != nil {
		return nil, 0, err
	}

	return articles, total, nil
}

// FindByTenantAndStatus 根据租户和状态查找文章
func (r *articleRepositoryImpl) FindByTenantAndStatus(ctx interface{}, tenantID, status uint, page, pageSize int) ([]entity.Article, int64, error) {
	var articles []entity.Article
	var total int64

	query := r.db.WithContext(ctx.(context.Context)).Model(&entity.Article{}).Where("tenant_id = ? AND status = ?", tenantID, status)
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := query.Where("tenant_id = ? AND status = ?", tenantID, status).Offset(offset).Limit(pageSize).Find(&articles).Error; err != nil {
		return nil, 0, err
	}

	return articles, total, nil
}

// videoRepositoryImpl VideoRepository的MySQL实现
type videoRepositoryImpl struct {
	db *gorm.DB
}

// NewVideoRepository 创建视频仓储实现
func NewVideoRepository(db *gorm.DB) repository.VideoRepository {
	return &videoRepositoryImpl{db: db}
}

// Create 创建视频
func (r *videoRepositoryImpl) Create(ctx interface{}, video *entity.Video) error {
	return r.db.WithContext(ctx.(context.Context)).Create(video).Error
}

// FindByID 根据ID查找视频
func (r *videoRepositoryImpl) FindByID(ctx interface{}, id uint) (*entity.Video, error) {
	var video entity.Video
	err := r.db.WithContext(ctx.(context.Context)).First(&video, id).Error
	if err != nil {
		return nil, err
	}
	return &video, nil
}

// Update 更新视频
func (r *videoRepositoryImpl) Update(ctx interface{}, video *entity.Video) error {
	return r.db.WithContext(ctx.(context.Context)).Save(video).Error
}

// Delete 删除视频
func (r *videoRepositoryImpl) Delete(ctx interface{}, id uint) error {
	return r.db.WithContext(ctx.(context.Context)).Delete(&entity.Video{}, id).Error
}

// ListPublished 列出已发布的视频（分页）
func (r *videoRepositoryImpl) ListPublished(ctx interface{}, page, pageSize int) ([]entity.Video, int64, error) {
	var videos []entity.Video
	var total int64

	query := r.db.WithContext(ctx.(context.Context)).Model(&entity.Video{}).Where("status = ?", 1)
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := query.Where("status = ?", 1).Offset(offset).Limit(pageSize).Find(&videos).Error; err != nil {
		return nil, 0, err
	}

	return videos, total, nil
}

// FindByTenant 根据租户查找视频
func (r *videoRepositoryImpl) FindByTenant(ctx interface{}, tenantID uint, page, pageSize int) ([]entity.Video, int64, error) {
	var videos []entity.Video
	var total int64

	query := r.db.WithContext(ctx.(context.Context)).Model(&entity.Video{}).Where("tenant_id = ?", tenantID)
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := query.Where("tenant_id = ?", tenantID).Offset(offset).Limit(pageSize).Find(&videos).Error; err != nil {
		return nil, 0, err
	}

	return videos, total, nil
}

// IncViewCount 增加浏览量
func (r *videoRepositoryImpl) IncViewCount(ctx interface{}, id uint) error {
	return r.db.WithContext(ctx.(context.Context)).Model(&entity.Video{}).Where("id = ?", id).Update("view_count", gorm.Expr("view_count + 1")).Error
}

// IncLikeCount 增加点赞数
func (r *videoRepositoryImpl) IncLikeCount(ctx interface{}, id uint) error {
	return r.db.WithContext(ctx.(context.Context)).Model(&entity.Video{}).Where("id = ?", id).Update("likes", gorm.Expr("likes + 1")).Error
}

// commentRepositoryImpl CommentRepository的MySQL实现
type commentRepositoryImpl struct {
	db *gorm.DB
}

// NewCommentRepository 创建评论仓储实现
func NewCommentRepository(db *gorm.DB) repository.CommentRepository {
	return &commentRepositoryImpl{db: db}
}

// Create 创建评论
func (r *commentRepositoryImpl) Create(ctx interface{}, comment *entity.Comment) error {
	return r.db.WithContext(ctx.(context.Context)).Create(comment).Error
}

// FindByID 根据ID查找评论
func (r *commentRepositoryImpl) FindByID(ctx interface{}, id uint) (*entity.Comment, error) {
	var comment entity.Comment
	err := r.db.WithContext(ctx.(context.Context)).First(&comment, id).Error
	if err != nil {
		return nil, err
	}
	return &comment, nil
}

// ListByTarget 根据目标列出评论（分页）
func (r *commentRepositoryImpl) ListByTarget(ctx interface{}, targetType, targetID string, page, pageSize int) ([]entity.Comment, int64, error) {
	var comments []entity.Comment
	var total int64

	query := r.db.WithContext(ctx.(context.Context)).Model(&entity.Comment{}).Where("target_type = ? AND CAST(target_id AS CHAR) = ?", targetType, targetID)
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := query.Where("target_type = ? AND CAST(target_id AS CHAR) = ?", targetType, targetID).Offset(offset).Limit(pageSize).Find(&comments).Error; err != nil {
		return nil, 0, err
	}

	return comments, total, nil
}

// Update 更新评论
func (r *commentRepositoryImpl) Update(ctx interface{}, comment *entity.Comment) error {
	return r.db.WithContext(ctx.(context.Context)).Save(comment).Error
}

// Delete 删除评论
func (r *commentRepositoryImpl) Delete(ctx interface{}, id uint) error {
	return r.db.WithContext(ctx.(context.Context)).Delete(&entity.Comment{}, id).Error
}

// CountByUser 统计用户评论数
func (r *commentRepositoryImpl) CountByUser(ctx interface{}, userID uint) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx.(context.Context)).Model(&entity.Comment{}).Where("user_id = ?", userID).Count(&count).Error
	return count, err
}

// oauth2RepositoryImpl OAuth2Repository的MySQL实现
type oauth2RepositoryImpl struct {
	db *gorm.DB
}

// NewOAuth2Repository 创建OAuth2仓储实现
func NewOAuth2Repository(db *gorm.DB) repository.OAuth2Repository {
	return &oauth2RepositoryImpl{db: db}
}

// CreateAuthCode 创建授权码
func (r *oauth2RepositoryImpl) CreateAuthCode(ctx interface{}, code *entity.OAuth2AuthCode) error {
	return r.db.WithContext(ctx.(context.Context)).Create(code).Error
}

// FindByCode 根据授权码查找
func (r *oauth2RepositoryImpl) FindByCode(ctx interface{}, codeStr string) (*entity.OAuth2AuthCode, error) {
	var authCode entity.OAuth2AuthCode
	err := r.db.WithContext(ctx.(context.Context)).Where("code = ?", codeStr).First(&authCode).Error
	if err != nil {
		return nil, err
	}
	return &authCode, nil
}

// MarkAsUsed 标记为已使用
func (r *oauth2RepositoryImpl) MarkAsUsed(ctx interface{}, codeStr string) error {
	return r.db.WithContext(ctx.(context.Context)).Model(&entity.OAuth2AuthCode{}).Where("code = ?", codeStr).Delete(&entity.OAuth2AuthCode{}).Error
}

// DeleteExpired 删除过期的授权码
func (r *oauth2RepositoryImpl) DeleteExpired(ctx interface{}) error {
	return r.db.WithContext(ctx.(context.Context)).Where("expires_at < ?", time.Now()).Delete(&entity.OAuth2AuthCode{}).Error
}

// tenantRepositoryImpl TenantRepository的MySQL实现
type tenantRepositoryImpl struct {
	db *gorm.DB
}

// NewTenantRepository 创建租户仓储实现
func NewTenantRepository(db *gorm.DB) repository.TenantRepository {
	return &tenantRepositoryImpl{db: db}
}

// Create 创建租户
func (r *tenantRepositoryImpl) Create(ctx interface{}, tenant *entity.Tenant) error {
	return r.db.WithContext(ctx.(context.Context)).Create(tenant).Error
}

// FindByID 根据ID查找租户
func (r *tenantRepositoryImpl) FindByID(ctx interface{}, id uint) (*entity.Tenant, error) {
	var tenant entity.Tenant
	err := r.db.WithContext(ctx.(context.Context)).First(&tenant, id).Error
	if err != nil {
		return nil, err
	}
	return &tenant, nil
}

// FindByDomain 根据域名查找租户
func (r *tenantRepositoryImpl) FindByDomain(ctx interface{}, domain string) (*entity.Tenant, error) {
	var tenant entity.Tenant
	err := r.db.WithContext(ctx.(context.Context)).Where("domain = ?", domain).First(&tenant).Error
	if err != nil {
		return nil, err
	}
	return &tenant, nil
}

// List 列出租户（分页）
func (r *tenantRepositoryImpl) List(ctx interface{}, page, pageSize int) ([]entity.Tenant, int64, error) {
	var tenants []entity.Tenant
	var total int64

	query := r.db.WithContext(ctx.(context.Context)).Model(&entity.Tenant{})
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Find(&tenants).Error; err != nil {
		return nil, 0, err
	}

	return tenants, total, nil
}

// Update 更新租户
func (r *tenantRepositoryImpl) Update(ctx interface{}, tenant *entity.Tenant) error {
	return r.db.WithContext(ctx.(context.Context)).Save(tenant).Error
}

// Delete 删除租户
func (r *tenantRepositoryImpl) Delete(ctx interface{}, id uint) error {
	return r.db.WithContext(ctx.(context.Context)).Delete(&entity.Tenant{}, id).Error
}

// contentVersionRepositoryImpl ContentVersionRepository的MySQL实现
type contentVersionRepositoryImpl struct {
	db *gorm.DB
}

// NewContentVersionRepository 创建内容版本仓储实现
func NewContentVersionRepository(db *gorm.DB) repository.ContentVersionRepository {
	return &contentVersionRepositoryImpl{db: db}
}

// Create 创建内容版本
func (r *contentVersionRepositoryImpl) Create(ctx interface{}, version *entity.ContentVersion) error {
	return r.db.WithContext(ctx.(context.Context)).Create(version).Error
}

// FindByTargetAndVersion 根据目标和版本号查找
func (r *contentVersionRepositoryImpl) FindByTargetAndVersion(ctx interface{}, targetType, targetID string, version int) (*entity.ContentVersion, error) {
	var versionEntity entity.ContentVersion
	err := r.db.WithContext(ctx.(context.Context)).Where("target_type = ? AND CAST(target_id AS CHAR) = ? AND version = ?", targetType, targetID, version).First(&versionEntity).Error
	if err != nil {
		return nil, err
	}
	return &versionEntity, nil
}

// ListByTarget 根据目标列出所有版本
func (r *contentVersionRepositoryImpl) ListByTarget(ctx interface{}, targetType, targetID string) ([]entity.ContentVersion, error) {
	var versions []entity.ContentVersion
	err := r.db.WithContext(ctx.(context.Context)).Where("target_type = ? AND CAST(target_id AS CHAR) = ?", targetType, targetID).Order("version ASC").Find(&versions).Error
	if err != nil {
		return nil, err
	}
	return versions, nil
}

// Delete 删除内容版本
func (r *contentVersionRepositoryImpl) Delete(ctx interface{}, id uint) error {
	return r.db.WithContext(ctx.(context.Context)).Delete(&entity.ContentVersion{}, id).Error
}

// cronJobRepositoryImpl CronJobRepository的MySQL实现
type cronJobRepositoryImpl struct {
	db *gorm.DB
}

// NewCronJobRepository 创建定时任务仓储实现
func NewCronJobRepository(db *gorm.DB) repository.CronJobRepository {
	return &cronJobRepositoryImpl{db: db}
}

// Create 创建定时任务
func (r *cronJobRepositoryImpl) Create(ctx interface{}, job *entity.CronJob) error {
	return r.db.WithContext(ctx.(context.Context)).Create(job).Error
}

// FindByID 根据ID查找定时任务
func (r *cronJobRepositoryImpl) FindByID(ctx interface{}, id uint) (*entity.CronJob, error) {
	var job entity.CronJob
	err := r.db.WithContext(ctx.(context.Context)).First(&job, id).Error
	if err != nil {
		return nil, err
	}
	return &job, nil
}

// FindPending 查找待执行的定时任务
func (r *cronJobRepositoryImpl) FindPending(ctx interface{}) ([]entity.CronJob, error) {
	var jobs []entity.CronJob
	err := r.db.WithContext(ctx.(context.Context)).Where("status = ? AND (next_run_at IS NULL OR next_run_at <= ?)", 1, time.Now()).Find(&jobs).Error
	if err != nil {
		return nil, err
	}
	return jobs, nil
}

// Update 更新定时任务
func (r *cronJobRepositoryImpl) Update(ctx interface{}, job *entity.CronJob) error {
	return r.db.WithContext(ctx.(context.Context)).Save(job).Error
}

// Delete 删除定时任务
func (r *cronJobRepositoryImpl) Delete(ctx interface{}, id uint) error {
	return r.db.WithContext(ctx.(context.Context)).Delete(&entity.CronJob{}, id).Error
}
