package schedule_db_data

import (
	"database/sql"
	"os"
)

func get_db_env(user_type string) (hostname string, name string, port string, user string, password string) {

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

	user_name := ""
	user_password := ""

	if user_type == "setter" {
		user_name = os.Getenv("DB_SETTER_NAME")
		if user_name == "" {
			user_name = "setter"
		}

		user_password := os.Getenv("DB_SETTER_PASSWORD")
		if user_password == "" {
			user_password = "123456"
		}

	} else {
		user_name = os.Getenv("DB_GETTER_NAME")
		if user_name == "" {
			user_name = "getter"
		}

		user_password := os.Getenv("DB_GETTER_PASSWORD")
		if user_password == "" {
			user_password = "123456"
		}

	}

	return db_hostname, db_name, user_name, user_password, db_port
}

func DB_connection(hostname string, db_name string, username string, password string, port string) (db_conn *sql.DB, err error) {

	connection_string := "host=" + hostname + " port=" + port + " user=" + username + " password=" + password + " dbname=" + db_name + " sslmode=disable"

	db, err := sql.Open("postgres", connection_string)
	if err != nil {
		return db, err
	}

	return db, nil
}
