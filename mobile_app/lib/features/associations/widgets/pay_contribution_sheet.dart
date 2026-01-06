import 'package:flutter/material.dart';
import 'package:intl/intl.dart';
import '../../../core/services/association_api_service.dart';
import '../../../core/services/wallet_api_service.dart';
import '../../../core/api/api_client.dart';

class PayContributionSheet extends StatefulWidget {
  final String associationId;
  final String currency;
  final VoidCallback onSuccess;

  const PayContributionSheet({
    super.key,
    required this.associationId,
    required this.currency,
    required this.onSuccess,
  });

  @override
  State<PayContributionSheet> createState() => _PayContributionSheetState();
}

class _PayContributionSheetState extends State<PayContributionSheet> {
  final AssociationApiService _associationApi = AssociationApiService(ApiClient().dio);
  final WalletApiService _walletApi = WalletApiService(); // No args needed
  
  List<dynamic> _wallets = [];
  String? _selectedWalletId;
  double _amount = 0;
  String _period = 'janvier_2026';
  String _pin = '';
  bool _loading = false;
  bool _loadingWallets = true;
  String? _error;

  @override
  void initState() {
    super.initState();
    _loadWallets();
  }

  Future<void> _loadWallets() async {
    try {
      // getWallets() returns List<Map<String, dynamic>> directly
      final allWallets = await _walletApi.getWallets();
      setState(() {
        _wallets = allWallets.where((w) => w['currency'] == widget.currency).toList();
        if (_wallets.isEmpty) _wallets = allWallets;
        if (_wallets.isNotEmpty) _selectedWalletId = _wallets[0]['id'];
      });
    } catch (e) {
      debugPrint('Failed to load wallets: $e');
    } finally {
      setState(() => _loadingWallets = false);
    }
  }

  String _formatCurrency(dynamic amount) {
    final value = (amount is num) ? amount.toDouble() : 0.0;
    return NumberFormat.currency(locale: 'fr_FR', symbol: widget.currency, decimalDigits: 0).format(value);
  }

  Future<void> _submit() async {
    if (_selectedWalletId == null || _amount <= 0 || _pin.length != 5) {
      setState(() => _error = 'Veuillez remplir tous les champs');
      return;
    }

    setState(() {
      _loading = true;
      _error = null;
    });

    try {
      await _associationApi.recordContribution(widget.associationId, {
        'wallet_id': _selectedWalletId,
        'pin': _pin,
        'amount': _amount,
        'period': _period,
        'description': 'Cotisation $_period',
      });

      if (mounted) {
        Navigator.pop(context);
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(
            content: Text('Cotisation de ${_formatCurrency(_amount)} payée avec succès!'),
            backgroundColor: const Color(0xFF10b981),
          ),
        );
        widget.onSuccess();
      }
    } catch (e) {
      setState(() => _error = 'Échec du paiement. Vérifiez votre solde et votre PIN.');
    } finally {
      setState(() => _loading = false);
    }
  }

  @override
  Widget build(BuildContext context) {
    return Container(
      padding: const EdgeInsets.all(24),
      decoration: const BoxDecoration(
        color: Color(0xFF1a1a2e),
        borderRadius: BorderRadius.vertical(top: Radius.circular(24)),
      ),
      child: Column(
        mainAxisSize: MainAxisSize.min,
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Row(
            mainAxisAlignment: MainAxisAlignment.spaceBetween,
            children: [
              const Text('Payer une cotisation', style: TextStyle(color: Colors.white, fontSize: 20, fontWeight: FontWeight.bold)),
              GestureDetector(
                onTap: () => Navigator.pop(context),
                child: Icon(Icons.close, color: Colors.white.withOpacity(0.6)),
              ),
            ],
          ),
          const SizedBox(height: 24),

          // Amount
          Text('Montant', style: TextStyle(color: Colors.white.withOpacity(0.7), fontSize: 14)),
          const SizedBox(height: 8),
          TextField(
            keyboardType: TextInputType.number,
            style: const TextStyle(color: Colors.white, fontSize: 18),
            decoration: InputDecoration(
              filled: true,
              fillColor: Colors.white.withOpacity(0.1),
              border: OutlineInputBorder(borderRadius: BorderRadius.circular(12), borderSide: BorderSide.none),
              suffixText: widget.currency,
              suffixStyle: TextStyle(color: Colors.white.withOpacity(0.5)),
            ),
            onChanged: (v) => _amount = double.tryParse(v) ?? 0,
          ),
          const SizedBox(height: 16),

          // Period
          Text('Période', style: TextStyle(color: Colors.white.withOpacity(0.7), fontSize: 14)),
          const SizedBox(height: 8),
          Container(
            padding: const EdgeInsets.symmetric(horizontal: 16),
            decoration: BoxDecoration(
              color: Colors.white.withOpacity(0.1),
              borderRadius: BorderRadius.circular(12),
            ),
            child: DropdownButton<String>(
              value: _period,
              isExpanded: true,
              dropdownColor: const Color(0xFF1a1a2e),
              underline: const SizedBox(),
              style: const TextStyle(color: Colors.white),
              items: const [
                DropdownMenuItem(value: 'janvier_2026', child: Text('Janvier 2026')),
                DropdownMenuItem(value: 'fevrier_2026', child: Text('Février 2026')),
                DropdownMenuItem(value: 'mars_2026', child: Text('Mars 2026')),
              ],
              onChanged: (v) => setState(() => _period = v!),
            ),
          ),
          const SizedBox(height: 16),

          // Wallet Selection
          Text('Portefeuille', style: TextStyle(color: Colors.white.withOpacity(0.7), fontSize: 14)),
          const SizedBox(height: 8),
          if (_loadingWallets)
            const Center(child: CircularProgressIndicator(color: Color(0xFF6366f1)))
          else
            SizedBox(
              height: 80,
              child: ListView.builder(
                scrollDirection: Axis.horizontal,
                itemCount: _wallets.length,
                itemBuilder: (context, index) {
                  final wallet = _wallets[index];
                  final isSelected = wallet['id'] == _selectedWalletId;
                  return GestureDetector(
                    onTap: () => setState(() => _selectedWalletId = wallet['id']),
                    child: Container(
                      width: 140,
                      margin: const EdgeInsets.only(right: 12),
                      padding: const EdgeInsets.all(12),
                      decoration: BoxDecoration(
                        color: isSelected ? const Color(0xFF6366f1).withOpacity(0.3) : Colors.white.withOpacity(0.05),
                        borderRadius: BorderRadius.circular(12),
                        border: Border.all(color: isSelected ? const Color(0xFF6366f1) : Colors.white.withOpacity(0.1)),
                      ),
                      child: Column(
                        crossAxisAlignment: CrossAxisAlignment.start,
                        mainAxisAlignment: MainAxisAlignment.center,
                        children: [
                          Text(wallet['currency'] ?? 'N/A', style: const TextStyle(color: Colors.white, fontWeight: FontWeight.bold)),
                          const SizedBox(height: 4),
                          Text(_formatCurrency(wallet['balance']), style: TextStyle(color: Colors.white.withOpacity(0.7), fontSize: 12)),
                        ],
                      ),
                    ),
                  );
                },
              ),
            ),
          const SizedBox(height: 16),

          // PIN
          Text('Code PIN', style: TextStyle(color: Colors.white.withOpacity(0.7), fontSize: 14)),
          const SizedBox(height: 8),
          TextField(
            obscureText: true,
            maxLength: 5,
            keyboardType: TextInputType.number,
            textAlign: TextAlign.center,
            style: const TextStyle(color: Colors.white, fontSize: 24, letterSpacing: 16),
            decoration: InputDecoration(
              filled: true,
              fillColor: Colors.white.withOpacity(0.1),
              border: OutlineInputBorder(borderRadius: BorderRadius.circular(12), borderSide: BorderSide.none),
              counterText: '',
            ),
            onChanged: (v) => _pin = v,
          ),
          const SizedBox(height: 16),

          // Error
          if (_error != null)
            Container(
              padding: const EdgeInsets.all(12),
              margin: const EdgeInsets.only(bottom: 16),
              decoration: BoxDecoration(
                color: Colors.red.withOpacity(0.2),
                borderRadius: BorderRadius.circular(8),
              ),
              child: Text(_error!, style: const TextStyle(color: Colors.red)),
            ),

          // Submit Button
          SizedBox(
            width: double.infinity,
            child: ElevatedButton(
              onPressed: _loading ? null : _submit,
              style: ElevatedButton.styleFrom(
                backgroundColor: const Color(0xFF6366f1),
                disabledBackgroundColor: Colors.grey,
                padding: const EdgeInsets.symmetric(vertical: 16),
                shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(12)),
              ),
              child: _loading
                  ? const SizedBox(width: 20, height: 20, child: CircularProgressIndicator(strokeWidth: 2, color: Colors.white))
                  : Text('Payer ${_formatCurrency(_amount)}', style: const TextStyle(fontSize: 16, fontWeight: FontWeight.bold)),
            ),
          ),
          const SizedBox(height: 16),
        ],
      ),
    );
  }
}

// Helper to show the sheet
void showPayContributionSheet(BuildContext context, String associationId, String currency, VoidCallback onSuccess) {
  showModalBottomSheet(
    context: context,
    isScrollControlled: true,
    backgroundColor: Colors.transparent,
    builder: (_) => Padding(
      padding: EdgeInsets.only(bottom: MediaQuery.of(context).viewInsets.bottom),
      child: PayContributionSheet(
        associationId: associationId,
        currency: currency,
        onSuccess: onSuccess,
      ),
    ),
  );
}
