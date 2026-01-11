import 'package:flutter/material.dart';
import 'dart:math';
import '../../../../core/services/pin_service.dart';
import '../../../../core/services/biometric_service.dart';

/// Dialog for verifying PIN before sensitive actions
/// Styled to match the Web "Premium Dark" design with Randomized Keypad
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
    final result = await showDialog<bool>(
      context: context,
      barrierDismissible: false,
      barrierColor: Colors.black.withOpacity(0.9), // Dark overlay like Web
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

class _PinVerifyDialogState extends State<PinVerifyDialog> with SingleTickerProviderStateMixin {
  final PinService _pinService = PinService();
  final BiometricService _biometricService = BiometricService();
  
  bool _isLoading = false;
  String? _errorMessage;
  int? _attemptsLeft;
  bool _isLocked = false;
  bool _biometricAvailable = false;
  
  // PIN State
  String _currentPin = '';
  final int _pinLength = 5;
  late List<int> _shuffledKeys;

  // Animation for error shake
  late AnimationController _shakeController;
  late Animation<double> _shakeAnimation;

  // Premium Dark Theme Colors
  static const Color _backgroundColor = Color(0xFF1A1A2E);
  static const Color _primaryColor = Color(0xFF6366F1); // Indigo
  static const Color _errorColor = Color(0xFFEF4444);
  static const Color _textColor = Colors.white;
  static const Color _subtitleColor = Color(0xFF9CA3AF);

  @override
  void initState() {
    super.initState();
    _shuffleKeys();
    _checkBiometric();
    
    _shakeController = AnimationController(
      duration: const Duration(milliseconds: 500),
      vsync: this,
    );
    _shakeAnimation = Tween<double>(begin: 0, end: 10)
        .chain(CurveTween(curve: Curves.elasticIn))
        .animate(_shakeController);
  }

  @override
  void dispose() {
    _shakeController.dispose();
    super.dispose();
  }

  void _shuffleKeys() {
    _shuffledKeys = [0, 1, 2, 3, 4, 5, 6, 7, 8, 9];
    _shuffledKeys.shuffle();
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
      debugPrint('Biometric authentication failed: $e');
    }
  }

  void _onKeyPressed(int key) {
    if (_currentPin.length < _pinLength && !_isLoading && !_isLocked) {
      setState(() {
        _currentPin += key.toString();
        _errorMessage = null;
      });

      if (_currentPin.length == _pinLength) {
        // Auto-verify
        _verifyPin(_currentPin);
      }
    }
  }

  void _onBackspace() {
    if (_currentPin.isNotEmpty && !_isLoading && !_isLocked) {
      setState(() {
        _currentPin = _currentPin.substring(0, _currentPin.length - 1);
        _errorMessage = null;
      });
    }
  }

  void _onClear() {
    setState(() {
      _currentPin = '';
      _errorMessage = null;
    });
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
      // Enhanced Security: Clear PIN and RESHUFFLE keys on error
      setState(() {
        _errorMessage = result.message ?? 'PIN incorrect';
        _attemptsLeft = result.attemptsLeft;
        _isLocked = result.isLocked;
        _currentPin = '';
        _shuffleKeys(); // Reshuffle to prevent pattern guessing
      });
      _shakeController.forward(from: 0);
    }
  }

  @override
  Widget build(BuildContext context) {
    return Theme(
      data: ThemeData.dark().copyWith(
        scaffoldBackgroundColor: _backgroundColor,
        colorScheme: const ColorScheme.dark(
          primary: _primaryColor,
          surface: _backgroundColor,
          onSurface: _textColor,
          error: _errorColor,
        ),
      ),
      child: AnimatedBuilder(
        animation: _shakeController,
        builder: (context, child) {
          return Transform.translate(
            offset: Offset(_shakeAnimation.value * sin(_shakeController.value * 3 * pi), 0),
            child: child,
          );
        },
        child: Dialog(
          backgroundColor: _backgroundColor,
          shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(24)),
          elevation: 24,
          insetPadding: const EdgeInsets.symmetric(horizontal: 16),
          child: Container(
            padding: const EdgeInsets.all(24),
            constraints: const BoxConstraints(maxWidth: 400),
            child: Column(
              mainAxisSize: MainAxisSize.min,
              children: [
                // Icon
                const Text(
                  '🔐',
                  style: TextStyle(fontSize: 40),
                ),
                
                const SizedBox(height: 12),
                
                // Title
                Text(
                  widget.title ?? 'Vérification requise',
                  style: const TextStyle(
                    color: _textColor,
                    fontSize: 20,
                    fontWeight: FontWeight.bold,
                  ),
                  textAlign: TextAlign.center,
                ),
                
                const SizedBox(height: 4),
                
                // Subtitle with security badge
                Text(
                  widget.subtitle ?? 'Entrez votre PIN',
                  style: const TextStyle(
                    color: _subtitleColor,
                    fontSize: 13,
                  ),
                  textAlign: TextAlign.center,
                ),
                const SizedBox(height: 4),
                Container(
                  padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 4),
                  decoration: BoxDecoration(
                    color: _primaryColor.withOpacity(0.1),
                    borderRadius: BorderRadius.circular(8),
                    border: Border.all(color: _primaryColor.withOpacity(0.2)),
                  ),
                  child: const Row(
                    mainAxisSize: MainAxisSize.min,
                    children: [
                      Icon(Icons.shield_outlined, size: 12, color: _primaryColor),
                      SizedBox(width: 4),
                      Text(
                        'Clavier sécurisé aléatoire',
                        style: TextStyle(color: _primaryColor, fontSize: 10, fontWeight: FontWeight.bold),
                      ),
                    ],
                  ),
                ),
                
                const SizedBox(height: 20),
                
                // PIN Display (Dots)
                Row(
                  mainAxisAlignment: MainAxisAlignment.center,
                  children: List.generate(_pinLength, (index) {
                    final isFilled = index < _currentPin.length;
                    return Container(
                      width: 45,
                      height: 55,
                      margin: const EdgeInsets.symmetric(horizontal: 4),
                      decoration: BoxDecoration(
                        color: isFilled ? _primaryColor.withOpacity(0.1) : Colors.white.withOpacity(0.05),
                        border: Border.all(
                          color: isFilled ? _primaryColor : _primaryColor.withOpacity(0.3),
                          width: 2,
                        ),
                        borderRadius: BorderRadius.circular(12),
                      ),
                      alignment: Alignment.center,
                      child: isFilled
                          ? const Icon(Icons.circle, size: 12, color: Colors.white)
                          : null,
                    );
                  }),
                ),

                const SizedBox(height: 20),

                // Error Message
                if (_errorMessage != null || _isLocked) ...[
                  Text(
                    _isLocked ? 'PIN bloqué. Attendez...' : _errorMessage!,
                    style: TextStyle(
                      color: _isLocked ? Colors.amber : _errorColor,
                      fontWeight: FontWeight.bold,
                      fontSize: 13,
                    ),
                    textAlign: TextAlign.center,
                  ),
                  const SizedBox(height: 12),
                ],
                
                // Randomized Keypad Grid
                if (!_isLocked)
                  Container(
                    width: 280,
                    child: GridView.builder(
                      shrinkWrap: true,
                      physics: const NeverScrollableScrollPhysics(),
                      gridDelegate: const SliverGridDelegateWithFixedCrossAxisCount(
                        crossAxisCount: 3,
                        mainAxisSpacing: 12,
                        crossAxisSpacing: 12,
                        childAspectRatio: 1.3,
                      ),
                      itemCount: 12,
                      itemBuilder: (context, index) {
                        // Layout:
                        // 0 1 2
                        // 3 4 5
                        // 6 7 8
                        // 9(C) 10(Digit) 11(Back)
                        
                        if (index == 9) {
                          // Clear Button
                          return _buildKeypadButton(
                            child: const Text('C', style: TextStyle(color: _errorColor, fontSize: 18, fontWeight: FontWeight.bold)),
                            onPressed: _onClear,
                            color: _errorColor.withOpacity(0.1),
                          );
                        } else if (index == 11) {
                          // Backspace Button
                          return _buildKeypadButton(
                            child: const Icon(Icons.backspace_outlined, size: 20),
                            onPressed: _onBackspace,
                          );
                        } else {
                          // Digit Button
                          // 0..8 map to _shuffledKeys[0..8]
                          // 10 maps to _shuffledKeys[9]
                          final keyIndex = index < 9 ? index : 9;
                          final digit = _shuffledKeys[keyIndex];
                          
                          return _buildKeypadButton(
                            child: Text(
                              digit.toString(),
                              style: const TextStyle(fontSize: 22, fontWeight: FontWeight.bold),
                            ),
                            onPressed: () => _onKeyPressed(digit),
                          );
                        }
                      },
                    ),
                  ),

                const SizedBox(height: 20),

                // Unlock Button (Visual only, since auto-submit)
                SizedBox(
                  width: double.infinity,
                  child: ElevatedButton(
                    onPressed: (_isLoading || _isLocked || _currentPin.length < 5) ? null : () => _verifyPin(_currentPin),
                    style: ElevatedButton.styleFrom(
                      backgroundColor: _primaryColor,
                      foregroundColor: Colors.white,
                      padding: const EdgeInsets.symmetric(vertical: 14),
                      shape: RoundedRectangleBorder(
                        borderRadius: BorderRadius.circular(12),
                      ),
                      elevation: 0,
                    ),
                    child: _isLoading 
                      ? const SizedBox(width: 20, height: 20, child: CircularProgressIndicator(color: Colors.white, strokeWidth: 2))
                      : const Text('Déverrouiller', style: TextStyle(fontSize: 16, fontWeight: FontWeight.bold)),
                  ),
                ),

                if (_biometricAvailable && !_isLocked)
                  Padding(
                    padding: const EdgeInsets.only(top: 8),
                    child: TextButton.icon(
                      onPressed: _authenticateWithBiometric,
                      icon: const Icon(Icons.fingerprint, color: _subtitleColor, size: 20),
                      label: const Text('Empreinte', style: TextStyle(color: _subtitleColor)),
                    ),
                  ),

                TextButton(
                  onPressed: () {
                    widget.onCancelled?.call();
                    Navigator.of(context).pop(false);
                  },
                  child: const Text('Retour', style: TextStyle(color: _subtitleColor)),
                ),
              ],
            ),
          ),
        ),
      ),
    );
  }

  Widget _buildKeypadButton({required Widget child, required VoidCallback onPressed, Color? color}) {
    return Material(
      color: Colors.transparent,
      child: InkWell(
        onTap: onPressed,
        borderRadius: BorderRadius.circular(12),
        child: Container(
          decoration: BoxDecoration(
            color: color ?? Colors.white.withOpacity(0.05),
            borderRadius: BorderRadius.circular(12),
            border: Border.all(color: Colors.white.withOpacity(0.1)),
          ),
          alignment: Alignment.center,
          child: child,
        ),
      ),
    );
  }
}

