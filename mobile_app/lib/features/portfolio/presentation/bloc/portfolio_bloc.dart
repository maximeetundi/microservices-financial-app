import 'package:flutter_bloc/flutter_bloc.dart';
import 'package:equatable/equatable.dart';

import '../../../../core/services/api_service.dart';

// Events
abstract class PortfolioEvent extends Equatable {
  const PortfolioEvent();

  @override
  List<Object?> get props => [];
}

class LoadPortfolioEvent extends PortfolioEvent {
  const LoadPortfolioEvent();
}

class RefreshPortfolioEvent extends PortfolioEvent {
  const RefreshPortfolioEvent();
}

// States
abstract class PortfolioState extends Equatable {
  const PortfolioState();

  @override
  List<Object?> get props => [];
}

class PortfolioInitialState extends PortfolioState {
  const PortfolioInitialState();
}

class PortfolioLoadingState extends PortfolioState {
  const PortfolioLoadingState();
}

class PortfolioLoadedState extends PortfolioState {
  final double totalValue;
  final double dailyChange;
  final double dailyChangePercent;
  final List<PortfolioAsset> assets;
  final List<Map<String, dynamic>> recentTransactions;

  const PortfolioLoadedState({
    required this.totalValue,
    required this.dailyChange,
    required this.dailyChangePercent,
    required this.assets,
    required this.recentTransactions,
  });

  @override
  List<Object> get props => [
        totalValue,
        dailyChange,
        dailyChangePercent,
        assets,
        recentTransactions,
      ];
}

class PortfolioErrorState extends PortfolioState {
  final String message;

  const PortfolioErrorState({required this.message});

  @override
  List<Object> get props => [message];
}

// Asset model
class PortfolioAsset extends Equatable {
  final String currency;
  final String name;
  final double balance;
  final double value;
  final double allocation;
  final double change24h;
  final String iconUrl;

  const PortfolioAsset({
    required this.currency,
    required this.name,
    required this.balance,
    required this.value,
    required this.allocation,
    this.change24h = 0.0,
    this.iconUrl = '',
  });

  @override
  List<Object> get props => [
        currency,
        name,
        balance,
        value,
        allocation,
        change24h,
        iconUrl,
      ];
}

// BLoC
class PortfolioBloc extends Bloc<PortfolioEvent, PortfolioState> {
  final ApiService _apiService;

  PortfolioBloc({
    required ApiService apiService,
  })  : _apiService = apiService,
        super(const PortfolioInitialState()) {
    on<LoadPortfolioEvent>(_onLoadPortfolio);
    on<RefreshPortfolioEvent>(_onRefreshPortfolio);
  }

  Future<void> _onLoadPortfolio(
    LoadPortfolioEvent event,
    Emitter<PortfolioState> emit,
  ) async {
    emit(const PortfolioLoadingState());

    try {
      // Get wallets to calculate portfolio
      final walletsData = await _apiService.wallet.getWallets();
      
      double totalValue = 0;
      double previousValue = 0;
      final List<PortfolioAsset> assets = [];
      
      // Calculate totals and create assets
      for (final wallet in walletsData) {
        final balance = (wallet['balance'] as num?)?.toDouble() ?? 0.0;
        final usdRate = (wallet['usd_rate'] as num?)?.toDouble() ?? 1.0;
        final value = balance * usdRate;
        final change = (wallet['daily_change'] as num?)?.toDouble() ?? 0.0;
        
        totalValue += value;
        previousValue += value / (1 + change / 100);
        
        if (balance > 0) {
          assets.add(PortfolioAsset(
            currency: wallet['currency'] ?? '',
            name: _getCurrencyName(wallet['currency'] ?? ''),
            balance: balance,
            value: value,
            allocation: 0, // Will be calculated after total
            change24h: change,
          ));
        }
      }
      
      // Calculate allocations
      final assetsWithAllocation = assets.map((asset) {
        return PortfolioAsset(
          currency: asset.currency,
          name: asset.name,
          balance: asset.balance,
          value: asset.value,
          allocation: totalValue > 0 ? (asset.value / totalValue) * 100 : 0,
          change24h: asset.change24h,
        );
      }).toList();
      
      // Sort by value
      assetsWithAllocation.sort((a, b) => b.value.compareTo(a.value));
      
      final dailyChange = totalValue - previousValue;
      final dailyChangePercent = previousValue > 0 
          ? (dailyChange / previousValue) * 100 
          : 0.0;
      
      // Get recent transactions
      List<Map<String, dynamic>> recentTransactions = [];
      if (walletsData.isNotEmpty) {
        try {
          final firstWalletId = walletsData.first['id'] ?? '';
          if (firstWalletId.isNotEmpty) {
            recentTransactions = await _apiService.wallet.getTransactions(
              firstWalletId,
              limit: 5,
            );
          }
        } catch (_) {}
      }
      
      emit(PortfolioLoadedState(
        totalValue: totalValue,
        dailyChange: dailyChange,
        dailyChangePercent: dailyChangePercent,
        assets: assetsWithAllocation,
        recentTransactions: recentTransactions,
      ));
    } catch (e) {
      emit(PortfolioErrorState(message: _getErrorMessage(e)));
    }
  }

  Future<void> _onRefreshPortfolio(
    RefreshPortfolioEvent event,
    Emitter<PortfolioState> emit,
  ) async {
    add(const LoadPortfolioEvent());
  }
  
  String _getCurrencyName(String currency) {
    final names = {
      'BTC': 'Bitcoin',
      'ETH': 'Ethereum',
      'USD': 'US Dollar',
      'EUR': 'Euro',
      'GBP': 'British Pound',
      'USDT': 'Tether',
      'USDC': 'USD Coin',
      'BNB': 'Binance Coin',
      'XRP': 'Ripple',
      'SOL': 'Solana',
      'ADA': 'Cardano',
      'DOGE': 'Dogecoin',
    };
    return names[currency.toUpperCase()] ?? currency;
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
