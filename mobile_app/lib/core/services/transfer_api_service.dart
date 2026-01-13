import '../api/api_client.dart';
import '../api/api_endpoints.dart';

/// Service des transferts
class TransferApiService {
  final ApiClient _client = ApiClient();
  
  /// Récupérer l'historique des transferts
  Future<List<Map<String, dynamic>>> getTransfers({
    int limit = 50,
    int offset = 0,
  }) async {
    final response = await _client.get(
      ApiEndpoints.transfersList,
      queryParameters: {
        'limit': limit,
        'offset': offset,
      },
    );
    if (response.statusCode == 200) {
      return List<Map<String, dynamic>>.from(response.data['transfers']);
    }
    throw Exception('Failed to get transfers');
  }
  
  /// Créer un nouveau transfert
  /// Créer un nouveau transfert
  Future<Map<String, dynamic>> createTransfer({
    required String fromWalletId,
    String? toWalletId,
    String? toEmail,
    String? toPhone,
    required double amount,
    required String currency,
    String? description,
  }) async {
    final data = {
      'from_wallet_id': fromWalletId,
      'amount': amount,
      'currency': currency,
      if (toWalletId != null) 'to_wallet_id': toWalletId,
      if (toEmail != null) 'to_email': toEmail,
      if (toPhone != null) 'to_phone': toPhone,
      if (description != null) 'description': description,
    };

    final response = await _client.post(
      ApiEndpoints.createTransfer,
      data: data,
    );
    if (response.statusCode == 201 || response.statusCode == 200) {
      return response.data['transfer'];
    }
    throw Exception(response.data['error'] ?? 'Transfer failed');
  }
  
  /// Récupérer un transfert par ID
  Future<Map<String, dynamic>> getTransfer(String transferId) async {
    final response = await _client.get(ApiEndpoints.transferById(transferId));
    if (response.statusCode == 200) {
      return response.data['transfer'];
    }
    throw Exception('Transfer not found');
  }
  
  /// Annuler un transfert
  /// Annuler un transfert
  Future<void> cancelTransfer(String transferId, {String? reason}) async {
    final response = await _client.post(
       ApiEndpoints.cancelTransfer(transferId),
       data: reason != null ? {'reason': reason} : {},
    );
    if (response.statusCode != 200) {
      throw Exception(response.data['error'] ?? 'Failed to cancel transfer');
    }
  }

  /// Inverser (Rembourser) un transfert
  Future<void> reverseTransfer(String transferId, {String? reason}) async {
    final response = await _client.post(
       '${ApiEndpoints.baseTransfer}/$transferId/reverse', // Assuming endpoint naming convention
       data: reason != null ? {'reason': reason} : {},
    );
    if (response.statusCode != 200) {
      throw Exception(response.data['error'] ?? 'Failed to reverse transfer');
    }
  }
  
  /// Transfert international
  Future<Map<String, dynamic>> createInternationalTransfer({
    required String fromWalletId,
    required double amount,
    required String currency,
    required String recipientName,
    required String recipientBank,
    required String recipientIBAN,
    required String recipientCountry,
    String? swiftCode,
    String? reference,
  }) async {
    final response = await _client.post(
      ApiEndpoints.internationalTransfer,
      data: {
        'from_wallet_id': fromWalletId,
        'amount': amount,
        'currency': currency,
        'recipient_name': recipientName,
        'recipient_bank': recipientBank,
        'recipient_iban': recipientIBAN,
        'recipient_country': recipientCountry,
        if (swiftCode != null) 'swift_code': swiftCode,
        if (reference != null) 'reference': reference,
      },
    );
    if (response.statusCode == 200 || response.statusCode == 201) {
      return response.data;
    }
    throw Exception(response.data['error'] ?? 'International transfer failed');
  }
  
  /// Récupérer les fournisseurs Mobile Money
  Future<List<String>> getMobileMoneyProviders() async {
    final response = await _client.get(ApiEndpoints.mobileMoneyProviders);
    if (response.statusCode == 200) {
      return List<String>.from(response.data['providers']);
    }
    throw Exception('Failed to get providers');
  }
  
  /// Envoyer via Mobile Money
  Future<Map<String, dynamic>> sendMobileMoney({
    required String fromWalletId,
    required String provider,
    required String phoneNumber,
    required double amount,
    required String currency,
  }) async {
    final response = await _client.post(
      ApiEndpoints.sendMobileMoney,
      data: {
        'from_wallet_id': fromWalletId,
        'provider': provider,
        'phone_number': phoneNumber,
        'amount': amount,
        'currency': currency,
      },
    );
    if (response.statusCode == 200) {
      return response.data;
    }
    throw Exception(response.data['error'] ?? 'Mobile money transfer failed');
  }
  
  /// Recevoir via Mobile Money
  Future<Map<String, dynamic>> receiveMobileMoney({
    required String toWalletId,
    required String provider,
    required String phoneNumber,
    required double amount,
    required String currency,
  }) async {
    final response = await _client.post(
      ApiEndpoints.receiveMobileMoney,
      data: {
        'to_wallet_id': toWalletId,
        'provider': provider,
        'phone_number': phoneNumber,
        'amount': amount,
        'currency': currency,
      },
    );
    if (response.statusCode == 200) {
      return response.data;
    }
    throw Exception(response.data['error'] ?? 'Mobile money receive failed');
  }
}
