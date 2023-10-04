package schedule_db_data

import (
	"errors"
	"fmt"
	"time"
)

type Generator_request struct {
	Is_base bool      `json:"is_base"`
	Start   time.Time `json:"start"`
	Year    int       `json:"year"`
}

type Generator_class_lesson struct {
	Time    time.Time `json:"time"`
	Room    Room      `json:"room`
	Subject Subject   `json:"subject"`
	Teacher Teacher   `json:"teacher"`
}

type Generator_class_lessons []Generator_class_lesson

func (cls *Generator_classes) Get_class_lessons(class Class) (Generator_class, error) {

	for _, cl := range *cls {
		if cl.Class.Character == class.Character && cl.Class.Number == cl.Class.Number {
			return cl, nil
		}
	}

	return Generator_class{}, errors.New(fmt.Sprint("No lessons for ", class.Number, class.Character, " has found"))
}

type Generator_class struct {
	Class   Class                   `json:"class"`
	Lessons Generator_class_lessons `json:"lessons"`
}

type Generator_classes []Generator_class

type Generator_teacher_lesson struct {
	Time    time.Time `json:"time"`
	Room    Room      `json:"room`
	Subject Subject   `json:"subject"`
	Class   Class     `json:"class"`
}

type Generator_teacher_lessons []Generator_class_lesson

type Generator_teacher struct {
	Teacher Teacher                   `json:"teacher"`
	Lessons Generator_teacher_lessons `json:"lessons"`
}

type Generator_teachers []Generator_teacher

type Generator_response struct {
	Classes_lessons  Generator_classes          `json:"classes_lessons"`
	Teachers_lessons map[string]Teacher_lessons `json:"teachers_lessons"`
}

type Class_lesson struct {
	Time    time.Time `json:"time"`
	Name    string    `json:"name"`
	Room    Room      `json:"room"`
	Teacher Teacher   `json:"teacher"`
}

type Class_lessons []Class_lesson

type Teacher_lesson struct {
	Time  time.Time `json:"time"`
	Name  string    `json:"name"`
	Room  Room      `json:"room"`
	Class Class     `json:"class"`
}

type Teacher_lessons []Teacher_lesson

// Универсальная функция для получения структуры ответа по структуре запроса
func Get_generator_data(req Generator_request) (Generator_response, error) {

	classes, err := Get_classes()
	if err != nil {
		return Generator_response{}, err
	}

	classes_schedules := Generator_classes{}

	for _, class := range classes {
		classes_schedules = append(classes_schedules, Generator_class{Class: class})
	}

	teachers, err := Get_teachers()
	if err != nil {
		return Generator_response{}, err
	}

	teachers_schedules := Generator_teachers{}

	for _, teacher := range teachers {
		teachers_schedules = append(teachers_schedules, Generator_teacher{Teacher: teacher})
	}

	rooms, err := Get_rooms()
	if err != nil {
		return Generator_response{}, err
	}

	subjects, err := Get_subjects()
	if err != nil {
		return Generator_response{}, err
	}

	timetable, err := Get_timetable()
	if err != nil {
		return Generator_response{}, err
	}

	// Подключаемся к БД
	db, err := DB_connection(Get_db_env("getter"))
	if err != nil {
		return Generator_response{}, err
	}
	defer db.Close()

	days := Days{}

	select_change_str := "SELECT data FROM schedule WHERE is_base = False AND start = '" + req.Start.Format("2006-01-02") + "';"

	select_base_str := "SELECT data FROM schedule WHERE is_base = True AND year = " + fmt.Sprint(req.Year) + ";"

	if req.Is_base {
		err = db.QueryRow(select_base_str).Scan(&days)
		if err != nil {
			return Generator_response{}, err
		}
	} else {

		select_str := "IF EXISTS (" + select_change_str + ") THEN ELSE " + select_base_str + " END IF;"

		/*err = db.QueryRow("SELECT data FROM schedule WHERE is_base = False AND start = '" + req.Start.Format("2006-01-02") + "' AND parallel = " + fmt.Sprint(req.Parallel) + ";").Scan(&days)
		if err == sql.ErrNoRows {
			err = db.QueryRow("SELECT data FROM schedule WHERE is_base = True AND year = " + fmt.Sprint(req.Year) + " AND parallel = " + fmt.Sprint(req.Parallel) + ";").Scan(&days)
			if err != nil {
				return Generator_response{}, err
			}
		} else if err != nil {
			return Generator_response{}, err
		}*/
		err = db.QueryRow(select_str).Scan(&days)
		if err != nil {
			return Generator_response{}, err
		}
	}

	classes_schedules := Generator_classes{}

	// Проверка существования классов, учителей, кабинетов и предметов
	for _, day := range days {
		for schedule_id, schedule := range day.Schedule {

			// Ищем класс в данных из БД по номеру и букве
			found_class, err := classes.Find(schedule.Class)
			if err != nil {
				// Отбрасываем расписание для несуществующего класса
				day.Schedule = append(day.Schedule[:schedule_id], day.Schedule[schedule_id+1:]...)
			} else {

				schedule.Class = found_class.ToString()

				for lesson_id, lesson := range schedule.Lessons {

					// Проверяем наличие предмета, кабинета и учителя в данных из БД
					if subjects.Contain(Subject{Name: lesson.Name}) && rooms.Contain(Room{Name: lesson.Room}) && teachers.Contain(lesson.Teacher) {
						// Передаём учителю данные из БД

						if lesson.Number < 1 || lesson.Number > 10 {
							return Generator_response{}, errors.New("wrong lesson number")
						}

						l_time, err := timetable.Get_lesson_time_by_number(lesson.Number)
						if err != nil {
							return Generator_response{}, err
						}

						l_subject, err := subjects.Find(lesson.Name)
						if err != nil {
							return Generator_response{}, err
						}

						l_room, err := rooms.Find(lesson.Room)
						if err != nil {
							return Generator_response{}, err
						}

						l_teacher, err := teachers.Find(lesson.Teacher.Login)
						if err != nil {
							return Generator_response{}, err
						}

						l_class, err := classes.Find(schedule.Class)
						if err != nil {
							return Generator_response{}, err
						}

						class_lesson := Generator_class_lesson{Time: l_time, Subject: l_subject, Room: l_room, Teacher: l_teacher}

						class_schedule, err := classes_schedules.Get_class_lessons(l_class)
						if err != nil {
							return Generator_response{}, err
						}

						class_schedule.Lessons = append(class_schedule.Lessons, class_lesson)

						teacher_lesson := Generator_teacher_lesson{Time: l_time, Subject: l_subject, Room: l_room, Class: l_class}

						teacher_schedule, err := teachers_schedules.Get_teacher_lessons(l_teacher)
						if err != nil {
							return Generator_response{}, err
						}

						teacher_schedule.Lessons = append(teacher_schedule.Lessons, teacher_lesson)

					} else {
						// Отбрасываем урок, если в БД нет данных о предмете или кабинете, или учителе
						schedule.Lessons = append(schedule.Lessons[:lesson_id], schedule.Lessons[lesson_id+1:]...)
					}
				}
			}
		}
	}

	return Generator_response{Classes: classes, Teachers: teachers, Classes_lessons: classes_schedules, Teachers_lessons: teachers_schedules}, nil
}

// Переработать маршрут получения данных для генератора так, чтобы он сам получал все расписания для всех параллелей, и раскладывал по картам
