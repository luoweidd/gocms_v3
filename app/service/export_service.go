package service

import (
	"encoding/csv"
	"fmt"
	"gocms_v3/app/db"
	"gocms_v3/app/model"
	"os"
)

// ExportService 数据导出服务
type ExportService struct{}

func NewExportService() *ExportService {
	return &ExportService{}
}

// ExportArticlesToCSV 导出文章列表为CSV文件
func (s *ExportService) ExportArticlesToCSV(page, pageSize int) (string, error) {
	file, err := os.Create("./storage/exports/articles.csv")
	if err != nil {
		return "", fmt.Errorf("创建CSV文件失败: %w", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// 写入表头
	headers := []string{"ID", "标题", "分类ID", "状态", "阅读量", "创建时间"}
	if err := writer.Write(headers); err != nil {
		return "", fmt.Errorf("写入CSV表头失败: %w", err)
	}

	var articles []model.Article
	db.GetDB().Order("created_at DESC").Limit(pageSize).Offset((page - 1) * pageSize).Find(&articles)

	for _, article := range articles {
		row := []string{
			fmt.Sprintf("%d", article.ID),
			article.Title,
			fmt.Sprintf("%d", article.CategoryID),
			fmt.Sprintf("%d", article.Status),
			fmt.Sprintf("%d", article.ViewCount),
			article.CreatedAt.Format("2006-01-02 15:04:05"),
		}
		if err := writer.Write(row); err != nil {
			return "", fmt.Errorf("写入CSV数据失败: %w", err)
		}
	}

	return "./storage/exports/articles.csv", nil
}

// ExportVideosToCSV 导出视频列表为CSV文件
func (s *ExportService) ExportVideosToCSV(page, pageSize int) (string, error) {
	file, err := os.Create("./storage/exports/videos.csv")
	if err != nil {
		return "", fmt.Errorf("创建CSV文件失败: %w", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	headers := []string{"ID", "标题", "URL", "分类ID", "年份", "状态", "播放量"}
	if err := writer.Write(headers); err != nil {
		return "", fmt.Errorf("写入CSV表头失败: %w", err)
	}

	var videos []model.Video
	db.GetDB().Order("created_at DESC").Limit(pageSize).Offset((page - 1) * pageSize).Find(&videos)

	for _, video := range videos {
		row := []string{
			fmt.Sprintf("%d", video.ID),
			video.Title,
			video.URL,
			fmt.Sprintf("%d", video.CategoryID),
			fmt.Sprintf("%d", video.Year),
			fmt.Sprintf("%d", video.Status),
			fmt.Sprintf("%d", video.ViewCount),
		}
		if err := writer.Write(row); err != nil {
			return "", fmt.Errorf("写入CSV数据失败: %w", err)
		}
	}

	return "./storage/exports/videos.csv", nil
}

// ExportCommentsToCSV 导出评论列表为CSV文件
func (s *ExportService) ExportCommentsToCSV() (string, error) {
	file, err := os.Create("./storage/exports/comments.csv")
	if err != nil {
		return "", fmt.Errorf("创建CSV文件失败: %w", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	headers := []string{"ID", "内容", "状态", "点赞数", "评论者", "IP地址", "创建时间"}
	if err := writer.Write(headers); err != nil {
		return "", fmt.Errorf("写入CSV表头失败: %w", err)
	}

	var comments []model.Comment
	db.GetDB().Order("created_at DESC").Find(&comments)

	for _, comment := range comments {
		row := []string{
			fmt.Sprintf("%d", comment.ID),
			comment.Content,
			fmt.Sprintf("%d", comment.Status),
			fmt.Sprintf("%d", comment.Likes),
			comment.Nickname,
			comment.IPAddress,
			comment.CreatedAt.Format("2006-01-02 15:04:05"),
		}
		if err := writer.Write(row); err != nil {
			return "", fmt.Errorf("写入CSV数据失败: %w", err)
		}
	}

	return "./storage/exports/comments.csv", nil
}
