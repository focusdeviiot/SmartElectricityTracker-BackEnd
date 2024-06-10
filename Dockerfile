# Use the official Golang image as the base image for building
FROM --platform=$BUILDPLATFORM golang:1.22.4 AS build

# Set the working directory
WORKDIR /app

# Copy go.mod and go.sum
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the rest of the application code
COPY . ./

# Set environment variables for the target platform
ARG TARGETOS
ARG TARGETARCH

# Build the Go application
RUN CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH go build -o backend ./cmd/smart_electricity_tracker_backend/main.go

# Use a minimal base image for the final image
FROM --platform=$TARGETPLATFORM alpine:3.20.0

# Install required packages including tzdata
RUN apk --no-cache add ca-certificates tzdata

# Set the timezone
ENV TZ=Asia/Bangkok
RUN ln -sf /usr/share/zoneinfo/Asia/Bangkok /etc/localtime && echo "Asia/Bangkok" > /etc/timezone

# Set the working directory
WORKDIR /root/

# Copy the Go binary from the build stage
COPY --from=build /app/configs/config.yaml .
COPY --from=build /app/backend .

# Expose port 8080
EXPOSE 8080

# Command to run the executable
CMD ["./backend"]
