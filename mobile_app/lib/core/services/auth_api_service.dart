import '../api/api_client.dart';
import '../api/api_endpoints.dart';

/// Service d'authentification
class AuthApiService {
  final ApiClient _client = ApiClient();
  
  /// Connexion avec email et mot de passe
  Future<Map<String, dynamic>> login(String email, String password) async {
    final response = await _client.post(
      ApiEndpoints.login,
      data: {
        'email': email,
        'password': password,
      },
    );
    
    if (response.statusCode == 200) {
      final data = response.data;
      await _client.saveTokens(
        data['access_token'],
        data['refresh_token'],
      );
      // Reset the logout flag so API calls work properly
      _client.resetLogoutFlag();
      return data;
    }
    throw Exception(response.data['error'] ?? 'Login failed');
  }
  
  /// Inscription d'un nouvel utilisateur
  Future<Map<String, dynamic>> register({
    required String email,
    required String password,
    required String firstName,
    required String lastName,
    String? phone,
    String? dateOfBirth,
    String? country,
    String? currency,
  }) async {
    final response = await _client.post(
      ApiEndpoints.register,
      data: {
        'email': email,
        'password': password,
        'first_name': firstName,
        'last_name': lastName,
        if (phone != null) 'phone': phone,
        if (dateOfBirth != null) 'date_of_birth': dateOfBirth,
        if (country != null) 'country': country,
        if (currency != null) 'currency': currency,
      },
    );
    
    if (response.statusCode == 201) {
      return response.data;
    }
    throw Exception(response.data['error'] ?? 'Registration failed');
  }
  
  /// Déconnexion
  Future<void> logout() async {
    try {
      await _client.post(ApiEndpoints.logout);
    } finally {
      await _client.clearTokens();
    }
  }
  
  /// Vérifier si l'utilisateur est connecté
  Future<bool> isAuthenticated() async {
    return await _client.hasValidToken();
  }
  
  /// Récupérer le profil utilisateur
  Future<Map<String, dynamic>> getProfile() async {
    final response = await _client.get(ApiEndpoints.profile);
    if (response.statusCode == 200) {
      return response.data['user'];
    }
    throw Exception('Failed to get profile');
  }
  
  /// Mettre à jour le profil
  Future<Map<String, dynamic>> updateProfile(Map<String, dynamic> data) async {
    final response = await _client.put(ApiEndpoints.updateProfile, data: data);
    if (response.statusCode == 200) {
      return response.data['user'];
    }
    throw Exception('Failed to update profile');
  }
  
  /// Changer le mot de passe
  Future<void> changePassword(String currentPassword, String newPassword) async {
    final response = await _client.post(
      ApiEndpoints.changePassword,
      data: {
        'current_password': currentPassword,
        'new_password': newPassword,
      },
    );
    if (response.statusCode != 200) {
      throw Exception(response.data['error'] ?? 'Failed to change password');
    }
  }
  
  /// Mot de passe oublié
  Future<void> forgotPassword(String email) async {
    final response = await _client.post(
      ApiEndpoints.forgotPassword,
      data: {'email': email},
    );
    if (response.statusCode != 200) {
      throw Exception('Failed to send reset email');
    }
  }
  
  /// Activer 2FA
  Future<Map<String, dynamic>> enable2FA() async {
    final response = await _client.post(ApiEndpoints.enable2FA);
    if (response.statusCode == 200) {
      return response.data;
    }
    throw Exception('Failed to enable 2FA');
  }
  
  /// Vérifier code 2FA
  Future<bool> verify2FA(String code) async {
    final response = await _client.post(
      ApiEndpoints.verify2FA,
      data: {'code': code},
    );
    return response.statusCode == 200;
  }

  /// Rechercher un utilisateur (Email ou Téléphone)
  Future<Map<String, dynamic>> lookupUser({String? email, String? phone}) async {
    final Map<String, dynamic> queryParams = {};
    if (email != null) queryParams['email'] = email;
    if (phone != null) queryParams['phone'] = phone;

    final response = await _client.get(
      ApiEndpoints.lookup,
      queryParameters: queryParams,
    );
    
    if (response.statusCode == 200) {
      return response.data;
    }
    throw Exception('User not found');
  }
}
