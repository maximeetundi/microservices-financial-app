import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import 'package:crypto/crypto.dart';
import 'dart:convert';

class PinVerificationDialog extends StatefulWidget {
  final String title;
  final String description;

  const PinVerificationDialog({
    Key? key,
    required this.title,
    required this.description,
  }) : super(key: key);

  @override
  State<PinVerificationDialog> createState() => _PinVerificationDialogState();
}

class _PinVerificationDialogState extends State<PinVerificationDialog> {
  final List<TextEditingController> _controllers = List.generate(5, (_) => TextEditingController());
  final List<FocusNode> _focusNodes = List.generate(5, (_) => FocusNode());
  bool _isVerifying = false;
  String? _error;

  @override
  void initState() {
    super.initState();
    // Focus first field after build
    WidgetsBinding.instance.addPostFrameCallback((_) {
      _focusNodes[0].requestFocus();
    });
  }

  @override
  void dispose() {
    for (var c in _controllers) {
      c.dispose();
    }
    for (var f in _focusNodes) {
      f.dispose();
    }
    super.dispose();
  }

  String get _pinValue => _controllers.map((c) => c.text).join();
  
  String _encryptPin(String pin) {
    final bytes = utf8.encode(pin);
    final hash = sha256.convert(bytes);
    return hash.toString();
  }

  void _onPinChange(int index, String value) {
    if (value.isNotEmpty && index < 4) {
      _focusNodes[index + 1].requestFocus();
    }
    
    // If all filled, verify
    if (_pinValue.length == 5) {
      _verify();
    }
    
    setState(() => _error = null);
  }

  void _onKeyDown(int index, RawKeyEvent event) {
    if (event is RawKeyDownEvent && 
        event.logicalKey == LogicalKeyboardKey.backspace &&
        _controllers[index].text.isEmpty &&
        index > 0) {
      _focusNodes[index - 1].requestFocus();
      _controllers[index - 1].clear();
    }
  }

  void _verify() async {
    if (_pinValue.length != 5) {
      setState(() => _error = 'Entrez votre code PIN à 5 chiffres');
      return;
    }

    setState(() { _isVerifying = true; _error = null; });
    
    // Simulate verification delay
    await Future.delayed(const Duration(milliseconds: 300));
    
    // Encrypt and return PIN
    final encryptedPin = _encryptPin(_pinValue);
    Navigator.pop(context, encryptedPin);
  }

  void _clear() {
    for (var c in _controllers) {
      c.clear();
    }
    _focusNodes[0].requestFocus();
    setState(() => _error = null);
  }

  @override
  Widget build(BuildContext context) {
    return Dialog(
      shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(20)),
      child: Padding(
        padding: const EdgeInsets.all(24),
        child: Column(
          mainAxisSize: MainAxisSize.min,
          children: [
            // Lock Icon
            Container(
              padding: const EdgeInsets.all(16),
              decoration: BoxDecoration(
                color: Colors.blue.shade50,
                shape: BoxShape.circle,
              ),
              child: Icon(Icons.lock, color: Colors.blue.shade700, size: 32),
            ),
            const SizedBox(height: 16),
            
            // Title
            Text(
              widget.title,
              style: const TextStyle(fontSize: 20, fontWeight: FontWeight.bold),
              textAlign: TextAlign.center,
            ),
            const SizedBox(height: 8),
            
            // Description
            Text(
              widget.description,
              style: TextStyle(color: Colors.grey[600], fontSize: 14),
              textAlign: TextAlign.center,
            ),
            const SizedBox(height: 24),
            
            // PIN Input Fields
            Row(
              mainAxisAlignment: MainAxisAlignment.center,
              children: List.generate(5, (index) => Container(
                width: 48,
                height: 56,
                margin: const EdgeInsets.symmetric(horizontal: 4),
                child: RawKeyboardListener(
                  focusNode: FocusNode(),
                  onKey: (event) => _onKeyDown(index, event),
                  child: TextField(
                    controller: _controllers[index],
                    focusNode: _focusNodes[index],
                    keyboardType: TextInputType.number,
                    textAlign: TextAlign.center,
                    maxLength: 1,
                    obscureText: true,
                    style: const TextStyle(fontSize: 24, fontWeight: FontWeight.bold),
                    decoration: InputDecoration(
                      counterText: '',
                      filled: true,
                      fillColor: Colors.grey[100],
                      border: OutlineInputBorder(
                        borderRadius: BorderRadius.circular(12),
                        borderSide: BorderSide.none,
                      ),
                      focusedBorder: OutlineInputBorder(
                        borderRadius: BorderRadius.circular(12),
                        borderSide: BorderSide(color: Colors.blue.shade500, width: 2),
                      ),
                    ),
                    inputFormatters: [FilteringTextInputFormatter.digitsOnly],
                    onChanged: (v) => _onPinChange(index, v),
                  ),
                ),
              )),
            ),
            
            // Error
            if (_error != null) ...[
              const SizedBox(height: 12),
              Text(_error!, style: const TextStyle(color: Colors.red, fontSize: 13)),
            ],
            
            const SizedBox(height: 24),
            
            // Buttons
            Row(
              children: [
                Expanded(
                  child: TextButton(
                    onPressed: () => Navigator.pop(context),
                    style: TextButton.styleFrom(
                      padding: const EdgeInsets.symmetric(vertical: 14),
                    ),
                    child: const Text('Annuler'),
                  ),
                ),
                const SizedBox(width: 12),
                Expanded(
                  child: ElevatedButton(
                    onPressed: _isVerifying ? null : _verify,
                    style: ElevatedButton.styleFrom(
                      backgroundColor: Colors.blue,
                      foregroundColor: Colors.white,
                      padding: const EdgeInsets.symmetric(vertical: 14),
                      shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(12)),
                    ),
                    child: _isVerifying
                        ? const SizedBox(width: 20, height: 20, child: CircularProgressIndicator(color: Colors.white, strokeWidth: 2))
                        : const Text('Confirmer'),
                  ),
                ),
              ],
            ),
            
            // Clear button
            TextButton(
              onPressed: _clear,
              child: const Text('Effacer', style: TextStyle(color: Colors.grey)),
            ),
          ],
        ),
      ),
    );
  }
}
