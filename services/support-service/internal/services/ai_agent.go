package services

import (
	"strings"
	"time"

	"github.com/crypto-bank/microservices-financial-app/services/support-service/internal/config"
	"github.com/crypto-bank/microservices-financial-app/services/support-service/internal/models"
)

// AIAgent handles AI-powered responses for customer support
type AIAgent struct {
	config *config.Config
}

func NewAIAgent(cfg *config.Config) *AIAgent {
	return &AIAgent{
		config: cfg,
	}
}

// GenerateResponse generates an AI response based on the message and context
func (a *AIAgent) GenerateResponse(message string, conversation *models.Conversation) string {
	// Simple keyword-based responses for now
	// In production, this would call an actual AI API (OpenAI, Claude, etc.)
	
	messageLower := strings.ToLower(message)
	
	// Greeting responses
	if containsAny(messageLower, []string{"bonjour", "salut", "hello", "hi", "bonsoir"}) {
		return "Bonjour ! 👋 Je suis l'assistant virtuel de Zekora. Comment puis-je vous aider aujourd'hui ?"
	}
	
	// Account related
	if containsAny(messageLower, []string{"compte", "account", "solde", "balance"}) {
		return "Pour toute question concernant votre compte ou votre solde, veuillez vous connecter à l'application et consulter la section 'Mon Compte'. Si vous rencontrez des difficultés, je peux transférer votre demande à un conseiller humain."
	}
	
	// Transfer related
	if containsAny(messageLower, []string{"virement", "transfer", "envoyer", "recevoir"}) {
		return "Pour effectuer un virement :\n1. Allez dans 'Transferts'\n2. Sélectionnez le type de transfert\n3. Entrez les informations du bénéficiaire\n4. Confirmez avec votre PIN\n\nAvez-vous besoin d'aide avec une étape spécifique ?"
	}
	
	// Card related
	if containsAny(messageLower, []string{"carte", "card", "visa", "mastercard"}) {
		return "Pour les services de carte :\n• Commander une nouvelle carte : Menu → Cartes → Commander\n• Bloquer votre carte : Menu → Cartes → Bloquer (disponible 24/7)\n• Consulter les transactions : Menu → Cartes → Historique\n\nQue souhaitez-vous faire ?"
	}
	
	// Security related
	if containsAny(messageLower, []string{"sécurité", "security", "password", "mot de passe", "pin", "hack", "fraude"}) {
		return "🔒 La sécurité de votre compte est notre priorité. Pour toute question liée à la sécurité (fraude suspectée, modification de mot de passe, etc.), je vous recommande de parler directement avec un conseiller humain. Voulez-vous que je transfère cette conversation ?"
	}
	
	// Crypto related
	if containsAny(messageLower, []string{"bitcoin", "btc", "ethereum", "eth", "crypto", "cryptomonnaie"}) {
		return "Zekora vous permet d'acheter, vendre et stocker des cryptomonnaies :\n• BTC, ETH, USDT, et plus encore\n• Frais compétitifs (0.5-0.75%)\n• Wallet sécurisé intégré\n\nPour acheter des cryptos, allez dans Exchange → Acheter Crypto."
	}
	
	// Fees related
	if containsAny(messageLower, []string{"frais", "fees", "commission", "tarif"}) {
		return "📊 Nos frais :\n• Transferts SEPA : Gratuit\n• Crypto-Crypto : 0.5%\n• Fiat-Crypto : 0.75%\n• Fiat-Fiat : 0.15-0.25%\n• Retraits ATM : 2€/retrait\n\nBesoin de plus de détails sur un type de frais ?"
	}
	
	// Help/escalation
	if containsAny(messageLower, []string{"humain", "human", "agent", "conseiller", "parler", "réel"}) {
		return "Je comprends que vous souhaitez parler à un conseiller humain. Utilisez le bouton 'Escalader vers un agent' en haut de la conversation pour être mis en relation avec un de nos conseillers. Le temps d'attente moyen est de 2-3 minutes."
	}
	
	// Thanks
	if containsAny(messageLower, []string{"merci", "thanks", "thank you"}) {
		return "Je vous en prie ! 😊 N'hésitez pas si vous avez d'autres questions. Bonne journée !"
	}
	
	// Goodbye
	if containsAny(messageLower, []string{"bye", "aurevoir", "au revoir", "bonne journée", "à bientôt"}) {
		return "Au revoir ! 👋 N'hésitez pas à revenir si vous avez d'autres questions. L'équipe Zekora vous souhaite une excellente journée !"
	}
	
	// Default response
	return "Je ne suis pas sûr de comprendre votre demande. Pourriez-vous reformuler ou choisir parmi ces options ?\n\n• 💳 Questions sur les cartes\n• 💸 Aide aux transferts\n• 📊 Informations sur les frais\n• 🔐 Sécurité du compte\n• ₿ Cryptomonnaies\n• 👤 Parler à un conseiller humain"
}

// ShouldEscalate determines if the conversation should be escalated to a human
func (a *AIAgent) ShouldEscalate(message string, messageCount int) (bool, string) {
	messageLower := strings.ToLower(message)
	
	// Explicit escalation request
	if containsAny(messageLower, []string{"humain", "human", "agent", "conseiller", "parler réel", "vraie personne"}) {
		return true, "Le client demande à parler à un conseiller humain."
	}
	
	// Security concerns should be escalated
	if containsAny(messageLower, []string{"fraude", "fraud", "hack", "volé", "stolen", "urgent"}) {
		return true, "Potentielle fraude ou problème de sécurité urgent."
	}
	
	// Complaint
	if containsAny(messageLower, []string{"plainte", "complaint", "insatisfait", "mécontent", "problème grave"}) {
		return true, "Le client exprime une insatisfaction nécessitant une intervention humaine."
	}
	
	// After many messages without resolution
	if messageCount > 10 {
		return true, "Conversation prolongée - escalade automatique après 10 messages."
	}
	
	return false, ""
}

// GetAvailableTopics returns topics the AI can help with
func (a *AIAgent) GetAvailableTopics() []string {
	return []string{
		"Compte et solde",
		"Transferts et virements",
		"Cartes bancaires",
		"Cryptomonnaies",
		"Frais et tarifs",
		"Sécurité",
	}
}

// Helper function to check if string contains any of the keywords
func containsAny(s string, keywords []string) bool {
	for _, keyword := range keywords {
		if strings.Contains(s, keyword) {
			return true
		}
	}
	return false
}

// GetWelcomeMessage returns the welcome message for new conversations
func (a *AIAgent) GetWelcomeMessage() string {
	hour := time.Now().Hour()
	greeting := "Bonjour"
	if hour >= 18 || hour < 6 {
		greeting = "Bonsoir"
	}
	
	return greeting + " ! 👋 Je suis l'assistant virtuel Zekora. Je suis là pour vous aider 24/7.\n\nVoici ce que je peux faire pour vous :\n• 💳 Assistance cartes bancaires\n• 💸 Aide aux transferts\n• ₿ Questions sur les cryptomonnaies\n• 📊 Informations sur les frais\n• 🔐 Sécurité du compte\n\nComment puis-je vous aider ?"
}
