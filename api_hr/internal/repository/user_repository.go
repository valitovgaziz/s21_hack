package repository

import (
    "api_hr/internal/models"

    "gorm.io/gorm"
)

type UserRepository struct {
    db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
    return &UserRepository{db: db}
}

func (r *UserRepository) Create(user *models.UserT) error {
    return r.db.Create(user).Error
}

func (r *UserRepository) FindByID(id uint) (*models.UserT, error) {
    var user models.UserT
    err := r.db.First(&user, id).Error
    return &user, err
}

func (r *UserRepository) FindByEmail(email string) (*models.UserT, error) {
    var user models.UserT
    err := r.db.Where("email = ?", email).First(&user).Error
    return &user, err
}

func (r *UserRepository) FindAll() ([]models.UserT, error) {
    var users []models.UserT
    err := r.db.Find(&users).Error
    return users, err
}

func (r *UserRepository) Update(user *models.UserT) error {
    return r.db.Save(user).Error
}

func (r *UserRepository) Delete(id uint) error {
    return r.db.Delete(&models.UserT{}, id).Error
}