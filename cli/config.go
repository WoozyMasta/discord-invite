package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/jessevdk/go-flags"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/woozymasta/discord-invite/internal/vars"
	"golang.org/x/term"
)

// Config holds the configuration settings for the application.
type Config struct {
	ChannelID string `short:"c" long:"channel_id" description:"Discord channel ID" env:"DINVITE_CHANNEL_ID" required:"true"`
	BotToken  string `short:"t" long:"bot_token" description:"Discord bot token" env:"DINVITE_BOT_TOKEN" required:"true"`
	Listen    string `short:"l" long:"listen" description:"Address to listen on" env:"DINVITE_LISTEN" default:":8080"`
	LogLevel  string `long:"log-level" description:"Log level" env:"DINVITE_LOG_LEVEL" default:"info"`
	MaxAge    int    `short:"a" long:"max_age" description:"Invite max age in seconds" env:"DINVITE_MAX_AGE" default:"3600"`
	MaxUses   int    `short:"u" long:"max_uses" description:"Invite max uses" env:"DINVITE_MAX_USES" default:"1"`
	Unique    bool   `short:"x" long:"unique" description:"Make every invite unique" env:"DINVITE_UNIQUE"`
	Version   bool   `short:"v" long:"version" description:"Show version, commit, and build time."`
}

// setup initializes the configuration by parsing command-line flags and environment variables.
// It also configures the logging settings.
func setup() *Config {
	zerolog.TimeFieldFormat = time.RFC3339
	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	noColor := true
	if term.IsTerminal(int(os.Stdout.Fd())) {
		noColor = false
	}

	// Setting up log output
	log.Logger = log.Output(zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: time.RFC3339,
		NoColor:    noColor,
	})

	cfg := &Config{}
	p := flags.NewParser(cfg, flags.HelpFlag|flags.IgnoreUnknown|flags.AllowBoolValues)

	p.Usage = "[OPTIONS]"
	p.LongDescription = "Invite links generator for Discord channel."
	p.Command.Name = filepath.Base(p.Command.Name)

	_, err := p.Parse()
	if err != nil {
		// p.WriteHelp(os.Stdout)
		fmt.Fprintf(os.Stderr, "\n%s\n", err)
		os.Exit(1)
	}
	if cfg.Version {
		printVersion()
	}

	if logLevel, err := zerolog.ParseLevel(cfg.LogLevel); err != nil || cfg.LogLevel == "" {
		log.Warn().Msgf("Log level '%s' is unknown or empty, falling back to 'info' level", cfg.LogLevel)
		log.Logger = log.Level(zerolog.InfoLevel)
	} else {
		log.Logger = log.Level(logLevel)
	}

	return cfg
}

// print version information message and exit
func printVersion() {
	fmt.Printf(`file:     %s
version:  %s
commit:   %s
built:    %s
project:  %s
`, os.Args[0], vars.Version, vars.Commit, vars.BuildTime, vars.URL)
	os.Exit(0)
}
