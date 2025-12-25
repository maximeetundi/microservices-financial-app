import 'package:flutter_bloc/flutter_bloc.dart';
import 'package:equatable/equatable.dart';

import '../../../../core/services/api_service.dart';

// Events
abstract class ExchangeEvent extends Equatable {
  const ExchangeEvent();

  @override
  List<Object?> get props => [];
}

class LoadExchangeRatesEvent extends ExchangeEvent {
  const LoadExchangeRatesEvent();
}

class GetExchangeRateEvent extends ExchangeEvent {
  final String fromCurrency;
  final String toCurrency;

  const GetExchangeRateEvent({
    required this.fromCurrency,
    required this.toCurrency,
  });

  @override
  List<Object> get props => [fromCurrency, toCurrency];
}

class ExecuteExchangeEvent extends ExchangeEvent {
  final String fromWalletId;
  final String toWalletId;
  final double amount;
  final String fromCurrency;
  final String toCurrency;

  const ExecuteExchangeEvent({
    required this.fromWalletId,
    required this.toWalletId,
    required this.amount,
    required this.fromCurrency,
    required this.toCurrency,
  });

  @override
  List<Object> get props => [fromWalletId, toWalletId, amount, fromCurrency, toCurrency];
}

class LoadExchangeHistoryEvent extends ExchangeEvent {
  const LoadExchangeHistoryEvent();
}

// Trading Events
class LoadMarketsEvent extends ExchangeEvent {
  const LoadMarketsEvent();
}

class LoadTradingPortfolioEvent extends ExchangeEvent {
  const LoadTradingPortfolioEvent();
}

class LoadOrdersEvent extends ExchangeEvent {
  const LoadOrdersEvent();
}

class PlaceOrderEvent extends ExchangeEvent {
  final String symbol;
  final String side;
  final String type;
  final double amount;
  final double? price;
  final double? stopPrice;

  const PlaceOrderEvent({
    required this.symbol,
    required this.side,
    required this.type,
    required this.amount,
    this.price,
    this.stopPrice,
  });

  @override
  List<Object?> get props => [symbol, side, type, amount, price, stopPrice];
}

// States
abstract class ExchangeState extends Equatable {
  const ExchangeState();

  @override
  List<Object?> get props => [];
}

class ExchangeInitialState extends ExchangeState {
  const ExchangeInitialState();
}

class ExchangeLoadingState extends ExchangeState {
  const ExchangeLoadingState();
}

class ExchangeRatesLoadedState extends ExchangeState {
  final Map<String, dynamic> rates;

  const ExchangeRatesLoadedState({required this.rates});

  @override
  List<Object> get props => [rates];
}

class ExchangeRateLoadedState extends ExchangeState {
  final String fromCurrency;
  final String toCurrency;
  final double rate;
  final double fee;

  const ExchangeRateLoadedState({
    required this.fromCurrency,
    required this.toCurrency,
    required this.rate,
    this.fee = 0.0,
  });

  @override
  List<Object> get props => [fromCurrency, toCurrency, rate, fee];
}

class ExchangeSuccessState extends ExchangeState {
  final String exchangeId;
  final double fromAmount;
  final double toAmount;

  const ExchangeSuccessState({
    required this.exchangeId,
    required this.fromAmount,
    required this.toAmount,
  });

  @override
  List<Object> get props => [exchangeId, fromAmount, toAmount];
}

class ExchangeHistoryLoadedState extends ExchangeState {
  final List<Map<String, dynamic>> exchanges;

  const ExchangeHistoryLoadedState({required this.exchanges});

  @override
  List<Object> get props => [exchanges];
}

class ExchangeErrorState extends ExchangeState {
  final String message;

  const ExchangeErrorState({required this.message});

  @override
  List<Object> get props => [message];
}

// Trading States
class MarketsLoadedState extends ExchangeState {
  final Map<String, dynamic> markets; // Adjust type if needed based on API response

  const MarketsLoadedState({required this.markets});

  @override
  List<Object> get props => [markets];
}

class TradingPortfolioLoadedState extends ExchangeState {
  final Map<String, dynamic> portfolio;

  const TradingPortfolioLoadedState({required this.portfolio});

  @override
  List<Object> get props => [portfolio];
}

class OrdersLoadedState extends ExchangeState {
  final List<Map<String, dynamic>> orders;

  const OrdersLoadedState({required this.orders});

  @override
  List<Object> get props => [orders];
}

class OrderPlacedState extends ExchangeState {
  final Map<String, dynamic> order;

  const OrderPlacedState({required this.order});

  @override
  List<Object> get props => [order];
}

// BLoC
class ExchangeBloc extends Bloc<ExchangeEvent, ExchangeState> {
  final ApiService _apiService;

  ExchangeBloc({
    required ApiService apiService,
  })  : _apiService = apiService,
        super(const ExchangeInitialState()) {
    on<LoadExchangeRatesEvent>(_onLoadRates);
    on<GetExchangeRateEvent>(_onGetRate);
    on<ExecuteExchangeEvent>(_onExecuteExchange);
    on<LoadExchangeHistoryEvent>(_onLoadHistory);
    // Trading handlers
    on<LoadMarketsEvent>(_onLoadMarkets);
    on<LoadTradingPortfolioEvent>(_onLoadTradingPortfolio);
    on<LoadOrdersEvent>(_onLoadOrders);
    on<PlaceOrderEvent>(_onPlaceOrder);
  }

  Future<void> _onLoadRates(
    LoadExchangeRatesEvent event,
    Emitter<ExchangeState> emit,
  ) async {
    emit(const ExchangeLoadingState());

    try {
      final rates = await _apiService.exchange.getRates();
      emit(ExchangeRatesLoadedState(rates: rates));
    } catch (e) {
      emit(ExchangeErrorState(message: _getErrorMessage(e)));
    }
  }

  Future<void> _onGetRate(
    GetExchangeRateEvent event,
    Emitter<ExchangeState> emit,
  ) async {
    emit(const ExchangeLoadingState());

    try {
      final rateData = await _apiService.exchange.getExchangeRate(
        event.fromCurrency,
        event.toCurrency,
      );
      
      emit(ExchangeRateLoadedState(
        fromCurrency: event.fromCurrency,
        toCurrency: event.toCurrency,
        rate: (rateData['rate'] as num).toDouble(),
        fee: (rateData['fee'] as num?)?.toDouble() ?? 0.0,
      ));
    } catch (e) {
      emit(ExchangeErrorState(message: _getErrorMessage(e)));
    }
  }

  Future<void> _onExecuteExchange(
    ExecuteExchangeEvent event,
    Emitter<ExchangeState> emit,
  ) async {
    emit(const ExchangeLoadingState());

    try {
      final result = await _apiService.exchange.executeExchange(
        fromWalletId: event.fromWalletId,
        toWalletId: event.toWalletId,
        amount: event.amount,
        fromCurrency: event.fromCurrency,
        toCurrency: event.toCurrency,
      );
      
      emit(ExchangeSuccessState(
        exchangeId: result['exchange_id'] ?? '',
        fromAmount: (result['from_amount'] as num?)?.toDouble() ?? event.amount,
        toAmount: (result['to_amount'] as num?)?.toDouble() ?? 0.0,
      ));
    } catch (e) {
      emit(ExchangeErrorState(message: _getErrorMessage(e)));
    }
  }

  Future<void> _onLoadHistory(
    LoadExchangeHistoryEvent event,
    Emitter<ExchangeState> emit,
  ) async {
    emit(const ExchangeLoadingState());

    try {
      final exchanges = await _apiService.exchange.getExchangeHistory();
      emit(ExchangeHistoryLoadedState(exchanges: exchanges));
    } catch (e) {
      emit(ExchangeErrorState(message: _getErrorMessage(e)));
    }
  }

  // Trading Handlers
  Future<void> _onLoadMarkets(
    LoadMarketsEvent event,
    Emitter<ExchangeState> emit,
  ) async {
    emit(const ExchangeLoadingState());
    try {
      final markets = await _apiService.exchange.getMarkets();
      emit(MarketsLoadedState(markets: markets));
    } catch (e) {
      emit(ExchangeErrorState(message: _getErrorMessage(e)));
    }
  }

  Future<void> _onLoadTradingPortfolio(
    LoadTradingPortfolioEvent event,
    Emitter<ExchangeState> emit,
  ) async {
    emit(const ExchangeLoadingState());
    try {
      final portfolio = await _apiService.exchange.getTradingPortfolio();
      emit(TradingPortfolioLoadedState(portfolio: portfolio));
    } catch (e) {
      emit(ExchangeErrorState(message: _getErrorMessage(e)));
    }
  }

  Future<void> _onLoadOrders(
    LoadOrdersEvent event,
    Emitter<ExchangeState> emit,
  ) async {
    emit(const ExchangeLoadingState());
    try {
      final orders = await _apiService.exchange.getOrders();
      emit(OrdersLoadedState(orders: orders));
    } catch (e) {
      emit(ExchangeErrorState(message: _getErrorMessage(e)));
    }
  }

  Future<void> _onPlaceOrder(
    PlaceOrderEvent event,
    Emitter<ExchangeState> emit,
  ) async {
    emit(const ExchangeLoadingState());
    try {
      final result = await _apiService.exchange.placeOrder(
        symbol: event.symbol,
        side: event.side,
        type: event.type,
        amount: event.amount,
        price: event.price,
        stopPrice: event.stopPrice,
      );
      emit(OrderPlacedState(order: result));
    } catch (e) {
      emit(ExchangeErrorState(message: _getErrorMessage(e)));
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
