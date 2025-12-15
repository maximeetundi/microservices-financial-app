import '../api/api_client.dart';
import '../api/api_endpoints.dart';

/// Service des cartes
class CardApiService {
  final ApiClient _client = ApiClient();
  
  /// Récupérer toutes les cartes de l'utilisateur
  Future<List<Map<String, dynamic>>> getCards() async {
    final response = await _client.get(ApiEndpoints.cardsList);
    if (response.statusCode == 200) {
      return List<Map<String, dynamic>>.from(response.data['cards']);
    }
    throw Exception('Failed to get cards');
  }
  
  /// Créer une nouvelle carte
  Future<Map<String, dynamic>> createCard({
    required String cardType,
    required String currency,
    String? name,
  }) async {
    final response = await _client.post(
      ApiEndpoints.createCard,
      data: {
        'card_type': cardType,
        'currency': currency,
        if (name != null) 'name': name,
      },
    );
    if (response.statusCode == 201) {
      return response.data['card'];
    }
    throw Exception(response.data['error'] ?? 'Failed to create card');
  }
  
  /// Récupérer une carte par ID
  Future<Map<String, dynamic>> getCard(String cardId) async {
    final response = await _client.get(ApiEndpoints.cardById(cardId));
    if (response.statusCode == 200) {
      return response.data['card'];
    }
    throw Exception('Card not found');
  }
  
  /// Créer une carte virtuelle
  Future<Map<String, dynamic>> createVirtualCard({
    required String currency,
    double? initialAmount,
    String? sourceWalletId,
    String? name,
  }) async {
    final response = await _client.post(
      ApiEndpoints.virtualCard,
      data: {
        'currency': currency,
        if (initialAmount != null) 'initial_amount': initialAmount,
        if (sourceWalletId != null) 'source_wallet_id': sourceWalletId,
        if (name != null) 'name': name,
      },
    );
    if (response.statusCode == 201) {
      return response.data['card'];
    }
    throw Exception(response.data['error'] ?? 'Failed to create virtual card');
  }
  
  /// Commander une carte physique
  Future<Map<String, dynamic>> orderPhysicalCard({
    required String currency,
    required String cardholderName,
    required Map<String, dynamic> shippingAddress,
    bool expressShipping = false,
  }) async {
    final response = await _client.post(
      ApiEndpoints.orderPhysicalCard,
      data: {
        'currency': currency,
        'cardholder_name': cardholderName,
        'shipping_address': shippingAddress,
        'express_shipping': expressShipping,
      },
    );
    if (response.statusCode == 201) {
      return response.data['card'];
    }
    throw Exception(response.data['error'] ?? 'Failed to order physical card');
  }
  
  /// Activer une carte
  Future<void> activateCard(String cardId) async {
    final response = await _client.post(ApiEndpoints.activateCard(cardId));
    if (response.statusCode != 200) {
      throw Exception(response.data['error'] ?? 'Failed to activate card');
    }
  }
  
  /// Désactiver une carte
  Future<void> deactivateCard(String cardId) async {
    final response = await _client.post(ApiEndpoints.deactivateCard(cardId));
    if (response.statusCode != 200) {
      throw Exception(response.data['error'] ?? 'Failed to deactivate card');
    }
  }
  
  /// Geler une carte
  Future<void> freezeCard(String cardId, String reason) async {
    final response = await _client.post(
      ApiEndpoints.freezeCard(cardId),
      data: {'reason': reason},
    );
    if (response.statusCode != 200) {
      throw Exception(response.data['error'] ?? 'Failed to freeze card');
    }
  }
  
  /// Dégeler une carte
  Future<void> unfreezeCard(String cardId) async {
    final response = await _client.post(ApiEndpoints.unfreezeCard(cardId));
    if (response.statusCode != 200) {
      throw Exception(response.data['error'] ?? 'Failed to unfreeze card');
    }
  }
  
  /// Bloquer une carte (définitif)
  Future<void> blockCard(String cardId, String reason) async {
    final response = await _client.post(
      ApiEndpoints.blockCard(cardId),
      data: {'reason': reason},
    );
    if (response.statusCode != 200) {
      throw Exception(response.data['error'] ?? 'Failed to block card');
    }
  }
  
  /// Charger une carte
  Future<void> loadCard({
    required String cardId,
    required double amount,
    required String sourceWalletId,
    String? description,
  }) async {
    final response = await _client.post(
      ApiEndpoints.loadCard(cardId),
      data: {
        'amount': amount,
        'source_wallet_id': sourceWalletId,
        if (description != null) 'description': description,
      },
    );
    if (response.statusCode != 200) {
      throw Exception(response.data['error'] ?? 'Failed to load card');
    }
  }
  
  /// Définir le PIN de la carte
  Future<void> setCardPIN(String cardId, String pin) async {
    final response = await _client.post(
      ApiEndpoints.setCardPIN(cardId),
      data: {'pin': pin},
    );
    if (response.statusCode != 200) {
      throw Exception(response.data['error'] ?? 'Failed to set PIN');
    }
  }
  
  /// Récupérer les limites de la carte
  Future<Map<String, dynamic>> getCardLimits(String cardId) async {
    final response = await _client.get(ApiEndpoints.cardLimits(cardId));
    if (response.statusCode == 200) {
      return response.data['limits'];
    }
    throw Exception('Failed to get card limits');
  }
  
  /// Mettre à jour les limites de la carte
  Future<void> updateCardLimits(String cardId, Map<String, dynamic> limits) async {
    final response = await _client.put(
      ApiEndpoints.cardLimits(cardId),
      data: limits,
    );
    if (response.statusCode != 200) {
      throw Exception(response.data['error'] ?? 'Failed to update limits');
    }
  }
  
  /// Récupérer les transactions de la carte
  Future<List<Map<String, dynamic>>> getCardTransactions(
    String cardId, {
    int limit = 50,
    int offset = 0,
  }) async {
    final response = await _client.get(
      ApiEndpoints.cardTransactions(cardId),
      queryParameters: {
        'limit': limit,
        'offset': offset,
      },
    );
    if (response.statusCode == 200) {
      return List<Map<String, dynamic>>.from(response.data['transactions']);
    }
    throw Exception('Failed to get transactions');
  }
  
  /// Récupérer le solde de la carte
  Future<double> getCardBalance(String cardId) async {
    final response = await _client.get(ApiEndpoints.cardBalance(cardId));
    if (response.statusCode == 200) {
      return (response.data['balance'] as num).toDouble();
    }
    throw Exception('Failed to get balance');
  }
  
  /// Récupérer le statut de livraison
  Future<Map<String, dynamic>> getShippingStatus(String cardId) async {
    final response = await _client.get(ApiEndpoints.shippingStatus(cardId));
    if (response.statusCode == 200) {
      return response.data['shipping_status'];
    }
    throw Exception('Failed to get shipping status');
  }
  
  /// Créer une carte cadeau
  Future<Map<String, dynamic>> createGiftCard({
    required double amount,
    required String currency,
    String? recipientEmail,
    String? recipientPhone,
    String? message,
    required String design,
    required String sourceWalletId,
  }) async {
    final response = await _client.post(
      ApiEndpoints.giftCards,
      data: {
        'amount': amount,
        'currency': currency,
        if (recipientEmail != null) 'recipient_email': recipientEmail,
        if (recipientPhone != null) 'recipient_phone': recipientPhone,
        if (message != null) 'message': message,
        'design': design,
        'source_wallet_id': sourceWalletId,
      },
    );
    if (response.statusCode == 201) {
      return response.data['gift_card'];
    }
    throw Exception(response.data['error'] ?? 'Failed to create gift card');
  }
  
  /// Récupérer les cartes cadeaux
  Future<List<Map<String, dynamic>>> getGiftCards() async {
    final response = await _client.get(ApiEndpoints.giftCards);
    if (response.statusCode == 200) {
      return List<Map<String, dynamic>>.from(response.data['gift_cards']);
    }
    throw Exception('Failed to get gift cards');
  }
  
  /// Échanger une carte cadeau
  Future<void> redeemGiftCard({
    required String code,
    required String targetWalletId,
  }) async {
    final response = await _client.post(
      ApiEndpoints.redeemGiftCard,
      data: {
        'code': code,
        'target_wallet_id': targetWalletId,
      },
    );
    if (response.statusCode != 200) {
      throw Exception(response.data['error'] ?? 'Failed to redeem gift card');
    }
  }
}
