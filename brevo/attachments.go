package brevo

import brevo_official "github.com/getbrevo/brevo-go/lib"

func CreateBrevoAttachment(url, content, name string) brevo_official.SendSmtpEmailAttachment {
	return brevo_official.SendSmtpEmailAttachment{
		Content: content,
		Name:    name,
		Url:     url,
	}
}