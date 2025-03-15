package games

import (
	"context"
	"database/sql"
	"encoding/json"

	"github.com/coopersmall/subswag/apm"
	"github.com/coopersmall/subswag/db"
	"github.com/coopersmall/subswag/domain"
	"github.com/coopersmall/subswag/domain/game"
	reposdomain "github.com/coopersmall/subswag/repos/domain"
	"github.com/coopersmall/subswag/utils"
)

type GameStateVersionRepo struct {
	*reposdomain.SharedRepo[game.GameStateVersionID, *game.GameStateVersion, db.GameStateVersion]
}

func NewGameStateVersionRepo(
	querier db.IQuerier,
	tracer apm.ITracer,
) *GameStateVersionRepo {
	return &GameStateVersionRepo{
		SharedRepo: reposdomain.NewSharedRepo[game.GameStateVersionID, *game.GameStateVersion, db.GameStateVersion](
			"game_state_version",
			querier,
			tracer,
			convertRowToGameStateVersion,
			convertGameStateVersionToRow,
			func(ctx context.Context, iqro db.ISharedQueriesReadOnly, gsvi game.GameStateVersionID) (db.GameStateVersion, error) {
				return iqro.GetGameStateVersion(ctx, int64(gsvi))
			},
			func(ctx context.Context, iqro db.ISharedQueriesReadOnly) ([]db.GameStateVersion, error) {
				return iqro.GetAllGameStateVersions(ctx)
			},
			func(ctx context.Context, iqrw db.ISharedQueriesReadWrite, gsv db.GameStateVersion) (sql.Result, error) {
				return iqrw.CreateGameStateVersion(ctx, db.CreateGameStateVersionParams{
					ID:          gsv.ID,
					GameStateID: gsv.GameStateID,
					CreatedAt:   gsv.CreatedAt,
					Data:        gsv.Data,
				})
			},
			nil,
			nil,
		),
	}
}

func (r *GameStateVersionRepo) GetVersionsForGameState(ctx context.Context, gameStateId game.GameStateID) ([]*game.GameStateVersion, error) {
	return r.SharedRepo.Query(ctx, func(ctx context.Context, queries db.ISharedQueriesReadOnly) ([]db.GameStateVersion, error) {
		return queries.GetGameStateVersionsByGameStateID(ctx, int64(gameStateId))
	})
}

func (r *GameStateVersionRepo) GetLatestVersion(ctx context.Context, gameStateId game.GameStateID) (*game.GameStateVersion, error) {
	found, err := r.SharedRepo.Query(ctx, func(ctx context.Context, queries db.ISharedQueriesReadOnly) ([]db.GameStateVersion, error) {
		found, err := queries.GetLatestGameStateVersionByGameStateID(ctx, int64(gameStateId))
		return []db.GameStateVersion{found}, err
	})
	if len(found) == 0 {
		return nil, utils.NewNotFoundError("game state version not found")
	}
	return found[0], err
}

func convertRowToGameStateVersion(result db.GameStateVersion) (*game.GameStateVersion, error) {
	var gameState *game.GameState
	err := utils.Unmarshal(result.Data, &gameState)
	if err != nil {
		return nil, err
	}

	return &game.GameStateVersion{
		ID:    game.GameStateVersionID(result.ID),
		State: gameState,
		Metadata: &domain.Metadata{
			CreatedAt: result.CreatedAt,
			UpdatedAt: result.UpdatedAt.Time,
		},
	}, nil
}

func convertGameStateVersionToRow(version *game.GameStateVersion) (db.GameStateVersion, error) {
	data, err := utils.Marshal(version.State)
	if err != nil {
		return db.GameStateVersion{}, err
	}

	return db.GameStateVersion{
		ID:          int64(version.ID),
		GameStateID: int64(version.State.ID),
		CreatedAt:   version.Metadata.CreatedAt,
		UpdatedAt: sql.NullTime{
			Time:  version.Metadata.UpdatedAt,
			Valid: !version.Metadata.UpdatedAt.IsZero(),
		},
		Data: json.RawMessage(data),
	}, err
}
