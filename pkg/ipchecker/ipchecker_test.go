package ipchecker

import (
	"fmt"
	"net"
	"testing"
	"time"

	"github.com/connormckelvey/cloudflare-spectrum-ddns/pkg/ipchecker/mocks"
	"github.com/golang/mock/gomock"
	"github.com/miekg/dns"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIPChecker(t *testing.T) {
	tests := []struct {
		DNSServer    string
		Domain       string
		ExpectedIP   string
		ExpoectedTTL uint32
	}{
		{"1.1.1.1", "example.com", "1.2.3.4", 1},
		{"8.8.8.8", "www.example.com", "1.2.3.5", 2},
	}

	for i, test := range tests {
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			m := mocks.NewMockDNSClient(ctrl)
			m.EXPECT().
				Exchange(gomock.Any(), gomock.Eq(test.DNSServer+":53")).
				DoAndReturn(func(msg *dns.Msg, serverAddr string) (*dns.Msg, time.Duration, error) {
					require.Equal(t, test.Domain+".", msg.Question[0].Name)
					msg.Answer = append(msg.Answer, &dns.A{
						Hdr: dns.RR_Header{Ttl: test.ExpoectedTTL},
						A:   net.ParseIP(test.ExpectedIP),
					})
					return msg, 0, nil
				})

			ipchecker := New(
				WithServerAddr(test.DNSServer),
				WithDNSClient(m),
			)

			ip, ttl, err := ipchecker.CheckIP(test.Domain)
			require.NoError(t, err)

			assert.Equal(t, test.ExpectedIP, ip)
			assert.EqualValues(t, test.ExpoectedTTL, ttl)
		})
	}
}
