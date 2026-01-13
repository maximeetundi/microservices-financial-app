import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';
import 'package:google_fonts/google_fonts.dart';

import '../../../../core/services/notification_api_service.dart';
import '../../../../core/theme/app_theme.dart';
import '../../../../core/widgets/glass_container.dart';

/// Notifications Page matching web design exactly
class NotificationsPage extends StatefulWidget {
  const NotificationsPage({super.key});

  @override
  State<NotificationsPage> createState() => _NotificationsPageState();
}

class _NotificationsPageState extends State<NotificationsPage> {
  final NotificationApiService _apiService = NotificationApiService();
  
  String _activeFilter = 'all';
  bool _loading = true;
  List<Map<String, dynamic>> _notifications = [];
  int _limit = 20;
  int _offset = 0;
  
  final List<Map<String, String>> _filters = [
    {'id': 'all', 'label': 'Toutes', 'icon': '📋'},
    {'id': 'transfer', 'label': 'Transferts', 'icon': '💸'},
    {'id': 'security', 'label': 'Sécurité', 'icon': '🔐'},
    {'id': 'card', 'label': 'Cartes', 'icon': '💳'},
  ];
  
  @override
  void initState() {
    super.initState();
    _fetchNotifications();
  }

  Future<void> _fetchNotifications() async {
    setState(() => _loading = true);
    try {
      final response = await _apiService.getNotifications(limit: _limit, offset: _offset);
      final List<dynamic> notifs = response['notifications'] ?? [];
      setState(() {
        _notifications = notifs.map((n) => Map<String, dynamic>.from(n)).toList();
        _loading = false;
      });
    } catch (e) {
      print('Error fetching notifications: $e');
      setState(() => _loading = false);
    }
  }

  List<Map<String, dynamic>> get _filteredNotifications {
    if (_activeFilter == 'all') return _notifications;
    return _notifications.where((n) => n['type'] == _activeFilter).toList();
  }
  
  int get _unreadCount => _notifications.where((n) => !(n['is_read'] ?? true)).length;

  @override
  Widget build(BuildContext context) {
    // Check if we are in dark mode (simple check using MediaQuery context brightness)
    // Adjust as per your ThemeProvider or app state management
    final isDark = Theme.of(context).brightness == Brightness.dark;
    
    return Scaffold(
      backgroundColor: isDark ? const Color(0xFF0F172A) : const Color(0xFFF8FAFC),
      body: Column(
        children: [
          // Gradient Header
          _buildHeader(isDark),
          
          // Content
          Expanded(
            child: _loading 
              ? Center(child: CircularProgressIndicator(color: const Color(0xFF6366F1)))
              : _notifications.isEmpty
                ? _buildEmptyState(isDark)
                : RefreshIndicator(
                    onRefresh: _refresh,
                    child: ListView(
                      padding: const EdgeInsets.all(16),
                      children: [
                        // Filter Pills
                        _buildFilters(isDark),
                        const SizedBox(height: 16),
                        
                        // Notifications List
                        if (_filteredNotifications.isEmpty)
                           Padding(
                             padding: const EdgeInsets.only(top: 40),
                             child: _buildEmptyState(isDark),
                           )
                        else
                           ..._filteredNotifications.map((n) => _buildNotificationItem(n, isDark)),
                      ],
                    ),
                  ),
          ),
        ],
      ),
    );
  }

  Widget _buildHeader(bool isDark) {
    return Container(
      padding: EdgeInsets.only(
        top: MediaQuery.of(context).padding.top + 16,
        left: 16,
        right: 16,
        bottom: 16,
      ),
      decoration: const BoxDecoration(
        gradient: LinearGradient(
          colors: [Color(0xFF6366F1), Color(0xFF4F46E5)],
        ),
      ),
      child: Row(
        children: [
          GestureDetector(
            onTap: () => context.go('/dashboard'),
            child: Container(
              width: 36,
              height: 36,
              decoration: BoxDecoration(
                color: Colors.white.withOpacity(0.2),
                borderRadius: BorderRadius.circular(10),
              ),
              child: const Icon(Icons.arrow_back_ios_new, color: Colors.white, size: 18),
            ),
          ),
          const SizedBox(width: 12),
          Expanded(
            child: Text(
              '🔔 Notifications',
              style: GoogleFonts.inter(
                fontSize: 18,
                fontWeight: FontWeight.w600,
                color: Colors.white,
              ),
            ),
          ),
          if (_unreadCount > 0) ...[
            Container(
              padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 2),
              decoration: BoxDecoration(
                color: const Color(0xFFEF4444),
                borderRadius: BorderRadius.circular(12),
              ),
              child: Text(
                '$_unreadCount',
                style: GoogleFonts.inter(
                  fontSize: 12,
                  fontWeight: FontWeight.bold,
                  color: Colors.white,
                ),
              ),
            ),
            const SizedBox(width: 8),
            GestureDetector(
              onTap: _markAllAsRead,
              child: Container(
                padding: const EdgeInsets.symmetric(horizontal: 12, vertical: 6),
                decoration: BoxDecoration(
                  color: Colors.white.withOpacity(0.2),
                  borderRadius: BorderRadius.circular(8),
                ),
                child: Text(
                  'Tout lire',
                  style: GoogleFonts.inter(
                    fontSize: 12,
                    color: Colors.white,
                  ),
                ),
              ),
            ),
          ],
          const SizedBox(width: 8),
          GestureDetector(
            onTap: _refresh,
            child: Container(
              width: 36,
              height: 36,
              decoration: BoxDecoration(
                color: Colors.white.withOpacity(0.2),
                borderRadius: BorderRadius.circular(10),
              ),
              child: _loading
                  ? const SizedBox(
                      width: 18,
                      height: 18,
                      child: CircularProgressIndicator(strokeWidth: 2, color: Colors.white),
                    )
                  : const Icon(Icons.refresh, color: Colors.white, size: 18),
            ),
          ),
        ],
      ),
    );
  }

  Widget _buildFilters(bool isDark) {
    return SizedBox(
      height: 36,
      child: ListView.separated(
        scrollDirection: Axis.horizontal,
        itemCount: _filters.length,
        separatorBuilder: (_, __) => const SizedBox(width: 8),
        itemBuilder: (context, index) {
          final filter = _filters[index];
          final isActive = _activeFilter == filter['id'];
          
          return GestureDetector(
            onTap: () => setState(() => _activeFilter = filter['id']!),
            child: AnimatedContainer(
              duration: const Duration(milliseconds: 200),
              padding: const EdgeInsets.symmetric(horizontal: 14, vertical: 8),
              decoration: BoxDecoration(
                color: isActive 
                    ? const Color(0xFF6366F1)
                    : (isDark ? Colors.white.withOpacity(0.05) : Colors.white),
                borderRadius: BorderRadius.circular(20),
                border: Border.all(
                  color: isActive 
                      ? const Color(0xFF6366F1)
                      : (isDark ? Colors.white.withOpacity(0.1) : const Color(0xFFE2E8F0)),
                ),
              ),
              child: Text(
                '${filter['icon']} ${filter['label']}',
                style: GoogleFonts.inter(
                  fontSize: 12,
                  fontWeight: FontWeight.w500,
                  color: isActive 
                      ? Colors.white
                      : (isDark ? const Color(0xFF94A3B8) : const Color(0xFF64748B)),
                ),
              ),
            ),
          );
        },
      ),
    );
  }

  Widget _buildNotificationItem(Map<String, dynamic> notification, bool isDark) {
    final isUnread = !(notification['is_read'] ?? true);
    final String id = notification['id']?.toString() ?? '';
    
    return GestureDetector(
      onTap: () => _handleTap(notification),
      child: Container(
        margin: const EdgeInsets.only(bottom: 8),
        padding: const EdgeInsets.all(14),
        decoration: BoxDecoration(
          color: isUnread 
              ? (isDark ? const Color(0xFF6366F1).withOpacity(0.1) : const Color(0xFFEEF2FF))
              : (isDark ? Colors.white.withOpacity(0.03) : Colors.white),
          borderRadius: BorderRadius.circular(12),
          border: Border.all(
            color: isUnread 
                ? (isDark ? const Color(0xFF6366F1).withOpacity(0.2) : const Color(0xFFC7D2FE))
                : (isDark ? Colors.white.withOpacity(0.08) : const Color(0xFFE2E8F0)),
          ),
        ),
        child: Row(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            // Icon
            Container(
              width: 40,
              height: 40,
              decoration: BoxDecoration(
                color: const Color(0xFF6366F1).withOpacity(0.15),
                borderRadius: BorderRadius.circular(10),
              ),
              child: Center(
                child: Text(
                  _getTypeIcon(notification['type'] ?? 'default'),
                  style: const TextStyle(fontSize: 18),
                ),
              ),
            ),
            const SizedBox(width: 12),
            
            // Content
            Expanded(
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  Text(
                    notification['title'] ?? 'Notification',
                    style: GoogleFonts.inter(
                      fontSize: 14,
                      fontWeight: FontWeight.w500,
                      color: isDark ? Colors.white : const Color(0xFF1E293B),
                    ),
                    maxLines: 1,
                    overflow: TextOverflow.ellipsis,
                  ),
                  const SizedBox(height: 2),
                  Text(
                    notification['message'] ?? '',
                    style: GoogleFonts.inter(
                      fontSize: 12,
                      color: isDark ? const Color(0xFF64748B) : const Color(0xFF94A3B8),
                    ),
                    maxLines: 2,
                    overflow: TextOverflow.ellipsis,
                  ),
                  const SizedBox(height: 4),
                  Text(
                    _formatTime(notification['created_at']),
                    style: GoogleFonts.inter(
                      fontSize: 10,
                      color: isDark ? const Color(0xFF475569) : const Color(0xFFCBD5E1),
                    ),
                  ),
                ],
              ),
            ),
            
            // Actions
            Column(
              children: [
                if (isUnread)
                  GestureDetector(
                    onTap: () => _markAsRead(id),
                    child: Container(
                      width: 28,
                      height: 28,
                      decoration: BoxDecoration(
                        color: const Color(0xFF22C55E).withOpacity(0.2),
                        borderRadius: BorderRadius.circular(6),
                      ),
                      child: const Center(
                        child: Text('✓', style: TextStyle(fontSize: 12, color: Color(0xFF22C55E))),
                      ),
                    ),
                  ),
                const SizedBox(height: 4),
                GestureDetector(
                  onTap: () => _deleteNotification(id),
                  child: Container(
                    width: 28,
                    height: 28,
                    decoration: BoxDecoration(
                      color: const Color(0xFFEF4444).withOpacity(0.2),
                      borderRadius: BorderRadius.circular(6),
                    ),
                    child: const Center(
                      child: Text('✕', style: TextStyle(fontSize: 12, color: Color(0xFFEF4444))),
                    ),
                  ),
                ),
              ],
            ),
          ],
        ),
      ),
    );
  }

  Widget _buildEmptyState(bool isDark) {
    return Center(
      child: Column(
        mainAxisAlignment: MainAxisAlignment.center,
        children: [
          const Text('🔔', style: TextStyle(fontSize: 64)),
          const SizedBox(height: 16),
          Text(
            'Aucune notification',
            style: GoogleFonts.inter(
              fontSize: 16,
              color: isDark ? const Color(0xFF64748B) : const Color(0xFF94A3B8),
            ),
          ),
        ],
      ),
    );
  }

  String _getTypeIcon(String type) {
    switch (type) {
      case 'transfer': return '💸';
      case 'card': return '💳';
      case 'security': return '🔐';
      case 'wallet': return '👛';
      case 'kyc': return '✅';
      case 'payment': return '💰';
      default: return '🔔';
    }
  }

  String _formatTime(dynamic date) {
    if (date == null) return '';
    try {
      final DateTime dt = date is DateTime ? date : DateTime.parse(date.toString());
      final diff = DateTime.now().difference(dt);
      
      if (diff.inMinutes < 1) return 'À l\'instant';
      if (diff.inMinutes < 60) return '${diff.inMinutes} min';
      if (diff.inHours < 24) return '${diff.inHours}h';
      return '${diff.inDays}j';
    } catch (e) {
      return '';
    }
  }

  Future<void> _handleTap(Map<String, dynamic> notification) async {
    final String id = notification['id']?.toString() ?? '';
    final bool isRead = notification['is_read'] ?? false;
    
    if (!isRead) {
      await _markAsRead(id);
    }
    
    // Navigation Logic
    final type = notification['type']?.toString().toLowerCase();
    var data = notification['data'];
    
    String? refId;
    if (notification['reference_id'] != null) {
       refId = notification['reference_id'].toString();
    } else if (data is Map) {
       refId = data['transfer_id']?.toString() ?? data['id']?.toString() ?? data['reference_id']?.toString();
    }
    
    if (type == 'transfer' || type == 'payment' || type == 'transaction') {
        if (refId != null) {
            // Navigate to transfer detail
            context.pushNamed('transfer-detail', pathParameters: {'transferId': refId});
        } else {
            context.go('/transactions');
        }
    } else if (type == 'card') {
        context.goNamed('cards');
    } else if (type == 'security') {
        context.goNamed('security');
    } else if (type == 'kyc') {
        context.goNamed('kyc');
    } else if (type == 'wallet') {
        context.go('/wallet');
    }
  }

  Future<void> _markAsRead(String id) async {
    if (id.isEmpty) return;
    
    // Optimistic UI update
    setState(() {
      final index = _notifications.indexWhere((n) => n['id'].toString() == id);
      if (index != -1) {
        _notifications[index]['is_read'] = true;
      }
    });

    try {
      await _apiService.markAsRead(id);
    } catch (e) {
      // Revert on failure? usually minor enough to ignore or just log
      print('Failed to mark as read: $e');
    }
  }

  Future<void> _markAllAsRead() async {
    // Optimistic UI update
    setState(() {
      for (var n in _notifications) {
        n['is_read'] = true;
      }
    });

    try {
      await _apiService.markAllAsRead();
    } catch (e) {
      print('Failed to mark all as read: $e');
    }
  }

  Future<void> _deleteNotification(String id) async {
    if (id.isEmpty) return;
    
    // Optimistic UI update
    setState(() {
      _notifications.removeWhere((n) => n['id'].toString() == id);
    });

    try {
      await _apiService.deleteNotification(id);
    } catch (e) {
      print('Failed to delete notification: $e');
      // Force refresh on error/mismatch
      _refresh(); 
    }
  }

  Future<void> _refresh() async {
    await _fetchNotifications();
  }
}
