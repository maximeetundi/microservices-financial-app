import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import 'package:go_router/go_router.dart';
import 'package:google_fonts/google_fonts.dart';
import 'package:share_plus/share_plus.dart';
import 'dart:convert';

import '../../core/theme/app_theme.dart';
import '../../core/widgets/glass_container.dart';
import '../../core/services/api_service.dart';

/// Merchant Screen matching web design exactly - with real API integration
class MerchantScreen extends StatefulWidget {
  const MerchantScreen({super.key});

  @override
  State<MerchantScreen> createState() => _MerchantScreenState();
}

class _MerchantScreenState extends State<MerchantScreen> {
  final ApiService _api = ApiService();
  
  String _activeTab = 'pending';
  bool _showCreateModal = false;
  bool _creating = false;
  bool _isLoading = true;
  
  // Stats
  Map<String, dynamic> _stats = {
    'totalAmount': 0.0,
    'totalPayments': 0,
    'pending': 0,
  };
  
  // Payments and wallets
  List<Map<String, dynamic>> _payments = [];
  List<Map<String, dynamic>> _wallets = [];
  
  // Create form
  String _paymentType = 'fixed';
  String? _selectedWalletId;
  final _titleController = TextEditingController();
  final _amountController = TextEditingController();
  final _minAmountController = TextEditingController();
  final _maxAmountController = TextEditingController();
  final _descriptionController = TextEditingController();
  int _expiresInMinutes = 60;
  bool _reusable = false;

  List<Map<String, dynamic>> get _filteredPayments {
    if (_activeTab == 'pending') {
      return _payments.where((p) => p['status'] == 'pending').toList();
    }
    return _payments.where((p) => p['status'] != 'pending').toList();
  }

  @override
  void initState() {
    super.initState();
    _loadData();
  }

  @override
  void dispose() {
    _titleController.dispose();
    _amountController.dispose();
    _minAmountController.dispose();
    _maxAmountController.dispose();
    _descriptionController.dispose();
    super.dispose();
  }

  Future<void> _loadData() async {
    setState(() => _isLoading = true);
    try {
      final paymentsRes = await _api.merchant.getPayments();
      final walletsRes = await _api.wallet.getWallets();
      
      setState(() {
        _payments = List<Map<String, dynamic>>.from(
          paymentsRes['payments'] ?? paymentsRes['data']?['payments'] ?? []
        );
        _wallets = List<Map<String, dynamic>>.from(
          walletsRes['wallets'] ?? walletsRes['data']?['wallets'] ?? walletsRes['data'] ?? []
        );
        
        // Calculate stats
        double totalPaid = 0;
        int pending = 0;
        for (var p in _payments) {
          if (p['status'] == 'paid') {
            totalPaid += (p['amount'] ?? 0).toDouble();
          }
          if (p['status'] == 'pending') {
            pending++;
          }
        }
        _stats = {
          'totalAmount': totalPaid,
          'totalPayments': _payments.length,
          'pending': pending,
        };
        _isLoading = false;
      });
    } catch (e) {
      debugPrint('Failed to load merchant data: $e');
      setState(() => _isLoading = false);
    }
  }

  Future<void> _createPayment() async {
    if (_selectedWalletId == null || _titleController.text.isEmpty) {
      ScaffoldMessenger.of(context).showSnackBar(
        const SnackBar(content: Text('Veuillez remplir tous les champs requis')),
      );
      return;
    }

    setState(() => _creating = true);

    try {
      final wallet = _wallets.firstWhere((w) => w['id'] == _selectedWalletId);
      final currency = wallet['currency'] ?? 'EUR';

      final data = {
        'type': _paymentType,
        'wallet_id': _selectedWalletId,
        'title': _titleController.text,
        'currency': currency,
        'description': _descriptionController.text,
        'expires_in_minutes': _expiresInMinutes,
        'reusable': _reusable,
      };

      if (_paymentType == 'fixed') {
        data['amount'] = double.tryParse(_amountController.text) ?? 0;
      } else {
        if (_minAmountController.text.isNotEmpty) {
          data['min_amount'] = double.tryParse(_minAmountController.text);
        }
        if (_maxAmountController.text.isNotEmpty) {
          data['max_amount'] = double.tryParse(_maxAmountController.text);
        }
      }

      final response = await _api.merchant.createPayment(data);
      
      final newPayment = response['payment_request'] ?? response['data']?['payment_request'];
      final qrCodeBase64 = response['qr_code_base64'] ?? response['data']?['qr_code_base64'];

      if (newPayment != null) {
        setState(() {
          _payments.insert(0, newPayment);
          _showCreateModal = false;
        });

        // Show QR code modal
        _showPaymentQR(newPayment, qrCodeBase64);

        // Reset form
        _titleController.clear();
        _amountController.clear();
        _minAmountController.clear();
        _maxAmountController.clear();
        _descriptionController.clear();
        _paymentType = 'fixed';
        _expiresInMinutes = 60;
        _reusable = false;
      }

      await _loadData();
    } catch (e) {
      ScaffoldMessenger.of(context).showSnackBar(
        SnackBar(content: Text('Erreur: $e'), backgroundColor: Colors.red),
      );
    } finally {
      setState(() => _creating = false);
    }
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
        child: Stack(
          children: [
            SafeArea(
              child: Column(
                children: [
                  _buildAppBar(isDark),
                  Expanded(
                    child: _isLoading 
                        ? const Center(child: CircularProgressIndicator())
                        : RefreshIndicator(
                            onRefresh: _loadData,
                            child: SingleChildScrollView(
                              physics: const AlwaysScrollableScrollPhysics(),
                              padding: const EdgeInsets.all(20),
                              child: Column(
                                crossAxisAlignment: CrossAxisAlignment.start,
                                children: [
                                  _buildStatsCards(isDark),
                                  const SizedBox(height: 24),
                                  _buildPaymentsList(isDark),
                                ],
                              ),
                            ),
                          ),
                  ),
                ],
              ),
            ),
            
            // Create modal
            if (_showCreateModal) _buildCreateModal(isDark),
          ],
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
          const SizedBox(width: 16),
          Expanded(
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Text(
                  '💼 Espace Marchand',
                  style: GoogleFonts.inter(
                    fontSize: 20,
                    fontWeight: FontWeight.bold,
                    color: isDark ? Colors.white : const Color(0xFF1E293B),
                  ),
                ),
                Text(
                  'Recevez des paiements via QR code',
                  style: GoogleFonts.inter(
                    fontSize: 12,
                    color: isDark ? const Color(0xFF94A3B8) : const Color(0xFF64748B),
                  ),
                ),
              ],
            ),
          ),
          ElevatedButton.icon(
            onPressed: () => setState(() => _showCreateModal = true),
            icon: const Icon(Icons.add, size: 18),
            label: const Text('Créer'),
            style: ElevatedButton.styleFrom(
              backgroundColor: const Color(0xFF6366F1),
              foregroundColor: Colors.white,
              elevation: 0,
              padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 10),
              shape: RoundedRectangleBorder(
                borderRadius: BorderRadius.circular(12),
              ),
            ),
          ),
        ],
      ),
    );
  }

  Widget _buildStatsCards(bool isDark) {
    return Row(
      children: [
        Expanded(
          child: _buildStatCard(
            '💰',
            _formatCurrency(_stats['totalAmount']?.toDouble() ?? 0),
            'Total reçu',
            const Color(0xFF3B82F6),
            isDark,
          ),
        ),
        const SizedBox(width: 12),
        Expanded(
          child: _buildStatCard(
            '📊',
            '${_stats['totalPayments']}',
            'Paiements',
            const Color(0xFFA855F7),
            isDark,
          ),
        ),
        const SizedBox(width: 12),
        Expanded(
          child: _buildStatCard(
            '⏳',
            '${_stats['pending']}',
            'En attente',
            const Color(0xFFF59E0B),
            isDark,
          ),
        ),
      ],
    );
  }

  Widget _buildStatCard(String emoji, String value, String label, Color color, bool isDark) {
    return GlassContainer(
      padding: const EdgeInsets.all(16),
      borderRadius: 16,
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Container(
            width: 40,
            height: 40,
            decoration: BoxDecoration(
              color: color.withOpacity(isDark ? 0.2 : 0.1),
              borderRadius: BorderRadius.circular(10),
            ),
            child: Center(
              child: Text(emoji, style: const TextStyle(fontSize: 20)),
            ),
          ),
          const SizedBox(height: 12),
          Text(
            value,
            style: GoogleFonts.inter(
              fontSize: 16,
              fontWeight: FontWeight.bold,
              color: isDark ? Colors.white : const Color(0xFF1E293B),
            ),
            maxLines: 1,
            overflow: TextOverflow.ellipsis,
          ),
          const SizedBox(height: 2),
          Text(
            label,
            style: GoogleFonts.inter(
              fontSize: 12,
              color: isDark ? const Color(0xFF64748B) : const Color(0xFF94A3B8),
            ),
          ),
        ],
      ),
    );
  }

  Widget _buildPaymentsList(bool isDark) {
    return GlassContainer(
      padding: EdgeInsets.zero,
      borderRadius: 20,
      child: Column(
        children: [
          // Header with tabs
          Padding(
            padding: const EdgeInsets.all(16),
            child: Row(
              mainAxisAlignment: MainAxisAlignment.spaceBetween,
              children: [
                Text(
                  'Demandes de paiement',
                  style: GoogleFonts.inter(
                    fontSize: 18,
                    fontWeight: FontWeight.bold,
                    color: isDark ? Colors.white : const Color(0xFF1E293B),
                  ),
                ),
                Container(
                  padding: const EdgeInsets.all(4),
                  decoration: BoxDecoration(
                    color: isDark ? const Color(0xFF1E293B) : const Color(0xFFF1F5F9),
                    borderRadius: BorderRadius.circular(12),
                  ),
                  child: Row(
                    children: [
                      _buildTab('En attente', 'pending', isDark),
                      _buildTab('Terminés', 'completed', isDark),
                    ],
                  ),
                ),
              ],
            ),
          ),
          
          // Payments list
          if (_filteredPayments.isEmpty)
            Padding(
              padding: const EdgeInsets.all(40),
              child: Column(
                children: [
                  const Text('📭', style: TextStyle(fontSize: 48)),
                  const SizedBox(height: 16),
                  Text(
                    'Aucune demande de paiement',
                    style: GoogleFonts.inter(
                      fontSize: 16,
                      color: isDark ? const Color(0xFF64748B) : const Color(0xFF94A3B8),
                    ),
                  ),
                  const SizedBox(height: 16),
                  ElevatedButton(
                    onPressed: () => setState(() => _showCreateModal = true),
                    style: ElevatedButton.styleFrom(
                      backgroundColor: const Color(0xFF6366F1),
                      foregroundColor: Colors.white,
                    ),
                    child: const Text('Créer votre première demande'),
                  ),
                ],
              ),
            )
          else
            ...(_filteredPayments.map((payment) => _buildPaymentItem(payment, isDark))),
        ],
      ),
    );
  }

  Widget _buildTab(String label, String value, bool isDark) {
    final isActive = _activeTab == value;
    return GestureDetector(
      onTap: () => setState(() => _activeTab = value),
      child: Container(
        padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 8),
        decoration: BoxDecoration(
          color: isActive 
              ? (isDark ? const Color(0xFF334155) : Colors.white)
              : Colors.transparent,
          borderRadius: BorderRadius.circular(8),
          boxShadow: isActive ? [
            BoxShadow(
              color: Colors.black.withOpacity(0.05),
              blurRadius: 4,
            ),
          ] : null,
        ),
        child: Text(
          label,
          style: GoogleFonts.inter(
            fontSize: 13,
            fontWeight: FontWeight.w600,
            color: isActive 
                ? const Color(0xFF6366F1)
                : (isDark ? const Color(0xFF64748B) : const Color(0xFF94A3B8)),
          ),
        ),
      ),
    );
  }

  Widget _buildPaymentItem(Map<String, dynamic> payment, bool isDark) {
    final createdAt = payment['created_at'] != null 
        ? DateTime.tryParse(payment['created_at'].toString()) 
        : null;
    
    return InkWell(
      onTap: () => _showPaymentQR(payment, null),
      child: Container(
        padding: const EdgeInsets.all(16),
        decoration: BoxDecoration(
          border: Border(
            top: BorderSide(
              color: isDark ? const Color(0xFF1E293B) : const Color(0xFFE2E8F0),
            ),
          ),
        ),
        child: Row(
          children: [
            Expanded(
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  Text(
                    payment['title'] ?? 'Sans titre',
                    style: GoogleFonts.inter(
                      fontSize: 16,
                      fontWeight: FontWeight.w600,
                      color: isDark ? Colors.white : const Color(0xFF1E293B),
                    ),
                  ),
                  const SizedBox(height: 4),
                  Row(
                    children: [
                      Container(
                        padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 2),
                        decoration: BoxDecoration(
                          color: isDark ? const Color(0xFF1E293B) : const Color(0xFFF1F5F9),
                          borderRadius: BorderRadius.circular(8),
                        ),
                        child: Text(
                          payment['type'] == 'fixed' ? 'Prix fixe' : 'Variable',
                          style: GoogleFonts.inter(
                            fontSize: 11,
                            color: isDark ? const Color(0xFF94A3B8) : const Color(0xFF64748B),
                          ),
                        ),
                      ),
                      const SizedBox(width: 8),
                      Text(
                        createdAt != null ? _formatDate(createdAt) : '',
                        style: GoogleFonts.inter(
                          fontSize: 12,
                          color: isDark ? const Color(0xFF64748B) : const Color(0xFF94A3B8),
                        ),
                      ),
                    ],
                  ),
                ],
              ),
            ),
            Column(
              crossAxisAlignment: CrossAxisAlignment.end,
              children: [
                Text(
                  payment['amount'] != null 
                      ? _formatCurrency((payment['amount'] as num).toDouble(), payment['currency'] ?? 'EUR')
                      : 'Variable',
                  style: GoogleFonts.inter(
                    fontSize: 16,
                    fontWeight: FontWeight.bold,
                    color: isDark ? Colors.white : const Color(0xFF1E293B),
                    fontStyle: payment['amount'] == null ? FontStyle.italic : FontStyle.normal,
                  ),
                ),
                const SizedBox(height: 4),
                _buildStatusBadge(payment['status'] ?? 'pending', isDark),
              ],
            ),
            const SizedBox(width: 12),
            Icon(
              Icons.qr_code_rounded,
              color: isDark ? const Color(0xFF64748B) : const Color(0xFF94A3B8),
            ),
          ],
        ),
      ),
    );
  }

  Widget _buildStatusBadge(String status, bool isDark) {
    Color bgColor;
    Color textColor;
    String label;
    
    switch (status) {
      case 'pending':
        bgColor = const Color(0xFFFEF3C7);
        textColor = const Color(0xFFF59E0B);
        label = 'En attente';
        break;
      case 'paid':
        bgColor = const Color(0xFFD1FAE5);
        textColor = const Color(0xFF10B981);
        label = 'Payé';
        break;
      case 'expired':
        bgColor = const Color(0xFFFEE2E2);
        textColor = const Color(0xFFEF4444);
        label = 'Expiré';
        break;
      default:
        bgColor = const Color(0xFFF1F5F9);
        textColor = const Color(0xFF64748B);
        label = status;
    }
    
    if (isDark) {
      bgColor = textColor.withOpacity(0.2);
    }
    
    return Container(
      padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 4),
      decoration: BoxDecoration(
        color: bgColor,
        borderRadius: BorderRadius.circular(12),
      ),
      child: Text(
        label,
        style: GoogleFonts.inter(
          fontSize: 11,
          fontWeight: FontWeight.w600,
          color: textColor,
        ),
      ),
    );
  }

  void _showPaymentQR(Map<String, dynamic> payment, String? qrCodeBase64) async {
    // If no QR code provided, fetch it
    String? qrCode = qrCodeBase64;
    if (qrCode == null) {
      try {
        final response = await _api.merchant.getQRCode(payment['id'].toString());
        qrCode = response['qr_code_base64'] ?? response['data']?['qr_code_base64'];
      } catch (e) {
        debugPrint('Failed to get QR code: $e');
      }
    }

    showModalBottomSheet(
      context: context,
      isScrollControlled: true,
      backgroundColor: Colors.transparent,
      builder: (context) => _buildQRModal(payment, qrCode),
    );
  }

  Widget _buildQRModal(Map<String, dynamic> payment, String? qrCodeBase64) {
    final isDark = Theme.of(context).brightness == Brightness.dark;
    
    return Container(
      margin: const EdgeInsets.all(16),
      decoration: BoxDecoration(
        color: isDark ? const Color(0xFF1E293B) : Colors.white,
        borderRadius: BorderRadius.circular(24),
      ),
      child: Column(
        mainAxisSize: MainAxisSize.min,
        children: [
          // Header
          Padding(
            padding: const EdgeInsets.all(20),
            child: Row(
              mainAxisAlignment: MainAxisAlignment.spaceBetween,
              children: [
                Expanded(
                  child: Text(
                    payment['title'] ?? 'Paiement',
                    style: GoogleFonts.inter(
                      fontSize: 20,
                      fontWeight: FontWeight.bold,
                      color: isDark ? Colors.white : const Color(0xFF1E293B),
                    ),
                    maxLines: 1,
                    overflow: TextOverflow.ellipsis,
                  ),
                ),
                IconButton(
                  onPressed: () => Navigator.pop(context),
                  icon: const Icon(Icons.close),
                ),
              ],
            ),
          ),
          
          // QR Code
          Container(
            width: 200,
            height: 200,
            margin: const EdgeInsets.all(20),
            decoration: BoxDecoration(
              color: Colors.white,
              borderRadius: BorderRadius.circular(16),
              border: Border.all(color: const Color(0xFFE2E8F0)),
            ),
            child: qrCodeBase64 != null
                ? ClipRRect(
                    borderRadius: BorderRadius.circular(16),
                    child: Image.memory(
                      base64Decode(qrCodeBase64.replaceFirst('data:image/png;base64,', '')),
                      fit: BoxFit.contain,
                    ),
                  )
                : const Center(
                    child: Icon(Icons.qr_code_2, size: 120, color: Color(0xFF1E293B)),
                  ),
          ),
          
          // Info
          Padding(
            padding: const EdgeInsets.symmetric(horizontal: 20),
            child: Column(
              children: [
                Row(
                  mainAxisAlignment: MainAxisAlignment.spaceBetween,
                  children: [
                    Text('Montant', style: GoogleFonts.inter(color: const Color(0xFF64748B))),
                    Text(
                      payment['amount'] != null 
                          ? _formatCurrency((payment['amount'] as num).toDouble(), payment['currency'] ?? 'EUR')
                          : 'À définir',
                      style: GoogleFonts.inter(
                        fontWeight: FontWeight.bold,
                        color: isDark ? Colors.white : const Color(0xFF1E293B),
                      ),
                    ),
                  ],
                ),
                const SizedBox(height: 12),
                Row(
                  mainAxisAlignment: MainAxisAlignment.spaceBetween,
                  children: [
                    Text('Statut', style: GoogleFonts.inter(color: const Color(0xFF64748B))),
                    _buildStatusBadge(payment['status'] ?? 'pending', isDark),
                  ],
                ),
              ],
            ),
          ),

          const SizedBox(height: 16),
          
          // Payment Code Display
          Padding(
            padding: const EdgeInsets.symmetric(horizontal: 20),
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Text(
                  'CODE DE PAIEMENT',
                  style: GoogleFonts.inter(
                    fontSize: 11,
                    fontWeight: FontWeight.bold,
                    color: const Color(0xFF64748B),
                    letterSpacing: 1,
                  ),
                ),
                const SizedBox(height: 8),
                Container(
                  width: double.infinity,
                  padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 14),
                  decoration: BoxDecoration(
                    color: isDark 
                        ? const Color(0xFF6366F1).withOpacity(0.15)
                        : const Color(0xFFEEF2FF),
                    borderRadius: BorderRadius.circular(12),
                    border: Border.all(
                      color: isDark 
                          ? const Color(0xFF6366F1).withOpacity(0.3)
                          : const Color(0xFFC7D2FE),
                      width: 2,
                    ),
                  ),
                  child: Row(
                    children: [
                      Expanded(
                        child: Text(
                          _getPaymentCode(payment),
                          style: GoogleFonts.sourceCodePro(
                            fontSize: 16,
                            fontWeight: FontWeight.bold,
                            color: const Color(0xFF6366F1),
                          ),
                          textAlign: TextAlign.center,
                        ),
                      ),
                      GestureDetector(
                        onTap: () {
                          final code = _getPaymentCode(payment);
                          Clipboard.setData(ClipboardData(text: code));
                          ScaffoldMessenger.of(context).showSnackBar(
                            SnackBar(content: Text('Code copié: $code')),
                          );
                        },
                        child: const Icon(Icons.copy, color: Color(0xFF6366F1), size: 22),
                      ),
                    ],
                  ),
                ),
              ],
            ),
          ),
          
          // Actions
          Padding(
            padding: const EdgeInsets.all(20),
            child: Row(
              children: [
                Expanded(
                  child: _buildActionButton('📋', 'Code', () {
                    final code = _getPaymentCode(payment);
                    Clipboard.setData(ClipboardData(text: code));
                    ScaffoldMessenger.of(context).showSnackBar(
                      SnackBar(content: Text('Code copié: $code')),
                    );
                  }, isDark),
                ),
                const SizedBox(width: 12),
                Expanded(
                  child: _buildActionButton('🔗', 'Lien', () {
                    final link = payment['payment_link'] ?? 'https://app.zekora.com/pay/${_getPaymentCode(payment)}';
                    Clipboard.setData(ClipboardData(text: link));
                    ScaffoldMessenger.of(context).showSnackBar(
                      const SnackBar(content: Text('Lien copié!')),
                    );
                  }, isDark),
                ),
                const SizedBox(width: 12),
                Expanded(
                  child: _buildActionButton('📤', 'Partager', () {
                    final code = _getPaymentCode(payment);
                    final link = payment['payment_link'] ?? 'https://app.zekora.com/pay/$code';
                    Share.share('Paiement: ${payment['title']}\nCode: $code\n$link');
                  }, isDark),
                ),
              ],
            ),
          ),
        ],
      ),
    );
  }

  String _getPaymentCode(Map<String, dynamic> payment) {
    // Try payment_id first
    final paymentId = payment['payment_id'] ?? payment['id'];
    if (paymentId != null) {
      final idStr = paymentId.toString();
      if (idStr.startsWith('pay_')) {
        return idStr;
      }
      return 'pay_$idStr';
    }
    // Try to extract from payment_link
    final link = payment['payment_link']?.toString() ?? '';
    if (link.contains('/pay/')) {
      return link.split('/pay/').last;
    }
    return '';
  }

  Widget _buildActionButton(String emoji, String label, VoidCallback onTap, bool isDark) {
    return InkWell(
      onTap: onTap,
      borderRadius: BorderRadius.circular(12),
      child: Container(
        padding: const EdgeInsets.symmetric(vertical: 12),
        decoration: BoxDecoration(
          color: isDark ? const Color(0xFF0F172A) : const Color(0xFFF8FAFC),
          borderRadius: BorderRadius.circular(12),
        ),
        child: Column(
          children: [
            Text(emoji, style: const TextStyle(fontSize: 20)),
            const SizedBox(height: 4),
            Text(
              label,
              style: GoogleFonts.inter(
                fontSize: 12,
                fontWeight: FontWeight.w600,
                color: isDark ? const Color(0xFF94A3B8) : const Color(0xFF64748B),
              ),
            ),
          ],
        ),
      ),
    );
  }

  Widget _buildCreateModal(bool isDark) {
    return GestureDetector(
      onTap: () => setState(() => _showCreateModal = false),
      child: Container(
        color: Colors.black.withOpacity(0.6),
        child: Center(
          child: GestureDetector(
            onTap: () {},
            child: Container(
              margin: const EdgeInsets.all(20),
              constraints: const BoxConstraints(maxHeight: 600),
              decoration: BoxDecoration(
                color: isDark ? const Color(0xFF1E293B) : Colors.white,
                borderRadius: BorderRadius.circular(24),
              ),
              child: Column(
                mainAxisSize: MainAxisSize.min,
                children: [
                  // Header
                  Padding(
                    padding: const EdgeInsets.all(20),
                    child: Row(
                      mainAxisAlignment: MainAxisAlignment.spaceBetween,
                      children: [
                        Text(
                          'Nouvelle demande',
                          style: GoogleFonts.inter(
                            fontSize: 20,
                            fontWeight: FontWeight.bold,
                            color: isDark ? Colors.white : const Color(0xFF1E293B),
                          ),
                        ),
                        IconButton(
                          onPressed: () => setState(() => _showCreateModal = false),
                          icon: const Icon(Icons.close),
                        ),
                      ],
                    ),
                  ),
                  
                  // Form
                  Flexible(
                    child: SingleChildScrollView(
                      padding: const EdgeInsets.symmetric(horizontal: 20),
                      child: Column(
                        crossAxisAlignment: CrossAxisAlignment.start,
                        children: [
                          // Type selection
                          Text('Type de paiement', style: GoogleFonts.inter(
                            fontWeight: FontWeight.w600,
                            color: isDark ? Colors.white : const Color(0xFF374151),
                          )),
                          const SizedBox(height: 12),
                          Row(
                            children: [
                              Expanded(child: _buildTypeButton('🏷️', 'Prix fixe', 'fixed', isDark)),
                              const SizedBox(width: 12),
                              Expanded(child: _buildTypeButton('🎁', 'Variable', 'variable', isDark)),
                            ],
                          ),
                          
                          const SizedBox(height: 20),
                          
                          // Wallet selection
                          Text('Portefeuille', style: GoogleFonts.inter(
                            fontWeight: FontWeight.w600,
                            color: isDark ? Colors.white : const Color(0xFF374151),
                          )),
                          const SizedBox(height: 8),
                          Container(
                            padding: const EdgeInsets.symmetric(horizontal: 16),
                            decoration: BoxDecoration(
                              color: isDark ? const Color(0xFF0F172A) : const Color(0xFFF8FAFC),
                              borderRadius: BorderRadius.circular(12),
                              border: Border.all(
                                color: isDark ? const Color(0xFF334155) : const Color(0xFFE2E8F0),
                              ),
                            ),
                            child: DropdownButtonHideUnderline(
                              child: DropdownButton<String>(
                                value: _selectedWalletId,
                                isExpanded: true,
                                hint: Text('Sélectionner...', style: GoogleFonts.inter(
                                  color: const Color(0xFF94A3B8),
                                )),
                                dropdownColor: isDark ? const Color(0xFF1E293B) : Colors.white,
                                items: _wallets.map((w) => DropdownMenuItem(
                                  value: w['id'].toString(),
                                  child: Text('${w['currency']} - ${_formatCurrency((w['balance'] as num).toDouble(), w['currency'] ?? 'EUR')}'),
                                )).toList(),
                                onChanged: (v) => setState(() => _selectedWalletId = v),
                                style: GoogleFonts.inter(
                                  color: isDark ? Colors.white : const Color(0xFF1E293B),
                                ),
                              ),
                            ),
                          ),
                          
                          const SizedBox(height: 20),
                          
                          // Title
                          _buildTextField('Titre', _titleController, 'Ex: iPhone 15 Pro', isDark),
                          
                          const SizedBox(height: 20),
                          
                          // Amount (for fixed)
                          if (_paymentType == 'fixed')
                            _buildTextField('Montant', _amountController, '0.00', isDark, isNumber: true),
                          
                          // Min/Max (for variable)
                          if (_paymentType == 'variable')
                            Row(
                              children: [
                                Expanded(child: _buildTextField('Min', _minAmountController, '0.00', isDark, isNumber: true)),
                                const SizedBox(width: 12),
                                Expanded(child: _buildTextField('Max', _maxAmountController, '0.00', isDark, isNumber: true)),
                              ],
                            ),
                          
                          const SizedBox(height: 20),
                          
                          // Description
                          _buildTextField('Description (optionnel)', _descriptionController, 'Détails...', isDark, maxLines: 3),
                          
                          const SizedBox(height: 20),
                        ],
                      ),
                    ),
                  ),
                  
                  // Actions
                  Padding(
                    padding: const EdgeInsets.all(20),
                    child: Row(
                      children: [
                        Expanded(
                          child: OutlinedButton(
                            onPressed: () => setState(() => _showCreateModal = false),
                            style: OutlinedButton.styleFrom(
                              padding: const EdgeInsets.symmetric(vertical: 16),
                              shape: RoundedRectangleBorder(
                                borderRadius: BorderRadius.circular(12),
                              ),
                            ),
                            child: const Text('Annuler'),
                          ),
                        ),
                        const SizedBox(width: 12),
                        Expanded(
                          child: ElevatedButton(
                            onPressed: _creating ? null : _createPayment,
                            style: ElevatedButton.styleFrom(
                              backgroundColor: const Color(0xFF6366F1),
                              foregroundColor: Colors.white,
                              padding: const EdgeInsets.symmetric(vertical: 16),
                              shape: RoundedRectangleBorder(
                                borderRadius: BorderRadius.circular(12),
                              ),
                            ),
                            child: _creating
                                ? const SizedBox(
                                    width: 20,
                                    height: 20,
                                    child: CircularProgressIndicator(color: Colors.white, strokeWidth: 2),
                                  )
                                : const Text('Créer le QR code'),
                          ),
                        ),
                      ],
                    ),
                  ),
                ],
              ),
            ),
          ),
        ),
      ),
    );
  }

  Widget _buildTypeButton(String emoji, String label, String value, bool isDark) {
    final isActive = _paymentType == value;
    return GestureDetector(
      onTap: () => setState(() => _paymentType = value),
      child: Container(
        padding: const EdgeInsets.all(16),
        decoration: BoxDecoration(
          color: isActive 
              ? const Color(0xFF6366F1).withOpacity(0.1)
              : (isDark ? const Color(0xFF0F172A) : const Color(0xFFF8FAFC)),
          borderRadius: BorderRadius.circular(12),
          border: Border.all(
            color: isActive ? const Color(0xFF6366F1) : (isDark ? const Color(0xFF334155) : const Color(0xFFE2E8F0)),
            width: isActive ? 2 : 1,
          ),
        ),
        child: Column(
          children: [
            Text(emoji, style: const TextStyle(fontSize: 28)),
            const SizedBox(height: 8),
            Text(
              label,
              style: GoogleFonts.inter(
                fontWeight: FontWeight.w600,
                color: isActive ? const Color(0xFF6366F1) : (isDark ? Colors.white : const Color(0xFF374151)),
              ),
            ),
          ],
        ),
      ),
    );
  }

  Widget _buildTextField(String label, TextEditingController controller, String hint, bool isDark, {bool isNumber = false, int maxLines = 1}) {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Text(label, style: GoogleFonts.inter(
          fontWeight: FontWeight.w600,
          color: isDark ? Colors.white : const Color(0xFF374151),
        )),
        const SizedBox(height: 8),
        TextField(
          controller: controller,
          keyboardType: isNumber ? TextInputType.number : TextInputType.text,
          maxLines: maxLines,
          style: GoogleFonts.inter(
            color: isDark ? Colors.white : const Color(0xFF1E293B),
          ),
          decoration: InputDecoration(
            hintText: hint,
            hintStyle: GoogleFonts.inter(color: const Color(0xFF94A3B8)),
            filled: true,
            fillColor: isDark ? const Color(0xFF0F172A) : const Color(0xFFF8FAFC),
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
      ],
    );
  }

  String _formatCurrency(double amount, [String currency = 'EUR']) {
    if (currency == 'XOF' || currency == 'XAF') {
      return '${amount.toStringAsFixed(0).replaceAllMapped(RegExp(r'(\d{1,3})(?=(\d{3})+(?!\d))'), (m) => '${m[1]} ')} FCFA';
    }
    return '${amount.toStringAsFixed(2)} $currency';
  }

  String _formatDate(DateTime date) {
    final now = DateTime.now();
    final diff = now.difference(date);
    
    if (diff.inDays == 0) {
      return 'Il y a ${diff.inHours}h';
    } else if (diff.inDays == 1) {
      return 'Hier';
    } else {
      return 'Il y a ${diff.inDays}j';
    }
  }
}
