import 'package:flutter/material.dart';
import 'dart:convert';
import '../../core/services/api_service.dart';
import '../../core/utils/currency_formatter.dart';

/// Screen for scanning QR codes and making payments
class ScanPayScreen extends StatefulWidget {
  const ScanPayScreen({Key? key}) : super(key: key);

  @override
  State<ScanPayScreen> createState() => _ScanPayScreenState();
}

class _ScanPayScreenState extends State<ScanPayScreen> {
  bool _scanning = true;
  dynamic _payment;
  String? _error;
  bool _processing = false;
  String? _selectedWalletId;
  double? _customAmount;
  List<dynamic> _wallets = [];

  @override
  void initState() {
    super.initState();
    _loadWallets();
  }

  Future<void> _loadWallets() async {
    try {
      final wallets = await ApiService.getWallets();
      setState(() => _wallets = wallets ?? []);
    } catch (e) {
      debugPrint('Error loading wallets: $e');
    }
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: const Color(0xFF0F0F23),
      appBar: AppBar(
        title: const Text('Scanner & Payer'),
        backgroundColor: Colors.transparent,
        elevation: 0,
      ),
      body: _scanning
          ? _buildScanner()
          : _payment != null
              ? _buildPaymentDetails()
              : _buildError(),
    );
  }

  Widget _buildScanner() {
    // In production, use mobile_scanner or qr_code_scanner package
    return Column(
      children: [
        Expanded(
          child: Container(
            margin: const EdgeInsets.all(24),
            decoration: BoxDecoration(
              color: Colors.black,
              borderRadius: BorderRadius.circular(24),
              border: Border.all(color: const Color(0xFF6366F1), width: 3),
            ),
            child: Stack(
              children: [
                // Camera preview placeholder
                Center(
                  child: Column(
                    mainAxisAlignment: MainAxisAlignment.center,
                    children: [
                      Icon(
                        Icons.qr_code_scanner,
                        size: 100,
                        color: Colors.white.withOpacity(0.3),
                      ),
                      const SizedBox(height: 16),
                      Text(
                        'Placez le QR code dans le cadre',
                        style: TextStyle(
                          color: Colors.white.withOpacity(0.6),
                          fontSize: 16,
                        ),
                      ),
                    ],
                  ),
                ),
                // Scanning animation
                Positioned.fill(
                  child: CustomPaint(
                    painter: _ScannerOverlayPainter(),
                  ),
                ),
              ],
            ),
          ),
        ),
        // Manual entry option
        Padding(
          padding: const EdgeInsets.all(24),
          child: Column(
            children: [
              const Text(
                'Ou entrez le code manuellement',
                style: TextStyle(color: Colors.grey),
              ),
              const SizedBox(height: 12),
              TextField(
                onSubmitted: _onCodeEntered,
                decoration: InputDecoration(
                  hintText: 'Code de paiement (pay_xxx)',
                  hintStyle: TextStyle(color: Colors.white.withOpacity(0.3)),
                  filled: true,
                  fillColor: Colors.white.withOpacity(0.05),
                  border: OutlineInputBorder(
                    borderRadius: BorderRadius.circular(12),
                    borderSide: BorderSide.none,
                  ),
                  suffixIcon: const Icon(Icons.arrow_forward, color: Color(0xFF6366F1)),
                ),
                style: const TextStyle(color: Colors.white),
              ),
              const SizedBox(height: 24),
              // Demo button for testing
              TextButton.icon(
                onPressed: () => _onCodeEntered('pay_demo'),
                icon: const Icon(Icons.bug_report),
                label: const Text('Test avec démo'),
                style: TextButton.styleFrom(foregroundColor: Colors.grey),
              ),
            ],
          ),
        ),
      ],
    );
  }

  Widget _buildPaymentDetails() {
    final amount = _payment['amount'];
    final currency = _payment['currency'] ?? 'EUR';
    final isVariable = amount == null;
    final compatibleWallets = _wallets.where((w) => w['currency'] == currency).toList();

    return SingleChildScrollView(
      padding: const EdgeInsets.all(24),
      child: Column(
        children: [
          // Merchant info
          Container(
            padding: const EdgeInsets.all(20),
            decoration: BoxDecoration(
              color: const Color(0xFF1A1A3E),
              borderRadius: BorderRadius.circular(16),
            ),
            child: Column(
              children: [
                const CircleAvatar(
                  radius: 30,
                  backgroundColor: Color(0xFF6366F1),
                  child: Icon(Icons.store, color: Colors.white, size: 30),
                ),
                const SizedBox(height: 12),
                Text(
                  _payment['title'] ?? 'Paiement',
                  style: const TextStyle(
                    color: Colors.white,
                    fontSize: 22,
                    fontWeight: FontWeight.bold,
                  ),
                ),
                if (_payment['description'] != null) ...[
                  const SizedBox(height: 8),
                  Text(
                    _payment['description'],
                    style: TextStyle(color: Colors.white.withOpacity(0.6)),
                    textAlign: TextAlign.center,
                  ),
                ],
              ],
            ),
          ),
          const SizedBox(height: 24),

          // Amount
          Container(
            padding: const EdgeInsets.all(20),
            decoration: BoxDecoration(
              color: const Color(0xFF1A1A3E),
              borderRadius: BorderRadius.circular(16),
            ),
            child: Column(
              children: [
                if (!isVariable)
                  Text(
                    CurrencyFormatter.format(amount, currency),
                    style: const TextStyle(
                      color: Color(0xFF22C55E),
                      fontSize: 36,
                      fontWeight: FontWeight.bold,
                    ),
                  )
                else ...[
                  const Text(
                    'Entrez le montant',
                    style: TextStyle(color: Colors.grey),
                  ),
                  const SizedBox(height: 12),
                  TextField(
                    onChanged: (v) => _customAmount = double.tryParse(v),
                    keyboardType: TextInputType.number,
                    textAlign: TextAlign.center,
                    style: const TextStyle(
                      color: Color(0xFF22C55E),
                      fontSize: 32,
                      fontWeight: FontWeight.bold,
                    ),
                    decoration: InputDecoration(
                      hintText: '0.00',
                      hintStyle: TextStyle(color: Colors.grey.withOpacity(0.5)),
                      border: InputBorder.none,
                      suffix: Text(
                        currency,
                        style: const TextStyle(color: Colors.grey, fontSize: 18),
                      ),
                    ),
                  ),
                  if (_payment['min_amount'] != null || _payment['max_amount'] != null)
                    Padding(
                      padding: const EdgeInsets.only(top: 8),
                      child: Text(
                        'Min: ${_payment['min_amount'] ?? 0} - Max: ${_payment['max_amount'] ?? '∞'}',
                        style: const TextStyle(color: Colors.grey, fontSize: 12),
                      ),
                    ),
                ],
              ],
            ),
          ),
          const SizedBox(height: 24),

          // Wallet selection
          Container(
            padding: const EdgeInsets.all(20),
            decoration: BoxDecoration(
              color: const Color(0xFF1A1A3E),
              borderRadius: BorderRadius.circular(16),
            ),
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                const Text(
                  'Payer depuis',
                  style: TextStyle(color: Colors.grey),
                ),
                const SizedBox(height: 12),
                if (compatibleWallets.isEmpty)
                  Text(
                    'Aucun portefeuille $currency disponible',
                    style: const TextStyle(color: Colors.red),
                  )
                else
                  ...compatibleWallets.map((wallet) => _WalletOption(
                        wallet: wallet,
                        selected: wallet['id'] == _selectedWalletId,
                        enabled: (wallet['balance'] ?? 0) >= (amount ?? _customAmount ?? 0),
                        onTap: () => setState(() => _selectedWalletId = wallet['id']),
                      )),
              ],
            ),
          ),
          const SizedBox(height: 32),

          // Pay button
          SizedBox(
            width: double.infinity,
            child: ElevatedButton(
              onPressed: _canPay ? _processPayment : null,
              style: ElevatedButton.styleFrom(
                backgroundColor: const Color(0xFF22C55E),
                padding: const EdgeInsets.symmetric(vertical: 18),
                shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(12)),
                disabledBackgroundColor: Colors.grey,
              ),
              child: _processing
                  ? const SizedBox(
                      width: 24,
                      height: 24,
                      child: CircularProgressIndicator(color: Colors.white, strokeWidth: 2),
                    )
                  : Text(
                      'Payer ${CurrencyFormatter.format(amount ?? _customAmount ?? 0, currency)}',
                      style: const TextStyle(
                        fontSize: 18,
                        fontWeight: FontWeight.bold,
                        color: Colors.white,
                      ),
                    ),
            ),
          ),
          const SizedBox(height: 16),

          // Cancel
          TextButton(
            onPressed: () => setState(() {
              _scanning = true;
              _payment = null;
            }),
            child: const Text('Annuler'),
          ),
        ],
      ),
    );
  }

  Widget _buildError() {
    return Center(
      child: Column(
        mainAxisAlignment: MainAxisAlignment.center,
        children: [
          const Icon(Icons.error_outline, size: 64, color: Colors.red),
          const SizedBox(height: 16),
          Text(
            _error ?? 'Paiement introuvable',
            style: const TextStyle(color: Colors.white, fontSize: 18),
          ),
          const SizedBox(height: 24),
          ElevatedButton(
            onPressed: () => setState(() {
              _scanning = true;
              _error = null;
            }),
            child: const Text('Réessayer'),
          ),
        ],
      ),
    );
  }

  bool get _canPay {
    if (_selectedWalletId == null) return false;
    final amount = _payment['amount'] ?? _customAmount;
    if (amount == null || amount <= 0) return false;
    final wallet = _wallets.firstWhere(
      (w) => w['id'] == _selectedWalletId,
      orElse: () => null,
    );
    if (wallet == null || (wallet['balance'] ?? 0) < amount) return false;
    return true;
  }

  Future<void> _onCodeEntered(String code) async {
    // Extract payment ID from QR data or direct input
    String paymentId = code;
    
    // Try to parse as JSON (QR code data)
    try {
      final data = jsonDecode(code);
      paymentId = data['payment_id'] ?? code;
    } catch (_) {
      // Not JSON, use as-is
    }

    // Remove prefix if present
    if (paymentId.contains('/pay/')) {
      paymentId = paymentId.split('/pay/').last;
    }

    try {
      final payment = await ApiService.getPaymentDetails(paymentId);
      setState(() {
        _payment = payment;
        _scanning = false;
      });
    } catch (e) {
      setState(() {
        _error = e.toString();
        _scanning = false;
      });
    }
  }

  Future<void> _processPayment() async {
    if (!_canPay) return;

    setState(() => _processing = true);
    try {
      await ApiService.payPayment(
        _payment['payment_id'],
        _selectedWalletId!,
        _payment['amount'] ?? _customAmount!,
      );
      
      if (mounted) {
        _showSuccessDialog();
      }
    } catch (e) {
      if (mounted) {
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(content: Text('Erreur: $e'), backgroundColor: Colors.red),
        );
      }
    } finally {
      setState(() => _processing = false);
    }
  }

  void _showSuccessDialog() {
    showDialog(
      context: context,
      barrierDismissible: false,
      builder: (context) => AlertDialog(
        backgroundColor: const Color(0xFF1A1A3E),
        shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(20)),
        content: Column(
          mainAxisSize: MainAxisSize.min,
          children: [
            const Icon(Icons.check_circle, color: Color(0xFF22C55E), size: 80),
            const SizedBox(height: 24),
            const Text(
              'Paiement réussi!',
              style: TextStyle(
                color: Colors.white,
                fontSize: 22,
                fontWeight: FontWeight.bold,
              ),
            ),
            const SizedBox(height: 8),
            Text(
              CurrencyFormatter.format(
                _payment['amount'] ?? _customAmount,
                _payment['currency'] ?? 'EUR',
              ),
              style: const TextStyle(color: Color(0xFF22C55E), fontSize: 28),
            ),
            const SizedBox(height: 24),
            SizedBox(
              width: double.infinity,
              child: ElevatedButton(
                onPressed: () {
                  Navigator.of(context).pop();
                  Navigator.of(context).pop();
                },
                style: ElevatedButton.styleFrom(
                  backgroundColor: const Color(0xFF6366F1),
                  padding: const EdgeInsets.symmetric(vertical: 14),
                ),
                child: const Text('Terminé'),
              ),
            ),
          ],
        ),
      ),
    );
  }
}

class _WalletOption extends StatelessWidget {
  final dynamic wallet;
  final bool selected;
  final bool enabled;
  final VoidCallback onTap;

  const _WalletOption({
    required this.wallet,
    required this.selected,
    required this.enabled,
    required this.onTap,
  });

  @override
  Widget build(BuildContext context) {
    return InkWell(
      onTap: enabled ? onTap : null,
      borderRadius: BorderRadius.circular(12),
      child: Container(
        padding: const EdgeInsets.all(12),
        margin: const EdgeInsets.only(bottom: 8),
        decoration: BoxDecoration(
          color: selected ? const Color(0xFF6366F1).withOpacity(0.1) : Colors.transparent,
          borderRadius: BorderRadius.circular(12),
          border: Border.all(
            color: selected ? const Color(0xFF6366F1) : Colors.white.withOpacity(0.1),
            width: 2,
          ),
        ),
        child: Row(
          children: [
            Radio<bool>(
              value: true,
              groupValue: selected,
              onChanged: enabled ? (_) => onTap() : null,
              activeColor: const Color(0xFF6366F1),
            ),
            const SizedBox(width: 12),
            Expanded(
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  Text(
                    wallet['currency'],
                    style: TextStyle(
                      color: enabled ? Colors.white : Colors.grey,
                      fontWeight: FontWeight.w600,
                    ),
                  ),
                  Text(
                    CurrencyFormatter.format(wallet['balance'], wallet['currency']),
                    style: TextStyle(
                      color: enabled ? Colors.white.withOpacity(0.6) : Colors.grey,
                      fontSize: 13,
                    ),
                  ),
                ],
              ),
            ),
            if (!enabled)
              const Text(
                'Solde insuffisant',
                style: TextStyle(color: Colors.red, fontSize: 12),
              ),
          ],
        ),
      ),
    );
  }
}

class _ScannerOverlayPainter extends CustomPainter {
  @override
  void paint(Canvas canvas, Size size) {
    final paint = Paint()
      ..color = const Color(0xFF6366F1)
      ..strokeWidth = 4
      ..style = PaintingStyle.stroke;

    final cornerLength = 40.0;
    final rect = Rect.fromCenter(
      center: Offset(size.width / 2, size.height / 2),
      width: 200,
      height: 200,
    );

    // Top left corner
    canvas.drawLine(rect.topLeft, rect.topLeft + Offset(cornerLength, 0), paint);
    canvas.drawLine(rect.topLeft, rect.topLeft + Offset(0, cornerLength), paint);

    // Top right corner
    canvas.drawLine(rect.topRight, rect.topRight + Offset(-cornerLength, 0), paint);
    canvas.drawLine(rect.topRight, rect.topRight + Offset(0, cornerLength), paint);

    // Bottom left corner
    canvas.drawLine(rect.bottomLeft, rect.bottomLeft + Offset(cornerLength, 0), paint);
    canvas.drawLine(rect.bottomLeft, rect.bottomLeft + Offset(0, -cornerLength), paint);

    // Bottom right corner
    canvas.drawLine(rect.bottomRight, rect.bottomRight + Offset(-cornerLength, 0), paint);
    canvas.drawLine(rect.bottomRight, rect.bottomRight + Offset(0, -cornerLength), paint);
  }

  @override
  bool shouldRepaint(covariant CustomPainter oldDelegate) => false;
}
