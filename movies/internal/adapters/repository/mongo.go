package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/caiqueborghese/sipubtech-challenge/movies/internal/domain"
	"github.com/caiqueborghese/sipubtech-challenge/movies/internal/ports"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var _ ports.MovieRepository = (*MongoRepository)(nil)

type MongoRepository struct {
	col *mongo.Collection
}

func NewMongoRepository(col *mongo.Collection) (*MongoRepository, error) {

	_, err := col.Indexes().CreateMany(context.Background(), []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "title", Value: 1}, {Key: "year", Value: 1}},
			Options: options.Index().SetUnique(true).SetName("uniq_title_year"),
		},
		{
			Keys:    bson.D{{Key: "legacy_id", Value: 1}},
			Options: options.Index().SetUnique(true).SetSparse(true).SetName("uniq_legacy_id"),
		},
	})
	if err != nil {
		return nil, fmt.Errorf("create indexes: %w", err)
	}
	return &MongoRepository{col: col}, nil
}

func (r *MongoRepository) List(ctx context.Context) ([]domain.Movie, error) {
	findOpts := options.Find().
		SetSort(bson.D{
			{Key: "legacy_id", Value: 1}, // 1 = asc
			{Key: "title", Value: 1},     // desempate
		}).
		SetCollation(&options.Collation{
			Locale:          "en",
			NumericOrdering: true, // ordenação numérica para strings "8", "10", ...
		})

	cur, err := r.col.Find(ctx, bson.D{}, findOpts)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	var out []domain.Movie
	for cur.Next(ctx) {
		var dbm dbMovie
		if err := cur.Decode(&dbm); err != nil {
			return nil, err
		}
		out = append(out, dbm.toDomain())
	}
	return out, cur.Err()
}

func (r *MongoRepository) Get(ctx context.Context, id string) (*domain.Movie, error) {

	if oid, err := primitive.ObjectIDFromHex(id); err == nil {
		var dbm dbMovie
		if err := r.col.FindOne(ctx, bson.M{"_id": oid}).Decode(&dbm); err != nil {
			if errors.Is(err, mongo.ErrNoDocuments) {
				return nil, domain.ErrNotFound
			}
			return nil, err
		}
		dm := dbm.toDomain()
		return &dm, nil
	}

	// fallback: legacy_id
	var dbm dbMovie
	if err := r.col.FindOne(ctx, bson.M{"legacy_id": id}).Decode(&dbm); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}
	dm := dbm.toDomain()
	return &dm, nil
}

func (r *MongoRepository) Create(ctx context.Context, m *domain.Movie) (*domain.Movie, error) {
	dbm := fromDomain(*m)
	res, err := r.col.InsertOne(ctx, dbm)
	if err != nil {

		var we mongo.WriteException
		if errors.As(err, &we) {
			for _, e := range we.WriteErrors {
				if e.Code == 11000 {
					return nil, fmt.Errorf("duplicate movie (title/year or legacy_id): %w", err)
				}
			}
		}
		return nil, err
	}

	cp := *m
	if cp.LegacyID != "" {
		cp.ID = cp.LegacyID
	} else {
		if oid, ok := res.InsertedID.(primitive.ObjectID); ok {
			cp.ID = oid.Hex()
		}
	}
	return &cp, nil
}

func (r *MongoRepository) Delete(ctx context.Context, id string) error {
	if oid, err := primitive.ObjectIDFromHex(id); err == nil {
		res, err := r.col.DeleteOne(ctx, bson.M{"_id": oid})
		if err != nil {
			return err
		}
		if res.DeletedCount == 0 {
			return domain.ErrNotFound
		}
		return nil
	}

	// fallback: por legacy_id
	res, err := r.col.DeleteOne(ctx, bson.M{"legacy_id": id})
	if err != nil {
		return err
	}
	if res.DeletedCount == 0 {
		return domain.ErrNotFound
	}
	return nil
}

func (r *MongoRepository) Count(ctx context.Context) (int64, error) {
	return r.col.CountDocuments(ctx, bson.D{})
}

func (r *MongoRepository) BulkInsertIgnoreDuplicates(ctx context.Context, ms []domain.Movie) (int, error) {
	if len(ms) == 0 {
		return 0, nil
	}
	models := make([]mongo.WriteModel, 0, len(ms))
	for _, m := range ms {
		models = append(models, mongo.NewInsertOneModel().SetDocument(fromDomain(m)))
	}
	res, err := r.col.BulkWrite(ctx, models, options.BulkWrite().SetOrdered(false))
	if err != nil {
		var bwe mongo.BulkWriteException
		if errors.As(err, &bwe) {

			for _, we := range bwe.WriteErrors {
				if we.Code != 11000 {

					return int(res.InsertedCount), err
				}
			}
			return int(res.InsertedCount), nil
		}
		return int(res.InsertedCount), err
	}
	return int(res.InsertedCount), nil
}

/************** mapeamentos **************/

type dbMovie struct {
	ID       primitive.ObjectID `bson:"_id,omitempty"`
	Title    string             `bson:"title"`
	Year     int                `bson:"year"`
	LegacyID string             `bson:"legacy_id,omitempty"`
	Created  time.Time          `bson:"created_at,omitempty"`
}

func (d dbMovie) toDomain() domain.Movie {
	id := d.LegacyID
	if id == "" {
		id = d.ID.Hex()
	}
	return domain.Movie{
		ID:    id,
		Title: d.Title,
		Year:  d.Year,
	}
}

func fromDomain(m domain.Movie) dbMovie {
	return dbMovie{
		Title:    m.Title,
		Year:     m.Year,
		LegacyID: m.LegacyID,
		Created:  time.Now().UTC(),
	}
}
