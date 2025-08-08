FROM golang:1.24
WORKDIR /app
COPY . .
RUN go mod download
RUN go build -o cvmatch ./cmd/main.go

CMD ["./cvmatch"]