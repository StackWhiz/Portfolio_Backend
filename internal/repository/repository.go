package repository

import (
	"errors"
	"stackwhiz-portfolio-backend/internal/models"

	"gorm.io/gorm"
)

// ProfileRepository handles profile data operations
type ProfileRepository struct {
	db *gorm.DB
}

func NewProfileRepository(db *gorm.DB) *ProfileRepository {
	return &ProfileRepository{db: db}
}

func (r *ProfileRepository) GetProfile() (*models.Profile, error) {
	var profile models.Profile
	err := r.db.First(&profile).Error
	if err != nil {
		return nil, err
	}
	return &profile, nil
}

func (r *ProfileRepository) UpdateProfile(profile *models.Profile) (*models.Profile, error) {
	// Update or create profile
	err := r.db.Save(profile).Error
	if err != nil {
		return nil, err
	}
	return profile, nil
}

// ExperienceRepository handles experience data operations
type ExperienceRepository struct {
	db *gorm.DB
}

func NewExperienceRepository(db *gorm.DB) *ExperienceRepository {
	return &ExperienceRepository{db: db}
}

func (r *ExperienceRepository) GetExperiences() ([]models.Experience, error) {
	var experiences []models.Experience
	err := r.db.Order("start_date DESC").Find(&experiences).Error
	if err != nil {
		return nil, err
	}
	return experiences, nil
}

func (r *ExperienceRepository) CreateExperience(experience *models.Experience) (*models.Experience, error) {
	err := r.db.Create(experience).Error
	if err != nil {
		return nil, err
	}
	return experience, nil
}

func (r *ExperienceRepository) UpdateExperience(id uint, experience *models.Experience) (*models.Experience, error) {
	var existingExperience models.Experience
	err := r.db.First(&existingExperience, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("experience not found")
		}
		return nil, err
	}

	experience.ID = id
	err = r.db.Save(experience).Error
	if err != nil {
		return nil, err
	}
	return experience, nil
}

func (r *ExperienceRepository) DeleteExperience(id uint) error {
	var experience models.Experience
	err := r.db.First(&experience, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("experience not found")
		}
		return err
	}

	err = r.db.Delete(&experience).Error
	if err != nil {
		return err
	}
	return nil
}

// SkillRepository handles skill data operations
type SkillRepository struct {
	db *gorm.DB
}

func NewSkillRepository(db *gorm.DB) *SkillRepository {
	return &SkillRepository{db: db}
}

func (r *SkillRepository) GetSkills() ([]models.Skill, error) {
	var skills []models.Skill
	err := r.db.Order("category, name").Find(&skills).Error
	if err != nil {
		return nil, err
	}
	return skills, nil
}

func (r *SkillRepository) CreateSkill(skill *models.Skill) (*models.Skill, error) {
	err := r.db.Create(skill).Error
	if err != nil {
		return nil, err
	}
	return skill, nil
}

func (r *SkillRepository) UpdateSkill(id uint, skill *models.Skill) (*models.Skill, error) {
	var existingSkill models.Skill
	err := r.db.First(&existingSkill, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("skill not found")
		}
		return nil, err
	}

	skill.ID = id
	err = r.db.Save(skill).Error
	if err != nil {
		return nil, err
	}
	return skill, nil
}

func (r *SkillRepository) DeleteSkill(id uint) error {
	var skill models.Skill
	err := r.db.First(&skill, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("skill not found")
		}
		return err
	}

	err = r.db.Delete(&skill).Error
	if err != nil {
		return err
	}
	return nil
}

// ProjectRepository handles project data operations
type ProjectRepository struct {
	db *gorm.DB
}

func NewProjectRepository(db *gorm.DB) *ProjectRepository {
	return &ProjectRepository{db: db}
}

func (r *ProjectRepository) GetProjects(featured *bool) ([]models.Project, error) {
	var projects []models.Project
	query := r.db.Order("created_at DESC")

	if featured != nil {
		query = query.Where("featured = ?", *featured)
	}

	err := query.Find(&projects).Error
	if err != nil {
		return nil, err
	}
	return projects, nil
}

func (r *ProjectRepository) CreateProject(project *models.Project) (*models.Project, error) {
	err := r.db.Create(project).Error
	if err != nil {
		return nil, err
	}
	return project, nil
}

func (r *ProjectRepository) UpdateProject(id uint, project *models.Project) (*models.Project, error) {
	var existingProject models.Project
	err := r.db.First(&existingProject, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("project not found")
		}
		return nil, err
	}

	project.ID = id
	err = r.db.Save(project).Error
	if err != nil {
		return nil, err
	}
	return project, nil
}

func (r *ProjectRepository) DeleteProject(id uint) error {
	var project models.Project
	err := r.db.First(&project, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("project not found")
		}
		return err
	}

	err = r.db.Delete(&project).Error
	if err != nil {
		return err
	}
	return nil
}

// ContactRepository handles contact data operations
type ContactRepository struct {
	db *gorm.DB
}

func NewContactRepository(db *gorm.DB) *ContactRepository {
	return &ContactRepository{db: db}
}

func (r *ContactRepository) CreateContact(contact *models.Contact) (*models.Contact, error) {
	err := r.db.Create(contact).Error
	if err != nil {
		return nil, err
	}
	return contact, nil
}

func (r *ContactRepository) GetContacts() ([]models.Contact, error) {
	var contacts []models.Contact
	err := r.db.Order("created_at DESC").Find(&contacts).Error
	if err != nil {
		return nil, err
	}
	return contacts, nil
}

func (r *ContactRepository) UpdateContactStatus(id uint, status string) (*models.Contact, error) {
	var contact models.Contact
	err := r.db.First(&contact, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("contact not found")
		}
		return nil, err
	}

	contact.Status = status
	err = r.db.Save(&contact).Error
	if err != nil {
		return nil, err
	}
	return &contact, nil
}
