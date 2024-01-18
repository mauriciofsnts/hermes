FROM golang:1.21 as builder

# STAGE: BUILD

WORKDIR /app
 
COPY go.mod .
COPY go.sum .
COPY Makefile .
 
COPY internal internal
COPY cmd cmd

# Adicione o comando abaixo para copiar o arquivo config.yml para o diretório /app
COPY config.yml .

RUN make dist

# STAGE: TARGET

FROM alpine:latest 

RUN addgroup -S user && adduser -S user -G user

USER user

WORKDIR /app 
COPY --from=builder /app/hermes /app/hermes

ENTRYPOINT ["/app/hermes"]
