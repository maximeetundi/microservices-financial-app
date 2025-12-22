import 'package:flutter/material.dart';

class MarketOverviewSection extends StatelessWidget {
  const MarketOverviewSection({super.key});

  @override
  Widget build(BuildContext context) {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        const Text(
          'Aperçu du marché',
          style: TextStyle(fontSize: 18, fontWeight: FontWeight.bold),
        ),
        const SizedBox(height: 16),
        Container(
          padding: const EdgeInsets.all(16),
          decoration: BoxDecoration(
            color: Colors.grey[100],
            borderRadius: BorderRadius.circular(16),
          ),
          child: Row(
            mainAxisAlignment: MainAxisAlignment.spaceAround,
            children: [
              _MarketStat(label: 'Cap. Marché', value: '\$1.2T', change: '+2.3%', isPositive: true),
              Container(width: 1, height: 40, color: Colors.grey[300]),
              _MarketStat(label: 'Volume 24h', value: '\$84B', change: '+5.1%', isPositive: true),
              Container(width: 1, height: 40, color: Colors.grey[300]),
              _MarketStat(label: 'BTC Dom.', value: '48.2%', change: '-0.5%', isPositive: false),
            ],
          ),
        ),
      ],
    );
  }
}

class _MarketStat extends StatelessWidget {
  final String label;
  final String value;
  final String change;
  final bool isPositive;

  const _MarketStat({
    required this.label,
    required this.value,
    required this.change,
    required this.isPositive,
  });

  @override
  Widget build(BuildContext context) {
    return Column(
      children: [
        Text(label, style: TextStyle(fontSize: 11, color: Colors.grey[600])),
        const SizedBox(height: 4),
        Text(value, style: const TextStyle(fontSize: 16, fontWeight: FontWeight.bold)),
        const SizedBox(height: 2),
        Text(
          change,
          style: TextStyle(
            fontSize: 12,
            color: isPositive ? Colors.green : Colors.red,
            fontWeight: FontWeight.w500,
          ),
        ),
      ],
    );
  }
}
