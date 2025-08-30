FROM golang:1.20-alpine
WORKDIR /app
COPY . .
RUN go build -o /bin/orchestrator ./cmd/orchestrator
EXPOSE 8080
CMD ["/bin/orchestrator"]