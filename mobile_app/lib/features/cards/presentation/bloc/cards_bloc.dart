import 'package:flutter_bloc/flutter_bloc.dart';
import 'package:equatable/equatable.dart';

import '../../../../core/services/api_service.dart';

// Events
abstract class CardsEvent extends Equatable {
  const CardsEvent();

  @override
  List<Object?> get props => [];
}

class LoadCardsEvent extends CardsEvent {
  const LoadCardsEvent();
}

class CreateVirtualCardEvent extends CardsEvent {
  final String currency;
  final double? initialAmount;
  final String? sourceWalletId;
  final String? name;

  const CreateVirtualCardEvent({
    required this.currency,
    this.initialAmount,
    this.sourceWalletId,
    this.name,
  });

  @override
  List<Object?> get props => [currency, initialAmount, sourceWalletId, name];
}

class FreezeCardEvent extends CardsEvent {
  final String cardId;
  final String reason;

  const FreezeCardEvent({required this.cardId, this.reason = 'user_requested'});

  @override
  List<Object> get props => [cardId, reason];
}

class UnfreezeCardEvent extends CardsEvent {
  final String cardId;

  const UnfreezeCardEvent({required this.cardId});

  @override
  List<Object> get props => [cardId];
}

class LoadCardEvent extends CardsEvent {
  final String cardId;
  final double amount;
  final String sourceWalletId;

  const LoadCardEvent({
    required this.cardId,
    required this.amount,
    required this.sourceWalletId,
  });

  @override
  List<Object> get props => [cardId, amount, sourceWalletId];
}

class SetCardPINEvent extends CardsEvent {
  final String cardId;
  final String pin;

  const SetCardPINEvent({required this.cardId, required this.pin});

  @override
  List<Object> get props => [cardId, pin];
}

class LoadCardTransactionsEvent extends CardsEvent {
  final String cardId;

  const LoadCardTransactionsEvent({required this.cardId});

  @override
  List<Object> get props => [cardId];
}

class TopUpCardEvent extends CardsEvent {
  final String cardId;
  final double amount;
  final String? sourceWalletId;

  const TopUpCardEvent({
    required this.cardId,
    required this.amount,
    this.sourceWalletId,
  });

  @override
  List<Object?> get props => [cardId, amount, sourceWalletId];
}

// States
abstract class CardsState extends Equatable {
  const CardsState();

  @override
  List<Object?> get props => [];
}

class CardsInitialState extends CardsState {
  const CardsInitialState();
}

class CardsLoadingState extends CardsState {
  const CardsLoadingState();
}

class CardsLoadedState extends CardsState {
  final List<Map<String, dynamic>> cards;

  const CardsLoadedState({required this.cards});

  @override
  List<Object> get props => [cards];
}

class CardCreatedState extends CardsState {
  final Map<String, dynamic> card;

  const CardCreatedState({required this.card});

  @override
  List<Object> get props => [card];
}

class CardActionSuccessState extends CardsState {
  final String message;

  const CardActionSuccessState({required this.message});

  @override
  List<Object> get props => [message];
}

class CardTransactionsLoadedState extends CardsState {
  final List<Map<String, dynamic>> transactions;

  const CardTransactionsLoadedState({required this.transactions});

  @override
  List<Object> get props => [transactions];
}

class CardsErrorState extends CardsState {
  final String message;

  const CardsErrorState({required this.message});

  @override
  List<Object> get props => [message];
}

// BLoC
class CardsBloc extends Bloc<CardsEvent, CardsState> {
  final ApiService _apiService;

  CardsBloc({
    required ApiService apiService,
  })  : _apiService = apiService,
        super(const CardsInitialState()) {
    on<LoadCardsEvent>(_onLoadCards);
    on<CreateVirtualCardEvent>(_onCreateVirtualCard);
    on<FreezeCardEvent>(_onFreezeCard);
    on<UnfreezeCardEvent>(_onUnfreezeCard);
    on<LoadCardEvent>(_onLoadCard);
    on<SetCardPINEvent>(_onSetCardPIN);
    on<LoadCardTransactionsEvent>(_onLoadCardTransactions);
  }

  Future<void> _onLoadCards(
    LoadCardsEvent event,
    Emitter<CardsState> emit,
  ) async {
    emit(const CardsLoadingState());

    try {
      final cards = await _apiService.card.getCards();
      emit(CardsLoadedState(cards: cards));
    } catch (e) {
      emit(CardsErrorState(message: _getErrorMessage(e)));
    }
  }

  Future<void> _onCreateVirtualCard(
    CreateVirtualCardEvent event,
    Emitter<CardsState> emit,
  ) async {
    emit(const CardsLoadingState());

    try {
      final card = await _apiService.card.createVirtualCard(
        currency: event.currency,
        initialAmount: event.initialAmount,
        sourceWalletId: event.sourceWalletId,
        name: event.name,
      );
      
      emit(CardCreatedState(card: card));
      add(const LoadCardsEvent());
    } catch (e) {
      emit(CardsErrorState(message: _getErrorMessage(e)));
    }
  }

  Future<void> _onFreezeCard(
    FreezeCardEvent event,
    Emitter<CardsState> emit,
  ) async {
    try {
      await _apiService.card.freezeCard(event.cardId, event.reason);
      emit(const CardActionSuccessState(message: 'Carte gelée avec succès'));
      add(const LoadCardsEvent());
    } catch (e) {
      emit(CardsErrorState(message: _getErrorMessage(e)));
    }
  }

  Future<void> _onUnfreezeCard(
    UnfreezeCardEvent event,
    Emitter<CardsState> emit,
  ) async {
    try {
      await _apiService.card.unfreezeCard(event.cardId);
      emit(const CardActionSuccessState(message: 'Carte dégelée avec succès'));
      add(const LoadCardsEvent());
    } catch (e) {
      emit(CardsErrorState(message: _getErrorMessage(e)));
    }
  }

  Future<void> _onLoadCard(
    LoadCardEvent event,
    Emitter<CardsState> emit,
  ) async {
    try {
      await _apiService.card.loadCard(
        cardId: event.cardId,
        amount: event.amount,
        sourceWalletId: event.sourceWalletId,
      );
      emit(const CardActionSuccessState(message: 'Carte rechargée avec succès'));
      add(const LoadCardsEvent());
    } catch (e) {
      emit(CardsErrorState(message: _getErrorMessage(e)));
    }
  }

  Future<void> _onSetCardPIN(
    SetCardPINEvent event,
    Emitter<CardsState> emit,
  ) async {
    try {
      await _apiService.card.setCardPIN(event.cardId, event.pin);
      emit(const CardActionSuccessState(message: 'PIN défini avec succès'));
    } catch (e) {
      emit(CardsErrorState(message: _getErrorMessage(e)));
    }
  }

  Future<void> _onLoadCardTransactions(
    LoadCardTransactionsEvent event,
    Emitter<CardsState> emit,
  ) async {
    emit(const CardsLoadingState());

    try {
      final transactions = await _apiService.card.getCardTransactions(event.cardId);
      emit(CardTransactionsLoadedState(transactions: transactions));
    } catch (e) {
      emit(CardsErrorState(message: _getErrorMessage(e)));
    }
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
