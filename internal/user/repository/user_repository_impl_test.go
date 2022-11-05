package repository

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

type TestSuiteUserRepository struct {
	suite.Suite
	mock           sqlmock.Sqlmock
	userRepository UserRepository
}

func (s *TestSuiteUserRepository) SetupTest() {
	dbMock, mock, err := sqlmock.New()
	s.NoError(err)
	s.mock = mock

	DB, _ := gorm.Open(mysql.New(mysql.Config{
		Conn:                      dbMock,
		SkipInitializeWithVersion: true,
	}), &gorm.Config{})

	s.userRepository = NewUserRepositoryImpl(DB)
}

func (s *TestSuiteUserRepository) TeardownTest() {
	s.mock = nil
	s.userRepository = nil
}

//goland:noinspection GoErrorStringFormat,GoErrorStringFormat,GoErrorStringFormat,GoErrorStringFormat,GoErrorStringFormat
func (s *TestSuiteUserRepository) TestCreateUser() {
	query := regexp.QuoteMeta("INSERT INTO `users` (`id`,`n_ip`,`nik`,`username`,`password`,`role`,`name`,`telp`,`sex`,`address`,`created_at`,`updated_at`,`deleted_at`) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?)")
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
			Name:        "Error duplicate username",
			Err:         errors.New("Error 1062: Duplicate entry '' for key 'username'"),
			ExpectedErr: utils.ErrUsernameAlreadyExist,
		},
		{
			Name:        "Error duplicate nip",
			Err:         errors.New("Error 1062: Duplicate entry '' for key 'n_ip'"),
			ExpectedErr: utils.ErrNIPAlreadyExist,
		},
		{
			Name:        "Error duplicate nik",
			Err:         errors.New("Error 1062: Duplicate entry '' for key 'nik'"),
			ExpectedErr: utils.ErrNIKAlreadyExist,
		},
		{
			Name:        "Generic error",
			Err:         errors.New("Generic error"),
			ExpectedErr: errors.New("Generic error"),
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

			err := s.userRepository.CreateUser(context.Background(), &entity.User{})

			s.Equal(tc.ExpectedErr, err)
		})
		s.TeardownTest()
	}
}

//goland:noinspection GoErrorStringFormat,GoErrorStringFormat
func (s *TestSuiteUserRepository) TestFindByUsername() {
	query := regexp.QuoteMeta("SELECT `id`,`username`,`password`,`role` FROM `users` WHERE username = ? AND `users`.`deleted_at` IS NULL ORDER BY `users`.`id` LIMIT 1")
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
			Name:        "Error no record found",
			Err:         gorm.ErrRecordNotFound,
			ExpectedErr: utils.ErrUserNotFound,
		},
		{
			Name:        "Generic error",
			Err:         errors.New("Generic error"),
			ExpectedErr: errors.New("Generic error"),
		},
	} {
		s.SetupTest()
		s.Run(tc.Name, func() {
			if tc.Err != nil {
				s.mock.ExpectQuery(query).WillReturnError(tc.Err)
			} else {
				s.mock.ExpectQuery(query).WillReturnRows(sqlmock.NewRows([]string{"id", "username", "password", "role"}).AddRow(1, "123", "123", 1))
			}

			_, err := s.userRepository.FindByUsername(context.Background(), "")

			s.Equal(tc.ExpectedErr, err)
		})
		s.TeardownTest()
	}
}

func TestUserRepository(t *testing.T) {
	suite.Run(t, new(TestSuiteUserRepository))
}
