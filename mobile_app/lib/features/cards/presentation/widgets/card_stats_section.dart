import 'package:flutter/material.dart';
import '../../../../core/theme/app_theme.dart';
import 'package:google_fonts/google_fonts.dart';

class CardStatsSection extends StatelessWidget {
  final double monthlySpend;
  final double limit;
  final int transactionCount;

  const CardStatsSection({
    super.key,
    required this.monthlySpend,
    required this.limit,
    required this.transactionCount,
  });

  @override
  Widget build(BuildContext context) {
    final percentUsed = (monthlySpend / limit * 100).clamp(0, 100);
    final isDark = Theme.of(context).brightness == Brightness.dark;

    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Text(
          'Statistiques du mois', 
          style: GoogleFonts.inter(
            fontSize: 16, 
            fontWeight: FontWeight.bold,
            color: isDark ? Colors.white : AppTheme.textPrimaryColor
          )
        ),
        const SizedBox(height: 16),
        
        // Progress bar
        Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Row(
              mainAxisAlignment: MainAxisAlignment.spaceBetween,
              children: [
                Text(
                  '\$${monthlySpend.toStringAsFixed(0)} dépensés', 
                  style: GoogleFonts.inter(
                    fontWeight: FontWeight.w500,
                    color: isDark ? Colors.white70 : AppTheme.textPrimaryColor
                  )
                ),
                Text(
                  'sur \$${limit.toStringAsFixed(0)}', 
                  style: GoogleFonts.inter(
                    color: isDark ? Colors.white38 : AppTheme.textSecondaryColor
                  )
                ),
              ],
            ),
            const SizedBox(height: 8),
            ClipRRect(
              borderRadius: BorderRadius.circular(4),
              child: LinearProgressIndicator(
                value: percentUsed / 100,
                minHeight: 8,
                backgroundColor: isDark ? Colors.white.withOpacity(0.1) : Colors.grey[200],
                valueColor: AlwaysStoppedAnimation<Color>(
                  percentUsed > 80 ? AppTheme.errorColor : (percentUsed > 50 ? Colors.orange : AppTheme.successColor),
                ),
              ),
            ),
            const SizedBox(height: 4),
            Text(
              '${percentUsed.toStringAsFixed(0)}% utilisé', 
              style: GoogleFonts.inter(
                fontSize: 12, 
                color: isDark ? Colors.white38 : AppTheme.textSecondaryColor
              )
            ),
          ],
        ),
        const SizedBox(height: 16),
        
        // Stats row
        Row(
          children: [
            Expanded(
              child: _StatCard(
                icon: Icons.receipt_long,
                label: 'Transactions',
                value: transactionCount.toString(),
                color: Colors.blue,
              ),
            ),
            const SizedBox(width: 12),
            Expanded(
              child: _StatCard(
                icon: Icons.trending_up,
                label: 'Moyenne',
                value: '\$${(monthlySpend / (transactionCount > 0 ? transactionCount : 1)).toStringAsFixed(0)}',
                color: AppTheme.successColor,
              ),
            ),
          ],
        ),
      ],
    );
  }
}

class _StatCard extends StatelessWidget {
  final IconData icon;
  final String label;
  final String value;
  final Color color;

  const _StatCard({
    required this.icon,
    required this.label,
    required this.value,
    required this.color,
  });

  @override
  Widget build(BuildContext context) {
    final isDark = Theme.of(context).brightness == Brightness.dark;

    return Container(
      padding: const EdgeInsets.all(16),
      decoration: BoxDecoration(
        color: color.withOpacity(0.1),
        borderRadius: BorderRadius.circular(16),
        border: Border.all(color: color.withOpacity(0.2)),
      ),
      child: Row(
        children: [
          Icon(icon, color: color, size: 24),
          const SizedBox(width: 12),
          Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              Text(
                label, 
                style: GoogleFonts.inter(
                  fontSize: 12, 
                  color: isDark ? Colors.white54 : AppTheme.textSecondaryColor
                )
              ),
              Text(
                value, 
                style: GoogleFonts.inter(
                  fontSize: 18, 
                  fontWeight: FontWeight.bold,
                  color: isDark ? Colors.white : AppTheme.textPrimaryColor
                )
              ),
            ],
          ),
        ],
      ),
    );
  }
}
