# 配置管理指南

Scaffold 支持多种配置方式，包括命令行参数、环境变量和默认配置。

## 配置优先级

配置加载遵循以下优先级（从高到低）：
1. 命令行参数
2. 环境变量
3. 默认配置

## 配置选项

### 服务器配置

| 环境变量 | 命令行参数 | 默认值 | 说明 |
|---------|-----------|--------|------|
| `PORT` | `--port` | `9090` | 服务器端口 |
| `HOST` | `--host` | `0.0.0.0` | 服务器主机地址 |
| `CACHE_DIR` | `--cache-dir` | `./cache` | 缓存目录路径 |
| `DEBUG` | `--debug` | `false` | 是否启用调试模式 |

### 数据库配置

| 环境变量 | 命令行参数 | 默认值 | 说明 |
|---------|-----------|--------|------|
| `DB_ENABLE` | `--db-enable` | `false` | 是否启用数据库 |
| `DB_DRIVER` | `--db-driver` | `sqlite` | 数据库驱动 (sqlite/mysql/postgres) |
| `DB_HOST` | `--db-host` | `localhost` | 数据库主机 |
| `DB_PORT` | `--db-port` | `3306` | 数据库端口 |
| `DB_USER` | `--db-user` | `root` | 数据库用户名 |
| `DB_PASSWORD` | `--db-password` | `` | 数据库密码 |
| `DB_NAME` | `--db-name` | `./data/scaffold.db` | 数据库名或文件路径 |

### 安全配置

| 环境变量 | 命令行参数 | 默认值 | 说明 |
|---------|-----------|--------|------|
| `ACCESS_KEY` | `--access-key` | `` | API 访问密钥（可选） |

### 日志配置

| 环境变量 | 默认值 | 说明 |
|---------|--------|------|
| `LOG_LEVEL` | `info` | 日志级别 (debug/info/warn/error) |
| `LOG_FORMAT` | `text` | 日志格式 (text/json) |
| `LOG_OUTPUT` | `stdout` | 日志输出目标 (stdout/file) |

## 使用示例

### 1. 命令行方式

```bash
# 基本启动
./scaffold --serve

# 自定义端口和数据库
./scaffold --serve --port 8080 --db-enable --db-driver mysql --db-name mydb

# 启用认证
./scaffold --serve --access-key my-secret-key
```

### 2. 环境变量方式

```bash
# 设置环境变量
export PORT=8080
export DB_ENABLE=true
export DB_DRIVER=mysql
export DB_NAME=mydb
export ACCESS_KEY=my-secret-key

# 启动服务
./scaffold --serve
```

### 3. Docker 方式

```bash
# 在 docker-compose.yml 中配置环境变量
environment:
  - PORT=9090
  - DB_ENABLE=true
  - DB_DRIVER=sqlite
  - DB_NAME=/app/data/scaffold.db
  - ACCESS_KEY=your-secret-key
```

### 4. 配置文件方式

创建 `.env` 文件：

```bash
# .env
PORT=9090
DB_ENABLE=true
DB_DRIVER=sqlite
DB_NAME=./data/scaffold.db
ACCESS_KEY=your-secret-key
DEBUG=false
```

然后加载配置：

```bash
source .env
./scaffold --serve
```

## 容器化部署配置

在 `deploy/.env` 文件中配置所有环境变量：

```bash
# 服务端口
FRONTEND_PORT=80
BACKEND_PORT=9090

# 服务器配置
PORT=9090
HOST=0.0.0.0
CACHE_DIR=./cache
DEBUG=false

# 数据库配置
DB_ENABLE=true
DB_DRIVER=sqlite
DB_NAME=/app/data/scaffold.db

# 安全配置
ACCESS_KEY=your-secret-key

# 日志配置
LOG_LEVEL=info
LOG_FORMAT=text
```

## 配置验证

启动服务时会显示当前使用的配置：

```bash
$ ./scaffold --serve
Scaffold server running at http://0.0.0.0:9090
Database connected successfully
```

## 故障排除

### 常见问题

1. **端口被占用**
   ```bash
   # 检查端口使用情况
   netstat -tlnp | grep 9090
   
   # 更换端口
   export PORT=8080
   ```

2. **数据库连接失败**
   ```bash
   # 检查数据库配置
   export DB_ENABLE=true
   export DB_DRIVER=sqlite
   export DB_NAME=./data/scaffold.db
   ```

3. **认证失败**
   ```bash
   # 确保 ACCESS_KEY 正确设置
   export ACCESS_KEY=correct-key
   
   # API 请求时携带正确的认证头
   curl -H "Authorization: Bearer correct-key" http://localhost:9090/api/templates
   ```

### 调试模式

启用调试模式查看更多详细信息：

```bash
export DEBUG=true
./scaffold --serve
```

这将显示：
- 详细的配置信息
- 数据库连接状态
- 请求日志
- 错误堆栈跟踪