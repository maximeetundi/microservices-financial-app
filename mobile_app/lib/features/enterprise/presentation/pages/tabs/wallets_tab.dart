import 'package:flutter/material.dart';
import '../../../../../core/services/api_service.dart';
import '../../../data/models/enterprise_model.dart';
import '../send_money_page.dart';

class WalletsTab extends StatefulWidget {
  final Enterprise enterprise;
  final VoidCallback onRefresh;

  const WalletsTab({Key? key, required this.enterprise, required this.onRefresh}) : super(key: key);

  @override
  State<WalletsTab> createState() => _WalletsTabState();
}

class _WalletsTabState extends State<WalletsTab> {
  final ApiService _api = ApiService();
  bool _isLoading = true;
  List<Map<String, dynamic>> _wallets = [];
  double _totalBalance = 0;

  @override
  void initState() {
    super.initState();
    _loadWallets();
  }

  Future<void> _loadWallets() async {
    setState(() => _isLoading = true);
    try {
      final response = await _api.enterprise.getEnterpriseWallets();
      List<dynamic> allWallets = [];
      
      if (response is List) {
        allWallets = response;
      } else if (response is Map && response['wallets'] != null) {
        allWallets = response['wallets'];
      }
      
      // Filter to enterprise wallets
      final walletIds = [
        if (widget.enterprise.defaultWalletId != null) widget.enterprise.defaultWalletId,
        ...widget.enterprise.walletIds,
      ];
      
      _wallets = allWallets
          .where((w) => walletIds.contains(w['id']))
          .map((w) => Map<String, dynamic>.from(w))
          .toList();
      
      // If no wallet IDs, show all fiat wallets
      if (_wallets.isEmpty) {
        _wallets = allWallets
            .where((w) => w['wallet_type'] == 'fiat')
            .map((w) => Map<String, dynamic>.from(w))
            .toList();
      }
      
      _totalBalance = _wallets.fold(0.0, (sum, w) => sum + (w['balance'] ?? 0).toDouble());
    } catch (e) {
      ScaffoldMessenger.of(context).showSnackBar(
        SnackBar(content: Text('Erreur: ${e.toString()}')),
      );
    } finally {
      setState(() => _isLoading = false);
    }
  }

  void _sendMoney(Map<String, dynamic> wallet) {
    Navigator.push(
      context,
      MaterialPageRoute(
        builder: (context) => SendMoneyPage(
          enterprise: widget.enterprise,
          wallet: wallet,
        ),
      ),
    ).then((result) {
      if (result == true) _loadWallets();
    });
  }

  String _formatAmount(dynamic amount, String? currency) {
    final value = (amount ?? 0).toDouble();
    final curr = currency ?? 'XOF';
    return '${value.toStringAsFixed(0)} $curr';
  }

  @override
  Widget build(BuildContext context) {
    return RefreshIndicator(
      onRefresh: _loadWallets,
      child: SingleChildScrollView(
        physics: const AlwaysScrollableScrollPhysics(),
        padding: const EdgeInsets.all(16),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            // Total Balance Card
            Container(
              width: double.infinity,
              padding: const EdgeInsets.all(20),
              decoration: BoxDecoration(
                gradient: LinearGradient(
                  colors: [Colors.green.shade700, Colors.green.shade500],
                  begin: Alignment.topLeft,
                  end: Alignment.bottomRight,
                ),
                borderRadius: BorderRadius.circular(16),
                boxShadow: [
                  BoxShadow(
                    color: Colors.green.withOpacity(0.3),
                    blurRadius: 12,
                    offset: const Offset(0, 6),
                  ),
                ],
              ),
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  const Text(
                    'Solde total',
                    style: TextStyle(color: Colors.white70, fontSize: 14),
                  ),
                  const SizedBox(height: 8),
                  Text(
                    _formatAmount(_totalBalance, 'XOF'),
                    style: const TextStyle(
                      color: Colors.white,
                      fontWeight: FontWeight.bold,
                      fontSize: 28,
                    ),
                  ),
                  const SizedBox(height: 8),
                  Text(
                    '${_wallets.length} portefeuille(s)',
                    style: const TextStyle(color: Colors.white60, fontSize: 12),
                  ),
                ],
              ),
            ),
            
            const SizedBox(height: 24),
            
            // Header
            Row(
              mainAxisAlignment: MainAxisAlignment.spaceBetween,
              children: [
                const Text(
                  'Portefeuilles',
                  style: TextStyle(fontSize: 18, fontWeight: FontWeight.bold),
                ),
                TextButton.icon(
                  onPressed: () {},
                  icon: const Icon(Icons.add, size: 18),
                  label: const Text('Nouveau'),
                ),
              ],
            ),
            
            const SizedBox(height: 12),
            
            // Wallets List
            if (_isLoading)
              const Center(child: CircularProgressIndicator())
            else if (_wallets.isEmpty)
              _buildEmptyState()
            else
              ..._wallets.map((wallet) => _WalletCard(
                wallet: wallet,
                isDefault: wallet['id'] == widget.enterprise.defaultWalletId,
                onSend: () => _sendMoney(wallet),
                onViewTransactions: () => _showTransactions(wallet),
              )).toList(),
          ],
        ),
      ),
    );
  }

  Widget _buildEmptyState() {
    return Container(
      padding: const EdgeInsets.all(32),
      child: Column(
        children: [
          Icon(Icons.account_balance_wallet_outlined, size: 64, color: Colors.grey[300]),
          const SizedBox(height: 16),
          Text('Aucun portefeuille', style: TextStyle(color: Colors.grey[600])),
        ],
      ),
    );
  }

  void _showTransactions(Map<String, dynamic> wallet) {
    showModalBottomSheet(
      context: context,
      isScrollControlled: true,
      backgroundColor: Colors.transparent,
      builder: (context) => Container(
        height: MediaQuery.of(context).size.height * 0.7,
        decoration: const BoxDecoration(
          color: Colors.white,
          borderRadius: BorderRadius.only(
            topLeft: Radius.circular(20),
            topRight: Radius.circular(20),
          ),
        ),
        child: Column(
          children: [
            Container(
              margin: const EdgeInsets.only(top: 12),
              width: 40,
              height: 4,
              decoration: BoxDecoration(
                color: Colors.grey[300],
                borderRadius: BorderRadius.circular(2),
              ),
            ),
            Padding(
              padding: const EdgeInsets.all(20),
              child: Row(
                children: [
                  const Text(
                    'Transactions',
                    style: TextStyle(fontSize: 20, fontWeight: FontWeight.bold),
                  ),
                  const Spacer(),
                  IconButton(
                    onPressed: () => Navigator.pop(context),
                    icon: const Icon(Icons.close),
                  ),
                ],
              ),
            ),
            const Expanded(
              child: Center(
                child: Text('Chargement des transactions...'),
              ),
            ),
          ],
        ),
      ),
    );
  }
}

class _WalletCard extends StatelessWidget {
  final Map<String, dynamic> wallet;
  final bool isDefault;
  final VoidCallback onSend;
  final VoidCallback onViewTransactions;

  const _WalletCard({
    required this.wallet,
    required this.isDefault,
    required this.onSend,
    required this.onViewTransactions,
  });

  String get _currencySymbol {
    final currency = wallet['currency'] ?? 'XOF';
    const symbols = {'XOF': 'F', 'EUR': '€', 'USD': '\$', 'XAF': 'F'};
    return symbols[currency] ?? currency[0];
  }

  @override
  Widget build(BuildContext context) {
    return Container(
      margin: const EdgeInsets.only(bottom: 12),
      decoration: BoxDecoration(
        color: Colors.white,
        borderRadius: BorderRadius.circular(16),
        boxShadow: [
          BoxShadow(
            color: Colors.grey.withOpacity(0.1),
            blurRadius: 10,
            offset: const Offset(0, 4),
          ),
        ],
      ),
      child: Column(
        children: [
          Padding(
            padding: const EdgeInsets.all(16),
            child: Row(
              children: [
                // Currency Icon
                Container(
                  width: 50,
                  height: 50,
                  decoration: BoxDecoration(
                    gradient: LinearGradient(
                      colors: [Colors.blue.shade600, Colors.blue.shade800],
                    ),
                    borderRadius: BorderRadius.circular(12),
                  ),
                  child: Center(
                    child: Text(
                      _currencySymbol,
                      style: const TextStyle(
                        color: Colors.white,
                        fontSize: 20,
                        fontWeight: FontWeight.bold,
                      ),
                    ),
                  ),
                ),
                const SizedBox(width: 16),
                
                // Info
                Expanded(
                  child: Column(
                    crossAxisAlignment: CrossAxisAlignment.start,
                    children: [
                      Row(
                        children: [
                          Text(
                            wallet['currency'] ?? 'XOF',
                            style: const TextStyle(fontWeight: FontWeight.w600),
                          ),
                          if (isDefault) ...[
                            const SizedBox(width: 8),
                            Container(
                              padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 2),
                              decoration: BoxDecoration(
                                color: Colors.blue.shade50,
                                borderRadius: BorderRadius.circular(10),
                              ),
                              child: Text(
                                'Défaut',
                                style: TextStyle(color: Colors.blue.shade700, fontSize: 10),
                              ),
                            ),
                          ],
                        ],
                      ),
                      const SizedBox(height: 4),
                      Text(
                        '${(wallet['balance'] ?? 0).toDouble().toStringAsFixed(0)} ${wallet['currency'] ?? 'XOF'}',
                        style: const TextStyle(
                          fontSize: 18,
                          fontWeight: FontWeight.bold,
                        ),
                      ),
                    ],
                  ),
                ),
              ],
            ),
          ),
          
          // Actions
          Container(
            decoration: BoxDecoration(
              border: Border(top: BorderSide(color: Colors.grey.shade100)),
            ),
            child: Row(
              children: [
                Expanded(
                  child: TextButton.icon(
                    onPressed: onViewTransactions,
                    icon: const Icon(Icons.receipt_long, size: 18),
                    label: const Text('Transactions'),
                  ),
                ),
                Container(width: 1, height: 30, color: Colors.grey.shade100),
                Expanded(
                  child: TextButton.icon(
                    onPressed: onSend,
                    icon: const Icon(Icons.send, size: 18, color: Colors.green),
                    label: const Text('Envoyer', style: TextStyle(color: Colors.green)),
                  ),
                ),
              ],
            ),
          ),
        ],
      ),
    );
  }
}
