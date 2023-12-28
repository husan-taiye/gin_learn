# 基础镜像
FROM ubuntu:latest
LABEL authors="fox_three"
# 编译后的打包进来镜像，放到工作目录 /app
COPY webook /app/webook
WORKDIR /app
# CMD 是执行命令
ENTRYPOINT ["/app/webook"]