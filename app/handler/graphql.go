package handler

import (
	"fmt"

	"gocms_v3/app/model"
	"gocms_v3/app/service"
)

// GraphQLHandler GraphQL请求处理器
type GraphQLHandler struct {
	articleService *service.ArticleService
	videoService   *service.VideoService
	userService    *service.UserService
	commentService *service.CommentService
}

// NewGraphQLHandler 创建GraphQL处理器
func NewGraphQLHandler(
	articleService *service.ArticleService,
	videoService *service.VideoService,
	userService *service.UserService,
	commentService *service.CommentService,
) *GraphQLHandler {
	return &GraphQLHandler{
		articleService: articleService,
		videoService:   videoService,
		userService:    userService,
		commentService: commentService,
	}
}

// QueryParams GraphQL查询参数
type QueryParams struct {
	Query     string  `form:"query" binding:"required"`
	Variables *string `form:"variables"`
	Operation *string `form:"operationName"`
	Page      int     `form:"page"`
	PageSize  int     `form:"pageSize"`
	ID        float64 `form:"id"`
	Keyword   string  `form:"keyword"`
}

// ExecuteGraphQL 执行GraphQL查询（简化版HTTP接口）
func (h *GraphQLHandler) ExecuteGraphQL(query string, variables map[string]interface{}) (interface{}, error) {
	// 根据查询类型分发到对应的service方法
	if contains(query, "articles") {
		page := 1
		pageSize := 20
		if p, ok := variables["page"]; ok {
			if val, ok := p.(float64); ok {
				page = int(val)
			}
		}
		if ps, ok := variables["pageSize"]; ok {
			if val, ok := ps.(float64); ok {
				pageSize = int(val)
			}
		}
		return h.articleService.GetArticleList(page, pageSize)
	}
	if contains(query, "videos") {
		page := 1
		pageSize := 20
		if p, ok := variables["page"]; ok {
			if val, ok := p.(float64); ok {
				page = int(val)
			}
		}
		if ps, ok := variables["pageSize"]; ok {
			if val, ok := ps.(float64); ok {
				pageSize = int(val)
			}
		}
		params := service.VideoQueryParams{Page: page, PageSize: pageSize}
		return h.videoService.GetVideoList(params)
	}
	if contains(query, "comments") {
		var articleID, videoID float64
		if v, ok := variables["articleID"]; ok {
			articleID = toFloat64(v)
		}
		if v, ok := variables["videoID"]; ok {
			videoID = toFloat64(v)
		}
		queryParams := model.CommentListQuery{
			Page:     1,
			PageSize: 20,
		}
		if articleID > 0 {
			queryParams.ArticleID = uint(articleID)
		}
		if videoID > 0 {
			queryParams.VideoID = uint(videoID)
		}
		return h.commentService.GetCommentList(queryParams)
	}

	return nil, fmt.Errorf("不支持的查询类型")
}

func contains(s, substr string) bool {
	for i := 0; i < len(s)-len(substr)+1; i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

func toFloat64(v interface{}) float64 {
	switch val := v.(type) {
	case float64:
		return val
	case int:
		return float64(val)
	default:
		return 0
	}
}
