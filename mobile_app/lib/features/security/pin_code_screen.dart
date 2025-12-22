import 'package:flutter/material.dart';
import 'package:flutter/services.dart';

import '../../core/services/biometric_service.dart';

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
  
  String _pin = '';
  String _confirmPin = '';
  bool _isConfirming = false;
  bool _hasError = false;
  String _errorMessage = '';
  bool _biometricAvailable = false;
  String _biometricName = '';
  
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
      if (_confirmPin.length < 6) {
        setState(() {
          _confirmPin += number;
          _hasError = false;
        });
        
        if (_confirmPin.length == 6) {
          _verifyConfirmPin();
        }
      }
    } else {
      if (_pin.length < 6) {
        setState(() {
          _pin += number;
          _hasError = false;
        });
        
        if (_pin.length == 6) {
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
      // Save the PIN
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
      _showError('Les codes ne correspondent pas');
      setState(() {
        _confirmPin = '';
      });
    }
  }

  Future<void> _verifyPin() async {
    final isValid = await _biometricService.verifyPin(_pin);
    
    if (isValid) {
      HapticFeedback.heavyImpact();
      widget.onSuccess?.call();
    } else {
      _showError('Code PIN incorrect');
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
        child: Column(
          children: [
            const Spacer(),
            
            // Lock icon
            Container(
              padding: const EdgeInsets.all(24),
              decoration: BoxDecoration(
                color: const Color(0xFF667eea).withOpacity(0.2),
                shape: BoxShape.circle,
              ),
              child: const Icon(
                Icons.lock_rounded,
                size: 48,
                color: Color(0xFF667eea),
              ),
            ),
            const SizedBox(height: 32),
            
            // Title
            Text(
              widget.title ?? (widget.isSetup
                  ? (_isConfirming ? 'Confirmez votre code' : 'Créez un code PIN')
                  : 'Entrez votre code PIN'),
              style: const TextStyle(
                color: Colors.white,
                fontSize: 24,
                fontWeight: FontWeight.bold,
              ),
            ),
            const SizedBox(height: 8),
            Text(
              widget.isSetup
                  ? (_isConfirming ? 'Entrez le même code' : '6 chiffres pour sécuriser l\'app')
                  : 'Pour accéder à votre compte',
              style: const TextStyle(
                color: Color(0xFF94A3B8),
                fontSize: 14,
              ),
            ),
            const SizedBox(height: 40),
            
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
                children: List.generate(6, (index) {
                  final currentPin = _isConfirming ? _confirmPin : _pin;
                  final isFilled = index < currentPin.length;
                  
                  return Container(
                    width: 20,
                    height: 20,
                    margin: const EdgeInsets.symmetric(horizontal: 8),
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
            _buildNumberPad(),
            
            const SizedBox(height: 24),
          ],
        ),
      ),
    );
  }

  Widget _buildNumberPad() {
    return Column(
      children: [
        Row(
          mainAxisAlignment: MainAxisAlignment.spaceEvenly,
          children: [
            _buildNumberButton('1'),
            _buildNumberButton('2'),
            _buildNumberButton('3'),
          ],
        ),
        const SizedBox(height: 16),
        Row(
          mainAxisAlignment: MainAxisAlignment.spaceEvenly,
          children: [
            _buildNumberButton('4'),
            _buildNumberButton('5'),
            _buildNumberButton('6'),
          ],
        ),
        const SizedBox(height: 16),
        Row(
          mainAxisAlignment: MainAxisAlignment.spaceEvenly,
          children: [
            _buildNumberButton('7'),
            _buildNumberButton('8'),
            _buildNumberButton('9'),
          ],
        ),
        const SizedBox(height: 16),
        Row(
          mainAxisAlignment: MainAxisAlignment.spaceEvenly,
          children: [
            // Biometric button
            _biometricAvailable
                ? _buildBiometricButton()
                : const SizedBox(width: 80),
            _buildNumberButton('0'),
            _buildBackspaceButton(),
          ],
        ),
      ],
    );
  }

  Widget _buildNumberButton(String number) {
    return GestureDetector(
      onTap: () => _onNumberPressed(number),
      child: Container(
        width: 80,
        height: 80,
        decoration: BoxDecoration(
          color: const Color(0xFF2d2d4a),
          shape: BoxShape.circle,
        ),
        child: Center(
          child: Text(
            number,
            style: const TextStyle(
              color: Colors.white,
              fontSize: 32,
              fontWeight: FontWeight.w500,
            ),
          ),
        ),
      ),
    );
  }

  Widget _buildBackspaceButton() {
    return GestureDetector(
      onTap: _onBackspace,
      child: Container(
        width: 80,
        height: 80,
        decoration: const BoxDecoration(
          shape: BoxShape.circle,
        ),
        child: const Center(
          child: Icon(
            Icons.backspace_outlined,
            color: Colors.white,
            size: 28,
          ),
        ),
      ),
    );
  }

  Widget _buildBiometricButton() {
    return GestureDetector(
      onTap: _authenticateWithBiometric,
      child: Container(
        width: 80,
        height: 80,
        decoration: const BoxDecoration(
          shape: BoxShape.circle,
        ),
        child: const Center(
          child: Icon(
            Icons.fingerprint,
            color: Color(0xFF667eea),
            size: 40,
          ),
        ),
      ),
    );
  }
}
