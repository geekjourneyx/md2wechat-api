package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"md2wechat-api/internal/model"
	"go.uber.org/zap"
)

// MD2WeChatService MD2WeChat 服务接口
type MD2WeChatService interface {
	Convert(req *model.MD2WeChatRequest, apiKey string) (*model.MD2WeChatData, error)
}

// md2wechatService MD2WeChat 服务实现
type md2wechatService struct {
	logger     *zap.Logger
	httpClient *http.Client
	baseURL    string
}

// NewMD2WeChatService 创建新的 MD2WeChat 服务
func NewMD2WeChatService(logger *zap.Logger) MD2WeChatService {
	return &md2wechatService{
		logger: logger,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		baseURL: "https://www.md2wechat.cn",
	}
}

// Convert 调用 MD2WeChat API 进行转换
func (s *md2wechatService) Convert(req *model.MD2WeChatRequest, apiKey string) (*model.MD2WeChatData, error) {
	// 记录请求日志
	s.logger.Info("开始调用 MD2WeChat API",
		zap.String("theme", req.Theme),
		zap.String("fontSize", req.FontSize),
		zap.Int("markdownLength", len(req.Markdown)))

	// 序列化请求数据
	jsonData, err := json.Marshal(req)
	if err != nil {
		s.logger.Error("序列化请求数据失败", zap.Error(err))
		return nil, fmt.Errorf("序列化请求数据失败: %w", err)
	}

	// 创建 HTTP 请求
	url := s.baseURL + "/api/convert"
	httpReq, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		s.logger.Error("创建 HTTP 请求失败", zap.Error(err))
		return nil, fmt.Errorf("创建 HTTP 请求失败: %w", err)
	}

	// 设置请求头
	httpReq.Header.Set("Content-Type", "application/json")
	if apiKey != "" {
		httpReq.Header.Set("X-API-Key", apiKey)
	}

	// 发送请求
	resp, err := s.httpClient.Do(httpReq)
	if err != nil {
		s.logger.Error("发送 HTTP 请求失败", zap.Error(err))
		return nil, fmt.Errorf("发送 HTTP 请求失败: %w", err)
	}
	defer resp.Body.Close()

	// 读取响应体
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		s.logger.Error("读取响应体失败", zap.Error(err))
		return nil, fmt.Errorf("读取响应体失败: %w", err)
	}

	// 检查 HTTP 状态码
	if resp.StatusCode != http.StatusOK {
		s.logger.Error("MD2WeChat API 返回错误状态码",
			zap.Int("statusCode", resp.StatusCode),
			zap.String("body", string(body)))
		return nil, fmt.Errorf("MD2WeChat API 返回错误状态码: %d", resp.StatusCode)
	}

	// 解析响应
	var apiResp model.MD2WeChatResponse
	if err := json.Unmarshal(body, &apiResp); err != nil {
		s.logger.Error("解析响应数据失败", zap.Error(err), zap.String("body", string(body)))
		return nil, fmt.Errorf("解析响应数据失败: %w", err)
	}

	// 检查业务状态码
	if apiResp.Code != 0 {
		s.logger.Error("MD2WeChat API 返回业务错误",
			zap.Int("code", apiResp.Code),
			zap.String("msg", apiResp.Msg))
		return nil, fmt.Errorf("MD2WeChat API 错误: %s", apiResp.Msg)
	}

	if apiResp.Data == nil {
		s.logger.Error("MD2WeChat API 返回空数据")
		return nil, fmt.Errorf("MD2WeChat API 返回空数据")
	}

	s.logger.Info("MD2WeChat API 调用成功",
		zap.String("theme", apiResp.Data.Theme),
		zap.String("fontSize", apiResp.Data.FontSize),
		zap.Int("wordCount", apiResp.Data.WordCount),
		zap.Int("estimatedReadTime", apiResp.Data.EstimatedReadTime))

	return apiResp.Data, nil
} 