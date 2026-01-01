import 'dart:async';
import 'package:flutter/material.dart';
import '../../../core/services/notification_api_service.dart';
import '../screens/notifications_screen.dart';

class NotificationBadge extends StatefulWidget {
  const NotificationBadge({super.key});

  @override
  State<NotificationBadge> createState() => _NotificationBadgeState();
}

class _NotificationBadgeState extends State<NotificationBadge> {
  late final NotificationApiService _notificationService;
  int _unreadCount = 0;
  Timer? _pollTimer;

  @override
  void initState() {
    super.initState();
    _notificationService = NotificationApiService();
    _fetchUnreadCount();
    // Poll every 30 seconds
    _pollTimer = Timer.periodic(
      const Duration(seconds: 30),
      (_) => _fetchUnreadCount(),
    );
  }

  @override
  void dispose() {
    _pollTimer?.cancel();
    super.dispose();
  }

  Future<void> _fetchUnreadCount() async {
    try {
      final response = await _notificationService.getUnreadCount();
      int count = response['unread_count'] ?? 0;
      
      // If user is on support screen, exclude support-related notifications from count
      if (NotificationApiService.isOnSupportScreen && count > 0) {
        // Get notifications to filter out support ones
        final notifResponse = await _notificationService.getNotifications(limit: 20);
        final notifications = (notifResponse['notifications'] ?? []) as List;
        final supportCount = notifications.where((n) {
          if (n['is_read'] == true) return false;
          final type = (n['type'] ?? '').toString().toLowerCase();
          return type == 'support' || type == 'conversation' || type == 'ticket' || type == 'message';
        }).length;
        count = (count - supportCount).clamp(0, count);
      }
      
      if (mounted) {
        setState(() {
          _unreadCount = count;
        });
      }
    } catch (e) {
      // Silently fail
    }
  }

  @override
  Widget build(BuildContext context) {
    return IconButton(
      icon: Stack(
        clipBehavior: Clip.none,
        children: [
          const Icon(Icons.notifications_outlined),
          if (_unreadCount > 0)
            Positioned(
              right: -6,
              top: -6,
              child: Container(
                padding: const EdgeInsets.all(4),
                decoration: const BoxDecoration(
                  color: Colors.red,
                  shape: BoxShape.circle,
                ),
                constraints: const BoxConstraints(
                  minWidth: 18,
                  minHeight: 18,
                ),
                child: Text(
                  _unreadCount > 99 ? '99+' : '$_unreadCount',
                  style: const TextStyle(
                    color: Colors.white,
                    fontSize: 10,
                    fontWeight: FontWeight.bold,
                  ),
                  textAlign: TextAlign.center,
                ),
              ),
            ),
        ],
      ),
      onPressed: () {
        Navigator.push(
          context,
          MaterialPageRoute(
            builder: (_) => const NotificationsScreen(),
          ),
        ).then((_) => _fetchUnreadCount()); // Refresh count when returning
      },
    );
  }
}
