# Build-Stage
FROM golang:1.24.5-alpine AS build
WORKDIR /app

# Install build dependencies
RUN apk add --no-cache gcc musl-dev curl

# Copy the source code
COPY . .

# Install templ
RUN go install github.com/a-h/templ/cmd/templ@latest

# Generate templ files
RUN templ generate

# Build Tailwind CSS (after templ generate so all classes are detected)
RUN curl -sLO https://github.com/tailwindlabs/tailwindcss/releases/latest/download/tailwindcss-linux-arm64 && \
    chmod +x tailwindcss-linux-arm64 && \
    ./tailwindcss-linux-arm64 -i ./assets/css/input.css -o ./assets/css/output.css --minify

# Build the application
RUN CGO_ENABLED=1 GOOS=linux go build -o main ./main.go

# Deploy-Stage
FROM alpine:3.20.2
WORKDIR /app

# Install ca-certificates
RUN apk add --no-cache ca-certificates

# Set environment variable for runtime
ENV APP_ENV=production

# Copy the binary from the build stage
COPY --from=build /app/main .

# Expose the port your application runs on
EXPOSE 8090

# Command to run the application
CMD ["./main"]
