import '../api/api_client.dart';
import '../api/api_endpoints.dart';

/// Service d'échange
class ExchangeApiService {
  final ApiClient _client = ApiClient();
  
  /// Récupérer tous les taux de change
  Future<Map<String, dynamic>> getRates() async {
    final response = await _client.get(ApiEndpoints.exchangeRates);
    if (response.statusCode == 200) {
      return response.data['rates'];
    }
    throw Exception('Failed to get rates');
  }
  
  /// Récupérer le taux entre deux devises
  Future<Map<String, dynamic>> getExchangeRate(String fromCurrency, String toCurrency) async {
    final response = await _client.get(
      ApiEndpoints.exchangePair(fromCurrency, toCurrency),
    );
    if (response.statusCode == 200) {
      return response.data;
    }
    throw Exception('Failed to get exchange rate');
  }
  
  /// Exécuter un échange
  Future<Map<String, dynamic>> executeExchange({
    required String fromWalletId,
    required String toWalletId,
    required double amount,
    required String fromCurrency,
    required String toCurrency,
  }) async {
    final response = await _client.post(
      ApiEndpoints.executeExchange,
      data: {
        'from_wallet_id': fromWalletId,
        'to_wallet_id': toWalletId,
        'amount': amount,
        'from_currency': fromCurrency,
        'to_currency': toCurrency,
      },
    );
    if (response.statusCode == 200) {
      return response.data;
    }
    throw Exception(response.data['error'] ?? 'Exchange failed');
  }

  /// Obtenir un devis (Quote)
  Future<Map<String, dynamic>> getQuote({
    required String fromCurrency,
    required String toCurrency,
    double? fromAmount,
    double? toAmount,
  }) async {
    final data = {
      'from_currency': fromCurrency,
      'to_currency': toCurrency,
      if (fromAmount != null) 'from_amount': fromAmount,
      if (toAmount != null) 'to_amount': toAmount,
    };

    final response = await _client.post(
      ApiEndpoints.quote,
      data: data,
    );

    if (response.statusCode == 200) {
      return response.data;
    }
    throw Exception(response.data['error'] ?? 'Failed to get quote');
  }

  /// Exécuter un devis (Quote)
  Future<Map<String, dynamic>> executeQuote({
    required String quoteId,
    required String fromWalletId,
    required String toWalletId,
  }) async {
    final response = await _client.post(
      ApiEndpoints.executeExchange,
      data: {
        'quote_id': quoteId,
        'from_wallet_id': fromWalletId,
        'to_wallet_id': toWalletId,
      },
    );

    if (response.statusCode == 200) {
      return response.data;
    }
    throw Exception(response.data['error'] ?? 'Failed to execute quote');
  }
  
  /// Récupérer l'historique des échanges
  Future<List<Map<String, dynamic>>> getExchangeHistory({
    int limit = 50,
    int offset = 0,
  }) async {
    final response = await _client.get(
      ApiEndpoints.exchangeHistory,
      queryParameters: {
        'limit': limit,
        'offset': offset,
      },
    );
    if (response.statusCode == 200) {
      return List<Map<String, dynamic>>.from(response.data['exchanges']);
    }
    throw Exception('Failed to get exchange history');
  }
  
  /// Récupérer un échange par ID
  Future<Map<String, dynamic>> getExchange(String exchangeId) async {
    final response = await _client.get(ApiEndpoints.exchangeById(exchangeId));
    if (response.statusCode == 200) {
      return response.data['exchange'];
    }
    throw Exception('Exchange not found');
  }
  
  /// Calculer le montant de sortie estimé
  Future<double> calculateExchange({
    required double amount,
    required String fromCurrency,
    required String toCurrency,
  }) async {
    final rate = await getExchangeRate(fromCurrency, toCurrency);
    final exchangeRate = (rate['rate'] as num).toDouble();
    final fee = (rate['fee'] as num?)?.toDouble() ?? 0;
    
    return (amount * exchangeRate) - fee;
  }
}
