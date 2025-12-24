import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';
import 'package:go_router/go_router.dart';

import '../../../../core/theme/app_theme.dart';
import '../../../../core/widgets/custom_button.dart';
import '../../../../core/widgets/loading_widget.dart';
import '../../../../core/widgets/glass_container.dart';
import 'package:google_fonts/google_fonts.dart';
import '../bloc/cards_bloc.dart';
import '../widgets/card_widget.dart';
import '../widgets/card_stats_section.dart';
import '../widgets/recent_card_transactions.dart';
import '../widgets/order_card_bottom_sheet.dart';

class CardsPage extends StatefulWidget {
  const CardsPage({Key? key}) : super(key: key);

  @override
  State<CardsPage> createState() => _CardsPageState();
}

class _CardsPageState extends State<CardsPage> {
  final PageController _pageController = PageController();
  int _currentCardIndex = 0;

  @override
  void initState() {
    super.initState();
    _loadCards();
  }

  void _loadCards() {
    context.read<CardsBloc>().add(LoadCardsEvent());
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
               // Custom App Bar Area
               Padding(
                 padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 8),
                 child: Row(
                   mainAxisAlignment: MainAxisAlignment.spaceBetween,
                   children: [
                     const SizedBox(width: 48), // Spacer
                     Text(
                        'Mes Cartes 💳',
                         style: GoogleFonts.inter(
                           fontSize: 20,
                           fontWeight: FontWeight.bold,
                           color: isDark ? Colors.white : AppTheme.textPrimaryColor,
                         ),
                      ),
                     GlassContainer(
                       padding: EdgeInsets.zero,
                       width: 48, 
                       height: 48,
                       borderRadius: 12,
                       child: IconButton(
                        icon: Icon(Icons.add_card, size: 24, color: isDark ? Colors.white : AppTheme.textPrimaryColor),
                        onPressed: _showOrderCardBottomSheet,
                      ),
                     ),
                   ],
                 ),
               ),
              Expanded(
                child: BlocBuilder<CardsBloc, CardsState>(
                  builder: (context, state) {
                    if (state is CardsLoadingState) {
                      return const LoadingWidget();
                    } else if (state is CardsLoadedState) {
                      return _buildCardsContent(state);
                    } else if (state is CardsErrorState) {
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

  Widget _buildCardsContent(CardsLoadedState state) {
    final cards = state.cards;
    
    if (cards.isEmpty) {
      return _buildEmptyCardsState();
    }

    return RefreshIndicator(
      onRefresh: () async {
        _loadCards();
      },
      child: SingleChildScrollView(
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            // Cards Carousel
            SizedBox(
              height: 220,
              child: PageView.builder(
                controller: _pageController,
                onPageChanged: (index) {
                  setState(() {
                    _currentCardIndex = index;
                  });
                },
                itemCount: cards.length,
                itemBuilder: (context, index) {
                  final card = cards[index];
                  return Padding(
                    padding: const EdgeInsets.symmetric(horizontal: 16),
                    child: CardWidget(
                      card: card,
                      onTap: () => context.push('/more/cards/${card['id']}'),
                    ),
                  );
                },
              ),
            ),

            // Page Indicators
            if (cards.length > 1)
              Container(
                height: 20,
                margin: const EdgeInsets.only(top: 16),
                child: Row(
                  mainAxisAlignment: MainAxisAlignment.center,
                  children: List.generate(
                    cards.length,
                    (index) => AnimatedContainer(
                      duration: const Duration(milliseconds: 300),
                      margin: const EdgeInsets.symmetric(horizontal: 4),
                      width: _currentCardIndex == index ? 24 : 8,
                      height: 8,
                      decoration: BoxDecoration(
                        color: _currentCardIndex == index
                            ? AppTheme.primaryColor
                            : Colors.grey.shade300,
                        borderRadius: BorderRadius.circular(4),
                      ),
                    ),
                  ),
                ),
              ),

            const SizedBox(height: 24),

            // Card Actions
            Padding(
              padding: const EdgeInsets.symmetric(horizontal: 16),
              child: Row(
                children: [
                  Expanded(
                    child: _buildActionButton(
                      icon: Icons.add,
                      label: 'Recharger',
                      onTap: () => _showTopUpDialog(cards[_currentCardIndex]),
                    ),
                  ),
                  const SizedBox(width: 12),
                  Expanded(
                    child: _buildActionButton(
                      icon: Icons.ac_unit,
                      label: 'Geler',
                      onTap: () => _freezeCard(cards[_currentCardIndex]),
                    ),
                  ),
                  const SizedBox(width: 12),
                  Expanded(
                    child: _buildActionButton(
                      icon: Icons.settings,
                      label: 'Réglages',
                      onTap: () => context.push('/more/cards/${cards[_currentCardIndex]['id']}'),
                    ),
                  ),
                ],
              ),
            ),

            const SizedBox(height: 32),

            // Card Statistics
            Padding(
              padding: const EdgeInsets.symmetric(horizontal: 16),
              child: CardStatsSection(
                monthlySpend: (cards[_currentCardIndex]['monthly_spend'] as num?)?.toDouble() ?? 0.0,
                limit: (cards[_currentCardIndex]['limit'] as num?)?.toDouble() ?? 1000.0,
                transactionCount: (cards[_currentCardIndex]['transaction_count'] as num?)?.toInt() ?? 0,
              ),
            ),

            const SizedBox(height: 32),

            // Recent Transactions
            Padding(
              padding: const EdgeInsets.symmetric(horizontal: 16),
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
                          color: Theme.of(context).brightness == Brightness.dark ? Colors.white : AppTheme.textPrimaryColor
                        ),
                      ),
                      TextButton(
                        onPressed: () {
                          // Navigate to all transactions
                        },
                        child: Text(
                          'Voir tout',
                          style: GoogleFonts.inter(color: AppTheme.primaryColor),
                        ),
                      ),
                    ],
                  ),
                  const SizedBox(height: 16),
                  RecentCardTransactions(
                    cardId: cards[_currentCardIndex]['id']?.toString() ?? '',
                    transactions: const []
                  ),
                ],
              ),
            ),

            const SizedBox(height: 100), // Bottom padding
          ],
        ),
      ),
    );
  }

  Widget _buildEmptyCardsState() {
    return Center(
      child: Padding(
        padding: const EdgeInsets.all(32),
        child: Column(
          mainAxisAlignment: MainAxisAlignment.center,
          children: [
            Container(
              width: 120,
              height: 80,
              decoration: BoxDecoration(
                color: Colors.grey.shade100,
                borderRadius: BorderRadius.circular(12),
                border: Border.all(
                  color: Colors.grey.shade300,
                  style: BorderStyle.solid,
                ),
              ),
              child: Icon(
                Icons.credit_card,
                size: 40,
                color: Colors.grey.shade400,
              ),
            ),
            const SizedBox(height: 24),
            Text(
              'No Cards Yet',
              style: Theme.of(context).textTheme.headlineMedium?.copyWith(
                color: Colors.grey.shade700,
              ),
            ),
            const SizedBox(height: 8),
            Text(
              'Order your first crypto card to start spending your digital assets anywhere',
              textAlign: TextAlign.center,
              style: Theme.of(context).textTheme.bodyMedium?.copyWith(
                color: Colors.grey.shade600,
              ),
            ),
            const SizedBox(height: 32),
            
            // Card Type Options
            _buildCardTypeOption(
              title: 'Virtual Card',
              description: 'Instant card for online payments',
              price: 'Free',
              color: AppTheme.primaryColor,
              onTap: () => _orderCard('virtual'),
            ),
            const SizedBox(height: 12),
            _buildCardTypeOption(
              title: 'Physical Card',
              description: 'Premium metal card for everywhere',
              price: '\$15',
              color: AppTheme.secondaryColor,
              onTap: () => _orderCard('physical'),
            ),
            const SizedBox(height: 12),
            _buildCardTypeOption(
              title: 'Premium Card',
              description: 'Exclusive benefits & cashback',
              price: '\$99',
              color: Colors.amber.shade700,
              onTap: () => _orderCard('premium'),
            ),
          ],
        ),
      ),
    );
  }

  Widget _buildCardTypeOption({
    required String title,
    required String description,
    required String price,
    required Color color,
    required VoidCallback onTap,
  }) {
    return GestureDetector(
      onTap: onTap,
      child: Container(
        padding: const EdgeInsets.all(16),
        decoration: BoxDecoration(
          color: Colors.white,
          borderRadius: BorderRadius.circular(12),
          border: Border.all(color: Colors.grey.shade200),
        ),
        child: Row(
          children: [
            Container(
              width: 48,
              height: 32,
              decoration: BoxDecoration(
                color: color,
                borderRadius: BorderRadius.circular(6),
              ),
              child: const Icon(
                Icons.credit_card,
                color: Colors.white,
                size: 20,
              ),
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
                      fontSize: 16,
                    ),
                  ),
                  Text(
                    description,
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
                  price,
                  style: const TextStyle(
                    fontWeight: FontWeight.bold,
                    fontSize: 16,
                  ),
                ),
                const Text(
                  'one-time',
                  style: TextStyle(
                    fontSize: 12,
                    color: Colors.grey,
                  ),
                ),
              ],
            ),
          ],
        ),
      ),
    );
  }

  Widget _buildActionButton({
    required IconData icon,
    required String label,
    required VoidCallback onTap,
  }) {
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
              Icon(
                icon,
                color: isDark ? Colors.white : AppTheme.primaryColor,
                size: 24,
              ),
              const SizedBox(height: 8),
              Text(
                label,
                style: GoogleFonts.inter(
                  fontSize: 14,
                  fontWeight: FontWeight.w500,
                  color: isDark ? Colors.white70 : AppTheme.textSecondaryColor,
                ),
              ),
            ],
          ),
        ),
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
            onPressed: _loadCards,
          ),
        ],
      ),
    );
  }

  void _showOrderCardBottomSheet() {
    showModalBottomSheet(
      context: context,
      isScrollControlled: true,
      backgroundColor: Colors.transparent,
      builder: (context) => const OrderCardBottomSheet(),
    );
  }

  void _orderCard(String cardType) {
    showModalBottomSheet(
      context: context,
      isScrollControlled: true,
      backgroundColor: Colors.transparent,
      builder: (context) => OrderCardBottomSheet(
        onOrder: () => context.read<CardsBloc>().add(CreateVirtualCardEvent(currency: 'USD')),
      ),
    );
  }

  void _showTopUpDialog(dynamic card) {
    showDialog(
      context: context,
      builder: (context) => AlertDialog(
        title: const Text('Top Up Card'),
        content: Column(
          mainAxisSize: MainAxisSize.min,
          children: [
            Text('Top up your ${card.type} card'),
            const SizedBox(height: 16),
            TextFormField(
              keyboardType: const TextInputType.numberWithOptions(decimal: true),
              decoration: const InputDecoration(
                labelText: 'Amount',
                prefixText: '\$ ',
                border: OutlineInputBorder(),
              ),
            ),
          ],
        ),
        actions: [
          TextButton(
            onPressed: () => Navigator.pop(context),
            child: const Text('Cancel'),
          ),
          ElevatedButton(
            onPressed: () {
              Navigator.pop(context);
              // Handle top up
              context.read<CardsBloc>().add(TopUpCardEvent(
                cardId: card.id,
                amount: 100.0, // Get from form
              ));
            },
            child: const Text('Top Up'),
          ),
        ],
      ),
    );
  }

  void _freezeCard(dynamic card) {
    showDialog(
      context: context,
      builder: (context) => AlertDialog(
        title: const Text('Freeze Card'),
        content: Text('Are you sure you want to freeze your ${card.type} card?'),
        actions: [
          TextButton(
            onPressed: () => Navigator.pop(context),
            child: const Text('Cancel'),
          ),
          ElevatedButton(
            onPressed: () {
              Navigator.pop(context);
              context.read<CardsBloc>().add(FreezeCardEvent(cardId: card.id));
            },
            style: ElevatedButton.styleFrom(
              backgroundColor: AppTheme.errorColor,
            ),
            child: const Text('Freeze'),
          ),
        ],
      ),
    );
  }

  @override
  void dispose() {
    _pageController.dispose();
    super.dispose();
  }
}