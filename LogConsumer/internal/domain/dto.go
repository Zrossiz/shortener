package domain

type RegisterRedirectEventDTO struct {
	Original string `json:"original"`
	Short    string `json:"short"`
	UserIP   string `json:"user_ip"`
	Os       string `json:"os"`
}
