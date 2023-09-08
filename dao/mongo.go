package dao

import (
	"fmt"
	"strconv"

	"github.com/companieshouse/chs.go/log"
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

	dbNames, _ := session.DatabaseNames()

	log.Info("MK: getMongoSession: len(session.DatabaseNames())) => " + strconv.Itoa(len(dbNames)))

	if session == nil {
		log.Info("MK: getMongoSession: establishing db connection...")
		var err error
		session, err = mgo.Dial(cfg.MongoDBURL)
		if err != nil {
			log.Info("MK: getMongoSession: error establishing db connection => " + err.Error())
			return nil, err
		}
	}
	log.Info("MK: getMongoSession: returning session copy...")
	return session.Copy(), nil
}

// GetLFPData fetches lfp data
func (m *Mongo) GetLFPData(reconciliationMetaData *models.ReconciliationMetaData) (models.PenaltyList, error) {

	log.Info("MK: GetLFPData and fetching lfp data.")

	var penalties []models.PayableResourceDao

	var penaltiesData models.PenaltyList

	log.Info("MK: opening Mongo session.")

	mongoSession, err := getMongoSession(m.Config)
	if err != nil {
		log.Info("MK: error opening Mongo session => " + fmt.Errorf("error connecting to MongoDB: %s", err).Error())
		return penaltiesData, fmt.Errorf("error connecting to MongoDB: %s", err)
	}

	defer mongoSession.Close()
	log.Info("MK: closing Mongo session.")
	err = mongoSession.DB(m.Config.Database).C(m.Config.LFPCollection).Find(bson.M{"data.created_at": bson.M{
		"$gt": reconciliationMetaData.StartTime,
		"$lt": reconciliationMetaData.EndTime,
	}, "e5_command_error": bson.M{"$ne": ""}}).All(&penalties)
	if err != nil {
		log.Info("MK: error retrieving lfp dat => " + fmt.Errorf("error connecting to MongoDB: %s", err).Error())
		return penaltiesData, fmt.Errorf("error retrieving lfp data: %s", err)
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
