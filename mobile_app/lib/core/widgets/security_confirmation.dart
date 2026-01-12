import 'package:flutter/material.dart';
import 'package:flutter/services.dart';

import '../services/biometric_service.dart';
import '../../features/auth/presentation/pages/pin_verify_dialog.dart';

/// Widget for requesting security confirmation (PIN or biometric) before sensitive actions
class SecurityConfirmation {
  static final BiometricService _biometricService = BiometricService();

  /// Show security confirmation and return true if authenticated
  static Future<bool> confirm(
    BuildContext context, {
    String? title,
    String? message,
    bool allowBiometric = true,
  }) async {
    final biometricAvailable = await _biometricService.isBiometricAvailable();
    final biometricEnabled = await _biometricService.isBiometricEnabled();
    final pinEnabled = await _biometricService.isPinEnabled();

    // If biometric is available and enabled, try biometric first
    if (allowBiometric && biometricAvailable && biometricEnabled) {
      final success = await _biometricService.authenticateWithBiometrics(
        reason: message ?? 'Authentifiez-vous pour confirmer cette action',
      );
      if (success) return true;
    }

    // If PIN is enabled, show PIN dialog
    if (pinEnabled) {
      return await _showPinDialog(context, title: title, message: message);
    }

    // If no security is configured, show warning and allow
    if (!biometricEnabled && !pinEnabled) {
      return await _showNoSecurityWarning(context);
    }

    return false;
  }

  /// Require security confirmation - shows error if fails
  static Future<bool> require(
    BuildContext context, {
    String? title,
    String? message,
  }) async {
    final result = await confirm(
      context,
      title: title,
      message: message,
    );

    if (!result) {
      ScaffoldMessenger.of(context).showSnackBar(
        const SnackBar(
          content: Text('Authentification requise pour cette action'),
          backgroundColor: Colors.red,
        ),
      );
    }

    return result;
  }

  static Future<bool> _showPinDialog(
    BuildContext context, {
    String? title,
    String? message,
  }) async {
    final result = await PinVerifyDialog.show(
      context,
      title: title ?? 'Confirmation de sécurité',
      subtitle: message ?? 'Entrez votre code PIN',
      returnEncryptedPin: false, // We just need boolean verification here
    );
    return result == true;
  }



  static Future<bool> _showNoSecurityWarning(BuildContext context) async {
    return await showDialog<bool>(
          context: context,
          builder: (context) => AlertDialog(
            backgroundColor: Colors.white,
            shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(20)),
            title: const Row(
              children: [
                Icon(Icons.warning_amber_rounded, color: Colors.orange),
                SizedBox(width: 8),
                Text('Sécurité non configurée'),
              ],
            ),
            content: const Text(
              'Vous n\'avez pas configuré de code PIN ou d\'empreinte. '
              'Nous vous recommandons de le faire dans Paramètres > Sécurité.',
            ),
            actions: [
              TextButton(
                onPressed: () => Navigator.of(context).pop(false),
                child: const Text('Annuler'),
              ),
              ElevatedButton(
                onPressed: () => Navigator.of(context).pop(true),
                style: ElevatedButton.styleFrom(
                  backgroundColor: const Color(0xFF667eea),
                ),
                child: const Text('Continuer quand même'),
              ),
            ],
          ),
        ) ??
        false;
  }
}

/// Mixin for screens that want to use security confirmation
mixin SecureActionMixin<T extends StatefulWidget> on State<T> {
  /// Confirm action with PIN/biometric before executing
  Future<void> secureAction(
    Future<void> Function() action, {
    String? title,
    String? message,
  }) async {
    final confirmed = await SecurityConfirmation.require(
      context,
      title: title,
      message: message,
    );

    if (confirmed) {
      await action();
    }
  }
}
