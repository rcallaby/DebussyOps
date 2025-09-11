FROM golang:1.20-alpine
WORKDIR /app
COPY . .
RUN go build -o /bin/orchestrator ./cmd/orchestrator ./examples/agents/todo ./examples/agents/todo
EXPOSE 8080
CMD ["/bin/orchestrator"]