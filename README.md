# 他妈的我们吃什么！

一个帮你随机决定午餐/晚餐吃什么的小工具，用 Go 写的 Web 服务 + HTML 模板。

## 本地运行

确保已安装 Go 1.24+，然后在项目根目录执行：

```bash
go run .
```

启动后在浏览器打开 [http://localhost:8080/](http://localhost:8080/)，点击按钮随机选餐厅。

## 线上部署（Render）

本项目已配置 `render.yaml`，支持一键部署到 [Render](https://render.com)：

1. Fork 或使用本仓库。
2. 登录 [Render Dashboard](https://dashboard.render.com/)，点击 **New → Blueprint**。
3. 连接你的 GitHub 账号并选择本仓库，Render 会自动读取 `render.yaml`。
4. 点击 **Apply** 创建服务，等待构建完成。
5. 构建完成后，Render 会生成一个公网 URL，格式为：

   ```
   https://what-the-fuck-do-we-eat.onrender.com
   ```

   （实际 URL 以 Render 分配的为准，可在 Dashboard 的服务页面顶部找到。）

后续每次合并到 `main` 分支，Render 会自动触发重新部署。

## 在线访问

> 部署完成后，将 Render 分配的公网 URL 填写到此处：
>
> **🌐 线上地址：**`<在此填入 Render 分配的 URL>`

## 运行测试

```bash
go test ./...
```
