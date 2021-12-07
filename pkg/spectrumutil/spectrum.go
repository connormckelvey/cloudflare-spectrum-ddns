package spectrumutil

import (
	"context"
	"errors"
	"fmt"

	"github.com/cloudflare/cloudflare-go"
	"github.com/sirupsen/logrus"
)

var SpectrumProtocols = map[SpectrumProtocol]int{
	SpectrumProtocolMinecraft: 25565,
	SpectrumProtocolSSH:       22,
}

type SpectrumProtocol string

func (p SpectrumProtocol) OriginDirect(ip string) string {
	return fmt.Sprintf("tcp://%s:%d", ip, SpectrumProtocols[p])
}

const (
	SpectrumProtocolMinecraft SpectrumProtocol = "minecraft"
	SpectrumProtocolSSH       SpectrumProtocol = "ssh"
)

func GetProtocol(v string) SpectrumProtocol {
	protocol := SpectrumProtocolSSH
	if v == "minecraft" {
		protocol = SpectrumProtocolMinecraft
	}
	return protocol
}

type SpectrumClient struct {
	cf     *cloudflare.API
	zoneID string
	log    *logrus.Entry
}

type Option func(*SpectrumClient)

func WithCouldflare(cf *cloudflare.API) Option {
	return func(s *SpectrumClient) {
		s.cf = cf
	}
}

func WithZoneID(zoneID string) Option {
	return func(s *SpectrumClient) {
		s.zoneID = zoneID
	}
}

func WithLogger(logger *logrus.Entry) Option {
	return func(s *SpectrumClient) {
		s.log = logger
	}
}

func New(opts ...Option) *SpectrumClient {
	s := &SpectrumClient{
		log: logrus.NewEntry(logrus.New()),
	}
	for _, apply := range opts {
		apply(s)
	}
	s.log = s.log.WithFields(logrus.Fields{
		"service": "spectrum_client",
		"zone_id": s.zoneID,
	})
	return s
}

var ErrNoAppWithDomain = errors.New("no app with domain")

func (s *SpectrumClient) AppByDomain(ctx context.Context, domain string) (*cloudflare.SpectrumApplication, error) {
	log := s.log.WithFields(logrus.Fields{
		"method": "app_by_domain",
		"domain": domain,
	})

	log.Debug("looking up spectrum application by DNS name")
	apps, err := s.cf.SpectrumApplications(ctx, s.zoneID)
	if err != nil {
		return nil, err
	}
	for _, app := range apps {
		if app.DNS.Name == domain {
			return &app, nil
		}
	}

	log.Debug("no spectrum app with matching domain")
	return nil, ErrNoAppWithDomain
}

func (s *SpectrumClient) Reconcile(ctx context.Context, protocol SpectrumProtocol, domain string, ip string) (*cloudflare.SpectrumApplication, error) {
	log := s.log.WithFields(logrus.Fields{
		"method":   "reconcile",
		"domain":   domain,
		"protocol": protocol,
		"ip":       ip,
	})

	existsingApp, err := s.AppByDomain(ctx, domain)
	if err != nil {
		if !errors.Is(err, ErrNoAppWithDomain) {
			return nil, err
		}
		log.Debug("no existing spectrum app detected")
		return s.CreateApp(ctx, protocol, domain, ip)
	}
	return s.UpdateAppIP(ctx, existsingApp, protocol, ip)
}

func (s *SpectrumClient) CreateApp(ctx context.Context, protocol SpectrumProtocol, domain string, ip string) (*cloudflare.SpectrumApplication, error) {
	log := s.log.WithFields(logrus.Fields{
		"method":   "create_app",
		"protocol": protocol,
		"domain":   domain,
		"ip":       ip,
	})

	log.Debug("creating new spectrum app")
	newApp, err := s.cf.CreateSpectrumApplication(ctx, s.zoneID, cloudflare.SpectrumApplication{
		Protocol: string(protocol),
		DNS: cloudflare.SpectrumApplicationDNS{
			Type: "CNAME",
			Name: domain,
		},
		OriginDirect: []string{protocol.OriginDirect(ip)},
	})
	if err != nil {
		return nil, err
	}
	return &newApp, nil
}

func (s *SpectrumClient) UpdateAppIP(ctx context.Context, app *cloudflare.SpectrumApplication, protocol SpectrumProtocol, ip string) (*cloudflare.SpectrumApplication, error) {
	log := s.log.WithFields(logrus.Fields{
		"method":   "update_app",
		"app_id":   app.ID,
		"protocol": protocol,
		"ip":       ip,
	})

	if len(app.OriginDirect) == 1 && app.OriginDirect[0] == protocol.OriginDirect(ip) {
		log.Info("spectrum app ip is already up-to-date, nothing to change")
		return app, nil
	}

	update := *app
	update.ID = ""
	update.CreatedOn = nil
	update.ModifiedOn = nil
	update.OriginDirect = []string{protocol.OriginDirect(ip)}

	log.Info("ip change detected, updating spectum application")
	updatedApp, err := s.cf.UpdateSpectrumApplication(ctx, s.zoneID, app.ID, update)
	if err != nil {
		return nil, err
	}

	log.Info("ip address updated")
	return &updatedApp, nil
}
