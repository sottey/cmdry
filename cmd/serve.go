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
		discovery := plugins.Discoverer{Directory: cfg.PluginDir, Timeout: cfg.PluginTimeout, Logger: logger, Filter: demoPluginAllowed(cfg.DemoMode)}
		initialReport, _ := discovery.DiscoverWithDiagnostics(context.Background(), registry)

		app, err := server.New(cfg, registry, logger, initialReport)
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

// Demo mode intentionally allows only transform-only plugins. This prevents a
// public demo from exposing information about the host, even if the normal
// plugin directory also contains server plugins.
func demoPluginAllowed(demoMode bool) func(plugins.Manifest) bool {
	if !demoMode {
		return nil
	}
	return func(manifest plugins.Manifest) bool {
		if manifest.Category == "server" || len(manifest.Permissions) != 1 {
			return false
		}
		return manifest.Permissions[0] == "data.transform"
	}
}

func init() {
	rootCmd.AddCommand(serveCmd)
	serveCmd.Flags().String("addr", "", "listen address (overrides CMDRY_ADDR)")
	serveCmd.Flags().String("plugin-dir", "", "plugin directory (overrides CMDRY_PLUGIN_DIR)")
	serveCmd.Example = fmt.Sprintf("  %s serve --addr 127.0.0.1:8080", "cmdry")
}
