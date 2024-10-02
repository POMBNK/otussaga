package event

import (
	"context"
	sq "github.com/Masterminds/squirrel"
	"github.com/POMBNK/orderservice/internal/entity"
	"github.com/POMBNK/orderservice/pkg/client/postgres"
)

type Repo struct {
	conn postgres.Client
}

func NewRepo(conn postgres.Client) *Repo {
	return &Repo{conn: conn}
}

func (r *Repo) InsertEvent(ctx context.Context, event entity.Event) (int, error) {

	insertBuilder := sq.Insert("events").PlaceholderFormat(sq.Dollar).
		Columns("event_type", "payload").
		Values(event.Info.Type, event.Payload).
		Suffix("RETURNING id")

	query, args, err := insertBuilder.ToSql()
	if err != nil {
		return 0, err
	}

	var eventID int
	err = r.conn.QueryRow(ctx, query, args...).Scan(&eventID)
	if err != nil {
		return 0, err
	}

	return eventID, nil
}

func (r *Repo) GetNewEvents(ctx context.Context) ([]entity.Event, error) {

	selectBuilder := sq.Select("id", "event_type", "payload").
		PlaceholderFormat(sq.Dollar).
		From("events").
		Where(sq.Eq{"status": "NEW"}).
		Limit(100)

	query, args, err := selectBuilder.ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := r.conn.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []entity.Event
	for rows.Next() {
		var event entity.Event
		if err := rows.Scan(&event.Info.ID, &event.Info.Type, &event.Payload); err != nil {
			return nil, err
		}
		events = append(events, event)
	}

	return events, nil
}

func (r *Repo) GetNewEventIDs(ctx context.Context) ([]int, error) {

	selectBuilder := sq.Select("id").
		PlaceholderFormat(sq.Dollar).
		From("events").
		Where(sq.Eq{"status": "NEW"}).
		Limit(100)

	query, args, err := selectBuilder.ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := r.conn.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ids []int
	for rows.Next() {
		var id int
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}

	return ids, nil
}

func (r *Repo) SetCompleted(ctx context.Context, ids []int) error {
	updBuilder := sq.Update("events").PlaceholderFormat(sq.Dollar).
		Set("status", "COMPLETED").
		Where(sq.Eq{"id": ids}).
		Suffix("RETURNING id")

	query, args, err := updBuilder.ToSql()
	if err != nil {
		return err
	}

	_, err = r.conn.Exec(ctx, query, args...)
	if err != nil {
		return err
	}

	return err
}

func (r *Repo) SetFailed(ctx context.Context, ids []int) error {
	updBuilder := sq.Update("events").PlaceholderFormat(sq.Dollar).
		Set("status", "FAILED").
		Where(sq.Eq{"id": ids}).
		Suffix("RETURNING id")

	query, args, err := updBuilder.ToSql()
	if err != nil {
		return err
	}

	_, err = r.conn.Exec(ctx, query, args...)
	if err != nil {
		return err
	}

	return err
}

func (r *Repo) GetNewEvent(ctx context.Context) (entity.Event, error) {

	selectBuilder := sq.Select("id", "event_type", "payload").
		PlaceholderFormat(sq.Dollar).
		From("events").
		Where(sq.Eq{"status": "NEW"}).
		Limit(1)

	query, args, err := selectBuilder.ToSql()
	if err != nil {
		return entity.Event{}, err
	}

	var event entity.Event
	err = r.conn.QueryRow(ctx, query, args...).Scan(&event.Info.ID, &event.Info.Type, &event.Payload)
	if err != nil {
		return entity.Event{}, err
	}

	return event, nil
}
