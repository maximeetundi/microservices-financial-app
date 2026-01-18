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
  bool _loading = true;
  bool _purchasing = false;
  final Map<String, TextEditingController> _formControllers = {};

  @override
  void initState() {
    super.initState();
    _loadWallets();
    _initFormControllers();
  }

  void _initFormControllers() {
    final formFields = widget.event['form_fields'] as List<dynamic>? ?? [];
    for (final field in formFields) {
      final name = field['name']?.toString() ?? '';
      _formControllers[name] = TextEditingController();
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

  int get _maxAllowed {
    int max = 100;
    int maxPerUser = widget.tier['max_per_user'] ?? 0;
    int total = widget.tier['quantity'] ?? -1;
    int sold = widget.tier['sold'] ?? 0;
    
    if (total != -1) {
      int remaining = total - sold;
      if (remaining < max) max = remaining;
    }
    
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

  Future<void> _purchaseTicket() async {
    if (!_formKey.currentState!.validate()) return;

    // Check we have wallets
    if (_wallets.isEmpty) {
      ScaffoldMessenger.of(context).showSnackBar(
        const SnackBar(content: Text('Aucun portefeuille disponible')),
      );
      return;
    }

    // Verify PIN and get the raw PIN
    final pinResult = await PinVerifyDialog.show(
      context,
      title: 'Confirmer l\'achat',
      subtitle: 'Entrez votre PIN pour valider',
      allowBiometric: true,
      returnRawPin: true,
    );

    if (pinResult == null || pinResult == false) return;
    final String pin = pinResult is String ? pinResult : '';
    if (pin.isEmpty) return;

    setState(() => _purchasing = true);
    try {
      // Prepare form data
      final formData = <String, String>{};
      for (final entry in _formControllers.entries) {
        formData[entry.key] = entry.value.text;
      }

      // Use first wallet with matching currency or first available
      final currency = widget.event['currency'] ?? 'XOF';
      final wallet = _wallets.firstWhere(
        (w) => w['currency'] == currency,
        orElse: () => _wallets.first,
      );

      await _ticketApi.purchaseTicket(
        eventId: widget.event['id'],
        tierId: widget.tier['id'],
        quantity: _quantity,
        formData: formData,
        walletId: wallet['id'].toString(),
        pin: pin,
      );

      if (!mounted) return;
      ScaffoldMessenger.of(context).showSnackBar(
        const SnackBar(content: Text('✅ Ticket(s) acheté(s) avec succès!'), backgroundColor: Colors.green),
      );
      Navigator.pop(context, true);
    } catch (e) {
      if (!mounted) return;
      ScaffoldMessenger.of(context).showSnackBar(
        SnackBar(content: Text('Erreur: ${e.toString().replaceAll("Exception: ", "")}')),
      );
    } finally {
      if (mounted) setState(() => _purchasing = false);
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
          child: _loading 
              ? const Center(child: CircularProgressIndicator())
              : Column(
                  children: [
                    _buildHeader(),
                    Expanded(
                      child: SingleChildScrollView(
                        padding: const EdgeInsets.all(20),
                        child: Form(
                          key: _formKey,
                          child: Column(
                            crossAxisAlignment: CrossAxisAlignment.start,
                            children: [
                              _buildTierCard(),
                              const SizedBox(height: 24),
                              _buildFormFields(),
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
          Expanded(
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                const Text(
                  '🎫 Acheter un ticket',
                  style: TextStyle(
                    fontSize: 20,
                    fontWeight: FontWeight.bold,
                    color: Colors.white,
                  ),
                ),
                Text(
                  widget.event['title'] ?? '',
                  style: TextStyle(
                    color: Colors.white.withOpacity(0.7),
                    fontSize: 14,
                  ),
                  maxLines: 1,
                  overflow: TextOverflow.ellipsis,
                ),
              ],
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

  Widget _buildFormFields() {
    final formFields = widget.event['form_fields'] as List<dynamic>? ?? [];
    if (formFields.isEmpty) return const SizedBox.shrink();

    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        const Text(
          '📝 Informations du participant',
          style: TextStyle(color: Colors.white, fontSize: 18, fontWeight: FontWeight.bold),
        ),
        const SizedBox(height: 12),
        ...formFields.map((field) {
          final name = field['name']?.toString() ?? '';
          final label = field['label']?.toString() ?? name;
          final required = field['required'] == true;
          
          return Padding(
            padding: const EdgeInsets.only(bottom: 16),
            child: TextFormField(
              controller: _formControllers[name],
              style: const TextStyle(color: Colors.white),
              decoration: InputDecoration(
                labelText: label + (required ? ' *' : ''),
                labelStyle: TextStyle(color: Colors.white.withOpacity(0.7)),
                filled: true,
                fillColor: Colors.white.withOpacity(0.1),
                border: OutlineInputBorder(
                  borderRadius: BorderRadius.circular(12),
                  borderSide: BorderSide.none,
                ),
              ),
              validator: required 
                  ? (v) => v == null || v.isEmpty ? 'Ce champ est requis' : null
                  : null,
            ),
          );
        }),
      ],
    );
  }

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

  @override
  void dispose() {
    for (final controller in _formControllers.values) {
      controller.dispose();
    }
    super.dispose();
  }
}
