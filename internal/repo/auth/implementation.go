package auth

import (
	"errors"
	"github.com/SwanHtetAungPhyo/backend/internal/model"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func (ar *Repository) SignUp(req *model.User) error {
	if err := ar.gorm.Create(req).Error; err != nil {
		ar.log.WithError(err).Error("SignUp failed")
		return err
	}
	ar.log.WithFields(logrus.Fields{
		"email": req.Email,
	}).Info("SignUp successful")
	return nil
}

func (ar *Repository) UpdateAccountVerificationStatus(email string) error {
	if err := ar.gorm.Model(&model.User{}).
		Where("email = ?", email).
		Update("verified", true).Error; err != nil {
		ar.log.WithError(err).Error("UpdateAccountVerificationStatus failed")

		return err
	}
	return nil
}

func (ar *Repository) UpdateTheKYCVerificationStatus(email string) error {
	if err := ar.gorm.Model(&model.User{}).
		Where("email = ?", email).
		Update("twoFactorVerified", true).Error; err != nil {
		ar.log.WithError(err).Error("UpdateTheKYCVerificationStatus failed")
		return err
	}
	return nil
}

func (ar *Repository) GetKYCVerifiedStatus(email string) (bool, error) {
	var user model.User

	if err := ar.gorm.Model(&model.User{}).Where("email = ?", email).First(&user).Error; err != nil {
		return false, err
	}
	return user.TwoFactorVerified, nil
}

func (ar *Repository) CheckUserExistence(email string) (*model.User, error) {
	var user model.User
	if err := ar.gorm.Model(&model.User{}).Where("email = ?", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
	}
	return &user, errors.New("user exists")
}

func (ar *Repository) SaveBiometrics(biometric model.Biometrics) error {
	if err := ar.gorm.Model(&model.Biometrics{}).
		Create(&biometric).Error; err != nil {
		ar.log.WithError(err).Error("SaveBiometrics failed")
		return err
	}
	return nil
}

func (ar *Repository) UpdateBioMetricsVerification(email string) error {
	tx := ar.gorm.Begin()

	var user model.User
	if err := tx.Model(&model.User{}).Where("email = ?", email).First(&user).Error; err != nil {
		ar.log.WithError(err).Error("UpdateBioMetricsVerification failed to find user")
		tx.Rollback()
		return err
	}

	if err := tx.Model(&model.Biometrics{}).
		Where("cognito_user_name = ?", user.CognitoUsername).
		Update("isVerified", true).Error; err != nil {
		ar.log.WithError(err).Error("UpdateBioMetricsVerification failed to update biometrics")
		tx.Rollback()
		return err
	}

	if err := tx.Commit().Error; err != nil {
		ar.log.WithError(err).Error("UpdateBioMetricsVerification commit failed")
		return err
	}

	return nil
}
