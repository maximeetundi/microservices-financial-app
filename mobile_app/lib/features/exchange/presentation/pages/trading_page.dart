import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';
import 'package:go_router/go_router.dart';
import 'package:google_fonts/google_fonts.dart';

import '../../../../core/theme/app_theme.dart';
import '../../../../core/widgets/custom_button.dart';
import '../../../../core/widgets/glass_container.dart';
import '../../../../core/widgets/loading_widget.dart';
import '../bloc/exchange_bloc.dart';

class TradingPage extends StatefulWidget {
  const TradingPage({super.key});

  @override
  State<TradingPage> createState() => _TradingPageState();
}

class _TradingPageState extends State<TradingPage> with SingleTickerProviderStateMixin {
  late TabController _tabController;
  final ScrollController _scrollController = ScrollController();
  
  // State
  String _selectedPair = 'BTC/USD';
  Map<String, dynamic> _marketData = {};
  Map<String, dynamic> _portfolio = {};
  List<Map<String, dynamic>> _recentOrders = [];
  bool _isLoading = true;

  @override
  void initState() {
    super.initState();
    _tabController = TabController(length: 2, vsync: this);
    _loadInitialData();
  }

  void _loadInitialData() {
    final bloc = context.read<ExchangeBloc>();
    bloc.add(const LoadMarketsEvent());
    bloc.add(const LoadTradingPortfolioEvent());
    bloc.add(const LoadOrdersEvent());
  }

  @override
  void dispose() {
    _tabController.dispose();
    _scrollController.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    final isDark = Theme.of(context).brightness == Brightness.dark;

    return Scaffold(
      backgroundColor: Colors.transparent,
      body: BlocListener<ExchangeBloc, ExchangeState>(
        listener: (context, state) {
          if (state is MarketsLoadedState) {
            setState(() {
              _marketData = state.markets;
              _isLoading = false;
            });
          } else if (state is TradingPortfolioLoadedState) {
            setState(() => _portfolio = state.portfolio);
          } else if (state is OrdersLoadedState) {
            setState(() => _recentOrders = state.orders);
          } else if (state is OrderPlacedState) {
            ScaffoldMessenger.of(context).showSnackBar(
              const SnackBar(content: Text('Ordre placé avec succès!'), backgroundColor: Colors.green),
            );
            context.read<ExchangeBloc>().add(const LoadOrdersEvent());
            context.read<ExchangeBloc>().add(const LoadTradingPortfolioEvent());
          } else if (state is ExchangeErrorState) {
             ScaffoldMessenger.of(context).showSnackBar(
              SnackBar(content: Text(state.message), backgroundColor: Colors.red),
            );
            setState(() => _isLoading = false);
          }
        },
        child: Container(
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
                _buildAppBar(isDark),
                if (_isLoading)
                  const Expanded(child: LoadingWidget())
                else
                  Expanded(
                    child: SingleChildScrollView(
                      controller: _scrollController,
                      padding: const EdgeInsets.all(16),
                      child: Column(
                        crossAxisAlignment: CrossAxisAlignment.start,
                        children: [
                          _buildMarketTickers(isDark),
                          const SizedBox(height: 24),
                          _buildTradingInterface(isDark),
                          const SizedBox(height: 24),
                          _buildRecentOrders(isDark),
                        ],
                      ),
                    ),
                  ),
              ],
            ),
          ),
        ),
      ),
    );
  }

  Widget _buildAppBar(bool isDark) {
    return Padding(
      padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 8),
      child: Row(
        mainAxisAlignment: MainAxisAlignment.spaceBetween,
        children: [
          Row(
            children: [
               GlassContainer(
                padding: EdgeInsets.zero,
                width: 40,
                height: 40,
                borderRadius: 12,
                child: IconButton(
                  icon: Icon(Icons.arrow_back_ios_new, size: 20, 
                      color: isDark ? Colors.white : AppTheme.textPrimaryColor),
                  onPressed: () => context.pop(),
                ),
              ),
              const SizedBox(width: 16),
              Text(
                'Trading Avancé',
                style: GoogleFonts.inter(
                  fontSize: 20,
                  fontWeight: FontWeight.bold,
                  color: isDark ? Colors.white : const Color(0xFF1E293B),
                ),
              ),
            ],
          ),
          if (_portfolio.isNotEmpty)
            Column(
              crossAxisAlignment: CrossAxisAlignment.end,
              children: [
                Text(
                  'Solde Total',
                  style: GoogleFonts.inter(
                    fontSize: 10,
                    color: isDark ? Colors.white70 : const Color(0xFF64748B),
                  ),
                ),
                Text(
                  '\$${(_portfolio['totalValue'] ?? 0).toString()}',
                  style: GoogleFonts.inter(
                    fontSize: 14,
                    fontWeight: FontWeight.bold,
                    color: const Color(0xFF22C55E),
                  ),
                ),
              ],
            ),
        ],
      ),
    );
  }

  Widget _buildMarketTickers(bool isDark) {
    // Mock tickers if _marketData is not formatted as list yet, or extract from it
    final markets = (_marketData['markets'] as List?) ?? [];
    
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Text(
          'Marchés',
          style: GoogleFonts.inter(
            fontWeight: FontWeight.w600,
            color: isDark ? Colors.white70 : const Color(0xFF64748B),
          ),
        ),
        const SizedBox(height: 12),
        SizedBox(
          height: 100,
          child: ListView.builder(
            scrollDirection: Axis.horizontal,
            itemCount: markets.length,
            itemBuilder: (context, index) {
              final market = markets[index];
              final change = (market['change_24h'] as num?)?.toDouble() ?? 0.0;
              final isPositive = change >= 0;
              final isSelected = market['symbol'] == _selectedPair;

              return GestureDetector(
                onTap: () => setState(() => _selectedPair = market['symbol']),
                child: Container(
                  width: 140,
                  margin: const EdgeInsets.only(right: 12),
                  padding: const EdgeInsets.all(12),
                  decoration: BoxDecoration(
                    color: isSelected 
                        ? (isDark ? Colors.white.withOpacity(0.1) : Colors.white)
                        : (isDark ? Colors.white.withOpacity(0.05) : Colors.white.withOpacity(0.6)),
                    borderRadius: BorderRadius.circular(16),
                    border: Border.all(
                      color: isSelected 
                          ? const Color(0xFF6366F1) 
                          : Colors.transparent,
                      width: 2,
                    ),
                    boxShadow: isSelected ? [
                      BoxShadow(
                        color: const Color(0xFF6366F1).withOpacity(0.2),
                        blurRadius: 8,
                        offset: const Offset(0, 4),
                      )
                    ] : null,
                  ),
                  child: Column(
                    crossAxisAlignment: CrossAxisAlignment.start,
                    mainAxisAlignment: MainAxisAlignment.spaceBetween,
                    children: [
                      Row(
                        mainAxisAlignment: MainAxisAlignment.spaceBetween,
                        children: [
                          Text(
                            market['symbol'] ?? '',
                            style: GoogleFonts.inter(
                              fontWeight: FontWeight.bold,
                              color: isDark ? Colors.white : const Color(0xFF1E293B),
                            ),
                          ),
                          if (isSelected)
                            const Icon(Icons.check_circle, size: 16, color: Color(0xFF6366F1))
                        ],
                      ),
                      Text(
                        '\$${market['price']}',
                        style: GoogleFonts.inter(
                          fontSize: 16,
                          fontWeight: FontWeight.w600,
                          color: isDark ? Colors.white : const Color(0xFF1E293B),
                        ),
                      ),
                      Container(
                        padding: const EdgeInsets.symmetric(horizontal: 6, vertical: 2),
                        decoration: BoxDecoration(
                          color: isPositive 
                              ? const Color(0xFF22C55E).withOpacity(0.15)
                              : const Color(0xFFEF4444).withOpacity(0.15),
                          borderRadius: BorderRadius.circular(4),
                        ),
                        child: Text(
                          '${isPositive ? '+' : ''}$change%',
                          style: GoogleFonts.inter(
                            fontSize: 10,
                            fontWeight: FontWeight.bold,
                            color: isPositive ? const Color(0xFF22C55E) : const Color(0xFFEF4444),
                          ),
                        ),
                      ),
                    ],
                  ),
                ),
              );
            },
          ),
        ),
      ],
    );
  }

  Widget _buildTradingInterface(bool isDark) {
    return Container(
      padding: const EdgeInsets.all(16),
      decoration: BoxDecoration(
        color: isDark ? Colors.white.withOpacity(0.05) : Colors.white,
        borderRadius: BorderRadius.circular(24),
        border: Border.all(
          color: isDark ? Colors.white.withOpacity(0.1) : const Color(0xFFE2E8F0),
        ),
      ),
      child: Column(
        children: [
          // Pair Header
          Row(
            mainAxisAlignment: MainAxisAlignment.spaceBetween,
            children: [
              Text(
                'Placer un ordre',
                style: GoogleFonts.inter(
                  fontSize: 16,
                  fontWeight: FontWeight.bold,
                  color: isDark ? Colors.white : const Color(0xFF1E293B),
                ),
              ),
              Text(
                _selectedPair,
                style: GoogleFonts.inter(
                  fontSize: 14,
                  fontWeight: FontWeight.w500,
                  color: const Color(0xFF6366F1),
                ),
              ),
            ],
          ),
          const SizedBox(height: 16),
          
          // Tabs
          Container(
            padding: const EdgeInsets.all(4),
            decoration: BoxDecoration(
              color: isDark ? Colors.black26 : Colors.grey.shade100,
              borderRadius: BorderRadius.circular(12),
            ),
            child: TabBar(
              controller: _tabController,
              indicator: BoxDecoration(
                color: isDark ? const Color(0xFF1E293B) : Colors.white,
                borderRadius: BorderRadius.circular(8),
                boxShadow: [
                  BoxShadow(
                    color: Colors.black.withOpacity(0.1),
                    blurRadius: 4,
                    offset: const Offset(0, 2),
                  ),
                ],
              ),
              labelColor: isDark ? Colors.white : AppTheme.textPrimaryColor,
              unselectedLabelColor: isDark ? Colors.white54 : Colors.grey,
              labelStyle: GoogleFonts.inter(fontWeight: FontWeight.w600),
              tabs: const [
                Tab(text: 'Acheter'),
                Tab(text: 'Vendre'),
              ],
            ),
          ),
          
          const SizedBox(height: 16),
          
          // Forms
          SizedBox(
            height: 350, // Fixed height for form
            child: TabBarView(
              controller: _tabController,
              children: [
                _TradingForm(pair: _selectedPair, isBuy: true),
                _TradingForm(pair: _selectedPair, isBuy: false),
              ],
            ),
          ),
        ],
      ),
    );
  }

  Widget _buildRecentOrders(bool isDark) {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Row(
          mainAxisAlignment: MainAxisAlignment.spaceBetween,
          children: [
            Text(
              'Ordres Récents',
              style: GoogleFonts.inter(
                fontSize: 16,
                fontWeight: FontWeight.bold,
                color: isDark ? Colors.white : const Color(0xFF1E293B),
              ),
            ),
            TextButton(
              onPressed: () {}, // Navigate to history
              child: const Text('Voir tout'),
            ),
          ],
        ),
        const SizedBox(height: 12),
        if (_recentOrders.isEmpty)
          Padding(
            padding: const EdgeInsets.symmetric(vertical: 20),
            child: Center(
              child: Text(
                'Aucun ordre récent',
                style: GoogleFonts.inter(color: Colors.grey),
              ),
            ),
          )
        else
          ListView.builder(
            shrinkWrap: true,
            physics: const NeverScrollableScrollPhysics(),
            itemCount: _recentOrders.length > 5 ? 5 : _recentOrders.length,
            itemBuilder: (context, index) {
              final order = _recentOrders[index];
              return _buildOrderItem(order, isDark);
            },
          ),
      ],
    );
  }

  Widget _buildOrderItem(Map<String, dynamic> order, bool isDark) {
    final type = order['side']?.toString().toUpperCase() ?? 'UNKNOWN';
    final isBuy = type == 'BUY';
    final status = order['status'] ?? 'pending';

    return Container(
      margin: const EdgeInsets.only(bottom: 8),
      padding: const EdgeInsets.all(12),
      decoration: BoxDecoration(
        color: isDark ? Colors.white.withOpacity(0.03) : Colors.white,
        borderRadius: BorderRadius.circular(12),
        border: Border.all(
          color: isDark ? Colors.white.withOpacity(0.05) : Colors.grey.shade200,
        ),
      ),
      child: Row(
        mainAxisAlignment: MainAxisAlignment.spaceBetween,
        children: [
          Row(
            children: [
              Container(
                padding: const EdgeInsets.all(8),
                decoration: BoxDecoration(
                  color: isBuy 
                      ? const Color(0xFF22C55E).withOpacity(0.15)
                      : const Color(0xFFEF4444).withOpacity(0.15),
                  shape: BoxShape.circle,
                ),
                child: Icon(
                  isBuy ? Icons.arrow_downward : Icons.arrow_upward,
                  size: 16,
                  color: isBuy ? const Color(0xFF22C55E) : const Color(0xFFEF4444),
                ),
              ),
              const SizedBox(width: 12),
              Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  Text(
                    order['pair'] ?? '',
                    style: GoogleFonts.inter(
                      fontWeight: FontWeight.w600,
                      color: isDark ? Colors.white : const Color(0xFF1E293B),
                    ),
                  ),
                  Text(
                    '${order['type']} • ${order['created_at']}',
                    style: GoogleFonts.inter(
                      fontSize: 10,
                      color: Colors.grey,
                    ),
                  ),
                ],
              ),
            ],
          ),
          Column(
            crossAxisAlignment: CrossAxisAlignment.end,
            children: [
              Text(
                '${order['amount']}',
                style: GoogleFonts.inter(
                  fontWeight: FontWeight.w600,
                  color: isDark ? Colors.white : const Color(0xFF1E293B),
                ),
              ),
              Text(
                status,
                style: GoogleFonts.inter(
                  fontSize: 10,
                  fontWeight: FontWeight.w500,
                  color: status == 'filled' ? const Color(0xFF22C55E) : Colors.orange,
                ),
              ),
            ],
          ),
        ],
      ),
    );
  }
}

class _TradingForm extends StatefulWidget {
  final String pair;
  final bool isBuy;

  const _TradingForm({required this.pair, required this.isBuy});

  @override
  State<_TradingForm> createState() => _TradingFormState();
}

class _TradingFormState extends State<_TradingForm> {
  final _amountController = TextEditingController();
  final _priceController = TextEditingController();
  final _stopPriceController = TextEditingController();
  
  String _orderType = 'market';
  double _estimatedValue = 0.0;
  
  // Get base/quote from pair (e.g., BTC/USD -> base=BTC, quote=USD)
  String get _baseCurrency => widget.pair.split('/').first;
  String get _quoteCurrency => widget.pair.length > 3 ? widget.pair.split('/').last : 'USD';

  void _calculateEstimate() {
    final amount = double.tryParse(_amountController.text) ?? 0.0;
    final price = double.tryParse(_priceController.text) ?? 43000.0; // Fallback to market price
    
    setState(() {
      _estimatedValue = amount * price;
    });
  }

  @override
  Widget build(BuildContext context) {
    final isDark = Theme.of(context).brightness == Brightness.dark;
    final primaryColor = widget.isBuy ? const Color(0xFF22C55E) : const Color(0xFFEF4444);

    return Column(
      children: [
        // Order Type Selector
        Container(
          height: 36,
          decoration: BoxDecoration(
            color: isDark ? Colors.white.withOpacity(0.05) : Colors.grey.shade100,
            borderRadius: BorderRadius.circular(8),
          ),
          child: Row(
            children: ['market', 'limit', 'stop_loss'].map((type) {
              final isSelected = _orderType == type;
              return Expanded(
                child: GestureDetector(
                  onTap: () => setState(() => _orderType = type),
                  child: Container(
                    decoration: BoxDecoration(
                      color: isSelected 
                          ? (isDark ? Colors.white.withOpacity(0.1) : Colors.white)
                          : Colors.transparent,
                      borderRadius: BorderRadius.circular(6),
                      boxShadow: isSelected ? [
                        BoxShadow(
                          color: Colors.black.withOpacity(0.05),
                          blurRadius: 2,
                        )
                      ] : null,
                    ),
                    alignment: Alignment.center,
                    child: Text(
                      type.toUpperCase().replaceAll('_', ' '),
                      style: GoogleFonts.inter(
                        fontSize: 10,
                        fontWeight: FontWeight.bold,
                        color: isSelected 
                            ? (isDark ? Colors.white : AppTheme.textPrimaryColor)
                            : (isDark ? Colors.white54 : Colors.grey),
                      ),
                    ),
                  ),
                ),
              );
            }).toList(),
          ),
        ),
        
        const SizedBox(height: 20),
        
        // Price Input (if not market)
        if (_orderType != 'market') ...[
          _buildInputField(
            label: 'Prix',
            controller: _priceController,
            suffix: _quoteCurrency,
            isDark: isDark,
            onChanged: (_) => _calculateEstimate(),
          ),
          const SizedBox(height: 16),
        ],

        // Stop Price Input (if stop loss)
        if (_orderType == 'stop_loss') ...[
          _buildInputField(
            label: 'Stop Prix',
            controller: _stopPriceController,
            suffix: _quoteCurrency,
            isDark: isDark,
          ),
          const SizedBox(height: 16),
        ],

        // Amount Input
        _buildInputField(
          label: 'Montant',
          controller: _amountController,
          suffix: _baseCurrency,
          isDark: isDark,
          onChanged: (_) => _calculateEstimate(),
        ),

        const SizedBox(height: 20),

        // Estimate
        Container(
          padding: const EdgeInsets.all(12),
          decoration: BoxDecoration(
            color: isDark ? Colors.white.withOpacity(0.05) : Colors.grey.shade50,
            borderRadius: BorderRadius.circular(12),
          ),
          child: Row(
            mainAxisAlignment: MainAxisAlignment.spaceBetween,
            children: [
              Text(
                'Total Estimé',
                style: GoogleFonts.inter(
                  color: isDark ? Colors.white70 : Colors.grey.shade700,
                ),
              ),
              Text(
                '$_estimatedValue $_quoteCurrency',
                style: GoogleFonts.inter(
                  fontWeight: FontWeight.bold,
                  color: isDark ? Colors.white : AppTheme.textPrimaryColor,
                ),
              ),
            ],
          ),
        ),

        const Spacer(),

        // Action Button
        CustomButton(
          text: '${widget.isBuy ? 'ACHETER' : 'VENDRE'} $_baseCurrency',
          onPressed: () {
            // Place Order
            final amount = double.tryParse(_amountController.text);
            if (amount == null || amount <= 0) return;

            final price = double.tryParse(_priceController.text);
            final stopPrice = double.tryParse(_stopPriceController.text);

            context.read<ExchangeBloc>().add(
              PlaceOrderEvent(
                symbol: widget.pair,
                side: widget.isBuy ? 'buy' : 'sell',
                type: _orderType,
                amount: amount,
                price: price,
                stopPrice: stopPrice,
              ),
            );
          },
          backgroundColor: primaryColor,
          textColor: Colors.white,
        ),
      ],
    );
  }

  Widget _buildInputField({
    required String label,
    required TextEditingController controller,
    required String suffix,
    required bool isDark,
    Function(String)? onChanged,
  }) {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Text(
          label,
          style: GoogleFonts.inter(
            fontSize: 12,
            color: isDark ? Colors.white60 : Colors.grey.shade700,
          ),
        ),
        const SizedBox(height: 6),
        Container(
          decoration: BoxDecoration(
            color: isDark ? Colors.white.withOpacity(0.05) : Colors.white,
            borderRadius: BorderRadius.circular(12),
            border: Border.all(
              color: isDark ? Colors.white.withOpacity(0.1) : Colors.grey.shade300,
            ),
          ),
          padding: const EdgeInsets.symmetric(horizontal: 12),
          child: TextField(
            controller: controller,
            keyboardType: const TextInputType.numberWithOptions(decimal: true),
            style: GoogleFonts.inter(
              color: isDark ? Colors.white : AppTheme.textPrimaryColor,
              fontWeight: FontWeight.w600,
            ),
            decoration: InputDecoration(
              border: InputBorder.none,
              isDense: true,
              contentPadding: const EdgeInsets.symmetric(vertical: 12),
              suffix: Text(
                suffix,
                style: GoogleFonts.inter(
                  color: Colors.grey,
                  fontWeight: FontWeight.bold,
                  fontSize: 12,
                ),
              ),
            ),
            onChanged: onChanged,
          ),
        ),
      ],
    );
  }
}
