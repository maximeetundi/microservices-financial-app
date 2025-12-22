// Send Transfer Use Case
class SendTransferUseCase {
  Future<TransferResult> execute(SendTransferParams params) async {
    // This would call the API in a real implementation
    await Future.delayed(const Duration(seconds: 1));
    
    return TransferResult(
      id: 'TRF-${DateTime.now().millisecondsSinceEpoch}',
      status: 'completed',
      amount: params.amount,
      currency: params.currency,
      recipientEmail: params.recipientEmail,
      createdAt: DateTime.now(),
    );
  }
}

class SendTransferParams {
  final String fromWalletId;
  final String? recipientEmail;
  final String? recipientPhone;
  final String? toWalletId;
  final double amount;
  final String currency;
  final String? description;

  SendTransferParams({
    required this.fromWalletId,
    this.recipientEmail,
    this.recipientPhone,
    this.toWalletId,
    required this.amount,
    required this.currency,
    this.description,
  });
}

class TransferResult {
  final String id;
  final String status;
  final double amount;
  final String currency;
  final String? recipientEmail;
  final DateTime createdAt;

  TransferResult({
    required this.id,
    required this.status,
    required this.amount,
    required this.currency,
    this.recipientEmail,
    required this.createdAt,
  });
}
