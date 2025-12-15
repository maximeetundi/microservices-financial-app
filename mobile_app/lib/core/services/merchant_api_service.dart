import 'package:http/http.dart' as http;
import 'dart:convert';
import 'base_api_service.dart';

/// Service API pour les paiements marchands
class MerchantApiService extends BaseApiService {
  /// Créer une demande de paiement
  Future<Map<String, dynamic>> createPayment(Map<String, dynamic> data) async {
    final response = await authenticatedPost('/merchant/payments', data);
    return response;
  }

  /// Obtenir toutes les demandes de paiement
  Future<Map<String, dynamic>> getPayments({int limit = 20, int offset = 0}) async {
    final response = await authenticatedGet('/merchant/payments?limit=$limit&offset=$offset');
    return response;
  }

  /// Obtenir l'historique des paiements
  Future<Map<String, dynamic>> getHistory({int limit = 20, int offset = 0}) async {
    final response = await authenticatedGet('/merchant/payments/history?limit=$limit&offset=$offset');
    return response;
  }

  /// Annuler une demande de paiement
  Future<void> cancelPayment(String paymentId) async {
    await authenticatedDelete('/merchant/payments/$paymentId');
  }

  /// Obtenir le QR code d'un paiement
  Future<Map<String, dynamic>> getQRCode(String paymentId) async {
    final response = await authenticatedGet('/payments/$paymentId/qr');
    return response;
  }

  /// Paiement rapide
  Future<Map<String, dynamic>> quickPay(Map<String, dynamic> data) async {
    final response = await authenticatedPost('/merchant/quick-pay', data);
    return response;
  }

  /// Obtenir les détails d'un paiement (public - pour scan)
  Future<Map<String, dynamic>> getPaymentDetails(String paymentId) async {
    final response = await get('/pay/$paymentId');
    return response;
  }

  /// Payer une demande de paiement
  Future<Map<String, dynamic>> payPayment(String paymentId, String walletId, double amount) async {
    final response = await authenticatedPost('/payments/$paymentId/pay', {
      'from_wallet_id': walletId,
      'amount': amount,
    });
    return response;
  }
}
