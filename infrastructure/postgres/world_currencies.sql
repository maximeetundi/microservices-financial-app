-- =====================================================
-- COMPREHENSIVE WORLD CURRENCIES
-- All 150+ ISO 4217 currencies
-- =====================================================

-- This file adds all world currencies to the payment providers
-- Run this after init.sql to expand currency support

-- Update Demo provider to support ALL currencies
UPDATE payment_providers 
SET supported_currencies = '[
  "AED", "AFN", "ALL", "AMD", "ANG", "AOA", "ARS", "AUD", "AWG", "AZN",
  "BAM", "BBD", "BDT", "BGN", "BHD", "BIF", "BMD", "BND", "BOB", "BRL", "BSD", "BTN", "BWP", "BYN", "BZD",
  "CAD", "CDF", "CHF", "CLP", "CNY", "COP", "CRC", "CUP", "CVE", "CZK",
  "DJF", "DKK", "DOP", "DZD",
  "EGP", "ERN", "ETB", "EUR",
  "FJD", "FKP",
  "GBP", "GEL", "GGP", "GHS", "GIP", "GMD", "GNF", "GTQ", "GYD",
  "HKD", "HNL", "HRK", "HTG", "HUF",
  "IDR", "ILS", "IMP", "INR", "IQD", "IRR", "ISK",
  "JEP", "JMD", "JOD", "JPY",
  "KES", "KGS", "KHR", "KMF", "KPW", "KRW", "KWD", "KYD", "KZT",
  "LAK", "LBP", "LKR", "LRD", "LSL", "LYD",
  "MAD", "MDL", "MGA", "MKD", "MMK", "MNT", "MOP", "MRU", "MUR", "MVR", "MWK", "MXN", "MYR", "MZN",
  "NAD", "NGN", "NIO", "NOK", "NPR", "NZD",
  "OMR",
  "PAB", "PEN", "PGK", "PHP", "PKR", "PLN", "PYG",
  "QAR",
  "RON", "RSD", "RUB", "RWF",
  "SAR", "SBD", "SCR", "SDG", "SEK", "SGD", "SHP", "SLL", "SOS", "SPL", "SRD", "STN", "SYP", "SZL",
  "THB", "TJS", "TMT", "TND", "TOP", "TRY", "TTD", "TVD", "TWD", "TZS",
  "UAH", "UGX", "USD", "UYU", "UZS",
  "VEF", "VND", "VUV",
  "WST",
  "XAF", "XCD", "XDR", "XOF", "XPF",
  "YER",
  "ZAR", "ZMW", "ZWD",
  "BTC", "ETH", "USDT", "USDC", "BNB", "XRP", "ADA", "SOL", "DOT", "DOGE", "MATIC", "LTC", "TRX", "AVAX", "UNI"
]'::jsonb
WHERE name = 'demo';

-- Update Flutterwave to support more currencies
UPDATE payment_providers 
SET supported_currencies = '[
  "NGN", "GHS", "KES", "ZAR", "TZS", "UGX", "RWF", "ZMW",
  "XOF", "XAF", "USD", "EUR", "GBP", "CAD", "AUD"
]'::jsonb
WHERE name = 'flutterwave';

-- Update CinetPay
UPDATE payment_providers 
SET supported_currencies = '[
  "XOF", "XAF", "GNF", "USD", "EUR"
]'::jsonb
WHERE name = 'cinetpay';

-- Update Paystack
UPDATE payment_providers 
SET supported_currencies = '[
  "NGN", "GHS", "ZAR", "KES", "USD", "EUR", "GBP"
]'::jsonb
WHERE name = 'paystack';

-- Update Orange Money
UPDATE payment_providers 
SET supported_currencies = '[
  "XOF", "XAF", "EUR"
]'::jsonb
WHERE name = 'orange_money';

-- Update MTN MoMo
UPDATE payment_providers 
SET supported_currencies = '[
  "XOF", "XAF", "GHS", "UGX", "RWF", "ZMW", "EUR"
]'::jsonb
WHERE name = 'mtn_momo';

-- Update Wave
UPDATE payment_providers 
SET supported_currencies = '[
  "XOF", "EUR", "USD"
]'::jsonb
WHERE name = 'wave';

-- Update Stripe
UPDATE payment_providers 
SET supported_currencies = '[
  "USD", "EUR", "GBP", "CAD", "AUD", "NZD", "CHF", "SEK", "NOK", "DKK",
  "JPY", "CNY", "HKD", "SGD", "MYR", "THB", "PHP", "INR", "KRW",
  "BRL", "MXN", "ARS", "CLP", "COP", "PEN",
  "AED", "SAR", "QAR", "KWD", "BHD", "OMR",
  "ILS", "TRY", "ZAR", "NGN", "KES", "GHS"
]'::jsonb
WHERE name = 'stripe';

-- Add more exchange rate pairs for popular currencies
INSERT INTO exchange_rates (from_currency, to_currency, rate, bid_price, ask_price, spread, source, volume_24h, change_24h, valid_until) VALUES
-- Major Fiat Currencies
('USD', 'GBP', 0.79, 0.7898, 0.7902, 0.0004, 'system', 3000000.00, 0.2, NOW() + INTERVAL '1 day'),
('GBP', 'USD', 1.27, 1.2698, 1.2702, 0.0004, 'system', 3000000.00, -0.2, NOW() + INTERVAL '1 day'),
('USD', 'JPY', 149.50, 149.45, 149.55, 0.10, 'system', 4000000.00, 0.3, NOW() + INTERVAL '1 day'),
('JPY', 'USD', 0.0067, 0.00669, 0.00671, 0.00002, 'system', 4000000.00, -0.3, NOW() + INTERVAL '1 day'),
('USD', 'CHF', 0.88, 0.8798, 0.8802, 0.0004, 'system', 2000000.00, 0.1, NOW() + INTERVAL '1 day'),
('USD', 'CAD', 1.35, 1.3498, 1.3502, 0.0004, 'system', 2500000.00, 0.2, NOW() + INTERVAL '1 day'),
('USD', 'AUD', 1.52, 1.5198, 1.5202, 0.0004, 'system', 2200000.00, 0.4, NOW() + INTERVAL '1 day'),
('USD', 'CNY', 7.24, 7.2398, 7.2402, 0.0004, 'system', 3500000.00, 0.1, NOW() + INTERVAL '1 day'),

-- African Currencies
('USD', 'XOF', 605.00, 604.50, 605.50, 1.00, 'system', 500000.00, 0.1, NOW() + INTERVAL '1 day'),
('XOF', 'USD', 0.00165, 0.001648, 0.001652, 0.000004, 'system', 500000.00, -0.1, NOW() + INTERVAL '1 day'),
('USD', 'XAF', 605.00, 604.50, 605.50, 1.00, 'system', 400000.00, 0.1, NOW() + INTERVAL '1 day'),
('XAF', 'USD', 0.00165, 0.001648, 0.001652, 0.000004, 'system', 400000.00, -0.1, NOW() + INTERVAL '1 day'),
('USD', 'NGN', 1450.00, 1448.00, 1452.00, 4.00, 'system', 800000.00, 0.5, NOW() + INTERVAL '1 day'),
('NGN', 'USD', 0.00069, 0.000688, 0.000692, 0.000004, 'system', 800000.00, -0.5, NOW() + INTERVAL '1 day'),
('USD', 'GHS', 15.50, 15.48, 15.52, 0.04, 'system', 300000.00, 0.3, NOW() + INTERVAL '1 day'),
('GHS', 'USD', 0.0645, 0.0644, 0.0646, 0.0002, 'system', 300000.00, -0.3, NOW() + INTERVAL '1 day'),
('USD', 'KES', 129.00, 128.85, 129.15, 0.30, 'system', 400000.00, 0.2, NOW() + INTERVAL '1 day'),
('KES', 'USD', 0.00775, 0.00774, 0.00776, 0.00002, 'system', 400000.00, -0.2, NOW() + INTERVAL '1 day'),
('USD', 'ZAR', 18.50, 18.48, 18.52, 0.04, 'system', 600000.00, 0.4, NOW() + INTERVAL '1 day'),
('ZAR', 'USD', 0.054, 0.0539, 0.0541, 0.0002, 'system', 600000.00, -0.4, NOW() + INTERVAL '1 day'),
('USD', 'MAD', 10.10, 10.08, 10.12, 0.04, 'system', 250000.00, 0.1, NOW() + INTERVAL '1 day'),
('USD', 'EGP', 48.50, 48.40, 48.60, 0.20, 'system', 350000.00, 0.3, NOW() + INTERVAL '1 day'),
('USD', 'TND', 3.15, 3.14, 3.16, 0.02, 'system', 200000.00, 0.1, NOW() + INTERVAL '1 day'),

-- Asian Currencies
('USD', 'INR', 83.20, 83.15, 83.25, 0.10, 'system', 1500000.00, 0.2, NOW() + INTERVAL '1 day'),
('USD', 'KRW', 1320.00, 1318.00, 1322.00, 4.00, 'system', 900000.00, 0.3, NOW() + INTERVAL '1 day'),
('USD', 'SGD', 1.34, 1.3398, 1.3402, 0.0004, 'system', 700000.00, 0.1, NOW() + INTERVAL '1 day'),
('USD', 'HKD', 7.82, 7.8198, 7.8202, 0.0004, 'system', 800000.00, 0.0, NOW() + INTERVAL '1 day'),
('USD', 'THB', 35.50, 35.48, 35.52, 0.04, 'system', 500000.00, 0.2, NOW() + INTERVAL '1 day'),
('USD', 'MYR', 4.65, 4.6498, 4.6502, 0.0004, 'system', 400000.00, 0.1, NOW() + INTERVAL '1 day'),
('USD', 'PHP', 56.50, 56.45, 56.55, 0.10, 'system', 350000.00, 0.2, NOW() + INTERVAL '1 day'),
('USD', 'IDR', 15800.00, 15790.00, 15810.00, 20.00, 'system', 600000.00, 0.3, NOW() + INTERVAL '1 day'),
('USD', 'VND', 24500.00, 24480.00, 24520.00, 40.00, 'system', 300000.00, 0.1, NOW() + INTERVAL '1 day'),

-- Latin American Currencies
('USD', 'BRL', 5.05, 5.0498, 5.0502, 0.0004, 'system', 900000.00, 0.5, NOW() + INTERVAL '1 day'),
('USD', 'MXN', 17.20, 17.18, 17.22, 0.04, 'system', 800000.00, 0.3, NOW() + INTERVAL '1 day'),
('USD', 'ARS', 850.00, 848.00, 852.00, 4.00, 'system', 400000.00, 1.2, NOW() + INTERVAL '1 day'),
('USD', 'CLP', 920.00, 918.00, 922.00, 4.00, 'system', 300000.00, 0.4, NOW() + INTERVAL '1 day'),
('USD', 'COP', 4100.00, 4095.00, 4105.00, 10.00, 'system', 250000.00, 0.3, NOW() + INTERVAL '1 day'),
('USD', 'PEN', 3.75, 3.7498, 3.7502, 0.0004, 'system', 200000.00, 0.2, NOW() + INTERVAL '1 day'),

-- Middle Eastern Currencies
('USD', 'AED', 3.67, 3.6698, 3.6702, 0.0004, 'system', 500000.00, 0.0, NOW() + INTERVAL '1 day'),
('USD', 'SAR', 3.75, 3.7498, 3.7502, 0.0004, 'system', 450000.00, 0.0, NOW() + INTERVAL '1 day'),
('USD', 'QAR', 3.64, 3.6398, 3.6402, 0.0004, 'system', 300000.00, 0.0, NOW() + INTERVAL '1 day'),
('USD', 'KWD', 0.31, 0.3098, 0.3102, 0.0004, 'system', 250000.00, 0.0, NOW() + INTERVAL '1 day'),
('USD', 'BHD', 0.38, 0.3798, 0.3802, 0.0004, 'system', 200000.00, 0.0, NOW() + INTERVAL '1 day'),
('USD', 'ILS', 3.65, 3.6498, 3.6502, 0.0004, 'system', 350000.00, 0.2, NOW() + INTERVAL '1 day'),
('USD', 'TRY', 32.50, 32.45, 32.55, 0.10, 'system', 600000.00, 0.8, NOW() + INTERVAL '1 day'),

-- European Currencies (non-EUR)
('USD', 'NOK', 10.80, 10.78, 10.82, 0.04, 'system', 400000.00, 0.3, NOW() + INTERVAL '1 day'),
('USD', 'SEK', 10.50, 10.48, 10.52, 0.04, 'system', 450000.00, 0.2, NOW() + INTERVAL '1 day'),
('USD', 'DKK', 6.90, 6.8998, 6.9002, 0.0004, 'system', 350000.00, 0.1, NOW() + INTERVAL '1 day'),
('USD', 'PLN', 4.05, 4.0498, 4.0502, 0.0004, 'system', 500000.00, 0.3, NOW() + INTERVAL '1 day'),
('USD', 'CZK', 23.50, 23.48, 23.52, 0.04, 'system', 300000.00, 0.2, NOW() + INTERVAL '1 day'),
('USD', 'HUF', 360.00, 359.50, 360.50, 1.00, 'system', 250000.00, 0.4, NOW() + INTERVAL '1 day'),
('USD', 'RON', 4.55, 4.5498, 4.5502, 0.0004, 'system', 200000.00, 0.2, NOW() + INTERVAL '1 day'),

-- Crypto to Fiat pairs
('BTC', 'EUR', 38250.00, 38240.00, 38260.00, 20.00, 'binance', 1234567.89, 2.5, NOW() + INTERVAL '1 day'),
('ETH', 'EUR', 2550.00, 2548.00, 2552.00, 4.00, 'binance', 987654.32, 1.8, NOW() + INTERVAL '1 day'),
('BTC', 'GBP', 33000.00, 32990.00, 33010.00, 20.00, 'binance', 654321.12, 2.5, NOW() + INTERVAL '1 day'),
('ETH', 'GBP', 2200.00, 2198.00, 2202.00, 4.00, 'binance', 432109.87, 1.8, NOW() + INTERVAL '1 day')

ON CONFLICT (from_currency, to_currency) DO UPDATE SET
    rate = EXCLUDED.rate,
    bid_price = EXCLUDED.bid_price,
    ask_price = EXCLUDED.ask_price,
    spread = EXCLUDED.spread,
    volume_24h = EXCLUDED.volume_24h,
    change_24h = EXCLUDED.change_24h,
    last_updated = NOW();
