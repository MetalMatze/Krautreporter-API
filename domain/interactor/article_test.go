package interactor

import (
	"github.com/MetalMatze/Krautreporter-API/domain/entity"
	"testing"
)

type testArticleRepository struct {
}

func (r testArticleRepository) SaveAll([]entity.Article) error {
	return nil
}

func (r testArticleRepository) FindOlderThan(id int, number int) ([]*entity.Article, error) {
	return nil, nil
}

func (r testArticleRepository) FindByID(int) (*entity.Article, error) {
	return nil, nil
}

func TestArticleInteractor_SaveAll(t *testing.T) {
	ai := ArticleInteractor{ArticleRepository: testArticleRepository{}}

	if err := ai.SaveAll([]entity.Article{}); err != nil {
		t.Errorf("ArticleInteractor.SaveAll() shouldn't return an err: %s", err.Error())
	}
}

func TestArticleInteractor_FindOlderThan(t *testing.T) {
	ai := ArticleInteractor{ArticleRepository: testArticleRepository{}}

	if _, err := ai.FindOlderThan(1000, 20); err != nil {
		t.Errorf("ArticleInteractor.FindOlderThan() shouldn't return an err: %s", err.Error())
	}
}
func TestArticleInteractor_FindByID(t *testing.T) {
	ai := ArticleInteractor{ArticleRepository: testArticleRepository{}}

	if _, err := ai.FindByID(1000); err != nil {
		t.Errorf("ArticleInteractor.FindByID() shouldn't return an err: %s", err.Error())
	}
}
