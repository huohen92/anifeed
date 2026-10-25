package config

// LabConfig holds experimental options. These settings are persisted before
// their full runtime behaviour is enabled.
type LabConfig struct {
	MultiSubgroup       bool
	MultiRssCoexist     bool
	WashEnabled         bool
	WashMode            string
	WashFolders         []string
	WashScheduleEnabled bool
	WashIntervalMinutes int
	PreferredSubgroups  []string
	InactiveSubgroups   []string
	AutoBackupSubgroup  bool
	AutoApplyRules      bool
	HideHints           bool
}

// DefaultLabConfig returns safe defaults for experimental features.
func DefaultLabConfig() LabConfig {
	return LabConfig{}
}
