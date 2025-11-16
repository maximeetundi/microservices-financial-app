-- Exchange Service Database Initialization
-- This file sets up the database schema for the exchange service

-- Enable UUID extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Create exchanges table
CREATE TABLE IF NOT EXISTS exchanges (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id VARCHAR NOT NULL,
    from_wallet_id VARCHAR,
    to_wallet_id VARCHAR,
    from_currency VARCHAR NOT NULL,
    to_currency VARCHAR NOT NULL,
    from_amount DECIMAL(20,8) NOT NULL,
    to_amount DECIMAL(20,8) NOT NULL,
    exchange_rate DECIMAL(20,8) NOT NULL,
    fee DECIMAL(20,8) NOT NULL,
    fee_percentage DECIMAL(5,2) NOT NULL,
    status VARCHAR NOT NULL DEFAULT 'pending',
    destination_amount DECIMAL(20,8),
    destination_currency VARCHAR,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    completed_at TIMESTAMP
);

-- Create exchange_rates table
CREATE TABLE IF NOT EXISTS exchange_rates (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    from_currency VARCHAR NOT NULL,
    to_currency VARCHAR NOT NULL,
    rate DECIMAL(20,8) NOT NULL,
    bid_price DECIMAL(20,8) NOT NULL,
    ask_price DECIMAL(20,8) NOT NULL,
    spread DECIMAL(20,8) NOT NULL,
    source VARCHAR NOT NULL,
    volume_24h DECIMAL(20,8) DEFAULT 0,
    change_24h DECIMAL(10,4) DEFAULT 0,
    last_updated TIMESTAMP DEFAULT NOW(),
    UNIQUE(from_currency, to_currency)
);

-- Create quotes table
CREATE TABLE IF NOT EXISTS quotes (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id VARCHAR NOT NULL,
    from_currency VARCHAR NOT NULL,
    to_currency VARCHAR NOT NULL,
    from_amount DECIMAL(20,8) NOT NULL,
    to_amount DECIMAL(20,8) NOT NULL,
    exchange_rate DECIMAL(20,8) NOT NULL,
    fee DECIMAL(20,8) NOT NULL,
    fee_percentage DECIMAL(5,2) NOT NULL,
    valid_until TIMESTAMP NOT NULL,
    estimated_delivery VARCHAR,
    created_at TIMESTAMP DEFAULT NOW()
);

-- Create trading_orders table
CREATE TABLE IF NOT EXISTS trading_orders (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id VARCHAR NOT NULL,
    order_type VARCHAR NOT NULL,
    pair VARCHAR NOT NULL,
    side VARCHAR NOT NULL,
    amount DECIMAL(20,8) NOT NULL,
    price DECIMAL(20,8),
    stop_price DECIMAL(20,8),
    status VARCHAR NOT NULL DEFAULT 'pending',
    filled_amount DECIMAL(20,8) DEFAULT 0,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    filled_at TIMESTAMP
);

-- Create indexes for better performance
CREATE INDEX IF NOT EXISTS idx_exchanges_user_id ON exchanges(user_id);
CREATE INDEX IF NOT EXISTS idx_exchanges_status ON exchanges(status);
CREATE INDEX IF NOT EXISTS idx_exchanges_created_at ON exchanges(created_at);
CREATE INDEX IF NOT EXISTS idx_rates_pair ON exchange_rates(from_currency, to_currency);
CREATE INDEX IF NOT EXISTS idx_rates_updated ON exchange_rates(last_updated);
CREATE INDEX IF NOT EXISTS idx_quotes_user_id ON quotes(user_id);
CREATE INDEX IF NOT EXISTS idx_quotes_valid_until ON quotes(valid_until);
CREATE INDEX IF NOT EXISTS idx_orders_user_id ON trading_orders(user_id);
CREATE INDEX IF NOT EXISTS idx_orders_status ON trading_orders(status);
CREATE INDEX IF NOT EXISTS idx_orders_pair ON trading_orders(pair);

-- Insert sample exchange rates
INSERT INTO exchange_rates (from_currency, to_currency, rate, bid_price, ask_price, spread, source, volume_24h, change_24h) VALUES
-- Crypto rates
('BTC', 'USD', 43500.00, 43485.00, 43515.00, 30.00, 'crypto_exchange_api', 1234567.89, 2.3),
('ETH', 'USD', 2450.00, 2448.00, 2452.00, 4.00, 'crypto_exchange_api', 987654.32, -1.2),
('LTC', 'USD', 72.00, 71.96, 72.04, 0.08, 'crypto_exchange_api', 543210.11, 0.8),
('ADA', 'USD', 0.52, 0.5198, 0.5202, 0.0004, 'crypto_exchange_api', 876543.21, 3.1),
('DOT', 'USD', 7.80, 7.798, 7.802, 0.004, 'crypto_exchange_api', 654321.45, -0.5),
('XRP', 'USD', 0.63, 0.6298, 0.6302, 0.0004, 'crypto_exchange_api', 432109.87, 1.7),

-- Fiat rates
('USD', 'EUR', 0.8456, 0.8454, 0.8458, 0.0004, 'fiat_exchange_api', 5000000.00, 0.1),
('EUR', 'USD', 1.1826, 1.1824, 1.1828, 0.0004, 'fiat_exchange_api', 5000000.00, -0.1),
('USD', 'GBP', 0.7821, 0.7819, 0.7823, 0.0004, 'fiat_exchange_api', 3000000.00, -0.2),
('GBP', 'USD', 1.2787, 1.2785, 1.2789, 0.0004, 'fiat_exchange_api', 3000000.00, 0.3),
('USD', 'JPY', 149.23, 149.20, 149.26, 0.06, 'fiat_exchange_api', 4000000.00, -0.1),
('JPY', 'USD', 0.0067, 0.00669, 0.00671, 0.00002, 'fiat_exchange_api', 4000000.00, 0.1),
('USD', 'CAD', 1.3567, 1.3565, 1.3569, 0.0004, 'fiat_exchange_api', 2000000.00, 0.2),
('CAD', 'USD', 0.7371, 0.7369, 0.7373, 0.0004, 'fiat_exchange_api', 2000000.00, -0.2),

-- Cross crypto-fiat rates
('BTC', 'EUR', 36785.40, 36770.40, 36800.40, 30.00, 'crypto_exchange_api', 876543.21, 2.4),
('ETH', 'EUR', 2071.82, 2069.82, 2073.82, 4.00, 'crypto_exchange_api', 654321.45, -1.1),
('BTC', 'GBP', 34018.35, 34003.35, 34033.35, 30.00, 'crypto_exchange_api', 543210.11, 2.1),
('ETH', 'GBP', 1915.45, 1913.45, 1917.45, 4.00, 'crypto_exchange_api', 432109.87, -1.4)
ON CONFLICT (from_currency, to_currency) DO UPDATE SET
    rate = EXCLUDED.rate,
    bid_price = EXCLUDED.bid_price,
    ask_price = EXCLUDED.ask_price,
    spread = EXCLUDED.spread,
    volume_24h = EXCLUDED.volume_24h,
    change_24h = EXCLUDED.change_24h,
    last_updated = NOW();

-- Create function to update updated_at timestamp
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Create triggers to automatically update updated_at
CREATE TRIGGER update_exchanges_updated_at BEFORE UPDATE ON exchanges
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_trading_orders_updated_at BEFORE UPDATE ON trading_orders
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- Grant necessary permissions (adjust as needed for your setup)
-- GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public TO exchange_user;
-- GRANT ALL PRIVILEGES ON ALL SEQUENCES IN SCHEMA public TO exchange_user;