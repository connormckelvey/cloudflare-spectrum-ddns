package ipchecker

import (
	"time"

	"github.com/miekg/dns"
)

type DNSClient interface {
	Exchange(*dns.Msg, string) (*dns.Msg, time.Duration, error)
}
