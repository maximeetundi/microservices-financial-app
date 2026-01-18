import 'package:flutter/material.dart';
import '../../../data/models/enterprise_model.dart';
import '../../../data/models/employee_model.dart';

class OverviewTab extends StatelessWidget {
  final Enterprise enterprise;
  final Employee? employee;

  const OverviewTab({Key? key, required this.enterprise, this.employee}) : super(key: key);

  @override
  Widget build(BuildContext context) {
    return SingleChildScrollView(
      padding: const EdgeInsets.all(16),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          // Welcome Card
          Container(
            width: double.infinity,
            padding: const EdgeInsets.all(20),
            decoration: BoxDecoration(
              gradient: LinearGradient(
                colors: [Colors.blue.shade700, Colors.blue.shade500],
                begin: Alignment.topLeft,
                end: Alignment.bottomRight,
              ),
              borderRadius: BorderRadius.circular(16),
              boxShadow: [
                BoxShadow(
                  color: Colors.blue.withOpacity(0.3),
                  blurRadius: 12,
                  offset: const Offset(0, 6),
                ),
              ],
            ),
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                const Text(
                  'Bienvenue',
                  style: TextStyle(color: Colors.white70, fontSize: 14),
                ),
                const SizedBox(height: 8),
                Text(
                  enterprise.name,
                  style: const TextStyle(
                    color: Colors.white,
                    fontWeight: FontWeight.bold,
                    fontSize: 24,
                  ),
                ),
                if (employee != null) ...[
                  const SizedBox(height: 8),
                  Container(
                    padding: const EdgeInsets.symmetric(horizontal: 12, vertical: 6),
                    decoration: BoxDecoration(
                      color: Colors.white.withOpacity(0.2),
                      borderRadius: BorderRadius.circular(20),
                    ),
                    child: Text(
                      employee!.roleLabel,
                      style: const TextStyle(color: Colors.white, fontSize: 12),
                    ),
                  ),
                ],
              ],
            ),
          ),
          
          const SizedBox(height: 24),
          
          // Stats
          const Text(
            'Statistiques',
            style: TextStyle(fontSize: 18, fontWeight: FontWeight.bold),
          ),
          const SizedBox(height: 12),
          
          Row(
            children: [
              Expanded(child: _StatCard(
                icon: Icons.people,
                label: 'Employés',
                value: '${enterprise.serviceGroups.length}',
                color: Colors.blue,
              )),
              const SizedBox(width: 12),
              Expanded(child: _StatCard(
                icon: Icons.account_balance_wallet,
                label: 'Portefeuilles',
                value: '${enterprise.walletIds.length + (enterprise.defaultWalletId != null ? 1 : 0)}',
                color: Colors.green,
              )),
            ],
          ),
          
          const SizedBox(height: 12),
          
          Row(
            children: [
              Expanded(child: _StatCard(
                icon: Icons.miscellaneous_services,
                label: 'Services',
                value: '${enterprise.serviceGroups.fold<int>(0, (sum, g) => sum + g.services.length)}',
                color: Colors.orange,
              )),
              const SizedBox(width: 12),
              Expanded(child: _StatCard(
                icon: Icons.trending_up,
                label: 'Statut',
                value: enterprise.status,
                color: enterprise.status == 'ACTIVE' ? Colors.teal : Colors.amber,
              )),
            ],
          ),
          
          const SizedBox(height: 24),
          
          // Quick Actions
          const Text(
            'Actions rapides',
            style: TextStyle(fontSize: 18, fontWeight: FontWeight.bold),
          ),
          const SizedBox(height: 12),
          
          GridView.count(
            shrinkWrap: true,
            physics: const NeverScrollableScrollPhysics(),
            crossAxisCount: 3,
            mainAxisSpacing: 12,
            crossAxisSpacing: 12,
            childAspectRatio: 1,
            children: [
              _QuickAction(
                icon: Icons.person_add,
                label: 'Inviter',
                onTap: () {},
              ),
              _QuickAction(
                icon: Icons.send,
                label: 'Envoyer',
                onTap: () {},
              ),
              _QuickAction(
                icon: Icons.qr_code,
                label: 'QR Code',
                onTap: () {},
              ),
              _QuickAction(
                icon: Icons.receipt,
                label: 'Facturer',
                onTap: () {},
              ),
              _QuickAction(
                icon: Icons.payments,
                label: 'Paie',
                onTap: () {},
              ),
              _QuickAction(
                icon: Icons.settings,
                label: 'Paramètres',
                onTap: () {},
              ),
            ],
          ),
          
          if (enterprise.description != null) ...[
            const SizedBox(height: 24),
            const Text(
              'Description',
              style: TextStyle(fontSize: 18, fontWeight: FontWeight.bold),
            ),
            const SizedBox(height: 8),
            Text(
              enterprise.description!,
              style: TextStyle(color: Colors.grey[600]),
            ),
          ],
        ],
      ),
    );
  }
}

class _StatCard extends StatelessWidget {
  final IconData icon;
  final String label;
  final String value;
  final Color color;

  const _StatCard({
    required this.icon,
    required this.label,
    required this.value,
    required this.color,
  });

  @override
  Widget build(BuildContext context) {
    return Container(
      padding: const EdgeInsets.all(16),
      decoration: BoxDecoration(
        color: Colors.white,
        borderRadius: BorderRadius.circular(12),
        boxShadow: [
          BoxShadow(
            color: Colors.grey.withOpacity(0.1),
            blurRadius: 8,
            offset: const Offset(0, 2),
          ),
        ],
      ),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Container(
            padding: const EdgeInsets.all(8),
            decoration: BoxDecoration(
              color: color.withOpacity(0.1),
              borderRadius: BorderRadius.circular(8),
            ),
            child: Icon(icon, color: color, size: 20),
          ),
          const SizedBox(height: 12),
          Text(
            value,
            style: const TextStyle(fontSize: 20, fontWeight: FontWeight.bold),
          ),
          Text(
            label,
            style: TextStyle(color: Colors.grey[600], fontSize: 12),
          ),
        ],
      ),
    );
  }
}

class _QuickAction extends StatelessWidget {
  final IconData icon;
  final String label;
  final VoidCallback onTap;

  const _QuickAction({
    required this.icon,
    required this.label,
    required this.onTap,
  });

  @override
  Widget build(BuildContext context) {
    return InkWell(
      onTap: onTap,
      borderRadius: BorderRadius.circular(12),
      child: Container(
        decoration: BoxDecoration(
          color: Colors.grey[100],
          borderRadius: BorderRadius.circular(12),
        ),
        child: Column(
          mainAxisAlignment: MainAxisAlignment.center,
          children: [
            Icon(icon, color: Colors.blue.shade700, size: 28),
            const SizedBox(height: 6),
            Text(
              label,
              style: const TextStyle(fontSize: 11, fontWeight: FontWeight.w500),
              textAlign: TextAlign.center,
            ),
          ],
        ),
      ),
    );
  }
}
