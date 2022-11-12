package repository

import (
	"context"
	"errors"
	error2 "github.com/suryaadi44/eAD-System/pkg/utils/error"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/suite"
	"github.com/suryaadi44/eAD-System/pkg/entity"
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

func (s *TestSuiteDocumentRepository) TestAddDocument() {
	query := regexp.QuoteMeta("INSERT INTO `documents` (`id`,`description`,`applicant_id`,`template_id`,`stage_id`,`created_at`,`updated_at`,`deleted_at`) VALUES (?,?,?,?,?,?,?,?)")
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
			ExpectedErr: error2.ErrDuplicateRegister,
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
				RegisterID:  123,
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
			ExpectedErr: error2.ErrDocumentNotFound,
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
				s.mock.ExpectQuery(query).WillReturnRows(sqlmock.NewRows([]string{"id", "register_id", "description", "applicant_id", "template_id", "stage_id", "verifier_id", "signer_id"}).
					AddRow(1, 123, "description", "1", 1, 3, "1", "1"))
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

func (s *TestSuiteDocumentRepository) TestGetBriefDocument() {
	query := regexp.QuoteMeta("SELECT id, register_id, description, created_at, applicant_id, template_id, stage_id FROM `documents` WHERE id = ? AND `documents`.`deleted_at` IS NULL ORDER BY created_at desc,`documents`.`id` LIMIT 1")
	queryPreloadAplicant := regexp.QuoteMeta("SELECT id, username, name FROM `users` WHERE `users`.`id` = ? AND `users`.`deleted_at` IS NULL")
	queryPreloadStage := regexp.QuoteMeta("SELECT id, status FROM `stages` WHERE `stages`.`id` = ? AND `stages`.`deleted_at` IS NULL")
	queryPreloadtemplate := regexp.QuoteMeta("SELECT name FROM `templates` WHERE `templates`.`id` = ? AND `templates`.`deleted_at` IS NULL")
	queryPreloadRegister := regexp.QuoteMeta("SELECT * FROM `registers` WHERE `registers`.`id` = ? AND `registers`.`deleted_at` IS NULL")

	for _, tc := range []struct {
		Name           string
		Err            error
		ExpectedErr    error
		ExpectedReturn *entity.Document
		ReturnedRows   *sqlmock.Rows
	}{
		{
			Name:        "Success",
			Err:         nil,
			ExpectedErr: nil,
			ExpectedReturn: &entity.Document{
				ID:         "1",
				RegisterID: 123,
				Register: entity.Register{
					Model: gorm.Model{
						ID: 123,
					},
					Description: "description",
				},
				Description: "description",
				CreatedAt:   time.Time{},
			},
			ReturnedRows: sqlmock.NewRows([]string{"id", "register_id", "description", "created_at"}).AddRow(1, 123, "description", time.Time{}),
		},
		{
			Name:           "Error No rows in result set",
			Err:            gorm.ErrRecordNotFound,
			ExpectedErr:    error2.ErrDocumentNotFound,
			ExpectedReturn: &entity.Document{},
			ReturnedRows:   sqlmock.NewRows([]string{"id", "register_id", "description", "created_at"}),
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
				s.mock.ExpectQuery(queryPreloadRegister).WillReturnRows(sqlmock.NewRows([]string{"id", "description"}).AddRow(123, "description"))
				s.mock.ExpectQuery(queryPreloadAplicant).WillReturnRows(sqlmock.NewRows([]string{"id", "username", "name"}).AddRow(1, "username", "name"))
				s.mock.ExpectQuery(queryPreloadStage).WillReturnRows(sqlmock.NewRows([]string{"id", "status"}).AddRow(3, "approved"))
				s.mock.ExpectQuery(queryPreloadtemplate).WillReturnRows(sqlmock.NewRows([]string{"name"}).AddRow("template"))
			}

			result, err := s.documentRepository.GetBriefDocument(context.Background(), "1")

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
	query := regexp.QuoteMeta("SELECT id, register_id, description, created_at, applicant_id, template_id, stage_id FROM `documents` WHERE `documents`.`deleted_at` IS NULL ORDER BY created_at desc")
	queryPreloadAplicant := regexp.QuoteMeta("SELECT id, username, name FROM `users` WHERE `users`.`id` = ? AND `users`.`deleted_at` IS NULL")
	queryPreloadStage := regexp.QuoteMeta("SELECT id, status FROM `stages` WHERE `stages`.`id` = ? AND `stages`.`deleted_at` IS NULL")
	queryPreloadtemplate := regexp.QuoteMeta("SELECT name FROM `templates` WHERE `templates`.`id` = ? AND `templates`.`deleted_at` IS NULL")
	queryPreloadRegister := regexp.QuoteMeta("SELECT * FROM `registers` WHERE `registers`.`id` = ? AND `registers`.`deleted_at` IS NULL")

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
					ID:         "1",
					RegisterID: 123,
					Register: entity.Register{
						Model: gorm.Model{
							ID: 123,
						},
						Description: "description",
					},
					Description: "description",
					CreatedAt:   time.Time{},
				},
			},
			ReturnedRows: sqlmock.NewRows([]string{"id", "register_id", "description", "created_at"}).AddRow(1, 123, "description", time.Time{}),
		},
		{
			Name:           "Error No rows in result set",
			Err:            nil,
			ExpectedErr:    error2.ErrDocumentNotFound,
			ExpectedReturn: &entity.Documents{},
			ReturnedRows:   sqlmock.NewRows([]string{"id", "register_id", "description", "created_at"}),
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
				s.mock.ExpectQuery(queryPreloadRegister).WillReturnRows(sqlmock.NewRows([]string{"id", "description"}).AddRow(123, "description"))
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
	query := regexp.QuoteMeta("SELECT id, register_id, description, created_at, applicant_id, template_id, stage_id FROM `documents` WHERE applicant_id = ? AND `documents`.`deleted_at` IS NULL ORDER BY created_at desc")
	queryPreloadAplicant := regexp.QuoteMeta("SELECT id, username, name FROM `users` WHERE `users`.`id` = ? AND `users`.`deleted_at` IS NULL")
	queryPreloadStage := regexp.QuoteMeta("SELECT id, status FROM `stages` WHERE `stages`.`id` = ? AND `stages`.`deleted_at` IS NULL")
	queryPreloadtemplate := regexp.QuoteMeta("SELECT name FROM `templates` WHERE `templates`.`id` = ? AND `templates`.`deleted_at` IS NULL")
	queryPreloadRegister := regexp.QuoteMeta("SELECT * FROM `registers` WHERE `registers`.`id` = ? AND `registers`.`deleted_at` IS NULL")

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
					ID:         "1",
					RegisterID: 123,
					Register: entity.Register{
						Model: gorm.Model{
							ID: 123,
						},
						Description: "description",
					},
					Description: "description",
					CreatedAt:   time.Time{},
				},
			},
			ReturnedRows: sqlmock.NewRows([]string{"id", "register_id", "description", "created_at"}).AddRow(1, 123, "description", time.Time{}),
		},
		{
			Name:           "Error No rows in result set",
			Err:            nil,
			ExpectedErr:    error2.ErrDocumentNotFound,
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
				s.mock.ExpectQuery(queryPreloadRegister).WillReturnRows(sqlmock.NewRows([]string{"id", "description"}).AddRow(123, "description"))
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
				RegisterID:  123,
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
			ExpectedErr: error2.ErrDocumentNotFound,
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
				s.mock.ExpectQuery(query).WillReturnRows(sqlmock.NewRows([]string{"id", "register_id", "description", "applicant_id", "template_id", "stage_id", "verifier_id", "signer_id"}).
					AddRow(1, 123, "description", "1", 1, 3, "1", "1"))
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
			ExpectedErr: error2.ErrDocumentNotFound,
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
			ExpectedErr: error2.ErrDocumentNotFound,
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
	query := regexp.QuoteMeta("UPDATE `documents` SET `updated_at`=? WHERE id = ? AND `documents`.`deleted_at` IS NULL")

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
			ExpectedErr:  error2.ErrDocumentNotFound,
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
			ExpectedErr:  error2.ErrDocumentNotFound,
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

func (s *TestSuiteDocumentRepository) TestDeleteDocument() {
	query := regexp.QuoteMeta("UPDATE `documents` SET `deleted_at`=? WHERE id = ? AND `documents`.`deleted_at` IS NULL")

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
			ExpectedErr:  error2.ErrDocumentNotFound,
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

			err := s.documentRepository.DeleteDocument(context.Background(), "1")

			if tc.ExpectedErr != nil {
				s.Equal(tc.ExpectedErr, err)
			}
		})
		s.TearDownTest()
	}
}

func (s *TestSuiteDocumentRepository) TestUpdateDocument() {
	query := regexp.QuoteMeta("UPDATE `documents` SET `updated_at`=? WHERE id = ? AND `documents`.`deleted_at` IS NULL")

	for _, tc := range []struct {
		Name         string
		Err          error
		ExpectedErr  error
		RowsAffected int64
	}{
		{
			Name:         "Success",
			Err:          nil,
			ExpectedErr:  nil,
			RowsAffected: 1,
		},
		{
			Name:         "Error No rows affected",
			Err:          nil,
			ExpectedErr:  error2.ErrDocumentNotFound,
			RowsAffected: 0,
		},
		{
			Name:         "Error generic error",
			Err:          errors.New("generic error"),
			ExpectedErr:  errors.New("generic error"),
			RowsAffected: 0,
		},
	} {
		s.SetupTest()
		s.Run(tc.Name, func() {
			if tc.Err != nil {
				s.mock.ExpectExec(query).WillReturnError(tc.Err)
			} else {
				s.mock.ExpectExec(query).WillReturnResult(sqlmock.NewResult(1, tc.RowsAffected))
			}

			err := s.documentRepository.UpdateDocument(context.Background(), &entity.Document{})

			if tc.ExpectedErr != nil {
				s.Equal(tc.ExpectedErr, err)
			}
		})
		s.TearDownTest()
	}
}

func (s *TestSuiteDocumentRepository) TestUpdateDocumentFields() {
	query := regexp.QuoteMeta("UPDATE `document_fields` SET `id`=?,`updated_at`=?,`document_id`=?,`template_field_id`=?,`value`=? WHERE id = ? AND document_id = ? AND `document_fields`.`deleted_at` IS NULL")

	for _, tc := range []struct {
		Name         string
		Err          error
		ExpectedErr  error
		RowsAffected int64
	}{
		{
			Name:         "Success",
			Err:          nil,
			ExpectedErr:  nil,
			RowsAffected: 1,
		},
		{
			Name:         "Error No rows affected",
			Err:          nil,
			ExpectedErr:  error2.ErrFieldNotFound,
			RowsAffected: 0,
		},
		{
			Name:         "Error generic error",
			Err:          errors.New("generic error"),
			ExpectedErr:  errors.New("generic error"),
			RowsAffected: 0,
		},
	} {
		s.SetupTest()
		s.Run(tc.Name, func() {
			if tc.Err != nil {
				s.mock.ExpectExec(query).WillReturnError(tc.Err)
			} else {
				s.mock.ExpectExec(query).WillReturnResult(sqlmock.NewResult(1, tc.RowsAffected))
			}

			err := s.documentRepository.UpdateDocumentFields(context.Background(), &entity.DocumentFields{
				{
					Model: gorm.Model{
						ID: 1,
					},
					DocumentID:      "123",
					TemplateFieldID: 1,
					Value:           "values",
				},
			})

			if tc.ExpectedErr != nil {
				s.Equal(tc.ExpectedErr, err)
			}
		})
		s.TearDownTest()
	}
}

func (s *TestSuiteDocumentRepository) TestAddDocumentRegister() {
	query := regexp.QuoteMeta("INSERT INTO `registers` (`created_at`,`updated_at`,`deleted_at`,`description`) VALUES (?,?,?,?)")

	for _, tc := range []struct {
		Name         string
		Err          error
		ExpectedErr  error
		RowsAffected int64
	}{
		{
			Name:         "Success",
			Err:          nil,
			ExpectedErr:  nil,
			RowsAffected: 1,
		},
		{
			Name:         "Error generic error",
			Err:          errors.New("generic error"),
			ExpectedErr:  errors.New("generic error"),
			RowsAffected: 0,
		},
	} {
		s.SetupTest()
		s.Run(tc.Name, func() {
			if tc.Err != nil {
				s.mock.ExpectExec(query).WillReturnError(tc.Err)
			} else {
				s.mock.ExpectExec(query).WillReturnResult(sqlmock.NewResult(1, tc.RowsAffected))
			}

			_, err := s.documentRepository.AddDocumentRegister(context.Background(), &entity.Register{
				Description: "test",
			})

			s.Equal(tc.ExpectedErr, err)
		})
		s.TearDownTest()
	}
}

func TestDocumentRepository(t *testing.T) {
	suite.Run(t, new(TestSuiteDocumentRepository))
}
