package schedule_db_data

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
)

type Group struct {
	Id         int    `json:"id"`
	YearNumber int    `json:"yearNumber"`
	Character  string `json:"character"`
}

func (g *Group) ToString() (class string) {
	return fmt.Sprint(g.YearNumber, g.Character)
}

func (g *Group) Parse(class_string string) (class Group, err error) {

	g.YearNumber, err = strconv.Atoi(class_string[:len(class_string)-2])
	if err != nil {
		return Group{}, err
	}

	g.Character = class_string[len(class_string)-2:]

	return *g, nil
}

type Groups []Group

func (groups *Groups) Contain(group Group) (res bool) {

	for _, g := range *groups {
		if g.Character == group.Character && g.YearNumber == group.YearNumber {
			return true
		}
	}

	return false
}

func (groups *Groups) Find(name string) (group Group, err error) {

	g, err := new(Group).Parse(name)
	if err != nil {
		return Group{}, err
	}

	for _, group := range *groups {
		if g.Character == group.Character && g.YearNumber == group.YearNumber {
			return group, nil
		}
	}

	return Group{}, errors.New("no Group with name " + name + " has found")
}

func Get_groups() (Groups, error) {

	db, err := DB_connection(Get_db_env("getter"))
	if err != nil {
		return Groups{}, err
	}
	defer db.Close()

	result, err := db.Query("SELECT * FROM groups;")
	if err != nil {
		return Groups{}, err
	}
	defer result.Close()

	groups := Groups{}

	for result.Next() {

		group := Group{}

		err := result.Scan(&group.Id, &group.YearNumber, &group.Character)
		if err != nil {
			return Groups{}, err
		}

		groups = append(groups, group)
	}

	return groups, nil
}

func Get_groups_by_year(year int) (Groups, error) {

	db, err := DB_connection(Get_db_env("getter"))
	if err != nil {
		return Groups{}, err
	}
	defer db.Close()

	result, err := db.Query(fmt.Sprint("SELECT * FROM groups WHERE yearNumber = ", year, ";"))
	if err != nil {
		return Groups{}, err
	}
	defer result.Close()

	groups := Groups{}

	for result.Next() {

		group := Group{}

		err := result.Scan(&group.Id, &group.YearNumber, &group.Character)
		if err != nil {
			return Groups{}, err
		}

		groups = append(groups, group)
	}

	return groups, nil
}

func Add_groups(buff []byte) error {

	groups := Groups{}
	err := json.Unmarshal(buff, &groups)
	if err != nil {
		return err
	}

	data_str := ""

	for i, group := range groups {
		data_str += fmt.Sprint("(", group.YearNumber, ", '", group.Character, "')")
		if i < len(groups)-1 {
			data_str += ", "
		}
	}

	db, err := DB_connection(Get_db_env("setter"))
	if err != nil {
		return err
	}
	defer db.Close()

	insert_string := "INSERT INTO groups (yearNumber, character) VALUES " + data_str + " ON CONFLICT (yearNumber, character) DO UPDATE SET yearNumber = EXCLUDED.yearNumber, character = EXCLUDED.character;"
	_, err = db.Exec(insert_string)
	if err != nil {
		return err
	}

	return nil
}
