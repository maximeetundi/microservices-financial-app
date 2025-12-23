import 'package:flutter/material.dart';
import '../../../../core/services/pin_service.dart';
import '../../../../core/services/biometric_service.dart';
import '../widgets/pin_input_widget.dart';

/// Dialog for verifying PIN before sensitive actions
class PinVerifyDialog extends StatefulWidget {
  final String? title;
  final String? subtitle;
  final bool allowBiometric;
  final VoidCallback? onVerified;
  final VoidCallback? onCancelled;

  const PinVerifyDialog({
    super.key,
    this.title,
    this.subtitle,
    this.allowBiometric = true,
    this.onVerified,
    this.onCancelled,
  });

  /// Shows the PIN verification dialog and returns true if verified
  static Future<bool> show(
    BuildContext context, {
    String? title,
    String? subtitle,
    bool allowBiometric = true,
  }) async {
    final result = await showModalBottomSheet<bool>(
      context: context,
      isScrollControlled: true,
      backgroundColor: Colors.transparent,
      builder: (context) => PinVerifyDialog(
        title: title,
        subtitle: subtitle,
        allowBiometric: allowBiometric,
      ),
    );
    return result ?? false;
  }

  @override
  State<PinVerifyDialog> createState() => _PinVerifyDialogState();
}

class _PinVerifyDialogState extends State<PinVerifyDialog> {
  final PinService _pinService = PinService();
  final BiometricService _biometricService = BiometricService();
  
  bool _isLoading = false;
  String? _errorMessage;
  int? _attemptsLeft;
  bool _isLocked = false;
  bool _biometricAvailable = false;

  @override
  void initState() {
    super.initState();
    _checkBiometric();
  }

  Future<void> _checkBiometric() async {
    if (widget.allowBiometric) {
      final available = await _biometricService.isBiometricAvailable();
      final enabled = await _biometricService.isBiometricEnabled();
      setState(() {
        _biometricAvailable = available && enabled;
      });
      
      // Auto-prompt biometric if available
      if (_biometricAvailable) {
        _authenticateWithBiometric();
      }
    }
  }

  Future<void> _authenticateWithBiometric() async {
    try {
      final authenticated = await _biometricService.authenticate(
        reason: 'Authentifiez-vous pour continuer',
      );
      
      if (authenticated && mounted) {
        widget.onVerified?.call();
        Navigator.of(context).pop(true);
      }
    } catch (e) {
      // Biometric failed, user can still use PIN
      debugPrint('Biometric authentication failed: $e');
    }
  }

  Future<void> _verifyPin(String pin) async {
    setState(() {
      _isLoading = true;
      _errorMessage = null;
    });

    final result = await _pinService.verifyPin(pin);

    setState(() {
      _isLoading = false;
    });

    if (result.valid) {
      widget.onVerified?.call();
      if (mounted) {
        Navigator.of(context).pop(true);
      }
    } else {
      setState(() {
        _errorMessage = result.message ?? 'PIN incorrect';
        _attemptsLeft = result.attemptsLeft;
        _isLocked = result.isLocked;
      });
    }
  }

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    
    return Container(
      decoration: BoxDecoration(
        color: theme.colorScheme.surface,
        borderRadius: const BorderRadius.vertical(top: Radius.circular(24)),
      ),
      padding: EdgeInsets.only(
        bottom: MediaQuery.of(context).viewInsets.bottom,
      ),
      child: SafeArea(
        child: Padding(
          padding: const EdgeInsets.all(24),
          child: Column(
            mainAxisSize: MainAxisSize.min,
            children: [
              // Handle bar
              Container(
                width: 40,
                height: 4,
                decoration: BoxDecoration(
                  color: theme.colorScheme.outline.withOpacity(0.3),
                  borderRadius: BorderRadius.circular(2),
                ),
              ),
              
              const SizedBox(height: 24),
              
              // Icon
              Container(
                width: 64,
                height: 64,
                decoration: BoxDecoration(
                  color: theme.colorScheme.primaryContainer,
                  shape: BoxShape.circle,
                ),
                child: Icon(
                  Icons.lock_outline,
                  size: 32,
                  color: theme.colorScheme.primary,
                ),
              ),
              
              const SizedBox(height: 16),
              
              // Title
              Text(
                widget.title ?? 'Vérification PIN',
                style: theme.textTheme.titleLarge?.copyWith(
                  fontWeight: FontWeight.bold,
                ),
              ),
              
              const SizedBox(height: 8),
              
              // Subtitle
              Text(
                widget.subtitle ?? 'Entrez votre code PIN pour continuer',
                textAlign: TextAlign.center,
                style: theme.textTheme.bodyMedium?.copyWith(
                  color: theme.colorScheme.onSurface.withOpacity(0.7),
                ),
              ),
              
              const SizedBox(height: 24),
              
              // Error/Lock message
              if (_isLocked)
                Container(
                  padding: const EdgeInsets.all(12),
                  margin: const EdgeInsets.only(bottom: 16),
                  decoration: BoxDecoration(
                    color: Colors.amber.withOpacity(0.2),
                    borderRadius: BorderRadius.circular(12),
                  ),
                  child: Row(
                    children: [
                      const Icon(Icons.timer, color: Colors.amber),
                      const SizedBox(width: 8),
                      const Expanded(
                        child: Text(
                          'PIN temporairement bloqué. Réessayez dans quelques minutes.',
                          style: TextStyle(color: Colors.amber),
                        ),
                      ),
                    ],
                  ),
                )
              else if (_errorMessage != null)
                Container(
                  padding: const EdgeInsets.all(12),
                  margin: const EdgeInsets.only(bottom: 16),
                  decoration: BoxDecoration(
                    color: theme.colorScheme.errorContainer,
                    borderRadius: BorderRadius.circular(12),
                  ),
                  child: Row(
                    children: [
                      Icon(Icons.error_outline, color: theme.colorScheme.error),
                      const SizedBox(width: 8),
                      Expanded(
                        child: Column(
                          crossAxisAlignment: CrossAxisAlignment.start,
                          children: [
                            Text(
                              _errorMessage!,
                              style: TextStyle(color: theme.colorScheme.error),
                            ),
                            if (_attemptsLeft != null)
                              Text(
                                '$_attemptsLeft tentative(s) restante(s)',
                                style: TextStyle(
                                  color: theme.colorScheme.error,
                                  fontSize: 12,
                                ),
                              ),
                          ],
                        ),
                      ),
                    ],
                  ),
                ),
              
              // PIN Input
              if (_isLoading)
                const Padding(
                  padding: EdgeInsets.all(24),
                  child: CircularProgressIndicator(),
                )
              else if (!_isLocked)
                PinInputWidget(
                  onCompleted: _verifyPin,
                  autofocus: !_biometricAvailable,
                  enabled: !_isLocked,
                ),
              
              const SizedBox(height: 24),
              
              // Biometric button
              if (_biometricAvailable && !_isLocked)
                TextButton.icon(
                  onPressed: _authenticateWithBiometric,
                  icon: const Icon(Icons.fingerprint),
                  label: const Text('Utiliser l\'empreinte digitale'),
                ),
              
              // Cancel button
              TextButton(
                onPressed: () {
                  widget.onCancelled?.call();
                  Navigator.of(context).pop(false);
                },
                child: const Text('Annuler'),
              ),
            ],
          ),
        ),
      ),
    );
  }
}
