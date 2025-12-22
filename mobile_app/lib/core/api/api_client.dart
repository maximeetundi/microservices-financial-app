import 'dart:io';
import 'package:dio/dio.dart';
import 'package:flutter/material.dart';
import 'package:flutter_secure_storage/flutter_secure_storage.dart';

class ApiClient {
  // Production API URL
  static const String _baseUrlProduction = 'https://api.app.maximeetundi.store';
  // For local development, use these:
  // static const String _baseUrlAndroid = 'http://10.0.2.2:8080';
  // static const String _baseUrlIOS = 'http://localhost:8080';
  
  late final Dio _dio;
  final FlutterSecureStorage _storage = const FlutterSecureStorage();
  
  // Global navigator key for logout navigation
  static final GlobalKey<NavigatorState> navigatorKey = GlobalKey<NavigatorState>();
  
  // Callback for logout - will be set by the app
  static VoidCallback? onLogout;
  
  // Flag to prevent multiple logout calls
  bool _isLoggingOut = false;
  
  static final ApiClient _instance = ApiClient._internal();
  factory ApiClient() => _instance;
  
  ApiClient._internal() {
    _dio = Dio(BaseOptions(
      baseUrl: _getBaseUrl(),
      connectTimeout: const Duration(seconds: 15),
      receiveTimeout: const Duration(seconds: 15),
      headers: {
        'Content-Type': 'application/json',
        'Accept': 'application/json',
      },
    ));
    
    _setupInterceptors();
  }
  
  static String _getBaseUrl() {
    // Always use production URL
    return _baseUrlProduction;
  }
  
  void _setupInterceptors() {
    // Auth Interceptor - Add JWT token to requests
    _dio.interceptors.add(InterceptorsWrapper(
      onRequest: (options, handler) async {
        final token = await _storage.read(key: 'access_token');
        if (token != null) {
          options.headers['Authorization'] = 'Bearer $token';
        }
        return handler.next(options);
      },
      onResponse: (response, handler) {
        return handler.next(response);
      },
      onError: (error, handler) async {
        if (error.response?.statusCode == 401 && !_isLoggingOut) {
          // Token expired, try to refresh
          final refreshed = await _refreshToken();
          if (refreshed) {
            // Retry the request
            final retryResponse = await _retry(error.requestOptions);
            return handler.resolve(retryResponse);
          } else {
            // Refresh failed - Auto logout
            await _handleLogout();
          }
        }
        return handler.next(error);
      },
    ));
    
    // Logging Interceptor (for development)
    _dio.interceptors.add(LogInterceptor(
      requestBody: true,
      responseBody: true,
      error: true,
    ));
  }
  
  Future<void> _handleLogout() async {
    if (_isLoggingOut) return;
    _isLoggingOut = true;
    
    try {
      // Clear all tokens
      await clearTokens();
      
      // Call logout callback if set (this will navigate to login)
      if (onLogout != null) {
        onLogout!();
      }
    } finally {
      // Reset flag after a delay to prevent rapid re-triggers
      Future.delayed(const Duration(seconds: 2), () {
        _isLoggingOut = false;
      });
    }
  }
  
  Future<bool> _refreshToken() async {
    try {
      final refreshToken = await _storage.read(key: 'refresh_token');
      if (refreshToken == null) return false;
      
      final response = await Dio().post(
        '${_getBaseUrl()}/auth-service/api/v1/auth/refresh',
        data: {'refresh_token': refreshToken},
      );
      
      if (response.statusCode == 200) {
        await _storage.write(
          key: 'access_token',
          value: response.data['access_token'],
        );
        await _storage.write(
          key: 'refresh_token',
          value: response.data['refresh_token'],
        );
        return true;
      }
    } catch (e) {
      // Refresh failed
      print('Token refresh failed: $e');
    }
    return false;
  }
  
  Future<Response> _retry(RequestOptions requestOptions) async {
    final token = await _storage.read(key: 'access_token');
    final options = Options(
      method: requestOptions.method,
      headers: {
        ...requestOptions.headers,
        'Authorization': 'Bearer $token',
      },
    );
    
    return _dio.request(
      requestOptions.path,
      data: requestOptions.data,
      queryParameters: requestOptions.queryParameters,
      options: options,
    );
  }
  
  // HTTP Methods
  Future<Response> get(String path, {Map<String, dynamic>? queryParameters}) {
    return _dio.get(path, queryParameters: queryParameters);
  }
  
  Future<Response> post(String path, {dynamic data}) {
    return _dio.post(path, data: data);
  }
  
  Future<Response> put(String path, {dynamic data}) {
    return _dio.put(path, data: data);
  }
  
  Future<Response> patch(String path, {dynamic data}) {
    return _dio.patch(path, data: data);
  }
  
  Future<Response> delete(String path, {dynamic data}) {
    return _dio.delete(path, data: data);
  }
  
  // Token Management
  Future<void> saveTokens(String accessToken, String refreshToken) async {
    await _storage.write(key: 'access_token', value: accessToken);
    await _storage.write(key: 'refresh_token', value: refreshToken);
  }
  
  Future<void> clearTokens() async {
    await _storage.delete(key: 'access_token');
    await _storage.delete(key: 'refresh_token');
  }
  
  Future<bool> hasValidToken() async {
    final token = await _storage.read(key: 'access_token');
    return token != null;
  }
}
