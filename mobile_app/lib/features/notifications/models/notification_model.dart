class NotificationModel {
  final String id;
  final String userId;
  final String type;
  final String title;
  final String message;
  final String? data;
  final bool isRead;
  final DateTime? readAt;
  final DateTime createdAt;

  NotificationModel({
    required this.id,
    required this.userId,
    required this.type,
    required this.title,
    required this.message,
    this.data,
    required this.isRead,
    this.readAt,
    required this.createdAt,
  });

  factory NotificationModel.fromJson(Map<String, dynamic> json) {
    return NotificationModel(
      id: json['id'] ?? '',
      userId: json['user_id'] ?? '',
      type: json['type'] ?? 'default',
      title: json['title'] ?? '',
      message: json['message'] ?? '',
      data: json['data'],
      isRead: json['is_read'] ?? false,
      readAt: json['read_at'] != null ? DateTime.parse(json['read_at']) : null,
      createdAt: json['created_at'] != null 
          ? DateTime.parse(json['created_at']) 
          : DateTime.now(),
    );
  }

  Map<String, dynamic> toJson() {
    return {
      'id': id,
      'user_id': userId,
      'type': type,
      'title': title,
      'message': message,
      'data': data,
      'is_read': isRead,
      'read_at': readAt?.toIso8601String(),
      'created_at': createdAt.toIso8601String(),
    };
  }

  /// Get icon based on notification type
  String get icon {
    switch (type) {
      case 'transfer':
        return '💸';
      case 'card':
        return '💳';
      case 'security':
        return '🔒';
      case 'promotion':
        return '🎁';
      case 'wallet':
        return '💰';
      default:
        return '🔔';
    }
  }
}
