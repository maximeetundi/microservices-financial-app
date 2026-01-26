import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';
import 'package:flutter_bloc/flutter_bloc.dart';
import 'package:flutter/services.dart';
import '../../../../core/theme/app_theme.dart';
import '../../../../core/widgets/glass_container.dart';
import 'package:google_fonts/google_fonts.dart';
import 'package:pretty_qr_code/pretty_qr_code.dart';
import '../../domain/entities/wallet.dart';
import '../bloc/wallet_bloc.dart';

class WalletDetailPage extends StatefulWidget {
  final String walletId;
  
  const WalletDetailPage({super.key, required this.walletId});

  @override
  State<WalletDetailPage> createState() => _WalletDetailPageState();
}

class _WalletDetailPageState extends State<WalletDetailPage> {
  @override
  void initState() {
    super.initState();
    context.read<WalletBloc>().add(LoadWalletTransactionsEvent(walletId: widget.walletId));
  }

  @override
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
               // Custom App Bar
               Padding(
                 padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 8),
                 child: Row(
                   mainAxisAlignment: MainAxisAlignment.spaceBetween,
                   children: [
                     GlassContainer(
                       padding: EdgeInsets.zero,
                       width: 40,
                       height: 40, 
                       borderRadius: 12,
                       child: IconButton(
                        icon: Icon(Icons.arrow_back_ios_new, size: 20, color: isDark ? Colors.white : AppTheme.textPrimaryColor),
                        onPressed: () => context.pop(),
                      ),
                     ),
                     Text(
                        'Détails du portefeuille',
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
                        icon: Icon(Icons.more_horiz, size: 24, color: isDark ? Colors.white : AppTheme.textPrimaryColor),
                        onPressed: () {},
                      ),
                     ),
                   ],
                 ),
               ),
              Expanded(
                child: BlocBuilder<WalletBloc, WalletState>(
                  builder: (context, state) {
                    if (state is WalletLoadingState) {
                      return const Center(child: CircularProgressIndicator());
                    }
                    
                    if (state is WalletErrorState) {
                      return Center(child: Text('Erreur: ${state.message}', style: TextStyle(color: AppTheme.errorColor)));
                    }
                    
                    if (state is WalletLoadedState) {
                      final wallet = state.wallets.firstWhere(
                        (w) => w.id == widget.walletId,
                        orElse: () => state.wallets.first,
                      );
                      
                      return SingleChildScrollView(
                        padding: const EdgeInsets.all(16),
                        child: Column(
                          crossAxisAlignment: CrossAxisAlignment.start,
                          children: [
                            // Balance Card
                            Container(
                              width: double.infinity,
                              padding: const EdgeInsets.all(24),
                              decoration: BoxDecoration(
                                gradient: AppTheme.primaryGradient,
                                borderRadius: BorderRadius.circular(24),
                                boxShadow: [
                                  BoxShadow(
                                    color: AppTheme.primaryColor.withOpacity(0.3),
                                    blurRadius: 20,
                                    offset: const Offset(0, 10),
                                  ),
                                ],
                                border: Border.all(
                                  color: Colors.white.withOpacity(0.1),
                                  width: 1,
                                ),
                              ),
                              child: Column(
                                crossAxisAlignment: CrossAxisAlignment.start,
                                children: [
                                  Row(
                                    children: [
                                      Container(
                                        padding: const EdgeInsets.all(8),
                                        decoration: BoxDecoration(
                                          color: Colors.white.withOpacity(0.1),
                                          shape: BoxShape.circle,
                                        ),
                                        child: Text(
                                          wallet.currency.substring(0, 1),
                                          style: GoogleFonts.inter(fontSize: 16, fontWeight: FontWeight.bold, color: Colors.white),
                                        ),
                                      ),
                                      const SizedBox(width: 12),
                                      Text(
                                        wallet.currency,
                                        style: GoogleFonts.inter(color: Colors.white70, fontSize: 16, fontWeight: FontWeight.w500),
                                      ),
                                    ],
                                  ),
                                  const SizedBox(height: 24),
                                  Text(
                                    '${wallet.balance.toStringAsFixed(2)} ${wallet.currency}',
                                    style: GoogleFonts.inter(
                                      color: Colors.white,
                                      fontSize: 32,
                                      fontWeight: FontWeight.bold,
                                      letterSpacing: -0.5,
                                    ),
                                  ),
                                  const SizedBox(height: 8),
                                  Text(
                                    'Solde Disponible',
                                    style: GoogleFonts.inter(
                                      color: Colors.white54,
                                      fontSize: 14,
                                    ),
                                  ),
                                ],
                              ),
                            ),
                            const SizedBox(height: 24),
                            
                            // Actions
                            Row(
                              mainAxisAlignment: MainAxisAlignment.spaceEvenly,
                              children: [
                                Expanded(
                                  child: _ActionButton(
                                    icon: Icons.arrow_upward_rounded,
                                    label: 'Envoyer',
                                    onTap: () => context.push('/transfer'),
                                  ),
                                ),
                                const SizedBox(width: 12),
                                Expanded(
                                  child: _ActionButton(
                                    icon: Icons.arrow_downward_rounded,
                                    label: 'Recevoir',
                                    onTap: () => _showReceiveDialog(context, wallet),
                                  ),
                                ),
                                const SizedBox(width: 12),
                                Expanded(
                                  child: _ActionButton(
                                    icon: Icons.swap_horiz_rounded,
                                    label: 'Échanger',
                                    onTap: () => context.push('/exchange'),
                                  ),
                                ),
                              ],
                            ),
                            const SizedBox(height: 32),
                            
                            // Transactions
                            Text(
                              'Transactions récentes',
                              style: GoogleFonts.inter(
                                fontSize: 18, 
                                fontWeight: FontWeight.bold,
                                color: isDark ? Colors.white : AppTheme.textPrimaryColor
                              ),
                            ),
                            const SizedBox(height: 16),
                            
                            if (state.recentTransactions.isEmpty)
                              Center(
                                child: Padding(
                                  padding: const EdgeInsets.all(32),
                                  child: Column(
                                    children: [
                                      Icon(Icons.history, size: 48, color: isDark ? Colors.white24 : Colors.black12),
                                      const SizedBox(height: 16),
                                      Text(
                                        'Aucune transaction', 
                                        style: GoogleFonts.inter(color: isDark ? Colors.white54 : AppTheme.textSecondaryColor)
                                      ),
                                    ],
                                  ),
                                ),
                              )
                            else
                              ListView.separated(
                                shrinkWrap: true,
                                physics: const NeverScrollableScrollPhysics(),
                                itemCount: state.recentTransactions.length,
                                separatorBuilder: (context, index) => const SizedBox(height: 12),
                                itemBuilder: (context, index) {
                                  final tx = state.recentTransactions[index];
                                  final isIncoming = tx.isIncoming;
                                  
                                  return GlassContainer(
                                    padding: const EdgeInsets.all(16),
                                    borderRadius: 16,
                                    child: Row(
                                      children: [
                                        Container(
                                          padding: const EdgeInsets.all(10),
                                          decoration: BoxDecoration(
                                            color: isIncoming ? const Color(0xFF10B981).withOpacity(0.1) : const Color(0xFFEF4444).withOpacity(0.1),
                                            shape: BoxShape.circle,
                                          ),
                                          child: Icon(
                                            isIncoming ? Icons.arrow_downward : Icons.arrow_upward,
                                            color: isIncoming ? const Color(0xFF10B981) : const Color(0xFFEF4444),
                                            size: 20,
                                          ),
                                        ),
                                        const SizedBox(width: 16),
                                        Expanded(
                                          child: Column(
                                            crossAxisAlignment: CrossAxisAlignment.start,
                                            children: [
                                              Text(
                                                tx.memo ?? tx.type.name,
                                                style: GoogleFonts.inter(
                                                  fontWeight: FontWeight.w600,
                                                  color: isDark ? Colors.white : AppTheme.textPrimaryColor,
                                                  fontSize: 16,
                                                ),
                                              ),
                                              Text(
                                                tx.createdAt.toString().substring(0, 16),
                                                style: GoogleFonts.inter(
                                                  color: isDark ? Colors.white54 : AppTheme.textSecondaryColor,
                                                  fontSize: 12,
                                                ),
                                              ),
                                            ],
                                          ),
                                        ),
                                        Text(
                                          '${isIncoming ? '+' : '-'}${tx.amount.toStringAsFixed(2)}',
                                          style: GoogleFonts.inter(
                                            color: isIncoming ? const Color(0xFF10B981) : const Color(0xFFEF4444),
                                            fontWeight: FontWeight.bold,
                                            fontSize: 16,
                                          ),
                                        ),
                                      ],
                                    ),
                                  );
                                },
                              ),
                          ],
                        ),
                      );
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

  void _showReceiveDialog(BuildContext context, Wallet wallet) {
    showModalBottomSheet(
      context: context,
      backgroundColor: Colors.transparent,
      builder: (context) {
        final isDark = Theme.of(context).brightness == Brightness.dark;
        return Container(
          decoration: BoxDecoration(
            color: isDark ? const Color(0xFF1E293B) : Colors.white,
            borderRadius: const BorderRadius.vertical(top: Radius.circular(24)),
          ),
          padding: const EdgeInsets.all(24),
          child: Column(
            mainAxisSize: MainAxisSize.min,
            children: [
              Container(
                width: 40,
                height: 4,
                decoration: BoxDecoration(
                  color: Colors.grey.withOpacity(0.3),
                  borderRadius: BorderRadius.circular(2),
                ),
              ),
              const SizedBox(height: 24),
              Text(
                wallet.isCrypto ? 'Adresse de dépôt ${wallet.currency}' : 'Recharger votre compte',
                style: GoogleFonts.inter(
                  fontSize: 20,
                  fontWeight: FontWeight.bold,
                  color: isDark ? Colors.white : AppTheme.textPrimaryColor,
                ),
              ),
              const SizedBox(height: 16),
              
              if (wallet.isCrypto) ...[
                Container(
                  padding: const EdgeInsets.all(16),
                  decoration: BoxDecoration(
                    color: Colors.white,
                    borderRadius: BorderRadius.circular(16),
                    border: Border.all(color: Colors.grey.shade200),
                  ),
                  child: PrettyQr(
                    data: wallet.address,
                    size: 200,
                    roundEdges: true,
                    elementColor: Colors.black,
                  ),
                ),
                const SizedBox(height: 24),
                Text(
                  'Adresse du portefeuille',
                  style: GoogleFonts.inter(
                    color: isDark ? Colors.white70 : AppTheme.textSecondaryColor,
                    fontSize: 14,
                  ),
                ),
                const SizedBox(height: 8),
                Container(
                  padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 12),
                  decoration: BoxDecoration(
                    color: isDark ? const Color(0xFF0F172A) : Colors.grey.shade100,
                    borderRadius: BorderRadius.circular(12),
                  ),
                  child: Row(
                    children: [
                      Expanded(
                        child: Text(
                          wallet.address,
                          style: GoogleFonts.sourceCodePro(
                            color: isDark ? Colors.white : AppTheme.textPrimaryColor,
                            fontWeight: FontWeight.w600,
                          ),
                          maxLines: 1,
                          overflow: TextOverflow.ellipsis,
                        ),
                      ),
                      IconButton(
                        icon: const Icon(Icons.copy, size: 20, color: AppTheme.primaryColor),
                        onPressed: () {
                          Clipboard.setData(ClipboardData(text: wallet.address));
                          Navigator.pop(context);
                          ScaffoldMessenger.of(context).showSnackBar(
                            const SnackBar(content: Text('Adresse copiée !')),
                          );
                        },
                      ),
                    ],
                  ),
                ),
                const SizedBox(height: 16),
                Container(
                  padding: const EdgeInsets.all(12),
                  decoration: BoxDecoration(
                    color: Colors.amber.withOpacity(0.1),
                    borderRadius: BorderRadius.circular(12),
                    border: Border.all(color: Colors.amber.withOpacity(0.3)),
                  ),
                  child: Row(
                    children: [
                      const Icon(Icons.warning_amber_rounded, color: Colors.amber, size: 20),
                      const SizedBox(width: 12),
                      Expanded(
                        child: Text(
                          'Envoyez uniquement du ${wallet.currency} sur cette adresse. Tout autre jeton sera perdu.',
                          style: GoogleFonts.inter(
                            color: Colors.amber.shade700,
                            fontSize: 12,
                          ),
                        ),
                      ),
                    ],
                  ),
                ),
              ] else ...[
                // Fiat Deposit Options
                _buildDepositOption(
                  context,
                  icon: Icons.credit_card,
                  title: 'Carte Bancaire',
                  subtitle: 'Visa, Mastercard',
                  onTap: () {
                    // TODO: Integrate Stripe/Card Service
                    Navigator.pop(context);
                    ScaffoldMessenger.of(context).showSnackBar(
                      const SnackBar(content: Text('Bientôt disponible')),
                    );
                  },
                ),
                const SizedBox(height: 12),
                _buildDepositOption(
                  context,
                  icon: Icons.phone_android,
                  title: 'Mobile Money',
                  subtitle: 'Orange, MTN, Wave',
                  onTap: () {
                    // TODO: Integrate Mobile Money
                    Navigator.pop(context);
                    ScaffoldMessenger.of(context).showSnackBar(
                      const SnackBar(content: Text('Bientôt disponible')),
                    );
                  },
                ),
                const SizedBox(height: 12),
                 _buildDepositOption(
                  context,
                  icon: Icons.account_balance,
                  title: 'Virement Bancaire',
                  subtitle: 'IBAN / SWIFT',
                  onTap: () {
                     Navigator.pop(context);
                     ScaffoldMessenger.of(context).showSnackBar(
                      const SnackBar(content: Text('Veuillez contacter le support pour les virements')),
                    );
                  },
                ),
              ],
              
              const SizedBox(height: 24),
              SizedBox(
                width: double.infinity,
                child: ElevatedButton(
                  onPressed: () => Navigator.pop(context),
                  style: ElevatedButton.styleFrom(
                    backgroundColor: isDark ? Colors.white10 : Colors.grey.shade200,
                    foregroundColor: isDark ? Colors.white : Colors.black87,
                    padding: const EdgeInsets.symmetric(vertical: 16),
                    shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(12)),
                    elevation: 0,
                  ),
                  child: const Text('Fermer'),
                ),
              ),
            ],
          ),
        );
      },
    );
  }

  Widget _buildDepositOption(BuildContext context, {
    required IconData icon,
    required String title,
    required String subtitle,
    required VoidCallback onTap,
  }) {
    final isDark = Theme.of(context).brightness == Brightness.dark;
    return InkWell(
      onTap: onTap,
      borderRadius: BorderRadius.circular(16),
      child: Container(
        padding: const EdgeInsets.all(16),
        decoration: BoxDecoration(
          border: Border.all(color: isDark ? Colors.white12 : Colors.grey.shade200),
          borderRadius: BorderRadius.circular(16),
        ),
        child: Row(
          children: [
            Container(
              padding: const EdgeInsets.all(10),
              decoration: BoxDecoration(
                color: AppTheme.primaryColor.withOpacity(0.1),
                shape: BoxShape.circle,
              ),
              child: Icon(icon, color: AppTheme.primaryColor),
            ),
            const SizedBox(width: 16),
            Expanded(
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  Text(
                    title,
                    style: GoogleFonts.inter(
                      fontWeight: FontWeight.bold,
                      color: isDark ? Colors.white : AppTheme.textPrimaryColor,
                    ),
                  ),
                  Text(
                    subtitle,
                    style: GoogleFonts.inter(
                      fontSize: 12,
                      color: isDark ? Colors.white54 : AppTheme.textSecondaryColor,
                    ),
                  ),
                ],
              ),
            ),
            Icon(Icons.chevron_right, color: isDark ? Colors.white24 : Colors.black12),
          ],
        ),
      ),
    );
  }
}

class _ActionButton extends StatelessWidget {
  final IconData icon;
  final String label;
  final VoidCallback onTap;

  const _ActionButton({required this.icon, required this.label, required this.onTap});

  @override
  Widget build(BuildContext context) {
    final isDark = Theme.of(context).brightness == Brightness.dark;
    
    return GlassContainer(
      padding: EdgeInsets.zero,
      borderRadius: 16,
      child: InkWell(
        onTap: onTap,
        borderRadius: BorderRadius.circular(16),
        child: Padding(
          padding: const EdgeInsets.symmetric(vertical: 16),
          child: Column(
            children: [
              Icon(icon, color: isDark ? Colors.white : AppTheme.primaryColor),
              const SizedBox(height: 8),
              Text(
                label, 
                style: GoogleFonts.inter(
                  fontSize: 12, 
                  fontWeight: FontWeight.w500,
                  color: isDark ? Colors.white70 : AppTheme.textSecondaryColor
                ),
              ),
            ],
          ),
        ),
      ),
    );
  }
}
