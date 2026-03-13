# 他妈的我们吃什么！

> 选不出今天吃什么？让它帮你决定！

## 线上访问

🌐 **[https://what-the-fuck-do-we-eat.onrender.com](https://what-the-fuck-do-we-eat.onrender.com)**

> **注意**：免费版 Render 服务在长时间无人访问后会进入休眠状态，首次访问可能需要等待约 30~60 秒冷启动。

## 本地运行

**前置条件**：已安装 [Go 1.21+](https://golang.org/dl/)

```bash
# 克隆仓库
git clone https://github.com/ZhouNan2020/What-the-fuck-do-we-eat-.git
cd What-the-fuck-do-we-eat-

# 运行
go run .
```

打开浏览器访问：`http://localhost:8080`

## 部署工作流说明

本仓库使用 **GitHub Actions + [Render.com](https://render.com)** 实现自动部署。

### 工作流文件

`.github/workflows/deploy.yml` 在每次向 `main` 分支推送时自动执行：

1. **Build & Test**：运行 `go test ./...` 和 `go build`，确保代码无误。
2. **Deploy to Render**：测试通过后，通过 Render Deploy Hook 触发线上自动部署。

### 首次部署步骤

1. 在 [Render.com](https://render.com) 注册账号并连接 GitHub 仓库。
2. 选择 **New → Blueprint**，Render 会自动读取 `render.yaml` 创建服务。
3. 进入 Render 服务页面 → **Settings → Deploy Hook**，复制 Deploy Hook URL。
4. 在 GitHub 仓库 → **Settings → Secrets and variables → Actions** 中添加：

   | Secret 名称              | 值                          |
   |--------------------------|-----------------------------|
   | `RENDER_DEPLOY_HOOK_URL` | 从 Render 复制的 Deploy Hook URL |

5. 推送代码到 `main` 分支，GitHub Actions 将自动构建并部署。

## 技术栈

- **语言**：Go 1.24
- **HTTP 框架**：标准库 `net/http`
- **模板**：`html/template`（模板文件通过 `//go:embed` 内嵌到二进制）
- **部署**：Docker + Render.com
