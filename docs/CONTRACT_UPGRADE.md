# 智能合约升级管理

## 📋 概述

智能合约升级管理系统支持合约版本管理、升级验证、升级执行等运维操作。

## 🔄 升级方式

### 1. 程序升级（Solana）

Solana 使用 BPF Loader 支持程序升级：

```bash
# 部署新版本程序
anchor deploy --provider.cluster devnet

# 升级程序
solana program deploy --program-id <program-id> <new-program.so>
```

### 2. 代理模式（EVM）

EVM 链使用代理合约模式：

- **实现合约**: 业务逻辑合约
- **代理合约**: 可升级的代理
- **升级合约**: 管理升级的合约

## 📊 版本管理

### 版本号规则

- **主版本号**: 不兼容的重大更新
- **次版本号**: 新增功能，向后兼容
- **修订版本号**: 问题修复

示例：`v1.2.3`

### 版本检查

```bash
GET /api/v1/admin/contract/info?program_id=<program-id>
```

响应：
```json
{
  "program_id": "...",
  "current_version": "1.0.0",
  "latest_version": "1.1.0",
  "upgrade_available": true
}
```

## 🔍 升级流程

### 1. 检查升级

```bash
GET /api/v1/admin/contract/upgrade/check?program_id=<program-id>
```

### 2. 验证升级

- 新程序账户存在
- 程序可执行
- 版本号正确
- 兼容性检查

### 3. 准备升级

```bash
POST /api/v1/admin/contract/upgrade/prepare
{
  "program_id": "...",
  "new_program_id": "..."
}
```

### 4. 执行升级

```bash
POST /api/v1/admin/contract/upgrade/execute
{
  "program_id": "...",
  "new_program_id": "..."
}
```

### 5. 验证升级结果

- 检查交易确认
- 验证新程序运行
- 回滚机制（如失败）

## 🔒 安全措施

### 权限控制

- 仅管理员可执行升级
- 多签验证（生产环境）
- 时间锁机制

### 升级前检查

- 代码审计
- 测试网验证
- 兼容性测试
- 数据迁移计划

### 回滚机制

- 保留旧版本程序
- 快速回滚能力
- 数据备份

## 📈 升级历史

### 记录内容

- 升级时间
- 版本变更
- 交易哈希
- 升级原因
- 升级结果

### 查询历史

```bash
GET /api/v1/admin/contract/upgrade/history?program_id=<program-id>
```

## 🎯 最佳实践

### 升级前

1. **充分测试**: 在测试网完整测试
2. **代码审计**: 专业审计机构审计
3. **通知用户**: 提前通知用户升级计划
4. **备份数据**: 完整备份链上数据

### 升级中

1. **低峰期执行**: 选择低峰期升级
2. **监控系统**: 实时监控升级过程
3. **准备回滚**: 随时准备回滚

### 升级后

1. **验证功能**: 验证所有功能正常
2. **监控异常**: 持续监控异常情况
3. **用户反馈**: 收集用户反馈

## 📝 升级检查清单

- [ ] 代码审计完成
- [ ] 测试网验证通过
- [ ] 兼容性测试通过
- [ ] 数据迁移计划制定
- [ ] 回滚方案准备
- [ ] 用户通知发送
- [ ] 升级时间确定
- [ ] 监控系统就绪
- [ ] 备份数据完成
- [ ] 升级文档更新

## 🔧 工具支持

### Solana

- `anchor upgrade`: Anchor 升级命令
- `solana program deploy`: Solana CLI 部署
- `solana program show`: 查看程序信息

### EVM

- Hardhat Upgrades Plugin
- OpenZeppelin Upgrades
- Truffle Migrations

## 📚 参考资料

- [Solana Program Upgrades](https://docs.solana.com/cli/deploy-a-program#upgrading-a-program)
- [OpenZeppelin Upgrades](https://docs.openzeppelin.com/upgrades-plugins/1.x/)
- [Hardhat Upgrades](https://docs.openzeppelin.com/upgrades-plugins/1.x/hardhat-upgrades)
