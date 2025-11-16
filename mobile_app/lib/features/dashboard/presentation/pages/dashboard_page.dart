import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';
import 'package:go_router/go_router.dart';

import '../../../../core/theme/app_theme.dart';
import '../../../../core/utils/constants.dart';
import '../../../auth/presentation/bloc/auth_bloc.dart';
import '../../../portfolio/presentation/bloc/portfolio_bloc.dart';
import '../../../wallet/presentation/bloc/wallet_bloc.dart';
import '../widgets/portfolio_summary_card.dart';
import '../widgets/quick_actions_section.dart';
import '../widgets/recent_transactions_section.dart';
import '../widgets/market_overview_section.dart';
import '../widgets/crypto_prices_section.dart';

class DashboardPage extends StatefulWidget {
  const DashboardPage({Key? key}) : super(key: key);

  @override
  State<DashboardPage> createState() => _DashboardPageState();
}

class _DashboardPageState extends State<DashboardPage> {
  final ScrollController _scrollController = ScrollController();
  bool _showGreeting = true;

  @override
  void initState() {
    super.initState();
    _loadDashboardData();
  }

  void _loadDashboardData() {
    context.read<PortfolioBloc>().add(LoadPortfolioEvent());
    context.read<WalletBloc>().add(LoadWalletsEvent());
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      body: RefreshIndicator(
        onRefresh: () async {
          _loadDashboardData();
        },
        child: CustomScrollView(
          controller: _scrollController,
          slivers: [
            // Custom App Bar with Greeting
            SliverAppBar(
              expandedHeight: 120.0,
              floating: false,
              pinned: true,
              backgroundColor: AppTheme.primaryColor,
              flexibleSpace: FlexibleSpaceBar(
                background: Container(
                  decoration: const BoxDecoration(
                    gradient: AppTheme.primaryGradient,
                  ),
                  child: SafeArea(
                    child: Padding(
                      padding: const EdgeInsets.all(20.0),
                      child: Column(
                        crossAxisAlignment: CrossAxisAlignment.start,
                        mainAxisAlignment: MainAxisAlignment.end,
                        children: [
                          BlocBuilder<AuthBloc, AuthState>(
                            builder: (context, state) {
                              String userName = 'User';
                              if (state is AuthenticatedState) {
                                userName = state.user.firstName ?? 'User';
                              }
                              return Text(
                                _getGreeting(userName),
                                style: const TextStyle(
                                  color: Colors.white,
                                  fontSize: 24,
                                  fontWeight: FontWeight.bold,
                                ),
                              );
                            },
                          ),
                          const SizedBox(height: 8),
                          const Text(
                            'Welcome to Crypto Bank',
                            style: TextStyle(
                              color: Colors.white70,
                              fontSize: 16,
                            ),
                          ),
                        ],
                      ),
                    ),
                  ),
                ),
              ),
              actions: [
                IconButton(
                  icon: const Icon(Icons.notifications_outlined, color: Colors.white),
                  onPressed: () => context.push('/dashboard/notifications'),
                ),
                IconButton(
                  icon: const Icon(Icons.person_outline, color: Colors.white),
                  onPressed: () => context.push('/more/profile'),
                ),
              ],
            ),

            // Dashboard Content
            SliverToBoxAdapter(
              child: Column(
                children: [
                  // Portfolio Summary
                  const PortfolioSummaryCard(),
                  
                  // Quick Actions
                  const QuickActionsSection(),
                  
                  // Market Overview
                  const MarketOverviewSection(),
                  
                  // Crypto Prices
                  const CryptoPricesSection(),
                  
                  // Recent Transactions
                  const RecentTransactionsSection(),
                  
                  const SizedBox(height: 100), // Bottom padding for navigation
                ],
              ),
            ),
          ],
        ),
      ),
    );
  }

  String _getGreeting(String userName) {
    final hour = DateTime.now().hour;
    String greeting;

    if (hour < 12) {
      greeting = 'Good morning';
    } else if (hour < 17) {
      greeting = 'Good afternoon';
    } else {
      greeting = 'Good evening';
    }

    return '$greeting, $userName!';
  }

  @override
  void dispose() {
    _scrollController.dispose();
    super.dispose();
  }
}