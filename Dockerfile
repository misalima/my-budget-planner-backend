FROM golang:1.22.4-alpine

WORKDIR /my-budget-planner-backend

COPY go.mod go.sum ./

RUN go mod download && go mod verify

COPY . .

WORKDIR ./cmd/app

RUN go build -o ./bin/mbp-app .

EXPOSE 8080
ENTRYPOINT [ "./bin/mbp-app" ]