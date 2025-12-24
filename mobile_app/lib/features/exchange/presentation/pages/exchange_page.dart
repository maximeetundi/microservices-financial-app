import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';
import 'package:go_router/go_router.dart';

import '../../../../core/theme/app_theme.dart';
import '../../../../core/widgets/custom_button.dart';
import '../../../../core/widgets/custom_text_field.dart';
import '../../../../core/widgets/loading_widget.dart';
import '../../../../core/widgets/glass_container.dart';
import '../bloc/exchange_bloc.dart';
import '../widgets/currency_selector.dart';
import '../widgets/exchange_rate_card.dart';
import '../widgets/exchange_history_list.dart';
import '../widgets/quick_exchange_amounts.dart';

class ExchangePage extends StatefulWidget {
  const ExchangePage({Key? key}) : super(key: key);

  @override
  State<ExchangePage> createState() => _ExchangePageState();
}

class _ExchangePageState extends State<ExchangePage>
    with SingleTickerProviderStateMixin {
  late TabController _tabController;
  
  final _fromAmountController = TextEditingController();
  final _toAmountController = TextEditingController();
  
  String _fromCurrency = 'BTC';
  String _toCurrency = 'USD';
  bool _isSwapping = false;

  @override
  void initState() {
    super.initState();
    _tabController = TabController(length: 3, vsync: this);
    _loadExchangeData();
  }

  void _loadExchangeData() {
    context.read<ExchangeBloc>().add(
      GetExchangeRateEvent(
        fromCurrency: _fromCurrency,
        toCurrency: _toCurrency,
      ),
    );
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
               // Custom App Bar
               Padding(
                 padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 8),
                 child: Row(
                   mainAxisAlignment: MainAxisAlignment.spaceBetween,
                   children: [
                     const SizedBox(width: 48), // Spacer
                     Text(
                       'Échange',
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
                          icon: Icon(Icons.history, color: isDark ? Colors.white : AppTheme.textPrimaryColor),
                          onPressed: () {
                             // Navigate to exchange history
                          },
                        ),
                     ),
                   ],
                 ),
               ),
               // Tab Bar
               Container(
                 margin: const EdgeInsets.symmetric(horizontal: 16),
                 decoration: BoxDecoration(
                   color: isDark ? Colors.white.withOpacity(0.05) : Colors.white.withOpacity(0.5),
                   borderRadius: BorderRadius.circular(16),
                 ),
                 child: TabBar(
                   controller: _tabController,
                   indicatorColor: AppTheme.primaryColor,
                   labelColor: AppTheme.primaryColor,
                   unselectedLabelColor: isDark ? Colors.white54 : Colors.grey,
                   labelStyle: GoogleFonts.inter(fontWeight: FontWeight.w600),
                   tabs: const [
                     Tab(text: 'Convertir'),
                     Tab(text: 'Trading'),
                     Tab(text: 'P2P'),
                   ],
                 ),
               ),
              Expanded(
                child: TabBarView(
                  controller: _tabController,
                  children: [
                    _buildExchangeTab(),
                    _buildTradingTab(),
                    _buildP2PTab(),
                  ],
                ),
              ),
            ],
          ),
        ),
      ),
    );
  }

  Widget _buildExchangeTab() {
    final isDark = Theme.of(context).brightness == Brightness.dark;
    
    return BlocListener<ExchangeBloc, ExchangeState>(
      listener: (context, state) {
        if (state is ExchangeSuccessState) {
          _showSuccessDialog(state.exchangeId);
        } else if (state is ExchangeErrorState) {
          _showErrorSnackBar(state.message);
        }
      },
      child: SingleChildScrollView(
        padding: const EdgeInsets.all(16),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.stretch,
          children: [
            // Exchange Rate Card
            BlocBuilder<ExchangeBloc, ExchangeState>(
              builder: (context, state) {
                if (state is ExchangeRateLoadedState) {
                  return ExchangeRateCard(
                    fromCurrency: _fromCurrency,
                    toCurrency: _toCurrency,
                    rate: state.rate,
                  );
                }
                return const SizedBox.shrink();
              },
            ),
            
            const SizedBox(height: 24),
            
            // From Currency Selection
            _buildCurrencySection(
              title: 'De',
              currency: _fromCurrency,
              controller: _fromAmountController,
              onCurrencyChanged: (currency) {
                setState(() {
                  _fromCurrency = currency;
                });
                _loadExchangeData();
              },
              onAmountChanged: _calculateToAmount,
            ),
            
            const SizedBox(height: 16),
            
            // Swap Button
            Center(
              child: GestureDetector(
                onTap: _swapCurrencies,
                child: AnimatedRotation(
                  turns: _isSwapping ? 0.5 : 0,
                  duration: const Duration(milliseconds: 300),
                  child: Container(
                    width: 48,
                    height: 48,
                    decoration: BoxDecoration(
                      gradient: AppTheme.primaryGradient,
                      shape: BoxShape.circle,
                      boxShadow: [
                        BoxShadow(
                          color: AppTheme.primaryColor.withOpacity(0.3),
                          blurRadius: 10,
                          offset: const Offset(0, 3),
                        ),
                      ],
                    ),
                    child: const Icon(
                      Icons.swap_vert,
                      color: Colors.white,
                      size: 24,
                    ),
                  ),
                ),
              ),
            ),
            
            const SizedBox(height: 16),
            
            // To Currency Selection
            _buildCurrencySection(
              title: 'Vers',
              currency: _toCurrency,
              controller: _toAmountController,
              onCurrencyChanged: (currency) {
                setState(() {
                  _toCurrency = currency;
                });
                _loadExchangeData();
              },
              isReadOnly: true,
            ),
            
            const SizedBox(height: 24),
            
            // Quick Amount Buttons
            QuickExchangeAmounts(
              amounts: [50.0, 100.0, 250.0, 500.0, 1000.0],
              currency: _fromCurrency,
              onSelected: (amount) {
                _fromAmountController.text = amount.toString();
                _calculateToAmount(amount.toString());
              },
            ),
            
            const SizedBox(height: 32),
            
            // Exchange Button
            BlocBuilder<ExchangeBloc, ExchangeState>(
              builder: (context, state) {
                final isLoading = state is ExchangeLoadingState;
                
                return Container(
                  decoration: BoxDecoration(
                    borderRadius: BorderRadius.circular(16),
                    boxShadow: [
                      BoxShadow(
                        color: AppTheme.primaryColor.withOpacity(0.3),
                        blurRadius: 20,
                        offset: const Offset(0, 8),
                      ),
                    ],
                  ),
                  child: CustomButton(
                    text: isLoading ? 'Traitement...' : 'Échanger',
                    onPressed: isLoading ? null : _handleExchange,
                    isLoading: isLoading,
                    backgroundColor: AppTheme.primaryColor,
                    textColor: Colors.white,
                  ),
                );
              },
            ),
            
            const SizedBox(height: 32),
            
            // Recent Exchanges
            _buildRecentExchanges(),
          ],
        ),
      ),
    );
  }

  Widget _buildTradingTab() {
    return Center(
      child: Column(
        mainAxisAlignment: MainAxisAlignment.center,
        children: [
          const Icon(
            Icons.trending_up,
            size: 80,
            color: AppTheme.primaryColor,
          ),
          const SizedBox(height: 16),
          Text(
            'Advanced Trading',
            style: Theme.of(context).textTheme.headlineMedium,
          ),
          const SizedBox(height: 8),
          const Text(
            'Professional trading tools with\ncharts and order management',
            textAlign: TextAlign.center,
          ),
          const SizedBox(height: 24),
          CustomButton(
            text: 'Open Trading',
            onPressed: () => context.push('/exchange/trading'),
          ),
        ],
      ),
    );
  }

  Widget _buildP2PTab() {
    return Center(
      child: Column(
        mainAxisAlignment: MainAxisAlignment.center,
        children: [
          const Icon(
            Icons.people,
            size: 80,
            color: AppTheme.secondaryColor,
          ),
          const SizedBox(height: 16),
          Text(
            'P2P Trading',
            style: Theme.of(context).textTheme.headlineMedium,
          ),
          const SizedBox(height: 8),
          const Text(
            'Buy and sell crypto directly\nwith other users',
            textAlign: TextAlign.center,
          ),
          const SizedBox(height: 24),
          CustomButton(
            text: 'Browse Offers',
            onPressed: () {
              // Navigate to P2P marketplace
            },
          ),
        ],
      ),
    );
  }

  Widget _buildCurrencySection({
    required String title,
    required String currency,
    required TextEditingController controller,
    required Function(String) onCurrencyChanged,
    Function(String)? onAmountChanged,
    bool isReadOnly = false,
  }) {
    final isDark = Theme.of(context).brightness == Brightness.dark;
    
    return Container(
      padding: const EdgeInsets.all(4),
      decoration: BoxDecoration(
        color: Colors.transparent, // Use glass container instead
        borderRadius: BorderRadius.circular(16),
      ),
      child: GlassContainer(
        padding: const EdgeInsets.all(16),
        borderRadius: 16,
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Text(
              title,
              style: GoogleFonts.inter(
                fontSize: 14,
                color: isDark ? Colors.white70 : AppTheme.textSecondaryColor,
              ),
            ),
            const SizedBox(height: 12),
            Row(
              children: [
                // Currency Selector
                GestureDetector(
                  onTap: () => _showCurrencyPicker(onCurrencyChanged),
                  child: Container(
                    padding: const EdgeInsets.symmetric(horizontal: 12, vertical: 8),
                    decoration: BoxDecoration(
                      color: isDark ? Colors.black26 : Colors.grey.shade100,
                      borderRadius: BorderRadius.circular(8),
                    ),
                    child: Row(
                      children: [
                        _getCurrencyIcon(currency),
                        const SizedBox(width: 8),
                        Text(
                          currency,
                          style: GoogleFonts.inter(
                            fontWeight: FontWeight.w600,
                            fontSize: 16,
                            color: isDark ? Colors.white : AppTheme.textPrimaryColor,
                          ),
                        ),
                        const SizedBox(width: 4),
                        Icon(
                          Icons.keyboard_arrow_down,
                          size: 20,
                          color: isDark ? Colors.white70 : AppTheme.textSecondaryColor,
                        ),
                      ],
                    ),
                  ),
                ),
                
                const SizedBox(width: 16),
                
                // Amount Input
                Expanded(
                  child: TextFormField(
                    controller: controller,
                    readOnly: isReadOnly,
                    keyboardType: const TextInputType.numberWithOptions(decimal: true),
                     style: GoogleFonts.inter(
                      fontSize: 20,
                      fontWeight: FontWeight.bold,
                      color: isDark ? Colors.white : AppTheme.textPrimaryColor,
                    ),
                    decoration: InputDecoration(
                      hintText: '0.00',
                      hintStyle: TextStyle(color: isDark ? Colors.white30 : Colors.grey.shade400),
                      border: InputBorder.none,
                      contentPadding: EdgeInsets.zero,
                    ),
                    onChanged: onAmountChanged,
                  ),
                ),
              ],
            ),
          ],
        ),
      ),
    );
  }

  Widget _buildRecentExchanges() {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Row(
          mainAxisAlignment: MainAxisAlignment.spaceBetween,
          children: [
            Text(
              'Recent Exchanges',
              style: Theme.of(context).textTheme.headlineSmall,
            ),
            TextButton(
              onPressed: () {
                // Navigate to full history
              },
              child: const Text('View All'),
            ),
          ],
        ),
        const SizedBox(height: 12),
        BlocBuilder<ExchangeBloc, ExchangeState>(
          builder: (context, state) {
            if (state is ExchangeHistoryLoadedState) {
              final historyItems = state.exchanges.map((e) => ExchangeHistoryItem(
                fromCurrency: e['from_currency']?.toString() ?? '',
                toCurrency: e['to_currency']?.toString() ?? '',
                fromAmount: (e['from_amount'] as num?)?.toDouble() ?? 0.0,
                toAmount: (e['to_amount'] as num?)?.toDouble() ?? 0.0,
                date: e['created_at']?.toString() ?? '',
                status: e['status']?.toString() ?? 'completed',
              )).toList();
              return ExchangeHistoryList(
                history: historyItems,
              );
            }
            return const SizedBox.shrink();
          },
        ),
      ],
    );
  }

  Widget _getCurrencyIcon(String currency) {
    // Return appropriate icon for currency
    final iconData = {
      'BTC': Icons.currency_bitcoin,
      'ETH': Icons.diamond,
      'USD': Icons.attach_money,
      'EUR': Icons.euro,
    }[currency] ?? Icons.monetization_on;

    return Icon(
      iconData,
      size: 20,
      color: AppTheme.primaryColor,
    );
  }

  void _showCurrencyPicker(Function(String) onCurrencyChanged) {
    showModalBottomSheet(
      context: context,
      builder: (context) => CurrencySelector(
        selectedCurrency: _fromCurrency,
        currencies: ['BTC', 'ETH', 'USD', 'EUR'],
        onChanged: (currency) {
          onCurrencyChanged(currency);
        },
      ),
    );
  }

  void _swapCurrencies() {
    setState(() {
      _isSwapping = true;
      final temp = _fromCurrency;
      _fromCurrency = _toCurrency;
      _toCurrency = temp;
      
      // Swap amounts too
      final tempAmount = _fromAmountController.text;
      _fromAmountController.text = _toAmountController.text;
      _toAmountController.text = tempAmount;
    });
    
    Future.delayed(const Duration(milliseconds: 300), () {
      setState(() {
        _isSwapping = false;
      });
    });
    
    _loadExchangeData();
  }

  void _calculateToAmount(String fromAmount) {
    if (fromAmount.isEmpty) {
      _toAmountController.clear();
      return;
    }
    
    final amount = double.tryParse(fromAmount);
    if (amount == null) return;
    
    // Recalculate on change
    context.read<ExchangeBloc>().add(
      GetExchangeRateEvent(
        fromCurrency: _fromCurrency,
        toCurrency: _toCurrency,
      ),
    );
  }

  void _handleExchange() {
    final amount = double.tryParse(_fromAmountController.text);
    if (amount == null || amount <= 0) {
      _showErrorSnackBar('Please enter a valid amount');
      return;
    }
    
    // Show confirmation dialog
    _showExchangeConfirmation(amount);
  }

  void _showExchangeConfirmation(double amount) {
    showDialog(
      context: context,
      builder: (context) => AlertDialog(
        title: const Text('Confirm Exchange'),
        content: Column(
          mainAxisSize: MainAxisSize.min,
          children: [
            Text('Exchange $amount $_fromCurrency to $_toCurrency?'),
            const SizedBox(height: 16),
            Container(
              padding: const EdgeInsets.all(12),
              decoration: BoxDecoration(
                color: Colors.grey.shade100,
                borderRadius: BorderRadius.circular(8),
              ),
              child: Column(
                children: [
                  Row(
                    mainAxisAlignment: MainAxisAlignment.spaceBetween,
                    children: [
                      const Text('Exchange Rate:'),
                      Text('1 $_fromCurrency = 43,500 $_toCurrency'),
                    ],
                  ),
                  const SizedBox(height: 8),
                  Row(
                    mainAxisAlignment: MainAxisAlignment.spaceBetween,
                    children: [
                      const Text('Fee:'),
                      Text('0.25% (${(amount * 0.0025).toStringAsFixed(8)} $_fromCurrency)'),
                    ],
                  ),
                ],
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
              context.read<ExchangeBloc>().add(
                ExecuteExchangeEvent(
                  fromWalletId: 'wallet-1', // TODO: Get from state
                  toWalletId: 'wallet-2', // TODO: Get from state
                  fromCurrency: _fromCurrency,
                  toCurrency: _toCurrency,
                  amount: amount,
                ),
              );
            },
            child: const Text('Confirm'),
          ),
        ],
      ),
    );
  }

  void _showSuccessDialog(dynamic transaction) {
    showDialog(
      context: context,
      builder: (context) => AlertDialog(
        icon: const Icon(
          Icons.check_circle,
          color: AppTheme.successColor,
          size: 48,
        ),
        title: const Text('Exchange Successful'),
        content: const Text('Your exchange has been completed successfully!'),
        actions: [
          ElevatedButton(
            onPressed: () => Navigator.pop(context),
            child: const Text('OK'),
          ),
        ],
      ),
    );
  }

  void _showErrorSnackBar(String message) {
    ScaffoldMessenger.of(context).showSnackBar(
      SnackBar(
        content: Text(message),
        backgroundColor: AppTheme.errorColor,
      ),
    );
  }

  @override
  void dispose() {
    _tabController.dispose();
    _fromAmountController.dispose();
    _toAmountController.dispose();
    super.dispose();
  }
}