FROM golang:1.21.13-alpine3.20 as go-builder
WORKDIR /ithings/
COPY ./go.mod ./go.mod
RUN go mod download
COPY ./ ./
RUN cd ./service/apisvr && go mod tidy && go build .

FROM alpine:3.16  as web-builder
WORKDIR /ithings/
ENV fileUrl=${frontFIle}
RUN wget -O front.tgz ${fileUrl}
RUN tar -xvzf front.tgz
RUN rm -rf front.tgz


FROM alpine:3.16
LABEL homepage="https://github.com/i-Things/iThings"
ENV TZ Asia/Shanghai
RUN apk add tzdata

WORKDIR /ithings/
COPY --from=go-builder /ithings/service/apisvr/apisvr ./apisvr
COPY --from=go-builder /ithings/service/deploy/conf/core/etc/ ./etc
COPY --from=web-builder /ithings/assets/front/* ./dist/app

ENTRYPOINT ["./apisvr"]
