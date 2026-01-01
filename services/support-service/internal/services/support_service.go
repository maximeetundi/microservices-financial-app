package services

import (
	"context"
	"fmt"
	"time"

	"github.com/crypto-bank/microservices-financial-app/services/support-service/internal/config"
	"github.com/crypto-bank/microservices-financial-app/services/support-service/internal/models"
	"github.com/crypto-bank/microservices-financial-app/services/support-service/internal/repository"
)

type SupportService struct {
	convRepo    *repository.ConversationRepository
	msgRepo     *repository.MessageRepository
	agentRepo   *repository.AgentRepository
	aiAgent     *AIAgent
	config      *config.Config
	publisher   *EventPublisher
	storage     *StorageService
}

func NewSupportService(
	convRepo *repository.ConversationRepository,
	msgRepo *repository.MessageRepository,
	agentRepo *repository.AgentRepository,
	cfg *config.Config,
	publisher *EventPublisher,
	storage *StorageService,
) *SupportService {
	return &SupportService{
		convRepo:  convRepo,
		msgRepo:   msgRepo,
		agentRepo: agentRepo,
		aiAgent:   NewAIAgent(cfg),
		config:    cfg,
		publisher: publisher,
		storage:   storage,
	}
}

// CreateConversation starts a new support conversation
func (s *SupportService) CreateConversation(userID, userName, userEmail string, req *models.CreateConversationRequest) (*models.Conversation, *models.Message, error) {
	// Create the conversation
	conv := &models.Conversation{
		UserID:     userID,
		UserName:   userName,
		UserEmail:  userEmail,
		AgentType:  req.AgentType,
		Subject:    req.Subject,
		Category:   req.Category,
		Priority:   req.Priority,
		Status:     models.ConversationStatusOpen,
	}

	if conv.Priority == "" {
		conv.Priority = models.PriorityMedium
	}

	if err := s.convRepo.Create(conv); err != nil {
		return nil, nil, fmt.Errorf("failed to create conversation: %w", err)
	}

	// Create the initial user message
	userMsg := &models.Message{
		ConversationID: conv.ID,
		SenderID:       userID,
		SenderName:     userName,
		SenderType:     models.MessageTypeUser,
		Content:        req.Message,
		ContentType:    "text",
		IsRead:         false,
	}

	if err := s.msgRepo.Create(userMsg); err != nil {
		return nil, nil, fmt.Errorf("failed to create message: %w", err)
	}

	// Update conversation with last message
	s.convRepo.UpdateLastMessage(conv.ID, req.Message)

	// If AI agent, generate automatic response
	if req.AgentType == models.AgentTypeAI {
		aiAgent, _ := s.agentRepo.GetAIAgent()
		if aiAgent != nil {
			conv.AgentID = &aiAgent.ID
			s.convRepo.AssignAgent(conv.ID, aiAgent.ID)

			// Generate welcome message
			welcomeMsg := &models.Message{
				ConversationID: conv.ID,
				SenderID:       aiAgent.ID,
				SenderName:     aiAgent.Name,
				SenderType:     models.MessageTypeAgent,
				Content:        s.aiAgent.GetWelcomeMessage(),
				ContentType:    "text",
				IsRead:         false,
			}
			s.msgRepo.Create(welcomeMsg)
			s.convRepo.UpdateLastMessage(conv.ID, welcomeMsg.Content)

			// Generate response to user's message
			aiResponse := s.aiAgent.GenerateResponse(req.Message, conv)
			responseMsg := &models.Message{
				ConversationID: conv.ID,
				SenderID:       aiAgent.ID,
				SenderName:     aiAgent.Name,
				SenderType:     models.MessageTypeAgent,
				Content:        aiResponse,
				ContentType:    "text",
				IsRead:         false,
			}
			s.msgRepo.Create(responseMsg)
			s.convRepo.UpdateLastMessage(conv.ID, aiResponse)

			return conv, responseMsg, nil
		}
	} else {
		// Human agent - set to pending
		s.convRepo.UpdateStatus(conv.ID, models.ConversationStatusPending)
		
		// Publish event for admin notification - human agent requested
		if s.publisher != nil {
			s.publisher.PublishTicketCreated(conv.ID, userID, userName, req.Subject)
		}
	}

	return conv, userMsg, nil
}

// SendMessage sends a message in a conversation
func (s *SupportService) SendMessage(conversationID, senderID, senderName string, senderType models.MessageType, req *models.SendMessageRequest) (*models.Message, *models.Message, error) {
	// Get conversation
	conv, err := s.convRepo.GetByID(conversationID)
	if err != nil {
		return nil, nil, fmt.Errorf("conversation not found: %w", err)
	}

	// Create the message
	msg := &models.Message{
		ConversationID: conversationID,
		SenderID:       senderID,
		SenderName:     senderName,
		SenderType:     senderType,
		Content:        req.Content,
		ContentType:    "text",
		Attachments:    req.Attachments,
		IsRead:         false,
	}

	// Ensure attachments is not nil
	if msg.Attachments == nil {
		msg.Attachments = []string{}
	}

	if req.ContentType != "" {
		msg.ContentType = req.ContentType
	}

	if err := s.msgRepo.Create(msg); err != nil {
		return nil, nil, fmt.Errorf("failed to create message: %w", err)
	}

	// Update conversation
	s.convRepo.UpdateLastMessage(conversationID, req.Content)

	// If AI agent and message is from user, generate response
	var aiResponseMsg *models.Message
	if conv.AgentType == models.AgentTypeAI && senderType == models.MessageTypeUser {
		// Check if should escalate
		messages, _ := s.msgRepo.GetByConversationID(conversationID, 100, 0)
		shouldEscalate, reason := s.aiAgent.ShouldEscalate(req.Content, len(messages))
		
		if shouldEscalate {
			// Escalate to human
			s.EscalateConversation(conversationID, reason)
			
			// Send system message
			systemMsg := &models.Message{
				ConversationID: conversationID,
				SenderID:       "system",
				SenderName:     "Système",
				SenderType:     models.MessageTypeSystem,
				Content:        "Votre demande a été transférée à un conseiller humain. Un agent va prendre en charge votre conversation sous peu. Temps d'attente estimé : 2-5 minutes.",
				ContentType:    "text",
				IsRead:         false,
			}
			s.msgRepo.Create(systemMsg)
			aiResponseMsg = systemMsg
		} else {
			// Generate AI response
			aiAgent, _ := s.agentRepo.GetAIAgent()
			if aiAgent != nil {
				aiResponse := s.aiAgent.GenerateResponse(req.Content, conv)
				aiResponseMsg = &models.Message{
					ConversationID: conversationID,
					SenderID:       aiAgent.ID,
					SenderName:     aiAgent.Name,
					SenderType:     models.MessageTypeAgent,
					Content:        aiResponse,
					ContentType:    "text",
					IsRead:         false,
				}
				s.msgRepo.Create(aiResponseMsg)
				s.convRepo.UpdateLastMessage(conversationID, aiResponse)
			}
		}
	}
	// If human agent conversation, publish message event for admin notification
	if conv.AgentType == models.AgentTypeHuman && senderType == models.MessageTypeUser {
		if s.publisher != nil {
			s.publisher.PublishUserMessage(conversationID, senderID, senderName, req.Content)
		}
	}

	return msg, aiResponseMsg, nil
}

// GetConversation returns a conversation with its messages
func (s *SupportService) GetConversation(id string) (*models.Conversation, error) {
	conv, err := s.convRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	// Get messages
	messages, _ := s.msgRepo.GetByConversationID(id, 100, 0)
	conv.Messages = make([]models.Message, len(messages))
	for i, m := range messages {
		conv.Messages[i] = *m
	}

	// Get agent info
	if conv.AgentID != nil {
		agent, _ := s.agentRepo.GetByID(*conv.AgentID)
		conv.Agent = agent
	}

	return conv, nil
}

// GetUserConversations returns all conversations for a user
func (s *SupportService) GetUserConversations(userID string, limit, offset int) ([]*models.Conversation, error) {
	return s.convRepo.GetByUserID(userID, limit, offset)
}

// GetMessages returns messages for a conversation
func (s *SupportService) GetMessages(conversationID string, limit, offset int) ([]*models.Message, error) {
	return s.msgRepo.GetByConversationID(conversationID, limit, offset)
}

// EscalateConversation escalates a conversation to human support
func (s *SupportService) EscalateConversation(id, reason string) error {
	// Get conversation info for notification
	conv, _ := s.convRepo.GetByID(id)
	
	// Update status
	if err := s.convRepo.UpdateStatus(id, models.ConversationStatusEscalated); err != nil {
		return err
	}

	// Find available human agent
	agents, err := s.agentRepo.GetAvailable(models.AgentTypeHuman)
	if err == nil && len(agents) > 0 {
		// Assign to first available agent
		s.convRepo.AssignAgent(id, agents[0].ID)
		s.agentRepo.IncrementActiveChats(agents[0].ID)
	}

	// Publish escalation event for admin notification
	if s.publisher != nil && conv != nil {
		s.publisher.PublishEscalation(id, conv.UserID, conv.UserName)
	}

	return nil
}

// CloseConversation closes a conversation with optional rating
func (s *SupportService) CloseConversation(id string, rating int, feedback string) error {
	conv, err := s.convRepo.GetByID(id)
	if err != nil {
		return err
	}

	// Decrement agent active chats
	if conv.AgentID != nil {
		s.agentRepo.DecrementActiveChats(*conv.AgentID)
	}

	return s.convRepo.Close(id, rating, feedback)
}

// GetStats returns support statistics
func (s *SupportService) GetStats() (*models.SupportStats, error) {
	stats, err := s.convRepo.GetStats()
	if err != nil {
		return nil, err
	}

	// Get active agents count
	agents, _ := s.agentRepo.GetAll()
	activeCount := 0
	for _, agent := range agents {
		if agent.IsAvailable && agent.Type == models.AgentTypeHuman {
			activeCount++
		}
	}
	stats.ActiveAgents = activeCount

	return stats, nil
}

// GetAllConversations returns all conversations (for admin)
func (s *SupportService) GetAllConversations(status string, limit, offset int) ([]*models.Conversation, error) {
	return s.convRepo.GetAll(status, limit, offset)
}

// AssignAgent assigns an agent to a conversation
func (s *SupportService) AssignAgent(conversationID, agentID string) error {
	if err := s.convRepo.AssignAgent(conversationID, agentID); err != nil {
		return err
	}
	return s.agentRepo.IncrementActiveChats(agentID)
}

// MarkMessagesAsRead marks all messages as read for a user
func (s *SupportService) MarkMessagesAsRead(conversationID, userID string) error {
	return s.msgRepo.MarkAsRead(conversationID, userID)
}

// GetAgents returns all agents
func (s *SupportService) GetAgents() ([]*models.Agent, error) {
	return s.agentRepo.GetAll()
}

// UpdateAgentAvailability updates an agent's availability
func (s *SupportService) UpdateAgentAvailability(agentID string, available bool) error {
	return s.agentRepo.UpdateAvailability(agentID, available)
}

// UploadAttachment uploads an attachment file
func (s *SupportService) UploadAttachment(ctx context.Context, filePath string, filename string, contentType string) (string, error) {
	if s.storage == nil {
		return "", fmt.Errorf("storage service not configured")
	}
	
	// Generate unique object name: date/timestamp_filename
	objectName := fmt.Sprintf("%s/%d_%s", time.Now().Format("2006-01-02"), time.Now().Unix(), filename)
	
	// Upload
	path, err := s.storage.UploadFile(ctx, objectName, filePath, contentType)
	if err != nil {
		return "", err
	}
	
	// Return full public URL
	fullURL := fmt.Sprintf("%s/%s", s.config.MinIOPublicURL, path)
	return fullURL, nil
}

// GetResponseTime calculates average response time
func (s *SupportService) GetResponseTime() time.Duration {
	// Placeholder - in production this would be calculated from actual data
	return 3 * time.Minute
}
