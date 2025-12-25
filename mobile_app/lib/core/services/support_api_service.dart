import '../api/api_client.dart';
import '../api/api_endpoints.dart';

/// Service de support client
class SupportApiService {
  final ApiClient _client = ApiClient();
  
  /// Récupérer toutes les conversations de support de l'utilisateur
  Future<List<Map<String, dynamic>>> getConversations({
    int limit = 50,
    int offset = 0,
  }) async {
    final response = await _client.get(
      ApiEndpoints.supportConversations,
      queryParameters: {
        'limit': limit,
        'offset': offset,
      },
    );
    if (response.statusCode == 200) {
      return List<Map<String, dynamic>>.from(response.data['conversations'] ?? []);
    }
    throw Exception('Failed to get conversations');
  }

  /// Alias for backwards compatibility
  Future<List<Map<String, dynamic>>> getTickets({int limit = 50, int offset = 0}) async {
    return getConversations(limit: limit, offset: offset);
  }
  
  /// Créer une nouvelle conversation de support
  Future<Map<String, dynamic>> createConversation({
    required String subject,
    required String category,
    required String message,
    String agentType = 'ai',
    String priority = 'normal',
  }) async {
    final response = await _client.post(
      ApiEndpoints.createConversation,
      data: {
        'subject': subject,
        'category': category,
        'message': message,
        'agent_type': agentType,
        'priority': priority,
      },
    );
    if (response.statusCode == 201 || response.statusCode == 200) {
      return response.data['conversation'] ?? response.data;
    }
    throw Exception(response.data['error'] ?? 'Failed to create conversation');
  }

  /// Alias for backwards compatibility
  Future<Map<String, dynamic>> createTicket({
    required String subject,
    required String category,
    required String description,
    String priority = 'normal',
  }) async {
    return createConversation(
      subject: subject,
      category: category,
      message: description,
      priority: priority,
    );
  }
  
  /// Récupérer une conversation par ID
  Future<Map<String, dynamic>> getConversation(String conversationId) async {
    final response = await _client.get(ApiEndpoints.conversationById(conversationId));
    if (response.statusCode == 200) {
      return response.data['conversation'] ?? response.data;
    }
    throw Exception('Conversation not found');
  }

  /// Alias for backwards compatibility
  Future<Map<String, dynamic>> getTicket(String ticketId) async {
    return getConversation(ticketId);
  }
  
  /// Récupérer les messages d'une conversation
  Future<List<Map<String, dynamic>>> getMessages(String conversationId) async {
    final response = await _client.get(ApiEndpoints.conversationMessages(conversationId));
    if (response.statusCode == 200) {
      return List<Map<String, dynamic>>.from(response.data['messages'] ?? []);
    }
    throw Exception('Failed to get messages');
  }
  
  /// Envoyer un message dans une conversation
  Future<Map<String, dynamic>> sendMessage({
    required String conversationId,
    required String content,
  }) async {
    final response = await _client.post(
      ApiEndpoints.sendConversationMessage(conversationId),
      data: {'content': content},
    );
    if (response.statusCode == 201 || response.statusCode == 200) {
      return response.data['message'] ?? response.data;
    }
    throw Exception(response.data['error'] ?? 'Failed to send message');
  }
  
  /// Fermer une conversation
  Future<void> closeConversation(String conversationId) async {
    final response = await _client.post(ApiEndpoints.closeConversation(conversationId));
    if (response.statusCode != 200) {
      throw Exception(response.data['error'] ?? 'Failed to close conversation');
    }
  }

  /// Alias for backwards compatibility
  Future<void> closeTicket(String ticketId) async {
    return closeConversation(ticketId);
  }
  
  /// Récupérer les statistiques de support
  Future<Map<String, dynamic>> getStats() async {
    final response = await _client.get(ApiEndpoints.supportStats);
    if (response.statusCode == 200) {
      return response.data;
    }
    throw Exception('Failed to get support stats');
  }
}

