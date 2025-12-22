import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';

class TradingPage extends StatefulWidget {
  const TradingPage({super.key});

  @override
  State<TradingPage> createState() => _TradingPageState();
}

class _TradingPageState extends State<TradingPage> with SingleTickerProviderStateMixin {
  late TabController _tabController;

  @override
  void initState() {
    super.initState();
    _tabController = TabController(length: 2, vsync: this);
  }

  @override
  void dispose() {
    _tabController.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: const Text('Trading'),
        leading: IconButton(
          icon: const Icon(Icons.arrow_back),
          onPressed: () => context.pop(),
        ),
        bottom: TabBar(
          controller: _tabController,
          tabs: const [
            Tab(text: 'Acheter'),
            Tab(text: 'Vendre'),
          ],
        ),
      ),
      body: TabBarView(
        controller: _tabController,
        children: [
          _TradingForm(isBuy: true),
          _TradingForm(isBuy: false),
        ],
      ),
    );
  }
}

class _TradingForm extends StatefulWidget {
  final bool isBuy;

  const _TradingForm({required this.isBuy});

  @override
  State<_TradingForm> createState() => _TradingFormState();
}

class _TradingFormState extends State<_TradingForm> {
  final _amountController = TextEditingController();
  String _selectedCrypto = 'BTC';
  String _selectedFiat = 'USD';
  double _estimatedAmount = 0;

  final List<String> _cryptos = ['BTC', 'ETH', 'USDT', 'SOL', 'XRP'];
  final List<String> _fiats = ['USD', 'EUR', 'GBP', 'XOF'];

  void _calculateEstimate() {
    final amount = double.tryParse(_amountController.text) ?? 0;
    // Mock prices
    final prices = {'BTC': 43000.0, 'ETH': 2200.0, 'USDT': 1.0, 'SOL': 100.0, 'XRP': 0.5};
    final price = prices[_selectedCrypto] ?? 1;
    
    setState(() {
      if (widget.isBuy) {
        _estimatedAmount = amount / price;
      } else {
        _estimatedAmount = amount * price;
      }
    });
  }

  @override
  Widget build(BuildContext context) {
    return SingleChildScrollView(
      padding: const EdgeInsets.all(16),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Text(
            widget.isBuy ? 'Acheter des crypto-monnaies' : 'Vendre des crypto-monnaies',
            style: const TextStyle(fontSize: 20, fontWeight: FontWeight.bold),
          ),
          const SizedBox(height: 24),
          
          // Crypto selector
          DropdownButtonFormField<String>(
            value: _selectedCrypto,
            decoration: const InputDecoration(
              labelText: 'Crypto-monnaie',
              border: OutlineInputBorder(),
            ),
            items: _cryptos.map((c) => DropdownMenuItem(value: c, child: Text(c))).toList(),
            onChanged: (v) {
              setState(() => _selectedCrypto = v!);
              _calculateEstimate();
            },
          ),
          const SizedBox(height: 16),
          
          // Amount input
          TextFormField(
            controller: _amountController,
            keyboardType: TextInputType.number,
            decoration: InputDecoration(
              labelText: widget.isBuy ? 'Montant à dépenser' : 'Quantité à vendre',
              suffixText: widget.isBuy ? _selectedFiat : _selectedCrypto,
              border: const OutlineInputBorder(),
            ),
            onChanged: (_) => _calculateEstimate(),
          ),
          const SizedBox(height: 16),
          
          // Fiat selector
          DropdownButtonFormField<String>(
            value: _selectedFiat,
            decoration: const InputDecoration(
              labelText: 'Devise',
              border: OutlineInputBorder(),
            ),
            items: _fiats.map((f) => DropdownMenuItem(value: f, child: Text(f))).toList(),
            onChanged: (v) {
              setState(() => _selectedFiat = v!);
              _calculateEstimate();
            },
          ),
          const SizedBox(height: 24),
          
          // Estimate
          Container(
            width: double.infinity,
            padding: const EdgeInsets.all(16),
            decoration: BoxDecoration(
              color: Colors.grey[100],
              borderRadius: BorderRadius.circular(12),
            ),
            child: Column(
              children: [
                Text(
                  widget.isBuy ? 'Vous recevrez environ' : 'Vous recevrez environ',
                  style: TextStyle(color: Colors.grey[600]),
                ),
                const SizedBox(height: 8),
                Text(
                  widget.isBuy 
                    ? '${_estimatedAmount.toStringAsFixed(6)} $_selectedCrypto'
                    : '${_estimatedAmount.toStringAsFixed(2)} $_selectedFiat',
                  style: const TextStyle(fontSize: 24, fontWeight: FontWeight.bold),
                ),
              ],
            ),
          ),
          const SizedBox(height: 32),
          
          SizedBox(
            width: double.infinity,
            child: ElevatedButton(
              onPressed: () {
                ScaffoldMessenger.of(context).showSnackBar(
                  const SnackBar(content: Text('Trading effectué!')),
                );
              },
              style: ElevatedButton.styleFrom(
                padding: const EdgeInsets.symmetric(vertical: 16),
                backgroundColor: widget.isBuy ? Colors.green : Colors.red,
              ),
              child: Text(
                widget.isBuy ? 'Acheter' : 'Vendre',
                style: const TextStyle(fontSize: 16, color: Colors.white),
              ),
            ),
          ),
        ],
      ),
    );
  }
}
