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

type GameStateRepo struct {
	*reposdomain.SharedRepo[game.GameStateID, *game.GameState, db.GameState]
}

func NewGameStateRepo(
	querier db.IQuerier,
	tracer apm.ITracer,
) *GameStateRepo {
	return &GameStateRepo{
		SharedRepo: reposdomain.NewSharedRepo[game.GameStateID, *game.GameState, db.GameState](
			"game_state",
			querier,
			tracer,
			convertRowToGameState,
			convertGameStateToRow,
			func(ctx context.Context, iqro db.ISharedQueriesReadOnly, gsi game.GameStateID) (db.GameState, error) {
				return iqro.GetGameState(ctx, int64(gsi))
			},
			func(ctx context.Context, iqro db.ISharedQueriesReadOnly) ([]db.GameState, error) {
				return iqro.GetAllGameStates(ctx)
			},
			func(ctx context.Context, iqrw db.ISharedQueriesReadWrite, gs db.GameState) (sql.Result, error) {
				return iqrw.CreateGameState(ctx, db.CreateGameStateParams{
					ID:        gs.ID,
					CreatedAt: gs.CreatedAt,
					Data:      gs.Data,
				})
			},
			func(ctx context.Context, iqrw db.ISharedQueriesReadWrite, gs db.GameState) (sql.Result, error) {
				return iqrw.UpdateGameState(ctx, db.UpdateGameStateParams{
					ID:        gs.ID,
					UpdatedAt: gs.UpdatedAt,
					Data:      gs.Data,
				})
			},
			func(ctx context.Context, iqrw db.ISharedQueriesReadWrite, gsi game.GameStateID) (sql.Result, error) {
				return iqrw.DeleteGameState(ctx, int64(gsi))
			},
		),
	}
}

func convertRowToGameState(result db.GameState) (*game.GameState, error) {
	var data game.GameStateData
	err := utils.Unmarshal(result.Data, &data)
	return &game.GameState{
		ID:            game.GameStateID(result.ID),
		GameStateData: data,
		Metadata: &domain.Metadata{
			CreatedAt: result.CreatedAt,
			UpdatedAt: result.UpdatedAt.Time,
		},
	}, err
}

func convertGameStateToRow(gameState *game.GameState) (db.GameState, error) {
	data, err := utils.Marshal(gameState.GameStateData)
	return db.GameState{
		ID:        int64(gameState.ID),
		CreatedAt: gameState.Metadata.CreatedAt,
		UpdatedAt: sql.NullTime{
			Time:  gameState.Metadata.UpdatedAt,
			Valid: !gameState.Metadata.UpdatedAt.IsZero(),
		},
		Data: json.RawMessage(data),
	}, err
}
