package helper

/**
 * Helper file to contain structs and constants used throught the repo
 */

type Account struct {
	Username string `json:"Username"`
	Password string `json:"Password"`
}

type KeyValue struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type BearerToken struct {
	Type, Token string
}

const (
	AdminUser     = "admin"
	AdminPassword = "admin"
	Url           = "localhost:5001"
	ServerUrl     = "https://" + Url
	MissingKeyMsg = "Key not found"
	LoginApi      = "/api/login"
	KeyValueApi   = "/api/key"
	Auth          = "Authorization"
)
