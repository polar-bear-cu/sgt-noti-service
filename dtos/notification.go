package dtos

type EmailMessage struct {
	ReminderID string `json:"reminder_id"`
	To         string `json:"to"`
	Title      string `json:"title"`
	Content    string `json:"content"`
}
