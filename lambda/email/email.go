package email

import (
	"fmt"
	"lambda/jwt"
	"lambda/types"
	"log"
	"net/url"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
)

const (
	Sender  = "pairwisenoreply@gmail.com"
	Subject = "Your PairWise Group Has Been Created"
	CharSet = "UTF-8"
)

type EmailService struct {
	emailService *ses.SES
}

func NewEmailService() (EmailService, error) {
	var emailService EmailService

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-east-1")},
	)
	if err != nil {
		return emailService, err
	}

	svc := ses.New(sess)

	return EmailService{
		emailService: svc,
	}, nil
}

func (es EmailService) SendVerificationEmail(email string) error {
	input := &ses.VerifyEmailIdentityInput{
		EmailAddress: aws.String(email),
	}

	_, err := es.emailService.VerifyEmailIdentity(input)

	if err != nil {
		return err
	}

	log.Printf("Verification sent to address: %s", email)
	return nil
}

func (es EmailService) SendConfirmationEmail(group types.Group) error {
	baseURL := os.Getenv("BASE_URL")
	endpoint := fmt.Sprintf(`group-details/%s`, group.GroupId)

	recipient := group.GroupMembers[0]
	jwt, jwtErr := jwt.CreateToken(group.GroupId)
	if jwtErr != nil {
		return jwtErr
	}

	magicLink := baseURL + endpoint + "?jwt=" + url.QueryEscape(jwt)

	htmlBody := fmt.Sprintf(
		`<p>Hi %s,</p>
		<p>Your PairWise group <b>%s</b> with group ID <b>%s</b> has been created. Please share this ID with anyone you would like to join your group.</p>
		<p>Click <a href="%s">here</a> to access your group and perform pairing.</p>
		<p>Happy pairing!</p>`, recipient.Name, group.GroupName, group.GroupId, magicLink)

	textBody := fmt.Sprintf(
		`Hi %s,\n
		Your PairWise group %s with group ID %s has been created. Please share this ID with anyone you would like to join your group.\n
        Click here to access your group and perform pairing: %s\n
		Happy pairing!`, recipient.Name, group.GroupName, group.GroupId, magicLink)

	input := &ses.SendEmailInput{
		Destination: &ses.Destination{
			CcAddresses: []*string{},
			ToAddresses: []*string{
				aws.String(recipient.Email),
			},
		},
		Message: &ses.Message{
			Body: &ses.Body{
				Html: &ses.Content{
					Charset: aws.String(CharSet),
					Data:    aws.String(htmlBody),
				},
				Text: &ses.Content{
					Charset: aws.String(CharSet),
					Data:    aws.String(textBody),
				},
			},
			Subject: &ses.Content{
				Charset: aws.String(CharSet),
				Data:    aws.String(Subject),
			},
		},
		Source: aws.String(Sender),
	}

	result, sendEmailErr := es.emailService.SendEmail(input)
	if sendEmailErr != nil {
		return sendEmailErr
	}

	log.Printf("Email Sent to address: %s", recipient.Email)
	log.Printf("%s", result)
	return nil
}

// func (es EmailService) SendAssignmentEmails(groupMembers map[types.GroupMember]types.GroupMember) error { // todo
// 	return nil
// }
