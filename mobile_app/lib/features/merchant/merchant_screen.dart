import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import 'dart:convert';
import '../../core/services/api_service.dart';
import '../../core/widgets/custom_button.dart';
import '../../core/utils/currency_formatter.dart';

class MerchantScreen extends StatefulWidget {
  const MerchantScreen({Key? key}) : super(key: key);

  @override
  State<MerchantScreen> createState() => _MerchantScreenState();
}

class _MerchantScreenState extends State<MerchantScreen> with SingleTickerProviderStateMixin {
  late TabController _tabController;
  List<dynamic> _payments = [];
  List<dynamic> _wallets = [];
  bool _loading = true;
  Map<String, dynamic> _stats = {};

  @override
  void initState() {
    super.initState();
    _tabController = TabController(length: 2, vsync: this);
    _loadData();
  }

  Future<void> _loadData() async {
    setState(() => _loading = true);
    try {
      final payments = await ApiService.getMerchantPayments();
      final wallets = await ApiService.getWallets();
      setState(() {
        _payments = payments['payments'] ?? [];
        _wallets = wallets ?? [];
      });
    } catch (e) {
      debugPrint('Error loading data: $e');
    } finally {
      setState(() => _loading = false);
    }
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: const Color(0xFF0F0F23),
      appBar: AppBar(
        title: const Text('💼 Espace Marchand'),
        backgroundColor: Colors.transparent,
        elevation: 0,
        bottom: TabBar(
          controller: _tabController,
          indicatorColor: const Color(0xFF6366F1),
          tabs: const [
            Tab(text: 'En attente'),
            Tab(text: 'Historique'),
          ],
        ),
      ),
      body: _loading
          ? const Center(child: CircularProgressIndicator())
          : TabBarView(
              controller: _tabController,
              children: [
                _buildPaymentsList(_payments.where((p) => p['status'] == 'pending').toList()),
                _buildPaymentsList(_payments.where((p) => p['status'] != 'pending').toList()),
              ],
            ),
      floatingActionButton: FloatingActionButton.extended(
        onPressed: () => _showCreatePaymentSheet(context),
        backgroundColor: const Color(0xFF6366F1),
        icon: const Icon(Icons.add),
        label: const Text('Nouveau'),
      ),
    );
  }

  Widget _buildPaymentsList(List<dynamic> payments) {
    if (payments.isEmpty) {
      return Center(
        child: Column(
          mainAxisAlignment: MainAxisAlignment.center,
          children: [
            const Icon(Icons.receipt_long, size: 64, color: Colors.grey),
            const SizedBox(height: 16),
            const Text(
              'Aucun paiement',
              style: TextStyle(color: Colors.grey, fontSize: 18),
            ),
          ],
        ),
      );
    }

    return ListView.builder(
      padding: const EdgeInsets.all(16),
      itemCount: payments.length,
      itemBuilder: (context, index) {
        final payment = payments[index];
        return _PaymentCard(
          payment: payment,
          onTap: () => _showQRCode(payment),
        );
      },
    );
  }

  void _showCreatePaymentSheet(BuildContext context) {
    showModalBottomSheet(
      context: context,
      isScrollControlled: true,
      backgroundColor: const Color(0xFF1A1A3E),
      shape: const RoundedRectangleBorder(
        borderRadius: BorderRadius.vertical(top: Radius.circular(20)),
      ),
      builder: (context) => _CreatePaymentSheet(
        wallets: _wallets,
        onCreated: (payment) {
          setState(() => _payments.insert(0, payment));
          Navigator.pop(context);
          _showQRCode(payment);
        },
      ),
    );
  }

  void _showQRCode(dynamic payment) {
    showModalBottomSheet(
      context: context,
      isScrollControlled: true,
      backgroundColor: const Color(0xFF1A1A3E),
      shape: const RoundedRectangleBorder(
        borderRadius: BorderRadius.vertical(top: Radius.circular(20)),
      ),
      builder: (context) => _QRCodeSheet(payment: payment),
    );
  }
}

class _PaymentCard extends StatelessWidget {
  final dynamic payment;
  final VoidCallback onTap;

  const _PaymentCard({required this.payment, required this.onTap});

  @override
  Widget build(BuildContext context) {
    final status = payment['status'] ?? 'pending';
    final amount = payment['amount'];
    final currency = payment['currency'] ?? 'EUR';

    return Card(
      color: const Color(0xFF1A1A3E),
      margin: const EdgeInsets.only(bottom: 12),
      shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(12)),
      child: InkWell(
        onTap: onTap,
        borderRadius: BorderRadius.circular(12),
        child: Padding(
          padding: const EdgeInsets.all(16),
          child: Row(
            children: [
              Container(
                width: 48,
                height: 48,
                decoration: BoxDecoration(
                  color: _getStatusColor(status).withOpacity(0.1),
                  borderRadius: BorderRadius.circular(12),
                ),
                child: Icon(
                  _getStatusIcon(status),
                  color: _getStatusColor(status),
                ),
              ),
              const SizedBox(width: 16),
              Expanded(
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    Text(
                      payment['title'] ?? 'Paiement',
                      style: const TextStyle(
                        color: Colors.white,
                        fontWeight: FontWeight.w600,
                        fontSize: 16,
                      ),
                    ),
                    const SizedBox(height: 4),
                    Text(
                      _getTypeLabel(payment['type']),
                      style: TextStyle(
                        color: Colors.white.withOpacity(0.6),
                        fontSize: 13,
                      ),
                    ),
                  ],
                ),
              ),
              Column(
                crossAxisAlignment: CrossAxisAlignment.end,
                children: [
                  Text(
                    amount != null
                        ? CurrencyFormatter.format(amount, currency)
                        : 'Variable',
                    style: const TextStyle(
                      color: Color(0xFF22C55E),
                      fontWeight: FontWeight.bold,
                      fontSize: 16,
                    ),
                  ),
                  const SizedBox(height: 4),
                  Container(
                    padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 4),
                    decoration: BoxDecoration(
                      color: _getStatusColor(status).withOpacity(0.1),
                      borderRadius: BorderRadius.circular(12),
                    ),
                    child: Text(
                      _getStatusLabel(status),
                      style: TextStyle(
                        color: _getStatusColor(status),
                        fontSize: 12,
                        fontWeight: FontWeight.w500,
                      ),
                    ),
                  ),
                ],
              ),
              const SizedBox(width: 8),
              const Icon(Icons.qr_code, color: Colors.grey),
            ],
          ),
        ),
      ),
    );
  }

  Color _getStatusColor(String status) {
    switch (status) {
      case 'paid':
        return const Color(0xFF22C55E);
      case 'pending':
        return const Color(0xFFFBBF24);
      case 'expired':
        return const Color(0xFFEF4444);
      default:
        return Colors.grey;
    }
  }

  IconData _getStatusIcon(String status) {
    switch (status) {
      case 'paid':
        return Icons.check_circle;
      case 'pending':
        return Icons.schedule;
      case 'expired':
        return Icons.timer_off;
      default:
        return Icons.receipt;
    }
  }

  String _getStatusLabel(String status) {
    switch (status) {
      case 'paid':
        return 'Payé';
      case 'pending':
        return 'En attente';
      case 'expired':
        return 'Expiré';
      default:
        return status;
    }
  }

  String _getTypeLabel(String? type) {
    switch (type) {
      case 'fixed':
        return 'Prix fixe';
      case 'variable':
        return 'Montant variable';
      case 'invoice':
        return 'Facture';
      default:
        return 'Paiement';
    }
  }
}

class _CreatePaymentSheet extends StatefulWidget {
  final List<dynamic> wallets;
  final Function(dynamic) onCreated;

  const _CreatePaymentSheet({required this.wallets, required this.onCreated});

  @override
  State<_CreatePaymentSheet> createState() => _CreatePaymentSheetState();
}

class _CreatePaymentSheetState extends State<_CreatePaymentSheet> {
  final _formKey = GlobalKey<FormState>();
  String _type = 'fixed';
  String? _walletId;
  double? _amount;
  String _title = '';
  String _description = '';
  int _expiresIn = 60;
  bool _reusable = false;
  bool _loading = false;

  @override
  Widget build(BuildContext context) {
    return Container(
      padding: EdgeInsets.only(
        bottom: MediaQuery.of(context).viewInsets.bottom,
      ),
      child: SingleChildScrollView(
        padding: const EdgeInsets.all(24),
        child: Form(
          key: _formKey,
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            mainAxisSize: MainAxisSize.min,
            children: [
              const Text(
                'Nouvelle demande de paiement',
                style: TextStyle(
                  color: Colors.white,
                  fontSize: 20,
                  fontWeight: FontWeight.bold,
                ),
              ),
              const SizedBox(height: 24),
              
              // Type selector
              Row(
                children: [
                  Expanded(
                    child: _TypeButton(
                      label: 'Prix fixe',
                      icon: Icons.sell,
                      selected: _type == 'fixed',
                      onTap: () => setState(() => _type = 'fixed'),
                    ),
                  ),
                  const SizedBox(width: 12),
                  Expanded(
                    child: _TypeButton(
                      label: 'Variable',
                      icon: Icons.volunteer_activism,
                      selected: _type == 'variable',
                      onTap: () => setState(() => _type = 'variable'),
                    ),
                  ),
                ],
              ),
              const SizedBox(height: 20),

              // Wallet selection
              DropdownButtonFormField<String>(
                value: _walletId,
                decoration: _inputDecoration('Portefeuille de réception'),
                dropdownColor: const Color(0xFF1A1A3E),
                style: const TextStyle(color: Colors.white),
                items: widget.wallets.map<DropdownMenuItem<String>>((w) {
                  return DropdownMenuItem(
                    value: w['id'],
                    child: Text('${w['currency']} - ${CurrencyFormatter.format(w['balance'], w['currency'])}'),
                  );
                }).toList(),
                onChanged: (v) => setState(() => _walletId = v),
                validator: (v) => v == null ? 'Sélectionnez un portefeuille' : null,
              ),
              const SizedBox(height: 16),

              // Title
              TextFormField(
                decoration: _inputDecoration('Titre (ex: iPhone 15)'),
                style: const TextStyle(color: Colors.white),
                onChanged: (v) => _title = v,
                validator: (v) => v?.isEmpty == true ? 'Titre requis' : null,
              ),
              const SizedBox(height: 16),

              // Amount (for fixed)
              if (_type == 'fixed') ...[
                TextFormField(
                  decoration: _inputDecoration('Montant'),
                  style: const TextStyle(color: Colors.white),
                  keyboardType: TextInputType.number,
                  onChanged: (v) => _amount = double.tryParse(v),
                  validator: (v) {
                    if (_type == 'fixed' && (v?.isEmpty == true || double.tryParse(v!) == null)) {
                      return 'Montant requis';
                    }
                    return null;
                  },
                ),
                const SizedBox(height: 16),
              ],

              // Description
              TextFormField(
                decoration: _inputDecoration('Description (optionnel)'),
                style: const TextStyle(color: Colors.white),
                maxLines: 2,
                onChanged: (v) => _description = v,
              ),
              const SizedBox(height: 16),

              // Expiration
              DropdownButtonFormField<int>(
                value: _expiresIn,
                decoration: _inputDecoration('Expiration'),
                dropdownColor: const Color(0xFF1A1A3E),
                style: const TextStyle(color: Colors.white),
                items: const [
                  DropdownMenuItem(value: 60, child: Text('1 heure')),
                  DropdownMenuItem(value: 1440, child: Text('24 heures')),
                  DropdownMenuItem(value: 10080, child: Text('7 jours')),
                  DropdownMenuItem(value: -1, child: Text('Jamais')),
                ],
                onChanged: (v) => setState(() => _expiresIn = v ?? 60),
              ),
              const SizedBox(height: 16),

              // Reusable
              SwitchListTile(
                value: _reusable,
                onChanged: (v) => setState(() => _reusable = v),
                title: const Text('Réutilisable', style: TextStyle(color: Colors.white)),
                subtitle: Text('Plusieurs clients peuvent payer', style: TextStyle(color: Colors.white.withOpacity(0.6))),
                activeColor: const Color(0xFF6366F1),
                contentPadding: EdgeInsets.zero,
              ),
              const SizedBox(height: 24),

              // Submit
              SizedBox(
                width: double.infinity,
                child: ElevatedButton(
                  onPressed: _loading ? null : _submit,
                  style: ElevatedButton.styleFrom(
                    backgroundColor: const Color(0xFF6366F1),
                    padding: const EdgeInsets.symmetric(vertical: 16),
                    shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(12)),
                  ),
                  child: _loading
                      ? const SizedBox(
                          width: 20,
                          height: 20,
                          child: CircularProgressIndicator(strokeWidth: 2, color: Colors.white),
                        )
                      : const Text('Créer le QR code', style: TextStyle(fontSize: 16, fontWeight: FontWeight.w600)),
                ),
              ),
            ],
          ),
        ),
      ),
    );
  }

  InputDecoration _inputDecoration(String label) {
    return InputDecoration(
      labelText: label,
      labelStyle: TextStyle(color: Colors.white.withOpacity(0.6)),
      filled: true,
      fillColor: Colors.white.withOpacity(0.05),
      border: OutlineInputBorder(
        borderRadius: BorderRadius.circular(12),
        borderSide: BorderSide(color: Colors.white.withOpacity(0.1)),
      ),
      enabledBorder: OutlineInputBorder(
        borderRadius: BorderRadius.circular(12),
        borderSide: BorderSide(color: Colors.white.withOpacity(0.1)),
      ),
      focusedBorder: OutlineInputBorder(
        borderRadius: BorderRadius.circular(12),
        borderSide: const BorderSide(color: Color(0xFF6366F1)),
      ),
    );
  }

  Future<void> _submit() async {
    if (!_formKey.currentState!.validate()) return;

    setState(() => _loading = true);
    try {
      final wallet = widget.wallets.firstWhere((w) => w['id'] == _walletId);
      
      final result = await ApiService.createMerchantPayment({
        'type': _type,
        'wallet_id': _walletId,
        'amount': _type == 'fixed' ? _amount : null,
        'currency': wallet['currency'],
        'title': _title,
        'description': _description,
        'expires_in_minutes': _expiresIn,
        'reusable': _reusable,
      });
      
      widget.onCreated(result['payment_request']);
    } catch (e) {
      ScaffoldMessenger.of(context).showSnackBar(
        SnackBar(content: Text('Erreur: $e'), backgroundColor: Colors.red),
      );
    } finally {
      setState(() => _loading = false);
    }
  }
}

class _TypeButton extends StatelessWidget {
  final String label;
  final IconData icon;
  final bool selected;
  final VoidCallback onTap;

  const _TypeButton({
    required this.label,
    required this.icon,
    required this.selected,
    required this.onTap,
  });

  @override
  Widget build(BuildContext context) {
    return InkWell(
      onTap: onTap,
      borderRadius: BorderRadius.circular(12),
      child: Container(
        padding: const EdgeInsets.all(16),
        decoration: BoxDecoration(
          color: selected ? const Color(0xFF6366F1).withOpacity(0.1) : Colors.white.withOpacity(0.05),
          borderRadius: BorderRadius.circular(12),
          border: Border.all(
            color: selected ? const Color(0xFF6366F1) : Colors.white.withOpacity(0.1),
            width: 2,
          ),
        ),
        child: Column(
          children: [
            Icon(icon, color: selected ? const Color(0xFF6366F1) : Colors.grey, size: 28),
            const SizedBox(height: 8),
            Text(
              label,
              style: TextStyle(
                color: selected ? const Color(0xFF6366F1) : Colors.white,
                fontWeight: FontWeight.w500,
              ),
            ),
          ],
        ),
      ),
    );
  }
}

class _QRCodeSheet extends StatelessWidget {
  final dynamic payment;

  const _QRCodeSheet({required this.payment});

  @override
  Widget build(BuildContext context) {
    return Container(
      padding: const EdgeInsets.all(24),
      child: Column(
        mainAxisSize: MainAxisSize.min,
        children: [
          Text(
            payment['title'] ?? 'Paiement',
            style: const TextStyle(
              color: Colors.white,
              fontSize: 20,
              fontWeight: FontWeight.bold,
            ),
          ),
          const SizedBox(height: 24),
          
          // QR Code placeholder - in real app, use qr_flutter package
          Container(
            width: 200,
            height: 200,
            decoration: BoxDecoration(
              color: Colors.white,
              borderRadius: BorderRadius.circular(16),
            ),
            child: Center(
              child: Column(
                mainAxisAlignment: MainAxisAlignment.center,
                children: [
                  const Icon(Icons.qr_code_2, size: 120, color: Color(0xFF0F0F23)),
                  Text(
                    payment['id'] ?? '',
                    style: const TextStyle(fontSize: 10, color: Colors.grey),
                  ),
                ],
              ),
            ),
          ),
          const SizedBox(height: 24),

          if (payment['amount'] != null)
            Text(
              CurrencyFormatter.format(payment['amount'], payment['currency'] ?? 'EUR'),
              style: const TextStyle(
                color: Color(0xFF22C55E),
                fontSize: 32,
                fontWeight: FontWeight.bold,
              ),
            )
          else
            const Text(
              'Montant variable',
              style: TextStyle(color: Colors.grey, fontSize: 18),
            ),
          const SizedBox(height: 24),

          Row(
            children: [
              Expanded(
                child: OutlinedButton.icon(
                  onPressed: () => _copyLink(context),
                  icon: const Icon(Icons.copy),
                  label: const Text('Copier'),
                  style: OutlinedButton.styleFrom(
                    foregroundColor: Colors.white,
                    side: BorderSide(color: Colors.white.withOpacity(0.2)),
                    padding: const EdgeInsets.symmetric(vertical: 12),
                  ),
                ),
              ),
              const SizedBox(width: 12),
              Expanded(
                child: OutlinedButton.icon(
                  onPressed: () => _share(context),
                  icon: const Icon(Icons.share),
                  label: const Text('Partager'),
                  style: OutlinedButton.styleFrom(
                    foregroundColor: Colors.white,
                    side: BorderSide(color: Colors.white.withOpacity(0.2)),
                    padding: const EdgeInsets.symmetric(vertical: 12),
                  ),
                ),
              ),
            ],
          ),
        ],
      ),
    );
  }

  void _copyLink(BuildContext context) {
    final link = payment['payment_link'] ?? '';
    Clipboard.setData(ClipboardData(text: link));
    ScaffoldMessenger.of(context).showSnackBar(
      const SnackBar(content: Text('Lien copié!')),
    );
  }

  void _share(BuildContext context) {
    // In real app, use share_plus package
    _copyLink(context);
  }
}
