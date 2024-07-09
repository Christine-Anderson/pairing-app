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
	SENDER                 = "pairwisenoreply@gmail.com"
	SUBJECT_GROUP_CREATION = "Your PairWise Group Has Been Created"
	SUBJECT_ASSIGNMENTS    = "Your Pairwise Assignment Has Been Generated"
	CHAR_SET               = "UTF-8"
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

func generateMailInput(email string, subject string, htmlBody string, textBody string) *ses.SendEmailInput {
	input := &ses.SendEmailInput{
		Destination: &ses.Destination{
			CcAddresses: []*string{},
			ToAddresses: []*string{
				aws.String(email),
			},
		},
		Message: &ses.Message{
			Body: &ses.Body{
				Html: &ses.Content{
					Charset: aws.String(CHAR_SET),
					Data:    aws.String(htmlBody),
				},
				Text: &ses.Content{
					Charset: aws.String(CHAR_SET),
					Data:    aws.String(textBody),
				},
			},
			Subject: &ses.Content{
				Charset: aws.String(CHAR_SET),
				Data:    aws.String(subject),
			},
		},
		Source: aws.String(SENDER),
	}
	return input
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

	input := generateMailInput(recipient.Email, SUBJECT_GROUP_CREATION, htmlBody, textBody)

	result, sendEmailErr := es.emailService.SendEmail(input)
	if sendEmailErr != nil {
		return sendEmailErr
	}

	log.Printf("Email Sent to address: %s", recipient.Email)
	log.Printf("%s", result)
	return nil
}

func generateMemberMap(group types.Group) map[string]types.GroupMember {
	memberMap := make(map[string]types.GroupMember)

	for _, member := range group.GroupMembers {
		memberMap[member.MemberId] = member
	}

	return memberMap
}

func (es EmailService) SendAssignmentEmails(assignments map[string]string, group types.Group) {
	memberMap := generateMemberMap(group)

	for giver, receiver := range assignments {
		htmlBody := fmt.Sprintf(
			`<p>Hi %s,</p>
			<p>The PairWise group <b>%s</b> with group ID <b>%s</b> has generated assignments.</p>
			<p>You were paired with %s.</p>
			<p>Happy pairing!</p>`, memberMap[giver].Name, group.GroupName, group.GroupId, memberMap[receiver].Name)

		textBody := fmt.Sprintf(
			`Hi %s,\n
			The PairWise group %s with group ID %s has generated assignments.\n
			You were paired with %s.\n
			Happy pairing!`, memberMap[giver].Name, group.GroupName, group.GroupId, memberMap[receiver].Name)

		input := generateMailInput(memberMap[giver].Email, SUBJECT_ASSIGNMENTS, htmlBody, textBody)
		result, sendEmailErr := es.emailService.SendEmail(input)
		if sendEmailErr != nil {
			log.Printf("Failed to send email to %s with error: %s", memberMap[giver].Email, sendEmailErr.Error())
		} else {
			log.Printf("Email Sent to address: %s", memberMap[giver].Email)
			log.Printf("%s", result)
		}
	}
}
