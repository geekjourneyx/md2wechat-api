package model

import "time"

// APIResponse 统一的 API 响应结构
type APIResponse struct {
	Code      int         `json:"code" example:"0"`
	Message   string      `json:"message" example:"success"`
	Data      interface{} `json:"data,omitempty"`
	Timestamp int64       `json:"timestamp" example:"1640995200"`
}

// ConvertAndDraftResponse 转换并创建草稿的响应结构
type ConvertAndDraftResponse struct {
	MediaID           string `json:"media_id" example:"media_id_123456"`
	HTML              string `json:"html" example:"<section>...</section>"`
	Theme             string `json:"theme" example:"default"`
	FontSize          string `json:"fontSize" example:"medium"`
	WordCount         int    `json:"wordCount" example:"156"`
	EstimatedReadTime int    `json:"estimatedReadTime" example:"1"`
}

// ErrorResponse 错误响应结构
type ErrorResponse struct {
	Code      int    `json:"code" example:"400"`
	Message   string `json:"message" example:"参数验证失败"`
	Details   string `json:"details,omitempty" example:"markdown 字段不能为空"`
	Timestamp int64  `json:"timestamp" example:"1640995200"`
}

// NewSuccessResponse 创建成功响应
func NewSuccessResponse(data interface{}) *APIResponse {
	return &APIResponse{
		Code:      0,
		Message:   "success",
		Data:      data,
		Timestamp: time.Now().Unix(),
	}
}

// NewErrorResponse 创建错误响应
func NewErrorResponse(code int, message string, details ...string) *ErrorResponse {
	resp := &ErrorResponse{
		Code:      code,
		Message:   message,
		Timestamp: time.Now().Unix(),
	}
	if len(details) > 0 {
		resp.Details = details[0]
	}
	return resp
} 