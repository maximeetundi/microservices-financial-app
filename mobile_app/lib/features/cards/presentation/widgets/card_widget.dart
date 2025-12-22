import 'package:flutter/material.dart';

class CardWidget extends StatelessWidget {
  final Map<String, dynamic> card;
  final VoidCallback? onTap;

  const CardWidget({super.key, required this.card, this.onTap});

  @override
  Widget build(BuildContext context) {
    final cardNumber = card['card_number'] ?? '**** **** **** ****';
    final holderName = card['holder_name'] ?? 'CARD HOLDER';
    final expiryDate = card['expiry_date'] ?? '--/--';
    final isActive = card['status'] == 'active';

    return GestureDetector(
      onTap: onTap,
      child: Container(
        width: double.infinity,
        height: 200,
        decoration: BoxDecoration(
          gradient: LinearGradient(
            colors: isActive 
                ? [Colors.purple[700]!, Colors.purple[400]!]
                : [Colors.grey[700]!, Colors.grey[500]!],
            begin: Alignment.topLeft,
            end: Alignment.bottomRight,
          ),
          borderRadius: BorderRadius.circular(20),
          boxShadow: [
            BoxShadow(
              color: (isActive ? Colors.purple : Colors.grey).withOpacity(0.3),
              blurRadius: 20,
              offset: const Offset(0, 10),
            ),
          ],
        ),
        padding: const EdgeInsets.all(24),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          mainAxisAlignment: MainAxisAlignment.spaceBetween,
          children: [
            Row(
              mainAxisAlignment: MainAxisAlignment.spaceBetween,
              children: [
                const Text(
                  'CryptoBank',
                  style: TextStyle(color: Colors.white, fontSize: 18, fontWeight: FontWeight.bold),
                ),
                Row(
                  children: [
                    if (!isActive)
                      Container(
                        padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 4),
                        decoration: BoxDecoration(
                          color: Colors.red.withOpacity(0.2),
                          borderRadius: BorderRadius.circular(8),
                        ),
                        child: const Text('BLOQUÉE', style: TextStyle(color: Colors.red, fontSize: 10)),
                      ),
                    const SizedBox(width: 8),
                    const Icon(Icons.contactless, color: Colors.white, size: 28),
                  ],
                ),
              ],
            ),
            Text(
              _formatCardNumber(cardNumber),
              style: const TextStyle(color: Colors.white, fontSize: 22, letterSpacing: 2),
            ),
            Row(
              mainAxisAlignment: MainAxisAlignment.spaceBetween,
              children: [
                Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    const Text('TITULAIRE', style: TextStyle(color: Colors.white70, fontSize: 10)),
                    Text(holderName.toUpperCase(), style: const TextStyle(color: Colors.white, fontSize: 14)),
                  ],
                ),
                Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    const Text('EXPIRE', style: TextStyle(color: Colors.white70, fontSize: 10)),
                    Text(expiryDate, style: const TextStyle(color: Colors.white, fontSize: 14)),
                  ],
                ),
                const Icon(Icons.credit_card, color: Colors.white, size: 32),
              ],
            ),
          ],
        ),
      ),
    );
  }

  String _formatCardNumber(String number) {
    return number.replaceAllMapped(RegExp(r'.{4}'), (m) => '${m.group(0)} ').trim();
  }
}
