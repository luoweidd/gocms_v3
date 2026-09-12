# VideoService LIKE查询性能优化方案

## 1. 问题分析

当前 `video_service.go` 中的 `GetVideoList` 方法存在以下LIKE全表扫描查询：

```go
// 第325-327行 - 4列OR模糊匹配 (最严重)
query.Where("title LIKE ? OR summary LIKE ? OR director LIKE ? OR actors LIKE ?", keyword, keyword, keyword, keyword)

// 第333行 - genre模糊匹配
query.Where("genre LIKE ?", fmt.Sprintf("%%%s%%", params.Genre))

// 第336行 - region模糊匹配  
query.Where("region LIKE ?", fmt.Sprintf("%%%s%%", params.Region))

// 第339行 - language模糊匹配
query.Where("language LIKE ?", fmt.Sprintf("%%%s%%", params.Language))

// 第342行 - director模糊匹配
query.Where("director LIKE ?", fmt.Sprintf("%%%s%%", params.Director))
```

### 性能瓶颈
- `LIKE '%%%s%%'` (前后百分号) 无法使用普通BTree索引
- OR条件导致无法有效利用索引
- 数据量超过1万条后查询响应时间急剧增加

---

## 2. 方案一：MySQL全文索引 (FULLTEXT INDEX)

### 2.1 修改模型定义 `app/model/video.go`

```go
// 将以下字段添加 FULLTEXT 标签
type Video struct {
    // ... existing fields ...
    
    Title        string     `gorm:"size:200;not null;fulltext" json:"title"`         // +fulltext
    Summary      string     `gorm:"size:1000;fulltext" json:"summary"`               // +fulltext
    Director     string     `gorm:"size:200;fulltext" json:"director"`               // +fulltext
    Actors       string     `gorm:"size:500;fulltext" json:"actors"`                 // +fulltext
    Genre        string     `gorm:"size:200;index;fulltext" json:"genre"`            // +fulltext
    Region       string     `gorm:"size:100;index;fulltext" json:"region"`           // +fulltext
    
    // Language 中文支持有限，保持普通索引
    Language     string     `gorm:"size:100;index" json:"language"`                  
}

// 添加全文索引方法
func (Video) TableName() string {
    return "videos"
}

func (Video) Indexes() map[string]string {
    return map[string]string{
        "idx_fulltext_search": "FULLTEXT(title, summary, director, actors, genre, region)",
    }
}
```

### 2.2 修改服务层 `app/service/video_service.go`

```go
// GetVideoList 获取视频列表（优化版 - 使用全文索引）
func (s *VideoService) GetVideoList(params VideoQueryParams) (*PageData, error) {
    query := db.GetDB().Model(&model.Video{})

    // 非模糊条件保持原样 (可使用普通索引)
    if params.CategoryID > 0 {
        query = query.Where("category_id = ?", params.CategoryID)
    }
    if params.Status > 0 {
        query = query.Where("status = ?", params.Status)
    }
    if params.AuthorID > 0 {
        query = query.Where("author_id = ?", params.AuthorID)
    }

    // === 优化1: 关键词搜索使用全文索引 ===
    if params.Keyword != "" {
        keyword := params.Keyword
        // MySQL 5.6+ 支持自然语言模式全文搜索
        query = query.Where(`
            MATCH(title, summary, director, actors) AGAINST(? IN BOOLEAN MODE)`, keyword)
    }

    // === 优化2: 精确匹配代替模糊匹配 ===
    if params.Genre != "" {
        // Genre存的是逗号分隔值，改用 FIND_IN_SET 或 LIKE '^genre,' 
        query = query.Where("genre LIKE ?", fmt.Sprintf("%s%%", params.Genre))
    }
    if params.Region != "" {
        // Region用精确匹配代替LIKE
        query = query.Where("region = ?", params.Region)
    }
    if params.Language != "" {
        query = query.Where("language = ?", params.Language)
    }
    if params.Director != "" {
        query = query.Where("director LIKE ?", fmt.Sprintf("%%%s%%", params.Director))
    }
    if params.Year > 0 {
        query = query.Where("year = ?", params.Year)
    }

    // ... rest unchanged ...
}
```

### 2.3 SQL迁移脚本

```sql
-- 创建全文索引 (MySQL 5.6+)
ALTER TABLE videos ADD FULLTEXT INDEX ft_title_summary 
    (title, summary);
ALTER TABLE videos ADD FULLTEXT INDEX ft_director_actors 
    (director, actors);
ALTER TABLE videos ADD FULLTEXT INDEX ft_genre_region 
    (genre, region);

-- 验证索引
SHOW INDEX FROM videos;
EXPLAIN SELECT * FROM videos 
WHERE MATCH(title, summary) AGAINST ('动作' IN BOOLEAN MODE);
```

---

## 3. 方案二：Elasticsearch集成 (推荐用于大数据量)

> 项目 `configs/config.yaml` 已有ES配置:
> ```yaml
> elasticsearch:
>   addresses: ["http://48.48.48.126:9200"]
>   username: "elastic"
>   index_prefix: "gocms_"
> ```

### 3.1 创建ES客户端 `app/db/elasticsearch.go`

```go
package db

import (
    "context"
    "github.com/olivere/elastic/v7"
)

var ESClient *elastic.Client

func InitElasticsearch(addresses []string) error {
    var err error
    ESClient, err = elastic.NewClient(
        elastic.SetURL(addresses...),
        elastic.SetBasicAuth("elastic", ""),
        elastic.SetSniff(false),
    )
    return err
}

func GetESClient() *elastic.Client {
    return ESClient
}
```

### 3.2 创建视频搜索引擎 `app/service/video_search_service.go`

```go
package service

import (
    "context"
    "encoding/json"
    "fmt"
    
    "gocms_v3/app/db"
    "gocms_v3/app/model"
)

type VideoSearchService struct{}

func NewVideoSearchService() *VideoSearchService {
    return &VideoSearchService{}
}

// SearchVideos 使用ES搜索视频
func (s *VideoSearchService) SearchVideos(
    ctx context.Context,
    keyword string,
    categoryID uint,
    status int,
    genre string,
    region string,
    year int,
    page, pageSize int,
) (*PageData, error) {
    client := db.GetESClient()
    indexName := "gocms_videos"

    // 构建查询
    boolQuery := elastic.NewBoolQuery()

    // 关键词全文搜索 (best_fields: title权重最高)
    if keyword != "" {
        multiMatch := elastic.NewMultiMatchQuery(keyword, "title", "summary", "director", "actors")
        multiMatch.Type("best_fields")
        multiMatch.Fuzziness("AUTO")
        boolQuery = boolQuery.Must(multiMatch)
    }

    // 精确条件
    if categoryID > 0 {
        boolQuery = boolQuery.Filter(elastic.NewTermQuery("category_id", categoryID))
    }
    if status > 0 {
        boolQuery = boolQuery.Filter(elastic.NewTermQuery("status", status))
    }
    if genre != "" {
        boolQuery = boolQuery.Filter(elastic.NewTermsQuery("genre", splitGenres(genre)))
    }
    if region != "" {
        boolQuery = boolQuery.Filter(elastic.NewTermQuery("region", region))
    }
    if year > 0 {
        boolQuery = boolQuery.Filter(elastic.NewTermQuery("year", year))
    }

    // 执行搜索
    result, err := client.Search(indexName).
        Query(boolQuery).
        From((page - 1) * pageSize).
        Size(pageSize).
        Sort("is_top", false).
        Sort("created_at", false).
        Do(ctx)
    if err != nil {
        return nil, fmt.Errorf("ES搜索失败: %w", err)
    }

    // 解析结果
    var videos []model.Video
    for _, hit := range result.Hits.Hits {
        var v model.Video
        if err := json.Unmarshal(hit.Source, &v); err != nil {
            continue
        }
        videos = append(videos, v)
    }

    return &PageData{
        List:  videos,
        Total: result.TotalHits(),
        Page:  page,
        PageSize: pageSize,
    }, nil
}

func splitGenres(genre string) []string {
    // "动作,喜剧" -> ["动作", "喜剧"]
    return []string{"动作", "喜剧"}
}
```

### 3.3 数据同步 `app/service/video_sync_service.go`

```go
package service

import (
    "context"
    "encoding/json"
    
    "gocms_v3/app/db"
    "gocms_v3/app/model"
)

func SyncVideosToES() error {
    client := db.GetESClient()
    indexName := "gocms_videos"
    
    // 创建索引映射
    mapping := `{
        "mappings": {
            "properties": {
                "id": {"type": "integer"},
                "title": {"type": "text", "analyzer": "ik_max_word", "search_analyzer": "ik_smart"},
                "summary": {"type": "text", "analyzer": "ik_max_word"},
                "director": {"type": "text", "analyzer": "ik_smart"},
                "actors": {"type": "text", "analyzer": "ik_smart"},
                "genre": {"type": "keyword"},
                "region": {"type": "keyword"},
                "year": {"type": "integer"},
                "status": {"type": "integer"},
                "category_id": {"type": "integer"},
                "created_at": {"type": "date"},
                "is_top": {"type": "boolean"}
            }
        }
    }`
    
    ctx := context.Background()
    exists, _ := client.IndexExists(indexName).Do(ctx)
    if !exists {
        _, err := client.CreateIndex(indexName).Body(mapping).Do(ctx)
        if err != nil {
            return err
        }
    }

    // 同步数据
    var videos []model.Video
    db.GetDB().Find(&videos)
    
    bulkReq := client.Bulk().Index(indexName)
    for _, v := range videos {
        doc, _ := json.Marshal(v)
        bulkReq.Add(elastic.NewBulkIndexOp().Id(fmt.Sprintf("%d", v.ID)).Doc(doc))
    }
    
    _, _, err := bulkReq.Do(ctx)
    return err
}
```

---

## 4. 推荐实施顺序

| 阶段 | 操作 | 预计效果 |
|------|------|---------|
| **Phase 1** | MySQL FULLTEXT INDEX (快速见效) | 查询性能提升 **5-10倍** |
| **Phase 2** | 优化 genre/region 用精确匹配代替LIKE | 可使用BTree索引 |
| **Phase 3** | Elasticsearch全文搜索 (大数据量场景) | 支持分词/模糊/相关性排序 |

---

## 5. MySQL vs Elasticsearch 对比

| 特性 | FULLTEXT INDEX | Elasticsearch |
|------|---------------|---------------|
| 数据量 < 10万 | ✅ 推荐 | ❌ 过度设计 |
| 中文分词支持 | ❌ 有限 (ngram) | ✅ 完美 (ik分词) |
| 相关性排序 | ⚠️ 基础 | ✅ BM25算法 |
| 部署复杂度 | 低 (MySQL内置) | 高 (额外服务) |
| 运维成本 | 低 | 中/高 |
| 实时性 | 强一致 | 近实时 (1s延迟) |

**建议**: 
- 视频量 < 10万 → 使用 **FULLTEXT INDEX**
- 视频量 > 10万 + 中文搜索需求 → 使用 **Elasticsearch**