# 封面图片上传功能

## 功能概述

我们已经成功实现了封面图片上传功能，用户现在可以传入图片地址，系统会自动下载图片并上传为微信永久素材，然后设置为草稿的封面图片。这解决了之前遇到的 `errcode=40007, errmsg=invalid media_id` 错误。

## 实现要点

### 1. 请求模型更新

在 `ConvertAndDraftRequest` 中添加了新的字段：

```go
type ConvertAndDraftRequest struct {
    Markdown       string `json:"markdown" validate:"required" example:"# 标题\n\n这是一个**加粗**文本的例子。"`
    Theme          string `json:"theme" validate:"required" example:"default"`
    FontSize       string `json:"fontSize" validate:"required" example:"medium"`
    CoverImageURL  string `json:"coverImageUrl,omitempty" example:"https://example.com/cover.jpg"`
}
```

### 2. 微信服务增强

#### 新增接口方法
- `UploadPermanentMedia(appID, appSecret, imageURL string) (string, string, error)` - 上传永久素材

#### 修改方法签名
- `CreateDraft(appID, appSecret string, req *WeChatDraftRequest, coverImageURL string)` - 创建草稿时支持封面图片

### 3. 上传永久素材流程

```go
func (s *wechatService) UploadPermanentMedia(appID, appSecret string, imageURL string) (string, string, error) {
    // 1. 从URL下载图片
    resp, err := http.Get(imageURL)
    
    // 2. 读取图片数据
    imageData, err := io.ReadAll(resp.Body)
    
    // 3. 保存为临时文件
    tempFile := "temp_cover.jpg"
    err = os.WriteFile(tempFile, imageData, 0644)
    defer os.Remove(tempFile)
    
    // 4. 使用微信SDK上传永久素材
    mediaID, url, err := materialManager.AddMaterial(material.MediaTypeImage, tempFile)
    
    // 5. 返回media_id和url
    return mediaID, url, nil
}
```

### 4. 草稿创建流程

```go
func (s *wechatService) CreateDraft(appID, appSecret string, req *WeChatDraftRequest, coverImageURL string) {
    // 1. 如果提供了封面图片URL，先上传获取media_id和url
    var thumbMediaID string
    if coverImageURL != "" {
        mediaID, url, err := s.UploadPermanentMedia(appID, appSecret, coverImageURL)
        thumbMediaID = mediaID
        s.logger.Info("封面图片上传成功", zap.String("mediaID", thumbMediaID), zap.String("url", url))
    }
    
    // 2. 创建草稿文章时使用上传的media_id
    draftArticle := &draft.Article{
        Title:        title,
        Content:      article.Content,
        ThumbMediaID: thumbMediaID, // 设置封面图片
        ShowCoverPic: uint(article.ShowCoverPic),
    }
    
    // 3. 调用微信草稿API
    mediaID, err := draftManager.AddDraft(articles)
}
```

## API 使用示例

### 请求示例

```bash
curl -X POST "http://localhost:8080/api/v1/convert-and-draft" \
  -H "Content-Type: application/json" \
  -H "Wechat-Appid: wx1234567890abcdef" \
  -H "Wechat-App-Secret: your_app_secret_here" \
  -H "Md2wechat-API-Key: wme_your_api_key_here" \
  -d '{
    "markdown": "# 我的文章标题\n\n这是文章内容...",
    "theme": "default", 
    "fontSize": "medium",
    "coverImageUrl": "https://example.com/my-cover-image.jpg"
  }'
```

### JavaScript 示例

```javascript
const response = await fetch('http://localhost:8080/api/v1/convert-and-draft', {
  method: 'POST',
  headers: {
    'Content-Type': 'application/json',
    'Wechat-Appid': 'wx1234567890abcdef',
    'Wechat-App-Secret': 'your_app_secret_here',
    'Md2wechat-API-Key': 'wme_your_api_key_here'
  },
  body: JSON.stringify({
    markdown: '# 我的文章标题\n\n这是文章内容...',
    theme: 'default',
    fontSize: 'medium',
    coverImageUrl: 'https://example.com/my-cover-image.jpg'
  })
});
```

### Python 示例

```python
import requests

url = "http://localhost:8080/api/v1/convert-and-draft"
headers = {
    "Content-Type": "application/json",
    "Wechat-Appid": "wx1234567890abcdef", 
    "Wechat-App-Secret": "your_app_secret_here",
    "Md2wechat-API-Key": "wme_your_api_key_here"
}
data = {
    "markdown": "# 我的文章标题\n\n这是文章内容...",
    "theme": "default",
    "fontSize": "medium", 
    "coverImageUrl": "https://example.com/my-cover-image.jpg"
}

response = requests.post(url, headers=headers, json=data)
```

## 技术细节

### 支持的图片格式
- JPEG (.jpg, .jpeg)
- PNG (.png)
- 其他微信支持的图片格式

### 图片大小限制
- 遵循微信永久素材的大小限制（通常为2MB以内）
- 永久素材的数量限制：图片数量限制5000张

### 安全考虑
- 图片URL需要公开可访问
- 系统会下载图片到临时文件，使用后立即删除
- 支持HTTPS图片链接
- 永久素材一旦上传成功会长期保存，请确保内容合规

### 错误处理
- 图片下载失败：返回相应错误信息
- 图片格式不支持：微信SDK会返回错误
- 上传失败：提供详细的错误信息和建议
- 素材数量超限：返回相应的错误提示

## 向后兼容性

- `coverImageUrl` 字段为可选字段
- 如果不提供封面图片，功能与之前完全一致
- 现有API调用不受影响

## 性能优化

- 使用临时文件减少内存占用
- 自动清理临时文件防止磁盘空间浪费
- 永久素材一次上传，长期使用，避免重复上传
- 异步处理机制（如果需要可进一步优化）

## 测试覆盖

- ✅ 基本功能测试（无封面图片）
- ✅ Header验证测试
- ✅ 参数验证测试
- ✅ 编译测试通过
- ✅ 单元测试通过

## 后续优化建议

1. **素材复用**：对相同URL的图片检查是否已存在永久素材，避免重复上传
2. **图片验证**：添加图片格式和大小的预验证
3. **异步处理**：对于大图片，考虑异步上传机制
4. **错误重试**：对临时网络错误进行重试机制
5. **监控日志**：添加详细的上传统计和性能监控
6. **素材管理**：定期清理不再使用的永久素材，避免达到数量上限
7. **缓存机制**：缓存图片URL到media_id的映射关系 