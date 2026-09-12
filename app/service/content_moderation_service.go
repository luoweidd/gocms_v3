package service

import (
	"fmt"
	"regexp"
	"strings"
	"sync"

	"gocms_v3/app/config"
	"gocms_v3/app/db"
	"gocms_v3/app/model"
)

// ContentModerationService 内容审核服务
type ContentModerationService struct {
	keywordCache []model.RiskKeyword
	cacheMu      sync.RWMutex
	cacheLoaded  bool
}

// NewContentModerationService 创建内容审核服务
func NewContentModerationService() *ContentModerationService {
	return &ContentModerationService{}
}

// loadKeywords 加载启用的风险关键词到缓存
func (s *ContentModerationService) loadKeywords() error {
	s.cacheMu.RLock()
	if s.cacheLoaded {
		s.cacheMu.RUnlock()
		return nil
	}
	s.cacheMu.RUnlock()

	s.cacheMu.Lock()
	defer s.cacheMu.Unlock()

	// 双重检查锁定
	if s.cacheLoaded {
		return nil
	}

	var keywords []model.RiskKeyword
	if err := db.GetDB().Where("enabled = ?", true).Find(&keywords).Error; err != nil {
		return fmt.Errorf("加载风险关键词失败: %w", err)
	}

	s.keywordCache = keywords
	s.cacheLoaded = true
	return nil
}

// CheckContent 检查内容风险
func (s *ContentModerationService) CheckContent(contentType string, contentID uint, content string) (*model.RiskResult, error) {
	if !config.GlobalConfig.Moderation.Enabled {
		return &model.RiskResult{IsSafe: true}, nil
	}

	if err := s.loadKeywords(); err != nil {
		return nil, err
	}

	s.cacheMu.RLock()
	keywords := s.keywordCache
	s.cacheMu.RUnlock()

	if len(keywords) == 0 {
		return &model.RiskResult{IsSafe: true}, nil
	}

	var hits []model.RiskKeywordHit
	maxLevel := int8(0)
	totalScore := 0

	for _, kw := range keywords {
		var matchedText string
		var found bool

		if kw.MatchType == 1 {
			// 精确匹配
			if strings.Contains(content, kw.Keyword) {
				matchedText = kw.Keyword
				found = true
			}
		} else {
			// 模糊匹配（正则）
			pattern := "*" + kw.Keyword + "*"
			pattern = strings.ReplaceAll(pattern, "*", ".*")
			match, _ := regexp.MatchString(strings.ToLower(pattern), strings.ToLower(content))
			if match {
				found = true
				matchedText = kw.Keyword
			}
		}

		if found {
			hit := model.RiskKeywordHit{
				ContentCategory: contentType,
				ContentID:       contentID,
				KeywordID:       kw.ID,
				Keyword:         kw.Keyword,
				MatchedText:     matchedText,
				Score:           kw.Level,
				Context:         s.getContext(content, kw.Keyword),
			}
			hits = append(hits, hit)

			if kw.Level > maxLevel {
				maxLevel = kw.Level
			}
			totalScore += int(kw.Level)
		}
	}

	// 确定建议操作
	moderation := config.GlobalConfig.Moderation
	suggestAction := getSuggestAction(maxLevel, totalScore, moderation)

	return &model.RiskResult{
		IsSafe:        suggestAction == "auto_pass",
		MaxLevel:      maxLevel,
		TotalScore:    totalScore,
		HitCount:      len(hits),
		Hits:          hits,
		SuggestAction: suggestAction,
	}, nil
}

// getContext 获取匹配关键词的上下文
func (s *ContentModerationService) getContext(content, keyword string) string {
	idx := strings.Index(content, keyword)
	if idx == -1 {
		return content
	}

	start := 0
	end := len(content)

	if idx > 50 {
		start = idx - 50
	}
	if idx+len(keyword)+50 < end {
		end = idx + len(keyword) + 50
	}

	ctx := content[start:end]
	if start > 0 {
		ctx = "..." + ctx
	}
	if end < len(content) {
		ctx = ctx + "..."
	}
	return ctx
}

// getSuggestAction 根据风险等级和评分获取建议操作
func getSuggestAction(maxLevel int8, totalScore int, moderation config.ModerationConfig) string {
	if maxLevel >= moderation.BlockLevel {
		return "block"
	}
	if maxLevel >= moderation.RequireManualReview {
		return "review"
	}
	if totalScore >= int(moderation.AutoPublishThreshold)+1 {
		return "review"
	}
	return "auto_pass"
}

// AutoAudit 自动审核并决定内容状态
func (s *ContentModerationService) AutoAudit(contentType string, contentID uint, content string) (int8, error) {
	result, err := s.CheckContent(contentType, contentID, content)
	if err != nil {
		return 0, err
	}

	// 保存命中记录
	if len(result.Hits) > 0 {
		for i := range result.Hits {
			result.Hits[i].ContentID = contentID
			db.GetDB().Create(&result.Hits[i])
		}
	}

	switch result.SuggestAction {
	case "block":
		return 2, nil // 拒绝
	case "review":
		return 0, nil // 待审核
	case "auto_pass":
		return 1, nil // 通过
	default:
		return 0, nil // 默认待审核
	}
}

// RecordKeywordHit 记录关键词命中（用于分析）
func (s *ContentModerationService) RecordKeywordHit(hit model.RiskKeywordHit) error {
	return db.GetDB().Create(&hit).Error
}

// GetRiskStats 获取风险统计
func (s *ContentModerationService) GetRiskStats() (map[string]interface{}, error) {
	totalKeywords := int64(0)
	politicsKeywords := int64(0)
	violenceKeywords := int64(0)
	sexKeywords := int64(0)
	totalHits := int64(0)

	db.GetDB().Model(&model.RiskKeyword{}).Where("enabled = ?", true).Count(&totalKeywords)
	db.GetDB().Model(&model.RiskKeyword{}).Where("enabled = ? AND category = ?", true, "politics").Count(&politicsKeywords)
	db.GetDB().Model(&model.RiskKeyword{}).Where("enabled = ? AND category = ?", true, "violence").Count(&violenceKeywords)
	db.GetDB().Model(&model.RiskKeyword{}).Where("enabled = ? AND category = ?", true, "sex").Count(&sexKeywords)
	db.GetDB().Model(&model.RiskKeywordHit{}).Count(&totalHits)

	return map[string]interface{}{
		"total_keywords":    totalKeywords,
		"politics_keywords": politicsKeywords,
		"violence_keywords": violenceKeywords,
		"sex_keywords":      sexKeywords,
		"total_hits":        totalHits,
	}, nil
}
