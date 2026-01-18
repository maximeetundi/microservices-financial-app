import 'package:flutter/material.dart';
import '../../../../core/services/api_service.dart';
import '../../data/models/enterprise_model.dart';
import '../../data/models/employee_model.dart';
import 'tabs/overview_tab.dart';
import 'tabs/employees_tab.dart';
import 'tabs/wallets_tab.dart';
import 'tabs/approvals_tab.dart';
import 'tabs/services_tab.dart';
import 'tabs/payroll_tab.dart';
import 'tabs/settings_tab.dart';

class EnterpriseDashboardPage extends StatefulWidget {
  final String? enterpriseId;
  
  const EnterpriseDashboardPage({Key? key, this.enterpriseId}) : super(key: key);

  @override
  State<EnterpriseDashboardPage> createState() => _EnterpriseDashboardPageState();
}

class _EnterpriseDashboardPageState extends State<EnterpriseDashboardPage> with SingleTickerProviderStateMixin {
  final ApiService _api = ApiService();
  
  bool _isLoading = true;
  String? _error;
  Enterprise? _enterprise;
  Employee? _currentEmployee;
  
  late TabController _tabController;
  
  // All possible tabs
  final List<_TabItem> _allTabs = [
    _TabItem(id: 'overview', title: 'Aperçu', icon: Icons.dashboard, adminOnly: false),
    _TabItem(id: 'employees', title: 'Employés', icon: Icons.people, adminOnly: true),
    _TabItem(id: 'wallets', title: 'Portefeuilles', icon: Icons.account_balance_wallet, adminOnly: true),
    _TabItem(id: 'approvals', title: 'Approbations', icon: Icons.check_circle, adminOnly: true),
    _TabItem(id: 'services', title: 'Services', icon: Icons.miscellaneous_services, adminOnly: true),
    _TabItem(id: 'payroll', title: 'Paie', icon: Icons.payments, adminOnly: true),
    _TabItem(id: 'settings', title: 'Paramètres', icon: Icons.settings, adminOnly: true),
  ];
  
  List<_TabItem> get _visibleTabs {
    if (_currentEmployee == null) return _allTabs.where((t) => !t.adminOnly).toList();
    if (_currentEmployee!.isAdmin) return _allTabs;
    return _allTabs.where((t) => !t.adminOnly).toList();
  }

  @override
  void initState() {
    super.initState();
    _tabController = TabController(length: 1, vsync: this); // Will be updated
    _loadData();
  }

  @override
  void dispose() {
    _tabController.dispose();
    super.dispose();
  }

  Future<void> _loadData() async {
    setState(() { _isLoading = true; _error = null; });
    
    try {
      // Get enterprises
      final response = await _api.enterprise.getEnterprises();
      List<dynamic> enterprises = [];
      
      if (response is List) {
        enterprises = response;
      } else if (response is Map && response['enterprises'] != null) {
        enterprises = response['enterprises'];
      }
      
      if (enterprises.isEmpty) {
        setState(() { 
          _isLoading = false; 
          _error = 'Aucune entreprise trouvée';
        });
        return;
      }
      
      // Select enterprise
      final enterpriseData = widget.enterpriseId != null
          ? enterprises.firstWhere((e) => e['id'] == widget.enterpriseId || e['_id'] == widget.enterpriseId, orElse: () => enterprises.first)
          : enterprises.first;
      
      _enterprise = Enterprise.fromJson(enterpriseData);
      
      // Get current employee role
      try {
        final empResponse = await _api.enterprise.getMyEmployee(_enterprise!.id);
        _currentEmployee = Employee.fromJson(empResponse);
      } catch (e) {
        // User might be owner without employee record
        _currentEmployee = Employee(
          id: '',
          enterpriseId: _enterprise!.id,
          userId: '',
          role: EmployeeRole.owner,
          status: EmployeeStatus.active,
        );
      }
      
      // Update tab controller
      _tabController.dispose();
      _tabController = TabController(length: _visibleTabs.length, vsync: this);
      
      setState(() { _isLoading = false; });
    } catch (e) {
      setState(() { 
        _isLoading = false; 
        _error = 'Erreur de chargement: ${e.toString()}';
      });
    }
  }

  Widget _buildTabContent(_TabItem tab) {
    if (_enterprise == null) return const SizedBox();
    
    switch (tab.id) {
      case 'overview':
        return OverviewTab(enterprise: _enterprise!, employee: _currentEmployee);
      case 'employees':
        return EmployeesTab(enterprise: _enterprise!, onRefresh: _loadData);
      case 'wallets':
        return WalletsTab(enterprise: _enterprise!, onRefresh: _loadData);
      case 'approvals':
        return ApprovalsTab(enterprise: _enterprise!, currentEmployee: _currentEmployee);
      case 'services':
        return ServicesTab(enterprise: _enterprise!, onRefresh: _loadData);
      case 'payroll':
        return PayrollTab(enterprise: _enterprise!);
      case 'settings':
        return SettingsTab(enterprise: _enterprise!, onRefresh: _loadData);
      default:
        return Center(child: Text('Onglet ${tab.title}'));
    }
  }

  @override
  Widget build(BuildContext context) {
    if (_isLoading) {
      return Scaffold(
        appBar: AppBar(title: const Text('Entreprise')),
        body: const Center(child: CircularProgressIndicator()),
      );
    }
    
    if (_error != null) {
      return Scaffold(
        appBar: AppBar(title: const Text('Entreprise')),
        body: Center(
          child: Column(
            mainAxisAlignment: MainAxisAlignment.center,
            children: [
              Icon(Icons.error_outline, size: 64, color: Colors.grey[400]),
              const SizedBox(height: 16),
              Text(_error!, style: TextStyle(color: Colors.grey[600])),
              const SizedBox(height: 16),
              ElevatedButton(
                onPressed: _loadData,
                child: const Text('Réessayer'),
              ),
            ],
          ),
        ),
      );
    }
    
    return Scaffold(
      appBar: AppBar(
        title: Row(
          children: [
            if (_enterprise?.logo != null)
              ClipRRect(
                borderRadius: BorderRadius.circular(8),
                child: Image.network(
                  _enterprise!.logo!,
                  width: 32,
                  height: 32,
                  fit: BoxFit.cover,
                  errorBuilder: (_, __, ___) => _buildLogoPlaceholder(),
                ),
              )
            else
              _buildLogoPlaceholder(),
            const SizedBox(width: 12),
            Expanded(
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  Text(
                    _enterprise?.name ?? 'Entreprise',
                    style: const TextStyle(fontSize: 16, fontWeight: FontWeight.bold),
                    overflow: TextOverflow.ellipsis,
                  ),
                  if (_currentEmployee != null)
                    Text(
                      _currentEmployee!.roleLabel,
                      style: TextStyle(fontSize: 12, color: Colors.grey[400]),
                    ),
                ],
              ),
            ),
          ],
        ),
        bottom: TabBar(
          controller: _tabController,
          isScrollable: true,
          tabs: _visibleTabs.map((tab) => Tab(
            icon: Icon(tab.icon, size: 20),
            text: tab.title,
          )).toList(),
        ),
      ),
      body: TabBarView(
        controller: _tabController,
        children: _visibleTabs.map(_buildTabContent).toList(),
      ),
    );
  }
  
  Widget _buildLogoPlaceholder() {
    return Container(
      width: 32,
      height: 32,
      decoration: BoxDecoration(
        gradient: LinearGradient(
          colors: [Colors.blue.shade600, Colors.blue.shade800],
        ),
        borderRadius: BorderRadius.circular(8),
      ),
      child: Center(
        child: Text(
          _enterprise?.name.substring(0, 1).toUpperCase() ?? 'E',
          style: const TextStyle(color: Colors.white, fontWeight: FontWeight.bold),
        ),
      ),
    );
  }
}

class _TabItem {
  final String id;
  final String title;
  final IconData icon;
  final bool adminOnly;
  
  const _TabItem({
    required this.id,
    required this.title,
    required this.icon,
    required this.adminOnly,
  });
}
