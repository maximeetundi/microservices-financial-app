import 'package:flutter/material.dart';

class CurrencySelector extends StatelessWidget {
  final String selectedCurrency;
  final List<String> currencies;
  final ValueChanged<String> onChanged;
  final bool isCrypto;

  const CurrencySelector({
    super.key,
    required this.selectedCurrency,
    required this.currencies,
    required this.onChanged,
    this.isCrypto = false,
  });

  @override
  Widget build(BuildContext context) {
    return Container(
      constraints: BoxConstraints(maxHeight: MediaQuery.of(context).size.height * 0.6),
      child: Column(
        mainAxisSize: MainAxisSize.min,
        children: [
          Container(
            padding: const EdgeInsets.all(16),
            child: Row(
              children: [
                const Text(
                  'Sélectionner une devise',
                  style: TextStyle(fontSize: 18, fontWeight: FontWeight.bold),
                ),
                const Spacer(),
                IconButton(
                  icon: const Icon(Icons.close),
                  onPressed: () => Navigator.pop(context),
                ),
              ],
            ),
          ),
          Expanded(
            child: ListView.builder(
              itemCount: currencies.length,
              itemBuilder: (context, index) {
                final currency = currencies[index];
                final isSelected = currency == selectedCurrency;
                
                return ListTile(
                  leading: Container(
                    width: 40,
                    height: 40,
                    decoration: BoxDecoration(
                      color: isSelected 
                          ? Theme.of(context).primaryColor.withOpacity(0.1)
                          : Colors.grey[100],
                      borderRadius: BorderRadius.circular(10),
                    ),
                    child: Center(
                      child: Text(
                        _getCurrencyIcon(currency),
                        style: const TextStyle(fontSize: 20),
                      ),
                    ),
                  ),
                  title: Text(currency, style: const TextStyle(fontWeight: FontWeight.bold)),
                  subtitle: Text(_getCurrencyName(currency)),
                  trailing: isSelected
                      ? Icon(Icons.check_circle, color: Theme.of(context).primaryColor)
                      : null,
                  onTap: () {
                    onChanged(currency);
                    Navigator.pop(context);
                  },
                );
              },
            ),
          ),
        ],
      ),
    );
  }

  String _getCurrencyIcon(String currency) {
    final icons = {
      'USD': '\$', 'EUR': '€', 'GBP': '£', 'XOF': 'F',
      'BTC': '₿', 'ETH': 'Ξ', 'USDT': '₮', 'SOL': '◎',
    };
    return icons[currency] ?? currency[0];
  }

  String _getCurrencyName(String currency) {
    final names = {
      'USD': 'Dollar américain', 'EUR': 'Euro', 'GBP': 'Livre sterling', 'XOF': 'Franc CFA',
      'BTC': 'Bitcoin', 'ETH': 'Ethereum', 'USDT': 'Tether', 'SOL': 'Solana',
    };
    return names[currency] ?? currency;
  }
}
