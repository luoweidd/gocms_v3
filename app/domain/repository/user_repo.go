package repository

import "gocms_v3/app/domain/entity"

// UserRepository 用户仓储接口
type UserRepository interface {
	// 创建用户
	Create(ctx interface{}, user *entity.User) error
	// 根据ID查找用户
	FindByID(ctx interface{}, id uint) (*entity.User, error)
	// 根据用户名查找用户
	FindByUsername(ctx interface{}, username string) (*entity.User, error)
	// 根据邮箱查找用户
	FindByEmail(ctx interface{}, email string) (*entity.User, error)
	// 根据手机号查找用户
	FindByPhone(ctx interface{}, phone string) (*entity.User, error)
	// 根据OpenID和Provider查找用户（OAuth2）
	FindByProvider(ctx interface{}, provider, openID string) (*entity.User, error)
	// 更新用户
	Update(ctx interface{}, user *entity.User) error
	// 删除用户
	Delete(ctx interface{}, id uint) error
	// 列出用户（分页）
	List(ctx interface{}, page, pageSize int) ([]entity.User, int64, error)
	// 根据租户查找用户
	FindByTenant(ctx interface{}, tenantID uint, page, pageSize int) ([]entity.User, int64, error)
	// 检查用户名是否存在
	UsernameExists(ctx interface{}, username string) (bool, error)
}

// ArticleRepository 文章仓储接口
type ArticleRepository interface {
	Create(ctx interface{}, article *entity.Article) error
	FindByID(ctx interface{}, id uint) (*entity.Article, error)
	FindBySlug(ctx interface{}, slug string) (*entity.Article, error)
	FindByTenant(ctx interface{}, tenantID uint, page, pageSize int) ([]entity.Article, int64, error)
	Update(ctx interface{}, article *entity.Article) error
	Delete(ctx interface{}, id uint) error
	ListPublished(ctx interface{}, page, pageSize int) ([]entity.Article, int64, error)
	ListByCategory(ctx interface{}, categoryID, page, pageSize int) ([]entity.Article, int64, error)
	FindByTenantAndStatus(ctx interface{}, tenantID, status uint, page, pageSize int) ([]entity.Article, int64, error)
}

// ArticleCategoryRepository 文章分类仓储接口
type ArticleCategoryRepository interface {
	Create(ctx interface{}, cat *entity.ArticleCategory) error
	FindByID(ctx interface{}, id uint) (*entity.ArticleCategory, error)
	GetTree(ctx interface{}) ([]entity.ArticleCategory, error)
	Update(ctx interface{}, cat *entity.ArticleCategory) error
	Delete(ctx interface{}, id uint) error
	ListByTenant(ctx interface{}, tenantID uint) ([]entity.ArticleCategory, error)
}

// VideoRepository 视频仓储接口
type VideoRepository interface {
	Create(ctx interface{}, video *entity.Video) error
	FindByID(ctx interface{}, id uint) (*entity.Video, error)
	Update(ctx interface{}, video *entity.Video) error
	Delete(ctx interface{}, id uint) error
	ListPublished(ctx interface{}, page, pageSize int) ([]entity.Video, int64, error)
	FindByTenant(ctx interface{}, tenantID uint, page, pageSize int) ([]entity.Video, int64, error)
	IncViewCount(ctx interface{}, id uint) error
	IncLikeCount(ctx interface{}, id uint) error
}

// CommentRepository 评论仓储接口
type CommentRepository interface {
	Create(ctx interface{}, comment *entity.Comment) error
	FindByID(ctx interface{}, id uint) (*entity.Comment, error)
	ListByTarget(ctx interface{}, targetType, targetID string, page, pageSize int) ([]entity.Comment, int64, error)
	Update(ctx interface{}, comment *entity.Comment) error
	Delete(ctx interface{}, id uint) error
	CountByUser(ctx interface{}, userID uint) (int64, error)
}

// OAuth2Repository OAuth2仓储接口
type OAuth2Repository interface {
	CreateAuthCode(ctx interface{}, code *entity.OAuth2AuthCode) error
	FindByCode(ctx interface{}, code string) (*entity.OAuth2AuthCode, error)
	MarkAsUsed(ctx interface{}, code string) error
	DeleteExpired(ctx interface{}) error
}

// TenantRepository 租户仓储接口
type TenantRepository interface {
	Create(ctx interface{}, tenant *entity.Tenant) error
	FindByID(ctx interface{}, id uint) (*entity.Tenant, error)
	FindByDomain(ctx interface{}, domain string) (*entity.Tenant, error)
	List(ctx interface{}, page, pageSize int) ([]entity.Tenant, int64, error)
	Update(ctx interface{}, tenant *entity.Tenant) error
	Delete(ctx interface{}, id uint) error
}

// ContentVersionRepository 内容版本仓储接口
type ContentVersionRepository interface {
	Create(ctx interface{}, version *entity.ContentVersion) error
	FindByTargetAndVersion(ctx interface{}, targetType, targetID string, version int) (*entity.ContentVersion, error)
	ListByTarget(ctx interface{}, targetType, targetID string) ([]entity.ContentVersion, error)
	Delete(ctx interface{}, id uint) error
}

// CronJobRepository 定时任务仓储接口
type CronJobRepository interface {
	Create(ctx interface{}, job *entity.CronJob) error
	FindByID(ctx interface{}, id uint) (*entity.CronJob, error)
	FindPending(ctx interface{}) ([]entity.CronJob, error)
	Update(ctx interface{}, job *entity.CronJob) error
	Delete(ctx interface{}, id uint) error
}

// ==================== 角色仓储接口（补充）====================

// RoleRepository 角色仓储接口
type RoleRepository interface {
	Create(ctx interface{}, role *entity.Role) error
	FindByID(ctx interface{}, id uint) (*entity.Role, error)
	FindByCode(ctx interface{}, code string) (*entity.Role, error)
	List(ctx interface{}, page, pageSize int) ([]entity.Role, int64, error)
	Update(ctx interface{}, role *entity.Role) error
	Delete(ctx interface{}, id uint) error
	ListByKeyword(ctx interface{}, keyword string, page, pageSize int) ([]entity.Role, int64, error)
}

// ==================== 标签仓储接口（补充）====================

// TagRepository 标签仓储接口
type TagRepository interface {
	Create(ctx interface{}, tag *entity.Tag) error
	FindByID(ctx interface{}, id uint) (*entity.Tag, error)
	FindByName(ctx interface{}, name string) (*entity.Tag, error)
	List(ctx interface{}, page, pageSize int) ([]entity.Tag, int64, error)
	Update(ctx interface{}, tag *entity.Tag) error
	Delete(ctx interface{}, id uint) error
	ListByKeyword(ctx interface{}, keyword string, page, pageSize int) ([]entity.Tag, int64, error)
}
