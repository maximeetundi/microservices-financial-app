import 'package:flutter/material.dart';
import '../../../../../core/services/api_service.dart';
import '../../../data/models/enterprise_model.dart';
import '../../../data/models/employee_model.dart';
import '../../../data/models/approval_model.dart';
import '../approval_detail_page.dart';

class ApprovalsTab extends StatefulWidget {
  final Enterprise enterprise;
  final Employee? currentEmployee;

  const ApprovalsTab({Key? key, required this.enterprise, this.currentEmployee}) : super(key: key);

  @override
  State<ApprovalsTab> createState() => _ApprovalsTabState();
}

class _ApprovalsTabState extends State<ApprovalsTab> {
  final ApiService _api = ApiService();
  bool _isLoading = true;
  List<ActionApproval> _approvals = [];

  @override
  void initState() {
    super.initState();
    _loadApprovals();
  }

  Future<void> _loadApprovals() async {
    setState(() => _isLoading = true);
    try {
      final response = await _api.enterprise.getPendingApprovals(widget.enterprise.id);
      List<dynamic> list = [];
      if (response is List) {
        list = response;
      } else if (response is Map && response['approvals'] != null) {
        list = response['approvals'];
      }
      _approvals = list.map((e) => ActionApproval.fromJson(e)).toList();
    } catch (e) {
      debugPrint('Error loading approvals: $e');
    } finally {
      setState(() => _isLoading = false);
    }
  }

  void _openApprovalDetail(ActionApproval approval) async {
    final result = await Navigator.push(
      context,
      MaterialPageRoute(
        builder: (context) => ApprovalDetailPage(
          approval: approval,
          currentEmployee: widget.currentEmployee,
        ),
      ),
    );
    if (result == true) _loadApprovals();
  }

  @override
  Widget build(BuildContext context) {
    return RefreshIndicator(
      onRefresh: _loadApprovals,
      child: _isLoading
          ? const Center(child: CircularProgressIndicator())
          : _approvals.isEmpty
              ? _buildEmptyState()
              : ListView.builder(
                  padding: const EdgeInsets.all(16),
                  itemCount: _approvals.length,
                  itemBuilder: (context, index) => _ApprovalCard(
                    approval: _approvals[index],
                    onTap: () => _openApprovalDetail(_approvals[index]),
                  ),
                ),
    );
  }

  Widget _buildEmptyState() {
    return Center(
      child: Column(
        mainAxisAlignment: MainAxisAlignment.center,
        children: [
          Icon(Icons.check_circle_outline, size: 80, color: Colors.green[200]),
          const SizedBox(height: 16),
          const Text(
            'Aucune approbation en attente',
            style: TextStyle(fontSize: 16, fontWeight: FontWeight.w500),
          ),
          const SizedBox(height: 8),
          Text(
            'Toutes les demandes ont été traitées',
            style: TextStyle(color: Colors.grey[600]),
          ),
        ],
      ),
    );
  }
}

class _ApprovalCard extends StatelessWidget {
  final ActionApproval approval;
  final VoidCallback onTap;

  const _ApprovalCard({required this.approval, required this.onTap});

  Color get _statusColor {
    switch (approval.status) {
      case ApprovalStatus.pending: return Colors.orange;
      case ApprovalStatus.approved: return Colors.green;
      case ApprovalStatus.rejected: return Colors.red;
      case ApprovalStatus.executed: return Colors.blue;
      default: return Colors.grey;
    }
  }

  IconData get _typeIcon {
    switch (approval.actionType) {
      case ActionType.transaction: return Icons.send;
      case ActionType.payroll: return Icons.payments;
      case ActionType.employeeUpdate: return Icons.person_add;
      default: return Icons.pending_actions;
    }
  }

  @override
  Widget build(BuildContext context) {
    return Card(
      margin: const EdgeInsets.only(bottom: 12),
      shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(16)),
      child: InkWell(
        onTap: onTap,
        borderRadius: BorderRadius.circular(16),
        child: Padding(
          padding: const EdgeInsets.all(16),
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              Row(
                children: [
                  Container(
                    padding: const EdgeInsets.all(10),
                    decoration: BoxDecoration(
                      color: Colors.blue.shade50,
                      borderRadius: BorderRadius.circular(12),
                    ),
                    child: Icon(_typeIcon, color: Colors.blue.shade700),
                  ),
                  const SizedBox(width: 12),
                  Expanded(
                    child: Column(
                      crossAxisAlignment: CrossAxisAlignment.start,
                      children: [
                        Text(
                          approval.actionName,
                          style: const TextStyle(fontWeight: FontWeight.bold, fontSize: 15),
                        ),
                        const SizedBox(height: 2),
                        Text(
                          approval.description,
                          style: TextStyle(color: Colors.grey[600], fontSize: 13),
                          maxLines: 1,
                          overflow: TextOverflow.ellipsis,
                        ),
                      ],
                    ),
                  ),
                  Container(
                    padding: const EdgeInsets.symmetric(horizontal: 10, vertical: 4),
                    decoration: BoxDecoration(
                      color: _statusColor.withOpacity(0.1),
                      borderRadius: BorderRadius.circular(12),
                    ),
                    child: Text(
                      approval.statusLabel,
                      style: TextStyle(
                        color: _statusColor,
                        fontSize: 12,
                        fontWeight: FontWeight.w500,
                      ),
                    ),
                  ),
                ],
              ),
              
              if (approval.amount != null) ...[
                const SizedBox(height: 12),
                Container(
                  padding: const EdgeInsets.all(12),
                  decoration: BoxDecoration(
                    color: Colors.grey[50],
                    borderRadius: BorderRadius.circular(8),
                  ),
                  child: Row(
                    mainAxisAlignment: MainAxisAlignment.spaceBetween,
                    children: [
                      const Text('Montant', style: TextStyle(color: Colors.grey)),
                      Text(
                        '${approval.amount!.toStringAsFixed(0)} ${approval.currency ?? 'XOF'}',
                        style: const TextStyle(fontWeight: FontWeight.bold, fontSize: 16),
                      ),
                    ],
                  ),
                ),
              ],
              
              const SizedBox(height: 12),
              
              // Progress
              Row(
                children: [
                  Expanded(
                    child: Column(
                      crossAxisAlignment: CrossAxisAlignment.start,
                      children: [
                        Text(
                          '${approval.approvedCount}/${approval.requiredApprovals} approbations',
                          style: TextStyle(color: Colors.grey[600], fontSize: 12),
                        ),
                        const SizedBox(height: 4),
                        ClipRRect(
                          borderRadius: BorderRadius.circular(4),
                          child: LinearProgressIndicator(
                            value: approval.progress,
                            backgroundColor: Colors.grey[200],
                            valueColor: AlwaysStoppedAnimation(
                              approval.progress >= 1 ? Colors.green : Colors.blue,
                            ),
                          ),
                        ),
                      ],
                    ),
                  ),
                  const SizedBox(width: 16),
                  Icon(Icons.chevron_right, color: Colors.grey[400]),
                ],
              ),
              
              // Initiator
              if (approval.initiatorName != null) ...[
                const SizedBox(height: 8),
                Text(
                  'Par ${approval.initiatorName}',
                  style: TextStyle(color: Colors.grey[500], fontSize: 11),
                ),
              ],
            ],
          ),
        ),
      ),
    );
  }
}
