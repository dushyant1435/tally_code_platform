# --- build stage ---------------------------------------------------------
FROM golang:1.25-alpine AS build

WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /out/server .

# --- runtime stage -------------------------------------------------------
# We need python3 inside the container because the server shells out to it
# to execute user-submitted code.
FROM python:3.12-alpine

RUN adduser -D -u 1000 tally
WORKDIR /app

COPY --from=build /out/server /app/server

USER tally
EXPOSE 8080
ENTRYPOINT ["/app/server"]
