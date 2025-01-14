package config

import (
	"github.com/nodeset-org/hyperdrive/modules/config"
	"github.com/nodeset-org/hyperdrive/shared/config/ids"
	"github.com/nodeset-org/hyperdrive/shared/logging"
)

// Configuration for the daemon loggers
type LoggingConfig struct {
	config.SectionHeader

	// The minimum record level that will be logged
	Level config.ChoiceParameter[logging.LogLevel]

	// The format to use when printing logs
	Format config.ChoiceParameter[logging.LogFormat]

	// True to include the source code position of the log statement in log messages
	AddSource config.BoolParameter

	// The maximum size (in megabytes) of the log file before it gets rotated
	MaxSize config.UintParameter

	// The maximum number of old log files to retain
	MaxBackups config.UintParameter

	// The maximum number of days to retain old log files based on the timestamp encoded in their filename
	MaxAge config.UintParameter

	// Toggle for saving rotated logs with local system time in the name vs. UTC
	LocalTime config.BoolParameter

	// Toggle for compressing rotated logs
	Compress config.BoolParameter
}

// Generates a new Logger configuration
func NewLoggingConfig() *LoggingConfig {
	cfg := &LoggingConfig{}
	cfg.SectionHeader.ID = config.Identifier(ids.LoggingSectionID)
	cfg.SectionHeader.Name = "Logging"
	cfg.SectionHeader.Description.Default = "Configure the logging options for the Hyperdrive sercive and any modules that support it."

	// Level
	cfg.Level.ID = config.Identifier(ids.LoggerLevelID)
	cfg.Level.Name = "Log Level"
	cfg.Level.Description.Default = "Select the minimum level for log messages. The lower it goes, the more verbose output the logs contain."
	cfg.Level.Default = logging.LogLevel_Info
	cfg.Level.AffectedContainers = []string{string(ContainerID_Daemon)}
	cfg.Level.Options = []config.ParameterOption[logging.LogLevel]{
		{
			Name: "Debug",
			Description: config.DynamicProperty[string]{
				Default: "Log debug messages - useful for development, or if something goes wrong and you need to provide extra information to supporters in order to track issues down.",
			},
			Value: logging.LogLevel_Debug,
		}, {
			Name: "Info",
			Description: config.DynamicProperty[string]{
				Default: "Log routine info messages.",
			},
			Value: logging.LogLevel_Info,
		}, {
			Name: "Warn",
			Description: config.DynamicProperty[string]{
				Default: "Only log warnings or higher, skipping info messages.",
			},
			Value: logging.LogLevel_Warn,
		}, {
			Name: "Error",
			Description: config.DynamicProperty[string]{
				Default: "Only log errors that prevent the daemon from running as expected.",
			},
			Value: logging.LogLevel_Error,
		},
	}

	// Format
	cfg.Format.ID = config.Identifier(ids.LoggerFormatID)
	cfg.Format.Name = "Format"
	cfg.Format.Description.Default = "Choose which format log messages will be printed in."
	cfg.Format.Default = logging.LogFormat_Logfmt
	cfg.Format.AffectedContainers = []string{string(ContainerID_Daemon)}
	cfg.Format.Options = []config.ParameterOption[logging.LogFormat]{
		{
			Name: "Logfmt",
			Description: config.DynamicProperty[string]{
				Default: "Use the logfmt format, which offers a good balance of human readability and parsability. See https://www.brandur.org/logfmt for more information on this format.",
			},
			Value: logging.LogFormat_Logfmt,
		}, {
			Name: "JSON",
			Description: config.DynamicProperty[string]{
				Default: "Log messages in JSON format. Useful if you want to process your logs through other tooling.",
			},
			Value: logging.LogFormat_Json,
		},
	}

	// AddSource
	cfg.AddSource.ID = config.Identifier(ids.LoggerAddSourceID)
	cfg.AddSource.Name = "Embed Source Location"
	cfg.AddSource.Description.Default = "Enable this to add the source location of where the logger was called to each log message. This is mostly for development use only."
	cfg.AddSource.Default = false
	cfg.AddSource.AffectedContainers = []string{string(ContainerID_Daemon)}

	// MaxSize
	cfg.MaxSize.NumberParameter.ID = config.Identifier(ids.LoggerMaxSizeID)
	cfg.MaxSize.NumberParameter.Name = "Max Log Size"
	cfg.MaxSize.NumberParameter.Description.Default = "The max size (in megabytes) of a log file before it gets rotated out and archived."
	cfg.MaxSize.NumberParameter.AffectedContainers = []string{string(ContainerID_Daemon)}
	cfg.MaxSize.Default = 20

	// MaxBackups
	cfg.MaxBackups.NumberParameter.ID = config.Identifier(ids.LoggerMaxBackupsID)
	cfg.MaxBackups.NumberParameter.Name = "Max Archived Logs"
	cfg.MaxBackups.NumberParameter.Description.Default = "The max number of archived logs to save before deleting old ones.\n\nUse 0 for no limit (preserve all archived logs)."
	cfg.MaxBackups.NumberParameter.AffectedContainers = []string{string(ContainerID_Daemon)}
	cfg.MaxBackups.Default = 3

	// MaxAge
	cfg.MaxAge.NumberParameter.ID = config.Identifier(ids.LoggerMaxAgeID)
	cfg.MaxAge.NumberParameter.Name = "Max Archive Age"
	cfg.MaxAge.NumberParameter.Description.Default = "The max number of days an archive log should be preserved for before being deleted.\n\nUse 0 for no limit (preserve all logs regardless of age)."
	cfg.MaxAge.NumberParameter.AffectedContainers = []string{string(ContainerID_Daemon)}
	cfg.MaxAge.Default = 90

	// LocalTime
	cfg.LocalTime.ID = config.Identifier(ids.LoggerLocalTimeID)
	cfg.LocalTime.Name = "Use Local Time"
	cfg.LocalTime.Description.Default = "When a log needs to be archived, by default the system will append the time of archiving to its filename in UTC. Enable this to use your local system's time in the filename instead."
	cfg.LocalTime.Default = false
	cfg.LocalTime.AffectedContainers = []string{string(ContainerID_Daemon)}

	// Compress
	cfg.Compress.ID = config.Identifier(ids.LoggerCompressID)
	cfg.Compress.Name = "Compress Archives"
	cfg.Compress.Description.Default = "Enable this to compress logs when they get archived to save space."
	cfg.Compress.Default = true
	cfg.Compress.AffectedContainers = []string{string(ContainerID_Daemon)}

	return cfg
}

// Get the title for the config
func (cfg *LoggingConfig) GetTitle() string {
	return "Logging"
}

// Get the parameters for this config
func (cfg *LoggingConfig) GetParameters() []config.IParameter {
	return []config.IParameter{
		&cfg.Level,
		&cfg.Format,
		&cfg.AddSource,
		&cfg.MaxSize,
		&cfg.MaxBackups,
		&cfg.MaxAge,
		&cfg.LocalTime,
		&cfg.Compress,
	}
}

// Get the sections underneath this one
func (cfg *LoggingConfig) GetSections() []config.ISection {
	return []config.ISection{}
}

// Convert the config into a LoggerOptions struct
func (cfg *LoggingConfig) GetOptions() logging.LoggerOptions {
	return logging.LoggerOptions{
		MaxSize:    int(cfg.MaxSize.Value),
		MaxBackups: int(cfg.MaxBackups.Value),
		MaxAge:     int(cfg.MaxAge.Value),
		LocalTime:  cfg.LocalTime.Value,
		Compress:   cfg.Compress.Value,
		Format:     cfg.Format.Value,
		Level:      cfg.Level.Value,
		AddSource:  cfg.AddSource.Value,
	}
}
