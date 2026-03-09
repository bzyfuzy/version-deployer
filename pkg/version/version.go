package version

import (
	"github.com/google/uuid"
)

type Version struct {
	ID uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key"`
}
