FROM gliderlabs/alpine:3.4

LABEL Author=macdfree

# 安装git
RUN apk add --no-cache git

# 复制blogci
COPY blogci /opt/bin/
COPY hugo /opt/bin/
COPY goblog /opt/goblog

# 设置环境变量
ENV PATH /opt/bin:$PATH
