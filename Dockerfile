FROM golang:1.17.1-alpine AS build
WORKDIR /src
COPY . .
RUN go build -o /out/example .
FROM scratch AS bin
COPY --from=build /out/example /