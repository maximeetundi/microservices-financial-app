import 'package:flutter/material.dart';
import '../../../../features/auth/presentation/pages/pin_verify_dialog.dart';
import '../../../../core/services/wallet_api_service.dart';
import '../../../../core/services/donation_api_service.dart';

class DonateModal extends StatefulWidget {
  final String campaignId;
  final String currency;
  final List<dynamic>? formSchema;
  
  // Advanced Options
  final String donationType; // free, fixed, tiers
  final double? fixedAmount;
  final List<dynamic>? donationTiers;
  final double? minAmount;
  final double? maxAmount;

  const DonateModal({
    super.key,
    required this.campaignId,
    required this.currency,
    this.formSchema,
    this.donationType = 'free',
    this.fixedAmount,
    this.donationTiers,
    this.minAmount,
    this.maxAmount,
  });

  @override
  State<DonateModal> createState() => _DonateModalState();
}

class _DonateModalState extends State<DonateModal> {
  final WalletApiService _walletApi = WalletApiService();
  final DonationApiService _donationApi = DonationApiService();
  final TextEditingController _amountController = TextEditingController();
  final TextEditingController _messageController = TextEditingController();

  List<dynamic> _wallets = [];
  String? _selectedWalletId;
  bool _loadingWallets = true;
  bool _submitting = false;
  
  String _frequency = 'one_time';
  bool _isAnonymous = false;
  final Map<String, dynamic> _formData = {};

  @override
  void initState() {
    super.initState();
    _loadWallets();
    
    // Init Amount
    if (widget.donationType == 'fixed' && widget.fixedAmount != null) {
      _amountController.text = widget.fixedAmount!.toStringAsFixed(0);
    } else if (widget.donationType == 'tiers' && widget.donationTiers != null && widget.donationTiers!.isNotEmpty) {
      // Default to first tier? Or let user pick. Let's let user pick.
    }
  }

  Future<void> _loadWallets() async {
    try {
      final wallets = await _walletApi.getWallets();
      setState(() {
        _wallets = wallets;
        if (_wallets.isNotEmpty) {
          _selectedWalletId = _wallets.first['id'];
        }
      });
    } catch (e) {
      debugPrint('Error loading wallets: $e');
    } finally {
      setState(() => _loadingWallets = false);
    }
  }

  Future<void> _donate() async {
    if (_amountController.text.isEmpty || _selectedWalletId == null) {
      ScaffoldMessenger.of(context).showSnackBar(const SnackBar(content: Text('Veuillez remplir tous les champs')));
      return;
    }

    final amount = double.tryParse(_amountController.text) ?? 0;
    if (amount <= 0) {
       ScaffoldMessenger.of(context).showSnackBar(const SnackBar(content: Text('Montant invalide')));
       return;
    }

    if (widget.donationType == 'free') {
       if (widget.minAmount != null && amount < widget.minAmount!) {
          ScaffoldMessenger.of(context).showSnackBar(SnackBar(content: Text('Le montant minimum est de ${widget.minAmount} ${widget.currency}')));
          return;
       }
       if (widget.maxAmount != null && amount > widget.maxAmount!) {
          ScaffoldMessenger.of(context).showSnackBar(SnackBar(content: Text('Le montant maximum est de ${widget.maxAmount} ${widget.currency}')));
          return;
       }
    }

    // Validate Dynamic Fields
    if (widget.formSchema != null) {
      for (var field in widget.formSchema!) {
        if (field['required'] == true) {
          final val = _formData[field['name']];
          if (val == null || val.toString().trim().isEmpty) {
             ScaffoldMessenger.of(context).showSnackBar(SnackBar(content: Text('Le champ "${field['label']}" est obligatoire')));
             return;
          }
        }
      }
    }

    // Secure PIN Verification
    final pin = await PinVerifyDialog.show(
      context, 
      returnRawPin: true,
      subtitle: 'Confirmez un don de ${_amountController.text} ${widget.currency}',
    );

    if (pin == null || pin is! String) return; // Cancelled or failed

    setState(() => _submitting = true);
    
    try {
      await _donationApi.initiateDonation(
        campaignId: widget.campaignId,
        amount: double.parse(_amountController.text),
        currency: widget.currency,
        walletId: _selectedWalletId!,
        pin: pin,
        message: _messageController.text,
        isAnonymous: _isAnonymous,
        frequency: _frequency,
        formData: _formData,
      );
      
      if (mounted) {
        Navigator.pop(context, true); // Return true on success
        ScaffoldMessenger.of(context).showSnackBar(const SnackBar(content: Text('Don effectué avec succès !')));
      }
    } catch (e) {
      if (mounted) {
        ScaffoldMessenger.of(context).showSnackBar(SnackBar(content: Text('Erreur: $e')));
      }
    } finally {
      if (mounted) setState(() => _submitting = false);
    }
  }

  @override
  Widget build(BuildContext context) {
    return Container(
      padding: EdgeInsets.only(
        bottom: MediaQuery.of(context).viewInsets.bottom,
        left: 20,
        right: 20,
        top: 20,
      ),
      decoration: const BoxDecoration(
        color: Color(0xFF1a1a2e),
        borderRadius: BorderRadius.vertical(top: Radius.circular(20)),
      ),
      child: SingleChildScrollView(
        child: Column(
          mainAxisSize: MainAxisSize.min,
          crossAxisAlignment: CrossAxisAlignment.stretch,
          children: [
            const Text(
              'Faire un don ❤️',
              textAlign: TextAlign.center,
              style: TextStyle(fontSize: 24, fontWeight: FontWeight.bold, color: Colors.white),
            ),
            const SizedBox(height: 20),
            
            // Amount
            // Amount Section
            if (widget.donationType == 'fixed') ...[
               Container(
                 width: double.infinity,
                 padding: const EdgeInsets.all(16),
                 decoration: BoxDecoration(
                   color: const Color(0xFF6366f1).withOpacity(0.1),
                   borderRadius: BorderRadius.circular(12),
                   border: Border.all(color: const Color(0xFF6366f1)),
                 ),
                 child: Column(
                   children: [
                     const Text('Montant Fixe', style: TextStyle(color: Colors.white70, fontSize: 12)),
                     Text(
                       '${widget.fixedAmount?.toStringAsFixed(0) ?? 0} ${widget.currency}',
                       style: const TextStyle(color: Color(0xFF6366f1), fontSize: 24, fontWeight: FontWeight.bold),
                     ),
                   ],
                 ),
               ),
            ] else if (widget.donationType == 'tiers' && widget.donationTiers != null) ...[
               Wrap(
                 spacing: 8,
                 runSpacing: 8,
                 children: widget.donationTiers!.map((tier) {
                   final amount = double.tryParse(tier['amount'].toString()) ?? 0;
                   final isSelected = _amountController.text == amount.toStringAsFixed(0);
                   return GestureDetector(
                     onTap: () {
                         setState(() {
                           _amountController.text = amount.toStringAsFixed(0);
                         });
                     },
                     child: Container(
                       padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 12),
                       decoration: BoxDecoration(
                         color: isSelected ? const Color(0xFF6366f1) : Colors.white.withOpacity(0.1),
                         borderRadius: BorderRadius.circular(12),
                         border: Border.all(color: isSelected ? const Color(0xFF6366f1) : Colors.transparent),
                       ),
                       child: Column(
                         children: [
                           Text(tier['label'], style: TextStyle(color: isSelected ? Colors.white : Colors.white70, fontWeight: FontWeight.bold)),
                           Text('${amount.toStringAsFixed(0)} ${widget.currency}', style: TextStyle(color: isSelected ? Colors.white : const Color(0xFF6366f1), fontSize: 12)),
                         ],
                       ),
                     ),
                   );
                 }).toList(),
               ),
            ] else ...[
                TextField(
                  controller: _amountController,
                  keyboardType: TextInputType.number,
                  style: const TextStyle(color: Colors.white),
                  decoration: InputDecoration(
                    labelText: 'Montant (${widget.currency})',
                    labelStyle: TextStyle(color: Colors.white.withOpacity(0.6)),
                    filled: true,
                    fillColor: Colors.white.withOpacity(0.1),
                    border: OutlineInputBorder(borderRadius: BorderRadius.circular(12), borderSide: BorderSide.none),
                    helperText: (widget.minAmount != null || widget.maxAmount != null) 
                        ? 'Min: ${widget.minAmount ?? 0} - Max: ${widget.maxAmount ?? "Illimité"}' 
                        : null,
                    helperStyle: const TextStyle(color: Colors.grey),
                  ),
                ),
            ],
            const SizedBox(height: 16),

            // Frequency
            const Text('Fréquence', style: TextStyle(color: Colors.white, fontWeight: FontWeight.bold)),
            const SizedBox(height: 8),
            Row(
              children: [
                _buildFrequencyOption('Une fois', 'one_time'),
                const SizedBox(width: 8),
                _buildFrequencyOption('Mensuel', 'monthly'),
                const SizedBox(width: 8),
                _buildFrequencyOption('Annuel', 'annually'),
              ],
            ),
            if (_frequency != 'one_time')
              Padding(
                padding: const EdgeInsets.only(top: 8),
                child: Text(
                  '⚠️ Vous recevrez une demande de paiement à valider manuellement à chaque période.',
                  style: TextStyle(color: Colors.amber[300], fontSize: 12),
                ),
              ),
            const SizedBox(height: 16),

            // Dynamic Fields
            if (widget.formSchema != null && widget.formSchema!.isNotEmpty) ...[
               const Text('Informations supplémentaires', style: TextStyle(color: Colors.white, fontWeight: FontWeight.bold)),
               const SizedBox(height: 12),
               ...widget.formSchema!.map((field) => Padding(
                 padding: const EdgeInsets.only(bottom: 12),
                 child: _buildDynamicField(field),
               )).toList(),
               const SizedBox(height: 4),
            ],

            // Wallet
            if (_loadingWallets)
              const Center(child: CircularProgressIndicator())
            else
              DropdownButtonFormField<String>(
                value: _selectedWalletId,
                dropdownColor: const Color(0xFF1a1a2e),
                style: const TextStyle(color: Colors.white),
                decoration: InputDecoration(
                  labelText: 'Payer avec',
                  labelStyle: TextStyle(color: Colors.white.withOpacity(0.6)),
                  filled: true,
                  fillColor: Colors.white.withOpacity(0.1),
                  border: OutlineInputBorder(borderRadius: BorderRadius.circular(12), borderSide: BorderSide.none),
                ),
                items: _wallets.map<DropdownMenuItem<String>>((wallet) {
                  return DropdownMenuItem<String>(
                    value: wallet['id'],
                    child: Text('${wallet['name']} (${wallet['balance']} ${wallet['currency']})'),
                  );
                }).toList(),
                onChanged: (val) => setState(() => _selectedWalletId = val),
              ),
            const SizedBox(height: 16),
            
            // Message
             TextField(
              controller: _messageController,
              style: const TextStyle(color: Colors.white),
              decoration: InputDecoration(
                labelText: 'Message (Optionnel)',
                labelStyle: TextStyle(color: Colors.white.withOpacity(0.6)),
                filled: true,
                fillColor: Colors.white.withOpacity(0.1),
                border: OutlineInputBorder(borderRadius: BorderRadius.circular(12), borderSide: BorderSide.none),
              ),
            ),
            const SizedBox(height: 16),

            // Anonymous
            SwitchListTile(
              title: const Text('Rester anonyme', style: TextStyle(color: Colors.white)),
              value: _isAnonymous,
              onChanged: (val) => setState(() => _isAnonymous = val),
              activeColor: const Color(0xFF6366f1),
            ),
            const SizedBox(height: 24),

            // Submit
            SizedBox(
              height: 50,
              child: ElevatedButton(
                onPressed: _submitting ? null : _donate,
                style: ElevatedButton.styleFrom(
                  backgroundColor: const Color(0xFF6366f1),
                  shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(12)),
                ),
                child: _submitting 
                  ? const CircularProgressIndicator(color: Colors.white)
                  : const Text('Confirmer le Don', style: TextStyle(fontSize: 16, fontWeight: FontWeight.bold)),
              ),
            ),
            const SizedBox(height: 24),
          ],
        ),
      ),
    );
  }

  Widget _buildDynamicField(Map<String, dynamic> field) {
    String label = field['label'] ?? field['name'];
    if (field['required'] == true) label += ' *';
    
    // For now simple text input for all types. Can be enhanced for 'select', 'date' etc.
    final type = field['type'] ?? 'text';

    if (type == 'select' && field['options'] != null) {
      final options = List<String>.from(field['options']);
      return DropdownButtonFormField<String>(
        decoration: InputDecoration(
          labelText: label,
          labelStyle: TextStyle(color: Colors.white.withOpacity(0.6)),
          filled: true,
          fillColor: Colors.white.withOpacity(0.1),
          border: OutlineInputBorder(borderRadius: BorderRadius.circular(12), borderSide: BorderSide.none),
        ),
        dropdownColor: const Color(0xFF1a1a2e),
        style: const TextStyle(color: Colors.white),
        items: options.map((opt) => DropdownMenuItem(value: opt, child: Text(opt))).toList(),
        onChanged: (val) => _formData[field['name']] = val,
      );
    }
    
    return TextField(
      onChanged: (val) => _formData[field['name']] = val,
      style: const TextStyle(color: Colors.white),
      decoration: InputDecoration(
        labelText: label,
        labelStyle: TextStyle(color: Colors.white.withOpacity(0.6)),
        filled: true,
        fillColor: Colors.white.withOpacity(0.1),
        border: OutlineInputBorder(borderRadius: BorderRadius.circular(12), borderSide: BorderSide.none),
      ),
    );
  }

  Widget _buildFrequencyOption(String label, String value) {
    final isSelected = _frequency == value;
    return Expanded(
      child: GestureDetector(
        onTap: () => setState(() => _frequency = value),
        child: Container(
          padding: const EdgeInsets.symmetric(vertical: 10),
          decoration: BoxDecoration(
            color: isSelected ? const Color(0xFF6366f1) : Colors.white.withOpacity(0.1),
            borderRadius: BorderRadius.circular(8),
            border: Border.all(color: isSelected ? const Color(0xFF6366f1) : Colors.transparent),
          ),
          child: Text(
            label,
            textAlign: TextAlign.center,
            style: TextStyle(
              color: isSelected ? Colors.white : Colors.white.withOpacity(0.7),
              fontWeight: FontWeight.bold,
              fontSize: 12,
            ),
          ),
        ),
      ),
    );
  }
}
