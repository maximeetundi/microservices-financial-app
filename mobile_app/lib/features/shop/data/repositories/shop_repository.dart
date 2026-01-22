import 'dart:convert';
import 'package:http/http.dart' as http;
import 'package:flutter_secure_storage/flutter_secure_storage.dart';
import '../models/shop_models.dart';

class ShopRepository {
  final String baseUrl;
  final FlutterSecureStorage _storage = const FlutterSecureStorage();

  ShopRepository({required this.baseUrl});

  Future<String?> _getToken() async {
    return await _storage.read(key: 'auth_token');
  }

  Future<Map<String, String>> _getHeaders() async {
    final token = await _getToken();
    return {
      'Content-Type': 'application/json',
      if (token != null) 'Authorization': 'Bearer $token',
    };
  }

  // Shops
  Future<List<Shop>> getShops({int page = 1, int pageSize = 20, String? search}) async {
    final params = 'page=$page&page_size=$pageSize${search != null ? '&search=$search' : ''}';
    final response = await http.get(
      Uri.parse('$baseUrl/shop-service/api/v1/shops?$params'),
      headers: await _getHeaders(),
    );

    if (response.statusCode == 200) {
      final data = json.decode(response.body);
      return (data['shops'] as List<dynamic>?)
              ?.map((e) => Shop.fromJson(e))
              .toList() ??
          [];
    }
    throw Exception('Failed to load shops');
  }

  Future<Shop> getShop(String slug) async {
    final response = await http.get(
      Uri.parse('$baseUrl/shop-service/api/v1/shops/$slug'),
      headers: await _getHeaders(),
    );

    if (response.statusCode == 200) {
      return Shop.fromJson(json.decode(response.body));
    }
    throw Exception('Failed to load shop');
  }

  // Products
  Future<List<Product>> getProducts(String shopSlug, {int page = 1, String? category}) async {
    final params = 'page=$page&page_size=50${category != null ? '&category=$category' : ''}';
    final response = await http.get(
      Uri.parse('$baseUrl/shop-service/api/v1/shops/$shopSlug/products?$params'),
      headers: await _getHeaders(),
    );

    if (response.statusCode == 200) {
      final data = json.decode(response.body);
      return (data['products'] as List<dynamic>?)
              ?.map((e) => Product.fromJson(e))
              .toList() ??
          [];
    }
    throw Exception('Failed to load products');
  }

  Future<Product> getProduct(String shopSlug, String productSlug) async {
    final response = await http.get(
      Uri.parse('$baseUrl/shop-service/api/v1/shops/$shopSlug/products/$productSlug'),
      headers: await _getHeaders(),
    );

    if (response.statusCode == 200) {
      return Product.fromJson(json.decode(response.body));
    }
    throw Exception('Failed to load product');
  }

  // Categories
  Future<List<Category>> getCategories(String shopSlug) async {
    final response = await http.get(
      Uri.parse('$baseUrl/shop-service/api/v1/shops/$shopSlug/categories'),
      headers: await _getHeaders(),
    );

    if (response.statusCode == 200) {
      final data = json.decode(response.body);
      return (data['categories'] as List<dynamic>?)
              ?.map((e) => Category.fromJson(e))
              .toList() ??
          [];
    }
    throw Exception('Failed to load categories');
  }

  // Orders
  Future<Order> createOrder({
    required String shopId,
    required List<CartItem> items,
    required String walletId,
    required String deliveryType,
  }) async {
    final response = await http.post(
      Uri.parse('$baseUrl/shop-service/api/v1/orders'),
      headers: await _getHeaders(),
      body: json.encode({
        'shop_id': shopId,
        'items': items.map((item) => item.toJson()).toList(),
        'wallet_id': walletId,
        'delivery_type': deliveryType,
      }),
    );

    if (response.statusCode == 201) {
      return Order.fromJson(json.decode(response.body));
    }
    throw Exception(json.decode(response.body)['error'] ?? 'Failed to create order');
  }

  Future<List<Order>> getMyOrders({int page = 1}) async {
    final response = await http.get(
      Uri.parse('$baseUrl/shop-service/api/v1/orders?page=$page'),
      headers: await _getHeaders(),
    );

    if (response.statusCode == 200) {
      final data = json.decode(response.body);
      return (data['orders'] as List<dynamic>?)
              ?.map((e) => Order.fromJson(e))
              .toList() ??
          [];
    }
    throw Exception('Failed to load orders');
  }

  Future<Order> getOrder(String orderId) async {
    final response = await http.get(
      Uri.parse('$baseUrl/shop-service/api/v1/orders/$orderId'),
      headers: await _getHeaders(),
    );

    if (response.statusCode == 200) {
      return Order.fromJson(json.decode(response.body));
    }
    throw Exception('Failed to load order');
  }
}
