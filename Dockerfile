FROM golang:1.16-alpine as BUILD

WORKDIR /app

COPY go.mod .
COPY go.sum .
COPY main.go .

RUN go mod tidy
RUN CGO_ENABLED=0 go build -o /product-service


FROM scratch
COPY --from=build /product-service /product-service
ENTRYPOINT [ "/product-service" ]
