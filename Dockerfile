FROM golang:1.18.8 as BUILDER

MAINTAINER zengchen1024<chenzeng765@gmail.com>

# build binary
WORKDIR /go/src/github.com/opensourceways/robot-gitee-software-package
COPY . .
RUN GO111MODULE=on CGO_ENABLED=0 go build -a -o robot-gitee-software-package .

WORKDIR /go/src/github.com/opensourceways/robot-gitee-software-package/download
RUN GO111MODULE=on CGO_ENABLED=0 go build -a -o download .

# copy binary config and utils
FROM alpine:3.14
COPY --from=BUILDER /go/src/github.com/opensourceways/robot-gitee-software-package/robot-gitee-software-package /opt/app/robot-gitee-software-package
COPY --from=BUILDER /go/src/github.com/opensourceways/robot-gitee-software-package/download/download /opt/app/download
COPY softwarepkg/infrastructure/pullrequestimpl/repo.sh /opt/app/repo.sh
COPY softwarepkg/infrastructure/codeimpl/code.sh /opt/app/code.sh
COPY softwarepkg/infrastructure/template /opt/app/template
RUN chmod +x /opt/app/repo.sh /opt/app/code.sh && apk update && apk add --no-cache git libc6-compat rpm curl

ENTRYPOINT ["/opt/app/robot-gitee-software-package"]
