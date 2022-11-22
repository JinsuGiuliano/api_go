package store

import (
	"context"

	"github.com/rs/zerolog"
	zlog "github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var log zerolog.Logger // nolint

func init() {
	log = zlog.With().Strs("gtag", []string{"mongodb"}).Logger().Level(zerolog.InfoLevel) // nolint
}

type Store struct {
	client     *mongo.Client
	db         *mongo.Database
	ctx        context.Context
	connection string
}

func New() *Store {
	s := &Store{
		connection: "mongodb://localhost:27017",
	}
	log.Info().Msgf("start database connection: %s", s.connection)
	client, err := mongo.NewClient(options.Client().ApplyURI(s.connection))
	if err != nil {
		log.Fatal().Msgf("error trying to connect database: %s", err.Error())
	}

	ctx := context.Background()

	err = client.Connect(ctx)
	if err != nil {
		log.Fatal().Msgf("error trying to connect database: %s", err.Error())
	}

	s.ctx = ctx
	s.client = client
	s.db = client.Database("apiGo")
	return s
}

func (s *Store) CloseDB() {
	err := s.client.Disconnect(s.ctx)
	if err != nil {
		log.Fatal().Msgf("error trying to disconnect database: %s", err.Error())
	}
	log.Info().Msgf("closing database connection.")
}
