package persist

import (
	"crawler/engine"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
	"gopkg.in/olivere/elastic.v5"
	"log"
)

func ItemSaver(index string) (chan engine.Item, error) {
	itemChan := make(chan engine.Item)
	client, err := elastic.NewClient(elastic.SetSniff(false))
	if err != nil {
		return nil, err
	}
	go func() {
		var count int
		for {
			item := <-itemChan
			log.Printf("Got Item#%d:%v", count, item)
			Save(item, client, index)
			count++
		}
	}()
	return itemChan, nil

}

func Save(item engine.Item, client *elastic.Client, index string) error {

	if item.Type == "" {
		return errors.New("must supply type")
	}
	indexService := client.Index().
		Index(index).
		Type(item.Type).
		BodyJson(item)
	if item.ID != "" {
		indexService.Id(item.ID)
	}

	_, err := indexService.
		Do(context.Background())
	if err != nil {
		return err
	}
	return nil

}
