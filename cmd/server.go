package cmd

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/huohen92/anifeed/request"
	"github.com/huohen92/anifeed/webui"
	"github.com/spf13/cobra"
)

var (
	port int
	host string

	serverCmd = &cobra.Command{
		Use:   "server",
		Short: "Start the Web UI and HTTP API",
		Long:  "Start the anifeed Web UI, JSON API and the legacy POST /add endpoint.",
		Run:   runWebServer,
	}

	webCmd = &cobra.Command{
		Use:   "web",
		Short: "Start the Web UI (alias of server)",
		Run:   runWebServer,
	}
)

func init() {
	for _, cmd := range []*cobra.Command{serverCmd, webCmd} {
		cmd.Flags().IntVarP(&port, "port", "p", 8115, "server port")
		cmd.Flags().StringVar(&host, "host", "", "listen address (default 0.0.0.0)")
		rootCmd.AddCommand(cmd)
	}
}

func runWebServer(cmd *cobra.Command, _ []string) {
	cfg, source := loadConfig(cmd)
	if err := request.ReloadConfig(); err != nil {
		log.Printf("[warning] 加载站点配置失败: %v", err)
	}

	srv := webui.New(cfg, source)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	if err := srv.Start(ctx); err != nil {
		log.Fatalln(err)
	}
	log.Println("anifeed Web UI 已退出")
}
