FROM golang:1.17.2-alpine AS build
RUN mkdir /app
WORKDIR /app
COPY . ./
RUN go mod download
RUN go build -o /main
# FROM scratch AS bin
# FROM gcr.io/distroless/base-debian10
FROM alpine
# FROM gcr.io/distroless/base-debian10
# FROM gcr.io/distroless/base
WORKDIR /
COPY --from=build /main /main
COPY .env /
EXPOSE 3000

# USER nonroot:nonroot
CMD [ "/main" ]