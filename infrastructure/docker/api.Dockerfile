FROM golang:1.24-alpine AS build
WORKDIR /src
COPY go.mod ./
COPY apps/api ./apps/api
COPY internal ./internal
RUN CGO_ENABLED=0 go build -trimpath -ldflags='-s -w' -o /api ./apps/api
FROM gcr.io/distroless/static-debian12:nonroot
COPY --from=build /api /api
USER nonroot:nonroot
EXPOSE 8080
ENTRYPOINT ["/api"]
