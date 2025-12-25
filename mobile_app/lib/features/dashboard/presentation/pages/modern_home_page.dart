import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';
import 'package:go_router/go_router.dart';
import 'package:google_fonts/google_fonts.dart';
import 'dart:ui';

import '../../../../core/widgets/animated_drawer.dart';
import '../../../../core/widgets/glass_container.dart';
import '../../../../core/widgets/stat_card.dart';
import '../../../../core/widgets/quick_action_button.dart';
import '../../../../core/widgets/credit_card_widget.dart';
import '../../../../core/theme/app_theme.dart';
import '../../../auth/presentation/bloc/auth_bloc.dart';
import '../../../wallet/presentation/bloc/wallet_bloc.dart';

/// Modern Home Page matching web frontend design exactly
class ModernHomePage extends StatefulWidget {
  const ModernHomePage({super.key});

  @override
  State<ModernHomePage> createState() => _ModernHomePageState();
}

class _ModernHomePageState extends State<ModernHomePage>
    with TickerProviderStateMixin {
  final GlobalKey<AnimatedDrawerState> _drawerKey = GlobalKey();
  late AnimationController _fabAnimationController;
  bool _refreshingRates = false;

  // Stats data
  double _totalBalance = 0.0;
  double _cryptoBalance = 0.0;
  double _cardsBalance = 0.0;
  int _activeCards = 0;
  int _monthlyTransfers = 0;
  double _monthlyVolume = 0.0;

  // Sample crypto markets data (replace with API)
  final List<Map<String, dynamic>> _cryptoMarkets = [
    {'name': 'Bitcoin', 'symbol': 'BTC', 'price': 43250.0, 'change': 2.4, 'bgColor': const Color(0xFFF7931A)},
    {'name': 'Ethereum', 'symbol': 'ETH', 'price': 2280.0, 'change': -1.2, 'bgColor': const Color(0xFF627EEA)},
    {'name': 'Solana', 'symbol': 'SOL', 'price': 98.50, 'change': 5.8, 'bgColor': const Color(0xFF9945FF)},
    {'name': 'BNB', 'symbol': 'BNB', 'price': 312.0, 'change': 0.8, 'bgColor': const Color(0xFFF0B90B)},
  ];

  // Sample fiat rates (replace with API)
  final List<Map<String, dynamic>> _fiatRates = [
    {'pair': 'EUR/USD', 'rate': 1.0832, 'change': 0.0012},
    {'pair': 'GBP/USD', 'rate': 1.2714, 'change': -0.0008},
    {'pair': 'USD/XOF', 'rate': 605.50, 'change': 0.15},
    {'pair': 'EUR/GBP', 'rate': 0.8520, 'change': 0.0005},
    {'pair': 'USD/CAD', 'rate': 1.3421, 'change': -0.0021},
    {'pair': 'USD/JPY', 'rate': 149.25, 'change': 0.35},
  ];

  @override
  void initState() {
    super.initState();
    _fabAnimationController = AnimationController(
      vsync: this,
      duration: const Duration(milliseconds: 300),
    )..forward();
    
    // Load wallet data
    context.read<WalletBloc>().add(LoadWalletsEvent());
  }

  @override
  void dispose() {
    _fabAnimationController.dispose();
    super.dispose();
  }

  String _formatMoney(double amount) {
    return '\$${amount.toStringAsFixed(2).replaceAllMapped(
      RegExp(r'(\d{1,3})(?=(\d{3})+(?!\d))'),
      (Match m) => '${m[1]},',
    )}';
  }

  @override
  Widget build(BuildContext context) {
    final isDark = Theme.of(context).brightness == Brightness.dark;
    
    return AnimatedDrawer(
      key: _drawerKey,
      child: Scaffold(
        backgroundColor: Colors.transparent,
        body: Container(
          decoration: BoxDecoration(
            // Match web background gradient exactly
            gradient: LinearGradient(
              begin: Alignment.topLeft,
              end: Alignment.bottomRight,
              colors: isDark 
                  ? [const Color(0xFF0F0F1A), const Color(0xFF1A1A2E), const Color(0xFF16213E)]
                  : [const Color(0xFFFAFBFC), const Color(0xFFF5F7F9), const Color(0xFFEEF1F5)],
            ),
          ),
          child: CustomScrollView(
            physics: const BouncingScrollPhysics(),
            slivers: [
              _buildAppBar(context),
              SliverToBoxAdapter(
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    // Header with greeting
                    _buildHeader(context),
                    const SizedBox(height: 24),
                    
                    // Stats Cards Grid (4 cards like web)
                    _buildStatsGrid(),
                    const SizedBox(height: 32),
                    
                    // Quick Actions
                    _buildQuickActionsSection(context),
                    const SizedBox(height: 32),
                    
                    // Crypto Markets & Recent Activity (side by side on web, stacked on mobile)
                    _buildCryptoMarketsSection(context),
                    const SizedBox(height: 24),
                    
                    _buildRecentActivitySection(context),
                    const SizedBox(height: 32),
                    
                    // My Cards Section
                    _buildCardsSection(context),
                    const SizedBox(height: 32),
                    
                    // Exchange Rates
                    _buildExchangeRatesSection(context),
                    
                    const SizedBox(height: 100),
                  ],
                ),
              ),
            ],
          ),
        ),
        
        // Floating Send Money Button
        floatingActionButton: ScaleTransition(
          scale: _fabAnimationController,
          child: Container(
            decoration: BoxDecoration(
              gradient: AppTheme.primaryGradient,
              borderRadius: BorderRadius.circular(30),
              boxShadow: [
                BoxShadow(
                  color: AppTheme.primaryColor.withOpacity(0.4),
                  blurRadius: 20,
                  offset: const Offset(0, 8),
                ),
              ],
            ),
            child: FloatingActionButton.extended(
              onPressed: () => context.push('/more/transfer'),
              backgroundColor: Colors.transparent,
              elevation: 0,
              icon: const Icon(Icons.send_rounded, color: Colors.white),
              label: Text(
                'Envoyer',
                style: GoogleFonts.inter(
                  color: Colors.white,
                  fontWeight: FontWeight.bold,
                ),
              ),
            ),
          ),
        ),
        floatingActionButtonLocation: FloatingActionButtonLocation.centerFloat,
      ),
    );
  }

  Widget _buildAppBar(BuildContext context) {
    final isDark = Theme.of(context).brightness == Brightness.dark;
    
    return SliverAppBar(
      expandedHeight: 0,
      floating: true,
      pinned: false,
      backgroundColor: Colors.transparent,
      elevation: 0,
      leading: Padding(
        padding: const EdgeInsets.all(8.0),
        child: GlassContainer(
          borderRadius: 12,
          blur: 10,
          showTopHighlight: false,
          child: IconButton(
            icon: Icon(
              Icons.menu_rounded,
              color: isDark ? Colors.white : AppTheme.textPrimaryColor,
            ),
            onPressed: () => _drawerKey.currentState?.toggleDrawer(),
          ),
        ),
      ),
      actions: [
        Padding(
          padding: const EdgeInsets.only(right: 8.0, top: 8.0, bottom: 8.0),
          child: GlassContainer(
            borderRadius: 12,
            blur: 10,
            width: 48,
            height: 48,
            showTopHighlight: false,
            child: Stack(
              children: [
                Center(
                  child: IconButton(
                    icon: Icon(
                      Icons.notifications_outlined,
                      color: isDark ? Colors.white : AppTheme.textPrimaryColor,
                    ),
                    onPressed: () => context.push('/dashboard/notifications'),
                  ),
                ),
                Positioned(
                  right: 12,
                  top: 12,
                  child: Container(
                    width: 8,
                    height: 8,
                    decoration: const BoxDecoration(
                      color: AppTheme.errorColor,
                      shape: BoxShape.circle,
                    ),
                  ),
                ),
              ],
            ),
          ),
        ),
        const SizedBox(width: 8),
      ],
    );
  }

  Widget _buildHeader(BuildContext context) {
    final isDark = Theme.of(context).brightness == Brightness.dark;
    
    return Padding(
      padding: const EdgeInsets.symmetric(horizontal: 20),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          BlocBuilder<AuthBloc, AuthState>(
            builder: (context, state) {
              String name = 'Utilisateur';
              if (state is AuthenticatedState) {
                name = state.user.fullName;
              }
              return Text(
                'Bonjour, $name 👋',
                style: GoogleFonts.inter(
                  fontSize: 28,
                  fontWeight: FontWeight.bold,
                  color: isDark ? Colors.white : const Color(0xFF1E293B),
                ),
              );
            },
          ),
          const SizedBox(height: 4),
          Text(
            _getFormattedDate(),
            style: GoogleFonts.inter(
              fontSize: 14,
              color: isDark ? const Color(0xFF94A3B8) : const Color(0xFF64748B),
            ),
          ),
        ],
      ),
    );
  }

  String _getFormattedDate() {
    final now = DateTime.now();
    final weekdays = ['dimanche', 'lundi', 'mardi', 'mercredi', 'jeudi', 'vendredi', 'samedi'];
    final months = ['janvier', 'février', 'mars', 'avril', 'mai', 'juin', 
                    'juillet', 'août', 'septembre', 'octobre', 'novembre', 'décembre'];
    return '${weekdays[now.weekday % 7]} ${now.day} ${months[now.month - 1]} ${now.year}';
  }

  Widget _buildStatsGrid() {
    return BlocBuilder<WalletBloc, WalletState>(
      builder: (context, state) {
        if (state is WalletLoadedState) {
          _totalBalance = 0.0;
          _cryptoBalance = 0.0;
          for (var wallet in state.wallets) {
            _totalBalance += wallet.balance;
            if (['BTC', 'ETH', 'USDT', 'USDC', 'SOL'].contains(wallet.currency)) {
              _cryptoBalance += wallet.balance;
            }
          }
        }
        _activeCards = 2; // Replace with actual API data
        _cardsBalance = 3500.0; // Replace with actual API data
        _monthlyTransfers = 12;
        _monthlyVolume = 4500.0;
        
        return Padding(
          padding: const EdgeInsets.symmetric(horizontal: 20),
          child: GridView.count(
            shrinkWrap: true,
            physics: const NeverScrollableScrollPhysics(),
            crossAxisCount: 2,
            crossAxisSpacing: 16,
            mainAxisSpacing: 16,
            childAspectRatio: 1.0,
            children: [
              StatCard(
                title: 'Solde Total',
                value: _formatMoney(_totalBalance),
                badge: '+5.2%',
                isBadgePositive: true,
                variant: StatCardVariant.blue,
                icon: const Icon(Icons.attach_money, color: Color(0xFF60A5FA), size: 24),
              ),
              StatCard(
                title: 'Crypto Portfolio',
                value: _formatMoney(_cryptoBalance),
                badge: '+12.8%',
                isBadgePositive: true,
                variant: StatCardVariant.green,
                icon: const Text('₿', style: TextStyle(fontSize: 24, color: Color(0xFF34D399))),
              ),
              StatCard(
                title: 'Cartes Actives',
                value: _activeCards.toString(),
                subtitle: 'Solde: ${_formatMoney(_cardsBalance)}',
                variant: StatCardVariant.purple,
                icon: const Icon(Icons.credit_card, color: Color(0xFFC084FC), size: 24),
              ),
              StatCard(
                title: 'Transferts ce mois',
                value: _monthlyTransfers.toString(),
                subtitle: 'Volume: ${_formatMoney(_monthlyVolume)}',
                variant: StatCardVariant.orange,
                icon: const Icon(Icons.swap_horiz, color: Color(0xFFFBBF24), size: 24),
              ),
            ],
          ),
        );
      },
    );
  }

  Widget _buildQuickActionsSection(BuildContext context) {
    final isDark = Theme.of(context).brightness == Brightness.dark;
    
    return Padding(
      padding: const EdgeInsets.symmetric(horizontal: 20),
      child: GlassContainer(
        padding: const EdgeInsets.all(20),
        borderRadius: 24,
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Text(
              '🚀 Actions Rapides',
              style: GoogleFonts.inter(
                fontSize: 18,
                fontWeight: FontWeight.bold,
                color: isDark ? Colors.white : const Color(0xFF1E293B),
              ),
            ),
            const SizedBox(height: 20),
            QuickActionsGrid(
              crossAxisCount: 3,
              actions: [
                QuickActionItem(
                  emoji: '📷',
                  label: 'Scanner',
                  onTap: () => context.push('/more/merchant/scan'),
                ),
                QuickActionItem(
                  emoji: '💸',
                  label: 'Envoyer',
                  onTap: () => context.push('/more/transfer'),
                ),
                QuickActionItem(
                  emoji: '💳',
                  label: 'Mes Cartes',
                  onTap: () => context.push('/more/cards'),
                ),
                QuickActionItem(
                  emoji: '👛',
                  label: 'Portefeuilles',
                  onTap: () => context.push('/wallet'),
                ),
                QuickActionItem(
                  emoji: '₿',
                  label: 'Acheter Crypto',
                  onTap: () => context.push('/exchange'),
                ),
                QuickActionItem(
                  emoji: '💱',
                  label: 'Convertir',
                  onTap: () => context.push('/exchange/fiat'),
                ),
              ],
            ),
          ],
        ),
      ),
    );
  }

  Widget _buildCryptoMarketsSection(BuildContext context) {
    final isDark = Theme.of(context).brightness == Brightness.dark;
    
    return Padding(
      padding: const EdgeInsets.symmetric(horizontal: 20),
      child: GlassContainer(
        padding: const EdgeInsets.all(20),
        borderRadius: 24,
        child: Column(
          children: [
            Row(
              mainAxisAlignment: MainAxisAlignment.spaceBetween,
              children: [
                Text(
                  '📊 Marchés Crypto',
                  style: GoogleFonts.inter(
                    fontSize: 18,
                    fontWeight: FontWeight.bold,
                    color: isDark ? Colors.white : const Color(0xFF1E293B),
                  ),
                ),
                GestureDetector(
                  onTap: () => context.push('/exchange'),
                  child: Text(
                    'Voir tout →',
                    style: GoogleFonts.inter(
                      fontSize: 14,
                      fontWeight: FontWeight.w500,
                      color: const Color(0xFF818CF8),
                    ),
                  ),
                ),
              ],
            ),
            const SizedBox(height: 16),
            ...List.generate(_cryptoMarkets.length, (index) {
              final crypto = _cryptoMarkets[index];
              return _buildCryptoMarketItem(crypto, isDark);
            }),
          ],
        ),
      ),
    );
  }

  Widget _buildCryptoMarketItem(Map<String, dynamic> crypto, bool isDark) {
    final change = crypto['change'] as double;
    final isPositive = change >= 0;
    
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
              color: (crypto['bgColor'] as Color).withOpacity(0.2),
              borderRadius: BorderRadius.circular(12),
            ),
            child: Center(
              child: Text(
                crypto['symbol'].toString().substring(0, 2),
                style: GoogleFonts.inter(
                  fontSize: 16,
                  fontWeight: FontWeight.bold,
                  color: crypto['bgColor'] as Color,
                ),
              ),
            ),
          ),
          const SizedBox(width: 16),
          Expanded(
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Text(
                  crypto['name'] as String,
                  style: GoogleFonts.inter(
                    fontSize: 16,
                    fontWeight: FontWeight.w600,
                    color: isDark ? Colors.white : const Color(0xFF1E293B),
                  ),
                ),
                Text(
                  crypto['symbol'] as String,
                  style: GoogleFonts.inter(
                    fontSize: 14,
                    color: isDark ? const Color(0xFF94A3B8) : const Color(0xFF64748B),
                  ),
                ),
              ],
            ),
          ),
          Column(
            crossAxisAlignment: CrossAxisAlignment.end,
            children: [
              Text(
                '\$${(crypto['price'] as double).toStringAsFixed(crypto['price'] >= 100 ? 0 : 2)}',
                style: GoogleFonts.inter(
                  fontSize: 16,
                  fontWeight: FontWeight.w600,
                  color: isDark ? Colors.white : const Color(0xFF1E293B),
                ),
              ),
              Text(
                '${isPositive ? '+' : ''}${change.toStringAsFixed(2)}%',
                style: GoogleFonts.inter(
                  fontSize: 14,
                  fontWeight: FontWeight.w500,
                  color: isPositive ? const Color(0xFF10B981) : const Color(0xFFEF4444),
                ),
              ),
            ],
          ),
        ],
      ),
    );
  }

  Widget _buildRecentActivitySection(BuildContext context) {
    final isDark = Theme.of(context).brightness == Brightness.dark;
    
    return Padding(
      padding: const EdgeInsets.symmetric(horizontal: 20),
      child: GlassContainer(
        padding: const EdgeInsets.all(20),
        borderRadius: 24,
        child: Column(
          children: [
            Row(
              mainAxisAlignment: MainAxisAlignment.spaceBetween,
              children: [
                Text(
                  '🕒 Activité Récente',
                  style: GoogleFonts.inter(
                    fontSize: 18,
                    fontWeight: FontWeight.bold,
                    color: isDark ? Colors.white : const Color(0xFF1E293B),
                  ),
                ),
                GestureDetector(
                  onTap: () => context.push('/transactions'),
                  child: Text(
                    'Voir tout →',
                    style: GoogleFonts.inter(
                      fontSize: 14,
                      fontWeight: FontWeight.w500,
                      color: const Color(0xFF818CF8),
                    ),
                  ),
                ),
              ],
            ),
            const SizedBox(height: 16),
            // Empty state for now - replace with actual transactions
            Container(
              padding: const EdgeInsets.all(24),
              child: Column(
                children: [
                  const Text('📊', style: TextStyle(fontSize: 40)),
                  const SizedBox(height: 12),
                  Text(
                    'Consultez vos transactions',
                    style: GoogleFonts.inter(
                      fontSize: 16,
                      fontWeight: FontWeight.bold,
                      color: isDark ? Colors.white : const Color(0xFF1E293B),
                    ),
                  ),
                  const SizedBox(height: 8),
                  Text(
                    'Voir l\'historique complet de vos transactions',
                    textAlign: TextAlign.center,
                    style: GoogleFonts.inter(
                      fontSize: 14,
                      color: AppTheme.textSecondaryColor,
                    ),
                  ),
                  const SizedBox(height: 16),
                  GestureDetector(
                    onTap: () => context.push('/transactions'),
                    child: Container(
                      padding: const EdgeInsets.symmetric(horizontal: 24, vertical: 12),
                      decoration: BoxDecoration(
                        gradient: AppTheme.primaryGradient,
                        borderRadius: BorderRadius.circular(12),
                      ),
                      child: Text(
                        'Voir les transactions',
                        style: GoogleFonts.inter(
                          color: Colors.white,
                          fontWeight: FontWeight.bold,
                        ),
                      ),
                    ),
                  ),
                ],
              ),
            ),
          ],
        ),
      ),
    );
  }

  Widget _buildCardsSection(BuildContext context) {
    final isDark = Theme.of(context).brightness == Brightness.dark;
    
    return Padding(
      padding: const EdgeInsets.symmetric(horizontal: 20),
      child: GlassContainer(
        padding: const EdgeInsets.all(20),
        borderRadius: 24,
        child: Column(
          children: [
            Row(
              mainAxisAlignment: MainAxisAlignment.spaceBetween,
              children: [
                Text(
                  '💳 Mes Cartes',
                  style: GoogleFonts.inter(
                    fontSize: 18,
                    fontWeight: FontWeight.bold,
                    color: isDark ? Colors.white : const Color(0xFF1E293B),
                  ),
                ),
                GestureDetector(
                  onTap: () => context.push('/more/cards'),
                  child: Text(
                    'Gérer →',
                    style: GoogleFonts.inter(
                      fontSize: 14,
                      fontWeight: FontWeight.w500,
                      color: const Color(0xFF818CF8),
                    ),
                  ),
                ),
              ],
            ),
            const SizedBox(height: 20),
            // Sample cards - replace with actual data
            SizedBox(
              height: 200,
              child: ListView(
                scrollDirection: Axis.horizontal,
                children: [
                  SizedBox(
                    width: 320,
                    child: CreditCardWidget(
                      type: CreditCardType.virtual,
                      cardNumber: '•••• •••• •••• 4532',
                      cardholderName: 'John Doe',
                      expiryDate: '12/28',
                      balance: 2500.00,
                      status: 'Active',
                      onTap: () => context.push('/more/cards'),
                    ),
                  ),
                  const SizedBox(width: 16),
                  // Add new card button
                  GestureDetector(
                    onTap: () => context.push('/more/cards'),
                    child: Container(
                      width: 200,
                      decoration: BoxDecoration(
                        borderRadius: BorderRadius.circular(24),
                        border: Border.all(
                          color: isDark ? const Color(0xFF334155) : const Color(0xFFE2E8F0),
                          width: 2,
                          style: BorderStyle.solid,
                        ),
                        color: isDark 
                            ? const Color(0xFF1E293B).withOpacity(0.3)
                            : const Color(0xFFF8FAFC).withOpacity(0.5),
                      ),
                      child: Column(
                        mainAxisAlignment: MainAxisAlignment.center,
                        children: [
                          Container(
                            width: 64,
                            height: 64,
                            decoration: BoxDecoration(
                              color: isDark 
                                  ? const Color(0xFF334155).withOpacity(0.5)
                                  : const Color(0xFFE2E8F0),
                              borderRadius: BorderRadius.circular(32),
                            ),
                            child: Icon(
                              Icons.add,
                              size: 32,
                              color: isDark ? const Color(0xFF94A3B8) : const Color(0xFF64748B),
                            ),
                          ),
                          const SizedBox(height: 12),
                          Text(
                            'Nouvelle Carte',
                            style: GoogleFonts.inter(
                              fontSize: 14,
                              fontWeight: FontWeight.w500,
                              color: isDark ? const Color(0xFF94A3B8) : const Color(0xFF64748B),
                            ),
                          ),
                        ],
                      ),
                    ),
                  ),
                ],
              ),
            ),
          ],
        ),
      ),
    );
  }

  Widget _buildExchangeRatesSection(BuildContext context) {
    final isDark = Theme.of(context).brightness == Brightness.dark;
    
    return Padding(
      padding: const EdgeInsets.symmetric(horizontal: 20),
      child: GlassContainer(
        padding: const EdgeInsets.all(20),
        borderRadius: 24,
        child: Column(
          children: [
            Row(
              mainAxisAlignment: MainAxisAlignment.spaceBetween,
              children: [
                Text(
                  '💱 Taux de Change',
                  style: GoogleFonts.inter(
                    fontSize: 18,
                    fontWeight: FontWeight.bold,
                    color: isDark ? Colors.white : const Color(0xFF1E293B),
                  ),
                ),
                GestureDetector(
                  onTap: _refreshRates,
                  child: Row(
                    children: [
                      if (_refreshingRates)
                        const SizedBox(
                          width: 16,
                          height: 16,
                          child: CircularProgressIndicator(strokeWidth: 2),
                        )
                      else
                        const Icon(Icons.refresh, size: 16, color: Color(0xFF818CF8)),
                      const SizedBox(width: 4),
                      Text(
                        'Actualiser',
                        style: GoogleFonts.inter(
                          fontSize: 14,
                          fontWeight: FontWeight.w500,
                          color: const Color(0xFF818CF8),
                        ),
                      ),
                    ],
                  ),
                ),
              ],
            ),
            const SizedBox(height: 16),
            GridView.builder(
              shrinkWrap: true,
              physics: const NeverScrollableScrollPhysics(),
              gridDelegate: const SliverGridDelegateWithFixedCrossAxisCount(
                crossAxisCount: 2,
                crossAxisSpacing: 12,
                mainAxisSpacing: 12,
                childAspectRatio: 1.8,
              ),
              itemCount: _fiatRates.length,
              itemBuilder: (context, index) {
                final rate = _fiatRates[index];
                final change = rate['change'] as double;
                final isPositive = change >= 0;
                
                return Container(
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
                  child: Column(
                    mainAxisAlignment: MainAxisAlignment.center,
                    children: [
                      Text(
                        rate['pair'] as String,
                        style: GoogleFonts.inter(
                          fontSize: 12,
                          fontWeight: FontWeight.w500,
                          color: isDark ? const Color(0xFF94A3B8) : const Color(0xFF64748B),
                        ),
                      ),
                      const SizedBox(height: 4),
                      Text(
                        (rate['rate'] as double).toStringAsFixed(4),
                        style: GoogleFonts.inter(
                          fontSize: 18,
                          fontWeight: FontWeight.bold,
                          color: isDark ? Colors.white : const Color(0xFF1E293B),
                        ),
                      ),
                      const SizedBox(height: 2),
                      Text(
                        '${isPositive ? '+' : ''}${(change * 100).toStringAsFixed(2)}%',
                        style: GoogleFonts.inter(
                          fontSize: 12,
                          fontWeight: FontWeight.w500,
                          color: isPositive ? const Color(0xFF10B981) : const Color(0xFFEF4444),
                        ),
                      ),
                    ],
                  ),
                );
              },
            ),
          ],
        ),
      ),
    );
  }

  void _refreshRates() {
    setState(() => _refreshingRates = true);
    // Simulate API call
    Future.delayed(const Duration(seconds: 1), () {
      if (mounted) {
        setState(() => _refreshingRates = false);
      }
    });
  }
}
