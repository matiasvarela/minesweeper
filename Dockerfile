FROM golang:1.14 AS builder

WORKDIR /go/src/github.com/matiasvarela/minesweeper

COPY . .

RUN go mod download
RUN go build -o /minesweeper main.go

FROM ubuntu

COPY --from=builder /minesweeper /app/minesweeper

ENV ENV=production

ENTRYPOINT ./app/minesweeper