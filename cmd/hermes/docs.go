// Package main Hermes Notification Service API
//
//	@title			Hermes API
//	@version		1.0
//	@description	Notification service API with support for email, Discord, and more
//	@termsOfService	http://swagger.io/terms/
//
//	@contact.name	API Support
//	@contact.email	support@hermes.local
//
//	@license.name	MIT
//	@license.url	https://opensource.org/licenses/MIT
//
//	@host		localhost:8080
//	@BasePath	/
//
//	@securityDefinitions.apikey	ApiKeyAuth
//	@in							header
//	@name						X-API-Key
//	@description				API Key authentication
//
//	@tag.name			Health
//	@tag.description	Health check endpoints
//	@tag.name			Notification
//	@tag.description	Notification sending endpoints
//	@tag.name			DLQ
//	@tag.description	Dead Letter Queue management
//	@tag.name			Monitoring
//	@tag.description	Metrics and monitoring

package main
