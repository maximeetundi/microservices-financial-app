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
/// 
/// // Merchant
/// await api.merchant.createPayment(...);
/// ```

export 'auth_api_service.dart';
export 'wallet_api_service.dart';
export 'transfer_api_service.dart';
export 'card_api_service.dart';
export 'exchange_api_service.dart';
export 'merchant_api_service.dart';
export 'enterprise_api_service.dart';
export 'shop_api_service.dart';

import 'auth_api_service.dart';
import 'wallet_api_service.dart';
import 'transfer_api_service.dart';
import 'card_api_service.dart';
import 'exchange_api_service.dart';
import 'merchant_api_service.dart';
import 'enterprise_api_service.dart';
import 'shop_api_service.dart';

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
  
  /// Service des paiements marchands
  final MerchantApiService merchant = MerchantApiService();

  /// Service entreprise
  final EnterpriseApiService enterprise = EnterpriseApiService();

  /// Service boutique (Shop)
  final ShopApiService shop = ShopApiService();
  
  // ========== Static methods for easy access ==========
  
  static Future<List<dynamic>> getWallets() async {
    return await _instance.wallet.getWallets();
  }
  
  static Future<Map<String, dynamic>> getMerchantPayments() async {
    return await _instance.merchant.getPayments();
  }
  
  static Future<Map<String, dynamic>> createMerchantPayment(Map<String, dynamic> data) async {
    return await _instance.merchant.createPayment(data);
  }
  
  static Future<Map<String, dynamic>> getPaymentDetails(String id) async {
    return await _instance.merchant.getPaymentDetails(id);
  }
  
  static Future<Map<String, dynamic>> payPayment(String id, String walletId, double amount) async {
    return await _instance.merchant.payPayment(id, walletId, amount);
  }
}

