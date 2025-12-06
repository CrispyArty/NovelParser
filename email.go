package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/smtp"
	"os"
	"strings"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/option"
)

// Retrieve a token, saves the token, then returns the generated client.
func getClient(config *oauth2.Config) *http.Client {
	// The file token.json stores the user's access and refresh tokens, and is
	// created automatically when the authorization flow completes for the first
	// time.
	tokFile := "token.json"
	tok, err := tokenFromFile(tokFile)
	if err != nil {
		tok = getTokenFromWeb(config)
		saveToken(tokFile, tok)
	}
	return config.Client(context.Background(), tok)
}

// Request a token from the web, then returns the retrieved token.
func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)

	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		log.Fatalf("Unable to read authorization code: %v", err)
	}

	tok, err := config.Exchange(context.TODO(), authCode)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web: %v", err)
	}
	return tok
}

// Retrieves a token from a local file.
func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

// Saves a token to a file path.
func saveToken(path string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", path)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}

func getCreds() {

	dbPassword := os.Getenv("DB_PASSWORD")
	if dbPassword == "" {
		fmt.Println("DB_PASSWORD environment variable not set.")
	} else {
		fmt.Printf("Database Password: %s\n", dbPassword)
	}
}

func sendEmail() {
	ctx := context.Background()
	b, err := os.ReadFile("client_secret.json")
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	// If modifying these scopes, delete your previously saved token.json.
	config, err := google.ConfigFromJSON(b, gmail.GmailReadonlyScope, gmail.GmailSendScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}
	client := getClient(config)

	srv, err := gmail.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		log.Fatalf("Unable to retrieve Gmail client: %v", err)
	}

	from := "crispykindle@gmai.com"
	to := "artemti0@gmail.com"
	subject := "Hello from Go!"
	body := "This is a test email sent from a Go program."

	attachmentFileName := "Chapter 111 - 120.epub"
	book, _ := os.ReadFile("uploads/my_simulated_road_to_immortality/" + attachmentFileName)
	encodedAttachment := base64.StdEncoding.EncodeToString(book)

	var messageBuilder strings.Builder
	messageBuilder.WriteString(fmt.Sprintf("From: %s\r\n", from))
	messageBuilder.WriteString(fmt.Sprintf("To: %s\r\n", to))
	messageBuilder.WriteString(fmt.Sprintf("Subject: %s\r\n", subject))
	messageBuilder.WriteString("MIME-Version: 1.0\r\n")
	messageBuilder.WriteString("Content-Type: multipart/mixed; boundary=\"foo_bar_baz\"\r\n\r\n")

	// Add text part
	messageBuilder.WriteString("--foo_bar_baz\r\n")
	messageBuilder.WriteString("Content-Type: text/plain; charset=\"UTF-8\"\r\n")
	messageBuilder.WriteString("Content-Transfer-Encoding: quoted-printable\r\n\r\n")
	messageBuilder.WriteString(body + "\r\n\r\n")

	// Add attachment part
	messageBuilder.WriteString("--foo_bar_baz\r\n")
	messageBuilder.WriteString(fmt.Sprintf("Content-Type: application/epub+zip; name=\"%s\"\r\n", attachmentFileName)) // Adjust MIME type for your file
	messageBuilder.WriteString("Content-Disposition: attachment; filename=\"" + attachmentFileName + "\"\r\n")
	messageBuilder.WriteString("Content-Transfer-Encoding: base64\r\n\r\n")
	messageBuilder.WriteString(encodedAttachment + "\r\n\r\n")

	messageBuilder.WriteString("--foo_bar_baz--")

	message := &gmail.Message{
		Raw: base64.StdEncoding.EncodeToString([]byte(messageBuilder.String())),
	}
	user := "me"

	// mediaOptions := googleapi.ContentType("uploads/my_simulated_road_to_immortality/Chapter 111 - 120.epub")

	m, err := srv.Users.Messages.Send(user, message).Do()

	if err != nil {
		log.Fatalf("Unable to retrieve labels: %v", err)
	}

	fmt.Println("Success!", m)

	// r, err := srv.Users.Labels.List(user).Do()
	// if err != nil {
	// 	log.Fatalf("Unable to retrieve labels: %v", err)
	// }
	// if len(r.Labels) == 0 {
	// 	fmt.Println("No labels found.")
	// 	return
	// }
	// fmt.Println("Labels:")
	// for _, l := range r.Labels {
	// 	fmt.Printf("- %s\n", l.Name)
	// }
}

func sendSmtp() {
	from := "crispykindle@gmail.com"
	password := "" // Use an app password for services like Gmail
	to := "artemti0@gmail.com"
	subject := "Hello from Go!"
	body := "This is a test email sent from a Go program."

	msg := "From: " + from + "\n" +
		"To: " + to + "\n" +
		"Subject: " + subject + "\n\n" +
		body

	err := smtp.SendMail(
		"smtp.gmail.com:587", // Replace with your SMTP server and port
		smtp.PlainAuth("", from, password, "smtp.gmail.com"), // Replace with your SMTP host
		from,
		[]string{to},
		[]byte(msg),
	)

	if err != nil {
		log.Printf("smtp error: %s", err)
		return
	}
	log.Println("Email sent successfully!")
}

func main() {
	sendEmail()
}
