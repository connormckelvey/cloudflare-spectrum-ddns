package ippoller

import (
	"context"
	"time"

	"github.com/sirupsen/logrus"
)

type IPChecker interface {
	CheckIP(domain string) (ip string, ttl uint32, err error)
}

type IPPoller struct {
	checker IPChecker
	log     *logrus.Entry
}

type Option func(*IPPoller)

func WithIPChecker(checker IPChecker) Option {
	return func(p *IPPoller) {
		p.checker = checker
	}
}

func WithLogger(logger *logrus.Entry) Option {
	return func(p *IPPoller) {
		p.log = logger
	}
}

func New(opts ...Option) *IPPoller {
	p := &IPPoller{
		log: logrus.NewEntry(logrus.New()),
	}
	for _, apply := range opts {
		apply(p)
	}
	p.log = p.log.WithFields(logrus.Fields{
		"service": "ip_poller",
	})
	return p
}

func (p *IPPoller) Start(ctx context.Context, domain string) (chan string, chan error) {
	log := p.log.WithFields(logrus.Fields{
		"method": "start",
		"domain": domain,
	})

	ips := make(chan string)
	errs := make(chan error)

	go func() {
		defer close(ips)
		defer close(errs)

		lastTTL := uint32(60)
		for {
			ip, ttl, err := p.checker.CheckIP(domain)
			if err != nil {
				errs <- err
				time.Sleep(time.Duration(lastTTL) * time.Second)
				continue
			}

			log.Debugf("ip is %v", ip)
			log.Debugf("ttl is %v", ttl)

			lastTTL = ttl
			ips <- ip

			log.Debugf("waiting for %v seconds", ttl)
			select {
			case <-ctx.Done():
				return
			case <-time.After(time.Duration(ttl) * time.Second):
			}
		}
	}()
	return ips, errs
}
