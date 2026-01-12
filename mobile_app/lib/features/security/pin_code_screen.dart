import 'package:flutter/material.dart';
import 'package:flutter/services.dart';

import '../../core/services/biometric_service.dart';
import '../../core/services/pin_service.dart';

/// PIN Code Screen for authentication and PIN setup
class PinCodeScreen extends StatefulWidget {
  final bool isSetup; // true = setting new PIN, false = verifying PIN
  final String? title;
  final VoidCallback? onSuccess;
  final VoidCallback? onCancel;
  
  const PinCodeScreen({
    super.key,
    this.isSetup = false,
    this.title,
    this.onSuccess,
    this.onCancel,
  });

  @override
  State<PinCodeScreen> createState() => _PinCodeScreenState();
}

class _PinCodeScreenState extends State<PinCodeScreen> with SingleTickerProviderStateMixin {
  final BiometricService _biometricService = BiometricService();
  final PinService _pinService = PinService();
  
  String _pin = '';
  String _confirmPin = '';
  bool _isConfirming = false;
  bool _hasError = false;
  String _errorMessage = '';
  bool _biometricAvailable = false;
  String _biometricName = '';
  bool _isLoading = false;
  
  late AnimationController _shakeController;
  late Animation<double> _shakeAnimation;

  @override
  void initState() {
    super.initState();
    _shakeController = AnimationController(
      duration: const Duration(milliseconds: 500),
      vsync: this,
    );
    _shakeAnimation = Tween<double>(begin: 0, end: 24)
        .chain(CurveTween(curve: Curves.elasticIn))
        .animate(_shakeController);
    _shakeController.addStatusListener((status) {
      if (status == AnimationStatus.completed) {
        _shakeController.reset();
      }
    });
    
    _checkBiometric();
  }

  Future<void> _checkBiometric() async {
    final available = await _biometricService.isBiometricAvailable();
    final enabled = await _biometricService.isBiometricEnabled();
    final name = await _biometricService.getBiometricTypeName();
    
    if (mounted) {
      setState(() {
        _biometricAvailable = available && enabled && !widget.isSetup;
        _biometricName = name;
      });
      
      // Auto-trigger biometric on verify screen
      if (_biometricAvailable && !widget.isSetup) {
        _authenticateWithBiometric();
      }
    }
  }

  Future<void> _authenticateWithBiometric() async {
    final success = await _biometricService.authenticateWithBiometrics(
      reason: 'Authentifiez-vous pour accéder à l\'application',
    );
    
    if (success && mounted) {
      HapticFeedback.lightImpact();
      widget.onSuccess?.call();
    }
  }

  @override
  void dispose() {
    _shakeController.dispose();
    super.dispose();
  }

  void _onNumberPressed(String number) {
    HapticFeedback.selectionClick();
    
    if (_isConfirming) {
      if (_confirmPin.length < 5) {
        setState(() {
          _confirmPin += number;
          _hasError = false;
        });
        
        if (_confirmPin.length == 5) {
          _verifyConfirmPin();
        }
      }
    } else {
      if (_pin.length < 5) {
        setState(() {
          _pin += number;
          _hasError = false;
        });
        
        if (_pin.length == 5) {
          if (widget.isSetup) {
            _proceedToConfirm();
          } else {
            _verifyPin();
          }
        }
      }
    }
  }

  void _onBackspace() {
    HapticFeedback.selectionClick();
    
    if (_isConfirming) {
      if (_confirmPin.isNotEmpty) {
        setState(() {
          _confirmPin = _confirmPin.substring(0, _confirmPin.length - 1);
          _hasError = false;
        });
      }
    } else {
      if (_pin.isNotEmpty) {
        setState(() {
          _pin = _pin.substring(0, _pin.length - 1);
          _hasError = false;
        });
      }
    }
  }

  void _proceedToConfirm() {
    setState(() {
      _isConfirming = true;
    });
  }

  Future<void> _verifyConfirmPin() async {
    if (_pin == _confirmPin) {
      setState(() => _isLoading = true);
      
      // 1. Save to Backend First
      final result = await _pinService.setupPin(_pin, _confirmPin);
      
      setState(() => _isLoading = false);

      if (result.success) {
        // 2. Sync to Local Storage (for offline unlock if needed)
        await _biometricService.setPin(_pin);
        HapticFeedback.heavyImpact();
        
        if (mounted) {
          ScaffoldMessenger.of(context).showSnackBar(
            const SnackBar(
              content: Text('Code PIN configuré avec succès!'),
              backgroundColor: Color(0xFF10B981),
            ),
          );
          widget.onSuccess?.call();
        }
      } else {
        _showError(result.message);
        setState(() {
          _confirmPin = '';
        });
      }
    } else {
      _showError('Les codes ne correspondent pas');
      setState(() {
        _confirmPin = '';
      });
    }
  }

  Future<void> _verifyPin() async {
    // Set loading state
    setState(() {
      _isLoading = true;
      _hasError = false;
    });
    
    // Use PinService to verify via API (not local BiometricService)
    final result = await _pinService.verifyPin(_pin);
    
    if (!mounted) return;
    
    setState(() {
      _isLoading = false;
    });
    
    if (result.valid) {
      HapticFeedback.heavyImpact();
      widget.onSuccess?.call();
    } else {
      // Build error message
      String message;
      if (result.isLocked) {
        message = result.message ?? 'PIN temporairement bloqué';
      } else if (result.attemptsLeft != null && result.attemptsLeft! > 0) {
        message = 'PIN incorrect (${result.attemptsLeft} essais restants)';
      } else {
        message = result.message ?? 'Code PIN incorrect';
      }
      _showError(message);
      setState(() {
        _pin = '';
      });
    }
  }

  void _showError(String message) {
    HapticFeedback.vibrate();
    _shakeController.forward();
    setState(() {
      _hasError = true;
      _errorMessage = message;
    });
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: const Color(0xFF1a1a2e),
      appBar: AppBar(
        backgroundColor: Colors.transparent,
        elevation: 0,
        leading: widget.onCancel != null
            ? IconButton(
                icon: const Icon(Icons.close, color: Colors.white),
                onPressed: widget.onCancel,
              )
            : null,
      ),
      body: SafeArea(
        child: LayoutBuilder(
          builder: (context, constraints) {
            // Adjust sizes based on available height
            final isCompact = constraints.maxHeight < 600;
            final iconSize = isCompact ? 36.0 : 48.0;
            final iconPadding = isCompact ? 16.0 : 24.0;
            final titleSize = isCompact ? 20.0 : 24.0;
            final dotSize = isCompact ? 16.0 : 20.0;
            final buttonSize = isCompact ? 64.0 : 80.0;
            final buttonFontSize = isCompact ? 26.0 : 32.0;
            
            return SingleChildScrollView(
              child: ConstrainedBox(
                constraints: BoxConstraints(minHeight: constraints.maxHeight),
                child: IntrinsicHeight(
                  child: Column(
                    children: [
                      SizedBox(height: isCompact ? 16 : 32),
                      
                      // Lock icon or loader
                      Container(
                        padding: EdgeInsets.all(iconPadding),
                        decoration: BoxDecoration(
                          color: const Color(0xFF667eea).withOpacity(0.2),
                          shape: BoxShape.circle,
                        ),
                        child: _isLoading 
                          ? SizedBox(
                              width: iconSize,
                              height: iconSize,
                              child: const CircularProgressIndicator(
                                color: Color(0xFF667eea),
                                strokeWidth: 3,
                              ),
                            )
                          : Icon(
                              Icons.lock_rounded,
                              size: iconSize,
                              color: const Color(0xFF667eea),
                            ),
                      ),
                      SizedBox(height: isCompact ? 16 : 32),
                      
                      // Title
                      Text(
                        widget.title ?? (widget.isSetup
                            ? (_isConfirming ? 'Confirmez votre code' : 'Créez un code PIN')
                            : 'Entrez votre code PIN'),
                        style: TextStyle(
                          color: Colors.white,
                          fontSize: titleSize,
                          fontWeight: FontWeight.bold,
                        ),
                      ),
                      const SizedBox(height: 8),
                      Text(
                        widget.isSetup
                            ? (_isConfirming ? 'Entrez le même code' : '5 chiffres pour sécuriser l\'app')
                            : 'Pour accéder à votre compte',
                        style: const TextStyle(
                          color: Color(0xFF94A3B8),
                          fontSize: 14,
                        ),
                      ),
                      SizedBox(height: isCompact ? 24 : 40),
                      
                      // PIN dots
                      AnimatedBuilder(
                        animation: _shakeAnimation,
                        builder: (context, child) {
                          return Transform.translate(
                            offset: Offset(_shakeAnimation.value * ((_shakeController.value * 10).toInt() % 2 == 0 ? 1 : -1), 0),
                            child: child,
                          );
                        },
                        child: Row(
                          mainAxisAlignment: MainAxisAlignment.center,
                          children: List.generate(5, (index) {
                            final currentPin = _isConfirming ? _confirmPin : _pin;
                            final isFilled = index < currentPin.length;
                            
                            return Container(
                              width: dotSize,
                              height: dotSize,
                              margin: EdgeInsets.symmetric(horizontal: isCompact ? 6 : 8),
                              decoration: BoxDecoration(
                                shape: BoxShape.circle,
                                color: _hasError
                                    ? const Color(0xFFEF4444)
                                    : isFilled
                                        ? const Color(0xFF667eea)
                                        : const Color(0xFF64748B).withOpacity(0.3),
                                border: Border.all(
                                  color: _hasError
                                      ? const Color(0xFFEF4444)
                                      : const Color(0xFF667eea),
                                  width: 2,
                                ),
                              ),
                            );
                          }),
                        ),
                      ),
                      
                      // Error message
                      if (_hasError) ...[
                        const SizedBox(height: 16),
                        Text(
                          _errorMessage,
                          style: const TextStyle(
                            color: Color(0xFFEF4444),
                            fontSize: 14,
                          ),
                        ),
                      ],
                      
                      const Spacer(),
                      
                      // Number pad
                      _buildNumberPad(buttonSize, buttonFontSize),
                      
                      SizedBox(height: isCompact ? 16 : 24),
                    ],
                  ),
                ),
              ),
            );
          },
        ),
      ),
    );
  }

  Widget _buildNumberPad(double buttonSize, double fontSize) {
    return Column(
      children: [
        Row(
          mainAxisAlignment: MainAxisAlignment.spaceEvenly,
          children: [
            _buildNumberButton('1', buttonSize, fontSize),
            _buildNumberButton('2', buttonSize, fontSize),
            _buildNumberButton('3', buttonSize, fontSize),
          ],
        ),
        const SizedBox(height: 12),
        Row(
          mainAxisAlignment: MainAxisAlignment.spaceEvenly,
          children: [
            _buildNumberButton('4', buttonSize, fontSize),
            _buildNumberButton('5', buttonSize, fontSize),
            _buildNumberButton('6', buttonSize, fontSize),
          ],
        ),
        const SizedBox(height: 12),
        Row(
          mainAxisAlignment: MainAxisAlignment.spaceEvenly,
          children: [
            _buildNumberButton('7', buttonSize, fontSize),
            _buildNumberButton('8', buttonSize, fontSize),
            _buildNumberButton('9', buttonSize, fontSize),
          ],
        ),
        const SizedBox(height: 12),
        Row(
          mainAxisAlignment: MainAxisAlignment.spaceEvenly,
          children: [
            // Biometric button
            _biometricAvailable
                ? _buildBiometricButton(buttonSize)
                : SizedBox(width: buttonSize),
            _buildNumberButton('0', buttonSize, fontSize),
            _buildBackspaceButton(buttonSize),
          ],
        ),
      ],
    );
  }

  Widget _buildNumberButton(String number, double size, double fontSize) {
    return GestureDetector(
      onTap: () => _onNumberPressed(number),
      child: Container(
        width: size,
        height: size,
        decoration: const BoxDecoration(
          color: Color(0xFF2d2d4a),
          shape: BoxShape.circle,
        ),
        child: Center(
          child: Text(
            number,
            style: TextStyle(
              color: Colors.white,
              fontSize: fontSize,
              fontWeight: FontWeight.w500,
            ),
          ),
        ),
      ),
    );
  }

  Widget _buildBackspaceButton(double size) {
    return GestureDetector(
      onTap: _onBackspace,
      child: Container(
        width: size,
        height: size,
        decoration: const BoxDecoration(
          shape: BoxShape.circle,
        ),
        child: Center(
          child: Icon(
            Icons.backspace_outlined,
            color: Colors.white,
            size: size * 0.35,
          ),
        ),
      ),
    );
  }

  Widget _buildBiometricButton(double size) {
    return GestureDetector(
      onTap: _authenticateWithBiometric,
      child: Container(
        width: size,
        height: size,
        decoration: const BoxDecoration(
          shape: BoxShape.circle,
        ),
        child: Center(
          child: Icon(
            Icons.fingerprint,
            color: const Color(0xFF667eea),
            size: size * 0.5,
          ),
        ),
      ),
    );
  }
}
