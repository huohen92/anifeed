package p115

import (
	"fmt"
	"time"

	"github.com/deadblue/elevengo"
	"github.com/deadblue/elevengo/option"
)

// QrcodeLogin drives a 115 QR-code login session for the Web UI.
type WebQrcodeLogin struct {
	agent   *elevengo.Agent
	session *elevengo.QrcodeSession
	started time.Time
}

// StartQrcodeLogin requests a new QR code from 115.
func StartQrcodeLogin() (login *WebQrcodeLogin, err error) {
	defer func() {
		if recovered := recover(); recovered != nil {
			login = nil
			err = fmt.Errorf("获取登录二维码异常: %v", recovered)
		}
	}()
	agent := newElevengoAgent()
	if agent == nil {
		return nil, fmt.Errorf("初始化 115 客户端失败")
	}
	session := &elevengo.QrcodeSession{}
	if err := agent.QrcodeStart(session, option.Qrcode().LoginTv()); err != nil {
		return nil, fmt.Errorf("start qrcode login: %w", err)
	}
	return &WebQrcodeLogin{agent: agent, session: session, started: time.Now()}, nil
}

// Image returns the raw PNG/JPEG image bytes of the QR code.
func (q *WebQrcodeLogin) Image() []byte {
	return q.session.Image
}

// Poll checks the QR code state. done is true once the user has scanned and
// confirmed the login; Cookies can then be persisted.
func (q *WebQrcodeLogin) Poll() (done bool, err error) {
	done, err = q.agent.QrcodePoll(q.session)
	if err == elevengo.ErrQrcodeCancelled {
		return false, err
	}
	return done, err
}

// Cookies exports the credentials as a "UID=...; CID=...; SEID=...; KID=..."
// string after a successful login.
func (q *WebQrcodeLogin) Cookies() string {
	cr := &elevengo.Credential{}
	q.agent.CredentialExport(cr)
	return fmt.Sprintf("UID=%s; CID=%s; SEID=%s; KID=%s", cr.UID, cr.CID, cr.SEID, cr.KID)
}

// Expired reports whether the session has outlived the 115 QR code lifetime.
func (q *WebQrcodeLogin) Expired() bool {
	return time.Since(q.started) > 3*time.Minute
}
