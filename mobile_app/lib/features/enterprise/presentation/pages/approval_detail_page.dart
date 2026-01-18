import 'package:flutter/material.dart';
import '../../../../core/services/api_service.dart';
import '../../data/models/approval_model.dart';
import '../../data/models/employee_model.dart';
import '../widgets/pin_verification_dialog.dart';

class ApprovalDetailPage extends StatefulWidget {
  final ActionApproval approval;
  final Employee? currentEmployee;

  const ApprovalDetailPage({Key? key, required this.approval, this.currentEmployee}) : super(key: key);

  @override
  State<ApprovalDetailPage> createState() => _ApprovalDetailPageState();
}

class _ApprovalDetailPageState extends State<ApprovalDetailPage> {
  final ApiService _api = ApiService();
  bool _isLoading = false;
  late ActionApproval _approval;

  @override
  void initState() {
    super.initState();
    _approval = widget.approval;
    _refreshApproval();
  }

  Future<void> _refreshApproval() async {
    try {
      final response = await _api.enterprise.getApprovalById(_approval.id);
      setState(() => _approval = ActionApproval.fromJson(response));
    } catch (e) {
      debugPrint('Error refreshing approval: $e');
    }
  }

  void _approveAction() async {
    final encryptedPin = await showDialog<String>(
      context: context,
      barrierDismissible: false,
      builder: (context) => const PinVerificationDialog(
        title: 'Approuver la transaction',
        description: 'Entrez votre code PIN pour approuver cette action.',
      ),
    );
    
    if (encryptedPin != null) {
      _executeApproval(encryptedPin);
    }
  }

  Future<void> _executeApproval(String encryptedPin) async {
    setState(() => _isLoading = true);
    try {
      await _api.enterprise.approveAction(_approval.id, encryptedPin);
      ScaffoldMessenger.of(context).showSnackBar(
        const SnackBar(content: Text('Action approuvée!'), backgroundColor: Colors.green),
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

  void _rejectAction() async {
    final reason = await showDialog<String>(
      context: context,
      builder: (context) => _RejectReasonDialog(),
    );
    
    if (reason != null) {
      final encryptedPin = await showDialog<String>(
        context: context,
        barrierDismissible: false,
        builder: (context) => const PinVerificationDialog(
          title: 'Rejeter la transaction',
          description: 'Entrez votre code PIN pour rejeter cette action.',
        ),
      );
      
      if (encryptedPin != null) {
        _executeRejection(encryptedPin, reason);
      }
    }
  }

  Future<void> _executeRejection(String encryptedPin, String reason) async {
    setState(() => _isLoading = true);
    try {
      await _api.enterprise.rejectAction(_approval.id, encryptedPin, reason: reason);
      ScaffoldMessenger.of(context).showSnackBar(
        const SnackBar(content: Text('Action rejetée'), backgroundColor: Colors.orange),
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

  Color get _statusColor {
    switch (_approval.status) {
      case ApprovalStatus.pending: return Colors.orange;
      case ApprovalStatus.approved: return Colors.green;
      case ApprovalStatus.rejected: return Colors.red;
      case ApprovalStatus.executed: return Colors.blue;
      default: return Colors.grey;
    }
  }

  bool get _canVote {
    if (!_approval.isPending) return false;
    if (widget.currentEmployee == null) return false;
    // Check if user already voted
    // Note: would need current user ID to check properly
    return true;
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: const Text('Détails de l\'approbation'),
      ),
      body: RefreshIndicator(
        onRefresh: _refreshApproval,
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
                    colors: [_statusColor.withOpacity(0.8), _statusColor],
                  ),
                  borderRadius: BorderRadius.circular(16),
                ),
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    Row(
                      children: [
                        Container(
                          padding: const EdgeInsets.symmetric(horizontal: 12, vertical: 6),
                          decoration: BoxDecoration(
                            color: Colors.white.withOpacity(0.2),
                            borderRadius: BorderRadius.circular(20),
                          ),
                          child: Text(
                            _approval.statusLabel,
                            style: const TextStyle(color: Colors.white, fontWeight: FontWeight.bold),
                          ),
                        ),
                      ],
                    ),
                    const SizedBox(height: 16),
                    Text(
                      _approval.actionName,
                      style: const TextStyle(
                        color: Colors.white,
                        fontSize: 22,
                        fontWeight: FontWeight.bold,
                      ),
                    ),
                    const SizedBox(height: 8),
                    Text(
                      _approval.description,
                      style: const TextStyle(color: Colors.white70),
                    ),
                  ],
                ),
              ),
              
              const SizedBox(height: 24),
              
              // Amount (if exists)
              if (_approval.amount != null) ...[
                Container(
                  padding: const EdgeInsets.all(20),
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
                  child: Row(
                    mainAxisAlignment: MainAxisAlignment.spaceBetween,
                    children: [
                      const Text('Montant', style: TextStyle(color: Colors.grey)),
                      Text(
                        '${_approval.amount!.toStringAsFixed(0)} ${_approval.currency ?? 'XOF'}',
                        style: const TextStyle(
                          fontSize: 24,
                          fontWeight: FontWeight.bold,
                        ),
                      ),
                    ],
                  ),
                ),
                const SizedBox(height: 16),
              ],
              
              // Progress Card
              Container(
                padding: const EdgeInsets.all(20),
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
                    const Text('Progression', style: TextStyle(fontWeight: FontWeight.bold, fontSize: 16)),
                    const SizedBox(height: 16),
                    Row(
                      children: [
                        Text(
                          '${_approval.approvedCount}',
                          style: const TextStyle(fontSize: 32, fontWeight: FontWeight.bold, color: Colors.green),
                        ),
                        Text(
                          ' / ${_approval.requiredApprovals}',
                          style: const TextStyle(fontSize: 24, color: Colors.grey),
                        ),
                        const Spacer(),
                        Text(
                          '${(_approval.progress * 100).toInt()}%',
                          style: TextStyle(fontSize: 18, fontWeight: FontWeight.bold, color: _statusColor),
                        ),
                      ],
                    ),
                    const SizedBox(height: 12),
                    ClipRRect(
                      borderRadius: BorderRadius.circular(4),
                      child: LinearProgressIndicator(
                        value: _approval.progress,
                        minHeight: 8,
                        backgroundColor: Colors.grey[200],
                        valueColor: AlwaysStoppedAnimation(_statusColor),
                      ),
                    ),
                  ],
                ),
              ),
              
              const SizedBox(height: 16),
              
              // Votes List
              if (_approval.approvals.isNotEmpty) ...[
                Container(
                  padding: const EdgeInsets.all(20),
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
                      const Text('Votes', style: TextStyle(fontWeight: FontWeight.bold, fontSize: 16)),
                      const SizedBox(height: 12),
                      ..._approval.approvals.map((vote) => Padding(
                        padding: const EdgeInsets.symmetric(vertical: 8),
                        child: Row(
                          children: [
                            CircleAvatar(
                              radius: 18,
                              backgroundColor: vote.isApproved ? Colors.green : Colors.red,
                              child: Icon(
                                vote.isApproved ? Icons.check : Icons.close,
                                color: Colors.white,
                                size: 18,
                              ),
                            ),
                            const SizedBox(width: 12),
                            Expanded(
                              child: Column(
                                crossAxisAlignment: CrossAxisAlignment.start,
                                children: [
                                  Text(
                                    vote.adminName ?? 'Admin',
                                    style: const TextStyle(fontWeight: FontWeight.w500),
                                  ),
                                  if (vote.reason != null)
                                    Text(vote.reason!, style: TextStyle(color: Colors.grey[600], fontSize: 12)),
                                ],
                              ),
                            ),
                            Text(
                              vote.isApproved ? 'Approuvé' : 'Rejeté',
                              style: TextStyle(
                                color: vote.isApproved ? Colors.green : Colors.red,
                                fontWeight: FontWeight.w500,
                              ),
                            ),
                          ],
                        ),
                      )).toList(),
                    ],
                  ),
                ),
                const SizedBox(height: 16),
              ],
              
              // Initiator Info
              Container(
                padding: const EdgeInsets.all(16),
                decoration: BoxDecoration(
                  color: Colors.grey[50],
                  borderRadius: BorderRadius.circular(12),
                ),
                child: Row(
                  children: [
                    const Icon(Icons.person_outline, color: Colors.grey),
                    const SizedBox(width: 12),
                    Expanded(
                      child: Text(
                        'Initié par ${_approval.initiatorName ?? 'un administrateur'}',
                        style: TextStyle(color: Colors.grey[600]),
                      ),
                    ),
                  ],
                ),
              ),
              
              const SizedBox(height: 32),
              
              // Action Buttons
              if (_canVote) ...[
                Row(
                  children: [
                    Expanded(
                      child: OutlinedButton(
                        onPressed: _isLoading ? null : _rejectAction,
                        style: OutlinedButton.styleFrom(
                          foregroundColor: Colors.red,
                          side: const BorderSide(color: Colors.red),
                          padding: const EdgeInsets.symmetric(vertical: 16),
                          shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(12)),
                        ),
                        child: const Text('Rejeter'),
                      ),
                    ),
                    const SizedBox(width: 16),
                    Expanded(
                      child: ElevatedButton(
                        onPressed: _isLoading ? null : _approveAction,
                        style: ElevatedButton.styleFrom(
                          backgroundColor: Colors.green,
                          foregroundColor: Colors.white,
                          padding: const EdgeInsets.symmetric(vertical: 16),
                          shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(12)),
                        ),
                        child: _isLoading
                            ? const SizedBox(width: 20, height: 20, child: CircularProgressIndicator(color: Colors.white, strokeWidth: 2))
                            : const Text('Approuver'),
                      ),
                    ),
                  ],
                ),
              ],
            ],
          ),
        ),
      ),
    );
  }
}

class _RejectReasonDialog extends StatefulWidget {
  @override
  State<_RejectReasonDialog> createState() => _RejectReasonDialogState();
}

class _RejectReasonDialogState extends State<_RejectReasonDialog> {
  final _controller = TextEditingController();

  @override
  void dispose() {
    _controller.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    return AlertDialog(
      title: const Text('Raison du rejet'),
      content: TextField(
        controller: _controller,
        maxLines: 3,
        decoration: const InputDecoration(
          hintText: 'Pourquoi rejetez-vous cette action?',
          border: OutlineInputBorder(),
        ),
      ),
      actions: [
        TextButton(
          onPressed: () => Navigator.pop(context),
          child: const Text('Annuler'),
        ),
        ElevatedButton(
          onPressed: () => Navigator.pop(context, _controller.text),
          child: const Text('Continuer'),
        ),
      ],
    );
  }
}
