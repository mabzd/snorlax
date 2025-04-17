FROM golang:1.24.2-alpine AS builder

WORKDIR /app

COPY go.mod go.sum .
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o ./build/dbm.exe ./cmd/dbm/
RUN CGO_ENABLED=0 GOOS=linux go build -o ./build/api.exe ./cmd/api/

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/build .

EXPOSE 8080

# For the purpose of this showcase, we run the migration tool (dbm) directly
# befre the service executable (api) to simplify setup and demonstration.
# In a real production environment, migrations should be handled as a
# separate deployment step in the CI/CD pipeline
CMD ["sh", "-c", "./dbm.exe && ./api.exe"]