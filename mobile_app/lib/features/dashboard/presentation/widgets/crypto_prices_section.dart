import 'package:flutter_bloc/flutter_bloc.dart';
import '../../exchange/presentation/bloc/exchange_bloc.dart';

class CryptoPricesSection extends StatelessWidget {
  const CryptoPricesSection({super.key});

  @override
  Widget build(BuildContext context) {
    return BlocBuilder<ExchangeBloc, ExchangeState>(
      builder: (context, state) {
        if (state is ExchangeLoadingState) {
          return const Center(child: CircularProgressIndicator());
        }

        List<_CryptoPrice> cryptos = [];
        
        if (state is ExchangeRatesLoadedState) {
           // Parse rates
           cryptos = state.rates.take(5).map((r) {
             final symbol = r['symbol'] ?? r['pair'] ?? 'UNK';
             // Extract base symbol e.g BTC from BTC/USD
             final baseSymbol = symbol.contains('/') ? symbol.split('/')[0] : symbol;
             
             return _CryptoPrice(
               symbol: baseSymbol, 
               name: _getName(baseSymbol), 
               price: (r['price'] ?? r['rate'] ?? 0).toDouble(), 
               change: (r['change_24h'] ?? 0).toDouble(), 
               icon: _getIcon(baseSymbol)
             );
           }).toList();
        } else {
           // Fallback or empty if not loaded yet (or error)
           // We could show previous hardcoded data as skeletons or nothing.
           // For now, let's return SizedBox if empty to avoid clutter, 
           // or keep the structure with empty list.
           if (cryptos.isEmpty) return const SizedBox.shrink();
        }

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
            if (cryptos.isEmpty) 
               const Text("Aucune donnée disponible")
            else
               ...cryptos.map((crypto) => _buildCryptoTile(crypto)),
          ],
        );
      },
    );
  }

  String _getName(String symbol) {
    switch (symbol) {
      case 'BTC': return 'Bitcoin';
      case 'ETH': return 'Ethereum';
      case 'SOL': return 'Solana';
      case 'XRP': return 'Ripple';
      case 'USDT': return 'Tether';
      case 'ADA': return 'Cardano';
      case 'DOGE': return 'Dogecoin';
      default: return symbol;
    }
  }

  String _getIcon(String symbol) {
     switch (symbol) {
      case 'BTC': return '₿';
      case 'ETH': return 'Ξ';
      case 'SOL': return '◎';
      case 'XRP': return '✕';
      case 'USDT': return '₮';
      case 'ADA': return '₳';
      case 'DOGE': return 'Ð';
      default: return '\$';
    }
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
