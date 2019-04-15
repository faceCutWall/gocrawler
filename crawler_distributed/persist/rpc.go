package persist

import (
	"crawler/engine"
	"crawler/persist"
	"gopkg.in/olivere/elastic.v5"
	"log"
)

type ItemSaverService struct {
	Client *elastic.Client
	Index  string
}

func (p *ItemSaverService) Save(item engine.Item, result *string) error {
	err:= persist.Save(item, p.Client, p.Index)
	log.Printf("Item %v saved.",item)
	if err==nil{
		*result = "ok"
	}else {
		log.Printf("Error saving item %v: %v",item,err)
	}
	return err
}
