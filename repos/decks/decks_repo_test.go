package deck_test

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/coopersmall/subswag/domain"
	"github.com/coopersmall/subswag/domain/card"
	"github.com/coopersmall/subswag/domain/user"
	"github.com/coopersmall/subswag/repos"
	"github.com/coopersmall/subswag/utils"
	"github.com/stretchr/testify/assert"
)

var (
	ctx       = context.Background()
	userId    = user.UserID(1)
	deckId    = card.SerializableDeckID(1)
	validUser = user.NewUser()
	validDeck = &card.SerializableDeck{
		ID: deckId,
		SerializableDeckData: card.SerializableDeckData{
			UserID:      userId,
			CardIDs:     []card.SerializableCardID{1, 2, 3},
			Name:        "Test Deck",
			Favorited:   true,
			GamesPlayed: 10,
			GamesWon:    7,
		},
		Metadata: &domain.Metadata{
			CreatedAt: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
			UpdatedAt: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
		},
	}
	repo repos.IDecksRepo
)

func (s *DecksRepoTestSuite) SetupSubTest() {
	s.Reset()
	repos, _ := s.GetRepos()
	repo = repos.DecksRepo(userId)
	validUser.ID = userId
	err := repos.UsersRepo().Create(ctx, validUser)
	assert.NoError(s.T(), err)
}

func (s *DecksRepoTestSuite) TestDecksRepo() {
	s.Run("it works", func() {
		err := repo.Create(ctx, validDeck)
		assert.NoError(s.T(), err)

		result, err := repo.Get(ctx, deckId)
		assert.NoError(s.T(), err)
		assert.NotNil(s.T(), result)
		assert.Equal(s.T(), deckId, result.ID)
		assert.Equal(s.T(), validDeck.SerializableDeckData, result.SerializableDeckData)

		deck2 := &card.SerializableDeck{}
		utils.DeepClone(validDeck, deck2)
		deck2.ID = 2

		err = repo.Create(ctx, deck2)
		assert.NoError(s.T(), err)

		results, err := repo.All(ctx)
		assert.NoError(s.T(), err)
		assert.Len(s.T(), results, 2)

		updatedDeck := &card.SerializableDeck{}
		utils.DeepClone(validDeck, updatedDeck)
		updatedDeck.Name = "Updated Deck"
		updatedDeck.Favorited = false
		updatedDeck.GamesPlayed = 15
		updatedDeck.GamesWon = 10

		err = repo.Update(ctx, updatedDeck)
		assert.NoError(s.T(), err)

		result, err = repo.Get(ctx, deckId)
		assert.NoError(s.T(), err)
		assert.NotNil(s.T(), result)
		assert.Equal(s.T(), updatedDeck.Name, result.Name)
		assert.Equal(s.T(), updatedDeck.Favorited, result.Favorited)
		assert.Equal(s.T(), updatedDeck.GamesPlayed, result.GamesPlayed)
		assert.Equal(s.T(), updatedDeck.GamesWon, result.GamesWon)

		err = repo.Delete(ctx, deckId)
		assert.NoError(s.T(), err)

		result, err = repo.Get(ctx, deckId)
		assert.Error(s.T(), err)
		assert.True(s.T(), errors.Is(err, sql.ErrNoRows))
		assert.Nil(s.T(), result)
	})

	s.Run("it deletes deck when user is deleted", func() {
		err := repo.Create(ctx, validDeck)
		assert.NoError(s.T(), err)

		repos, _ := s.GetRepos()
		err = repos.UsersRepo().Delete(ctx, userId)
		assert.NoError(s.T(), err)

		result, err := repo.Get(ctx, deckId)
		assert.Error(s.T(), err)
		assert.Nil(s.T(), result)
	})
}

func (s *DecksRepoTestSuite) TestDecksRepoFailure() {
	s.Run("it fails for non-existent user", func() {
		repos, _ := s.GetRepos()
		err := repos.UsersRepo().Delete(ctx, userId)
		assert.NoError(s.T(), err)

		err = repo.Create(ctx, validDeck)
		assert.Error(s.T(), err)
	})

	s.Run("it fails for non-existent deck", func() {
		result, err := repo.Get(ctx, deckId)
		assert.Error(s.T(), err)
		assert.True(s.T(), errors.Is(err, sql.ErrNoRows))
		assert.Nil(s.T(), result)

		err = repo.Update(ctx, validDeck)
		assert.Error(s.T(), err)

		err = repo.Delete(ctx, deckId)
		assert.Error(s.T(), err)
	})

	s.Run("it fails for duplicate deck", func() {
		err := repo.Create(ctx, validDeck)
		assert.NoError(s.T(), err)

		err = repo.Create(ctx, validDeck)
		assert.Error(s.T(), err)
	})
}
