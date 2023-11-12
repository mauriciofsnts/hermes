<div align="center">
  <h1 align="center">HERMES</h1>

  <h3>This service makes it easy to add email features to your applications, without the need to handle the technical details of SMTP (Simple Mail Transfer Protocol).</h3>

  <p align="center">
  <img src="https://img.shields.io/badge/HTML5-E34F26.svg?style=flat-square&logo=HTML5&logoColor=white" alt="HTML5" />
  <img src="https://img.shields.io/badge/YAML-CB171E.svg?style=flat-square&logo=YAML&logoColor=white" alt="YAML" />
  <img src="https://img.shields.io/badge/Docker-2496ED.svg?style=flat-square&logo=Docker&logoColor=white" alt="Docker" />
  <img src="https://img.shields.io/badge/Go-00ADD8.svg?style=flat-square&logo=Go&logoColor=white" alt="Go" />
  
  </p>
  <img src="https://img.shields.io/github/license/mauriciofsnts/hermes?style=flat-square&color=5D6D7E" alt="GitHub license" />
  <img src="https://img.shields.io/github/last-commit/mauriciofsnts/hermes?style=flat-square&color=5D6D7E" alt="git-last-commit" />
  <img src="https://img.shields.io/github/commit-activity/m/mauriciofsnts/hermes?style=flat-square&color=5D6D7E" alt="GitHub commit activity" />
  <img src="https://img.shields.io/github/languages/top/mauriciofsnts/hermes?style=flat-square&color=5D6D7E" alt="GitHub top language" />
</div>

---

## 📖 Table of Contents

- [📖 Table of Contents](#-table-of-contents)
- [📍 Overview](#-overview)
- [📦 Features](#-features)
- [🚀 Getting Started](#-getting-started)
  - [📝 Environment Variables](#-environment-variables)
  - [🔧 Installation](#-installation)
  - [🤖 Running hermes](#-running-hermes)
- [🤝 Contributing](#-contributing)
- [📄 License](#-license)

---

## 📍 Overview

This service allows you to seamlessly integrate email functionalities into your applications, without the hassle of dealing directly with the complexities of the Simple Mail Transfer Protocol (SMTP). Our remarkable API enables you to send emails effortlessly using your preferred SMTP server.

---

## 📦 Features

### Rate Limit

This service has a optinal rate limit to prevent abuse. You can make a maximum of 1 request every 30 seconds. If you exceed this limit, you will receive a response with a status of 429 (Too Many Requests) and the following `Retry-After` header indicating the time in seconds to wait before making the next request.

### Templates

Simplifies the setup of customizable email templates through Hermes APIs. With the Go templating language, users can easily include dynamic variables in their templates, ensuring adaptability for specific email content needs.

### Send emails via API

Send emails through a simple API request. The API supports the `POST /api/send-email` method and requires the following parameters in the request body (in JSON format):

```json
{
	"to": "example@tdl.com",
	"subject": "Email subject",
	"body": "Hello, hermes is awesome!"
}
```

or with template:

```json
{
	"to": "example@tdl.com",
	"subject": "Email subject",
	"templateName": "welcome",
	"content": {
		"Content": "Hello, hermes is awesome!",
		"Link": "https://github.com/mauriciofsnts/hermes"
	}
}
```

#### Response example

```json
{
	"message": "Email sent successfully"
}
```

```json
{
	"error": "Failed to send email: <error message>"
}
```

---

### Queue system

The service has a queue system to send emails asynchronously. This ensures that the API response is fast and that the email is sent in the background.

#### Kafka

To use the Kafka queue system, you need to configure the Kafka server. For that, you should set the following environment variables:

```yaml
kafka:
  enabled: true
  topic: string
  brokers:
    - kafka1:19092
```

#### Redis

To use the Redis queue system, you need to configure the Redis server. For that, you should set the following environment variables:

```yaml
redis:
  enabled: true
  topic: string
  port: 6379
  host: localhost
  password: string
```

#### Memory cache

If you prefer not to utilize Redis or Kafka, Memory cache is an alternative option. To enable Memory cache, simply set the "enabled" environment variables for both Redis and Kafka to false.

---

## 📝 Environment Variables

The following environment variables are required to run this service:

| Name          | Description                                               |
| ------------- | --------------------------------------------------------- |
| defaultFrom   | The default email address to use as the sender of emails. |
| allowedOrigin | A list of allowed origins for CORS.                       |
| location      | The folder where the templates will be saved              |
| smtp.host     | The SMTP server host.                                     |
| smtp.port     | The SMTP server port.                                     |
| smtp.username | The SMTP server username.                                 |
| smtp.password | The SMTP server password.                                 |

You can check all the available configurations in the [config_example.yml]() file.

---

## 🔧 Installation

1. Clone the hermes repository:

```sh
git clone https://github.com/mauriciofsnts/hermes
```

2. Change to the project directory:

```sh
cd hermes
```

3. Install the dependencies:

```sh
go mod download
```

### 🤖 Running hermes

Development:

```sh
make dev
```

Prodution:

```sh
make build
```

<details closed>
<summary>Docker</summary>

To run this service using Docker Compose, follow the instructions below:

1. Ensure that Docker and Docker Compose are installed in your environment.

2. In the terminal, navigate to the root directory of your project containing the `docker-compose.yml` and `Dockerfile` files.

3. Execute the following command to start the Email API service:

```bash
docker-compose up
```

4. Wait until Docker Compose builds the images and starts the containers. You will see the service logs in the terminal.

5. The API will be available at `http://127.0.0.1:8293/api/send-email`. You can send POST requests to this endpoint to send emails.

6. To stop the service, press `Ctrl+C` in the terminal and execute the following command to stop and remove the containers:

```bash
docker-compose down
```

## </details>

## 🤝 Contributing

Contributions are welcome! Here are several ways you can contribute:

- **[Submit Pull Requests](https://github.com/mauriciofsnts/hermes/blob/main/CONTRIBUTING.md)**: Review open PRs, and submit your own PRs.
- **[Join the Discussions](https://github.com/mauriciofsnts/hermes/discussions)**: Share your insights, provide feedback, or ask questions.
- **[Report Issues](https://github.com/mauriciofsnts/hermes/issues)**: Submit bugs found or log feature requests for MAURICIOFSNTS.

#### _Contributing Guidelines_

<details closed>
<summary>Click to expand</summary>

1. **Fork the Repository**: Start by forking the project repository to your GitHub account.
2. **Clone Locally**: Clone the forked repository to your local machine using a Git client.
   ```sh
   git clone <your-forked-repo-url>
   ```
3. **Create a New Branch**: Always work on a new branch, giving it a descriptive name.
   ```sh
   git checkout -b new-feature-x
   ```
4. **Make Your Changes**: Develop and test your changes locally.
5. **Commit Your Changes**: Commit with a clear and concise message describing your updates.
   ```sh
   git commit -m 'Implemented new feature x.'
   ```
6. **Push to GitHub**: Push the changes to your forked repository.
   ```sh
   git push origin new-feature-x
   ```
7. **Submit a Pull Request**: Create a PR against the original project repository. Clearly describe the changes and their motivations.

Once your PR is reviewed and approved, it will be merged into the main branch.

</details>

---

## 📄 License

This project is protected under the MIT License. For more details, refer to the [LICENSE](https://github.com/mauriciofsnts/hermes/blob/master/LICENSE) file.
