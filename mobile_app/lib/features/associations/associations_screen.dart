import 'package:flutter/material.dart';
import 'package:intl/intl.dart';
import 'association_details_screen.dart';
import 'create_association_screen.dart';
import '../../core/services/association_api_service.dart';
import '../../core/api/api_client.dart';

class AssociationsScreen extends StatefulWidget {
  const AssociationsScreen({super.key});

  @override
  State<AssociationsScreen> createState() => _AssociationsScreenState();
}

class _AssociationsScreenState extends State<AssociationsScreen> {
  final AssociationApiService _api = AssociationApiService(ApiClient().dio);
  List<dynamic> _associations = [];
  bool _loading = true;

  @override
  void initState() {
    super.initState();
    _loadAssociations();
  }

  Future<void> _loadAssociations() async {
    setState(() => _loading = true);
    try {
      final response = await _api.getAssociations();
      setState(() {
        _associations = response.data is List ? response.data : [];
      });
    } catch (e) {
      debugPrint('Failed to load associations: $e');
      // Mock data for demo
      setState(() {
        _associations = [
          {
            'id': '1',
            'name': 'Tontine Famille Toure',
            'type': 'tontine',
            'description': 'Tontine mensuelle familiale',
            'total_members': 12,
            'treasury_balance': 1200000,
            'currency': 'XOF',
            'status': 'active'
          },
          {
            'id': '2',
            'name': 'Épargne Commerçants',
            'type': 'savings',
            'description': 'Caisse de solidarité du marché',
            'total_members': 45,
            'treasury_balance': 5500000,
            'currency': 'XOF',
            'status': 'active'
          }
        ];
      });
    } finally {
      setState(() => _loading = false);
    }
  }

  String _formatCurrency(dynamic amount, String currency) {
    final value = (amount is num) ? amount.toDouble() : 0.0;
    final formatter = NumberFormat.currency(locale: 'fr_FR', symbol: currency, decimalDigits: 0);
    return formatter.format(value);
  }

  IconData _getTypeIcon(String type) {
    switch (type) {
      case 'tontine':
        return Icons.loop;
      case 'savings':
        return Icons.savings;
      case 'credit':
        return Icons.account_balance;
      default:
        return Icons.groups;
    }
  }

  String _getTypeLabel(String type) {
    switch (type) {
      case 'tontine':
        return 'Tontine';
      case 'savings':
        return 'Épargne';
      case 'credit':
        return 'Crédit';
      default:
        return 'Association';
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
              // Header
              Padding(
                padding: const EdgeInsets.all(20),
                child: Row(
                  mainAxisAlignment: MainAxisAlignment.spaceBetween,
                  children: [
                    Column(
                      crossAxisAlignment: CrossAxisAlignment.start,
                      children: [
                        const Text(
                          'Mes Associations',
                          style: TextStyle(
                            color: Colors.white,
                            fontSize: 24,
                            fontWeight: FontWeight.bold,
                          ),
                        ),
                        Text(
                          '${_associations.length} association(s)',
                          style: TextStyle(color: Colors.white.withOpacity(0.7)),
                        ),
                      ],
                    ),
                    IconButton(
                      onPressed: _loadAssociations,
                      icon: const Icon(Icons.refresh, color: Colors.white),
                    ),
                  ],
                ),
              ),

              // Quick Stats
              Padding(
                padding: const EdgeInsets.symmetric(horizontal: 20),
                child: Row(
                  children: [
                    _buildStatCard('Total Épargné', '2.5M FCFA', Icons.savings, const Color(0xFF6366f1)),
                    const SizedBox(width: 12),
                    _buildStatCard('Prêts', '150K FCFA', Icons.account_balance, const Color(0xFF10b981)),
                  ],
                ),
              ),

              const SizedBox(height: 20),

              // Associations List
              Expanded(
                child: _loading
                    ? const Center(child: CircularProgressIndicator(color: Color(0xFF6366f1)))
                    : _associations.isEmpty
                        ? _buildEmptyState()
                        : RefreshIndicator(
                            onRefresh: _loadAssociations,
                            child: ListView.builder(
                              padding: const EdgeInsets.symmetric(horizontal: 20),
                              itemCount: _associations.length,
                              itemBuilder: (context, index) => _buildAssociationCard(_associations[index]),
                            ),
                          ),
              ),
            ],
          ),
        ),
      ),
      floatingActionButton: FloatingActionButton.extended(
        onPressed: () => Navigator.push(
          context,
          MaterialPageRoute(builder: (_) => const CreateAssociationScreen()),
        ).then((_) => _loadAssociations()),
        backgroundColor: const Color(0xFF6366f1),
        icon: const Icon(Icons.add),
        label: const Text('Créer'),
      ),
    );
  }

  Widget _buildStatCard(String title, String value, IconData icon, Color color) {
    return Expanded(
      child: Container(
        padding: const EdgeInsets.all(16),
        decoration: BoxDecoration(
          color: color.withOpacity(0.15),
          borderRadius: BorderRadius.circular(16),
          border: Border.all(color: color.withOpacity(0.3)),
        ),
        child: Row(
          children: [
            Container(
              padding: const EdgeInsets.all(8),
              decoration: BoxDecoration(
                color: color.withOpacity(0.2),
                borderRadius: BorderRadius.circular(8),
              ),
              child: Icon(icon, color: color, size: 20),
            ),
            const SizedBox(width: 12),
            Expanded(
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  Text(title, style: TextStyle(color: Colors.white.withOpacity(0.7), fontSize: 11)),
                  Text(value, style: const TextStyle(color: Colors.white, fontWeight: FontWeight.bold, fontSize: 14)),
                ],
              ),
            ),
          ],
        ),
      ),
    );
  }

  Widget _buildEmptyState() {
    return Center(
      child: Column(
        mainAxisAlignment: MainAxisAlignment.center,
        children: [
          Icon(Icons.groups, size: 64, color: Colors.white.withOpacity(0.3)),
          const SizedBox(height: 16),
          Text(
            'Aucune association',
            style: TextStyle(color: Colors.white.withOpacity(0.7), fontSize: 18),
          ),
          const SizedBox(height: 8),
          Text(
            'Créez ou rejoignez une association',
            style: TextStyle(color: Colors.white.withOpacity(0.5)),
          ),
        ],
      ),
    );
  }

  Widget _buildAssociationCard(Map<String, dynamic> association) {
    final type = association['type'] ?? 'general';
    final currency = association['currency'] ?? 'XOF';

    return GestureDetector(
      onTap: () => Navigator.push(
        context,
        MaterialPageRoute(
          builder: (_) => AssociationDetailsScreen(associationId: association['id']),
        ),
      ).then((_) => _loadAssociations()),
      child: Container(
        margin: const EdgeInsets.only(bottom: 16),
        decoration: BoxDecoration(
          color: Colors.white.withOpacity(0.05),
          borderRadius: BorderRadius.circular(16),
          border: Border.all(color: Colors.white.withOpacity(0.1)),
        ),
        child: Column(
          children: [
            Padding(
              padding: const EdgeInsets.all(16),
              child: Row(
                children: [
                  Container(
                    width: 48,
                    height: 48,
                    decoration: BoxDecoration(
                      gradient: const LinearGradient(
                        colors: [Color(0xFF6366f1), Color(0xFF8b5cf6)],
                      ),
                      borderRadius: BorderRadius.circular(12),
                    ),
                    child: Icon(_getTypeIcon(type), color: Colors.white),
                  ),
                  const SizedBox(width: 12),
                  Expanded(
                    child: Column(
                      crossAxisAlignment: CrossAxisAlignment.start,
                      children: [
                        Text(
                          association['name'] ?? 'Sans nom',
                          style: const TextStyle(
                            color: Colors.white,
                            fontWeight: FontWeight.bold,
                            fontSize: 16,
                          ),
                        ),
                        const SizedBox(height: 4),
                        Text(
                          _getTypeLabel(type),
                          style: TextStyle(color: Colors.white.withOpacity(0.6), fontSize: 12),
                        ),
                      ],
                    ),
                  ),
                  Container(
                    padding: const EdgeInsets.symmetric(horizontal: 10, vertical: 4),
                    decoration: BoxDecoration(
                      color: const Color(0xFF10b981).withOpacity(0.2),
                      borderRadius: BorderRadius.circular(12),
                    ),
                    child: const Text(
                      'Actif',
                      style: TextStyle(color: Color(0xFF10b981), fontSize: 12, fontWeight: FontWeight.w500),
                    ),
                  ),
                ],
              ),
            ),
            Container(
              padding: const EdgeInsets.all(16),
              decoration: BoxDecoration(
                color: Colors.white.withOpacity(0.03),
                borderRadius: const BorderRadius.only(
                  bottomLeft: Radius.circular(16),
                  bottomRight: Radius.circular(16),
                ),
              ),
              child: Row(
                mainAxisAlignment: MainAxisAlignment.spaceBetween,
                children: [
                  Row(
                    children: [
                      Icon(Icons.people, size: 16, color: Colors.white.withOpacity(0.5)),
                      const SizedBox(width: 4),
                      Text(
                        '${association['total_members'] ?? 0} membres',
                        style: TextStyle(color: Colors.white.withOpacity(0.7), fontSize: 13),
                      ),
                    ],
                  ),
                  Text(
                    _formatCurrency(association['treasury_balance'], currency),
                    style: const TextStyle(
                      color: Color(0xFF6366f1),
                      fontWeight: FontWeight.bold,
                      fontSize: 14,
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
}
