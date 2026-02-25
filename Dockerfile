FROM golang:1.26-alpine AS compile
WORKDIR /app
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 go build -o /bin/app /app/.


FROM scratch AS execute
WORKDIR /bin
COPY --from=compile /bin/app /bin/app
COPY --from=compile /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
ENTRYPOINT [ "/bin/app" ]
