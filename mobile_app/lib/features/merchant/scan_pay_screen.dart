import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';
import 'package:google_fonts/google_fonts.dart';
import 'package:image_picker/image_picker.dart';
import 'package:mobile_scanner/mobile_scanner.dart';
import 'dart:io';

import '../../core/theme/app_theme.dart';
import '../../core/widgets/glass_container.dart';
import '../../core/services/api_service.dart';

/// Scan & Pay Screen matching web design exactly - with real QR scanning
class ScanPayScreen extends StatefulWidget {
  const ScanPayScreen({super.key});

  @override
  State<ScanPayScreen> createState() => _ScanPayScreenState();
}

class _ScanPayScreenState extends State<ScanPayScreen> {
  final ApiService _api = ApiService();
  final ImagePicker _picker = ImagePicker();
  
  String _activeTab = 'camera';
  final _codeController = TextEditingController();
  bool _loading = false;
  String? _error;
  bool _cameraActive = true;
  MobileScannerController? _scannerController;
  bool _hasScanned = false;

  @override
  void initState() {
    super.initState();
    _initScanner();
  }

  void _initScanner() {
    _scannerController = MobileScannerController(
      detectionSpeed: DetectionSpeed.normal,
      facing: CameraFacing.back,
      torchEnabled: false,
    );
  }

  @override
  void dispose() {
    _codeController.dispose();
    _scannerController?.dispose();
    super.dispose();
  }

  void _onDetect(BarcodeCapture capture) {
    if (_hasScanned) return;
    
    final List<Barcode> barcodes = capture.barcodes;
    for (final barcode in barcodes) {
      final code = barcode.rawValue;
      if (code != null && code.isNotEmpty) {
        setState(() => _hasScanned = true);
        _processPaymentCode(code);
        break;
      }
    }
  }

  Future<void> _processPaymentCode(String code) async {
    setState(() {
      _loading = true;
      _error = null;
    });

    try {
      // Extract payment ID from URL or code
      String paymentId = code;
      if (code.contains('/pay/')) {
        paymentId = code.split('/pay/').last.split('?').first;
      } else if (code.startsWith('pay_')) {
        paymentId = code;
      }

      // Fetch payment details
      final paymentDetails = await _api.merchant.getPaymentDetails(paymentId);
      
      if (mounted) {
        setState(() => _loading = false);
        _showPaymentConfirmation(paymentDetails, paymentId);
      }
    } catch (e) {
      if (mounted) {
        setState(() {
          _loading = false;
          _error = 'Code de paiement invalide ou expiré';
          _hasScanned = false;
        });
      }
    }
  }

  void _showPaymentConfirmation(Map<String, dynamic> payment, String paymentId) {
    final isDark = Theme.of(context).brightness == Brightness.dark;
    
    showModalBottomSheet(
      context: context,
      isScrollControlled: true,
      backgroundColor: Colors.transparent,
      builder: (context) => _PaymentConfirmationSheet(
        payment: payment,
        paymentId: paymentId,
        api: _api,
        onComplete: () {
          Navigator.pop(context);
          context.go('/dashboard');
        },
        onCancel: () {
          Navigator.pop(context);
          setState(() => _hasScanned = false);
        },
      ),
    );
  }

  @override
  Widget build(BuildContext context) {
    final isDark = Theme.of(context).brightness == Brightness.dark;
    
    return Scaffold(
      backgroundColor: Colors.transparent,
      body: Container(
        decoration: BoxDecoration(
          gradient: LinearGradient(
            begin: Alignment.topLeft,
            end: Alignment.bottomRight,
            colors: isDark 
                ? [const Color(0xFF020617), const Color(0xFF0F172A)]
                : [const Color(0xFFFAFBFC), const Color(0xFFEFF6FF)],
          ),
        ),
        child: SafeArea(
          child: Column(
            children: [
              _buildAppBar(isDark),
              Expanded(
                child: SingleChildScrollView(
                  padding: const EdgeInsets.all(20),
                  child: Column(
                    children: [
                      _buildHeader(isDark),
                      const SizedBox(height: 24),
                      _buildMethodTabs(isDark),
                      const SizedBox(height: 20),
                      _buildTabContent(isDark),
                      
                      if (_error != null) ...[
                        const SizedBox(height: 16),
                        Container(
                          padding: const EdgeInsets.all(16),
                          decoration: BoxDecoration(
                            color: const Color(0xFFEF4444).withOpacity(0.1),
                            border: Border.all(color: const Color(0xFFEF4444).withOpacity(0.2)),
                            borderRadius: BorderRadius.circular(16),
                          ),
                          child: Row(
                            children: [
                              const Icon(Icons.error_outline, color: Color(0xFFEF4444)),
                              const SizedBox(width: 12),
                              Expanded(
                                child: Text(
                                  _error!,
                                  style: GoogleFonts.inter(
                                    color: const Color(0xFFEF4444),
                                    fontSize: 14,
                                  ),
                                ),
                              ),
                            ],
                          ),
                        ),
                      ],
                      
                      const SizedBox(height: 24),
                      _buildHelpText(isDark),
                    ],
                  ),
                ),
              ),
            ],
          ),
        ),
      ),
    );
  }

  Widget _buildAppBar(bool isDark) {
    return Padding(
      padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 8),
      child: Row(
        children: [
          GlassContainer(
            padding: EdgeInsets.zero,
            width: 40,
            height: 40,
            borderRadius: 12,
            child: IconButton(
              icon: Icon(Icons.arrow_back_ios_new, size: 20, 
                  color: isDark ? Colors.white : AppTheme.textPrimaryColor),
              onPressed: () => context.go('/dashboard'),
            ),
          ),
          const Spacer(),
          GlassContainer(
            padding: EdgeInsets.zero,
            width: 40,
            height: 40,
            borderRadius: 12,
            child: IconButton(
              icon: Icon(Icons.home_rounded, size: 20, 
                  color: isDark ? Colors.white : AppTheme.textPrimaryColor),
              onPressed: () => context.go('/dashboard'),
            ),
          ),
        ],
      ),
    );
  }

  Widget _buildHeader(bool isDark) {
    return Column(
      children: [
        Container(
          width: 80,
          height: 80,
          decoration: BoxDecoration(
            gradient: const LinearGradient(
              colors: [Color(0xFF6366F1), Color(0xFF8B5CF6)],
            ),
            borderRadius: BorderRadius.circular(20),
            boxShadow: [
              BoxShadow(
                color: const Color(0xFF6366F1).withOpacity(0.3),
                blurRadius: 20,
                offset: const Offset(0, 8),
              ),
            ],
          ),
          child: const Center(
            child: Text('📱', style: TextStyle(fontSize: 36)),
          ),
        ),
        const SizedBox(height: 16),
        Text(
          'Payer un marchand',
          style: GoogleFonts.inter(
            fontSize: 24,
            fontWeight: FontWeight.bold,
            color: isDark ? Colors.white : const Color(0xFF1E293B),
          ),
        ),
        const SizedBox(height: 8),
        Text(
          'Scannez ou entrez le code de paiement',
          style: GoogleFonts.inter(
            fontSize: 14,
            color: isDark ? const Color(0xFF94A3B8) : const Color(0xFF64748B),
          ),
        ),
      ],
    );
  }

  Widget _buildMethodTabs(bool isDark) {
    return Row(
      children: [
        Expanded(child: _buildMethodTab('📷', 'Scanner', 'camera', isDark)),
        const SizedBox(width: 8),
        Expanded(child: _buildMethodTab('⌨️', 'Code', 'code', isDark)),
        const SizedBox(width: 8),
        Expanded(child: _buildMethodTab('🖼️', 'Image', 'image', isDark)),
      ],
    );
  }

  Widget _buildMethodTab(String emoji, String label, String value, bool isDark) {
    final isActive = _activeTab == value;
    return GestureDetector(
      onTap: () {
        setState(() {
          _activeTab = value;
          _error = null;
          _hasScanned = false;
        });
      },
      child: AnimatedContainer(
        duration: const Duration(milliseconds: 200),
        padding: const EdgeInsets.symmetric(vertical: 12, horizontal: 16),
        decoration: BoxDecoration(
          color: isActive 
              ? const Color(0xFF6366F1)
              : (isDark ? const Color(0xFF1E293B) : const Color(0xFFF1F5F9)),
          borderRadius: BorderRadius.circular(12),
          boxShadow: isActive ? [
            BoxShadow(
              color: const Color(0xFF6366F1).withOpacity(0.3),
              blurRadius: 12,
              offset: const Offset(0, 4),
            ),
          ] : null,
        ),
        child: Row(
          mainAxisAlignment: MainAxisAlignment.center,
          children: [
            Text(emoji, style: const TextStyle(fontSize: 16)),
            const SizedBox(width: 6),
            Text(
              label,
              style: GoogleFonts.inter(
                fontSize: 14,
                fontWeight: FontWeight.w600,
                color: isActive 
                    ? Colors.white
                    : (isDark ? const Color(0xFF94A3B8) : const Color(0xFF64748B)),
              ),
            ),
          ],
        ),
      ),
    );
  }

  Widget _buildTabContent(bool isDark) {
    switch (_activeTab) {
      case 'camera':
        return _buildCameraTab(isDark);
      case 'code':
        return _buildCodeTab(isDark);
      case 'image':
        return _buildImageTab(isDark);
      default:
        return const SizedBox();
    }
  }

  Widget _buildCameraTab(bool isDark) {
    return GlassContainer(
      padding: const EdgeInsets.all(20),
      borderRadius: 20,
      child: Column(
        children: [
          // Camera view with scanner
          ClipRRect(
            borderRadius: BorderRadius.circular(16),
            child: SizedBox(
              height: 300,
              child: Stack(
                children: [
                  // Real QR Scanner
                  MobileScanner(
                    controller: _scannerController,
                    onDetect: _onDetect,
                  ),
                  
                  // Scanner frame overlay
                  Center(
                    child: SizedBox(
                      width: 200,
                      height: 200,
                      child: CustomPaint(
                        painter: ScannerFramePainter(),
                      ),
                    ),
                  ),
                  
                  // Loading overlay
                  if (_loading)
                    Container(
                      color: Colors.black.withOpacity(0.7),
                      child: const Center(
                        child: Column(
                          mainAxisAlignment: MainAxisAlignment.center,
                          children: [
                            CircularProgressIndicator(color: Color(0xFF6366F1)),
                            SizedBox(height: 16),
                            Text(
                              'Chargement...',
                              style: TextStyle(color: Colors.white70),
                            ),
                          ],
                        ),
                      ),
                    ),
                ],
              ),
            ),
          ),
          
          const SizedBox(height: 16),
          
          // Camera controls
          Row(
            mainAxisAlignment: MainAxisAlignment.center,
            children: [
              _buildCameraControl(
                Icons.flash_on, 
                'Flash', 
                isDark,
                onTap: () => _scannerController?.toggleTorch(),
              ),
              const SizedBox(width: 20),
              _buildCameraControl(
                Icons.cameraswitch, 
                'Changer', 
                isDark,
                onTap: () => _scannerController?.switchCamera(),
              ),
            ],
          ),
          
          const SizedBox(height: 16),
          Text(
            'Placez le QR code dans le cadre pour scanner automatiquement',
            style: GoogleFonts.inter(
              fontSize: 12,
              color: isDark ? const Color(0xFF64748B) : const Color(0xFF94A3B8),
            ),
            textAlign: TextAlign.center,
          ),
        ],
      ),
    );
  }

  Widget _buildCameraControl(IconData icon, String label, bool isDark, {VoidCallback? onTap}) {
    return GestureDetector(
      onTap: onTap,
      child: Container(
        padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 12),
        decoration: BoxDecoration(
          color: isDark ? const Color(0xFF1E293B) : const Color(0xFFF1F5F9),
          borderRadius: BorderRadius.circular(12),
        ),
        child: Row(
          children: [
            Icon(icon, size: 20, color: isDark ? Colors.white70 : const Color(0xFF64748B)),
            const SizedBox(width: 8),
            Text(
              label,
              style: GoogleFonts.inter(
                fontSize: 14,
                color: isDark ? Colors.white70 : const Color(0xFF64748B),
              ),
            ),
          ],
        ),
      ),
    );
  }

  Widget _buildCodeTab(bool isDark) {
    return GlassContainer(
      padding: const EdgeInsets.all(20),
      borderRadius: 20,
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Text(
            'Code de paiement',
            style: GoogleFonts.inter(
              fontSize: 14,
              fontWeight: FontWeight.w500,
              color: isDark ? const Color(0xFF94A3B8) : const Color(0xFF64748B),
            ),
          ),
          const SizedBox(height: 8),
          TextField(
            controller: _codeController,
            style: GoogleFonts.sourceCodePro(
              fontSize: 16,
              color: isDark ? Colors.white : const Color(0xFF1E293B),
            ),
            textAlign: TextAlign.center,
            decoration: InputDecoration(
              hintText: 'pay_abc123...',
              hintStyle: GoogleFonts.sourceCodePro(
                color: isDark ? const Color(0xFF475569) : const Color(0xFFCBD5E1),
              ),
              filled: true,
              fillColor: isDark ? const Color(0xFF1E293B) : const Color(0xFFF8FAFC),
              border: OutlineInputBorder(
                borderRadius: BorderRadius.circular(12),
                borderSide: BorderSide(
                  color: isDark ? const Color(0xFF334155) : const Color(0xFFE2E8F0),
                ),
              ),
              enabledBorder: OutlineInputBorder(
                borderRadius: BorderRadius.circular(12),
                borderSide: BorderSide(
                  color: isDark ? const Color(0xFF334155) : const Color(0xFFE2E8F0),
                ),
              ),
              focusedBorder: OutlineInputBorder(
                borderRadius: BorderRadius.circular(12),
                borderSide: const BorderSide(color: Color(0xFF6366F1)),
              ),
            ),
          ),
          const SizedBox(height: 12),
          Text(
            'Le code se trouve sous le QR code du marchand',
            style: GoogleFonts.inter(
              fontSize: 12,
              color: isDark ? const Color(0xFF64748B) : const Color(0xFF94A3B8),
            ),
            textAlign: TextAlign.center,
          ),
          const SizedBox(height: 20),
          SizedBox(
            width: double.infinity,
            child: ElevatedButton(
              onPressed: _loading ? null : () => _processPaymentCode(_codeController.text),
              style: ElevatedButton.styleFrom(
                backgroundColor: const Color(0xFF6366F1),
                foregroundColor: Colors.white,
                padding: const EdgeInsets.symmetric(vertical: 16),
                shape: RoundedRectangleBorder(
                  borderRadius: BorderRadius.circular(12),
                ),
                elevation: 0,
              ),
              child: _loading 
                  ? const SizedBox(
                      width: 20,
                      height: 20,
                      child: CircularProgressIndicator(
                        strokeWidth: 2,
                        color: Colors.white,
                      ),
                    )
                  : Text(
                      'Continuer →',
                      style: GoogleFonts.inter(
                        fontSize: 16,
                        fontWeight: FontWeight.w600,
                      ),
                    ),
            ),
          ),
        ],
      ),
    );
  }

  Widget _buildImageTab(bool isDark) {
    return GlassContainer(
      padding: const EdgeInsets.all(20),
      borderRadius: 20,
      child: Column(
        children: [
          GestureDetector(
            onTap: _pickImage,
            child: Container(
              padding: const EdgeInsets.all(40),
              decoration: BoxDecoration(
                border: Border.all(
                  color: isDark ? const Color(0xFF334155) : const Color(0xFFE2E8F0),
                  width: 2,
                  style: BorderStyle.solid,
                ),
                borderRadius: BorderRadius.circular(16),
              ),
              child: Column(
                children: [
                  const Text('🖼️', style: TextStyle(fontSize: 48)),
                  const SizedBox(height: 16),
                  Text(
                    'Cliquez pour sélectionner une image',
                    style: GoogleFonts.inter(
                      fontSize: 16,
                      fontWeight: FontWeight.w500,
                      color: isDark ? Colors.white : const Color(0xFF1E293B),
                    ),
                  ),
                  const SizedBox(height: 8),
                  Text(
                    'PNG, JPG jusqu\'à 5MB',
                    style: GoogleFonts.inter(
                      fontSize: 12,
                      color: isDark ? const Color(0xFF64748B) : const Color(0xFF94A3B8),
                    ),
                  ),
                ],
              ),
            ),
          ),
        ],
      ),
    );
  }

  Future<void> _pickImage() async {
    final XFile? image = await _picker.pickImage(source: ImageSource.gallery);
    if (image != null) {
      setState(() => _loading = true);
      
      try {
        // Analyze image for QR code
        final result = await _scannerController?.analyzeImage(image.path);
        if (result != null && result.barcodes.isNotEmpty) {
          final code = result.barcodes.first.rawValue;
          if (code != null) {
            await _processPaymentCode(code);
            return;
          }
        }
        setState(() {
          _loading = false;
          _error = 'Aucun QR code trouvé dans l\'image';
        });
      } catch (e) {
        setState(() {
          _loading = false;
          _error = 'Erreur lors de l\'analyse de l\'image';
        });
      }
    }
  }

  Widget _buildHelpText(bool isDark) {
    return GestureDetector(
      onTap: () => context.go('/more/merchant'),
      child: RichText(
        textAlign: TextAlign.center,
        text: TextSpan(
          style: GoogleFonts.inter(
            fontSize: 14,
            color: isDark ? const Color(0xFF64748B) : const Color(0xFF94A3B8),
          ),
          children: [
            const TextSpan(text: 'Vous n\'avez pas de code ? '),
            TextSpan(
              text: 'Créer un paiement marchand',
              style: GoogleFonts.inter(
                color: const Color(0xFF6366F1),
                fontWeight: FontWeight.w500,
              ),
            ),
          ],
        ),
      ),
    );
  }
}

/// Scanner frame painter
class ScannerFramePainter extends CustomPainter {
  @override
  void paint(Canvas canvas, Size size) {
    final paint = Paint()
      ..color = const Color(0xFF6366F1).withOpacity(0.9)
      ..strokeWidth = 4
      ..style = PaintingStyle.stroke;
    
    const cornerLength = 30.0;
    const cornerRadius = 12.0;
    
    // Top left corner
    final topLeftPath = Path()
      ..moveTo(0, cornerLength)
      ..lineTo(0, cornerRadius)
      ..quadraticBezierTo(0, 0, cornerRadius, 0)
      ..lineTo(cornerLength, 0);
    canvas.drawPath(topLeftPath, paint);
    
    // Top right corner
    final topRightPath = Path()
      ..moveTo(size.width - cornerLength, 0)
      ..lineTo(size.width - cornerRadius, 0)
      ..quadraticBezierTo(size.width, 0, size.width, cornerRadius)
      ..lineTo(size.width, cornerLength);
    canvas.drawPath(topRightPath, paint);
    
    // Bottom left corner
    final bottomLeftPath = Path()
      ..moveTo(0, size.height - cornerLength)
      ..lineTo(0, size.height - cornerRadius)
      ..quadraticBezierTo(0, size.height, cornerRadius, size.height)
      ..lineTo(cornerLength, size.height);
    canvas.drawPath(bottomLeftPath, paint);
    
    // Bottom right corner
    final bottomRightPath = Path()
      ..moveTo(size.width - cornerLength, size.height)
      ..lineTo(size.width - cornerRadius, size.height)
      ..quadraticBezierTo(size.width, size.height, size.width, size.height - cornerRadius)
      ..lineTo(size.width, size.height - cornerLength);
    canvas.drawPath(bottomRightPath, paint);
  }

  @override
  bool shouldRepaint(covariant CustomPainter oldDelegate) => false;
}

/// Payment confirmation bottom sheet
class _PaymentConfirmationSheet extends StatefulWidget {
  final Map<String, dynamic> payment;
  final String paymentId;
  final ApiService api;
  final VoidCallback onComplete;
  final VoidCallback onCancel;

  const _PaymentConfirmationSheet({
    required this.payment,
    required this.paymentId,
    required this.api,
    required this.onComplete,
    required this.onCancel,
  });

  @override
  State<_PaymentConfirmationSheet> createState() => _PaymentConfirmationSheetState();
}

class _PaymentConfirmationSheetState extends State<_PaymentConfirmationSheet> {
  List<Map<String, dynamic>> _wallets = [];
  String? _selectedWalletId;
  final _amountController = TextEditingController();
  bool _loading = false;
  bool _loadingWallets = true;
  String? _error;

  @override
  void initState() {
    super.initState();
    _loadWallets();
    // Pre-fill amount if fixed
    if (widget.payment['amount'] != null) {
      _amountController.text = widget.payment['amount'].toString();
    }
  }

  Future<void> _loadWallets() async {
    try {
      final res = await widget.api.wallet.getWallets();
      setState(() {
        _wallets = List<Map<String, dynamic>>.from(
          res['wallets'] ?? res['data']?['wallets'] ?? res['data'] ?? []
        );
        _loadingWallets = false;
      });
    } catch (e) {
      setState(() => _loadingWallets = false);
    }
  }

  Future<void> _confirmPayment() async {
    if (_selectedWalletId == null) {
      setState(() => _error = 'Veuillez sélectionner un portefeuille');
      return;
    }

    final amount = double.tryParse(_amountController.text);
    if (amount == null || amount <= 0) {
      setState(() => _error = 'Veuillez entrer un montant valide');
      return;
    }

    setState(() {
      _loading = true;
      _error = null;
    });

    try {
      await widget.api.merchant.payPayment(
        widget.paymentId,
        _selectedWalletId!,
        amount,
      );
      
      if (mounted) {
        ScaffoldMessenger.of(context).showSnackBar(
          const SnackBar(
            content: Text('Paiement effectué avec succès!'),
            backgroundColor: Colors.green,
          ),
        );
        widget.onComplete();
      }
    } catch (e) {
      setState(() {
        _loading = false;
        _error = e.toString();
      });
    }
  }

  @override
  Widget build(BuildContext context) {
    final isDark = Theme.of(context).brightness == Brightness.dark;
    
    return Container(
      margin: const EdgeInsets.all(16),
      decoration: BoxDecoration(
        color: isDark ? const Color(0xFF1E293B) : Colors.white,
        borderRadius: BorderRadius.circular(24),
      ),
      child: Padding(
        padding: const EdgeInsets.all(24),
        child: Column(
          mainAxisSize: MainAxisSize.min,
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            // Header
            Row(
              children: [
                const Text('💳', style: TextStyle(fontSize: 32)),
                const SizedBox(width: 12),
                Expanded(
                  child: Column(
                    crossAxisAlignment: CrossAxisAlignment.start,
                    children: [
                      Text(
                        widget.payment['title'] ?? 'Paiement',
                        style: GoogleFonts.inter(
                          fontSize: 20,
                          fontWeight: FontWeight.bold,
                          color: isDark ? Colors.white : const Color(0xFF1E293B),
                        ),
                      ),
                      if (widget.payment['merchant_name'] != null)
                        Text(
                          widget.payment['merchant_name'],
                          style: GoogleFonts.inter(
                            fontSize: 14,
                            color: const Color(0xFF64748B),
                          ),
                        ),
                    ],
                  ),
                ),
                IconButton(
                  onPressed: widget.onCancel,
                  icon: const Icon(Icons.close),
                ),
              ],
            ),
            
            const SizedBox(height: 24),
            
            // Amount
            Text(
              'Montant à payer',
              style: GoogleFonts.inter(
                fontSize: 14,
                fontWeight: FontWeight.w500,
                color: const Color(0xFF64748B),
              ),
            ),
            const SizedBox(height: 8),
            TextField(
              controller: _amountController,
              keyboardType: TextInputType.number,
              enabled: widget.payment['type'] != 'fixed',
              style: GoogleFonts.inter(
                fontSize: 24,
                fontWeight: FontWeight.bold,
                color: isDark ? Colors.white : const Color(0xFF1E293B),
              ),
              textAlign: TextAlign.center,
              decoration: InputDecoration(
                suffixText: widget.payment['currency'] ?? 'EUR',
                filled: true,
                fillColor: isDark ? const Color(0xFF0F172A) : const Color(0xFFF8FAFC),
                border: OutlineInputBorder(
                  borderRadius: BorderRadius.circular(12),
                  borderSide: BorderSide.none,
                ),
              ),
            ),
            
            const SizedBox(height: 20),
            
            // Wallet selection
            Text(
              'Payer depuis',
              style: GoogleFonts.inter(
                fontSize: 14,
                fontWeight: FontWeight.w500,
                color: const Color(0xFF64748B),
              ),
            ),
            const SizedBox(height: 8),
            if (_loadingWallets)
              const Center(child: CircularProgressIndicator())
            else
              Container(
                padding: const EdgeInsets.symmetric(horizontal: 16),
                decoration: BoxDecoration(
                  color: isDark ? const Color(0xFF0F172A) : const Color(0xFFF8FAFC),
                  borderRadius: BorderRadius.circular(12),
                ),
                child: DropdownButtonHideUnderline(
                  child: DropdownButton<String>(
                    value: _selectedWalletId,
                    isExpanded: true,
                    hint: Text('Sélectionner un portefeuille', style: GoogleFonts.inter(
                      color: const Color(0xFF94A3B8),
                    )),
                    dropdownColor: isDark ? const Color(0xFF1E293B) : Colors.white,
                    items: _wallets.map((w) => DropdownMenuItem(
                      value: w['id'].toString(),
                      child: Text('${w['currency']} - ${(w['balance'] as num).toStringAsFixed(2)} ${w['currency']}'),
                    )).toList(),
                    onChanged: (v) => setState(() => _selectedWalletId = v),
                    style: GoogleFonts.inter(
                      color: isDark ? Colors.white : const Color(0xFF1E293B),
                    ),
                  ),
                ),
              ),
            
            // Error
            if (_error != null) ...[
              const SizedBox(height: 16),
              Text(
                _error!,
                style: GoogleFonts.inter(color: const Color(0xFFEF4444)),
              ),
            ],
            
            const SizedBox(height: 24),
            
            // Confirm button
            SizedBox(
              width: double.infinity,
              child: ElevatedButton(
                onPressed: _loading ? null : _confirmPayment,
                style: ElevatedButton.styleFrom(
                  backgroundColor: const Color(0xFF22C55E),
                  foregroundColor: Colors.white,
                  padding: const EdgeInsets.symmetric(vertical: 16),
                  shape: RoundedRectangleBorder(
                    borderRadius: BorderRadius.circular(12),
                  ),
                ),
                child: _loading
                    ? const SizedBox(
                        width: 20,
                        height: 20,
                        child: CircularProgressIndicator(color: Colors.white, strokeWidth: 2),
                      )
                    : Text(
                        '✓ Confirmer le paiement',
                        style: GoogleFonts.inter(
                          fontSize: 16,
                          fontWeight: FontWeight.w600,
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
