import 'package:flutter/material.dart';
import '../../../../../core/services/api_service.dart';
import '../../../data/models/enterprise_model.dart';
import '../../../data/models/employee_model.dart';
import '../invite_employee_page.dart';

class EmployeesTab extends StatefulWidget {
  final Enterprise enterprise;
  final VoidCallback onRefresh;

  const EmployeesTab({Key? key, required this.enterprise, required this.onRefresh}) : super(key: key);

  @override
  State<EmployeesTab> createState() => _EmployeesTabState();
}

class _EmployeesTabState extends State<EmployeesTab> {
  final ApiService _api = ApiService();
  bool _isLoading = true;
  List<Employee> _employees = [];

  @override
  void initState() {
    super.initState();
    _loadEmployees();
  }

  Future<void> _loadEmployees() async {
    setState(() => _isLoading = true);
    try {
      final response = await _api.enterprise.getEmployees(widget.enterprise.id);
      List<dynamic> list = [];
      if (response is List) {
        list = response;
      } else if (response is Map && response['employees'] != null) {
        list = response['employees'];
      }
      _employees = list.map((e) => Employee.fromJson(e)).toList();
    } catch (e) {
      ScaffoldMessenger.of(context).showSnackBar(
        SnackBar(content: Text('Erreur: ${e.toString()}')),
      );
    } finally {
      setState(() => _isLoading = false);
    }
  }

  void _inviteEmployee() async {
    final result = await Navigator.push(
      context,
      MaterialPageRoute(
        builder: (context) => InviteEmployeePage(enterprise: widget.enterprise),
      ),
    );
    if (result == true) {
      _loadEmployees();
    }
  }

  @override
  Widget build(BuildContext context) {
    return Column(
      children: [
        // Header
        Padding(
          padding: const EdgeInsets.all(16),
          child: Row(
            mainAxisAlignment: MainAxisAlignment.spaceBetween,
            children: [
              Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  const Text(
                    'Employés',
                    style: TextStyle(fontSize: 20, fontWeight: FontWeight.bold),
                  ),
                  Text(
                    '${_employees.length} membres',
                    style: TextStyle(color: Colors.grey[600]),
                  ),
                ],
              ),
              ElevatedButton.icon(
                onPressed: _inviteEmployee,
                icon: const Icon(Icons.person_add, size: 18),
                label: const Text('Inviter'),
                style: ElevatedButton.styleFrom(
                  backgroundColor: Colors.blue,
                  foregroundColor: Colors.white,
                  shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(12)),
                ),
              ),
            ],
          ),
        ),
        
        // List
        Expanded(
          child: _isLoading
              ? const Center(child: CircularProgressIndicator())
              : _employees.isEmpty
                  ? _buildEmptyState()
                  : RefreshIndicator(
                      onRefresh: _loadEmployees,
                      child: ListView.builder(
                        padding: const EdgeInsets.symmetric(horizontal: 16),
                        itemCount: _employees.length,
                        itemBuilder: (context, index) => _EmployeeCard(
                          employee: _employees[index],
                          onTap: () => _showEmployeeDetails(_employees[index]),
                        ),
                      ),
                    ),
        ),
      ],
    );
  }

  Widget _buildEmptyState() {
    return Center(
      child: Column(
        mainAxisAlignment: MainAxisAlignment.center,
        children: [
          Icon(Icons.people_outline, size: 64, color: Colors.grey[300]),
          const SizedBox(height: 16),
          Text('Aucun employé', style: TextStyle(color: Colors.grey[600])),
          const SizedBox(height: 8),
          TextButton(
            onPressed: _inviteEmployee,
            child: const Text('Inviter un employé'),
          ),
        ],
      ),
    );
  }

  void _showEmployeeDetails(Employee employee) {
    showModalBottomSheet(
      context: context,
      isScrollControlled: true,
      backgroundColor: Colors.transparent,
      builder: (context) => Container(
        height: MediaQuery.of(context).size.height * 0.6,
        decoration: const BoxDecoration(
          color: Colors.white,
          borderRadius: BorderRadius.only(
            topLeft: Radius.circular(20),
            topRight: Radius.circular(20),
          ),
        ),
        child: Column(
          children: [
            Container(
              margin: const EdgeInsets.only(top: 12),
              width: 40,
              height: 4,
              decoration: BoxDecoration(
                color: Colors.grey[300],
                borderRadius: BorderRadius.circular(2),
              ),
            ),
            Padding(
              padding: const EdgeInsets.all(20),
              child: Column(
                children: [
                  CircleAvatar(
                    radius: 40,
                    backgroundColor: Colors.blue.shade100,
                    child: Text(
                      employee.fullName.isNotEmpty 
                          ? employee.fullName[0].toUpperCase() 
                          : 'E',
                      style: TextStyle(fontSize: 32, color: Colors.blue.shade700),
                    ),
                  ),
                  const SizedBox(height: 16),
                  Text(
                    employee.fullName.isNotEmpty ? employee.fullName : 'Employé',
                    style: const TextStyle(fontSize: 20, fontWeight: FontWeight.bold),
                  ),
                  const SizedBox(height: 4),
                  Text(
                    employee.position ?? 'Poste non défini',
                    style: TextStyle(color: Colors.grey[600]),
                  ),
                  const SizedBox(height: 12),
                  Row(
                    mainAxisAlignment: MainAxisAlignment.center,
                    children: [
                      _StatusChip(label: employee.roleLabel, color: Colors.blue),
                      const SizedBox(width: 8),
                      _StatusChip(
                        label: employee.statusLabel,
                        color: employee.isActive ? Colors.green : Colors.orange,
                      ),
                    ],
                  ),
                  const SizedBox(height: 20),
                  
                  // Info rows
                  if (employee.email != null)
                    _InfoRow(icon: Icons.email, label: 'Email', value: employee.email!),
                  if (employee.phone != null)
                    _InfoRow(icon: Icons.phone, label: 'Téléphone', value: employee.phone!),
                  if (employee.department != null)
                    _InfoRow(icon: Icons.business, label: 'Département', value: employee.department!),
                  if (employee.salary != null)
                    _InfoRow(
                      icon: Icons.payments, 
                      label: 'Salaire', 
                      value: '${employee.salary} ${employee.salaryCurrency ?? 'XOF'}',
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

class _EmployeeCard extends StatelessWidget {
  final Employee employee;
  final VoidCallback onTap;

  const _EmployeeCard({required this.employee, required this.onTap});

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
              CircleAvatar(
                backgroundColor: Colors.blue.shade100,
                child: Text(
                  employee.fullName.isNotEmpty 
                      ? employee.fullName[0].toUpperCase() 
                      : 'E',
                  style: TextStyle(color: Colors.blue.shade700, fontWeight: FontWeight.bold),
                ),
              ),
              const SizedBox(width: 12),
              Expanded(
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    Text(
                      employee.fullName.isNotEmpty ? employee.fullName : 'Invitation en attente',
                      style: const TextStyle(fontWeight: FontWeight.w600),
                    ),
                    const SizedBox(height: 2),
                    Text(
                      employee.position ?? employee.email ?? 'Poste non défini',
                      style: TextStyle(color: Colors.grey[600], fontSize: 13),
                    ),
                  ],
                ),
              ),
              Column(
                crossAxisAlignment: CrossAxisAlignment.end,
                children: [
                  Container(
                    padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 4),
                    decoration: BoxDecoration(
                      color: employee.isAdmin ? Colors.blue.shade50 : Colors.grey.shade100,
                      borderRadius: BorderRadius.circular(12),
                    ),
                    child: Text(
                      employee.roleLabel,
                      style: TextStyle(
                        color: employee.isAdmin ? Colors.blue.shade700 : Colors.grey[700],
                        fontSize: 11,
                        fontWeight: FontWeight.w500,
                      ),
                    ),
                  ),
                  const SizedBox(height: 4),
                  Row(
                    mainAxisSize: MainAxisSize.min,
                    children: [
                      Container(
                        width: 8,
                        height: 8,
                        decoration: BoxDecoration(
                          color: employee.isActive ? Colors.green : Colors.orange,
                          shape: BoxShape.circle,
                        ),
                      ),
                      const SizedBox(width: 4),
                      Text(
                        employee.statusLabel,
                        style: TextStyle(color: Colors.grey[500], fontSize: 11),
                      ),
                    ],
                  ),
                ],
              ),
            ],
          ),
        ),
      ),
    );
  }
}

class _StatusChip extends StatelessWidget {
  final String label;
  final Color color;

  const _StatusChip({required this.label, required this.color});

  @override
  Widget build(BuildContext context) {
    return Container(
      padding: const EdgeInsets.symmetric(horizontal: 12, vertical: 6),
      decoration: BoxDecoration(
        color: color.withOpacity(0.1),
        borderRadius: BorderRadius.circular(20),
      ),
      child: Text(
        label,
        style: TextStyle(color: color, fontSize: 12, fontWeight: FontWeight.w500),
      ),
    );
  }
}

class _InfoRow extends StatelessWidget {
  final IconData icon;
  final String label;
  final String value;

  const _InfoRow({required this.icon, required this.label, required this.value});

  @override
  Widget build(BuildContext context) {
    return Padding(
      padding: const EdgeInsets.symmetric(vertical: 8),
      child: Row(
        children: [
          Icon(icon, size: 20, color: Colors.grey[600]),
          const SizedBox(width: 12),
          Text(label, style: TextStyle(color: Colors.grey[600])),
          const Spacer(),
          Text(value, style: const TextStyle(fontWeight: FontWeight.w500)),
        ],
      ),
    );
  }
}
