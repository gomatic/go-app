package app

import (
	slogx "github.com/skykernel/go-log"
	"github.com/urfave/cli/v3"
)

// LogLevelFlag returns the standard --log-level global flag, sourced from
// <envPrefix>LOG_LEVEL and bound to cfg.
func LogLevelFlag(cfg *slogx.LoggerConfig, envPrefix string) cli.Flag {
	return &cli.StringFlag{
		Name:        "log-level",
		Sources:     cli.EnvVars(envPrefix + "LOG_LEVEL"),
		Value:       "info",
		Usage:       "Logging level (debug, info, warn, error)",
		Destination: (*string)(&cfg.LogLevel),
	}
}

// LogFormatFlag returns the standard --log-format global flag, sourced from
// <envPrefix>LOG_FORMAT, bound to cfg, defaulting to def (consumers differ:
// a CLI defaults to text, a daemon to json).
func LogFormatFlag(cfg *slogx.LoggerConfig, envPrefix string, def slogx.LogFormat) cli.Flag {
	return &cli.StringFlag{
		Name:        "log-format",
		Sources:     cli.EnvVars(envPrefix + "LOG_FORMAT"),
		Value:       string(def),
		Usage:       "Log output format (text, json)",
		Destination: (*string)(&cfg.LogFormat),
	}
}

// OutputFlag returns the standard --output/-o global flag selecting the result
// encoding, sourced from <envPrefix>OUTPUT and defaulting to json. The action
// combinator reads it to encode each command's result.
func OutputFlag(envPrefix string) cli.Flag {
	return &cli.StringFlag{
		Name:    outputFlagName,
		Aliases: []string{"o"},
		Sources: cli.EnvVars(envPrefix + "OUTPUT"),
		Value:   "json",
		Usage:   "Result output format (json, yaml)",
	}
}
