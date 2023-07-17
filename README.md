# Serviço de API de Email

## Visão Geral

Bem-vindo ao Serviço de API de Email! 📧✉️

Este serviço permite que você integre facilmente funcionalidades de email em suas aplicações, sem a complicação de lidar diretamente com as complexidades do SMTP (Simple Mail Transfer Protocol). Nossa API incrível permite que você envie emails utilizando seu servidor SMTP preferido com facilidade.

## Configuração

Para começar, por favor forneça as seguintes configurações:

-   **SMTP Host**: Especifique o endereço do seu servidor SMTP.
-   **SMTP Port**: Forneça o número da porta usada pelo seu servidor SMTP.
-   **SMTP Username**: Compartilhe o nome de usuário ou endereço de email associado à sua conta SMTP.
-   **SMTP Password**: Forneça a senha associada à sua conta SMTP.
-   **Default From**: Especifique o endereço de email padrão para o campo "De" nos emails enviados.
-   **Allowed Origin**: Defina o cabeçalho de origem permitida para verificar solicitações feitas a partir de URLs específicas.

## Uso da API

Uma vez que o serviço esteja configurado, você pode utilizar nossa API para enviar emails através de requisições HTTP. Vamos aos detalhes:

### Endpoint

O endpoint base para a nossa API é: `https://api.example.com/api/send-email`

### Método de Requisição

A API suporta apenas o método `POST` para enviar os dados necessários para enviar o email.

### Parâmetros da Requisição

A requisição `POST` deve incluir os seguintes parâmetros no corpo (em formato JSON):

-   `"to"`: O endereço de email do destinatário.
-   `"subject"`: O assunto do email.
-   `"body"`: O corpo do email.

Exemplo de requisição:
 

`{
  "to": "exemplo@dominio.com",
  "subject": "Assunto do email",
  "body": "Conteúdo do email"
}` 

### Resposta da Requisição

Após enviar a requisição, nossa API responderá com um objeto JSON indicando o resultado da operação:

-   Em caso de envio bem-sucedido do email, a resposta seguirá este formato:
 

`{
  "message": "Email enviado com sucesso"
}` 

-   Em caso de falha, a resposta será assim:
 
`{
  "error": "Falha ao enviar o email: <mensagem de erro>"
}` 

## Considerações Finais

Nosso Serviço de API de Email foi criado para simplificar o processo de envio de emails utilizando um servidor SMTP. Certifique-se de fornecer configurações precisas do servidor SMTP para operações sem complicações. Se encontrar qualquer problema durante o uso da API, consulte as mensagens de erro na resposta para solucionar problemas.
