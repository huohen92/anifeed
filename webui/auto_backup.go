package webui

import (
	"log"
	"strings"

	"github.com/huohen92/anifeed/config"
	"github.com/huohen92/anifeed/p115"
	"github.com/huohen92/anifeed/store"
)

// recordAutoBackupFailures counts consecutive main-RSS add failures from an
// automatic run and schedules backup discovery for the affected feed.
func (s *Server) recordAutoBackupFailures(cfg *config.Config, results []p115.FeedRunResult) {
	if cfg == nil || !cfg.Lab.AutoBackupSubgroup {
		return
	}
	byURL := make(map[string]p115.FeedRunResult, len(results))
	for _, result := range results {
		byURL[strings.TrimSpace(result.URL)] = result
	}
	threshold := cfg.Lab.AutoBackupFailures
	if threshold <= 0 {
		threshold = 2
	}
	for _, feeds := range cfg.RSS {
		for _, feed := range feeds {
			if !feed.IsEnabled() {
				continue
			}
			result, ok := byURL[strings.TrimSpace(feed.Url)]
			if !ok {
				continue
			}
			addFailed := strings.Contains(strings.ToLower(result.Error), "add offline task")
			s.backupMu.Lock()
			if s.backupFailures == nil {
				s.backupFailures = map[string]int{}
				s.backupTriggered = map[string]bool{}
			}
			if addFailed {
				s.backupFailures[feed.Url]++
				count := s.backupFailures[feed.Url]
				trigger := count >= threshold && !s.backupTriggered[feed.Url]
				if trigger {
					s.backupTriggered[feed.Url] = true
				}
				s.backupMu.Unlock()
				log.Printf("[auto-backup] 主 RSS 连续添加失败 %d/%d: %s", count, threshold, feed.Name)
				if trigger {
					go s.autoBackupSubgroups(cfg, feed.Url)
				}
			} else if result.Error == "" {
				s.backupFailures[feed.Url] = 0
				s.backupTriggered[feed.Url] = false
				s.backupMu.Unlock()
			} else {
				s.backupMu.Unlock()
			}
		}
	}
}

// autoBackupSubgroups adds a preferred backup RSS when the main feed is missing
// episodes and that subgroup already has one of those episodes.
func (s *Server) autoBackupSubgroups(cfg *config.Config, onlyURL string) {
	if cfg == nil || !cfg.Lab.AutoBackupSubgroup {
		return
	}
	st, err := store.NewWithPath(s.databasePath())
	if err != nil {
		return
	}
	defer st.Close()
	changed := false
	for site, feeds := range cfg.RSS {
		for i := range feeds {
			feed := &feeds[i]
			if !feed.AutoBackupEnabled() || (onlyURL != "" && feed.Url != onlyURL) || !feed.IsEnabled() || feed.TotalEpisodes <= 0 {
				continue
			}
			if cfg.Lab.AutoBackupSkipOffset && feed.RenameOffset != 0 {
				continue
			}
			source := config.NormalizeSource(feed.Source, feed.Url)
			sourceID := feedSourceID(source, feed.Url)
			if sourceID == "" || (source != config.SourceAniBT && source != config.SourceMikan) {
				continue
			}
			history, err := st.ListHistory(2000, feed.Name)
			if err != nil {
				continue
			}
			have := map[int]bool{}
			for _, item := range history {
				if episode, epErr := p115.ExtractEpisode(item.Title); epErr == nil && episode > 0 {
					have[int(episode)] = true
				}
			}
			missing := map[int]bool{}
			startEpisode := 0
			for episode := 1; episode <= feed.TotalEpisodes; episode++ {
				if !have[episode] {
					missing[episode] = true
					if startEpisode == 0 {
						startEpisode = episode
					}
				}
			}
			if len(missing) == 0 {
				continue
			}
			var groups *browseGroupsResponse
			if source == config.SourceAniBT {
				groups, err = s.browseAniBTGroups(sourceID)
			} else {
				groups, err = s.browseMikanGroups(sourceID)
			}
			if err != nil || groups == nil {
				continue
			}
			existing := map[string]bool{strings.TrimSpace(feed.Url): true}
			existingNames := map[string]bool{}
			for _, group := range feed.Subgroups {
				existing[strings.TrimSpace(group.URL)] = true
				existingNames[strings.ToLower(strings.TrimSpace(group.Name))] = true
			}
			if feed.SubgroupName != "" {
				existingNames[strings.ToLower(strings.TrimSpace(feed.SubgroupName))] = true
			}
			var additions []config.RssSubgroup
			for _, group := range groups.Groups {
				if strings.TrimSpace(group.RSS) == "" || existing[group.RSS] || existingNames[strings.ToLower(strings.TrimSpace(group.Name))] {
					continue
				}
				if !preferredBackupGroup(cfg.Lab.PreferredSubgroups, group.Name) {
					continue
				}
				for _, item := range group.Items {
					episode, epErr := p115.ExtractEpisode(item.Title)
					if epErr != nil || episode <= 0 || !missing[int(episode)] {
						continue
					}
					backupSeason := feed.RenameSeason
					backupOffset := 0
					mainOffset := feed.RenameOffset
					if mainOffset == 0 {
						mainOffset = cfg.Rename.Offset
					}
					mainNeedsOffset := mainOffset != 0 || detectSeason(feed.Name) > 1
					if (feed.AutoOffset || cfg.Lab.AutoOffset) && mainNeedsOffset {
						if suggestion, err := s.suggestOffset(cfg, feed, group.RSS); err == nil {
							backupSeason = suggestion.Season
							backupOffset = suggestion.Offset
							log.Printf("[auto-offset] %s / %s → TMDB S%02d, offset %+d (%d/%d)", feed.Name, group.Name, suggestion.Season, suggestion.Offset, suggestion.Samples, suggestion.Total)
						} else {
							log.Printf("[auto-offset] %s / %s 无法自动判定，保留 offset 0: %v", feed.Name, group.Name, err)
						}
					}
					additions = append(additions, config.RssSubgroup{
						Name: group.Name, URL: group.RSS, Enabled: ptrBool(true),
						Filter: "", Exclude: "",
						StartEpisode: startEpisode,
						Season:       backupSeason, Offset: backupOffset,
					})
					break
				}
			}
			if len(additions) == 0 {
				continue
			}
			feed.Subgroups = append(feed.Subgroups, additions...)
			cfg.RSS[site] = feeds
			changed = true
			log.Printf("[auto-backup] %s 添加 %d 个备用字幕组", feed.Name, len(additions))
		}
	}
	if !changed {
		return
	}
	resp, err := s.loadEditable()
	if err != nil {
		return
	}
	byURL := map[string]config.RssConfig{}
	for _, feeds := range cfg.RSS {
		for _, feed := range feeds {
			byURL[feed.Url] = feed
		}
	}
	for i := range resp.Config.RSS {
		if feed, ok := byURL[resp.Config.RSS[i].URL]; ok {
			for _, group := range feed.Subgroups {
				found := false
				for _, existing := range resp.Config.RSS[i].Subgroups {
					if existing.URL == group.URL {
						found = true
						break
					}
				}
				if !found {
					resp.Config.RSS[i].Subgroups = append(resp.Config.RSS[i].Subgroups, EditableSubgroup{
						Name: group.Name, URL: group.URL, Enabled: func() bool { return group.Enabled == nil || *group.Enabled }(),
						Filter: group.Filter, Exclude: group.Exclude, StartEpisode: group.StartEpisode, Season: group.Season, Offset: group.Offset,
					})
				}
			}
		}
	}
	if _, err := s.saveEditable(&resp.Config); err != nil {
		log.Printf("[auto-backup] 保存备用字幕组失败: %v", err)
	}
}

func preferredBackupGroup(preferred []string, name string) bool {
	if len(preferred) == 0 {
		return true
	}
	lower := strings.ToLower(name)
	for _, group := range preferred {
		if strings.Contains(lower, strings.ToLower(strings.TrimSpace(group))) {
			return true
		}
	}
	return false
}

func ptrBool(v bool) *bool { return &v }
