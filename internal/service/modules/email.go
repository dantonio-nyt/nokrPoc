//package main
//
//import (
//	"bytes"
//	"context"
//	"encoding/json"
//	"fmt"
//	"github.com/pkg/errors"
//	"html/template"
//	"strings"
//
//	"github.com/nytm/go-common"
//)
//
//// EmailTemplate
//type EmailTemplate struct {
//	UpdateFPALink string
//	FPAData       *sugar.FPA
//}
//
//// sendFPAEmail formats the FPA email and forwards it to Helix.
//func (s Service) sendFPAEmail(msg []byte) error {
//
//	var fpaTemplate sugar.FPA
//	var emails []string
//	ctx := context.Background()
//
//	err := json.Unmarshal(msg, &fpaTemplate)
//	if err != nil {
//		s.logError(fmt.Sprintf("Couldn't unmarshal message: %v, %s\n", string(msg), err.Error()))
//		return errors.New("Couldn't unmarshal message")
//	}
//	if fpaTemplate.Escalation.SubjectLine == "" {
//		s.logError(fmt.Sprintf("SubjectLine is empty."))
//		return errors.New("SubjectLine is empty.")
//	}
//
//	templateData, err := prepData(fpaTemplate, s.config)
//	if err != nil {
//		s.logError(fmt.Sprintf("Problem formatting email: %v, %s\n", string(msg), err.Error()))
//		return errors.New("Problem formatting email.")
//	}
//
//	if s.config.Env == "stg" || s.config.IsTesting {
//		if s.config.EmailTo == "" {
//			return errors.New("EmailTo is not specified.")
//		}
//		s.logInfo(fmt.Sprintf("Sending test emails to: %s\n", s.config.EmailTo))
//		emails = strings.Split(s.config.EmailTo, ",")
//	} else {
//		s.logInfo(fmt.Sprintf("Sending prd emails to: %s\n", s.config.EmailTo))
//		emails = fpaTemplate.DeliveryPartner.Emails
//		// if the email to address is specified append it to the list of emails to send out
//		// this is useful for debugging.
//		if s.config.EmailTo != "" {
//			emails = append(emails, s.config.EmailTo)
//		}
//	}
//
//	// make sure that we are sending emails to something
//	if len(emails) == 0 {
//		return errors.New("Email list is empty for escalation id: " + fpaTemplate.DeliveryPartner.EscalationID)
//	}
//
//	resp, err := s.client.SendRealtimeToAddressListByTemplateId(ctx, emails, s.config.HelixTemplateID, templateData, s.config.HelixTrackingTag)
//	if err != nil {
//		s.logError(fmt.Sprintf("escalationId: %s, accountId: %s. Unable to send email: %v",
//			fpaTemplate.Escalation.EscalationID,
//			fpaTemplate.Customer.AccountNumber,
//			err.Error()))
//		return err
//	} else {
//		s.logInfo(fmt.Sprintf("escalationId: %s, accountId: %s. Successfully dispatched email. Response: %v",
//			fpaTemplate.Escalation.EscalationID,
//			fpaTemplate.Customer.AccountNumber,
//			resp))
//	}
//
//	return nil
//}
//
//// prepData
//func prepData(fpaTemplate sugar.FPA, config Config) (map[string]string, error) {
//
//	fmap := template.FuncMap{
//		"formatTime":     formatTime,
//		"formatDate":     formatDate,
//		"formatISODate":  formatISODate,
//		"formatDateTime": formatDateTime,
//	}
//
//	// format the data
//	FormatFPAData(&fpaTemplate)
//
//	// Using minified version of table-email.tmpl
//	t := template.Must(template.New("table-email.min.tmpl").Funcs(fmap).ParseFiles("./templates/table-email.min.tmpl"))
//	// Keep the next line also, to test changes in non-minified version.
//	// t := template.Must(template.New("table-email.tmpl").Funcs(fmap).ParseFiles("./templates/table-email.tmpl"))
//	var tpl bytes.Buffer
//	err := t.Execute(&tpl, EmailTemplate{
//		UpdateFPALink: buildEscalationLink(config.FpaServiceHost, fpaTemplate.Escalation.LinkUUID),
//		FPAData:       &fpaTemplate,
//	})
//	if err != nil {
//		return nil, err
//	}
//
//	templateData := map[string]string{
//		"SubjectLine":  fpaTemplate.Escalation.SubjectLine,
//		"TemplateData": tpl.String(),
//	}
//
//	if config.Env == "prd" {
//		templateData["Cc"] = strings.Join(fpaTemplate.DeliveryPartner.Emails, ",")
//	} else if config.Env == "dev" {
//		// don't want to use real email addresses when cc'ing in dev.
//		templateData["Cc"] = "nytcsdev@nytimes.com"
//	}
//
//	return templateData, nil
//}
//
//// buildEscalationLink
//func buildEscalationLink(fpaServiceHost string, linkUUID string) string {
//	return fmt.Sprintf("%s/escalation/%s", fpaServiceHost, linkUUID)
//}
