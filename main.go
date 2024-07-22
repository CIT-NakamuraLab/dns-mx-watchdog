package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/miekg/dns"
	"github.com/slack-go/slack"
)

func main() {
	token := os.Getenv("SLACK_BOT_TOKEN")
	channelId := os.Getenv("CHANNEL_ID")

	priDnsServer := os.Getenv("PRI_DNS_SERVER")
	secDnsServer := os.Getenv("SEC_DNS_SERVER")
	domain := os.Getenv("DOMAIN")

	client := slack.New(token)

	ticker := time.NewTicker(time.Millisecond * 1000 * 60 * 60)
	defer ticker.Stop()

	log.Printf("DNS Watchdog task has been started")

	count := 24
	execLookup(*client, channelId, priDnsServer, secDnsServer, domain, &count)
	for {
		select {
		case <-ticker.C:
			log.Printf("count=%d\n", count)
			execLookup(*client, channelId, priDnsServer, secDnsServer, domain, &count)
		}
	}
}

func lookupMXRecords(dnsServer string, domain string) bool {
	m := new(dns.Msg)
	m.SetQuestion(domain+".", dns.TypeMX)

	c := new(dns.Client)
	in, _, err := c.Exchange(m, dnsServer)
	if err != nil {
		fmt.Println("DNS Query Failed:", err)
		return false
	}

	for _, answer := range in.Answer {
		if mx, ok := answer.(*dns.MX); ok {
			fmt.Printf("Host: %s, Priority: %d\n", mx.Mx, mx.Preference)
		}
	}
	return true
}

func execLookup(client slack.Client, channelId string, priDnsServer string, secDnsServer string, domain string, count *int) {
	pri_ok := lookupMXRecords(priDnsServer, domain)
	sec_ok := lookupMXRecords(secDnsServer, domain)
	if pri_ok && sec_ok {
		*count++
		if *count > 24 {
			*count = 0
			SendDailyNotification(client, channelId, priDnsServer, secDnsServer, domain)
		}
	} else {
		*count = 0
		SendHourlyNotification(client, channelId, priDnsServer, secDnsServer, domain, pri_ok, sec_ok)
	}
}
