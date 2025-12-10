package gmailapi

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	appConfig "github.com/crispyarty/novelparser/internal/config"
	"github.com/crispyarty/novelparser/internal/mailer/common"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/option"
)

var (
	tokenFilepath   = appConfig.AssetPath("token.json")
	secretsFilepath = appConfig.AssetPath("client_secret.json")
	_token          *oauth2.Token
)

func getToken() *oauth2.Token {
	if _token != nil {
		return _token
	}

	var err error
	_token, err = tokenFromFile()

	if err != nil {
		_token = nil
	}

	return _token
}

func setToken(t *oauth2.Token) {
	_token = t
}

func ValidateToken() error {
	token := getToken()

	if token == nil {
		return fmt.Errorf("cant retrive token from file \"%v\"", tokenFilepath)
	}

	ctx := context.Background()
	config := oauth2Config()

	newToken, err := config.TokenSource(ctx, token).Token()

	if err != nil {
		return err
		// 	log.Printf("Error refreshing token: %v", err)
		// 	// The refresh token is likely invalid or revoked
		// 	fmt.Println("Refresh token is NOT valid.")
		// } else {
		// 	fmt.Printf("%T\n", newToken)
		// 	fmt.Println("new token: ", newToken)
		// 	fmt.Println("AccessToken: ", newToken.AccessToken)
		// 	fmt.Println("RefreshToken: ", newToken.RefreshToken)
		// 	fmt.Println("Expiry: ", newToken.Expiry)
		// 	fmt.Println("ExpiresIn: ", newToken.ExpiresIn)
	}

	saveToken(newToken)

	return nil
}

func CreateNewToken() {
	config := oauth2Config()
	token := getTokenFromWeb(config)
	saveToken(token)
}

func SendBook(message *common.Message) {
	srv := gmailService()

	// fmt.Println(message.Raw())
	msg := &gmail.Message{
		Raw: string(message.Raw()),
	}

	user := "me"

	s := srv.Users.Messages.Send(user, msg)
	// fmt.Println("Send request!", s)

	_, err := s.Do()

	if err != nil {
		log.Panicf("Send book via email error: %v", err)

		// var apiErr *oauth2.RetrieveError

		// if ok := errors.As(err, &apiErr); ok {
		// 	log.Println("ErrorCode:", apiErr.ErrorCode)

		// 	log.Println("Header:", apiErr.Response.Header)
		// 	log.Println("Body:", string(apiErr.Body))
		// 	log.Panicf("Google Api error: %v", apiErr)
		// } else {
		// 	log.Panicf("Send book via email error: %v", err)
		// }
	}

	// fmt.Println("Success!")
}

// Request a token from the web, then returns the retrieved token.
func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)

	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		log.Panicf("Unable to read authorization code: %v", err)
	}

	tok, err := config.Exchange(context.TODO(), authCode)
	if err != nil {
		log.Panicf("Unable to retrieve token from web: %v", err)
	}
	return tok
}

// Retrieves a token from a local file.
func tokenFromFile() (*oauth2.Token, error) {
	f, err := os.Open(tokenFilepath)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

// Saves a token to a file path.
func saveToken(token *oauth2.Token) {
	setToken(token)

	// fmt.Printf("Saving credential file to: %s\n", path)
	f, err := os.OpenFile(tokenFilepath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Panicf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	encoder := json.NewEncoder(f)
	encoder.SetIndent("", "  ")
	encoder.Encode(token)
}

func oauth2Config() *oauth2.Config {
	b, err := os.ReadFile(secretsFilepath)
	if err != nil {
		log.Panicf("Unable to read client secret file: %v", err)
	}

	// gmail.GmailReadonlyScope
	config, err := google.ConfigFromJSON(b, gmail.GmailSendScope)
	if err != nil {
		log.Panicf("Unable to parse client secret file to config: %v", err)
	}

	return config
}

func gmailService() *gmail.Service {
	ctx := context.Background()

	config := oauth2Config()

	// fmt.Println(getToken())

	client := config.Client(context.Background(), getToken())

	srv, err := gmail.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		log.Panicf("Unable to retrieve Gmail client: %v", err)
	}

	return srv
}
