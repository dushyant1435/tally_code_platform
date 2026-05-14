# --- build stage ---------------------------------------------------------
FROM golang:1.25-alpine AS build

WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /out/server .

# --- runtime stage -------------------------------------------------------
# Debian slim because it has good package coverage for all four toolchains.
# Image size is ~700 MB; if you only need Python+JS, drop g++ and openjdk.
FROM debian:bookworm-slim

ENV DEBIAN_FRONTEND=noninteractive \
    PATH="/usr/lib/jvm/default-java/bin:${PATH}"

RUN apt-get update && apt-get install -y --no-install-recommends \
        python3 \
        nodejs \
        g++ \
        default-jdk-headless \
        ca-certificates \
    && rm -rf /var/lib/apt/lists/* /usr/share/doc /usr/share/man \
    && useradd --create-home --uid 1000 tally

WORKDIR /app
COPY --from=build /out/server /app/server

USER tally
EXPOSE 8080
ENTRYPOINT ["/app/server"]
