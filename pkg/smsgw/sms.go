// smsgw is an abstract implementation of SMS messages and a SMS gateway.
package smsgw

import (
	"golang.org/x/text/currency"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

// Message represents a single SMS.
type Message struct {
	// From is the sender of the SMS typically formatted as an E.164
	// international phone number such as +18005551234.
	From string
	// To is the recipient of the SMS formatted as an E.164 phone number.
	To string
	// Message is the text content of the SMS displayed to the receiving user.
	Message string
}

// ApiResponse allows a generic parsing of the response for a sent message
// to be displayed for the user.
type ApiResponse struct {
	// RemoteId is optionally an unique identifier that can be used to track
	// the sent message.
	RemoteId string
	// Cost uses the Go currency format to represent the cost of the sent
	// message. If nil, cost is assumed to not be relevant per message.
	Cost *currency.Amount
}

// SmsApi is a gateway to send SMS messages.
type SmsApi interface {
	// Send takes a message and sends it to the remote gateway getting
	// an ApiResponse if successfully sent, otherwise displays an error.
	Send(Message) (*ApiResponse, error)
}

func (apiResponse ApiResponse) CostToString() string {
	if apiResponse.Cost == nil {
		return ""
	}

	formatter := message.NewPrinter(language.AmericanEnglish)
	return formatter.Sprint(apiResponse.Cost)
}
