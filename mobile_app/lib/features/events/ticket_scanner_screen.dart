import 'package:flutter/material.dart';
import 'package:mobile_scanner/mobile_scanner.dart';
import '../../core/services/ticket_api_service.dart';

class TicketScannerScreen extends StatefulWidget {
  final String eventId;
  final String eventTitle;

  const TicketScannerScreen({
    super.key,
    required this.eventId,
    required this.eventTitle,
  });

  @override
  State<TicketScannerScreen> createState() => _TicketScannerScreenState();
}

class _TicketScannerScreenState extends State<TicketScannerScreen> {
  final TicketApiService _ticketApi = TicketApiService();
  final MobileScannerController _scannerController = MobileScannerController();
  final TextEditingController _manualCodeController = TextEditingController();
  
  bool _isScanning = true;
  bool _isProcessing = false;
  Map<String, dynamic>? _scanResult;
  String? _errorMessage;

  @override
  void dispose() {
    _scannerController.dispose();
    _manualCodeController.dispose();
    super.dispose();
  }

  Future<void> _onScan(BarcodeCapture capture) async {
    if (_isProcessing || !_isScanning) return;
    
    final barcode = capture.barcodes.firstOrNull;
    if (barcode?.rawValue == null) return;

    String code = barcode!.rawValue!.trim();
    
    // Handle ZEKORA_TICKET: prefix
    if (code.startsWith('ZEKORA_TICKET:')) {
      code = code.replaceFirst('ZEKORA_TICKET:', '');
    }

    await _verifyTicket(code);
  }

  Future<void> _verifyTicket(String code) async {
    setState(() {
      _isProcessing = true;
      _isScanning = false;
      _errorMessage = null;
    });

    try {
      final result = await _ticketApi.verifyTicket(code);
      setState(() {
        _scanResult = result;
      });
    } catch (e) {
      setState(() {
        _errorMessage = e.toString().replaceFirst('Exception: ', '');
      });
    } finally {
      setState(() => _isProcessing = false);
    }
  }

  Future<void> _markAsUsed() async {
    if (_scanResult == null || _scanResult!['ticket'] == null) return;
    
    final ticketId = _scanResult!['ticket']['id'];
    
    setState(() => _isProcessing = true);
    
    try {
      await _ticketApi.useTicket(ticketId);
      setState(() {
        _scanResult!['ticket']['used'] = true;
      });
      if (mounted) {
        ScaffoldMessenger.of(context).showSnackBar(
          const SnackBar(
            content: Text('✅ Ticket validé avec succès!'),
            backgroundColor: Colors.green,
          ),
        );
      }
    } catch (e) {
      if (mounted) {
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(
            content: Text('Erreur: ${e.toString()}'),
            backgroundColor: Colors.red,
          ),
        );
      }
    } finally {
      setState(() => _isProcessing = false);
    }
  }

  void _resetScanner() {
    setState(() {
      _scanResult = null;
      _errorMessage = null;
      _isScanning = true;
      _manualCodeController.clear();
    });
  }

  void _verifyManualCode() {
    final code = _manualCodeController.text.trim();
    if (code.isEmpty) return;
    _verifyTicket(code);
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: const Color(0xFF1a1a2e),
      appBar: AppBar(
        backgroundColor: Colors.transparent,
        elevation: 0,
        leading: IconButton(
          icon: const Icon(Icons.arrow_back, color: Colors.white),
          onPressed: () => Navigator.pop(context),
        ),
        title: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            const Text(
              'Scanner les tickets',
              style: TextStyle(color: Colors.white, fontSize: 18),
            ),
            Text(
              widget.eventTitle,
              style: const TextStyle(color: Colors.white54, fontSize: 12),
            ),
          ],
        ),
        actions: [
          IconButton(
            icon: Icon(
              _scannerController.torchEnabled ? Icons.flash_on : Icons.flash_off,
              color: Colors.white,
            ),
            onPressed: () => _scannerController.toggleTorch(),
          ),
        ],
      ),
      body: Column(
        children: [
          // Scanner Area
          Expanded(
            flex: 3,
            child: _buildScannerArea(),
          ),
          
          // Result or Manual Entry Area
          Expanded(
            flex: 2,
            child: _buildBottomArea(),
          ),
        ],
      ),
    );
  }

  Widget _buildScannerArea() {
    if (_scanResult != null || _errorMessage != null) {
      return _buildResultDisplay();
    }

    return Stack(
      children: [
        MobileScanner(
          controller: _scannerController,
          onDetect: _onScan,
        ),
        // Scanner Overlay
        Center(
          child: Container(
            width: 250,
            height: 250,
            decoration: BoxDecoration(
              border: Border.all(color: const Color(0xFFf59e0b), width: 3),
              borderRadius: BorderRadius.circular(20),
            ),
            child: _isProcessing
                ? const Center(
                    child: CircularProgressIndicator(color: Color(0xFFf59e0b)),
                  )
                : null,
          ),
        ),
        // Hint text
        Positioned(
          bottom: 20,
          left: 0,
          right: 0,
          child: Text(
            'Placez le QR code du ticket dans le cadre',
            textAlign: TextAlign.center,
            style: TextStyle(
              color: Colors.white.withOpacity(0.7),
              fontSize: 14,
            ),
          ),
        ),
      ],
    );
  }

  Widget _buildResultDisplay() {
    final isValid = _scanResult != null && _scanResult!['valid'] == true;
    final ticket = _scanResult?['ticket'] as Map<String, dynamic>?;
    final isUsed = ticket?['used'] == true;

    return Container(
      margin: const EdgeInsets.all(20),
      padding: const EdgeInsets.all(24),
      decoration: BoxDecoration(
        color: Colors.white.withOpacity(0.05),
        borderRadius: BorderRadius.circular(20),
        border: Border.all(
          color: _errorMessage != null
              ? Colors.red
              : (isUsed ? Colors.orange : Colors.green),
          width: 2,
        ),
      ),
      child: Column(
        mainAxisAlignment: MainAxisAlignment.center,
        children: [
          // Icon
          Icon(
            _errorMessage != null
                ? Icons.cancel
                : (isUsed ? Icons.warning : Icons.check_circle),
            size: 64,
            color: _errorMessage != null
                ? Colors.red
                : (isUsed ? Colors.orange : Colors.green),
          ),
          const SizedBox(height: 16),
          
          // Title
          Text(
            _errorMessage != null
                ? 'Ticket Invalide'
                : (isUsed ? 'Déjà Utilisé' : 'Ticket Valide!'),
            style: const TextStyle(
              color: Colors.white,
              fontSize: 24,
              fontWeight: FontWeight.bold,
            ),
          ),
          const SizedBox(height: 12),
          
          // Details
          if (_errorMessage != null)
            Text(
              _errorMessage!,
              style: const TextStyle(color: Colors.white70),
              textAlign: TextAlign.center,
            )
          else if (ticket != null) ...[
            _buildInfoRow('Code', ticket['ticket_code'] ?? 'N/A'),
            _buildInfoRow('Type', ticket['tier_name'] ?? 'Standard'),
            _buildInfoRow('Statut', isUsed ? 'Utilisé' : 'Non utilisé'),
          ],
        ],
      ),
    );
  }

  Widget _buildInfoRow(String label, String value) {
    return Padding(
      padding: const EdgeInsets.symmetric(vertical: 4),
      child: Row(
        mainAxisAlignment: MainAxisAlignment.center,
        children: [
          Text(
            '$label: ',
            style: const TextStyle(color: Colors.white54, fontSize: 14),
          ),
          Text(
            value,
            style: const TextStyle(color: Colors.white, fontSize: 14, fontWeight: FontWeight.w600),
          ),
        ],
      ),
    );
  }

  Widget _buildBottomArea() {
    return Container(
      padding: const EdgeInsets.all(20),
      decoration: BoxDecoration(
        color: const Color(0xFF16213e),
        borderRadius: const BorderRadius.vertical(top: Radius.circular(24)),
      ),
      child: Column(
        children: [
          // Manual code entry
          if (_scanResult == null && _errorMessage == null) ...[
            const Text(
              'Ou entrez le code manuellement',
              style: TextStyle(color: Colors.white70, fontSize: 14),
            ),
            const SizedBox(height: 12),
            Row(
              children: [
                Expanded(
                  child: TextField(
                    controller: _manualCodeController,
                    style: const TextStyle(color: Colors.white, fontFamily: 'monospace'),
                    decoration: InputDecoration(
                      hintText: 'TKT-XXXXX',
                      hintStyle: const TextStyle(color: Colors.white38),
                      filled: true,
                      fillColor: Colors.white.withOpacity(0.1),
                      border: OutlineInputBorder(
                        borderRadius: BorderRadius.circular(12),
                        borderSide: BorderSide.none,
                      ),
                    ),
                  ),
                ),
                const SizedBox(width: 12),
                ElevatedButton(
                  onPressed: _isProcessing ? null : _verifyManualCode,
                  style: ElevatedButton.styleFrom(
                    backgroundColor: const Color(0xFF6366f1),
                    padding: const EdgeInsets.symmetric(horizontal: 24, vertical: 16),
                    shape: RoundedRectangleBorder(
                      borderRadius: BorderRadius.circular(12),
                    ),
                  ),
                  child: const Text('Vérifier', style: TextStyle(color: Colors.white)),
                ),
              ],
            ),
          ],
          
          // Action buttons when result is shown
          if (_scanResult != null || _errorMessage != null) ...[
            const Spacer(),
            if (_scanResult != null && 
                _scanResult!['valid'] == true && 
                _scanResult!['ticket']?['used'] != true)
              SizedBox(
                width: double.infinity,
                child: ElevatedButton.icon(
                  onPressed: _isProcessing ? null : _markAsUsed,
                  icon: _isProcessing
                      ? const SizedBox(
                          width: 20,
                          height: 20,
                          child: CircularProgressIndicator(color: Colors.white, strokeWidth: 2),
                        )
                      : const Icon(Icons.check, color: Colors.white),
                  label: const Text('Marquer comme utilisé', style: TextStyle(color: Colors.white)),
                  style: ElevatedButton.styleFrom(
                    backgroundColor: Colors.green,
                    padding: const EdgeInsets.symmetric(vertical: 16),
                    shape: RoundedRectangleBorder(
                      borderRadius: BorderRadius.circular(12),
                    ),
                  ),
                ),
              ),
            const SizedBox(height: 12),
            SizedBox(
              width: double.infinity,
              child: OutlinedButton.icon(
                onPressed: _resetScanner,
                icon: const Icon(Icons.qr_code_scanner, color: Colors.white),
                label: const Text('Scanner un autre ticket', style: TextStyle(color: Colors.white)),
                style: OutlinedButton.styleFrom(
                  padding: const EdgeInsets.symmetric(vertical: 16),
                  side: const BorderSide(color: Colors.white38),
                  shape: RoundedRectangleBorder(
                    borderRadius: BorderRadius.circular(12),
                  ),
                ),
              ),
            ),
          ],
        ],
      ),
    );
  }
}
