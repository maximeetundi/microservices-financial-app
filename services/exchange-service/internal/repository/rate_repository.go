package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/crypto-bank/exchange-service/internal/models"
	"github.com/go-redis/redis/v8"
)

type RateRepository struct {
	db    *sql.DB
	redis *redis.Client
}

func NewRateRepository(db *sql.DB, redis *redis.Client) *RateRepository {
	return &RateRepository{
		db:    db,
		redis: redis,
	}
}

func (r *RateRepository) GetRate(fromCurrency, toCurrency string) (*models.ExchangeRate, error) {
	// Try Redis cache first
	key := fmt.Sprintf("rate:%s:%s", fromCurrency, toCurrency)
	val, err := r.redis.Get(context.Background(), key).Result()
	if err == nil {
		var rate models.ExchangeRate
		if err := json.Unmarshal([]byte(val), &rate); err == nil {
			return &rate, nil
		}
	}

	// Fallback to database
	rate := &models.ExchangeRate{}
	query := `
		SELECT from_currency, to_currency, rate, bid_price, ask_price, spread,
		       source, volume_24h, change_24h, last_updated
		FROM exchange_rates 
		WHERE from_currency = $1 AND to_currency = $2`

	err = r.db.QueryRow(query, fromCurrency, toCurrency).Scan(
		&rate.FromCurrency, &rate.ToCurrency, &rate.Rate,
		&rate.BidPrice, &rate.AskPrice, &rate.Spread,
		&rate.Source, &rate.Volume24h, &rate.Change24h, &rate.LastUpdated)

	if err != nil {
		return nil, err
	}

	// Cache in Redis for 1 minute
	r.cacheRate(rate, 1*time.Minute)

	return rate, nil
}

func (r *RateRepository) SaveRate(rate *models.ExchangeRate) error {
	query := `
		INSERT INTO exchange_rates (from_currency, to_currency, rate, bid_price, ask_price,
		                           spread, source, volume_24h, change_24h, last_updated)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		ON CONFLICT (from_currency, to_currency)
		DO UPDATE SET 
			rate = EXCLUDED.rate,
			bid_price = EXCLUDED.bid_price,
			ask_price = EXCLUDED.ask_price,
			spread = EXCLUDED.spread,
			source = EXCLUDED.source,
			volume_24h = EXCLUDED.volume_24h,
			change_24h = EXCLUDED.change_24h,
			last_updated = EXCLUDED.last_updated`

	_, err := r.db.Exec(query,
		rate.FromCurrency, rate.ToCurrency, rate.Rate,
		rate.BidPrice, rate.AskPrice, rate.Spread,
		rate.Source, rate.Volume24h, rate.Change24h, rate.LastUpdated)

	if err != nil {
		return err
	}

	// Cache in Redis for 1 minute
	r.cacheRate(rate, 1*time.Minute)

	return nil
}

func (r *RateRepository) GetAllRates() ([]*models.ExchangeRate, error) {
	query := `
		SELECT from_currency, to_currency, rate, bid_price, ask_price, spread,
		       source, volume_24h, change_24h, last_updated
		FROM exchange_rates 
		ORDER BY from_currency, to_currency`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var rates []*models.ExchangeRate
	for rows.Next() {
		rate := &models.ExchangeRate{}
		err := rows.Scan(
			&rate.FromCurrency, &rate.ToCurrency, &rate.Rate,
			&rate.BidPrice, &rate.AskPrice, &rate.Spread,
			&rate.Source, &rate.Volume24h, &rate.Change24h, &rate.LastUpdated)
		if err != nil {
			return nil, err
		}
		rates = append(rates, rate)
	}

	return rates, nil
}

func (r *RateRepository) cacheRate(rate *models.ExchangeRate, ttl time.Duration) {
	key := fmt.Sprintf("rate:%s:%s", rate.FromCurrency, rate.ToCurrency)
	data, err := json.Marshal(rate)
	if err == nil {
		r.redis.Set(context.Background(), key, data, ttl)
	}
}

func (r *RateRepository) InvalidateCache(fromCurrency, toCurrency string) {
	key := fmt.Sprintf("rate:%s:%s", fromCurrency, toCurrency)
	r.redis.Del(context.Background(), key)
}