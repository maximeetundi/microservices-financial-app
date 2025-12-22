import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';

class NotificationsPage extends StatefulWidget {
  const NotificationsPage({super.key});

  @override
  State<NotificationsPage> createState() => _NotificationsPageState();
}

class _NotificationsPageState extends State<NotificationsPage> {
  final List<_NotificationItem> _notifications = [
    _NotificationItem(
      id: '1',
      type: 'transfer',
      title: 'Transfert reçu',
      message: 'Vous avez reçu 500 USD de John Doe',
      time: DateTime.now().subtract(const Duration(minutes: 5)),
      isRead: false,
    ),
    _NotificationItem(
      id: '2',
      type: 'card',
      title: 'Paiement carte',
      message: 'Paiement de 45.99 USD chez Amazon',
      time: DateTime.now().subtract(const Duration(hours: 2)),
      isRead: false,
    ),
    _NotificationItem(
      id: '3',
      type: 'security',
      title: 'Nouvelle connexion',
      message: 'Connexion détectée depuis Paris, France',
      time: DateTime.now().subtract(const Duration(days: 1)),
      isRead: true,
    ),
    _NotificationItem(
      id: '4',
      type: 'promotion',
      title: 'Offre spéciale',
      message: 'Profitez de 0% de frais sur vos transferts!',
      time: DateTime.now().subtract(const Duration(days: 2)),
      isRead: true,
    ),
  ];

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: const Text('Notifications'),
        leading: IconButton(
          icon: const Icon(Icons.arrow_back),
          onPressed: () => context.pop(),
        ),
        actions: [
          if (_notifications.any((n) => !n.isRead))
            TextButton(
              onPressed: _markAllAsRead,
              child: const Text('Tout marquer lu'),
            ),
        ],
      ),
      body: _notifications.isEmpty
          ? _buildEmptyState()
          : ListView.separated(
              itemCount: _notifications.length,
              separatorBuilder: (_, __) => const Divider(height: 1),
              itemBuilder: (context, index) {
                final notif = _notifications[index];
                return _buildNotificationTile(notif);
              },
            ),
    );
  }

  Widget _buildEmptyState() {
    return const Center(
      child: Column(
        mainAxisAlignment: MainAxisAlignment.center,
        children: [
          Text('🔔', style: TextStyle(fontSize: 64)),
          SizedBox(height: 16),
          Text('Aucune notification', style: TextStyle(fontSize: 18, color: Colors.grey)),
        ],
      ),
    );
  }

  Widget _buildNotificationTile(_NotificationItem notif) {
    return Dismissible(
      key: Key(notif.id),
      direction: DismissDirection.endToStart,
      background: Container(
        color: Colors.red,
        alignment: Alignment.centerRight,
        padding: const EdgeInsets.only(right: 16),
        child: const Icon(Icons.delete, color: Colors.white),
      ),
      onDismissed: (_) {
        setState(() => _notifications.remove(notif));
      },
      child: ListTile(
        onTap: () => _markAsRead(notif),
        tileColor: notif.isRead ? null : Theme.of(context).primaryColor.withOpacity(0.05),
        leading: Container(
          width: 48,
          height: 48,
          decoration: BoxDecoration(
            color: _getTypeColor(notif.type).withOpacity(0.2),
            borderRadius: BorderRadius.circular(12),
          ),
          child: Center(child: Text(_getTypeIcon(notif.type), style: const TextStyle(fontSize: 24))),
        ),
        title: Text(
          notif.title,
          style: TextStyle(fontWeight: notif.isRead ? FontWeight.normal : FontWeight.bold),
        ),
        subtitle: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Text(notif.message, maxLines: 2, overflow: TextOverflow.ellipsis),
            const SizedBox(height: 4),
            Text(_formatTime(notif.time), style: const TextStyle(fontSize: 12, color: Colors.grey)),
          ],
        ),
        trailing: notif.isRead ? null : Container(
          width: 8,
          height: 8,
          decoration: BoxDecoration(
            color: Theme.of(context).primaryColor,
            shape: BoxShape.circle,
          ),
        ),
      ),
    );
  }

  void _markAsRead(_NotificationItem notif) {
    if (!notif.isRead) {
      setState(() => notif.isRead = true);
    }
  }

  void _markAllAsRead() {
    setState(() {
      for (var n in _notifications) {
        n.isRead = true;
      }
    });
  }

  String _getTypeIcon(String type) {
    switch (type) {
      case 'transfer': return '💸';
      case 'card': return '💳';
      case 'security': return '🔒';
      case 'promotion': return '🎁';
      default: return '🔔';
    }
  }

  Color _getTypeColor(String type) {
    switch (type) {
      case 'transfer': return Colors.green;
      case 'card': return Colors.purple;
      case 'security': return Colors.red;
      case 'promotion': return Colors.amber;
      default: return Colors.grey;
    }
  }

  String _formatTime(DateTime time) {
    final diff = DateTime.now().difference(time);
    if (diff.inMinutes < 60) return 'Il y a ${diff.inMinutes} min';
    if (diff.inHours < 24) return 'Il y a ${diff.inHours}h';
    return 'Il y a ${diff.inDays} jour(s)';
  }
}

class _NotificationItem {
  final String id;
  final String type;
  final String title;
  final String message;
  final DateTime time;
  bool isRead;

  _NotificationItem({
    required this.id,
    required this.type,
    required this.title,
    required this.message,
    required this.time,
    required this.isRead,
  });
}
