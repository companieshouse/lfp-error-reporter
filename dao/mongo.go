package dao

import (
	"context"
	"fmt"
	"time"

	"github.com/companieshouse/chs.go/log"
	"github.com/companieshouse/lfp-error-reporter/config"
	"github.com/companieshouse/lfp-error-reporter/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Mongo provides a MongoDB implementation of the DAO
type Mongo struct {
	Config *config.Config
	Client *mongo.Client
	DAO
}

// New returns a new Mongo object using the provided config
func New(cfg *config.Config) DAO {
	return &Mongo{
		Config: cfg,
		Client: nil,
	}
}

// getMongoClient returns a mongo
func (m *Mongo) getMongoClient() error {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	var err error
	m.Client, err = mongo.Connect(ctx, options.Client().ApplyURI(m.Config.MongoDBURL))
	if err != nil {
		return err
	}
	return nil
}

// GetLFPData fetches lfp data
func (m *Mongo) GetLFPData(reconciliationMetaData *models.ReconciliationMetaData) (models.PenaltyList, error) {
	ctx := context.Background()
	var (
		penalties     []models.PayableResourceDao
		penaltiesData models.PenaltyList
	)

	err := m.getMongoClient()
	if err != nil {
		return penaltiesData, fmt.Errorf("error connecting to MongoDB: %s", err)
	}

	collection := m.Client.Database(m.Config.Database).Collection(m.Config.LFPCollection)

	log.Info("GetLFPData: reconciliationMetaData.StartTime (start time): " + reconciliationMetaData.StartTime.String())
	log.Info("GetLFPData: reconciliationMetaData.EndTime (end time): " + reconciliationMetaData.EndTime.String())

	filter := bson.M{"data.created_at": bson.M{
		"$gt": reconciliationMetaData.StartTime.AddDate(0, 0, -1),
		"$lt": reconciliationMetaData.EndTime.AddDate(0, 0, -1),
	}, "e5_command_error": bson.M{"$ne": ""}}
	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		return penaltiesData, fmt.Errorf("error retrieving lfp data: %s", err)
	}
	if err = cursor.All(ctx, &penalties); err != nil {
		return penaltiesData, fmt.Errorf("error storing lfp data: %s", err)
	}
	penaltiesData = models.PenaltyList{
		Penalties: penalties,
	}
	return penaltiesData, err
}
