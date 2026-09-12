package controller

import (
	"fmt"
	"net/http"
	"strconv"

	"gocms_v3/app/model"
	"gocms_v3/app/response"
	"gocms_v3/app/service"

	"github.com/gin-gonic/gin"
)

// SeoController SEO控制器
type SeoController struct {
	seoService *service.SeoService
}

// NewSeoController 创建SEO控制器
func NewSeoController() *SeoController {
	return &SeoController{}
}

// SetService 设置服务（避免循环依赖）
func (c *SeoController) SetService(svc *service.SeoService) {
	c.seoService = svc
}

// CreateSeoSetting 创建SEO设置
func (c *SeoController) CreateSeoSetting(ginCtx *gin.Context) {
	var req model.SeoCreateRequest
	if err := ginCtx.ShouldBindJSON(&req); err != nil {
		response.Error(ginCtx, 400, "参数错误: "+err.Error())
		return
	}

	setting, err := c.seoService.CreateSeoSetting(ginCtx.Request.Context(), req)
	if err != nil {
		response.Error(ginCtx, 500, err.Error())
		return
	}

	response.SuccessWithMsg(ginCtx, "SEO设置创建成功", setting)
}

// UpdateSeoSetting 更新SEO设置
func (c *SeoController) UpdateSeoSetting(ginCtx *gin.Context) {
	idStr := ginCtx.Param("id")
	var id uint
	if _, err := fmt.Sscanf(idStr, "%d", &id); err != nil || id == 0 {
		response.Error(ginCtx, 400, "无效的SEO设置ID")
		return
	}

	var req model.SeoUpdateRequest
	if err := ginCtx.ShouldBindJSON(&req); err != nil {
		response.Error(ginCtx, 400, "参数错误: "+err.Error())
		return
	}

	setting, err := c.seoService.UpdateSeoSetting(id, req)
	if err != nil {
		response.Error(ginCtx, 500, err.Error())
		return
	}

	response.SuccessWithMsg(ginCtx, "SEO设置更新成功", setting)
}

// GetSeoList 获取SEO列表
func (c *SeoController) GetSeoList(ginCtx *gin.Context) {
	query := model.SeoSettingQuery{
		Page:         1,
		PageSize:     10,
		ResourceType: ginCtx.Query("resource_type"),
		Keyword:      ginCtx.Query("keyword"),
	}

	if pageStr := ginCtx.Query("page"); pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			query.Page = p
		}
	}

	if pageSizeStr := ginCtx.Query("page_size"); pageSizeStr != "" {
		if ps, err := strconv.Atoi(pageSizeStr); err == nil && ps > 0 {
			query.PageSize = ps
		}
	}

	list, total, err := c.seoService.GetSeoList(query)
	if err != nil {
		response.Error(ginCtx, 500, err.Error())
		return
	}

	response.Success(ginCtx, gin.H{
		"list":     list,
		"total":    total,
		"page":     query.Page,
		"pageSize": query.PageSize,
	})
}

// GetSeoDetail 获取SEO详情
func (c *SeoController) GetSeoDetail(ginCtx *gin.Context) {
	idStr := ginCtx.Param("id")
	var id uint
	if _, err := fmt.Sscanf(idStr, "%d", &id); err != nil || id == 0 {
		response.Error(ginCtx, 400, "无效的SEO设置ID")
		return
	}

	detail, err := c.seoService.GetSeoSetting(id)
	if err != nil {
		response.Error(ginCtx, 404, err.Error())
		return
	}

	response.Success(ginCtx, detail)
}

// DeleteSeoSetting 删除SEO设置
func (c *SeoController) DeleteSeoSetting(ginCtx *gin.Context) {
	idStr := ginCtx.Param("id")
	var id uint
	if _, err := fmt.Sscanf(idStr, "%d", &id); err != nil || id == 0 {
		response.Error(ginCtx, 400, "无效的SEO设置ID")
		return
	}

	if err := c.seoService.DeleteSeoSetting(id); err != nil {
		response.Error(ginCtx, 500, err.Error())
		return
	}

	response.SuccessWithMsg(ginCtx, "SEO设置删除成功", nil)
}

// GetSitemapConfigs 获取站点地图配置
func (c *SeoController) GetSitemapConfigs(ginCtx *gin.Context) {
	configs, err := c.seoService.GetSitemapConfigList()
	if err != nil {
		response.Error(ginCtx, 500, err.Error())
		return
	}

	response.Success(ginCtx, configs)
}

// GenerateSitemap 生成站点地图
func (c *SeoController) GenerateSitemap(ginCtx *gin.Context) {
	xml, err := c.seoService.GenerateSitemap()
	if err != nil {
		response.Error(ginCtx, 500, err.Error())
		return
	}

	ginCtx.Header("Content-Type", "application/xml")
	ginCtx.String(http.StatusOK, xml)
}

// GetSitemapStats 获取站点地图统计
func (c *SeoController) GetSitemapStats(ginCtx *gin.Context) {
	stats, err := c.seoService.GetSitemapStats()
	if err != nil {
		response.Error(ginCtx, 500, err.Error())
		return
	}

	response.Success(ginCtx, stats)
}
