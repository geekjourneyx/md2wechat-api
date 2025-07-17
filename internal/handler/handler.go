package handler

import (
	"net/http"
	"strings"

	"md2wechat-api/internal/model"
	"md2wechat-api/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

// Handler HTTP 处理器
type Handler struct {
	logger           *zap.Logger
	validator        *validator.Validate
	md2wechatService service.MD2WeChatService
	wechatService    service.WeChatService
}

// New 创建新的处理器
func New(logger *zap.Logger, md2wechatService service.MD2WeChatService, wechatService service.WeChatService) *Handler {
	return &Handler{
		logger:           logger,
		validator:        validator.New(),
		md2wechatService: md2wechatService,
		wechatService:    wechatService,
	}
}

// ConvertAndCreateDraft 转换 Markdown 并创建微信草稿
// @Summary 转换 Markdown 并创建微信草稿
// @Description 将 Markdown 内容转换为 HTML，然后创建微信公众号草稿。支持上传封面图片，传入图片URL会自动上传为临时素材并设置为草稿封面。
// @Tags 草稿管理
// @Accept json
// @Produce json
// @Param Wechat-Appid header string true "微信公众号AppID" example("wx1234567890abcdef")
// @Param Wechat-App-Secret header string true "微信公众号AppSecret" example("your_app_secret_here")
// @Param Md2wechat-API-Key header string false "MD2WeChat API密钥" example("wme_your_api_key_here")
// @Param request body model.ConvertAndDraftRequest true "请求参数，包含Markdown内容、主题、字体大小和可选的封面图片URL"
// @Success 200 {object} model.APIResponse{data=model.ConvertAndDraftResponse} "成功响应"
// @Failure 400 {object} model.ErrorResponse "参数错误"
// @Failure 401 {object} model.ErrorResponse "认证失败"
// @Failure 500 {object} model.ErrorResponse "服务器错误"
// @Router /convert-and-draft [post]
func (h *Handler) ConvertAndCreateDraft(c *gin.Context) {
	var req model.ConvertAndDraftRequest

	// 绑定请求数据
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("请求参数绑定失败", zap.Error(err))
		c.JSON(http.StatusBadRequest, model.NewErrorResponse(400, "请求参数格式错误", err.Error()))
		return
	}

	// 验证请求参数
	if err := h.validator.Struct(&req); err != nil {
		h.logger.Error("请求参数验证失败", zap.Error(err))
		c.JSON(http.StatusBadRequest, model.NewErrorResponse(400, "参数验证失败", err.Error()))
		return
	}

	// 从 Header 中获取微信配置
	wechatAppID := c.GetHeader("Wechat-Appid")
	wechatAppSecret := c.GetHeader("Wechat-App-Secret")
	md2wechatAPIKey := c.GetHeader("Md2wechat-API-Key")

	// 验证必填的 Header 参数
	if wechatAppID == "" {
		h.logger.Error("缺少微信AppID")
		c.JSON(http.StatusBadRequest, model.NewErrorResponse(400, "认证失败", "缺少 Wechat-Appid 请求头"))
		return
	}
	if wechatAppSecret == "" {
		h.logger.Error("缺少微信AppSecret")
		c.JSON(http.StatusBadRequest, model.NewErrorResponse(400, "认证失败", "缺少 Wechat-App-Secret 请求头"))
		return
	}

	h.logger.Info("收到转换并创建草稿请求",
		zap.String("appID", wechatAppID),
		zap.String("theme", req.Theme),
		zap.String("fontSize", req.FontSize),
		zap.String("coverImageURL", req.CoverImageURL),
		zap.Int("markdownLength", len(req.Markdown)))

	// 第一步：调用 MD2WeChat API 转换 Markdown
	md2wechatReq := &model.MD2WeChatRequest{
		Markdown: req.Markdown,
		Theme:    req.Theme,
		FontSize: req.FontSize,
	}

	convertData, err := h.md2wechatService.Convert(md2wechatReq, md2wechatAPIKey)
	if err != nil {
		h.logger.Error("MD2WeChat 转换失败", zap.Error(err))
		c.JSON(http.StatusInternalServerError, model.NewErrorResponse(500, "Markdown 转换失败", err.Error()))
		return
	}

	// 第二步：创建微信草稿
	// 从 HTML 中提取或生成标题
	title := extractTitleFromMarkdown(req.Markdown)
	if title == "" {
		title = "未命名文章"
	}

	wechatReq := &model.WeChatDraftRequest{
		Articles: []model.WeChatArticle{
			{
				Title:   title,
				Content: convertData.HTML,
				Author:  "", // 可以根据需要设置作者
			},
		},
	}

	draftResp, err := h.wechatService.CreateDraft(wechatAppID, wechatAppSecret, wechatReq, req.CoverImageURL)
	if err != nil {
		h.logger.Error("创建微信草稿失败", zap.Error(err))
		c.JSON(http.StatusInternalServerError, model.NewErrorResponse(500, "创建微信草稿失败", err.Error()))
		return
	}

	// 构建响应
	resp := &model.ConvertAndDraftResponse{
		MediaID:           draftResp.MediaID,
		HTML:              convertData.HTML,
		Theme:             convertData.Theme,
		FontSize:          convertData.FontSize,
		WordCount:         convertData.WordCount,
		EstimatedReadTime: convertData.EstimatedReadTime,
	}

	h.logger.Info("转换并创建草稿成功",
		zap.String("mediaID", draftResp.MediaID),
		zap.String("title", title),
		zap.Int("wordCount", convertData.WordCount))

	c.JSON(http.StatusOK, model.NewSuccessResponse(resp))
}

// extractTitleFromMarkdown 从 Markdown 中提取标题
func extractTitleFromMarkdown(markdown string) string {
	lines := strings.Split(markdown, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "# ") {
			title := strings.TrimPrefix(line, "# ")
			title = strings.TrimSpace(title)
			if len([]rune(title)) <= 64 { // 微信标题限制
				return title
			}
		}
	}
	return ""
} 