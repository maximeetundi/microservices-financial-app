import '../../domain/entities/transfer.dart';
import '../../../../core/api/api_client.dart';
import '../../../../core/api/api_endpoints.dart';

/// Repository implementation for transfer operations
class TransferRepositoryImpl {
  final ApiClient _apiClient = ApiClient();

  /// Send a transfer
  Future<Transfer> sendTransfer({
    required String type,
    required String fromWalletId,
    required String recipient,
    required double amount,
    String? memo,
    String? bankCode,
    String? recipientName,
  }) async {
    final response = await _apiClient.post(
      ApiEndpoints.transfers,
      data: {
        'type': type,
        'from_wallet_id': fromWalletId,
        'recipient': recipient,
        'amount': amount,
        if (memo != null) 'description': memo,
        if (bankCode != null) 'bank_code': bankCode,
        if (recipientName != null) 'recipient_name': recipientName,
      },
    );
    return Transfer.fromJson(response.data);
  }

  /// Get transfer history
  Future<List<Transfer>> getTransferHistory({int limit = 20, int offset = 0}) async {
    final response = await _apiClient.get(
      '${ApiEndpoints.transfers}?limit=$limit&offset=$offset',
    );
    
    final List<dynamic> data = response.data['transfers'] ?? [];
    return data.map((json) => Transfer.fromJson(json)).toList();
  }

  /// Get transfer by ID
  Future<Transfer> getTransferById(String transferId) async {
    final response = await _apiClient.get('${ApiEndpoints.transfers}/$transferId');
    return Transfer.fromJson(response.data);
  }

  /// Estimate transfer fee
  Future<Map<String, dynamic>> estimateTransferFee({
    required String type,
    required double amount,
    required String currency,
  }) async {
    final response = await _apiClient.post(
      '${ApiEndpoints.transfers}/estimate',
      data: {
        'type': type,
        'amount': amount,
        'currency': currency,
      },
    );
    return response.data;
  }
}
