package client

import (
	"crawler/engine"
	"crawler_distributed/config"
	"crawler_distributed/rpcsupport"
	"log"
)


func ItemSaver(host string) (chan engine.Item, error) {
	client, err := rpcsupport.NewClient(host)
	if err != nil {
		return nil, err
	}
	itemChan := make(chan engine.Item)
	go func() {
		var count int
		for {
			item := <-itemChan
			log.Printf("Got Item#%d:%v", count, item)
			count++
			result := ""
			err := client.Call(config.ItemSaverRpc, item, &result)
			if err != nil {
				log.Printf("Item Saver: error saving item %v: %v", item, err)
			}

		}
	}()
	return itemChan, nil

}
