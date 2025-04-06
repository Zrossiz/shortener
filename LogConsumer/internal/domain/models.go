package domain

import "time"

type RedirectEventDAO struct {
	ID        uint      `json:"id"`
	Original  string    `json:"original"`
	Short     string    `json:"short"`
	UserIP    string    `json:"user_ip"`
	Os        string    `json:"os"`
	CreatedAt time.Time `json:"created_at"`
}
