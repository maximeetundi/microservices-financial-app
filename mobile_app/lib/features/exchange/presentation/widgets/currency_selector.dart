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
      // Fiat
      'USD': '\$', 'EUR': '€', 'GBP': '£', 'JPY': '¥', 'CHF': '₣',
      'CAD': '\$', 'AUD': '\$', 'NZD': '\$', 'MXN': '\$', 'BRL': 'R\$',
      'CNY': '¥', 'HKD': '\$', 'SGD': '\$', 'KRW': '₩', 'INR': '₹',
      'IDR': 'Rp', 'MYR': 'RM', 'THB': '฿', 'PHP': '₱', 'VND': '₫',
      'AED': 'د', 'SAR': 'ر', 'QAR': 'ر', 'KWD': 'د', 'EGP': '£',
      'XAF': 'F', 'XOF': 'F', 'NGN': '₦', 'ZAR': 'R', 'KES': 'Sh',
      'GHS': '₵', 'MAD': 'د', 'TND': 'د', 'DZD': 'د', 'UGX': 'Sh',
      'TZS': 'Sh', 'RWF': 'Fr', 'ETB': 'Br',
      'NOK': 'kr', 'SEK': 'kr', 'DKK': 'kr', 'PLN': 'zł', 'CZK': 'Kč',
      'HUF': 'Ft', 'RON': 'lei', 'TRY': '₺', 'RUB': '₽',
      // Crypto
      'BTC': '₿', 'ETH': 'Ξ', 'USDT': '₮', 'USDC': '\$', 'SOL': '◎',
      'XRP': '✕', 'BNB': '◆', 'ADA': '₳', 'DOGE': 'Ð', 'DOT': '●',
      'LTC': 'Ł', 'AVAX': '▲', 'MATIC': '◇', 'LINK': '⬡', 'UNI': '🦄',
      'ATOM': '⚛', 'ALGO': '⌬', 'VET': '⌘', 'XLM': '★', 'FIL': '⌨',
    };
    return icons[currency] ?? currency[0];
  }

  String _getCurrencyName(String currency) {
    final names = {
      // Major Fiat
      'USD': 'Dollar américain', 'EUR': 'Euro', 'GBP': 'Livre sterling',
      'JPY': 'Yen japonais', 'CHF': 'Franc suisse',
      // Americas
      'CAD': 'Dollar canadien', 'MXN': 'Peso mexicain', 'BRL': 'Réal brésilien',
      'ARS': 'Peso argentin', 'CLP': 'Peso chilien', 'COP': 'Peso colombien',
      'PEN': 'Sol péruvien', 'AUD': 'Dollar australien', 'NZD': 'Dollar néo-zélandais',
      // Europe
      'NOK': 'Couronne norvégienne', 'SEK': 'Couronne suédoise',
      'DKK': 'Couronne danoise', 'PLN': 'Zloty polonais', 'CZK': 'Couronne tchèque',
      'HUF': 'Forint hongrois', 'RON': 'Leu roumain', 'TRY': 'Livre turque',
      'RUB': 'Rouble russe', 'UAH': 'Hryvnia ukrainien',
      // Asia
      'CNY': 'Yuan chinois', 'HKD': 'Dollar de Hong Kong', 'SGD': 'Dollar de Singapour',
      'KRW': 'Won sud-coréen', 'INR': 'Roupie indienne', 'IDR': 'Roupie indonésienne',
      'MYR': 'Ringgit malaisien', 'THB': 'Baht thaïlandais', 'PHP': 'Peso philippin',
      'VND': 'Dong vietnamien', 'PKR': 'Roupie pakistanaise', 'BDT': 'Taka bangladais',
      // Middle East
      'AED': 'Dirham des EAU', 'SAR': 'Riyal saoudien', 'QAR': 'Riyal qatari',
      'KWD': 'Dinar koweïtien', 'BHD': 'Dinar bahreïni', 'OMR': 'Rial omanais',
      'ILS': 'Shekel israélien', 'EGP': 'Livre égyptienne', 'JOD': 'Dinar jordanien',
      // Africa
      'XAF': 'Franc CFA (CEMAC)', 'XOF': 'Franc CFA (UEMOA)', 'NGN': 'Naira nigérian',
      'ZAR': 'Rand sud-africain', 'KES': 'Shilling kényan', 'GHS': 'Cédi ghanéen',
      'MAD': 'Dirham marocain', 'TND': 'Dinar tunisien', 'DZD': 'Dinar algérien',
      'UGX': 'Shilling ougandais', 'TZS': 'Shilling tanzanien', 'RWF': 'Franc rwandais',
      'ETB': 'Birr éthiopien', 'MUR': 'Roupie mauricienne',
      // Crypto
      'BTC': 'Bitcoin', 'ETH': 'Ethereum', 'USDT': 'Tether USD', 'USDC': 'USD Coin',
      'SOL': 'Solana', 'XRP': 'Ripple', 'BNB': 'BNB Chain', 'ADA': 'Cardano',
      'DOGE': 'Dogecoin', 'DOT': 'Polkadot', 'LTC': 'Litecoin', 'AVAX': 'Avalanche',
      'MATIC': 'Polygon', 'LINK': 'Chainlink', 'UNI': 'Uniswap', 'ATOM': 'Cosmos',
      'ALGO': 'Algorand', 'VET': 'VeChain', 'XLM': 'Stellar', 'FIL': 'Filecoin',
    };
    return names[currency] ?? currency;
  }
}
