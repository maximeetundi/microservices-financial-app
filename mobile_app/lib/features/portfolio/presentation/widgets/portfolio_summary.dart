import 'package:flutter/material.dart';

class PortfolioSummary extends StatelessWidget {
  final PortfolioData portfolio;

  const PortfolioSummary({super.key, required this.portfolio});

  @override
  Widget build(BuildContext context) {
    final isPositive = portfolio.totalChange >= 0;

    return Container(
      width: double.infinity,
      padding: const EdgeInsets.all(24),
      decoration: BoxDecoration(
        gradient: LinearGradient(
          colors: [Colors.indigo[700]!, Colors.indigo[500]!],
          begin: Alignment.topLeft,
          end: Alignment.bottomRight,
        ),
        borderRadius: BorderRadius.circular(20),
      ),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          const Text('Valeur totale', style: TextStyle(color: Colors.white70, fontSize: 14)),
          const SizedBox(height: 8),
          Text(
            '\$${portfolio.totalValue.toStringAsFixed(2)}',
            style: const TextStyle(color: Colors.white, fontSize: 36, fontWeight: FontWeight.bold),
          ),
          const SizedBox(height: 16),
          Row(
            children: [
              Container(
                padding: const EdgeInsets.symmetric(horizontal: 12, vertical: 6),
                decoration: BoxDecoration(
                  color: isPositive ? Colors.green.withOpacity(0.2) : Colors.red.withOpacity(0.2),
                  borderRadius: BorderRadius.circular(20),
                ),
                child: Row(
                  children: [
                    Icon(
                      isPositive ? Icons.trending_up : Icons.trending_down,
                      color: isPositive ? Colors.greenAccent : Colors.redAccent,
                      size: 18,
                    ),
                    const SizedBox(width: 4),
                    Text(
                      '${isPositive ? '+' : ''}\$${portfolio.totalChange.abs().toStringAsFixed(2)} (${portfolio.totalChangePercent.toStringAsFixed(2)}%)',
                      style: TextStyle(
                        color: isPositive ? Colors.greenAccent : Colors.redAccent,
                        fontWeight: FontWeight.bold,
                      ),
                    ),
                  ],
                ),
              ),
              const SizedBox(width: 8),
              const Text('24h', style: TextStyle(color: Colors.white70, fontSize: 12)),
            ],
          ),
        ],
      ),
    );
  }
}

class PortfolioData {
  final double totalValue;
  final double totalChange;
  final double totalChangePercent;
  final List<dynamic> holdings;
  final dynamic performance;

  PortfolioData({
    required this.totalValue,
    required this.totalChange,
    required this.totalChangePercent,
    this.holdings = const [],
    this.performance,
  });
}
