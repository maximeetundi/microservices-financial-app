import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';
import 'package:go_router/go_router.dart';
import 'package:google_fonts/google_fonts.dart';
import 'dart:ui';

import '../../../../core/services/pin_service.dart';
import '../../../auth/presentation/pages/pin_setup_screen.dart';
import '../../../../core/widgets/animated_drawer.dart';
import '../../../../core/widgets/glass_container.dart';
import '../../../../core/widgets/stat_card.dart';
import '../../../../core/widgets/quick_action_button.dart';
import '../../../../core/widgets/credit_card_widget.dart';
import '../../../../core/theme/app_theme.dart';
import '../../../auth/presentation/bloc/auth_bloc.dart';
import '../../../wallet/presentation/bloc/wallet_bloc.dart';
import '../../../cards/presentation/bloc/cards_bloc.dart';
import '../../../transfer/presentation/bloc/transfer_bloc.dart';
import '../../../exchange/presentation/bloc/exchange_bloc.dart';

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

  // Real data lists
  List<Map<String, dynamic>> _cryptoMarkets = [];
  List<Map<String, dynamic>> _fiatRates = [];

  @override
  void initState() {
    super.initState();
    _fabAnimationController = AnimationController(
      vsync: this,
      duration: const Duration(milliseconds: 300),
    )..forward();
    
    // Load all required data
    context.read<WalletBloc>().add(LoadWalletsEvent());
    context.read<ExchangeBloc>().add(const LoadMarketsEvent());
    context.read<ExchangeBloc>().add(const LoadExchangeRatesEvent());
    context.read<CardsBloc>().add(const LoadCardsEvent());
    context.read<TransferBloc>().add(const LoadTransferDataEvent());
    
    // Check PIN status and show setup modal if needed
    _checkPinStatus();
  }

  Future<void> _refreshRates() async {
    setState(() => _refreshingRates = true);
    context.read<ExchangeBloc>().add(const LoadExchangeRatesEvent());
    // Simulate delay or wait for state change (could be improved with BlocListener)
    await Future.delayed(const Duration(seconds: 1)); 
    if (mounted) setState(() => _refreshingRates = false);
  }

  /// Check if user has set their PIN and show setup modal if not
  Future<void> _checkPinStatus() async {
    final pinService = PinService();
    
    try {
      // First check locally (fast)
      final hasLocalPin = await pinService.isPinSetLocally();
      if (hasLocalPin) return; // PIN already set
      
      // Check with API
      final hasPin = await pinService.checkPinStatus();
      
      if (!hasPin && mounted) {
        // Show PIN setup modal - user must set PIN before continuing
        await showModalBottomSheet(
          context: context,
          isScrollControlled: true,
          isDismissible: false,
          enableDrag: false,
          backgroundColor: Colors.transparent,
          builder: (context) => Container(
            height: MediaQuery.of(context).size.height * 0.85,
            decoration: BoxDecoration(
              color: Theme.of(context).scaffoldBackgroundColor,
              borderRadius: const BorderRadius.vertical(top: Radius.circular(24)),
            ),
            child: PinSetupScreen(
              canSkip: false,
              onPinSet: () {
                // PIN has been set successfully
                debugPrint('PIN défini avec succès depuis ModernHomePage');
              },
            ),
          ),
        );
      }
    } catch (e) {
      debugPrint('Error checking PIN status: $e');
      // Don't block user on error
    }
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
        body: MultiBlocListener(
          listeners: [
            BlocListener<ExchangeBloc, ExchangeState>(
              listener: (context, state) {
                if (state is MarketsLoadedState) {
                  setState(() {
                    List<dynamic> rawMarkets = [];
                    if (state.markets['markets'] is List) {
                       rawMarkets = state.markets['markets'] as List;
                    } else if (state.markets is List) {
                       rawMarkets = state.markets as List;
                    }
                    
                    // Map backend fields to UI fields
                    _cryptoMarkets = rawMarkets.map((market) {
                      final m = Map<String, dynamic>.from(market);
                      // Backend uses: symbol (BTC/USDT), base_asset (BTC), price, change_24h
                      // UI expects: symbol, name, price, change, bgColor
                      final baseAsset = m['base_asset']?.toString() ?? m['symbol']?.toString().split('/').first ?? '??';
                      
                      // Set bgColor based on symbol
                      Color bgColor;
                      switch (baseAsset) {
                        case 'BTC': bgColor = const Color(0xFFF7931A); break;
                        case 'ETH': bgColor = const Color(0xFF627EEA); break;
                        case 'SOL': bgColor = const Color(0xFF9945FF); break;
                        case 'BNB': bgColor = const Color(0xFFF0B90B); break;
                        case 'XRP': bgColor = const Color(0xFF23292F); break;
                        case 'ADA': bgColor = const Color(0xFF0033AD); break;
                        case 'DOGE': bgColor = const Color(0xFFC2A633); break;
                        default: bgColor = Colors.blueGrey;
                      }
                      
                      return {
                        'symbol': baseAsset,
                        'name': baseAsset,
                        'price': (m['price'] as num?)?.toDouble() ?? 0.0,
                        'change': (m['change_24h'] as num?)?.toDouble() ?? 0.0,
                        'bgColor': bgColor,
                      };
                    }).toList();
                  });
                } else if (state is ExchangeRatesLoadedState) {
                  setState(() {
                    if (state.rates.isNotEmpty) {
                      _fiatRates = [];
                      state.rates.forEach((key, value) {
                        // value could be a Map with rate info or just a number
                        double rateValue = 0.0;
                        double changeValue = 0.0;
                        String pair = '';
                        
                        if (value is Map) {
                          rateValue = (value['rate'] as num?)?.toDouble() ?? 0.0;
                          changeValue = (value['change_24h'] as num?)?.toDouble() ?? 0.0;
                          pair = '${value['from_currency'] ?? 'USD'}/${value['to_currency'] ?? key}';
                        } else if (value is num) {
                          rateValue = value.toDouble();
                          pair = 'USD/$key';
                        }
                        
                        _fiatRates.add({
                          'pair': pair,
                          'rate': rateValue,
                          'change': changeValue,
                        });
                      });
                    }
                  });
                }
              },
            ),
            BlocListener<CardsBloc, CardsState>(
              listener: (context, state) {
                if (state is CardsLoadedState) {
                  setState(() {
                    _activeCards = state.cards.where((c) => (c['status'] == 'ACTIVE' || c['status'] == 'active')).length;
                    _cardsBalance = state.cards.fold(0.0, (sum, c) => sum + ((c['balance'] is num) ? (c['balance'] as num).toDouble() : 0.0));
                  });
                }
              },
            ),
            BlocListener<TransferBloc, TransferState>(
              listener: (context, state) {
                if (state is TransferLoadedState) {
                  final now = DateTime.now();
                  final currentMonthTransfers = state.recentTransfers.where((t) => 
                    t.createdAt.month == now.month && t.createdAt.year == now.year
                  ).toList();
                  
                  setState(() {
                    _monthlyTransfers = currentMonthTransfers.length;
                    _monthlyVolume = currentMonthTransfers.fold(0.0, (sum, t) => sum + t.amount);
                  });
                }
              },
            ),
          ],
          child: Container(
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
        // Values are now set by BlocListeners from CardsBloc and TransferBloc
        
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
    // Null-safe parsing of values
    final change = (crypto['change'] as num?)?.toDouble() ?? 0.0;
    final price = (crypto['price'] as num?)?.toDouble() ?? 0.0;
    final symbol = crypto['symbol']?.toString() ?? '??';
    final name = crypto['name']?.toString() ?? 'Unknown';
    final bgColor = crypto['bgColor'] as Color? ?? Colors.blueGrey;
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
              color: bgColor.withOpacity(0.2),
              borderRadius: BorderRadius.circular(12),
            ),
            child: Center(
              child: Text(
                symbol.length >= 2 ? symbol.substring(0, 2) : symbol,
                style: GoogleFonts.inter(
                  fontSize: 16,
                  fontWeight: FontWeight.bold,
                  color: bgColor,
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
                  name,
                  style: GoogleFonts.inter(
                    fontSize: 16,
                    fontWeight: FontWeight.w600,
                    color: isDark ? Colors.white : const Color(0xFF1E293B),
                  ),
                ),
                Text(
                  symbol,
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
                '\$${price.toStringAsFixed(price >= 100 ? 0 : 2)}',
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
            SizedBox(
              height: 200,
              child: BlocBuilder<CardsBloc, CardsState>(
                builder: (context, state) {
                  List<Widget> cardWidgets = [];
                  
                  if (state is CardsLoadedState && state.cards.isNotEmpty) {
                     cardWidgets = state.cards.map((card) {
                       return Padding(
                         padding: const EdgeInsets.only(right: 16),
                         child: SizedBox(
                           width: 320,
                           child: CreditCardWidget(
                             type: (card['type'] == 'virtual') ? CreditCardType.virtual : CreditCardType.physical,
                             cardNumber: card['card_number'] ?? '•••• •••• •••• ••••',
                             cardholderName: card['cardholder_name'] ?? 'Utilisateur',
                             expiryDate: card['expiry_date'] ?? '12/28',
                             balance: (card['balance'] is num) ? (card['balance'] as num).toDouble() : 0.0,
                             status: card['status'] ?? 'Active',
                             onTap: () => context.push('/more/cards'),
                           ),
                         ),
                       );
                     }).toList();
                  }

                  // Add "New Card" button at the end
                  cardWidgets.add(_buildAddCardButton(context, isDark));

                  return ListView(
                    scrollDirection: Axis.horizontal,
                    children: cardWidgets,
                  );
                },
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
                // Null-safe parsing
                final change = (rate['change'] as num?)?.toDouble() ?? 0.0;
                final rateValue = (rate['rate'] as num?)?.toDouble() ?? 1.0;
                final pair = rate['pair']?.toString() ?? 'N/A';
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
                        pair,
                        style: GoogleFonts.inter(
                          fontSize: 12,
                          fontWeight: FontWeight.w500,
                          color: isDark ? const Color(0xFF94A3B8) : const Color(0xFF64748B),
                        ),
                      ),
                      const SizedBox(height: 4),
                      Text(
                        rateValue.toStringAsFixed(4),
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

  Widget _buildAddCardButton(BuildContext context, bool isDark) {
    return GestureDetector(
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
    );
  }

}
