import 'dart:io';
import 'package:flutter/foundation.dart';
import 'package:flutter/services.dart';
import 'package:local_auth/local_auth.dart';
import 'package:flutter_secure_storage/flutter_secure_storage.dart';

/// Service for biometric and PIN authentication
class BiometricService {
  static final BiometricService _instance = BiometricService._internal();
  factory BiometricService() => _instance;
  BiometricService._internal();

  final LocalAuthentication _localAuth = LocalAuthentication();
  final FlutterSecureStorage _storage = const FlutterSecureStorage();

  // Storage keys
  static const String _biometricEnabledKey = 'biometric_enabled';
  static const String _pinCodeKey = 'pin_code';
  static const String _pinEnabledKey = 'pin_enabled';
  static const String _lockTimeoutKey = 'lock_timeout';

  /// Check if biometrics are available on this device
  Future<bool> isBiometricAvailable() async {
    try {
      final canCheck = await _localAuth.canCheckBiometrics;
      final isDeviceSupported = await _localAuth.isDeviceSupported();
      return canCheck && isDeviceSupported;
    } on PlatformException {
      return false;
    }
  }

  /// Get list of available biometric types
  Future<List<BiometricType>> getAvailableBiometrics() async {
    try {
      return await _localAuth.getAvailableBiometrics();
    } on PlatformException {
      return [];
    }
  }

  /// Check if fingerprint is available
  Future<bool> isFingerprintAvailable() async {
    final biometrics = await getAvailableBiometrics();
    return biometrics.contains(BiometricType.fingerprint);
  }

  /// Check if Face ID is available
  Future<bool> isFaceIdAvailable() async {
    final biometrics = await getAvailableBiometrics();
    return biometrics.contains(BiometricType.face);
  }

  /// Get user-friendly biometric type name
  Future<String> getBiometricTypeName() async {
    final biometrics = await getAvailableBiometrics();
    if (biometrics.contains(BiometricType.face)) {
      return Platform.isIOS ? 'Face ID' : 'Reconnaissance faciale';
    } else if (biometrics.contains(BiometricType.fingerprint)) {
      return 'Empreinte digitale';
    } else if (biometrics.contains(BiometricType.iris)) {
      return 'Scan de l\'iris';
    }
    return 'Biométrie';
  }

  /// Authenticate with biometrics
  Future<bool> authenticateWithBiometrics({String? reason}) async {
    try {
      return await _localAuth.authenticate(
        localizedReason: reason ?? 'Authentifiez-vous pour continuer',
        options: const AuthenticationOptions(
          stickyAuth: true,
          biometricOnly: true,
        ),
      );
    } on PlatformException catch (e) {
      debugPrint('Biometric auth error: ${e.message}');
      return false;
    }
  }

  /// Authenticate with any available method (biometric or PIN)
  Future<bool> authenticate({String? reason}) async {
    try {
      return await _localAuth.authenticate(
        localizedReason: reason ?? 'Authentifiez-vous pour continuer',
        options: const AuthenticationOptions(
          stickyAuth: true,
          biometricOnly: false,
        ),
      );
    } on PlatformException catch (e) {
      debugPrint('Auth error: ${e.message}');
      return false;
    }
  }

  /// Check if biometric login is enabled by user
  Future<bool> isBiometricEnabled() async {
    final value = await _storage.read(key: _biometricEnabledKey);
    return value == 'true';
  }

  /// Enable or disable biometric login
  Future<void> setBiometricEnabled(bool enabled) async {
    await _storage.write(key: _biometricEnabledKey, value: enabled.toString());
  }

  /// Check if PIN is enabled
  Future<bool> isPinEnabled() async {
    final value = await _storage.read(key: _pinEnabledKey);
    return value == 'true';
  }

  /// Set PIN code
  Future<void> setPin(String pin) async {
    // In production, hash the PIN before storing
    await _storage.write(key: _pinCodeKey, value: pin);
    await _storage.write(key: _pinEnabledKey, value: 'true');
  }

  /// Verify PIN code
  Future<bool> verifyPin(String pin) async {
    final storedPin = await _storage.read(key: _pinCodeKey);
    return storedPin == pin;
  }

  /// Remove PIN code
  Future<void> removePin() async {
    await _storage.delete(key: _pinCodeKey);
    await _storage.write(key: _pinEnabledKey, value: 'false');
  }

  /// Get lock timeout in minutes
  Future<int> getLockTimeout() async {
    final value = await _storage.read(key: _lockTimeoutKey);
    return int.tryParse(value ?? '5') ?? 5;
  }

  /// Set lock timeout in minutes
  Future<void> setLockTimeout(int minutes) async {
    await _storage.write(key: _lockTimeoutKey, value: minutes.toString());
  }

  /// Cancel any running authentication
  Future<void> cancelAuthentication() async {
    await _localAuth.stopAuthentication();
  }
}
