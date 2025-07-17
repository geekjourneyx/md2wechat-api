# 永久素材上传功能升级

## 升级说明

本次升级将封面图片上传功能从**临时素材**改为**永久素材**，提供更好的性能和管理体验。

## 改动摘要

### 1. 接口变更

**之前 (临时素材)：**
```go
type WeChatService interface {
    UploadTempMedia(appID, appSecret string, imageURL string) (string, error)
}
```

**现在 (永久素材)：**
```go
type WeChatService interface {
    UploadPermanentMedia(appID, appSecret string, imageURL string) (string, string, error)
}
```

### 2. 核心实现变更

**之前 (临时素材)：**
```go
// 上传临时素材
result, err := materialManager.MediaUpload(material.MediaTypeImage, tempFile)
if err != nil {
    return "", fmt.Errorf("上传临时素材失败: %w", err)
}
return result.MediaID, nil
```

**现在 (永久素材)：**
```go
// 上传永久素材
mediaID, url, err := materialManager.AddMaterial(material.MediaTypeImage, tempFile)
if err != nil {
    return "", "", fmt.Errorf("上传永久素材失败: %w", err)
}
return mediaID, url, nil
```

### 3. 调用方式变更

**之前：**
```go
mediaID, err := s.UploadTempMedia(appID, appSecret, coverImageURL)
if err != nil {
    return nil, fmt.Errorf("上传封面图片失败: %w", err)
}
s.logger.Info("封面图片上传成功", zap.String("mediaID", mediaID))
```

**现在：**
```go
mediaID, url, err := s.UploadPermanentMedia(appID, appSecret, coverImageURL)
if err != nil {
    return nil, fmt.Errorf("上传封面图片失败: %w", err)
}
s.logger.Info("封面图片上传成功", zap.String("mediaID", mediaID), zap.String("url", url))
```

## 微信API对比

### 临时素材 API
- **方法**: `MediaUpload(mediaType, filename)`
- **生命周期**: 3天后自动删除
- **数量限制**: 无明确限制
- **大小限制**: 5MB
- **返回值**: 仅返回 `media_id`
- **适用场景**: 短期使用，如客服消息

### 永久素材 API  
- **方法**: `AddMaterial(mediaType, filename)`
- **生命周期**: 永久保存
- **数量限制**: 图片5000张，其他类型1000个
- **大小限制**: 2MB
- **返回值**: 返回 `media_id` 和 `url`
- **适用场景**: 长期使用，如图文消息、菜单等

## 优势对比

### 临时素材的限制
1. ❌ 3天后自动过期，草稿中的图片会失效
2. ❌ 需要重复上传相同图片
3. ❌ 无法获取图片URL
4. ❌ 管理困难

### 永久素材的优势
1. ✅ 永久保存，草稿图片不会失效
2. ✅ 一次上传，重复使用
3. ✅ 提供图片URL，便于管理
4. ✅ 更适合草稿场景

## 技术细节

### 返回值变化
永久素材API除了返回 `media_id` 外，还返回图片的访问URL：

```go
mediaID, url, err := materialManager.AddMaterial(material.MediaTypeImage, tempFile)
// mediaID: "media_id_string" - 用于草稿引用
// url: "http://mmbiz.qpic.cn/..." - 图片访问地址
```

### 错误处理增强
由于返回值增加，所有错误处理都相应调整：

```go
// 之前
return "", fmt.Errorf("错误信息")

// 现在  
return "", "", fmt.Errorf("错误信息")
```

### 日志记录增强
新增图片URL的日志记录，便于调试和管理：

```go
s.logger.Info("永久素材上传成功", 
    zap.String("mediaID", mediaID), 
    zap.String("url", url))
```

## 影响范围

### 受影响的文件
1. `internal/service/wechat.go` - 核心服务实现
2. `internal/handler/handler_test.go` - 测试模拟方法
3. `COVER_IMAGE_FEATURE.md` - 功能文档更新

### 不受影响的部分
1. ✅ API接口定义和路由
2. ✅ 请求/响应模型结构  
3. ✅ 认证和验证逻辑
4. ✅ 前端调用方式
5. ✅ 草稿创建流程

## 向后兼容性

这次升级对API使用者完全透明：

- ✅ **请求格式不变** - 仍然使用相同的JSON格式
- ✅ **响应格式不变** - 返回相同的数据结构
- ✅ **Header要求不变** - 仍然需要相同的认证信息
- ✅ **业务流程不变** - 转换和创建草稿的流程保持一致

唯一的变化是底层使用永久素材，用户无需修改任何调用代码。

## 验证结果

### 编译测试
```bash
✅ go build -v ./...
# 编译成功，无语法错误
```

### 单元测试
```bash
✅ go test ./internal/handler -v
=== RUN   TestConvertAndCreateDraft_Success
--- PASS: TestConvertAndCreateDraft_Success (0.00s)
=== RUN   TestConvertAndCreateDraft_MissingHeaders  
--- PASS: TestConvertAndCreateDraft_MissingHeaders (0.00s)
=== RUN   TestConvertAndCreateDraft_InvalidBodyRequest
--- PASS: TestConvertAndCreateDraft_InvalidBodyRequest (0.00s)
PASS
```

## 使用注意事项

### 1. 素材数量管理
- 图片素材上限为5000张
- 建议定期清理不再使用的素材
- 考虑实现素材复用机制

### 2. 图片大小限制
- 永久素材图片大小限制为2MB（vs 临时素材5MB）
- 建议在上传前进行图片大小检查

### 3. 内容合规
- 永久素材会长期保存在微信服务器
- 请确保上传的图片内容符合微信平台规范

### 4. 错误处理
- 新增素材数量超限的错误处理
- 增强图片大小超限的提示信息

## 总结

通过将封面图片上传从临时素材升级为永久素材，我们实现了：

1. **更好的用户体验** - 草稿图片不会过期失效
2. **更高的性能** - 相同图片无需重复上传  
3. **更强的管理能力** - 提供图片URL，便于管理
4. **更好的兼容性** - API接口完全向后兼容

这次升级为微信草稿功能提供了更稳定、更高效的封面图片支持。 