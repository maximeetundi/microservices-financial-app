import 'package:flutter/material.dart';
import '../../../../core/services/api_service.dart';
import '../../data/models/enterprise_model.dart';

class InviteEmployeePage extends StatefulWidget {
  final Enterprise enterprise;

  const InviteEmployeePage({Key? key, required this.enterprise}) : super(key: key);

  @override
  State<InviteEmployeePage> createState() => _InviteEmployeePageState();
}

class _InviteEmployeePageState extends State<InviteEmployeePage> {
  final ApiService _api = ApiService();
  final _formKey = GlobalKey<FormState>();
  
  bool _isLoading = false;
  
  final _firstNameController = TextEditingController();
  final _lastNameController = TextEditingController();
  final _emailController = TextEditingController();
  final _phoneController = TextEditingController();
  final _positionController = TextEditingController();
  final _departmentController = TextEditingController();
  final _salaryController = TextEditingController();
  
  String _role = 'EMPLOYEE';
  String _salaryCurrency = 'XOF';
  String _salaryFrequency = 'MONTHLY';

  @override
  void dispose() {
    _firstNameController.dispose();
    _lastNameController.dispose();
    _emailController.dispose();
    _phoneController.dispose();
    _positionController.dispose();
    _departmentController.dispose();
    _salaryController.dispose();
    super.dispose();
  }

  Future<void> _invite() async {
    if (!_formKey.currentState!.validate()) return;
    
    setState(() => _isLoading = true);
    
    try {
      await _api.enterprise.inviteEmployee(widget.enterprise.id, {
        'first_name': _firstNameController.text,
        'last_name': _lastNameController.text,
        'email': _emailController.text,
        'phone_number': _phoneController.text,
        'job_title': _positionController.text,
        'department': _departmentController.text,
        'role': _role,
        'salary': double.tryParse(_salaryController.text) ?? 0,
        'salary_currency': _salaryCurrency,
        'salary_frequency': _salaryFrequency,
      });
      
      ScaffoldMessenger.of(context).showSnackBar(
        const SnackBar(
          content: Text('Invitation envoyée avec succès!'),
          backgroundColor: Colors.green,
        ),
      );
      
      Navigator.pop(context, true);
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
        title: const Text('Inviter un employé'),
      ),
      body: Form(
        key: _formKey,
        child: SingleChildScrollView(
          padding: const EdgeInsets.all(16),
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              // Info Banner
              Container(
                padding: const EdgeInsets.all(16),
                decoration: BoxDecoration(
                  color: Colors.blue.shade50,
                  borderRadius: BorderRadius.circular(12),
                  border: Border.all(color: Colors.blue.shade100),
                ),
                child: Row(
                  children: [
                    Icon(Icons.info_outline, color: Colors.blue.shade700),
                    const SizedBox(width: 12),
                    const Expanded(
                      child: Text(
                        'Une invitation sera envoyée par email et/ou SMS.',
                        style: TextStyle(fontSize: 13),
                      ),
                    ),
                  ],
                ),
              ),
              
              const SizedBox(height: 24),
              
              // Personal Info
              const Text('Informations personnelles', style: TextStyle(fontWeight: FontWeight.bold, fontSize: 16)),
              const SizedBox(height: 12),
              
              Row(
                children: [
                  Expanded(
                    child: TextFormField(
                      controller: _firstNameController,
                      decoration: _inputDecoration('Prénom'),
                      validator: (v) => v!.isEmpty ? 'Requis' : null,
                    ),
                  ),
                  const SizedBox(width: 12),
                  Expanded(
                    child: TextFormField(
                      controller: _lastNameController,
                      decoration: _inputDecoration('Nom'),
                      validator: (v) => v!.isEmpty ? 'Requis' : null,
                    ),
                  ),
                ],
              ),
              
              const SizedBox(height: 12),
              
              TextFormField(
                controller: _emailController,
                decoration: _inputDecoration('Email', icon: Icons.email),
                keyboardType: TextInputType.emailAddress,
                validator: (v) {
                  if (v!.isEmpty && _phoneController.text.isEmpty) {
                    return 'Email ou téléphone requis';
                  }
                  return null;
                },
              ),
              
              const SizedBox(height: 12),
              
              TextFormField(
                controller: _phoneController,
                decoration: _inputDecoration('Téléphone', icon: Icons.phone),
                keyboardType: TextInputType.phone,
              ),
              
              const SizedBox(height: 24),
              
              // Work Info
              const Text('Informations professionnelles', style: TextStyle(fontWeight: FontWeight.bold, fontSize: 16)),
              const SizedBox(height: 12),
              
              TextFormField(
                controller: _positionController,
                decoration: _inputDecoration('Poste / Titre'),
              ),
              
              const SizedBox(height: 12),
              
              TextFormField(
                controller: _departmentController,
                decoration: _inputDecoration('Département'),
              ),
              
              const SizedBox(height: 12),
              
              // Role
              DropdownButtonFormField<String>(
                value: _role,
                decoration: _inputDecoration('Rôle'),
                items: const [
                  DropdownMenuItem(value: 'EMPLOYEE', child: Text('Employé')),
                  DropdownMenuItem(value: 'ADMIN', child: Text('Administrateur')),
                ],
                onChanged: (v) => setState(() => _role = v!),
              ),
              
              const SizedBox(height: 24),
              
              // Salary
              const Text('Rémunération', style: TextStyle(fontWeight: FontWeight.bold, fontSize: 16)),
              const SizedBox(height: 12),
              
              Row(
                children: [
                  Expanded(
                    flex: 2,
                    child: TextFormField(
                      controller: _salaryController,
                      decoration: _inputDecoration('Salaire'),
                      keyboardType: TextInputType.number,
                    ),
                  ),
                  const SizedBox(width: 12),
                  Expanded(
                    child: DropdownButtonFormField<String>(
                      value: _salaryCurrency,
                      decoration: _inputDecoration('Devise'),
                      items: const [
                        DropdownMenuItem(value: 'XOF', child: Text('XOF')),
                        DropdownMenuItem(value: 'EUR', child: Text('EUR')),
                        DropdownMenuItem(value: 'USD', child: Text('USD')),
                      ],
                      onChanged: (v) => setState(() => _salaryCurrency = v!),
                    ),
                  ),
                ],
              ),
              
              const SizedBox(height: 12),
              
              DropdownButtonFormField<String>(
                value: _salaryFrequency,
                decoration: _inputDecoration('Fréquence'),
                items: const [
                  DropdownMenuItem(value: 'MONTHLY', child: Text('Mensuel')),
                  DropdownMenuItem(value: 'WEEKLY', child: Text('Hebdomadaire')),
                  DropdownMenuItem(value: 'DAILY', child: Text('Journalier')),
                ],
                onChanged: (v) => setState(() => _salaryFrequency = v!),
              ),
              
              const SizedBox(height: 32),
              
              // Submit
              SizedBox(
                width: double.infinity,
                child: ElevatedButton(
                  onPressed: _isLoading ? null : _invite,
                  style: ElevatedButton.styleFrom(
                    backgroundColor: Colors.blue,
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
                      : const Text('Envoyer l\'invitation'),
                ),
              ),
            ],
          ),
        ),
      ),
    );
  }

  InputDecoration _inputDecoration(String label, {IconData? icon}) {
    return InputDecoration(
      labelText: label,
      prefixIcon: icon != null ? Icon(icon, size: 20) : null,
      border: OutlineInputBorder(borderRadius: BorderRadius.circular(12)),
      contentPadding: const EdgeInsets.symmetric(horizontal: 16, vertical: 14),
    );
  }
}
