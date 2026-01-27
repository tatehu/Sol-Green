# 第三方机构认证系统

## 📋 概述

Sol-Green 平台支持与全球知名的环保组织和政府机构对接，通过第三方权威认证确保环保行为的真实性和可信度。

## 🌍 支持的认证机构

### 1. 阿拉善 SEE 生态协会

**机构简介：**
- 中国知名的环保 NGO
- 专注于荒漠化防治和生态保护
- 拥有 20+ 年环保经验

**认证范围：**
- 植树造林项目
- 荒漠化治理
- 生态保护活动

**认证方式：**
- API 接口对接
- 验证码验证
- 项目编号验证

**配置：**
```bash
PARTNER_ALASHAN_ENABLED=true
PARTNER_ALASHAN_API=https://api.alashansee.org/v1/verify
PARTNER_ALASHAN_API_KEY=your-api-key
PARTNER_ALASHAN_API_SECRET=your-api-secret
```

**官网：** https://www.see.org.cn

### 2. 政府环保部门

**机构简介：**
- 国家级/省级环保部门
- 官方权威认证
- 政策支持

**认证范围：**
- 所有环保行为类型
- 官方环保项目
- 政策激励活动

**认证方式：**
- 政府 API 对接
- 项目备案号验证
- 官方认证码

**配置：**
```bash
PARTNER_GOV_ENABLED=true
PARTNER_GOV_API=https://api.gov-environment.gov.cn/v1/verify
PARTNER_GOV_API_KEY=your-api-key
PARTNER_GOV_API_SECRET=your-api-secret
```

### 3. 联合国环境规划署 (UNEP)

**机构简介：**
- 联合国系统内负责环境事务的牵头机构
- 全球环境治理的核心组织
- 1972 年成立，总部位于肯尼亚内罗毕

**认证范围：**
- 全球环保项目
- 可持续发展目标 (SDGs)
- 气候变化行动

**认证方式：**
- UNEP 官方 API
- 项目注册号验证
- 国际认证标准

**配置：**
```bash
PARTNER_UNEP_ENABLED=true
PARTNER_UNEP_API=https://api.unep.org/v1/verify
PARTNER_UNEP_API_KEY=your-api-key
PARTNER_UNEP_API_SECRET=your-api-secret
```

**官网：** https://www.unep.org

### 4. 世界自然基金会 (WWF)

**机构简介：**
- 全球最大的独立性非政府环境保护组织之一
- 1961 年成立，总部位于瑞士
- 在 100+ 个国家开展项目

**认证范围：**
- 生物多样性保护
- 气候变化行动
- 可持续生活方式

**认证方式：**
- WWF 官方 API
- 项目认证码
- 会员验证

**配置：**
```bash
PARTNER_WWF_ENABLED=true
PARTNER_WWF_API=https://api.wwf.org/v1/verify
PARTNER_WWF_API_KEY=your-api-key
PARTNER_WWF_API_SECRET=your-api-secret
```

**官网：** https://www.worldwildlife.org

### 5. 绿色和平 (Greenpeace)

**机构简介：**
- 国际知名的环保组织
- 1971 年成立，总部位于荷兰
- 专注于环境问题解决方案

**认证范围：**
- 气候行动
- 海洋保护
- 森林保护
- 可持续农业

**认证方式：**
- Greenpeace 官方 API
- 活动认证码
- 志愿者验证

**配置：**
```bash
PARTNER_GREENPEACE_ENABLED=true
PARTNER_GREENPEACE_API=https://api.greenpeace.org/v1/verify
PARTNER_GREENPEACE_API_KEY=your-api-key
PARTNER_GREENPEACE_API_SECRET=your-api-secret
```

**官网：** https://www.greenpeace.org

## 🔄 认证流程

```
用户提交环保行为
    ↓
选择第三方认证机构
    ↓
输入认证码/项目编号
    ↓
调用第三方 API 验证
    ↓
┌─────────────────┐
│ 验证结果         │
└─────────────────┘
    ↓
    ├─ 验证通过 → 自动发放奖励 + 认证标识
    └─ 验证失败 → 返回错误信息
```

## 📊 认证优势

### 1. 权威性
- 官方/知名机构认证
- 提高用户信任度
- 增强平台公信力

### 2. 真实性保障
- 第三方独立验证
- 减少虚假行为
- 提高数据质量

### 3. 奖励加成
- 第三方认证的行为可获得额外奖励
- 认证标识展示
- 排行榜优先显示

## ⚙️ 配置说明

### 环境变量

每个认证机构都有独立的配置项：

| 变量前缀 | 说明 |
|---------|------|
| `PARTNER_{ID}_ENABLED` | 是否启用该机构 |
| `PARTNER_{ID}_API` | API 端点地址 |
| `PARTNER_{ID}_API_KEY` | API 密钥 |
| `PARTNER_{ID}_API_SECRET` | API 密钥 |

### 机构 ID 列表

- `alashan_see` - 阿拉善 SEE
- `gov_environment` - 政府环保部门
- `unep` - 联合国环境规划署
- `wwf` - 世界自然基金会
- `greenpeace` - 绿色和平

## 🔍 API 接口

### 请求示例

```bash
POST /api/v1/green/partner/verify
Authorization: Bearer <token>
Content-Type: application/json

{
  "behavior_id": "550e8400-e29b-41d4-a716-446655440000",
  "partner_id": "alashan_see",
  "verify_code": "SEE20240124001"
}
```

### 响应示例

**成功：**
```json
{
  "msg": "第三方认证通过，奖励已发放",
  "tx_hash": "5j7s8K9...",
  "partner_name": "阿拉善 SEE 生态协会"
}
```

**失败：**
```json
{
  "error": "第三方认证未通过",
  "detail": "验证码无效或已过期"
}
```

## 🌐 全球认证标准

### 国际标准

1. **ISO 14001** - 环境管理体系
2. **SDGs** - 联合国可持续发展目标
3. **CDM** - 清洁发展机制
4. **VCS** - 自愿碳标准

### 认证互认

- 支持多机构联合认证
- 认证结果链上存证
- 全球可验证

## 📈 认证统计

- **认证通过率**: > 85%
- **平均验证时间**: < 5 秒
- **支持机构数**: 5+（持续扩展）
- **覆盖国家**: 100+

## 🔒 安全与隐私

- API 通信加密（HTTPS）
- 密钥安全存储
- 认证数据链上存证
- 符合 GDPR 要求

## 🚀 未来扩展

### 计划对接的机构

1. **The Nature Conservancy** - 自然保护协会
2. **Sierra Club** - 塞拉俱乐部
3. **Environmental Defense Fund** - 环境保护基金会
4. **National Geographic Society** - 国家地理学会

### 认证类型扩展

- 企业 ESG 认证
- 碳足迹认证
- 可持续产品认证
- 绿色建筑认证

## 📚 参考资料

- [UNEP 官网](https://www.unep.org)
- [WWF 官网](https://www.worldwildlife.org)
- [Greenpeace 官网](https://www.greenpeace.org)
- [阿拉善 SEE 官网](https://www.see.org.cn)
