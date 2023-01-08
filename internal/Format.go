package internal

import (
	"fmt"
	"sort"
	"strings"

	"github.com/emersion/go-imap"
)

func SortMessages(messages []*MyMessage) {
	sort.Slice(messages, func(i, j int) bool {
		m1 := messages[i]
		m2 := messages[j]

		if m1.Envelope.Date.Before(m2.Envelope.Date) {
			return false
		}
		if m1.Envelope.Date.After(m2.Envelope.Date) {
			return true
		}
		return false
	})
}

func SortFolders(folders []string) {
	sort.Slice(folders, func(i, j int) bool {
		m1 := folders[i]
		m2 := folders[j]

		if m1 < m2 {
			return true
		}
		if m2 < m1 {
			return false
		}
		return false
	})
}

func MessageToString(idx int, message *MyMessage) string {
	return fmt.Sprintf("%d %s %s %s", idx,
		message.Envelope.Date.Local().Format("Mon, 02 Jan 2006 15:04"),
		max_len(short_describe_email(message.Envelope.From[0]), 25),
		max_len(message.Envelope.Subject, 50))
}

func MessageSummary(message *MyMessage) string {
	res := make([]string, 0)
	for _, addr := range message.Envelope.From {
		res = append(res, "from: "+LongDescribeEmail(addr))
	}
	for _, addr := range message.Envelope.To {
		res = append(res, "to: "+LongDescribeEmail(addr))
	}
	for _, addr := range message.Envelope.Cc {
		res = append(res, "cc: "+LongDescribeEmail(addr))
	}
	for _, addr := range message.Envelope.Bcc {
		res = append(res, "bcc: "+LongDescribeEmail(addr))
	}
	for _, addr := range message.Envelope.ReplyTo {
		res = append(res, "replyto: "+LongDescribeEmail(addr))
	}
	res = append(res, "subject: "+message.Envelope.Subject)

	return strings.Join(res, "\n")
}

func GetIndex[K comparable](data []K, item K, otherwise int) int {
	for idx, it := range data {
		if it == item {
			return idx
		}
	}
	return otherwise
}

func max_len(str string, max_len int) string {
	if len(str) > max_len {
		str = str[:max_len]
	}
	return str
}
func short_describe_emails(addrs []*imap.Address) string {
	res := ""
	for _, addr := range addrs {
		if res != "" {
			res = res + ","
		}
		res = res + short_describe_email(addr)
	}
	return res
}
func short_describe_email(addr *imap.Address) string {
	if len(addr.PersonalName) > 4 {
		return max_len(addr.PersonalName, 20)
	}
	return fmt.Sprintf("%s@%s", addr.MailboxName, addr.HostName)
}
func LongDescribeEmail(addr *imap.Address) string {
	return fmt.Sprintf("\"%s\" <%s@%s>", addr.PersonalName, addr.MailboxName, addr.HostName)
}

func short_describe_message(idx int, m *MyMessage) string {
	return fmt.Sprintf(
		"|%03d/%s/%s/%s/%s|",
		idx,
		max_len(short_describe_emails(m.Envelope.From), 15),
		max_len(m.Envelope.Subject, 20),
		max_len(short_describe_emails(m.Envelope.Cc), 15),
		max_len(short_describe_emails(m.Envelope.ReplyTo), 15),
	)
}
