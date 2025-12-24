import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import 'package:go_router/go_router.dart';
import 'package:flutter_bloc/flutter_bloc.dart';
import 'dart:ui';

import '../../../../core/widgets/animated_drawer.dart';
import '../../../../core/widgets/glass_container.dart';
import '../../../../core/theme/app_theme.dart';
import '../../../auth/presentation/bloc/auth_bloc.dart';
import '../../../wallet/presentation/bloc/wallet_bloc.dart';

/// Modern Home Page with animated drawer and transfer focus
class ModernHomePage extends StatefulWidget {
  const ModernHomePage({super.key});

  @override
  State<ModernHomePage> createState() => _ModernHomePageState();
}

class _ModernHomePageState extends State<ModernHomePage>
    with TickerProviderStateMixin {
  final GlobalKey<AnimatedDrawerState> _drawerKey = GlobalKey();
  late AnimationController _fabAnimationController;

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

  @override
  Widget build(BuildContext context) {
    final isDark = Theme.of(context).brightness == Brightness.dark;
    
    return AnimatedDrawer(
      key: _drawerKey,
      child: Scaffold(
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
          child: CustomScrollView(
            physics: const BouncingScrollPhysics(),
            slivers: [
              // Modern App Bar
              _buildModernAppBar(context),
              
              // Content
              SliverToBoxAdapter(
                child: Column(
                  children: [
                    // Balance Card
                    _buildBalanceCard(),
                    
                    // Quick Actions - TRANSFER FOCUSED
                    _buildQuickActions(context),
                    
                    // Recent Activity
                    _buildRecentActivity(),
                    
                    // Services Section
                    _buildServicesGrid(context),
                    
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
              label: const Text(
                'Envoyer',
                style: TextStyle(
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

  Widget _buildModernAppBar(BuildContext context) {
    final isDark = Theme.of(context).brightness == Brightness.dark;
    
    return SliverAppBar(
      expandedHeight: 120,
      floating: false,
      pinned: true,
      backgroundColor: Colors.transparent,
      elevation: 0,
      leading: Padding(
        padding: const EdgeInsets.all(8.0),
        child: GlassContainer(
          borderRadius: 12,
          blur: 10,
          color: isDark ? Colors.white.withOpacity(0.1) : Colors.white.withOpacity(0.8),
          child: IconButton(
            icon: Icon(Icons.menu_rounded, color: isDark ? Colors.white : AppTheme.textPrimaryColor),
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
            color: isDark ? Colors.white.withOpacity(0.1) : Colors.white.withOpacity(0.8),
            child: Stack(
              children: [
                Center(
                  child: IconButton(
                    icon: Icon(Icons.notifications_outlined, color: isDark ? Colors.white : AppTheme.textPrimaryColor),
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
      flexibleSpace: FlexibleSpaceBar(
        background: Container(
          decoration: BoxDecoration(
            gradient: AppTheme.primaryGradient,
          ),
          child: SafeArea(
            child: Padding(
              padding: const EdgeInsets.fromLTRB(60, 20, 60, 20),
              child: Column(
                mainAxisAlignment: MainAxisAlignment.end,
                crossAxisAlignment: CrossAxisAlignment.center,
                children: [
                  BlocBuilder<AuthBloc, AuthState>(
                    builder: (context, state) {
                      String greeting = 'Bonjour';
                      final hour = DateTime.now().hour;
                      if (hour < 12) {
                        greeting = 'Bonjour ☀️';
                      } else if (hour < 18) {
                        greeting = 'Bon après-midi 👋';
                      } else {
                        greeting = 'Bonsoir 🌙';
                      }
                      
                      String name = 'Utilisateur';
                      if (state is AuthenticatedState) {
                        name = state.user.firstName ?? 'Utilisateur';
                      }
                      
                      return Text(
                        '$greeting, $name!',
                        style: GoogleFonts.inter(
                          color: Colors.white,
                          fontSize: 22,
                          fontWeight: FontWeight.bold,
                        ),
                      );
                    },
                  ),
                ],
              ),
            ),
          ),
        ),
      ),
    );
  }

  Widget _buildBalanceCard() {
    return Container(
      margin: const EdgeInsets.all(20),
      child: GlassContainer(
        gradient: AppTheme.cardGradient,
        padding: const EdgeInsets.all(24),
        borderRadius: 24,
        child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Row(
            mainAxisAlignment: MainAxisAlignment.spaceBetween,
            children: [
              Text(
                'Solde Total',
                style: GoogleFonts.inter(
                  color: Colors.white.withOpacity(0.8),
                  fontSize: 14,
                  fontWeight: FontWeight.w500,
                ),
              ),
              Container(
                padding: const EdgeInsets.symmetric(horizontal: 12, vertical: 4),
                decoration: BoxDecoration(
                  color: Colors.white.withOpacity(0.2),
                  borderRadius: BorderRadius.circular(20),
                ),
                child: Row(
                  children: [
                    Container(
                      width: 6,
                      height: 6,
                      decoration: const BoxDecoration(
                        color: AppTheme.successColor,
                        shape: BoxShape.circle,
                      ),
                    ),
                    const SizedBox(width: 6),
                    Text(
                      'Actif',
                      style: GoogleFonts.inter(
                        color: Colors.white,
                        fontSize: 12,
                        fontWeight: FontWeight.w600,
                      ),
                    ),
                  ],
                ),
              ),
            ],
          ),
          const SizedBox(height: 16),
          BlocBuilder<WalletBloc, WalletState>(
            builder: (context, state) {
              double totalBalance = 0.0;
              if (state is WalletLoadedState) {
                for (var wallet in state.wallets) {
                  totalBalance += wallet.balance;
                }
              }
              return Row(
                crossAxisAlignment: CrossAxisAlignment.end,
                children: [
                  Text(
                    '\$${totalBalance.toStringAsFixed(2)}',
                    style: GoogleFonts.inter(
                      color: Colors.white,
                      fontSize: 36,
                      fontWeight: FontWeight.bold,
                      letterSpacing: -1,
                    ),
                  ),
                  const SizedBox(width: 8),
                  Padding(
                    padding: const EdgeInsets.only(bottom: 6),
                    child: Text(
                      'USD',
                      style: GoogleFonts.inter(
                        color: Colors.white.withOpacity(0.8),
                        fontSize: 16,
                        fontWeight: FontWeight.w500,
                      ),
                    ),
                  ),
                ],
              );
            },
          ),
          const SizedBox(height: 24),
          Row(
            children: [
              Expanded(
                child: _buildBalanceAction(
                  icon: Icons.add_rounded,
                  label: 'Dépôt',
                  onTap: () {},
                ),
              ),
              const SizedBox(width: 12),
              Expanded(
                child: _buildBalanceAction(
                  icon: Icons.qr_code_rounded,
                  label: 'Recevoir',
                  onTap: () {},
                ),
              ),
              const SizedBox(width: 12),
              Expanded(
                child: _buildBalanceAction(
                  icon: Icons.history_rounded,
                  label: 'Historique',
                  onTap: () {},
                ),
              ),
            ],
          ),
        ],
      ),
      ),
    );
  }

  Widget _buildBalanceAction({
    required IconData icon,
    required String label,
    required VoidCallback onTap,
  }) {
    return GestureDetector(
      onTap: onTap,
      child: Container(
        padding: const EdgeInsets.symmetric(vertical: 12),
        decoration: BoxDecoration(
          color: Colors.white.withOpacity(0.1),
          borderRadius: BorderRadius.circular(12),
        ),
        child: Column(
          children: [
            Icon(icon, color: Colors.white70, size: 24),
            const SizedBox(height: 4),
            Text(
              label,
              style: const TextStyle(
                color: Colors.white70,
                fontSize: 12,
              ),
            ),
          ],
        ),
      ),
    );
  }

  Widget _buildQuickActions(BuildContext context) {
    return Padding(
      padding: const EdgeInsets.symmetric(horizontal: 20),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          const Text(
            'Actions Rapides',
            style: TextStyle(
              fontSize: 18,
              fontWeight: FontWeight.bold,
              color: Color(0xFF1a1a2e),
            ),
          ),
          const SizedBox(height: 16),
          Row(
            children: [
              Expanded(
                child: _buildQuickActionCard(
                  icon: Icons.send_rounded,
                  label: 'Envoyer',
                  color: const Color(0xFF667eea),
                  onTap: () => context.push('/more/transfer'),
                ),
              ),
              const SizedBox(width: 12),
              Expanded(
                child: _buildQuickActionCard(
                  icon: Icons.qr_code_scanner_rounded,
                  label: 'Scanner',
                  color: const Color(0xFF10B981),
                  onTap: () => context.push('/more/merchant/scan'),
                ),
              ),
              const SizedBox(width: 12),
              Expanded(
                child: _buildQuickActionCard(
                  icon: Icons.credit_card_rounded,
                  label: 'Cartes',
                  color: const Color(0xFFF59E0B),
                  onTap: () => context.push('/more/cards'),
                ),
              ),
              const SizedBox(width: 12),
              Expanded(
                child: _buildQuickActionCard(
                  icon: Icons.swap_horiz_rounded,
                  label: 'Exchange',
                  color: const Color(0xFF8B5CF6),
                  onTap: () => context.push('/exchange'),
                ),
              ),
            ],
          ),
        ],
      ),
    );
  }

  Widget _buildQuickActionCard({
    required IconData icon,
    required String label,
    required Color color,
    required VoidCallback onTap,
  }) {
    return GestureDetector(
      onTap: onTap,
      child: GlassContainer(
        padding: const EdgeInsets.symmetric(vertical: 16),
        borderRadius: 16,
        child: Column(
          children: [
            Container(
              width: 48,
              height: 48,
              decoration: BoxDecoration(
                color: color.withOpacity(0.1),
                borderRadius: BorderRadius.circular(14),
              ),
              child: Icon(icon, color: color, size: 24),
            ),
            const SizedBox(height: 8),
            Text(
              label,
              style: GoogleFonts.inter(
                fontSize: 12,
                fontWeight: FontWeight.w600,
                color: Theme.of(context).brightness == Brightness.dark 
                    ? Colors.white 
                    : AppTheme.textPrimaryColor,
              ),
            ),
          ],
        ),
      ),
    );
  }

  Widget _buildRecentActivity() {
    return Padding(
      padding: const EdgeInsets.all(20),
      child: Column(
        children: [
          Row(
            mainAxisAlignment: MainAxisAlignment.spaceBetween,
            children: [
              Text(
                'Activité Récente',
                style: GoogleFonts.inter(
                  fontSize: 18,
                  fontWeight: FontWeight.bold,
                  color: Theme.of(context).brightness == Brightness.dark ? Colors.white : AppTheme.textPrimaryColor,
                ),
              ),
              TextButton(
                onPressed: () => context.push('/wallet'),
                child: const Text('Voir tout'),
              ),
            ],
          ),
          const SizedBox(height: 12),
          BlocBuilder<WalletBloc, WalletState>(
            builder: (context, state) {
              if (state is WalletLoadedState) {
                // Show empty state with link to wallet for transactions
                return GlassContainer(
                  padding: const EdgeInsets.all(24),
                  borderRadius: 16,
                  child: Column(
                    children: [
                      const Text('📊', style: TextStyle(fontSize: 40)),
                      const SizedBox(height: 12),
                      Text(
                        'Consultez vos transactions',
                        style: GoogleFonts.inter(
                          color: Theme.of(context).brightness == Brightness.dark ? Colors.white : AppTheme.textPrimaryColor,
                          fontWeight: FontWeight.bold,
                          fontSize: 16,
                        ),
                      ),
                      const SizedBox(height: 8),
                      Text(
                        'Accédez à vos portefeuilles pour voir l\'historique complet',
                        textAlign: TextAlign.center,
                        style: GoogleFonts.inter(color: AppTheme.textSecondaryColor),
                      ),
                      const SizedBox(height: 16),
                      GestureDetector(
                        onTap: () => context.push('/wallet'),
                        child: Container(
                          padding: const EdgeInsets.symmetric(horizontal: 24, vertical: 12),
                          decoration: BoxDecoration(
                            gradient: AppTheme.primaryGradient,
                            borderRadius: BorderRadius.circular(12),
                          ),
                          child: Text(
                            'Voir les portefeuilles',
                            style: GoogleFonts.inter(color: Colors.white, fontWeight: FontWeight.bold),
                          ),
                        ),
                      ),
                    ],
                  ),
                );
              }
              
              // Loading state
              return GlassContainer(
                padding: const EdgeInsets.all(32),
                borderRadius: 16,
                child: const Center(
                  child: CircularProgressIndicator(),
                ),
              );
            },
          ),
        ],
      ),
    );
  }

  Widget _buildRealTransactionItem(dynamic tx, String currency) {
    final isCredit = tx.isIncoming ?? false;
    final amount = tx.amount ?? 0.0;
    final type = tx.type?.toString().split('.').last ?? 'transfer';
    
    IconData icon;
    Color color;
    String title;
    String subtitle;
    
    switch (type) {
      case 'deposit':
      case 'credit':
        icon = Icons.arrow_downward_rounded;
        color = const Color(0xFF10B981);
        title = tx.memo ?? 'Dépôt';
        subtitle = 'Dépôt';
        break;
      case 'withdrawal':
      case 'debit':
        icon = Icons.arrow_upward_rounded;
        color = const Color(0xFFEF4444);
        title = tx.memo ?? 'Retrait';
        subtitle = 'Retrait';
        break;
      case 'transfer':
        if (isCredit) {
          icon = Icons.arrow_downward_rounded;
          color = const Color(0xFF10B981);
          title = tx.senderName ?? tx.memo ?? 'Reçu';
          subtitle = 'Transfert reçu';
        } else {
          icon = Icons.arrow_upward_rounded;
          color = const Color(0xFFEF4444);
          title = tx.recipientName ?? tx.memo ?? 'Envoi';
          subtitle = 'Transfert envoyé';
        }
        break;
      case 'payment':
        icon = Icons.shopping_bag_rounded;
        color = const Color(0xFFF59E0B);
        title = tx.merchantName ?? tx.memo ?? 'Paiement';
        subtitle = 'Paiement';
        break;
      case 'exchange':
        icon = Icons.swap_horiz_rounded;
        color = const Color(0xFF3B82F6);
        title = tx.memo ?? 'Échange';
        subtitle = 'Conversion';
        break;
      default:
        icon = isCredit ? Icons.arrow_downward_rounded : Icons.arrow_upward_rounded;
        color = isCredit ? const Color(0xFF10B981) : const Color(0xFFEF4444);
        title = tx.memo ?? (isCredit ? 'Crédit' : 'Débit');
        subtitle = type;
    }
    
    final amountStr = '${isCredit ? '+' : '-'}${amount.toStringAsFixed(0)} $currency';
    
    return _buildTransactionItem(
      icon: icon,
      color: color,
      title: title,
      subtitle: subtitle,
      amount: amountStr,
      isCredit: isCredit,
    );
  }

  Widget _buildTransactionItem({
    required IconData icon,
    required Color color,
    required String title,
    required String subtitle,
    required String amount,
    required bool isCredit,
  }) {
    return Padding(
      padding: const EdgeInsets.all(16),
      child: Row(
        children: [
          Container(
            width: 48,
            height: 48,
            decoration: BoxDecoration(
              color: color.withOpacity(0.1),
              borderRadius: BorderRadius.circular(12),
            ),
            child: Icon(icon, color: color, size: 24),
          ),
          const SizedBox(width: 16),
          Expanded(
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Text(
                  title,
                  style: const TextStyle(
                    fontWeight: FontWeight.w600,
                    fontSize: 15,
                    color: Color(0xFF1a1a2e),
                  ),
                ),
                Text(
                  subtitle,
                  style: const TextStyle(
                    color: Color(0xFF94A3B8),
                    fontSize: 13,
                  ),
                ),
              ],
            ),
          ),
          Text(
            amount,
            style: TextStyle(
              fontWeight: FontWeight.bold,
              fontSize: 15,
              color: isCredit ? const Color(0xFF10B981) : const Color(0xFF1a1a2e),
            ),
          ),
        ],
      ),
    );
  }

  Widget _buildServicesGrid(BuildContext context) {
    return Padding(
      padding: const EdgeInsets.symmetric(horizontal: 20),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Text(
            'Services',
            style: GoogleFonts.inter(
              fontSize: 18,
              fontWeight: FontWeight.bold,
              color: Theme.of(context).brightness == Brightness.dark ? Colors.white : AppTheme.textPrimaryColor,
            ),
          ),
          const SizedBox(height: 16),
          GridView.count(
            shrinkWrap: true,
            physics: const NeverScrollableScrollPhysics(),
            crossAxisCount: 2,
            crossAxisSpacing: 12,
            mainAxisSpacing: 12,
            childAspectRatio: 1.5,
            children: [
              _buildServiceCard(
                icon: '🏦',
                title: 'Virement Bancaire',
                subtitle: 'SEPA & Swift',
                color: const Color(0xFF3B82F6),
                onTap: () {},
              ),
              _buildServiceCard(
                icon: '📱',
                title: 'Mobile Money',
                subtitle: 'Orange, MTN...',
                color: const Color(0xFFF97316),
                onTap: () {},
              ),
              _buildServiceCard(
                icon: '💳',
                title: 'Carte Virtuelle',
                subtitle: 'Instantané',
                color: const Color(0xFF8B5CF6),
                onTap: () => context.push('/more/cards'),
              ),
              _buildServiceCard(
                icon: '🛒',
                title: 'Paiement QR',
                subtitle: 'Scan & Pay',
                color: const Color(0xFF10B981),
                onTap: () => context.push('/more/merchant/scan'),
              ),
            ],
          ),
        ],
      ),
    );
  }

  Widget _buildServiceCard({
    required String icon,
    required String title,
    required String subtitle,
    required Color color,
    required VoidCallback onTap,
  }) {
    return GestureDetector(
      onTap: onTap,
      child: GlassContainer(
        padding: const EdgeInsets.all(16),
        borderRadius: 16,
        borderColor: color.withOpacity(0.2),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          mainAxisAlignment: MainAxisAlignment.center,
          children: [
            Text(icon, style: const TextStyle(fontSize: 28)),
            const SizedBox(height: 8),
            Text(
              title,
              style: const TextStyle(
                fontWeight: FontWeight.w600,
                fontSize: 14,
                color: Color(0xFF1a1a2e),
              ),
            ),
            Text(
              subtitle,
              style: const TextStyle(
                color: Color(0xFF94A3B8),
                fontSize: 12,
              ),
            ),
          ],
        ),
      ),
    );
  }
}
