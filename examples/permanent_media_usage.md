# 永久素材上传使用示例

本文档展示如何使用升级后的永久素材上传功能。

## 基本使用

### 1. 包含封面图片的草稿创建

```bash
curl -X POST "http://localhost:8080/api/v1/convert-and-draft" \
  -H "Content-Type: application/json" \
  -H "Wechat-Appid: wx1234567890abcdef" \
  -H "Wechat-App-Secret: your_app_secret_here" \
  -H "Md2wechat-API-Key: wme_your_api_key_here" \
  -d '{
    "markdown": "# 使用永久素材的文章\n\n这篇文章会使用永久素材作为封面图片，确保图片长期有效。\n\n## 主要优势\n\n- 图片永久保存\n- 避免重复上传\n- 提供图片URL\n- 更好的管理体验",
    "theme": "modern",
    "fontSize": "medium",
    "coverImageUrl": "https://example.com/permanent-cover.jpg"
  }'
```

### 2. 成功响应示例

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "media_id": "media_ABC123456789",
    "html": "<section style=\"...\">\n  <h1>使用永久素材的文章</h1>\n  <p>这篇文章会使用永久素材作为封面图片...</p>\n  ...\n</section>",
    "theme": "modern",
    "fontSize": "medium",
    "wordCount": 89,
    "estimatedReadTime": 1
  },
  "timestamp": 1640995200
}
```

## 编程语言示例

### JavaScript/Node.js

```javascript
const axios = require('axios');

async function createDraftWithPermanentMedia() {
  try {
    const response = await axios.post('http://localhost:8080/api/v1/convert-and-draft', {
      markdown: `# 技术分享：微信开发最佳实践

## 永久素材的优势

使用永久素材能够：

1. **长期保存** - 不会因为时间过期而失效
2. **减少重复** - 相同图片只需上传一次
3. **便于管理** - 提供图片URL，便于查看和管理
4. **性能优化** - 减少不必要的上传操作

> 推荐在生产环境中使用永久素材功能

\`\`\`javascript
// 示例代码
console.log("使用永久素材创建草稿");
\`\`\`

这样的文章内容会被转换为美观的HTML，并配上永久有效的封面图片。`,
      theme: 'default',
      fontSize: 'medium',
      coverImageUrl: 'https://cdn.example.com/tech-sharing-cover.png'
    }, {
      headers: {
        'Content-Type': 'application/json',
        'Wechat-Appid': 'wx1234567890abcdef',
        'Wechat-App-Secret': 'your_app_secret_here',
        'Md2wechat-API-Key': 'wme_your_api_key_here'
      }
    });

    console.log('✅ 草稿创建成功！');
    console.log('📄 草稿ID:', response.data.data.media_id);
    console.log('📝 字数统计:', response.data.data.wordCount);
    console.log('⏱️ 预计阅读时间:', response.data.data.estimatedReadTime, '分钟');
    
  } catch (error) {
    console.error('❌ 请求失败:', error.response?.data || error.message);
    
    // 处理常见错误
    if (error.response?.status === 400) {
      console.log('💡 提示：请检查请求参数和认证信息');
    } else if (error.response?.status === 500) {
      console.log('💡 提示：可能是图片下载失败或微信API调用失败');
    }
  }
}

createDraftWithPermanentMedia();
```

### Python

```python
import requests
import json

def create_draft_with_permanent_media():
    """使用永久素材创建微信草稿"""
    
    url = "http://localhost:8080/api/v1/convert-and-draft"
    
    headers = {
        "Content-Type": "application/json",
        "Wechat-Appid": "wx1234567890abcdef",
        "Wechat-App-Secret": "your_app_secret_here",
        "Md2wechat-API-Key": "wme_your_api_key_here"
    }
    
    # 准备Markdown内容
    markdown_content = """# 产品发布公告

## 新功能上线

我们很高兴地宣布，新版本已经正式发布！

### 主要更新

- ✨ 全新的用户界面设计
- 🚀 性能提升30%
- 🔒 增强的安全性
- 📱 更好的移动端体验

### 使用说明

1. **登录账户** - 使用您的现有凭据
2. **体验新功能** - 在主菜单中查找新增选项
3. **反馈建议** - 通过客服渠道告诉我们您的想法

> 💡 **小贴士**：您可以在设置中切换到经典模式

感谢您的支持！

---

*产品团队*  
*2024年1月*"""
    
    data = {
        "markdown": markdown_content,
        "theme": "modern",
        "fontSize": "large",
        "coverImageUrl": "https://cdn.example.com/product-announcement.jpg"
    }
    
    try:
        response = requests.post(url, json=data, headers=headers, timeout=30)
        response.raise_for_status()
        
        result = response.json()
        
        print("✅ 草稿创建成功！")
        print(f"📄 草稿ID: {result['data']['media_id']}")
        print(f"📝 字数统计: {result['data']['wordCount']}")
        print(f"⏱️ 预计阅读时间: {result['data']['estimatedReadTime']} 分钟")
        print(f"🎨 使用主题: {result['data']['theme']}")
        print(f"📏 字体大小: {result['data']['fontSize']}")
        
        return result['data']['media_id']
        
    except requests.exceptions.RequestException as e:
        print(f"❌ 请求失败: {e}")
        
        if hasattr(e, 'response') and e.response is not None:
            try:
                error_data = e.response.json()
                print(f"💡 错误详情: {error_data.get('message', '未知错误')}")
                if 'details' in error_data:
                    print(f"📋 详细信息: {error_data['details']}")
            except:
                print(f"💡 HTTP状态码: {e.response.status_code}")
                
        return None

if __name__ == "__main__":
    media_id = create_draft_with_permanent_media()
    if media_id:
        print(f"\n🎉 请登录微信公众号后台查看草稿：{media_id}")
    else:
        print("\n😔 草稿创建失败，请检查配置和网络连接")
```

### Go 语言

```go
package main

import (
    "bytes"
    "encoding/json"
    "fmt"
    "io"
    "net/http"
    "time"
)

type ConvertAndDraftRequest struct {
    Markdown      string `json:"markdown"`
    Theme         string `json:"theme"`
    FontSize      string `json:"fontSize"`
    CoverImageURL string `json:"coverImageUrl"`
}

type APIResponse struct {
    Code      int                 `json:"code"`
    Message   string              `json:"message"`
    Data      *ConvertAndDraftResponse `json:"data"`
    Timestamp int64               `json:"timestamp"`
}

type ConvertAndDraftResponse struct {
    MediaID           string `json:"media_id"`
    HTML              string `json:"html"`
    Theme             string `json:"theme"`
    FontSize          string `json:"fontSize"`
    WordCount         int    `json:"wordCount"`
    EstimatedReadTime int    `json:"estimatedReadTime"`
}

func createDraftWithPermanentMedia() error {
    // 准备请求数据
    req := ConvertAndDraftRequest{
        Markdown: `# Go语言实战指南

## 为什么选择Go

Go语言凭借其简洁的语法和强大的性能，正在成为云原生时代的首选语言。

### 核心优势

1. **高性能** - 编译型语言，接近C++的执行效率
2. **并发友好** - 内置goroutine，轻松处理并发
3. **简单易学** - 语法简洁，上手容易
4. **生态丰富** - 丰富的标准库和第三方包

### 适用场景

- 🌐 Web服务开发
- 🔄 微服务架构
- 🏗️ 基础设施工具
- 📊 数据处理

` + "```go\n" + `package main

import "fmt"

func main() {
    fmt.Println("Hello, Go!")
    
    // 并发示例
    go func() {
        fmt.Println("这是一个goroutine")
    }()
}
` + "```" + `

> 💡 **建议**：从小项目开始，逐步掌握Go的精髓

开始您的Go语言之旅吧！`,
        Theme:         "modern",
        FontSize:      "medium",
        CoverImageURL: "https://cdn.example.com/golang-guide.png",
    }

    // 序列化请求数据
    jsonData, err := json.Marshal(req)
    if err != nil {
        return fmt.Errorf("序列化请求失败: %w", err)
    }

    // 创建HTTP请求
    httpReq, err := http.NewRequest("POST", "http://localhost:8080/api/v1/convert-and-draft", bytes.NewBuffer(jsonData))
    if err != nil {
        return fmt.Errorf("创建请求失败: %w", err)
    }

    // 设置请求头
    httpReq.Header.Set("Content-Type", "application/json")
    httpReq.Header.Set("Wechat-Appid", "wx1234567890abcdef")
    httpReq.Header.Set("Wechat-App-Secret", "your_app_secret_here")
    httpReq.Header.Set("Md2wechat-API-Key", "wme_your_api_key_here")

    // 发送请求
    client := &http.Client{Timeout: 30 * time.Second}
    resp, err := client.Do(httpReq)
    if err != nil {
        return fmt.Errorf("发送请求失败: %w", err)
    }
    defer resp.Body.Close()

    // 读取响应
    body, err := io.ReadAll(resp.Body)
    if err != nil {
        return fmt.Errorf("读取响应失败: %w", err)
    }

    // 检查状态码
    if resp.StatusCode != http.StatusOK {
        return fmt.Errorf("HTTP错误 %d: %s", resp.StatusCode, string(body))
    }

    // 解析响应
    var apiResp APIResponse
    if err := json.Unmarshal(body, &apiResp); err != nil {
        return fmt.Errorf("解析响应失败: %w", err)
    }

    // 输出结果
    fmt.Println("✅ 草稿创建成功！")
    fmt.Printf("📄 草稿ID: %s\n", apiResp.Data.MediaID)
    fmt.Printf("📝 字数统计: %d\n", apiResp.Data.WordCount)
    fmt.Printf("⏱️ 预计阅读时间: %d 分钟\n", apiResp.Data.EstimatedReadTime)
    fmt.Printf("🎨 使用主题: %s\n", apiResp.Data.Theme)
    fmt.Printf("📏 字体大小: %s\n", apiResp.Data.FontSize)

    return nil
}

func main() {
    if err := createDraftWithPermanentMedia(); err != nil {
        fmt.Printf("❌ 操作失败: %v\n", err)
    }
}
```

## 最佳实践

### 1. 图片选择建议

```bash
# 推荐的图片规格
- 分辨率：900x500 像素（适合微信文章封面）
- 格式：JPEG 或 PNG
- 大小：≤ 2MB（永久素材限制）
- 比例：16:9 或 3:2（视觉效果较好）
```

### 2. 错误处理示例

```javascript
// 完整的错误处理示例
async function handleDraftCreation() {
  try {
    const response = await createDraft();
    console.log('成功创建草稿');
  } catch (error) {
    // 根据错误类型给出不同的处理建议
    switch (error.response?.status) {
      case 400:
        console.error('参数错误：请检查Markdown内容和图片URL');
        break;
      case 401:
        console.error('认证失败：请检查AppID和AppSecret');
        break;
      case 413:
        console.error('图片过大：请使用小于2MB的图片');
        break;
      case 429:
        console.error('请求频率过高：请稍后重试');
        break;
      case 500:
        console.error('服务器错误：可能是图片下载失败或微信API异常');
        break;
      default:
        console.error('未知错误：', error.message);
    }
  }
}
```

### 3. 批量处理示例

```python
import asyncio
import aiohttp

async def batch_create_drafts(articles):
    """批量创建草稿"""
    async with aiohttp.ClientSession() as session:
        tasks = []
        for article in articles:
            task = create_single_draft(session, article)
            tasks.append(task)
        
        # 并发处理，但限制并发数
        semaphore = asyncio.Semaphore(3)  # 最多3个并发请求
        
        async def limited_create(article):
            async with semaphore:
                return await create_single_draft(session, article)
        
        results = await asyncio.gather(*[limited_create(article) for article in articles])
        return results

# 使用示例
articles = [
    {
        "title": "文章1",
        "content": "# 标题1\n内容1...",
        "cover": "https://example.com/cover1.jpg"
    },
    {
        "title": "文章2", 
        "content": "# 标题2\n内容2...",
        "cover": "https://example.com/cover2.jpg"
    }
]

results = asyncio.run(batch_create_drafts(articles))
```

## 性能优化建议

### 1. 图片复用检查

在实际应用中，建议实现图片复用机制：

```python
# 伪代码：检查图片是否已上传
def get_or_upload_media(image_url):
    # 1. 检查本地缓存
    cached_media_id = cache.get(f"media:{hash(image_url)}")
    if cached_media_id:
        return cached_media_id
    
    # 2. 调用上传接口
    media_id = upload_permanent_media(image_url)
    
    # 3. 缓存结果
    cache.set(f"media:{hash(image_url)}", media_id, expire=86400)
    
    return media_id
```

### 2. 异步处理

对于大批量处理，建议使用异步模式：

```javascript
// 使用队列处理大量请求
const Queue = require('bull');
const draftQueue = new Queue('draft creation');

draftQueue.process(async (job) => {
  const { markdown, theme, fontSize, coverImageUrl } = job.data;
  return await createDraft(markdown, theme, fontSize, coverImageUrl);
});

// 添加任务到队列
draftQueue.add('create', {
  markdown: content,
  theme: 'modern',
  fontSize: 'medium',
  coverImageUrl: 'https://example.com/cover.jpg'
});
```

通过这些示例和最佳实践，您可以充分利用永久素材上传功能，创建高质量的微信公众号草稿。 