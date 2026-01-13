import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';
import '../../../../core/services/donation_api_service.dart';
import '../../../auth/presentation/bloc/auth_bloc.dart';

class DonorsListPage extends StatefulWidget {
  final String campaignId;

  const DonorsListPage({super.key, required this.campaignId});

  @override
  State<DonorsListPage> createState() => _DonorsListPageState();
}

class _DonorsListPageState extends State<DonorsListPage> {
  final DonationApiService _api = DonationApiService();
  List<dynamic> _donations = [];
  Map<String, dynamic>? _campaign;
  bool _loading = true;
  bool _isCreator = false;

  @override
  void initState() {
    super.initState();
    _loadData();
  }

  Future<void> _loadData() async {
    setState(() => _loading = true);
    try {
      final results = await Future.wait([
        _api.getCampaign(widget.campaignId),
        _api.getDonations(widget.campaignId, limit: 100),
      ]);
      
      _campaign = results[0] as Map<String, dynamic>;
      _donations = results[1] as List<dynamic>;

      // Check Creator Status
      final authState = context.read<AuthBloc>().state;
      if (authState is AuthenticatedState) {
        _isCreator = _campaign!['creator_id'] == authState.user.id;
      }

    } catch (e) {
      debugPrint('Error: $e');
    } finally {
      if (mounted) setState(() => _loading = false);
    }
  }

  Future<void> _refundDonation(String donationId) async {
    // Confirm Dialog
    final confirm = await showDialog<bool>(
      context: context,
      builder: (ctx) => AlertDialog(
        title: const Text('Rembourser ce don ?'),
        content: const Text('Cette action est irréversible.'),
        actions: [
          TextButton(onPressed: () => Navigator.pop(ctx, false), child: const Text('Annuler')),
          TextButton(onPressed: () => Navigator.pop(ctx, true), child: const Text('Rembourser', style: TextStyle(color: Colors.red))),
        ],
      ),
    );

    if (confirm != true) return;

    try {
      await _api.refundDonation(donationId);
      if (mounted) {
        ScaffoldMessenger.of(context).showSnackBar(const SnackBar(content: Text('Remboursement initié')));
        _loadData(); // Reload
      }
    } catch (e) {
      if (mounted) ScaffoldMessenger.of(context).showSnackBar(SnackBar(content: Text('Erreur: $e')));
    }
  }

  @override
  Widget build(BuildContext context) {
    if (_loading) {
        return const Scaffold(
            backgroundColor: Color(0xFF1a1a2e),
            body: Center(child: CircularProgressIndicator()),
        );
    }

    final bool showDonors = _isCreator || (_campaign?['show_donors'] == true);

    return Scaffold(
      backgroundColor: const Color(0xFF1a1a2e),
      appBar: AppBar(
        title: const Text('Donateurs'),
        backgroundColor: Colors.transparent,
        elevation: 0,
      ),
      body: !showDonors 
          ? Center(
              child: Column(
                mainAxisAlignment: MainAxisAlignment.center,
                children: const [
                  Icon(Icons.lock, size: 64, color: Colors.grey),
                  SizedBox(height: 16),
                  Text('Liste privée', style: TextStyle(color: Colors.white, fontSize: 18, fontWeight: FontWeight.bold)),
                  SizedBox(height: 8),
                  Text('Seul l\'organisateur peut voir les donateurs.', style: TextStyle(color: Colors.grey)),
                ],
              ),
            )
          : _donations.isEmpty 
              ? const Center(child: Text('Aucun don pour le moment', style: TextStyle(color: Colors.grey)))
              : ListView.builder(
                  padding: const EdgeInsets.all(16),
                  itemCount: _donations.length,
                  itemBuilder: (context, index) {
                    return _buildDonorRow(_donations[index]);
                  },
                ),
    );
  }

  Widget _buildDonorRow(dynamic donation) {
    final isAnon = donation['is_anonymous'] == true;
    final name = isAnon ? 'Donateur Anonyme' : (donation['donor_name'] ?? 'Bienfaiteur');
    final status = donation['status'];

    return Container(
      margin: const EdgeInsets.only(bottom: 12),
      padding: const EdgeInsets.all(12),
      decoration: BoxDecoration(
        color: Colors.white.withOpacity(0.05),
        borderRadius: BorderRadius.circular(12),
        border: Border.all(color: Colors.white.withOpacity(0.1)),
      ),
      child: Column(
        children: [
          Row(
            children: [
              CircleAvatar(
                backgroundColor: isAnon ? Colors.grey[800] : Colors.indigo[900],
                child: Text(isAnon ? '?' : name[0].toUpperCase(), style: const TextStyle(color: Colors.white)),
              ),
              const SizedBox(width: 12),
              Expanded(
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    Text(name, style: const TextStyle(color: Colors.white, fontWeight: FontWeight.bold)),
                    if (donation['message'] != null && donation['message'].isNotEmpty)
                       Text('"${donation['message']}"', style: TextStyle(color: Colors.white.withOpacity(0.6), fontSize: 12, fontStyle: FontStyle.italic)),
                     Text(
                        _formatDate(donation['created_at']), 
                        style: TextStyle(color: Colors.white.withOpacity(0.4), fontSize: 10),
                     ),
                  ],
                ),
              ),
              Column(
                crossAxisAlignment: CrossAxisAlignment.end,
                children: [
                    Text(
                        _formatAmount((donation['amount'] ?? 0).toDouble(), donation['currency']),
                        style: const TextStyle(color: Color(0xFF6366f1), fontWeight: FontWeight.bold),
                    ),
                    if (status == 'refunded' || status == 'refunding')
                        const Text('Remboursé', style: TextStyle(color: Colors.orange, fontSize: 10)),
                ],
              ),
            ],
          ),
          
          // Refund Action
          if (_isCreator && status == 'paid') ...[
              const Divider(color: Colors.white12),
              Align(
                  alignment: Alignment.centerRight,
                  child: TextButton.icon(
                      icon: const Icon(Icons.undo, size: 16, color: Colors.red),
                      label: const Text('Rembourser', style: TextStyle(color: Colors.red, fontSize: 12)),
                      onPressed: () => _refundDonation(donation['id']),
                  ),
              )
          ]
        ],
      ),
    );
  }

  String _formatAmount(double amount, String? currency) {
    return '${amount.toStringAsFixed(0)} ${currency ?? "XOF"}';
  }

  String _formatDate(String? dateStr) {
      if (dateStr == null) return '';
      try {
          final date = DateTime.parse(dateStr);
          return '${date.day}/${date.month}/${date.year}';
      } catch (e) {
          return '';
      }
  }
}
