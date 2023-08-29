package schedule_db_data

import (
	"database/sql"
	"os"
)

func get_db_env() (hostname string, name string, port string, user string, password string) {

	db_hostname := os.Getenv("DB_HOSTNAME")
	if db_hostname == "" {
		db_hostname = "localhost"
	}

	db_name := os.Getenv("DB_NAME")
	if db_name == "" {
		db_name = "scheduler"
	}

	db_port := os.Getenv("DB_PORT")
	if db_port == "" {
		db_port = "5432"
	}

	getter_name := os.Getenv("DB_GETTER_NAME")
	if getter_name == "" {
		getter_name = "getter"
	}

	getter_password := os.Getenv("DB_GETTER_PASSWORD")
	if getter_password == "" {
		getter_password = "123456"
	}

	return db_hostname, db_name, db_port, getter_name, getter_password
}

func DB_connection(hostname string, db_name string, username string, password string, port string) (db_conn *sql.DB, err error) {

	connection_string := "host=" + hostname + " port=" + port + " user=" + username + " password=" + password + " dbname=" + db_name + " sslmode=disable"

	db, err := sql.Open("postgres", connection_string)
	if err != nil {
		return db, err
	}

	return db, nil
}
