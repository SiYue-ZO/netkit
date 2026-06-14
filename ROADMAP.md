# NetKit 开发路线图

## Phase 1 - 基础框架 + 核心扫描 (v0.1.0) ✅

**目标**：可运行的 CLI 骨架 + 端口扫描 + 存活探测

- [x] 项目初始化 (Go module + Cobra CLI + Viper 配置)
- [x] 全局 flag (timeout/threads/output/format/verbose/no-color)
- [x] 公共类型定义 (pkg/types) + 协程池 (pkg/pool)
- [x] 输出层 (JSON/CSV/Table 三种格式)
- [x] `netkit ping` - ICMP/TCP 存活探测，支持 CIDR/文件输入
- [x] `netkit portscan` - TCP Connect 扫描，Top100/自定义端口
- [x] `netkit cidr` - CIDR 展开/信息/聚合
- [x] `netkit version` - 版本信息
- [x] GoReleaser 跨平台编译 + Makefile
- [x] GitHub Actions CI (lint + 多平台测试 + 自动发布)
- [x] golangci-lint 配置

---

## Phase 2 - 侦察增强 (v0.2.0)

**目标**：信息收集能力闭环

### `netkit dns` - DNS 枚举与查询
- [ ] DNS 记录查询 (A/AAAA/CNAME/MX/TXT/NS/SOA)
- [ ] 子域名爆破 (字典 + 并发)
- [ ] DNSSEC 检测
- [ ] 自定义 DNS 服务器
- [ ] 反向 DNS 查询

### `netkit fingerprint` - 服务指纹识别
- [ ] Banner 抓取 (TCP 连接后读取响应)
- [ ] TLS 指纹识别 (JA3/JA3S)
- [ ] 协议探针 (HTTP/SSH/FTP/SMTP/RDP 等特征匹配)
- [ ] 服务版本识别规则库
- [ ] 与 portscan 集成的自动指纹模式

### `netkit web` - Web 侦察
- [ ] HTTP 探针 (状态码/标题/重定向/响应时间/Content-Length)
- [ ] 基础技术栈识别 (Server header / X-Powered-By / cookie 特征)
- [ ] 支持 HTTP/HTTPS 自动跟随
- [ ] 自定义 Header / User-Agent
- [ ] 并发请求 + 限速

### `netkit whois` - Whois 查询
- [ ] 域名 Whois 查询
- [ ] IP Whois 查询 (ASN/归属/网段)
- [ ] Whois 服务器自动选择
- [ ] 结果结构化输出

### 管道模式
- [ ] stdin/stdout 管道组合支持
- [ ] 示例: `netkit dns -d example.com | netkit portscan | netkit fingerprint`
- [ ] 管道模式自动检测 (stdin 是否有数据)

---

## Phase 3 - 安全检查 (v0.3.0)

**目标**：安全评估能力

### `netkit ssl` - SSL/TLS 安全检查
- [ ] 证书信息 (有效期/颁发者/SAN/链)
- [ ] 协议版本检测 (TLS 1.0/1.1/1.2/1.3)
- [ ] 弱密码套件检测
- [ ] 证书透明度日志检查
- [ ] HSTS 检测
- [ ] 安全评级 (A/B/C/D/F)

### `netkit trace` - 路由追踪
- [ ] MTR 模式路由追踪
- [ ] 逐跳延迟统计
- [ ] 丢包率计算
- [ ] 支持 TCP/ICMP 探测
- [ ] 路径可视化 (ASCII)

### SYN 扫描增强
- [ ] TCP SYN 扫描 (需要 root/管理员权限)
- [ ] 使用 raw socket / gopacket 实现
- [ ] Windows Npcap 支持
- [ ] 自动降级到 Connect 扫描

### UDP 扫描增强
- [ ] 常见 UDP 端口探测 (DNS/SNMP/NTP/TFTP 等)
- [ ] UDP 协议探针
- [ ] ICMP Port Unreachable 判断

---

## Phase 4 - 体验优化 + 生产就绪 (v0.4.0)

**目标**：生产可用

### `netkit report` - 报告生成
- [ ] HTML 报告生成 (内嵌 CSS，无外部依赖)
- [ ] Markdown 报告
- [ ] 汇总统计 (开放端口分布/服务统计/风险概览)
- [ ] 报告模板可自定义

### 配置管理增强
- [ ] 预设扫描模板 (quick/full/stealth/web)
- [ ] 配置文件验证
- [ ] 环境变量覆盖
- [ ] 多配置文件合并

### 用户体验
- [ ] 进度条显示 (扫描进度)
- [ ] Shell 自动补全 (Bash/Zsh/Fish/PowerShell)
- [ ] 彩色输出优化 (支持 NO_COLOR 环境变量)
- [ ] 交互式模式 (TUI，可选)

### 测试与质量
- [ ] 核心模块单元测试覆盖 > 70%
- [ ] 集成测试 (TestMain + 本地测试服务器)
- [ ] Benchmark 基准测试
- [ ] 代码覆盖率报告 (Codecov 集成)

---

## Phase 5 - 高级功能 (v0.5.0+)

**目标**：差异化竞争力

### 指纹规则库
- [ ] CPE 匹配规则
- [ ] Web 指纹库 (CMS/框架/WAF/CDN)
- [ ] 社区贡献规则机制
- [ ] 规则自动更新

### 漏洞检测
- [ ] 基于指纹的 CVE 关联
- [ ] 常见弱口令检测 (SSH/FTP/MySQL/Redis)
- [ ] 基础 Web 漏洞检测 (目录遍历/信息泄露)
- [ ] 自定义 PoC 插件

### 协作与集成
- [ ] 结果数据库存储 (SQLite)
- [ ] REST API 服务模式
- [ ] 与 Nuclei 模板兼容
- [ ] 与 Nmap XML 互操作
- [ ] Elasticsearch 导出

### 性能优化
- [ ] 扫描引擎异步 I/O 优化
- [ ] 内存池复用
- [ ] 大规模扫描 (10万+ IP) 支持
- [ ] 分布式扫描 (可选)

---

## 版本发布节奏

| 版本 | 里程碑 | 核心交付 |
|---|---|---|
| v0.1.0 | Phase 1 完成 | CLI 骨架 + ping + portscan + cidr |
| v0.2.0 | Phase 2 完成 | dns + fingerprint + web + whois + 管道 |
| v0.3.0 | Phase 3 完成 | ssl + trace + SYN/UDP 扫描 |
| v0.4.0 | Phase 4 完成 | 报告 + 补全 + 测试覆盖 |
| v0.5.0 | Phase 5 开始 | 指纹库 + 弱口令 + CVE 关联 |
| v1.0.0 | 生产就绪 | 全功能稳定版 |
