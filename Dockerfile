# This is a standard Dockerfile for building a Go app.
# It is a multi-stage build: the first stage compiles the Go source into a binary, and
#   the second stage copies only the binary into an alpine base.

# -- Stage 1 -- #
# Compile the app.
FROM golang:1.20-alpine as builder
WORKDIR /app
# The build context is set to the directory where the repo is cloned.
# This will copy all files in the repo to /app inside the container.
# If your app requires the build context to be set to a subdirectory inside the repo, you
#   can use the source_dir app spec option, see: https://www.digitalocean.com/docs/app-platform/references/app-specification-reference/
COPY . .
RUN go build -ldflags "-s -w" -o bin/html2png

# -- Stage 2 -- #
# Create the final environment with the compiled binary.
FROM alpine
# Install any required dependencies.
RUN apk --no-cache add ca-certificates
WORKDIR /root/

# Install headless chrome
RUN set -x \
    && apk update \
    && apk upgrade \
  # replacing default repositories with edge ones \
    && echo "https://dl-cdn.alpinelinux.org/alpine/edge/community" >> /etc/apk/repositories \
    && echo "https://dl-cdn.alpinelinux.org/alpine/edge/main" >> /etc/apk/repositories \
  # Add the packages \
    && apk add --no-cache dumb-init curl make gcc g++ linux-headers binutils-gold gnupg libstdc++ nss chromium \
  # Do some cleanup \
    && apk del --no-cache make gcc g++ binutils-gold gnupg libstdc++ \
    && rm -rf /usr/include \
    && rm -rf /var/cache/apk/* /root/.node-gyp /usr/share/man /tmp/*

# Copy the binary from the builder stage and set it as the default command.
COPY --from=builder /app/bin/html2png /usr/local/bin/
CMD ["html2png"]
