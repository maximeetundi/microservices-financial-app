import 'package:flutter/material.dart';

class QuickExchangeAmounts extends StatelessWidget {
  final List<double> amounts;
  final double? selectedAmount;
  final String currency;
  final ValueChanged<double> onSelected;

  const QuickExchangeAmounts({
    super.key,
    this.amounts = const [100, 250, 500, 1000],
    this.selectedAmount,
    required this.currency,
    required this.onSelected,
  });

  @override
  Widget build(BuildContext context) {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Text('Montants rapides', style: TextStyle(color: Colors.grey[600])),
        const SizedBox(height: 8),
        Wrap(
          spacing: 8,
          runSpacing: 8,
          children: amounts.map((amount) {
            final isSelected = amount == selectedAmount;
            return InkWell(
              onTap: () => onSelected(amount),
              borderRadius: BorderRadius.circular(8),
              child: Container(
                padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 8),
                decoration: BoxDecoration(
                  color: isSelected 
                      ? Theme.of(context).primaryColor 
                      : Colors.grey[100],
                  borderRadius: BorderRadius.circular(8),
                ),
                child: Text(
                  '${amount.toStringAsFixed(0)} $currency',
                  style: TextStyle(
                    color: isSelected ? Colors.white : Colors.black,
                    fontWeight: FontWeight.w500,
                  ),
                ),
              ),
            );
          }).toList(),
        ),
      ],
    );
  }
}
