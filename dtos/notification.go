package dtos

type EmailMessage struct {
	EventID        string `json:"event_id"`
	UserID         string `json:"user_id"`
	SubscriptionID string `json:"subscription_id"`
	To             string `json:"to"`
	Subject        string `json:"subject"`
	Body           string `json:"body"`
}
