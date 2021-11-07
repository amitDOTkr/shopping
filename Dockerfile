FROM golang:1.17.2-alpine AS build
RUN mkdir /app
COPY . ./app
WORKDIR /app
RUN go mod download
RUN go build -o /main
# FROM scratch AS bin
# FROM gcr.io/distroless/base-debian10
# FROM alpine
FROM gcr.io/distroless/base-debian10
# FROM gcr.io/distroless/base
COPY --from=build /main /
EXPOSE 3000
CMD ["/main"]
# FROM golang:1.17.1-alpine AS build
# WORKDIR /src
# COPY . .
# RUN go build -o /out/example .
# FROM scratch AS bin
# COPY --from=build /out/example /
# ENTRYPOINT [ "main" ]