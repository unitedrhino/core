FROM docker.unitedrhino.com/unitedrhino/golang:1.23.4-alpine3.21 as go-builder
ARG frontFile
WORKDIR /unitedrhino/
COPY ./ ./
RUN echo "Front file URL: $frontFile"
RUN mkdir -p front
RUN cd front&&wget -O front.tgz $frontFile || true
RUN cd front&&tar -xvzf front.tgz
RUN cd front&&ls -l
RUN cd front&&rm -rf front.tgz
RUN go env -w GOPROXY=https://goproxy.cn,direct
RUN go mod download
RUN cd ./service/apisvr && go mod tidy && go build  -tags no_k8s -ldflags="-s -w" -o coresvr .

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
