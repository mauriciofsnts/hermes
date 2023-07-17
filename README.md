# Serviço de API de Email

## Visão Geral

Bem-vindo ao Serviço de API de Email! 📧✉️

Este serviço permite que você integre facilmente funcionalidades de email em suas aplicações, sem a complicação de lidar diretamente com as complexidades do SMTP (Simple Mail Transfer Protocol). Nossa API incrível permite que você envie emails utilizando seu servidor SMTP preferido com facilidade.

## Configuração

Para começar, por favor forneça as seguintes configurações:

- **SMTP Host**: Especifique o endereço do seu servidor SMTP.
- **SMTP Port**: Forneça o número da porta usada pelo seu servidor SMTP.
- **SMTP Username**: Compartilhe o nome de usuário ou endereço de email associado à sua conta SMTP.
- **SMTP Password**: Forneça a senha associada à sua conta SMTP.
- **Default From**: Especifique o endereço de email padrão para o campo "De" nos emails enviados.
- **Allowed Origin**: Defina o cabeçalho de origem permitida para verificar solicitações feitas a partir de URLs específicas.

## Uso da API

Uma vez que o serviço esteja configurado, você pode utilizar nossa API para enviar emails através de requisições HTTP. Vamos aos detalhes:

### Endpoint

O endpoint base para a nossa API é: `https://api.example.com/api/send-email`

### Método de Requisição

A API suporta apenas o método `POST` para enviar os dados necessários para enviar o email.

### Parâmetros da Requisição

A requisição `POST` deve incluir os seguintes parâmetros no corpo (em formato JSON):

- `"to"`: O endereço de email do destinatário.
- `"subject"`: O assunto do email.
- `"body"`: O corpo do email.

Exemplo de requisição:

```json
{
  "to": "exemplo@dominio.com",
  "subject": "Assunto do email",
  "body": "Conteúdo do email"
}
```

### Resposta da Requisição

Após enviar a requisição, nossa API responderá com um objeto JSON indicando o resultado da operação:

- Em caso de envio bem-sucedido do email, a resposta seguirá este formato:

```json
{
  "message": "Email enviado com sucesso"
}
```

- Em caso de falha, a resposta será assim:

```json
{
  "error": "Falha ao enviar o email: <mensagem de erro>"
}
```

## Rate Limit

Esta API possui um limite de taxa (rate limit) para evitar abusos. Você pode fazer no máximo 1 requisição a cada 30 segundos. Caso atinja o limite, você receberá uma resposta com status 429 (Muitas Requisições) e o seguinte cabeçalho `Retry-After` com o tempo em segundos para aguardar antes de fazer a próxima requisição.

## Executando com Docker Compose

Para executar este serviço usando Docker Compose, siga as instruções abaixo:

1. Certifique-se de que o Docker e o Docker Compose estejam instalados em seu ambiente.

2. No terminal, navegue até o diretório raiz do seu projeto que contém os arquivos `docker-compose.yml` e `Dockerfile`.

3. Execute o seguinte comando para iniciar o serviço de API de Email:

```bash
docker-compose up
```

4. Aguarde até que o Docker Compose construa as imagens e inicie os contêineres. Você verá os logs do serviço no terminal.

5. A API estará disponível em `http://localhost:8293/api/send-email`. Você pode enviar requisições POST para este endpoint para enviar emails.

6. Para interromper a execução do serviço, pressione `Ctrl+C` no terminal e execute o seguinte comando para parar e remover os contêineres:

```bash
docker-compose down
```

## Considerações Finais

Nosso Serviço de API de Email foi criado para simplificar o processo de envio de emails utilizando um servidor SMTP. Certifique-se de fornecer configurações precisas do servidor SMTP para operações sem complicações. Se encontrar qualquer problema durante o uso da API, consulte as mensagens de erro na resposta para solucionar problemas.

Desejo a você uma experiência fantástica ao utilizar nosso Serviço de API de Email! 😎📧✉️
