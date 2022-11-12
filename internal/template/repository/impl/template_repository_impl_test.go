package impl

import (
	"context"
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/suite"
	"github.com/suryaadi44/eAD-System/pkg/entity"
	"github.com/suryaadi44/eAD-System/pkg/utils"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"regexp"
	"testing"
)

type TestSuiteTemplateRepository struct {
	suite.Suite
	mock                   sqlmock.Sqlmock
	templateRepositoryImpl *TemplateRepositoryImpl
}

func (s *TestSuiteTemplateRepository) SetupTest() {
	dbMock, mock, err := sqlmock.New()
	s.NoError(err)
	s.mock = mock

	DB, _ := gorm.Open(mysql.New(mysql.Config{
		Conn:                      dbMock,
		SkipInitializeWithVersion: true,
	}), &gorm.Config{
		SkipDefaultTransaction: true,
	})

	s.templateRepositoryImpl = &TemplateRepositoryImpl{db: DB}
}

func (s *TestSuiteTemplateRepository) TearDownTest() {
	s.mock = nil
	s.templateRepositoryImpl = nil
}

func (s *TestSuiteTemplateRepository) TestAddTemplate() {
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

			err := s.templateRepositoryImpl.AddTemplate(context.Background(), &entity.Template{})

			s.Equal(tc.ExpectedErr, err)
		})
		s.TearDownTest()
	}
}

func (s *TestSuiteTemplateRepository) TestGetAllTemplate() {
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

			result, err := s.templateRepositoryImpl.GetAllTemplate(context.Background())

			if tc.ExpectedErr != nil {
				s.Equal(tc.ExpectedErr, err)
			} else {
				s.Equal(tc.ExpectedReturn, result)
			}
		})
		s.TearDownTest()
	}
}

func (s *TestSuiteTemplateRepository) TestGetTemplateDetail() {
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

			result, err := s.templateRepositoryImpl.GetTemplateDetail(context.Background(), 1)

			if tc.ExpectedErr != nil {
				s.Equal(tc.ExpectedErr, err)
			} else {
				s.Equal(tc.ExpectedReturn, result)
			}
		})
		s.TearDownTest()
	}
}

func (s *TestSuiteTemplateRepository) TestGetTemplateFields() {
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

			result, err := s.templateRepositoryImpl.GetTemplateFields(context.Background(), 1)

			if tc.ExpectedErr != nil {
				s.Equal(tc.ExpectedErr, err)
			} else {
				s.Equal(tc.ExpectedReturn, result)
			}
		})
		s.TearDownTest()
	}
}

func TestTemplateRepository(t *testing.T) {
	suite.Run(t, new(TestSuiteTemplateRepository))
}
