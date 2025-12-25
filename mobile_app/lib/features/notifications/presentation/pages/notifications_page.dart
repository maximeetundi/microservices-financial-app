import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';
import 'package:google_fonts/google_fonts.dart';

import '../../../../core/theme/app_theme.dart';
import '../../../../core/widgets/glass_container.dart';

/// Notifications Page matching web design exactly
class NotificationsPage extends StatefulWidget {
  const NotificationsPage({super.key});

  @override
  State<NotificationsPage> createState() => _NotificationsPageState();
}

class _NotificationsPageState extends State<NotificationsPage> {
  String _activeFilter = 'all';
  bool _loading = false;
  
  final List<Map<String, String>> _filters = [
    {'id': 'all', 'label': 'Toutes', 'icon': '📋'},
    {'id': 'transfer', 'label': 'Transferts', 'icon': '💸'},
    {'id': 'security', 'label': 'Sécurité', 'icon': '🔐'},
    {'id': 'card', 'label': 'Cartes', 'icon': '💳'},
  ];
  
  // Sample notifications
  final List<Map<String, dynamic>> _notifications = [
    {
      'id': '1',
      'type': 'transfer',
      'title': 'Transfert reçu',
      'message': 'Vous avez reçu 50 000 XOF de John Doe',
      'created_at': DateTime.now().subtract(const Duration(minutes: 5)),
      'is_read': false,
    },
    {
      'id': '2',
      'type': 'security',
      'title': 'Nouvelle connexion',
      'message': 'Connexion détectée depuis Abidjan, Côte d\'Ivoire',
      'created_at': DateTime.now().subtract(const Duration(hours: 2)),
      'is_read': false,
    },
    {
      'id': '3',
      'type': 'card',
      'title': 'Carte activée',
      'message': 'Votre carte virtuelle a été activée avec succès',
      'created_at': DateTime.now().subtract(const Duration(days: 1)),
      'is_read': true,
    },
    {
      'id': '4',
      'type': 'transfer',
      'title': 'Transfert envoyé',
      'message': 'Votre transfert de 25 000 XOF a été effectué',
      'created_at': DateTime.now().subtract(const Duration(days: 2)),
      'is_read': true,
    },
  ];
  
  List<Map<String, dynamic>> get _filteredNotifications {
    if (_activeFilter == 'all') return _notifications;
    return _notifications.where((n) => n['type'] == _activeFilter).toList();
  }
  
  int get _unreadCount => _notifications.where((n) => !n['is_read']).length;

  @override
  Widget build(BuildContext context) {
    final isDark = Theme.of(context).brightness == Brightness.dark;
    
    return Scaffold(
      backgroundColor: isDark ? const Color(0xFF0F172A) : const Color(0xFFF8FAFC),
      body: Column(
        children: [
          // Gradient Header
          _buildHeader(isDark),
          
          // Content
          Expanded(
            child: _notifications.isEmpty
                ? _buildEmptyState(isDark)
                : ListView(
                    padding: const EdgeInsets.all(16),
                    children: [
                      // Filter Pills
                      _buildFilters(isDark),
                      const SizedBox(height: 16),
                      
                      // Notifications List
                      ..._filteredNotifications.map((n) => _buildNotificationItem(n, isDark)),
                    ],
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
    final isUnread = !notification['is_read'];
    
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
                  _getTypeIcon(notification['type']),
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
                    notification['title'],
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
                    notification['message'],
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
                    onTap: () => _markAsRead(notification['id']),
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
                  onTap: () => _deleteNotification(notification['id']),
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
      default: return '🔔';
    }
  }

  String _formatTime(DateTime date) {
    final diff = DateTime.now().difference(date);
    if (diff.inMinutes < 1) return 'À l\'instant';
    if (diff.inMinutes < 60) return '${diff.inMinutes} min';
    if (diff.inHours < 24) return '${diff.inHours}h';
    return '${diff.inDays}j';
  }

  void _handleTap(Map<String, dynamic> notification) {
    if (!notification['is_read']) {
      _markAsRead(notification['id']);
    }
  }

  void _markAsRead(String id) {
    setState(() {
      final notif = _notifications.firstWhere((n) => n['id'] == id);
      notif['is_read'] = true;
    });
  }

  void _markAllAsRead() {
    setState(() {
      for (var n in _notifications) {
        n['is_read'] = true;
      }
    });
  }

  void _deleteNotification(String id) {
    setState(() {
      _notifications.removeWhere((n) => n['id'] == id);
    });
  }

  void _refresh() {
    setState(() => _loading = true);
    Future.delayed(const Duration(seconds: 1), () {
      if (mounted) setState(() => _loading = false);
    });
  }
}
