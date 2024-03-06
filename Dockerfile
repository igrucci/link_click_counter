FROM golang:1.22-bookworm

COPY ./ ./

RUN go build -o main .

EXPOSE 8080

CMD ["./main"]