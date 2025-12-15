import 'package:intl/intl.dart';

/// Utility class for formatting currency values
class CurrencyFormatter {
  // Currency symbols
  static const Map<String, String> _symbols = {
    'EUR': '€',
    'USD': '\$',
    'GBP': '£',
    'XOF': 'CFA',
    'XAF': 'CFA',
    'NGN': '₦',
    'GHS': '₵',
    'KES': 'KSh',
    'ZAR': 'R',
    'BTC': '₿',
    'ETH': 'Ξ',
    'USDT': '\$',
    'USDC': '\$',
  };

  /// Format amount with currency symbol
  static String format(dynamic amount, String currency) {
    if (amount == null) return '${_symbols[currency] ?? currency} 0.00';
    
    double value = amount is num ? amount.toDouble() : 0.0;
    
    // Crypto currencies have more decimal places
    int decimalPlaces = _isCrypto(currency) ? 8 : 2;
    
    final formatter = NumberFormat.currency(
      symbol: _symbols[currency] ?? '$currency ',
      decimalDigits: decimalPlaces,
      locale: 'fr_FR',
    );
    
    return formatter.format(value);
  }

  /// Format amount with compact notation (K, M, B)
  static String formatCompact(dynamic amount, String currency) {
    if (amount == null) return '${_symbols[currency] ?? currency} 0';
    
    double value = amount is num ? amount.toDouble() : 0.0;
    
    final formatter = NumberFormat.compactCurrency(
      symbol: _symbols[currency] ?? '$currency ',
      locale: 'fr_FR',
    );
    
    return formatter.format(value);
  }

  /// Format amount without symbol
  static String formatAmount(dynamic amount, {int decimals = 2}) {
    if (amount == null) return '0.00';
    
    double value = amount is num ? amount.toDouble() : 0.0;
    
    return value.toStringAsFixed(decimals);
  }

  /// Get currency symbol
  static String getSymbol(String currency) {
    return _symbols[currency] ?? currency;
  }

  /// Check if currency is crypto
  static bool _isCrypto(String currency) {
    const cryptoCurrencies = ['BTC', 'ETH', 'USDT', 'USDC', 'BNB', 'SOL', 'XRP'];
    return cryptoCurrencies.contains(currency.toUpperCase());
  }

  /// Parse amount from string
  static double parse(String amountString) {
    try {
      // Remove currency symbols and spaces
      String cleaned = amountString.replaceAll(RegExp(r'[^\d.,\-]'), '');
      // Handle European number format (comma as decimal separator)
      cleaned = cleaned.replaceAll(',', '.');
      return double.parse(cleaned);
    } catch (_) {
      return 0.0;
    }
  }
}
