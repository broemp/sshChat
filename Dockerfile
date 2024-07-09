FROM golang:1.22.5 AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o sshChat

FROM scratch

WORKDIR /app
COPY --from=builder /app/sshChat .

ENTRYPOINT [ "./sshChat" ]
