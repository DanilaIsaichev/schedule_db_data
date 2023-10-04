package schedule_db_data

import (
	"encoding/json"
	"errors"
	"fmt"
)

type Room struct {
	Id         int    `json:"id"`
	Name       string `json:"name"`
	Short_name string `json:"short_name"`
	Wing       int    `json:"wing"`
	Floor      int    `json:"floor"`
}

type Rooms []Room

func (rooms *Rooms) Contain(room Room) (res bool) {

	for _, r := range *rooms {
		if r.Name == room.Name {
			return true
		}
	}

	return false
}

func (rooms *Rooms) Find(name string) (room Room, err error) {

	for _, room := range *rooms {
		if name == room.Name {
			return room, nil
		}
	}

	return Room{}, errors.New("no room with name " + name + " has found")
}

func Get_rooms() (Rooms, error) {

	db, err := DB_connection(Get_db_env("getter"))
	if err != nil {
		return Rooms{}, err
	}
	defer db.Close()

	result, err := db.Query("SELECT * FROM room;")
	if err != nil {
		return Rooms{}, err
	}
	defer result.Close()

	rooms := Rooms{}

	for result.Next() {

		room := Room{}

		err := result.Scan(&room.Id, &room.Name, &room.Short_name, &room.Wing, &room.Floor)
		if err != nil {
			return Rooms{}, err
		}

		rooms = append(rooms, room)
	}

	return rooms, nil
}

func Add_rooms(buff []byte) error {
	rooms := Rooms{}
	err := json.Unmarshal(buff, &rooms)
	if err != nil {
		return err
	}

	data_str := ""

	for i, room := range rooms {
		data_str += fmt.Sprint("('", room.Name, "', '", room.Short_name, "', ", room.Wing, ", ", room.Floor, ")")
		if i < len(rooms)-1 {
			data_str += ", "
		}
	}

	db, err := DB_connection(Get_db_env("setter"))
	if err != nil {
		return err
	}
	defer db.Close()

	insert_string := "INSERT INTO room (name, short_name, wing, floor) VALUES " + data_str + " ON CONFLICT (name, short_name) DO UPDATE SET name = EXCLUDED.name, short_name = EXCLUDED.short_name, wing = EXCLUDED.wing, floor = EXCLUDED.floor;"
	_, err = db.Exec(insert_string)
	if err != nil {
		return err
	}

	return nil
}
