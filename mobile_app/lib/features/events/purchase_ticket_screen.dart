import 'package:flutter/material.dart';
import '../../core/services/ticket_api_service.dart';
import '../../core/services/wallet_api_service.dart';

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
  String? _selectedWalletId;
  String _pin = '';
  Map<String, String> _formData = {};
  bool _loading = false;
  bool _purchasing = false;

  @override
  void initState() {
    super.initState();
    _loadWallets();
    _initializeFormData();
  }

  void _initializeFormData() {
    final fields = widget.event['form_fields'] as List? ?? [];
    for (var field in fields) {
      _formData[field['name'] ?? ''] = '';
    }
  }

  Future<void> _loadWallets() async {
    setState(() => _loading = true);
    try {
      _wallets = await _walletApi.getWallets();
    } catch (e) {
      debugPrint('Error loading wallets: $e');
    } finally {
      setState(() => _loading = false);
    }
  }

  Future<void> _purchaseTicket() async {
    if (!_formKey.currentState!.validate()) return;
    if (_selectedWalletId == null) {
      ScaffoldMessenger.of(context).showSnackBar(
        const SnackBar(content: Text('Veuillez sélectionner un portefeuille')),
      );
      return;
    }

    setState(() => _purchasing = true);
    try {
      final result = await _ticketApi.purchaseTicket(
        eventId: widget.event['id'],
        tierId: widget.tier['id'],
        formData: _formData,
        walletId: _selectedWalletId!,
        pin: _pin,
      );

      if (!mounted) return;
      
      // Show success modal
      showDialog(
        context: context,
        barrierDismissible: false,
        builder: (_) => _buildSuccessDialog(result),
      );
    } catch (e) {
      if (!mounted) return;
      ScaffoldMessenger.of(context).showSnackBar(
        SnackBar(
          content: Text(e.toString().replaceAll('Exception: ', '')),
          backgroundColor: Colors.red,
        ),
      );
    } finally {
      if (mounted) setState(() => _purchasing = false);
    }
  }

  Widget _buildSuccessDialog(Map<String, dynamic> result) {
    final ticket = result['ticket'];
    
    return Dialog(
      shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(20)),
      child: Container(
        padding: const EdgeInsets.all(24),
        decoration: BoxDecoration(
          gradient: const LinearGradient(
            begin: Alignment.topLeft,
            end: Alignment.bottomRight,
            colors: [Color(0xFF1a1a2e), Color(0xFF16213e)],
          ),
          borderRadius: BorderRadius.circular(20),
        ),
        child: Column(
          mainAxisSize: MainAxisSize.min,
          children: [
            Container(
              width: 80,
              height: 80,
              decoration: BoxDecoration(
                color: Colors.green.withOpacity(0.2),
                shape: BoxShape.circle,
              ),
              child: const Icon(Icons.check, size: 50, color: Colors.green),
            ),
            const SizedBox(height: 20),
            const Text(
              'Achat réussi !',
              style: TextStyle(
                color: Colors.white,
                fontSize: 24,
                fontWeight: FontWeight.bold,
              ),
            ),
            const SizedBox(height: 12),
            Text(
              'Votre ticket ${widget.tier['name']} a été acheté',
              style: const TextStyle(color: Colors.white70),
              textAlign: TextAlign.center,
            ),
            const SizedBox(height: 24),
            if (ticket?['qr_code'] != null)
              Container(
                padding: const EdgeInsets.all(12),
                decoration: BoxDecoration(
                  color: Colors.white,
                  borderRadius: BorderRadius.circular(12),
                ),
                child: Image.network(
                  ticket['qr_code'],
                  width: 150,
                  height: 150,
                  errorBuilder: (_, __, ___) => const SizedBox(
                    width: 150,
                    height: 150,
                    child: Icon(Icons.qr_code, size: 100),
                  ),
                ),
              ),
            const SizedBox(height: 12),
            if (ticket?['ticket_code'] != null)
              Container(
                padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 8),
                decoration: BoxDecoration(
                  color: Colors.white.withOpacity(0.1),
                  borderRadius: BorderRadius.circular(8),
                ),
                child: Text(
                  ticket['ticket_code'],
                  style: const TextStyle(
                    color: Colors.white,
                    fontSize: 16,
                    fontWeight: FontWeight.bold,
                    fontFamily: 'monospace',
                  ),
                ),
              ),
            const SizedBox(height: 24),
            SizedBox(
              width: double.infinity,
              child: ElevatedButton(
                onPressed: () {
                  Navigator.of(context).pop(); // Close dialog
                  Navigator.of(context).pop(); // Go back to event details
                },
                style: ElevatedButton.styleFrom(
                  backgroundColor: const Color(0xFF6366f1),
                  padding: const EdgeInsets.all(16),
                  shape: RoundedRectangleBorder(
                    borderRadius: BorderRadius.circular(12),
                  ),
                ),
                child: const Text(
                  'Fermer',
                  style: TextStyle(fontSize: 16, fontWeight: FontWeight.w600),
                ),
              ),
            ),
          ],
        ),
      ),
    );
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
                child: _loading
                    ? const Center(child: CircularProgressIndicator())
                    : SingleChildScrollView(
                        padding: const EdgeInsets.all(20),
                        child: Form(
                          key: _formKey,
                          child: Column(
                            crossAxisAlignment: CrossAxisAlignment.start,
                            children: [
                              _buildTierCard(),
                              const SizedBox(height: 24),
                              _buildFormFields(),
                              const SizedBox(height: 24),
                              _buildPaymentSection(),
                              const SizedBox(height: 32),
                              _buildPurchaseButton(),
                            ],
                          ),
                        ),
                      ),
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
            onPressed: () => Navigator.pop(context),
            icon: const Icon(Icons.arrow_back, color: Colors.white),
          ),
          const SizedBox(width: 8),
          const Expanded(
            child: Text(
              '🎫 Acheter un ticket',
              style: TextStyle(
                fontSize: 22,
                fontWeight: FontWeight.bold,
                color: Colors.white,
              ),
            ),
          ),
        ],
      ),
    );
  }

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
      child: Row(
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
    );
  }

  Widget _buildFormFields() {
    final fields = widget.event['form_fields'] as List? ?? [];
    
    if (fields.isEmpty) {
      return const SizedBox.shrink();
    }

    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        const Text(
          '📝 Informations',
          style: TextStyle(
            color: Colors.white,
            fontSize: 18,
            fontWeight: FontWeight.bold,
          ),
        ),
        const SizedBox(height: 12),
        ...fields.map((field) => Padding(
          padding: const EdgeInsets.only(bottom: 16),
          child: TextFormField(
            style: const TextStyle(color: Colors.white),
            decoration: InputDecoration(
              labelText: field['label'] ?? '',
              labelStyle: TextStyle(color: Colors.white.withOpacity(0.7)),
              filled: true,
              fillColor: const Color(0xFF1e293b),
              border: OutlineInputBorder(
                borderRadius: BorderRadius.circular(12),
                borderSide: const BorderSide(color: Color(0xFF475569), width: 1.5),
              ),
              enabledBorder: OutlineInputBorder(
                borderRadius: BorderRadius.circular(12),
                borderSide: const BorderSide(color: Color(0xFF475569), width: 1.5),
              ),
              focusedBorder: OutlineInputBorder(
                borderRadius: BorderRadius.circular(12),
                borderSide: const BorderSide(color: Color(0xFF6366f1), width: 2),
              ),
              errorBorder: OutlineInputBorder(
                borderRadius: BorderRadius.circular(12),
                borderSide: const BorderSide(color: Colors.red, width: 1.5),
              ),
              contentPadding: const EdgeInsets.symmetric(horizontal: 16, vertical: 16),
            ),
            keyboardType: field['type'] == 'email'
                ? TextInputType.emailAddress
                : field['type'] == 'phone'
                    ? TextInputType.phone
                    : field['type'] == 'number'
                        ? TextInputType.number
                        : TextInputType.text,
            validator: (value) {
              if (field['required'] == true && (value == null || value.isEmpty)) {
                return 'Ce champ est requis';
              }
              if (field['type'] == 'email' && value != null && value.isNotEmpty) {
                if (!value.contains('@')) {
                  return 'Email invalide';
                }
              }
              return null;
            },
            onChanged: (value) {
              setState(() {
                _formData[field['name'] ?? ''] = value;
              });
            },
          ),
        )),
      ],
    );
  }

  Widget _buildPaymentSection() {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        const Text(
          '💳 Paiement',
          style: TextStyle(
            color: Colors.white,
            fontSize: 18,
            fontWeight: FontWeight.bold,
          ),
        ),
        const SizedBox(height: 12),
        DropdownButtonFormField<String>(
          value: _selectedWalletId,
          dropdownColor: const Color(0xFF1e293b),
          decoration: InputDecoration(
            labelText: 'Sélectionner un portefeuille',
            labelStyle: TextStyle(color: Colors.white.withOpacity(0.7)),
            filled: true,
            fillColor: const Color(0xFF1e293b),
            border: OutlineInputBorder(
              borderRadius: BorderRadius.circular(12),
              borderSide: const BorderSide(color: Color(0xFF475569), width: 1.5),
            ),
            enabledBorder: OutlineInputBorder(
              borderRadius: BorderRadius.circular(12),
              borderSide: const BorderSide(color: Color(0xFF475569), width: 1.5),
            ),
            focusedBorder: OutlineInputBorder(
              borderRadius: BorderRadius.circular(12),
              borderSide: const BorderSide(color: Color(0xFF6366f1), width: 2),
            ),
            contentPadding: const EdgeInsets.symmetric(horizontal: 16, vertical: 16),
          ),
          style: const TextStyle(color: Colors.white),
          items: _wallets.map<DropdownMenuItem<String>>((wallet) {
            return DropdownMenuItem<String>(
              value: wallet['id'],
              child: Text(
                '${wallet['currency']} - ${_formatAmount(wallet['balance'])}',
                style: const TextStyle(color: Colors.white),
              ),
            );
          }).toList(),
          onChanged: (value) {
            setState(() => _selectedWalletId = value);
          },
          validator: (value) => value == null ? 'Sélectionnez un portefeuille' : null,
        ),
        const SizedBox(height: 16),
        TextFormField(
          style: const TextStyle(color: Colors.white),
          obscureText: true,
          maxLength: 5,
          keyboardType: TextInputType.number,
          decoration: InputDecoration(
            labelText: 'Code PIN',
            labelStyle: TextStyle(color: Colors.white.withOpacity(0.7)),
            filled: true,
            fillColor: const Color(0xFF1e293b),
            border: OutlineInputBorder(
              borderRadius: BorderRadius.circular(12),
              borderSide: const BorderSide(color: Color(0xFF475569), width: 1.5),
            ),
            enabledBorder: OutlineInputBorder(
              borderRadius: BorderRadius.circular(12),
              borderSide: const BorderSide(color: Color(0xFF475569), width: 1.5),
            ),
            focusedBorder: OutlineInputBorder(
              borderRadius: BorderRadius.circular(12),
              borderSide: const BorderSide(color: Color(0xFF6366f1), width: 2),
            ),
            errorBorder: OutlineInputBorder(
              borderRadius: BorderRadius.circular(12),
              borderSide: const BorderSide(color: Colors.red, width: 1.5),
            ),
            counterText: '',
            contentPadding: const EdgeInsets.symmetric(horizontal: 16, vertical: 16),
          ),
          validator: (value) {
            if (value == null || value.isEmpty) {
              return 'Code PIN requis';
            }
            if (value.length != 5) {
              return 'Le PIN doit contenir 5 chiffres';
            }
            return null;
          },
          onChanged: (value) {
            setState(() => _pin = value);
          },
        ),
      ],
    );
  }

  Widget _buildPurchaseButton() {
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
                'Acheter - ${_formatAmount(widget.tier['price'] ?? 0)} ${widget.event['currency'] ?? 'XOF'}',
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
