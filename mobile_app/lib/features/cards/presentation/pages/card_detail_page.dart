import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';
import 'package:go_router/go_router.dart';
import 'package:google_fonts/google_fonts.dart';
import 'package:intl/intl.dart';

import '../../../../core/widgets/glass_container.dart';
import '../../../../core/widgets/loading_widget.dart';
import '../bloc/cards_bloc.dart';

class CardDetailPage extends StatefulWidget {
  final String cardId;

  const CardDetailPage({super.key, required this.cardId});

  @override
  State<CardDetailPage> createState() => _CardDetailPageState();
}

class _CardDetailPageState extends State<CardDetailPage> {
  Map<String, dynamic>? _card;
  bool _isLoading = true;

  @override
  void initState() {
    super.initState();
    _loadCardDetails();
  }

  void _loadCardDetails() {
    // Get card from the bloc state
    final cardsState = context.read<CardsBloc>().state;
    if (cardsState is CardsLoadedState) {
      final card = cardsState.cards.firstWhere(
        (c) => c['id'] == widget.cardId,
        orElse: () => <String, dynamic>{},
      );
      if (card.isNotEmpty) {
        setState(() {
          _card = card;
          _isLoading = false;
        });
        return;
      }
    }
    // If not found, refresh cards
    context.read<CardsBloc>().add(LoadCardsEvent());
  }

  String _formatCurrency(dynamic amount, String currency) {
    final value = (amount is num) ? amount.toDouble() : double.tryParse(amount?.toString() ?? '0') ?? 0;
    final formatter = NumberFormat.currency(
      locale: 'fr_FR',
      symbol: currency == 'USD' ? '\$' : currency == 'EUR' ? '€' : currency,
      decimalDigits: 2,
    );
    return formatter.format(value);
  }

  String _getCardholderName() {
    if (_card == null) return 'N/A';
    final firstName = _card!['first_name'] ?? _card!['cardholder_name'] ?? '';
    final lastName = _card!['last_name'] ?? '';
    if (firstName.isEmpty && lastName.isEmpty) return 'N/A';
    return '$firstName $lastName'.trim().toUpperCase();
  }

  @override
  Widget build(BuildContext context) {
    final isDark = Theme.of(context).brightness == Brightness.dark;

    return Scaffold(
      backgroundColor: isDark ? const Color(0xFF0F172A) : const Color(0xFFF8FAFC),
      body: BlocListener<CardsBloc, CardsState>(
        listener: (context, state) {
          if (state is CardsLoadedState) {
            final card = state.cards.firstWhere(
              (c) => c['id'] == widget.cardId,
              orElse: () => <String, dynamic>{},
            );
            if (card.isNotEmpty) {
              setState(() {
                _card = card;
                _isLoading = false;
              });
            }
          } else if (state is CardsErrorState) {
            setState(() => _isLoading = false);
          }
        },
        child: SafeArea(
          child: _isLoading
              ? const LoadingWidget()
              : _card == null
                  ? Center(
                      child: Column(
                        mainAxisAlignment: MainAxisAlignment.center,
                        children: [
                          const Icon(Icons.credit_card_off, size: 64, color: Colors.grey),
                          const SizedBox(height: 16),
                          const Text('Carte non trouvée'),
                          const SizedBox(height: 16),
                          ElevatedButton(
                            onPressed: () => context.pop(),
                            child: const Text('Retour'),
                          ),
                        ],
                      ),
                    )
                  : _buildContent(isDark),
        ),
      ),
    );
  }

  Widget _buildContent(bool isDark) {
    final balance = _card!['balance'] ?? _card!['available_balance'] ?? 0;
    final currency = _card!['currency'] ?? 'USD';
    final cardNumber = _card!['card_number'] ?? '**** **** **** ****';
    final expiryMonth = _card!['expiry_month'] ?? 12;
    final expiryYear = _card!['expiry_year'] ?? 28;
    final status = _card!['status'] ?? 'inactive';
    final cardType = _card!['card_type'] ?? 'virtual';
    final isVirtual = _card!['is_virtual'] == true || cardType == 'virtual';

    return SingleChildScrollView(
      padding: const EdgeInsets.all(16),
      child: Column(
        children: [
          // Header
          Row(
            children: [
              GlassContainer(
                padding: EdgeInsets.zero,
                width: 40,
                height: 40,
                borderRadius: 12,
                child: IconButton(
                  icon: Icon(Icons.arrow_back_ios_new, size: 18,
                      color: isDark ? Colors.white : const Color(0xFF1E293B)),
                  onPressed: () => context.pop(),
                ),
              ),
              const SizedBox(width: 16),
              Text(
                'Détails de la carte',
                style: GoogleFonts.inter(
                  fontSize: 20,
                  fontWeight: FontWeight.bold,
                  color: isDark ? Colors.white : const Color(0xFF1E293B),
                ),
              ),
            ],
          ),
          const SizedBox(height: 24),

          // Card visualization
          Container(
            width: double.infinity,
            height: 200,
            decoration: BoxDecoration(
              gradient: LinearGradient(
                colors: isVirtual
                    ? [const Color(0xFF6366F1), const Color(0xFF8B5CF6)]
                    : [const Color(0xFF1E293B), const Color(0xFF334155)],
                begin: Alignment.topLeft,
                end: Alignment.bottomRight,
              ),
              borderRadius: BorderRadius.circular(20),
              boxShadow: [
                BoxShadow(
                  color: (isVirtual ? const Color(0xFF6366F1) : const Color(0xFF1E293B))
                      .withOpacity(0.3),
                  blurRadius: 20,
                  offset: const Offset(0, 10),
                ),
              ],
            ),
            padding: const EdgeInsets.all(24),
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              mainAxisAlignment: MainAxisAlignment.spaceBetween,
              children: [
                Row(
                  mainAxisAlignment: MainAxisAlignment.spaceBetween,
                  children: [
                    Text(
                      'Zekora',
                      style: GoogleFonts.inter(
                        color: Colors.white,
                        fontSize: 18,
                        fontWeight: FontWeight.bold,
                      ),
                    ),
                    Row(
                      children: [
                        if (isVirtual)
                          Container(
                            padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 4),
                            decoration: BoxDecoration(
                              color: Colors.white.withOpacity(0.2),
                              borderRadius: BorderRadius.circular(6),
                            ),
                            child: Text(
                              'VIRTUEL',
                              style: GoogleFonts.inter(
                                color: Colors.white,
                                fontSize: 10,
                                fontWeight: FontWeight.bold,
                              ),
                            ),
                          ),
                        const SizedBox(width: 8),
                        const Icon(Icons.contactless, color: Colors.white, size: 24),
                      ],
                    ),
                  ],
                ),
                Text(
                  cardNumber,
                  style: GoogleFonts.robotoMono(
                    color: Colors.white,
                    fontSize: 22,
                    letterSpacing: 3,
                  ),
                ),
                Row(
                  mainAxisAlignment: MainAxisAlignment.spaceBetween,
                  children: [
                    Column(
                      crossAxisAlignment: CrossAxisAlignment.start,
                      children: [
                        Text(
                          'TITULAIRE',
                          style: GoogleFonts.inter(color: Colors.white70, fontSize: 10),
                        ),
                        Text(
                          _getCardholderName(),
                          style: GoogleFonts.inter(
                            color: Colors.white,
                            fontSize: 12,
                            fontWeight: FontWeight.w600,
                          ),
                        ),
                      ],
                    ),
                    Column(
                      crossAxisAlignment: CrossAxisAlignment.start,
                      children: [
                        Text(
                          'EXPIRE',
                          style: GoogleFonts.inter(color: Colors.white70, fontSize: 10),
                        ),
                        Text(
                          '${expiryMonth.toString().padLeft(2, '0')}/${expiryYear.toString().substring(expiryYear.toString().length - 2)}',
                          style: GoogleFonts.inter(
                            color: Colors.white,
                            fontSize: 12,
                            fontWeight: FontWeight.w600,
                          ),
                        ),
                      ],
                    ),
                    Container(
                      padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 4),
                      decoration: BoxDecoration(
                        color: status == 'active'
                            ? const Color(0xFF22C55E).withOpacity(0.2)
                            : const Color(0xFFEF4444).withOpacity(0.2),
                        borderRadius: BorderRadius.circular(6),
                      ),
                      child: Text(
                        status == 'active' ? 'ACTIF' : status.toUpperCase(),
                        style: GoogleFonts.inter(
                          color: status == 'active' ? const Color(0xFF22C55E) : const Color(0xFFEF4444),
                          fontSize: 10,
                          fontWeight: FontWeight.bold,
                        ),
                      ),
                    ),
                  ],
                ),
              ],
            ),
          ),
          const SizedBox(height: 24),

          // Balance
          Container(
            width: double.infinity,
            padding: const EdgeInsets.all(20),
            decoration: BoxDecoration(
              color: isDark ? Colors.white.withOpacity(0.05) : Colors.white,
              borderRadius: BorderRadius.circular(16),
              border: Border.all(
                color: isDark ? Colors.white.withOpacity(0.1) : const Color(0xFFE2E8F0),
              ),
            ),
            child: Column(
              children: [
                Text(
                  'Solde disponible',
                  style: GoogleFonts.inter(
                    color: isDark ? Colors.white70 : const Color(0xFF64748B),
                  ),
                ),
                const SizedBox(height: 8),
                Text(
                  _formatCurrency(balance, currency),
                  style: GoogleFonts.inter(
                    fontSize: 32,
                    fontWeight: FontWeight.bold,
                    color: isDark ? Colors.white : const Color(0xFF1E293B),
                  ),
                ),
              ],
            ),
          ),
          const SizedBox(height: 24),

          // Actions
          Row(
            children: [
              Expanded(
                child: _ActionCard(
                  icon: Icons.add,
                  label: 'Recharger',
                  onTap: () => _showLoadDialog(),
                  isDark: isDark,
                ),
              ),
              const SizedBox(width: 12),
              Expanded(
                child: _ActionCard(
                  icon: status == 'active' ? Icons.lock : Icons.lock_open,
                  label: status == 'active' ? 'Bloquer' : 'Débloquer',
                  onTap: () => _toggleCardStatus(),
                  isDark: isDark,
                ),
              ),
              const SizedBox(width: 12),
              Expanded(
                child: _ActionCard(
                  icon: Icons.visibility,
                  label: 'Voir PIN',
                  onTap: () => _showPinDialog(),
                  isDark: isDark,
                ),
              ),
            ],
          ),
          const SizedBox(height: 24),

          // Settings
          Align(
            alignment: Alignment.centerLeft,
            child: Text(
              'Paramètres',
              style: GoogleFonts.inter(
                fontSize: 18,
                fontWeight: FontWeight.bold,
                color: isDark ? Colors.white : const Color(0xFF1E293B),
              ),
            ),
          ),
          const SizedBox(height: 16),

          _SettingsTile(
            icon: Icons.shopping_bag,
            title: 'Limites de dépenses',
            subtitle: _formatCurrency(_card!['daily_limit'] ?? 1000, currency) + '/jour',
            isDark: isDark,
          ),
          _SettingsTile(
            icon: Icons.language,
            title: 'Paiements en ligne',
            subtitle: _card!['allow_online'] == true ? 'Activé' : 'Désactivé',
            isDark: isDark,
          ),
          _SettingsTile(
            icon: Icons.atm,
            title: 'Retraits DAB',
            subtitle: _card!['allow_atm'] == true ? 'Activé' : 'Désactivé',
            isDark: isDark,
          ),
          _SettingsTile(
            icon: Icons.public,
            title: 'Paiements internationaux',
            subtitle: _card!['allow_international'] == true ? 'Activé' : 'Désactivé',
            isDark: isDark,
          ),
        ],
      ),
    );
  }

  void _showLoadDialog() {
    // TODO: Implement load card dialog
    ScaffoldMessenger.of(context).showSnackBar(
      const SnackBar(content: Text('Fonctionnalité de rechargement à venir')),
    );
  }

  void _toggleCardStatus() {
    final status = _card!['status'];
    if (status == 'active') {
      context.read<CardsBloc>().add(FreezeCardEvent(widget.cardId));
    } else {
      context.read<CardsBloc>().add(UnfreezeCardEvent(widget.cardId));
    }
  }

  void _showPinDialog() {
    // TODO: Implement PIN reveal with authentication
    ScaffoldMessenger.of(context).showSnackBar(
      const SnackBar(content: Text('Authentification requise pour voir le PIN')),
    );
  }
}

class _ActionCard extends StatelessWidget {
  final IconData icon;
  final String label;
  final VoidCallback onTap;
  final bool isDark;

  const _ActionCard({
    required this.icon,
    required this.label,
    required this.onTap,
    required this.isDark,
  });

  @override
  Widget build(BuildContext context) {
    return InkWell(
      onTap: onTap,
      borderRadius: BorderRadius.circular(12),
      child: Container(
        padding: const EdgeInsets.symmetric(vertical: 16),
        decoration: BoxDecoration(
          color: isDark ? Colors.white.withOpacity(0.05) : Colors.white,
          border: Border.all(
            color: isDark ? Colors.white.withOpacity(0.1) : const Color(0xFFE2E8F0),
          ),
          borderRadius: BorderRadius.circular(12),
        ),
        child: Column(
          children: [
            Icon(icon, color: const Color(0xFF6366F1)),
            const SizedBox(height: 8),
            Text(
              label,
              style: GoogleFonts.inter(
                fontSize: 11,
                color: isDark ? Colors.white70 : const Color(0xFF64748B),
              ),
            ),
          ],
        ),
      ),
    );
  }
}

class _SettingsTile extends StatelessWidget {
  final IconData icon;
  final String title;
  final String subtitle;
  final bool isDark;

  const _SettingsTile({
    required this.icon,
    required this.title,
    required this.subtitle,
    required this.isDark,
  });

  @override
  Widget build(BuildContext context) {
    return Container(
      margin: const EdgeInsets.only(bottom: 8),
      decoration: BoxDecoration(
        color: isDark ? Colors.white.withOpacity(0.03) : Colors.white,
        borderRadius: BorderRadius.circular(12),
        border: Border.all(
          color: isDark ? Colors.white.withOpacity(0.05) : const Color(0xFFE2E8F0),
        ),
      ),
      child: ListTile(
        leading: Container(
          padding: const EdgeInsets.all(8),
          decoration: BoxDecoration(
            color: const Color(0xFF6366F1).withOpacity(0.1),
            borderRadius: BorderRadius.circular(8),
          ),
          child: Icon(icon, size: 20, color: const Color(0xFF6366F1)),
        ),
        title: Text(
          title,
          style: GoogleFonts.inter(
            color: isDark ? Colors.white : const Color(0xFF1E293B),
          ),
        ),
        subtitle: Text(
          subtitle,
          style: GoogleFonts.inter(
            color: isDark ? Colors.white54 : const Color(0xFF64748B),
          ),
        ),
        trailing: Icon(
          Icons.chevron_right,
          color: isDark ? Colors.white30 : Colors.grey,
        ),
        onTap: () {},
      ),
    );
  }
}
