import 'package:flutter/material.dart';

class TransferConfirmationSheet extends StatelessWidget {
  final String? transferType;
  final String? recipientName;
  final String? recipientEmail;
  final String? recipient;
  final double amount;
  final String currency;
  final dynamic fee;
  final VoidCallback onConfirm;

  const TransferConfirmationSheet({
    super.key,
    this.transferType,
    this.recipientName,
    this.recipientEmail,
    this.recipient,
    required this.amount,
    required this.currency,
    this.fee,
    required this.onConfirm,
  });

  double get feeValue {
    if (fee == null) return 0;
    if (fee is double) return fee;
    if (fee is int) return fee.toDouble();
    if (fee is String) {
      // Parse fee string like "Free" or "$2.50"
      if (fee == 'Free') return 0;
      final cleaned = fee.replaceAll(RegExp(r'[^\d.]'), '');
      return double.tryParse(cleaned) ?? 0;
    }
    return 0;
  }

  @override
  Widget build(BuildContext context) {
    return Container(
      padding: const EdgeInsets.all(24),
      decoration: const BoxDecoration(
        borderRadius: BorderRadius.vertical(top: Radius.circular(20)),
      ),
      child: Column(
        mainAxisSize: MainAxisSize.min,
        children: [
          Container(
            width: 40,
            height: 4,
            decoration: BoxDecoration(
              color: Colors.grey[300],
              borderRadius: BorderRadius.circular(2),
            ),
          ),
          const SizedBox(height: 24),
          
          const Text(
            'Confirm Transfer',
            style: TextStyle(fontSize: 20, fontWeight: FontWeight.bold),
          ),
          const SizedBox(height: 24),
          
          // Amount
          Container(
            width: double.infinity,
            padding: const EdgeInsets.all(20),
            decoration: BoxDecoration(
              color: Colors.blue[50],
              borderRadius: BorderRadius.circular(16),
            ),
            child: Column(
              children: [
                Text(
                  '${amount.toStringAsFixed(2)} $currency',
                  style: TextStyle(
                    fontSize: 32,
                    fontWeight: FontWeight.bold,
                    color: Colors.blue[700],
                  ),
                ),
              ],
            ),
          ),
          const SizedBox(height: 24),
          
          // Details
          if (transferType != null)
            _DetailRow(label: 'Type', value: _formatTransferType(transferType!)),
          _DetailRow(label: 'Recipient', value: recipientName ?? recipient ?? recipientEmail ?? 'Unknown'),
          if (recipientEmail != null && recipientName != null)
            _DetailRow(label: 'Email', value: recipientEmail!),
          _DetailRow(label: 'Fee', value: fee?.toString() ?? 'Free'),
          const Divider(),
          _DetailRow(
            label: 'Total', 
            value: '${(amount + feeValue).toStringAsFixed(2)} $currency',
            isBold: true,
          ),
          const SizedBox(height: 24),
          
          // Buttons
          Row(
            children: [
              Expanded(
                child: OutlinedButton(
                  onPressed: () => Navigator.pop(context),
                  style: OutlinedButton.styleFrom(
                    padding: const EdgeInsets.symmetric(vertical: 16),
                  ),
                  child: const Text('Cancel'),
                ),
              ),
              const SizedBox(width: 16),
              Expanded(
                child: ElevatedButton(
                  onPressed: onConfirm,
                  style: ElevatedButton.styleFrom(
                    padding: const EdgeInsets.symmetric(vertical: 16),
                  ),
                  child: const Text('Confirm'),
                ),
              ),
            ],
          ),
        ],
      ),
    );
  }

  String _formatTransferType(String type) {
    switch (type) {
      case 'crypto': return 'Crypto Transfer';
      case 'fiat': return 'Bank Transfer';
      case 'instant': return 'Instant Transfer';
      default: return type;
    }
  }
}

class _DetailRow extends StatelessWidget {
  final String label;
  final String value;
  final bool isBold;

  const _DetailRow({required this.label, required this.value, this.isBold = false});

  @override
  Widget build(BuildContext context) {
    return Padding(
      padding: const EdgeInsets.symmetric(vertical: 8),
      child: Row(
        mainAxisAlignment: MainAxisAlignment.spaceBetween,
        children: [
          Text(label, style: TextStyle(color: Colors.grey[600])),
          Text(
            value,
            style: TextStyle(fontWeight: isBold ? FontWeight.bold : FontWeight.normal),
          ),
        ],
      ),
    );
  }
}
