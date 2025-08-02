package response

type ErrorResponse struct {
	Error string `json:"error"`
}

type UserRegisterResponse struct {
	ID       string `json:"id"`
	Nickname string `json:"nickname"`
	Email    string `json:"email"`
}

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type UserProfileResponse struct {
	ID       string `json:"id"`
	Nickname string `json:"nickname"`
	Email    string `json:"email"`
}

type ParsedResumeDTO struct {
	FullName   string          `json:"full_name"`
	Email      string          `json:"email"`
	Phone      string          `json:"phone"`
	Location   string          `json:"location"`
	Skills     []string        `json:"skills"`
	Experience []ExperienceDTO `json:"experience"`
	Education  []EducationDTO  `json:"education"`
}

type ExperienceDTO struct {
	Company     string `json:"company"`
	Position    string `json:"position"`
	StartDate   string `json:"start_date"`
	EndDate     string `json:"end_date"`
	Description string `json:"description"`
}

type EducationDTO struct {
	Institution string `json:"institution"`
	Degree      string `json:"degree"`
	Field       string `json:"field"`
	StartDate   string `json:"start_date"`
	EndDate     string `json:"end_date"`
}
