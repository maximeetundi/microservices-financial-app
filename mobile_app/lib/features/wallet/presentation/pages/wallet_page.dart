import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';
import 'package:go_router/go_router.dart';
import 'package:qr_flutter/qr_flutter.dart';

import '../../../../core/theme/app_theme.dart';
import '../../../../core/widgets/custom_button.dart';
import '../../../../core/widgets/loading_widget.dart';
import '../bloc/wallet_bloc.dart';
import '../widgets/wallet_card.dart';
import '../widgets/wallet_actions.dart';
import '../widgets/recent_transactions_list.dart';
import '../widgets/create_wallet_bottom_sheet.dart';

class WalletPage extends StatefulWidget {
  const WalletPage({Key? key}) : super(key: key);

  @override
  State<WalletPage> createState() => _WalletPageState();
}

class _WalletPageState extends State<WalletPage> {
  @override
  void initState() {
    super.initState();
    _loadWallets();
  }

  void _loadWallets() {
    context.read<WalletBloc>().add(LoadWalletsEvent());
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: const Text('My Wallets'),
        centerTitle: true,
        elevation: 0,
        actions: [
          IconButton(
            icon: const Icon(Icons.add),
            onPressed: _showCreateWalletBottomSheet,
          ),
        ],
      ),
      body: BlocBuilder<WalletBloc, WalletState>(
        builder: (context, state) {
          if (state is WalletLoadingState) {
            return const LoadingWidget();
          } else if (state is WalletLoadedState) {
            return _buildWalletContent(state);
          } else if (state is WalletErrorState) {
            return _buildErrorState(state.message);
          }
          return const SizedBox.shrink();
        },
      ),
    );
  }

  Widget _buildWalletContent(WalletLoadedState state) {
    final wallets = state.wallets;
    
    return RefreshIndicator(
      onRefresh: () async {
        _loadWallets();
      },
      child: CustomScrollView(
        slivers: [
          // Total Portfolio Value
          SliverToBoxAdapter(
            child: Container(
              margin: const EdgeInsets.all(16),
              padding: const EdgeInsets.all(20),
              decoration: BoxDecoration(
                gradient: AppTheme.primaryGradient,
                borderRadius: BorderRadius.circular(20),
                boxShadow: [
                  BoxShadow(
                    color: AppTheme.primaryColor.withOpacity(0.3),
                    blurRadius: 15,
                    offset: const Offset(0, 5),
                  ),
                ],
              ),
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  const Text(
                    'Total Portfolio Value',
                    style: TextStyle(
                      color: Colors.white70,
                      fontSize: 16,
                    ),
                  ),
                  const SizedBox(height: 8),
                  Text(
                    '\$${state.totalValue.toStringAsFixed(2)}',
                    style: const TextStyle(
                      color: Colors.white,
                      fontSize: 32,
                      fontWeight: FontWeight.bold,
                    ),
                  ),
                  const SizedBox(height: 16),
                  Row(
                    children: [
                      Icon(
                        state.dailyChange >= 0 ? Icons.trending_up : Icons.trending_down,
                        color: state.dailyChange >= 0 ? Colors.green : Colors.red,
                        size: 20,
                      ),
                      const SizedBox(width: 4),
                      Text(
                        '${state.dailyChange >= 0 ? '+' : ''}${state.dailyChange.toStringAsFixed(2)}%',
                        style: TextStyle(
                          color: state.dailyChange >= 0 ? Colors.green : Colors.red,
                          fontSize: 16,
                          fontWeight: FontWeight.w600,
                        ),
                      ),
                      const Text(
                        ' Today',
                        style: TextStyle(
                          color: Colors.white70,
                          fontSize: 16,
                        ),
                      ),
                    ],
                  ),
                ],
              ),
            ),
          ),

          // Quick Actions
          SliverToBoxAdapter(
            child: Padding(
              padding: const EdgeInsets.symmetric(horizontal: 16),
              child: WalletActions(
                onSendPressed: _showSendOptions,
                onReceivePressed: _showReceiveOptions,
                onBuyPressed: _showBuyOptions,
                onExchangePressed: () => context.push('/exchange'),
              ),
            ),
          ),

          // Wallets List
          if (wallets.isEmpty)
            SliverFillRemaining(
              child: _buildEmptyWalletState(),
            )
          else
            SliverPadding(
              padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 8),
              sliver: SliverList(
                delegate: SliverChildBuilderDelegate(
                  (context, index) {
                    final wallet = wallets[index];
                    return Padding(
                      padding: const EdgeInsets.only(bottom: 12),
                      child: WalletCard(
                        wallet: wallet,
                        onTap: () => context.push('/wallet/${wallet.id}'),
                      ),
                    );
                  },
                  childCount: wallets.length,
                ),
              ),
            ),

          // Recent Transactions
          SliverToBoxAdapter(
            child: Padding(
              padding: const EdgeInsets.all(16),
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  Row(
                    mainAxisAlignment: MainAxisAlignment.spaceBetween,
                    children: [
                      Text(
                        'Recent Transactions',
                        style: Theme.of(context).textTheme.headlineSmall,
                      ),
                      TextButton(
                        onPressed: () {
                          // Navigate to all transactions
                        },
                        child: const Text('View All'),
                      ),
                    ],
                  ),
                  const SizedBox(height: 12),
                  RecentTransactionsList(
                    transactions: state.recentTransactions,
                  ),
                ],
              ),
            ),
          ),

          // Bottom padding
          const SliverToBoxAdapter(
            child: SizedBox(height: 100),
          ),
        ],
      ),
    );
  }

  Widget _buildEmptyWalletState() {
    return Center(
      child: Column(
        mainAxisAlignment: MainAxisAlignment.center,
        children: [
          Icon(
            Icons.account_balance_wallet_outlined,
            size: 80,
            color: Colors.grey.shade400,
          ),
          const SizedBox(height: 20),
          Text(
            'No Wallets Yet',
            style: Theme.of(context).textTheme.headlineMedium?.copyWith(
              color: Colors.grey.shade600,
            ),
          ),
          const SizedBox(height: 8),
          Text(
            'Create your first wallet to start managing your crypto and fiat currencies',
            textAlign: TextAlign.center,
            style: Theme.of(context).textTheme.bodyMedium?.copyWith(
              color: Colors.grey.shade500,
            ),
          ),
          const SizedBox(height: 32),
          CustomButton(
            text: 'Create Wallet',
            onPressed: _showCreateWalletBottomSheet,
            icon: Icons.add,
          ),
        ],
      ),
    );
  }

  Widget _buildErrorState(String message) {
    return Center(
      child: Column(
        mainAxisAlignment: MainAxisAlignment.center,
        children: [
          Icon(
            Icons.error_outline,
            size: 80,
            color: AppTheme.errorColor,
          ),
          const SizedBox(height: 20),
          Text(
            'Something went wrong',
            style: Theme.of(context).textTheme.headlineMedium,
          ),
          const SizedBox(height: 8),
          Text(
            message,
            textAlign: TextAlign.center,
            style: Theme.of(context).textTheme.bodyMedium,
          ),
          const SizedBox(height: 32),
          CustomButton(
            text: 'Try Again',
            onPressed: _loadWallets,
          ),
        ],
      ),
    );
  }

  void _showCreateWalletBottomSheet() {
    showModalBottomSheet(
      context: context,
      isScrollControlled: true,
      backgroundColor: Colors.transparent,
      builder: (context) => const CreateWalletBottomSheet(),
    );
  }

  void _showSendOptions() {
    showModalBottomSheet(
      context: context,
      builder: (context) => Container(
        padding: const EdgeInsets.all(20),
        child: Column(
          mainAxisSize: MainAxisSize.min,
          children: [
            const Text(
              'Send Options',
              style: TextStyle(
                fontSize: 20,
                fontWeight: FontWeight.bold,
              ),
            ),
            const SizedBox(height: 20),
            ListTile(
              leading: const Icon(Icons.qr_code_scanner),
              title: const Text('Scan QR Code'),
              onTap: () {
                Navigator.pop(context);
                // Navigate to QR scanner
              },
            ),
            ListTile(
              leading: const Icon(Icons.person),
              title: const Text('Send to Contact'),
              onTap: () {
                Navigator.pop(context);
                // Navigate to contacts
              },
            ),
            ListTile(
              leading: const Icon(Icons.link),
              title: const Text('Enter Address'),
              onTap: () {
                Navigator.pop(context);
                // Navigate to manual address entry
              },
            ),
          ],
        ),
      ),
    );
  }

  void _showReceiveOptions() {
    showModalBottomSheet(
      context: context,
      builder: (context) => Container(
        padding: const EdgeInsets.all(20),
        child: Column(
          mainAxisSize: MainAxisSize.min,
          children: [
            const Text(
              'Receive Crypto',
              style: TextStyle(
                fontSize: 20,
                fontWeight: FontWeight.bold,
              ),
            ),
            const SizedBox(height: 20),
            
            // QR Code
            Container(
              padding: const EdgeInsets.all(20),
              decoration: BoxDecoration(
                color: Colors.white,
                borderRadius: BorderRadius.circular(12),
                border: Border.all(color: Colors.grey.shade300),
              ),
              child: QrImageView(
                data: 'bitcoin:1A1zP1eP5QGefi2DMPTfTL5SLmv7DivfNa',
                version: QrVersions.auto,
                size: 200.0,
              ),
            ),
            
            const SizedBox(height: 16),
            
            // Address
            Container(
              padding: const EdgeInsets.all(12),
              decoration: BoxDecoration(
                color: Colors.grey.shade100,
                borderRadius: BorderRadius.circular(8),
              ),
              child: Row(
                children: [
                  const Expanded(
                    child: Text(
                      '1A1zP1eP5QGefi2DMPTfTL5SLmv7DivfNa',
                      style: TextStyle(fontFamily: 'monospace'),
                    ),
                  ),
                  IconButton(
                    icon: const Icon(Icons.copy),
                    onPressed: () {
                      // Copy to clipboard
                    },
                  ),
                ],
              ),
            ),
            
            const SizedBox(height: 16),
            CustomButton(
              text: 'Share Address',
              onPressed: () {
                // Share address
              },
              icon: Icons.share,
            ),
          ],
        ),
      ),
    );
  }

  void _showBuyOptions() {
    showModalBottomSheet(
      context: context,
      builder: (context) => Container(
        padding: const EdgeInsets.all(20),
        child: Column(
          mainAxisSize: MainAxisSize.min,
          children: [
            const Text(
              'Buy Crypto',
              style: TextStyle(
                fontSize: 20,
                fontWeight: FontWeight.bold,
              ),
            ),
            const SizedBox(height: 20),
            ListTile(
              leading: const Icon(Icons.credit_card),
              title: const Text('Credit/Debit Card'),
              subtitle: const Text('Instant purchase'),
              onTap: () {
                Navigator.pop(context);
                // Navigate to card purchase
              },
            ),
            ListTile(
              leading: const Icon(Icons.account_balance),
              title: const Text('Bank Transfer'),
              subtitle: const Text('Lower fees'),
              onTap: () {
                Navigator.pop(context);
                // Navigate to bank transfer
              },
            ),
            ListTile(
              leading: const Icon(Icons.swap_horiz),
              title: const Text('P2P Exchange'),
              subtitle: const Text('Buy from other users'),
              onTap: () {
                Navigator.pop(context);
                // Navigate to P2P
              },
            ),
          ],
        ),
      ),
    );
  }
}