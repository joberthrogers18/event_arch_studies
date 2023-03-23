FROM golang:1.19.6-alpine

WORKDIR /usr/app

COPY ../go.mod ./

COPY ../go.sum ./

RUN go mod download

COPY . ./

RUN ls && go build ./producer/web-crawler.go

CMD [ "./web-crawler" ]