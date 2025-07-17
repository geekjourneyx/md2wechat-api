# API 使用示例

本文档提供了如何使用 MD2WeChat API 的详细示例。

## 基本使用

### 转换 Markdown 并创建微信草稿

```bash
curl -X POST "http://localhost:8080/api/v1/convert-and-draft" \
  -H "Content-Type: application/json" \
  -H "Wechat-Appid: wx1234567890abcdef" \
  -H "Wechat-App-Secret: your_app_secret_here" \
  -H "Md2wechat-API-Key: wme_your_api_key_here" \
  -d '{
    "markdown": "# 我的第一篇文章\n\n这是一个**重要**的公告。\n\n## 主要内容\n\n- 功能特性\n- 使用说明\n- 注意事项\n\n> 这是一个引用块\n\n```javascript\nconsole.log(\"Hello World!\");\n```",
    "theme": "default",
    "fontSize": "medium"
  }'
```

### 成功响应示例

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "media_id": "media_ABC123456789",
    "html": "<section style=\"...\">\n  <h1>我的第一篇文章</h1>\n  <p>这是一个<strong>重要</strong>的公告。</p>\n  ...\n</section>",
    "theme": "default",
    "fontSize": "medium",
    "wordCount": 156,
    "estimatedReadTime": 1
  },
  "timestamp": 1640995200
}
```

## 高级用法

### 使用不同主题

支持的主题包括：`default`, `classic`, `simple`, `modern` 等。

```bash
curl -X POST "http://localhost:8080/api/v1/convert-and-draft" \
  -H "Content-Type: application/json" \
  -H "Wechat-Appid: wx1234567890abcdef" \
  -H "Wechat-App-Secret: your_app_secret_here" \
  -H "Md2wechat-API-Key: wme_your_api_key_here" \
  -d '{
    "markdown": "# 技术文章\n\n使用 modern 主题的文章内容...",
    "theme": "modern",
    "fontSize": "large"
  }'
```

### 复杂 Markdown 内容

```bash
curl -X POST "http://localhost:8080/api/v1/convert-and-draft" \
  -H "Content-Type: application/json" \
  -H "Wechat-Appid: wx1234567890abcdef" \
  -H "Wechat-App-Secret: your_app_secret_here" \
  -H "Md2wechat-API-Key: wme_your_api_key_here" \
  -d '{
    "markdown": "# 完整功能演示\n\n## 文本格式\n\n这是**粗体**文本，这是*斜体*文本，这是~~删除线~~文本。\n\n## 列表\n\n### 无序列表\n- 项目一\n- 项目二\n  - 子项目\n  - 另一个子项目\n\n### 有序列表\n1. 第一步\n2. 第二步\n3. 第三步\n\n## 链接和图片\n\n[访问 GitHub](https://github.com)\n\n## 代码\n\n内联代码：`console.log(\"Hello\")`\n\n```javascript\n// 代码块\nfunction greet(name) {\n  console.log(`Hello, ${name}!`);\n}\n\ngreet(\"World\");\n```\n\n## 引用\n\n> 这是一个重要的引用内容。\n> \n> 可以包含多行。\n\n## 表格\n\n| 功能 | 状态 | 描述 |\n|------|------|------|\n| 转换 | ✅ | Markdown 到 HTML |\n| 草稿 | ✅ | 微信公众号草稿 |\n| 主题 | ✅ | 多种样式主题 |\n\n## 结论\n\n这个 API 提供了完整的 Markdown 转换功能！",
    "theme": "classic",
    "fontSize": "medium"
  }'
```

## 错误处理

### 参数验证错误

```bash
# 缺少必填 Header 的请求
curl -X POST "http://localhost:8080/api/v1/convert-and-draft" \
  -H "Content-Type: application/json" \
  -d '{
    "markdown": "# 标题",
    "theme": "default",
    "fontSize": "medium"
  }'
```

错误响应：
```json
{
  "code": 400,
  "message": "认证失败",
  "details": "缺少 Wechat-Appid 请求头",
  "timestamp": 1640995200
}
```

### 微信 API 错误

```json
{
  "code": 500,
  "message": "创建微信草稿失败",
  "details": "40001: invalid credential, access_token is invalid or not latest",
  "timestamp": 1640995200
}
```

## 健康检查

```bash
curl -X GET "http://localhost:8080/health"
```

响应：
```json
{
  "status": "OK",
  "time": 1640995200
}
```

## 使用 JavaScript/Node.js

```javascript
const axios = require('axios');

async function convertAndCreateDraft() {
  try {
    const response = await axios.post('http://localhost:8080/api/v1/convert-and-draft', {
      markdown: '# 我的文章\n\n这是文章内容...',
      theme: 'default',
      fontSize: 'medium'
    }, {
      headers: {
        'Content-Type': 'application/json',
        'Wechat-Appid': 'wx1234567890abcdef',
        'Wechat-App-Secret': 'your_app_secret_here',
        'Md2wechat-API-Key': 'wme_your_api_key_here'
      }
    });

    console.log('草稿创建成功！');
    console.log('Media ID:', response.data.data.media_id);
    console.log('字数统计:', response.data.data.wordCount);
  } catch (error) {
    console.error('请求失败:', error.response?.data || error.message);
  }
}

convertAndCreateDraft();
```

## 使用 Python

```python
import requests
import json

def convert_and_create_draft():
    url = "http://localhost:8080/api/v1/convert-and-draft"
    
    headers = {
        "Content-Type": "application/json",
        "Wechat-Appid": "wx1234567890abcdef",
        "Wechat-App-Secret": "your_app_secret_here",
        "Md2wechat-API-Key": "wme_your_api_key_here"
    }
    
    data = {
        "markdown": "# 我的文章\n\n这是文章内容...",
        "theme": "default",
        "fontSize": "medium"
    }
    
    try:
        response = requests.post(url, json=data, headers=headers)
        response.raise_for_status()
        
        result = response.json()
        print("草稿创建成功！")
        print(f"Media ID: {result['data']['media_id']}")
        print(f"字数统计: {result['data']['wordCount']}")
        
    except requests.exceptions.RequestException as e:
        print(f"请求失败: {e}")

if __name__ == "__main__":
    convert_and_create_draft()
```

## 配置说明

### 主题选项
- `default`: 默认主题
- `classic`: 经典主题  
- `simple`: 简洁主题
- `modern`: 现代主题

### 字体大小选项
- `small`: 小字体
- `medium`: 中等字体
- `large`: 大字体

## Header 参数说明

### 必填 Header
- `Wechat-Appid`: 微信公众号的 AppID
- `Wechat-App-Secret`: 微信公众号的 AppSecret

### 可选 Header
- `Md2wechat-API-Key`: MD2WeChat 服务的 API 密钥（如果服务需要认证）

## 注意事项

1. **安全性**: 敏感信息（AppID、AppSecret、API Key）放在 Header 中，不要在 URL 或日志中暴露
2. **微信配置**: 确保使用的 AppID 和 AppSecret 是有效的，且公众号已通过认证
3. **内容限制**: 微信公众号对文章内容有一定限制，请确保内容符合规范
4. **频率限制**: 注意 API 调用频率，避免超出微信接口限制
5. **错误重试**: 建议在生产环境中实现适当的错误重试机制
6. **Header 验证**: 确保在请求中包含所有必需的 Header 参数

## 故障排除

### 常见问题

1. **40001 错误**: 检查 AppID 和 AppSecret 是否正确
2. **40003 错误**: 检查公众号是否有创建草稿的权限
3. **网络超时**: 检查网络连接和服务可用性

### 调试建议

1. 使用健康检查接口验证服务状态
2. 查看服务日志获取详细错误信息
3. 使用简单的 Markdown 内容测试基本功能
4. 确认微信公众号配置正确 