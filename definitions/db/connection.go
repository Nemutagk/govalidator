package db

type DbConnection struct {
	Driver        string
	Host          string
	Port          string
	User          string
	Password      string
	Database      string
	AnotherConfig *map[string]interface{}
}
