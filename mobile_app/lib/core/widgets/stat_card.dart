import 'dart:ui';
import 'package:flutter/material.dart';
import 'package:google_fonts/google_fonts.dart';
import '../theme/app_theme.dart';

/// Stat card variants matching web design (.stat-card-blue, .stat-card-green, etc.)
enum StatCardVariant { blue, green, purple, orange }

/// A stat card widget that exactly matches the web frontend design
class StatCard extends StatelessWidget {
  final String title;
  final String value;
  final String? subtitle;
  final String? badge;
  final bool isBadgePositive;
  final Widget? icon;
  final StatCardVariant variant;

  const StatCard({
    super.key,
    required this.title,
    required this.value,
    this.subtitle,
    this.badge,
    this.isBadgePositive = true,
    this.icon,
    this.variant = StatCardVariant.blue,
  });

  @override
  Widget build(BuildContext context) {
    final isDark = Theme.of(context).brightness == Brightness.dark;
    
    return Container(
      padding: const EdgeInsets.all(16),
      decoration: BoxDecoration(
        // Match web gradient backgrounds
        gradient: _getGradient(isDark),
        borderRadius: BorderRadius.circular(20),
        // Colored left border (3px) matching web design
        border: Border(
          left: BorderSide(
            color: _getBorderColor(),
            width: 3,
          ),
        ),
        boxShadow: [
          BoxShadow(
            color: isDark 
                ? Colors.black.withOpacity(0.3)
                : Colors.black.withOpacity(0.03),
            blurRadius: 20,
            offset: const Offset(0, 4),
          ),
        ],
      ),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        mainAxisSize: MainAxisSize.min,
        children: [
          // Top row: Icon and Badge
          Row(
            mainAxisAlignment: MainAxisAlignment.spaceBetween,
            children: [
              // Icon container
              Container(
                width: 40,
                height: 40,
                decoration: BoxDecoration(
                  color: _getIconBgColor(),
                  borderRadius: BorderRadius.circular(10),
                ),
                child: icon ?? Icon(
                  _getDefaultIcon(),
                  color: _getIconColor(),
                  size: 20,
                ),
              ),
              // Badge
              if (badge != null)
                Container(
                  padding: const EdgeInsets.symmetric(horizontal: 12, vertical: 4),
                  decoration: BoxDecoration(
                    color: isBadgePositive
                        ? const Color(0xFF10B981).withOpacity(0.2)
                        : const Color(0xFFEF4444).withOpacity(0.2),
                    borderRadius: BorderRadius.circular(20),
                  ),
                  child: Text(
                    badge!,
                    style: GoogleFonts.inter(
                      fontSize: 12,
                      fontWeight: FontWeight.w600,
                      color: isBadgePositive
                          ? const Color(0xFF10B981)
                          : const Color(0xFFEF4444),
                    ),
                  ),
                ),
            ],
          ),
          const SizedBox(height: 12),
          // Title
          Text(
            title,
            style: GoogleFonts.inter(
              fontSize: 14,
              fontWeight: FontWeight.w500,
              color: isDark ? const Color(0xFF94A3B8) : const Color(0xFF64748B),
            ),
          ),
          const SizedBox(height: 4),
          // Value
          Text(
            value,
            style: GoogleFonts.inter(
              fontSize: 20,
              fontWeight: FontWeight.bold,
              color: isDark ? Colors.white : const Color(0xFF1E293B),
            ),
            maxLines: 1,
            overflow: TextOverflow.ellipsis,
          ),
          // Subtitle
          if (subtitle != null) ...[
            const SizedBox(height: 4),
            Text(
              subtitle!,
              style: GoogleFonts.inter(
                fontSize: 14,
                fontWeight: FontWeight.w500,
                color: isDark ? const Color(0xFF64748B) : const Color(0xFF9CA3AF),
              ),
            ),
          ],
        ],
      ),
    );
  }

  LinearGradient _getGradient(bool isDark) {
    switch (variant) {
      case StatCardVariant.blue:
        return isDark
            ? const LinearGradient(
                begin: Alignment.topLeft,
                end: Alignment.bottomRight,
                colors: [Color(0xFF1E2850), Color(0xFF324282)],
              )
            : const LinearGradient(
                begin: Alignment.topLeft,
                end: Alignment.bottomRight,
                colors: [Color(0xFFF0F5FF), Color(0xFFE6EDFF)],
              );
      case StatCardVariant.green:
        return isDark
            ? const LinearGradient(
                begin: Alignment.topLeft,
                end: Alignment.bottomRight,
                colors: [Color(0xFF143C32), Color(0xFF1E5A46)],
              )
            : const LinearGradient(
                begin: Alignment.topLeft,
                end: Alignment.bottomRight,
                colors: [Color(0xFFF0FDF6), Color(0xFFDCFCE7)],
              );
      case StatCardVariant.purple:
        return isDark
            ? const LinearGradient(
                begin: Alignment.topLeft,
                end: Alignment.bottomRight,
                colors: [Color(0xFF321E50), Color(0xFF503278)],
              )
            : const LinearGradient(
                begin: Alignment.topLeft,
                end: Alignment.bottomRight,
                colors: [Color(0xFFFAF6FF), Color(0xFFF3E8FF)],
              );
      case StatCardVariant.orange:
        return isDark
            ? const LinearGradient(
                begin: Alignment.topLeft,
                end: Alignment.bottomRight,
                colors: [Color(0xFF3C2814), Color(0xFF643C1E)],
              )
            : const LinearGradient(
                begin: Alignment.topLeft,
                end: Alignment.bottomRight,
                colors: [Color(0xFFFFFCF5), Color(0xFFFEF6E6)],
              );
    }
  }

  Color _getBorderColor() {
    switch (variant) {
      case StatCardVariant.blue:
        return const Color(0xFF6366F1);
      case StatCardVariant.green:
        return const Color(0xFF22C55E);
      case StatCardVariant.purple:
        return const Color(0xFFA855F7);
      case StatCardVariant.orange:
        return const Color(0xFFF59E0B);
    }
  }

  Color _getIconBgColor() {
    switch (variant) {
      case StatCardVariant.blue:
        return const Color(0xFF3B82F6).withOpacity(0.3);
      case StatCardVariant.green:
        return const Color(0xFF10B981).withOpacity(0.3);
      case StatCardVariant.purple:
        return const Color(0xFFA855F7).withOpacity(0.3);
      case StatCardVariant.orange:
        return const Color(0xFFF59E0B).withOpacity(0.3);
    }
  }

  Color _getIconColor() {
    switch (variant) {
      case StatCardVariant.blue:
        return const Color(0xFF60A5FA);
      case StatCardVariant.green:
        return const Color(0xFF34D399);
      case StatCardVariant.purple:
        return const Color(0xFFC084FC);
      case StatCardVariant.orange:
        return const Color(0xFFFBBF24);
    }
  }

  IconData _getDefaultIcon() {
    switch (variant) {
      case StatCardVariant.blue:
        return Icons.attach_money;
      case StatCardVariant.green:
        return Icons.currency_bitcoin;
      case StatCardVariant.purple:
        return Icons.credit_card;
      case StatCardVariant.orange:
        return Icons.swap_horiz;
    }
  }
}
