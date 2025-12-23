import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';
import 'package:mobile_scanner/mobile_scanner.dart';
import 'package:flutter_bloc/flutter_bloc.dart';
import 'package:image_picker/image_picker.dart';
import 'dart:convert';

import '../../../core/widgets/security_confirmation.dart';
import '../wallet/presentation/bloc/wallet_bloc.dart';

/// Screen for scanning QR codes and making payments
class ScanPayScreen extends StatefulWidget {
  const ScanPayScreen({Key? key}) : super(key: key);

  @override
  State<ScanPayScreen> createState() => _ScanPayScreenState();
}

class _ScanPayScreenState extends State<ScanPayScreen> with WidgetsBindingObserver {
  MobileScannerController? _cameraController;
  final ImagePicker _imagePicker = ImagePicker();
  bool _scanning = true;
  bool _pickingImage = false;
  Map<String, dynamic>? _payment;
  String? _error;
  bool _processing = false;
  String? _selectedWalletId;
  double? _customAmount;
  bool _hasScanned = false;

  @override
  void initState() {
    super.initState();
    WidgetsBinding.instance.addObserver(this);
    _initializeCamera();
    context.read<WalletBloc>().add(LoadWalletsEvent());
  }

  void _initializeCamera() {
    _cameraController = MobileScannerController(
      facing: CameraFacing.back,
      torchEnabled: false,
    );
  }

  @override
  void dispose() {
    WidgetsBinding.instance.removeObserver(this);
    _cameraController?.dispose();
    super.dispose();
  }

  @override
  void didChangeAppLifecycleState(AppLifecycleState state) {
    if (state == AppLifecycleState.inactive) {
      _cameraController?.stop();
    } else if (state == AppLifecycleState.resumed && _scanning) {
      _cameraController?.start();
    }
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: const Color(0xFFF5F7FA),
      appBar: AppBar(
        backgroundColor: Colors.transparent,
        elevation: 0,
        leading: IconButton(
          icon: const Icon(Icons.arrow_back_ios, color: Color(0xFF1a1a2e)),
          onPressed: () => context.pop(),
        ),
        title: const Text(
          'Scanner & Payer 📷',
          style: TextStyle(
            color: Color(0xFF1a1a2e),
            fontWeight: FontWeight.bold,
            fontSize: 20,
          ),
        ),
        centerTitle: true,
        actions: [
          if (_scanning && _cameraController != null)
            IconButton(
              icon: ValueListenableBuilder(
                valueListenable: _cameraController!.torchState,
                builder: (context, state, child) {
                  return Icon(
                    state == TorchState.on ? Icons.flash_on : Icons.flash_off,
                    color: const Color(0xFF1a1a2e),
                  );
                },
              ),
              onPressed: () => _cameraController?.toggleTorch(),
            ),
        ],
      ),
      body: _scanning
          ? _buildScanner()
          : _payment != null
              ? _buildPaymentDetails()
              : _buildError(),
    );
  }

  Widget _buildScanner() {
    return Column(
      children: [
        Expanded(
          child: Container(
            margin: const EdgeInsets.all(24),
            decoration: BoxDecoration(
              borderRadius: BorderRadius.circular(24),
              boxShadow: [
                BoxShadow(
                  color: const Color(0xFF667eea).withOpacity(0.3),
                  blurRadius: 20,
                  offset: const Offset(0, 8),
                ),
              ],
            ),
            child: ClipRRect(
              borderRadius: BorderRadius.circular(24),
              child: Stack(
                children: [
                  // Real Camera Scanner
                  MobileScanner(
                    controller: _cameraController,
                    onDetect: (capture) {
                      if (_hasScanned) return;
                      final List<Barcode> barcodes = capture.barcodes;
                      for (final barcode in barcodes) {
                        if (barcode.rawValue != null) {
                          _hasScanned = true;
                          _onCodeScanned(barcode.rawValue!);
                          break;
                        }
                      }
                    },
                  ),
                  // Overlay with corners
                  CustomPaint(
                    painter: _ScannerOverlayPainter(),
                    child: Container(),
                  ),
                  // Center text
                  const Positioned(
                    bottom: 40,
                    left: 0,
                    right: 0,
                    child: Text(
                      'Placez le QR code dans le cadre',
                      textAlign: TextAlign.center,
                      style: TextStyle(
                        color: Colors.white,
                        fontSize: 16,
                        fontWeight: FontWeight.w500,
                        shadows: [Shadow(blurRadius: 10, color: Colors.black)],
                      ),
                    ),
                  ),
                ],
              ),
            ),
          ),
        ),
        // Manual entry section
        Container(
          padding: const EdgeInsets.all(24),
          decoration: const BoxDecoration(
            color: Colors.white,
            borderRadius: BorderRadius.vertical(top: Radius.circular(24)),
          ),
          child: Column(
            children: [
              const Text(
                'Autres options',
                style: TextStyle(color: Color(0xFF64748B), fontSize: 14),
              ),
              const SizedBox(height: 12),
              // Manual code entry
              TextField(
                onSubmitted: _onCodeScanned,
                decoration: InputDecoration(
                  hintText: 'Entrer le code manuellement (pay_xxx)',
                  hintStyle: const TextStyle(color: Color(0xFFCBD5E1), fontSize: 14),
                  filled: true,
                  fillColor: const Color(0xFFF8FAFC),
                  border: OutlineInputBorder(
                    borderRadius: BorderRadius.circular(12),
                    borderSide: const BorderSide(color: Color(0xFFE2E8F0)),
                  ),
                  enabledBorder: OutlineInputBorder(
                    borderRadius: BorderRadius.circular(12),
                    borderSide: const BorderSide(color: Color(0xFFE2E8F0)),
                  ),
                  prefixIcon: const Icon(Icons.keyboard, color: Color(0xFF64748B)),
                  suffixIcon: const Icon(Icons.arrow_forward, color: Color(0xFF667eea)),
                ),
              ),
              const SizedBox(height: 12),
              // Upload QR image button
              SizedBox(
                width: double.infinity,
                child: OutlinedButton.icon(
                  onPressed: _pickingImage ? null : _pickAndScanImage,
                  icon: _pickingImage 
                    ? const SizedBox(
                        width: 18, height: 18,
                        child: CircularProgressIndicator(strokeWidth: 2),
                      )
                    : const Icon(Icons.photo_library),
                  label: Text(_pickingImage ? 'Analyse en cours...' : 'Choisir une image QR'),
                  style: OutlinedButton.styleFrom(
                    foregroundColor: const Color(0xFF667eea),
                    side: const BorderSide(color: Color(0xFF667eea)),
                    padding: const EdgeInsets.symmetric(vertical: 14),
                    shape: RoundedRectangleBorder(
                      borderRadius: BorderRadius.circular(12),
                    ),
                  ),
                ),
              ),
              const SizedBox(height: 12),
              // Demo button for testing
              TextButton.icon(
                onPressed: () => _onCodeScanned('pay_demo'),
                icon: const Icon(Icons.science, size: 18),
                label: const Text('Tester avec une démo'),
                style: TextButton.styleFrom(
                  foregroundColor: const Color(0xFF64748B),
                ),
              ),
            ],
          ),
        ),
      ],
    );
  }

  /// Pick an image from gallery and scan for QR code
  Future<void> _pickAndScanImage() async {
    setState(() => _pickingImage = true);
    
    try {
      final XFile? image = await _imagePicker.pickImage(
        source: ImageSource.gallery,
        maxWidth: 1024,
        maxHeight: 1024,
      );
      
      if (image == null) {
        setState(() => _pickingImage = false);
        return;
      }

      // Use MobileScanner to analyze the image
      final BarcodeCapture? result = await _cameraController?.analyzeImage(image.path);
      
      if (result != null && result.barcodes.isNotEmpty) {
        final code = result.barcodes.first.rawValue;
        if (code != null) {
          _hasScanned = true;
          _onCodeScanned(code);
          return;
        }
      }
      
      // No QR found
      if (mounted) {
        ScaffoldMessenger.of(context).showSnackBar(
          const SnackBar(
            content: Text('Aucun QR code trouvé dans cette image'),
            backgroundColor: Colors.orange,
          ),
        );
      }
    } catch (e) {
      debugPrint('Error picking/scanning image: $e');
      if (mounted) {
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(
            content: Text('Erreur: $e'),
            backgroundColor: Colors.red,
          ),
        );
      }
    } finally {
      if (mounted) {
        setState(() => _pickingImage = false);
      }
    }
  }

  Widget _buildPaymentDetails() {
    final amount = _payment!['amount'];
    final currency = _payment!['currency'] ?? 'XOF';
    final isVariable = amount == null;

    return BlocBuilder<WalletBloc, WalletState>(
      builder: (context, state) {
        List<dynamic> wallets = [];
        if (state is WalletLoadedState) {
          wallets = state.wallets.where((w) => w.currency == currency).toList();
          if (_selectedWalletId == null && wallets.isNotEmpty) {
            _selectedWalletId = wallets.first.id;
          }
        }

        return SingleChildScrollView(
          padding: const EdgeInsets.all(24),
          child: Column(
            children: [
              // Merchant Info Card
              Container(
                padding: const EdgeInsets.all(24),
                decoration: BoxDecoration(
                  color: Colors.white,
                  borderRadius: BorderRadius.circular(20),
                  boxShadow: [
                    BoxShadow(
                      color: Colors.black.withOpacity(0.05),
                      blurRadius: 20,
                    ),
                  ],
                ),
                child: Column(
                  children: [
                    Container(
                      width: 70,
                      height: 70,
                      decoration: BoxDecoration(
                        gradient: const LinearGradient(
                          colors: [Color(0xFF667eea), Color(0xFF764ba2)],
                        ),
                        borderRadius: BorderRadius.circular(20),
                      ),
                      child: const Icon(Icons.store, color: Colors.white, size: 35),
                    ),
                    const SizedBox(height: 16),
                    Text(
                      _payment!['title'] ?? 'Paiement',
                      style: const TextStyle(
                        fontSize: 22,
                        fontWeight: FontWeight.bold,
                        color: Color(0xFF1a1a2e),
                      ),
                    ),
                    if (_payment!['description'] != null) ...[
                      const SizedBox(height: 8),
                      Text(
                        _payment!['description'],
                        style: const TextStyle(color: Color(0xFF64748B)),
                        textAlign: TextAlign.center,
                      ),
                    ],
                  ],
                ),
              ),
              const SizedBox(height: 20),

              // Amount Card
              Container(
                padding: const EdgeInsets.all(24),
                decoration: BoxDecoration(
                  color: Colors.white,
                  borderRadius: BorderRadius.circular(20),
                  boxShadow: [
                    BoxShadow(
                      color: Colors.black.withOpacity(0.05),
                      blurRadius: 20,
                    ),
                  ],
                ),
                child: Column(
                  children: [
                    if (!isVariable) ...[
                      Text(
                        '${amount.toStringAsFixed(0)} $currency',
                        style: const TextStyle(
                          color: Color(0xFF10B981),
                          fontSize: 40,
                          fontWeight: FontWeight.bold,
                        ),
                      ),
                    ] else ...[
                      const Text(
                        'Entrez le montant',
                        style: TextStyle(color: Color(0xFF64748B)),
                      ),
                      const SizedBox(height: 12),
                      TextField(
                        onChanged: (v) => setState(() => _customAmount = double.tryParse(v)),
                        keyboardType: TextInputType.number,
                        textAlign: TextAlign.center,
                        style: const TextStyle(
                          color: Color(0xFF10B981),
                          fontSize: 36,
                          fontWeight: FontWeight.bold,
                        ),
                        decoration: InputDecoration(
                          hintText: '0',
                          hintStyle: TextStyle(color: const Color(0xFF10B981).withOpacity(0.3)),
                          border: InputBorder.none,
                          suffix: Text(
                            currency,
                            style: const TextStyle(color: Color(0xFF64748B), fontSize: 20),
                          ),
                        ),
                      ),
                    ],
                  ],
                ),
              ),
              const SizedBox(height: 20),

              // Wallet Selection
              Container(
                padding: const EdgeInsets.all(20),
                decoration: BoxDecoration(
                  color: Colors.white,
                  borderRadius: BorderRadius.circular(20),
                  boxShadow: [
                    BoxShadow(
                      color: Colors.black.withOpacity(0.05),
                      blurRadius: 20,
                    ),
                  ],
                ),
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    const Text(
                      'Payer depuis',
                      style: TextStyle(
                        color: Color(0xFF64748B),
                        fontWeight: FontWeight.w500,
                      ),
                    ),
                    const SizedBox(height: 12),
                    if (wallets.isEmpty)
                      Text(
                        'Aucun portefeuille $currency disponible',
                        style: const TextStyle(color: Colors.red),
                      )
                    else
                      ...wallets.map((wallet) => _WalletOption(
                            wallet: wallet,
                            selected: wallet.id == _selectedWalletId,
                            enabled: wallet.balance >= (amount ?? _customAmount ?? 0),
                            onTap: () => setState(() => _selectedWalletId = wallet.id),
                          )),
                  ],
                ),
              ),
              const SizedBox(height: 24),

              // Pay Button
              GestureDetector(
                onTap: _canPay ? _processPayment : null,
                child: Container(
                  width: double.infinity,
                  padding: const EdgeInsets.symmetric(vertical: 18),
                  decoration: BoxDecoration(
                    gradient: _canPay
                        ? const LinearGradient(colors: [Color(0xFF10B981), Color(0xFF059669)])
                        : null,
                    color: _canPay ? null : const Color(0xFFCBD5E1),
                    borderRadius: BorderRadius.circular(16),
                    boxShadow: _canPay
                        ? [
                            BoxShadow(
                              color: const Color(0xFF10B981).withOpacity(0.3),
                              blurRadius: 12,
                              offset: const Offset(0, 4),
                            ),
                          ]
                        : null,
                  ),
                  child: _processing
                      ? const Center(
                          child: SizedBox(
                            width: 24,
                            height: 24,
                            child: CircularProgressIndicator(color: Colors.white, strokeWidth: 2),
                          ),
                        )
                      : Row(
                          mainAxisAlignment: MainAxisAlignment.center,
                          children: [
                            Text(
                              'Payer ${(amount ?? _customAmount ?? 0).toStringAsFixed(0)} $currency',
                              style: TextStyle(
                                color: _canPay ? Colors.white : const Color(0xFF94A3B8),
                                fontSize: 18,
                                fontWeight: FontWeight.bold,
                              ),
                            ),
                            const SizedBox(width: 8),
                            Icon(
                              Icons.check_circle,
                              color: _canPay ? Colors.white : const Color(0xFF94A3B8),
                            ),
                          ],
                        ),
                ),
              ),
              const SizedBox(height: 16),

              // Cancel Button
              TextButton(
                onPressed: () => setState(() {
                  _scanning = true;
                  _payment = null;
                  _hasScanned = false;
                  _cameraController?.start();
                }),
                child: const Text('Annuler et rescanner'),
              ),
            ],
          ),
        );
      },
    );
  }

  Widget _buildError() {
    return Center(
      child: Padding(
        padding: const EdgeInsets.all(24),
        child: Column(
          mainAxisAlignment: MainAxisAlignment.center,
          children: [
            Container(
              padding: const EdgeInsets.all(24),
              decoration: BoxDecoration(
                color: Colors.red.withOpacity(0.1),
                borderRadius: BorderRadius.circular(20),
              ),
              child: const Icon(Icons.error_outline, size: 64, color: Colors.red),
            ),
            const SizedBox(height: 24),
            Text(
              _error ?? 'QR Code invalide ou expiré',
              style: const TextStyle(
                color: Color(0xFF1a1a2e),
                fontSize: 18,
                fontWeight: FontWeight.w500,
              ),
              textAlign: TextAlign.center,
            ),
            const SizedBox(height: 32),
            GestureDetector(
              onTap: () => setState(() {
                _scanning = true;
                _error = null;
                _hasScanned = false;
                _cameraController?.start();
              }),
              child: Container(
                padding: const EdgeInsets.symmetric(horizontal: 32, vertical: 16),
                decoration: BoxDecoration(
                  gradient: const LinearGradient(
                    colors: [Color(0xFF667eea), Color(0xFF764ba2)],
                  ),
                  borderRadius: BorderRadius.circular(12),
                ),
                child: const Text(
                  'Réessayer',
                  style: TextStyle(
                    color: Colors.white,
                    fontWeight: FontWeight.bold,
                    fontSize: 16,
                  ),
                ),
              ),
            ),
          ],
        ),
      ),
    );
  }

  bool get _canPay {
    if (_selectedWalletId == null) return false;
    final amount = _payment?['amount'] ?? _customAmount;
    if (amount == null || amount <= 0) return false;
    
    final walletState = context.read<WalletBloc>().state;
    if (walletState is WalletLoadedState) {
      final wallet = walletState.wallets.where((w) => w.id == _selectedWalletId).firstOrNull;
      if (wallet == null || wallet.balance < amount) return false;
    }
    return true;
  }

  void _onCodeScanned(String code) async {
    _cameraController?.stop();
    
    // Parse QR code data
    String paymentId = code;
    
    try {
      // Try to parse as JSON
      final data = jsonDecode(code);
      paymentId = data['payment_id'] ?? data['id'] ?? code;
    } catch (_) {
      // Not JSON, try URL parsing
      if (code.contains('/pay/')) {
        paymentId = code.split('/pay/').last;
      }
    }

    // Simulate API call - in real app, fetch from backend
    await Future.delayed(const Duration(milliseconds: 500));

    if (paymentId == 'pay_demo' || paymentId.startsWith('pay_')) {
      setState(() {
        _payment = {
          'payment_id': paymentId,
          'title': 'Boutique CryptoBank',
          'description': 'Achat en magasin',
          'amount': 5000.0,
          'currency': 'XOF',
          'merchant': 'CryptoBank Store',
        };
        _scanning = false;
      });
    } else {
      setState(() {
        _error = 'QR Code non reconnu: $paymentId';
        _scanning = false;
      });
    }
  }

  void _processPayment() async {
    if (!_canPay) return;

    // Require security confirmation before payment
    final confirmed = await SecurityConfirmation.require(
      context,
      title: 'Confirmer le paiement',
      message: 'Authentifiez-vous pour valider ce paiement',
    );
    if (!confirmed) return;

    setState(() => _processing = true);

    try {
      // Simulate payment - in real app, call API
      await Future.delayed(const Duration(seconds: 2));

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
      if (mounted) {
        setState(() => _processing = false);
      }
    }
  }

  void _showSuccessDialog() {
    showDialog(
      context: context,
      barrierDismissible: false,
      builder: (context) => AlertDialog(
        backgroundColor: Colors.white,
        shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(24)),
        content: Column(
          mainAxisSize: MainAxisSize.min,
          children: [
            Container(
              padding: const EdgeInsets.all(16),
              decoration: BoxDecoration(
                color: const Color(0xFF10B981).withOpacity(0.1),
                shape: BoxShape.circle,
              ),
              child: const Icon(Icons.check_circle, color: Color(0xFF10B981), size: 64),
            ),
            const SizedBox(height: 24),
            const Text(
              'Paiement réussi! 🎉',
              style: TextStyle(
                fontSize: 22,
                fontWeight: FontWeight.bold,
                color: Color(0xFF1a1a2e),
              ),
            ),
            const SizedBox(height: 8),
            Text(
              '${(_payment!['amount'] ?? _customAmount).toStringAsFixed(0)} ${_payment!['currency']}',
              style: const TextStyle(
                color: Color(0xFF10B981),
                fontSize: 32,
                fontWeight: FontWeight.bold,
              ),
            ),
            const SizedBox(height: 24),
            GestureDetector(
              onTap: () {
                Navigator.of(context).pop();
                context.go('/wallet');
              },
              child: Container(
                width: double.infinity,
                padding: const EdgeInsets.symmetric(vertical: 16),
                decoration: BoxDecoration(
                  gradient: const LinearGradient(
                    colors: [Color(0xFF667eea), Color(0xFF764ba2)],
                  ),
                  borderRadius: BorderRadius.circular(12),
                ),
                child: const Text(
                  'Terminé',
                  textAlign: TextAlign.center,
                  style: TextStyle(
                    color: Colors.white,
                    fontWeight: FontWeight.bold,
                    fontSize: 16,
                  ),
                ),
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
    return GestureDetector(
      onTap: enabled ? onTap : null,
      child: Container(
        padding: const EdgeInsets.all(16),
        margin: const EdgeInsets.only(bottom: 8),
        decoration: BoxDecoration(
          color: selected ? const Color(0xFF667eea).withOpacity(0.1) : const Color(0xFFF8FAFC),
          borderRadius: BorderRadius.circular(12),
          border: Border.all(
            color: selected ? const Color(0xFF667eea) : const Color(0xFFE2E8F0),
            width: 2,
          ),
        ),
        child: Row(
          children: [
            Container(
              width: 40,
              height: 40,
              decoration: BoxDecoration(
                color: selected ? const Color(0xFF667eea) : const Color(0xFFE2E8F0),
                borderRadius: BorderRadius.circular(10),
              ),
              child: Icon(
                Icons.account_balance_wallet,
                color: selected ? Colors.white : const Color(0xFF64748B),
                size: 20,
              ),
            ),
            const SizedBox(width: 12),
            Expanded(
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  Text(
                    wallet.name ?? wallet.currency,
                    style: TextStyle(
                      color: enabled ? const Color(0xFF1a1a2e) : const Color(0xFF94A3B8),
                      fontWeight: FontWeight.w600,
                    ),
                  ),
                  Text(
                    '${wallet.balance.toStringAsFixed(0)} ${wallet.currency}',
                    style: TextStyle(
                      color: enabled ? const Color(0xFF64748B) : const Color(0xFFCBD5E1),
                      fontSize: 13,
                    ),
                  ),
                ],
              ),
            ),
            if (!enabled)
              Container(
                padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 4),
                decoration: BoxDecoration(
                  color: Colors.red.withOpacity(0.1),
                  borderRadius: BorderRadius.circular(8),
                ),
                child: const Text(
                  'Insuffisant',
                  style: TextStyle(color: Colors.red, fontSize: 11),
                ),
              )
            else if (selected)
              const Icon(Icons.check_circle, color: Color(0xFF667eea)),
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
      ..color = const Color(0xFF667eea)
      ..strokeWidth = 4
      ..style = PaintingStyle.stroke
      ..strokeCap = StrokeCap.round;

    const cornerLength = 50.0;
    final rect = Rect.fromCenter(
      center: Offset(size.width / 2, size.height / 2),
      width: 220,
      height: 220,
    );

    // Top left corner
    canvas.drawLine(rect.topLeft, rect.topLeft + const Offset(cornerLength, 0), paint);
    canvas.drawLine(rect.topLeft, rect.topLeft + const Offset(0, cornerLength), paint);

    // Top right corner
    canvas.drawLine(rect.topRight, rect.topRight + const Offset(-cornerLength, 0), paint);
    canvas.drawLine(rect.topRight, rect.topRight + const Offset(0, cornerLength), paint);

    // Bottom left corner
    canvas.drawLine(rect.bottomLeft, rect.bottomLeft + const Offset(cornerLength, 0), paint);
    canvas.drawLine(rect.bottomLeft, rect.bottomLeft + const Offset(0, -cornerLength), paint);

    // Bottom right corner
    canvas.drawLine(rect.bottomRight, rect.bottomRight + const Offset(-cornerLength, 0), paint);
    canvas.drawLine(rect.bottomRight, rect.bottomRight + const Offset(0, -cornerLength), paint);
  }

  @override
  bool shouldRepaint(covariant CustomPainter oldDelegate) => false;
}
