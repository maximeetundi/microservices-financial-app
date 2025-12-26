// Get Transfer History Use Case

import '../../../../core/services/transfer_api_service.dart';

// Get Transfer History Use Case
class GetTransferHistoryUseCase {
  final TransferApiService _transferApi = TransferApiService();
  
  Future<List<TransferModel>> execute({int limit = 20, int offset = 0}) async {
    try {
      final transfers = await _transferApi.getTransfers(limit: limit, offset: offset);
      
      return transfers.map((t) => TransferModel(
        id: t['id']?.toString() ?? '',
        fromWalletId: t['from_wallet_id']?.toString() ?? '',
        toWalletId: t['to_wallet_id']?.toString(),
        toEmail: t['to_email']?.toString(),
        toPhone: t['to_phone']?.toString(),
        amount: (t['amount'] as num?)?.toDouble() ?? 0.0,
        currency: t['currency']?.toString() ?? 'USD',
        status: t['status']?.toString() ?? 'pending',
        description: t['description']?.toString(),
        createdAt: t['created_at'] != null 
            ? DateTime.tryParse(t['created_at'].toString()) ?? DateTime.now()
            : DateTime.now(),
      )).toList();
    } catch (e) {
      // Return empty list on API error instead of mock data
      return [];
    }
  }
}

class TransferModel {
  final String id;
  final String fromWalletId;
  final String? toWalletId;
  final String? toEmail;
  final String? toPhone;
  final double amount;
  final String currency;
  final String status;
  final String? description;
  final DateTime createdAt;

  TransferModel({
    required this.id,
    required this.fromWalletId,
    this.toWalletId,
    this.toEmail,
    this.toPhone,
    required this.amount,
    required this.currency,
    required this.status,
    this.description,
    required this.createdAt,
  });
}

