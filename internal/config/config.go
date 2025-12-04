package config

import "os"

type Config struct {
	AppName    string
	AppEnv     string
	AppURL     string
	AppTagline string

	UmamiWebsiteID  string
	UmamiHost       string
	PlausibleDomain string
	PlausibleHost   string

	ResendAPIKey     string
	ResendAudienceID string
}

func Load() *Config {
	return &Config{
		AppName:    envString("APP_NAME", "Axel Adrian"),
		AppEnv:     envString("APP_ENV", "development"),
		AppURL:     envString("APP_URL", "http://localhost:8090"),
		AppTagline: envString("APP_TAGLINE", "I do stuff."),

		UmamiWebsiteID:  os.Getenv("UMAMI_WEBSITE_ID"),
		UmamiHost:       envString("UMAMI_HOST", "cloud.umami.is"),
		PlausibleDomain: os.Getenv("PLAUSIBLE_DOMAIN"),
		PlausibleHost:   envString("PLAUSIBLE_HOST", "plausible.io"),

		ResendAPIKey:     os.Getenv("RESEND_API_KEY"),
		ResendAudienceID: os.Getenv("RESEND_AUDIENCE_ID"),
	}
}

func (c *Config) IsDev() bool {
	return c.AppEnv != "production"
}

func envString(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}
