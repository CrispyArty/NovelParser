package smtp

import (
	// "crypto/tls"
	"crypto/tls"
	"errors"
	"fmt"
	"net"
	"net/smtp"

	"github.com/crispyarty/novelparser/internal/mailer/common"
)

// func getCreds() {
// 	dbPassword := os.Getenv("DB_PASSWORD")
// 	if dbPassword == "" {
// 		fmt.Println("DB_PASSWORD environment variable not set.")
// 	} else {
// 		fmt.Printf("Database Password: %s\n", dbPassword)
// 	}
// }

type loginAuth struct {
	username, password string
}

func LoginAuth(username, password string) smtp.Auth {
	return &loginAuth{username, password}
}

func (a *loginAuth) Start(server *smtp.ServerInfo) (string, []byte, error) {
	return "LOGIN", []byte(a.username), nil
}

func (a *loginAuth) Next(fromServer []byte, more bool) ([]byte, error) {
	if more {
		switch string(fromServer) {
		case "Username:":
			return []byte(a.username), nil
		case "Password:":
			return []byte(a.password), nil
		default:
			return nil, errors.New("Unknown from server")
		}
	}
	return nil, nil
}

func SendSmtp(message *common.Message) {
	password := "" // Use an app password for services like Gmail

	// smtp server configuration.
	smtpHost := "smtp.office365.com"
	smtpPort := "587"

	conn, err := net.Dial("tcp", "smtp.office365.com:587")
	if err != nil {
		println(err)
	}

	c, err := smtp.NewClient(conn, smtpHost)
	if err != nil {
		println(err)
	}

	tlsconfig := &tls.Config{
		ServerName: smtpHost,
	}

	if err = c.StartTLS(tlsconfig); err != nil {
		println(err)
	}

	auth := LoginAuth(message.From, password)

	if err = c.Auth(auth); err != nil {
		println(err)
	}

	// Sending email.
	err = smtp.SendMail(smtpHost+":"+smtpPort, auth, message.From, []string{message.To}, message.Raw())
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Email Sent!")

	// conn, err := tls.Dial("tcp", "smtp-mail.outlook.com:587", &tls.Config{
	// 	ServerName: "smtp-mail.outlook.com",
	// })

	// client, err := smtp.NewClient(conn, "smtp-mail.outlook.com")

	// defer client.Close()

	// auth := LoginAuth(from, password)

	// err := smtp.SendMail(
	// 	"smtp.gmail.com:587", // Replace with your SMTP server and port
	// 	// auth,                 // Replace with your SMTP host
	// 	smtp.PlainAuth("", from, password, "smtp.gmail.com"),
	// 	from,
	// 	[]string{to},
	// 	message(bookPath),
	// )

	// // err := smtp.SendMail(
	// // 	"smtp-mail.outlook.com:587",                                 // Replace with your SMTP server and port
	// // 	smtp.PlainAuth("", from, password, "smtp-mail.outlook.com"), // Replace with your SMTP host
	// // 	from,
	// // 	[]string{to},
	// // 	message(bookPath),
	// // )

	// if err != nil {
	// 	log.Printf("smtp error: %s", err)
	// 	return
	// }
	// log.Println("Email sent successfully!")
}
