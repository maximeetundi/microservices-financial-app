import 'package:flutter/material.dart';

class CryptoPricesSection extends StatelessWidget {
  const CryptoPricesSection({super.key});

  @override
  Widget build(BuildContext context) {
    final cryptos = [
      _CryptoPrice(symbol: 'BTC', name: 'Bitcoin', price: 43250.00, change: 2.45, icon: '₿'),
      _CryptoPrice(symbol: 'ETH', name: 'Ethereum', price: 2280.50, change: -1.2, icon: 'Ξ'),
      _CryptoPrice(symbol: 'SOL', name: 'Solana', price: 98.75, change: 5.8, icon: '◎'),
      _CryptoPrice(symbol: 'XRP', name: 'Ripple', price: 0.52, change: 0.5, icon: '✕'),
    ];

    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Row(
          mainAxisAlignment: MainAxisAlignment.spaceBetween,
          children: [
            const Text(
              'Prix Crypto',
              style: TextStyle(fontSize: 18, fontWeight: FontWeight.bold),
            ),
            TextButton(
              onPressed: () {},
              child: const Text('Voir plus'),
            ),
          ],
        ),
        const SizedBox(height: 8),
        ...cryptos.map((crypto) => _buildCryptoTile(crypto)),
      ],
    );
  }

  Widget _buildCryptoTile(_CryptoPrice crypto) {
    final isPositive = crypto.change >= 0;
    
    return ListTile(
      contentPadding: EdgeInsets.zero,
      leading: Container(
        width: 48,
        height: 48,
        decoration: BoxDecoration(
          color: Colors.grey[100],
          borderRadius: BorderRadius.circular(12),
        ),
        child: Center(
          child: Text(crypto.icon, style: const TextStyle(fontSize: 24)),
        ),
      ),
      title: Text(crypto.symbol, style: const TextStyle(fontWeight: FontWeight.bold)),
      subtitle: Text(crypto.name, style: const TextStyle(fontSize: 12)),
      trailing: Column(
        mainAxisAlignment: MainAxisAlignment.center,
        crossAxisAlignment: CrossAxisAlignment.end,
        children: [
          Text(
            '\$${crypto.price.toStringAsFixed(2)}',
            style: const TextStyle(fontWeight: FontWeight.bold),
          ),
          Row(
            mainAxisSize: MainAxisSize.min,
            children: [
              Icon(
                isPositive ? Icons.arrow_drop_up : Icons.arrow_drop_down,
                color: isPositive ? Colors.green : Colors.red,
                size: 20,
              ),
              Text(
                '${isPositive ? '+' : ''}${crypto.change.toStringAsFixed(2)}%',
                style: TextStyle(
                  fontSize: 12,
                  color: isPositive ? Colors.green : Colors.red,
                ),
              ),
            ],
          ),
        ],
      ),
    );
  }
}

class _CryptoPrice {
  final String symbol;
  final String name;
  final double price;
  final double change;
  final String icon;

  _CryptoPrice({
    required this.symbol,
    required this.name,
    required this.price,
    required this.change,
    required this.icon,
  });
}
