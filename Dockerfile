FROM golang:1.19-alpine AS build
WORKDIR /src
RUN apk update && apk add git
RUN git clone --depth 1 -b develop https://github.com/ponyo877/folks.git /src
RUN go mod download
RUN go build -o /folks api/main.go

FROM alpine:latest
WORKDIR /
COPY --from=build /folks /folks
ENTRYPOINT ["/folks"]