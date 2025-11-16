import 'package:flutter_bloc/flutter_bloc.dart';
import 'package:equatable/equatable.dart';

import '../../domain/entities/transfer.dart';
import '../../domain/usecases/send_transfer_usecase.dart';
import '../../domain/usecases/get_transfer_history_usecase.dart';
import '../../../wallet/domain/entities/wallet.dart';

// Events
abstract class TransferEvent extends Equatable {
  const TransferEvent();

  @override
  List<Object?> get props => [];
}

class LoadTransferDataEvent extends TransferEvent {
  const LoadTransferDataEvent();
}

class SendTransferEvent extends TransferEvent {
  final String type;
  final String fromWallet;
  final String recipient;
  final double amount;
  final String? memo;
  final String? bankCode;
  final String? recipientName;

  const SendTransferEvent({
    required this.type,
    required this.fromWallet,
    required this.recipient,
    required this.amount,
    this.memo,
    this.bankCode,
    this.recipientName,
  });

  @override
  List<Object?> get props => [
        type,
        fromWallet,
        recipient,
        amount,
        memo,
        bankCode,
        recipientName,
      ];
}

class GetTransferHistoryEvent extends TransferEvent {
  final int page;
  final int limit;
  final String? status;
  final String? type;

  const GetTransferHistoryEvent({
    this.page = 1,
    this.limit = 20,
    this.status,
    this.type,
  });

  @override
  List<Object?> get props => [page, limit, status, type];
}

class EstimateTransferFeeEvent extends TransferEvent {
  final String type;
  final String fromCurrency;
  final double amount;

  const EstimateTransferFeeEvent({
    required this.type,
    required this.fromCurrency,
    required this.amount,
  });

  @override
  List<Object> get props => [type, fromCurrency, amount];
}

// States
abstract class TransferState extends Equatable {
  const TransferState();

  @override
  List<Object?> get props => [];
}

class TransferInitialState extends TransferState {
  const TransferInitialState();
}

class TransferLoadingState extends TransferState {
  const TransferLoadingState();
}

class TransferLoadedState extends TransferState {
  final List<Wallet> wallets;
  final List<Transfer> recentTransfers;
  final List<Contact> contacts;
  final TransferFee? estimatedFee;

  const TransferLoadedState({
    required this.wallets,
    required this.recentTransfers,
    required this.contacts,
    this.estimatedFee,
  });

  @override
  List<Object?> get props => [wallets, recentTransfers, contacts, estimatedFee];

  TransferLoadedState copyWith({
    List<Wallet>? wallets,
    List<Transfer>? recentTransfers,
    List<Contact>? contacts,
    TransferFee? estimatedFee,
  }) {
    return TransferLoadedState(
      wallets: wallets ?? this.wallets,
      recentTransfers: recentTransfers ?? this.recentTransfers,
      contacts: contacts ?? this.contacts,
      estimatedFee: estimatedFee ?? this.estimatedFee,
    );
  }
}

class TransferSuccessState extends TransferState {
  final Transfer transfer;

  const TransferSuccessState({required this.transfer});

  @override
  List<Object> get props => [transfer];
}

class TransferErrorState extends TransferState {
  final String message;

  const TransferErrorState({required this.message});

  @override
  List<Object> get props => [message];
}

class TransferFeeEstimatedState extends TransferState {
  final TransferFee fee;

  const TransferFeeEstimatedState({required this.fee});

  @override
  List<Object> get props => [fee];
}

// BLoC
class TransferBloc extends Bloc<TransferEvent, TransferState> {
  final SendTransferUseCase _sendTransferUseCase;
  final GetTransferHistoryUseCase _getTransferHistoryUseCase;

  TransferBloc({
    required SendTransferUseCase sendTransferUseCase,
    required GetTransferHistoryUseCase getTransferHistoryUseCase,
  })  : _sendTransferUseCase = sendTransferUseCase,
        _getTransferHistoryUseCase = getTransferHistoryUseCase,
        super(const TransferInitialState()) {
    on<LoadTransferDataEvent>(_onLoadTransferData);
    on<SendTransferEvent>(_onSendTransfer);
    on<GetTransferHistoryEvent>(_onGetTransferHistory);
    on<EstimateTransferFeeEvent>(_onEstimateTransferFee);
  }

  Future<void> _onLoadTransferData(
    LoadTransferDataEvent event,
    Emitter<TransferState> emit,
  ) async {
    emit(const TransferLoadingState());

    try {
      // Load wallets, recent transfers, and contacts
      // This would normally call multiple use cases
      
      // Mock data for now
      final wallets = [
        Wallet(
          id: '1',
          userId: 'user1',
          currency: 'BTC',
          balance: 0.5,
          availableBalance: 0.5,
          address: '1A1zP1eP5QGefi2DMPTfTL5SLmv7DivfNa',
          type: WalletType.crypto,
          createdAt: DateTime.now(),
          updatedAt: DateTime.now(),
          networkInfo: const WalletNetworkInfo(network: 'Bitcoin'),
        ),
        Wallet(
          id: '2',
          userId: 'user1',
          currency: 'USD',
          balance: 5000.0,
          availableBalance: 5000.0,
          address: 'USD-WALLET',
          type: WalletType.fiat,
          createdAt: DateTime.now(),
          updatedAt: DateTime.now(),
          networkInfo: const WalletNetworkInfo(network: 'Bank'),
        ),
      ];

      final recentTransfers = [
        Transfer(
          id: '1',
          userId: 'user1',
          type: TransferType.crypto,
          fromWallet: '1',
          toAddress: '3FUpjxWpPGqxGSzeLdZHamksAPtJ3EGcjh',
          amount: 0.1,
          currency: 'BTC',
          status: TransferStatus.completed,
          createdAt: DateTime.now().subtract(const Duration(hours: 2)),
        ),
        Transfer(
          id: '2',
          userId: 'user1',
          type: TransferType.fiat,
          fromWallet: '2',
          toAddress: 'US1234567890123456',
          amount: 500.0,
          currency: 'USD',
          status: TransferStatus.pending,
          createdAt: DateTime.now().subtract(const Duration(days: 1)),
          recipientName: 'John Doe',
          bankCode: 'CHASUS33',
        ),
      ];

      final contacts = [
        Contact(
          id: '1',
          name: 'Alice Johnson',
          email: 'alice@example.com',
          address: '1BvBMSEYstWetqTFn5Au4m4GFg7xJaNVN2',
          currency: 'BTC',
          type: ContactType.crypto,
        ),
        Contact(
          id: '2',
          name: 'Bob Smith',
          email: 'bob@example.com',
          address: 'bob@cryptobank.com',
          currency: 'USD',
          type: ContactType.instant,
        ),
      ];

      emit(TransferLoadedState(
        wallets: wallets,
        recentTransfers: recentTransfers,
        contacts: contacts,
      ));
    } catch (e) {
      emit(TransferErrorState(message: e.toString()));
    }
  }

  Future<void> _onSendTransfer(
    SendTransferEvent event,
    Emitter<TransferState> emit,
  ) async {
    emit(const TransferLoadingState());

    try {
      final result = await _sendTransferUseCase(SendTransferParams(
        type: event.type,
        fromWallet: event.fromWallet,
        recipient: event.recipient,
        amount: event.amount,
        memo: event.memo,
        bankCode: event.bankCode,
        recipientName: event.recipientName,
      ));

      result.fold(
        (failure) => emit(TransferErrorState(message: failure.message)),
        (transfer) {
          emit(TransferSuccessState(transfer: transfer));
          // Reload transfer data to show the new transfer
          add(const LoadTransferDataEvent());
        },
      );
    } catch (e) {
      emit(TransferErrorState(message: e.toString()));
    }
  }

  Future<void> _onGetTransferHistory(
    GetTransferHistoryEvent event,
    Emitter<TransferState> emit,
  ) async {
    try {
      final result = await _getTransferHistoryUseCase(
        GetTransferHistoryParams(
          page: event.page,
          limit: event.limit,
          status: event.status,
          type: event.type,
        ),
      );

      result.fold(
        (failure) => emit(TransferErrorState(message: failure.message)),
        (transfers) {
          if (state is TransferLoadedState) {
            emit((state as TransferLoadedState).copyWith(
              recentTransfers: transfers,
            ));
          }
        },
      );
    } catch (e) {
      emit(TransferErrorState(message: e.toString()));
    }
  }

  Future<void> _onEstimateTransferFee(
    EstimateTransferFeeEvent event,
    Emitter<TransferState> emit,
  ) async {
    try {
      // Calculate fee based on type
      final fee = TransferFee.calculate(
        type: event.type,
        amount: event.amount,
        currency: event.fromCurrency,
      );

      emit(TransferFeeEstimatedState(fee: fee));

      if (state is TransferLoadedState) {
        emit((state as TransferLoadedState).copyWith(estimatedFee: fee));
      }
    } catch (e) {
      emit(TransferErrorState(message: e.toString()));
    }
  }
}

// Use Case Parameters
class SendTransferParams extends Equatable {
  final String type;
  final String fromWallet;
  final String recipient;
  final double amount;
  final String? memo;
  final String? bankCode;
  final String? recipientName;

  const SendTransferParams({
    required this.type,
    required this.fromWallet,
    required this.recipient,
    required this.amount,
    this.memo,
    this.bankCode,
    this.recipientName,
  });

  @override
  List<Object?> get props => [
        type,
        fromWallet,
        recipient,
        amount,
        memo,
        bankCode,
        recipientName,
      ];
}

class GetTransferHistoryParams extends Equatable {
  final int page;
  final int limit;
  final String? status;
  final String? type;

  const GetTransferHistoryParams({
    this.page = 1,
    this.limit = 20,
    this.status,
    this.type,
  });

  @override
  List<Object?> get props => [page, limit, status, type];
}

// Models
class Contact extends Equatable {
  final String id;
  final String name;
  final String email;
  final String address;
  final String currency;
  final ContactType type;
  final DateTime? lastUsed;

  const Contact({
    required this.id,
    required this.name,
    required this.email,
    required this.address,
    required this.currency,
    required this.type,
    this.lastUsed,
  });

  @override
  List<Object?> get props => [id, name, email, address, currency, type, lastUsed];
}

enum ContactType { crypto, fiat, instant }

class TransferFee extends Equatable {
  final double amount;
  final String currency;
  final String type;
  final double percentage;

  const TransferFee({
    required this.amount,
    required this.currency,
    required this.type,
    required this.percentage,
  });

  static TransferFee calculate({
    required String type,
    required double amount,
    required String currency,
  }) {
    double feePercentage;
    
    switch (type) {
      case 'crypto':
        feePercentage = 0.0025; // 0.25%
        break;
      case 'fiat':
        feePercentage = 0.001; // 0.1% minimum $2.50
        break;
      case 'instant':
        feePercentage = 0.0; // Free
        break;
      default:
        feePercentage = 0.005; // 0.5%
    }

    double feeAmount = amount * feePercentage;
    
    // Apply minimum fees for fiat
    if (type == 'fiat' && feeAmount < 2.50) {
      feeAmount = 2.50;
    }

    return TransferFee(
      amount: feeAmount,
      currency: currency,
      type: type,
      percentage: feePercentage,
    );
  }

  @override
  List<Object> get props => [amount, currency, type, percentage];
}