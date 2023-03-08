FROM golang:1.18.8 as BUILDER

MAINTAINER zengchen1024<chenzeng765@gmail.com>

# build binary
WORKDIR /go/src/github.com/opensourceways/robot-gitee-software-package
COPY . .
RUN GO111MODULE=on CGO_ENABLED=0 go build -a -o robot-gitee-software-package .

# copy binary config and utils
FROM alpine:3.14
COPY --from=BUILDER /go/src/github.com/opensourceways/robot-gitee-software-package/robot-gitee-software-package /opt/app/robot-gitee-software-package
COPY event/repo.sh /opt/app/repo.sh
COPY event/template /opt/app/template
RUN chmod +x /opt/app/repo.sh && apk update && apk add --no-cache git libc6-compat

ENTRYPOINT ["/opt/app/robot-gitee-software-package"]
