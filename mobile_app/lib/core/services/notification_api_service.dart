import '../api/api_client.dart';

class NotificationApiService {
  final ApiClient _client = ApiClient();
  
  /// Static flag to track if user is currently on support chat screen
  /// When true, support-related notifications will be excluded from count
  static bool isOnSupportScreen = false;

  /// Get list of notifications
  Future<Map<String, dynamic>> getNotifications({int limit = 20, int offset = 0}) async {
    final response = await _client.get(
      '/notification-service/api/v1/notifications',
      queryParameters: {'limit': limit, 'offset': offset},
    );
    if (response.statusCode == 200) {
      return response.data is Map<String, dynamic> 
          ? response.data 
          : {'notifications': []};
    }
    throw Exception('Failed to load notifications');
  }

  /// Get unread notification count
  Future<Map<String, dynamic>> getUnreadCount() async {
    final response = await _client.get('/notification-service/api/v1/notifications/unread-count');
    if (response.statusCode == 200) {
      return response.data is Map<String, dynamic> 
          ? response.data 
          : {'unread_count': 0};
    }
    return {'unread_count': 0};
  }

  /// Mark a notification as read
  Future<void> markAsRead(String notificationId) async {
    await _client.post('/notification-service/api/v1/notifications/$notificationId/read');
  }

  /// Mark all notifications as read
  Future<void> markAllAsRead() async {
    await _client.post('/notification-service/api/v1/notifications/read-all');
  }

  /// Delete a notification
  Future<void> deleteNotification(String notificationId) async {
    await _client.delete('/notification-service/api/v1/notifications/$notificationId');
  }
}
