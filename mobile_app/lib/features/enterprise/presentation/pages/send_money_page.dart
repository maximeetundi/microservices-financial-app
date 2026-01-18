import 'package:flutter/material.dart';
import '../../../../core/services/api_service.dart';
import '../../data/models/enterprise_model.dart';
import '../widgets/pin_verification_dialog.dart';

class SendMoneyPage extends StatefulWidget {
  final Enterprise enterprise;
  final Map<String, dynamic> wallet;

  const SendMoneyPage({Key? key, required this.enterprise, required this.wallet}) : super(key: key);

  @override
  State<SendMoneyPage> createState() => _SendMoneyPageState();
}

class _SendMoneyPageState extends State<SendMoneyPage> {
  final ApiService _api = ApiService();
  final _formKey = GlobalKey<FormState>();
  
  bool _isLoading = false;
  bool _isLookingUp = false;
  Map<String, dynamic>? _recipient;
  String? _lookupError;
  
  final _amountController = TextEditingController();
  final _recipientController = TextEditingController();
  final _descriptionController = TextEditingController();

  String get _currency => widget.wallet['currency'] ?? 'XOF';
  double get _balance => (widget.wallet['balance'] ?? 0).toDouble();

  @override
  void dispose() {
    _amountController.dispose();
    _recipientController.dispose();
    _descriptionController.dispose();
    super.dispose();
  }

  Future<void> _lookupRecipient() async {
    if (_recipientController.text.length < 3) return;
    
    setState(() { _isLookingUp = true; _lookupError = null; _recipient = null; });
    
    try {
      final isEmail = _recipientController.text.contains('@');
      final query = isEmail 
          ? {'email': _recipientController.text}
          : {'phone': _recipientController.text};
      
      // Use user lookup API
      final response = await _api.user.lookup(query);
      setState(() => _recipient = response);
    } catch (e) {
      setState(() => _lookupError = 'Utilisateur introuvable');
    } finally {
      setState(() => _isLookingUp = false);
    }
  }

  void _confirmSend() async {
    if (!_formKey.currentState!.validate()) return;
    
    final amount = double.tryParse(_amountController.text) ?? 0;
    if (amount > _balance) {
      ScaffoldMessenger.of(context).showSnackBar(
        const SnackBar(content: Text('Solde insuffisant')),
      );
      return;
    }
    
    // Show PIN dialog
    final encryptedPin = await showDialog<String>(
      context: context,
      barrierDismissible: false,
      builder: (context) => const PinVerificationDialog(
        title: 'Confirmer le transfert',
        description: 'Entrez votre code PIN pour initier le processus d\'approbation.',
      ),
    );
    
    if (encryptedPin != null) {
      _executeSend(encryptedPin);
    }
  }

  Future<void> _executeSend(String encryptedPin) async {
    setState(() => _isLoading = true);
    
    try {
      final amount = double.tryParse(_amountController.text) ?? 0;
      
      final response = await _api.enterprise.initiateAction(widget.enterprise.id, {
        'action_type': 'TRANSACTION',
        'action_name': 'Transfert ${amount.toStringAsFixed(0)} $_currency',
        'description': _descriptionController.text.isEmpty 
            ? 'Transfert vers ${_recipientController.text}'
            : _descriptionController.text,
        'amount': amount,
        'currency': _currency,
        'payload': {
          'recipient_identifier': _recipientController.text,
          'amount': amount,
          'currency': _currency,
          'description': _descriptionController.text,
          'source_wallet_id': widget.wallet['id'],
          'initiator_pin': encryptedPin,
        },
      });
      
      ScaffoldMessenger.of(context).showSnackBar(
        const SnackBar(
          content: Text('Demande de transfert créée! En attente d\'approbation.'),
          backgroundColor: Colors.green,
        ),
      );
      
      Navigator.pop(context, true);
      
      // Navigate to approval page if exists
      if (response['approval'] != null) {
        final approvalId = response['approval']['id'] ?? response['approval']['_id'];
        if (approvalId != null) {
          // Could navigate to approval detail page
        }
      }
    } catch (e) {
      ScaffoldMessenger.of(context).showSnackBar(
        SnackBar(content: Text('Erreur: ${e.toString()}')),
      );
    } finally {
      setState(() => _isLoading = false);
    }
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: const Text('Envoyer de l\'argent'),
      ),
      body: Form(
        key: _formKey,
        child: SingleChildScrollView(
          padding: const EdgeInsets.all(16),
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              // Balance Card
              Container(
                width: double.infinity,
                padding: const EdgeInsets.all(20),
                decoration: BoxDecoration(
                  gradient: LinearGradient(
                    colors: [Colors.green.shade600, Colors.green.shade800],
                  ),
                  borderRadius: BorderRadius.circular(16),
                ),
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    const Text('Solde disponible', style: TextStyle(color: Colors.white70)),
                    const SizedBox(height: 4),
                    Text(
                      '${_balance.toStringAsFixed(0)} $_currency',
                      style: const TextStyle(
                        color: Colors.white,
                        fontSize: 28,
                        fontWeight: FontWeight.bold,
                      ),
                    ),
                  ],
                ),
              ),
              
              const SizedBox(height: 24),
              
              // Amount
              const Text('Montant à envoyer', style: TextStyle(fontWeight: FontWeight.bold)),
              const SizedBox(height: 8),
              Container(
                padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 8),
                decoration: BoxDecoration(
                  color: Colors.grey[100],
                  borderRadius: BorderRadius.circular(12),
                ),
                child: Row(
                  children: [
                    Expanded(
                      child: TextFormField(
                        controller: _amountController,
                        keyboardType: TextInputType.number,
                        style: const TextStyle(fontSize: 24, fontWeight: FontWeight.bold),
                        decoration: const InputDecoration(
                          border: InputBorder.none,
                          hintText: '0',
                        ),
                        validator: (v) {
                          if (v!.isEmpty) return 'Requis';
                          final amount = double.tryParse(v);
                          if (amount == null || amount <= 0) return 'Montant invalide';
                          return null;
                        },
                      ),
                    ),
                    Container(
                      padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 8),
                      decoration: BoxDecoration(
                        color: Colors.green.shade100,
                        borderRadius: BorderRadius.circular(8),
                      ),
                      child: Text(
                        _currency,
                        style: TextStyle(
                          color: Colors.green.shade700,
                          fontWeight: FontWeight.bold,
                          fontSize: 18,
                        ),
                      ),
                    ),
                  ],
                ),
              ),
              
              const SizedBox(height: 24),
              
              // Recipient
              const Text('Destinataire', style: TextStyle(fontWeight: FontWeight.bold)),
              const SizedBox(height: 8),
              Row(
                children: [
                  Expanded(
                    child: TextFormField(
                      controller: _recipientController,
                      decoration: InputDecoration(
                        hintText: 'Email ou téléphone',
                        border: OutlineInputBorder(borderRadius: BorderRadius.circular(12)),
                        prefixIcon: const Icon(Icons.person),
                      ),
                      validator: (v) => v!.isEmpty ? 'Requis' : null,
                      onFieldSubmitted: (_) => _lookupRecipient(),
                    ),
                  ),
                  const SizedBox(width: 8),
                  IconButton(
                    onPressed: _isLookingUp ? null : _lookupRecipient,
                    style: IconButton.styleFrom(
                      backgroundColor: Colors.blue,
                      foregroundColor: Colors.white,
                    ),
                    icon: _isLookingUp 
                        ? const SizedBox(width: 20, height: 20, child: CircularProgressIndicator(color: Colors.white, strokeWidth: 2))
                        : const Icon(Icons.search),
                  ),
                ],
              ),
              
              // Recipient Result
              if (_recipient != null) ...[
                const SizedBox(height: 12),
                Container(
                  padding: const EdgeInsets.all(12),
                  decoration: BoxDecoration(
                    color: Colors.green.shade50,
                    borderRadius: BorderRadius.circular(12),
                    border: Border.all(color: Colors.green.shade200),
                  ),
                  child: Row(
                    children: [
                      CircleAvatar(
                        backgroundColor: Colors.green,
                        child: Text(
                          '${_recipient!['first_name']?[0] ?? 'U'}',
                          style: const TextStyle(color: Colors.white),
                        ),
                      ),
                      const SizedBox(width: 12),
                      Expanded(
                        child: Column(
                          crossAxisAlignment: CrossAxisAlignment.start,
                          children: [
                            Text(
                              '${_recipient!['first_name'] ?? ''} ${_recipient!['last_name'] ?? ''}',
                              style: const TextStyle(fontWeight: FontWeight.bold),
                            ),
                            const Text('✓ Utilisateur vérifié', style: TextStyle(color: Colors.green, fontSize: 12)),
                          ],
                        ),
                      ),
                    ],
                  ),
                ),
              ],
              
              if (_lookupError != null) ...[
                const SizedBox(height: 8),
                Text(_lookupError!, style: const TextStyle(color: Colors.red, fontSize: 13)),
              ],
              
              const SizedBox(height: 24),
              
              // Description
              const Text('Note (optionnel)', style: TextStyle(fontWeight: FontWeight.bold)),
              const SizedBox(height: 8),
              TextFormField(
                controller: _descriptionController,
                decoration: InputDecoration(
                  hintText: 'Ex: Paiement fournisseur',
                  border: OutlineInputBorder(borderRadius: BorderRadius.circular(12)),
                ),
              ),
              
              const SizedBox(height: 24),
              
              // Summary
              Container(
                padding: const EdgeInsets.all(16),
                decoration: BoxDecoration(
                  color: Colors.grey[50],
                  borderRadius: BorderRadius.circular(12),
                  border: Border.all(color: Colors.grey[200]!),
                ),
                child: Column(
                  children: [
                    Row(
                      mainAxisAlignment: MainAxisAlignment.spaceBetween,
                      children: [
                        const Text('Montant', style: TextStyle(color: Colors.grey)),
                        Text(
                          '${_amountController.text.isEmpty ? '0' : _amountController.text} $_currency',
                          style: const TextStyle(fontWeight: FontWeight.w500),
                        ),
                      ],
                    ),
                    const Divider(height: 24),
                    Row(
                      mainAxisAlignment: MainAxisAlignment.spaceBetween,
                      children: [
                        const Text('Total à débiter', style: TextStyle(fontWeight: FontWeight.bold)),
                        Text(
                          '${_amountController.text.isEmpty ? '0' : _amountController.text} $_currency',
                          style: const TextStyle(fontWeight: FontWeight.bold, fontSize: 18),
                        ),
                      ],
                    ),
                  ],
                ),
              ),
              
              // Multi-admin notice
              const SizedBox(height: 16),
              Container(
                padding: const EdgeInsets.all(12),
                decoration: BoxDecoration(
                  color: Colors.blue.shade50,
                  borderRadius: BorderRadius.circular(12),
                ),
                child: Row(
                  children: [
                    Icon(Icons.info_outline, color: Colors.blue.shade700, size: 20),
                    const SizedBox(width: 12),
                    const Expanded(
                      child: Text(
                        'Cette transaction nécessitera l\'approbation d\'un autre administrateur.',
                        style: TextStyle(fontSize: 13),
                      ),
                    ),
                  ],
                ),
              ),
              
              const SizedBox(height: 24),
              
              // Submit
              SizedBox(
                width: double.infinity,
                child: ElevatedButton(
                  onPressed: _isLoading ? null : _confirmSend,
                  style: ElevatedButton.styleFrom(
                    backgroundColor: Colors.green,
                    foregroundColor: Colors.white,
                    padding: const EdgeInsets.symmetric(vertical: 16),
                    shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(12)),
                  ),
                  child: _isLoading
                      ? const SizedBox(
                          width: 20,
                          height: 20,
                          child: CircularProgressIndicator(color: Colors.white, strokeWidth: 2),
                        )
                      : const Text('Confirmer →', style: TextStyle(fontSize: 16)),
                ),
              ),
            ],
          ),
        ),
      ),
    );
  }
}
