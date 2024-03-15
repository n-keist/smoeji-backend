package db_plugins

import "gorm.io/gorm"

type UuidDbPlugin struct{}

func (this *UuidDbPlugin) Name() string {
	return "uuid generation enabler plugin"
}

func (this *UuidDbPlugin) Initialize(db *gorm.DB) error {
	db.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\";")
	return nil
}
