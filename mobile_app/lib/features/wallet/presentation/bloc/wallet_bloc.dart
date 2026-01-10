import 'package:flutter_bloc/flutter_bloc.dart';
import 'package:equatable/equatable.dart';

import '../../../../core/services/api_service.dart';
import '../../../../core/services/secure_storage_service.dart';
import '../../domain/entities/wallet.dart';
import '../../domain/entities/transaction.dart' as tx;

// Events
abstract class WalletEvent extends Equatable {
  const WalletEvent();

  @override
  List<Object?> get props => [];
}

class LoadWalletsEvent extends WalletEvent {
  const LoadWalletsEvent();
}

class CreateWalletEvent extends WalletEvent {
  final String currency;
  final String walletType;
  final String? name;

  const CreateWalletEvent({
    required this.currency,
    this.walletType = 'crypto',
    this.name,
  });

  @override
  List<Object?> get props => [currency, walletType, name];
}

class SendTransactionEvent extends WalletEvent {
  final String walletId;
  final String toAddress;
  final double amount;
  final String? memo;

  const SendTransactionEvent({
    required this.walletId,
    required this.toAddress,
    required this.amount,
    this.memo,
  });

  @override
  List<Object?> get props => [walletId, toAddress, amount, memo];
}

class DepositEvent extends WalletEvent {
  final String walletId;
  final double amount;
  final String? paymentMethod;

  const DepositEvent({
    required this.walletId,
    required this.amount,
    this.paymentMethod,
  });

  @override
  List<Object?> get props => [walletId, amount, paymentMethod];
}

class LoadWalletTransactionsEvent extends WalletEvent {
  final String walletId;
  final int limit;
  final int offset;

  const LoadWalletTransactionsEvent({
    required this.walletId,
    this.limit = 50,
    this.offset = 0,
  });

  @override
  List<Object> get props => [walletId, limit, offset];
}

class RefreshWalletEvent extends WalletEvent {
  final String walletId;

  const RefreshWalletEvent({required this.walletId});

  @override
  List<Object> get props => [walletId];
}

// States
abstract class WalletState extends Equatable {
  const WalletState();

  @override
  List<Object?> get props => [];
}

class WalletInitialState extends WalletState {
  const WalletInitialState();
}

class WalletLoadingState extends WalletState {
  const WalletLoadingState();
}

class WalletLoadedState extends WalletState {
  final List<Wallet> wallets;
  final double totalValue;
  final double dailyChange;
  final List<tx.Transaction> recentTransactions;

  const WalletLoadedState({
    required this.wallets,
    required this.totalValue,
    required this.dailyChange,
    required this.recentTransactions,
  });

  @override
  List<Object> get props => [wallets, totalValue, dailyChange, recentTransactions];

  WalletLoadedState copyWith({
    List<Wallet>? wallets,
    double? totalValue,
    double? dailyChange,
    List<tx.Transaction>? recentTransactions,
  }) {
    return WalletLoadedState(
      wallets: wallets ?? this.wallets,
      totalValue: totalValue ?? this.totalValue,
      dailyChange: dailyChange ?? this.dailyChange,
      recentTransactions: recentTransactions ?? this.recentTransactions,
    );
  }
}

class WalletErrorState extends WalletState {
  final String message;

  const WalletErrorState({required this.message});

  @override
  List<Object> get props => [message];
}

class WalletTransactionSentState extends WalletState {
  final String transactionId;

  const WalletTransactionSentState({required this.transactionId});

  @override
  List<Object> get props => [transactionId];
}

class WalletCreatedState extends WalletState {
  final Wallet wallet;

  const WalletCreatedState({required this.wallet});

  @override
  List<Object> get props => [wallet];
}

// BLoC
class WalletBloc extends Bloc<WalletEvent, WalletState> {
  final ApiService _apiService;
  final SecureStorageService _secureStorage; // Add secure storage

  WalletBloc({
    required ApiService apiService,
    required SecureStorageService secureStorage, // Add to constructor
  })  : _apiService = apiService,
        _secureStorage = secureStorage,
        super(const WalletInitialState()) {
    on<LoadWalletsEvent>(_onLoadWallets);
    on<CreateWalletEvent>(_onCreateWallet);
    on<SendTransactionEvent>(_onSendTransaction);
    on<DepositEvent>(_onDeposit);
    on<LoadWalletTransactionsEvent>(_onLoadWalletTransactions);
    on<RefreshWalletEvent>(_onRefreshWallet);
  }

  Future<void> _onLoadWallets(
    LoadWalletsEvent event,
    Emitter<WalletState> emit,
  ) async {
    emit(const WalletLoadingState());

    try {
      // 1. Try local cache first
      final cachedWallets = await _secureStorage.getWallets();
      if (cachedWallets != null) {
        try {
          final wallets = cachedWallets.map((w) => Wallet.fromJson(w as Map<String, dynamic>)).toList();
          if (wallets.isNotEmpty) {
             final totalValue = wallets.fold<double>(
              0, 
              (sum, wallet) => sum + (wallet.balance * wallet.usdRate),
            );
            emit(WalletLoadedState(
              wallets: wallets,
              totalValue: totalValue,
              dailyChange: 0,
              recentTransactions: const [],
            ));
          }
        } catch (_) {}
      }

      final walletsData = await _apiService.wallet.getWallets();
      
      // Cache fresh data
      await _secureStorage.saveWallets(walletsData);

      final wallets = walletsData.map((w) => Wallet.fromJson(w)).toList();
      
      final totalValue = wallets.fold<double>(
        0, 
        (sum, wallet) => sum + (wallet.balance * wallet.usdRate),
      );
      
      // Get recent transactions from first wallet (or empty)
      List<tx.Transaction> recentTransactions = [];
      if (wallets.isNotEmpty) {
        try {
          final txData = await _apiService.wallet.getTransactions(
            wallets.first.id,
            limit: 10,
          );
          recentTransactions = txData.map((t) => tx.Transaction.fromJson(t)).toList();
        } catch (_) {}
      }
      
      emit(WalletLoadedState(
        wallets: wallets,
        totalValue: totalValue,
        dailyChange: 0, // TODO: Calculate from API
        recentTransactions: recentTransactions,
      ));
    } catch (e) {
      if (state is! WalletLoadedState) {
         emit(WalletErrorState(message: _getErrorMessage(e)));
      }
    }
  }

  Future<void> _onCreateWallet(
    CreateWalletEvent event,
    Emitter<WalletState> emit,
  ) async {
    try {
      final walletData = await _apiService.wallet.createWallet(
        currency: event.currency,
        walletType: event.walletType,
        name: event.name,
      );
      
      final wallet = Wallet.fromJson(walletData);
      emit(WalletCreatedState(wallet: wallet));
      
      // Reload wallets
      add(const LoadWalletsEvent());
    } catch (e) {
      emit(WalletErrorState(message: _getErrorMessage(e)));
    }
  }

  Future<void> _onSendTransaction(
    SendTransactionEvent event,
    Emitter<WalletState> emit,
  ) async {
    try {
      final result = await _apiService.wallet.withdraw(
        walletId: event.walletId,
        amount: event.amount,
        destinationAddress: event.toAddress,
      );
      
      emit(WalletTransactionSentState(
        transactionId: result['transaction_id'] ?? '',
      ));
      
      // Reload wallets
      add(const LoadWalletsEvent());
    } catch (e) {
      emit(WalletErrorState(message: _getErrorMessage(e)));
    }
  }
  
  Future<void> _onDeposit(
    DepositEvent event,
    Emitter<WalletState> emit,
  ) async {
    try {
      await _apiService.wallet.deposit(
        walletId: event.walletId,
        amount: event.amount,
        paymentMethod: event.paymentMethod,
      );
      
      // Reload wallets
      add(const LoadWalletsEvent());
    } catch (e) {
      emit(WalletErrorState(message: _getErrorMessage(e)));
    }
  }

  Future<void> _onLoadWalletTransactions(
    LoadWalletTransactionsEvent event,
    Emitter<WalletState> emit,
  ) async {
    try {
      final txData = await _apiService.wallet.getTransactions(
        event.walletId,
        limit: event.limit,
        offset: event.offset,
      );
      
      final transactions = txData.map((t) => tx.Transaction.fromJson(t)).toList();
      
      if (state is WalletLoadedState) {
        emit((state as WalletLoadedState).copyWith(
          recentTransactions: transactions,
        ));
      }
    } catch (e) {
      emit(WalletErrorState(message: _getErrorMessage(e)));
    }
  }

  Future<void> _onRefreshWallet(
    RefreshWalletEvent event,
    Emitter<WalletState> emit,
  ) async {
    add(const LoadWalletsEvent());
  }
  
  String _getErrorMessage(dynamic error) {
    if (error is Exception) {
      final message = error.toString();
      if (message.contains('Exception: ')) {
        return message.replaceFirst('Exception: ', '');
      }
      return message;
    }
    return 'Une erreur est survenue';
  }
}