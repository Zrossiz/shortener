package domain

import "time"

type RegisterRedirectEvent struct {
	Original string
	Short    string
	UserIP   string `json:"user_ip"`
	Os       string
}

type RedirectEventDAO struct {
	ID        uint
	Original  string
	Short     string
	UserIP    string `json:"user_ip"`
	Os        string
	CreatedAt time.Time `json:"created_at"`
}
