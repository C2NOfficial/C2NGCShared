package brevo

import (
	"context"
	"io"
	"log"
	"net/http"

	"github.com/C2NOfficial/C2NGCShared/constants"
	gcp_shared "github.com/C2NOfficial/C2NGCShared/gcp"
	brevo_official "github.com/getbrevo/brevo-go/lib"
)

func SendEmail(ctx context.Context, ep *Payload) {

	sendSmtpEmail := brevo_official.SendSmtpEmail{
		To:         ep.To,
		Bcc:        ep.BCC,
		TemplateId: ep.TemplateID,
		Params:     ep.Params,
		Sender: &brevo_official.SendSmtpEmailSender{
			Name:  constants.CompanyTradeName,
			Email: constants.CompanyInfoEmail,
		},
		Attachment: ep.Attachments,
	}

	_, resp, err := Client.TransactionalEmailsApi.SendTransacEmail(ctx, sendSmtpEmail)
	if err != nil {
		gcp_shared.LogError("BREVO MAIL ERROR", "Error occurred while sending the brevo email")
	}
	if resp.StatusCode != http.StatusCreated { // brevo server returns 201 on success
		respBodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			gcp_shared.LogError("BREVO MAIL ERROR", "Error occurred while sending the brevo email. Response: "+string(respBodyBytes))
		}
		gcp_shared.LogError("BREVO MAIL ERROR", "Error occurred while sending the brevo email. Response: "+string(respBodyBytes))
	}
	log.Println("Email sent successfully, status code:", resp.StatusCode)
}
