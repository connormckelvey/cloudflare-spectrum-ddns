package config

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"

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
	DDNS_DOMAIN           = "DDNS_DOMAIN"
	DNS_SERVER            = "DNS_SERVER"
)

type Config struct {
	CloudflareAPIKey    string
	CloudflareAPIEmail  string
	CloudflareZoneName  string
	SpectrumAppDomain   string
	SpectrumAppProtocol spectrumutil.SpectrumProtocol
	DDNSDomain          string
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
	if c.DDNSDomain == "" {
		return FieldRequiredError(DDNS_DOMAIN)
	}
	if c.DNSServer == "" {
		return FieldRequiredError(DNS_SERVER)
	}
	return nil
}

func (c *Config) DebugEntry(log *logrus.Entry) *logrus.Entry {
	return log.WithFields(logrus.Fields{
		"cloudflare_api_key":    redact(c.CloudflareAPIKey),
		"cloudflare_api_email":  redact(c.CloudflareAPIEmail),
		"cloudflare_zone_name":  c.CloudflareZoneName,
		"spectrum_app_domain":   c.SpectrumAppDomain,
		"spectrum_app_protocol": c.SpectrumAppProtocol,
		"ddns_domain":           c.DDNSDomain,
		"dns_server":            c.DNSServer,
		"polling":               c.Polling,
		"debug":                 c.Debug,
	})
}

var (
	cloudflareZoneName  = flag.String("zone-name", os.Getenv(CLOUDFLARE_ZONE_NAME), "-zone example.com")
	spectrumAppDomain   = flag.String("app-domain", os.Getenv(SPECTRUM_APP_DOMAIN), "-app-domain minecraft.example.com")
	spectrumAppProtocol = flag.String("app-protocol", os.Getenv(SPECTRUM_APP_PROTOCOL), "-app-protocol ssh")
	ddnsDomain          = flag.String("ddns-domain", os.Getenv(DDNS_DOMAIN), "-ddns-domain example.noip.com")
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
		DDNSDomain:          *ddnsDomain,
		DNSServer:           *dnsServer,
		Polling:             *poll,
		Debug:               *debug,
	}
}

func redact(v string) string {
	b := strings.Builder{}
	for i, ch := range v {
		if i < 4 {
			b.WriteRune(ch)
		} else {
			b.WriteRune('x')
		}
	}
	return b.String()
}
