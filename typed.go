package godex

import "context"

// Model provides a typed facade over Codex read helpers.
type Model[T any] struct {
	codex *Codex
}

// For returns a typed facade for the provided Codex.
func For[T any](c *Codex) Model[T] {
	return Model[T]{codex: c}
}

// QueryOne loads one value from an arbitrary query.
func (m Model[T]) QueryOne(query string, arg any) (T, error) {
	return m.QueryOneContext(context.Background(), query, arg)
}

// QueryOneContext loads one value from an arbitrary query with context.
func (m Model[T]) QueryOneContext(ctx context.Context, query string, arg any) (T, error) {
	var value T
	err := m.codex.QueryOneContextInto(ctx, &value, query, arg)
	return value, err
}

// QueryMany loads many values from an arbitrary query.
func (m Model[T]) QueryMany(query string, arg any) ([]T, error) {
	return m.QueryManyContext(context.Background(), query, arg)
}

// QueryManyContext loads many values from an arbitrary query with context.
func (m Model[T]) QueryManyContext(ctx context.Context, query string, arg any) ([]T, error) {
	var values []T
	err := m.codex.QueryContextInto(ctx, &values, query, arg)
	return values, err
}

// QueryNamedOne loads one value from a registered named query.
func (m Model[T]) QueryNamedOne(name string, arg any) (T, error) {
	return m.QueryNamedOneContext(context.Background(), name, arg)
}

// QueryNamedOneContext loads one value from a registered named query with context.
func (m Model[T]) QueryNamedOneContext(ctx context.Context, name string, arg any) (T, error) {
	var value T
	err := m.codex.QueryNamedOneContextInto(ctx, name, &value, arg)
	return value, err
}

// QueryNamedMany loads many values from a registered named query.
func (m Model[T]) QueryNamedMany(name string, arg any) ([]T, error) {
	return m.QueryNamedManyContext(context.Background(), name, arg)
}

// QueryNamedManyContext loads many values from a registered named query with context.
func (m Model[T]) QueryNamedManyContext(ctx context.Context, name string, arg any) ([]T, error) {
	var values []T
	err := m.codex.QueryNamedContextInto(ctx, name, &values, arg)
	return values, err
}

// SelectByID loads one value using the default SelectById query.
func (m Model[T]) SelectByID(arg any) (T, error) {
	return m.SelectByIDContext(context.Background(), arg)
}

// SelectByIDContext loads one value using the default SelectById query with context.
func (m Model[T]) SelectByIDContext(ctx context.Context, arg any) (T, error) {
	var value T
	err := m.codex.SelectByIDContextInto(ctx, &value, arg)
	return value, err
}

// SelectOne loads one value using the default SelectOne query.
func (m Model[T]) SelectOne(arg any) (T, error) {
	return m.SelectOneContext(context.Background(), arg)
}

// SelectOneContext loads one value using the default SelectOne query with context.
func (m Model[T]) SelectOneContext(ctx context.Context, arg any) (T, error) {
	var value T
	err := m.codex.SelectOneContextInto(ctx, &value, arg)
	return value, err
}

// Select loads many values using the default Select query.
func (m Model[T]) Select(arg any) ([]T, error) {
	return m.SelectContext(context.Background(), arg)
}

// SelectContext loads many values using the default Select query with context.
func (m Model[T]) SelectContext(ctx context.Context, arg any) ([]T, error) {
	var values []T
	err := m.codex.SelectContextInto(ctx, &values, arg)
	return values, err
}
