package lms

import (
	"gorm.io/gorm"
)

type LMSRepository interface {
	GetByName(name string) (*LMS, error)
	Create(lms *LMS) error
	CreateMany(lmsList []*LMS) error // <-- new method
	Update(lms *LMS) error
	ListAll() ([]LMS, error)
}

type GormLMSRepository struct {
	db *gorm.DB
}

func NewGormLMSRepository(db *gorm.DB) *GormLMSRepository {
	return &GormLMSRepository{db: db}
}

func (r *GormLMSRepository) GetByName(name string) (*LMS, error) {
	var lms LMS
	if err := r.db.Where("name = ?", name).First(&lms).Error; err != nil {
		return nil, err
	}
	return &lms, nil
}

func (r *GormLMSRepository) Create(lms *LMS) error {
	return r.db.Create(lms).Error
}

// CreateMany inserts multiple LMS entries at once
func (r *GormLMSRepository) CreateMany(lmsList []*LMS) error {
	if len(lmsList) == 0 {
		return nil
	}
	return r.db.Create(&lmsList).Error
}

func (r *GormLMSRepository) Update(lms *LMS) error {
	return r.db.Save(lms).Error
}

func (r *GormLMSRepository) ListAll() ([]LMS, error) {
	var list []LMS
	err := r.db.Find(&list).Error
	return list, err
}
