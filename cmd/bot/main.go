package main

import (
    "log"
    "example.com/bot/internal/bot"
    "example.com/bot/internal/config"
)

func main() {
    cfg, err := config.LoadConfig("config.yml")
    if err != nil {
        log.Fatalf("failed to load config: %v", err)
    }

    b, err := bot.NewBot(cfg)
    if err != nil {
        log.Fatalf("failed to initialize bot: %v", err)
    }

    if err := b.Run(); err != nil {
        log.Fatalf("bot returned error: %v", err)
    }
}
