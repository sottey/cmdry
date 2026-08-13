package cmd

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/sottey/cmdry/internal/config"
	"github.com/sottey/cmdry/internal/plugins"
	"github.com/sottey/cmdry/internal/server"
	"github.com/spf13/cobra"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the Cmdry web server",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load()
		if err != nil {
			return err
		}
		if addr, _ := cmd.Flags().GetString("addr"); addr != "" {
			cfg.Addr = addr
		}
		if dir, _ := cmd.Flags().GetString("plugin-dir"); dir != "" {
			cfg.PluginDir = dir
		}

		logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: cfg.LogLevel}))
		registry := plugins.NewRegistry()
		discovery := plugins.Discoverer{Directory: cfg.PluginDir, Timeout: cfg.PluginTimeout, Logger: logger}
		discovery.Discover(context.Background(), registry)

		app, err := server.New(cfg, registry, logger)
		if err != nil {
			return err
		}
		httpServer := &http.Server{Addr: cfg.Addr, Handler: app, ReadHeaderTimeout: 5 * time.Second}
		go func() {
			logger.Info("Cmdry started", "addr", cfg.Addr, "plugins", registry.Len())
			if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				logger.Error("server stopped", "error", err)
			}
		}()
		stop := make(chan os.Signal, 1)
		signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
		<-stop
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		return httpServer.Shutdown(ctx)
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
	serveCmd.Flags().String("addr", "", "listen address (overrides CMDRY_ADDR)")
	serveCmd.Flags().String("plugin-dir", "", "plugin directory (overrides CMDRY_PLUGIN_DIR)")
	serveCmd.Example = fmt.Sprintf("  %s serve --addr 127.0.0.1:8080", "cmdry")
}
