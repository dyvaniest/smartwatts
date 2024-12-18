package repository

import (
	"a21hc3NpZ25tZW50/model"
	"fmt"

	"gorm.io/gorm"
)

type SessionsRepository interface {
	AddSessions(session model.Session) error
	DeleteSession(token string) error
	UpdateSessions(session model.Session) error
	SessionAvailEmail(email string) error
	SessionAvailToken(token string) (model.Session, error)
}

type sessionsRepoImpl struct {
	db *gorm.DB
}

func NewSessionRepo(db *gorm.DB) *sessionsRepoImpl {
	return &sessionsRepoImpl{db}
}

func (s *sessionsRepoImpl) AddSessions(session model.Session) error {
	// menyimpan sessions ke database
	if result := s.db.Create(&session); result.Error != nil {
		return result.Error
	}
	return nil
}

func (s *sessionsRepoImpl) DeleteSession(token string) error {
	// menghapus sessions dari database
	sesi := model.Session{}
	if result := s.db.Where("token = ?", token).Delete(&sesi); result.Error != nil {
		return result.Error
	}
	return nil
}

func (s *sessionsRepoImpl) UpdateSessions(session model.Session) error {
	// mengupdate sessions ke database sesuai parameter
	if result := s.db.Model(&model.Session{}).Where("email", session.Email).Updates(session); result.Error != nil {
		return result.Error
	}

	return nil
}

func (s *sessionsRepoImpl) SessionAvailEmail(email string) error {
	// cek apakah token sessions sudah ada di database dengan kolom name
	sessi := model.Session{}
	if result := s.db.Where("email = ?", email).First(&sessi); result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return fmt.Errorf("sessions not found for user %v", email)
		}
		return result.Error
	}
	return nil
}

func (s *sessionsRepoImpl) SessionAvailToken(token string) (model.Session, error) {
	// cek apakah token sessions sudah ada di database
	sessi := model.Session{}
	if result := s.db.Where("token = ?", token).First(&sessi); result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return model.Session{}, fmt.Errorf("sessions not found for user %v", token)
		}
		return model.Session{}, result.Error
	}
	return sessi, nil
}
