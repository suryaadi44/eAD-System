package repository

import (
	"context"
	"errors"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/suite"
	"github.com/suryaadi44/eAD-System/pkg/entity"
	"github.com/suryaadi44/eAD-System/pkg/utils"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type TestSuiteDocumentRepository struct {
	suite.Suite
	mock               sqlmock.Sqlmock
	documentRepository *DocumentRepositoryImpl
}

func (s *TestSuiteDocumentRepository) SetupTest() {
	dbMock, mock, err := sqlmock.New()
	s.NoError(err)
	s.mock = mock

	DB, _ := gorm.Open(mysql.New(mysql.Config{
		Conn:                      dbMock,
		SkipInitializeWithVersion: true,
	}), &gorm.Config{
		SkipDefaultTransaction: true,
	})

	s.documentRepository = &DocumentRepositoryImpl{db: DB}
}

func (s *TestSuiteDocumentRepository) TearDownTest() {
	s.mock = nil
	s.documentRepository = nil
}

func (s *TestSuiteDocumentRepository) TestAddTemplate() {
	query := regexp.QuoteMeta("INSERT INTO `templates` (`created_at`,`updated_at`,`deleted_at`,`name`,`path`,`margin_top`,`margin_bottom`,`margin_left`,`margin_right`,`is_active`) VALUES (?,?,?,?,?,?,?,?,?,?)")
	for _, tc := range []struct {
		Name        string
		Err         error
		ExpectedErr error
	}{
		{
			Name:        "Success",
			Err:         nil,
			ExpectedErr: nil,
		},
		{
			Name:        "Error duplicate file name",
			Err:         errors.New("Error 1062: Duplicate entry '' for key 'name'"),
			ExpectedErr: utils.ErrDuplicateTemplateName,
		},
		{
			Name:        "Error generic error",
			Err:         errors.New("generic error"),
			ExpectedErr: errors.New("generic error"),
		},
	} {
		s.SetupTest()
		s.Run(tc.Name, func() {
			if tc.Err != nil {
				s.mock.ExpectExec(query).WillReturnError(tc.Err)
			} else {
				s.mock.ExpectExec(query).WillReturnResult(sqlmock.NewResult(1, 1))
			}

			err := s.documentRepository.AddTemplate(context.Background(), &entity.Template{})

			s.Equal(tc.ExpectedErr, err)
		})
		s.TearDownTest()
	}
}

func (s *TestSuiteDocumentRepository) TestGetAllTemplate() {
	query := regexp.QuoteMeta("SELECT * FROM `templates` WHERE `templates`.`deleted_at` IS NULL")
	preloadField := regexp.QuoteMeta("SELECT * FROM `template_fields` WHERE `template_fields`.`template_id` = ? AND `template_fields`.`deleted_at` IS NULL")
	for _, tc := range []struct {
		Name             string
		Err              error
		ExpectedErr      error
		ExpectedReturn   *entity.Templates
		ReturnedRow      *sqlmock.Rows
		ReturnedRowField *sqlmock.Rows
	}{
		{
			Name:        "Success",
			Err:         nil,
			ExpectedErr: nil,
			ExpectedReturn: &entity.Templates{
				{
					Model: gorm.Model{
						ID: 1,
					},
					Name:         "template1",
					Path:         "path1",
					MarginTop:    1,
					MarginBottom: 1,
					MarginLeft:   1,
					MarginRight:  1,
					IsActive:     true,
					Fields: entity.TemplateFields{
						{
							Model: gorm.Model{
								ID: 1,
							},
							TemplateID: 1,
							Key:        "key1",
						},
					},
				},
			},
			ReturnedRow: sqlmock.NewRows([]string{"id", "name", "path", "margin_top", "margin_bottom", "margin_left", "margin_right", "is_active"}).
				AddRow(1, "template1", "path1", 1, 1, 1, 1, 1),
			ReturnedRowField: sqlmock.NewRows([]string{"id", "template_id", "key"}).
				AddRow(1, 1, "key1"),
		},
		{
			Name:             "Error No rows in result set",
			Err:              nil,
			ExpectedErr:      utils.ErrTemplateNotFound,
			ReturnedRow:      sqlmock.NewRows([]string{"id", "name", "path", "margin_top", "margin_bottom", "margin_left", "margin_right", "is_active"}),
			ReturnedRowField: sqlmock.NewRows([]string{"id", "template_id", "key"}),
		},
		{
			Name:        "Error generic error",
			Err:         errors.New("generic error"),
			ExpectedErr: errors.New("generic error"),
		},
	} {
		s.SetupTest()
		s.Run(tc.Name, func() {
			if tc.Err != nil {
				s.mock.ExpectQuery(query).WillReturnError(tc.Err)
			} else {
				s.mock.ExpectQuery(query).WillReturnRows(tc.ReturnedRow)
				s.mock.ExpectQuery(preloadField).WillReturnRows(tc.ReturnedRowField)
			}

			result, err := s.documentRepository.GetAllTemplate(context.Background())

			if tc.ExpectedErr != nil {
				s.Equal(tc.ExpectedErr, err)
			} else {
				s.Equal(tc.ExpectedReturn, result)
			}
		})
		s.TearDownTest()
	}
}

func (s *TestSuiteDocumentRepository) TestGetTemplateDetail() {
	query := regexp.QuoteMeta("SELECT * FROM `templates` WHERE id = ? AND `templates`.`deleted_at` IS NULL ORDER BY `templates`.`id` LIMIT 1")
	preloadField := regexp.QuoteMeta("SELECT * FROM `template_fields` WHERE `template_fields`.`template_id` = ? AND `template_fields`.`deleted_at` IS NULL")
	for _, tc := range []struct {
		Name             string
		Err              error
		ExpectedErr      error
		ExpectedReturn   *entity.Template
		ReturnedRow      *sqlmock.Rows
		ReturnedRowField *sqlmock.Rows
	}{
		{
			Name:        "Success",
			Err:         nil,
			ExpectedErr: nil,
			ExpectedReturn: &entity.Template{
				Model: gorm.Model{
					ID: 1,
				},
				Name:         "template1",
				Path:         "path1",
				MarginTop:    1,
				MarginBottom: 1,
				MarginLeft:   1,
				MarginRight:  1,
				IsActive:     true,
				Fields: entity.TemplateFields{
					{
						Model: gorm.Model{
							ID: 1,
						},
						TemplateID: 1,
						Key:        "key1",
					},
				},
			},
			ReturnedRow: sqlmock.NewRows([]string{"id", "name", "path", "margin_top", "margin_bottom", "margin_left", "margin_right", "is_active"}).
				AddRow(1, "template1", "path1", 1, 1, 1, 1, 1),
			ReturnedRowField: sqlmock.NewRows([]string{"id", "template_id", "key"}).
				AddRow(1, 1, "key1"),
		},
		{
			Name:             "Error No rows in result set",
			Err:              nil,
			ExpectedErr:      utils.ErrTemplateNotFound,
			ReturnedRow:      sqlmock.NewRows([]string{"id", "name", "path", "margin_top", "margin_bottom", "margin_left", "margin_right", "is_active"}),
			ReturnedRowField: sqlmock.NewRows([]string{"id", "template_id", "key"}),
		},
		{
			Name:        "Error generic error",
			Err:         errors.New("generic error"),
			ExpectedErr: errors.New("generic error"),
		},
	} {
		s.SetupTest()
		s.Run(tc.Name, func() {
			if tc.Err != nil {
				s.mock.ExpectQuery(query).WillReturnError(tc.Err)
			} else {
				s.mock.ExpectQuery(query).WillReturnRows(tc.ReturnedRow)
				s.mock.ExpectQuery(preloadField).WillReturnRows(tc.ReturnedRowField)
			}

			result, err := s.documentRepository.GetTemplateDetail(context.Background(), 1)

			if tc.ExpectedErr != nil {
				s.Equal(tc.ExpectedErr, err)
			} else {
				s.Equal(tc.ExpectedReturn, result)
			}
		})
		s.TearDownTest()
	}
}

func (s *TestSuiteDocumentRepository) TestGetTemplateFields() {
	query := regexp.QuoteMeta("SELECT * FROM `template_fields` WHERE template_id = ? AND `template_fields`.`deleted_at` IS NULL")
	for _, tc := range []struct {
		Name           string
		Err            error
		ExpectedErr    error
		ExpectedReturn *entity.TemplateFields
		ReturnedRow    *sqlmock.Rows
	}{
		{
			Name:        "Success",
			Err:         nil,
			ExpectedErr: nil,
			ExpectedReturn: &entity.TemplateFields{
				{
					Model: gorm.Model{
						ID: 1,
					},
					TemplateID: 1,
					Key:        "key1",
				},
			},

			ReturnedRow: sqlmock.NewRows([]string{"id", "template_id", "key"}).
				AddRow(1, 1, "key1"),
		},
		{
			Name:        "Error No rows in result set",
			Err:         nil,
			ExpectedErr: utils.ErrTemplateFieldNotFound,
			ReturnedRow: sqlmock.NewRows([]string{"id", "template_id", "key"}),
		},
		{
			Name:        "Error generic error",
			Err:         errors.New("generic error"),
			ExpectedErr: errors.New("generic error"),
		},
	} {
		s.SetupTest()
		s.Run(tc.Name, func() {
			if tc.Err != nil {
				s.mock.ExpectQuery(query).WillReturnError(tc.Err)
			} else {
				s.mock.ExpectQuery(query).WillReturnRows(tc.ReturnedRow)
			}

			result, err := s.documentRepository.GetTemplateFields(context.Background(), 1)

			if tc.ExpectedErr != nil {
				s.Equal(tc.ExpectedErr, err)
			} else {
				s.Equal(tc.ExpectedReturn, result)
			}
		})
		s.TearDownTest()
	}
}

func (s *TestSuiteDocumentRepository) TestAddDocument() {
	query := regexp.QuoteMeta("INSERT INTO `documents` (`id`,`register`,`description`,`applicant_id`,`template_id`,`stage_id`,`created_at`,`updated_at`,`deleted_at`) VALUES (?,?,?,?,?,?,?,?,?)")
	for _, tc := range []struct {
		Name        string
		Err         error
		ExpectedErr error
	}{
		{
			Name:        "Success",
			Err:         nil,
			ExpectedErr: nil,
		},
		{
			Name:        "Error duplicate register",
			Err:         errors.New("Error 1062: Duplicate entry '' for key ''"),
			ExpectedErr: utils.ErrDuplicateRegister,
		},
		{
			Name:        "Error generic error",
			Err:         errors.New("generic error"),
			ExpectedErr: errors.New("generic error"),
		},
	} {
		s.SetupTest()
		s.Run(tc.Name, func() {
			if tc.Err != nil {
				s.mock.ExpectExec(query).WillReturnError(tc.Err)
			} else {
				s.mock.ExpectExec(query).WillReturnResult(sqlmock.NewResult(1, 1))
			}

			_, err := s.documentRepository.AddDocument(context.Background(), &entity.Document{})

			s.Equal(tc.ExpectedErr, err)
		})
		s.TearDownTest()
	}
}

func (s *TestSuiteDocumentRepository) TestGetDocument() {
	query := regexp.QuoteMeta("SELECT * FROM `documents` WHERE id = ? AND `documents`.`deleted_at` IS NULL ORDER BY `documents`.`id` LIMIT 1")
	queryPreloadUser := regexp.QuoteMeta("SELECT id, username, name FROM `users` WHERE `users`.`id` = ? AND `users`.`deleted_at` IS NULL")
	queryPreloadStage := regexp.QuoteMeta("SELECT * FROM `stages` WHERE `stages`.`id` = ?")
	queryPreloadTemplate := regexp.QuoteMeta("SELECT * FROM `templates` WHERE `templates`.`id` = ? AND `templates`.`deleted_at` IS NULL")
	queryPreloadTemplateFields := regexp.QuoteMeta("SELECT * FROM `template_fields` WHERE `template_fields`.`id` = ? AND `template_fields`.`deleted_at` IS NULL")
	queryPreloadDocumentFields := regexp.QuoteMeta("SELECT * FROM `document_fields` WHERE `document_fields`.`document_id` = ? AND `document_fields`.`deleted_at` IS NULL")
	queryPreloadEmployee := regexp.QuoteMeta("SELECT id, username, name, n_ip, position FROM `users` WHERE `users`.`id` = ? AND `users`.`deleted_at` IS NULL")

	for _, tc := range []struct {
		Name           string
		Err            error
		ExpectedErr    error
		ExpectedReturn *entity.Document
		ReturnedRow    *sqlmock.Rows
	}{
		{
			Name:        "Success",
			Err:         nil,
			ExpectedErr: nil,
			ExpectedReturn: &entity.Document{
				ID:          "1",
				Register:    "register",
				Description: "description",
				ApplicantID: "1",
				Applicant: entity.User{
					ID:       "1",
					Username: "username",
					Name:     "name",
				},
				TemplateID: 1,
				Template: entity.Template{
					Model: gorm.Model{
						ID: 1,
					},
					Name:         "template",
					Path:         "path",
					MarginTop:    0,
					MarginBottom: 0,
					MarginLeft:   0,
					MarginRight:  0,
					IsActive:     false,
					Fields:       nil,
				},
				Fields: entity.DocumentFields{
					{
						Model: gorm.Model{
							ID: 1,
						},
						DocumentID:      "1",
						TemplateFieldID: 1,
						TemplateField: entity.TemplateField{
							Model: gorm.Model{
								ID: 1,
							},
							TemplateID: 1,
							Key:        "key",
						},
						Value: "value",
					},
				},
				StageID: 3,
				Stage: entity.Stage{
					ID:     3,
					Status: "approved",
				},
				VerifierID: "1",
				Verifier: entity.User{
					ID:       "1",
					Username: "username",
					Name:     "name",
					NIP:      "123",
					Position: "position",
				},
				VerifiedAt: time.Time{},
				SignerID:   "1",
				Signer: entity.User{
					ID:       "1",
					Username: "username",
					Name:     "name",
					NIP:      "123",
					Position: "position",
				},
				SignedAt:  time.Time{},
				CreatedAt: time.Time{},
				UpdatedAt: time.Time{},
				DeletedAt: gorm.DeletedAt{},
			},
		},
		{
			Name:        "Error No rows in result set",
			Err:         gorm.ErrRecordNotFound,
			ExpectedErr: utils.ErrDocumentNotFound,
		},
		{
			Name:        "Error generic error",
			Err:         errors.New("generic error"),
			ExpectedErr: errors.New("generic error"),
		},
	} {
		s.SetupTest()
		s.Run(tc.Name, func() {
			if tc.Err != nil {
				s.mock.ExpectQuery(query).WillReturnError(tc.Err)
			} else {
				s.mock.ExpectQuery(query).WillReturnRows(sqlmock.NewRows([]string{"id", "register", "description", "applicant_id", "template_id", "stage_id", "verifier_id", "signer_id"}).
					AddRow(1, "register", "description", "1", 1, 3, "1", "1"))
				s.mock.ExpectQuery(queryPreloadUser).WillReturnRows(sqlmock.NewRows([]string{"id", "username", "name"}).AddRow(1, "username", "name"))
				s.mock.ExpectQuery(queryPreloadDocumentFields).WillReturnRows(sqlmock.NewRows([]string{"id", "document_id", "template_field_id", "value"}).AddRow(1, 1, 1, "value"))
				s.mock.ExpectQuery(queryPreloadTemplateFields).WillReturnRows(sqlmock.NewRows([]string{"id", "template_id", "key"}).AddRow(1, 1, "key"))
				s.mock.ExpectQuery(queryPreloadEmployee).WillReturnRows(sqlmock.NewRows([]string{"id", "username", "name", "n_ip", "position"}).AddRow(1, "username", "name", "123", "position"))
				s.mock.ExpectQuery(queryPreloadStage).WillReturnRows(sqlmock.NewRows([]string{"id", "status"}).AddRow(3, "approved"))
				s.mock.ExpectQuery(queryPreloadTemplate).WillReturnRows(sqlmock.NewRows([]string{"id", "name", "path", "margin_top", "margin_bottom", "margin_left", "margin_right", "is_active"}).AddRow(1, "template", "path", 0, 0, 0, 0, false))
				s.mock.ExpectQuery(queryPreloadEmployee).WillReturnRows(sqlmock.NewRows([]string{"id", "username", "name", "n_ip", "position"}).AddRow(1, "username", "name", "123", "position"))

			}

			result, err := s.documentRepository.GetDocument(context.Background(), "")

			if tc.ExpectedErr != nil {
				s.Equal(tc.ExpectedErr, err)
			} else {
				s.Equal(tc.ExpectedReturn, result)
			}
		})
		s.TearDownTest()
	}
}

func (s *TestSuiteDocumentRepository) TestGetBriefDocuments() {
	query := regexp.QuoteMeta("SELECT id, register, description, created_at, applicant_id, template_id, stage_id FROM `documents` WHERE `documents`.`deleted_at` IS NULL ORDER BY created_at desc")
	queryPreloadAplicant := regexp.QuoteMeta("SELECT id, username, name FROM `users` WHERE `users`.`id` = ? AND `users`.`deleted_at` IS NULL")
	queryPreloadStage := regexp.QuoteMeta("SELECT id, status FROM `stages` WHERE `stages`.`id` = ? AND `stages`.`deleted_at` IS NULL")
	queryPreloadtemplate := regexp.QuoteMeta("SELECT name FROM `templates` WHERE `templates`.`id` = ? AND `templates`.`deleted_at` IS NULL")

	for _, tc := range []struct {
		Name           string
		Err            error
		ExpectedErr    error
		ExpectedReturn *entity.Documents
		ReturnedRows   *sqlmock.Rows
	}{
		{
			Name:        "Success",
			Err:         nil,
			ExpectedErr: nil,
			ExpectedReturn: &entity.Documents{
				{
					ID:          "1",
					Register:    "register",
					Description: "description",
					CreatedAt:   time.Time{},
				},
			},
			ReturnedRows: sqlmock.NewRows([]string{"id", "register", "description", "created_at"}).AddRow(1, "register", "description", time.Time{}),
		},
		{
			Name:           "Error No rows in result set",
			Err:            nil,
			ExpectedErr:    utils.ErrDocumentNotFound,
			ExpectedReturn: &entity.Documents{},
			ReturnedRows:   sqlmock.NewRows([]string{"id", "register", "description", "created_at"}),
		},
		{
			Name:           "Error generic error",
			Err:            errors.New("generic error"),
			ExpectedErr:    errors.New("generic error"),
			ExpectedReturn: nil,
			ReturnedRows:   nil,
		},
	} {
		s.SetupTest()
		s.Run(tc.Name, func() {
			if tc.Err != nil {
				s.mock.ExpectQuery(query).WillReturnError(tc.Err)
			} else {
				s.mock.ExpectQuery(query).WillReturnRows(tc.ReturnedRows)
				s.mock.ExpectQuery(queryPreloadAplicant).WillReturnRows(sqlmock.NewRows([]string{"id", "username", "name"}).AddRow(1, "username", "name"))
				s.mock.ExpectQuery(queryPreloadStage).WillReturnRows(sqlmock.NewRows([]string{"id", "status"}).AddRow(3, "approved"))
				s.mock.ExpectQuery(queryPreloadtemplate).WillReturnRows(sqlmock.NewRows([]string{"name"}).AddRow("template"))
			}

			result, err := s.documentRepository.GetBriefDocuments(context.Background(), 0, 0)

			if tc.ExpectedErr != nil {
				s.Equal(tc.ExpectedErr, err)
			} else {
				s.Equal(tc.ExpectedReturn, result)
			}
		})
		s.TearDownTest()
	}
}

func (s *TestSuiteDocumentRepository) TestGetBriefDocumentsByApplicant() {
	query := regexp.QuoteMeta("SELECT id, register, description, created_at, applicant_id, template_id, stage_id FROM `documents` WHERE applicant_id = ? AND `documents`.`deleted_at` IS NULL ORDER BY created_at desc")
	queryPreloadAplicant := regexp.QuoteMeta("SELECT id, username, name FROM `users` WHERE `users`.`id` = ? AND `users`.`deleted_at` IS NULL")
	queryPreloadStage := regexp.QuoteMeta("SELECT id, status FROM `stages` WHERE `stages`.`id` = ? AND `stages`.`deleted_at` IS NULL")
	queryPreloadtemplate := regexp.QuoteMeta("SELECT name FROM `templates` WHERE `templates`.`id` = ? AND `templates`.`deleted_at` IS NULL")

	for _, tc := range []struct {
		Name           string
		Err            error
		ExpectedErr    error
		ExpectedReturn *entity.Documents
		ReturnedRows   *sqlmock.Rows
	}{
		{
			Name:        "Success",
			Err:         nil,
			ExpectedErr: nil,
			ExpectedReturn: &entity.Documents{
				{
					ID:          "1",
					Register:    "register",
					Description: "description",
					CreatedAt:   time.Time{},
				},
			},
			ReturnedRows: sqlmock.NewRows([]string{"id", "register", "description", "created_at"}).AddRow(1, "register", "description", time.Time{}),
		},
		{
			Name:           "Error No rows in result set",
			Err:            nil,
			ExpectedErr:    utils.ErrDocumentNotFound,
			ExpectedReturn: &entity.Documents{},
			ReturnedRows:   sqlmock.NewRows([]string{"id", "register", "description", "created_at"}),
		},
		{
			Name:           "Error generic error",
			Err:            errors.New("generic error"),
			ExpectedErr:    errors.New("generic error"),
			ExpectedReturn: nil,
			ReturnedRows:   nil,
		},
	} {
		s.SetupTest()
		s.Run(tc.Name, func() {
			if tc.Err != nil {
				s.mock.ExpectQuery(query).WillReturnError(tc.Err)
			} else {
				s.mock.ExpectQuery(query).WillReturnRows(tc.ReturnedRows)
				s.mock.ExpectQuery(queryPreloadAplicant).WillReturnRows(sqlmock.NewRows([]string{"id", "username", "name"}).AddRow(1, "username", "name"))
				s.mock.ExpectQuery(queryPreloadStage).WillReturnRows(sqlmock.NewRows([]string{"id", "status"}).AddRow(3, "approved"))
				s.mock.ExpectQuery(queryPreloadtemplate).WillReturnRows(sqlmock.NewRows([]string{"name"}).AddRow("template"))
			}

			result, err := s.documentRepository.GetBriefDocumentsByApplicant(context.Background(), "1", 0, 0)

			if tc.ExpectedErr != nil {
				s.Equal(tc.ExpectedErr, err)
			} else {
				s.Equal(tc.ExpectedReturn, result)
			}
		})
		s.TearDownTest()
	}
}

func (s *TestSuiteDocumentRepository) TestGetDocumentStatus() {
	query := regexp.QuoteMeta("SELECT * FROM `documents` WHERE id = ? AND `documents`.`deleted_at` IS NULL ORDER BY `documents`.`id` LIMIT 1")
	queryPreloadStage := regexp.QuoteMeta("SELECT * FROM `stages` WHERE `stages`.`id` = ?")
	queryPreloadEmployee := regexp.QuoteMeta("SELECT id, username, name, n_ip, position FROM `users` WHERE `users`.`id` = ? AND `users`.`deleted_at` IS NULL")

	for _, tc := range []struct {
		Name           string
		Err            error
		ExpectedErr    error
		ExpectedReturn *entity.Document
		ReturnedRow    *sqlmock.Rows
	}{
		{
			Name: "Success",
			ExpectedReturn: &entity.Document{
				ID:          "1",
				Register:    "register",
				Description: "description",
				ApplicantID: "1",
				TemplateID:  1,
				StageID:     3,
				Stage: entity.Stage{
					ID:     3,
					Status: "approved",
				},
				VerifierID: "1",
				Verifier: entity.User{
					ID:       "1",
					Username: "username",
					Name:     "name",
					NIP:      "123",
					Position: "position",
				},
				VerifiedAt: time.Time{},
				SignerID:   "1",
				Signer: entity.User{
					ID:       "1",
					Username: "username",
					Name:     "name",
					NIP:      "123",
					Position: "position",
				},
				SignedAt:  time.Time{},
				CreatedAt: time.Time{},
				UpdatedAt: time.Time{},
				DeletedAt: gorm.DeletedAt{},
			},
		},
		{
			Name:        "Error No rows in result set",
			Err:         gorm.ErrRecordNotFound,
			ExpectedErr: utils.ErrDocumentNotFound,
		},
		{
			Name:        "Error generic error",
			Err:         errors.New("generic error"),
			ExpectedErr: errors.New("generic error"),
		},
	} {
		s.SetupTest()
		s.Run(tc.Name, func() {
			if tc.Err != nil {
				s.mock.ExpectQuery(query).WillReturnError(tc.Err)
			} else {
				s.mock.ExpectQuery(query).WillReturnRows(sqlmock.NewRows([]string{"id", "register", "description", "applicant_id", "template_id", "stage_id", "verifier_id", "signer_id"}).
					AddRow(1, "register", "description", "1", 1, 3, "1", "1"))
				s.mock.ExpectQuery(queryPreloadEmployee).WillReturnRows(sqlmock.NewRows([]string{"id", "username", "name", "n_ip", "position"}).AddRow(1, "username", "name", "123", "position"))
				s.mock.ExpectQuery(queryPreloadStage).WillReturnRows(sqlmock.NewRows([]string{"id", "status"}).AddRow(3, "approved"))
				s.mock.ExpectQuery(queryPreloadEmployee).WillReturnRows(sqlmock.NewRows([]string{"id", "username", "name", "n_ip", "position"}).AddRow(1, "username", "name", "123", "position"))
			}

			result, err := s.documentRepository.GetDocumentStatus(context.Background(), "")

			if tc.ExpectedErr != nil {
				s.Equal(tc.ExpectedErr, err)
			} else {
				s.Equal(tc.ExpectedReturn, result)
			}
		})
		s.TearDownTest()
	}
}

func (s *TestSuiteDocumentRepository) TestGetApplicantID() {
	query := regexp.QuoteMeta("SELECT `applicant_id` FROM `documents` WHERE id = ? AND `documents`.`deleted_at` IS NULL ORDER BY `documents`.`id` LIMIT 1")

	res := "1"

	for _, tc := range []struct {
		Name           string
		Err            error
		ExpectedErr    error
		ExpectedReturn *string
		ReturnedRow    *sqlmock.Rows
	}{
		{
			Name:           "Success",
			ExpectedReturn: &res,
		},
		{
			Name:        "Error No rows in result set",
			Err:         gorm.ErrRecordNotFound,
			ExpectedErr: utils.ErrDocumentNotFound,
		},
		{
			Name:        "Error generic error",
			Err:         errors.New("generic error"),
			ExpectedErr: errors.New("generic error"),
		},
	} {
		s.SetupTest()
		s.Run(tc.Name, func() {
			if tc.Err != nil {
				s.mock.ExpectQuery(query).WillReturnError(tc.Err)
			} else {
				s.mock.ExpectQuery(query).WillReturnRows(sqlmock.NewRows([]string{"applicant_id"}).AddRow("1"))
			}

			result, err := s.documentRepository.GetApplicantID(context.Background(), "")

			if tc.ExpectedErr != nil {
				s.Equal(tc.ExpectedErr, err)
			} else {
				s.Equal(tc.ExpectedReturn, result)
			}
		})
		s.TearDownTest()
	}
}

func (s *TestSuiteDocumentRepository) TestGetDocumentStage() {
	query := regexp.QuoteMeta("SELECT `stage_id` FROM `documents` WHERE id = ? AND `documents`.`deleted_at` IS NULL ORDER BY `documents`.`id` LIMIT 1")

	res := 1

	for _, tc := range []struct {
		Name           string
		Err            error
		ExpectedErr    error
		ExpectedReturn *int
		ReturnedRow    *sqlmock.Rows
	}{
		{
			Name:           "Success",
			ExpectedReturn: &res,
		},
		{
			Name:        "Error No rows in result set",
			Err:         gorm.ErrRecordNotFound,
			ExpectedErr: utils.ErrDocumentNotFound,
		},
		{
			Name:        "Error generic error",
			Err:         errors.New("generic error"),
			ExpectedErr: errors.New("generic error"),
		},
	} {
		s.SetupTest()
		s.Run(tc.Name, func() {
			if tc.Err != nil {
				s.mock.ExpectQuery(query).WillReturnError(tc.Err)
			} else {
				s.mock.ExpectQuery(query).WillReturnRows(sqlmock.NewRows([]string{"stage_id"}).AddRow(1))
			}

			result, err := s.documentRepository.GetDocumentStage(context.Background(), "")

			if tc.ExpectedErr != nil {
				s.Equal(tc.ExpectedErr, err)
			} else {
				s.Equal(tc.ExpectedReturn, result)
			}
		})
		s.TearDownTest()
	}
}

func (s *TestSuiteDocumentRepository) TestVerifyDocument() {
	query := regexp.QuoteMeta("UPDATE `documents` SET `stage_id`=?,`verifier_id`=?,`verified_at`=?,`updated_at`=? WHERE id = ? AND `documents`.`deleted_at` IS NULL")

	for _, tc := range []struct {
		Name         string
		Err          error
		ExpectedErr  error
		RowsAffected int64
	}{
		{
			Name:         "Success",
			RowsAffected: 1,
		},
		{
			Name:         "Error No rows affected",
			RowsAffected: 0,
			ExpectedErr:  utils.ErrDocumentNotFound,
		},
		{
			Name:        "Error generic error",
			Err:         errors.New("generic error"),
			ExpectedErr: errors.New("generic error"),
		},
	} {
		s.SetupTest()
		s.Run(tc.Name, func() {
			if tc.Err != nil {
				s.mock.ExpectExec(query).WillReturnError(tc.Err)
			} else {
				s.mock.ExpectExec(query).WillReturnResult(sqlmock.NewResult(1, tc.RowsAffected))
			}

			err := s.documentRepository.VerifyDocument(context.Background(), &entity.Document{})

			if tc.ExpectedErr != nil {
				s.Equal(tc.ExpectedErr, err)
			}
		})
		s.TearDownTest()
	}
}

func (s *TestSuiteDocumentRepository) TestSignDocument() {
	query := regexp.QuoteMeta("UPDATE `documents` SET `stage_id`=?,`signer_id`=?,`signed_at`=?,`updated_at`=? WHERE id = ? AND `documents`.`deleted_at` IS NULL")

	for _, tc := range []struct {
		Name         string
		Err          error
		ExpectedErr  error
		RowsAffected int64
	}{
		{
			Name:         "Success",
			RowsAffected: 1,
		},
		{
			Name:         "Error No rows affected",
			RowsAffected: 0,
			ExpectedErr:  utils.ErrDocumentNotFound,
		},
		{
			Name:        "Error generic error",
			Err:         errors.New("generic error"),
			ExpectedErr: errors.New("generic error"),
		},
	} {
		s.SetupTest()
		s.Run(tc.Name, func() {
			if tc.Err != nil {
				s.mock.ExpectExec(query).WillReturnError(tc.Err)
			} else {
				s.mock.ExpectExec(query).WillReturnResult(sqlmock.NewResult(1, tc.RowsAffected))
			}

			err := s.documentRepository.SignDocument(context.Background(), &entity.Document{})

			if tc.ExpectedErr != nil {
				s.Equal(tc.ExpectedErr, err)
			}
		})
		s.TearDownTest()
	}
}

func TestDocumentRepository(t *testing.T) {
	suite.Run(t, new(TestSuiteDocumentRepository))
}
