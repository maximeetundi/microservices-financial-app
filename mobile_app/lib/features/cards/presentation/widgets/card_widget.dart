import 'package:flutter/material.dart';
import 'package:google_fonts/google_fonts.dart';

/// Card widget matching web frontend credit card design exactly
/// Supports virtual, physical card types with web-matching gradients
class CardWidget extends StatelessWidget {
  final Map<String, dynamic> card;
  final VoidCallback? onTap;

  const CardWidget({super.key, required this.card, this.onTap});

  @override
  Widget build(BuildContext context) {
    final cardNumber = card['card_number'] ?? '**** **** **** ****';
    final holderName = card['cardholder_name'] ?? card['holder_name'] ?? 'CARD HOLDER';
    final expiryMonth = card['expiry_month']?.toString().padLeft(2, '0') ?? '12';
    final expiryYear = card['expiry_year']?.toString() ?? '28';
    final expiryDate = '$expiryMonth/${expiryYear.length > 2 ? expiryYear.substring(expiryYear.length - 2) : expiryYear}';
    final isActive = card['status'] == 'active';
    final isVirtual = card['is_virtual'] == true || card['card_type'] == 'virtual';
    final balance = (card['balance'] as num?)?.toDouble() ?? 0.0;
    final currency = card['currency'] ?? 'USD';

    return GestureDetector(
      onTap: onTap,
      child: AspectRatio(
        aspectRatio: 16 / 10,
        child: Container(
          decoration: BoxDecoration(
            // Match web: credit-card-virtual or credit-card-physical
            gradient: isVirtual
                ? const LinearGradient(
                    begin: Alignment.topLeft,
                    end: Alignment.bottomRight,
                    colors: [Color(0xFF6366F1), Color(0xFF8B5CF6)],
                  )
                : const LinearGradient(
                    begin: Alignment.topLeft,
                    end: Alignment.bottomRight,
                    colors: [Color(0xFF1E1E2F), Color(0xFF3D3D5C)],
                  ),
            borderRadius: BorderRadius.circular(24),
            boxShadow: [
              BoxShadow(
                color: Colors.black.withOpacity(0.3),
                blurRadius: 40,
                offset: const Offset(0, 20),
              ),
            ],
          ),
          child: Stack(
            children: [
              // Radial glow effect (matching web ::before)
              Positioned(
                top: -50,
                right: -50,
                child: Container(
                  width: 150,
                  height: 150,
                  decoration: BoxDecoration(
                    shape: BoxShape.circle,
                    gradient: RadialGradient(
                      colors: [
                        const Color(0xFF6366F1).withOpacity(0.3),
                        Colors.transparent,
                      ],
                    ),
                  ),
                ),
              ),
              
              // Shine effect on hover (simulated with gradient)
              Positioned.fill(
                child: Container(
                  decoration: BoxDecoration(
                    borderRadius: BorderRadius.circular(24),
                    gradient: LinearGradient(
                      begin: Alignment.topLeft,
                      end: Alignment.bottomRight,
                      colors: [
                        Colors.white.withOpacity(0.0),
                        Colors.white.withOpacity(0.05),
                        Colors.white.withOpacity(0.0),
                      ],
                    ),
                  ),
                ),
              ),
              
              // Card content
              Padding(
                padding: const EdgeInsets.all(24),
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    // Top row: Status badge and card type icon
                    Row(
                      mainAxisAlignment: MainAxisAlignment.spaceBetween,
                      children: [
                        // Status badge (matching web backdrop-blur style)
                        Container(
                          padding: const EdgeInsets.symmetric(
                            horizontal: 10,
                            vertical: 5,
                          ),
                          decoration: BoxDecoration(
                            color: Colors.white.withOpacity(0.2),
                            borderRadius: BorderRadius.circular(20),
                            border: Border.all(
                              color: Colors.white.withOpacity(0.1),
                            ),
                          ),
                          child: Text(
                            isActive ? 'Active' : 'Gelée',
                            style: GoogleFonts.inter(
                              fontSize: 12,
                              fontWeight: FontWeight.w600,
                              color: Colors.white,
                            ),
                          ),
                        ),
                        Text(
                          isVirtual ? '🌐' : '💳',
                          style: const TextStyle(
                            fontSize: 24,
                            color: Colors.white70,
                          ),
                        ),
                      ],
                    ),
                    
                    const Spacer(),
                    
                    // Card number (with letter spacing matching web tracking-[0.15em])
                    Text(
                      _formatCardNumber(cardNumber),
                      style: GoogleFonts.robotoMono(
                        fontSize: 20,
                        fontWeight: FontWeight.w500,
                        letterSpacing: 3.0,
                        color: Colors.white,
                        shadows: [
                          Shadow(
                            color: Colors.black.withOpacity(0.2),
                            blurRadius: 4,
                          ),
                        ],
                      ),
                    ),
                    const SizedBox(height: 4),
                    // Card holder name
                    Text(
                      holderName.isEmpty ? 'Ma Carte' : holderName.toUpperCase(),
                      style: GoogleFonts.inter(
                        fontSize: 10,
                        fontWeight: FontWeight.w500,
                        letterSpacing: 2.0,
                        color: Colors.white.withOpacity(0.6),
                      ),
                    ),
                    
                    const Spacer(),
                    
                    // Bottom row: Balance and Expiry
                    Row(
                      mainAxisAlignment: MainAxisAlignment.spaceBetween,
                      crossAxisAlignment: CrossAxisAlignment.end,
                      children: [
                        Column(
                          crossAxisAlignment: CrossAxisAlignment.start,
                          children: [
                            Text(
                              'Solde',
                              style: GoogleFonts.inter(
                                fontSize: 10,
                                fontWeight: FontWeight.w500,
                                letterSpacing: 1.0,
                                color: Colors.white.withOpacity(0.5),
                              ),
                            ),
                            const SizedBox(height: 2),
                            Text(
                              _formatBalance(balance, currency),
                              style: GoogleFonts.inter(
                                fontSize: 18,
                                fontWeight: FontWeight.bold,
                                color: Colors.white,
                              ),
                            ),
                          ],
                        ),
                        Column(
                          crossAxisAlignment: CrossAxisAlignment.end,
                          children: [
                            Text(
                              'Expire',
                              style: GoogleFonts.inter(
                                fontSize: 10,
                                fontWeight: FontWeight.w500,
                                letterSpacing: 1.0,
                                color: Colors.white.withOpacity(0.5),
                              ),
                            ),
                            const SizedBox(height: 2),
                            Text(
                              expiryDate,
                              style: GoogleFonts.robotoMono(
                                fontSize: 16,
                                fontWeight: FontWeight.w500,
                                color: Colors.white,
                              ),
                            ),
                          ],
                        ),
                      ],
                    ),
                  ],
                ),
              ),
            ],
          ),
        ),
      ),
    );
  }

  String _formatCardNumber(String number) {
    // Show last 4 digits with masked prefix
    if (number.length >= 4) {
      final last4 = number.substring(number.length - 4);
      return '•••• •••• •••• $last4';
    }
    return number.replaceAllMapped(RegExp(r'.{4}'), (m) => '${m.group(0)} ').trim();
  }

  String _formatBalance(double balance, String currency) {
    if (currency == 'USD' || currency == 'EUR' || currency == 'GBP') {
      final symbol = currency == 'USD' ? '\$' : (currency == 'EUR' ? '€' : '£');
      return '$symbol${balance.toStringAsFixed(2)}';
    }
    return '${balance.toStringAsFixed(2)} $currency';
  }
}
