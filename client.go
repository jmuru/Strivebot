package client

import (
	"fmt"
	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
	"github.com/strivebot/cache"
	"log"
	"os"
	"os/signal"
	"syscall"
)

type Config struct {
	ConsumerKey string `yaml:"consumerKey"`
	ConsumerSecret string `yaml:"consumerSecret"`
	AccessSecret string `yaml:"accessSecret"`
	AccessToken string `yaml:"accessToken"`
	BearerToken string `yaml:"bearerToken"`
}

type TwitterClient struct {
	client *twitter.Client
	tCache *cache.Cache
}

func NewTwitterClient(c *Config, tca *cache.Cache) *TwitterClient {
	consumerKey := c.ConsumerKey
	consumerSecret := c.ConsumerSecret
	accessToken := c.AccessToken
	accessSecret := c.AccessSecret

	if consumerKey == "" || consumerSecret == "" || accessToken == "" || accessSecret == "" {
		log.Fatal("Consumer key/secret and Access token/secret required")
	}

	config := oauth1.NewConfig(consumerKey, consumerSecret)
	token := oauth1.NewToken(accessToken, accessSecret)
	// OAuth1 http.Client will automatically authorize Requests
	httpClient := config.Client(oauth1.NoContext, token)

	// Twitter Client
	tc := twitter.NewClient(httpClient)
	return &TwitterClient{
		client: tc,
		tCache: tca,
	}
}

type TweetMediaInfo struct {
	Url string
	Type string
	Duration int
	MediaUrl string
}

type TweetHashTagInfo struct {
	Hashtag string
}

func (tc *TwitterClient) HelloWorld() {
	fmt.Println("Hello world")
}

func (tc *TwitterClient) StartStream() {
	filterParams := &twitter.StreamFilterParams{
		Track:         []string{"#strivebot1"},
		StallWarnings: twitter.Bool(true),
	}
	// Convenience Demux demultiplexed stream messages
	demux := twitter.NewSwitchDemux()
	demux.Tweet = func(tweet *twitter.Tweet) {
		fmt.Println(tweet.Text)
		fmt.Println(tweet.QuotedStatusIDStr)
		if !tweet.PossiblySensitive {
			err := tc.tCache.StoreTweetFromStream(tweet.QuotedStatusIDStr); if err != nil {
				fmt.Errorf("unable to store tweet")
			}
		}
	}
	demux.DM = func(dm *twitter.DirectMessage) {
		fmt.Println(dm.SenderID)
	}
	demux.Event = func(event *twitter.Event) {
		fmt.Printf("%#v\n", event)
	}

	stream, err := tc.client.Streams.Filter(filterParams)
	if err != nil {
		log.Fatal(err)
	}

	go demux.HandleChan(stream.Messages)

	// Wait for SIGINT and SIGTERM (HIT CTRL-C)
	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	log.Println(<-ch)

	fmt.Println("Stopping Stream...")
	stream.Stop()
}

func (tc *TwitterClient) getTweetInfo() string {
	return ""
}


