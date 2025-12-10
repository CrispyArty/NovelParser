package common

import (
	"encoding/base64"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type Message struct {
	From     string
	To       string
	Subject  string
	Text     string
	Bookpath string
}

func (m *Message) Raw() []byte {
	book, err := os.ReadFile(m.Bookpath)
	if err != nil {
		log.Panicf("Error reading book: %v", err)
	}

	encodedAttachment := base64.StdEncoding.EncodeToString(book)
	attachmentFileName := filepath.Base(m.Bookpath)

	var messageBuilder strings.Builder
	messageBuilder.WriteString(fmt.Sprintf("From: %s\r\n", m.From))
	messageBuilder.WriteString(fmt.Sprintf("To: %s\r\n", m.To))
	messageBuilder.WriteString(fmt.Sprintf("Subject: %s\r\n", m.Subject))
	messageBuilder.WriteString("MIME-Version: 1.0\r\n")
	messageBuilder.WriteString("Content-Type: multipart/mixed; boundary=\"CrispyArtyNovelParser\"\r\n\r\n")

	// Add text part
	messageBuilder.WriteString("--CrispyArtyNovelParser\r\n")
	messageBuilder.WriteString("Content-Type: text/plain; charset=\"UTF-8\"\r\n")
	messageBuilder.WriteString("Content-Transfer-Encoding: quoted-printable\r\n\r\n")
	messageBuilder.WriteString(m.Text + "\r\n\r\n")

	// Add attachment part
	messageBuilder.WriteString("--CrispyArtyNovelParser\r\n")
	messageBuilder.WriteString(fmt.Sprintf("Content-Type: application/epub+zip; name=\"%s\"\r\n", attachmentFileName))
	messageBuilder.WriteString("Content-Disposition: attachment; filename=\"" + attachmentFileName + "\"\r\n")
	messageBuilder.WriteString("Content-Transfer-Encoding: base64\r\n\r\n")
	messageBuilder.WriteString(encodedAttachment + "\r\n\r\n")

	messageBuilder.WriteString("--CrispyArtyNovelParser--")

	msg := []byte(messageBuilder.String())
	dst := make([]byte, base64.StdEncoding.EncodedLen(len(msg)))

	base64.StdEncoding.Encode(dst, msg)

	return dst
	// return []byte(base64.StdEncoding.EncodeToString([]byte(messageBuilder.String())))
}
