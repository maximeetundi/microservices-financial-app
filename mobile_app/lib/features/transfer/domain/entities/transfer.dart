import 'package:equatable/equatable.dart';

class Transfer extends Equatable {
  final String id;
  final String userId;
  final TransferType type;
  final String fromWallet;
  final String? toWallet;
  final String toAddress;
  final double amount;
  final double fee;
  final String currency;
  final TransferStatus status;
  final String? txHash;
  final String? memo;
  final String? recipientName;
  final String? bankCode;
  final String? recipientEmail;
  final DateTime createdAt;
  final DateTime? processedAt;
  final DateTime? completedAt;
  final String? failureReason;
  final int? estimatedConfirmationTime;

  const Transfer({
    required this.id,
    required this.userId,
    required this.type,
    required this.fromWallet,
    this.toWallet,
    required this.toAddress,
    required this.amount,
    this.fee = 0.0,
    required this.currency,
    required this.status,
    this.txHash,
    this.memo,
    this.recipientName,
    this.bankCode,
    this.recipientEmail,
    required this.createdAt,
    this.processedAt,
    this.completedAt,
    this.failureReason,
    this.estimatedConfirmationTime,
  });

  bool get isPending => status == TransferStatus.pending;
  bool get isProcessing => status == TransferStatus.processing;
  bool get isCompleted => status == TransferStatus.completed;
  bool get isFailed => status == TransferStatus.failed;
  bool get isCancelled => status == TransferStatus.cancelled;

  bool get isCryptoTransfer => type == TransferType.crypto;
  bool get isFiatTransfer => type == TransferType.fiat;
  bool get isInstantTransfer => type == TransferType.instant;

  String get statusText {
    switch (status) {
      case TransferStatus.pending:
        return 'Pending';
      case TransferStatus.processing:
        return 'Processing';
      case TransferStatus.completed:
        return 'Completed';
      case TransferStatus.failed:
        return 'Failed';
      case TransferStatus.cancelled:
        return 'Cancelled';
    }
  }

  String get typeText {
    switch (type) {
      case TransferType.crypto:
        return 'Crypto Transfer';
      case TransferType.fiat:
        return 'Bank Transfer';
      case TransferType.instant:
        return 'Instant Transfer';
      case TransferType.p2p:
        return 'P2P Transfer';
    }
  }

  String get displayAmount {
    return '${amount.toStringAsFixed(isCryptoTransfer ? 8 : 2)} $currency';
  }

  String get displayRecipient {
    if (recipientName != null) {
      return recipientName!;
    } else if (recipientEmail != null) {
      return recipientEmail!;
    } else if (toAddress.length > 20) {
      return '${toAddress.substring(0, 10)}...${toAddress.substring(toAddress.length - 6)}';
    } else {
      return toAddress;
    }
  }

  Duration? get estimatedTimeRemaining {
    if (estimatedConfirmationTime == null || isCompleted || isFailed || isCancelled) {
      return null;
    }

    final elapsedTime = DateTime.now().difference(createdAt);
    final totalTime = Duration(minutes: estimatedConfirmationTime!);
    final remaining = totalTime - elapsedTime;

    return remaining.isNegative ? null : remaining;
  }

  Transfer copyWith({
    String? id,
    String? userId,
    TransferType? type,
    String? fromWallet,
    String? toWallet,
    String? toAddress,
    double? amount,
    double? fee,
    String? currency,
    TransferStatus? status,
    String? txHash,
    String? memo,
    String? recipientName,
    String? bankCode,
    String? recipientEmail,
    DateTime? createdAt,
    DateTime? processedAt,
    DateTime? completedAt,
    String? failureReason,
    int? estimatedConfirmationTime,
  }) {
    return Transfer(
      id: id ?? this.id,
      userId: userId ?? this.userId,
      type: type ?? this.type,
      fromWallet: fromWallet ?? this.fromWallet,
      toWallet: toWallet ?? this.toWallet,
      toAddress: toAddress ?? this.toAddress,
      amount: amount ?? this.amount,
      fee: fee ?? this.fee,
      currency: currency ?? this.currency,
      status: status ?? this.status,
      txHash: txHash ?? this.txHash,
      memo: memo ?? this.memo,
      recipientName: recipientName ?? this.recipientName,
      bankCode: bankCode ?? this.bankCode,
      recipientEmail: recipientEmail ?? this.recipientEmail,
      createdAt: createdAt ?? this.createdAt,
      processedAt: processedAt ?? this.processedAt,
      completedAt: completedAt ?? this.completedAt,
      failureReason: failureReason ?? this.failureReason,
      estimatedConfirmationTime: estimatedConfirmationTime ?? this.estimatedConfirmationTime,
    );
  }

  @override
  List<Object?> get props => [
        id,
        userId,
        type,
        fromWallet,
        toWallet,
        toAddress,
        amount,
        fee,
        currency,
        status,
        txHash,
        memo,
        recipientName,
        bankCode,
        recipientEmail,
        createdAt,
        processedAt,
        completedAt,
        failureReason,
        estimatedConfirmationTime,
      ];
}

enum TransferType {
  crypto,
  fiat,
  instant,
  p2p,
}

enum TransferStatus {
  pending,
  processing,
  completed,
  failed,
  cancelled,
}

class TransferReceipt extends Equatable {
  final String transferId;
  final String confirmationCode;
  final DateTime timestamp;
  final TransferDetails details;
  final List<TransferStep> steps;

  const TransferReceipt({
    required this.transferId,
    required this.confirmationCode,
    required this.timestamp,
    required this.details,
    required this.steps,
  });

  @override
  List<Object> get props => [
        transferId,
        confirmationCode,
        timestamp,
        details,
        steps,
      ];
}

class TransferDetails extends Equatable {
  final String fromAccount;
  final String toAccount;
  final double amount;
  final double fee;
  final double exchangeRate;
  final String currency;
  final String? memo;

  const TransferDetails({
    required this.fromAccount,
    required this.toAccount,
    required this.amount,
    required this.fee,
    required this.exchangeRate,
    required this.currency,
    this.memo,
  });

  @override
  List<Object?> get props => [
        fromAccount,
        toAccount,
        amount,
        fee,
        exchangeRate,
        currency,
        memo,
      ];
}

class TransferStep extends Equatable {
  final String title;
  final String description;
  final bool isCompleted;
  final DateTime? completedAt;
  final bool isCurrent;

  const TransferStep({
    required this.title,
    required this.description,
    required this.isCompleted,
    this.completedAt,
    this.isCurrent = false,
  });

  @override
  List<Object?> get props => [
        title,
        description,
        isCompleted,
        completedAt,
        isCurrent,
      ];
}

class TransferLimit extends Equatable {
  final double dailyLimit;
  final double monthlyLimit;
  final double remainingDaily;
  final double remainingMonthly;
  final String currency;
  final DateTime resetTime;

  const TransferLimit({
    required this.dailyLimit,
    required this.monthlyLimit,
    required this.remainingDaily,
    required this.remainingMonthly,
    required this.currency,
    required this.resetTime,
  });

  bool canTransfer(double amount) {
    return amount <= remainingDaily && amount <= remainingMonthly;
  }

  double get maxTransferAmount {
    return remainingDaily < remainingMonthly ? remainingDaily : remainingMonthly;
  }

  @override
  List<Object> get props => [
        dailyLimit,
        monthlyLimit,
        remainingDaily,
        remainingMonthly,
        currency,
        resetTime,
      ];
}