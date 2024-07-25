package config

type DBSettings struct {
	Username string `yaml:"user"`
	Password string `yaml:"pass"`
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	DBName   string `yaml:"dbname"`
}

type Settings struct {
	Database DBSettings `yaml:"database"`
}

func ReadSettings() *Settings {
	result := &Settings{
		Database: DBSettings{
			Username: "root",
			Password: "secret",
			Host:     "localhost",
			Port:     "3306",
			DBName:   "finances",
		},
	}

	return result
}
