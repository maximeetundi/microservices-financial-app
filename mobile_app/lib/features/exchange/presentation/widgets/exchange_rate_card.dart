import 'package:flutter/material.dart';

class ExchangeRateCard extends StatelessWidget {
  final String fromCurrency;
  final String toCurrency;
  final double rate;
  final double fee;

  const ExchangeRateCard({
    super.key,
    required this.fromCurrency,
    required this.toCurrency,
    required this.rate,
    this.fee = 0.5,
  });

  @override
  Widget build(BuildContext context) {
    return Container(
      padding: const EdgeInsets.all(16),
      decoration: BoxDecoration(
        color: Colors.grey[100],
        borderRadius: BorderRadius.circular(16),
      ),
      child: Column(
        children: [
          Row(
            mainAxisAlignment: MainAxisAlignment.spaceBetween,
            children: [
              Text('Taux de change', style: TextStyle(color: Colors.grey[600])),
              Row(
                children: [
                  const Icon(Icons.access_time, size: 14, color: Colors.grey),
                  const SizedBox(width: 4),
                  Text('Mise à jour il y a 5 sec', style: TextStyle(fontSize: 12, color: Colors.grey[600])),
                ],
              ),
            ],
          ),
          const SizedBox(height: 12),
          Row(
            mainAxisAlignment: MainAxisAlignment.center,
            children: [
              Text(
                '1 $fromCurrency = ${rate.toStringAsFixed(rate < 1 ? 6 : 2)} $toCurrency',
                style: const TextStyle(fontSize: 20, fontWeight: FontWeight.bold),
              ),
            ],
          ),
          const SizedBox(height: 12),
          Row(
            mainAxisAlignment: MainAxisAlignment.spaceBetween,
            children: [
              Text('Frais de conversion', style: TextStyle(color: Colors.grey[600])),
              Text('$fee%', style: const TextStyle(fontWeight: FontWeight.w500)),
            ],
          ),
        ],
      ),
    );
  }
}
