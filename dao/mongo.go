package dao

import (
	"fmt"

	"github.com/companieshouse/lfp-error-reporter/config"
	"github.com/companieshouse/lfp-error-reporter/models"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

// Mongo provides a MongoDB implementation of the DAO
type Mongo struct {
	Config *config.Config
}

// New returns a new Mongo struct using the provided config
func New(cfg *config.Config) *Mongo {

	return &Mongo{
		Config: cfg,
	}
}

var session *mgo.Session

// getMongoSession retrieves a fresh MongoDB session
func getMongoSession(cfg *config.Config) (*mgo.Session, error) {

	if session == nil {

		var err error
		session, err = mgo.Dial(cfg.MongoDBURL)
		if err != nil {
			return nil, err
		}
	}

	return session.Copy(), nil
}

// GetLFPData fetches lfp data
func (m *Mongo) GetLFPData(reconciliationMetaData *models.ReconciliationMetaData) (models.PenaltyList, error) {

	var penalties []models.PayableResourceDao

	var penaltiesData models.PenaltyList

	mongoSession, err := getMongoSession(m.Config)
	if err != nil {
		return penaltiesData, fmt.Errorf("error connecting to MongoDB: %s", err)
	}
	defer mongoSession.Close()

	err = mongoSession.DB(m.Config.Database).C(m.Config.LFPCollection).Find(bson.M{"data.created_at": bson.M{
		"$gt": reconciliationMetaData.StartTime,
		"$lt": reconciliationMetaData.EndTime,
	}, "e5_command_error": bson.M{"$ne": ""}}).All(&penalties)
	if err != nil {
		return penaltiesData, fmt.Errorf("error retrieving lfp data: %s", err)
	}

	penaltiesData = models.PenaltyList{
		Penalties: penalties,
	}

	return penaltiesData, err
}
