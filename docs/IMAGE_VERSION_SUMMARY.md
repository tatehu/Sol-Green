# 图片版本说明文档

## 📸 图片版本概览

项目包含**36张功能展示图片**，分为中文版和英文版两个版本。

### 图片统计
- **中文版图片**：18张（用于中文文档）
- **英文版图片**：18张（用于英文文档，带_EN后缀）
- **总计**：36张图片

## 🎯 图片命名规则

### 中文版图片
- 命名格式：`XX_description.png`
- 示例：`01_homepage.png`、`02_wallet_connect.png`
- 用途：在 `PROJECT_INTRO_CN.md` 中使用

### 英文版图片
- 命名格式：`XX_description_EN.png`
- 示例：`01_homepage_EN.png`、`02_wallet_connect_EN.png`
- 用途：在 `PROJECT_INTRO_EN.md` 中使用

## 📋 图片列表

### 基础功能展示（8张）

| 编号 | 中文版文件名 | 英文版文件名 | 说明 |
|-----|------------|------------|------|
| 01 | 01_homepage.png | 01_homepage_EN.png | 平台首页 |
| 02 | 02_wallet_connect.png | 02_wallet_connect_EN.png | 钱包连接 |
| 03 | 03_green_submit.png | 03_green_submit_EN.png | 提交环保行为 |
| 04 | 04_challenges.png | 04_challenges_EN.png | 挑战活动列表 |
| 05 | 05_create_challenge.png | 05_create_challenge_EN.png | 创建挑战 |
| 06 | 06_marketing_activities.png | 06_marketing_activities_EN.png | 营销活动 |
| 07 | 07_behavior_status.png | 07_behavior_status_EN.png | 查询行为状态 |
| 08 | 08_workflow.png | 08_workflow_EN.png | 完整操作流程 |

### 详细功能展示（10张）

| 编号 | 中文版文件名 | 英文版文件名 | 说明 |
|-----|------------|------------|------|
| 09 | 09_waste_sorting_detail.png | 09_waste_sorting_detail_EN.png | 垃圾分类详情 |
| 10 | 10_tree_planting_detail.png | 10_tree_planting_detail_EN.png | 植树造林详情 |
| 11 | 11_low_carbon_travel_detail.png | 11_low_carbon_travel_detail_EN.png | 低碳出行详情 |
| 12 | 12_challenge_detail.png | 12_challenge_detail_EN.png | 挑战活动详情 |
| 13 | 13_sign_in_activity.png | 13_sign_in_activity_EN.png | 每日签到活动 |
| 14 | 14_user_dashboard.png | 14_user_dashboard_EN.png | 用户个人中心 |
| 15 | 15_ai_detection_flow.png | 15_ai_detection_flow_EN.png | AI反欺诈检测 |
| 16 | 16_blockchain_transaction.png | 16_blockchain_transaction_EN.png | 区块链交易 |
| 17 | 17_third_party_verification.png | 17_third_party_verification_EN.png | 第三方认证 |
| 18 | 18_mobile_responsive.png | 18_mobile_responsive_EN.png | 移动端界面 |

## 📄 文档中的图片引用

### 中文文档 (`PROJECT_INTRO_CN.md`)
```markdown
![首页界面](./assets/01_homepage.png)
![钱包连接](./assets/02_wallet_connect.png)
...
```

### 英文文档 (`PROJECT_INTRO_EN.md`)
```markdown
![Homepage Interface](./assets/01_homepage_EN.png)
![Wallet Connection](./assets/02_wallet_connect_EN.png)
...
```

## ✅ 验证检查

### 图片文件检查
```bash
# 检查中文版图片（18张）
ls -1 assets/*.png | grep -v "_EN" | wc -l
# 应该输出：18

# 检查英文版图片（18张）
ls -1 assets/*_EN.png | wc -l
# 应该输出：18

# 检查总图片数（36张）
ls -1 assets/*.png | wc -l
# 应该输出：36
```

### 文档引用检查
```bash
# 检查中文文档中的图片引用（应该没有_EN）
grep "_EN" PROJECT_INTRO_CN.md
# 应该没有输出

# 检查英文文档中的图片引用（应该都有_EN）
grep "assets.*\.png" PROJECT_INTRO_EN.md | grep -v "_EN"
# 应该没有输出（所有引用都应该有_EN）
```

## 🔍 图片内容差异

### 中文版图片特点
- 界面文字为中文
- 按钮、标签、提示信息均为中文
- 符合中文用户使用习惯

### 英文版图片特点
- 界面文字为英文
- 按钮、标签、提示信息均为英文
- 符合英文用户使用习惯

## 📝 使用建议

1. **文档编写**：
   - 中文文档始终使用中文版图片（不带_EN后缀）
   - 英文文档始终使用英文版图片（带_EN后缀）

2. **图片更新**：
   - 更新功能时，需要同时更新中英文两个版本的图片
   - 保持图片内容与文档语言一致

3. **路径引用**：
   - 使用相对路径：`./assets/XX_description.png`
   - 确保图片文件存在于 `assets/` 目录

4. **预览验证**：
   - 使用 `README_IMAGES.md` 快速预览所有图片
   - 验证图片路径是否正确

## 🚀 快速预览

查看所有图片的快速预览：
- [图片预览页面](../README_IMAGES.md)

查看详细文档：
- [中文版项目介绍](../PROJECT_INTRO_CN.md)
- [English Project Introduction](../PROJECT_INTRO_EN.md)

---

**版本**: v3.0.0  
**更新日期**: 2026-01-27
