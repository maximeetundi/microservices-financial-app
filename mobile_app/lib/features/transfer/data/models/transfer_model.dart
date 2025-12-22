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

  factory TransferModel.fromJson(Map<String, dynamic> json) {
    return TransferModel(
      id: json['id'] ?? '',
      fromWalletId: json['from_wallet_id'] ?? '',
      toWalletId: json['to_wallet_id'],
      toEmail: json['to_email'],
      toPhone: json['to_phone'],
      amount: (json['amount'] ?? 0).toDouble(),
      currency: json['currency'] ?? 'USD',
      status: json['status'] ?? 'pending',
      description: json['description'],
      createdAt: json['created_at'] != null 
          ? DateTime.parse(json['created_at']) 
          : DateTime.now(),
    );
  }

  Map<String, dynamic> toJson() {
    return {
      'id': id,
      'from_wallet_id': fromWalletId,
      'to_wallet_id': toWalletId,
      'to_email': toEmail,
      'to_phone': toPhone,
      'amount': amount,
      'currency': currency,
      'status': status,
      'description': description,
      'created_at': createdAt.toIso8601String(),
    };
  }
}
