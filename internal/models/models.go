package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// User — пользователь системы
type User struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey"`
	Email     string    `gorm:"type:varchar(255);unique;not null"`
	Nickname  string    `gorm:"type:varchar(255)"`
	Password  string    `gorm:"type:varchar(255);not null"`
	Role      string    `gorm:"type:varchar(50);default:user"`
	Resumes   []Resume  `gorm:"foreignKey:UserID"`
	Vacancies []Vacancy `gorm:"foreignKey:UserID"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	u.ID = uuid.New()
	return
}

// Resume — информация о загруженном резюме
type Resume struct {
	ID         uuid.UUID    `gorm:"type:uuid;primaryKey"`
	UserID     uuid.UUID    `gorm:"type:uuid;not null;index"`
	User       User         `gorm:"foreignKey:UserID"`
	FullName   string       `gorm:"type:varchar(255);not null"`
	Email      string       `gorm:"type:varchar(255)"`
	Phone      string       `gorm:"type:varchar(50)"`
	Location   string       `gorm:"type:varchar(255)"`
	Skills     []Skill      `gorm:"many2many:resume_skills;"`
	Experience []Experience `gorm:"foreignKey:ResumeID;constraint:OnDelete:CASCADE"`
	Education  []Education  `gorm:"foreignKey:ResumeID;constraint:OnDelete:CASCADE"`
	File       ResumeFile   `gorm:"foreignKey:ResumeID;constraint:OnDelete:CASCADE"`
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
	StartDate   string    `gorm:"type:varchar(32)"`
	EndDate     string    `gorm:"type:varchar(32)"`
	Description string    `gorm:"type:text"`
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
	StartDate   string    `gorm:"type:varchar(32)"`
	EndDate     string    `gorm:"type:varchar(32)"`
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
	UserID      uuid.UUID `gorm:"type:uuid;not null;index"`
	User        User      `gorm:"foreignKey:UserID"`
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
