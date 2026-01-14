import 'package:dio/dio.dart';
import 'base_api_service.dart';

class SubscriptionApiService extends BaseApiService {
  static const String _basePath = '/enterprise-service/api/v1';

  // Fetch my subscriptions
  Future<dynamic> getMySubscriptions() async {
    return get('$_basePath/subscriptions/me');
  }

  // Link a subscription (Join by Matricule)
  Future<dynamic> linkSubscription(Map<String, dynamic> data) async {
    return post('$_basePath/subscriptions/link', data: data);
  }

  // Get Pending Bills
  Future<dynamic> getPendingBills() async {
    return get('$_basePath/billing/my-bills');
  }

  // Pay a Bill
  Future<dynamic> payBill(Map<String, dynamic> data) async {
    // This might route to wallet-service eventually, but enterprise-service 
    // can orchestrate the B2B payment request
    return post('$_basePath/billing/pay', data: data);
  }
}
