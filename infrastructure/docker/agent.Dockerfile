FROM golang:1.24-alpine AS build
WORKDIR /src
COPY go.mod ./
COPY apps/agent ./apps/agent
RUN CGO_ENABLED=0 go build -trimpath -ldflags='-s -w' -o /agent ./apps/agent
FROM gcr.io/distroless/static-debian12:nonroot
COPY --from=build /agent /agent
USER nonroot:nonroot
EXPOSE 8090
ENTRYPOINT ["/agent"]
