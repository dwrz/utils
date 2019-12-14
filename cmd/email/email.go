package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
)

var (
	awsregion = flag.String("aws-region", "us-east-1", "aws region")
	bodyHTML  = flag.String("body-html", "", "email html body")
	bodyText  = flag.String("body-text", "", "email text body")
	charset   = flag.String("charset", "UTF-8", "email charset")
	// NB: Sender must be verified in AWS SES.
	from    = flag.String("from", "", "email sender (required)")
	subject = flag.String("subject", "", "email subject")
	// NB: Recipient may need to be verified in AWS SES.
	to = flag.String("to", "", "email recipient(s) (required)")
)

func main() {
	flag.Parse()

	switch {
	case *from == "":
		fmt.Fprintf(os.Stderr, "missing from\n")
		return
	case *to == "":
		fmt.Fprintf(os.Stderr, "missing to\n")
		return
	}

	// NB: assume that AWS credentials are handled by AWS SDK.
	// E.g., shared credential files or environment variables.
	sess, err := session.NewSession(&aws.Config{Region: awsregion})
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to create AWS session: %v", err)
		return
	}

	svc := ses.New(sess)

	// Compose the message.
	email := &ses.SendEmailInput{
		Destination: &ses.Destination{
			CcAddresses: []*string{},
			ToAddresses: []*string{to},
		},
		Message: &ses.Message{
			Body: &ses.Body{},
			Subject: &ses.Content{
				Charset: charset,
				Data:    subject,
			},
		},
		Source: from,
	}
	if *bodyHTML != "" {
		email.Message.Body.Html = &ses.Content{
			Charset: charset,
			Data:    bodyHTML,
		}
	}
	if *bodyText != "" {
		email.Message.Body.Text = &ses.Content{
			Charset: charset,
			Data:    bodyText,
		}
	}

	result, err := svc.SendEmail(email)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to send email: %v", err)
		return
	}

	if result.MessageId != nil {
		fmt.Println(*result.MessageId)
	}
}
