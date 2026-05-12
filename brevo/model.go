package brevo

import brevo_official "github.com/getbrevo/brevo-go/lib"

// used to pass data to the function in order to send brevo mails
type Payload struct {
	To          []brevo_official.SendSmtpEmailTo
	BCC         []brevo_official.SendSmtpEmailBcc	
	TemplateID  int64
	Params      map[string]any
	Attachments []brevo_official.SendSmtpEmailAttachment
}

func (ep *Payload) AddAttachment(attachment brevo_official.SendSmtpEmailAttachment) {
	ep.Attachments = append(ep.Attachments, attachment)
}

type AddContactPayload struct {
	Email         string                 `json:"email"`
	ListIds       []int                  `json:"listIds"`
	UpdateEnabled bool                   `json:"updateEnabled"`
}

