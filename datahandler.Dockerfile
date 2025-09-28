# --- Build stage ---
FROM golang:1.25.1-alpine3.22 AS build

WORKDIR /src

COPY datahandler_source/ ./
RUN go mod tidy

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o datahandler ./main.go


# --- Runtime stage ---
FROM redhat/ubi10-minimal

WORKDIR /interface

COPY --from=build /src/datahandler .

EXPOSE 9876

ENTRYPOINT ["./datahandler"]
