package service

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"strings"

	"md2wechat-api/internal/model"

	"github.com/silenceper/wechat/v2"
	"github.com/silenceper/wechat/v2/cache"
	"github.com/silenceper/wechat/v2/officialaccount/config"
	"github.com/silenceper/wechat/v2/officialaccount/draft"
	"github.com/silenceper/wechat/v2/officialaccount/material"
	"go.uber.org/zap"
)

// WeChatService 微信服务接口
type WeChatService interface {
	CreateDraft(appID, appSecret string, req *model.WeChatDraftRequest, coverImageURL string) (*model.WeChatDraftResponse, error)
	UploadPermanentMedia(appID, appSecret string, imageURL string) (string, string, error)
}

// wechatService 微信服务实现
type wechatService struct {
	logger *zap.Logger
}

// NewWeChatService 创建新的微信服务
func NewWeChatService(logger *zap.Logger) WeChatService {
	return &wechatService{
		logger: logger,
	}
}

// CreateDraft 创建微信草稿
func (s *wechatService) CreateDraft(appID, appSecret string, req *model.WeChatDraftRequest, coverImageURL string) (*model.WeChatDraftResponse, error) {
	s.logger.Info("开始创建微信草稿",
		zap.String("appID", appID),
		zap.Int("articleCount", len(req.Articles)),
		zap.String("coverImageURL", coverImageURL))

	// 初始化微信 SDK
	wc := wechat.NewWechat()
	
	// 使用内存缓存
	memory := cache.NewMemory()

	// 创建公众号配置
	cfg := &config.Config{
		AppID:     appID,
		AppSecret: appSecret,
		Cache:     memory,
	}

	// 获取公众号实例
	officialAccount := wc.GetOfficialAccount(cfg)

	// 获取草稿管理器（使用 Draft 包而不是 Material 包）
	draftManager := officialAccount.GetDraft()

	// 上传封面图片获取 media_id
	var thumbMediaID string
	if coverImageURL != "" {
		mediaID, url, err := s.UploadPermanentMedia(appID, appSecret, coverImageURL)
		if err != nil {
			s.logger.Error("上传封面图片失败", zap.Error(err), zap.String("imageURL", coverImageURL))
			return nil, fmt.Errorf("上传封面图片失败: %w", err)
		}
		thumbMediaID = mediaID
		s.logger.Info("封面图片上传成功", zap.String("mediaID", thumbMediaID), zap.String("url", url))
	}

	// 准备草稿数据
	articles := make([]*draft.Article, 0, len(req.Articles))
	for i, article := range req.Articles {
		s.logger.Debug("处理文章",
			zap.Int("index", i),
			zap.String("title", article.Title),
			zap.Int("contentLength", len(article.Content)))

		// 提取文章标题（如果没有提供，从内容中提取）
		title := article.Title
		if title == "" {
			title = s.extractTitleFromHTML(article.Content)
		}

		// 生成文章摘要（如果没有提供）
		digest := article.Digest
		if digest == "" {
			digest = s.generateDigest(article.Content)
		}

		// 使用上传的封面图片 media_id，如果没有则使用文章自带的
		mediaID := thumbMediaID
		if mediaID == "" {
			mediaID = article.ThumbMediaID
		}

		draftArticle := &draft.Article{
			Title:            title,
			Content:          article.Content,
			Author:           article.Author,
			Digest:           digest,
			ContentSourceURL: article.ContentSourceURL,
			ThumbMediaID:     mediaID,
			ShowCoverPic:     uint(article.ShowCoverPic),
		}

		s.logger.Info("草稿文章", zap.Any("draftArticle", draftArticle))

		articles = append(articles, draftArticle)
	}

	// 创建草稿（使用 AddDraft 方法而不是 AddNews）
	mediaID, err := draftManager.AddDraft(articles)
	if err != nil {
		s.logger.Error("创建微信草稿失败", zap.Error(err))
		return nil, fmt.Errorf("创建微信草稿失败: %w", err)
	}

	s.logger.Info("微信草稿创建成功", zap.String("mediaID", mediaID))

	return &model.WeChatDraftResponse{
		ErrCode: 0,
		ErrMsg:  "ok",
		MediaID: mediaID,
	}, nil
}

// extractTitleFromHTML 从 HTML 内容中提取标题
func (s *wechatService) extractTitleFromHTML(html string) string {
	// 匹配 h1-h6 标签
	re := regexp.MustCompile(`<h[1-6][^>]*>(.*?)</h[1-6]>`)
	matches := re.FindStringSubmatch(html)
	if len(matches) > 1 {
		// 去除 HTML 标签
		title := regexp.MustCompile(`<[^>]*>`).ReplaceAllString(matches[1], "")
		return strings.TrimSpace(title)
	}

	// 如果没有找到标题标签，尝试提取第一行文本
	lines := strings.Split(strings.TrimSpace(html), "\n")
	for _, line := range lines {
		text := regexp.MustCompile(`<[^>]*>`).ReplaceAllString(line, "")
		text = strings.TrimSpace(text)
		if text != "" && len(text) <= 64 { // 微信标题限制 64 字符
			return text
		}
	}

	return "未命名文章"
}

// generateDigest 生成文章摘要
func (s *wechatService) generateDigest(html string) string {
	// 去除 HTML 标签
	text := regexp.MustCompile(`<[^>]*>`).ReplaceAllString(html, "")
	
	// 去除多余空白字符
	text = regexp.MustCompile(`\s+`).ReplaceAllString(text, " ")
	text = strings.TrimSpace(text)

	// 限制摘要长度（微信摘要限制 120 字符）
	if len([]rune(text)) > 100 {
		runes := []rune(text)
		text = string(runes[:100]) + "..."
	}

	return text
}

// UploadPermanentMedia 上传永久素材
func (s *wechatService) UploadPermanentMedia(appID, appSecret string, imageURL string) (string, string, error) {
	s.logger.Info("开始上传永久素材", zap.String("imageURL", imageURL))

	// 下载图片
	resp, err := http.Get(imageURL)
	if err != nil {
		s.logger.Error("下载图片失败", zap.Error(err))
		return "", "", fmt.Errorf("下载图片失败: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		s.logger.Error("下载图片失败，状态码不正确", zap.Int("statusCode", resp.StatusCode))
		return "", "", fmt.Errorf("下载图片失败，状态码: %d", resp.StatusCode)
	}

	// 读取图片数据
	imageData, err := io.ReadAll(resp.Body)
	if err != nil {
		s.logger.Error("读取图片数据失败", zap.Error(err))
		return "", "", fmt.Errorf("读取图片数据失败: %w", err)
	}

	// 初始化微信 SDK
	wc := wechat.NewWechat()
	
	// 使用内存缓存
	memory := cache.NewMemory()

	// 创建公众号配置
	cfg := &config.Config{
		AppID:     appID,
		AppSecret: appSecret,
		Cache:     memory,
	}

	// 获取公众号实例
	officialAccount := wc.GetOfficialAccount(cfg)

	// 获取素材管理器
	materialManager := officialAccount.GetMaterial()

	// 上传永久素材 - 需要先保存为临时文件
	tempFile := "temp_cover.jpg"
	err = os.WriteFile(tempFile, imageData, 0644)
	if err != nil {
		s.logger.Error("保存临时文件失败", zap.Error(err))
		return "", "", fmt.Errorf("保存临时文件失败: %w", err)
	}
	defer os.Remove(tempFile) // 清理临时文件

	// 上传永久素材
	mediaID, url, err := materialManager.AddMaterial(material.MediaTypeImage, tempFile)
	if err != nil {
		s.logger.Error("上传永久素材失败", zap.Error(err))
		return "", "", fmt.Errorf("上传永久素材失败: %w", err)
	}

	s.logger.Info("永久素材上传成功", zap.String("mediaID", mediaID), zap.String("url", url))
	return mediaID, url, nil
} 