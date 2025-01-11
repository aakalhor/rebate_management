package main

import (
	"awesomeProject2/rebate/cmd"
	_ "awesomeProject2/rebate/docs"
)

// @title Rebate Program Swagger API
// @version 1.0
// @description Comprehensive API documentation for the Rebate Program service. This service handles rebate creation, transaction management, and claims processing.
// @termsOfService http://swagger.io/terms/
// @contact.name Amirali Kalhor
// @contact.url http://www.swagger.io/support
// @contact.email aakalhor2000@gmail.com
// @host localhost:8080
// @BasePath /
// @schemes http
func main() {

	cmd.Boot()
}
