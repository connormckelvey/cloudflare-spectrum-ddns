package config

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"os"

	"github.com/connormckelvey/cloudflare-spectrum-ddns/pkg/spectrumutil"
	_ "github.com/joho/godotenv/autoload"
	"github.com/sirupsen/logrus"
)

const (
	CLOUDFLARE_API_KEY    = "CLOUDFLARE_API_KEY"
	CLOUDFLARE_API_EMAIL  = "CLOUDFLARE_API_EMAIL"
	CLOUDFLARE_ZONE_NAME  = "CLOUDFLARE_ZONE_NAME"
	SPECTRUM_APP_DOMAIN   = "SPECTRUM_APP_DOMAIN"
	SPECTRUM_APP_PROTOCOL = "SPECTRUM_APP_PROTOCOL"
	DDNS_HOSTNAME         = "DDNS_HOSTNAME"
	DNS_SERVER            = "DNS_SERVER"
)

type Config struct {
	CloudflareAPIKey    string
	CloudflareAPIEmail  string
	CloudflareZoneName  string
	SpectrumAppDomain   string
	SpectrumAppProtocol spectrumutil.SpectrumProtocol
	DDNSHostname        string
	DNSServer           string
	Polling             bool
	Debug               bool
}

var ErrFieldRequired = errors.New("missing required configuration")

func FieldRequiredError(env string) error {
	return fmt.Errorf("%w: '%v'", ErrFieldRequired, env)
}

func (c *Config) Validate() error {
	if c.CloudflareAPIKey == "" {
		return FieldRequiredError(CLOUDFLARE_API_KEY)
	}
	if c.CloudflareAPIEmail == "" {
		return FieldRequiredError(CLOUDFLARE_API_EMAIL)
	}
	if c.CloudflareZoneName == "" {
		return FieldRequiredError(CLOUDFLARE_ZONE_NAME)
	}
	if c.SpectrumAppDomain == "" {
		return FieldRequiredError(SPECTRUM_APP_DOMAIN)
	}
	if c.SpectrumAppProtocol == "" {
		return FieldRequiredError(SPECTRUM_APP_PROTOCOL)
	}
	if c.DDNSHostname == "" {
		return FieldRequiredError(DDNS_HOSTNAME)
	}
	if c.DNSServer == "" {
		return FieldRequiredError(DNS_SERVER)
	}
	return nil
}

func (c *Config) DebugEntry(log *logrus.Entry) *logrus.Entry {
	return log.WithFields(logrus.Fields{
		"cloudflare_api_key":    c.CloudflareAPIKey,
		"cloudflare_api_email":  c.CloudflareAPIEmail,
		"cloudflare_zone_name":  c.CloudflareZoneName,
		"spectrum_app_domain":   c.SpectrumAppDomain,
		"spectrum_app_protocol": c.SpectrumAppProtocol,
		"ddns_hostname":         c.DDNSHostname,
		"dns_server":            c.DNSServer,
		"polling":               c.Polling,
		"debug":                 c.Debug,
	})
}

func (c *Config) MarshalJSON() ([]byte, error) {
	redacted := *c
	redacted.CloudflareAPIKey = redacted.CloudflareAPIKey[:3] + "xxxxxxxxxx"
	redacted.CloudflareAPIEmail = redacted.CloudflareAPIEmail[:3] + "xxxxxxxxxx"
	return json.Marshal(redacted)
}

func (c *Config) String() string {
	configJSON, _ := json.MarshalIndent(c, "", "\t")
	return string(configJSON)
}

var (
	cloudflareZoneName  = flag.String("zone-name", os.Getenv(CLOUDFLARE_ZONE_NAME), "-zone example.com")
	spectrumAppDomain   = flag.String("app-domain", os.Getenv(SPECTRUM_APP_DOMAIN), "-app-domain minecraft.example.com")
	spectrumAppProtocol = flag.String("app-protocol", os.Getenv(SPECTRUM_APP_PROTOCOL), "-app-protocol ssh")
	ddnsHostname        = flag.String("ddns-hostname", os.Getenv(DDNS_HOSTNAME), "-ddns-hostname example.noip.com")
	dnsServer           = flag.String("dns-server", os.Getenv(DNS_SERVER), "-dns-server 8.8.8.8")
	poll                = flag.Bool("poll", false, "-poll")
	debug               = flag.Bool("debug", false, "-debug")
)

func Load() *Config {
	flag.Parse()
	return &Config{
		CloudflareAPIKey:    os.Getenv(CLOUDFLARE_API_KEY),
		CloudflareAPIEmail:  os.Getenv(CLOUDFLARE_API_EMAIL),
		CloudflareZoneName:  *cloudflareZoneName,
		SpectrumAppDomain:   *spectrumAppDomain,
		SpectrumAppProtocol: spectrumutil.GetProtocol(*spectrumAppProtocol),
		DDNSHostname:        *ddnsHostname,
		DNSServer:           *dnsServer,
		Polling:             *poll,
		Debug:               *debug,
	}
}
