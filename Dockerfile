FROM golang:1.13.1

LABEL maintainer="Incognito Chain <dev@incognito.org>"

WORKDIR /go/src/app
COPY go.mod .
RUN go mod download

COPY . .

RUN go build -o portalfeeders

CMD ["./portalfeeders"]
