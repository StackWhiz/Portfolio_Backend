package api

import (
	"net/http"
	"stackwhiz-portfolio-backend/internal/service"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Handlers struct {
	profileService    *service.ProfileService
	experienceService *service.ExperienceService
	skillService      *service.SkillService
	projectService    *service.ProjectService
	contactService    *service.ContactService
	authService       *service.AuthService
}

func NewHandlers(
	profileService *service.ProfileService,
	experienceService *service.ExperienceService,
	skillService *service.SkillService,
	projectService *service.ProjectService,
	contactService *service.ContactService,
	authService *service.AuthService,
) *Handlers {
	return &Handlers{
		profileService:    profileService,
		experienceService: experienceService,
		skillService:      skillService,
		projectService:    projectService,
		contactService:    contactService,
		authService:       authService,
	}
}

// HealthCheck returns the health status of the API
// @Summary Health check endpoint
// @Description Returns the health status of the API
// @Tags health
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /health [get]
func (h *Handlers) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "healthy",
		"service": "stackwhiz-portfolio-backend",
		"version": "1.0.0",
	})
}

// GetProfile returns the main profile information
// @Summary Get profile information
// @Description Returns the main profile information
// @Tags profile
// @Accept json
// @Produce json
// @Success 200 {object} models.Profile
// @Router /profile [get]
func (h *Handlers) GetProfile(c *gin.Context) {
	profile, err := h.profileService.GetProfile()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get profile"})
		return
	}
	c.JSON(http.StatusOK, profile)
}

// UpdateProfile updates the main profile information
// @Summary Update profile information
// @Description Updates the main profile information (admin only)
// @Tags profile
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param profile body models.Profile true "Profile data"
// @Success 200 {object} models.Profile
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /admin/profile [put]
func (h *Handlers) UpdateProfile(c *gin.Context) {
	var profile service.ProfileUpdateRequest
	if err := c.ShouldBindJSON(&profile); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updatedProfile, err := h.profileService.UpdateProfile(&profile)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update profile"})
		return
	}

	c.JSON(http.StatusOK, updatedProfile)
}

// GetExperiences returns all work experiences
// @Summary Get work experiences
// @Description Returns all work experiences ordered by start date
// @Tags experiences
// @Accept json
// @Produce json
// @Success 200 {array} models.Experience
// @Router /experiences [get]
func (h *Handlers) GetExperiences(c *gin.Context) {
	experiences, err := h.experienceService.GetExperiences()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get experiences"})
		return
	}
	c.JSON(http.StatusOK, experiences)
}

// CreateExperience creates a new work experience
// @Summary Create work experience
// @Description Creates a new work experience entry (admin only)
// @Tags experiences
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param experience body service.ExperienceCreateRequest true "Experience data"
// @Success 201 {object} models.Experience
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /admin/experiences [post]
func (h *Handlers) CreateExperience(c *gin.Context) {
	var req service.ExperienceCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	experience, err := h.experienceService.CreateExperience(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create experience"})
		return
	}

	c.JSON(http.StatusCreated, experience)
}

// UpdateExperience updates an existing work experience
// @Summary Update work experience
// @Description Updates an existing work experience entry (admin only)
// @Tags experiences
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Experience ID"
// @Param experience body service.ExperienceUpdateRequest true "Experience data"
// @Success 200 {object} models.Experience
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /admin/experiences/{id} [put]
func (h *Handlers) UpdateExperience(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid experience ID"})
		return
	}

	var req service.ExperienceUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	experience, err := h.experienceService.UpdateExperience(uint(id), &req)
	if err != nil {
		if err.Error() == "experience not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "Experience not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update experience"})
		return
	}

	c.JSON(http.StatusOK, experience)
}

// DeleteExperience deletes a work experience
// @Summary Delete work experience
// @Description Deletes a work experience entry (admin only)
// @Tags experiences
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Experience ID"
// @Success 204
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /admin/experiences/{id} [delete]
func (h *Handlers) DeleteExperience(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid experience ID"})
		return
	}

	err = h.experienceService.DeleteExperience(uint(id))
	if err != nil {
		if err.Error() == "experience not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "Experience not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete experience"})
		return
	}

	c.Status(http.StatusNoContent)
}

// GetSkills returns all skills
// @Summary Get skills
// @Description Returns all skills grouped by category
// @Tags skills
// @Accept json
// @Produce json
// @Success 200 {array} models.Skill
// @Router /skills [get]
func (h *Handlers) GetSkills(c *gin.Context) {
	skills, err := h.skillService.GetSkills()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get skills"})
		return
	}
	c.JSON(http.StatusOK, skills)
}

// CreateSkill creates a new skill
// @Summary Create skill
// @Description Creates a new skill entry (admin only)
// @Tags skills
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param skill body service.SkillCreateRequest true "Skill data"
// @Success 201 {object} models.Skill
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /admin/skills [post]
func (h *Handlers) CreateSkill(c *gin.Context) {
	var req service.SkillCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	skill, err := h.skillService.CreateSkill(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create skill"})
		return
	}

	c.JSON(http.StatusCreated, skill)
}

// UpdateSkill updates an existing skill
// @Summary Update skill
// @Description Updates an existing skill entry (admin only)
// @Tags skills
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Skill ID"
// @Param skill body service.SkillUpdateRequest true "Skill data"
// @Success 200 {object} models.Skill
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /admin/skills/{id} [put]
func (h *Handlers) UpdateSkill(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid skill ID"})
		return
	}

	var req service.SkillUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	skill, err := h.skillService.UpdateSkill(uint(id), &req)
	if err != nil {
		if err.Error() == "skill not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "Skill not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update skill"})
		return
	}

	c.JSON(http.StatusOK, skill)
}

// DeleteSkill deletes a skill
// @Summary Delete skill
// @Description Deletes a skill entry (admin only)
// @Tags skills
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Skill ID"
// @Success 204
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /admin/skills/{id} [delete]
func (h *Handlers) DeleteSkill(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid skill ID"})
		return
	}

	err = h.skillService.DeleteSkill(uint(id))
	if err != nil {
		if err.Error() == "skill not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "Skill not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete skill"})
		return
	}

	c.Status(http.StatusNoContent)
}

// GetProjects returns all projects
// @Summary Get projects
// @Description Returns all projects, optionally filtered by featured status
// @Tags projects
// @Accept json
// @Produce json
// @Param featured query bool false "Filter by featured status"
// @Success 200 {array} models.Project
// @Router /projects [get]
func (h *Handlers) GetProjects(c *gin.Context) {
	featured := c.Query("featured")
	var featuredFilter *bool
	if featured != "" {
		if featured == "true" {
			featuredFilter = &[]bool{true}[0]
		} else if featured == "false" {
			featuredFilter = &[]bool{false}[0]
		}
	}

	projects, err := h.projectService.GetProjects(featuredFilter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get projects"})
		return
	}
	c.JSON(http.StatusOK, projects)
}

// CreateProject creates a new project
// @Summary Create project
// @Description Creates a new project entry (admin only)
// @Tags projects
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param project body service.ProjectCreateRequest true "Project data"
// @Success 201 {object} models.Project
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /admin/projects [post]
func (h *Handlers) CreateProject(c *gin.Context) {
	var req service.ProjectCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	project, err := h.projectService.CreateProject(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create project"})
		return
	}

	c.JSON(http.StatusCreated, project)
}

// UpdateProject updates an existing project
// @Summary Update project
// @Description Updates an existing project entry (admin only)
// @Tags projects
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Project ID"
// @Param project body service.ProjectUpdateRequest true "Project data"
// @Success 200 {object} models.Project
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /admin/projects/{id} [put]
func (h *Handlers) UpdateProject(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID"})
		return
	}

	var req service.ProjectUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	project, err := h.projectService.UpdateProject(uint(id), &req)
	if err != nil {
		if err.Error() == "project not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "Project not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update project"})
		return
	}

	c.JSON(http.StatusOK, project)
}

// DeleteProject deletes a project
// @Summary Delete project
// @Description Deletes a project entry (admin only)
// @Tags projects
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Project ID"
// @Success 204
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /admin/projects/{id} [delete]
func (h *Handlers) DeleteProject(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID"})
		return
	}

	err = h.projectService.DeleteProject(uint(id))
	if err != nil {
		if err.Error() == "project not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "Project not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete project"})
		return
	}

	c.Status(http.StatusNoContent)
}

// CreateContact creates a new contact form submission
// @Summary Create contact submission
// @Description Creates a new contact form submission
// @Tags contact
// @Accept json
// @Produce json
// @Param contact body service.ContactCreateRequest true "Contact data"
// @Success 201 {object} models.Contact
// @Failure 400 {object} map[string]interface{}
// @Router /contact [post]
func (h *Handlers) CreateContact(c *gin.Context) {
	var req service.ContactCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Add client information
	req.IPAddress = c.ClientIP()
	req.UserAgent = c.GetHeader("User-Agent")

	contact, err := h.contactService.CreateContact(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create contact"})
		return
	}

	c.JSON(http.StatusCreated, contact)
}

// GetContacts returns all contact submissions (admin only)
// @Summary Get contact submissions
// @Description Returns all contact form submissions (admin only)
// @Tags contact
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {array} models.Contact
// @Failure 401 {object} map[string]interface{}
// @Router /admin/contacts [get]
func (h *Handlers) GetContacts(c *gin.Context) {
	contacts, err := h.contactService.GetContacts()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get contacts"})
		return
	}
	c.JSON(http.StatusOK, contacts)
}

// UpdateContactStatus updates the status of a contact submission
// @Summary Update contact status
// @Description Updates the status of a contact form submission (admin only)
// @Tags contact
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Contact ID"
// @Param status body service.ContactStatusUpdateRequest true "Status data"
// @Success 200 {object} models.Contact
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /admin/contacts/{id}/status [put]
func (h *Handlers) UpdateContactStatus(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid contact ID"})
		return
	}

	var req service.ContactStatusUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	contact, err := h.contactService.UpdateContactStatus(uint(id), req.Status)
	if err != nil {
		if err.Error() == "contact not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "Contact not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update contact status"})
		return
	}

	c.JSON(http.StatusOK, contact)
}

// Login authenticates a user and returns a JWT token
// @Summary User login
// @Description Authenticates a user and returns a JWT token
// @Tags auth
// @Accept json
// @Produce json
// @Param credentials body service.LoginRequest true "Login credentials"
// @Success 200 {object} service.LoginResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /auth/login [post]
func (h *Handlers) Login(c *gin.Context) {
	var req service.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := h.authService.Login(&req)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	c.JSON(http.StatusOK, response)
}
