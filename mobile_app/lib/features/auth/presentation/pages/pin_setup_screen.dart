import 'package:flutter/material.dart';
import '../../../../core/services/pin_service.dart';
import '../../../../core/services/biometric_service.dart';
import '../widgets/pin_input_widget.dart';

/// Screen for setting up the 5-digit PIN (shown after registration)
class PinSetupScreen extends StatefulWidget {
  final VoidCallback? onPinSet;
  final bool canSkip;

  const PinSetupScreen({
    super.key,
    this.onPinSet,
    this.canSkip = false,
  });

  @override
  State<PinSetupScreen> createState() => _PinSetupScreenState();
}

class _PinSetupScreenState extends State<PinSetupScreen> {
  final PinService _pinService = PinService();
  final GlobalKey<FormState> _formKey = GlobalKey<FormState>();
  
  String _pin = '';
  String _confirmPin = '';
  bool _isConfirmStep = false;
  bool _isLoading = false;
  String? _errorMessage;

  void _onPinCompleted(String pin) {
    if (!_isConfirmStep) {
      setState(() {
        _pin = pin;
        _isConfirmStep = true;
        _errorMessage = null;
      });
    } else {
      setState(() {
        _confirmPin = pin;
      });
      _submitPin();
    }
  }

  Future<void> _submitPin() async {
    if (_pin.isEmpty || _confirmPin.isEmpty) return;

    setState(() {
      _isLoading = true;
      _errorMessage = null;
    });

    final result = await _pinService.setupPin(_pin, _confirmPin);

    setState(() {
      _isLoading = false;
    });

    if (result.success) {
      // Show success and proceed
      if (mounted) {
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(
            content: Text(result.message),
            backgroundColor: Colors.green,
          ),
        );
        widget.onPinSet?.call();
        Navigator.of(context).pop(true);
      }
    } else {
      setState(() {
        _errorMessage = result.message;
        _pin = '';
        _confirmPin = '';
        _isConfirmStep = false;
      });
    }
  }

  void _reset() {
    setState(() {
      _pin = '';
      _confirmPin = '';
      _isConfirmStep = false;
      _errorMessage = null;
    });
  }

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    
    return Scaffold(
      appBar: AppBar(
        title: const Text('Définir votre PIN'),
        automaticallyImplyLeading: widget.canSkip,
      ),
      body: SafeArea(
        child: Padding(
          padding: const EdgeInsets.all(24),
          child: Column(
            mainAxisAlignment: MainAxisAlignment.center,
            children: [
              // Icon
              Container(
                width: 80,
                height: 80,
                decoration: BoxDecoration(
                  color: theme.colorScheme.primaryContainer,
                  shape: BoxShape.circle,
                ),
                child: Icon(
                  Icons.lock_outline,
                  size: 40,
                  color: theme.colorScheme.primary,
                ),
              ),
              
              const SizedBox(height: 24),
              
              // Title
              Text(
                _isConfirmStep ? 'Confirmer le PIN' : 'Créer votre PIN',
                style: theme.textTheme.headlineSmall?.copyWith(
                  fontWeight: FontWeight.bold,
                ),
              ),
              
              const SizedBox(height: 8),
              
              // Subtitle
              Text(
                _isConfirmStep 
                    ? 'Entrez à nouveau votre code PIN'
                    : 'Créez un code PIN à 5 chiffres pour sécuriser vos transactions',
                textAlign: TextAlign.center,
                style: theme.textTheme.bodyMedium?.copyWith(
                  color: theme.colorScheme.onSurface.withOpacity(0.7),
                ),
              ),
              
              const SizedBox(height: 32),
              
              // Error message
              if (_errorMessage != null)
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
                        child: Text(
                          _errorMessage!,
                          style: TextStyle(color: theme.colorScheme.error),
                        ),
                      ),
                    ],
                  ),
                ),
              
              // PIN Input
              if (_isLoading)
                const CircularProgressIndicator()
              else
                PinInputWidget(
                  key: ValueKey(_isConfirmStep),
                  onCompleted: _onPinCompleted,
                  autofocus: true,
                ),
              
              const SizedBox(height: 24),
              
              // Reset button
              if (_isConfirmStep)
                TextButton(
                  onPressed: _reset,
                  child: const Text('Modifier le PIN'),
                ),
              
              const Spacer(),
              
              // Info text
              Text(
                'Ce PIN sera requis pour toutes les transactions sensibles',
                textAlign: TextAlign.center,
                style: theme.textTheme.bodySmall?.copyWith(
                  color: theme.colorScheme.onSurface.withOpacity(0.5),
                ),
              ),
            ],
          ),
        ),
      ),
    );
  }
}
