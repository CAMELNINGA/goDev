FROM golang:1.17 AS build


COPY go.mod   /goDev
WORKDIR /goDev

RUN go mod download
RUN go mod tidy
COPY . .

RUN CGO_ENABLED=0 go build -o ./bin/goDev ./cmd/Devops

FROM scratch

COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

COPY --from=build /goDev/migrations /migrations
COPY --from=build /goDev/bin/goDev /goDev

CMD ["./goDev" ]
