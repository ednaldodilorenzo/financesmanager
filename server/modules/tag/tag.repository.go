package tag

import (
	"strings"

	"github.com/ednaldo-dilorenzo/iappointment/config"
	"github.com/ednaldo-dilorenzo/iappointment/model"
)

type TagRepository interface {
	FindAllWithFilter(filter string) ([]model.TransactionTag, error)
}

type TagRespositoryStruct struct {
	dbConfig *config.Database
}

func NewTagRepository(database *config.Database) TagRepository {
	return &TagRespositoryStruct{
		database,
	}
}

func (t *TagRespositoryStruct) FindAllWithFilter(filter string) ([]model.TransactionTag, error) {
	var items []model.TransactionTag

	if filter == "" {
		return []model.TransactionTag{}, nil
	}

	query := t.dbConfig.DB.Model(&model.TransactionTag{}).
		Distinct("tag"). // Ensure only distinct tags are selected
		Where("LOWER(tag) LIKE ?", "%"+strings.ToLower(filter)+"%")

	if err := query.Find(&items).Error; err != nil {
		return nil, err
	}

	return items, nil
}
