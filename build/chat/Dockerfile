FROM golang:alpine AS build_step
ENV GO111MODULE=on
WORKDIR /app
COPY . .
RUN go mod tidy
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o chat cmd/chat/main.go

FROM alpine AS run_step
WORKDIR /app
COPY --from=build_step app/chat .
ENTRYPOINT [ "./chat" ]