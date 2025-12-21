import '../api/api_client.dart';
import '../api/api_endpoints.dart';

/// Service de support client
class SupportApiService {
  final ApiClient _client = ApiClient();
  
  /// Récupérer tous les tickets de support de l'utilisateur
  Future<List<Map<String, dynamic>>> getTickets({
    int limit = 50,
    int offset = 0,
  }) async {
    final response = await _client.get(
      ApiEndpoints.supportTickets,
      queryParameters: {
        'limit': limit,
        'offset': offset,
      },
    );
    if (response.statusCode == 200) {
      return List<Map<String, dynamic>>.from(response.data['tickets'] ?? []);
    }
    throw Exception('Failed to get support tickets');
  }
  
  /// Créer un nouveau ticket de support
  Future<Map<String, dynamic>> createTicket({
    required String subject,
    required String category,
    required String description,
    String priority = 'normal',
  }) async {
    final response = await _client.post(
      ApiEndpoints.createTicket,
      data: {
        'subject': subject,
        'category': category,
        'description': description,
        'priority': priority,
      },
    );
    if (response.statusCode == 201 || response.statusCode == 200) {
      return response.data['ticket'] ?? response.data;
    }
    throw Exception(response.data['error'] ?? 'Failed to create ticket');
  }
  
  /// Récupérer un ticket par ID
  Future<Map<String, dynamic>> getTicket(String ticketId) async {
    final response = await _client.get(ApiEndpoints.ticketById(ticketId));
    if (response.statusCode == 200) {
      return response.data['ticket'] ?? response.data;
    }
    throw Exception('Ticket not found');
  }
  
  /// Récupérer les messages d'un ticket
  Future<List<Map<String, dynamic>>> getMessages(String ticketId) async {
    final response = await _client.get(ApiEndpoints.ticketMessages(ticketId));
    if (response.statusCode == 200) {
      return List<Map<String, dynamic>>.from(response.data['messages'] ?? []);
    }
    throw Exception('Failed to get messages');
  }
  
  /// Envoyer un message dans un ticket
  Future<Map<String, dynamic>> sendMessage({
    required String ticketId,
    required String content,
  }) async {
    final response = await _client.post(
      ApiEndpoints.sendMessage(ticketId),
      data: {'content': content},
    );
    if (response.statusCode == 201 || response.statusCode == 200) {
      return response.data['message'] ?? response.data;
    }
    throw Exception(response.data['error'] ?? 'Failed to send message');
  }
  
  /// Fermer un ticket
  Future<void> closeTicket(String ticketId) async {
    final response = await _client.post(ApiEndpoints.closeTicket(ticketId));
    if (response.statusCode != 200) {
      throw Exception(response.data['error'] ?? 'Failed to close ticket');
    }
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
