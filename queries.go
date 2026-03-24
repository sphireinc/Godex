package godex

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"reflect"

	"github.com/jmoiron/sqlx"
)

// DefaultQueries stores specific queries to be used by Codex
type DefaultQueries struct {
	SelectById string `json:"selectById,omitempty"`
	SelectOne  string `json:"selectOne,omitempty"`
	Select     string `json:"select,omitempty"`
	Insert     string `json:"insert,omitempty"`
	Update     string `json:"update,omitempty"`
	Delete     string `json:"delete,omitempty"`
	SoftDelete string `json:"softDelete,omitempty"`
}

var (
	ErrDBNotConfigured     = errors.New("godex: database handle is not configured")
	ErrResultNotBound      = errors.New("godex: no result prototype is bound")
	ErrQueryNotFound       = errors.New("godex: query not found")
	ErrDestinationRequired = errors.New("godex: destination must be a non-nil pointer")
)

// Bind registers the result prototype used by the legacy Select* helpers.
func (c *Codex) Bind(result any) *Codex {
	c.result = result
	return c
}

// RegisterQuery stores a custom named query on the codex.
func (c *Codex) RegisterQuery(name, query string) {
	if c.Queries == nil {
		c.Queries = map[string]string{}
	}
	c.Queries[name] = query
}

// LookupQuery resolves a custom named query.
func (c *Codex) LookupQuery(name string) (string, error) {
	query, ok := c.Queries[name]
	if !ok {
		return "", fmt.Errorf("%w: %s", ErrQueryNotFound, name)
	}
	return query, nil
}

// RawQuery performs a query with given args
func (c *Codex) RawQuery(query string, args ...any) (*sql.Rows, error) {
	return c.RawQueryContext(context.Background(), query, args...)
}

// RawQueryContext performs a query with the supplied context.
func (c *Codex) RawQueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error) {
	if err := c.ensureDB(); err != nil {
		return nil, err
	}
	return c.DB.QueryContext(ctx, query, args...)
}

// QueryOneInto runs a query and scans one row into dest.
func (c *Codex) QueryOneInto(dest any, query string, arg any) error {
	return c.QueryOneContextInto(context.Background(), dest, query, arg)
}

// QueryOneContextInto runs a query and scans one row into dest with the supplied context.
func (c *Codex) QueryOneContextInto(ctx context.Context, dest any, query string, arg any) error {
	if err := c.ensureDB(); err != nil {
		return err
	}
	if err := validateDestination(dest); err != nil {
		return err
	}
	boundQuery, args, err := c.bindQuery(query, arg)
	if err != nil {
		return err
	}
	return c.DB.GetContext(ctx, dest, boundQuery, args...)
}

// QueryInto runs a query and scans all rows into dest.
func (c *Codex) QueryInto(dest any, query string, arg any) error {
	return c.QueryContextInto(context.Background(), dest, query, arg)
}

// QueryContextInto runs a query and scans all rows into dest with the supplied context.
func (c *Codex) QueryContextInto(ctx context.Context, dest any, query string, arg any) error {
	if err := c.ensureDB(); err != nil {
		return err
	}
	if err := validateDestination(dest); err != nil {
		return err
	}
	boundQuery, args, err := c.bindQuery(query, arg)
	if err != nil {
		return err
	}
	return c.DB.SelectContext(ctx, dest, boundQuery, args...)
}

// QueryNamedOneInto runs a registered custom query and scans one row into dest.
func (c *Codex) QueryNamedOneInto(name string, dest any, arg any) error {
	return c.QueryNamedOneContextInto(context.Background(), name, dest, arg)
}

// QueryNamedOneContextInto runs a registered custom query and scans one row into dest with the supplied context.
func (c *Codex) QueryNamedOneContextInto(ctx context.Context, name string, dest any, arg any) error {
	query, err := c.LookupQuery(name)
	if err != nil {
		return err
	}
	return c.QueryOneContextInto(ctx, dest, query, arg)
}

// QueryNamedInto runs a registered custom query and scans all rows into dest.
func (c *Codex) QueryNamedInto(name string, dest any, arg any) error {
	return c.QueryNamedContextInto(context.Background(), name, dest, arg)
}

// QueryNamedContextInto runs a registered custom query and scans all rows into dest with the supplied context.
func (c *Codex) QueryNamedContextInto(ctx context.Context, name string, dest any, arg any) error {
	query, err := c.LookupQuery(name)
	if err != nil {
		return err
	}
	return c.QueryContextInto(ctx, dest, query, arg)
}

// SelectByIDInto runs the default SelectById query and scans one row into dest.
func (c *Codex) SelectByIDInto(dest any, arg any) error {
	return c.SelectByIDContextInto(context.Background(), dest, arg)
}

// SelectByIDContextInto runs the default SelectById query and scans one row into dest with the supplied context.
func (c *Codex) SelectByIDContextInto(ctx context.Context, dest any, arg any) error {
	return c.QueryOneContextInto(ctx, dest, c.DefaultQueries.SelectById, arg)
}

// SelectOneInto runs the default SelectOne query and scans one row into dest.
func (c *Codex) SelectOneInto(dest any, arg any) error {
	return c.SelectOneContextInto(context.Background(), dest, arg)
}

// SelectOneContextInto runs the default SelectOne query and scans one row into dest with the supplied context.
func (c *Codex) SelectOneContextInto(ctx context.Context, dest any, arg any) error {
	return c.QueryOneContextInto(ctx, dest, c.DefaultQueries.SelectOne, arg)
}

// SelectInto runs the default Select query and scans all rows into dest.
func (c *Codex) SelectInto(dest any, arg any) error {
	return c.SelectContextInto(context.Background(), dest, arg)
}

// SelectContextInto runs the default Select query and scans all rows into dest with the supplied context.
func (c *Codex) SelectContextInto(ctx context.Context, dest any, arg any) error {
	return c.QueryContextInto(ctx, dest, c.DefaultQueries.Select, arg)
}

// InsertWith runs the default Insert query with the supplied named argument payload.
func (c *Codex) InsertWith(arg any) (int64, error) {
	return c.InsertWithContext(context.Background(), arg)
}

// InsertWithContext runs the default Insert query with the supplied context and named argument payload.
func (c *Codex) InsertWithContext(ctx context.Context, arg any) (int64, error) {
	result, err := c.execNamed(ctx, c.DefaultQueries.Insert, arg)
	if err != nil {
		return -1, err
	}
	return result.LastInsertId()
}

// UpdateWith runs the default Update query with the supplied named argument payload.
func (c *Codex) UpdateWith(arg any) (int64, error) {
	return c.UpdateWithContext(context.Background(), arg)
}

// UpdateWithContext runs the default Update query with the supplied context and named argument payload.
func (c *Codex) UpdateWithContext(ctx context.Context, arg any) (int64, error) {
	result, err := c.execNamed(ctx, c.DefaultQueries.Update, arg)
	if err != nil {
		return -1, err
	}
	return result.RowsAffected()
}

// DeleteWith runs the default Delete query with the supplied named argument payload.
func (c *Codex) DeleteWith(arg any) error {
	return c.DeleteWithContext(context.Background(), arg)
}

// DeleteWithContext runs the default Delete query with the supplied context and named argument payload.
func (c *Codex) DeleteWithContext(ctx context.Context, arg any) error {
	_, err := c.execNamed(ctx, c.DefaultQueries.Delete, arg)
	return err
}

// SoftDeleteWith runs the default SoftDelete query with the supplied named argument payload.
func (c *Codex) SoftDeleteWith(arg any) (int64, error) {
	return c.SoftDeleteWithContext(context.Background(), arg)
}

// SoftDeleteWithContext runs the default SoftDelete query with the supplied context and named argument payload.
func (c *Codex) SoftDeleteWithContext(ctx context.Context, arg any) (int64, error) {
	result, err := c.execNamed(ctx, c.DefaultQueries.SoftDelete, arg)
	if err != nil {
		return -1, err
	}
	return result.RowsAffected()
}

// SelectById one item.
//
// Deprecated: prefer SelectByIDInto or typed helpers from For[T].
func (c *Codex) SelectById(args ...any) (any, error) {
	dest, err := c.newLegacyResult()
	if err != nil {
		return nil, err
	}
	if err := c.SelectByIDInto(dest, normalizeArgs(args)); err != nil {
		return nil, err
	}
	return dest, nil
}

// SelectOne item.
//
// Deprecated: prefer SelectOneInto or typed helpers from For[T].
func (c *Codex) SelectOne(args ...any) (any, error) {
	dest, err := c.newLegacyResult()
	if err != nil {
		return nil, err
	}
	if err := c.SelectOneInto(dest, normalizeArgs(args)); err != nil {
		return nil, err
	}
	return dest, nil
}

// Select one.
//
// Deprecated: prefer SelectInto or typed helpers from For[T].
func (c *Codex) Select(args ...any) ([]any, error) {
	resultType, err := c.legacyResultType()
	if err != nil {
		return nil, err
	}
	slicePtr := reflect.New(reflect.SliceOf(resultType))
	if err := c.SelectInto(slicePtr.Interface(), normalizeArgs(args)); err != nil {
		return nil, err
	}
	values := slicePtr.Elem()
	res := make([]any, 0, values.Len())
	for i := 0; i < values.Len(); i++ {
		res = append(res, values.Index(i).Interface())
	}
	return res, nil
}

// Insert a new record.
//
// Deprecated: prefer InsertWith.
func (c *Codex) Insert(args ...any) (int64, error) {
	return c.InsertWith(normalizeLegacyExecArg(c, args))
}

// Update a specific record.
//
// Deprecated: prefer UpdateWith.
func (c *Codex) Update(args ...any) (int64, error) {
	return c.UpdateWith(normalizeLegacyExecArg(c, args))
}

// Delete a specific record.
//
// Deprecated: prefer DeleteWith.
func (c *Codex) Delete(args ...any) error {
	return c.DeleteWith(normalizeLegacyExecArg(c, args))
}

// SoftDelete a specific record.
//
// Deprecated: prefer SoftDeleteWith.
func (c *Codex) SoftDelete(args ...any) (int64, error) {
	return c.SoftDeleteWith(normalizeLegacyExecArg(c, args))
}

func (c *Codex) ensureDB() error {
	if c == nil || c.DB == nil {
		return ErrDBNotConfigured
	}
	return nil
}

func (c *Codex) execNamed(ctx context.Context, query string, arg any) (sql.Result, error) {
	if err := c.ensureDB(); err != nil {
		return nil, err
	}
	if arg == nil {
		return nil, errors.New("godex: named exec requires a payload")
	}
	return c.DB.NamedExecContext(ctx, query, arg)
}

func (c *Codex) bindQuery(query string, arg any) (string, []any, error) {
	if arg == nil {
		return query, nil, nil
	}
	if args, ok := arg.([]any); ok {
		return query, args, nil
	}
	if hasNamedParameters(query) {
		switch value := arg.(type) {
		case CxArgs:
			arg = map[string]any(value)
		}
		namedQuery, args, err := sqlx.Named(query, arg)
		if err != nil {
			return "", nil, err
		}
		return c.DB.Rebind(namedQuery), args, nil
	}
	return query, []any{arg}, nil
}

func hasNamedParameters(query string) bool {
	for i := 0; i < len(query); i++ {
		if query[i] == ':' && i+1 < len(query) {
			next := query[i+1]
			if next == '_' || ('a' <= next && next <= 'z') || ('A' <= next && next <= 'Z') {
				return true
			}
		}
	}
	return false
}

func validateDestination(dest any) error {
	if dest == nil {
		return ErrDestinationRequired
	}
	value := reflect.ValueOf(dest)
	if value.Kind() != reflect.Pointer || value.IsNil() {
		return ErrDestinationRequired
	}
	return nil
}

func normalizeArgs(args []any) any {
	switch len(args) {
	case 0:
		return nil
	case 1:
		return args[0]
	default:
		return args
	}
}

func normalizeLegacyExecArg(fallback any, args []any) any {
	if len(args) == 0 {
		return fallback
	}
	if len(args) == 1 {
		return args[0]
	}
	return args
}

func (c *Codex) legacyResultType() (reflect.Type, error) {
	if c.result == nil {
		return nil, ErrResultNotBound
	}
	resultType := reflect.TypeOf(c.result)
	if resultType.Kind() != reflect.Pointer {
		return nil, fmt.Errorf("%w: result prototype must be a pointer", ErrResultNotBound)
	}
	return resultType.Elem(), nil
}

func (c *Codex) newLegacyResult() (any, error) {
	resultType, err := c.legacyResultType()
	if err != nil {
		return nil, err
	}
	return reflect.New(resultType).Interface(), nil
}
