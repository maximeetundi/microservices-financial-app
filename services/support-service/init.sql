-- Support Service Database Schema

-- Agents table
CREATE TABLE IF NOT EXISTS agents (
    id VARCHAR(255) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255),
    type VARCHAR(50) NOT NULL DEFAULT 'human',
    avatar VARCHAR(255),
    is_available BOOLEAN DEFAULT true,
    max_chats INTEGER DEFAULT 5,
    active_chats INTEGER DEFAULT 0,
    rating DECIMAL(3,2) DEFAULT 5.00,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Conversations table
CREATE TABLE IF NOT EXISTS conversations (
    id VARCHAR(255) PRIMARY KEY,
    user_id VARCHAR(255) NOT NULL,
    user_name VARCHAR(255) NOT NULL,
    user_email VARCHAR(255),
    agent_id VARCHAR(255) REFERENCES agents(id),
    agent_type VARCHAR(50) NOT NULL DEFAULT 'ai',
    subject VARCHAR(500) NOT NULL,
    category VARCHAR(100) NOT NULL,
    status VARCHAR(50) NOT NULL DEFAULT 'open',
    priority VARCHAR(50) DEFAULT 'medium',
    last_message TEXT,
    last_message_at TIMESTAMP,
    unread_count INTEGER DEFAULT 0,
    message_count INTEGER DEFAULT 0,
    rating INTEGER,
    feedback TEXT,
    resolved_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Messages table
CREATE TABLE IF NOT EXISTS messages (
    id VARCHAR(255) PRIMARY KEY,
    conversation_id VARCHAR(255) NOT NULL REFERENCES conversations(id) ON DELETE CASCADE,
    sender_id VARCHAR(255) NOT NULL,
    sender_name VARCHAR(255) NOT NULL,
    sender_type VARCHAR(50) NOT NULL,
    content TEXT NOT NULL,
    content_type VARCHAR(50) DEFAULT 'text',
    is_read BOOLEAN DEFAULT false,
    read_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Quick replies table (for common responses)
CREATE TABLE IF NOT EXISTS quick_replies (
    id VARCHAR(255) PRIMARY KEY,
    label VARCHAR(255) NOT NULL,
    response TEXT NOT NULL,
    category VARCHAR(100),
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Indexes for performance
CREATE INDEX IF NOT EXISTS idx_conversations_user_id ON conversations(user_id);
CREATE INDEX IF NOT EXISTS idx_conversations_status ON conversations(status);
CREATE INDEX IF NOT EXISTS idx_conversations_agent_id ON conversations(agent_id);
CREATE INDEX IF NOT EXISTS idx_messages_conversation_id ON messages(conversation_id);
CREATE INDEX IF NOT EXISTS idx_messages_created_at ON messages(created_at);

-- Insert default AI agent
INSERT INTO agents (id, name, email, type, avatar, is_available, max_chats, rating)
VALUES ('ai-agent-001', 'Assistant IA CryptoBank', 'ai@cryptobank.com', 'ai', '🤖', true, 1000, 5.0)
ON CONFLICT (id) DO NOTHING;

-- Insert some quick replies
INSERT INTO quick_replies (id, label, response, category, is_active) VALUES
('qr-001', 'Bonjour', 'Bonjour ! Comment puis-je vous aider aujourd''hui ?', 'greeting', true),
('qr-002', 'Transfert en cours', 'Votre transfert est en cours de traitement. Vous recevrez une notification dès qu''il sera terminé.', 'transfer', true),
('qr-003', 'Vérification requise', 'Pour des raisons de sécurité, nous devons vérifier quelques informations. Un email vous a été envoyé.', 'security', true),
('qr-004', 'Temps d''attente', 'Nous traitons votre demande. Le temps de traitement moyen est de 24-48h ouvrées.', 'general', true),
('qr-005', 'Fermeture ticket', 'Votre demande a été traitée. N''hésitez pas à nous recontacter si nécessaire. Bonne journée !', 'closing', true)
ON CONFLICT (id) DO NOTHING;
