package service

import (
	"arbak-portfolio-backend/internal/models"
	"arbak-portfolio-backend/internal/repository"
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/redis/go-redis/v9"
)

// ProfileService handles profile-related operations
type ProfileService struct {
	repo  *repository.ProfileRepository
	redis *redis.Client
}

func NewProfileService(repo *repository.ProfileRepository, redis *redis.Client) *ProfileService {
	return &ProfileService{
		repo:  repo,
		redis: redis,
	}
}

func (s *ProfileService) GetProfile() (*models.Profile, error) {
	// Try to get from cache first
	ctx := context.Background()
	cached, err := s.redis.Get(ctx, "profile").Result()
	if err == nil {
		var profile models.Profile
		if err := json.Unmarshal([]byte(cached), &profile); err == nil {
			return &profile, nil
		}
	}

	// Get from database
	profile, err := s.repo.GetProfile()
	if err != nil {
		return nil, err
	}

	// Cache the result
	profileJSON, _ := json.Marshal(profile)
	s.redis.Set(ctx, "profile", profileJSON, time.Hour)

	return profile, nil
}

type ProfileUpdateRequest struct {
	Name      string `json:"name" binding:"required"`
	Title     string `json:"title" binding:"required"`
	Location  string `json:"location"`
	Email     string `json:"email" binding:"required,email"`
	Phone     string `json:"phone"`
	Telegram  string `json:"telegram"`
	GitHub    string `json:"github"`
	LinkedIn  string `json:"linkedin"`
	Summary   string `json:"summary"`
	Avatar    string `json:"avatar"`
	ResumeURL string `json:"resume_url"`
}

func (s *ProfileService) UpdateProfile(req *ProfileUpdateRequest) (*models.Profile, error) {
	profile := &models.Profile{
		Name:      req.Name,
		Title:     req.Title,
		Location:  req.Location,
		Email:     req.Email,
		Phone:     req.Phone,
		Telegram:  req.Telegram,
		GitHub:    req.GitHub,
		LinkedIn:  req.LinkedIn,
		Summary:   req.Summary,
		Avatar:    req.Avatar,
		ResumeURL: req.ResumeURL,
	}

	updatedProfile, err := s.repo.UpdateProfile(profile)
	if err != nil {
		return nil, err
	}

	// Invalidate cache
	ctx := context.Background()
	s.redis.Del(ctx, "profile")

	return updatedProfile, nil
}

// ExperienceService handles experience-related operations
type ExperienceService struct {
	repo  *repository.ExperienceRepository
	redis *redis.Client
}

func NewExperienceService(repo *repository.ExperienceRepository, redis *redis.Client) *ExperienceService {
	return &ExperienceService{
		repo:  repo,
		redis: redis,
	}
}

func (s *ExperienceService) GetExperiences() ([]models.Experience, error) {
	// Try to get from cache first
	ctx := context.Background()
	cached, err := s.redis.Get(ctx, "experiences").Result()
	if err == nil {
		var experiences []models.Experience
		if err := json.Unmarshal([]byte(cached), &experiences); err == nil {
			return experiences, nil
		}
	}

	// Get from database
	experiences, err := s.repo.GetExperiences()
	if err != nil {
		return nil, err
	}

	// Cache the result
	experiencesJSON, _ := json.Marshal(experiences)
	s.redis.Set(ctx, "experiences", experiencesJSON, time.Hour)

	return experiences, nil
}

type ExperienceCreateRequest struct {
	Company      string     `json:"company" binding:"required"`
	Position     string     `json:"position" binding:"required"`
	Location     string     `json:"location"`
	StartDate    time.Time  `json:"start_date" binding:"required"`
	EndDate      *time.Time `json:"end_date"`
	Current      bool       `json:"current"`
	Description  string     `json:"description"`
	Achievements []string   `json:"achievements"`
	Technologies []string   `json:"technologies"`
}

func (s *ExperienceService) CreateExperience(req *ExperienceCreateRequest) (*models.Experience, error) {
	experience := &models.Experience{
		Company:      req.Company,
		Position:     req.Position,
		Location:     req.Location,
		StartDate:    req.StartDate,
		EndDate:      req.EndDate,
		Current:      req.Current,
		Description:  req.Description,
		Achievements: req.Achievements,
		Technologies: req.Technologies,
	}

	createdExperience, err := s.repo.CreateExperience(experience)
	if err != nil {
		return nil, err
	}

	// Invalidate cache
	ctx := context.Background()
	s.redis.Del(ctx, "experiences")

	return createdExperience, nil
}

type ExperienceUpdateRequest struct {
	Company      string     `json:"company"`
	Position     string     `json:"position"`
	Location     string     `json:"location"`
	StartDate    time.Time  `json:"start_date"`
	EndDate      *time.Time `json:"end_date"`
	Current      bool       `json:"current"`
	Description  string     `json:"description"`
	Achievements []string   `json:"achievements"`
	Technologies []string   `json:"technologies"`
}

func (s *ExperienceService) UpdateExperience(id uint, req *ExperienceUpdateRequest) (*models.Experience, error) {
	experience := &models.Experience{
		Company:      req.Company,
		Position:     req.Position,
		Location:     req.Location,
		StartDate:    req.StartDate,
		EndDate:      req.EndDate,
		Current:      req.Current,
		Description:  req.Description,
		Achievements: req.Achievements,
		Technologies: req.Technologies,
	}

	updatedExperience, err := s.repo.UpdateExperience(id, experience)
	if err != nil {
		return nil, err
	}

	// Invalidate cache
	ctx := context.Background()
	s.redis.Del(ctx, "experiences")

	return updatedExperience, nil
}

func (s *ExperienceService) DeleteExperience(id uint) error {
	err := s.repo.DeleteExperience(id)
	if err != nil {
		return err
	}

	// Invalidate cache
	ctx := context.Background()
	s.redis.Del(ctx, "experiences")

	return nil
}

// SkillService handles skill-related operations
type SkillService struct {
	repo  *repository.SkillRepository
	redis *redis.Client
}

func NewSkillService(repo *repository.SkillRepository, redis *redis.Client) *SkillService {
	return &SkillService{
		repo:  repo,
		redis: redis,
	}
}

func (s *SkillService) GetSkills() ([]models.Skill, error) {
	// Try to get from cache first
	ctx := context.Background()
	cached, err := s.redis.Get(ctx, "skills").Result()
	if err == nil {
		var skills []models.Skill
		if err := json.Unmarshal([]byte(cached), &skills); err == nil {
			return skills, nil
		}
	}

	// Get from database
	skills, err := s.repo.GetSkills()
	if err != nil {
		return nil, err
	}

	// Cache the result
	skillsJSON, _ := json.Marshal(skills)
	s.redis.Set(ctx, "skills", skillsJSON, time.Hour)

	return skills, nil
}

type SkillCreateRequest struct {
	Name        string `json:"name" binding:"required"`
	Category    string `json:"category" binding:"required"`
	Level       int    `json:"level" binding:"min=1,max=10"`
	Description string `json:"description"`
	Icon        string `json:"icon"`
}

func (s *SkillService) CreateSkill(req *SkillCreateRequest) (*models.Skill, error) {
	skill := &models.Skill{
		Name:        req.Name,
		Category:    req.Category,
		Level:       req.Level,
		Description: req.Description,
		Icon:        req.Icon,
	}

	createdSkill, err := s.repo.CreateSkill(skill)
	if err != nil {
		return nil, err
	}

	// Invalidate cache
	ctx := context.Background()
	s.redis.Del(ctx, "skills")

	return createdSkill, nil
}

type SkillUpdateRequest struct {
	Name        string `json:"name"`
	Category    string `json:"category"`
	Level       int    `json:"level" binding:"min=1,max=10"`
	Description string `json:"description"`
	Icon        string `json:"icon"`
}

func (s *SkillService) UpdateSkill(id uint, req *SkillUpdateRequest) (*models.Skill, error) {
	skill := &models.Skill{
		Name:        req.Name,
		Category:    req.Category,
		Level:       req.Level,
		Description: req.Description,
		Icon:        req.Icon,
	}

	updatedSkill, err := s.repo.UpdateSkill(id, skill)
	if err != nil {
		return nil, err
	}

	// Invalidate cache
	ctx := context.Background()
	s.redis.Del(ctx, "skills")

	return updatedSkill, nil
}

func (s *SkillService) DeleteSkill(id uint) error {
	err := s.repo.DeleteSkill(id)
	if err != nil {
		return err
	}

	// Invalidate cache
	ctx := context.Background()
	s.redis.Del(ctx, "skills")

	return nil
}

// ProjectService handles project-related operations
type ProjectService struct {
	repo  *repository.ProjectRepository
	redis *redis.Client
}

func NewProjectService(repo *repository.ProjectRepository, redis *redis.Client) *ProjectService {
	return &ProjectService{
		repo:  repo,
		redis: redis,
	}
}

func (s *ProjectService) GetProjects(featured *bool) ([]models.Project, error) {
	// Try to get from cache first
	ctx := context.Background()
	cacheKey := "projects"
	if featured != nil {
		if *featured {
			cacheKey = "projects:featured"
		} else {
			cacheKey = "projects:non-featured"
		}
	}

	cached, err := s.redis.Get(ctx, cacheKey).Result()
	if err == nil {
		var projects []models.Project
		if err := json.Unmarshal([]byte(cached), &projects); err == nil {
			return projects, nil
		}
	}

	// Get from database
	projects, err := s.repo.GetProjects(featured)
	if err != nil {
		return nil, err
	}

	// Cache the result
	projectsJSON, _ := json.Marshal(projects)
	s.redis.Set(ctx, cacheKey, projectsJSON, time.Hour)

	return projects, nil
}

type ProjectCreateRequest struct {
	Name            string   `json:"name" binding:"required"`
	Description     string   `json:"description" binding:"required"`
	LongDescription string   `json:"long_description"`
	Technologies    []string `json:"technologies"`
	GitHubURL       string   `json:"github_url"`
	LiveURL         string   `json:"live_url"`
	ImageURL        string   `json:"image_url"`
	Featured        bool     `json:"featured"`
	Category        string   `json:"category"`
	Status          string   `json:"status"`
}

func (s *ProjectService) CreateProject(req *ProjectCreateRequest) (*models.Project, error) {
	project := &models.Project{
		Name:            req.Name,
		Description:     req.Description,
		LongDescription: req.LongDescription,
		Technologies:    req.Technologies,
		GitHubURL:       req.GitHubURL,
		LiveURL:         req.LiveURL,
		ImageURL:        req.ImageURL,
		Featured:        req.Featured,
		Category:        req.Category,
		Status:          req.Status,
	}

	createdProject, err := s.repo.CreateProject(project)
	if err != nil {
		return nil, err
	}

	// Invalidate cache
	ctx := context.Background()
	s.redis.Del(ctx, "projects", "projects:featured", "projects:non-featured")

	return createdProject, nil
}

type ProjectUpdateRequest struct {
	Name            string   `json:"name"`
	Description     string   `json:"description"`
	LongDescription string   `json:"long_description"`
	Technologies    []string `json:"technologies"`
	GitHubURL       string   `json:"github_url"`
	LiveURL         string   `json:"live_url"`
	ImageURL        string   `json:"image_url"`
	Featured        bool     `json:"featured"`
	Category        string   `json:"category"`
	Status          string   `json:"status"`
}

func (s *ProjectService) UpdateProject(id uint, req *ProjectUpdateRequest) (*models.Project, error) {
	project := &models.Project{
		Name:            req.Name,
		Description:     req.Description,
		LongDescription: req.LongDescription,
		Technologies:    req.Technologies,
		GitHubURL:       req.GitHubURL,
		LiveURL:         req.LiveURL,
		ImageURL:        req.ImageURL,
		Featured:        req.Featured,
		Category:        req.Category,
		Status:          req.Status,
	}

	updatedProject, err := s.repo.UpdateProject(id, project)
	if err != nil {
		return nil, err
	}

	// Invalidate cache
	ctx := context.Background()
	s.redis.Del(ctx, "projects", "projects:featured", "projects:non-featured")

	return updatedProject, nil
}

func (s *ProjectService) DeleteProject(id uint) error {
	err := s.repo.DeleteProject(id)
	if err != nil {
		return err
	}

	// Invalidate cache
	ctx := context.Background()
	s.redis.Del(ctx, "projects", "projects:featured", "projects:non-featured")

	return nil
}

// ContactService handles contact-related operations
type ContactService struct {
	repo  *repository.ContactRepository
	redis *redis.Client
}

func NewContactService(repo *repository.ContactRepository, redis *redis.Client) *ContactService {
	return &ContactService{
		repo:  repo,
		redis: redis,
	}
}

type ContactCreateRequest struct {
	Name      string `json:"name" binding:"required"`
	Email     string `json:"email" binding:"required,email"`
	Subject   string `json:"subject"`
	Message   string `json:"message" binding:"required"`
	IPAddress string `json:"ip_address"`
	UserAgent string `json:"user_agent"`
}

type ContactStatusUpdateRequest struct {
	Status string `json:"status" binding:"required"`
}

func (s *ContactService) CreateContact(req *ContactCreateRequest) (*models.Contact, error) {
	contact := &models.Contact{
		Name:      req.Name,
		Email:     req.Email,
		Subject:   req.Subject,
		Message:   req.Message,
		IPAddress: req.IPAddress,
		UserAgent: req.UserAgent,
		Status:    "new",
	}

	createdContact, err := s.repo.CreateContact(contact)
	if err != nil {
		return nil, err
	}

	return createdContact, nil
}

func (s *ContactService) GetContacts() ([]models.Contact, error) {
	return s.repo.GetContacts()
}

func (s *ContactService) UpdateContactStatus(id uint, status string) (*models.Contact, error) {
	return s.repo.UpdateContactStatus(id, status)
}

// AuthService handles authentication-related operations
type AuthService struct {
	jwtSecret string
}

func NewAuthService(jwtSecret string) *AuthService {
	return &AuthService{
		jwtSecret: jwtSecret,
	}
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	Token string `json:"token"`
	User  struct {
		ID       uint   `json:"id"`
		Username string `json:"username"`
		Email    string `json:"email"`
		Role     string `json:"role"`
	} `json:"user"`
}

func (s *AuthService) Login(req *LoginRequest) (*LoginResponse, error) {
	// This is a simplified implementation
	// In a real application, you would:
	// 1. Hash the password
	// 2. Compare with stored hash
	// 3. Generate JWT token

	// For demo purposes, accept any username/password
	if req.Username == "" || req.Password == "" {
		return nil, errors.New("invalid credentials")
	}

	// Generate JWT token (simplified)
	token := "demo-jwt-token-" + req.Username

	response := &LoginResponse{
		Token: token,
		User: struct {
			ID       uint   `json:"id"`
			Username string `json:"username"`
			Email    string `json:"email"`
			Role     string `json:"role"`
		}{
			ID:       1,
			Username: req.Username,
			Email:    "admin@example.com",
			Role:     "admin",
		},
	}

	return response, nil
}
