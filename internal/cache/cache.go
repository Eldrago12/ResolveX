package cache

import (
	"context"
	"crypto/tls"
	"log"
	"net/url"
	"time"

	"github.com/go-redis/redis/v8"
)

type Cache struct {
	client *redis.Client
}

func NewCache(redisURL string) *Cache {
	parsedURL, err := url.Parse(redisURL)
	if err != nil {
		log.Fatalf("Invalid Redis URL: %v", err)
	}

	password, _ := parsedURL.User.Password()
	host := parsedURL.Host

	client := redis.NewClient(&redis.Options{
		Addr:      host,
		Password:  password,
		DB:        0,
		TLSConfig: &tls.Config{InsecureSkipVerify: true}, // Set to true if Aiven uses self-signed certs; otherwise, omit or set false
	})

	return &Cache{client: client}
}

func (c *Cache) Get(ctx context.Context, key string) (string, error) {
	return c.client.Get(ctx, key).Result()
}

func (c *Cache) Set(ctx context.Context, key string, value string, ttl time.Duration) error {
	return c.client.Set(ctx, key, value, ttl).Err()
}

func (c *Cache) Increment(ctx context.Context, key string, ttl time.Duration) (int64, error) {
	count, err := c.client.Incr(ctx, key).Result()
	if err != nil {
		return 0, err
	}
	if count == 1 {
		err = c.client.Expire(ctx, key, ttl).Err()
	}
	return count, err
}

func (c *Cache) TrackFrequency(ctx context.Context, domain string) error {
	return c.client.ZIncrBy(ctx, "domain_frequency", 1, domain).Err()
}

func (c *Cache) GetTopDomains(ctx context.Context, limit int) ([]string, error) {
	domains, err := c.client.ZRevRange(ctx, "domain_frequency", 0, int64(limit-1)).Result()
	return domains, err
}

func (c *Cache) Close() error {
	return c.client.Close()
}
