package dtos

type EmailMessage struct {
	ReminderID     string `json:"reminder_id"`
	UserID         string `json:"user_id"`
	SubscriptionID string `json:"subscription_id"`
	To             string `json:"to"`
	Title          string `json:"title"`
	Content        string `json:"content"`
	InAppTitle     string `json:"inapp_title"`
	InAppContent   string `json:"inapp_content"`
}
