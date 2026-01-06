import 'package:flutter/material.dart';
import 'package:intl/intl.dart';
import '../../core/services/association_api_service.dart';
import '../../core/api/api_client.dart';
import 'widgets/pay_contribution_sheet.dart';

class AssociationDetailsScreen extends StatefulWidget {
  final String associationId;

  const AssociationDetailsScreen({super.key, required this.associationId});

  @override
  State<AssociationDetailsScreen> createState() => _AssociationDetailsScreenState();
}

class _AssociationDetailsScreenState extends State<AssociationDetailsScreen> with SingleTickerProviderStateMixin {
  final AssociationApiService _api = AssociationApiService(ApiClient().dio);
  late TabController _tabController;
  
  Map<String, dynamic>? _association;
  List<dynamic> _members = [];
  Map<String, dynamic> _treasury = {};
  List<dynamic> _meetings = [];
  bool _loading = true;

  @override
  void initState() {
    super.initState();
    _tabController = TabController(length: 4, vsync: this);
    _loadData();
  }

  @override
  void dispose() {
    _tabController.dispose();
    super.dispose();
  }

  Future<void> _loadData() async {
    setState(() => _loading = true);
    try {
      final responses = await Future.wait([
        _api.getAssociation(widget.associationId),
        _api.getMembers(widget.associationId),
        _api.getTreasury(widget.associationId),
        _api.getMeetings(widget.associationId),
      ]);
      setState(() {
        _association = responses[0].data;
        _members = responses[1].data is List ? responses[1].data : [];
        _treasury = responses[2].data is Map ? responses[2].data : {};
        _meetings = responses[3].data is List ? responses[3].data : [];
      });
    } catch (e) {
      debugPrint('Failed to load association: $e');
      // Mock data
      setState(() {
        _association = {
          'id': widget.associationId,
          'name': 'Tontine Famille Toure',
          'type': 'tontine',
          'total_members': 12,
          'treasury_balance': 1200000,
          'currency': 'XOF',
          'status': 'active'
        };
        _members = [
          {'id': '1', 'user_name': 'Mamadou Toure', 'role': 'president', 'contributions_paid': 150000},
          {'id': '2', 'user_name': 'Fatou Diallo', 'role': 'treasurer', 'contributions_paid': 150000},
          {'id': '3', 'user_name': 'Ibrahim Kone', 'role': 'member', 'contributions_paid': 100000},
        ];
        _treasury = {
          'total_balance': 1200000,
          'total_contributions': 1500000,
          'total_loans': 300000,
          'transactions': [
            {'id': '1', 'type': 'contribution', 'amount': 50000, 'description': 'Cotisation Janvier', 'created_at': '2024-01-15'},
          ]
        };
      });
    } finally {
      setState(() => _loading = false);
    }
  }

  String _formatCurrency(dynamic amount) {
    final value = (amount is num) ? amount.toDouble() : 0.0;
    final currency = _association?['currency'] ?? 'XOF';
    return NumberFormat.currency(locale: 'fr_FR', symbol: currency, decimalDigits: 0).format(value);
  }

  String _getRoleLabel(String role) {
    switch (role) {
      case 'president': return 'Président';
      case 'treasurer': return 'Trésorier';
      case 'secretary': return 'Secrétaire';
      default: return 'Membre';
    }
  }

  String _getInitials(String? name) {
    if (name == null || name.isEmpty) return '?';
    return name.split(' ').map((n) => n.isNotEmpty ? n[0] : '').join().toUpperCase().substring(0, 2.clamp(0, 2));
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
          child: _loading
              ? const Center(child: CircularProgressIndicator(color: Color(0xFF6366f1)))
              : Column(
                  children: [
                    _buildHeader(),
                    _buildStatsBar(),
                    _buildTabBar(),
                    Expanded(child: _buildTabContent()),
                  ],
                ),
        ),
      ),
      floatingActionButton: FloatingActionButton.extended(
        onPressed: () => showPayContributionSheet(
          context,
          widget.associationId,
          _association?['currency'] ?? 'XOF',
          _loadData,
        ),
        backgroundColor: const Color(0xFF6366f1),
        icon: const Icon(Icons.add),
        label: const Text('Cotisation'),
      ),
    );
  }

  Widget _buildHeader() {
    return Container(
      padding: const EdgeInsets.all(20),
      decoration: BoxDecoration(
        gradient: LinearGradient(
          colors: [const Color(0xFF6366f1), const Color(0xFF8b5cf6).withOpacity(0.8)],
        ),
      ),
      child: Row(
        children: [
          GestureDetector(
            onTap: () => Navigator.pop(context),
            child: Container(
              padding: const EdgeInsets.all(8),
              decoration: BoxDecoration(
                color: Colors.white.withOpacity(0.2),
                borderRadius: BorderRadius.circular(8),
              ),
              child: const Icon(Icons.arrow_back, color: Colors.white),
            ),
          ),
          const SizedBox(width: 16),
          Expanded(
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Text(
                  _association?['name'] ?? 'Association',
                  style: const TextStyle(color: Colors.white, fontSize: 20, fontWeight: FontWeight.bold),
                ),
                Text(
                  _association?['type'] == 'tontine' ? 'Tontine Rotative' : 'Association',
                  style: TextStyle(color: Colors.white.withOpacity(0.8)),
                ),
              ],
            ),
          ),
          Container(
            padding: const EdgeInsets.symmetric(horizontal: 12, vertical: 6),
            decoration: BoxDecoration(
              color: Colors.white.withOpacity(0.2),
              borderRadius: BorderRadius.circular(20),
            ),
            child: const Text('Actif', style: TextStyle(color: Colors.white, fontWeight: FontWeight.w500)),
          ),
        ],
      ),
    );
  }

  Widget _buildStatsBar() {
    return Container(
      padding: const EdgeInsets.symmetric(vertical: 16),
      decoration: BoxDecoration(
        color: Colors.white.withOpacity(0.05),
        border: Border(bottom: BorderSide(color: Colors.white.withOpacity(0.1))),
      ),
      child: Row(
        mainAxisAlignment: MainAxisAlignment.spaceEvenly,
        children: [
          _buildStatItem('${_association?['total_members'] ?? 0}', 'Membres'),
          _buildStatItem(_formatCurrency(_association?['treasury_balance']), 'Caisse'),
          _buildStatItem('${_members.where((m) => m['role'] == 'president').length}', 'Responsables'),
        ],
      ),
    );
  }

  Widget _buildStatItem(String value, String label) {
    return Column(
      children: [
        Text(value, style: const TextStyle(color: Color(0xFF6366f1), fontSize: 18, fontWeight: FontWeight.bold)),
        const SizedBox(height: 4),
        Text(label, style: TextStyle(color: Colors.white.withOpacity(0.6), fontSize: 12)),
      ],
    );
  }

  Widget _buildTabBar() {
    return TabBar(
      controller: _tabController,
      indicatorColor: const Color(0xFF6366f1),
      labelColor: const Color(0xFF6366f1),
      unselectedLabelColor: Colors.white.withOpacity(0.5),
      tabs: const [
        Tab(icon: Icon(Icons.people), text: 'Membres'),
        Tab(icon: Icon(Icons.account_balance_wallet), text: 'Caisse'),
        Tab(icon: Icon(Icons.event), text: 'Réunions'),
        Tab(icon: Icon(Icons.monetization_on), text: 'Prêts'),
      ],
    );
  }

  Widget _buildTabContent() {
    return TabBarView(
      controller: _tabController,
      children: [
        _buildMembersTab(),
        _buildTreasuryTab(),
        _buildMeetingsTab(),
        _buildLoansTab(),
      ],
    );
  }

  Widget _buildMembersTab() {
    return ListView.builder(
      padding: const EdgeInsets.all(16),
      itemCount: _members.length,
      itemBuilder: (context, index) {
        final member = _members[index];
        return Container(
          margin: const EdgeInsets.only(bottom: 12),
          padding: const EdgeInsets.all(16),
          decoration: BoxDecoration(
            color: Colors.white.withOpacity(0.05),
            borderRadius: BorderRadius.circular(12),
            border: Border.all(color: Colors.white.withOpacity(0.1)),
          ),
          child: Row(
            children: [
              CircleAvatar(
                backgroundColor: const Color(0xFF6366f1).withOpacity(0.3),
                child: Text(_getInitials(member['user_name']), style: const TextStyle(color: Color(0xFF6366f1))),
              ),
              const SizedBox(width: 12),
              Expanded(
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    Text(member['user_name'] ?? 'Membre', style: const TextStyle(color: Colors.white, fontWeight: FontWeight.w500)),
                    Text(_getRoleLabel(member['role'] ?? 'member'), style: TextStyle(color: Colors.white.withOpacity(0.6), fontSize: 12)),
                  ],
                ),
              ),
              Column(
                crossAxisAlignment: CrossAxisAlignment.end,
                children: [
                  Text(_formatCurrency(member['contributions_paid']), style: const TextStyle(color: Color(0xFF10b981), fontWeight: FontWeight.bold)),
                  Text('cotisé', style: TextStyle(color: Colors.white.withOpacity(0.5), fontSize: 11)),
                ],
              ),
            ],
          ),
        );
      },
    );
  }

  Widget _buildTreasuryTab() {
    final transactions = (_treasury['transactions'] as List?) ?? [];
    return SingleChildScrollView(
      padding: const EdgeInsets.all(16),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Row(
            children: [
              Expanded(child: _buildTreasuryCard('Cotisations', _formatCurrency(_treasury['total_contributions']), const Color(0xFF10b981))),
              const SizedBox(width: 12),
              Expanded(child: _buildTreasuryCard('Prêts', _formatCurrency(_treasury['total_loans']), const Color(0xFFef4444))),
            ],
          ),
          const SizedBox(height: 16),
          _buildTreasuryCard('Solde Total', _formatCurrency(_treasury['total_balance']), const Color(0xFF6366f1), full: true),
          const SizedBox(height: 24),
          const Text('Dernières transactions', style: TextStyle(color: Colors.white, fontSize: 16, fontWeight: FontWeight.bold)),
          const SizedBox(height: 12),
          ...transactions.take(10).map((tx) => _buildTransactionItem(tx)),
        ],
      ),
    );
  }

  Widget _buildTreasuryCard(String label, String value, Color color, {bool full = false}) {
    return Container(
      width: full ? double.infinity : null,
      padding: const EdgeInsets.all(16),
      decoration: BoxDecoration(
        color: color.withOpacity(0.15),
        borderRadius: BorderRadius.circular(12),
        border: Border.all(color: color.withOpacity(0.3)),
      ),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Text(label, style: TextStyle(color: color, fontSize: 12)),
          const SizedBox(height: 4),
          Text(value, style: TextStyle(color: color, fontSize: 20, fontWeight: FontWeight.bold)),
        ],
      ),
    );
  }

  Widget _buildTransactionItem(Map<String, dynamic> tx) {
    final isCredit = tx['type'] == 'contribution';
    return Container(
      margin: const EdgeInsets.only(bottom: 8),
      padding: const EdgeInsets.all(12),
      decoration: BoxDecoration(
        color: Colors.white.withOpacity(0.03),
        borderRadius: BorderRadius.circular(8),
      ),
      child: Row(
        children: [
          Container(
            padding: const EdgeInsets.all(8),
            decoration: BoxDecoration(
              color: (isCredit ? const Color(0xFF10b981) : const Color(0xFFef4444)).withOpacity(0.2),
              borderRadius: BorderRadius.circular(8),
            ),
            child: Icon(isCredit ? Icons.arrow_upward : Icons.arrow_downward, color: isCredit ? const Color(0xFF10b981) : const Color(0xFFef4444), size: 16),
          ),
          const SizedBox(width: 12),
          Expanded(
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Text(tx['description'] ?? 'Transaction', style: const TextStyle(color: Colors.white, fontSize: 13)),
                Text(tx['created_at'] ?? '', style: TextStyle(color: Colors.white.withOpacity(0.5), fontSize: 11)),
              ],
            ),
          ),
          Text(
            '${isCredit ? '+' : '-'}${_formatCurrency(tx['amount'])}',
            style: TextStyle(color: isCredit ? const Color(0xFF10b981) : const Color(0xFFef4444), fontWeight: FontWeight.bold),
          ),
        ],
      ),
    );
  }

  Widget _buildMeetingsTab() {
    if (_meetings.isEmpty) {
      return Center(
        child: Column(
          mainAxisAlignment: MainAxisAlignment.center,
          children: [
            Icon(Icons.event_busy, size: 48, color: Colors.white.withOpacity(0.3)),
            const SizedBox(height: 16),
            Text('Aucune réunion', style: TextStyle(color: Colors.white.withOpacity(0.6))),
          ],
        ),
      );
    }
    return ListView.builder(
      padding: const EdgeInsets.all(16),
      itemCount: _meetings.length,
      itemBuilder: (context, index) {
        final meeting = _meetings[index];
        return Container(
          margin: const EdgeInsets.only(bottom: 12),
          padding: const EdgeInsets.all(16),
          decoration: BoxDecoration(
            color: Colors.white.withOpacity(0.05),
            borderRadius: BorderRadius.circular(12),
          ),
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              Text(meeting['title'] ?? 'Réunion', style: const TextStyle(color: Colors.white, fontWeight: FontWeight.bold)),
              const SizedBox(height: 8),
              Row(
                children: [
                  Icon(Icons.location_on, size: 14, color: Colors.white.withOpacity(0.5)),
                  const SizedBox(width: 4),
                  Text(meeting['location'] ?? 'Lieu non défini', style: TextStyle(color: Colors.white.withOpacity(0.6), fontSize: 12)),
                ],
              ),
            ],
          ),
        );
      },
    );
  }

  Widget _buildLoansTab() {
    return Center(
      child: Column(
        mainAxisAlignment: MainAxisAlignment.center,
        children: [
          Icon(Icons.money_off, size: 48, color: Colors.white.withOpacity(0.3)),
          const SizedBox(height: 16),
          Text('Aucun prêt en cours', style: TextStyle(color: Colors.white.withOpacity(0.6))),
          const SizedBox(height: 24),
          ElevatedButton.icon(
            onPressed: () {},
            icon: const Icon(Icons.add),
            label: const Text('Demander un prêt'),
            style: ElevatedButton.styleFrom(backgroundColor: const Color(0xFF6366f1)),
          ),
        ],
      ),
    );
  }
}
