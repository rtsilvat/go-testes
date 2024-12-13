# Etapa de build
FROM golang:1.21 AS builder

WORKDIR /app

# Copia os arquivos do projeto
COPY . .

# Instala as dependências
RUN go mod tidy

# Compila a aplicação
# Gera o binário na pasta correta
RUN go build -o /app/main .

# Etapa final
FROM debian:bookworm-slim

WORKDIR /app

# Copia o binário compilado da etapa anterior
COPY --from=builder /app/main /app/main
COPY .env .env

# Garante permissões de execução
RUN chmod +x /app/main

# Porta exposta
EXPOSE 8080

# Comando para rodar o binário
CMD ["/app/main"]
