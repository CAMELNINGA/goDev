FROM golang:1.17 AS build


RUN mkdir /goDev
WORKDIR /goDev
COPY go.mod ./

RUN go mod download
COPY . .

RUN go build -o ./bin/goDev ./cmd/Devops

FROM scratch

COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

COPY --from=build /goDev/migrations /migrations
COPY --from=build /goDev/bin/goDev /goDev

CMD ["./goDev" ]