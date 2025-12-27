import 'dart:async';
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
  
  // Mutex for token refresh - prevents multiple concurrent refresh calls
  Completer<bool>? _refreshCompleter;
  bool _isRefreshing = false;
  
  // Auth endpoints that should NOT trigger token refresh on 401
  static const List<String> _authEndpoints = [
    '/auth/login',
    '/auth/register',
    '/auth/refresh',
    '/auth/forgot-password',
    '/auth/reset-password',
  ];
  
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
  
  /// Check if a URL is an auth endpoint (should skip refresh on 401)
  bool _isAuthEndpoint(String url) {
    return _authEndpoints.any((endpoint) => url.contains(endpoint));
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
        final requestPath = error.requestOptions.path;
        
        // Skip refresh for auth endpoints to prevent loops
        if (_isAuthEndpoint(requestPath)) {
          debugPrint('ApiClient: Auth endpoint 401, not attempting refresh: $requestPath');
          return handler.next(error);
        }
        
        if (error.response?.statusCode == 401 && !_isLoggingOut) {
          debugPrint('ApiClient: Got 401 on $requestPath, attempting token refresh...');
          
          // Token expired, try to refresh
          final refreshed = await _refreshTokenWithMutex();
          if (refreshed) {
            debugPrint('ApiClient: Token refreshed successfully, retrying request...');
            // Retry the request
            try {
              final retryResponse = await _retry(error.requestOptions);
              return handler.resolve(retryResponse);
            } catch (retryError) {
              debugPrint('ApiClient: Retry failed after refresh: $retryError');
              return handler.next(error);
            }
          } else {
            debugPrint('ApiClient: Token refresh failed, logging out...');
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
      logPrint: (log) => debugPrint('API: $log'),
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
  
  /// Refresh token with mutex to prevent concurrent refresh calls
  Future<bool> _refreshTokenWithMutex() async {
    // If already refreshing, wait for the result
    if (_isRefreshing && _refreshCompleter != null) {
      debugPrint('ApiClient: Already refreshing, waiting for result...');
      return await _refreshCompleter!.future;
    }
    
    // Start refresh
    _isRefreshing = true;
    _refreshCompleter = Completer<bool>();
    
    try {
      final result = await _doRefreshToken();
      _refreshCompleter!.complete(result);
      return result;
    } catch (e) {
      _refreshCompleter!.complete(false);
      return false;
    } finally {
      _isRefreshing = false;
      _refreshCompleter = null;
    }
  }
  
  Future<bool> _doRefreshToken() async {
    try {
      final refreshToken = await _storage.read(key: 'refresh_token');
      if (refreshToken == null) {
        debugPrint('ApiClient: No refresh token available');
        return false;
      }
      
      debugPrint('ApiClient: Sending refresh token request...');
      
      // Use a new Dio instance to avoid interceptors
      final response = await Dio().post(
        '${_getBaseUrl()}/auth-service/api/v1/auth/refresh',
        data: {'refresh_token': refreshToken},
        options: Options(
          headers: {
            'Content-Type': 'application/json',
          },
        ),
      );
      
      if (response.statusCode == 200) {
        final newAccessToken = response.data['access_token'];
        final newRefreshToken = response.data['refresh_token'];
        
        if (newAccessToken != null && newRefreshToken != null) {
          await _storage.write(key: 'access_token', value: newAccessToken);
          await _storage.write(key: 'refresh_token', value: newRefreshToken);
          debugPrint('ApiClient: Tokens saved successfully');
          return true;
        }
      }
      debugPrint('ApiClient: Refresh response invalid: ${response.statusCode}');
    } catch (e) {
      debugPrint('ApiClient: Token refresh error: $e');
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
  
  /// Upload a file using multipart form data
  Future<Response> uploadFile(
    String path,
    dynamic file, {
    String fieldName = 'file',
    Map<String, String>? extraFields,
  }) async {
    FormData formData;
    
    if (file is File) {
      final fileName = file.path.split('/').last;
      final extension = fileName.split('.').last.toLowerCase();
      
      // Determine content type based on file extension
      String? contentType;
      switch (extension) {
        case 'jpg':
        case 'jpeg':
          contentType = 'image/jpeg';
          break;
        case 'png':
          contentType = 'image/png';
          break;
        case 'gif':
          contentType = 'image/gif';
          break;
        case 'pdf':
          contentType = 'application/pdf';
          break;
        case 'heic':
        case 'heif':
          contentType = 'image/heic';
          break;
        default:
          contentType = 'application/octet-stream';
      }
      
      formData = FormData.fromMap({
        fieldName: await MultipartFile.fromFile(
          file.path,
          filename: fileName,
          contentType: DioMediaType.parse(contentType),
        ),
        if (extraFields != null) ...extraFields,
      });
    } else {
      throw ArgumentError('File must be of type File');
    }
    
    return _dio.post(
      path,
      data: formData,
      options: Options(
        contentType: 'multipart/form-data',
      ),
    );
  }
  
  // Token Management
  Future<void> saveTokens(String accessToken, String refreshToken) async {
    await _storage.write(key: 'access_token', value: accessToken);
    await _storage.write(key: 'refresh_token', value: refreshToken);
    debugPrint('ApiClient: Tokens saved after login');
  }
  
  Future<void> clearTokens() async {
    await _storage.delete(key: 'access_token');
    await _storage.delete(key: 'refresh_token');
    debugPrint('ApiClient: Tokens cleared');
  }
  
  Future<bool> hasValidToken() async {
    final token = await _storage.read(key: 'access_token');
    return token != null;
  }
  
  /// Force refresh tokens - can be called manually
  Future<bool> forceRefreshTokens() async {
    return await _refreshTokenWithMutex();
  }
  
  /// Reset logout flag - call after successful login
  void resetLogoutFlag() {
    _isLoggingOut = false;
  }
}

