package webui

import "net/http"

type ruleOption struct {
	ID          string `json:"id"`
	Label       string `json:"label"`
	Pattern     string `json:"pattern"`
	Description string `json:"description,omitempty"`
}

type ruleCategory struct {
	ID      string       `json:"id"`
	Label   string       `json:"label"`
	Options []ruleOption `json:"options"`
}

type groupRule struct {
	ID          string `json:"id"`
	Label       string `json:"label"`
	Pattern     string `json:"pattern"`
	Description string `json:"description,omitempty"`
}

type rulesResponse struct {
	Source     string         `json:"source"`
	Updated    string         `json:"updated"`
	Categories []ruleCategory `json:"categories"`
	Groups     []groupRule    `json:"groups"`
}

const rulesSourceURL = "https://github.com/jacksen168sub/AB_AR_exclude_regular_expressions"

func (s *Server) handleRules(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, rulesResponse{
		Source:  rulesSourceURL,
		Updated: "2024.10.15",
		Categories: []ruleCategory{
			{
				ID:    "language",
				Label: "语言 / 字幕语言",
				Options: []ruleOption{
					{ID: "chs", Label: "简体", Pattern: `(?i)(?:简体|简中|简日|CHS|GB|SC)`, Description: "简中、简日、CHS、GB、SC"},
					{ID: "cht", Label: "繁体", Pattern: `(?i)(?:繁体|繁中|繁日|CHT|BIG5|TC)`, Description: "繁中、繁日、CHT、BIG5、TC"},
					{ID: "bilingual", Label: "简繁双语", Pattern: `(?i)(?:简繁|繁简|CHS&CHT|CHT&CHS|SC&TC|TC&SC|简日双语|繁日双语|简繁日)`, Description: "简繁内封、简日/繁日双语"},
				},
			},
			{
				ID:    "subtitle",
				Label: "字幕类型",
				Options: []ruleOption{
					{ID: "softsub", Label: "内封 / 软字幕", Pattern: `(?i)(?:内封|内挂|软字幕|字幕内封|CHS&CHT|SC&TC)`, Description: "字幕作为独立轨道封装"},
					{ID: "hardsub", Label: "内嵌 / 硬字幕", Pattern: `(?i)(?:内嵌|硬字幕|硬字)`, Description: "字幕烧录进画面"},
					{ID: "external", Label: "外挂字幕", Pattern: `(?i)(?:外挂|外置)`, Description: "字幕文件与视频分离"},
				},
			},
			{
				ID:    "container",
				Label: "容器格式",
				Options: []ruleOption{
					{ID: "mkv", Label: "MKV", Pattern: `(?i)(?:MKV|\.mkv)`, Description: "Matroska，常见内封多字幕"},
					{ID: "mp4", Label: "MP4", Pattern: `(?i)(?:MP4|\.mp4)`, Description: "MP4 容器"},
				},
			},
			{
				ID:    "resolution",
				Label: "分辨率",
				Options: []ruleOption{
					{ID: "r1080", Label: "1080p", Pattern: `(?i)(?:1080p|1920x1080)`, Description: "1080p / 1920x1080"},
					{ID: "r720", Label: "720p", Pattern: `(?i)(?:720p|1280x720)`, Description: "720p / 1280x720"},
					{ID: "r2160", Label: "4K / 2160p", Pattern: `(?i)(?:2160p|4K|3840x2160)`, Description: "2160p / 4K"},
				},
			},
		},
		Groups: []groupRule{
			{ID: "sakurato-cht-hard", Label: "桜都 · 繁体内嵌", Pattern: `(?i)^(?=.*\[桜都字幕组\])(?=.*(?:繁体内嵌|繁日内嵌)).*`, Description: "排除桜都繁体内嵌/繁日内嵌"},
			{ID: "moemoe-cht", Label: "喵萌奶茶屋 · 繁体", Pattern: `(?i)^(?=.*喵萌奶茶屋)(?=.*(?:繁日双语|繁体))(?!.*LoliHouse).*`, Description: "排除喵萌繁体，保留 LoliHouse 合作版"},
			{ID: "kitauji-cht-hard", Label: "北宇治 · 繁日内嵌", Pattern: `(?i)^(?=.*\[北宇治字幕组\])(?=.*(?:繁体内嵌|繁日内嵌)).*`, Description: "排除北宇治繁日内嵌"},
			{ID: "mingy-cht-hard", Label: "MingY · 繁日内嵌", Pattern: `(?i)^(?=.*\[MingY\])(?=.*(?:繁体内嵌|繁日内嵌)).*`, Description: "排除 MingY 繁日内嵌"},
			{ID: "bohuixunchun-cht-hard", Label: "拨雪寻春 · 繁日内嵌", Pattern: `(?i)^(?=.*\[❀拨雪寻春❀\])(?=.*繁日内嵌).*`, Description: "排除拨雪寻春繁日内嵌"},
			{ID: "manmao-cht", Label: "漫猫/猫恋 · 繁日双语", Pattern: `(?i)^(?=.*(?:\[爱恋&漫猫字幕组\]|\[漫猫字幕社&猫恋汉化组\]))(?=.*繁日双语).*`, Description: "排除漫猫/猫恋繁日双语"},
			{ID: "up21-not-bglobal", Label: "Up to 21°C · 非 B-Global", Pattern: `(?i)^(?=.*\[Up to 21\°C\])(?!.*\b(?:B-Global)\b).*`, Description: "排除 Up to 21°C 的 ABEMA/Baha/CR 源，保留 B-Global"},
			{ID: "yuha-cht-hard", Label: "悠哈璃羽 · 繁日内嵌", Pattern: `(?i)^(?=.*悠哈璃羽字幕社)(?=.*繁日内嵌).*`, Description: "排除悠哈璃羽繁日内嵌"},
		},
	})
}
