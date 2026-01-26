import 'package:dio/dio.dart';
import '../api_client.dart';

class ShopApiService {
  final ApiClient _client = ApiClient();

  // Products
  Future<List<dynamic>> getProducts({String? categoryId, String? shopId}) async {
    try {
      final Map<String, dynamic> query = {};
      if (categoryId != null) query['category_id'] = categoryId;
      if (shopId != null) query['shop_id'] = shopId;

      final response = await _client.get('/shop-service/api/v1/products', queryParameters: query);
      return response.data['products'] ?? [];
    } catch (e) {
      throw _handleError(e);
    }
  }

  Future<Map<String, dynamic>> getProduct(String id) async {
    try {
      final response = await _client.get('/shop-service/api/v1/products/$id');
      return response.data;
    } catch (e) {
      throw _handleError(e);
    }
  }

  Future<Map<String, dynamic>> createProduct(Map<String, dynamic> data) async {
    try {
      // Handle image upload if present
      if (data['image_path'] != null) {
        final formData = FormData.fromMap({
          ...data,
          'image': await MultipartFile.fromFile(data['image_path']),
        });
        final response = await _client.post('/shop-service/api/v1/products', data: formData);
        return response.data;
      }
      
      final response = await _client.post('/shop-service/api/v1/products', data: data);
      return response.data;
    } catch (e) {
      throw _handleError(e);
    }
  }

  Future<void> deleteProduct(String id) async {
    try {
      await _client.delete('/shop-service/api/v1/products/$id');
    } catch (e) {
      throw _handleError(e);
    }
  }

  // Orders
  Future<List<dynamic>> getOrders({String? status}) async {
    try {
      final Map<String, dynamic> query = {};
      if (status != null) query['status'] = status;
      
      final response = await _client.get('/shop-service/api/v1/orders', queryParameters: query);
      return response.data['orders'] ?? [];
    } catch (e) {
      throw _handleError(e);
    }
  }
  
  // Categories
  Future<List<dynamic>> getCategories() async {
    try {
      final response = await _client.get('/shop-service/api/v1/categories');
      return response.data['categories'] ?? [];
    } catch (e) {
      throw _handleError(e);
    }
  }

  String _handleError(dynamic e) {
    if (e is DioException) {
      return e.response?.data['message'] ?? e.message ?? 'Unknown error';
    }
    return e.toString();
  }
}
