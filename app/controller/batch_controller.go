package controller

import (
	"fmt"
	"net/http"

	"gocms_v3/app/db"
	"gocms_v3/app/model"
	"gocms_v3/app/response"

	"github.com/gin-gonic/gin"
)

type BatchController struct{}

func NewBatchController() *BatchController {
	return &BatchController{}
}

// BatchCreateArticles API: POST /articles/batch
func (ctrl *BatchController) BatchCreateArticles(c *gin.Context) {
	var req struct {
		Items []struct {
			Title      string `json:"title" binding:"required"`
			Summary    string `json:"summary"`
			Content    string `json:"content"`
			CategoryID uint   `json:"category_id"`
			Status     int    `json:"status"`
		} `json:"items" binding:"required,min=1,max=50"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "请求参数错误："+err.Error())
		return
	}

	results := make([]map[string]interface{}, 0)
	failed := make([]string, 0)

	for i, item := range req.Items {
		article := &model.Article{
			Title:      item.Title,
			Summary:    item.Summary,
			Content:    item.Content,
			CategoryID: item.CategoryID,
			Status:     item.Status,
		}

		err := db.GetDB().Create(article).Error
		if err != nil {
			failed = append(failed, fmt.Sprintf("第%d条失败：%v", i+1, err))
			continue
		}
		results = append(results, map[string]interface{}{
			"success": true,
			"data":    article,
		})
	}

	if len(failed) > 0 {
		response.Success(c, gin.H{
			"created": len(results),
			"failed":  len(failed),
			"results": results,
			"errors":  failed,
		})
		return
	}

	response.SuccessWithMsg(c, fmt.Sprintf("成功创建%d篇文章", len(results)), gin.H{"count": len(results)})
}

// BatchCreateVideos API: POST /videos/batch
func (ctrl *BatchController) BatchCreateVideos(c *gin.Context) {
	var req struct {
		Items []struct {
			Title      string `json:"title" binding:"required"`
			URL        string `json:"url" binding:"required"`
			CategoryID uint   `json:"category_id"`
			Year       int    `json:"year"`
			Status     int    `json:"status"`
		} `json:"items" binding:"required,min=1,max=50"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "请求参数错误："+err.Error())
		return
	}

	results := make([]map[string]interface{}, 0)
	failed := make([]string, 0)

	for i, item := range req.Items {
		video := &model.Video{
			Title:      item.Title,
			URL:        item.URL,
			CategoryID: item.CategoryID,
			Year:       item.Year,
			Status:     item.Status,
		}

		err := db.GetDB().Create(video).Error
		if err != nil {
			failed = append(failed, fmt.Sprintf("第%d条失败：%v", i+1, err))
			continue
		}
		results = append(results, map[string]interface{}{
			"success": true,
			"data":    video,
		})
	}

	if len(failed) > 0 {
		response.Success(c, gin.H{
			"created": len(results),
			"failed":  len(failed),
			"results": results,
			"errors":  failed,
		})
		return
	}

	response.SuccessWithMsg(c, fmt.Sprintf("成功创建%d个视频", len(results)), gin.H{"count": len(results)})
}
