package dao

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/companieshouse/chs.go/log"
	"github.com/companieshouse/lfp-error-reporter/config"
	"github.com/companieshouse/lfp-error-reporter/models"
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

// GetPenaltyPaymentData fetches penalty payment data
func (m *Mongo) GetPenaltyPaymentData(reconciliationMetaData *models.ReconciliationMetaData) (models.PenaltyList, error) {
	ctx := context.Background()
	var (
		penalties       []models.PayableResourceDao
		penaltiesData   models.PenaltyList
		lambdaStartTime = reconciliationMetaData.StartTime.AddDate(0, 0, -1)
		lambdaEndTime   = reconciliationMetaData.EndTime.AddDate(0, 0, -1)
	)

	err := m.getMongoClient()
	if err != nil {
		return penaltiesData, fmt.Errorf("error connecting to MongoDB: %s", err)
	}

	collection := m.Client.Database(m.Config.Database).Collection(m.Config.PayableResourcesCollection)

	log.Info("GetPenaltyPaymentData: lambda start time: " + lambdaStartTime.String())
	log.Info("GetPenaltyPaymentData: lambda end time: " + lambdaEndTime.String())

	filter := bson.M{"data.created_at": bson.M{
		"$gt": lambdaStartTime,
		"$lt": lambdaEndTime,
	}, "e5_command_error": bson.M{"$ne": ""}}
	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		return penaltiesData, fmt.Errorf("error retrieving penalty payment data: %s", err)
	}
	if err = cursor.All(ctx, &penalties); err != nil {
		return penaltiesData, fmt.Errorf("error storing penalty payment data: %s", err)
	}
	penaltiesData = models.PenaltyList{
		Penalties: penalties,
	}
	return penaltiesData, err
}
