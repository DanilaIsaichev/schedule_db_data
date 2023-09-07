package schedule_db_data

import (
	"errors"
	"fmt"
	"time"
)

type Lesson_time struct {
	Id   int       `json:"id"`
	Time time.Time `jsqon: "time"`
}

type Timetable []Lesson_time

func (t *Timetable) Get_lesson_time_by_number(id int) (time.Time, error) {

	if id > 0 {
		return (*t)[id-1].Time, nil
	} else {
		return time.Time{}, errors.New(fmt.Sprint("Lesson's number can't be ", id))
	}
}

func Get_timetable() (Timetable, error) {
	db, err := DB_connection(Get_db_env("getter"))
	if err != nil {
		return Timetable{}, err
	}
	defer db.Close()

	result, err := db.Query("SELECT * FROM timetable;")
	if err != nil {
		return Timetable{}, err
	}
	defer result.Close()

	timetable := Timetable{}

	for result.Next() {

		lesson_time := Lesson_time{}

		err := result.Scan(&lesson_time.Id, &lesson_time.Time)
		if err != nil {
			return Timetable{}, err
		}

		timetable = append(timetable, lesson_time)
	}

	return timetable, nil
}

func Get_time_by_number(id int) (Lesson_time, error) {
	db, err := DB_connection(Get_db_env("getter"))
	if err != nil {
		return Lesson_time{}, err
	}
	defer db.Close()

	result, err := db.Query(fmt.Sprint("SELECT * FROM timetable WHERE id = ", id, ";"))
	if err != nil {
		return Lesson_time{}, err
	}
	defer result.Close()

	lesson_time := Lesson_time{}

	err = result.Scan(&lesson_time.Id, &lesson_time.Time)
	if err != nil {
		return Lesson_time{}, err
	}

	return lesson_time, nil
}
