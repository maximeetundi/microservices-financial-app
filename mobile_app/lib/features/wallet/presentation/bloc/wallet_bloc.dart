import 'package:flutter_bloc/flutter_bloc.dart';
import 'package:equatable/equatable.dart';

import '../../domain/entities/wallet.dart';
import '../../domain/entities/transaction.dart';
import '../../domain/usecases/get_wallets_usecase.dart';
import '../../domain/usecases/create_wallet_usecase.dart';
import '../../domain/usecases/send_transaction_usecase.dart';
import '../../domain/usecases/get_wallet_transactions_usecase.dart';

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
  final String? name;

  const CreateWalletEvent({
    required this.currency,
    this.name,
  });

  @override
  List<Object?> get props => [currency, name];
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

class LoadWalletTransactionsEvent extends WalletEvent {
  final String walletId;
  final int page;
  final int limit;

  const LoadWalletTransactionsEvent({
    required this.walletId,
    this.page = 1,
    this.limit = 20,
  });

  @override
  List<Object> get props => [walletId, page, limit];
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
  final List<Transaction> recentTransactions;

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
    List<Transaction>? recentTransactions,
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
  final Transaction transaction;

  const WalletTransactionSentState({required this.transaction});

  @override
  List<Object> get props => [transaction];
}

class WalletCreatedState extends WalletState {
  final Wallet wallet;

  const WalletCreatedState({required this.wallet});

  @override
  List<Object> get props => [wallet];
}

// BLoC
class WalletBloc extends Bloc<WalletEvent, WalletState> {
  final GetWalletsUseCase _getWalletsUseCase;
  final CreateWalletUseCase _createWalletUseCase;
  final SendTransactionUseCase _sendTransactionUseCase;
  final GetWalletTransactionsUseCase _getWalletTransactionsUseCase;

  WalletBloc({
    required GetWalletsUseCase getWalletsUseCase,
    required CreateWalletUseCase createWalletUseCase,
    required SendTransactionUseCase sendTransactionUseCase,
    required GetWalletTransactionsUseCase getWalletTransactionsUseCase,
  })  : _getWalletsUseCase = getWalletsUseCase,
        _createWalletUseCase = createWalletUseCase,
        _sendTransactionUseCase = sendTransactionUseCase,
        _getWalletTransactionsUseCase = getWalletTransactionsUseCase,
        super(const WalletInitialState()) {
    on<LoadWalletsEvent>(_onLoadWallets);
    on<CreateWalletEvent>(_onCreateWallet);
    on<SendTransactionEvent>(_onSendTransaction);
    on<LoadWalletTransactionsEvent>(_onLoadWalletTransactions);
    on<RefreshWalletEvent>(_onRefreshWallet);
  }

  Future<void> _onLoadWallets(
    LoadWalletsEvent event,
    Emitter<WalletState> emit,
  ) async {
    emit(const WalletLoadingState());

    try {
      final result = await _getWalletsUseCase(NoParams());

      result.fold(
        (failure) => emit(WalletErrorState(message: failure.message)),
        (walletsData) {
          final totalValue = walletsData.wallets
              .fold<double>(0, (sum, wallet) => sum + wallet.balance * wallet.usdRate);
          
          // Calculate daily change (mock calculation)
          final dailyChange = 2.3; // This should come from the API
          
          emit(WalletLoadedState(
            wallets: walletsData.wallets,
            totalValue: totalValue,
            dailyChange: dailyChange,
            recentTransactions: walletsData.recentTransactions,
          ));
        },
      );
    } catch (e) {
      emit(WalletErrorState(message: e.toString()));
    }
  }

  Future<void> _onCreateWallet(
    CreateWalletEvent event,
    Emitter<WalletState> emit,
  ) async {
    try {
      final result = await _createWalletUseCase(CreateWalletParams(
        currency: event.currency,
        name: event.name,
      ));

      result.fold(
        (failure) => emit(WalletErrorState(message: failure.message)),
        (wallet) {
          emit(WalletCreatedState(wallet: wallet));
          // Reload wallets to show the new one
          add(const LoadWalletsEvent());
        },
      );
    } catch (e) {
      emit(WalletErrorState(message: e.toString()));
    }
  }

  Future<void> _onSendTransaction(
    SendTransactionEvent event,
    Emitter<WalletState> emit,
  ) async {
    try {
      final result = await _sendTransactionUseCase(SendTransactionParams(
        walletId: event.walletId,
        toAddress: event.toAddress,
        amount: event.amount,
        memo: event.memo,
      ));

      result.fold(
        (failure) => emit(WalletErrorState(message: failure.message)),
        (transaction) {
          emit(WalletTransactionSentState(transaction: transaction));
          // Reload wallets to update balances
          add(const LoadWalletsEvent());
        },
      );
    } catch (e) {
      emit(WalletErrorState(message: e.toString()));
    }
  }

  Future<void> _onLoadWalletTransactions(
    LoadWalletTransactionsEvent event,
    Emitter<WalletState> emit,
  ) async {
    try {
      final result = await _getWalletTransactionsUseCase(
        GetWalletTransactionsParams(
          walletId: event.walletId,
          page: event.page,
          limit: event.limit,
        ),
      );

      result.fold(
        (failure) => emit(WalletErrorState(message: failure.message)),
        (transactions) {
          // Update the current state with new transactions
          if (state is WalletLoadedState) {
            emit((state as WalletLoadedState).copyWith(
              recentTransactions: transactions,
            ));
          }
        },
      );
    } catch (e) {
      emit(WalletErrorState(message: e.toString()));
    }
  }

  Future<void> _onRefreshWallet(
    RefreshWalletEvent event,
    Emitter<WalletState> emit,
  ) async {
    // Just reload all wallets for now
    add(const LoadWalletsEvent());
  }
}

// Use Case Parameters
class CreateWalletParams extends Equatable {
  final String currency;
  final String? name;

  const CreateWalletParams({
    required this.currency,
    this.name,
  });

  @override
  List<Object?> get props => [currency, name];
}

class SendTransactionParams extends Equatable {
  final String walletId;
  final String toAddress;
  final double amount;
  final String? memo;

  const SendTransactionParams({
    required this.walletId,
    required this.toAddress,
    required this.amount,
    this.memo,
  });

  @override
  List<Object?> get props => [walletId, toAddress, amount, memo];
}

class GetWalletTransactionsParams extends Equatable {
  final String walletId;
  final int page;
  final int limit;

  const GetWalletTransactionsParams({
    required this.walletId,
    this.page = 1,
    this.limit = 20,
  });

  @override
  List<Object> get props => [walletId, page, limit];
}

class NoParams extends Equatable {
  @override
  List<Object> get props => [];
}