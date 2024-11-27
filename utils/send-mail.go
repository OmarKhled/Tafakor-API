package utils

import (
	"fmt"
	"net/smtp"
)

// Sends mails given sender and reciever
func SendMail(senderMail string, senderPass string, host string, port string, recieverEmail string, subject string, body string) error {

	// Set up authentication information.
	auth := smtp.PlainAuth("", senderMail, senderPass, host)

	fmt.Println(auth)

	// Array of recievers
	recipientList := []string{recieverEmail}

	// Email headers - RFC 822 Message Format
	recieverHeader := fmt.Sprintf("To: %v\r\n", recieverEmail)
	subjectHeader := fmt.Sprintf("Subject: %v!\r\n", subject)
	MIMEHeader := fmt.Sprintf("MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\r\n")
	content := body

	msg := []byte(recieverHeader + subjectHeader + MIMEHeader + "\r\n" + content)
	err := smtp.SendMail(host+":"+port, auth, senderMail, recipientList, msg)
	fmt.Println("sendemail")

	return err
}
