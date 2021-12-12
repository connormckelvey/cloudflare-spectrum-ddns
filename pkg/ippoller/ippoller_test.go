package ippoller

import (
	"fmt"
	"net"
	"testing"
	"time"

	"github.com/connormckelvey/cloudflare-spectrum-ddns/pkg/ipchecker"
	"github.com/connormckelvey/cloudflare-spectrum-ddns/pkg/ipchecker/mocks"
	"github.com/golang/mock/gomock"
	"github.com/miekg/dns"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/net/context"
)

func TestIPPoller(t *testing.T) {
	tests := []struct {
		DNSServer     string
		Domain        string
		ExpectedLoops int
		ExpectedIP    string
		ExpoectedTTL  uint32
	}{
		{"1.1.1.1", "example.com", 2, "1.2.3.4", 1},
		{"8.8.8.8", "www.example.com", 2, "1.2.3.5", 1},
	}

	for i, test := range tests {
		t.Run(fmt.Sprint(i), func(t *testing.T) {

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			m := mocks.NewMockDNSClient(ctrl)
			m.EXPECT().
				Exchange(gomock.Any(), gomock.Eq(test.DNSServer+":53")).MinTimes(test.ExpectedLoops).
				DoAndReturn(func(msg *dns.Msg, serverAddr string) (*dns.Msg, time.Duration, error) {
					require.Equal(t, test.Domain+".", msg.Question[0].Name)
					msg.Answer = append(msg.Answer, &dns.A{
						Hdr: dns.RR_Header{Ttl: test.ExpoectedTTL},
						A:   net.ParseIP(test.ExpectedIP),
					})
					return msg, 0, nil
				})

			ipchecker := ipchecker.New(
				ipchecker.WithServerAddr(test.DNSServer),
				ipchecker.WithDNSClient(m),
			)
			ippoller := New(WithIPChecker(ipchecker))

			deadline := test.ExpectedLoops * int(test.ExpoectedTTL)
			ctx, cancel := context.WithTimeout(context.Background(), time.Duration(deadline)*time.Second)
			defer cancel()

			ips, errs := ippoller.Start(ctx, test.Domain)
			var actualLoops int

		Loop:
			for {
				select {
				case ip, more := <-ips:
					if !more {
						break Loop
					}
					assert.Equal(t, test.ExpectedIP, ip)
					actualLoops++

				case err, more := <-errs:
					if err != nil && more {
						assert.FailNow(t, err.Error())
					}
				}
			}

			assert.Equal(t, test.ExpectedLoops, actualLoops)
		})
	}
}
