package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/huohen92/anifeed/config"
	"github.com/huohen92/anifeed/p115"
	"github.com/huohen92/anifeed/request"
	"github.com/huohen92/anifeed/rsssite"
	"github.com/huohen92/anifeed/version"
	"github.com/spf13/cobra"
)

var (
	pAgent        *p115.Agent
	rssUrl        string
	cookies       string
	rssJsonPath   string
	qrLogin       bool
	disableCache  bool
	chunkDelay    int
	chunkSize     int
	cooldownMinMs int
	cooldownMaxMs int
	clearTaskNum  int
	rootCmd       = &cobra.Command{
		Use:     "anifeed",
		Short:   `Add offline tasks to 115`,
		Version: version.Version,
		Run: func(_cmd *cobra.Command, _args []string) {
			initAgent(_cmd)
			if rssJsonPath != "" {
				rsssite.SetRssJsonPath(rssJsonPath)
			}
			if rssUrl != "" {
				pAgent.AddRssUrlTask(rssUrl)
				return
			}
			if clearTaskNum > 0 {
				err := pAgent.OfflineClear(clearTaskNum - 1)
				if err != nil {
					log.Fatalln(err)
				}
				return
			}
			if qrLogin {
				// -q without other actions: just login and save cookies
				log.Println("login successful")
				return
			}
			pAgent.ExecuteAllRssTask()
		},
	}
	// magnet link
	linkUrl   string
	cid       string
	savepath  string
	textFile  string
	magnetCmd = &cobra.Command{
		Use:   "magnet",
		Short: `Add magnet tasks to 115`,
		Run: func(_cmd *cobra.Command, _args []string) {
			initAgent(_cmd)
			magnets := []string{}
			if textFile != "" {
				var err error
				magnets, err = rsssite.GetMagnetsFromText(textFile)
				if err != nil {
					log.Fatalln(err)
				}
			} else if linkUrl != "" {
				magnets = append(magnets, linkUrl)
			}
			if len(magnets) == 0 {
				log.Fatalln("magnets is empty")
			}
			pAgent.AddMagnetTask(magnets, cid, savepath)
		},
	}
)

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&rssUrl, "url", "u", "", "rss url")
	rootCmd.PersistentFlags().StringVar(&cookies, "cookies", "", "115 cookies")
	rootCmd.PersistentFlags().StringVarP(&rssJsonPath, "rss", "r", "", "rss json path")
	rootCmd.PersistentFlags().BoolVarP(&qrLogin, "qrcode", "q", false, "login 115 by qrcode")
	magnetCmd.Flags().StringVarP(&linkUrl, "link", "l", "", "magnet link")
	magnetCmd.Flags().StringVar(&cid, "cid", "", "cid")
	magnetCmd.Flags().StringVar(&savepath, "savepath", "", "save path")
	magnetCmd.Flags().StringVar(&textFile, "text", "", "text file")
	rootCmd.PersistentFlags().BoolVar(&disableCache, "no-cache", false, "skip checking cache in db.sqlite")
	rootCmd.PersistentFlags().IntVar(&chunkDelay, "chunk-delay", 0, "chunk delay. default 2")
	rootCmd.PersistentFlags().IntVar(&chunkSize, "chunk-size", 0, "chunk size. default 200")
	rootCmd.PersistentFlags().IntVar(&cooldownMinMs, "cooldown-min-ms", 1000, "minimum cooldown between 115 API calls in milliseconds. default 1000")
	rootCmd.PersistentFlags().IntVar(&cooldownMaxMs, "cooldown-max-ms", 1100, "maximum cooldown between 115 API calls in milliseconds. default 1100")
	rootCmd.Flags().IntVar(&clearTaskNum, "clear-task-type", 0, "clear offline task type: 1-6.\n 1: OfflineClearDone\n 2: OfflineClearAll\n 3: OfflineClearFailed\n 4: OfflineClearRunning\n 5: OfflineClearDoneAndDelete\n 6: OfflineClearAllAndDelete")
	rootCmd.AddCommand(magnetCmd)
}

func buildCLIParams(cmd *cobra.Command) config.CLIParams {
	cliParams := config.CLIParams{
		Cookies: cookies,
	}

	if commandFlagChanged(cmd, "no-cache") {
		cliParams.DisableCache = disableCache
		cliParams.DisableCacheSet = true
	}
	if commandFlagChanged(cmd, "chunk-delay") {
		cliParams.ChunkDelay = chunkDelay
		cliParams.ChunkDelaySet = true
	}
	if commandFlagChanged(cmd, "chunk-size") {
		cliParams.ChunkSize = chunkSize
		cliParams.ChunkSizeSet = true
	}
	if commandFlagChanged(cmd, "cooldown-min-ms") {
		cliParams.CooldownMinMs = cooldownMinMs
		cliParams.CooldownMinMsSet = true
	}
	if commandFlagChanged(cmd, "cooldown-max-ms") {
		cliParams.CooldownMaxMs = cooldownMaxMs
		cliParams.CooldownMaxMsSet = true
	}
	if cmd != nil && cmd.Flags().Changed("port") {
		cliParams.Port = port
		cliParams.PortSet = true
	}
	if cmd != nil && cmd.Flags().Changed("host") {
		cliParams.Host = host
		cliParams.HostSet = true
	}
	return cliParams
}

func commandFlagChanged(cmd *cobra.Command, name string) bool {
	if cmd == nil {
		return false
	}
	return cmd.Flags().Changed(name) || cmd.InheritedFlags().Changed(name) || cmd.PersistentFlags().Changed(name)
}

func loadConfig(cmd *cobra.Command) (*config.Config, *config.ConfigSource) {
	cliParams := buildCLIParams(cmd)

	cfg, source, err := config.LoadWithOptions(cliParams, config.LoadAllOptions())
	if err != nil {
		log.Fatalln(err)
	}

	p115.SetOption(p115.Option{
		DisableCache:       cfg.P115.DisableCache,
		ChunkDelay:         cfg.P115.ChunkDelay,
		ChunkSize:          cfg.P115.ChunkSize,
		CooldownMinMs:      cfg.P115.CooldownMinMs,
		CooldownMaxMs:      cfg.P115.CooldownMaxMs,
		Retry:              cfg.Offline.Retry,
		PostProcess:        cfg.Rename.Enabled || cfg.Offline.DualFolder,
		OnePerEpisode:      cfg.Offline.OnePerEpisode,
		DatabasePath:       cfg.Database.Path,
		PreferredSubgroups: cfg.Lab.PreferredSubgroups,
	})
	if err := request.ReloadConfig(); err != nil {
		log.Printf("[warning] 加载站点配置失败: %v", err)
	}
	return cfg, source
}

func initAgent(cmd *cobra.Command) *config.Config {
	cfg, _ := loadConfig(cmd)

	var agentErr error
	if qrLogin {
		// -q: try existing cookies first, fall back to qrcode login
		pAgent, agentErr = p115.NewAgentByQrcode()
	} else if cfg.Auth.Cookies != "" {
		pAgent, agentErr = p115.NewAgent(cfg.Auth.Cookies)
	} else {
		pAgent, agentErr = p115.New()
	}
	if agentErr != nil {
		log.Fatalln(agentErr)
	}
	return cfg
}
