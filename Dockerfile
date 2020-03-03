FROM golang:1.13

WORKDIR /go/src/humanitec.io/walhallapiadaptor

# Ideally we would only include ./cmd and ./internal but the Dockerfile does not allow for directories to be copied in one go
COPY . .

RUN go build -o /bin/walhallapiadaptor humanitec.io/walhallapiadaptor/cmd/walhallapiadaptor

ENTRYPOINT ["/bin/walhallapiadaptor"]
