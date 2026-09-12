package model

// RiskKeyword 风险关键词
type RiskKeyword struct {
	BaseModel
	Keyword     string `gorm:"size:100;uniqueIndex:idx_keyword_enabled;not null" json:"keyword"` // 关键词
	Level       int8   `gorm:"default:2;index:idx_level" json:"level"`                           // 风险等级: 1低 2中 3高 4阻断
	Category    string `gorm:"size:50;index:idx_category" json:"category"`                       // 分类: politics/violence/sex/adult/spam/etc
	Description string `gorm:"size:255" json:"description"`                                      // 描述
	Enabled     bool   `gorm:"default:true;index:idx_keyword_enabled" json:"enabled"`            // 是否启用
	MatchType   int8   `gorm:"default:1" json:"match_type"`                                      // 匹配方式: 1精确 2模糊
	CreatedBy   uint   `json:"created_by"`                                                       // 创建人ID
}

// RiskKeywordHit 风险关键词命中记录
type RiskKeywordHit struct {
	BaseModel
	ContentCategory string `gorm:"size:20" json:"content_category"` // article/video/comment
	ContentID       uint   `json:"content_id"`                      // 内容ID
	KeywordID       uint   `json:"keyword_id"`                      // 命中的关键词ID
	Keyword         string `gorm:"size:100" json:"keyword"`         // 命中的关键词
	MatchedText     string `gorm:"size:500" json:"matched_text"`    // 匹配到的文本
	Score           int8   `json:"score"`                           // 风险评分(关键词等级)
	Context         string `gorm:"size:1000" json:"context"`        // 上下文
}

func (RiskKeyword) TableName() string {
	return "risk_keywords"
}

func (RiskKeywordHit) TableName() string {
	return "risk_keyword_hits"
}

// RiskResult 风险检测结果
type RiskResult struct {
	IsSafe        bool             `json:"is_safe"`        // 是否安全
	MaxLevel      int8             `json:"max_level"`      // 最高风险等级
	TotalScore    int              `json:"total_score"`    // 总评分
	HitCount      int              `json:"hit_count"`      // 命中数量
	Hits          []RiskKeywordHit `json:"hits"`           // 命中详情
	SuggestAction string           `json:"suggest_action"` // 建议操作: auto_pass/review/block
}
