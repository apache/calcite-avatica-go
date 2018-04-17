FROM golang:1.10-alpine

WORKDIR /go/src/github.com/apache/calcite-avatica-go
COPY . .
RUN apk --no-cache --update add git \
    && go get -u github.com/golang/dep/cmd/dep \
    && dep ensure -v

CMD ["python", "app.py"]