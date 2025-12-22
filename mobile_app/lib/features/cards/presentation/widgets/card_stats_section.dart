import 'package:flutter/material.dart';

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

    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        const Text('Statistiques du mois', style: TextStyle(fontSize: 16, fontWeight: FontWeight.bold)),
        const SizedBox(height: 16),
        
        // Progress bar
        Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Row(
              mainAxisAlignment: MainAxisAlignment.spaceBetween,
              children: [
                Text('\$${monthlySpend.toStringAsFixed(0)} dépensés', style: const TextStyle(fontWeight: FontWeight.w500)),
                Text('sur \$${limit.toStringAsFixed(0)}', style: TextStyle(color: Colors.grey[600])),
              ],
            ),
            const SizedBox(height: 8),
            ClipRRect(
              borderRadius: BorderRadius.circular(4),
              child: LinearProgressIndicator(
                value: percentUsed / 100,
                minHeight: 8,
                backgroundColor: Colors.grey[200],
                valueColor: AlwaysStoppedAnimation<Color>(
                  percentUsed > 80 ? Colors.red : (percentUsed > 50 ? Colors.orange : Colors.green),
                ),
              ),
            ),
            const SizedBox(height: 4),
            Text('${percentUsed.toStringAsFixed(0)}% utilisé', style: TextStyle(fontSize: 12, color: Colors.grey[600])),
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
                color: Colors.green,
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
    return Container(
      padding: const EdgeInsets.all(16),
      decoration: BoxDecoration(
        color: color.withOpacity(0.1),
        borderRadius: BorderRadius.circular(12),
      ),
      child: Row(
        children: [
          Icon(icon, color: color, size: 24),
          const SizedBox(width: 12),
          Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              Text(label, style: TextStyle(fontSize: 12, color: Colors.grey[600])),
              Text(value, style: const TextStyle(fontSize: 18, fontWeight: FontWeight.bold)),
            ],
          ),
        ],
      ),
    );
  }
}
