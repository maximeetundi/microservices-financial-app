
import 'package:flutter/material.dart';
import '../../data/repositories/subscription_repository.dart';
import '../../../../core/di/injection_container.dart';
import '../../../../features/auth/presentation/pages/pin_verify_dialog.dart';

class LinkSubscriptionModal extends StatefulWidget {
  final VoidCallback onSuccess;

  const LinkSubscriptionModal({Key? key, required this.onSuccess}) : super(key: key);

  @override
  _LinkSubscriptionModalState createState() => _LinkSubscriptionModalState();
}

class _LinkSubscriptionModalState extends State<LinkSubscriptionModal> {
  final _matriculeController = TextEditingController();
  // For MVP, we might hardcode enterprise selection or assume user knows ID, 
  // but better to search. Simplest: Enter Matricule + Select Enterprise from list?
  // Or just "Enter Subscription Code" if we made a unified system.
  // Let's assume user must select Enterprise first.
  
  // Actually, searching enterprises requires EnterpriseRepository or similar.
  // For MVP, let's keep it simple: Just Matricule? No, ID collision possible.
  // We need Enterprise ID.
  // Let's implement a simple dropdown of "Known Enterprises" (mock or fetched).
  // Ideally, search by name.
  
  bool _isLoading = false;
  
  // Dynamic Form Data
  final Map<String, dynamic> _formData = {};

  final SubscriptionRepository _repository = sl<SubscriptionRepository>();

  @override
  Widget build(BuildContext context) {
    return Container(
      padding: EdgeInsets.all(16),
      child: Column(
        mainAxisSize: MainAxisSize.min,
        children: [
          Text("Lier un Compte", style: TextStyle(fontSize: 18, fontWeight: FontWeight.bold)),
          SizedBox(height: 16),
          // TODO: Enterprise Search/Select
          TextField(
             decoration: InputDecoration(
               labelText: "ID Entreprise (Demo)",
               hintText: "Entrer l'ID de l'école/service",
               border: OutlineInputBorder(),
             ),
             onChanged: (val) => _formData['enterprise_id'] = val,
          ),
          SizedBox(height: 12),
          TextField(
            controller: _matriculeController,
            decoration: InputDecoration(
              labelText: "Votre Matricule / N° Compteur",
              border: OutlineInputBorder(),
            ),
          ),
          // Dynamic fields would be rendered here if we fetched form schema first...
          // For MVP 1.0, let's just send matricule.
          SizedBox(height: 20),
          ElevatedButton(
            onPressed: _isLoading ? null : _submit,
            child: _isLoading ? CircularProgressIndicator() : Text("Lier"),
            style: ElevatedButton.styleFrom(
              minimumSize: Size(double.infinity, 50),
            ),
          )
        ],
      ),
    );
  }

  Future<void> _submit() async {
    // 1. Verify PIN first
    final pinVerified = await PinVerifyDialog.show(
      context,
      title: 'Confirmation requise',
      subtitle: 'Entrez votre code PIN pour valider',
      allowBiometric: true,
    );

    if (pinVerified != true) {
      // User cancelled or failed
      return;
    }

    setState(() => _isLoading = true);
    try {
      await _repository.linkSubscription(
        enterpriseId: _formData['enterprise_id'] ?? '', 
        matricule: _matriculeController.text,
        formData: {}, // Empty for now
      );
      Navigator.pop(context);
      widget.onSuccess();
      ScaffoldMessenger.of(context).showSnackBar(SnackBar(content: Text('Succès!')));
    } catch (e) {
      ScaffoldMessenger.of(context).showSnackBar(SnackBar(content: Text('Erreur: $e')));
    } finally {
      if (mounted) setState(() => _isLoading = false);
    }
  }
}
