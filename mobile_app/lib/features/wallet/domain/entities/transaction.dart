import 'package:equatable/equatable.dart';

class Transaction extends Equatable {
  final String id;
  final String walletId;
  final String? fromAddress;
  final String? toAddress;
  final double amount;
  final double fee;
  final String currency;
  final TransactionType type;
  final TransactionStatus status;
  final String? txHash;
  final String? memo;
  final DateTime createdAt;
  final DateTime? confirmedAt;
  final int confirmations;
  final int requiredConfirmations;

  const Transaction({
    required this.id,
    required this.walletId,
    this.fromAddress,
    this.toAddress,
    required this.amount,
    this.fee = 0.0,
    required this.currency,
    required this.type,
    required this.status,
    this.txHash,
    this.memo,
    required this.createdAt,
    this.confirmedAt,
    this.confirmations = 0,
    this.requiredConfirmations = 1,
  });
  
  factory Transaction.fromJson(Map<String, dynamic> json) {
    return Transaction(
      id: json['id'] ?? json['transaction_id'] ?? '',
      walletId: json['wallet_id'] ?? '',
      fromAddress: json['from_address'] ?? json['source_address'],
      toAddress: json['to_address'] ?? json['destination_address'],
      amount: (json['amount'] as num?)?.toDouble() ?? 0.0,
      fee: (json['fee'] as num?)?.toDouble() ?? 0.0,
      currency: json['currency'] ?? '',
      type: _parseTransactionType(json['type'] ?? json['transaction_type']),
      status: _parseTransactionStatus(json['status']),
      txHash: json['tx_hash'] ?? json['hash'],
      memo: json['memo'] ?? json['description'],
      createdAt: json['created_at'] != null 
          ? DateTime.parse(json['created_at']) 
          : DateTime.now(),
      confirmedAt: json['confirmed_at'] != null 
          ? DateTime.parse(json['confirmed_at']) 
          : null,
      confirmations: json['confirmations'] ?? 0,
      requiredConfirmations: json['required_confirmations'] ?? 1,
    );
  }
  
  static TransactionType _parseTransactionType(String? type) {
    switch (type?.toLowerCase()) {
      case 'send':
      case 'withdrawal':
        return TransactionType.send;
      case 'receive':
      case 'deposit':
        return TransactionType.receive;
      case 'exchange':
      case 'swap':
        return TransactionType.exchange;
      case 'buy':
        return TransactionType.buy;
      case 'sell':
        return TransactionType.sell;
      case 'fee':
        return TransactionType.fee;
      case 'reward':
        return TransactionType.reward;
      case 'staking':
        return TransactionType.staking;
      default:
        return TransactionType.send;
    }
  }
  
  static TransactionStatus _parseTransactionStatus(String? status) {
    switch (status?.toLowerCase()) {
      case 'confirmed':
      case 'completed':
      case 'success':
        return TransactionStatus.confirmed;
      case 'failed':
      case 'error':
        return TransactionStatus.failed;
      case 'cancelled':
      case 'canceled':
        return TransactionStatus.cancelled;
      default:
        return TransactionStatus.pending;
    }
  }

  bool get isIncoming => type == TransactionType.receive;

  bool get isOutgoing => type == TransactionType.send;

  bool get isConfirmed => status == TransactionStatus.confirmed;

  bool get isPending => status == TransactionStatus.pending;

  bool get isFailed => status == TransactionStatus.failed;

  String get displayAmount {
    final prefix = isIncoming ? '+' : '-';
    return '$prefix${amount.toStringAsFixed(8)} $currency';
  }

  String get statusText {
    switch (status) {
      case TransactionStatus.pending:
        return 'En attente';
      case TransactionStatus.confirmed:
        return 'Confirmé';
      case TransactionStatus.failed:
        return 'Échoué';
      case TransactionStatus.cancelled:
        return 'Annulé';
    }
  }

  @override
  List<Object?> get props => [
        id,
        walletId,
        fromAddress,
        toAddress,
        amount,
        fee,
        currency,
        type,
        status,
        txHash,
        memo,
        createdAt,
        confirmedAt,
        confirmations,
        requiredConfirmations,
      ];
}

enum TransactionType {
  send,
  receive,
  exchange,
  buy,
  sell,
  fee,
  reward,
  staking,
}

enum TransactionStatus {
  pending,
  confirmed,
  failed,
  cancelled,
}
