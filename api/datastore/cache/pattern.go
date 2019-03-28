package cache

import (
	"encoding/json"
	"strconv"

	"github.com/go-redis/redis"
	"github.com/smilga/analyzer/api"
)

const redisKey = "inspect:patterns"

// PatternCache wraps storage and caches
type PatternCache struct {
	store  api.PatternStorage
	client *redis.Client
}

// Stores all patterns in redis to be accessible to other docker services
func (c *PatternCache) Save(p *api.Pattern) error {
	err := c.store.Save(p)
	if err != nil {
		return err
	}

	bs, err := json.Marshal(p)
	if err != nil {
		return err
	}

	_, err = c.client.HSet(redisKey, strconv.Itoa(int(p.ID)), string(bs)).Result()
	if err != nil {
		return err
	}

	return nil
}

func (c *PatternCache) Delete(id api.PatternID) error {
	err := c.store.Delete(id)
	if err != nil {
		return err
	}

	_, err = c.client.HDel(redisKey, strconv.Itoa(int(id))).Result()
	if err != nil {
		return err
	}
	return nil
}

func (c *PatternCache) Get(id api.PatternID) (*api.Pattern, error) {
	return c.store.Get(id)
}

func (c *PatternCache) All() ([]*api.Pattern, error) {
	return c.store.All()
}

func NewPatternCache(store api.PatternStorage) api.PatternStorage {
	return &PatternCache{
		store: store,
		client: redis.NewClient(&redis.Options{
			Addr: "redis:6379",
		}),
	}
}
