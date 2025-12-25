import 'package:flutter/material.dart';
import 'package:google_fonts/google_fonts.dart';

/// Quick action button matching web .quick-action-btn styling
class QuickActionButton extends StatelessWidget {
  final String emoji;
  final String label;
  final VoidCallback onTap;

  const QuickActionButton({
    super.key,
    required this.emoji,
    required this.label,
    required this.onTap,
  });

  @override
  Widget build(BuildContext context) {
    final isDark = Theme.of(context).brightness == Brightness.dark;
    
    return GestureDetector(
      onTap: onTap,
      child: AnimatedContainer(
        duration: const Duration(milliseconds: 300),
        padding: const EdgeInsets.all(16),
        decoration: BoxDecoration(
          // Light mode: softer, warmer - Dark mode: glass effect
          color: isDark 
              ? Colors.white.withOpacity(0.05)
              : const Color(0xFFF1F5F9).withOpacity(0.5),
          borderRadius: BorderRadius.circular(16),
          border: Border.all(
            color: isDark 
                ? Colors.white.withOpacity(0.1)
                : const Color(0xFFE2E8F0).withOpacity(0.8),
            width: 1,
          ),
        ),
        child: Column(
          mainAxisAlignment: MainAxisAlignment.center,
          children: [
            Text(
              emoji,
              style: const TextStyle(fontSize: 28),
            ),
            const SizedBox(height: 8),
            Text(
              label,
              textAlign: TextAlign.center,
              style: GoogleFonts.inter(
                fontSize: 12,
                fontWeight: FontWeight.w500,
                color: isDark ? const Color(0xFFE2E8F0) : const Color(0xFF475569),
              ),
            ),
          ],
        ),
      ),
    );
  }
}

/// Quick actions grid matching web design
class QuickActionsGrid extends StatelessWidget {
  final List<QuickActionItem> actions;
  final int crossAxisCount;

  const QuickActionsGrid({
    super.key,
    required this.actions,
    this.crossAxisCount = 3,
  });

  @override
  Widget build(BuildContext context) {
    return GridView.builder(
      shrinkWrap: true,
      physics: const NeverScrollableScrollPhysics(),
      gridDelegate: SliverGridDelegateWithFixedCrossAxisCount(
        crossAxisCount: crossAxisCount,
        crossAxisSpacing: 12,
        mainAxisSpacing: 12,
        childAspectRatio: 1.0,
      ),
      itemCount: actions.length,
      itemBuilder: (context, index) {
        final action = actions[index];
        return QuickActionButton(
          emoji: action.emoji,
          label: action.label,
          onTap: action.onTap,
        );
      },
    );
  }
}

/// Quick action item data model
class QuickActionItem {
  final String emoji;
  final String label;
  final VoidCallback onTap;

  const QuickActionItem({
    required this.emoji,
    required this.label,
    required this.onTap,
  });
}
