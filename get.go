package schedule_db_data

import (
	"database/sql"
	"fmt"
	"time"
)

type Request struct {
	Is_base  bool      `json:"is_base"`
	Start    time.Time `json:"start"`
	Year     int       `json:"year"`
	Parallel int       `json:"parallel"`
}

type Response struct {
	Groups   Groups   `json:"groups"`
	Teachers Teachers `json:"teachers"`
	Rooms    Rooms    `json:"rooms"`
	Subjects Subjects `json:"subjects"`
	Days     Days     `json:"days"`
}

// Универсальная функция для получения структуры ответа по структуре запроса
func Get_editor_data(req Request) (Response, error) {

	groups, err := Get_groups_by_year(req.Parallel)
	if err != nil {
		return Response{}, err
	}

	teachers, err := Get_teachers()
	if err != nil {
		return Response{}, err
	}

	rooms, err := Get_rooms()
	if err != nil {
		return Response{}, err
	}

	subjects, err := Get_subjects()
	if err != nil {
		return Response{}, err
	}

	// Подключаемся к БД
	db, err := DB_connection(Get_db_env("getter"))
	if err != nil {
		return Response{}, err
	}
	defer db.Close()

	var buff []byte

	//TODO: сканировать в буффер, после чего использовать Get_days() для парсинга
	if req.Is_base {
		err = db.QueryRow("SELECT data FROM schedule WHERE is_base = True AND year = " + fmt.Sprint(req.Year) + " AND parallel = " + fmt.Sprint(req.Parallel) + ";").Scan(&buff)
		if err == sql.ErrNoRows {
			return Response{Teachers: teachers, Groups: groups, Rooms: rooms, Subjects: subjects, Days: Days{}}, err
		} else if err != nil {
			return Response{}, err
		}
	} else {
		err = db.QueryRow("SELECT data FROM schedule WHERE is_base = False AND start = '" + req.Start.Format("2006-01-02") + "' AND parallel = " + fmt.Sprint(req.Parallel) + ";").Scan(&buff)
		if err == sql.ErrNoRows {
			err = db.QueryRow("SELECT data FROM schedule WHERE is_base = True AND year = " + fmt.Sprint(req.Year) + " AND parallel = " + fmt.Sprint(req.Parallel) + ";").Scan(&buff)
			if err == sql.ErrNoRows {
				return Response{Teachers: teachers, Groups: groups, Rooms: rooms, Subjects: subjects, Days: Days{}}, err
			} else if err != nil {
				return Response{}, err
			}
		} else if err != nil {
			return Response{}, err
		}
	}

	days, err := Get_days(buff)
	if err != nil {
		return Response{}, err
	}

	// Проверка существования классов, учителей, кабинетов и предметов
	for day_id, day := range days {
		for schedule_id, schedule := range day.Schedule {

			// Ищем класс в данных из БД по номеру и букве
			found_group, err := groups.Find(schedule.Group)
			if err != nil {
				// Отбрасываем расписание для несуществующего класса
				day.Schedule = append(day.Schedule[:schedule_id], day.Schedule[schedule_id+1:]...)
			} else {

				schedule.Group = found_group.ToString()
				for lesson_id, lesson := range schedule.Lessons {

					for l_data_id := range lesson.Lesson_data {

						// Проверяем наличие предмета, кабинета и учителя в данных из БД
						if subjects.Contain(lesson.Subject) && rooms.Contain(lesson.Lesson_data[l_data_id].Room) && teachers.Contain(lesson.Lesson_data[l_data_id].Teacher) {

							// Передаём учителю данные из БД
							days[day_id].Schedule[schedule_id].Lessons[lesson_id].Lesson_data[l_data_id].Teacher, err = teachers.Find_by_login(lesson.Lesson_data[l_data_id].Teacher.Login)
							if err != nil {
								// Отбрасываем урок, если в БД нет данных об учителе
								days[day_id].Schedule[schedule_id].Lessons = append(schedule.Lessons[:lesson_id], schedule.Lessons[lesson_id+1:]...)
							}
						} else {
							// Отбрасываем урок, если в БД нет данных о предмете или кабинете, или учителе
							days[day_id].Schedule[schedule_id].Lessons = append(schedule.Lessons[:lesson_id], schedule.Lessons[lesson_id+1:]...)
						}
					}

				}
			}
		}
	}

	return Response{Teachers: teachers, Groups: groups, Rooms: rooms, Subjects: subjects, Days: days}, nil
}
