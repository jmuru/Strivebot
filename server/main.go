package main

import (
	client "github.com/strivebot"
	cache "github.com/strivebot/cache"
	yaml "gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
)

type config struct {
	ConsumerKey string `yaml:"consumerKey"`
	ConsumerSecret string `yaml:"consumerSecret"`
	AccessSecret string `yaml:"accessSecret"`
	AccessToken string `yaml:"accessToken"`
	BearerToken string `yaml:"bearerToken"`
}

func (c *config) getConf() *config {
	yamlFile, err := ioutil.ReadFile("./server/conf.yml")
	if err != nil {
		log.Printf("yamlFile.Get err   #%v ", err)
	}
	err = yaml.Unmarshal(yamlFile, c)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}

	return c
}

func main() {
	var conf config
	conf.getConf()

	tCache := cache.NewCache()

	tClient := client.NewTwitterClient(&client.Config{
		ConsumerKey: conf.ConsumerKey,
		ConsumerSecret: conf.ConsumerSecret,
		AccessToken: conf.AccessToken,
		AccessSecret: conf.AccessSecret,
		BearerToken: conf.BearerToken,
	}, tCache)
	tClient.StartStream()

}