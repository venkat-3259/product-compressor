package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"

	"zocket/app/controllers"
	"zocket/app/queries"
	"zocket/pkg/configs"
	imageprocessor "zocket/pkg/image_processor"
	"zocket/pkg/middleware"
	"zocket/pkg/routes"
	"zocket/platform/apiserver"
	"zocket/platform/database"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/streadway/amqp"
)

// @title Zocket image process API
// @version 1.0
// @description This is an auto-generated API Docs.
// @termsOfService http://swagger.io/terms/
// @venkat API Support
// @contact.email venkateshwarachinnasamy@gmail.com
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @BasePath /api/v1
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func main() {

	ctx, cancelCtx := signal.NotifyContext(context.Background(), os.Interrupt)

	defer cancelCtx()

	config, err := configs.GetConfig()
	if err != nil {
		log.Println("Failed to load config! Reason: ", err)
		return
	}

	// Connecting to postgres database
	db, err := database.ConnectPostgres(ctx, config.Postgres)
	if err != nil {
		log.Println("Failed to connect postgres! Reason: ", err)
		return
	}
	defer db.Close()

	// Initialize RabbitMQ connection
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Println("Failed to connect with Rabbit MQ: ", err)
	}
	defer conn.Close()

	// Create a channel for communication
	ch, err := conn.Channel()
	if err != nil {
		panic(err)
	}
	defer ch.Close()

	queries := queries.NewQueries(db)

	// Define a new Fiber app with config.
	app := fiber.New(configs.GetFiberConfig(config.Fiber))
	app.Use(recover.New(recover.Config{EnableStackTrace: true}))

	// Middlewares.
	middleware.FiberMiddleware(app) // Register Fiber's middleware for app.

	image := imageprocessor.ImageProcessor{
		Ctx:     ctx,
		Channel: ch,
		Queries: queries,
	}

	go imageprocessor.InitChannel(&image)

	h := controllers.NewHandler(config, &image)

	// Routes.
	routes.RegisterSwaggerRoute(app)
	routes.RegisterPublicRoutes(app, h)
	routes.RegisterNotFoundRoute(app)

	go log.Fatal(app.Listen(fmt.Sprintf("%v:%v", config.Fiber.Host, config.Fiber.Port)))

	apiserver.StartFiberWithGracefulShutdown(app, config.Fiber)
}
