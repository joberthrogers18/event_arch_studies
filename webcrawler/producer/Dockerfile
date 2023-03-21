FROM golang:1.19.6-alpine

WORKDIR /usr/app

COPY ../go.mod ./

COPY ../go.sum ./

RUN go mod download

COPY *.go ./

RUN go build -o ./web-crawler.go

CMD [ "./web-crawler" ]