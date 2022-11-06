package repository

import (
	"context"
	"errors"
	"regexp"
	"testing"

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
	}), &gorm.Config{})

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
			s.mock.ExpectBegin()
			if tc.Err != nil {
				s.mock.ExpectExec(query).WillReturnError(tc.Err)
				s.mock.ExpectRollback()
			} else {
				s.mock.ExpectExec(query).WillReturnResult(sqlmock.NewResult(1, 1))
				s.mock.ExpectCommit()
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
		ReturnedRow      *sqlmock.Rows
		ReturnedRowField *sqlmock.Rows
	}{
		{
			Name:        "Success",
			Err:         nil,
			ExpectedErr: nil,
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

			_, err := s.documentRepository.GetAllTemplate(context.Background())

			s.Equal(tc.ExpectedErr, err)
		})
		s.TearDownTest()
	}
}

func TestDocumentRepository(t *testing.T) {
	suite.Run(t, new(TestSuiteDocumentRepository))
}
