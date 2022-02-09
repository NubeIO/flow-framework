package migration

import (
	"github.com/NubeIO/flow-framework/migration/versions"
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func initGorMigrate(db *gorm.DB) *gormigrate.Gormigrate {
	options := &gormigrate.Options{
		TableName:                 "_GMigrationsHistory",
		IDColumnName:              "id",
		IDColumnSize:              255,
		UseTransaction:            true,
		ValidateUnknownMigrations: true,
	}
	m := gormigrate.New(db, options, []*gormigrate.Migration{
		// Migration example
		// versions.Migration2021020770243(),
	})
	m.InitSchema(func(tx *gorm.DB) error {
		interfaces := versions.GetInitInterfaces()
		for _, i := range interfaces {
			if err := tx.AutoMigrate(
				i,
			); err != nil {
				return err
			}
		}
		return nil
	})
	return m
}

func Upgrade(db *gorm.DB) error {
	m := initGorMigrate(db)
	return m.Migrate()
}

func UpgradeTo(db *gorm.DB, id string) error {
	m := initGorMigrate(db)
	return m.MigrateTo(id)
}
