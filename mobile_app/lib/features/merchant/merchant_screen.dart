import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import 'package:go_router/go_router.dart';
import 'package:flutter_bloc/flutter_bloc.dart';
import 'package:qr_flutter/qr_flutter.dart';

import '../wallet/presentation/bloc/wallet_bloc.dart';

/// Merchant Screen - Create payment requests and generate QR codes
class MerchantScreen extends StatefulWidget {
  const MerchantScreen({Key? key}) : super(key: key);

  @override
  State<MerchantScreen> createState() => _MerchantScreenState();
}

class _MerchantScreenState extends State<MerchantScreen> with SingleTickerProviderStateMixin {
  late TabController _tabController;
  List<Map<String, dynamic>> _payments = [];
  bool _loading = false;

  @override
  void initState() {
    super.initState();
    _tabController = TabController(length: 2, vsync: this);
    context.read<WalletBloc>().add(LoadWalletsEvent());
    _loadDemoPayments();
  }

  void _loadDemoPayments() {
    // Demo payments for testing
    _payments = [
      {
        'id': 'pay_001',
        'title': 'iPhone 15 Pro',
        'amount': 599000,
        'currency': 'XOF',
        'status': 'pending',
        'type': 'fixed',
        'created_at': DateTime.now().subtract(const Duration(hours: 2)),
      },
      {
        'id': 'pay_002',
        'title': 'Donation',
        'amount': null,
        'currency': 'XOF',
        'status': 'pending',
        'type': 'variable',
        'created_at': DateTime.now().subtract(const Duration(hours: 5)),
      },
      {
        'id': 'pay_003',
        'title': 'Commande #1234',
        'amount': 25000,
        'currency': 'XOF',
        'status': 'paid',
        'type': 'fixed',
        'created_at': DateTime.now().subtract(const Duration(days: 1)),
      },
    ];
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
          'Espace Marchand 💼',
          style: TextStyle(
            color: Color(0xFF1a1a2e),
            fontWeight: FontWeight.bold,
            fontSize: 20,
          ),
        ),
        centerTitle: true,
        bottom: TabBar(
          controller: _tabController,
          labelColor: const Color(0xFF667eea),
          unselectedLabelColor: const Color(0xFF64748B),
          indicatorColor: const Color(0xFF667eea),
          indicatorWeight: 3,
          tabs: const [
            Tab(text: 'En attente'),
            Tab(text: 'Historique'),
          ],
        ),
      ),
      body: Column(
        children: [
          // Stats Card
          _buildStatsCard(),
          
          // Payments List
          Expanded(
            child: TabBarView(
              controller: _tabController,
              children: [
                _buildPaymentsList(_payments.where((p) => p['status'] == 'pending').toList()),
                _buildPaymentsList(_payments.where((p) => p['status'] != 'pending').toList()),
              ],
            ),
          ),
        ],
      ),
      floatingActionButton: Row(
        mainAxisAlignment: MainAxisAlignment.end,
        children: [
          // Scan button
          FloatingActionButton(
            heroTag: 'scan',
            onPressed: () => context.push('/more/merchant/scan'),
            backgroundColor: const Color(0xFF10B981),
            child: const Icon(Icons.qr_code_scanner, color: Colors.white),
          ),
          const SizedBox(width: 12),
          // Create payment button
          FloatingActionButton.extended(
            heroTag: 'create',
            onPressed: () => _showCreatePaymentSheet(context),
            backgroundColor: const Color(0xFF667eea),
            icon: const Icon(Icons.add, color: Colors.white),
            label: const Text('Nouveau', style: TextStyle(color: Colors.white, fontWeight: FontWeight.bold)),
          ),
        ],
      ),
    );
  }

  Widget _buildStatsCard() {
    final pendingCount = _payments.where((p) => p['status'] == 'pending').length;
    final paidToday = _payments.where((p) => p['status'] == 'paid').fold<double>(0, (sum, p) => sum + (p['amount'] ?? 0));

    return Container(
      margin: const EdgeInsets.all(16),
      padding: const EdgeInsets.all(20),
      decoration: BoxDecoration(
        gradient: const LinearGradient(
          colors: [Color(0xFF667eea), Color(0xFF764ba2)],
        ),
        borderRadius: BorderRadius.circular(20),
        boxShadow: [
          BoxShadow(
            color: const Color(0xFF667eea).withOpacity(0.3),
            blurRadius: 15,
            offset: const Offset(0, 8),
          ),
        ],
      ),
      child: Row(
        children: [
          Expanded(
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                const Text(
                  'En attente',
                  style: TextStyle(color: Colors.white70, fontSize: 14),
                ),
                const SizedBox(height: 4),
                Text(
                  '$pendingCount paiements',
                  style: const TextStyle(
                    color: Colors.white,
                    fontSize: 24,
                    fontWeight: FontWeight.bold,
                  ),
                ),
              ],
            ),
          ),
          Container(
            width: 1,
            height: 50,
            color: Colors.white24,
          ),
          Expanded(
            child: Padding(
              padding: const EdgeInsets.only(left: 20),
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  const Text(
                    'Reçu aujourd\'hui',
                    style: TextStyle(color: Colors.white70, fontSize: 14),
                  ),
                  const SizedBox(height: 4),
                  Text(
                    '${paidToday.toStringAsFixed(0)} XOF',
                    style: const TextStyle(
                      color: Colors.white,
                      fontSize: 24,
                      fontWeight: FontWeight.bold,
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

  Widget _buildPaymentsList(List<Map<String, dynamic>> payments) {
    if (payments.isEmpty) {
      return Center(
        child: Column(
          mainAxisAlignment: MainAxisAlignment.center,
          children: [
            Container(
              padding: const EdgeInsets.all(24),
              decoration: BoxDecoration(
                color: const Color(0xFFE2E8F0),
                borderRadius: BorderRadius.circular(20),
              ),
              child: const Icon(Icons.receipt_long, size: 48, color: Color(0xFF94A3B8)),
            ),
            const SizedBox(height: 16),
            const Text(
              'Aucun paiement',
              style: TextStyle(
                color: Color(0xFF64748B),
                fontSize: 18,
                fontWeight: FontWeight.w500,
              ),
            ),
            const SizedBox(height: 8),
            const Text(
              'Créez une demande de paiement',
              style: TextStyle(color: Color(0xFF94A3B8)),
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
      backgroundColor: Colors.transparent,
      builder: (context) => _CreatePaymentSheet(
        onCreated: (payment) {
          setState(() => _payments.insert(0, payment));
          Navigator.pop(context);
          _showQRCode(payment);
        },
      ),
    );
  }

  void _showQRCode(Map<String, dynamic> payment) {
    showModalBottomSheet(
      context: context,
      isScrollControlled: true,
      backgroundColor: Colors.transparent,
      builder: (context) => _QRCodeSheet(payment: payment),
    );
  }
}

class _PaymentCard extends StatelessWidget {
  final Map<String, dynamic> payment;
  final VoidCallback onTap;

  const _PaymentCard({required this.payment, required this.onTap});

  @override
  Widget build(BuildContext context) {
    final status = payment['status'] ?? 'pending';
    final amount = payment['amount'];
    final currency = payment['currency'] ?? 'XOF';

    return GestureDetector(
      onTap: onTap,
      child: Container(
        margin: const EdgeInsets.only(bottom: 12),
        padding: const EdgeInsets.all(16),
        decoration: BoxDecoration(
          color: Colors.white,
          borderRadius: BorderRadius.circular(16),
          boxShadow: [
            BoxShadow(
              color: Colors.black.withOpacity(0.05),
              blurRadius: 10,
            ),
          ],
        ),
        child: Row(
          children: [
            Container(
              width: 50,
              height: 50,
              decoration: BoxDecoration(
                color: _getStatusColor(status).withOpacity(0.1),
                borderRadius: BorderRadius.circular(14),
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
                      color: Color(0xFF1a1a2e),
                      fontWeight: FontWeight.w600,
                      fontSize: 16,
                    ),
                  ),
                  const SizedBox(height: 4),
                  Text(
                    _getTypeLabel(payment['type']),
                    style: const TextStyle(
                      color: Color(0xFF64748B),
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
                  amount != null ? '${amount.toStringAsFixed(0)} $currency' : 'Variable',
                  style: const TextStyle(
                    color: Color(0xFF10B981),
                    fontWeight: FontWeight.bold,
                    fontSize: 16,
                  ),
                ),
                const SizedBox(height: 4),
                Container(
                  padding: const EdgeInsets.symmetric(horizontal: 10, vertical: 4),
                  decoration: BoxDecoration(
                    color: _getStatusColor(status).withOpacity(0.1),
                    borderRadius: BorderRadius.circular(12),
                  ),
                  child: Text(
                    _getStatusLabel(status),
                    style: TextStyle(
                      color: _getStatusColor(status),
                      fontSize: 11,
                      fontWeight: FontWeight.w600,
                    ),
                  ),
                ),
              ],
            ),
            const SizedBox(width: 8),
            const Icon(Icons.qr_code, color: Color(0xFF94A3B8)),
          ],
        ),
      ),
    );
  }

  Color _getStatusColor(String status) {
    switch (status) {
      case 'paid':
        return const Color(0xFF10B981);
      case 'pending':
        return const Color(0xFFF59E0B);
      case 'expired':
        return const Color(0xFFEF4444);
      default:
        return const Color(0xFF64748B);
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
  final Function(Map<String, dynamic>) onCreated;

  const _CreatePaymentSheet({required this.onCreated});

  @override
  State<_CreatePaymentSheet> createState() => _CreatePaymentSheetState();
}

class _CreatePaymentSheetState extends State<_CreatePaymentSheet> {
  String _type = 'fixed';
  String? _walletId;
  final _amountController = TextEditingController();
  final _titleController = TextEditingController();
  bool _loading = false;

  @override
  Widget build(BuildContext context) {
    return Container(
      decoration: const BoxDecoration(
        color: Colors.white,
        borderRadius: BorderRadius.vertical(top: Radius.circular(24)),
      ),
      padding: EdgeInsets.only(
        bottom: MediaQuery.of(context).viewInsets.bottom,
      ),
      child: SingleChildScrollView(
        padding: const EdgeInsets.all(24),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          mainAxisSize: MainAxisSize.min,
          children: [
            // Handle
            Center(
              child: Container(
                width: 40,
                height: 4,
                decoration: BoxDecoration(
                  color: const Color(0xFFE2E8F0),
                  borderRadius: BorderRadius.circular(2),
                ),
              ),
            ),
            const SizedBox(height: 24),
            const Text(
              'Nouvelle demande de paiement',
              style: TextStyle(
                fontSize: 22,
                fontWeight: FontWeight.bold,
                color: Color(0xFF1a1a2e),
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

            // Title
            const Text('Titre', style: TextStyle(color: Color(0xFF64748B), fontWeight: FontWeight.w500)),
            const SizedBox(height: 8),
            TextField(
              controller: _titleController,
              decoration: InputDecoration(
                hintText: 'ex: iPhone 15',
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
              ),
            ),
            const SizedBox(height: 16),

            // Amount (for fixed)
            if (_type == 'fixed') ...[
              const Text('Montant', style: TextStyle(color: Color(0xFF64748B), fontWeight: FontWeight.w500)),
              const SizedBox(height: 8),
              TextField(
                controller: _amountController,
                keyboardType: TextInputType.number,
                decoration: InputDecoration(
                  hintText: '0',
                  suffix: const Text('XOF'),
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
                ),
              ),
              const SizedBox(height: 16),
            ],

            const SizedBox(height: 8),

            // Submit
            GestureDetector(
              onTap: _loading ? null : _submit,
              child: Container(
                width: double.infinity,
                padding: const EdgeInsets.symmetric(vertical: 18),
                decoration: BoxDecoration(
                  gradient: const LinearGradient(
                    colors: [Color(0xFF667eea), Color(0xFF764ba2)],
                  ),
                  borderRadius: BorderRadius.circular(16),
                ),
                child: _loading
                    ? const Center(
                        child: SizedBox(
                          width: 24,
                          height: 24,
                          child: CircularProgressIndicator(color: Colors.white, strokeWidth: 2),
                        ),
                      )
                    : const Row(
                        mainAxisAlignment: MainAxisAlignment.center,
                        children: [
                          Icon(Icons.qr_code, color: Colors.white),
                          SizedBox(width: 8),
                          Text(
                            'Créer le QR code',
                            style: TextStyle(
                              color: Colors.white,
                              fontWeight: FontWeight.bold,
                              fontSize: 16,
                            ),
                          ),
                        ],
                      ),
              ),
            ),
          ],
        ),
      ),
    );
  }

  void _submit() {
    if (_titleController.text.isEmpty) return;
    if (_type == 'fixed' && _amountController.text.isEmpty) return;

    setState(() => _loading = true);

    // Simulate API call
    Future.delayed(const Duration(milliseconds: 500), () {
      final payment = {
        'id': 'pay_${DateTime.now().millisecondsSinceEpoch}',
        'title': _titleController.text,
        'amount': _type == 'fixed' ? double.tryParse(_amountController.text) : null,
        'currency': 'XOF',
        'status': 'pending',
        'type': _type,
        'created_at': DateTime.now(),
      };
      widget.onCreated(payment);
    });
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
    return GestureDetector(
      onTap: onTap,
      child: Container(
        padding: const EdgeInsets.all(16),
        decoration: BoxDecoration(
          color: selected ? const Color(0xFF667eea).withOpacity(0.1) : const Color(0xFFF8FAFC),
          borderRadius: BorderRadius.circular(14),
          border: Border.all(
            color: selected ? const Color(0xFF667eea) : const Color(0xFFE2E8F0),
            width: 2,
          ),
        ),
        child: Column(
          children: [
            Icon(
              icon,
              color: selected ? const Color(0xFF667eea) : const Color(0xFF64748B),
              size: 28,
            ),
            const SizedBox(height: 8),
            Text(
              label,
              style: TextStyle(
                color: selected ? const Color(0xFF667eea) : const Color(0xFF1a1a2e),
                fontWeight: FontWeight.w600,
              ),
            ),
          ],
        ),
      ),
    );
  }
}

class _QRCodeSheet extends StatelessWidget {
  final Map<String, dynamic> payment;

  const _QRCodeSheet({required this.payment});

  @override
  Widget build(BuildContext context) {
    final amount = payment['amount'];
    final currency = payment['currency'] ?? 'XOF';

    return Container(
      decoration: const BoxDecoration(
        color: Colors.white,
        borderRadius: BorderRadius.vertical(top: Radius.circular(24)),
      ),
      padding: const EdgeInsets.all(24),
      child: Column(
        mainAxisSize: MainAxisSize.min,
        children: [
          // Handle
          Container(
            width: 40,
            height: 4,
            decoration: BoxDecoration(
              color: const Color(0xFFE2E8F0),
              borderRadius: BorderRadius.circular(2),
            ),
          ),
          const SizedBox(height: 24),
          Text(
            payment['title'] ?? 'Paiement',
            style: const TextStyle(
              fontSize: 22,
              fontWeight: FontWeight.bold,
              color: Color(0xFF1a1a2e),
            ),
          ),
          const SizedBox(height: 24),

          // QR Code
          Container(
            padding: const EdgeInsets.all(16),
            decoration: BoxDecoration(
              color: Colors.white,
              borderRadius: BorderRadius.circular(20),
              border: Border.all(color: const Color(0xFFE2E8F0)),
            ),
            child: QrImageView(
              data: 'cryptobank://pay/${payment['id']}',
              version: QrVersions.auto,
              size: 200,
              backgroundColor: Colors.white,
              eyeStyle: const QrEyeStyle(
                eyeShape: QrEyeShape.square,
                color: Color(0xFF1a1a2e),
              ),
              dataModuleStyle: const QrDataModuleStyle(
                dataModuleShape: QrDataModuleShape.square,
                color: Color(0xFF1a1a2e),
              ),
            ),
          ),
          const SizedBox(height: 24),

          // Amount
          if (amount != null)
            Text(
              '${amount.toStringAsFixed(0)} $currency',
              style: const TextStyle(
                color: Color(0xFF10B981),
                fontSize: 36,
                fontWeight: FontWeight.bold,
              ),
            )
          else
            const Text(
              'Montant variable',
              style: TextStyle(color: Color(0xFF64748B), fontSize: 18),
            ),
          const SizedBox(height: 24),

          // Actions
          Row(
            children: [
              Expanded(
                child: GestureDetector(
                  onTap: () => _copyLink(context),
                  child: Container(
                    padding: const EdgeInsets.symmetric(vertical: 14),
                    decoration: BoxDecoration(
                      color: const Color(0xFFF8FAFC),
                      borderRadius: BorderRadius.circular(12),
                      border: Border.all(color: const Color(0xFFE2E8F0)),
                    ),
                    child: const Row(
                      mainAxisAlignment: MainAxisAlignment.center,
                      children: [
                        Icon(Icons.copy, color: Color(0xFF64748B)),
                        SizedBox(width: 8),
                        Text('Copier', style: TextStyle(color: Color(0xFF1a1a2e), fontWeight: FontWeight.w600)),
                      ],
                    ),
                  ),
                ),
              ),
              const SizedBox(width: 12),
              Expanded(
                child: GestureDetector(
                  onTap: () => _share(context),
                  child: Container(
                    padding: const EdgeInsets.symmetric(vertical: 14),
                    decoration: BoxDecoration(
                      gradient: const LinearGradient(
                        colors: [Color(0xFF667eea), Color(0xFF764ba2)],
                      ),
                      borderRadius: BorderRadius.circular(12),
                    ),
                    child: const Row(
                      mainAxisAlignment: MainAxisAlignment.center,
                      children: [
                        Icon(Icons.share, color: Colors.white),
                        SizedBox(width: 8),
                        Text('Partager', style: TextStyle(color: Colors.white, fontWeight: FontWeight.w600)),
                      ],
                    ),
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
    final link = 'https://cryptobank.app/pay/${payment['id']}';
    Clipboard.setData(ClipboardData(text: link));
    ScaffoldMessenger.of(context).showSnackBar(
      const SnackBar(
        content: Text('Lien copié!'),
        backgroundColor: Color(0xFF10B981),
      ),
    );
  }

  void _share(BuildContext context) {
    _copyLink(context);
    // In real app, use share_plus package
  }
}
