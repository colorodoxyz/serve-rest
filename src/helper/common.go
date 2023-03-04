package helper

type Account struct {
	Username string `json:"Username"`
	Password string `json:"Password"`
}

type KeyValue struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

const (
	AdminUser     = "admin"
	AdminPassword = "admin"
	Url           = "localhost:5001"
	ServerUrl     = "http://" + Url
)
