import '../models/subscription.dart';
import '../../../../core/services/subscription_api_service.dart';

class SubscriptionRepository {
  final SubscriptionApiService _apiService;

  SubscriptionRepository(this._apiService);

  // Get My Subscriptions
  Future<List<Subscription>> getSubscriptions() async {
    final response = await _apiService.getMySubscriptions();
    // Expecting response.data to be a List or { "subscriptions": [] }
    // Adjust based on actual API response format (usually unified in BaseApiService)
    final data = response.data; 
    if (data is List) {
      return data.map((json) => Subscription.fromJson(json)).toList();
    } else if (data is Map && data.containsKey('subscriptions')) {
       return (data['subscriptions'] as List).map((json) => Subscription.fromJson(json)).toList();
    }
    return [];
  }

  // Link Subscription
  Future<void> linkSubscription({
    required String enterpriseId,
    required String matricule,
    required Map<String, dynamic> formData,
  }) async {
    await _apiService.linkSubscription({
      'enterprise_id': enterpriseId,
      'external_id': matricule,
      'form_data': formData,
    });
  }

  // Get Pending Bills (Invoices)
  Future<List<Map<String, dynamic>>> getPendingBills() async {
    final response = await _apiService.getPendingBills();
    final data = response.data;
    // Assuming a list of invoices
    if (data is List) {
      return List<Map<String, dynamic>>.from(data);
    }
    return [];
  }

  // Pay Bill
  Future<void> payBill(String invoiceId) async {
    await _apiService.payBill({'invoice_id': invoiceId});
  }
}
