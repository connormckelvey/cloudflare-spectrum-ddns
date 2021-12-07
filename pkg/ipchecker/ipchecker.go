package ipchecker

import (
	"errors"

	"github.com/miekg/dns"
	"github.com/sirupsen/logrus"
)

type IPChecker struct {
	serverAddr string
	client     *dns.Client
	log        *logrus.Entry
}

type Option func(*IPChecker)

func WithServerAddr(serverAddr string) Option {
	return func(c *IPChecker) {
		c.serverAddr = serverAddr
	}
}

func WithLogger(logger *logrus.Entry) Option {
	return func(c *IPChecker) {
		c.log = logger
	}
}

func New(opts ...Option) *IPChecker {
	c := &IPChecker{
		log:    logrus.NewEntry(logrus.New()),
		client: &dns.Client{},
	}
	for _, apply := range opts {
		apply(c)
	}

	c.log = c.log.WithFields(logrus.Fields{
		"service": "ip_checker",
		"server":  c.serverAddr,
	})

	return c
}

func (c *IPChecker) CheckIP(domain string) (ip string, ttl uint32, err error) {
	log := c.log.WithFields(logrus.Fields{
		"method": "check_ip",
		"domain": domain,
	})

	log.Debug("querying dns records for domain")
	message := new(dns.Msg).SetQuestion(domain+".", dns.TypeA)
	r, _, err := c.client.Exchange(message, c.serverAddr+":53")
	if err != nil {
		return "", 0, err
	}
	if len(r.Answer) == 0 {
		return "", 0, errors.New("no answers returned")
	}
	for _, ans := range r.Answer {
		if a, ok := ans.(*dns.A); ok {
			ip = a.A.String()
			ttl = ans.Header().Ttl
		}
	}
	log.Debug("received dns answer record for domain")
	return ip, ttl, nil
}
