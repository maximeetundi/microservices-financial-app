import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';
import 'package:go_router/go_router.dart';
import 'package:google_fonts/google_fonts.dart';
import 'package:qr_flutter/qr_flutter.dart';

import '../../../../core/theme/app_theme.dart';
import '../../../../core/widgets/custom_button.dart';
import '../../../../core/widgets/loading_widget.dart';
import '../../../../core/widgets/glass_container.dart';
import '../bloc/wallet_bloc.dart';
import '../widgets/deposit_bottom_sheet.dart';
import '../widgets/wallet_widgets.dart';

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
    final isDark = Theme.of(context).brightness == Brightness.dark;

    return Scaffold(
      backgroundColor: Colors.transparent, // Allow gradient to show
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
               // Custom App Bar Area
               Padding(
                 padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 12),
                 child: Row(
                   mainAxisAlignment: MainAxisAlignment.spaceBetween,
                   children: [
                     const SizedBox(width: 48), // Spacer for centering
                     Text(
                       'Mes Portefeuilles',
                       style: GoogleFonts.inter(
                         fontSize: 20,
                         fontWeight: FontWeight.bold,
                         color: isDark ? Colors.white : AppTheme.textPrimaryColor,
                       ),
                     ),
                     GlassContainer(
                        padding: EdgeInsets.zero,
                        width: 40,
                        height: 40,
                        borderRadius: 12,
                        child: IconButton(
                          icon: Icon(Icons.add, color: isDark ? Colors.white : AppTheme.textPrimaryColor),
                          onPressed: _showCreateWalletBottomSheet,
                        ),
                     ),
                   ],
                 ),
               ),
               Expanded(
                 child: BlocBuilder<WalletBloc, WalletState>(
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
               ),
            ],
          ),
        ),
      ),
    );
  }

  Widget _buildWalletContent(WalletLoadedState state) {
    final wallets = state.wallets;
    final isDark = Theme.of(context).brightness == Brightness.dark;
    
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
              child: GlassContainer(
                gradient: AppTheme.cardGradient,
                padding: const EdgeInsets.all(24),
                borderRadius: 24,
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    Text(
                      'Valeur Totale',
                      style: GoogleFonts.inter(
                        color: Colors.white.withOpacity(0.8),
                        fontSize: 16,
                      ),
                    ),
                    const SizedBox(height: 8),
                    Text(
                      '\$${state.totalValue.toStringAsFixed(2)}',
                      style: GoogleFonts.inter(
                        color: Colors.white,
                        fontSize: 32,
                        fontWeight: FontWeight.bold,
                      ),
                    ),
                    const SizedBox(height: 16),
                    Row(
                      children: [
                        Container(
                          padding: const EdgeInsets.all(4),
                          decoration: BoxDecoration(
                            color: state.dailyChange >= 0 ? Colors.green.withOpacity(0.2) : Colors.red.withOpacity(0.2),
                            shape: BoxShape.circle,
                          ),
                          child: Icon(
                            state.dailyChange >= 0 ? Icons.trending_up : Icons.trending_down,
                            color: state.dailyChange >= 0 ? Colors.greenAccent : Colors.redAccent,
                            size: 16,
                          ),
                        ),
                        const SizedBox(width: 8),
                        Text(
                          '${state.dailyChange >= 0 ? '+' : ''}${state.dailyChange.toStringAsFixed(2)}%',
                          style: GoogleFonts.inter(
                            color: state.dailyChange >= 0 ? Colors.greenAccent : Colors.redAccent,
                            fontSize: 14,
                            fontWeight: FontWeight.w600,
                          ),
                        ),
                        Text(
                          ' Aujourd\'hui',
                          style: GoogleFonts.inter(
                            color: Colors.white70,
                            fontSize: 14,
                          ),
                        ),
                      ],
                    ),
                  ],
                ),
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
                        'Transactions Récentes',
                        style: GoogleFonts.inter(
                          fontSize: 18,
                          fontWeight: FontWeight.bold,
                          color: isDark ? Colors.white : AppTheme.textPrimaryColor,
                        ),
                      ),
                      TextButton(
                        onPressed: () {
                          // Navigate to all transactions
                        },
                        child: Text(
                          'Voir tout',
                          style: GoogleFonts.inter(
                            color: AppTheme.primaryColor,
                            fontWeight: FontWeight.w600,
                          ),
                        ),
                      ),
                    ],
                  ),
                  const SizedBox(height: 12),
                  GlassContainer(
                    padding: const EdgeInsets.all(12),
                    borderRadius: 16,
                     child: RecentTransactionsList(
                      transactions: state.recentTransactions,
                    ),
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
    // Use deposit bottom sheet if wallet is loaded
    final state = context.read<WalletBloc>().state;
    if (state is WalletLoadedState && state.wallets.isNotEmpty) {
      final wallet = state.wallets.first;
      showModalBottomSheet(
        context: context,
        isScrollControlled: true,
        backgroundColor: Colors.transparent,
        builder: (context) => DepositBottomSheet(
          walletId: wallet.id,
          walletCurrency: wallet.currency,
          onSuccess: () {
            _loadWallets(); // Refresh wallets after deposit
          },
        ),
      );
    } else {
      // Fallback to original options if no wallet
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
                },
              ),
              ListTile(
                leading: const Icon(Icons.account_balance),
                title: const Text('Bank Transfer'),
                subtitle: const Text('Lower fees'),
                onTap: () {
                  Navigator.pop(context);
                },
              ),
              ListTile(
                leading: const Icon(Icons.swap_horiz),
                title: const Text('P2P Exchange'),
                subtitle: const Text('Buy from other users'),
                onTap: () {
                  Navigator.pop(context);
                },
              ),
            ],
          ),
        ),
      );
    }
  }
}