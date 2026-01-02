import 'package:flutter/material.dart';
import '../../core/services/association_api_service.dart';
import '../../core/services/api_client.dart';

class CreateAssociationScreen extends StatefulWidget {
  const CreateAssociationScreen({super.key});

  @override
  State<CreateAssociationScreen> createState() => _CreateAssociationScreenState();
}

class _CreateAssociationScreenState extends State<CreateAssociationScreen> {
  final AssociationApiService _api = AssociationApiService(ApiClient().dio);
  final _formKey = GlobalKey<FormState>();
  
  int _currentStep = 0;
  bool _loading = false;

  // Form data
  String _name = '';
  String _description = '';
  String _currency = 'XOF';
  String _type = 'tontine';
  double _contributionAmount = 0;
  String _frequency = 'monthly';
  bool _loansEnabled = false;
  double _loanInterestRate = 5;

  final List<Map<String, dynamic>> _types = [
    {'value': 'tontine', 'label': 'Tontine Rotative', 'icon': Icons.loop, 'description': 'Chaque membre cotise, un bénéficiaire à tour de rôle'},
    {'value': 'savings', 'label': 'Groupe d\'Épargne', 'icon': Icons.savings, 'description': 'Cotisations accumulées dans une caisse commune'},
    {'value': 'credit', 'label': 'Crédit Mutuel', 'icon': Icons.account_balance, 'description': 'Octroi de crédits aux membres avec intérêts'},
    {'value': 'general', 'label': 'Association Générale', 'icon': Icons.groups, 'description': 'Gestion simple pour clubs et amicales'},
  ];

  final List<Map<String, String>> _currencies = [
    {'value': 'XOF', 'label': 'XOF - Franc CFA (BCEAO)'},
    {'value': 'XAF', 'label': 'XAF - Franc CFA (BEAC)'},
    {'value': 'GNF', 'label': 'GNF - Franc Guinéen'},
    {'value': 'USD', 'label': 'USD - Dollar'},
    {'value': 'EUR', 'label': 'EUR - Euro'},
  ];

  Future<void> _submit() async {
    if (!_formKey.currentState!.validate()) return;
    _formKey.currentState!.save();

    setState(() => _loading = true);
    try {
      await _api.createAssociation({
        'name': _name,
        'description': _description,
        'currency': _currency,
        'type': _type,
        'rules': {
          'contribution_amount': _contributionAmount,
          'frequency': _frequency,
          'loans_enabled': _loansEnabled,
          'loan_interest_rate': _loanInterestRate,
        }
      });
      if (mounted) {
        ScaffoldMessenger.of(context).showSnackBar(
          const SnackBar(content: Text('Association créée avec succès!'), backgroundColor: Color(0xFF10b981)),
        );
        Navigator.pop(context, true);
      }
    } catch (e) {
      debugPrint('Failed to create association: $e');
      if (mounted) {
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(content: Text('Erreur: $e'), backgroundColor: Colors.red),
        );
      }
    } finally {
      setState(() => _loading = false);
    }
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      body: Container(
        decoration: const BoxDecoration(
          gradient: LinearGradient(
            begin: Alignment.topLeft,
            end: Alignment.bottomRight,
            colors: [Color(0xFF1a1a2e), Color(0xFF16213e)],
          ),
        ),
        child: SafeArea(
          child: Column(
            children: [
              _buildHeader(),
              _buildProgressIndicator(),
              Expanded(
                child: Form(
                  key: _formKey,
                  child: SingleChildScrollView(
                    padding: const EdgeInsets.all(20),
                    child: _buildCurrentStep(),
                  ),
                ),
              ),
              _buildBottomButtons(),
            ],
          ),
        ),
      ),
    );
  }

  Widget _buildHeader() {
    return Padding(
      padding: const EdgeInsets.all(20),
      child: Row(
        children: [
          GestureDetector(
            onTap: () => Navigator.pop(context),
            child: Container(
              padding: const EdgeInsets.all(8),
              decoration: BoxDecoration(
                color: Colors.white.withOpacity(0.1),
                borderRadius: BorderRadius.circular(8),
              ),
              child: const Icon(Icons.arrow_back, color: Colors.white),
            ),
          ),
          const SizedBox(width: 16),
          const Text('Créer une association', style: TextStyle(color: Colors.white, fontSize: 20, fontWeight: FontWeight.bold)),
        ],
      ),
    );
  }

  Widget _buildProgressIndicator() {
    return Padding(
      padding: const EdgeInsets.symmetric(horizontal: 20),
      child: Row(
        children: List.generate(3, (index) {
          final isActive = index <= _currentStep;
          final isCompleted = index < _currentStep;
          return Expanded(
            child: Row(
              children: [
                Container(
                  width: 32,
                  height: 32,
                  decoration: BoxDecoration(
                    color: isActive ? const Color(0xFF6366f1) : Colors.white.withOpacity(0.1),
                    shape: BoxShape.circle,
                  ),
                  child: Center(
                    child: isCompleted
                        ? const Icon(Icons.check, color: Colors.white, size: 16)
                        : Text('${index + 1}', style: TextStyle(color: isActive ? Colors.white : Colors.white.withOpacity(0.5))),
                  ),
                ),
                if (index < 2)
                  Expanded(
                    child: Container(
                      height: 2,
                      color: isCompleted ? const Color(0xFF6366f1) : Colors.white.withOpacity(0.1),
                    ),
                  ),
              ],
            ),
          );
        }),
      ),
    );
  }

  Widget _buildCurrentStep() {
    switch (_currentStep) {
      case 0:
        return _buildStep1();
      case 1:
        return _buildStep2();
      case 2:
        return _buildStep3();
      default:
        return const SizedBox();
    }
  }

  Widget _buildStep1() {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        const Text('Informations de base', style: TextStyle(color: Colors.white, fontSize: 18, fontWeight: FontWeight.bold)),
        const SizedBox(height: 24),
        _buildTextField('Nom de l\'association', (v) => _name = v, validator: (v) => v!.isEmpty ? 'Requis' : null),
        const SizedBox(height: 16),
        _buildTextField('Description', (v) => _description = v, maxLines: 3),
        const SizedBox(height: 16),
        const Text('Devise', style: TextStyle(color: Colors.white70, fontSize: 14)),
        const SizedBox(height: 8),
        Container(
          padding: const EdgeInsets.symmetric(horizontal: 16),
          decoration: BoxDecoration(
            color: Colors.white.withOpacity(0.05),
            borderRadius: BorderRadius.circular(12),
            border: Border.all(color: Colors.white.withOpacity(0.1)),
          ),
          child: DropdownButton<String>(
            value: _currency,
            isExpanded: true,
            dropdownColor: const Color(0xFF1a1a2e),
            underline: const SizedBox(),
            style: const TextStyle(color: Colors.white),
            items: _currencies.map((c) => DropdownMenuItem(value: c['value'], child: Text(c['label']!))).toList(),
            onChanged: (v) => setState(() => _currency = v!),
          ),
        ),
      ],
    );
  }

  Widget _buildStep2() {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        const Text('Type d\'association', style: TextStyle(color: Colors.white, fontSize: 18, fontWeight: FontWeight.bold)),
        const SizedBox(height: 24),
        ...List.generate(_types.length, (index) {
          final type = _types[index];
          final isSelected = _type == type['value'];
          return GestureDetector(
            onTap: () => setState(() => _type = type['value']),
            child: Container(
              margin: const EdgeInsets.only(bottom: 12),
              padding: const EdgeInsets.all(16),
              decoration: BoxDecoration(
                color: isSelected ? const Color(0xFF6366f1).withOpacity(0.2) : Colors.white.withOpacity(0.05),
                borderRadius: BorderRadius.circular(12),
                border: Border.all(color: isSelected ? const Color(0xFF6366f1) : Colors.white.withOpacity(0.1), width: isSelected ? 2 : 1),
              ),
              child: Row(
                children: [
                  Container(
                    padding: const EdgeInsets.all(12),
                    decoration: BoxDecoration(
                      color: const Color(0xFF6366f1).withOpacity(0.2),
                      borderRadius: BorderRadius.circular(12),
                    ),
                    child: Icon(type['icon'], color: const Color(0xFF6366f1)),
                  ),
                  const SizedBox(width: 16),
                  Expanded(
                    child: Column(
                      crossAxisAlignment: CrossAxisAlignment.start,
                      children: [
                        Text(type['label'], style: const TextStyle(color: Colors.white, fontWeight: FontWeight.bold)),
                        const SizedBox(height: 4),
                        Text(type['description'], style: TextStyle(color: Colors.white.withOpacity(0.6), fontSize: 12)),
                      ],
                    ),
                  ),
                  if (isSelected) const Icon(Icons.check_circle, color: Color(0xFF6366f1)),
                ],
              ),
            ),
          );
        }),
      ],
    );
  }

  Widget _buildStep3() {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        const Text('Configuration', style: TextStyle(color: Colors.white, fontSize: 18, fontWeight: FontWeight.bold)),
        const SizedBox(height: 24),
        if (_type == 'tontine' || _type == 'savings') ...[
          _buildTextField('Montant de cotisation', (v) => _contributionAmount = double.tryParse(v) ?? 0, keyboardType: TextInputType.number),
          const SizedBox(height: 16),
          const Text('Fréquence', style: TextStyle(color: Colors.white70, fontSize: 14)),
          const SizedBox(height: 8),
          Container(
            padding: const EdgeInsets.symmetric(horizontal: 16),
            decoration: BoxDecoration(
              color: Colors.white.withOpacity(0.05),
              borderRadius: BorderRadius.circular(12),
              border: Border.all(color: Colors.white.withOpacity(0.1)),
            ),
            child: DropdownButton<String>(
              value: _frequency,
              isExpanded: true,
              dropdownColor: const Color(0xFF1a1a2e),
              underline: const SizedBox(),
              style: const TextStyle(color: Colors.white),
              items: const [
                DropdownMenuItem(value: 'weekly', child: Text('Hebdomadaire')),
                DropdownMenuItem(value: 'biweekly', child: Text('Bi-mensuelle')),
                DropdownMenuItem(value: 'monthly', child: Text('Mensuelle')),
              ],
              onChanged: (v) => setState(() => _frequency = v!),
            ),
          ),
          const SizedBox(height: 24),
        ],
        SwitchListTile(
          value: _loansEnabled,
          onChanged: (v) => setState(() => _loansEnabled = v),
          title: const Text('Autoriser les prêts', style: TextStyle(color: Colors.white)),
          subtitle: Text('Permettre aux membres d\'emprunter', style: TextStyle(color: Colors.white.withOpacity(0.6))),
          activeColor: const Color(0xFF6366f1),
          contentPadding: EdgeInsets.zero,
        ),
        if (_loansEnabled) ...[
          const SizedBox(height: 16),
          _buildTextField('Taux d\'intérêt (%)', (v) => _loanInterestRate = double.tryParse(v) ?? 5, keyboardType: TextInputType.number, initialValue: '5'),
        ],
      ],
    );
  }

  Widget _buildTextField(String label, Function(String) onChanged, {String? Function(String?)? validator, int maxLines = 1, TextInputType? keyboardType, String? initialValue}) {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Text(label, style: TextStyle(color: Colors.white.withOpacity(0.7), fontSize: 14)),
        const SizedBox(height: 8),
        TextFormField(
          initialValue: initialValue,
          maxLines: maxLines,
          keyboardType: keyboardType,
          style: const TextStyle(color: Colors.white),
          decoration: InputDecoration(
            filled: true,
            fillColor: Colors.white.withOpacity(0.05),
            border: OutlineInputBorder(borderRadius: BorderRadius.circular(12), borderSide: BorderSide(color: Colors.white.withOpacity(0.1))),
            enabledBorder: OutlineInputBorder(borderRadius: BorderRadius.circular(12), borderSide: BorderSide(color: Colors.white.withOpacity(0.1))),
            focusedBorder: OutlineInputBorder(borderRadius: BorderRadius.circular(12), borderSide: const BorderSide(color: Color(0xFF6366f1))),
          ),
          validator: validator,
          onChanged: onChanged,
        ),
      ],
    );
  }

  Widget _buildBottomButtons() {
    return Container(
      padding: const EdgeInsets.all(20),
      child: Row(
        children: [
          if (_currentStep > 0)
            Expanded(
              child: OutlinedButton(
                onPressed: () => setState(() => _currentStep--),
                style: OutlinedButton.styleFrom(
                  side: BorderSide(color: Colors.white.withOpacity(0.3)),
                  padding: const EdgeInsets.symmetric(vertical: 16),
                  shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(12)),
                ),
                child: Text('Précédent', style: TextStyle(color: Colors.white.withOpacity(0.8))),
              ),
            ),
          if (_currentStep > 0) const SizedBox(width: 16),
          Expanded(
            flex: 2,
            child: ElevatedButton(
              onPressed: _loading ? null : () {
                if (_currentStep < 2) {
                  setState(() => _currentStep++);
                } else {
                  _submit();
                }
              },
              style: ElevatedButton.styleFrom(
                backgroundColor: _currentStep == 2 ? const Color(0xFF10b981) : const Color(0xFF6366f1),
                padding: const EdgeInsets.symmetric(vertical: 16),
                shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(12)),
              ),
              child: _loading
                  ? const SizedBox(width: 20, height: 20, child: CircularProgressIndicator(strokeWidth: 2, color: Colors.white))
                  : Text(_currentStep == 2 ? 'Créer' : 'Suivant'),
            ),
          ),
        ],
      ),
    );
  }
}
