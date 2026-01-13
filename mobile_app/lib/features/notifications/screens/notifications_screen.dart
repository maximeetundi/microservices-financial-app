import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';
import '../../../core/services/notification_api_service.dart';
import '../models/notification_model.dart';
import '../../support/support_screen.dart';

class NotificationsScreen extends StatefulWidget {
  const NotificationsScreen({super.key});

  @override
  State<NotificationsScreen> createState() => _NotificationsScreenState();
}

class _NotificationsScreenState extends State<NotificationsScreen> {
  late final NotificationApiService _notificationService;
  List<NotificationModel> _notifications = [];
  bool _isLoading = true;
  String? _error;
  String _activeFilter = 'all';

  // Filters matching web frontend
  static const List<Map<String, String>> _filters = [
    {'id': 'all', 'label': 'Toutes', 'icon': '📋'},
    {'id': 'transfer', 'label': 'Transferts', 'icon': '💸'},
    {'id': 'security', 'label': 'Sécurité', 'icon': '🔐'},
    {'id': 'card', 'label': 'Cartes', 'icon': '💳'},
  ];

  List<NotificationModel> get _filteredNotifications {
    if (_activeFilter == 'all') return _notifications;
    return _notifications.where((n) => n.type == _activeFilter).toList();
  }

  @override
  void initState() {
    super.initState();
    _notificationService = NotificationApiService();
    _loadNotifications();
  }

  Future<void> _loadNotifications() async {
    setState(() {
      _isLoading = true;
      _error = null;
    });

    try {
      final response = await _notificationService.getNotifications();
      final list = response['notifications'] as List? ?? [];
      setState(() {
        _notifications = list.map((n) => NotificationModel.fromJson(n)).toList();
        _isLoading = false;
      });
    } catch (e) {
      setState(() {
        _error = 'Impossible de charger les notifications';
        _isLoading = false;
      });
    }
  }

  Future<void> _handleNotificationTap(NotificationModel notif) async {
    await _markAsRead(notif);
    _navigateToRelatedScreen(notif);
  }

  void _navigateToRelatedScreen(NotificationModel notif) {
    // 1. Prioritize Action URL
    if (notif.actionUrl != null) {
      final url = notif.actionUrl!;
      
      if (url.startsWith('/transactions/')) {
        final id = url.split('/').last;
        // Try to open transfer detail
        // Note: We assume it's a transfer for now as 'transactions' usually implies wallet history
        // If we had a generic transaction detail, we'd use that. 
        // For now, let's open the wallet page which lists transactions.
        context.pushNamed('wallet');
        return;
      }
      
      if (url.startsWith('/tickets/')) {
        final eventId = notif.data?['event_id']?.toString();
        if (eventId != null) {
           context.pushNamed('event-detail', pathParameters: {'eventId': eventId});
        } else {
           context.pushNamed('dashboard');
        }
        return;
      }

      if (url.startsWith('/payment-requests/')) {
          context.pushNamed('wallet');
          return;
      }
      
      // Fallback
      context.pushNamed('dashboard');
      return;
    }

    // 2. Fallback to Type-based navigation
    final type = notif.type.toLowerCase();
    final refId = notif.data?['reference_id']?.toString() ?? 
                  notif.data?['id']?.toString() ?? 
                  notif.data?['ticket_id']?.toString();

    if (type == 'support' || type == 'conversation' || type == 'ticket') {
      Navigator.push(
        context,
        MaterialPageRoute(
          builder: (_) => ChatScreen(agentType: 'human', ticketId: refId),
        ),
      );
    } else if (type == 'transfer' || type == 'transaction') {
      context.pushNamed('wallet');
    } else if (type == 'card') {
      context.pushNamed('cards');
    } else if (type == 'security' || type == 'kyc') {
      context.pushNamed('settings');
    } else if (type == 'wallet' || type == 'payment') {
      context.pushNamed('wallet');
    } else {
       context.pushNamed('dashboard');
    }
  }

  Future<void> _markAsRead(NotificationModel notif) async {
    try {
      await _notificationService.markAsRead(notif.id);
      setState(() {
        final index = _notifications.indexWhere((n) => n.id == notif.id);
        if (index != -1) {
          _notifications[index] = NotificationModel(
            id: notif.id,
            userId: notif.userId,
            type: notif.type,
            title: notif.title,
            message: notif.message,
            data: notif.data,
            isRead: true,
            readAt: DateTime.now(),
            createdAt: notif.createdAt,
          );
        }
      });
    } catch (e) {
      // Silently fail
    }
  }

  Future<void> _markAllAsRead() async {
    try {
      await _notificationService.markAllAsRead();
      setState(() {
        _notifications = _notifications.map((n) => NotificationModel(
          id: n.id,
          userId: n.userId,
          type: n.type,
          title: n.title,
          message: n.message,
          data: n.data,
          isRead: true,
          readAt: DateTime.now(),
          createdAt: n.createdAt,
        )).toList();
      });
    } catch (e) {
      ScaffoldMessenger.of(context).showSnackBar(
        const SnackBar(content: Text('Erreur lors du marquage')),
      );
    }
  }

  Future<void> _deleteNotification(NotificationModel notif) async {
    try {
      await _notificationService.deleteNotification(notif.id);
      setState(() {
        _notifications.removeWhere((n) => n.id == notif.id);
      });
    } catch (e) {
      ScaffoldMessenger.of(context).showSnackBar(
        const SnackBar(content: Text('Erreur lors de la suppression')),
      );
    }
  }

  String _formatTime(DateTime date) {
    final now = DateTime.now();
    final diff = now.difference(date);

    if (diff.inMinutes < 1) return 'À l\'instant';
    if (diff.inMinutes < 60) return 'Il y a ${diff.inMinutes} min';
    if (diff.inHours < 24) return 'Il y a ${diff.inHours} h';
    return '${date.day}/${date.month}';
  }

  Color _getTypeColor(String type) {
    switch (type) {
      case 'transfer':
        return Colors.green;
      case 'card':
        return Colors.purple;
      case 'security':
        return Colors.red;
      case 'promotion':
        return Colors.amber;
      case 'wallet':
        return Colors.blue;
      default:
        return Colors.grey;
    }
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: const Text('Notifications'),
        actions: [
          if (_notifications.any((n) => !n.isRead))
            TextButton(
              onPressed: _markAllAsRead,
              child: const Text('Tout marquer lu'),
            ),
        ],
      ),
      body: _buildBody(),
    );
  }

  Widget _buildBody() {
    if (_isLoading) {
      return const Center(child: CircularProgressIndicator());
    }

    if (_error != null) {
      return Center(
        child: Column(
          mainAxisAlignment: MainAxisAlignment.center,
          children: [
            const Icon(Icons.error_outline, size: 48, color: Colors.red),
            const SizedBox(height: 16),
            Text(_error!, style: Theme.of(context).textTheme.titleMedium),
            const SizedBox(height: 16),
            ElevatedButton(
              onPressed: _loadNotifications,
              child: const Text('Réessayer'),
            ),
          ],
        ),
      );
    }

    if (_notifications.isEmpty) {
      return Center(
        child: Column(
          mainAxisAlignment: MainAxisAlignment.center,
          children: [
            const Text('🔔', style: TextStyle(fontSize: 64)),
            const SizedBox(height: 16),
            Text(
              'Aucune notification',
              style: Theme.of(context).textTheme.titleMedium?.copyWith(
                color: Colors.grey,
              ),
            ),
          ],
        ),
      );
    }

    return Column(
      children: [
        // Filter Pills matching web frontend
        _buildFilterPills(),
        // Notifications List
        Expanded(
          child: _filteredNotifications.isEmpty
              ? Center(
                  child: Text(
                    'Aucune notification dans cette catégorie',
                    style: TextStyle(color: Colors.grey[500]),
                  ),
                )
              : RefreshIndicator(
                  onRefresh: _loadNotifications,
                  child: _buildNotificationsList(),
                ),
        ),
      ],
    );
  }

  Widget _buildFilterPills() {
    return Container(
      height: 50,
      padding: const EdgeInsets.symmetric(horizontal: 8),
      child: ListView.builder(
        scrollDirection: Axis.horizontal,
        itemCount: _filters.length,
        itemBuilder: (context, index) {
          final filter = _filters[index];
          final isActive = _activeFilter == filter['id'];
          return Padding(
            padding: const EdgeInsets.symmetric(horizontal: 4, vertical: 8),
            child: GestureDetector(
              onTap: () {
                setState(() {
                  _activeFilter = filter['id']!;
                });
              },
              child: Container(
                padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 8),
                decoration: BoxDecoration(
                  color: isActive
                      ? Theme.of(context).primaryColor
                      : Theme.of(context).primaryColor.withOpacity(0.1),
                  borderRadius: BorderRadius.circular(20),
                  border: Border.all(
                    color: isActive
                        ? Theme.of(context).primaryColor
                        : Colors.transparent,
                  ),
                ),
                child: Row(
                  mainAxisSize: MainAxisSize.min,
                  children: [
                    Text(filter['icon']!, style: const TextStyle(fontSize: 14)),
                    const SizedBox(width: 6),
                    Text(
                      filter['label']!,
                      style: TextStyle(
                        fontSize: 13,
                        fontWeight: FontWeight.w500,
                        color: isActive ? Colors.white : Colors.grey[700],
                      ),
                    ),
                  ],
                ),
              ),
            ),
          );
        },
      ),
    );
  }

  Widget _buildNotificationsList() {
    return ListView.separated(
      itemCount: _filteredNotifications.length,
      separatorBuilder: (_, __) => const Divider(height: 1),
      itemBuilder: (context, index) {
        final notif = _filteredNotifications[index];
          return Dismissible(
            key: Key(notif.id),
            direction: DismissDirection.endToStart,
            background: Container(
              color: Colors.red,
              alignment: Alignment.centerRight,
              padding: const EdgeInsets.only(right: 16),
              child: const Icon(Icons.delete, color: Colors.white),
            ),
            onDismissed: (_) => _deleteNotification(notif),
            child: ListTile(
              onTap: () => _handleNotificationTap(notif),
              tileColor: notif.isRead ? null : Theme.of(context).primaryColor.withOpacity(0.05),
              leading: Container(
                width: 48,
                height: 48,
                decoration: BoxDecoration(
                  color: _getTypeColor(notif.type).withOpacity(0.2),
                  borderRadius: BorderRadius.circular(12),
                ),
                child: Center(
                  child: Text(notif.icon, style: const TextStyle(fontSize: 24)),
                ),
              ),
              title: Text(
                notif.title,
                style: TextStyle(
                  fontWeight: notif.isRead ? FontWeight.normal : FontWeight.bold,
                ),
              ),
              subtitle: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  Text(
                    notif.message,
                    maxLines: 2,
                    overflow: TextOverflow.ellipsis,
                    style: TextStyle(color: Colors.grey[600]),
                  ),
                  const SizedBox(height: 4),
                  Text(
                    _formatTime(notif.createdAt),
                    style: TextStyle(fontSize: 12, color: Colors.grey[400]),
                  ),
                ],
              ),
              trailing: notif.isRead
                  ? null
                  : Container(
                      width: 8,
                      height: 8,
                      decoration: BoxDecoration(
                        color: Theme.of(context).primaryColor,
                        shape: BoxShape.circle,
                      ),
                    ),
            ),
          );
      },
    );
  }
}
