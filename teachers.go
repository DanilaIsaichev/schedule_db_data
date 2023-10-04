package schedule_db_data

import (
	"encoding/json"
	"errors"
	"fmt"
	//_ "github.com/lib/pq"
)

type Teacher struct {
	Id         int    `json:"id"`
	Login      string `json:"login"`
	First_name string `json:"first_name"`
	Last_name  string `json:"last_name"`
	Patronymic string `json:"patronymic"`
	Short_name string `json:"short_name"`
}

type Teachers []Teacher

func (teachers *Teachers) Contain(teacher Teacher) (res bool) {

	for _, t := range *teachers {
		if t.Login == teacher.Login {
			return true
		}
	}

	return false
}

func (teachers *Teachers) Find_by_id(id int) (teacher Teacher, err error) {

	for _, teacher := range *teachers {
		if teacher.Id == id {
			return teacher, nil
		}
	}

	return Teacher{}, errors.New(fmt.Sprint("no teacher with id ", id, " has found"))
}

func (teachers *Teachers) Find_by_login(login string) (teacher Teacher, err error) {

	for _, teacher := range *teachers {
		if teacher.Login == login {
			return teacher, nil
		}
	}

	return Teacher{}, errors.New("no teacher with login " + login + " has found")
}

func Get_teachers() (Teachers, error) {

	db, err := DB_connection(Get_db_env("getter"))
	if err != nil {
		return Teachers{}, err
	}
	defer db.Close()

	result, err := db.Query("SELECT * FROM teacher;")
	if err != nil {
		return Teachers{}, err
	}
	defer result.Close()

	teachers := Teachers{}

	for result.Next() {

		teacher := Teacher{}

		err := result.Scan(&teacher.Id, &teacher.Login, &teacher.First_name, &teacher.Last_name, &teacher.Patronymic, &teacher.Short_name)
		if err != nil {
			return Teachers{}, err
		}

		teachers = append(teachers, teacher)
	}

	return teachers, nil
}

func Add_teachers(buff []byte) error {

	teachers := Teachers{}
	err := json.Unmarshal(buff, &teachers)
	if err != nil {
		return err
	}

	data_str := ""

	for i, teacher := range teachers {
		data_str += fmt.Sprint("('", teacher.Login, "', '", teacher.First_name, "', '", teacher.Last_name, "', '", teacher.Patronymic, "', '", teacher.Short_name, "')")
		if i < len(teachers)-1 {
			data_str += ", "
		}
	}

	db, err := DB_connection(Get_db_env("setter"))
	if err != nil {
		return err
	}
	defer db.Close()

	insert_string := "INSERT INTO teacher (login, first_name, last_name, patronymic, short_name) VALUES " + data_str + " ON CONFLICT (login) DO UPDATE SET login = EXCLUDED.login, first_name = EXCLUDED.first_name, last_name = EXCLUDED.last_name, patronymic = EXCLUDED.patronymic, short_name = EXCLUDED.short_name;"
	_, err = db.Exec(insert_string)
	if err != nil {
		return err
	}

	return nil
}
