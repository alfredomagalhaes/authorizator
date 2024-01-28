package repository

import (
	"context"
	"errors"
	"fmt"
	"log"
	"reflect"
	"strings"
	"time"

	"github.com/alfredomagalhaes/authorizator/types"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsoncodec"
	"go.mongodb.org/mongo-driver/bson/bsonrw"
	"go.mongodb.org/mongo-driver/bson/bsontype"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var appCollectionName string = "applications"
var ErrNoRecordsFound = errors.New("no records found with given parameters")
var (
	tUUID       = reflect.TypeOf(uuid.UUID{})
	uuidSubtype = byte(0x04)

	mongoRegistry = bson.NewRegistry()
)

// MongoRepository base struct from a mongo repository
// contains the client and database config after a
// successful connection
type MongoRepository struct {
	client *mongo.Client
	db     *mongo.Database
	log    *zerolog.Logger
}

// MongoRepositoryConnConfig struct to config the connection to mongoDB
type MongoRepositoryConnConfig struct {
	Username     string
	Password     string
	Host         string
	Port         string
	DatabaseName string
}

// NewMongoRepository creates a new connections with MongoDB database and returns a new
// MongoRepository
func NewMongoRepository(mrc MongoRepositoryConnConfig, ctx context.Context, log *zerolog.Logger) *MongoRepository {

	var mongoRepo MongoRepository
	var err error
	mongoCredentials := options.Credential{
		AuthMechanism: "SCRAM-SHA-256",
		Username:      mrc.Username,
		Password:      mrc.Password,
	}

	mongoConnUrl := fmt.Sprintf("mongodb://%s:%s/",
		mrc.Host,
		mrc.Port,
	)
	mongoRegistry.RegisterTypeEncoder(tUUID, bsoncodec.ValueEncoderFunc(uuidEncodeValue))
	mongoRegistry.RegisterTypeDecoder(tUUID, bsoncodec.ValueDecoderFunc(uuidDecodeValue))
	clientOptions := options.Client().ApplyURI(mongoConnUrl).SetAuth(mongoCredentials).SetRegistry(mongoRegistry)
	mongoRepo.client, err = mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal().Err(err)
	}

	// Test the connection with MongoDB
	err = mongoRepo.client.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Fatal().Err(err)
	}

	mongoRepo.db = mongoRepo.client.Database(mrc.DatabaseName)
	mongoRepo.log = log

	return &mongoRepo

}

// CloseConn executes "Disconnect" from mongodb client
// to not let unused connections open
func (mr *MongoRepository) CloseConn(ctx context.Context) error {
	return mr.client.Disconnect(ctx)
}

// CreateIndexes create indexes on collections
func (mr *MongoRepository) CreateIndexes() {
	//Value attribute on bson.D struct
	//defines the order of the index
	//1 - ascending
	//-1 - descending
	indexModel := mongo.IndexModel{
		Keys: bson.D{
			{Key: "external_id", Value: 1}},
		Options: options.Index().SetUnique(true),
	}
	coll := mr.db.Collection(appCollectionName)
	name, err := coll.Indexes().CreateOne(context.TODO(), indexModel)
	if err != nil {
		log.Printf("Failed to create index: %v on collection %s\n", err, appCollectionName)
	}
	log.Printf("Index %s created at collection %s", name, appCollectionName)

}

// GetApplications get all valid applications from the database
// deleted applications are ignored
func (mr *MongoRepository) GetApplications(useCache bool) ([]types.Application, error) {
	return nil, nil
}

// GetApplicationsFromCache check if the applications are in cache server
// and return valid items.
func (mr *MongoRepository) GetApplicationsFromCache() ([]types.Application, error) {
	return nil, nil
}

// GetApplication search for a single application with the given ID,
// deleted applications should not return
func (mr *MongoRepository) GetApplication(id uuid.UUID) (types.Application, error) {
	filter := bson.M{
		"_id":    id,
		"active": true,
	}
	coll := mr.db.Collection(appCollectionName)
	mongoResult := coll.FindOne(context.TODO(), filter)

	var result types.Application

	err := mongoResult.Decode(&result)

	if err != nil && err == mongo.ErrNoDocuments {
		return result, ErrNoRecordsFound
	} else if err != nil {
		mr.log.Error().Err(err).Msg("error while trying to get applications")
		return result, ErrNoRecordsFound
	}
	return result, nil
}

func (mr *MongoRepository) GetApplicationFromCache(id uuid.UUID) (types.Application, error) {
	return types.Application{}, nil
}

// SaveApplication save a new application in the database
func (mr *MongoRepository) SaveApplication(app types.Application) (uuid.UUID, error) {
	ctx := context.Background()
	coll := mr.db.Collection(appCollectionName)

	app.Created_At = time.Now()
	app.Updated_At = app.Created_At
	app.Active = true
	app.ID = uuid.New()

	result, err := coll.InsertOne(ctx, app)

	if err != nil {
		errString := err.Error()
		if strings.Contains(errString, "duplicate") {
			err = errors.New("application already exists, try another `external_id`")
		} else {
			err = errors.New("error while trying to create an application, try again later")
			mr.log.Error().Err(err).Msg(err.Error())
		}

		return uuid.Nil, err //errors.New("could not create the application")
	}

	oId, ok := result.InsertedID.(primitive.Binary)

	if !ok {
		mr.log.Error().Err(err).Msg("error while trying to parse create id")
		return uuid.Nil, errors.New("could not get inserted ID")
	}

	parsedID, err := uuid.FromBytes(oId.Data)

	if err != nil {
		mr.log.Error().Err(err).Msg("error while trying to parse create id")
		return uuid.Nil, errors.New("could not get inserted ID")
	}
	return parsedID, nil
}

func uuidEncodeValue(ec bsoncodec.EncodeContext, vw bsonrw.ValueWriter, val reflect.Value) error {
	if !val.IsValid() || val.Type() != tUUID {
		return bsoncodec.ValueEncoderError{Name: "uuidEncodeValue", Types: []reflect.Type{tUUID}, Received: val}
	}
	b := val.Interface().(uuid.UUID)
	return vw.WriteBinaryWithSubtype(b[:], uuidSubtype)
}

func uuidDecodeValue(dc bsoncodec.DecodeContext, vr bsonrw.ValueReader, val reflect.Value) error {
	if !val.CanSet() || val.Type() != tUUID {
		return bsoncodec.ValueDecoderError{Name: "uuidDecodeValue", Types: []reflect.Type{tUUID}, Received: val}
	}

	var data []byte
	var subtype byte
	var err error
	switch vrType := vr.Type(); vrType {
	case bsontype.Binary:
		data, subtype, err = vr.ReadBinary()
		if subtype != uuidSubtype {
			return fmt.Errorf("unsupported binary subtype %v for UUID", subtype)
		}
	case bsontype.Null:
		err = vr.ReadNull()
	case bsontype.Undefined:
		err = vr.ReadUndefined()
	default:
		return fmt.Errorf("cannot decode %v into a UUID", vrType)
	}

	if err != nil {
		return err
	}
	uuid2, err := uuid.FromBytes(data)
	if err != nil {
		return err
	}
	val.Set(reflect.ValueOf(uuid2))
	return nil
}
