import 'package:flutter/material.dart';
import '../../../../core/theme/app_theme.dart';
import 'package:google_fonts/google_fonts.dart';

class RecentTransactionsSection extends StatelessWidget {
  const RecentTransactionsSection({super.key});

  @override
  Widget build(BuildContext context) {
    final isDark = Theme.of(context).brightness == Brightness.dark;

    final transactions = [
      _TransactionItem(type: 'received', title: 'De John Doe', amount: 500, currency: 'USD', date: 'Aujourd\'hui'),
      _TransactionItem(type: 'sent', title: 'À Marie Martin', amount: 150, currency: 'EUR', date: 'Hier'),
      _TransactionItem(type: 'exchange', title: 'BTC → USD', amount: 1200, currency: 'USD', date: 'Il y a 2j'),
      _TransactionItem(type: 'card', title: 'Amazon', amount: 45.99, currency: 'USD', date: 'Il y a 3j'),
    ];

    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Row(
          mainAxisAlignment: MainAxisAlignment.spaceBetween,
          children: [
            Text(
              'Transactions récentes',
              style: GoogleFonts.inter(
                fontSize: 18, 
                fontWeight: FontWeight.bold,
                color: isDark ? Colors.white : AppTheme.textPrimaryColor,
              ),
            ),
            TextButton(
              onPressed: () {},
              child: Text(
                'Voir tout',
                style: GoogleFonts.inter(color: AppTheme.primaryColor, fontWeight: FontWeight.w600),
              ),
            ),
          ],
        ),
        const SizedBox(height: 8),
        ...transactions.map((tx) => _buildTransactionTile(context, tx, isDark)),
      ],
    );
  }

  Widget _buildTransactionTile(BuildContext context, _TransactionItem tx, bool isDark) {
    IconData icon;
    Color color;
    String prefix;

    switch (tx.type) {
      case 'received':
        icon = Icons.arrow_downward;
        color = const Color(0xFF10B981);
        prefix = '+';
        break;
      case 'sent':
        icon = Icons.arrow_upward;
        color = const Color(0xFFEF4444);
        prefix = '-';
        break;
      case 'exchange':
        icon = Icons.swap_horiz;
        color = const Color(0xFF8B5CF6);
        prefix = '';
        break;
      default:
        icon = Icons.credit_card;
        color = const Color(0xFFF59E0B);
        prefix = '-';
    }

    return Padding(
      padding: const EdgeInsets.symmetric(vertical: 8.0),
      child: Row(
        children: [
          Container(
            width: 48,
            height: 48,
            decoration: BoxDecoration(
              color: color.withOpacity(0.1),
              shape: BoxShape.circle,
            ),
            child: Icon(icon, color: color, size: 24),
          ),
          const SizedBox(width: 16),
          Expanded(
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Text(
                  tx.title,
                  style: GoogleFonts.inter(
                    fontWeight: FontWeight.w600,
                    fontSize: 16,
                    color: isDark ? Colors.white : AppTheme.textPrimaryColor,
                  ),
                ),
                const SizedBox(height: 4),
                Text(
                  tx.date, 
                  style: GoogleFonts.inter(
                    fontSize: 12,
                    color: isDark ? Colors.white54 : AppTheme.textSecondaryColor,
                  ),
                ),
              ],
            ),
          ),
          Text(
            '$prefix${tx.amount.toStringAsFixed(2)} ${tx.currency}',
            style: GoogleFonts.inter(
              color: tx.type == 'received' ? const Color(0xFF10B981) : (tx.type == 'exchange' ? (isDark ? Colors.white : AppTheme.textPrimaryColor) : const Color(0xFFEF4444)),
              fontWeight: FontWeight.bold,
              fontSize: 16,
            ),
          ),
        ],
      ),
    );
  }
}

class _TransactionItem {
  final String type;
  final String title;
  final double amount;
  final String currency;
  final String date;

  _TransactionItem({
    required this.type,
    required this.title,
    required this.amount,
    required this.currency,
    required this.date,
  });
}
