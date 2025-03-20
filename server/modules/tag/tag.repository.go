package tag

import (
	"strings"

	"github.com/ednaldo-dilorenzo/iappointment/config"
	"github.com/ednaldo-dilorenzo/iappointment/model"
)

type TagRepository interface {
	FindAllWithFilter(filter string, userId int) ([]model.TransactionTag, error)
}

type tagRespository struct {
	dbConfig *config.Database
}

func NewTagRepository(database *config.Database) TagRepository {
	return &tagRespository{
		database,
	}
}

func (t *tagRespository) FindAllWithFilter(filter string, userId int) ([]model.TransactionTag, error) {
	var items []model.TransactionTag

	if filter == "" {
		return []model.TransactionTag{}, nil
	}

	query := t.dbConfig.DB.Model(&model.TransactionTag{}).
		Distinct("tag"). // Ensure only distinct tags are selected
		Where("user_id = ? AND LOWER(tag) LIKE ?", userId, "%"+strings.ToLower(filter)+"%")

	if err := query.Find(&items).Error; err != nil {
		return nil, err
	}

	return items, nil
}
