FROM golang:1.19 AS build-stage

WORKDIR /weather-api

COPY go.mod go.sum ./
RUN go mod download

COPY src/ ./src

RUN CGO_ENABLED=0 GOOS=linux go build -o /main src/main.go

# Run the tests in the container
# FROM build-stage AS run-test-stage
# RUN go test -v ./...

# Deploy the application binary into a lean image
FROM gcr.io/distroless/base-debian11 AS build-release-stage

WORKDIR /

COPY --from=build-stage /main /main

# Copy the template files
COPY --from=build-stage weather-api/src ./src

EXPOSE 8080

USER nonroot:nonroot

ENTRYPOINT ["/main"]
