import 'package:flutter/material.dart';
import '../../../../../core/services/api_service.dart';
import '../../../data/models/enterprise_model.dart';

class PayrollTab extends StatefulWidget {
  final Enterprise enterprise;

  const PayrollTab({Key? key, required this.enterprise}) : super(key: key);

  @override
  State<PayrollTab> createState() => _PayrollTabState();
}

class _PayrollTabState extends State<PayrollTab> {
  final ApiService _api = ApiService();
  bool _isLoading = true;
  Map<String, dynamic>? _payrollPreview;

  @override
  void initState() {
    super.initState();
    _loadPayroll();
  }

  Future<void> _loadPayroll() async {
    setState(() => _isLoading = true);
    try {
      final response = await _api.enterprise.getPayrollPreview(widget.enterprise.id);
      _payrollPreview = response is Map ? Map<String, dynamic>.from(response) : null;
    } catch (e) {
      debugPrint('Error loading payroll: $e');
    } finally {
      setState(() => _isLoading = false);
    }
  }

  @override
  Widget build(BuildContext context) {
    return RefreshIndicator(
      onRefresh: _loadPayroll,
      child: SingleChildScrollView(
        physics: const AlwaysScrollableScrollPhysics(),
        padding: const EdgeInsets.all(16),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            // Header Card
            Container(
              width: double.infinity,
              padding: const EdgeInsets.all(20),
              decoration: BoxDecoration(
                gradient: LinearGradient(
                  colors: [Colors.purple.shade700, Colors.purple.shade500],
                  begin: Alignment.topLeft,
                  end: Alignment.bottomRight,
                ),
                borderRadius: BorderRadius.circular(16),
                boxShadow: [
                  BoxShadow(
                    color: Colors.purple.withOpacity(0.3),
                    blurRadius: 12,
                    offset: const Offset(0, 6),
                  ),
                ],
              ),
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  const Row(
                    children: [
                      Icon(Icons.payments, color: Colors.white, size: 28),
                      SizedBox(width: 12),
                      Text(
                        'Gestion de la paie',
                        style: TextStyle(
                          color: Colors.white,
                          fontWeight: FontWeight.bold,
                          fontSize: 20,
                        ),
                      ),
                    ],
                  ),
                  const SizedBox(height: 16),
                  if (_isLoading)
                    const Center(child: CircularProgressIndicator(color: Colors.white))
                  else if (_payrollPreview != null) ...[
                    Row(
                      children: [
                        Expanded(
                          child: _PayrollStat(
                            label: 'Employés',
                            value: '${_payrollPreview!['employee_count'] ?? 0}',
                          ),
                        ),
                        Expanded(
                          child: _PayrollStat(
                            label: 'Total mensuel',
                            value: '${(_payrollPreview!['total_amount'] ?? 0).toStringAsFixed(0)}',
                          ),
                        ),
                      ],
                    ),
                  ] else
                    const Text(
                      'Configurez les salaires de vos employés',
                      style: TextStyle(color: Colors.white70),
                    ),
                ],
              ),
            ),
            
            const SizedBox(height: 24),
            
            // Actions
            const Text(
              'Actions',
              style: TextStyle(fontSize: 18, fontWeight: FontWeight.bold),
            ),
            const SizedBox(height: 12),
            
            _ActionCard(
              icon: Icons.play_arrow,
              title: 'Exécuter la paie',
              subtitle: 'Lancer le paiement des salaires',
              color: Colors.green,
              onTap: _runPayroll,
            ),
            
            _ActionCard(
              icon: Icons.history,
              title: 'Historique',
              subtitle: 'Voir les paiements passés',
              color: Colors.blue,
              onTap: _viewHistory,
            ),
            
            _ActionCard(
              icon: Icons.settings,
              title: 'Paramètres de paie',
              subtitle: 'Configurer les salaires',
              color: Colors.orange,
              onTap: _configurePayroll,
            ),
            
            const SizedBox(height: 24),
            
            // Next Payroll Info
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
                      'La prochaine paie nécessitera l\'approbation d\'un autre administrateur.',
                      style: TextStyle(fontSize: 13),
                    ),
                  ),
                ],
              ),
            ),
          ],
        ),
      ),
    );
  }

  void _runPayroll() {
    showDialog(
      context: context,
      builder: (context) => AlertDialog(
        title: const Text('Exécuter la paie'),
        content: const Text(
          'Cette action créera une demande d\'approbation pour le paiement des salaires. '
          'Un autre administrateur devra approuver la transaction.',
        ),
        actions: [
          TextButton(
            onPressed: () => Navigator.pop(context),
            child: const Text('Annuler'),
          ),
          ElevatedButton(
            onPressed: () {
              Navigator.pop(context);
              _initiatePayroll();
            },
            child: const Text('Continuer'),
          ),
        ],
      ),
    );
  }

  void _initiatePayroll() async {
    // TODO: Implement payroll initiation with PIN
    ScaffoldMessenger.of(context).showSnackBar(
      const SnackBar(content: Text('Demande de paie créée')),
    );
  }

  void _viewHistory() {
    // TODO: Implement history view
  }

  void _configurePayroll() {
    // TODO: Implement payroll configuration
  }
}

class _PayrollStat extends StatelessWidget {
  final String label;
  final String value;

  const _PayrollStat({required this.label, required this.value});

  @override
  Widget build(BuildContext context) {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Text(label, style: const TextStyle(color: Colors.white60, fontSize: 12)),
        const SizedBox(height: 4),
        Text(
          value,
          style: const TextStyle(
            color: Colors.white,
            fontWeight: FontWeight.bold,
            fontSize: 24,
          ),
        ),
      ],
    );
  }
}

class _ActionCard extends StatelessWidget {
  final IconData icon;
  final String title;
  final String subtitle;
  final Color color;
  final VoidCallback onTap;

  const _ActionCard({
    required this.icon,
    required this.title,
    required this.subtitle,
    required this.color,
    required this.onTap,
  });

  @override
  Widget build(BuildContext context) {
    return Card(
      margin: const EdgeInsets.only(bottom: 12),
      shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(12)),
      child: InkWell(
        onTap: onTap,
        borderRadius: BorderRadius.circular(12),
        child: Padding(
          padding: const EdgeInsets.all(16),
          child: Row(
            children: [
              Container(
                padding: const EdgeInsets.all(12),
                decoration: BoxDecoration(
                  color: color.withOpacity(0.1),
                  borderRadius: BorderRadius.circular(12),
                ),
                child: Icon(icon, color: color),
              ),
              const SizedBox(width: 16),
              Expanded(
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    Text(title, style: const TextStyle(fontWeight: FontWeight.w600)),
                    Text(subtitle, style: TextStyle(color: Colors.grey[600], fontSize: 13)),
                  ],
                ),
              ),
              Icon(Icons.chevron_right, color: Colors.grey[400]),
            ],
          ),
        ),
      ),
    );
  }
}
