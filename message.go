package main

import (
	"fmt"
	"log"
	"time"

	"github.com/slack-go/slack"
)

func SendDailyNotification(client slack.Client, channelId string, priDnsServer string, secDnsServer string, domain string) {
	attachment := slack.Attachment{
		Title: "DNS MX Query Successful (daily report)",
		Color: "good",

		Fields: []slack.AttachmentField{
			{
				Title: "Domain",
				Value: "*" + domain + "*",
				Short: false,
			}, {
				Title: "Primary DNS Server",
				Value: "*" + priDnsServer + "*",
				Short: false,
			}, {
				Title: "Secondary DNS Server",
				Value: "*" + secDnsServer + "*",
				Short: false,
			},
			{
				Title: "Date & Time",
				Value: time.Now().Format("2006/1/2 15:04"),
				Short: false,
			},
		},
	}

	_, _, err := client.PostMessage(
		channelId,
		slack.MsgOptionAttachments(attachment),
		slack.MsgOptionAsUser(true),
	)
	if err != nil {
		fmt.Printf("%s\n", err)
		return
	}
	log.Printf("[daily report] Message successfully sent to channel")
}

func SendHourlyNotification(client slack.Client, channelId string, priDnsServer string, secDnsServer string, domain string, pri_ok bool, sec_ok bool) {
	attachment := slack.Attachment{
		Pretext: "<!channel>",
		Title:   "DNS MX Query FAILED (hourly report)",
		Color:   "danger",

		Fields: []slack.AttachmentField{
			{
				Title: "Domain",
				Value: "*" + domain + "*",
				Short: false,
			}, {
				Title: "Primary DNS Server",
				Value: "*" + priDnsServer + "*",
				Short: false,
			}, {
				Title: "Secondary DNS Server",
				Value: "*" + secDnsServer + "*",
				Short: false,
			}, {
				Title: "Status",
				Value: fmt.Sprintf("*Primary: %t, Secondary: %t*", pri_ok, sec_ok),
				Short: false,
			},
			{
				Title: "Date & Time",
				Value: time.Now().Format("2006/1/2 15:04"),
				Short: false,
			},
		},
	}

	_, _, err := client.PostMessage(
		channelId,
		slack.MsgOptionAttachments(attachment),
		slack.MsgOptionAsUser(true),
	)
	if err != nil {
		fmt.Printf("%s\n", err)
		return
	}
	fmt.Printf("[hourly report] Message successfully sent to channel")
}
