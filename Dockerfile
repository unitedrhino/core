FROM golang:1.21.13-alpine3.20 as go-builder
ARG frontFile
WORKDIR /ithings/
COPY ./ ./
RUN go env -w GOPROXY=https://goproxy.cn,direct
RUN go mod download
RUN cd ./service/apisvr && go mod tidy && go build -ldflags="-s -w" .
RUN echo "Front file URL: $frontFile"
RUN wget -O front.tgz $frontFile
RUN tar -xvzf front.tgz
RUN ls -l
RUN rm -rf front.tgz

#FROM alpine:3.20  as web-builder
#ARG frontFile
#WORKDIR /ithings/
#RUN echo "Front file URL: $frontFile"
#RUN wget -O front.tgz $frontFile
#RUN tar -xvzf front.tgz
#RUN ls -l
#RUN rm -rf front.tgz

FROM alpine:3.20
LABEL homepage="https://github.com/i-Things/iThings"
ENV TZ Asia/Shanghai
RUN apk add tzdata

WORKDIR /ithings/
COPY --from=go-builder /ithings/service/apisvr/apisvr ./apisvr
COPY --from=go-builder /ithings/deploy/conf/core/etc/ ./etc
COPY --from=go-builder /ithings/front/* ./dist/app

ENTRYPOINT ["./apisvr"]
