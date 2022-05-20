package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"

	"github.com/chunked-app/cortex/chunk"
	"github.com/chunked-app/cortex/gql"
	"github.com/chunked-app/cortex/stores/inmem"
	"github.com/chunked-app/cortex/user"
)

var (
	Commit  = "N/A"
	Version = "N/A"
	BuiltOn = "N/A"

	rootCmd = &cobra.Command{
		Short:   "cortex",
		Version: fmt.Sprintf("%s (commit: %s, build_time: %s)", Version, Commit, BuiltOn),
	}
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	rootCmd.AddCommand(
		cmdServe(ctx),
		cmdShowConfigs(),
	)

	_ = rootCmd.Execute()
}

func cmdShowConfigs() *cobra.Command {
	return &cobra.Command{
		Use:   "configs",
		Short: "Show currently loaded configurations",
		Run: func(cmd *cobra.Command, args []string) {
			cfg := loadConfig(cmd)
			_ = yaml.NewEncoder(os.Stdout).Encode(cfg)
		},
	}
}

func cmdServe(ctx context.Context) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "serve",
		Short: "Start HTTP server",
	}

	cmd.Run = func(cmd *cobra.Command, args []string) {
		cfg := loadConfig(cmd)

		store := &inmem.Store{}
		usersAPI := &user.API{Store: store}
		chunksAPI := &chunk.API{Store: store, Users: usersAPI}

		srv := gql.Server{
			UsersAPI:  usersAPI,
			ChunksAPI: chunksAPI,
			SystemInfo: map[string]interface{}{
				"version":    Version,
				"commit_sha": Commit,
				"build_time": BuiltOn,
			},
		}

		addr := fmt.Sprintf("0.0.0.0:%d", cfg.Port)
		log.Infof("starting server at '%s'...", addr)
		if err := srv.Serve(ctx, addr); err != nil {
			log.Fatalf("server stopped with error: %v", err)
		}
		log.Infof("server exited gracefully")
	}
	return cmd
}
