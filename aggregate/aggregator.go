package aggregate

import (
	"context"
	"errors"
	"fmt"
	"github.com/dghubble/go-twitter/twitter"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)


var (
	DatabaseName = ""
)

type TweetMediaInfo struct {
	Url string
	Type string
	Duration int
	MediaUrl string
}

type TweetHashTagInfo struct {
	Hashtag string
}

type DBClient struct {
	db           *mongo.Client
	databaseName string
	kctx         context.Context
}

// assumption: source video will always be first medial element
func parseTweetMedia(tweet *twitter.Tweet) (*TweetMediaInfo, error) {
	var ti []*TweetMediaInfo
	for _, media := range tweet.Entities.Media {
		if media.Type == "video" {
			ti = append(ti, &TweetMediaInfo{
				Url: media.URL,
				Type:media.Type,
				Duration: media.VideoInfo.DurationMillis,
				MediaUrl:media.MediaURL,
			})
		}
	}

	if len(ti) == 0 {
		return nil, fmt.Errorf("unable tweet info does not include a video result list length == %d", len(ti))
	}
	return ti[0], nil
}

func parseTweetHashTags(tweet * twitter.Tweet) []*TweetHashTagInfo {
	var th []*TweetHashTagInfo
	for _, ht := range tweet.Entities.Hashtags{
		th = append(th, &TweetHashTagInfo{ht.Text})
	}
	return th
}

func NewDBClient(DatabaseURI string) (*DBClient, error) {
	dbClient, err := mongo.NewClient(options.Client().ApplyURI(DatabaseURI))
	if err != nil {
		return nil, errors.New("unable to create db client")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Minute)
	defer cancel()
	if err := dbClient.Connect(ctx); err != nil {
		panic("unable to connect to db")
	} else {
		fmt.Println("connected to db")
	}
	return &DBClient{
		db:           dbClient,
		databaseName: "cryptokat",
		kctx:         ctx,
	}, nil
}

// SAVE TO DB EXAMPLE FROM KAT PROJECT

//func (k KatClient) CreateAccessoryData(a *AccessoryData) error {
//	katDB := k.db.Database(k.databaseName)
//	accessoryCollection := katDB.Collection("accessory")
//	result, err := accessoryCollection.InsertOne(context.TODO(), a)
//	if err != nil {
//		errors.New(fmt.Sprintf("unable to insert in accessory collection, error: %v\n", err))
//	}
//	fmt.Printf(fmt.Sprintf("Inserted: %v into db\n", result))
//	return nil
//}

//func (k KatClient) GetAccessoryData(placement string) (*AccessoryResponse, error) {
//	var accessories []AccessoryData
//	fmt.Printf("placement in client %s\n", placement)
//	nctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
//	defer cancel()
//	katDB := k.db.Database(k.databaseName)
//	accessoryCollection := katDB.Collection("accessory")
//
//	cursor, err := accessoryCollection.Find(nctx, bson.D{{"placement", placement}})
//	defer cursor.Close(nctx)
//	if err != nil {
//		log.Fatalf("cursor err %v", err)
//		return nil, errors.New(fmt.Sprintf("error: %v", err))
//	} else {
//		for cursor.Next(nctx) {
//			var result AccessoryData
//			if err := cursor.Decode(&result); err != nil {
//				fmt.Println("cursor.Next() error:", err)
//				os.Exit(1)
//			} else {
//				accessories = append(accessories, result)
//			}
//		}
//	}