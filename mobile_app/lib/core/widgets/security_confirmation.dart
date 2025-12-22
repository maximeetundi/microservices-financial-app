import 'package:flutter/material.dart';
import 'package:flutter/services.dart';

import '../services/biometric_service.dart';

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
    String pin = '';
    bool hasError = false;
    bool loading = false;

    return await showModalBottomSheet<bool>(
          context: context,
          isScrollControlled: true,
          backgroundColor: Colors.transparent,
          builder: (context) => StatefulBuilder(
            builder: (context, setState) => Container(
              decoration: const BoxDecoration(
                color: Color(0xFF1a1a2e),
                borderRadius: BorderRadius.vertical(top: Radius.circular(24)),
              ),
              padding: EdgeInsets.only(
                bottom: MediaQuery.of(context).viewInsets.bottom,
              ),
              child: SafeArea(
                child: Column(
                  mainAxisSize: MainAxisSize.min,
                  children: [
                    const SizedBox(height: 16),
                    Container(
                      width: 40,
                      height: 4,
                      decoration: BoxDecoration(
                        color: Colors.white24,
                        borderRadius: BorderRadius.circular(2),
                      ),
                    ),
                    const SizedBox(height: 24),

                    // Lock icon
                    Container(
                      padding: const EdgeInsets.all(16),
                      decoration: BoxDecoration(
                        color: const Color(0xFF667eea).withOpacity(0.2),
                        shape: BoxShape.circle,
                      ),
                      child: const Icon(
                        Icons.lock_rounded,
                        size: 32,
                        color: Color(0xFF667eea),
                      ),
                    ),
                    const SizedBox(height: 16),

                    Text(
                      title ?? 'Confirmation requise',
                      style: const TextStyle(
                        color: Colors.white,
                        fontSize: 20,
                        fontWeight: FontWeight.bold,
                      ),
                    ),
                    const SizedBox(height: 8),
                    Text(
                      message ?? 'Entrez votre code PIN',
                      style: const TextStyle(
                        color: Color(0xFF94A3B8),
                        fontSize: 14,
                      ),
                    ),
                    const SizedBox(height: 24),

                    // PIN dots
                    Row(
                      mainAxisAlignment: MainAxisAlignment.center,
                      children: List.generate(6, (index) {
                        final isFilled = index < pin.length;
                        return Container(
                          width: 16,
                          height: 16,
                          margin: const EdgeInsets.symmetric(horizontal: 6),
                          decoration: BoxDecoration(
                            shape: BoxShape.circle,
                            color: hasError
                                ? const Color(0xFFEF4444)
                                : isFilled
                                    ? const Color(0xFF667eea)
                                    : Colors.transparent,
                            border: Border.all(
                              color: hasError
                                  ? const Color(0xFFEF4444)
                                  : const Color(0xFF667eea),
                              width: 2,
                            ),
                          ),
                        );
                      }),
                    ),

                    if (hasError) ...[
                      const SizedBox(height: 12),
                      const Text(
                        'Code incorrect',
                        style: TextStyle(color: Color(0xFFEF4444), fontSize: 14),
                      ),
                    ],

                    const SizedBox(height: 24),

                    // Number pad
                    _buildCompactNumberPad(
                      onNumber: (number) async {
                        if (pin.length < 6) {
                          HapticFeedback.selectionClick();
                          setState(() {
                            pin += number;
                            hasError = false;
                          });

                          if (pin.length == 6) {
                            setState(() => loading = true);
                            final valid = await _biometricService.verifyPin(pin);
                            if (valid) {
                              Navigator.of(context).pop(true);
                            } else {
                              HapticFeedback.vibrate();
                              setState(() {
                                hasError = true;
                                pin = '';
                                loading = false;
                              });
                            }
                          }
                        }
                      },
                      onBackspace: () {
                        if (pin.isNotEmpty) {
                          HapticFeedback.selectionClick();
                          setState(() {
                            pin = pin.substring(0, pin.length - 1);
                            hasError = false;
                          });
                        }
                      },
                      onBiometric: () async {
                        final success = await _biometricService.authenticateWithBiometrics();
                        if (success) {
                          Navigator.of(context).pop(true);
                        }
                      },
                      showBiometric: true,
                    ),

                    const SizedBox(height: 16),

                    TextButton(
                      onPressed: () => Navigator.of(context).pop(false),
                      child: const Text(
                        'Annuler',
                        style: TextStyle(color: Color(0xFF94A3B8)),
                      ),
                    ),
                    const SizedBox(height: 16),
                  ],
                ),
              ),
            ),
          ),
        ) ??
        false;
  }

  static Widget _buildCompactNumberPad({
    required Function(String) onNumber,
    required VoidCallback onBackspace,
    required VoidCallback onBiometric,
    required bool showBiometric,
  }) {
    return Column(
      children: [
        for (var row in [
          ['1', '2', '3'],
          ['4', '5', '6'],
          ['7', '8', '9'],
          ['bio', '0', 'del'],
        ])
          Padding(
            padding: const EdgeInsets.symmetric(vertical: 6),
            child: Row(
              mainAxisAlignment: MainAxisAlignment.center,
              children: row.map((key) {
                if (key == 'bio') {
                  return _buildPadButton(
                    child: showBiometric
                        ? const Icon(Icons.fingerprint, color: Color(0xFF667eea), size: 28)
                        : const SizedBox(width: 60),
                    onTap: showBiometric ? onBiometric : null,
                  );
                } else if (key == 'del') {
                  return _buildPadButton(
                    child: const Icon(Icons.backspace_outlined, color: Colors.white, size: 24),
                    onTap: onBackspace,
                  );
                } else {
                  return _buildPadButton(
                    child: Text(
                      key,
                      style: const TextStyle(
                        color: Colors.white,
                        fontSize: 28,
                        fontWeight: FontWeight.w500,
                      ),
                    ),
                    onTap: () => onNumber(key),
                  );
                }
              }).toList(),
            ),
          ),
      ],
    );
  }

  static Widget _buildPadButton({
    required Widget child,
    VoidCallback? onTap,
  }) {
    return GestureDetector(
      onTap: onTap,
      child: Container(
        width: 70,
        height: 60,
        margin: const EdgeInsets.symmetric(horizontal: 10),
        decoration: BoxDecoration(
          color: const Color(0xFF2d2d4a),
          borderRadius: BorderRadius.circular(16),
        ),
        child: Center(child: child),
      ),
    );
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
