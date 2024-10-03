FROM golang:1.22-alpine AS builder

WORKDIR /usr/local/src

COPY ["go.mod", "go.sum", "./"]
RUN go mod download

COPY . ./
RUN go build -o ./bin/paymentService cmd/main.go

FROM alpine AS runner

RUN apk --no-cache add bash make

COPY --from=builder /usr/local/src/bin/paymentService /


ENTRYPOINT /paymentService