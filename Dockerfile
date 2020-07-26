FROM golang AS builder
WORKDIR /app
COPY . .
RUN GOPROXY="https://goproxy.cn" go mod tidy && \
	GO111MODULE=on CGO_ENABLED=0 GOOS=linux go build -v -a -o bin/drone-git-merge ./src/...

FROM alpine:latest
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.ustc.edu.cn/g' /etc/apk/repositories && \
	apk update && \
	apk add ca-certificates && \
	rm -rf /var/cache/apk/*

COPY --from=builder /app/bin/drone-git-merge /bin

ENTRYPOINT ["/bin/drone-git-merge"]
