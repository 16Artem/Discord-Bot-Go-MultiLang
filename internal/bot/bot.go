package bot

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"sync"
	"syscall"
	"time"

	"example.com/bot/internal/commands"
	"example.com/bot/internal/config"
	"example.com/bot/internal/locale"
	"example.com/bot/internal/services"
	"github.com/bwmarrin/discordgo"
)

type Bot struct {
	cfg        *config.Config
	dg         *discordgo.Session
	locale     *locale.Loader
	settings   map[string]string
	settingsMu sync.RWMutex
}

func NewBot(cfg *config.Config) (*Bot, error) {
	if cfg.Token == "" {
		cfg.Token = os.Getenv("DISCORD_TOKEN")
	}
	if cfg.Token == "" {
		return nil, fmt.Errorf("discord token not provided in config.yml or DISCORD_TOKEN env")
	}

	l, err := locale.NewLoader("internal/locale/languages")
	if err != nil {
		return nil, err
	}

	dg, err := discordgo.New("Bot " + cfg.Token)
	if err != nil {
		return nil, err
	}

	b := &Bot{cfg: cfg, dg: dg, locale: l, settings: make(map[string]string)}

	_ = b.loadSettings()

	dg.AddHandler(b.onInteractionCreate)

	dg.Identify.Intents = discordgo.IntentsGuilds

	var _ services.Service = (*Bot)(nil)

	return b, nil
}

func (b *Bot) Run() error {
	if err := b.dg.Open(); err != nil {
		return err
	}
	log.Println("Bot is running")

	if err := b.registerCommands(); err != nil {
		return err
	}

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	log.Println("shutdown signal received")

	b.unregisterCommands()
	b.dg.Close()
	time.Sleep(200 * time.Millisecond)
	return nil
}

func (b *Bot) registerCommands() error {
	appID := b.dg.State.User.ID
	_, err := b.dg.ApplicationCommandCreate(appID, "", &discordgo.ApplicationCommand{
		Name:        "help",
		Description: "Show help information",
	})
	if err != nil {
		return err
	}

	_, err = b.dg.ApplicationCommandCreate(appID, "", &discordgo.ApplicationCommand{
		Name:        "settings",
		Description: "Server settings",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionSubCommand,
				Name:        "lang",
				Description: "Set server language",
				Options: []*discordgo.ApplicationCommandOption{
					{
						Type:        discordgo.ApplicationCommandOptionString,
						Name:        "lang",
						Description: "Language code (ru or en)",
						Required:    true,
						Choices: []*discordgo.ApplicationCommandOptionChoice{
							{Name: "Русский", Value: "ru"},
							{Name: "English", Value: "en"},
						},
					},
				},
			},
		},
	})
	return err
}

func (b *Bot) unregisterCommands() {
	appID := b.dg.State.User.ID
	cmds, _ := b.dg.ApplicationCommands(appID, "")
	for _, c := range cmds {
		_ = b.dg.ApplicationCommandDelete(appID, "", c.ID)
	}
}

func (b *Bot) onInteractionCreate(s *discordgo.Session, ic *discordgo.InteractionCreate) {
	if ic.Type != discordgo.InteractionApplicationCommand {
		return
	}
	name := ic.ApplicationCommandData().Name
	switch name {
	case "help":
		cmd := commands.NewHelpCommand(b)
		cmd.Execute(s, ic, nil)
	case "settings":
		cmd := commands.NewSettingsCommand(b)
		cmd.Execute(s, ic, nil)
	default:
		_ = s.InteractionRespond(ic.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{Content: "Unknown command."},
		})
	}
}

func (b *Bot) settingsPath() string {
	return filepath.Join("data", "settings.json")
}

func (b *Bot) loadSettings() error {
	b.settingsMu.Lock()
	defer b.settingsMu.Unlock()
	p := b.settingsPath()
	if _, err := os.Stat(p); os.IsNotExist(err) {
		b.settings = make(map[string]string)
		return nil
	}
	bts, err := ioutil.ReadFile(p)
	if err != nil {
		return err
	}
	return json.Unmarshal(bts, &b.settings)
}

func (b *Bot) saveSettings() error {
	b.settingsMu.RLock()
	defer b.settingsMu.RUnlock()
	bts, err := json.MarshalIndent(b.settings, "", "  ")
	if err != nil {
		return err
	}
	dir := filepath.Dir(b.settingsPath())
	_ = os.MkdirAll(dir, 0o755)
	return ioutil.WriteFile(b.settingsPath(), bts, 0o644)
}

func (b *Bot) GetGuildLang(guildID string) string {
	b.settingsMu.RLock()
	defer b.settingsMu.RUnlock()

	if v, ok := b.settings[guildID]; ok {
		return v
	}

	if guildID != "" {
		guild, err := b.dg.Guild(guildID)
		if err == nil && guild.PreferredLocale != "" {
			switch guild.PreferredLocale {
			case "ru":
				return "ru"
			default:
				return "en"
			}
		}
	}

	return b.cfg.DefaultLang
}

func (b *Bot) SetGuildLang(guildID, lang string) error {
	b.settingsMu.Lock()
	b.settings[guildID] = lang
	b.settingsMu.Unlock()
	return b.saveSettings()
}

func (b *Bot) Config() *config.Config { return b.cfg }
func (b *Bot) Locale() *locale.Loader { return b.locale }
