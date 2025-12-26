import '../api/api_client.dart';

/// Base class for API services providing common HTTP methods
abstract class BaseApiService {
  final ApiClient _client = ApiClient();

  /// Perform authenticated GET request (auto-prepends /api/v1)
  Future<Map<String, dynamic>> authenticatedGet(String path) async {
    final response = await _client.get('/api/v1$path');
    if (response.statusCode == 200) {
      return response.data is Map<String, dynamic> 
          ? response.data 
          : {'data': response.data};
    }
    throw Exception(response.data['error'] ?? 'Request failed');
  }

  /// Perform authenticated GET request with raw path (no auto-prefix)
  Future<Map<String, dynamic>> authenticatedGetRaw(String path) async {
    final response = await _client.get(path);
    if (response.statusCode == 200) {
      return response.data is Map<String, dynamic> 
          ? response.data 
          : {'data': response.data};
    }
    throw Exception(response.data['error'] ?? 'Request failed');
  }

  /// Perform authenticated POST request (auto-prepends /api/v1)
  Future<Map<String, dynamic>> authenticatedPost(String path, Map<String, dynamic> data) async {
    final response = await _client.post('/api/v1$path', data: data);
    if (response.statusCode == 200 || response.statusCode == 201) {
      return response.data is Map<String, dynamic> 
          ? response.data 
          : {'data': response.data};
    }
    throw Exception(response.data['error'] ?? 'Request failed');
  }

  /// Perform authenticated POST request with raw path (no auto-prefix)
  Future<Map<String, dynamic>> authenticatedPostRaw(String path, Map<String, dynamic> data) async {
    final response = await _client.post(path, data: data);
    if (response.statusCode == 200 || response.statusCode == 201) {
      return response.data is Map<String, dynamic> 
          ? response.data 
          : {'data': response.data};
    }
    throw Exception(response.data['error'] ?? 'Request failed');
  }

  /// Perform authenticated DELETE request (auto-prepends /api/v1)
  Future<void> authenticatedDelete(String path) async {
    final response = await _client.delete('/api/v1$path');
    if (response.statusCode != 200 && response.statusCode != 204) {
      throw Exception(response.data['error'] ?? 'Delete failed');
    }
  }

  /// Perform authenticated DELETE request with raw path (no auto-prefix)
  Future<void> authenticatedDeleteRaw(String path) async {
    final response = await _client.delete(path);
    if (response.statusCode != 200 && response.statusCode != 204) {
      throw Exception(response.data['error'] ?? 'Delete failed');
    }
  }

  /// Perform unauthenticated GET request (public endpoints, auto-prepends /api/v1)
  Future<Map<String, dynamic>> get(String path) async {
    final response = await _client.get('/api/v1$path');
    if (response.statusCode == 200) {
      return response.data is Map<String, dynamic> 
          ? response.data 
          : {'data': response.data};
    }
    throw Exception(response.data['error'] ?? 'Request failed');
  }

  /// Perform unauthenticated GET request with raw path (no auto-prefix)
  Future<Map<String, dynamic>> getRaw(String path) async {
    final response = await _client.get(path);
    if (response.statusCode == 200) {
      return response.data is Map<String, dynamic> 
          ? response.data 
          : {'data': response.data};
    }
    throw Exception(response.data['error'] ?? 'Request failed');
  }
}

