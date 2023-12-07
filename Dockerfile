FROM golang:1.21.2-alpine

WORKDIR /src
COPY go.mod go.sum ./

RUN go mod download
COPY . .

RUN CGO_ENABLED=0 go build -o /bin/app .


FROM alpine
WORKDIR /src
COPY --from=0 /bin/app /bin/app
COPY --from=0 /src/index.html /src/index.html
ENTRYPOINT [ "/bin/app" ]