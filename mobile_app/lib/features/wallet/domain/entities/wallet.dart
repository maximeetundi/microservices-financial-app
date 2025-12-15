import 'package:equatable/equatable.dart';

class Wallet extends Equatable {
  final String id;
  final String userId;
  final String currency;
  final String? name;
  final double balance;
  final double availableBalance;
  final double pendingBalance;
  final String address;
  final String? privateKey;
  final WalletType type;
  final WalletStatus status;
  final DateTime createdAt;
  final DateTime updatedAt;
  final double usdRate;
  final double dailyChange;
  final WalletNetworkInfo networkInfo;

  const Wallet({
    required this.id,
    required this.userId,
    required this.currency,
    this.name,
    required this.balance,
    required this.availableBalance,
    this.pendingBalance = 0.0,
    required this.address,
    this.privateKey,
    required this.type,
    this.status = WalletStatus.active,
    required this.createdAt,
    required this.updatedAt,
    this.usdRate = 1.0,
    this.dailyChange = 0.0,
    required this.networkInfo,
  });
  
  factory Wallet.fromJson(Map<String, dynamic> json) {
    return Wallet(
      id: json['id'] ?? json['wallet_id'] ?? '',
      userId: json['user_id'] ?? '',
      currency: json['currency'] ?? '',
      name: json['name'],
      balance: (json['balance'] as num?)?.toDouble() ?? 0.0,
      availableBalance: (json['available_balance'] as num?)?.toDouble() ?? 
          (json['balance'] as num?)?.toDouble() ?? 0.0,
      pendingBalance: (json['pending_balance'] as num?)?.toDouble() ?? 0.0,
      address: json['address'] ?? '',
      privateKey: json['private_key'],
      type: _parseWalletType(json['wallet_type'] ?? json['type']),
      status: _parseWalletStatus(json['status']),
      createdAt: json['created_at'] != null 
          ? DateTime.parse(json['created_at']) 
          : DateTime.now(),
      updatedAt: json['updated_at'] != null 
          ? DateTime.parse(json['updated_at']) 
          : DateTime.now(),
      usdRate: (json['usd_rate'] as num?)?.toDouble() ?? 1.0,
      dailyChange: (json['daily_change'] as num?)?.toDouble() ?? 0.0,
      networkInfo: json['network_info'] != null 
          ? WalletNetworkInfo.fromJson(json['network_info'])
          : const WalletNetworkInfo(network: 'mainnet'),
    );
  }
  
  static WalletType _parseWalletType(String? type) {
    switch (type?.toLowerCase()) {
      case 'fiat':
        return WalletType.fiat;
      case 'stablecoin':
        return WalletType.stablecoin;
      default:
        return WalletType.crypto;
    }
  }
  
  static WalletStatus _parseWalletStatus(String? status) {
    switch (status?.toLowerCase()) {
      case 'frozen':
        return WalletStatus.frozen;
      case 'suspended':
        return WalletStatus.suspended;
      case 'closed':
        return WalletStatus.closed;
      default:
        return WalletStatus.active;
    }
  }

  String get displayName => name ?? '${currency} Wallet';

  double get totalValue => balance * usdRate;

  bool get isCrypto => type == WalletType.crypto;

  bool get isFiat => type == WalletType.fiat;

  bool get canSend => status == WalletStatus.active && availableBalance > 0;

  bool get canReceive => status == WalletStatus.active;

  String get formattedBalance {
    if (isFiat) {
      return '\$${balance.toStringAsFixed(2)}';
    } else {
      return '${balance.toStringAsFixed(8)} $currency';
    }
  }

  String get formattedAvailableBalance {
    if (isFiat) {
      return '\$${availableBalance.toStringAsFixed(2)}';
    } else {
      return '${availableBalance.toStringAsFixed(8)} $currency';
    }
  }

  Wallet copyWith({
    String? id,
    String? userId,
    String? currency,
    String? name,
    double? balance,
    double? availableBalance,
    double? pendingBalance,
    String? address,
    String? privateKey,
    WalletType? type,
    WalletStatus? status,
    DateTime? createdAt,
    DateTime? updatedAt,
    double? usdRate,
    double? dailyChange,
    WalletNetworkInfo? networkInfo,
  }) {
    return Wallet(
      id: id ?? this.id,
      userId: userId ?? this.userId,
      currency: currency ?? this.currency,
      name: name ?? this.name,
      balance: balance ?? this.balance,
      availableBalance: availableBalance ?? this.availableBalance,
      pendingBalance: pendingBalance ?? this.pendingBalance,
      address: address ?? this.address,
      privateKey: privateKey ?? this.privateKey,
      type: type ?? this.type,
      status: status ?? this.status,
      createdAt: createdAt ?? this.createdAt,
      updatedAt: updatedAt ?? this.updatedAt,
      usdRate: usdRate ?? this.usdRate,
      dailyChange: dailyChange ?? this.dailyChange,
      networkInfo: networkInfo ?? this.networkInfo,
    );
  }

  @override
  List<Object?> get props => [
        id,
        userId,
        currency,
        name,
        balance,
        availableBalance,
        pendingBalance,
        address,
        privateKey,
        type,
        status,
        createdAt,
        updatedAt,
        usdRate,
        dailyChange,
        networkInfo,
      ];
}

enum WalletType {
  crypto,
  fiat,
  stablecoin,
}

enum WalletStatus {
  active,
  frozen,
  suspended,
  closed,
}

class WalletNetworkInfo extends Equatable {
  final String network;
  final String? contractAddress;
  final int confirmations;
  final double networkFee;
  final String feeUnit;
  final int blockHeight;
  final Duration averageBlockTime;

  const WalletNetworkInfo({
    required this.network,
    this.contractAddress,
    this.confirmations = 1,
    this.networkFee = 0.0,
    this.feeUnit = '',
    this.blockHeight = 0,
    this.averageBlockTime = const Duration(minutes: 10),
  });
  
  factory WalletNetworkInfo.fromJson(Map<String, dynamic> json) {
    return WalletNetworkInfo(
      network: json['network'] ?? 'mainnet',
      contractAddress: json['contract_address'],
      confirmations: json['confirmations'] ?? 1,
      networkFee: (json['network_fee'] as num?)?.toDouble() ?? 0.0,
      feeUnit: json['fee_unit'] ?? '',
      blockHeight: json['block_height'] ?? 0,
      averageBlockTime: Duration(
        seconds: json['average_block_time_seconds'] ?? 600,
      ),
    );
  }

  @override
  List<Object?> get props => [
        network,
        contractAddress,
        confirmations,
        networkFee,
        feeUnit,
        blockHeight,
        averageBlockTime,
      ];
}

class WalletsData extends Equatable {
  final List<Wallet> wallets;
  final List<Transaction> recentTransactions;
  final double totalValue;
  final double dailyChange;

  const WalletsData({
    required this.wallets,
    required this.recentTransactions,
    required this.totalValue,
    required this.dailyChange,
  });

  @override
  List<Object> get props => [wallets, recentTransactions, totalValue, dailyChange];
}

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
        return 'Pending';
      case TransactionStatus.confirmed:
        return 'Confirmed';
      case TransactionStatus.failed:
        return 'Failed';
      case TransactionStatus.cancelled:
        return 'Cancelled';
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