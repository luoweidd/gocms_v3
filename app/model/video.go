package model

import (
	"time"
)

// Video 视频内容模型
type Video struct {
	BaseModel
	// 基本信息
	Title    string `gorm:"size:60;not null;column:title;index:ftx_video_combined" json:"title"`
	Summary  string `gorm:"size:100;column:summary" json:"summary"`      // 视频摘要（不加入联合索引）
	Content  string `gorm:"type:text;column:content" json:"content"`     // 详细描述
	URL      string `gorm:"size:500;not null;column:url" json:"url"`     // 视频文件URL
	Cover    string `gorm:"size:500;column:cover" json:"cover"`          // 封面图URL
	Format   string `gorm:"size:20;column=format" json:"format"`         // 视频格式(mp4, webm, avi, mkv)
	Duration int    `gorm:"default:0;column:duration" json:"duration"`   // 时长(秒)
	FileSize int64  `gorm:"default:0;column:file_size" json:"file_size"` // 文件大小(字节)
	Width    int    `gorm:"default:0;column:width" json:"width"`         // 视频宽度(px)
	Height   int    `gorm:"default:0;column:height" json:"height"`       // 视频高度(px)
	Bitrate  int    `gorm:"default:0;column:bitrate" json:"bitrate"`     // 码率(kbps)

	// 影视属性
	Year           int        `gorm:"default:0;index:idx_video_year;column:year" json:"year"`                // 年份
	ReleaseDate    *time.Time `gorm:"column:release_date" json:"release_date"`                               // 上映/发布日期
	Director       string     `gorm:"size:30;column:director;index:ftx_video_combined" json:"director"`      // 导演
	Actors         string     `gorm:"size:60;column:actors;index:ftx_video_combined" json:"actors"`          // 演员(逗号分隔)
	Genre          string     `gorm:"size:30;index:idx_video_genre;index:ftx_video_combined" json:"genre"`   // 类型(逗号分隔: 动作,喜剧,剧情)
	Region         string     `gorm:"size:30;index:idx_video_region;index:ftx_video_combined" json:"region"` // 地区(如: 中国大陆, 美国, 日本)
	Language       string     `gorm:"size:100;column:language" json:"language"`                              // 语言(如: 普通话, 英语, 日语)
	Rating         float64    `gorm:"type:decimal(3,2);default:0;column:rating" json:"rating"`               // 评分(0-10)
	RatingCount    int        `gorm:"default:0;column:rating_count" json:"rating_count"`                     // 评分人数
	TotalEpisodes  int        `gorm:"default:0;column:total_episodes" json:"total_episodes"`                 // 总集数(0表示电影)
	CurrentEpisode int        `gorm:"default:0;column:current_episode" json:"current_episode"`               // 已更新集数

	// 分类与管理
	CategoryID  uint       `gorm:"index:idx_video_category;column:category_id" json:"category_id"`           // 视频分类ID
	AuthorID    uint       `gorm:"index:idx_video_author;column:author_id" json:"author_id"`                 // 作者ID
	Status      int        `gorm:"default:0;index:idx_video_status;column:status" json:"status"`             // 0草稿 1发布 2下架
	ViewCount   int        `gorm:"default:0;index:idx_video_view_count;column:view_count" json:"view_count"` // 播放次数
	IsTop       bool       `gorm:"default:false;index:idx_video_is_top;column:is_top" json:"is_top"`         // 是否置顶
	PublishedAt *time.Time `gorm:"index:idx_video_published;column:published_at" json:"published_at"`        // 发布时间

	// 关联
	Tags     []Tag          `gorm:"many2many:video_tags;" json:"tags"`                                                                              // 标签
	Category *VideoCategory `gorm:"foreignKey:CategoryID" json:"category"`                                                                          // 所属分类
	Audits   []ContentAudit `gorm:"foreignKey:ContentID;references:id;constraint:where content_type = 'video';preloadable" json:"audits,omitempty"` // 审核记录
}

func (Video) TableName() string {
	return "videos"
}
