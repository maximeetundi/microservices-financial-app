import 'package:flutter/material.dart';

class PerformanceMetrics extends StatelessWidget {
  final PortfolioPerformance performance;

  const PerformanceMetrics({super.key, required this.performance});

  @override
  Widget build(BuildContext context) {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        const Text('Performance', style: TextStyle(fontSize: 16, fontWeight: FontWeight.bold)),
        const SizedBox(height: 16),
        Row(
          children: [
            Expanded(child: _MetricCard(label: 'Rendement 24h', value: '${performance.dailyReturn >= 0 ? '+' : ''}${performance.dailyReturn.toStringAsFixed(2)}%', isPositive: performance.dailyReturn >= 0)),
            const SizedBox(width: 12),
            Expanded(child: _MetricCard(label: 'Rendement 7j', value: '${performance.weeklyReturn >= 0 ? '+' : ''}${performance.weeklyReturn.toStringAsFixed(2)}%', isPositive: performance.weeklyReturn >= 0)),
          ],
        ),
        const SizedBox(height: 12),
        Row(
          children: [
            Expanded(child: _MetricCard(label: 'Rendement 30j', value: '${performance.monthlyReturn >= 0 ? '+' : ''}${performance.monthlyReturn.toStringAsFixed(2)}%', isPositive: performance.monthlyReturn >= 0)),
            const SizedBox(width: 12),
            Expanded(child: _MetricCard(label: 'Total', value: '${performance.totalReturn >= 0 ? '+' : ''}${performance.totalReturn.toStringAsFixed(2)}%', isPositive: performance.totalReturn >= 0)),
          ],
        ),
      ],
    );
  }
}

class _MetricCard extends StatelessWidget {
  final String label;
  final String value;
  final bool isPositive;

  const _MetricCard({required this.label, required this.value, required this.isPositive});

  @override
  Widget build(BuildContext context) {
    return Container(
      padding: const EdgeInsets.all(16),
      decoration: BoxDecoration(
        color: Colors.grey[100],
        borderRadius: BorderRadius.circular(12),
      ),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Text(label, style: TextStyle(fontSize: 12, color: Colors.grey[600])),
          const SizedBox(height: 4),
          Text(
            value,
            style: TextStyle(
              fontSize: 18,
              fontWeight: FontWeight.bold,
              color: isPositive ? Colors.green : Colors.red,
            ),
          ),
        ],
      ),
    );
  }
}

class PortfolioPerformance {
  final double dailyReturn;
  final double weeklyReturn;
  final double monthlyReturn;
  final double totalReturn;

  PortfolioPerformance({
    required this.dailyReturn,
    required this.weeklyReturn,
    required this.monthlyReturn,
    required this.totalReturn,
  });
}
