import 'package:flutter/material.dart';
import '../../../data/models/enterprise_model.dart';

class ServicesTab extends StatefulWidget {
  final Enterprise enterprise;
  final VoidCallback onRefresh;

  const ServicesTab({Key? key, required this.enterprise, required this.onRefresh}) : super(key: key);

  @override
  State<ServicesTab> createState() => _ServicesTabState();
}

class _ServicesTabState extends State<ServicesTab> {
  @override
  Widget build(BuildContext context) {
    final groups = widget.enterprise.serviceGroups;
    
    return SingleChildScrollView(
      padding: const EdgeInsets.all(16),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          // Header
          Row(
            mainAxisAlignment: MainAxisAlignment.spaceBetween,
            children: [
              const Text(
                'Services',
                style: TextStyle(fontSize: 20, fontWeight: FontWeight.bold),
              ),
              ElevatedButton.icon(
                onPressed: _addServiceGroup,
                icon: const Icon(Icons.add, size: 18),
                label: const Text('Nouveau groupe'),
                style: ElevatedButton.styleFrom(
                  backgroundColor: Colors.blue,
                  foregroundColor: Colors.white,
                  shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(12)),
                ),
              ),
            ],
          ),
          
          const SizedBox(height: 16),
          
          if (groups.isEmpty)
            _buildEmptyState()
          else
            ...groups.map((group) => _ServiceGroupCard(
              group: group,
              onEdit: () => _editGroup(group),
              onAddService: () => _addService(group),
            )).toList(),
        ],
      ),
    );
  }

  Widget _buildEmptyState() {
    return Container(
      padding: const EdgeInsets.all(32),
      child: Column(
        children: [
          Icon(Icons.miscellaneous_services, size: 64, color: Colors.grey[300]),
          const SizedBox(height: 16),
          Text('Aucun service configuré', style: TextStyle(color: Colors.grey[600])),
          const SizedBox(height: 8),
          TextButton(
            onPressed: _addServiceGroup,
            child: const Text('Créer un groupe de services'),
          ),
        ],
      ),
    );
  }

  void _addServiceGroup() {
    // TODO: Implement add service group
    ScaffoldMessenger.of(context).showSnackBar(
      const SnackBar(content: Text('Fonctionnalité en cours de développement')),
    );
  }

  void _editGroup(ServiceGroup group) {
    // TODO: Implement edit group
  }

  void _addService(ServiceGroup group) {
    // TODO: Implement add service
  }
}

class _ServiceGroupCard extends StatelessWidget {
  final ServiceGroup group;
  final VoidCallback onEdit;
  final VoidCallback onAddService;

  const _ServiceGroupCard({
    required this.group,
    required this.onEdit,
    required this.onAddService,
  });

  @override
  Widget build(BuildContext context) {
    return Container(
      margin: const EdgeInsets.only(bottom: 16),
      decoration: BoxDecoration(
        color: Colors.white,
        borderRadius: BorderRadius.circular(16),
        boxShadow: [
          BoxShadow(
            color: Colors.grey.withOpacity(0.1),
            blurRadius: 10,
            offset: const Offset(0, 4),
          ),
        ],
      ),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          // Header
          Container(
            padding: const EdgeInsets.all(16),
            decoration: BoxDecoration(
              color: Colors.blue.shade50,
              borderRadius: const BorderRadius.only(
                topLeft: Radius.circular(16),
                topRight: Radius.circular(16),
              ),
            ),
            child: Row(
              children: [
                Container(
                  padding: const EdgeInsets.all(8),
                  decoration: BoxDecoration(
                    color: Colors.white,
                    borderRadius: BorderRadius.circular(8),
                  ),
                  child: Icon(Icons.folder, color: Colors.blue.shade700),
                ),
                const SizedBox(width: 12),
                Expanded(
                  child: Column(
                    crossAxisAlignment: CrossAxisAlignment.start,
                    children: [
                      Text(
                        group.name,
                        style: const TextStyle(fontWeight: FontWeight.bold, fontSize: 16),
                      ),
                      Text(
                        '${group.services.length} service(s) • ${group.currency}',
                        style: TextStyle(color: Colors.grey[600], fontSize: 12),
                      ),
                    ],
                  ),
                ),
                if (group.isPrivate)
                  Container(
                    padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 4),
                    decoration: BoxDecoration(
                      color: Colors.orange.shade100,
                      borderRadius: BorderRadius.circular(8),
                    ),
                    child: Row(
                      children: [
                        Icon(Icons.lock, size: 12, color: Colors.orange.shade700),
                        const SizedBox(width: 4),
                        Text(
                          'Privé',
                          style: TextStyle(color: Colors.orange.shade700, fontSize: 11),
                        ),
                      ],
                    ),
                  ),
                IconButton(
                  onPressed: onEdit,
                  icon: const Icon(Icons.edit, size: 18),
                ),
              ],
            ),
          ),
          
          // Services List
          if (group.services.isEmpty)
            Padding(
              padding: const EdgeInsets.all(16),
              child: Row(
                children: [
                  Icon(Icons.info_outline, size: 16, color: Colors.grey[400]),
                  const SizedBox(width: 8),
                  Text('Aucun service', style: TextStyle(color: Colors.grey[500])),
                ],
              ),
            )
          else
            ...group.services.map((service) => Container(
              padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 12),
              decoration: BoxDecoration(
                border: Border(bottom: BorderSide(color: Colors.grey.shade100)),
              ),
              child: Row(
                children: [
                  Container(
                    width: 8,
                    height: 8,
                    decoration: const BoxDecoration(
                      color: Colors.blue,
                      shape: BoxShape.circle,
                    ),
                  ),
                  const SizedBox(width: 12),
                  Expanded(
                    child: Column(
                      crossAxisAlignment: CrossAxisAlignment.start,
                      children: [
                        Text(service.name, style: const TextStyle(fontWeight: FontWeight.w500)),
                        if (service.description != null)
                          Text(
                            service.description!,
                            style: TextStyle(color: Colors.grey[500], fontSize: 12),
                          ),
                      ],
                    ),
                  ),
                  Text(
                    '${service.price.toStringAsFixed(0)} ${group.currency}',
                    style: const TextStyle(fontWeight: FontWeight.bold),
                  ),
                ],
              ),
            )).toList(),
          
          // Add Service Button
          Padding(
            padding: const EdgeInsets.all(12),
            child: TextButton.icon(
              onPressed: onAddService,
              icon: const Icon(Icons.add, size: 18),
              label: const Text('Ajouter un service'),
            ),
          ),
        ],
      ),
    );
  }
}
