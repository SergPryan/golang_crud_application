FROM golang:1.23.1
WORKDIR /app
COPY . .
#RUN go mod init app && go mod tidy
RUN go build -o main
CMD ["./main"]