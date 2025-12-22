import 'package:flutter/material.dart';

class ExchangeHistoryList extends StatelessWidget {
  final List<ExchangeHistoryItem> history;

  const ExchangeHistoryList({super.key, required this.history});

  @override
  Widget build(BuildContext context) {
    if (history.isEmpty) {
      return const Center(
        child: Padding(
          padding: EdgeInsets.all(32),
          child: Text('Aucun échange récent', style: TextStyle(color: Colors.grey)),
        ),
      );
    }

    return ListView.builder(
      shrinkWrap: true,
      physics: const NeverScrollableScrollPhysics(),
      itemCount: history.length,
      itemBuilder: (context, index) {
        final item = history[index];
        return ListTile(
          contentPadding: EdgeInsets.zero,
          leading: Container(
            width: 48,
            height: 48,
            decoration: BoxDecoration(
              color: Colors.purple.withOpacity(0.1),
              borderRadius: BorderRadius.circular(12),
            ),
            child: const Icon(Icons.swap_horiz, color: Colors.purple),
          ),
          title: Text('${item.fromAmount} ${item.fromCurrency} → ${item.toAmount} ${item.toCurrency}'),
          subtitle: Text(item.date),
          trailing: Container(
            padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 4),
            decoration: BoxDecoration(
              color: item.status == 'completed' ? Colors.green[50] : Colors.orange[50],
              borderRadius: BorderRadius.circular(8),
            ),
            child: Text(
              item.status == 'completed' ? 'Terminé' : 'En cours',
              style: TextStyle(
                fontSize: 12,
                color: item.status == 'completed' ? Colors.green : Colors.orange,
              ),
            ),
          ),
        );
      },
    );
  }
}

class ExchangeHistoryItem {
  final String fromCurrency;
  final String toCurrency;
  final double fromAmount;
  final double toAmount;
  final String date;
  final String status;

  ExchangeHistoryItem({
    required this.fromCurrency,
    required this.toCurrency,
    required this.fromAmount,
    required this.toAmount,
    required this.date,
    this.status = 'completed',
  });
}
