import 'package:flutter/material.dart';
import '../../../../core/services/wallet_api_service.dart';
import '../../../../core/services/donation_api_service.dart';

class DonateModal extends StatefulWidget {
  final String campaignId;
  final String currency;

  const DonateModal({
    super.key,
    required this.campaignId,
    required this.currency,
  });

  @override
  State<DonateModal> createState() => _DonateModalState();
}

class _DonateModalState extends State<DonateModal> {
  final WalletApiService _walletApi = WalletApiService();
  final DonationApiService _donationApi = DonationApiService();
  final TextEditingController _amountController = TextEditingController();
  final TextEditingController _pinController = TextEditingController();
  final TextEditingController _messageController = TextEditingController();

  List<dynamic> _wallets = [];
  String? _selectedWalletId;
  bool _loadingWallets = true;
  bool _submitting = false;
  
  String _frequency = 'one_time';
  bool _isAnonymous = false;

  @override
  void initState() {
    super.initState();
    _loadWallets();
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
    if (_amountController.text.isEmpty || _selectedWalletId == null || _pinController.text.isEmpty) {
      ScaffoldMessenger.of(context).showSnackBar(const SnackBar(content: Text('Veuillez remplir tous les champs')));
      return;
    }

    setState(() => _submitting = true);
    
    try {
      await _donationApi.initiateDonation(
        campaignId: widget.campaignId,
        amount: double.parse(_amountController.text),
        currency: widget.currency,
        walletId: _selectedWalletId!,
        pin: _pinController.text,
        message: _messageController.text,
        isAnonymous: _isAnonymous,
        frequency: _frequency,
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
              ),
            ),
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
            
            // PIN
            TextField(
              controller: _pinController,
              keyboardType: TextInputType.number,
              obscureText: true,
              style: const TextStyle(color: Colors.white),
              decoration: InputDecoration(
                labelText: 'Code PIN',
                labelStyle: TextStyle(color: Colors.white.withOpacity(0.6)),
                filled: true,
                fillColor: Colors.white.withOpacity(0.1),
                border: OutlineInputBorder(borderRadius: BorderRadius.circular(12), borderSide: BorderSide.none),
              ),
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
