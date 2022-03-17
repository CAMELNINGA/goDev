FROM golang:1.15 AS build

RUN useradd -u 10001 gopher

RUN mkdir /goDev
WORKDIR /goDev
COPY go.mod go.sum ./

RUN go mod download
COPY . .

RUN make build

FROM scratch

COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build /etc/passwd /etc/passwd

USER gopher

COPY --from=build /goDev/migrations /migrations
COPY --from=build /goDev/bin/goDev /goDev

CMD ["./credit-history"]