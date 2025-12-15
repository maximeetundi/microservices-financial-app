/// Unified API Service - Point d'entrée unique pour tous les services API
/// 
/// Usage:
/// ```dart
/// final api = ApiService();
/// 
/// // Auth
/// await api.auth.login('email@example.com', 'password');
/// 
/// // Wallets
/// final wallets = await api.wallet.getWallets();
/// 
/// // Transfers
/// await api.transfer.createTransfer(...);
/// 
/// // Cards
/// final cards = await api.card.getCards();
/// 
/// // Exchange
/// final rate = await api.exchange.getExchangeRate('BTC', 'USD');
/// ```

export 'auth_api_service.dart';
export 'wallet_api_service.dart';
export 'transfer_api_service.dart';
export 'card_api_service.dart';
export 'exchange_api_service.dart';

import 'auth_api_service.dart';
import 'wallet_api_service.dart';
import 'transfer_api_service.dart';
import 'card_api_service.dart';
import 'exchange_api_service.dart';

/// Service API unifié regroupant tous les microservices
class ApiService {
  static final ApiService _instance = ApiService._internal();
  factory ApiService() => _instance;
  
  ApiService._internal();
  
  /// Service d'authentification
  final AuthApiService auth = AuthApiService();
  
  /// Service des portefeuilles
  final WalletApiService wallet = WalletApiService();
  
  /// Service des transferts
  final TransferApiService transfer = TransferApiService();
  
  /// Service des cartes
  final CardApiService card = CardApiService();
  
  /// Service d'échange
  final ExchangeApiService exchange = ExchangeApiService();
}
