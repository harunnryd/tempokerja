package cmd

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/harunnryd/tempokerja/config"
	"github.com/harunnryd/tempokerja/internal/app/handler"
	"github.com/harunnryd/tempokerja/internal/app/repo"
	"github.com/harunnryd/tempokerja/internal/app/server"
	"github.com/harunnryd/tempokerja/internal/app/usecase"
	"github.com/spf13/cobra"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
)

func init() {
	cobra.OnInitialize()
}

var rootCmd = &cobra.Command{
	Use:   "tempodoloe",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
			examples and usage of using your application.`,
	Run: func(cmd *cobra.Command, args []string) {
		start()
	},
}

// Execute executes the root command.
func Execute() (err error) {
	if err = rootCmd.Execute(); err != nil {
		fmt.Println(err.Error())
		return
	}
	return
}

func start() {
	logger := log.New(os.Stdout, "http: ", log.LstdFlags)
	cfg := config.NewConfig()
	done := make(chan bool, 1)
	quit := make(chan os.Signal, 1)

	signal.Notify(quit, os.Interrupt)

	temporalClient := initTemporalClient(cfg)
	repo := repo.New(cfg)
	usecase := usecase.New(repo)
	worker := initWorker(temporalClient, usecase, repo)

	s := server.NewServer(
		net.JoinHostPort(cfg.GetString("server.host"), cfg.GetString("server.port")),
		handler.New(cfg, temporalClient, usecase),
		time.Duration(cfg.GetInt("server.read_timeout"))*time.Second,
		time.Duration(cfg.GetInt("server.write_timeout"))*time.Second,
		time.Duration(cfg.GetInt("server.idle_timeout"))*time.Second,
	)

	httpServer := s.GetHTTPServer()
	go s.GracefullShutdown(httpServer, logger, quit, done)

	logger.Println("=> http server started on", net.JoinHostPort(cfg.GetString("server.host"), cfg.GetString("server.port")))
	if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Fatalf("Could not listen on %s: %v\n", cfg.GetString("server.port"), err)
	}

	worker.Stop()

	<-done

	logger.Println("Server stopped")
}

func initTemporalClient(cfg config.Config) client.Client {
	temporalClientOptions := client.Options{HostPort: net.JoinHostPort(cfg.GetString("temporal.host"), cfg.GetString("temporal.port"))}
	temporalClient, err := client.NewClient(temporalClientOptions)
	if err != nil {
		log.Fatal("cannot start temporal client: " + err.Error())
	}
	return temporalClient
}

func initWorker(temporalClient client.Client, usecase usecase.Usecase, repo repo.Repo) worker.Worker {
	workerOptions := worker.Options{
		MaxConcurrentWorkflowTaskExecutionSize: 4,
	}
	worker := worker.New(temporalClient, "CREATE_ORDER", workerOptions)
	worker.RegisterWorkflow(usecase.Order().CreateOrder)
	worker.RegisterActivity(repo.Product().GetProductByID)
	worker.RegisterActivity(repo.Product().DeductQuantityByID)
	worker.RegisterActivity(repo.Transaction().CreateOrder)
	worker.RegisterActivity(repo.Transaction().CancelOrderByID)

	err := worker.Start()
	if err != nil {
		log.Fatal("cannot start temporal worker: " + err.Error())
	}

	return worker
}
