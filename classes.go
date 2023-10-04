package schedule_db_data

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
)

type Class struct {
	Id        int    `json:"id"`
	Number    int    `json:"number"`
	Character string `json:"character"`
}

func (c *Class) ToString() (class string) {
	return fmt.Sprint(c.Number, c.Character)
}

func (c *Class) Parse(class_string string) (class Class, err error) {

	c.Number, err = strconv.Atoi(class_string[:len(class_string)-2])
	if err != nil {
		return Class{}, err
	}

	c.Character = class_string[len(class_string)-2:]

	return *c, nil
}

type Classes []Class

func (classes *Classes) Contain(class Class) (res bool) {

	for _, c := range *classes {
		if c.Character == class.Character && c.Number == class.Number {
			return true
		}
	}

	return false
}

func (classes *Classes) Find(name string) (class Class, err error) {

	c, err := new(Class).Parse(name)
	if err != nil {
		return Class{}, err
	}

	for _, class := range *classes {
		if c.Character == class.Character && c.Number == class.Number {
			return class, nil
		}
	}

	return Class{}, errors.New("no class with name " + name + " has found")
}

func Get_classes() (Classes, error) {

	db, err := DB_connection(Get_db_env("getter"))
	if err != nil {
		return Classes{}, err
	}
	defer db.Close()

	result, err := db.Query("SELECT * FROM class;")
	if err != nil {
		return Classes{}, err
	}
	defer result.Close()

	classes := Classes{}

	for result.Next() {

		class := Class{}

		err := result.Scan(&class.Id, &class.Number, &class.Character)
		if err != nil {
			return Classes{}, err
		}

		classes = append(classes, class)
	}

	return classes, nil
}

func Get_classes_by_parallel(parallel int) (Classes, error) {

	db, err := DB_connection(Get_db_env("getter"))
	if err != nil {
		return Classes{}, err
	}
	defer db.Close()

	result, err := db.Query(fmt.Sprint("SELECT * FROM class WHERE number = ", parallel, ";"))
	if err != nil {
		return Classes{}, err
	}
	defer result.Close()

	classes := Classes{}

	for result.Next() {

		class := Class{}

		err := result.Scan(&class.Id, &class.Number, &class.Character)
		if err != nil {
			return Classes{}, err
		}

		classes = append(classes, class)
	}

	return classes, nil
}

func Add_classes(buff []byte) error {

	classes := Classes{}
	err := json.Unmarshal(buff, &classes)
	if err != nil {
		return err
	}

	data_str := ""

	for i, class := range classes {
		data_str += fmt.Sprint("(", class.Number, ", '", class.Character, "')")
		if i < len(classes)-1 {
			data_str += ", "
		}
	}

	db, err := DB_connection(Get_db_env("setter"))
	if err != nil {
		return err
	}
	defer db.Close()

	insert_string := "INSERT INTO teacher (number, character) VALUES " + data_str + " ON CONFLICT (number, character) DO UPDATE SET number = EXCLUDED.number, character = EXCLUDED.character;"
	_, err = db.Exec(insert_string)
	if err != nil {
		return err
	}

	return nil
}
