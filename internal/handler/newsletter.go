package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"regexp"

	"github.com/axadrn/axeladrian/internal/config"
)

type NewsletterHandler struct {
	cfg *config.Config
}

func NewNewsletterHandler(cfg *config.Config) *NewsletterHandler {
	return &NewsletterHandler{cfg: cfg}
}

type subscribeRequest struct {
	Email string `json:"email"`
}

type subscribeResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error,omitempty"`
}

func (h *NewsletterHandler) Subscribe(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var req subscribeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		json.NewEncoder(w).Encode(subscribeResponse{Error: "Invalid request"})
		return
	}

	if !isValidEmail(req.Email) {
		json.NewEncoder(w).Encode(subscribeResponse{Error: "Invalid email"})
		return
	}

	// Add to Resend audience
	apiKey := h.cfg.ResendAPIKey
	audienceID := h.cfg.ResendAudienceID

	if apiKey == "" || audienceID == "" {
		json.NewEncoder(w).Encode(subscribeResponse{Error: "Newsletter not configured"})
		return
	}

	payload := map[string]string{"email": req.Email}
	body, _ := json.Marshal(payload)

	url := "https://api.resend.com/audiences/" + audienceID + "/contacts"
	httpReq, _ := http.NewRequest("POST", url, bytes.NewBuffer(body))
	httpReq.Header.Set("Authorization", "Bearer "+apiKey)
	httpReq.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(httpReq)
	if err != nil {
		// Return success anyway to prevent email enumeration
		json.NewEncoder(w).Encode(subscribeResponse{Success: true})
		return
	}
	defer resp.Body.Close()

	// Always return success to prevent email enumeration
	json.NewEncoder(w).Encode(subscribeResponse{Success: true})
}

func isValidEmail(email string) bool {
	re := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return re.MatchString(email)
}
