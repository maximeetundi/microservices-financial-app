import 'package:flutter/material.dart';
import '../../core/services/ticket_api_service.dart';
import '../../core/services/wallet_api_service.dart';
import '../../features/auth/presentation/pages/pin_verify_dialog.dart';

class PurchaseTicketScreen extends StatefulWidget {
  final Map<String, dynamic> event;
  final Map<String, dynamic> tier;

  const PurchaseTicketScreen({
    super.key,
    required this.event,
    required this.tier,
  });

  @override
  State<PurchaseTicketScreen> createState() => _PurchaseTicketScreenState();
}

class _PurchaseTicketScreenState extends State<PurchaseTicketScreen> {
  final TicketApiService _ticketApi = TicketApiService();
  final WalletApiService _walletApi = WalletApiService();
  final _formKey = GlobalKey<FormState>();
  
  List<dynamic> _wallets = [];
  int _quantity = 1;

  // ... (existing initState)

  // Helper to check limits
  int get _maxAllowed {
    int max = 100; // safe default max
    int maxPerUser = widget.tier['max_per_user'] ?? 0;
    int total = widget.tier['quantity'] ?? -1;
    int sold = widget.tier['sold'] ?? 0;
    
    // Check remaining
    if (total != -1) {
      int remaining = total - sold;
      if (remaining < max) max = remaining;
    }
    
    // Check per user limit (local batch limit)
    // We don't know existing purchases here easily without API call.
    // So we just limit the BATCH size to MaxPerUser.
    if (maxPerUser > 0) {
      if (maxPerUser < max) max = maxPerUser;
    }
    
    return max;
  }

  void _incrementQuantity() {
    if (_quantity < _maxAllowed) {
      setState(() => _quantity++);
    }
  }

  void _decrementQuantity() {
    if (_quantity > 1) {
      setState(() => _quantity--);
    }
  }

 // ... existing methods

  Widget _buildTierCard() {
    return Container(
      padding: const EdgeInsets.all(20),
      decoration: BoxDecoration(
        color: Colors.white.withOpacity(0.05),
        borderRadius: BorderRadius.circular(16),
        border: Border.all(
          color: _hexToColor(widget.tier['color'] ?? '#6366f1'),
          width: 2,
        ),
      ),
      child: Column(
        children: [
          Row(
            children: [
              Text(widget.tier['icon'] ?? '🎫', style: const TextStyle(fontSize: 40)),
              const SizedBox(width: 16),
              Expanded(
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    Text(
                      widget.tier['name'] ?? 'Ticket',
                      style: const TextStyle(
                        color: Colors.white,
                        fontSize: 20,
                        fontWeight: FontWeight.bold,
                      ),
                    ),
                    if (widget.tier['description']?.isNotEmpty ?? false)
                      Text(
                        widget.tier['description'],
                        style: TextStyle(color: Colors.white.withOpacity(0.7)),
                      ),
                  ],
                ),
              ),
              Column(
                crossAxisAlignment: CrossAxisAlignment.end,
                children: [
                  Text(
                    _formatAmount(widget.tier['price'] ?? 0),
                    style: const TextStyle(
                      color: Colors.white,
                      fontSize: 24,
                      fontWeight: FontWeight.bold,
                    ),
                  ),
                  Text(
                    widget.event['currency'] ?? 'XOF',
                    style: TextStyle(color: Colors.white.withOpacity(0.6)),
                  ),
                ],
              ),
            ],
          ),
          const SizedBox(height: 16),
          // Quantity Selector
          Container(
            padding: const EdgeInsets.symmetric(vertical: 8, horizontal: 16),
            decoration: BoxDecoration(
              color: Colors.black.withOpacity(0.2),
              borderRadius: BorderRadius.circular(12),
            ),
            child: Row(
              mainAxisAlignment: MainAxisAlignment.spaceBetween,
              children: [
                const Text(
                  'Quantité',
                  style: TextStyle(color: Colors.white, fontSize: 16),
                ),
                Row(
                  children: [
                    _buildQuantityButton(Icons.remove, _decrementQuantity, _quantity > 1),
                    SizedBox(
                      width: 40,
                      child: Text(
                        '$_quantity',
                        textAlign: TextAlign.center,
                        style: const TextStyle(color: Colors.white, fontSize: 18, fontWeight: FontWeight.bold),
                      ),
                    ),
                    _buildQuantityButton(Icons.add, _incrementQuantity, _quantity < _maxAllowed),
                  ],
                ),
              ],
            ),
          ),
        ],
      ),
    );
  }

  Widget _buildQuantityButton(IconData icon, VoidCallback onPressed, bool enabled) {
    return Container(
      margin: const EdgeInsets.symmetric(horizontal: 4),
      decoration: BoxDecoration(
        color: enabled ? const Color(0xFF6366f1) : Colors.white.withOpacity(0.1),
        shape: BoxShape.circle,
      ),
      child: IconButton(
        icon: Icon(icon, color: enabled ? Colors.white : Colors.white38, size: 20),
        onPressed: enabled ? onPressed : null,
        constraints: const BoxConstraints(minWidth: 36, minHeight: 36),
        padding: EdgeInsets.zero,
      ),
    );
  }

  // ... _buildFormFields ...

  Widget _buildPurchaseButton() {
    final totalPrice = (widget.tier['price'] ?? 0) * _quantity;
    
    return SizedBox(
      width: double.infinity,
      child: ElevatedButton(
        onPressed: _purchasing ? null : _purchaseTicket,
        style: ElevatedButton.styleFrom(
          backgroundColor: const Color(0xFF6366f1),
          disabledBackgroundColor: Colors.grey,
          padding: const EdgeInsets.all(18),
          shape: RoundedRectangleBorder(
            borderRadius: BorderRadius.circular(12),
          ),
        ),
        child: _purchasing
            ? const SizedBox(
                height: 20,
                width: 20,
                child: CircularProgressIndicator(
                  strokeWidth: 2,
                  valueColor: AlwaysStoppedAnimation<Color>(Colors.white),
                ),
              )
            : Text(
                'Acheter $_quantity ticket${_quantity > 1 ? 's' : ''} - ${_formatAmount(totalPrice)} ${widget.event['currency'] ?? 'XOF'}',
                style: const TextStyle(
                  fontSize: 18,
                  fontWeight: FontWeight.bold,
                  color: Colors.white,
                ),
              ),
      ),
    );
  }

  String _formatAmount(dynamic amount) {
    final num = (amount is int) ? amount : (amount as double).toInt();
    return num.toString().replaceAllMapped(
      RegExp(r'(\d{1,3})(?=(\d{3})+(?!\d))'),
      (m) => '${m[1]} ',
    );
  }

  Color _hexToColor(String hex) {
    hex = hex.replaceFirst('#', '');
    if (hex.length == 6) hex = 'FF$hex';
    return Color(int.parse(hex, radix: 16));
  }
}
