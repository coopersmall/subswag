package domain

import (
	"context"
	"database/sql"

	"github.com/coopersmall/subswag/apm"
	"github.com/coopersmall/subswag/db"
	"github.com/coopersmall/subswag/utils"
)

// SharedRepo provides a generic implementation of shared repository operations
type SharedRepo[ID any, Data any, DBData any] struct {
	name       string
	querier    db.IQuerier
	tracer     apm.ITracer
	convertRow func(DBData) (Data, error)
	toRow      func(Data) (DBData, error)
	getFunc    func(context.Context, db.ISharedQueriesReadOnly, ID) (DBData, error)
	getAllFunc func(context.Context, db.ISharedQueriesReadOnly) ([]DBData, error)
	createFunc func(context.Context, db.ISharedQueriesReadWrite, DBData) (sql.Result, error)
	updateFunc func(context.Context, db.ISharedQueriesReadWrite, DBData) (sql.Result, error)
	deleteFunc func(context.Context, db.ISharedQueriesReadWrite, ID) (sql.Result, error)
}

// NewSharedRepo creates a new SharedRepo instance
func NewSharedRepo[ID any, Data any, DBData any](
	name string,
	querier db.IQuerier,
	tracer apm.ITracer,
	convertRow func(DBData) (Data, error),
	toRow func(Data) (DBData, error),
	getFunc func(context.Context, db.ISharedQueriesReadOnly, ID) (DBData, error),
	getAllFunc func(context.Context, db.ISharedQueriesReadOnly) ([]DBData, error),
	createFunc func(context.Context, db.ISharedQueriesReadWrite, DBData) (sql.Result, error),
	updateFunc func(context.Context, db.ISharedQueriesReadWrite, DBData) (sql.Result, error),
	deleteFunc func(context.Context, db.ISharedQueriesReadWrite, ID) (sql.Result, error),
) *SharedRepo[ID, Data, DBData] {
	return &SharedRepo[ID, Data, DBData]{
		name:       name,
		querier:    querier,
		tracer:     tracer,
		convertRow: convertRow,
		toRow:      toRow,
		getFunc:    getFunc,
		getAllFunc: getAllFunc,
		createFunc: createFunc,
		updateFunc: updateFunc,
		deleteFunc: deleteFunc,
	}
}

// Get retrieves a single item by ID
func (r *SharedRepo[ID, Data, DBData]) Get(
	ctx context.Context,
	id ID,
) (Data, error) {
	var (
		result DBData
		data   Data
		err    error
	)
	r.tracer.Trace(ctx, r.name+".get", func(ctx context.Context, span apm.ISpan) error {
		err := r.querier.Shared(ctx, func(d db.ISharedQueriesReadOnly) error {
			result, err = r.getFunc(ctx, d, id)
			return err
		})
		if err != nil {
			return utils.ErrorOrNil("not found", utils.NewNotFoundError, err)
		}
		data, err = r.convertRow(result)
		return err
	})
	return data, err
}

// All retrieves all items
func (r *SharedRepo[ID, Data, DBData]) All(
	ctx context.Context,
) ([]Data, error) {
	var (
		results []DBData
		items   []Data
		err     error
	)
	r.tracer.Trace(ctx, r.name+".all", func(ctx context.Context, span apm.ISpan) error {
		err := r.querier.Shared(ctx, func(d db.ISharedQueriesReadOnly) error {
			var err error
			results, err = r.getAllFunc(ctx, d)
			return err
		})
		if err != nil {
			return utils.ErrorOrNil("not found", utils.NewNotFoundError, err)
		}
		items = make([]Data, len(results))
		for i, result := range results {
			item, err := r.convertRow(result)
			if err != nil {
				return utils.NewInternalError("failed to convert row", err)
			}
			items[i] = item
		}
		return nil
	})
	return items, err
}

// Create creates a new item
func (r *SharedRepo[ID, Data, DBData]) Create(
	ctx context.Context,
	data Data,
) error {
	var err error
	r.tracer.Trace(ctx, r.name+".create", func(ctx context.Context, span apm.ISpan) error {
		var row DBData
		row, err = r.toRow(data)
		if err != nil {
			return utils.NewInternalError("failed to convert to row", err)
		}

		var result sql.Result
		err = r.querier.SharedWrite(ctx, func(d db.ISharedQueriesReadWrite) error {
			result, err = r.createFunc(ctx, d, row)
			if err != nil {
				return err
			}
			if rowsAffected, _ := result.RowsAffected(); rowsAffected == 0 {
				return utils.NewInvalidStateError("item already exists")
			}
			return nil
		})
		return err
	})
	return utils.ErrorOrNil("item already exists", utils.NewInvalidStateError, err)
}

// Update updates an existing item
func (r *SharedRepo[ID, Data, DBData]) Update(
	ctx context.Context,
	data Data,
) error {
	var err error
	r.tracer.Trace(ctx, r.name+".update", func(ctx context.Context, span apm.ISpan) error {
		var row DBData
		row, err := r.toRow(data)
		if err != nil {
			return utils.NewInternalError("failed to convert to row", err)
		}

		var result sql.Result
		err = r.querier.SharedWrite(ctx, func(d db.ISharedQueriesReadWrite) error {
			result, err = r.updateFunc(ctx, d, row)
			if err != nil {
				return err
			}
			if rowsAffected, _ := result.RowsAffected(); rowsAffected == 0 {
				return utils.NewNotFoundError("not found")
			}
			return nil
		})
		return utils.ErrorOrNil("not found", utils.NewNotFoundError, err)
	})
	return err
}

// Delete deletes an item by ID
func (r *SharedRepo[ID, Data, DBData]) Delete(
	ctx context.Context,
	id ID,
) error {
	var err error
	r.tracer.Trace(ctx, r.name+".delete", func(ctx context.Context, span apm.ISpan) error {
		var result sql.Result
		err = r.querier.SharedWrite(ctx, func(d db.ISharedQueriesReadWrite) error {
			result, err = r.deleteFunc(ctx, d, id)
			if err != nil {
				return err
			}
			if rowsAffected, _ := result.RowsAffected(); rowsAffected == 0 {
				return utils.NewNotFoundError("not found")
			}
			return nil
		})
		return utils.ErrorOrNil("not found", utils.NewNotFoundError, err)
	})
	return err
}

// // Query performs a custom read query that returns multiple items
func (r *SharedRepo[ID, Data, DBData]) Query(
	ctx context.Context,
	queryFunc func(context.Context, db.ISharedQueriesReadOnly) ([]DBData, error),
) ([]Data, error) {
	var (
		results []DBData
		items   []Data
		err     error
	)
	r.tracer.Trace(ctx, r.name+".query", func(ctx context.Context, span apm.ISpan) error {
		err := r.querier.Shared(ctx, func(d db.ISharedQueriesReadOnly) error {
			var err error
			results, err = queryFunc(ctx, d)
			return err
		})
		if err != nil {
			return utils.NewInternalError("failed to query", err)
		}
		items = make([]Data, len(results))
		for i, result := range results {
			item, err := r.convertRow(result)
			if err != nil {
				return utils.NewInternalError("failed to convert row", err)
			}
			items[i] = item
		}
		return nil
	})
	return items, err
}

// // Execute performs a custom write operation
func (r *SharedRepo[ID, Data, DBData]) Execute(
	ctx context.Context,
	executeFunc func(context.Context, db.ISharedQueriesReadWrite) error,
) error {
	var err error
	r.tracer.Trace(ctx, r.name+".execute", func(ctx context.Context, span apm.ISpan) error {
		err := r.querier.SharedWrite(ctx, func(d db.ISharedQueriesReadWrite) error {
			err := executeFunc(ctx, d)
			return err
		})
		return utils.ErrorOrNil("operation failed", utils.NewInternalError, err)
	})
	return err
}
