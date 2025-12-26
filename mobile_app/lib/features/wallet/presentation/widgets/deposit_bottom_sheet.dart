import 'package:flutter/material.dart';
import '../../../../core/services/wallet_api_service.dart';
import '../../../../core/theme/app_theme.dart';
import '../../../../core/widgets/custom_button.dart';

class DepositBottomSheet extends StatefulWidget {
  final String walletId;
  final String walletCurrency;
  final VoidCallback? onSuccess;

  const DepositBottomSheet({
    Key? key,
    required this.walletId,
    required this.walletCurrency,
    this.onSuccess,
  }) : super(key: key);

  @override
  State<DepositBottomSheet> createState() => _DepositBottomSheetState();
}

class _DepositBottomSheetState extends State<DepositBottomSheet> {
  final _walletApi = WalletApiService();
  final _amountController = TextEditingController(text: '5000');
  String _selectedMethod = 'orange';
  bool _isLoading = false;
  String? _error;
  bool _success = false;

  final List<Map<String, dynamic>> _paymentMethods = [
    // Mobile Money Section
    {
      'id': 'orange',
      'name': 'Orange Money',
      'icon': '🟠',
      'color': Colors.orange,
      'subtitle': 'Paiement instantané',
      'section': 'mobile',
    },
    {
      'id': 'mtn',
      'name': 'MTN Mobile Money',
      'icon': '🟡',
      'color': Colors.yellow.shade700,
      'subtitle': 'Paiement instantané',
      'section': 'mobile',
    },
    {
      'id': 'wave',
      'name': 'Wave',
      'icon': '🌊',
      'color': Colors.blue,
      'subtitle': 'Paiement instantané',
      'section': 'mobile',
    },
    // Bank Transfer Section
    {
      'id': 'bank',
      'name': 'Virement Bancaire',
      'icon': '🏦',
      'color': const Color(0xFF10B981), // Emerald green
      'subtitle': 'IBAN / RIB • 1-3 jours',
      'section': 'bank',
    },
    // Card Section
    {
      'id': 'card',
      'name': 'Carte Bancaire',
      'icon': '💳',
      'color': const Color(0xFF8B5CF6), // Purple
      'subtitle': 'Visa, Mastercard • Instantané',
      'section': 'card',
    },
  ];

  final List<int> _quickAmounts = [1000, 5000, 10000, 50000];

  @override
  void dispose() {
    _amountController.dispose();
    super.dispose();
  }

  Future<void> _submitDeposit() async {
    final amount = double.tryParse(_amountController.text);
    if (amount == null || amount <= 0) {
      setState(() => _error = 'Veuillez entrer un montant valide');
      return;
    }

    setState(() {
      _isLoading = true;
      _error = null;
    });

    try {
      await _walletApi.deposit(
        walletId: widget.walletId,
        amount: amount,
        paymentMethod: _selectedMethod,
      );

      setState(() {
        _success = true;
        _isLoading = false;
      });

      // Close after 2 seconds and call success callback
      await Future.delayed(const Duration(seconds: 2));
      if (mounted) {
        widget.onSuccess?.call();
        Navigator.pop(context);
      }
    } catch (e) {
      setState(() {
        _error = e.toString().replaceAll('Exception: ', '');
        _isLoading = false;
      });
    }
  }

  @override
  Widget build(BuildContext context) {
    return Container(
      decoration: BoxDecoration(
        color: Theme.of(context).scaffoldBackgroundColor,
        borderRadius: const BorderRadius.vertical(top: Radius.circular(20)),
      ),
      child: SingleChildScrollView(
        padding: const EdgeInsets.all(20),
        child: Column(
          mainAxisSize: MainAxisSize.min,
          crossAxisAlignment: CrossAxisAlignment.stretch,
          children: [
            // Handle
            Center(
              child: Container(
                width: 40,
                height: 4,
                margin: const EdgeInsets.only(bottom: 20),
                decoration: BoxDecoration(
                  color: Colors.grey.shade300,
                  borderRadius: BorderRadius.circular(2),
                ),
              ),
            ),

            // Title
            const Text(
              'Recharger',
              style: TextStyle(
                fontSize: 24,
                fontWeight: FontWeight.bold,
              ),
              textAlign: TextAlign.center,
            ),
            const SizedBox(height: 8),
            Text(
              'Portefeuille ${widget.walletCurrency}',
              style: TextStyle(
                fontSize: 14,
                color: Colors.grey.shade600,
              ),
              textAlign: TextAlign.center,
            ),
            const SizedBox(height: 24),

            // Success State
            if (_success) ...[
              Container(
                padding: const EdgeInsets.all(24),
                decoration: BoxDecoration(
                  color: Colors.green.shade50,
                  borderRadius: BorderRadius.circular(16),
                  border: Border.all(color: Colors.green.shade200),
                ),
                child: Column(
                  children: [
                    Icon(
                      Icons.check_circle,
                      size: 64,
                      color: Colors.green.shade600,
                    ),
                    const SizedBox(height: 16),
                    const Text(
                      'Dépôt réussi!',
                      style: TextStyle(
                        fontSize: 20,
                        fontWeight: FontWeight.bold,
                        color: Colors.green,
                      ),
                    ),
                    const SizedBox(height: 8),
                    Text(
                      '${_amountController.text} ${widget.walletCurrency} ajoutés à votre portefeuille',
                      style: TextStyle(color: Colors.grey.shade700),
                      textAlign: TextAlign.center,
                    ),
                  ],
                ),
              ),
            ] else ...[
              // Amount Input
              Text(
                'Montant',
                style: TextStyle(
                  fontWeight: FontWeight.w600,
                  color: Theme.of(context).brightness == Brightness.dark 
                      ? Colors.white 
                      : Colors.black87,
                ),
              ),
              const SizedBox(height: 8),
              TextField(
                controller: _amountController,
                keyboardType: TextInputType.number,
                style: TextStyle(
                  fontSize: 24,
                  fontWeight: FontWeight.bold,
                  color: Theme.of(context).brightness == Brightness.dark 
                      ? Colors.white 
                      : Colors.black87,
                ),
                decoration: InputDecoration(
                  suffixText: widget.walletCurrency,
                  suffixStyle: TextStyle(
                    fontSize: 18,
                    fontWeight: FontWeight.w600,
                    color: Theme.of(context).brightness == Brightness.dark 
                        ? Colors.white70 
                        : Colors.grey.shade600,
                  ),
                  border: OutlineInputBorder(
                    borderRadius: BorderRadius.circular(12),
                    borderSide: BorderSide(
                      color: Theme.of(context).brightness == Brightness.dark 
                          ? Colors.white24 
                          : Colors.grey.shade300,
                    ),
                  ),
                  enabledBorder: OutlineInputBorder(
                    borderRadius: BorderRadius.circular(12),
                    borderSide: BorderSide(
                      color: Theme.of(context).brightness == Brightness.dark 
                          ? Colors.white24 
                          : Colors.grey.shade300,
                    ),
                  ),
                  focusedBorder: OutlineInputBorder(
                    borderRadius: BorderRadius.circular(12),
                    borderSide: const BorderSide(
                      color: Color(0xFF6366F1),
                      width: 2,
                    ),
                  ),
                  filled: true,
                  fillColor: Theme.of(context).brightness == Brightness.dark 
                      ? const Color(0xFF1E293B) // Slate-800
                      : Colors.grey.shade100,
                ),
              ),
              const SizedBox(height: 12),

              // Quick Amount Buttons
              Row(
                children: _quickAmounts.map((amt) {
                  return Expanded(
                    child: Padding(
                      padding: const EdgeInsets.symmetric(horizontal: 4),
                      child: OutlinedButton(
                        onPressed: () {
                          _amountController.text = amt.toString();
                        },
                        style: OutlinedButton.styleFrom(
                          padding: const EdgeInsets.symmetric(vertical: 12),
                          shape: RoundedRectangleBorder(
                            borderRadius: BorderRadius.circular(8),
                          ),
                        ),
                        child: Text(
                          amt.toString(),
                          style: const TextStyle(fontWeight: FontWeight.w600),
                        ),
                      ),
                    ),
                  );
                }).toList(),
              ),
              const SizedBox(height: 24),

              // Payment Methods
              Text(
                'Méthode de paiement',
                style: TextStyle(
                  fontWeight: FontWeight.w600,
                  color: Theme.of(context).brightness == Brightness.dark 
                      ? Colors.white 
                      : Colors.black87,
                ),
              ),
              const SizedBox(height: 12),
              ...List.generate(_paymentMethods.length, (index) {
                final method = _paymentMethods[index];
                final isSelected = _selectedMethod == method['id'];
                return GestureDetector(
                  onTap: () {
                    setState(() => _selectedMethod = method['id']);
                  },
                  child: Container(
                    margin: const EdgeInsets.only(bottom: 8),
                    padding: const EdgeInsets.all(16),
                    decoration: BoxDecoration(
                      border: Border.all(
                        color: isSelected
                            ? method['color'] as Color
                            : Colors.grey.shade300,
                        width: isSelected ? 2 : 1,
                      ),
                      borderRadius: BorderRadius.circular(12),
                      color: isSelected
                          ? (method['color'] as Color).withOpacity(0.1)
                          : null,
                    ),
                    child: Row(
                      children: [
                        Text(
                          method['icon'] as String,
                          style: const TextStyle(fontSize: 28),
                        ),
                        const SizedBox(width: 16),
                        Expanded(
                          child: Column(
                            crossAxisAlignment: CrossAxisAlignment.start,
                            children: [
                              Text(
                                method['name'] as String,
                                style: const TextStyle(
                                  fontWeight: FontWeight.w600,
                                  fontSize: 16,
                                ),
                              ),
                              Text(
                                method['subtitle'] as String? ?? 'Paiement instantané',
                                style: TextStyle(
                                  color: Colors.grey.shade600,
                                  fontSize: 12,
                                ),
                              ),
                            ],
                          ),
                        ),
                        if (isSelected)
                          Container(
                            width: 24,
                            height: 24,
                            decoration: BoxDecoration(
                              color: method['color'] as Color,
                              shape: BoxShape.circle,
                            ),
                            child: const Icon(
                              Icons.check,
                              color: Colors.white,
                              size: 16,
                            ),
                          ),
                      ],
                    ),
                  ),
                );
              }),

              // Error Message
              if (_error != null) ...[
                const SizedBox(height: 16),
                Container(
                  padding: const EdgeInsets.all(12),
                  decoration: BoxDecoration(
                    color: Colors.red.shade50,
                    borderRadius: BorderRadius.circular(8),
                    border: Border.all(color: Colors.red.shade200),
                  ),
                  child: Text(
                    _error!,
                    style: TextStyle(color: Colors.red.shade700),
                  ),
                ),
              ],

              const SizedBox(height: 24),

              // Submit Button
              CustomButton(
                text: _isLoading
                    ? 'Traitement...'
                    : 'Déposer ${_amountController.text} ${widget.walletCurrency}',
                onPressed: _isLoading ? null : _submitDeposit,
                isLoading: _isLoading,
              ),
            ],

            const SizedBox(height: 20),
          ],
        ),
      ),
    );
  }
}
