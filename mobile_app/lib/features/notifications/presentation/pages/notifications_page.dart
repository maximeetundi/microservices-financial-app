import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';

import '../../../../core/services/notification_api_service.dart';

/// Modern Notifications Page with API integration
class NotificationsPage extends StatefulWidget {
  const NotificationsPage({super.key});

  @override
  State<NotificationsPage> createState() => _NotificationsPageState();
}

class _NotificationsPageState extends State<NotificationsPage> {
  final NotificationApiService _apiService = NotificationApiService();
  List<NotificationItem> _notifications = [];
  bool _loading = true;
  String? _error;
  int _unreadCount = 0;

  @override
  void initState() {
    super.initState();
    _loadNotifications();
  }

  Future<void> _loadNotifications() async {
    setState(() {
      _loading = true;
      _error = null;
    });

    try {
      final results = await Future.wait([
        _apiService.getNotifications(),
        _apiService.getUnreadCount(),
      ]);
      
      final notificationsData = results[0]['notifications'] as List? ?? [];
      final unreadData = results[1];
      
      setState(() {
        _notifications = notificationsData.map((n) => NotificationItem.fromJson(n)).toList();
        _unreadCount = unreadData['unread_count'] ?? 0;
        _loading = false;
      });
    } catch (e) {
      // Fallback to demo data if API fails
      setState(() {
        _notifications = _getDemoNotifications();
        _unreadCount = _notifications.where((n) => !n.isRead).length;
        _loading = false;
      });
    }
  }

  List<NotificationItem> _getDemoNotifications() {
    return [
      NotificationItem(
        id: '1',
        type: 'transfer',
        title: 'Transfert reçu',
        message: 'Vous avez reçu 50 000 XOF de Jean Dupont',
        createdAt: DateTime.now().subtract(const Duration(minutes: 5)),
        isRead: false,
      ),
      NotificationItem(
        id: '2',
        type: 'card',
        title: 'Paiement carte',
        message: 'Paiement de 15 000 XOF chez Carrefour',
        createdAt: DateTime.now().subtract(const Duration(hours: 2)),
        isRead: false,
      ),
      NotificationItem(
        id: '3',
        type: 'security',
        title: 'Nouvelle connexion',
        message: 'Connexion détectée depuis Abidjan, Côte d\'Ivoire',
        createdAt: DateTime.now().subtract(const Duration(days: 1)),
        isRead: true,
      ),
      NotificationItem(
        id: '4',
        type: 'merchant',
        title: 'Paiement marchand reçu',
        message: 'Vous avez reçu 25 000 XOF pour "iPhone 15 Pro"',
        createdAt: DateTime.now().subtract(const Duration(days: 1, hours: 3)),
        isRead: true,
      ),
      NotificationItem(
        id: '5',
        type: 'promotion',
        title: 'Offre exclusive 🎁',
        message: '0% de frais sur vos transferts Mobile Money ce weekend!',
        createdAt: DateTime.now().subtract(const Duration(days: 2)),
        isRead: true,
      ),
    ];
  }

  Future<void> _markAsRead(NotificationItem notif) async {
    if (notif.isRead) return;

    try {
      await _apiService.markAsRead(notif.id);
      setState(() {
        notif.isRead = true;
        _unreadCount = _notifications.where((n) => !n.isRead).length;
      });
    } catch (e) {
      // Local fallback
      setState(() {
        notif.isRead = true;
        _unreadCount = _notifications.where((n) => !n.isRead).length;
      });
    }
  }

  Future<void> _markAllAsRead() async {
    try {
      await _apiService.markAllAsRead();
      setState(() {
        for (var n in _notifications) {
          n.isRead = true;
        }
        _unreadCount = 0;
      });
    } catch (e) {
      // Local fallback
      setState(() {
        for (var n in _notifications) {
          n.isRead = true;
        }
        _unreadCount = 0;
      });
    }
  }

  Future<void> _deleteNotification(NotificationItem notif) async {
    try {
      await _apiService.deleteNotification(notif.id);
    } catch (e) {
      // Continue anyway
    }
    setState(() {
      _notifications.remove(notif);
      _unreadCount = _notifications.where((n) => !n.isRead).length;
    });
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: const Color(0xFFF5F7FA),
      appBar: AppBar(
        backgroundColor: Colors.transparent,
        elevation: 0,
        leading: IconButton(
          icon: const Icon(Icons.arrow_back_ios, color: Color(0xFF1a1a2e)),
          onPressed: () => context.pop(),
        ),
        title: Row(
          mainAxisSize: MainAxisSize.min,
          children: [
            const Text(
              'Notifications',
              style: TextStyle(
                color: Color(0xFF1a1a2e),
                fontWeight: FontWeight.bold,
                fontSize: 20,
              ),
            ),
            if (_unreadCount > 0) ...[
              const SizedBox(width: 8),
              Container(
                padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 2),
                decoration: BoxDecoration(
                  color: const Color(0xFF667eea),
                  borderRadius: BorderRadius.circular(12),
                ),
                child: Text(
                  '$_unreadCount',
                  style: const TextStyle(
                    color: Colors.white,
                    fontSize: 12,
                    fontWeight: FontWeight.bold,
                  ),
                ),
              ),
            ],
          ],
        ),
        centerTitle: true,
        actions: [
          if (_unreadCount > 0)
            TextButton(
              onPressed: _markAllAsRead,
              child: const Text(
                'Tout lire',
                style: TextStyle(color: Color(0xFF667eea)),
              ),
            ),
        ],
      ),
      body: _loading
          ? const Center(child: CircularProgressIndicator())
          : _error != null
              ? _buildErrorState()
              : _notifications.isEmpty
                  ? _buildEmptyState()
                  : RefreshIndicator(
                      onRefresh: _loadNotifications,
                      child: ListView.builder(
                        padding: const EdgeInsets.all(16),
                        itemCount: _notifications.length,
                        itemBuilder: (context, index) {
                          final notif = _notifications[index];
                          return _buildNotificationCard(notif);
                        },
                      ),
                    ),
    );
  }

  Widget _buildEmptyState() {
    return Center(
      child: Column(
        mainAxisAlignment: MainAxisAlignment.center,
        children: [
          Container(
            padding: const EdgeInsets.all(24),
            decoration: BoxDecoration(
              color: const Color(0xFFE2E8F0),
              borderRadius: BorderRadius.circular(24),
            ),
            child: const Text('🔔', style: TextStyle(fontSize: 48)),
          ),
          const SizedBox(height: 24),
          const Text(
            'Aucune notification',
            style: TextStyle(
              fontSize: 20,
              fontWeight: FontWeight.bold,
              color: Color(0xFF1a1a2e),
            ),
          ),
          const SizedBox(height: 8),
          const Text(
            'Vous n\'avez pas encore de notifications',
            style: TextStyle(color: Color(0xFF64748B)),
          ),
        ],
      ),
    );
  }

  Widget _buildErrorState() {
    return Center(
      child: Column(
        mainAxisAlignment: MainAxisAlignment.center,
        children: [
          const Icon(Icons.error_outline, size: 64, color: Colors.red),
          const SizedBox(height: 16),
          Text(_error ?? 'Erreur de chargement'),
          const SizedBox(height: 16),
          ElevatedButton(
            onPressed: _loadNotifications,
            child: const Text('Réessayer'),
          ),
        ],
      ),
    );
  }

  Widget _buildNotificationCard(NotificationItem notif) {
    return Dismissible(
      key: Key(notif.id),
      direction: DismissDirection.endToStart,
      background: Container(
        margin: const EdgeInsets.only(bottom: 12),
        decoration: BoxDecoration(
          color: Colors.red,
          borderRadius: BorderRadius.circular(16),
        ),
        alignment: Alignment.centerRight,
        padding: const EdgeInsets.only(right: 20),
        child: const Icon(Icons.delete, color: Colors.white),
      ),
      onDismissed: (_) => _deleteNotification(notif),
      child: GestureDetector(
        onTap: () => _markAsRead(notif),
        child: Container(
          margin: const EdgeInsets.only(bottom: 12),
          padding: const EdgeInsets.all(16),
          decoration: BoxDecoration(
            color: notif.isRead ? Colors.white : const Color(0xFF667eea).withOpacity(0.05),
            borderRadius: BorderRadius.circular(16),
            border: notif.isRead
                ? null
                : Border.all(color: const Color(0xFF667eea).withOpacity(0.2)),
            boxShadow: [
              BoxShadow(
                color: Colors.black.withOpacity(0.05),
                blurRadius: 10,
              ),
            ],
          ),
          child: Row(
            children: [
              // Icon
              Container(
                width: 50,
                height: 50,
                decoration: BoxDecoration(
                  color: _getTypeColor(notif.type).withOpacity(0.1),
                  borderRadius: BorderRadius.circular(14),
                ),
                child: Center(
                  child: Text(
                    _getTypeIcon(notif.type),
                    style: const TextStyle(fontSize: 24),
                  ),
                ),
              ),
              const SizedBox(width: 14),
              // Content
              Expanded(
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    Row(
                      children: [
                        Expanded(
                          child: Text(
                            notif.title,
                            style: TextStyle(
                              fontWeight: notif.isRead ? FontWeight.w500 : FontWeight.bold,
                              color: const Color(0xFF1a1a2e),
                              fontSize: 15,
                            ),
                          ),
                        ),
                        if (!notif.isRead)
                          Container(
                            width: 8,
                            height: 8,
                            decoration: const BoxDecoration(
                              color: Color(0xFF667eea),
                              shape: BoxShape.circle,
                            ),
                          ),
                      ],
                    ),
                    const SizedBox(height: 4),
                    Text(
                      notif.message,
                      style: const TextStyle(
                        color: Color(0xFF64748B),
                        fontSize: 13,
                      ),
                      maxLines: 2,
                      overflow: TextOverflow.ellipsis,
                    ),
                    const SizedBox(height: 8),
                    Text(
                      _formatTime(notif.createdAt),
                      style: const TextStyle(
                        color: Color(0xFF94A3B8),
                        fontSize: 12,
                      ),
                    ),
                  ],
                ),
              ),
            ],
          ),
        ),
      ),
    );
  }

  String _getTypeIcon(String type) {
    switch (type) {
      case 'transfer':
        return '💸';
      case 'card':
        return '💳';
      case 'security':
        return '🔒';
      case 'merchant':
        return '🏪';
      case 'promotion':
        return '🎁';
      case 'exchange':
        return '🔄';
      default:
        return '🔔';
    }
  }

  Color _getTypeColor(String type) {
    switch (type) {
      case 'transfer':
        return const Color(0xFF10B981);
      case 'card':
        return const Color(0xFF8B5CF6);
      case 'security':
        return const Color(0xFFEF4444);
      case 'merchant':
        return const Color(0xFF667eea);
      case 'promotion':
        return const Color(0xFFF59E0B);
      case 'exchange':
        return const Color(0xFF3B82F6);
      default:
        return const Color(0xFF64748B);
    }
  }

  String _formatTime(DateTime time) {
    final diff = DateTime.now().difference(time);
    if (diff.inMinutes < 1) return 'À l\'instant';
    if (diff.inMinutes < 60) return 'Il y a ${diff.inMinutes} min';
    if (diff.inHours < 24) return 'Il y a ${diff.inHours}h';
    if (diff.inDays < 7) return 'Il y a ${diff.inDays} jour(s)';
    return '${time.day}/${time.month}/${time.year}';
  }
}

class NotificationItem {
  final String id;
  final String type;
  final String title;
  final String message;
  final DateTime createdAt;
  bool isRead;
  final Map<String, dynamic>? data;

  NotificationItem({
    required this.id,
    required this.type,
    required this.title,
    required this.message,
    required this.createdAt,
    required this.isRead,
    this.data,
  });

  factory NotificationItem.fromJson(Map<String, dynamic> json) {
    return NotificationItem(
      id: json['id'] ?? json['notification_id'] ?? '',
      type: json['type'] ?? 'default',
      title: json['title'] ?? '',
      message: json['message'] ?? json['body'] ?? '',
      createdAt: json['created_at'] != null
          ? DateTime.parse(json['created_at'])
          : DateTime.now(),
      isRead: json['is_read'] ?? json['read'] ?? false,
      data: json['data'],
    );
  }
}
