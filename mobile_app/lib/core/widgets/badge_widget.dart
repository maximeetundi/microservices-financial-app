import 'package:flutter/material.dart';
import 'package:google_fonts/google_fonts.dart';

/// Badge variants matching web design (.badge-success, .badge-warning, etc.)
enum BadgeVariant { success, warning, danger, info }

/// A badge widget that exactly matches the web frontend design
class BadgeWidget extends StatelessWidget {
  final String text;
  final BadgeVariant variant;

  const BadgeWidget({
    super.key,
    required this.text,
    this.variant = BadgeVariant.success,
  });

  @override
  Widget build(BuildContext context) {
    return Container(
      padding: const EdgeInsets.symmetric(horizontal: 12, vertical: 4),
      decoration: BoxDecoration(
        color: _getBackgroundColor(),
        borderRadius: BorderRadius.circular(20),
      ),
      child: Text(
        text,
        style: GoogleFonts.inter(
          fontSize: 12,
          fontWeight: FontWeight.w600,
          color: _getTextColor(),
        ),
      ),
    );
  }

  Color _getBackgroundColor() {
    switch (variant) {
      case BadgeVariant.success:
        return const Color(0xFF10B981).withOpacity(0.2);
      case BadgeVariant.warning:
        return const Color(0xFFF59E0B).withOpacity(0.2);
      case BadgeVariant.danger:
        return const Color(0xFFEF4444).withOpacity(0.2);
      case BadgeVariant.info:
        return const Color(0xFF3B82F6).withOpacity(0.2);
    }
  }

  Color _getTextColor() {
    switch (variant) {
      case BadgeVariant.success:
        return const Color(0xFF10B981);
      case BadgeVariant.warning:
        return const Color(0xFFF59E0B);
      case BadgeVariant.danger:
        return const Color(0xFFEF4444);
      case BadgeVariant.info:
        return const Color(0xFF3B82F6);
    }
  }
}
