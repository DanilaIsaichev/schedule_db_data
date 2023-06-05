package schedule_db_data

import (
	"bytes"
	"database/sql"
	"errors"
	"io"
)

func scan_prepare(src interface{}) (prepared_bytes []byte, err error) {

	// Массив байтов
	data := []byte{}

	// Приведение к байтам и запись в массив
	if val, ok := src.([]byte); ok {
		data = val
	} else if val, ok := src.([]byte); ok {
		data = []byte(val)
	} else if src == nil {
		return []byte{}, errors.New("couldn't convert db data to []byte")
	}

	// Новый reader для массива
	reader := bytes.NewReader(data)

	// Считываем байты
	bdata, err := io.ReadAll(reader)
	if err != nil {
		return []byte{}, err
	}

	return bdata, nil
}

func DB_connection(hostname string, db_name string, username string, password string, port string) (db_conn *sql.DB, err error) {

	connection_string := "host=" + hostname + " port=" + port + " user=" + username + " password=" + password + " dbname=" + db_name + " sslmode=disable"

	db, err := sql.Open("postgres", connection_string)
	if err != nil {
		return db, err
	}

	return db, nil
}
