package modules

import (
	"encoding/json"
	"fmt"
	"github.com/nokrPOC/internal/service"
	"net/http"
)

type PubSubService struct {
	Logger
}

func NewBugSubService() {
	return PubSubService{}
}
// pushRequest is the expect format of the message
// received from PubSub.
type pushRequest struct {
	Message struct {
		Attributes map[string]string
		Data       []byte
		ID         string `json:"messageId"`
	}
	Subscription string
}

// pushHandler handles push requests from PubSub, receiving the
// message and forwarding the email. Requests must contain
// the configured token.
func (s PubSubService) pushHandler(w http.ResponseWriter, r *http.Request) {

	l := getLogger(r.Context())

	// Verify the token.
	if r.URL.Query().Get("token") != s.config.PubSubVerificationToken {
		http.Error(w, "bad token", http.StatusUnauthorized)
		return
	}

	if r.Body == nil {
		http.Error(w, "missing request body", http.StatusBadRequest)
		return
	}

	// Get the message
	msg := &pushRequest{}
	if err := json.NewDecoder(r.Body).Decode(msg); err != nil {
		http.Error(w, fmt.Sprintf("could not decode body: %v", err), http.StatusBadRequest)
		return
	}

	// Send email
	err := s.sendFPAEmail(msg.Message.Data)
	if err != nil {
		l.logError(err.Error())
		http.Error(w, "email send failure", http.StatusInternalServerError)
	}
}
