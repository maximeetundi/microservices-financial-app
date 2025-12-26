import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';
import 'package:google_fonts/google_fonts.dart';

import '../../../../core/theme/app_theme.dart';
import '../../../../core/widgets/custom_button.dart';
import '../../../../core/widgets/glass_container.dart';
import '../../domain/entities/wallet.dart';
import '../../domain/entities/transaction.dart' as tx;
import '../bloc/wallet_bloc.dart';

/// Quick action buttons for wallet operations - matching web design
class WalletActions extends StatelessWidget {
  final VoidCallback? onSendPressed;
  final VoidCallback? onReceivePressed;
  final VoidCallback? onBuyPressed;
  final VoidCallback? onExchangePressed;

  const WalletActions({
    Key? key,
    this.onSendPressed,
    this.onReceivePressed,
    this.onBuyPressed,
    this.onExchangePressed,
  }) : super(key: key);

  @override
  Widget build(BuildContext context) {
    final isDark = Theme.of(context).brightness == Brightness.dark;
    
    return Padding(
      padding: const EdgeInsets.symmetric(vertical: 16),
      child: Row(
        mainAxisAlignment: MainAxisAlignment.spaceEvenly,
        children: [
          _buildActionButton(
            context: context,
            emoji: '💸',
            label: 'Envoyer',
            color: AppTheme.primaryColor,
            onTap: onSendPressed,
            isDark: isDark,
          ),
          _buildActionButton(
            context: context,
            emoji: '📥',
            label: 'Recevoir',
            color: AppTheme.secondaryColor,
            onTap: onReceivePressed,
            isDark: isDark,
          ),
          _buildActionButton(
            context: context,
            emoji: '💳',
            label: 'Recharger',
            color: const Color(0xFFF59E0B),
            onTap: onBuyPressed,
            isDark: isDark,
          ),
          _buildActionButton(
            context: context,
            emoji: '💱',
            label: 'Exchange',
            color: const Color(0xFF3B82F6),
            onTap: onExchangePressed,
            isDark: isDark,
          ),
        ],
      ),
    );
  }

  Widget _buildActionButton({
    required BuildContext context,
    required String emoji,
    required String label,
    required Color color,
    required bool isDark,
    VoidCallback? onTap,
  }) {
    return GestureDetector(
      onTap: onTap,
      child: Column(
        children: [
          Container(
            width: 60,
            height: 60,
            decoration: BoxDecoration(
              color: isDark 
                  ? Colors.white.withOpacity(0.05)
                  : Colors.white.withOpacity(0.8),
              shape: BoxShape.circle,
              border: Border.all(
                color: isDark 
                    ? Colors.white.withOpacity(0.1)
                    : const Color(0xFFE2E8F0),
              ),
              boxShadow: [
                BoxShadow(
                  color: Colors.black.withOpacity(0.05),
                  blurRadius: 10,
                  offset: const Offset(0, 4),
                ),
              ],
            ),
            child: Center(
              child: Text(emoji, style: const TextStyle(fontSize: 24)),
            ),
          ),
          const SizedBox(height: 8),
          Text(
            label,
            style: GoogleFonts.inter(
              fontSize: 12,
              fontWeight: FontWeight.w500,
              color: isDark ? Colors.white70 : const Color(0xFF64748B),
            ),
          ),
        ],
      ),
    );
  }
}

/// Card displaying wallet information - matching web design exactly
class WalletCard extends StatelessWidget {
  final Wallet wallet;
  final VoidCallback? onTap;
  final VoidCallback? onRechargePressed;
  final VoidCallback? onSendPressed;

  const WalletCard({
    Key? key,
    required this.wallet,
    this.onTap,
    this.onRechargePressed,
    this.onSendPressed,
  }) : super(key: key);

  @override
  Widget build(BuildContext context) {
    final isDark = Theme.of(context).brightness == Brightness.dark;
    
    return GestureDetector(
      onTap: onTap,
      child: GlassContainer(
        padding: const EdgeInsets.all(20),
        borderRadius: 20,
        color: isDark 
            ? Colors.white.withOpacity(0.05)
            : Colors.white.withOpacity(0.9),
        child: Row(
          children: [
            // Currency icon with colored background
            Container(
              width: 56,
              height: 56,
              decoration: BoxDecoration(
                color: _getCurrencyColor(wallet.currency).withOpacity(0.2),
                borderRadius: BorderRadius.circular(16),
              ),
              child: Center(
                child: Text(
                  _getCurrencyEmoji(wallet.currency),
                  style: const TextStyle(fontSize: 28),
                ),
              ),
            ),
            const SizedBox(width: 16),
            
            // Wallet info
            Expanded(
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  Row(
                    children: [
                      Flexible(
                        child: Text(
                          wallet.name ?? '${_getCurrencyName(wallet.currency)} Wallet',
                          style: GoogleFonts.inter(
                            fontWeight: FontWeight.w600,
                            fontSize: 16,
                            color: isDark ? Colors.white : const Color(0xFF1E293B),
                          ),
                          overflow: TextOverflow.ellipsis,
                        ),
                      ),
                      const SizedBox(width: 8),
                      // Status badge
                      Container(
                        padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 2),
                        decoration: BoxDecoration(
                          color: wallet.status == WalletStatus.active 
                              ? const Color(0xFF10B981).withOpacity(0.2)
                              : const Color(0xFFEF4444).withOpacity(0.2),
                          borderRadius: BorderRadius.circular(12),
                        ),
                        child: Text(
                          wallet.status == WalletStatus.active ? 'Actif' : 'Inactif',
                          style: GoogleFonts.inter(
                            fontSize: 10,
                            fontWeight: FontWeight.w600,
                            color: wallet.status == WalletStatus.active 
                                ? const Color(0xFF10B981)
                                : const Color(0xFFEF4444),
                          ),
                        ),
                      ),
                    ],
                  ),
                  const SizedBox(height: 4),
                  Text(
                    _formatBalance(wallet.balance, wallet.currency),
                    style: GoogleFonts.inter(
                      color: isDark ? const Color(0xFF94A3B8) : const Color(0xFF64748B),
                      fontSize: 14,
                    ),
                  ),
                  // Wallet address (truncated)
                  if (wallet.address != null && wallet.address!.isNotEmpty)
                    Text(
                      '${wallet.address!.substring(0, 6)}...${wallet.address!.substring(wallet.address!.length - 4)}',
                      style: GoogleFonts.robotoMono(
                        color: isDark ? const Color(0xFF64748B) : const Color(0xFF94A3B8),
                        fontSize: 12,
                      ),
                    ),
                ],
              ),
            ),
            
            // USD value and actions
            Column(
              crossAxisAlignment: CrossAxisAlignment.end,
              children: [
                Text(
                  '≈ \$${(wallet.balance * wallet.usdRate).toStringAsFixed(2)}',
                  style: GoogleFonts.inter(
                    fontWeight: FontWeight.bold,
                    fontSize: 18,
                    color: isDark ? Colors.white : const Color(0xFF1E293B),
                  ),
                ),
                const SizedBox(height: 8),
                // Action buttons
                Row(
                  children: [
                    GestureDetector(
                      onTap: onRechargePressed,
                      child: _buildMiniActionButton(
                        context, 
                        '💳', 
                        'Recharger',
                        isDark,
                      ),
                    ),
                    const SizedBox(width: 8),
                    GestureDetector(
                      onTap: onSendPressed,
                      child: _buildMiniActionButton(
                        context, 
                        '💸', 
                        'Envoyer',
                        isDark,
                      ),
                    ),
                  ],
                ),
              ],
            ),
          ],
        ),
      ),
    );
  }

  Widget _buildMiniActionButton(BuildContext context, String emoji, String label, bool isDark) {
    return Container(
      padding: const EdgeInsets.symmetric(horizontal: 10, vertical: 6),
      decoration: BoxDecoration(
        color: isDark 
            ? const Color(0xFF334155).withOpacity(0.5)
            : const Color(0xFFF1F5F9),
        borderRadius: BorderRadius.circular(8),
        border: Border.all(
          color: isDark 
              ? Colors.white.withOpacity(0.1)
              : const Color(0xFFE2E8F0),
        ),
      ),
      child: Row(
        children: [
          Text(emoji, style: const TextStyle(fontSize: 12)),
          const SizedBox(width: 4),
          Text(
            label,
            style: GoogleFonts.inter(
              fontSize: 11,
              fontWeight: FontWeight.w500,
              color: isDark ? Colors.white70 : const Color(0xFF64748B),
            ),
          ),
        ],
      ),
    );
  }

  String _getCurrencyEmoji(String currency) {
    switch (currency.toUpperCase()) {
      case 'BTC': return '₿';
      case 'ETH': return 'Ξ';
      case 'SOL': return '◎';
      case 'USD': return '💵';
      case 'EUR': return '💶';
      case 'GBP': return '💷';
      case 'XOF': return '🇨🇮';
      case 'USDT': return '₮';
      case 'USDC': return '💲';
      default: return '💰';
    }
  }

  String _getCurrencyName(String currency) {
    switch (currency.toUpperCase()) {
      case 'BTC': return 'Bitcoin';
      case 'ETH': return 'Ethereum';
      case 'SOL': return 'Solana';
      case 'USD': return 'US Dollar';
      case 'EUR': return 'Euro';
      case 'GBP': return 'British Pound';
      case 'XOF': return 'Franc CFA';
      case 'USDT': return 'Tether';
      case 'USDC': return 'USD Coin';
      default: return currency;
    }
  }

  Color _getCurrencyColor(String currency) {
    switch (currency.toUpperCase()) {
      case 'BTC': return const Color(0xFFF7931A);
      case 'ETH': return const Color(0xFF627EEA);
      case 'SOL': return const Color(0xFF9945FF);
      case 'USD': return const Color(0xFF22C55E);
      case 'EUR': return const Color(0xFF3B82F6);
      case 'GBP': return const Color(0xFF4C1D95);
      case 'XOF': return const Color(0xFFFF6900);
      case 'USDT': return const Color(0xFF26A17B);
      case 'USDC': return const Color(0xFF2775CA);
      default: return AppTheme.primaryColor;
    }
  }

  String _formatBalance(double balance, String currency) {
    if (['BTC', 'ETH', 'SOL'].contains(currency.toUpperCase())) {
      return '${balance.toStringAsFixed(8)} $currency';
    }
    return '${balance.toStringAsFixed(2)} $currency';
  }
}

/// List of recent transactions - matching web design
class RecentTransactionsList extends StatelessWidget {
  final List<tx.Transaction> transactions;

  const RecentTransactionsList({
    Key? key,
    required this.transactions,
  }) : super(key: key);

  @override
  Widget build(BuildContext context) {
    final isDark = Theme.of(context).brightness == Brightness.dark;
    
    if (transactions.isEmpty) {
      return Container(
        padding: const EdgeInsets.all(32),
        child: Center(
          child: Column(
            children: [
              const Text('📊', style: TextStyle(fontSize: 40)),
              const SizedBox(height: 12),
              Text(
                'Aucune transaction',
                style: GoogleFonts.inter(
                  color: isDark ? Colors.white : const Color(0xFF1E293B),
                  fontSize: 16,
                  fontWeight: FontWeight.bold,
                ),
              ),
              const SizedBox(height: 4),
              Text(
                'Vos transactions apparaîtront ici',
                style: GoogleFonts.inter(
                  color: isDark ? const Color(0xFF64748B) : const Color(0xFF94A3B8),
                  fontSize: 14,
                ),
              ),
            ],
          ),
        ),
      );
    }

    return Column(
      children: transactions.map((transaction) {
        final isReceive = transaction.isIncoming;
        return Container(
          margin: const EdgeInsets.only(bottom: 12),
          padding: const EdgeInsets.all(16),
          decoration: BoxDecoration(
            color: isDark 
                ? const Color(0xFF1E293B).withOpacity(0.5)
                : const Color(0xFFF8FAFC),
            borderRadius: BorderRadius.circular(16),
            border: Border.all(
              color: isDark 
                  ? const Color(0xFF334155).withOpacity(0.5)
                  : const Color(0xFFE2E8F0),
            ),
          ),
          child: Row(
            children: [
              Container(
                width: 48,
                height: 48,
                decoration: BoxDecoration(
                  color: (isReceive ? const Color(0xFF10B981) : const Color(0xFFEF4444))
                      .withOpacity(0.2),
                  borderRadius: BorderRadius.circular(12),
                ),
                child: Center(
                  child: Text(
                    isReceive ? '📥' : '📤',
                    style: const TextStyle(fontSize: 20),
                  ),
                ),
              ),
              const SizedBox(width: 16),
              Expanded(
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    Text(
                      _getTransactionTitle(transaction),
                      style: GoogleFonts.inter(
                        fontWeight: FontWeight.w600,
                        fontSize: 15,
                        color: isDark ? Colors.white : const Color(0xFF1E293B),
                      ),
                    ),
                    Text(
                      _formatDate(transaction.createdAt),
                      style: GoogleFonts.inter(
                        color: isDark ? const Color(0xFF64748B) : const Color(0xFF94A3B8),
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
                    '${isReceive ? '+' : '-'}${transaction.amount.toStringAsFixed(4)} ${transaction.currency}',
                    style: GoogleFonts.inter(
                      fontWeight: FontWeight.bold,
                      fontSize: 15,
                      color: isReceive ? const Color(0xFF10B981) : const Color(0xFFEF4444),
                    ),
                  ),
                  Text(
                    '≈ \$${(transaction.amount * 1.0).toStringAsFixed(2)}', // Replace with actual USD rate
                    style: GoogleFonts.inter(
                      color: isDark ? const Color(0xFF64748B) : const Color(0xFF94A3B8),
                      fontSize: 12,
                    ),
                  ),
                ],
              ),
            ],
          ),
        );
      }).toList(),
    );
  }

  String _getTransactionTitle(tx.Transaction transaction) {
    switch (transaction.type) {
      case tx.TransactionType.receive:
        return 'Dépôt';
      case tx.TransactionType.send:
        return 'Envoyé';
      case tx.TransactionType.exchange:
        return 'Échange';
      case tx.TransactionType.buy:
        return 'Achat';
      case tx.TransactionType.sell:
        return 'Vente';
      case tx.TransactionType.fee:
        return 'Frais';
      case tx.TransactionType.reward:
        return 'Récompense';
      case tx.TransactionType.staking:
        return 'Staking';
    }
  }

  String _formatDate(DateTime date) {
    final now = DateTime.now();
    final diff = now.difference(date);
    if (diff.inDays == 0) {
      return 'Aujourd\'hui';
    } else if (diff.inDays == 1) {
      return 'Hier';
    } else if (diff.inDays < 7) {
      return 'Il y a ${diff.inDays} jours';
    }
    return '${date.day}/${date.month}/${date.year}';
  }
}

/// Bottom sheet for creating a new wallet - matching web design
class CreateWalletBottomSheet extends StatefulWidget {
  const CreateWalletBottomSheet({Key? key}) : super(key: key);

  @override
  State<CreateWalletBottomSheet> createState() => _CreateWalletBottomSheetState();
}

class _CreateWalletBottomSheetState extends State<CreateWalletBottomSheet> {
  String _selectedCurrency = 'BTC';
  String _selectedType = 'crypto';
  final _nameController = TextEditingController();

  final List<Map<String, dynamic>> _cryptoCurrencies = [
    {'code': 'BTC', 'name': 'Bitcoin', 'emoji': '₿', 'color': const Color(0xFFF7931A)},
    {'code': 'ETH', 'name': 'Ethereum', 'emoji': 'Ξ', 'color': const Color(0xFF627EEA)},
    {'code': 'SOL', 'name': 'Solana', 'emoji': '◎', 'color': const Color(0xFF9945FF)},
    {'code': 'USDT', 'name': 'Tether', 'emoji': '₮', 'color': const Color(0xFF26A17B)},
  ];

  final List<Map<String, dynamic>> _fiatCurrencies = [
    {'code': 'USD', 'name': 'US Dollar', 'emoji': '💵', 'color': const Color(0xFF22C55E)},
    {'code': 'EUR', 'name': 'Euro', 'emoji': '💶', 'color': const Color(0xFF3B82F6)},
    {'code': 'XOF', 'name': 'Franc CFA', 'emoji': '🇨🇮', 'color': const Color(0xFFFF6900)},
  ];

  @override
  Widget build(BuildContext context) {
    final isDark = Theme.of(context).brightness == Brightness.dark;
    
    return Container(
      padding: const EdgeInsets.all(24),
      decoration: BoxDecoration(
        color: isDark ? const Color(0xFF1E293B) : Colors.white,
        borderRadius: const BorderRadius.vertical(top: Radius.circular(24)),
      ),
      child: Column(
        mainAxisSize: MainAxisSize.min,
        crossAxisAlignment: CrossAxisAlignment.stretch,
        children: [
          // Handle
          Center(
            child: Container(
              width: 40,
              height: 4,
              decoration: BoxDecoration(
                color: isDark ? Colors.white24 : Colors.grey.shade300,
                borderRadius: BorderRadius.circular(2),
              ),
            ),
          ),
          const SizedBox(height: 24),
          
          // Title with emoji
          Text(
            '🆕 Nouveau Portefeuille',
            style: GoogleFonts.inter(
              fontSize: 24,
              fontWeight: FontWeight.bold,
              color: isDark ? Colors.white : const Color(0xFF1E293B),
            ),
            textAlign: TextAlign.center,
          ),
          const SizedBox(height: 24),
          
          // Wallet Type Toggle
          Row(
            children: [
              Expanded(
                child: _buildTypeButton(
                  'crypto',
                  '₿ Crypto',
                  isDark,
                ),
              ),
              const SizedBox(width: 12),
              Expanded(
                child: _buildTypeButton(
                  'fiat',
                  '💵 Fiat',
                  isDark,
                ),
              ),
            ],
          ),
          const SizedBox(height: 24),
          
          // Currency Selection
          Text(
            'Sélectionner la devise',
            style: GoogleFonts.inter(
              fontSize: 14,
              fontWeight: FontWeight.w600,
              color: isDark ? Colors.white70 : const Color(0xFF64748B),
            ),
          ),
          const SizedBox(height: 12),
          Wrap(
            spacing: 8,
            runSpacing: 8,
            children: (_selectedType == 'crypto' ? _cryptoCurrencies : _fiatCurrencies)
                .map((currency) {
              final isSelected = _selectedCurrency == currency['code'];
              return GestureDetector(
                onTap: () {
                  setState(() => _selectedCurrency = currency['code'] as String);
                },
                child: Container(
                  padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 12),
                  decoration: BoxDecoration(
                    color: isSelected 
                        ? (currency['color'] as Color).withOpacity(0.2)
                        : (isDark ? const Color(0xFF334155) : const Color(0xFFF1F5F9)),
                    borderRadius: BorderRadius.circular(12),
                    border: Border.all(
                      color: isSelected 
                          ? currency['color'] as Color
                          : (isDark ? Colors.white.withOpacity(0.1) : const Color(0xFFE2E8F0)),
                      width: isSelected ? 2 : 1,
                    ),
                  ),
                  child: Row(
                    mainAxisSize: MainAxisSize.min,
                    children: [
                      Text(
                        currency['emoji'] as String,
                        style: const TextStyle(fontSize: 18),
                      ),
                      const SizedBox(width: 8),
                      Text(
                        currency['code'] as String,
                        style: GoogleFonts.inter(
                          fontWeight: FontWeight.w600,
                          color: isSelected 
                              ? (currency['color'] as Color)
                              : (isDark ? Colors.white : const Color(0xFF1E293B)),
                        ),
                      ),
                    ],
                  ),
                ),
              );
            }).toList(),
          ),
          
          const SizedBox(height: 24),
          
          // Wallet Name
          Text(
            'Nom du portefeuille (optionnel)',
            style: GoogleFonts.inter(
              fontSize: 14,
              fontWeight: FontWeight.w600,
              color: isDark ? Colors.white70 : const Color(0xFF64748B),
            ),
          ),
          const SizedBox(height: 8),
          TextField(
            controller: _nameController,
            style: GoogleFonts.inter(
              color: isDark ? Colors.white : const Color(0xFF1E293B),
            ),
            decoration: InputDecoration(
              hintText: 'Mon portefeuille Bitcoin',
              hintStyle: GoogleFonts.inter(
                color: isDark ? const Color(0xFF64748B) : const Color(0xFF94A3B8),
              ),
              filled: true,
              fillColor: isDark 
                  ? const Color(0xFF334155)
                  : const Color(0xFFF1F5F9),
              border: OutlineInputBorder(
                borderRadius: BorderRadius.circular(12),
                borderSide: BorderSide.none,
              ),
              contentPadding: const EdgeInsets.symmetric(horizontal: 16, vertical: 16),
            ),
          ),
          
          const SizedBox(height: 32),
          
          // Create Button
          BlocConsumer<WalletBloc, WalletState>(
            listener: (context, state) {
              if (state is WalletCreatedState) {
                Navigator.pop(context);
                ScaffoldMessenger.of(context).showSnackBar(
                  const SnackBar(
                    content: Text('Portefeuille créé avec succès!'),
                    backgroundColor: Color(0xFF10B981),
                  ),
                );
              } else if (state is WalletErrorState) {
                ScaffoldMessenger.of(context).showSnackBar(
                  SnackBar(
                    content: Text(state.message),
                    backgroundColor: const Color(0xFFEF4444),
                  ),
                );
              }
            },
            builder: (context, state) {
              final isLoading = state is WalletLoadingState;
              return Container(
                decoration: BoxDecoration(
                  borderRadius: BorderRadius.circular(12),
                  gradient: const LinearGradient(
                    colors: [Color(0xFF6366F1), Color(0xFF8B5CF6)],
                  ),
                  boxShadow: [
                    BoxShadow(
                      color: const Color(0xFF6366F1).withOpacity(0.3),
                      blurRadius: 20,
                      offset: const Offset(0, 8),
                    ),
                  ],
                ),
                child: ElevatedButton(
                  style: ElevatedButton.styleFrom(
                    backgroundColor: Colors.transparent,
                    shadowColor: Colors.transparent,
                    padding: const EdgeInsets.symmetric(vertical: 16),
                    shape: RoundedRectangleBorder(
                      borderRadius: BorderRadius.circular(12),
                    ),
                  ),
                  onPressed: isLoading ? null : _createWallet,
                  child: isLoading 
                    ? const SizedBox(
                        height: 20,
                        width: 20,
                        child: CircularProgressIndicator(color: Colors.white, strokeWidth: 2),
                      )
                    : Text(
                        'Créer le portefeuille',
                        style: GoogleFonts.inter(
                          fontSize: 16,
                          fontWeight: FontWeight.bold,
                          color: Colors.white,
                        ),
                      ),
                ),
              );
            },
          ),
          const SizedBox(height: 16),
        ],
      ),
    );
  }

  Widget _buildTypeButton(String type, String label, bool isDark) {
    final isSelected = _selectedType == type;
    
    return GestureDetector(
      onTap: () {
        setState(() {
          _selectedType = type;
          _selectedCurrency = type == 'crypto' ? 'BTC' : 'USD';
        });
      },
      child: Container(
        padding: const EdgeInsets.symmetric(vertical: 16),
        decoration: BoxDecoration(
          color: isSelected 
              ? const Color(0xFF6366F1).withOpacity(0.2)
              : (isDark ? const Color(0xFF334155) : const Color(0xFFF1F5F9)),
          borderRadius: BorderRadius.circular(12),
          border: Border.all(
            color: isSelected 
                ? const Color(0xFF6366F1)
                : (isDark ? Colors.white.withOpacity(0.1) : const Color(0xFFE2E8F0)),
            width: isSelected ? 2 : 1,
          ),
        ),
        child: Text(
          label,
          textAlign: TextAlign.center,
          style: GoogleFonts.inter(
            fontWeight: FontWeight.w600,
            color: isSelected 
                ? const Color(0xFF6366F1)
                : (isDark ? Colors.white : const Color(0xFF1E293B)),
          ),
        ),
      ),
    );
  }

  void _createWallet() {
    context.read<WalletBloc>().add(CreateWalletEvent(
      currency: _selectedCurrency,
      walletType: _selectedType,
      name: _nameController.text.isEmpty ? null : _nameController.text,
    ));
  }

  @override
  void dispose() {
    _nameController.dispose();
    super.dispose();
  }
}
