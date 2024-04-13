package service

import (
	dto "classrooms/internal/DTO"
	"classrooms/internal/models"
	"classrooms/internal/repository"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

type ScheduleService struct {
	RoomRepo   repository.Room
	LessonRepo repository.Lesson
}

func NewScheduleService(roomRepo repository.Room, lessonRepo repository.Lesson) *ScheduleService {
	return &ScheduleService{RoomRepo: roomRepo, LessonRepo: lessonRepo}
}

type ScheduleItem struct {
	Time_start string
	Time_end   string
	Name       string
	Groups     []string
	Types      []string
	Rooms      map[string]string
}

type ScheduleDate struct {
	Date  string
	Day   string
	Pairs map[string]ScheduleItem
}

type LectorSchedule struct {
	Name     string
	Groups   map[string]int
	Schedule map[string]ScheduleDate
}

func (s *ScheduleService) GetScheduleByRoomIdAndDate(roomId int, date string) ([]dto.ScheduleLessonDto, error) {

	//var lessons []models.Lesson

	// room, err := s.RoomRepo.GetRoomById(roomId)
	// if err != nil {
	// 	return nil, err
	// }

	parsedDate, err := time.Parse("02.01.2006", date)

	lessons, err := s.LessonRepo.GetScheduleLessons(roomId, parsedDate)

	if err != nil {

		return nil, err
	}

	lessonsDTO := make([]dto.ScheduleLessonDto, len(lessons))

	for index, lesson := range lessons {
		lessonsDTO[index].Lector = lesson.Lector
		lessonsDTO[index].Groups = lesson.Groups
		lessonsDTO[index].Subject = lesson.Name
		lessonsDTO[index].Type = lesson.Type
		lessonsDTO[index].StartTime = lesson.TimeStart[:5]
		lessonsDTO[index].EndTime = lesson.TimeEnd[:5]
	}
	return lessonsDTO, nil
	// folderPath := "../internal/localStorage/LectorsDataJSON/"
	//
	// files, err := os.ReadDir(folderPath)
	// if err != nil {
	// 	fmt.Println("Error reading directory:", err)
	// 	return nil, fmt.Errorf("ошибка при получении данных")
	// }

	// for _, file := range files {
	// 	filePath := folderPath + file.Name()
	// 	fileContent, err := os.ReadFile(filePath)
	// 	if err != nil {
	// 		fmt.Println("Error reading file:", err)
	// 		continue
	// 	}
	// 	var schedule LectorSchedule
	// 	if err := json.Unmarshal(fileContent, &schedule); err != nil {
	// 		fmt.Println("Ошибка при декодировании JSON:", err)
	// 		continue
	// 	}

	// 	var lesson models.Lesson
	// 	if day, ok := schedule.Schedule[date]; ok {
	// 		for _, pair := range day.Pairs {
	// 			for _, roomName := range pair.Rooms {
	// 				if roomName == room.Name {
	// 					// lesson.Room = roomName
	// 					// lesson.Date = date
	// 					// parsedTime, _ := time.Parse("15:04:00", pair.Time_start)
	// 					// lesson.TimeStart = parsedTime.Format("15:04")
	// 					// parsedTime, _ = time.Parse("15:04:00", pair.Time_end)
	// 					// lesson.TimeEnd = parsedTime.Format("15:04")
	// 					// lesson.Groups = pair.Groups
	// 					lesson.Name = pair.Name
	// 					lesson.Type = pair.Types[0]
	// 					lesson.Lector = schedule.Name
	// 					lessons = append(lessons, lesson)
	// 				}

	// 			}

	// 		}
	// 	}

	// }

	//return lessons, nil
}

func (s *ScheduleService) WriteScheduleToDB(lectorSchedule LectorSchedule) error {

	lectorName := lectorSchedule.Name

	for _, schedule := range lectorSchedule.Schedule {
		//day, _ := time.Parse("02.01.2006", date)
		for _, pair := range schedule.Pairs {
			//startTime, _ := time.Parse("02.01.2006 15:04:05", date+" "+pair.Time_start)
			//endTime, _ := time.Parse("02.01.2006 15:04:05", date+" "+pair.Time_end)
			subject := pair.Name
			groups := strings.Join(pair.Groups, ", ")
			types := pair.Types
			for _, roomName := range pair.Rooms {

				if !strings.Contains(roomName, "-") || roomName[0] == '-' {
					continue
				}

				roomId, _ := s.RoomRepo.GetRoomId(roomName)

				lesson := models.Lesson{Lector: lectorName, TimeStart: "startTime",
					TimeEnd: "endTime", Name: subject, Groups: groups, Type: types[0], RoomId: roomId}

				err := s.LessonRepo.InsertLessonToDB(lesson)
				if err != nil {
					return err
				}

			}

		}
	}

	return nil
}

func (s *ScheduleService) RefreshLectorSchedule() error {
	s.DeleteAllLessons()

	filePathWithLectorsInfo := "../internal/localStorage/allLectors.json"

	fileContent, err := os.ReadFile(filePathWithLectorsInfo)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return err
	}

	var lectors map[string]string
	if err := json.Unmarshal(fileContent, &lectors); err != nil {
		fmt.Println("Ошибка при декодировании JSON:", err)
	}

	for lectorHash, _ := range lectors {

		fmt.Println(lectorHash)
		lectorSchedule, err := GetLectorSchedule(lectorHash)
		if err != nil {
			fmt.Println(err)
		}

		err = s.WriteScheduleToDB(lectorSchedule)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *ScheduleService) DeleteAllLessons() error {
	return s.LessonRepo.DeleteAll()
}

func GetLectorSchedule(LectorHash string) (LectorSchedule, error) {
	response, err := http.Get(fmt.Sprintf("https://public.mai.ru/schedule/data/%s.json", LectorHash))
	if err != nil {
		fmt.Printf("Ошибка при выполнении запроса: %s\n", err)
		return LectorSchedule{}, err
	}

	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return LectorSchedule{}, err
	}

	var lectorSchedule LectorSchedule
	if err := json.Unmarshal(body, &lectorSchedule); err != nil {
		return LectorSchedule{}, err
	}

	return lectorSchedule, nil
}
