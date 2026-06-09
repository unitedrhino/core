FROM docker.unitedrhino.com/unitedrhino/golang:1.26.1-alpine3.23 as go-builder
ARG frontFile
WORKDIR /unitedrhino/
COPY ./ ./
RUN go version
RUN echo "Front file URL: $frontFile"
RUN mkdir -p front
RUN cd front&&wget -O front.tgz $frontFile || true
RUN cd front&&tar -xvzf front.tgz
RUN cd front&&ls -l
RUN cd front&&rm -rf front.tgz
RUN go env -w GOPROXY=https://goproxy.cn,direct
ENV GOPRIVATE=*.gitee.com,gitee.com/*
ENV GONOSUMCHECK=*
ENV GONOSUMDB=*
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.ustc.edu.cn/g' /etc/apk/repositories
RUN apk add git
RUN go mod download
RUN cd ./service/apisvr && go mod tidy && go build  -tags no_k8s  -o coresvr .

FROM docker.unitedrhino.com/unitedrhino/alpine:3.20
LABEL homepage="https://gitee.com/unitedrhino"
ENV TZ Asia/Shanghai
RUN apk add tzdata

WORKDIR /unitedrhino/
COPY --from=go-builder /unitedrhino/service/apisvr/coresvr ./coresvr
COPY --from=go-builder /unitedrhino/service/apisvr/etc ./etc
RUN mkdir -p ./dist/app
COPY --from=go-builder /unitedrhino/front/ ./dist/app

ENTRYPOINT ["./coresvr"]
