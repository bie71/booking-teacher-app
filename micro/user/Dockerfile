### Multi-stage Dockerfile for the user service
##
# Stage 1: Build the Go binary. We use an official Golang image and
# compile the user service from source. The binary is statically linked
# (CGO disabled) so it can run in a minimal container.
FROM golang:1.24-alpine AS builder
WORKDIR /app
# Copy the entire repository into the build context. This allows the
# service to reference modules in other directories during compilation.
COPY . .
# Change to the user service directory and download dependencies
WORKDIR /app/micro/user
RUN go mod tidy
# Build the user service binary
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/app

# Stage 2: Create a small runtime image. We use Alpine Linux and install
# tzdata to support time zone configuration. The Jakarta time zone is
# configured by copying the relevant zoneinfo file and setting the TZ
# environment variable.
FROM alpine:3.18
RUN apk add --no-cache tzdata && \
    cp /usr/share/zoneinfo/Asia/Jakarta /etc/localtime && \
    echo "Asia/Jakarta" > /etc/timezone
WORKDIR /root/
COPY --from=builder /app/micro/user/main ./
# Copy email templates into the container. These templates are used by
# the email service to render password reset emails. Placing them in
# /templates allows the SMTP_TEMPLATE_PATH environment variable to
# reference files relative to this directory.
COPY micro/user/templates /templates
# Set the container time zone
ENV TZ=Asia/Jakarta
# Expose the default application port
EXPOSE 8081
# Launch the compiled binary
CMD ["./main"]