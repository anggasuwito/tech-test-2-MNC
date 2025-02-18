package model

import "database/sql"

type BaseModel struct {
	ID        string       `gorm:"column:id;primaryKey;size:36"`
	CreatedAt sql.NullTime `gorm:"column:created_at;autoCreateTime:milli"`
	CreatedBy string       `gorm:"column:created_by;size:36;"`
	UpdatedAt sql.NullTime `gorm:"column:updated_at;autoUpdateTime:milli"`
	UpdatedBy string       `gorm:"column:updated_by;size:36;"`
	DeletedAt sql.NullTime `gorm:"column:deleted_at;index;"`
	DeletedBy string       `gorm:"column:deleted_by;size:36;"`
	Note      string       `gorm:"column:note;size:255;"`
}
