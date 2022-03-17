FROM golang:1.15 AS build

RUN useradd -u 10001 gopher

RUN mkdir /goDev
WORKDIR /goDev
COPY go.mod  ./

RUN go mod download
COPY . .

RUN go build s -w -X goDev/internal/version.Version=1 -X goDev/internal/version.Commit=1 -o ./bin/goDev ./cmd/goDev

FROM scratch

COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build /etc/passwd /etc/passwd

USER gopher

COPY --from=build /goDev/migrations /migrations
COPY --from=build /goDev/bin/goDev /goDev

CMD ["./credit-history"]