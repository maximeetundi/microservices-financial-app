import 'package:flutter/material.dart';
import 'package:mobile_scanner/mobile_scanner.dart';
import '../../core/services/ticket_api_service.dart';

class VerifyTicketScreen extends StatefulWidget {
  const VerifyTicketScreen({super.key});

  @override
  State<VerifyTicketScreen> createState() => _VerifyTicketScreenState();
}

class _VerifyTicketScreenState extends State<VerifyTicketScreen> {
  final TicketApiService _ticketApi = TicketApiService();
  final TextEditingController _codeController = TextEditingController();
  
  bool _scanning = false;
  bool _verifying = false;
  Map<String, dynamic>? _result;

  void _onQRScanned(String code) async {
    if (_verifying) return;
    
    setState(() {
      _scanning = false;
      _verifying = true;
    });

    try {
      final result = await _ticketApi.verifyTicket(code);
      setState(() => _result = result);
    } catch (e) {
      setState(() => _result = {
        'valid': false,
        'message': 'Erreur: $e',
      });
    } finally {
      setState(() => _verifying = false);
    }
  }

  void _manualVerify() async {
    final code = _codeController.text.trim();
    if (code.isEmpty) return;
    
    _onQRScanned(code);
  }

  void _markAsUsed() async {
    final ticketId = _result?['ticket']?['id'];
    if (ticketId == null) return;

    try {
      await _ticketApi.useTicket(ticketId);
      ScaffoldMessenger.of(context).showSnackBar(
        const SnackBar(
          content: Text('✅ Ticket marqué comme utilisé'),
          backgroundColor: Colors.green,
        ),
      );
      setState(() => _result = null);
    } catch (e) {
      ScaffoldMessenger.of(context).showSnackBar(
        SnackBar(
          content: Text('Erreur: $e'),
          backgroundColor: Colors.red,
        ),
      );
    }
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      body: Container(
        decoration: const BoxDecoration(
          gradient: LinearGradient(
            begin: Alignment.topLeft,
            end: Alignment.bottomRight,
            colors: [Color(0xFF1a1a2e), Color(0xFF16213e)],
          ),
        ),
        child: SafeArea(
          child: Column(
            children: [
              _buildHeader(),
              Expanded(
                child: _scanning
                    ? _buildScanner()
                    : _result != null
                        ? _buildResult()
                        : _buildMainContent(),
              ),
            ],
          ),
        ),
      ),
    );
  }

  Widget _buildHeader() {
    return Padding(
      padding: const EdgeInsets.all(20),
      child: Row(
        children: [
          IconButton(
            onPressed: () {
              if (_scanning) {
                setState(() => _scanning = false);
              } else {
                Navigator.pop(context);
              }
            },
            icon: const Icon(Icons.arrow_back, color: Colors.white),
          ),
          const SizedBox(width: 8),
          const Text(
            '🔍 Vérifier un ticket',
            style: TextStyle(
              fontSize: 22,
              fontWeight: FontWeight.bold,
              color: Colors.white,
            ),
          ),
        ],
      ),
    );
  }

  Widget _buildMainContent() {
    return Padding(
      padding: const EdgeInsets.all(20),
      child: Column(
        mainAxisAlignment: MainAxisAlignment.center,
        children: [
          // Scan button
          GestureDetector(
            onTap: () => setState(() => _scanning = true),
            child: Container(
              width: 200,
              height: 200,
              decoration: BoxDecoration(
                color: const Color(0xFF6366f1).withOpacity(0.2),
                borderRadius: BorderRadius.circular(100),
                border: Border.all(
                  color: const Color(0xFF6366f1),
                  width: 3,
                ),
              ),
              child: const Column(
                mainAxisAlignment: MainAxisAlignment.center,
                children: [
                  Icon(Icons.qr_code_scanner, size: 64, color: Colors.white),
                  SizedBox(height: 12),
                  Text(
                    'Scanner QR',
                    style: TextStyle(
                      color: Colors.white,
                      fontSize: 18,
                      fontWeight: FontWeight.w600,
                    ),
                  ),
                ],
              ),
            ),
          ),
          const SizedBox(height: 40),
          
          // Divider
          Row(
            children: [
              Expanded(child: Divider(color: Colors.white.withOpacity(0.2))),
              Padding(
                padding: const EdgeInsets.symmetric(horizontal: 16),
                child: Text(
                  'ou',
                  style: TextStyle(color: Colors.white.withOpacity(0.5)),
                ),
              ),
              Expanded(child: Divider(color: Colors.white.withOpacity(0.2))),
            ],
          ),
          const SizedBox(height: 24),
          
          // Manual input
          TextField(
            controller: _codeController,
            style: const TextStyle(color: Colors.white),
            decoration: InputDecoration(
              hintText: 'Entrer le code du ticket',
              hintStyle: TextStyle(color: Colors.white.withOpacity(0.4)),
              filled: true,
              fillColor: Colors.white.withOpacity(0.1),
              border: OutlineInputBorder(
                borderRadius: BorderRadius.circular(12),
                borderSide: BorderSide.none,
              ),
              suffixIcon: IconButton(
                onPressed: _verifying ? null : _manualVerify,
                icon: _verifying
                    ? const SizedBox(
                        width: 20,
                        height: 20,
                        child: CircularProgressIndicator(strokeWidth: 2),
                      )
                    : const Icon(Icons.search, color: Colors.white),
              ),
            ),
            onSubmitted: (_) => _manualVerify(),
          ),
        ],
      ),
    );
  }

  Widget _buildScanner() {
    return Stack(
      children: [
        MobileScanner(
          onDetect: (capture) {
            final barcodes = capture.barcodes;
            if (barcodes.isNotEmpty && barcodes.first.rawValue != null) {
              _onQRScanned(barcodes.first.rawValue!);
            }
          },
        ),
        // Overlay
        Center(
          child: Container(
            width: 250,
            height: 250,
            decoration: BoxDecoration(
              border: Border.all(color: const Color(0xFF6366f1), width: 3),
              borderRadius: BorderRadius.circular(20),
            ),
          ),
        ),
        Positioned(
          bottom: 40,
          left: 0,
          right: 0,
          child: Center(
            child: Container(
              padding: const EdgeInsets.symmetric(horizontal: 20, vertical: 12),
              decoration: BoxDecoration(
                color: Colors.black.withOpacity(0.7),
                borderRadius: BorderRadius.circular(20),
              ),
              child: const Text(
                'Placez le QR code dans le cadre',
                style: TextStyle(color: Colors.white),
              ),
            ),
          ),
        ),
      ],
    );
  }

  Widget _buildResult() {
    final result = _result!;
    final valid = result['valid'] == true;
    final canUse = result['can_use'] == true;
    final alreadyUsed = result['already_used'] == true;
    final ticket = result['ticket'];
    final event = result['event'];

    return SingleChildScrollView(
      padding: const EdgeInsets.all(20),
      child: Column(
        children: [
          // Status Icon
          Container(
            width: 100,
            height: 100,
            decoration: BoxDecoration(
              color: valid
                  ? (alreadyUsed ? Colors.orange : Colors.green).withOpacity(0.2)
                  : Colors.red.withOpacity(0.2),
              shape: BoxShape.circle,
            ),
            child: Icon(
              valid
                  ? (alreadyUsed ? Icons.history : Icons.check)
                  : Icons.close,
              size: 50,
              color: valid
                  ? (alreadyUsed ? Colors.orange : Colors.green)
                  : Colors.red,
            ),
          ),
          const SizedBox(height: 20),
          
          // Message
          Text(
            result['message'] ?? (valid ? 'Ticket valide' : 'Ticket invalide'),
            style: TextStyle(
              color: valid ? Colors.white : Colors.red,
              fontSize: 22,
              fontWeight: FontWeight.bold,
            ),
            textAlign: TextAlign.center,
          ),
          const SizedBox(height: 24),
          
          // Ticket details
          if (ticket != null) ...[
            Container(
              padding: const EdgeInsets.all(20),
              decoration: BoxDecoration(
                color: Colors.white.withOpacity(0.05),
                borderRadius: BorderRadius.circular(16),
                border: Border.all(color: Colors.white.withOpacity(0.1)),
              ),
              child: Column(
                children: [
                  Text(
                    ticket['tier_icon'] ?? '🎫',
                    style: const TextStyle(fontSize: 48),
                  ),
                  const SizedBox(height: 12),
                  Text(
                    event?['title'] ?? ticket['event_title'] ?? 'Événement',
                    style: const TextStyle(
                      color: Colors.white,
                      fontSize: 18,
                      fontWeight: FontWeight.bold,
                    ),
                    textAlign: TextAlign.center,
                  ),
                  const SizedBox(height: 8),
                  Text(
                    ticket['tier_name'] ?? 'Standard',
                    style: TextStyle(
                      color: Colors.white.withOpacity(0.7),
                      fontSize: 16,
                    ),
                  ),
                  const SizedBox(height: 16),
                  Text(
                    ticket['ticket_code'] ?? '',
                    style: const TextStyle(
                      color: Colors.white70,
                      fontSize: 14,
                      fontFamily: 'monospace',
                    ),
                  ),
                ],
              ),
            ),
          ],
          const SizedBox(height: 24),
          
          // Actions
          Row(
            children: [
              Expanded(
                child: OutlinedButton(
                  onPressed: () => setState(() => _result = null),
                  style: OutlinedButton.styleFrom(
                    foregroundColor: Colors.white,
                    side: const BorderSide(color: Colors.white54),
                    padding: const EdgeInsets.all(16),
                    shape: RoundedRectangleBorder(
                      borderRadius: BorderRadius.circular(12),
                    ),
                  ),
                  child: const Text('Scanner un autre'),
                ),
              ),
              if (canUse) ...[
                const SizedBox(width: 12),
                Expanded(
                  child: ElevatedButton(
                    onPressed: _markAsUsed,
                    style: ElevatedButton.styleFrom(
                      backgroundColor: Colors.green,
                      padding: const EdgeInsets.all(16),
                      shape: RoundedRectangleBorder(
                        borderRadius: BorderRadius.circular(12),
                      ),
                    ),
                    child: const Text('✓ Valider l\'entrée'),
                  ),
                ),
              ],
            ],
          ),
        ],
      ),
    );
  }

  @override
  void dispose() {
    _codeController.dispose();
    super.dispose();
  }
}
