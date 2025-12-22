import 'package:flutter/material.dart';

class AssetAllocation extends StatelessWidget {
  final List<Holding> holdings;
  final double totalValue;

  const AssetAllocation({super.key, required this.holdings, required this.totalValue});

  @override
  Widget build(BuildContext context) {
    final colors = [Colors.blue, Colors.green, Colors.purple, Colors.orange, Colors.red, Colors.teal];

    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        const Text('Répartition des actifs', style: TextStyle(fontSize: 16, fontWeight: FontWeight.bold)),
        const SizedBox(height: 16),
        
        // Visual bar
        ClipRRect(
          borderRadius: BorderRadius.circular(8),
          child: SizedBox(
            height: 16,
            child: Row(
              children: holdings.asMap().entries.map((entry) {
                final percent = (entry.value.value / totalValue * 100).clamp(0, 100);
                return Expanded(
                  flex: percent.toInt().clamp(1, 100),
                  child: Container(color: colors[entry.key % colors.length]),
                );
              }).toList(),
            ),
          ),
        ),
        const SizedBox(height: 16),
        
        // Legend
        Wrap(
          spacing: 16,
          runSpacing: 8,
          children: holdings.asMap().entries.map((entry) {
            final percent = (entry.value.value / totalValue * 100);
            return Row(
              mainAxisSize: MainAxisSize.min,
              children: [
                Container(
                  width: 12,
                  height: 12,
                  decoration: BoxDecoration(
                    color: colors[entry.key % colors.length],
                    borderRadius: BorderRadius.circular(3),
                  ),
                ),
                const SizedBox(width: 6),
                Text('${entry.value.symbol} ${percent.toStringAsFixed(1)}%', style: const TextStyle(fontSize: 12)),
              ],
            );
          }).toList(),
        ),
      ],
    );
  }
}

class Holding {
  final String symbol;
  final String name;
  final double quantity;
  final double value;
  final double change24h;

  Holding({
    required this.symbol,
    required this.name,
    required this.quantity,
    required this.value,
    this.change24h = 0,
  });
}
