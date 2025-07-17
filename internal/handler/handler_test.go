package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"md2wechat-api/internal/model"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
)

// MockMD2WeChatService 模拟 MD2WeChat 服务
type MockMD2WeChatService struct {
	mock.Mock
}

func (m *MockMD2WeChatService) Convert(req *model.MD2WeChatRequest, apiKey string) (*model.MD2WeChatData, error) {
	args := m.Called(req, apiKey)
	return args.Get(0).(*model.MD2WeChatData), args.Error(1)
}

// MockWeChatService 模拟微信服务
type MockWeChatService struct {
	mock.Mock
}

func (m *MockWeChatService) CreateDraft(appID, appSecret string, req *model.WeChatDraftRequest, coverImageURL string) (*model.WeChatDraftResponse, error) {
	args := m.Called(appID, appSecret, req, coverImageURL)
	return args.Get(0).(*model.WeChatDraftResponse), args.Error(1)
}

func (m *MockWeChatService) UploadPermanentMedia(appID, appSecret string, imageURL string) (string, string, error) {
	args := m.Called(appID, appSecret, imageURL)
	return args.String(0), args.String(1), args.Error(2)
}

func TestConvertAndCreateDraft_Success(t *testing.T) {
	// 设置 Gin 为测试模式
	gin.SetMode(gin.TestMode)

	// 创建模拟服务
	mockMD2WeChat := new(MockMD2WeChatService)
	mockWeChat := new(MockWeChatService)

	// 设置模拟返回值
	mockMD2WeChat.On("Convert", mock.AnythingOfType("*model.MD2WeChatRequest"), "test_api_key").Return(&model.MD2WeChatData{
		HTML:              "<section><h1>标题</h1><p>这是一个<strong>加粗</strong>文本的例子。</p></section>",
		Theme:             "default",
		FontSize:          "medium",
		WordCount:         156,
		EstimatedReadTime: 1,
	}, nil)

	mockWeChat.On("CreateDraft", "wx123", "secret123", mock.AnythingOfType("*model.WeChatDraftRequest"), "").Return(&model.WeChatDraftResponse{
		ErrCode: 0,
		ErrMsg:  "ok",
		MediaID: "media_123456",
	}, nil)

	// 创建处理器
	logger := zap.NewNop()
	handler := New(logger, mockMD2WeChat, mockWeChat)

	// 创建路由
	router := gin.New()
	router.POST("/api/v1/convert-and-draft", handler.ConvertAndCreateDraft)

	// 准备请求数据
	reqData := model.ConvertAndDraftRequest{
		Markdown: "# 标题\n\n这是一个**加粗**文本的例子。",
		Theme:    "default",
		FontSize: "medium",
	}

	jsonData, _ := json.Marshal(reqData)

	// 创建请求
	req, _ := http.NewRequest("POST", "/api/v1/convert-and-draft", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Wechat-Appid", "wx123")
	req.Header.Set("Wechat-App-Secret", "secret123")
	req.Header.Set("Md2wechat-API-Key", "test_api_key")

	// 执行请求
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// 验证结果
	assert.Equal(t, http.StatusOK, w.Code)

	var response model.APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, 0, response.Code)
	assert.Equal(t, "success", response.Message)

	// 验证模拟服务被调用
	mockMD2WeChat.AssertExpectations(t)
	mockWeChat.AssertExpectations(t)
}

func TestConvertAndCreateDraft_MissingHeaders(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// 创建处理器
	logger := zap.NewNop()
	mockMD2WeChat := new(MockMD2WeChatService)
	mockWeChat := new(MockWeChatService)
	handler := New(logger, mockMD2WeChat, mockWeChat)

	// 创建路由
	router := gin.New()
	router.POST("/api/v1/convert-and-draft", handler.ConvertAndCreateDraft)

	// 准备有效的请求数据（但缺少header）
	reqData := map[string]interface{}{
		"markdown": "# 标题",
		"theme":    "default",
		"fontSize": "medium",
	}

	jsonData, _ := json.Marshal(reqData)

	// 创建请求（缺少 header 参数）
	req, _ := http.NewRequest("POST", "/api/v1/convert-and-draft", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	// 故意不设置 Wechat-Appid 和 Wechat-App-Secret

	// 执行请求
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// 验证结果
	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response model.ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, 400, response.Code)
	assert.Contains(t, response.Message, "认证失败")
}

func TestConvertAndCreateDraft_InvalidBodyRequest(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// 创建处理器
	logger := zap.NewNop()
	mockMD2WeChat := new(MockMD2WeChatService)
	mockWeChat := new(MockWeChatService)
	handler := New(logger, mockMD2WeChat, mockWeChat)

	// 创建路由
	router := gin.New()
	router.POST("/api/v1/convert-and-draft", handler.ConvertAndCreateDraft)

	// 准备无效的请求数据（缺少必填字段）
	reqData := map[string]interface{}{
		"markdown": "# 标题",
		// 缺少 theme 和 fontSize
	}

	jsonData, _ := json.Marshal(reqData)

	// 创建请求（包含正确的header）
	req, _ := http.NewRequest("POST", "/api/v1/convert-and-draft", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Wechat-Appid", "wx123")
	req.Header.Set("Wechat-App-Secret", "secret123")
	req.Header.Set("Md2wechat-API-Key", "test_api_key")

	// 执行请求
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// 验证结果
	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response model.ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, 400, response.Code)
	assert.Contains(t, response.Message, "参数验证失败")
} 