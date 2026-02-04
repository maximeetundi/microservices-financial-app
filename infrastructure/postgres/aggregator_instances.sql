-- =====================================================
-- AGGREGATOR INSTANCES & MULTI-HOT WALLET LINKING
-- Une instance peut avoir PLUSIEURS hot wallets
-- =====================================================

-- Table des instances d'agrégateurs
CREATE TABLE IF NOT EXISTS aggregator_instances (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    aggregator_id UUID NOT NULL REFERENCES aggregator_settings(id) ON DELETE CASCADE,
    instance_name VARCHAR(100) NOT NULL, -- Ex: "Flutterwave CI-1", "Stripe Europe"
    
    -- Clés API (stockées encryptées ou référence Vault)
    api_credentials JSONB NOT NULL, -- Clés spécifiques à cette instance
    
    -- Configuration spécifique
    enabled BOOLEAN DEFAULT true,
    priority INT DEFAULT 50, -- Plus élevé = préféré
    
    -- Limites de transaction
    daily_limit DECIMAL(20, 8), -- Limite journalière (NULL = illimité)
    monthly_limit DECIMAL(20, 8), -- Limite mensuelle (NULL = illimité)
    daily_usage DECIMAL(20, 8) DEFAULT 0,
    monthly_usage DECIMAL(20, 8) DEFAULT 0,
    usage_reset_date DATE DEFAULT CURRENT_DATE,
    
    -- Restricted to specific countries or ALL
    restricted_countries TEXT[], -- NULL = tous pays, sinon liste ISO codes
    
    -- Mode test/production
    is_test_mode BOOLEAN DEFAULT true,
    
    -- Statistiques
    total_transactions INT DEFAULT 0,
    total_volume DECIMAL(20, 8) DEFAULT 0,
    last_used_at TIMESTAMP,
    
    -- Metadata
    notes TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_by UUID REFERENCES admins(id)
);

-- Table de liaison MANY-TO-MANY entre instances et hot wallets
CREATE TABLE IF NOT EXISTS aggregator_instance_wallets (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    instance_id UUID NOT NULL REFERENCES aggregator_instances(id) ON DELETE CASCADE,
    hot_wallet_id UUID NOT NULL REFERENCES platform_accounts(id) ON DELETE CASCADE,
    
    -- Configuration par wallet
    is_primary BOOLEAN DEFAULT false, -- Si true = wallet préféré pour cette instance
    priority INT DEFAULT 50, -- Priorité d'utilisation de ce wallet
    
    -- Limites spécifiques au wallet
    min_balance DECIMAL(20, 8) DEFAULT 0, -- Balance minimum pour utiliser ce wallet
    max_balance DECIMAL(20, 8), -- Balance maximum (NULL = illimité)
    
    -- Auto-recharge configuration
    auto_recharge_enabled BOOLEAN DEFAULT false,
    recharge_threshold DECIMAL(20, 8), -- Déclencher recharge si balance < seuil
    recharge_target DECIMAL(20, 8), -- Recharger jusqu'à ce montant
    recharge_source_wallet_id UUID REFERENCES platform_accounts(id), -- D'où prendre les fonds
    
    -- Statistiques par wallet
    total_deposits DECIMAL(20, 8) DEFAULT 0,
    total_withdrawals DECIMAL(20, 8) DEFAULT 0,
    transaction_count INT DEFAULT 0,
    last_used_at TIMESTAMP,
    
    enabled BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    -- Contrainte: une seule combinaison instance+wallet
    UNIQUE(instance_id, hot_wallet_id)
);

-- Index
CREATE INDEX IF NOT EXISTS idx_aggregator_instances_aggregator ON aggregator_instances(aggregator_id);
CREATE INDEX IF NOT EXISTS idx_aggregator_instances_enabled ON aggregator_instances(enabled);
CREATE INDEX IF NOT EXISTS idx_aggregator_instances_priority ON aggregator_instances(priority DESC);

CREATE INDEX IF NOT EXISTS idx_aggregator_instance_wallets_instance ON aggregator_instance_wallets(instance_id);
CREATE INDEX IF NOT EXISTS idx_aggregator_instance_wallets_wallet ON aggregator_instance_wallets(hot_wallet_id);
CREATE INDEX IF NOT EXISTS idx_aggregator_instance_wallets_primary ON aggregator_instance_wallets(instance_id, is_primary) WHERE is_primary = true;

-- Trigger pour updated_at
CREATE OR REPLACE FUNCTION update_aggregator_instance_timestamp()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER aggregator_instance_updated_at
    BEFORE UPDATE ON aggregator_instances
    FOR EACH ROW
    EXECUTE FUNCTION update_aggregator_instance_timestamp();

CREATE TRIGGER aggregator_instance_wallet_updated_at
    BEFORE UPDATE ON aggregator_instance_wallets
    FOR EACH ROW
    EXECUTE FUNCTION update_aggregator_instance_timestamp();

-- Table de transactions par instance (pour tracking)
CREATE TABLE IF NOT EXISTS aggregator_instance_transactions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    instance_id UUID NOT NULL REFERENCES aggregator_instances(id),
    wallet_id UUID REFERENCES platform_accounts(id), -- Quel hot wallet a été utilisé
    transaction_id UUID NOT NULL, -- Référence à la transaction principale
    transaction_type VARCHAR(20) NOT NULL, -- 'deposit' ou 'withdrawal'
    amount DECIMAL(20, 8) NOT NULL,
    currency VARCHAR(10) NOT NULL,
    status VARCHAR(20) NOT NULL,
    provider_reference VARCHAR(255),
    user_id UUID, -- Utilisateur concerné
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_aggregator_instance_txn_instance ON aggregator_instance_transactions(instance_id);
CREATE INDEX IF NOT EXISTS idx_aggregator_instance_txn_wallet ON aggregator_instance_transactions(wallet_id);
CREATE INDEX IF NOT EXISTS idx_aggregator_instance_txn_transaction ON aggregator_instance_transactions(transaction_id);
CREATE INDEX IF NOT EXISTS idx_aggregator_instance_txn_created ON aggregator_instance_transactions(created_at DESC);

-- =====================================================
-- VUE POUR SÉLECTION INTELLIGENTE DE WALLET
-- =====================================================

CREATE OR REPLACE VIEW aggregator_instance_wallets_available AS
SELECT 
    aiw.id,
    aiw.instance_id,
    aiw.hot_wallet_id,
    aiw.is_primary,
    aiw.priority,
    aiw.min_balance,
    aiw.max_balance,
    aiw.enabled,
    -- Instance info
    ai.aggregator_id,
    ai.instance_name,
    ai.enabled AS instance_enabled,
    agg.provider_code,
    agg.provider_name,
    agg.is_enabled AS aggregator_enabled,
    -- Wallet info
    pa.currency AS wallet_currency,
    pa.balance AS wallet_balance,
    pa.account_type,
    -- Disponibilité
    CASE 
        WHEN NOT aiw.enabled THEN 'wallet_disabled'
        WHEN NOT ai.enabled THEN 'instance_disabled'
        WHEN NOT agg.is_enabled THEN 'aggregator_disabled'
        WHEN pa.balance < COALESCE(aiw.min_balance, 0) THEN 'insufficient_balance'
        WHEN aiw.max_balance IS NOT NULL AND pa.balance > aiw.max_balance THEN 'balance_too_high'
        ELSE 'available'
    END AS availability_status
FROM aggregator_instance_wallets aiw
JOIN aggregator_instances ai ON aiw.instance_id = ai.id
JOIN aggregator_settings agg ON ai.aggregator_id = agg.id
JOIN platform_accounts pa ON aiw.hot_wallet_id = pa.id;

-- =====================================================
-- FONCTION: Sélectionner le meilleur wallet pour une instance
-- =====================================================

CREATE OR REPLACE FUNCTION select_best_wallet_for_instance(
    p_instance_id UUID,
    p_currency VARCHAR(10),
    p_amount DECIMAL(20, 8) DEFAULT 0
)
RETURNS UUID AS $$
DECLARE
    v_wallet_id UUID;
BEGIN
    -- Sélectionner le meilleur wallet disponible:
    -- 1. Primary wallet en premier
    -- 2. Sinon par priorité
    -- 3. Sinon par balance (plus élevée d'abord)
    SELECT hot_wallet_id INTO v_wallet_id
    FROM aggregator_instance_wallets_available
    WHERE instance_id = p_instance_id
      AND wallet_currency = p_currency
      AND availability_status = 'available'
      AND wallet_balance >= p_amount -- Assez pour couvrir le montant
    ORDER BY 
        is_primary DESC,
        priority DESC,
        wallet_balance DESC
    LIMIT 1;
    
    RETURN v_wallet_id;
END;
$$ LANGUAGE plpgsql;

-- =====================================================
-- EXEMPLES (Seed Data)
-- =====================================================

DO $$
DECLARE
    v_flutterwave_id UUID;
    v_instance_id UUID;
    v_hot_wallet_ngn UUID;
    v_hot_wallet_xof UUID;
    v_cold_wallet_ngn UUID;
    v_cold_wallet_xof UUID;
BEGIN
    -- Récupérer Flutterwave
    SELECT id INTO v_flutterwave_id 
    FROM aggregator_settings 
    WHERE provider_code = 'flutterwave' 
    LIMIT 1;
    
    IF v_flutterwave_id IS NULL THEN
        RETURN;
    END IF;
    
    -- Créer une instance Flutterwave
    INSERT INTO aggregator_instances (
        aggregator_id, 
        instance_name, 
        api_credentials,
        enabled,
        priority,
        is_test_mode,
        notes
    ) VALUES (
        v_flutterwave_id,
        'Flutterwave Multi-Currency',
        jsonb_build_object(
            'public_key', 'FLWPUBK_TEST-CHANGE_ME',
            'secret_key', 'FLWSECK_TEST-CHANGE_ME',
            'encryption_key', 'FLWSECK_TEST-CHANGE_ME',
            'webhook_secret', ''
        ),
        true,
        100,
        true,
        'Instance principale avec multi-wallets'
    )
    ON CONFLICT DO NOTHING
    RETURNING id INTO v_instance_id;
    
    IF v_instance_id IS NULL THEN
        SELECT id INTO v_instance_id
        FROM aggregator_instances
        WHERE instance_name = 'Flutterwave Multi-Currency'
        LIMIT 1;
    END IF;
    
    IF v_instance_id IS NOT NULL THEN
        -- Récupérer hot wallets existants (operations) et cold wallets (réserve)
        -- Hot wallet NGN (compte operations)
        SELECT id INTO v_hot_wallet_ngn 
        FROM platform_accounts 
        WHERE account_type = 'operations' 
        AND currency = 'NGN'
        LIMIT 1;
        
        -- Cold wallet NGN (réserve)
        SELECT id INTO v_cold_wallet_ngn 
        FROM platform_accounts 
        WHERE account_type = 'reserve' 
        AND currency = 'NGN'
        LIMIT 1;
        
        -- Hot wallet XOF
        SELECT id INTO v_hot_wallet_xof 
        FROM platform_accounts 
        WHERE account_type = 'operations' 
        AND currency = 'XOF'
        LIMIT 1;
        
        -- Cold wallet XOF
        SELECT id INTO v_cold_wallet_xof 
        FROM platform_accounts 
        WHERE account_type = 'reserve' 
        AND currency = 'XOF'
        LIMIT 1;
        
        -- Lier wallet NGN (primary) si existe
        IF v_hot_wallet_ngn IS NOT NULL THEN
            INSERT INTO aggregator_instance_wallets (
                instance_id,
                hot_wallet_id,
                is_primary,
                priority,
                min_balance,
                auto_recharge_enabled,
                recharge_threshold,
                recharge_target,
                recharge_source_wallet_id
            ) VALUES (
                v_instance_id,
                v_hot_wallet_ngn,
                true,
                100,
                10000, -- Minimum 10k NGN
                true,
                50000, -- Recharger si < 50k
                200000, -- Recharger jusqu'à 200k
                v_cold_wallet_ngn
            ) ON CONFLICT (instance_id, hot_wallet_id) DO NOTHING;
        END IF;
        
        -- Lier wallet XOF (secondary) si existe
        IF v_hot_wallet_xof IS NOT NULL THEN
            INSERT INTO aggregator_instance_wallets (
                instance_id,
                hot_wallet_id,
                is_primary,
                priority,
                min_balance,
                auto_recharge_enabled,
                recharge_threshold,
                recharge_target,
                recharge_source_wallet_id
            ) VALUES (
                v_instance_id,
                v_hot_wallet_xof,
                false,
                90,
                5000, -- Minimum 5k XOF
                true,
                25000,
                100000,
                v_cold_wallet_xof
            ) ON CONFLICT (instance_id, hot_wallet_id) DO NOTHING;
        END IF;
    END IF;
    
END $$;

