FROM golang:1.22.0-alpine

# 设置工作目录
WORKDIR /app

# 安装依赖
COPY go.mod go.sum ./
RUN go mod download && go mod verify

# 拷贝整个项目并执行编译命令
COPY . .
RUN go build -v -o myapp .

# 运行
EXPOSE 443
CMD ["./myapp"]

# 进入终端使用docker命令打包镜像和运行
# $ docker build -t your_image_name .
# $ docker run -it --rm --name your_container_name your_image_name
# -it选项表示以交互模式运行容器，--rm选项表示容器停止后立即删除
# 正常运行后可以查看container的ip地址，使用ip:443访问
# docker inspect -f '{{range .NetworkSettings.Networks}}{{.IPAddress}}{{end}}' <container_name_or_id>
# curl <container_ip>:443