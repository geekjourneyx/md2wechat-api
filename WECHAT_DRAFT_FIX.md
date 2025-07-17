# 微信草稿 API 修复说明

## 问题描述

用户遇到了微信草稿创建失败的错误：
```
errcode=45106, errmsg=This API has been unsupported
```

这个错误表明我们之前使用的微信 API 已经被官方废弃了。

## 问题分析

### 原因
我们之前使用的是 `github.com/silenceper/wechat/v2/officialaccount/material` 包中的 `AddNews` 方法来创建草稿，但这个 API 已经被微信官方废弃。

### 解决方案
需要使用 `github.com/silenceper/wechat/v2/officialaccount/draft` 包中的 `AddDraft` 方法，这是微信官方推荐的新的草稿创建 API。

## 修改内容

### 1. 导入包修改
**修改前：**
```go
import (
    "github.com/silenceper/wechat/v2/officialaccount/material"
)
```

**修改后：**
```go
import (
    "github.com/silenceper/wechat/v2/officialaccount/draft"
)
```

### 2. 获取管理器修改
**修改前：**
```go
// 获取素材管理
materialManager := officialAccount.GetMaterial()
```

**修改后：**
```go
// 获取草稿管理器（使用 Draft 包而不是 Material 包）
draftManager := officialAccount.GetDraft()
```

### 3. 文章结构体类型修改
**修改前：**
```go
articles := make([]*material.Article, 0, len(req.Articles))
// ...
materialArticle := &material.Article{
    // ... 字段赋值
}
```

**修改后：**
```go
articles := make([]*draft.Article, 0, len(req.Articles))
// ...
draftArticle := &draft.Article{
    // ... 字段赋值
}
```

### 4. API 调用方法修改
**修改前：**
```go
// 创建草稿
mediaID, err := materialManager.AddNews(articles)
```

**修改后：**
```go
// 创建草稿（使用 AddDraft 方法而不是 AddNews）
mediaID, err := draftManager.AddDraft(articles)
```

### 5. 类型转换修复
由于 `draft.Article` 结构体中的 `ShowCoverPic` 字段类型是 `uint`，而我们的模型中是 `int`，需要进行类型转换：

**修改前：**
```go
ShowCoverPic: article.ShowCoverPic,
```

**修改后：**
```go
ShowCoverPic: uint(article.ShowCoverPic),
```

## 完整的修改后代码

```go
// CreateDraft 创建微信草稿
func (s *wechatService) CreateDraft(appID, appSecret string, req *model.WeChatDraftRequest) (*model.WeChatDraftResponse, error) {
    s.logger.Info("开始创建微信草稿",
        zap.String("appID", appID),
        zap.Int("articleCount", len(req.Articles)))

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

    // 准备草稿数据
    articles := make([]*draft.Article, 0, len(req.Articles))
    for i, article := range req.Articles {
        // ... 处理逻辑 ...

        draftArticle := &draft.Article{
            Title:            title,
            Content:          article.Content,
            Author:           article.Author,
            Digest:           digest,
            ContentSourceURL: article.ContentSourceURL,
            ThumbMediaID:     article.ThumbMediaID,
            ShowCoverPic:     uint(article.ShowCoverPic),
        }

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
```

## 验证结果

1. **编译测试**：✅ 通过
   ```bash
   go build -v ./...
   ```

2. **单元测试**：✅ 通过
   ```bash
   go test ./internal/handler -v
   === RUN   TestConvertAndCreateDraft_Success
   --- PASS: TestConvertAndCreateDraft_Success (0.00s)
   === RUN   TestConvertAndCreateDraft_MissingHeaders
   --- PASS: TestConvertAndCreateDraft_MissingHeaders (0.00s)
   === RUN   TestConvertAndCreateDraft_InvalidBodyRequest
   --- PASS: TestConvertAndCreateDraft_InvalidBodyRequest (0.00s)
   PASS
   ```

## API 接口对比

### 旧 API（已废弃）
```go
// Material 包 - AddNews 方法（errcode=45106）
func (material *Material) AddNews(articles []*Article) (mediaID string, err error)
```

### 新 API（推荐使用）
```go
// Draft 包 - AddDraft 方法
func (draft *Draft) AddDraft(articles []*Article) (mediaID string, err error)
```

## 影响范围

这次修改只影响了微信草稿创建的核心逻辑，不会影响：
- API 接口定义和路由
- 请求/响应模型结构
- 认证和验证逻辑
- 其他业务功能

## 注意事项

1. **向后兼容性**：新的 Draft API 与旧的 Material API 在接口签名上保持一致，因此不需要修改上层调用代码
2. **字段类型**：需要注意 `ShowCoverPic` 字段的类型转换（int → uint）
3. **错误处理**：新 API 的错误处理方式保持不变
4. **性能影响**：使用新 API 不会有性能上的负面影响

## 总结

通过将微信 SDK 从废弃的 Material 包切换到官方推荐的 Draft 包，成功解决了 `errcode=45106` 错误，确保了微信草稿创建功能的正常运行。修改后的代码已通过所有测试验证。 