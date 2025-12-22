import 'package:flutter/material.dart';

class RecentTransactionsSection extends StatelessWidget {
  const RecentTransactionsSection({super.key});

  @override
  Widget build(BuildContext context) {
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
            const Text(
              'Transactions récentes',
              style: TextStyle(fontSize: 18, fontWeight: FontWeight.bold),
            ),
            TextButton(
              onPressed: () {},
              child: const Text('Voir tout'),
            ),
          ],
        ),
        const SizedBox(height: 8),
        ...transactions.map((tx) => _buildTransactionTile(tx)),
      ],
    );
  }

  Widget _buildTransactionTile(_TransactionItem tx) {
    IconData icon;
    Color color;
    String prefix;

    switch (tx.type) {
      case 'received':
        icon = Icons.arrow_downward;
        color = Colors.green;
        prefix = '+';
        break;
      case 'sent':
        icon = Icons.arrow_upward;
        color = Colors.red;
        prefix = '-';
        break;
      case 'exchange':
        icon = Icons.swap_horiz;
        color = Colors.purple;
        prefix = '';
        break;
      default:
        icon = Icons.credit_card;
        color = Colors.orange;
        prefix = '-';
    }

    return ListTile(
      contentPadding: EdgeInsets.zero,
      leading: Container(
        width: 48,
        height: 48,
        decoration: BoxDecoration(
          color: color.withOpacity(0.1),
          borderRadius: BorderRadius.circular(12),
        ),
        child: Icon(icon, color: color),
      ),
      title: Text(tx.title),
      subtitle: Text(tx.date, style: const TextStyle(fontSize: 12)),
      trailing: Text(
        '$prefix${tx.amount.toStringAsFixed(2)} ${tx.currency}',
        style: TextStyle(
          color: tx.type == 'received' ? Colors.green : (tx.type == 'exchange' ? Colors.black : Colors.red),
          fontWeight: FontWeight.bold,
        ),
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
