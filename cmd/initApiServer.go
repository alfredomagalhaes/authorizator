/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/alfredomagalhaes/authorizator/infrastructure/db"
	"github.com/alfredomagalhaes/authorizator/repository"
	"github.com/alfredomagalhaes/authorizator/routes"
	"github.com/alfredomagalhaes/authorizator/services"
	"github.com/gofiber/contrib/fiberzerolog"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
)

// initApiServerCmd represents the initApiServer command
var initApiServerCmd = &cobra.Command{
	Use:   "initApiServer",
	Short: "Initialize the API Server to Access control service.",
	Long: `This application will manage users and their roles to grant/deny
	access to other microservices `,
	Run: func(cmd *cobra.Command, args []string) {

		logger := zerolog.New(os.Stderr).With().Timestamp().Logger()

		app := fiber.New()

		app.Use(fiberzerolog.New(fiberzerolog.Config{
			Logger: &logger,
		}))

		log.SetOutput(logger)

		pgConnConfig := db.PgConnConfig{
			Username:     os.Getenv("DB_USER"),
			Password:     os.Getenv("DB_PASSWORD"),
			Host:         os.Getenv("DB_HOST"),
			Port:         os.Getenv("DB_PORT"),
			DatabaseName: os.Getenv("DB_NAME"),
			TimeZone:     os.Getenv("DB_TIMEZONE"),
		}

		pgDB, err := db.NewPgConnection(pgConnConfig)

		if err != nil {
			log.Fatalf("could not connect to database %v", err)
		}
		//Initialize all repositories
		appRepo, _ := repository.NewPgAppRepository(pgDB)
		roleRepo, _ := repository.NewPgRoleRepository(pgDB)

		//Initialize permission service
		permService, _ := services.NewPermissionService(appRepo, roleRepo)

		//Create all the tables in the database
		appRepo.MigrateTable()

		//Initialize all the routes
		//Creates the sub route for the api version
		apiV1 := app.Group("/api/v1")

		//Initialize "applications" routes
		routes.ApplicationRoute(apiV1, appRepo, permService)
		routes.RolesRoute(apiV1, roleRepo)

		// Listen from a different goroutine
		go func() {
			if err := app.Listen(":3000"); err != nil {
				log.Panic(err)
			}
		}()

		c := make(chan os.Signal, 1)                    // Create channel to signify a signal being sent
		signal.Notify(c, os.Interrupt, syscall.SIGTERM) // When an interrupt or termination signal is sent, notify the channel

		_ = <-c // This blocks the main thread until an interrupt is received
		fmt.Println("Gracefully shutting down...")
		_ = app.Shutdown()

		fmt.Println("Running cleanup tasks...")

		// Your cleanup tasks go here
		// db.Close()
		// redisConn.Close()
		fmt.Println("Fiber was successful shutdown.")
	},
}

func init() {
	rootCmd.AddCommand(initApiServerCmd)

	if err := godotenv.Load(".env"); err != nil {
		log.Fatal("No .env file found")
	}

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// initApiServerCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// initApiServerCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
