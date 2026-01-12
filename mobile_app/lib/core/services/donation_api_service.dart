import '../api/api_client.dart';

class DonationApiService {
  final ApiClient _client = ApiClient();

  // === Campaigns ===

  /// Create a new campaign
  Future<Map<String, dynamic>> createCampaign(Map<String, dynamic> data) async {
    final response = await _client.post(
      '/donation-service/api/v1/campaigns',
      data: data,
    );
    if (response.statusCode == 201 || response.statusCode == 200) {
      return response.data['campaign'] ?? response.data;
    }
    throw Exception(response.data['error'] ?? 'Failed to create campaign');
  }

  /// Get campaigns
  Future<List<dynamic>> getCampaigns({int limit = 20, int offset = 0, String? creatorId}) async {
    final Map<String, dynamic> params = {'limit': limit, 'offset': offset};
    if (creatorId != null && creatorId.isNotEmpty) {
      params['creator_id'] = creatorId;
    }

    final response = await _client.get(
      '/donation-service/api/v1/campaigns',
      queryParameters: params,
    );
    if (response.statusCode == 200) {
      return response.data['campaigns'] ?? [];
    }
    throw Exception('Failed to load campaigns');
  }

  /// Get campaign by ID
  Future<Map<String, dynamic>> getCampaign(String id) async {
    final response = await _client.get('/donation-service/api/v1/campaigns/$id');
    if (response.statusCode == 200) {
      return response.data;
    }
    throw Exception('Campaign not found');
  }

  /// Update campaign
  Future<void> updateCampaign(String id, Map<String, dynamic> data) async {
    final response = await _client.put(
      '/donation-service/api/v1/campaigns/$id',
      data: data,
    );
    if (response.statusCode != 200) {
      throw Exception(response.data['error'] ?? 'Failed to update campaign');
    }
  }

  // === Donations ===

  /// Initiate a donation
  Future<Map<String, dynamic>> initiateDonation({
    required String campaignId,
    required double amount,
    required String currency,
    required String walletId,
    required String pin,
    String message = '',
    bool isAnonymous = false,
    String frequency = 'one_time',
  }) async {
    final response = await _client.post(
      '/donation-service/api/v1/donations',
      data: {
        'campaign_id': campaignId,
        'amount': amount,
        'currency': currency,
        'wallet_id': walletId,
        'pin': pin,
        'message': message,
        'is_anonymous': isAnonymous,
        'frequency': frequency,
      },
    );
    if (response.statusCode == 201 || response.statusCode == 200) {
      return response.data;
    }
    throw Exception(response.data['error'] ?? 'Failed to donate');
  }

  /// Get donations for a campaign
  Future<List<dynamic>> getDonations(String campaignId, {int limit = 20, int offset = 0}) async {
    final response = await _client.get(
      '/donation-service/api/v1/donations',
      queryParameters: {'campaign_id': campaignId, 'limit': limit, 'offset': offset},
    );
    if (response.statusCode == 200) {
      return response.data['donations'] ?? [];
    }
    throw Exception('Failed to load donations');
  }

  /// Upload Campaign Media (Image or Video)
  Future<String> uploadMedia(dynamic file) async {
    // Assuming file is File from dart:io or similar handled by ApiClient
    final response = await _client.uploadFile(
      '/donation-service/api/v1/upload', 
      file,
      fieldName: 'file',
    );
    if (response.statusCode == 200) {
      return response.data['url'] ?? '';
    }
     // Fallback if image upload service is shared or mocked
    return '';
  }
}
