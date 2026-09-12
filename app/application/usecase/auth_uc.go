package usecase

import (
	"context"
	"gocms_v3/app/domain/entity"
	"gocms_v3/app/domain/repository"
	"mime/multipart"
)

// AuthUseCase 认证用例
type AuthUseCase struct {
	userRepo   repository.UserRepository
	oauth2Repo repository.OAuth2Repository
	config     *OAuthConfig
}

// OAuthConfig OAuth配置
type OAuthConfig struct {
	ClientID     string
	ClientSecret string
	RedirectURL  string
	TokenURL     string
	AuthURL      string
}

// ArticleUseCase 文章用例
type ArticleUseCase struct {
	articleRepo  repository.ArticleRepository
	versionRepo  repository.ContentVersionRepository
	mediaService MediaService
}

// NewArticleUseCase 创建文章用例
func NewArticleUseCase(articleRepo repository.ArticleRepository, versionRepo repository.ContentVersionRepository) *ArticleUseCase {
	return &ArticleUseCase{
		articleRepo: articleRepo,
		versionRepo: versionRepo,
	}
}

// MediaService 媒体服务接口
type MediaService interface {
	UploadFile(ctx context.Context, file *multipart.FileHeader, userID uint, tenantID uint) (*entity.FileUploadRecord, error)
	DownloadFile(ctx context.Context, recordID uint) (string, error)
	DeleteFile(ctx context.Context, recordID uint) error
}

// WebhookService 回调服务接口
type WebhookService interface {
	RegisterWebhook(ctx context.Context, webhook *entity.Webhook) error
	TriggerWebhook(ctx context.Context, event string, payload map[string]interface{}) error
	ListWebhooks(ctx context.Context, eventType string) ([]*entity.Webhook, error)
}

// CronScheduleService 定时调度服务接口
type CronScheduleService interface {
	Start() error
	Stop() error
	RegisterJob(ctx context.Context, job *entity.CronJob) error
	RemoveJob(ctx context.Context, jobID uint) error
	ListJobs(ctx context.Context) ([]*entity.CronJob, error)
}

// i18nService 国际化服务接口
type i18nService interface {
	GetTranslation(ctx context.Context, key, lang string) (string, error)
	SaveTranslation(ctx context.Context, t *entity.Translation) error
	ListLanguages(ctx context.Context) ([]entity.Language, error)
}

// CDNService CDN服务接口
type CDNService interface {
	PurgeCache(ctx context.Context, urls []string) error
	GetCDNURL(ctx context.Context, path string) string
	UploadToCDN(ctx context.Context, file *multipart.FileHeader) (string, error)
}
