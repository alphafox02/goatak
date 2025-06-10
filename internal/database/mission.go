package database

import (
	"fmt"
	"strings"
	"time"

	"gorm.io/gorm"

	"github.com/kdudkov/goatak/pkg/model"
)

func (mm *DatabaseManager) CreateMission(m *model.Mission) error {
	if mm == nil || mm.db == nil {
		return nil
	}

	if m == nil {
		return fmt.Errorf("null mission")
	}

	if m.Name == "" {
		return fmt.Errorf("null mission name")
	}

	err := mm.db.Transaction(func(tx *gorm.DB) error {
		if NewMissionQuery(tx).Scope(m.Scope).Name(m.Name).One() != nil {
			return fmt.Errorf("mission %s exists", m.Name)
		}

		if err := tx.Create(m).Error; err != nil {
			return err
		}

		c := &model.Change{
			CreatedAt:  time.Now(),
			Type:       "CREATE_MISSION",
			MissionID:  m.ID,
			CreatorUID: m.CreatorUID,
		}

		return tx.Create(c).Error
	})

	if err != nil {
		return err
	}

	_, err = mm.subscribe(m.ID, m.CreatorUID, m.Creator, true)

	return err
}

func (mm *DatabaseManager) UpdateKw(name, scope string, kw []string) error {
	return mm.MissionQuery().Name(name).Scope(scope).Update(map[string]any{"keywords": strings.Join(kw, ",")})
}
