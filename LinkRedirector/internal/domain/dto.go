package domain

type UrlKafkaDTO struct {
	Original string `json:"original"`
	Short    string `json:"short"`
	UserIP   string `json:"user_ip"`
	OS       string `json:"os"`
}

type GetUrlDTO struct {
	Short  string `json:"short"`
	UserIP string `json:"user_ip"`
	OS     string `json:"os"`
}
