import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';
import 'package:go_router/go_router.dart';
import 'package:google_fonts/google_fonts.dart';
import 'package:intl/intl.dart';

import '../../../../core/theme/app_theme.dart';
import '../../../../core/widgets/glass_container.dart';
import '../../../wallet/presentation/bloc/wallet_bloc.dart';
import '../../../wallet/domain/entities/transaction.dart';
import '../../../../core/services/transfer_api_service.dart';

/// Transactions Page matching web design with filters

/// Transactions Page matching web design with filters
class TransactionsPage extends StatefulWidget {
  const TransactionsPage({super.key});

  @override
  State<TransactionsPage> createState() => _TransactionsPageState();
}

class _TransactionsPageState extends State<TransactionsPage> {
  final TransferApiService _transferApi = TransferApiService();
  String _filterType = '';
  String _filterPeriod = 'all';
  
  final List<Map<String, String>> _typeFilters = [
    {'value': '', 'label': 'Tous types'},
    {'value': 'deposit', 'label': 'Dépôts'},
    {'value': 'withdraw', 'label': 'Retraits'},
    {'value': 'transfer', 'label': 'Transferts'},
    {'value': 'exchange', 'label': 'Échanges'},
  ];
  
  final List<Map<String, String>> _periodFilters = [
    {'value': 'all', 'label': 'Toujours'},
    {'value': 'today', 'label': "Aujourd'hui"},
    {'value': 'week', 'label': '7 jours'},
    {'value': 'month', 'label': '30 jours'},
  ];

  @override
  void initState() {
    super.initState();
    // Load wallets which will also load transactions
    context.read<WalletBloc>().add(LoadWalletsEvent());
  }

  List<Transaction> _getFilteredTransactions(List<Transaction> transactions) {
    var result = List<Transaction>.from(transactions);
    
    // Filter by type
    if (_filterType.isNotEmpty) {
      result = result.where((tx) => tx.type.name.toLowerCase() == _filterType).toList();
    }
    
    // Filter by period
    if (_filterPeriod != 'all') {
      final now = DateTime.now();
      DateTime startDate;
      
      switch (_filterPeriod) {
        case 'today':
          startDate = DateTime(now.year, now.month, now.day);
          break;
        case 'week':
          startDate = now.subtract(const Duration(days: 7));
          break;
        case 'month':
          startDate = DateTime(now.year, now.month - 1, now.day);
          break;
        default:
          startDate = DateTime(2000);
      }
      
      result = result.where((tx) => tx.createdAt.isAfter(startDate)).toList();
    }
    
    // Sort by date descending
    result.sort((a, b) => b.createdAt.compareTo(a.createdAt));
    
    return result;
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
                child: BlocBuilder<WalletBloc, WalletState>(
                  builder: (context, state) {
                    if (state is WalletLoadingState) {
                      return const Center(child: CircularProgressIndicator());
                    }
                    
                    if (state is WalletErrorState) {
                      return _buildErrorState(state.message, isDark);
                    }
                    
                    if (state is WalletLoadedState) {
                      final transactions = _getFilteredTransactions(state.recentTransactions);
                      return _buildContent(transactions, isDark);
                    }
                    
                    return const SizedBox();
                  },
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
          const SizedBox(width: 16),
          Expanded(
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Text(
                  '📊 Transactions',
                  style: GoogleFonts.inter(
                    fontSize: 20,
                    fontWeight: FontWeight.bold,
                    color: isDark ? Colors.white : const Color(0xFF1E293B),
                  ),
                ),
                Text(
                  'Historique de vos transactions',
                  style: GoogleFonts.inter(
                    fontSize: 12,
                    color: isDark ? const Color(0xFF94A3B8) : const Color(0xFF64748B),
                  ),
                ),
              ],
            ),
          ),
        ],
      ),
    );
  }

  Widget _buildContent(List<Transaction> transactions, bool isDark) {
    return Column(
      children: [
        // Filters
        Padding(
          padding: const EdgeInsets.symmetric(horizontal: 16),
          child: Row(
            children: [
              Expanded(child: _buildFilterDropdown(_filterType, _typeFilters, (v) {
                setState(() => _filterType = v);
              }, isDark)),
              const SizedBox(width: 12),
              Expanded(child: _buildFilterDropdown(_filterPeriod, _periodFilters, (v) {
                setState(() => _filterPeriod = v);
              }, isDark)),
            ],
          ),
        ),
        
        const SizedBox(height: 16),
        
        // Transactions list
        Expanded(
          child: transactions.isEmpty
              ? _buildEmptyState(isDark)
              : ListView.builder(
                  padding: const EdgeInsets.symmetric(horizontal: 16),
                  itemCount: transactions.length,
                  itemBuilder: (context, index) => 
                      _buildTransactionItem(transactions[index], isDark),
                ),
        ),
      ],
    );
  }

  Widget _buildFilterDropdown(
    String value,
    List<Map<String, String>> options,
    Function(String) onChanged,
    bool isDark,
  ) {
    return Container(
      padding: const EdgeInsets.symmetric(horizontal: 12),
      decoration: BoxDecoration(
        color: isDark 
            ? Colors.white.withOpacity(0.05)
            : Colors.white,
        borderRadius: BorderRadius.circular(12),
        border: Border.all(
          color: isDark 
              ? Colors.white.withOpacity(0.1)
              : const Color(0xFFE2E8F0),
        ),
      ),
      child: DropdownButtonHideUnderline(
        child: DropdownButton<String>(
          value: value,
          isExpanded: true,
          dropdownColor: isDark ? const Color(0xFF1E293B) : Colors.white,
          items: options.map((opt) => DropdownMenuItem(
            value: opt['value'],
            child: Text(
              opt['label']!,
              style: GoogleFonts.inter(
                fontSize: 14,
                color: isDark ? Colors.white : const Color(0xFF1E293B),
              ),
            ),
          )).toList(),
          onChanged: (v) => onChanged(v ?? ''),
        ),
      ),
    );
  }

  Widget _buildTransactionItem(Transaction tx, bool isDark) {
    final isIncoming = tx.isIncoming;
    
    return Padding(
      padding: const EdgeInsets.only(bottom: 12),
      child: GestureDetector(
        onTap: () => _showTransactionDetails(tx),
        child: GlassContainer(
        padding: const EdgeInsets.all(16),
        borderRadius: 16,
        child: Row(
          children: [
            // Icon
            Container(
              width: 48,
              height: 48,
              decoration: BoxDecoration(
                color: _getTypeColor(tx.type.name, isIncoming).withOpacity(0.15),
                borderRadius: BorderRadius.circular(14),
              ),
              child: Center(
                child: Text(
                  _getTypeIcon(tx.type.name, isIncoming),
                  style: const TextStyle(fontSize: 20),
                ),
              ),
            ),
            const SizedBox(width: 16),
            
            // Info
            Expanded(
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  Text(
                    _getTransactionTitle(tx),
                    style: GoogleFonts.inter(
                      fontSize: 15,
                      fontWeight: FontWeight.w600,
                      color: isDark ? Colors.white : const Color(0xFF1E293B),
                    ),
                    maxLines: 1,
                    overflow: TextOverflow.ellipsis,
                  ),
                  const SizedBox(height: 4),
                  Text(
                    tx.memo ?? tx.currency,
                    style: GoogleFonts.inter(
                      fontSize: 13,
                      color: isDark ? const Color(0xFF64748B) : const Color(0xFF94A3B8),
                    ),
                    maxLines: 1,
                    overflow: TextOverflow.ellipsis,
                  ),
                ],
              ),
            ),
            
            // Amount & Date
            Column(
              crossAxisAlignment: CrossAxisAlignment.end,
              children: [
                Text(
                  '${isIncoming ? '+' : '-'}${_formatMoney(tx.amount, tx.currency)}',
                  style: GoogleFonts.inter(
                    fontSize: 15,
                    fontWeight: FontWeight.bold,
                    color: isIncoming 
                        ? const Color(0xFF22C55E)
                        : const Color(0xFFEF4444),
                  ),
                ),
                const SizedBox(height: 4),
                Text(
                  _formatDate(tx.createdAt),
                  style: GoogleFonts.inter(
                    fontSize: 11,
                    color: isDark ? const Color(0xFF64748B) : const Color(0xFF94A3B8),
                  ),
                ),
              ],
            ),
          ],
        ),
      ),
    );
  }

  void _showTransactionDetails(Transaction tx) {
    showModalBottomSheet(
      context: context,
      backgroundColor: Colors.transparent,
      isScrollControlled: true,
      builder: (context) => Container(
        padding: const EdgeInsets.all(24),
        decoration: BoxDecoration(
          color: Theme.of(context).scaffoldBackgroundColor,
          borderRadius: const BorderRadius.vertical(top: Radius.circular(24)),
        ),
        child: Column(
          mainAxisSize: MainAxisSize.min,
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Center(child: Container(width: 40, height: 4, decoration: BoxDecoration(color: Colors.grey.withOpacity(0.3), borderRadius: BorderRadius.circular(2)))),
            const SizedBox(height: 24),
            Text('Détails de la transaction', style: GoogleFonts.inter(fontSize: 18, fontWeight: FontWeight.bold)),
            const SizedBox(height: 24),
            _buildDetailRow('Type', _getTransactionTitle(tx)),
            _buildDetailRow('Montant', '${tx.isIncoming ? '+' : ''}${_formatMoney(tx.amount, tx.currency)}'),
            _buildDetailRow('Date', _formatDate(tx.createdAt)),
            _buildDetailRow('Statut', tx.status),
            if (tx.metadata != null && tx.metadata!.isNotEmpty) ...[
                const SizedBox(height: 16),
                const Text('MÉTADONNÉES', style: TextStyle(fontSize: 12, fontWeight: FontWeight.bold, color: Colors.grey)),
                const SizedBox(height: 8),
                Text(tx.metadata.toString(), style: const TextStyle(fontSize: 12)),
            ],

            const SizedBox(height: 32),
            
            // Actions for Transfers
            if (tx.type.name.toLowerCase() == 'transfer') ...[
                // Cancel (Sender)
                if (!tx.isIncoming && (tx.status == 'pending' || tx.status == 'processing')) 
                    SizedBox(
                        width: double.infinity,
                        child: ElevatedButton.icon(
                            onPressed: () {
                                Navigator.pop(context);
                                _confirmAction(tx, 'cancel');
                            },
                            icon: const Icon(Icons.cancel),
                            label: const Text('Annuler le transfert'),
                            style: ElevatedButton.styleFrom(backgroundColor: Colors.red),
                        ),
                    ),
                
                // Reverse (Beneficiary)
                if (tx.isIncoming && tx.status == 'completed')
                    SizedBox(
                        width: double.infinity,
                        child: ElevatedButton.icon(
                            onPressed: () {
                                Navigator.pop(context);
                                _confirmAction(tx, 'reverse');
                            },
                            icon: const Icon(Icons.undo),
                            label: const Text('Rembourser (Retourner les fonds)'),
                            style: ElevatedButton.styleFrom(backgroundColor: Colors.orange),
                        ),
                    ),
            ],
            const SizedBox(height: 24),
          ],
        ),
      ),
    );
  }

  Future<void> _confirmAction(Transaction tx, String action) async {
      final isCancel = action == 'cancel';
      final reasonController = TextEditingController();

      final result = await showDialog<String>(
          context: context,
          builder: (ctx) => AlertDialog(
              backgroundColor: const Color(0xFF1e293b),
              title: Text(isCancel ? 'Annuler le transfert ?' : 'Retourner les fonds ?', style: const TextStyle(color: Colors.white)),
              content: Column(
                  mainAxisSize: MainAxisSize.min,
                  children: [
                      Text(
                          isCancel 
                             ? 'Voulez-vous vraiment annuler ce transfert en cours ?'
                             : 'Voulez-vous retourner les fonds à l\'expéditeur ? Cette action est irréversible.',
                          style: const TextStyle(color: Colors.white70),
                      ),
                      const SizedBox(height: 16),
                      TextField(
                          controller: reasonController,
                          style: const TextStyle(color: Colors.white),
                          decoration: InputDecoration(
                              filled: true,
                              fillColor: Colors.black26,
                              hintText: 'Motif (Optionnel)',
                              hintStyle: TextStyle(color: Colors.white.withOpacity(0.3)),
                              border: OutlineInputBorder(borderRadius: BorderRadius.circular(12), borderSide: BorderSide.none),
                          ),
                      )
                  ],
              ),
              actions: [
                  TextButton(onPressed: () => Navigator.pop(ctx), child: const Text('Retour', style: TextStyle(color: Colors.grey))),
                  TextButton(
                      onPressed: () => Navigator.pop(ctx, reasonController.text), 
                      child: Text('Confirmer', style: TextStyle(color: isCancel ? Colors.red : Colors.orange)),
                  ),
              ],
          ),
      );

      if (result != null) {
          try {
              if (isCancel) {
                  await _transferApi.cancelTransfer(tx.id, reason: result);
              } else {
                  await _transferApi.reverseTransfer(tx.id, reason: result);
              }
              if (mounted) {
                  ScaffoldMessenger.of(context).showSnackBar(SnackBar(content: Text(isCancel ? 'Transfert annulé' : 'Fonds retournés')));
                  context.read<WalletBloc>().add(LoadWalletsEvent()); // Reload
              }
          } catch (e) {
              if (mounted) ScaffoldMessenger.of(context).showSnackBar(SnackBar(content: Text('Erreur: $e')));
          }
      }
  }

  Widget _buildDetailRow(String label, String value) {
      return Padding(
          padding: const EdgeInsets.only(bottom: 12),
          child: Row(
              mainAxisAlignment: MainAxisAlignment.spaceBetween,
              children: [
                  Text(label, style: const TextStyle(color: Colors.grey)),
                  Text(value, style: const TextStyle(fontWeight: FontWeight.bold)),
              ],
          ),
      );
  }
    return Center(
      child: Column(
        mainAxisAlignment: MainAxisAlignment.center,
        children: [
          const Text('📜', style: TextStyle(fontSize: 64)),
          const SizedBox(height: 16),
          Text(
            'Aucune transaction',
            style: GoogleFonts.inter(
              fontSize: 16,
              color: isDark ? const Color(0xFF64748B) : const Color(0xFF94A3B8),
            ),
          ),
        ],
      ),
    );
  }

  Widget _buildErrorState(String message, bool isDark) {
    return Center(
      child: Column(
        mainAxisAlignment: MainAxisAlignment.center,
        children: [
          const Icon(Icons.error_outline, size: 48, color: Color(0xFFEF4444)),
          const SizedBox(height: 16),
          Text(
            'Erreur: $message',
            style: GoogleFonts.inter(
              color: const Color(0xFFEF4444),
            ),
          ),
          const SizedBox(height: 16),
          ElevatedButton(
            onPressed: () => context.read<WalletBloc>().add(LoadWalletsEvent()),
            child: const Text('Réessayer'),
          ),
        ],
      ),
    );
  }

  String _getTypeIcon(String type, bool isIncoming) {
    switch (type.toLowerCase()) {
      case 'deposit':
        return '↓';
      case 'withdraw':
        return '↑';
      case 'transfer':
        return '💸';
      case 'exchange':
        return '💱';
      case 'payment':
        return '💳';
      default:
        return isIncoming ? '↓' : '↑';
    }
  }

  Color _getTypeColor(String type, bool isIncoming) {
    switch (type.toLowerCase()) {
      case 'deposit':
        return const Color(0xFF22C55E);
      case 'withdraw':
        return const Color(0xFFEF4444);
      case 'transfer':
        return const Color(0xFFA855F7);
      case 'exchange':
        return const Color(0xFF3B82F6);
      case 'payment':
        return const Color(0xFFF97316);
      default:
        return isIncoming 
            ? const Color(0xFF22C55E)
            : const Color(0xFFEF4444);
    }
  }

  String _getTransactionTitle(Transaction tx) {
    switch (tx.type.name.toLowerCase()) {
      case 'deposit':
        return 'Dépôt';
      case 'withdraw':
        return 'Retrait';
      case 'transfer':
        return tx.isIncoming ? 'Reçu' : 'Envoyé';
      case 'exchange':
        return 'Échange';
      case 'payment':
        return 'Paiement';
      default:
        return 'Transaction';
    }
  }

  String _formatMoney(double amount, String currency) {
    final formatter = NumberFormat.currency(locale: 'fr_FR', symbol: currency, decimalDigits: 2);
    return formatter.format(amount.abs());
  }

  String _formatDate(DateTime date) {
    return DateFormat('dd MMM', 'fr_FR').format(date);
  }
}
