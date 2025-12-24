import 'package:flutter/material.dart';
import '../../../../core/theme/app_theme.dart';
import 'package:google_fonts/google_fonts.dart';

class RecentCardTransactions extends StatelessWidget {
  final String cardId;
  final List<CardTransaction> transactions;

  const RecentCardTransactions({
    super.key,
    required this.cardId,
    required this.transactions,
  });

  @override
  @override
  Widget build(BuildContext context) {
    if (transactions.isEmpty) {
      return Center(
        child: Padding(
          padding: const EdgeInsets.all(32),
          child: Text(
            'Aucune transaction récente', 
            style: GoogleFonts.inter(color: Theme.of(context).brightness == Brightness.dark ? Colors.white54 : AppTheme.textSecondaryColor)
          ),
        ),
      );
    }

    final isDark = Theme.of(context).brightness == Brightness.dark;

    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Row(
          mainAxisAlignment: MainAxisAlignment.spaceBetween,
          children: [
            Text(
              'Transactions récentes', 
              style: GoogleFonts.inter(fontSize: 16, fontWeight: FontWeight.bold, color: isDark ? Colors.white : AppTheme.textPrimaryColor)
            ),
            TextButton(
              onPressed: () {}, 
              child: Text(
                'Voir tout',
                style: GoogleFonts.inter(color: AppTheme.primaryColor),
              )
            ),
          ],
        ),
        ...transactions.take(5).map((tx) => _buildTransactionTile(context, tx, isDark)),
      ],
    );
  }

  Widget _buildTransactionTile(BuildContext context, CardTransaction tx, bool isDark) {
    return Padding(
      padding: const EdgeInsets.symmetric(vertical: 8.0),
      child: Row(
        children: [
          Container(
            width: 44,
            height: 44,
            decoration: BoxDecoration(
              color: isDark ? Colors.white.withOpacity(0.1) : Colors.grey[100],
              borderRadius: BorderRadius.circular(10),
            ),
            child: Center(child: Text(_getMerchantIcon(tx.merchant), style: const TextStyle(fontSize: 20))),
          ),
          const SizedBox(width: 16),
          Expanded(
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Text(
                  tx.merchant, 
                  style: GoogleFonts.inter(
                    fontWeight: FontWeight.w500,
                    color: isDark ? Colors.white : AppTheme.textPrimaryColor,
                  )
                ),
                Text(
                  tx.date, 
                  style: GoogleFonts.inter(
                    fontSize: 12,
                    color: isDark ? Colors.white54 : AppTheme.textSecondaryColor,
                  )
                ),
              ],
            ),
          ),
          Text(
            '-\$${tx.amount.toStringAsFixed(2)}',
            style: GoogleFonts.inter(fontWeight: FontWeight.bold, color: const Color(0xFFEF4444)),
          ),
        ],
      ),
    );
  }

  String _getMerchantIcon(String merchant) {
    final m = merchant.toLowerCase();
    if (m.contains('amazon')) return '📦';
    if (m.contains('uber')) return '🚗';
    if (m.contains('netflix')) return '🎬';
    if (m.contains('spotify')) return '🎵';
    if (m.contains('carrefour') || m.contains('grocery')) return '🛒';
    if (m.contains('restaurant') || m.contains('food')) return '🍔';
    return '💳';
  }
}

class CardTransaction {
  final String id;
  final String cardId;
  final String merchant;
  final double amount;
  final String date;
  final String status;

  CardTransaction({
    required this.id,
    required this.cardId,
    required this.merchant,
    required this.amount,
    required this.date,
    this.status = 'completed',
  });
}
