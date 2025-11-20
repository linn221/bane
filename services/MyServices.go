package services

import (
	"errors"
	"time"

	"github.com/linn221/bane/config"
	"github.com/linn221/bane/models"
	"github.com/linn221/bane/utils"
	"github.com/linn221/bane/validate"
	"gorm.io/gorm"
)

// MyServices contains all service instances
type MyServices struct {
	EndpointService       *endpointService
	TagService            *tagService
	NoteService           *noteService
	ProgramService        *programService
	MyRequestService      *myRequestService
	MemorySheetService    *memorySheetService
	WordService           *wordService
	VulnConnectionService *vulnConnectionService
	MySheetService        *mySheetService
	ProjectService        *projectService
	TaskService           *taskService
	AliasService          *aliasService
}

// NewMyServices creates a new MyServices instance with all services initialized
func NewMyServices(db *gorm.DB, cache config.CacheService) *MyServices {
	aliasService := &aliasService{db: db}

	endpointService := &endpointService{
		GeneralCrud: GeneralCrud[models.EndpointInput, models.Endpoint]{
			transform: func(input *models.EndpointInput) models.Endpoint {
				return models.Endpoint{
					Name:        input.Name,
					Description: input.Description,
					HttpSchema:  input.HttpSchema,
					HttpMethod:  input.HttpMethod,
					HttpDomain:  input.HttpDomain,
					HttpPath:    input.HttpPath,
					HttpQueries: input.HttpQueries,
					HttpHeaders: input.HttpHeaders,
					HttpCookies: input.HttpCookies,
					HttpBody:    input.HttpBody,
				}
			},
		},
		db:           db,
		aliasService: aliasService,
	}

	tagService := &tagService{
		db:           db,
		aliasService: aliasService,
	}

	noteService := &noteService{
		GeneralCrud: GeneralCrud[models.NoteInput, models.Note]{
			transform: func(input *models.NoteInput) models.Note {
				today := utils.Today()
				return models.Note{
					ReferenceType: input.ReferenceType,
					ReferenceID:   input.ReferenceId,
					Value:         input.Value,
					NoteDate:      models.MyDate{Time: today},
				}
			},
		},
		db: db,
	}

	programService := &programService{
		GeneralCrud: GeneralCrud[models.ProgramInput, models.Program]{
			transform: func(input *models.ProgramInput) models.Program {
				return models.Program{
					Name:        input.Name,
					Url:         input.URL,
					Description: utils.SafeDeref(input.Description),
					Domain:      input.Domain,
				}
			},
			updates: func(existing models.Program, input *models.ProgramInput) map[string]any {
				return map[string]any{
					"Name":        input.Name,
					"Url":         input.URL,
					"Description": utils.SafeDeref(input.Description),
					"Domain":      input.Domain,
				}
			},
			validateWrite: func(db *gorm.DB, input *models.ProgramInput, id int) error {
				// Check alias uniqueness using AliasService
				if input.Alias != "" {
					existingId, err := aliasService.GetId(input.Alias)
					if err == nil && existingId != id {
						return errors.New("duplicate program alias")
					}
				}
				return validate.Validate(db,
					validate.NewUniqueRule("programs", "name", input.Name, nil).Except(id).Say("duplicate program name"))
			},
			validateDelete: func(db *gorm.DB, existing models.Program) error {
				return nil
			},
		},
		db:           db,
		aliasService: aliasService,
	}

	myRequestService := &myRequestService{
		db: db,
	}

		memorySheetService := &memorySheetService{
			GeneralCrud: GeneralCrud[models.MemorySheetInput, models.MemorySheet]{
				transform: func(input *models.MemorySheetInput) models.MemorySheet {
					result := models.MemorySheet{
						Value: input.Value,
					}
					today := utils.Today()
					result.CreateDate = models.MyDate{Time: today}
					result.CurrentDate = result.CreateDate
					result.NextDate = models.MyDate{Time: result.CurrentDate.Time.AddDate(0, 0, 1)}
					return result
				},
				updates: func(existing models.MemorySheet, input *models.MemorySheetInput) map[string]any {
					updates := map[string]any{}

					if input.UpdateNextDate {
						currentDate := existing.NextDate
						nextDate := GetNextDate(currentDate.Time, existing.Index+1)
						// NextDate has moved for the note
						updates["CurrentDate"] = currentDate
						updates["NextDate"] = models.MyDate{Time: nextDate}
						updates["Index"] = existing.Index + 1
					} else { // normal update coming from graphql
						if input.Value != "" {
							updates["Value"] = input.Value
						}
					}
					return updates
				},
			},
			db:           db,
			aliasService: aliasService,
		}

	wordService := &wordService{
		db:           db,
		aliasService: aliasService,
	}

	vulnConnectionService := &vulnConnectionService{
		db: db,
	}

		mySheetService := &mySheetService{
			GeneralCrud: GeneralCrud[models.MySheetInput, models.MySheet]{
				transform: func(input *models.MySheetInput) models.MySheet {
					result := models.MySheet{
						Title: input.Title,
						Body:  input.Body,
					}
					today := utils.Today()
					if input.Date != nil {
						today = input.Date.Time
					}
					result.Created = models.MyDate{Time: today}
					result.NextDate = models.MyDate{Time: today.AddDate(0, 0, 1)}
					result.PreviousDate = models.MyDate{Time: time.Time{}}
					return result
				},
				updates: func(existing models.MySheet, input *models.MySheetInput) map[string]any {
					updates := map[string]any{}

					if input.UpdateNextDate {
						currentDate := existing.NextDate
						nextDate := GetNextDate(currentDate.Time, existing.Index+1)
						// NextDate has moved for the sheet
						updates["PreviousDate"] = currentDate
						updates["NextDate"] = models.MyDate{Time: nextDate}
						updates["Index"] = existing.Index + 1
					} else { // normal update coming from graphql
						if input.Title != "" {
							updates["Title"] = input.Title
						}
						if input.Body != "" {
							updates["Body"] = input.Body
						}
					}
					return updates
				},
			},
			db:           db,
			aliasService: aliasService,
		}

	projectService := &projectService{
		GeneralCrud: GeneralCrud[models.ProjectInput, models.Project]{
			transform: func(input *models.ProjectInput) models.Project {
				return models.Project{
					Name:        input.Name,
					Description: input.Description,
				}
			},
			updates: func(existing models.Project, input *models.ProjectInput) map[string]any {
				updates := map[string]any{}
				if input.Name != "" {
					updates["Name"] = input.Name
				}
				if input.Description != "" {
					updates["Description"] = input.Description
				}
				return updates
			},
			validateWrite: func(db *gorm.DB, input *models.ProjectInput, id int) error {
				return input.Validate(db, id)
			},
		},
		db:           db,
		aliasService: aliasService,
	}

		taskService := &taskService{
			GeneralCrud: GeneralCrud[models.TaskInput, models.Task]{
				transform: func(input *models.TaskInput) models.Task {
					result := models.Task{
						Title:       input.Title,
						Description: input.Description,
						Priority:    input.Priority,
						Status:      models.TaskStatusInProgress,
						Created:     models.MyDate{Time: utils.Today()},
					}
					if input.Deadline != nil {
						result.Deadline = *input.Deadline
					}
					if input.RemindDate != nil {
						result.RemindDate = *input.RemindDate
					}
					return result
				},
				updates: func(existing models.Task, input *models.TaskInput) map[string]any {
					updates := map[string]any{}
					if input.Title != "" {
						updates["Title"] = input.Title
					}
					if input.Description != "" {
						updates["Description"] = input.Description
					}
					if input.Priority != 0 {
						updates["Priority"] = input.Priority
					}
					if input.Deadline != nil {
						updates["Deadline"] = *input.Deadline
					}
					if input.RemindDate != nil {
						updates["RemindDate"] = *input.RemindDate
					}
					return updates
				},
				validateWrite: func(db *gorm.DB, input *models.TaskInput, id int) error {
					return input.Validate(db, id)
				},
			},
			db:           db,
			aliasService: aliasService,
		}

	return &MyServices{
		AliasService:          aliasService,
		EndpointService:       endpointService,
		TagService:            tagService,
		NoteService:           noteService,
		ProgramService:        programService,
		MyRequestService:      myRequestService,
		MemorySheetService:    memorySheetService,
		WordService:           wordService,
		VulnConnectionService: vulnConnectionService,
		MySheetService:        mySheetService,
		ProjectService:        projectService,
		TaskService:           taskService,
	}
}
