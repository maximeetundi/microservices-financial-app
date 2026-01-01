import '../api/api_client.dart';
import '../api/api_endpoints.dart';

class TicketApiService {
  final ApiClient _client = ApiClient();

  // === Event Management (Organizer) ===
  
  /// Create a new event
  Future<Map<String, dynamic>> createEvent(Map<String, dynamic> data) async {
    final response = await _client.post(
      '/ticket-service/api/v1/events',
      data: data,
    );
    if (response.statusCode == 201 || response.statusCode == 200) {
      return response.data;
    }
    throw Exception(response.data['error'] ?? 'Failed to create event');
  }

  /// Get my events (as organizer)
  Future<List<dynamic>> getMyEvents({int limit = 20, int offset = 0}) async {
    final response = await _client.get(
      '/ticket-service/api/v1/events',
      queryParameters: {'limit': limit, 'offset': offset},
    );
    if (response.statusCode == 200) {
      return response.data['events'] ?? [];
    }
    throw Exception('Failed to load events');
  }

  /// Get event by ID
  Future<Map<String, dynamic>> getEvent(String id) async {
    final response = await _client.get('/ticket-service/api/v1/events/$id');
    if (response.statusCode == 200) {
      return response.data['event'] ?? response.data;
    }
    throw Exception('Event not found');
  }

  /// Get event by code
  Future<Map<String, dynamic>> getEventByCode(String code) async {
    final response = await _client.get('/ticket-service/api/v1/events/code/$code');
    if (response.statusCode == 200) {
      return response.data['event'] ?? response.data;
    }
    throw Exception('Event not found');
  }

  /// Update event
  Future<Map<String, dynamic>> updateEvent(String id, Map<String, dynamic> data) async {
    final response = await _client.put(
      '/ticket-service/api/v1/events/$id',
      data: data,
    );
    if (response.statusCode == 200) {
      return response.data['event'] ?? response.data;
    }
    throw Exception(response.data['error'] ?? 'Failed to update event');
  }

  /// Publish event
  Future<void> publishEvent(String id) async {
    final response = await _client.post('/ticket-service/api/v1/events/$id/publish');
    if (response.statusCode != 200) {
      throw Exception(response.data['error'] ?? 'Failed to publish event');
    }
  }

  /// Delete event
  Future<void> deleteEvent(String id) async {
    final response = await _client.delete('/ticket-service/api/v1/events/$id');
    if (response.statusCode != 200) {
      throw Exception(response.data['error'] ?? 'Failed to delete event');
    }
  }

  /// Get event statistics
  Future<Map<String, dynamic>> getEventStats(String id) async {
    final response = await _client.get('/ticket-service/api/v1/events/$id/stats');
    if (response.statusCode == 200) {
      return response.data['stats'] ?? response.data;
    }
    throw Exception('Failed to load stats');
  }

  /// Get sold tickets for event
  Future<List<dynamic>> getEventTickets(String id, {int limit = 50, int offset = 0}) async {
    final response = await _client.get(
      '/ticket-service/api/v1/events/$id/tickets',
      queryParameters: {'limit': limit, 'offset': offset},
    );
    if (response.statusCode == 200) {
      return response.data['tickets'] ?? [];
    }
    throw Exception('Failed to load tickets');
  }

  // === Public Events ===

  /// Get active events (for browsing)
  Future<List<dynamic>> getActiveEvents({int limit = 20, int offset = 0}) async {
    final response = await _client.get(
      '/ticket-service/api/v1/events/active',
      queryParameters: {'limit': limit, 'offset': offset},
    );
    if (response.statusCode == 200) {
      return response.data['events'] ?? [];
    }
    throw Exception('Failed to load events');
  }

  // === Ticket Purchase ===

  /// Purchase a ticket
  Future<Map<String, dynamic>> purchaseTicket({
    required String eventId,
    required String tierId,
    required Map<String, String> formData,
    required String walletId,
    required String pin,
  }) async {
    final response = await _client.post(
      '/ticket-service/api/v1/tickets/purchase',
      data: {
        'event_id': eventId,
        'tier_id': tierId,
        'form_data': formData,
        'wallet_id': walletId,
        'pin': pin,
      },
    );
    if (response.statusCode == 201 || response.statusCode == 200) {
      return response.data;
    }
    throw Exception(response.data['error'] ?? 'Failed to purchase ticket');
  }

  /// Get my purchased tickets
  Future<List<dynamic>> getMyTickets({int limit = 20, int offset = 0}) async {
    final response = await _client.get(
      '/ticket-service/api/v1/tickets',
      queryParameters: {'limit': limit, 'offset': offset},
    );
    if (response.statusCode == 200) {
      return response.data['tickets'] ?? [];
    }
    throw Exception('Failed to load tickets');
  }

  /// Get ticket by ID
  Future<Map<String, dynamic>> getTicket(String id) async {
    final response = await _client.get('/ticket-service/api/v1/tickets/$id');
    if (response.statusCode == 200) {
      return response.data['ticket'] ?? response.data;
    }
    throw Exception('Ticket not found');
  }

  // === Ticket Verification ===

  /// Verify a ticket by code
  Future<Map<String, dynamic>> verifyTicket(String code) async {
    final response = await _client.post(
      '/ticket-service/api/v1/tickets/verify',
      data: {'ticket_code': code},
    );
    if (response.statusCode == 200) {
      return response.data;
    }
    throw Exception(response.data['error'] ?? 'Failed to verify ticket');
  }

  /// Mark ticket as used
  Future<void> useTicket(String id) async {
    final response = await _client.post('/ticket-service/api/v1/tickets/$id/use');
    if (response.statusCode != 200) {
      throw Exception(response.data['error'] ?? 'Failed to use ticket');
    }
  }

  // === Utilities ===

  /// Get available icons for tiers
  Future<List<String>> getAvailableIcons() async {
    try {
      final response = await _client.get('/ticket-service/api/v1/icons');
      if (response.statusCode == 200) {
        return List<String>.from(response.data['icons'] ?? []);
      }
    } catch (e) {
      // Return default icons
    }
    return [
      '⭐', '🌟', '✨', '💎', '👑', '🏆', '🎖️', '🥇', '🥈', '🥉',
      '🎫', '🎟️', '🎪', '🎭', '🎬', '🎵', '🎸', '🎤', '🎧', '🎹',
      '🔥', '💫', '⚡', '🌈', '🎯', '🚀', '💥', '🎉', '🎊', '🎁',
    ];
  }
}
