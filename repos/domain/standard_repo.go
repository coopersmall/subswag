package domain

import (
	"context"
	"database/sql"

	"github.com/coopersmall/subswag/apm"
	"github.com/coopersmall/subswag/db"
	"github.com/coopersmall/subswag/domain/user"
	"github.com/coopersmall/subswag/utils"
)

// StandardRepo provides a standard implementation of common repository operations
type StandardRepo[ID any, Data any, DBData any] struct {
	name       string
	querier    db.IQuerier
	tracer     apm.ITracer
	userId     user.UserID
	convertRow func(DBData) (Data, error)
	toRow      func(Data) (DBData, error)
	isEmpty    func(Data) bool
	getFunc    func(context.Context, db.IStandardQueriesReadOnly, ID) (DBData, error)
	getAllFunc func(context.Context, db.IStandardQueriesReadOnly) ([]DBData, error)
	createFunc func(context.Context, db.IStandardQueriesReadWrite, DBData) (sql.Result, error)
	updateFunc func(context.Context, db.IStandardQueriesReadWrite, DBData) (sql.Result, error)
	deleteFunc func(context.Context, db.IStandardQueriesReadWrite, ID) (sql.Result, error)
}

// NewStandardRepo creates a new StandardRepo instance
func NewStandardRepo[ID any, Data any, DBData any](
	name string,
	querier db.IQuerier,
	tracer apm.ITracer,
	userId user.UserID,
	convertRow func(DBData) (Data, error),
	toRow func(Data) (DBData, error),
	isEmpty func(Data) bool,
	getFunc func(context.Context, db.IStandardQueriesReadOnly, ID) (DBData, error),
	getAllFunc func(context.Context, db.IStandardQueriesReadOnly) ([]DBData, error),
	createFunc func(context.Context, db.IStandardQueriesReadWrite, DBData) (sql.Result, error),
	updateFunc func(context.Context, db.IStandardQueriesReadWrite, DBData) (sql.Result, error),
	deleteFunc func(context.Context, db.IStandardQueriesReadWrite, ID) (sql.Result, error),
) *StandardRepo[ID, Data, DBData] {
	return &StandardRepo[ID, Data, DBData]{
		name:       name,
		querier:    querier,
		tracer:     tracer,
		userId:     userId,
		convertRow: convertRow,
		toRow:      toRow,
		isEmpty:    isEmpty,
		getFunc:    getFunc,
		getAllFunc: getAllFunc,
		createFunc: createFunc,
		updateFunc: updateFunc,
		deleteFunc: deleteFunc,
	}
}

// Get retrieves a single item by ID
func (r *StandardRepo[ID, Data, DBData]) Get(
	ctx context.Context,
	id ID,
) (Data, error) {
	var (
		result DBData
		data   Data
		err    error
	)
	r.tracer.Trace(ctx, r.name+".get", func(ctx context.Context, span apm.ISpan) error {
		err = r.querier.Standard(ctx, r.userId, func(d db.IStandardQueriesReadOnly) error {
			result, err = r.getFunc(ctx, d, id)
			return err
		})
		if err != nil {
			return utils.NewNotFoundError("not found", err)
		}
		data, err = r.convertRow(result)
		if err != nil {
			err = utils.NewInternalError("failed to convert row", err)
			return err
		}
		if r.isEmpty(data) {
			err = utils.NewNotFoundError("not found", nil)
			return err
		}
		return nil
	})
	return data, err
}

// All retrieves all items
func (r *StandardRepo[ID, Data, DBData]) All(ctx context.Context) ([]Data, error) {
	var (
		results []DBData
		items   []Data
		err     error
	)
	r.tracer.Trace(ctx, r.name+".all", func(ctx context.Context, span apm.ISpan) error {
		err = r.querier.Standard(ctx, r.userId, func(d db.IStandardQueriesReadOnly) error {
			var err error
			results, err = r.getAllFunc(ctx, d)
			return err
		})
		if err != nil {
			return utils.ErrorOrNil("not found", utils.NewNotFoundError, err)
		}

		items = make([]Data, len(results))
		for i, result := range results {
			var item Data
			item, err = r.convertRow(result)
			if err != nil {
				err = utils.NewInternalError("failed to convert row", err)
				return err
			}
			if r.isEmpty(item) {
				err = utils.NewNotFoundError("not found", nil)
				return err
			}
			items[i] = item
		}
		return nil
	})
	return items, err
}

// Create creates a new item
func (r *StandardRepo[ID, Data, DBData]) Create(ctx context.Context, data Data) error {
	var err error
	r.tracer.Trace(ctx, r.name+".create", func(ctx context.Context, span apm.ISpan) error {
		var row DBData
		row, err = r.toRow(data)
		if err != nil {
			err = utils.NewInternalError("failed to convert to row", err)
			return err
		}
		err = r.querier.StandardWrite(ctx, r.userId, func(d db.IStandardQueriesReadWrite) error {
			var rows sql.Result
			rows, err = r.createFunc(ctx, d, row)
			if err != nil {
				return err
			}
			if rowsAffected, _ := rows.RowsAffected(); rowsAffected == 0 {
				return utils.NewNotFoundError("not found", nil)
			}
			return nil
		})
		return err
	})
	return err
}

// Update updates an existing item
func (r *StandardRepo[ID, Data, DBData]) Update(ctx context.Context, data Data) error {
	var err error
	r.tracer.Trace(ctx, r.name+".update", func(ctx context.Context, span apm.ISpan) error {
		var row DBData
		row, err = r.toRow(data)
		if err != nil {
			return utils.NewInternalError("failed to convert to row", err)
		}
		err = r.querier.StandardWrite(ctx, r.userId, func(d db.IStandardQueriesReadWrite) error {
			var rows sql.Result
			rows, err = r.updateFunc(ctx, d, row)
			if err != nil {
				return err
			}
			if rowsAffected, _ := rows.RowsAffected(); rowsAffected == 0 {
				return utils.NewNotFoundError("not found", nil)
			}
			return nil
		})
		return err
	})
	return err
}

// Delete deletes an item by ID
func (r *StandardRepo[ID, Data, DBData]) Delete(ctx context.Context, id ID) error {
	var err error
	r.tracer.Trace(ctx, r.name+".delete", func(ctx context.Context, span apm.ISpan) error {
		err = r.querier.StandardWrite(ctx, r.userId, func(d db.IStandardQueriesReadWrite) error {
			var rows sql.Result
			rows, err = r.deleteFunc(ctx, d, id)
			if err != nil {
				return err
			}
			if rowsAffected, _ := rows.RowsAffected(); rowsAffected == 0 {
				return utils.NewNotFoundError("not found", nil)
			}
			return nil
		})
		return err
	})
	return err
}

// Query performs a custom read query that returns multiple items
func (r *StandardRepo[ID, Data, DBData]) Query(
	ctx context.Context,
	queryFunc func(context.Context, db.IStandardQueriesReadOnly, user.UserID) ([]DBData, error),
) ([]Data, error) {
	var (
		results []DBData
		items   []Data
		err     error
	)
	r.tracer.Trace(ctx, r.name+".query", func(ctx context.Context, span apm.ISpan) error {
		err = r.querier.Standard(ctx, r.userId, func(d db.IStandardQueriesReadOnly) error {
			var err error
			results, err = queryFunc(ctx, d, r.userId)
			return err
		})
		if err != nil {
			return utils.ErrorOrNil("not found", utils.NewNotFoundError, err)
		}
		items = make([]Data, len(results))
		for i, result := range results {
			var item Data
			item, err = r.convertRow(result)
			if err != nil {
				err = utils.NewInternalError("failed to convert row", err)
				return err
			}
			if r.isEmpty(item) {
				err = utils.NewNotFoundError("not found", nil)
				return err
			}
			items[i] = item
		}
		return nil
	})
	return items, err
}

// Execute performs a custom write operation
func (r *StandardRepo[ID, Data, DBData]) Execute(
	ctx context.Context,
	executeFunc func(context.Context, db.IStandardQueriesReadWrite, user.UserID) error,
) error {
	var err error
	r.tracer.Trace(ctx, r.name+".execute", func(ctx context.Context, span apm.ISpan) error {
		err := r.querier.StandardWrite(ctx, r.userId, func(d db.IStandardQueriesReadWrite) error {
			err := executeFunc(ctx, d, r.userId)
			return err
		})
		return utils.ErrorOrNil("operation failed", utils.NewNotFoundError, err)
	})
	return err
}
