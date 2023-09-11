package dao

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/companieshouse/chs.go/log"
	"github.com/companieshouse/lfp-error-reporter/config"
	"github.com/companieshouse/lfp-error-reporter/models"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gopkg.in/mgo.v2/bson"
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

	log.Info("MK: getMongoClient getting mongo session.")

	if m.Client == nil {
		log.Info("MK: getMongoClient: establishing db connection...")
		var err error
		m.Client, err = mongo.Connect(ctx, options.Client().ApplyURI(m.Config.MongoDBURL))
		if err != nil {
			log.Info("MK: getMongoClient: error establishing db connection => " + err.Error())
			return err
		}
	}
	log.Info("MK: getMongoClient: returning session copy...")
	return nil
}

// GetLFPData fetches lfp data
func (m *Mongo) GetLFPData(reconciliationMetaData *models.ReconciliationMetaData) (models.PenaltyList, error) {

	log.Info("MK: GetLFPData and fetching lfp data.")
	ctx := context.Background()
	var (
		penalties     []models.PayableResourceDao
		penaltiesData models.PenaltyList
	)

	log.Info("MK: opening Mongo session.")

	err := m.getMongoClient()
	if err != nil {
		log.Info("MK: error opening Mongo session")
		return penaltiesData, fmt.Errorf("error connecting to MongoDB: %s", err)
	}

	log.Info("MK: closing Mongo session.")
	collection := m.Client.Database(m.Config.Database).Collection(m.Config.LFPCollection)

	dataItem, err := collection.FindOne(context.TODO(), bson.D{}).DecodeBytes()
	if err != nil {
		log.Info("MK: GetLFPData: m.Client.ListDatabaseNames error returned => " + err.Error())
	}
	if dataItem != nil {
		log.Info("MK: dataItem => : " + string(dataItem))
	}

	filter := bson.M{"data.created_at": bson.M{
		"$gt": reconciliationMetaData.StartTime,
		"$lt": reconciliationMetaData.EndTime,
	}, "e5_command_error": bson.M{"$ne": ""}}
	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		log.Info("MK: error retrieving lfp data")
		return penaltiesData, fmt.Errorf("error retrieving lfp data: %s", err)
	}
	if err = cursor.All(ctx, &penalties); err != nil {
		log.Info("MK: error putting lfp data in array")
		return penaltiesData, fmt.Errorf("error storing lfp data: %s", err)
	}
	penaltiesData = models.PenaltyList{
		Penalties: penalties,
	}

	if len(penaltiesData.Penalties) > 0 {
		log.Info("MK: GetLFPData: len(penalties) => " + strconv.Itoa(len(penalties)))
		log.Info("MK: GetLFPData: penaltiesData.Penalties[0].CompanyNumber => " + penaltiesData.Penalties[0].CompanyNumber)
		log.Info("MK: GetLFPData: penaltiesData.Penalties[0].E5CommandError => " + penaltiesData.Penalties[0].E5CommandError)
		log.Info("MK: GetLFPData: penaltiesData.Penalties[0].Reference => " + penaltiesData.Penalties[0].Reference)
		log.Info("MK: GetLFPData: penaltiesData.Penalties[0]..Data.Payment.Amount => " + penaltiesData.Penalties[0].Data.Payment.Amount)
	}

	return penaltiesData, err
}
