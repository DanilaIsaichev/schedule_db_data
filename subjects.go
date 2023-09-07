package schedule_db_data

import (
	"encoding/json"
	"errors"
	"fmt"
	//_ "github.com/lib/pq"
)

type Subject struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	Short_name  string `json:"short_name"`
	Groups      int    `json:"groups"`
	Description string `json:"description"`
}

type Subjects []Subject

func (subjects *Subjects) Contain(subject Subject) (res bool) {

	for _, s := range *subjects {
		if s.Name == subject.Name {
			return true
		}
	}

	return false
}

func (subjects *Subjects) Find(name string) (class Subject, err error) {

	for _, subject := range *subjects {
		if name == subject.Name {
			return subject, nil
		}
	}

	return Subject{}, errors.New("no subject with name " + name + " has found")
}

func Get_subjects() (Subjects, error) {

	db, err := DB_connection(Get_db_env("getter"))
	if err != nil {
		return Subjects{}, err
	}
	defer db.Close()

	result, err := db.Query("SELECT * FROM subject;")
	if err != nil {
		return Subjects{}, nil
	}
	defer result.Close()

	subjects := Subjects{}

	for result.Next() {

		subject := Subject{}

		err := result.Scan(&subject.Id, &subject.Name, &subject.Short_name, &subject.Groups, &subject.Description)
		if err != nil {
			return Subjects{}, err
		}

		subjects = append(subjects, subject)
	}

	return subjects, nil
}

func Add_subjects(buff []byte) error {

	subjects := Subjects{}
	err := json.Unmarshal(buff, &subjects)
	if err != nil {
		return err
	}

	data_str := ""

	for i, subject := range subjects {
		data_str += fmt.Sprint("('", subject.Name, "', '", subject.Short_name, "', ", subject.Groups, ", '", subject.Description, "')")
		if i < len(subjects)-1 {
			data_str += ", "
		}
	}

	db, err := DB_connection(Get_db_env("setter"))
	if err != nil {
		return err
	}
	defer db.Close()

	insert_string := "INSERT INTO subject (name, short_name, groups, description) VALUES " + data_str + " ON CONFLICT (name, short_name) DO UPDATE SET name = EXCLUDED.name, short_name = EXCLUDED.short_name, groups = EXCLUDED.groups, description = EXCLUDED.description;"
	_, err = db.Exec(insert_string)
	if err != nil {
		return err
	}

	return nil
}
