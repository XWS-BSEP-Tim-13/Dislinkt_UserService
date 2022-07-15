package domain

type NotificationDto struct {
	Username         string `json:"username"`
	NotificationType string `json:"notification_type"`
}
