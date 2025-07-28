package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Resume — информация о загруженном резюме
type Resume struct {
	ID         uuid.UUID    `gorm:"type:uuid;primaryKey"`
	FullName   string       `gorm:"type:varchar(255);not null"`
	Email      string       `gorm:"type:varchar(255)"`
	Phone      string       `gorm:"type:varchar(50)"`
	Location   string       `gorm:"type:varchar(255)"`
	Skills     []Skill      `gorm:"many2many:resume_skills;"`
	Experience []Experience `gorm:"foreignKey:ResumeID"`
	Education  []Education  `gorm:"foreignKey:ResumeID"`
	File       ResumeFile   `gorm:"foreignKey:ResumeID"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DeletedAt  gorm.DeletedAt `gorm:"index"`
}

func (m *Resume) BeforeCreate(tx *gorm.DB) (err error) {
	m.ID = uuid.New()
	return
}

// ResumeFile — файл с резюме
type ResumeFile struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey"`
	ResumeID  uuid.UUID `gorm:"type:uuid;not null;uniqueIndex"`
	Path      string    `gorm:"type:varchar(512);not null"`
	MimeType  string    `gorm:"type:varchar(100)"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (m *ResumeFile) BeforeCreate(tx *gorm.DB) (err error) {
	m.ID = uuid.New()
	return
}

// Skill — отдельный навык (используется для резюме и вакансии)
type Skill struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey"`
	Name      string    `gorm:"type:varchar(100);unique;not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (m *Skill) BeforeCreate(tx *gorm.DB) (err error) {
	m.ID = uuid.New()
	return
}

// Experience — опыт работы
type Experience struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey"`
	ResumeID    uuid.UUID `gorm:"type:uuid;not null;index"`
	Company     string    `gorm:"type:varchar(255)"`
	Position    string    `gorm:"type:varchar(255)"`
	StartDate   time.Time
	EndDate     *time.Time
	Description string `gorm:"type:text"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}

func (m *Experience) BeforeCreate(tx *gorm.DB) (err error) {
	m.ID = uuid.New()
	return
}

// Education — образование
type Education struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey"`
	ResumeID    uuid.UUID `gorm:"type:uuid;not null;index"`
	Institution string    `gorm:"type:varchar(255)"`
	Degree      string    `gorm:"type:varchar(255)"`
	Field       string    `gorm:"type:varchar(255)"`
	StartDate   time.Time
	EndDate     *time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}

func (m *Education) BeforeCreate(tx *gorm.DB) (err error) {
	m.ID = uuid.New()
	return
}

// Vacancy — вакансия (Job Description)
type Vacancy struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey"`
	Title       string    `gorm:"type:varchar(255);not null"`
	Description string    `gorm:"type:text"`
	Location    string    `gorm:"type:varchar(255)"`
	Skills      []Skill   `gorm:"many2many:vacancy_skills;"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}

func (m *Vacancy) BeforeCreate(tx *gorm.DB) (err error) {
	m.ID = uuid.New()
	return
}

// MatchingResult — результат сравнения резюме и вакансии
type MatchingResult struct {
	ID              uuid.UUID `gorm:"type:uuid;primaryKey"`
	ResumeID        uuid.UUID `gorm:"type:uuid;not null;index"`
	VacancyID       uuid.UUID `gorm:"type:uuid;not null;index"`
	Score           float64
	MatchedSkills   string `gorm:"type:text"` // JSON-строка с совпавшими навыками
	UnmatchedSkills string `gorm:"type:text"` // JSON-строка с несовпавшими
	Recommendations string `gorm:"type:text"` // JSON-строка с советами
	CreatedAt       time.Time
	UpdatedAt       time.Time
	DeletedAt       gorm.DeletedAt `gorm:"index"`
}

func (m *MatchingResult) BeforeCreate(tx *gorm.DB) (err error) {
	m.ID = uuid.New()
	return
}
