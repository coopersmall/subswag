package card_test

import (
	"context"
	"time"

	"github.com/coopersmall/subswag/domain"
	"github.com/coopersmall/subswag/domain/card"
	"github.com/coopersmall/subswag/repos"
	"github.com/coopersmall/subswag/utils"
	"github.com/stretchr/testify/assert"
)

var (
	ctx           = context.Background()
	cardId        = card.SerializableCardID(1)
	validFaceCard = &card.SerializableFaceCard{
		SerializableCardBaseData: card.SerializableCardBaseData{
			ID: cardId,
			SerializableCardData: card.SerializableCardData{
				ArtworkURL: "https://example.com/art.jpg",
				Suite:      card.CardSuiteHearts,
				Rarity:     card.CardRarityCommon,
				Tribe:      card.CardTribeMilitary,
			},
			Metadata: &domain.Metadata{
				CreatedAt: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
				UpdatedAt: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
			},
		},
		Type: card.SerializableCardTypeFace,
		Face: "King",
	}
	validNumberCard = &card.SerializableNumberCard{
		SerializableCardBaseData: card.SerializableCardBaseData{
			ID: cardId + 1,
			SerializableCardData: card.SerializableCardData{
				ArtworkURL: "https://example.com/art2.jpg",
				Suite:      card.CardSuiteDiamonds,
				Rarity:     card.CardRarityRare,
				Tribe:      card.CardTribeTech,
			},
			Metadata: &domain.Metadata{
				CreatedAt: time.Date(2021, 1, 1, 0, 0, 0, 1, time.UTC),
				UpdatedAt: time.Date(2021, 1, 1, 0, 0, 0, 1, time.UTC),
			},
		},
		Type:   card.SerializableCardTypeNumber,
		Number: 7,
	}
	repo repos.ICardsRepo
)

func (s *CardsRepoTestSuite) SetupSubTest() {
	s.Reset()
	repos, _ := s.GetRepos()
	repo = repos.CardsRepo()
}

func (s *CardsRepoTestSuite) TestCardsRepoSuccess() {
	s.Run("it works", func() {
		err := repo.Create(ctx, validFaceCard)
		assert.NoError(s.T(), err)

		result, err := repo.Get(ctx, cardId)
		assert.NoError(s.T(), err)
		assert.Equal(s.T(), validFaceCard.ID, result.GetID())
		assert.Equal(s.T(), validFaceCard.Type, result.GetType())

		err = repo.Create(ctx, validNumberCard)
		assert.NoError(s.T(), err)

		results, err := repo.All(ctx)
		assert.NoError(s.T(), err)
		assert.Len(s.T(), results, 2)
		assert.Equal(s.T(), validFaceCard.ID, results[0].GetID())
		assert.Equal(s.T(), validNumberCard.ID, results[1].GetID())

		updatedFaceCard := &card.SerializableFaceCard{}
		utils.DeepClone(validFaceCard, updatedFaceCard)
		updatedFaceCard.Face = "Queen"

		err = repo.Update(ctx, updatedFaceCard)
		assert.NoError(s.T(), err)

		result, err = repo.Get(ctx, cardId)
		assert.NoError(s.T(), err)
		assert.Equal(s.T(), updatedFaceCard.ID, result.GetID())

		err = repo.Delete(ctx, cardId)
		assert.NoError(s.T(), err)

		results, err = repo.All(ctx)
		assert.NoError(s.T(), err)
		assert.Len(s.T(), results, 1)

		err = repo.Delete(ctx, validNumberCard.ID)
		assert.NoError(s.T(), err)

		results, err = repo.All(ctx)
		assert.NoError(s.T(), err)
		assert.Len(s.T(), results, 0)
	})
}

func (s *CardsRepoTestSuite) TestCardsRepoFailure() {
	s.Run("it fails with non-existent card", func() {
		result, err := repo.Get(ctx, cardId)
		assert.Error(s.T(), err)
		assert.Nil(s.T(), result)
	})

	s.Run("it fails with duplicate card", func() {
		err := repo.Create(ctx, validFaceCard)
		assert.NoError(s.T(), err)

		err = repo.Create(ctx, validFaceCard)
		assert.Error(s.T(), err)
	})
}
