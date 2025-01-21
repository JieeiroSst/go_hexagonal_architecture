package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/JIeeiroSst/hex/internal/config"
	"github.com/JIeeiroSst/hex/internal/core/ports"
	"github.com/JIeeiroSst/hex/internal/core/services"
	"github.com/JIeeiroSst/hex/internal/handlers"
	"github.com/JIeeiroSst/hex/internal/infrastructure/cache"
	"github.com/JIeeiroSst/hex/internal/infrastructure/database"
	"github.com/JIeeiroSst/hex/internal/jobs"
	"github.com/JIeeiroSst/hex/internal/messaging"
	"github.com/JIeeiroSst/hex/internal/repositories"
	"github.com/JIeeiroSst/utils/consul"
	"github.com/gin-gonic/gin"
	"github.com/rabbitmq/amqp091-go"
	"github.com/robfig/cron/v3"
	"github.com/urfave/cli/v2"
	"go.uber.org/fx"
)

type AppConfig struct {
	Port     string
	Env      string
	LogLevel string
}

func main() {
	app := &cli.App{
		Name:  "hexapp",
		Usage: "Hexagonal Architecture Application",
		Commands: []*cli.Command{
			{
				Name:    "api",
				Usage:   "Run the API server",
				Aliases: []string{"serve", "server"},
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "port",
						Aliases: []string{"p"},
						Value:   "8080",
						Usage:   "Port number for the API server",
					},
					&cli.StringFlag{
						Name:    "env",
						Aliases: []string{"e"},
						Value:   "development",
						Usage:   "Environment (development/staging/production)",
					},
				},
				Action: runAPI,
			},
			{
				Name:    "consumer",
				Usage:   "Run the message consumer",
				Aliases: []string{"cons", "worker"},
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "queue",
						Aliases: []string{"q"},
						Value:   "default",
						Usage:   "Queue name to consume messages from",
					},
				},
				Action: runConsumer,
			},
			{
				Name:    "cron",
				Usage:   "Run scheduled jobs",
				Aliases: []string{"scheduler", "jobs"},
				Action:  runCron,
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func runAPI(c *cli.Context) error {
	config := AppConfig{
		Port:     c.String("port"),
		Env:      c.String("env"),
		LogLevel: "debug",
	}

	app := fx.New(
		fx.Provide(
			func() *AppConfig { return &config },
			NewConfig,
			NewGinEngine,
			NewDatabase,
			NewRedisCache,
			repositories.NewUserRepository,
			services.NewUserService,
			handlers.NewUserHandler,
		),
		fx.Invoke(registerRoutes),
	)

	return startApp(app)
}

func runConsumer(c *cli.Context) error {
	config := AppConfig{
		Env:      c.String("env"),
		LogLevel: "debug",
	}

	app := fx.New(
		fx.Provide(
			func() *AppConfig { return &config },
			NewConfig,
			NewDatabase,
			NewRabbitMQConnection,
			fx.Annotate(
				messaging.NewRabbitMQConsumer,
				fx.As(new(ports.MessageConsumer)),
			),
			messaging.NewUserEventHandler,
		),
		fx.Invoke(registerMessageConsumers),
	)

	return startApp(app)
}

func runCron(c *cli.Context) error {
	config := AppConfig{
		Env:      c.String("env"),
		LogLevel: "debug",
	}

	app := fx.New(
		fx.Provide(
			func() *AppConfig { return &config },
			NewConfig,
			NewDatabase,
			NewCron,
			NewRedisCache,
			repositories.NewUserRepository,
			services.NewUserService,
			jobs.NewUserCleanupJob,
		),
		fx.Invoke(registerCronJobs),
	)

	return startApp(app)
}

func startApp(app *fx.App) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle shutdown signals
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigChan
		cancel()
	}()

	return app.Start(ctx)
}

func NewGinEngine(appConfig *AppConfig) *gin.Engine {
	if appConfig.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()
	r.Use(gin.Recovery(), gin.Logger())
	return r
}

func NewCron() *cron.Cron {
	return cron.New(cron.WithSeconds())
}

func NewRabbitMQConnection(config *config.Config) (*amqp091.Connection, error) {
	return amqp091.Dial(config.RabbitMQ.URL)
}

func NewConfig() *config.Config {
	dirEnv, err := config.ReadFileEnv(".env")
	if err != nil {
		return nil
	}
	consul := consul.NewConfigConsul(dirEnv.HostConsul,
		dirEnv.KeyConsul, dirEnv.ServiceConsul)
	var config config.Config
	conf, err := consul.ConnectConfigConsul()
	if err != nil {
		return nil
	}
	if err := json.Unmarshal(conf, &config); err != nil {
		return nil
	}
	return &config
}

func NewDatabase(config *config.Config) (*database.PostgresDB, error) {
	masterConfig := database.DBConfig{
		Host:     config.Database.Master.Host,
		Port:     config.Database.Master.Port,
		User:     config.Database.Master.User,
		Password: config.Database.Master.Password,
		DBName:   config.Database.Master.DBName,
	}

	slaveConfig := database.DBConfig{
		Host:     config.Database.Slave.Host,
		Port:     config.Database.Slave.Port,
		User:     config.Database.Slave.User,
		Password: config.Database.Slave.Password,
		DBName:   config.Database.Slave.DBName,
	}

	return database.NewPostgresDB(masterConfig, slaveConfig)
}

func NewRedisCache(config *config.Config) ports.CacheRepository {
	return cache.NewRedisCache(config.Redis.URL)
}

func registerRoutes(
	lifecycle fx.Lifecycle,
	router *gin.Engine,
	userHandler *handlers.UserHandler,
	appConfig *AppConfig,
) {
	lifecycle.Append(fx.Hook{
		OnStart: func(context.Context) error {
			api := router.Group("/api/v1")
			{
				users := api.Group("/users")
				{
					users.POST("", userHandler.Create)
					users.GET("/:id", userHandler.Get)
					users.PUT("/:id", userHandler.Update)
					users.DELETE("/:id", userHandler.Delete)
				}
			}

			router.GET("/health", func(c *gin.Context) {
				c.JSON(200, gin.H{
					"status": "ok",
					"time":   time.Now().UTC(),
				})
			})

			go func() {
				if err := router.Run(":" + appConfig.Port); err != nil {
					log.Printf("Error starting server: %v", err)
				}
			}()
			return nil
		},
	})
}

func registerMessageConsumers(
	lifecycle fx.Lifecycle,
	consumer ports.MessageConsumer,
	userEventHandler *messaging.UserEventHandler,
) {
	lifecycle.Append(fx.Hook{
		OnStart: func(context.Context) error {
			go func() {
				if err := consumer.StartConsuming("user_events", userEventHandler.HandleMessage); err != nil {
					log.Printf("Error in consumer: %v", err)
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			return consumer.Close()
		},
	})
}

func registerCronJobs(
	lifecycle fx.Lifecycle,
	cronJob *cron.Cron,
	userCleanupJob *jobs.UserCleanupJob,
) {
	lifecycle.Append(fx.Hook{
		OnStart: func(context.Context) error {
			_, err := cronJob.AddFunc("* * * * * *", userCleanupJob.Execute)
			if err != nil {
				return fmt.Errorf("failed to register cleanup job: %w", err)
			}
			cronJob.Start()
			return nil
		},
		OnStop: func(context.Context) error {
			cronJob.Stop()
			return nil
		},
	})
}
