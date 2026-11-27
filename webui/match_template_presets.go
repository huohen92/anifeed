package webui

// presetMatchTemplates is derived from the public rule collection at
// https://github.com/jacksen168sub/AB_AR_exclude_regular_expressions and the
// ANI-RSS group-scoped rule model. Rules are intentionally unscoped here:
// when a template is applied, the add-subscription page prepends the current
// {{subgroup}} scope automatically.
func presetMatchTemplates() []SubgroupMatchTemplate {
	return []SubgroupMatchTemplate{
		// 桜都字幕组
		presetTemplate("sakurato-chs-hard", "桜都字幕组", "1080P 简 内嵌", []string{"1080P", "简", "内嵌"}, "1080[Pp]\n(?:简体内嵌|简日内嵌)", ""),
		presetTemplate("sakurato-cht-hard", "桜都字幕组", "1080P 繁 内嵌", []string{"1080P", "繁", "内嵌"}, "1080[Pp]\n(?:繁体内嵌|繁日内嵌)", ""),
		presetTemplate("sakurato-soft", "桜都字幕组", "1080P 简繁 内封", []string{"1080P", "简繁", "内封"}, "1080[Pp]\n内封", ""),

		// 今晚月色真美（AniBT 1080P WebRip AVC 内嵌）
		presetTemplate("tsukigakirei-chs-hard", "今晚月色真美", "1080P 简日 内嵌", []string{"1080P", "简日", "内嵌"}, "1080[Pp]\n简日内嵌", "", "Tsukigakirei", "tsukigakirei"),
		presetTemplate("tsukigakirei-cht-hard", "今晚月色真美", "1080P 繁日 内嵌", []string{"1080P", "繁日", "内嵌"}, "1080[Pp]\n繁日內嵌", "", "Tsukigakirei", "tsukigakirei"),

		// 喵萌奶茶屋
		presetTemplate("moemoe-soft", "喵萌奶茶屋", "1080P 简繁 内封", []string{"1080P", "简繁日", "内封"}, "1080[Pp]\n简繁日?内封", ""),
		presetTemplate("moemoe-chs-jp", "喵萌奶茶屋", "1080P 简日 双语", []string{"1080P", "简日", "双语"}, "1080[Pp]\n简日双语", ""),
		presetTemplate("moemoe-cht-jp", "喵萌奶茶屋", "1080P 繁日 双语", []string{"1080P", "繁日", "双语"}, "1080[Pp]\n繁日双语", ""),

		// 北宇治字幕组
		presetTemplate("kitauji-chs-hard", "北宇治字幕组", "1080P 简 内嵌", []string{"1080P", "简日", "内嵌"}, "1080[Pp]\n简日内嵌", ""),
		presetTemplate("kitauji-cht-hard", "北宇治字幕组", "1080P 繁 内嵌", []string{"1080P", "繁日", "内嵌"}, "1080[Pp]\n繁日内嵌", ""),
		presetTemplate("kitauji-soft", "北宇治字幕组", "1080P 简繁 内封", []string{"1080P", "简繁日", "内封"}, "1080[Pp]\n简繁日内封", ""),

		// MingY
		presetTemplate("mingy-chs-hard", "MingY", "1080P 简 内嵌", []string{"1080P", "简日", "内嵌"}, "1080[Pp]\n简日内嵌", ""),
		presetTemplate("mingy-cht-hard", "MingY", "1080P 繁 内嵌", []string{"1080P", "繁日", "内嵌"}, "1080[Pp]\n繁日内嵌", ""),
		presetTemplate("mingy-soft", "MingY", "1080P 简繁 内封", []string{"1080P", "简繁日", "内封"}, "1080[Pp]\n简繁日内封", ""),

		// ❀拨雪寻春❀
		presetTemplate("bohui-chs-hard", "拨雪寻春", "1080P 简 内嵌", []string{"1080P", "简日", "内嵌"}, "1080[Pp]\n简日内嵌", ""),
		presetTemplate("bohui-cht-hard", "拨雪寻春", "1080P 繁 内嵌", []string{"1080P", "繁日", "内嵌"}, "1080[Pp]\n繁日内嵌", ""),
		presetTemplate("bohui-soft", "拨雪寻春", "1080P 简繁 内封", []string{"1080P", "简繁日", "内封"}, "1080[Pp]\n简繁日内封", ""),

		// 悠哈璃羽字幕社
		presetTemplate("yuha-chs-hard", "悠哈璃羽字幕社", "1080P 简 内嵌", []string{"1080P", "简日", "内嵌"}, "1080[Pp]\n简日内嵌", ""),
		presetTemplate("yuha-cht-hard", "悠哈璃羽字幕社", "1080P 繁 内嵌", []string{"1080P", "繁日", "内嵌"}, "1080[Pp]\n繁日内嵌", ""),

		// 漫猫 / 猫恋
		presetTemplate("manmao-chs-jp", "爱恋&漫猫字幕组", "1080P 简日 双语", []string{"1080P", "简日", "双语"}, "1080[Pp]\n简日双语", ""),
		presetTemplate("manmao-cht-jp", "爱恋&漫猫字幕组", "1080P 繁日 双语", []string{"1080P", "繁日", "双语"}, "1080[Pp]\n繁日双语", ""),
		presetTemplate("manmao2-chs-jp", "漫猫字幕社&猫恋汉化组", "1080P 简日 双语", []string{"1080P", "简日", "双语"}, "1080[Pp]\n简日双语", ""),
		presetTemplate("manmao2-cht-jp", "漫猫字幕社&猫恋汉化组", "1080P 繁日 双语", []string{"1080P", "繁日", "双语"}, "1080[Pp]\n繁日双语", ""),

		// LoliHouse / ANi 常用内封
		presetTemplate("lolihouse-soft", "LoliHouse", "1080P 简繁 内封", []string{"1080P", "简繁", "内封"}, "1080[Pp]\n内封", ""),
		presetTemplate("ani-soft", "ANi", "1080P 简繁 内封", []string{"1080P", "简繁", "内封"}, "1080[Pp]\n内封", ""),

		// Up to 21°C（B-Global 4K / 非 B-Global 1080p）
		presetTemplate("up21-not-bglobal", "Up to 21°C", "1080P 非 B-Global", []string{"1080P", "ABEMA", "Baha", "CR"}, "1080[Pp]", "B-Global"),
		presetTemplate("up21-bglobal-4k", "Up to 21°C", "4K B-Global", []string{"4K", "B-Global"}, "(?:2160[Pp]|4K|3840x2160)", ""),
	}
}

func presetTemplate(id, subgroup, name string, tags []string, filter, exclude string, aliases ...string) SubgroupMatchTemplate {
	return SubgroupMatchTemplate{ID: "preset-" + id, Subgroup: subgroup, Name: name, Tags: tags, Aliases: aliases, Filter: filter, Exclude: exclude}
}
