package comm

import (
	"encoding/json"
	"fmt"

	"github.com/go-redis/redis"
	"github.com/smilga/analyzer/api"
)

type list string

type userList string

const (
	PendingList        userList = "pending:websites:user:"
	TimeoutedList      userList = "timeouted:websites:user:"
	ListsList          list     = "pending:lists"
	TimeoutedListsList list     = "timeouted:lists"
	ResultsList        list     = "inspect:results"
	PatternsList       list     = "inspect:patterns"
)

type Comm struct {
	Client *redis.Client
}

func (c *Comm) CollectResults() (*api.Result, error) {
	ss, err := c.Client.RPop(string(ResultsList)).Result()
	if err == redis.Nil {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}
	var result = &api.Result{}
	err = json.Unmarshal([]byte(ss), result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (c *Comm) Inspect(w *api.Website) error {
	ws, err := json.Marshal(w)
	if err != nil {
		return err
	}

	list := fmt.Sprintf("%s%v", PendingList, w.UserID)

	_, err = c.Client.SAdd(string(ListsList), list).Result()
	if err != nil {
		return err
	}

	_, err = c.Client.LPush(list, string(ws)).Result()
	if err != nil {
		return err
	}

	return nil
}

func NewComm() *Comm {
	return &Comm{
		Client: redis.NewClient(&redis.Options{
			Addr: "redis:6379",
		}),
	}
}
