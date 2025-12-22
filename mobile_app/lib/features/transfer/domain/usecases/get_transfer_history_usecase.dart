// Get Transfer History Use Case (no external import needed - TransferModel is defined below)

// Get Transfer History Use Case
class GetTransferHistoryUseCase {
  Future<List<TransferModel>> execute({int limit = 20, int offset = 0}) async {
    // This would call the API in a real implementation
    await Future.delayed(const Duration(milliseconds: 500));
    
    return [
      TransferModel(
        id: 'TRF-001',
        fromWalletId: 'wallet-1',
        toEmail: 'john@example.com',
        amount: 500.0,
        currency: 'USD',
        status: 'completed',
        createdAt: DateTime.now().subtract(const Duration(hours: 2)),
      ),
      TransferModel(
        id: 'TRF-002',
        fromWalletId: 'wallet-1',
        toEmail: 'jane@example.com',
        amount: 250.0,
        currency: 'EUR',
        status: 'completed',
        createdAt: DateTime.now().subtract(const Duration(days: 1)),
      ),
    ];
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
