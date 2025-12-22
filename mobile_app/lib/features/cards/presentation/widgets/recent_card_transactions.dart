import 'package:flutter/material.dart';

class RecentCardTransactions extends StatelessWidget {
  final String cardId;
  final List<CardTransaction> transactions;

  const RecentCardTransactions({
    super.key,
    required this.cardId,
    required this.transactions,
  });

  @override
  Widget build(BuildContext context) {
    if (transactions.isEmpty) {
      return const Center(
        child: Padding(
          padding: EdgeInsets.all(32),
          child: Text('Aucune transaction récente', style: TextStyle(color: Colors.grey)),
        ),
      );
    }

    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Row(
          mainAxisAlignment: MainAxisAlignment.spaceBetween,
          children: [
            const Text('Transactions récentes', style: TextStyle(fontSize: 16, fontWeight: FontWeight.bold)),
            TextButton(onPressed: () {}, child: const Text('Voir tout')),
          ],
        ),
        ...transactions.take(5).map((tx) => _buildTransactionTile(tx)),
      ],
    );
  }

  Widget _buildTransactionTile(CardTransaction tx) {
    return ListTile(
      contentPadding: EdgeInsets.zero,
      leading: Container(
        width: 44,
        height: 44,
        decoration: BoxDecoration(
          color: Colors.grey[100],
          borderRadius: BorderRadius.circular(10),
        ),
        child: Center(child: Text(_getMerchantIcon(tx.merchant), style: const TextStyle(fontSize: 20))),
      ),
      title: Text(tx.merchant, style: const TextStyle(fontWeight: FontWeight.w500)),
      subtitle: Text(tx.date, style: const TextStyle(fontSize: 12)),
      trailing: Text(
        '-\$${tx.amount.toStringAsFixed(2)}',
        style: const TextStyle(fontWeight: FontWeight.bold, color: Colors.red),
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
