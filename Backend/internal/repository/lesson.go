package repository

import (
	"classrooms/internal/models"
	"errors"
	"time"

	"gorm.io/gorm"
)

type LessonPostgres struct {
	db *gorm.DB
}

func NewLessonPostgres(db *gorm.DB) *LessonPostgres {
	return &LessonPostgres{db: db}
}

func (l *LessonPostgres) GetScheduleLessons(roomId int, date time.Time) ([]models.Lesson, error) {
	var lessons []models.Lesson

	result := l.db.Raw(`
		SELECT 
			MAX(lector) as lector,
			date,
			time_start,
			time_end,
			name,
			STRING_AGG(groups::TEXT, ', ') AS groups,
			type,
			room_id
		FROM 
			lessons
		WHERE 
			room_id = ? AND date = ?
		GROUP BY 
			date, time_start, time_end, name, type, room_id
	`, roomId, date).Scan(&lessons)

	if result.Error != nil {
		return nil, result.Error
	}
	return lessons, nil
}

func (l *LessonPostgres) InsertLessonToDB(lesson models.Lesson) error {
	l.db.AutoMigrate(&models.Lesson{})

	if l.db.Create(&lesson).Error != nil {
		return errors.New("ошибка при добавлении записи")
	}
	return nil
}

func (l *LessonPostgres) DeleteAll() error {
	err := l.db.Where("id > 0").Delete(&models.Lesson{}).Error

	if err != nil {
		return err
	}
	return nil
}
