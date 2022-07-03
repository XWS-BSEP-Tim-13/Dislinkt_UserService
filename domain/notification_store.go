package domain

type NotificationStore interface {
	Insert(company *Notification) error
	DeleteAll()
	GetByUsername(username string) ([]*Notification, error)
}
