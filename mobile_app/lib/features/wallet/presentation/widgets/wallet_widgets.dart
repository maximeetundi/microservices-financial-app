import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';

import '../../../../core/theme/app_theme.dart';
import '../../../../core/widgets/custom_button.dart';
import '../../domain/entities/wallet.dart';
import '../../domain/entities/transaction.dart' as tx;
import '../bloc/wallet_bloc.dart';

/// Quick action buttons for wallet operations
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
    return Row(
      mainAxisAlignment: MainAxisAlignment.spaceEvenly,
      children: [
        _buildActionButton(
          icon: Icons.arrow_upward,
          label: 'Send',
          color: AppTheme.primaryColor,
          onTap: onSendPressed,
        ),
        _buildActionButton(
          icon: Icons.arrow_downward,
          label: 'Receive',
          color: AppTheme.secondaryColor,
          onTap: onReceivePressed,
        ),
        _buildActionButton(
          icon: Icons.add,
          label: 'Buy',
          color: AppTheme.warningColor,
          onTap: onBuyPressed,
        ),
        _buildActionButton(
          icon: Icons.swap_horiz,
          label: 'Exchange',
          color: AppTheme.infoColor,
          onTap: onExchangePressed,
        ),
      ],
    );
  }

  Widget _buildActionButton({
    required IconData icon,
    required String label,
    required Color color,
    VoidCallback? onTap,
  }) {
    return GestureDetector(
      onTap: onTap,
      child: Column(
        children: [
          Container(
            width: 56,
            height: 56,
            decoration: BoxDecoration(
              color: color.withOpacity(0.1),
              shape: BoxShape.circle,
            ),
            child: Icon(
              icon,
              color: color,
              size: 24,
            ),
          ),
          const SizedBox(height: 8),
          Text(
            label,
            style: TextStyle(
              fontSize: 12,
              fontWeight: FontWeight.w500,
              color: Colors.grey.shade700,
            ),
          ),
        ],
      ),
    );
  }
}

/// Card displaying wallet information
class WalletCard extends StatelessWidget {
  final Wallet wallet;
  final VoidCallback? onTap;

  const WalletCard({
    Key? key,
    required this.wallet,
    this.onTap,
  }) : super(key: key);

  @override
  Widget build(BuildContext context) {
    return GestureDetector(
      onTap: onTap,
      child: Container(
        padding: const EdgeInsets.all(16),
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
        child: Row(
          children: [
            Container(
              width: 48,
              height: 48,
              decoration: BoxDecoration(
                color: _getCurrencyColor(wallet.currency).withOpacity(0.1),
                shape: BoxShape.circle,
              ),
              child: Icon(
                _getCurrencyIcon(wallet.currency),
                color: _getCurrencyColor(wallet.currency),
              ),
            ),
            const SizedBox(width: 16),
            Expanded(
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  Text(
                    wallet.name ?? '${wallet.currency} Wallet',
                    style: const TextStyle(
                      fontWeight: FontWeight.w600,
                      fontSize: 16,
                    ),
                  ),
                  const SizedBox(height: 4),
                  Text(
                    '${wallet.balance.toStringAsFixed(8)} ${wallet.currency}',
                    style: TextStyle(
                      color: Colors.grey.shade600,
                      fontSize: 14,
                    ),
                  ),
                ],
              ),
            ),
            Column(
              crossAxisAlignment: CrossAxisAlignment.end,
              children: [
                Text(
                  '\$${(wallet.balance * wallet.usdRate).toStringAsFixed(2)}',
                  style: const TextStyle(
                    fontWeight: FontWeight.w600,
                    fontSize: 16,
                  ),
                ),
                const SizedBox(height: 4),
                Icon(
                  Icons.chevron_right,
                  color: Colors.grey.shade400,
                ),
              ],
            ),
          ],
        ),
      ),
    );
  }

  IconData _getCurrencyIcon(String currency) {
    switch (currency.toUpperCase()) {
      case 'BTC':
        return Icons.currency_bitcoin;
      case 'ETH':
        return Icons.diamond;
      case 'USD':
        return Icons.attach_money;
      case 'EUR':
        return Icons.euro;
      default:
        return Icons.monetization_on;
    }
  }

  Color _getCurrencyColor(String currency) {
    switch (currency.toUpperCase()) {
      case 'BTC':
        return AppTheme.bitcoinColor;
      case 'ETH':
        return AppTheme.ethereumColor;
      case 'USD':
        return AppTheme.usdColor;
      case 'EUR':
        return AppTheme.eurColor;
      default:
        return AppTheme.primaryColor;
    }
  }
}

/// List of recent transactions
class RecentTransactionsList extends StatelessWidget {
  final List<tx.Transaction> transactions;

  const RecentTransactionsList({
    Key? key,
    required this.transactions,
  }) : super(key: key);

  @override
  Widget build(BuildContext context) {
    if (transactions.isEmpty) {
      return Container(
        padding: const EdgeInsets.all(24),
        child: Center(
          child: Column(
            children: [
              Icon(
                Icons.receipt_long,
                size: 48,
                color: Colors.grey.shade400,
              ),
              const SizedBox(height: 12),
              Text(
                'No transactions yet',
                style: TextStyle(
                  color: Colors.grey.shade600,
                  fontSize: 16,
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
          margin: const EdgeInsets.only(bottom: 8),
          padding: const EdgeInsets.all(12),
          decoration: BoxDecoration(
            color: Colors.white,
            borderRadius: BorderRadius.circular(12),
          ),
          child: Row(
            children: [
              Container(
                width: 40,
                height: 40,
                decoration: BoxDecoration(
                  color: (isReceive ? AppTheme.successColor : AppTheme.errorColor)
                      .withOpacity(0.1),
                  shape: BoxShape.circle,
                ),
                child: Icon(
                  isReceive ? Icons.arrow_downward : Icons.arrow_upward,
                  color: isReceive ? AppTheme.successColor : AppTheme.errorColor,
                  size: 20,
                ),
              ),
              const SizedBox(width: 12),
              Expanded(
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    Text(
                      transaction.type.name.toUpperCase(),
                      style: const TextStyle(
                        fontWeight: FontWeight.w600,
                        fontSize: 14,
                      ),
                    ),
                    Text(
                      _formatDate(transaction.createdAt),
                      style: TextStyle(
                        color: Colors.grey.shade600,
                        fontSize: 12,
                      ),
                    ),
                  ],
                ),
              ),
              Text(
                '${isReceive ? '+' : '-'}${transaction.amount.toStringAsFixed(4)} ${transaction.currency}',
                style: TextStyle(
                  fontWeight: FontWeight.w600,
                  color: isReceive ? AppTheme.successColor : AppTheme.errorColor,
                ),
              ),
            ],
          ),
        );
      }).toList(),
    );
  }

  String _formatDate(DateTime date) {
    final now = DateTime.now();
    final diff = now.difference(date);
    if (diff.inDays == 0) {
      return 'Today';
    } else if (diff.inDays == 1) {
      return 'Yesterday';
    } else if (diff.inDays < 7) {
      return '${diff.inDays} days ago';
    }
    return '${date.day}/${date.month}/${date.year}';
  }
}

/// Bottom sheet for creating a new wallet
class CreateWalletBottomSheet extends StatefulWidget {
  const CreateWalletBottomSheet({Key? key}) : super(key: key);

  @override
  State<CreateWalletBottomSheet> createState() => _CreateWalletBottomSheetState();
}

class _CreateWalletBottomSheetState extends State<CreateWalletBottomSheet> {
  String _selectedCurrency = 'BTC';
  String _selectedType = 'crypto';
  final _nameController = TextEditingController();

  final List<Map<String, dynamic>> _currencies = [
    {'code': 'BTC', 'name': 'Bitcoin', 'icon': Icons.currency_bitcoin},
    {'code': 'ETH', 'name': 'Ethereum', 'icon': Icons.diamond},
    {'code': 'USD', 'name': 'US Dollar', 'icon': Icons.attach_money},
    {'code': 'EUR', 'name': 'Euro', 'icon': Icons.euro},
  ];

  @override
  Widget build(BuildContext context) {
    return Container(
      padding: const EdgeInsets.all(20),
      decoration: const BoxDecoration(
        color: Colors.white,
        borderRadius: BorderRadius.vertical(top: Radius.circular(20)),
      ),
      child: Column(
        mainAxisSize: MainAxisSize.min,
        crossAxisAlignment: CrossAxisAlignment.stretch,
        children: [
          Center(
            child: Container(
              width: 40,
              height: 4,
              decoration: BoxDecoration(
                color: Colors.grey.shade300,
                borderRadius: BorderRadius.circular(2),
              ),
            ),
          ),
          const SizedBox(height: 20),
          const Text(
            'Create New Wallet',
            style: TextStyle(
              fontSize: 24,
              fontWeight: FontWeight.bold,
            ),
            textAlign: TextAlign.center,
          ),
          const SizedBox(height: 24),
          
          // Currency Selection
          const Text(
            'Select Currency',
            style: TextStyle(
              fontSize: 16,
              fontWeight: FontWeight.w600,
            ),
          ),
          const SizedBox(height: 12),
          Wrap(
            spacing: 8,
            runSpacing: 8,
            children: _currencies.map((currency) {
              final isSelected = _selectedCurrency == currency['code'];
              return ChoiceChip(
                label: Row(
                  mainAxisSize: MainAxisSize.min,
                  children: [
                    Icon(
                      currency['icon'] as IconData,
                      size: 18,
                      color: isSelected ? Colors.white : AppTheme.primaryColor,
                    ),
                    const SizedBox(width: 4),
                    Text(currency['code'] as String),
                  ],
                ),
                selected: isSelected,
                onSelected: (selected) {
                  if (selected) {
                    setState(() => _selectedCurrency = currency['code'] as String);
                  }
                },
                selectedColor: AppTheme.primaryColor,
                labelStyle: TextStyle(
                  color: isSelected ? Colors.white : null,
                ),
              );
            }).toList(),
          ),
          
          const SizedBox(height: 24),
          
          // Wallet Name (optional)
          TextFormField(
            controller: _nameController,
            decoration: const InputDecoration(
              labelText: 'Wallet Name (Optional)',
              hintText: 'My Bitcoin Wallet',
              prefixIcon: Icon(Icons.label_outline),
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
                    content: Text('Wallet created successfully!'),
                    backgroundColor: AppTheme.successColor,
                  ),
                );
              } else if (state is WalletErrorState) {
                ScaffoldMessenger.of(context).showSnackBar(
                  SnackBar(
                    content: Text(state.message),
                    backgroundColor: AppTheme.errorColor,
                  ),
                );
              }
            },
            builder: (context, state) {
              final isLoading = state is WalletLoadingState;
              return CustomButton(
                text: isLoading ? 'Creating...' : 'Create Wallet',
                onPressed: isLoading ? null : _createWallet,
                isLoading: isLoading,
              );
            },
          ),
          const SizedBox(height: 16),
        ],
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
