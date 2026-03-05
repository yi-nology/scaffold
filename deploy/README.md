# Scaffold 项目部署指南

现代化的全栈项目脚手架部署解决方案，基于 Docker Compose 实现一键部署，专为中国开发者优化。

## 🎯 项目特色

- **全栈一体化**：集成 Vue 3 前端 + Go 后端 + SQLite 数据库
- **零配置部署**：Docker Compose 一键启动所有服务
- **中国网络优化**：内置 goproxy.cn、npmmirror.com 等国内镜像源
- **安全认证机制**：可选的 API 访问密钥认证
- **数据持久化**：SQLite 数据库外部挂载，数据永不丢失
- **生产就绪**：包含健康检查、日志管理、资源监控等企业级特性

## 🏗️ 技术架构

```mermaid
graph TD
    A[用户请求] --> B[Nginx反向代理]
    B --> C[Vue 3 前端应用]
    B --> D[Go 后端API]
    D --> E[SQLite 数据库]
    
    subgraph "容器化部署"
        B
        C
        D
        E
    end
    
    subgraph "数据持久化"
        F[./data/] --> E
        G[./templates/] --> D
    end
```

## 🚀 快速开始

### 1. 环境准备

确保已安装以下工具：
- Docker Engine 20.10+
- Docker Compose 1.29+

### 2. 获取项目

```bash
git clone https://github.com/yi-nology/scaffold.git
cd scaffold/deploy
```

### 3. 配置环境

```bash
# 复制环境配置文件
cp .env.example .env

# 编辑配置（可选）
vim .env
```

### 4. 启动服务

**标准部署方式：**
```bash
./deploy.sh up
```

**中国网络优化部署（推荐）：**
```bash
# 首次运行：配置国内镜像源
./deploy-cn.sh optimize

# 启动服务
./deploy-cn.sh up
```

## 📁 项目结构

```
deploy/
├── backend/           # 后端服务配置
│   └── Dockerfile     # Go 应用构建配置
├── frontend/          # 前端服务配置
│   ├── Dockerfile     # Vue 应用构建配置
│   └── nginx.conf     # Nginx 反向代理配置
├── data/             # 数据持久化目录
│   └── scaffold.db    # SQLite 数据库文件
├── templates/        # 模板文件目录
├── docker-compose.yml # 服务编排配置
├── .env.example      # 环境变量模板
├── deploy.sh         # 标准部署脚本
├── deploy-cn.sh      # 中国网络优化脚本
└── README.md         # 本文档
```

## ⚙️ 核心配置

### 环境变量配置

在 `.env` 文件中可配置：

```bash
# 服务端口配置
FRONTEND_PORT=80        # 前端服务端口
BACKEND_PORT=9090       # 后端服务端口

# 安全配置
ACCESS_KEY=your-secret-key  # API 访问密钥（可选）
TZ=Asia/Shanghai           # 时区设置

# 调试配置
DEBUG=false                # 调试模式开关
```

### 访问认证配置

**启用认证：**
```bash
# 在 .env 文件中设置
ACCESS_KEY=your-secure-key-here
```

**API 请求认证：**
```bash
# Bearer Token 方式（推荐）
curl -H "Authorization: Bearer your-secure-key-here" \
     http://localhost/api/templates

# 直接传递方式
curl -H "Authorization: your-secure-key-here" \
     http://localhost/api/templates
```

## 🛠️ 部署命令详解

### 标准部署脚本 (deploy.sh)

```bash
# 启动所有服务
./deploy.sh up

# 停止所有服务
./deploy.sh down

# 重建并启动服务
./deploy.sh rebuild

# 查看服务状态
./deploy.sh status

# 查看实时日志
./deploy.sh logs

# 清理所有容器和数据卷
./deploy.sh clean

# 显示帮助信息
./deploy.sh help
```

### 中国网络优化脚本 (deploy-cn.sh)

```bash
# 配置国内镜像源（首次运行）
./deploy-cn.sh optimize

# 使用国内镜像源构建和启动
./deploy-cn.sh rebuild

# 启动服务
./deploy-cn.sh up

# 其他命令同标准脚本
```

## 🔧 进阶配置

### 自定义 Nginx 配置

修改 `frontend/nginx.conf` 文件来自定义反向代理规则：

```nginx
server {
    listen 80;
    server_name localhost;
    
    # SPA 应用支持
    location / {
        try_files $uri $uri/ /index.html;
    }
    
    # API 转发
    location /api/ {
        proxy_pass http://backend:9090;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
    }
}
```

### 资源限制配置

在 `docker-compose.yml` 中添加资源限制：

```yaml
services:
  frontend:
    deploy:
      resources:
        limits:
          memory: 512M
          cpus: '0.5'
  backend:
    deploy:
      resources:
        limits:
          memory: 1G
          cpus: '1.0'
```

## 📊 监控与维护

### 查看服务状态

```bash
# Docker Compose 状态
docker-compose ps

# 资源使用情况
docker stats

# 容器日志
docker-compose logs -f
```

### 数据库管理

```bash
# 直接操作 SQLite 数据库
sqlite3 ./data/scaffold.db

# 备份数据库
cp ./data/scaffold.db ./data/scaffold_backup.db

# 恢复数据库
cp ./data/scaffold_backup.db ./data/scaffold.db
```

### 模板文件管理

```bash
# 查看模板目录
ls ./templates/

# 在容器中查看挂载情况
docker exec scaffold-backend ls /app/templates
```

## 🔒 安全最佳实践

### 1. 访问控制
- 启用 API 访问密钥认证
- 后端服务不直接对外暴露端口
- 使用非 root 用户运行容器

### 2. 数据保护
- 定期备份 SQLite 数据库
- 使用 HTTPS（生产环境）
- 配置适当的文件权限

### 3. 网络安全
- 限制容器间网络通信
- 启用 Docker 内容信任
- 定期更新基础镜像

## 🐛 常见问题排查

### 端口冲突
```bash
# 检查端口占用
netstat -tlnp | grep :80
netstat -tlnp | grep :9090

# 修改 .env 文件中的端口配置
```

### 构建失败
```bash
# 清理构建缓存
docker-compose build --no-cache

# 查看详细构建日志
docker-compose build --progress=plain
```

### 权限问题
```bash
# 给脚本添加执行权限
chmod +x deploy.sh deploy-cn.sh

# 修复数据目录权限
sudo chown -R $(id -u):$(id -g) data/ templates/
```

## 🌏 中国网络优化

本项目针对中国网络环境进行了深度优化：

### 镜像源配置
- **Docker Hub**：阿里云、网易、中科大镜像源
- **Go Modules**：goproxy.cn 国内代理
- **NPM**：npmmirror.com 淘宝镜像
- **Alpine Packages**：阿里云镜像源

### 使用优化脚本
```bash
# 首次配置镜像源
./deploy-cn.sh optimize

# 后续构建将自动使用国内镜像
./deploy-cn.sh rebuild
```

## 🔄 升级与维护

### 应用升级
```bash
# 拉取最新代码
git pull origin main

# 重建服务
./deploy-cn.sh rebuild
```

### 数据迁移
```bash
# 备份当前数据
cp -r data/ data_backup_$(date +%Y%m%d)

# 执行升级后的数据迁移（如有需要）
docker exec scaffold-backend ./scaffold migrate
```

## 📈 生产环境部署

### 推荐配置
1. **SSL/TLS**：配置 HTTPS 证书
2. **负载均衡**：使用外部负载均衡器
3. **监控告警**：集成 Prometheus + Grafana
4. **日志收集**：配置 ELK 栈
5. **备份策略**：自动化数据库备份

### 性能调优
```bash
# 启用 Docker BuildKit
export DOCKER_BUILDKIT=1

# 并行构建
docker-compose build --parallel

# 资源限制
# 在 docker-compose.yml 中配置 CPU 和内存限制
```

## 🤝 贡献指南

欢迎提交 Issue 和 Pull Request 来改进项目！

### 开发环境搭建
```bash
# 克隆项目
git clone https://github.com/yi-nology/scaffold.git

# 安装依赖
cd web && npm install

# 启动开发服务器
make dev-web

# 启动后端服务
make dev
```

## 📄 许可证

MIT License - 详见 [LICENSE](../LICENSE) 文件

## 🙏 致谢

- [Vue.js](https://vuejs.org/) - 前端框架
- [Go](https://golang.org/) - 后端语言
- [Docker](https://www.docker.com/) - 容器化平台
- [SQLite](https://www.sqlite.org/) - 嵌入式数据库

---
**项目主页**：[https://github.com/yi-nology/scaffold](https://github.com/yi-nology/scaffold)  
**问题反馈**：[Issues](https://github.com/yi-nology/scaffold/issues)  
**贡献代码**：[Pull Requests](https://github.com/yi-nology/scaffold/pulls)