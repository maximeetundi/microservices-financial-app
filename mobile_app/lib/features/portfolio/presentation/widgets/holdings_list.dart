import 'package:flutter/material.dart';
import 'asset_allocation.dart';

class HoldingsList extends StatelessWidget {
  final List<Holding> holdings;

  const HoldingsList({super.key, required this.holdings});

  @override
  Widget build(BuildContext context) {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        const Text('Mes actifs', style: TextStyle(fontSize: 16, fontWeight: FontWeight.bold)),
        const SizedBox(height: 16),
        ...holdings.map((h) => _buildHoldingTile(h)),
      ],
    );
  }

  Widget _buildHoldingTile(Holding holding) {
    final isPositive = holding.change24h >= 0;

    return ListTile(
      contentPadding: EdgeInsets.zero,
      leading: Container(
        width: 48,
        height: 48,
        decoration: BoxDecoration(
          color: Colors.grey[100],
          borderRadius: BorderRadius.circular(12),
        ),
        child: Center(
          child: Text(_getSymbolIcon(holding.symbol), style: const TextStyle(fontSize: 24)),
        ),
      ),
      title: Text(holding.symbol, style: const TextStyle(fontWeight: FontWeight.bold)),
      subtitle: Text('${holding.quantity.toStringAsFixed(holding.quantity < 1 ? 6 : 2)} ${holding.symbol}'),
      trailing: Column(
        mainAxisAlignment: MainAxisAlignment.center,
        crossAxisAlignment: CrossAxisAlignment.end,
        children: [
          Text('\$${holding.value.toStringAsFixed(2)}', style: const TextStyle(fontWeight: FontWeight.bold)),
          Row(
            mainAxisSize: MainAxisSize.min,
            children: [
              Icon(
                isPositive ? Icons.arrow_drop_up : Icons.arrow_drop_down,
                color: isPositive ? Colors.green : Colors.red,
                size: 18,
              ),
              Text(
                '${isPositive ? '+' : ''}${holding.change24h.toStringAsFixed(2)}%',
                style: TextStyle(
                  fontSize: 12,
                  color: isPositive ? Colors.green : Colors.red,
                ),
              ),
            ],
          ),
        ],
      ),
    );
  }

  String _getSymbolIcon(String symbol) {
    final icons = {
      'BTC': '₿', 'ETH': 'Ξ', 'SOL': '◎', 'XRP': '✕',
      'USD': '\$', 'EUR': '€', 'GBP': '£',
    };
    return icons[symbol] ?? symbol[0];
  }
}
