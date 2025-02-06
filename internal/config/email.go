package config

type EmailConfig struct {
	Host     string
	Port     int
	Username string
	Password string
	Sender   string
}

func GetEmailConfig() *EmailConfig {
	return &EmailConfig{
		Host:     "smtp.example.com",
		Port:     587,
		Username: "your-email@example.com",
		Password: "your-email-password",
		Sender:   "no-reply@example.com",
	}
}
