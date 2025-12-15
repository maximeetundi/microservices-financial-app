import '../api/api_client.dart';
import '../api/api_endpoints.dart';

/// Service des portefeuilles
class WalletApiService {
  final ApiClient _client = ApiClient();
  
  /// Récupérer tous les portefeuilles de l'utilisateur
  Future<List<Map<String, dynamic>>> getWallets() async {
    final response = await _client.get(ApiEndpoints.walletsList);
    if (response.statusCode == 200) {
      return List<Map<String, dynamic>>.from(response.data['wallets']);
    }
    throw Exception('Failed to get wallets');
  }
  
  /// Créer un nouveau portefeuille
  Future<Map<String, dynamic>> createWallet({
    required String currency,
    required String walletType,
    String? name,
  }) async {
    final response = await _client.post(
      ApiEndpoints.createWallet,
      data: {
        'currency': currency,
        'wallet_type': walletType,
        if (name != null) 'name': name,
      },
    );
    if (response.statusCode == 201) {
      return response.data['wallet'];
    }
    throw Exception(response.data['error'] ?? 'Failed to create wallet');
  }
  
  /// Récupérer un portefeuille par ID
  Future<Map<String, dynamic>> getWallet(String walletId) async {
    final response = await _client.get(ApiEndpoints.walletById(walletId));
    if (response.statusCode == 200) {
      return response.data['wallet'];
    }
    throw Exception('Wallet not found');
  }
  
  /// Récupérer le solde d'un portefeuille
  Future<double> getBalance(String walletId) async {
    final response = await _client.get(ApiEndpoints.walletBalance(walletId));
    if (response.statusCode == 200) {
      return (response.data['balance'] as num).toDouble();
    }
    throw Exception('Failed to get balance');
  }
  
  /// Récupérer les transactions d'un portefeuille
  Future<List<Map<String, dynamic>>> getTransactions(
    String walletId, {
    int limit = 50,
    int offset = 0,
    String? type,
  }) async {
    final response = await _client.get(
      ApiEndpoints.walletTransactions(walletId),
      queryParameters: {
        'limit': limit,
        'offset': offset,
        if (type != null) 'type': type,
      },
    );
    if (response.statusCode == 200) {
      return List<Map<String, dynamic>>.from(response.data['transactions']);
    }
    throw Exception('Failed to get transactions');
  }
  
  /// Déposer des fonds
  Future<Map<String, dynamic>> deposit({
    required String walletId,
    required double amount,
    String? paymentMethod,
  }) async {
    final response = await _client.post(
      ApiEndpoints.walletDeposit(walletId),
      data: {
        'amount': amount,
        if (paymentMethod != null) 'payment_method': paymentMethod,
      },
    );
    if (response.statusCode == 200) {
      return response.data;
    }
    throw Exception(response.data['error'] ?? 'Deposit failed');
  }
  
  /// Retirer des fonds
  Future<Map<String, dynamic>> withdraw({
    required String walletId,
    required double amount,
    required String destinationAddress,
    String? network,
  }) async {
    final response = await _client.post(
      ApiEndpoints.walletWithdraw(walletId),
      data: {
        'amount': amount,
        'destination_address': destinationAddress,
        if (network != null) 'network': network,
      },
    );
    if (response.statusCode == 200) {
      return response.data;
    }
    throw Exception(response.data['error'] ?? 'Withdrawal failed');
  }
}
