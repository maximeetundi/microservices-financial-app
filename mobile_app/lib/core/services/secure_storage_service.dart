import 'package:flutter_secure_storage/flutter_secure_storage.dart';
import 'package:local_auth/local_auth.dart';
import 'dart:convert';

/// Service de stockage sécurisé pour les données sensibles
class SecureStorageService {
  static final SecureStorageService _instance = SecureStorageService._internal();
  factory SecureStorageService() => _instance;
  
  SecureStorageService._internal();
  
  final FlutterSecureStorage _storage = const FlutterSecureStorage(
    aOptions: AndroidOptions(
      encryptedSharedPreferences: true,
    ),
    iOptions: IOSOptions(
      accessibility: KeychainAccessibility.first_unlock_this_device,
    ),
  );
  
  final LocalAuthentication _localAuth = LocalAuthentication();
  
  // Keys
  static const String _accessTokenKey = 'access_token';
  static const String _refreshTokenKey = 'refresh_token';
  static const String _userIdKey = 'user_id';
  static const String _biometricEnabledKey = 'biometric_enabled';
  static const String _pinCodeKey = 'pin_code';
  static const String _lastLoginKey = 'last_login';
  
  // Token Management
  Future<void> saveTokens(String accessToken, String refreshToken) async {
    await _storage.write(key: _accessTokenKey, value: accessToken);
    await _storage.write(key: _refreshTokenKey, value: refreshToken);
    await _storage.write(key: _lastLoginKey, value: DateTime.now().toIso8601String());
  }
  
  Future<String?> getAccessToken() async {
    return await _storage.read(key: _accessTokenKey);
  }
  
  Future<String?> getRefreshToken() async {
    return await _storage.read(key: _refreshTokenKey);
  }
  
  Future<void> clearTokens() async {
    await _storage.delete(key: _accessTokenKey);
    await _storage.delete(key: _refreshTokenKey);
  }
  
  Future<bool> hasValidSession() async {
    final token = await getAccessToken();
    return token != null;
  }
  
  // User ID
  Future<void> saveUserId(String userId) async {
    await _storage.write(key: _userIdKey, value: userId);
  }
  
  Future<String?> getUserId() async {
    return await _storage.read(key: _userIdKey);
  }
  
  // Biometric Authentication
  Future<bool> isBiometricAvailable() async {
    try {
      final isAvailable = await _localAuth.canCheckBiometrics;
      final isDeviceSupported = await _localAuth.isDeviceSupported();
      return isAvailable && isDeviceSupported;
    } catch (e) {
      return false;
    }
  }
  
  Future<List<BiometricType>> getAvailableBiometrics() async {
    try {
      return await _localAuth.getAvailableBiometrics();
    } catch (e) {
      return [];
    }
  }
  
  Future<bool> authenticateWithBiometrics({
    String reason = 'Authentifiez-vous pour continuer',
  }) async {
    try {
      return await _localAuth.authenticate(
        localizedReason: reason,
        options: const AuthenticationOptions(
          stickyAuth: true,
          biometricOnly: false,
        ),
      );
    } catch (e) {
      return false;
    }
  }
  
  Future<void> enableBiometric(bool enabled) async {
    await _storage.write(key: _biometricEnabledKey, value: enabled.toString());
  }
  
  Future<bool> isBiometricEnabled() async {
    final value = await _storage.read(key: _biometricEnabledKey);
    return value == 'true';
  }
  
  // PIN Code
  Future<void> savePinCode(String pin) async {
    await _storage.write(key: _pinCodeKey, value: pin);
  }
  
  Future<bool> verifyPinCode(String pin) async {
    final storedPin = await _storage.read(key: _pinCodeKey);
    return storedPin == pin;
  }
  
  Future<bool> hasPinCode() async {
    final pin = await _storage.read(key: _pinCodeKey);
    return pin != null;
  }
  
  // Clear all data (logout)
  Future<void> clearAll() async {
    await _storage.deleteAll();
  }
  
  // Last login
  Future<DateTime?> getLastLogin() async {
    final value = await _storage.read(key: _lastLoginKey);
    if (value != null) {
      return DateTime.parse(value);
    }
    return null;
  }

  // User Profile Caching
  static const String _userProfileKey = 'user_profile_cache';

  Future<void> saveUserProfile(Map<String, dynamic> json) async {
    await _storage.write(key: _userProfileKey, value: jsonEncode(json));
  }

  Future<Map<String, dynamic>?> getUserProfile() async {
    final jsonStr = await _storage.read(key: _userProfileKey);
    if (jsonStr != null) {
      try {
        return jsonDecode(jsonStr) as Map<String, dynamic>;
      } catch (e) {
        return null;
      }
    }
    return null;
  }

  // Wallet Caching
  static const String _walletsKey = 'wallets_cache';

  Future<void> saveWallets(List<dynamic> json) async {
    await _storage.write(key: _walletsKey, value: jsonEncode(json));
  }

  Future<List<dynamic>?> getWallets() async {
    final jsonStr = await _storage.read(key: _walletsKey);
    if (jsonStr != null) {
      try {
        return jsonDecode(jsonStr) as List<dynamic>;
      } catch (e) {
        return null;
      }
    }
    return null;
  }
}
