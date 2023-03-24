FROM golang:1.19.6-alpine

WORKDIR /usr/app

COPY ../go.mod ./

COPY ../go.sum ./

RUN go mod download

COPY . ./

RUN go build ./producer/web_crawler.go

CMD [ "./web_crawler" ]