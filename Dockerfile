FROM golang:1.22 as build 

WORKDIR /app
COPY cmd/vport/vport.go /app/
COPY go.mod /app/
COPY go.sum /app/

RUN go mod tidy
RUN go build -o vport vport.go 

FROM ubuntu as runtime 

COPY --from=build /app/vport /bin/vport
RUN apt update && apt install iproute2 iputils-ping -y

COPY entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh

ENTRYPOINT [ "sh", "-c", "/entrypoint.sh" ]