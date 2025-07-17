package model

// ConvertAndDraftRequest 转换并创建草稿的请求结构
type ConvertAndDraftRequest struct {
	Markdown       string `json:"markdown" validate:"required" example:"# 标题\n\n这是一个**加粗**文本的例子。"`
	Theme          string `json:"theme" validate:"required" example:"default"`
	FontSize       string `json:"fontSize" validate:"required" example:"medium"`
	CoverImageURL  string `json:"coverImageUrl,omitempty" example:"https://example.com/cover.jpg"`
}

// MD2WeChatRequest MD2WeChat API 请求结构
type MD2WeChatRequest struct {
	Markdown string `json:"markdown"`
	Theme    string `json:"theme"`
	FontSize string `json:"fontSize"`
}

// MD2WeChatResponse MD2WeChat API 响应结构
type MD2WeChatResponse struct {
	Code int                 `json:"code"`
	Msg  string              `json:"msg"`
	Data *MD2WeChatData      `json:"data,omitempty"`
}

// MD2WeChatData MD2WeChat 响应数据
type MD2WeChatData struct {
	HTML               string `json:"html"`
	Theme              string `json:"theme"`
	FontSize           string `json:"fontSize"`
	WordCount          int    `json:"wordCount"`
	EstimatedReadTime  int    `json:"estimatedReadTime"`
}

// WeChatDraftRequest 微信草稿请求结构
type WeChatDraftRequest struct {
	Articles []WeChatArticle `json:"articles"`
}

// WeChatArticle 微信文章结构
type WeChatArticle struct {
	Title            string `json:"title"`
	Author           string `json:"author,omitempty"`
	Digest           string `json:"digest,omitempty"`
	Content          string `json:"content"`
	ContentSourceURL string `json:"content_source_url,omitempty"`
	ThumbMediaID     string `json:"thumb_media_id,omitempty"`
	ShowCoverPic     int    `json:"show_cover_pic,omitempty"`
}

// WeChatDraftResponse 微信草稿响应结构
type WeChatDraftResponse struct {
	ErrCode int    `json:"errcode"`
	ErrMsg  string `json:"errmsg"`
	MediaID string `json:"media_id,omitempty"`
} 