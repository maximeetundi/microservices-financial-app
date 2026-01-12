import 'package:dio/dio.dart';
import 'package:flutter/foundation.dart';
import '../api/api_client.dart';
import '../api/api_endpoints.dart';
import 'secure_storage_service.dart';
import 'package:encrypt/encrypt.dart' as enc;

/// Service for managing 5-digit transaction PIN
class PinService {
  static final PinService _instance = PinService._internal();
  factory PinService() => _instance;
  PinService._internal();

  final ApiClient _client = ApiClient();
  final ApiClient _client = ApiClient();
  final SecureStorageService _storage = SecureStorageService();

  static const String _pinSetKey = 'pin_set';
  
  // Cache for RSA Public Key
  String? _publicKeyPem;

  /// Fetch Public Key for Encryption
  Future<String?> _getPublicKey() async {
    if (_publicKeyPem != null) return _publicKeyPem;
    try {
      final response = await _client.get(ApiEndpoints.publicKey);
      if (response.statusCode == 200) {
        _publicKeyPem = response.data['public_key'];
        return _publicKeyPem;
      }
    } catch (e) {
      debugPrint('Failed to fetch public key: $e');
    }
    return null;
  }

  /// Encrypt PIN using RSA-OAEP
  Future<String> _encryptPin(String pin) async {
    final pem = await _getPublicKey();
    if (pem == null) return pin; // Fallback to plain text if key fetch fails
    
    try {
       final parser = enc.RSAKeyParser();
       final publicKey = parser.parse(pem);
       final encrypter = enc.Encrypter(enc.RSA(publicKey: publicKey as dynamic, encoding: enc.RSAEncoding.OAEP, digest: enc.RSADigest.SHA256));
       final encrypted = encrypter.encrypt(pin);
       return encrypted.base64;
    } catch (e) {
       debugPrint('Encryption failed: $e');
       return pin;
    }
  }

  /// Check if user has set their PIN
  Future<bool> checkPinStatus() async {
    try {
      final response = await _client.get(ApiEndpoints.checkPinStatus);
      if (response.statusCode == 200) {
        final hasPin = response.data['has_pin'] ?? false;
        await _storage.write(_pinSetKey, hasPin.toString());
        return hasPin;
      }
      return false;
    } catch (e) {
      debugPrint('Error checking PIN status: $e');
      // Check local storage as fallback
      final stored = await _storage.read(_pinSetKey);
      return stored == 'true';
    }
  }

  /// Set up a new PIN (first time)
  Future<PinResult> setupPin(String pin, String confirmPin) async {
    try {
      if (pin != confirmPin) {
        return PinResult(success: false, message: 'Les PINs ne correspondent pas');
      }

      if (!_isValidPin(pin)) {
        return PinResult(success: false, message: 'Le PIN doit contenir exactement 5 chiffres');
      }

      final response = await _client.post(
        ApiEndpoints.setupPin,
        data: {
          'pin': pin,
          'confirm_pin': confirmPin,
        },
      );

      if (response.statusCode == 200) {
        await _storage.write(_pinSetKey, 'true');
        return PinResult(success: true, message: 'PIN défini avec succès');
      }
      
      return PinResult(
        success: false, 
        message: response.data['error'] ?? 'Échec de la définition du PIN',
      );
    } catch (e) {
      debugPrint('Error setting up PIN: $e');
      return PinResult(success: false, message: 'Erreur de connexion');
    }
  }

  /// Verify the PIN
  Future<PinVerifyResult> verifyPin(String pin) async {
    try {
      // Encrypt PIN before sending
      final encryptedPin = await _encryptPin(pin);
      
      final response = await _client.post(
        ApiEndpoints.verifyPin,
        data: {'pin': encryptedPin},
      );

      if (response.statusCode == 200) {
        return PinVerifyResult(
          valid: response.data['valid'] ?? false,
          attemptsLeft: response.data['attempts_left'],
          message: response.data['message'],
          encryptedPin: encryptedPin, // Pass the encrypted PIN
        );
      }

      return PinVerifyResult(
        valid: false,
        attemptsLeft: response.data['attempts_left'],
        message: response.data['message'] ?? 'PIN incorrect',
      );
    } on DioException catch (e) {
      // Handle 429 rate limit error (too many attempts)
      if (e.response?.statusCode == 429) {
        final data = e.response?.data;
        return PinVerifyResult(
          valid: false,
          message: data?['message'] ?? 'Trop de tentatives. Réessayez plus tard.',
          lockedUntil: data?['locked_until'],
        );
      }
      debugPrint('Error verifying PIN: $e');
      return PinVerifyResult(valid: false, message: 'Erreur de connexion');
    } catch (e) {
      debugPrint('Error verifying PIN: $e');
      return PinVerifyResult(valid: false, message: 'Erreur de connexion');
    }
  }

  /// Change the PIN
  Future<PinResult> changePin(String currentPin, String newPin, String confirmPin) async {
    try {
      if (newPin != confirmPin) {
        return PinResult(success: false, message: 'Les nouveaux PINs ne correspondent pas');
      }

      if (!_isValidPin(newPin)) {
        return PinResult(success: false, message: 'Le PIN doit contenir exactement 5 chiffres');
      }

      final response = await _client.post(
        ApiEndpoints.changePin,
        data: {
          'current_pin': currentPin,
          'new_pin': newPin,
          'confirm_pin': confirmPin,
        },
      );

      if (response.statusCode == 200) {
        return PinResult(success: true, message: 'PIN modifié avec succès');
      }
      
      return PinResult(
        success: false, 
        message: response.data['error'] ?? 'Échec de la modification du PIN',
      );
    } catch (e) {
      debugPrint('Error changing PIN: $e');
      return PinResult(success: false, message: 'Erreur de connexion');
    }
  }

  /// Check if PIN is valid format (5 digits)
  bool _isValidPin(String pin) {
    return pin.length == 5 && RegExp(r'^\d{5}$').hasMatch(pin);
  }

  /// Check if PIN is set locally (faster than API call)
  Future<bool> isPinSetLocally() async {
    final stored = await _storage.read(_pinSetKey);
    return stored == 'true';
  }

  /// Clear PIN status (on logout)
  Future<void> clearPinStatus() async {
    await _storage.delete(_pinSetKey);
  }
}

/// Result of PIN setup/change operations
class PinResult {
  final bool success;
  final String message;

  PinResult({required this.success, required this.message});
}

class PinVerifyResult {
  final bool valid;
  final int? attemptsLeft;
  final String? message;
  final String? lockedUntil;
  final String? encryptedPin; // Added field

  PinVerifyResult({
    required this.valid,
    this.attemptsLeft,
    this.message,
    this.lockedUntil,
    this.encryptedPin,
  });

  bool get isLocked => lockedUntil != null || (attemptsLeft != null && attemptsLeft == 0);
}
