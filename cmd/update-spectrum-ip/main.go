package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/cloudflare/cloudflare-go"
	"github.com/connormckelvey/cloudflare-spectrum-ddns/internal/config"
	"github.com/connormckelvey/cloudflare-spectrum-ddns/pkg/ipchecker"
	"github.com/connormckelvey/cloudflare-spectrum-ddns/pkg/ippoller"
	"github.com/connormckelvey/cloudflare-spectrum-ddns/pkg/spectrumutil"
	"github.com/sirupsen/logrus"
)

func main() {
	config := config.Load()
	err := config.Validate()
	if err != nil {
		panic(err)
	}

	logger := logrus.New()
	if config.Debug {
		logger.SetLevel(logrus.DebugLevel)
	}

	log := logger.WithField("process", "update-spectrum-ip")

	log.Debug(config)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-c
		cancel()
	}()

	fmt.Println(config.CloudflareAPIEmail, config.CloudflareAPIKey)
	api, err := cloudflare.New(config.CloudflareAPIKey, config.CloudflareAPIEmail)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(config.CloudflareAPIEmail, config.CloudflareAPIKey, config.CloudflareZoneName)

	zoneID, err := api.ZoneIDByName(config.CloudflareZoneName)
	if err != nil {
		log.Fatal(err)
	}

	checker := ipchecker.New(
		ipchecker.WithServerAddr(config.DNSServer),
		ipchecker.WithLogger(log),
	)

	poller := ippoller.New(
		ippoller.WithIPChecker(checker),
		ippoller.WithLogger(log),
	)

	client := spectrumutil.New(
		spectrumutil.WithCouldflare(api),
		spectrumutil.WithZoneID(zoneID),
		spectrumutil.WithLogger(log),
	)

	ips, errs := poller.Start(ctx, config.DDNSHostname)
	for {
		select {
		case ip, more := <-ips:
			if !more {
				log.Debug("exiting")
				return
			}
			_, err := client.Reconcile(ctx, config.SpectrumAppProtocol, config.SpectrumAppDomain, ip)
			if err != nil {
				log.Fatalf("reconcile failed: %v", err)
				continue
			}
			if !config.Polling {
				return
			}

		case err := <-errs:
			log.Fatalf("poller failed: %v", err)
		}
	}
}
