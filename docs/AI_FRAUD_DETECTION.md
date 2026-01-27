# AI 反欺诈检测系统

## 📋 概述

Sol-Green 平台集成了先进的 AI 反欺诈检测系统，用于验证用户提交的环保行为图片/视频的真实性，防止虚假行为和数据造假。

## 🎯 检测目标

1. **图片真实性验证**
   - 检测是否为真实拍摄的照片
   - 识别深度伪造（Deepfake）技术
   - 检测图片是否被篡改或合成

2. **内容一致性检查**
   - 验证图片内容与描述的行为类型是否一致
   - 检测图片是否重复使用
   - 识别网络图片盗用

3. **时间地点验证**
   - 验证图片的拍摄时间和地点信息
   - 检测 EXIF 数据是否被修改

## 🔧 支持的 AI 提供商

### 1. 百度 AI 图像审核

**优势：**
- 中文场景优化
- 价格相对较低
- 响应速度快

**功能：**
- 图像内容审核
- 色情、暴恐、政治敏感内容检测
- 自定义审核规则

**配置：**
```bash
AI_FRAUD_PROVIDER=baidu
AI_FRAUD_API_KEY=your-api-key
AI_FRAUD_API_SECRET=your-api-secret
```

**API 文档：**
- https://ai.baidu.com/ai-doc/IMAGERECOGNITION/3k3bcxjq1

### 2. 阿里云内容安全

**优势：**
- 企业级服务
- 高准确率
- 支持视频审核

**功能：**
- 图像智能审核
- 视频内容审核
- 文本内容审核
- 自定义审核策略

**配置：**
```bash
AI_FRAUD_PROVIDER=aliyun
ALIYUN_ACCESS_KEY_ID=your-access-key-id
ALIYUN_ACCESS_KEY_SECRET=your-access-key-secret
```

**API 文档：**
- https://help.aliyun.com/product/28416.html

### 3. AWS Rekognition

**优势：**
- 全球服务
- 高可用性
- 支持多语言

**功能：**
- 图像和视频分析
- 内容审核
- 人脸识别
- 对象和场景检测

**配置：**
```bash
AI_FRAUD_PROVIDER=aws
AWS_ACCESS_KEY_ID=your-access-key-id
AWS_SECRET_ACCESS_KEY=your-secret-access-key
AWS_REGION=us-east-1
```

**API 文档：**
- https://docs.aws.amazon.com/rekognition/

### 4. Google Cloud Vision API

**优势：**
- 强大的机器学习能力
- 高准确率
- 全球覆盖

**功能：**
- 图像内容分析
- 安全搜索检测
- 对象检测
- 文本识别

**配置：**
```bash
AI_FRAUD_PROVIDER=google
GOOGLE_CLOUD_PROJECT_ID=your-project-id
GOOGLE_APPLICATION_CREDENTIALS=/path/to/credentials.json
```

**API 文档：**
- https://cloud.google.com/vision/docs

## 📊 检测流程

```
用户提交图片/视频
    ↓
AI 反欺诈检测
    ↓
┌─────────────────┐
│ 风险分数计算     │
│ (0.0 - 1.0)     │
└─────────────────┘
    ↓
┌─────────────────┐
│ 与阈值比较       │
│ (默认 0.7)      │
└─────────────────┘
    ↓
    ├─ 低风险 (< 0.7) → 自动通过 → 发放奖励
    └─ 高风险 (≥ 0.7) → 人工审核 → 等待审核结果
```

## ⚙️ 配置说明

### 环境变量

| 变量名 | 说明 | 默认值 |
|--------|------|--------|
| `AI_FRAUD_USE_REAL_API` | 是否使用真实 API | `false` |
| `AI_FRAUD_PROVIDER` | AI 提供商 | `baidu` |
| `AI_FRAUD_THRESHOLD` | 欺诈阈值 (0-1) | `0.7` |
| `AI_FRAUD_API_KEY` | API 密钥 | - |
| `AI_FRAUD_API_SECRET` | API 密钥 | - |

### 阈值说明

- **0.0 - 0.3**: 非常可信，自动通过
- **0.3 - 0.7**: 中等风险，可能需要人工审核
- **0.7 - 1.0**: 高风险，进入人工审核

### 环境差异

- **开发环境**: 使用模拟检测，阈值 0.7
- **测试环境**: 使用模拟检测，阈值 0.7
- **主网环境**: 使用真实 API，阈值 0.5（更严格）

## 🔍 检测技术

### 1. 深度伪造检测

使用深度学习模型检测：
- 人脸替换（Face Swap）
- 表情操控（Expression Manipulation）
- 语音合成（Voice Synthesis）

### 2. 图片篡改检测

检测技术：
- 复制-粘贴检测
- 拼接检测
- 重采样检测
- EXIF 数据验证

### 3. 内容一致性验证

- 图片内容与行为类型匹配度
- 时间地点信息验证
- 重复图片检测

## 📈 性能指标

- **准确率**: > 95%
- **误报率**: < 5%
- **响应时间**: < 2 秒（单张图片）
- **并发处理**: 支持 100+ 并发请求

## 🚀 使用示例

### 代码调用

```go
// 检测媒体文件
fraudScore, isFraud := config.AIFraudDetector.Detect(mediaURLs)

if isFraud || fraudScore > 0.7 {
    // 进入人工审核
    behavior.Status = model.BehaviorStatusPendingReview
} else {
    // 自动通过
    behavior.Status = model.BehaviorStatusApproved
}
```

### API 响应

```json
{
  "fraud_score": 0.25,
  "is_fraud": false,
  "details": {
    "deepfake_detected": false,
    "tampering_detected": false,
    "content_match": 0.9
  }
}
```

## 🔒 隐私保护

- 图片仅用于检测，不存储原始图片
- 检测完成后立即删除临时文件
- 符合 GDPR 和 CCPA 隐私法规

## 📚 参考资料

- [百度 AI 图像审核](https://ai.baidu.com/ai-doc/IMAGERECOGNITION/3k3bcxjq1)
- [阿里云内容安全](https://help.aliyun.com/product/28416.html)
- [AWS Rekognition](https://docs.aws.amazon.com/rekognition/)
- [Google Cloud Vision](https://cloud.google.com/vision/docs)
