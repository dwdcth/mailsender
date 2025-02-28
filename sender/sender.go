package sender

import (
	"context"
	"fmt"
	"log"
	"net/smtp"
	"time"

	"github.com/dwdcth/mailsender/g"
	"github.com/dwdcth/mailsender/proc"
	"github.com/jordan-wright/email"
	"golang.org/x/sync/semaphore"
)

var (
	mailChan chan *MailObject
	sem      *semaphore.Weighted
	pool     *email.Pool
)

func Start() {
	sem = semaphore.NewWeighted(int64(g.GetConfig().Mail.SendConcurrent))
	mailChan = make(chan *MailObject, g.GetConfig().Mail.MaxQueueSize)

	mcfg := g.GetConfig().Mail
	var err error
	auth := smtp.PlainAuth("", mcfg.MailServerAccount, mcfg.MailServerPasswd, mcfg.MailServerHost)
	pool, err = email.NewPool(
		fmt.Sprintf("%s:%d", mcfg.MailServerHost, mcfg.MailServerPort),
		mcfg.SendConcurrent,
		auth,
	)
	if err != nil {
		log.Fatalf("Failed to create email pool: %v", err)
	}

	go startSender()
}

// try pushing one mail into sender queue, maybe failed
func AddMail(r []string, subject string, content string, from ...string) bool {
	mcfg := g.GetConfig().Mail
	fromUserName := mcfg.FromUser
	if len(from) == 1 {
		fromUserName = from[0]
	}

	nm := NewMailObject(r, subject, content, fromUserName)
	select {
	case mailChan <- nm:
		return true
	default:
		return false
	}
}

// sender cron
func startSender() {
	ctx := context.Background()
	for mail := range mailChan {
		if err := sem.Acquire(ctx, 1); err != nil {
			log.Printf("Failed to acquire semaphore: %v", err)
			continue
		}
		go func(mailObject *MailObject) {
			defer sem.Release(1)
			sendMail(mailObject)
		}(mail)
	}
}

func sendMail(mo *MailObject) {
	mcfg := g.GetConfig().Mail
	e := &email.Email{
		To:      mo.Receivers,
		From:    fmt.Sprintf("%s <%s>", mo.FromUser, mcfg.MailServerAccount),
		Subject: mo.Subject,
		Text:    []byte(mo.Content),
	}

	// statistics
	proc.MailSendCnt.Incr()

	err := pool.Send(e, 10*time.Second)
	if err != nil {
		// statistics
		proc.MailSendErrCnt.Incr()
		log.Println(err, ", mailObject:", mo)
	} else {
		// statistics
		proc.MailSendOkCnt.Incr()
	}
}

// Mail Content Struct
type MailObject struct {
	Receivers []string
	Subject   string
	Content   string
	FromUser  string
}

func NewMailObject(receivers []string, subject string, content string, fromUserName string) *MailObject {
	return &MailObject{Receivers: receivers, Subject: subject, Content: content,
		FromUser: fromUserName}
}
