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
  
  // Search
  String _searchQuery = '';
  List<dynamic> get _filteredDonations {
    if (_searchQuery.isEmpty) return _donations;
    final q = _searchQuery.toLowerCase();
    return _donations.where((d) {
      final name = (d['donor_name'] ?? '').toString().toLowerCase();
      final message = (d['message'] ?? '').toString().toLowerCase();
      return name.contains(q) || message.contains(q);
    }).toList();
  }

  // Tier Stats
  List<Map<String, dynamic>> get _tierStats {
    if (_campaign == null || _campaign!['donation_type'] != 'tiers' || _campaign!['donation_tiers'] == null) return [];
    
    final tiers = List<dynamic>.from(_campaign!['donation_tiers']);
    return tiers.map((tier) {
         final amount = double.tryParse(tier['amount'].toString()) ?? 0;
         final matching = _donations.where((d) => ((d['amount'] ?? 0) - amount).abs() < 0.1);
         final total = matching.fold(0.0, (sum, d) => sum + (d['amount'] ?? 0));
         return {
           'label': tier['label'],
           'count': matching.length,
           'total': total,
         };
    }).toList();
  }

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
    final reasonController = TextEditingController();
    
    // Confirm Dialog with Reason
    final confirm = await showDialog<bool>(
      context: context,
      builder: (ctx) => AlertDialog(
        backgroundColor: const Color(0xFF1e1e2e),
        shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(20)),
        title: const Text('Rembourser ce don ?', style: TextStyle(color: Colors.white)),
        content: Column(
            mainAxisSize: MainAxisSize.min,
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
                const Text('Cette action est irréversible. Le montant sera retourné au donateur.', style: TextStyle(color: Colors.grey, fontSize: 13)),
                const SizedBox(height: 16),
                TextField(
                    controller: reasonController,
                    style: const TextStyle(color: Colors.white),
                    maxLines: 2,
                    decoration: InputDecoration(
                        filled: true,
                        fillColor: Colors.black26,
                        hintText: 'Motif (Optionnel)',
                        hintStyle: TextStyle(color: Colors.white.withOpacity(0.3)),
                        border: OutlineInputBorder(borderRadius: BorderRadius.circular(12), borderSide: BorderSide.none),
                    ),
                )
            ],
        ),
        actions: [
          TextButton(onPressed: () => Navigator.pop(ctx, false), child: const Text('Annuler', style: TextStyle(color: Colors.grey))),
          TextButton(
              onPressed: () => Navigator.pop(ctx, true), 
              style: TextButton.styleFrom(backgroundColor: Colors.red.withOpacity(0.1)),
              child: const Text('Rembourser', style: TextStyle(color: Colors.red, fontWeight: FontWeight.bold))
          ),
        ],
      ),
    );

    if (confirm != true) return;

    try {
      await _api.refundDonation(donationId, reason: reasonController.text);
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
              : Column(
                children: [
                   // Search Bar
                   Padding(
                     padding: const EdgeInsets.symmetric(horizontal: 16),
                     child: TextField(
                       style: const TextStyle(color: Colors.white),
                       decoration: InputDecoration(
                         filled: true,
                         fillColor: Colors.white.withOpacity(0.05),
                         hintText: 'Rechercher...',
                         hintStyle: TextStyle(color: Colors.white.withOpacity(0.4)),
                         prefixIcon: const Icon(Icons.search, color: Colors.grey),
                         border: OutlineInputBorder(borderRadius: BorderRadius.circular(12), borderSide: BorderSide.none),
                         contentPadding: const EdgeInsets.symmetric(horizontal: 16, vertical: 12),
                       ),
                       onChanged: (v) => setState(() => _searchQuery = v),
                     ),
                   ),
                   
                   // Tier Stats
                   if (_tierStats.isNotEmpty)
                     SizedBox(
                       height: 100,
                       child: ListView.separated(
                         padding: const EdgeInsets.all(16),
                         scrollDirection: Axis.horizontal,
                         itemCount: _tierStats.length,
                         separatorBuilder: (c, i) => const SizedBox(width: 12),
                         itemBuilder: (context, index) {
                           final stat = _tierStats[index];
                           return Container(
                             padding: const EdgeInsets.all(12),
                             decoration: BoxDecoration(
                               color: Colors.white.withOpacity(0.05),
                               borderRadius: BorderRadius.circular(12),
                               border: Border.all(color: Colors.white10),
                             ),
                             child: Column(
                               crossAxisAlignment: CrossAxisAlignment.start,
                               mainAxisAlignment: MainAxisAlignment.center,
                               children: [
                                 Text(stat['label'], style: const TextStyle(color: Colors.grey, fontSize: 10, fontWeight: FontWeight.bold)),
                                 const SizedBox(height: 4),
                                 Text('${stat['count']} dons', style: const TextStyle(color: Colors.white, fontWeight: FontWeight.bold)),
                                 Text(_formatAmount(stat['total'], _campaign!['currency']), style: const TextStyle(color: Color(0xFF6366f1), fontSize: 10)),
                               ],
                             ),
                           );
                         },
                       ),
                     ),

                   Expanded(
                     child: _filteredDonations.isEmpty 
                      ? const Center(child: Text('Aucun résultat', style: TextStyle(color: Colors.grey)))
                      : ListView.builder(
                          padding: const EdgeInsets.all(16),
                          itemCount: _filteredDonations.length,
                          itemBuilder: (context, index) {
                            return _buildDonorRow(_filteredDonations[index]);
                          },
                        ),
                   ),
                ],
              ),
    );
  }

  void _showDonationDetails(dynamic donation) {
     final isAnon = donation['is_anonymous'] == true;
     final name = isAnon ? 'Donateur Anonyme' : (donation['donor_name'] ?? 'Bienfaiteur');
     
     showModalBottomSheet(
       context: context,
       backgroundColor: Colors.transparent,
       builder: (context) => Container(
         decoration: const BoxDecoration(
            color: Color(0xFF1a1a2e),
            borderRadius: BorderRadius.vertical(top: Radius.circular(20)),
         ),
         padding: const EdgeInsets.all(24),
         child: SingleChildScrollView(
           child: Column(
             mainAxisSize: MainAxisSize.min,
             children: [
                CircleAvatar(
                   radius: 30,
                   backgroundColor: isAnon ? Colors.grey[800] : Colors.indigo[900],
                   child: Text(isAnon ? '?' : name[0].toUpperCase(), style: const TextStyle(color: Colors.white, fontSize: 24)),
                ),
                const SizedBox(height: 16),
                Text(name, style: const TextStyle(color: Colors.white, fontSize: 20, fontWeight: FontWeight.bold)),
                Text(_formatDate(donation['created_at']), style: const TextStyle(color: Colors.grey)),
                const SizedBox(height: 24),
                
                // Amount
                Container(
                   padding: const EdgeInsets.all(16),
                   decoration: BoxDecoration(
                      color: Colors.white.withOpacity(0.05),
                      borderRadius: BorderRadius.circular(16)
                   ),
                   child: Column(
                      children: [
                         const Text('Montant du Don', style: TextStyle(color: Colors.grey, fontSize: 12)),
                         const SizedBox(height: 4),
                         Text(
                             _formatAmount((donation['amount'] ?? 0).toDouble(), donation['currency']),
                             style: const TextStyle(color: Color(0xFF6366f1), fontSize: 24, fontWeight: FontWeight.bold),
                         ),

                         if (donation['frequency'] != null && donation['frequency'] != 'one_time')
                            Container(
                               margin: const EdgeInsets.only(top: 8),
                               padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 4),
                               decoration: BoxDecoration(color: Colors.indigo.withOpacity(0.2), borderRadius: BorderRadius.circular(8)),
                               child: Text('🔄 ${donation['frequency']}', style: const TextStyle(color: Colors.indigoAccent, fontSize: 12)),
                            )
                      ],
                   ),
                ),
                const SizedBox(height: 16),

                // Message
                if (donation['message'] != null && donation['message'].isNotEmpty)
                   Column(
                     children: [
                        const Text('MESSAGE', style: TextStyle(color: Colors.grey, fontSize: 10, fontWeight: FontWeight.bold)),
                        const SizedBox(height: 8),
                        Container(
                           padding: const EdgeInsets.all(12),
                           decoration: BoxDecoration(color: Colors.white.withOpacity(0.05), borderRadius: BorderRadius.circular(12)),
                           child: Text('"${donation['message']}"', style: const TextStyle(color: Colors.white, fontStyle: FontStyle.italic), textAlign: TextAlign.center),
                        ),
                        const SizedBox(height: 16),
                     ],
                   ),

                // Custom Fields
                if (donation['form_data'] != null && (donation['form_data'] as Map).isNotEmpty)
                   Column(
                     children: [
                        const Text('INFORMATIONS', style: TextStyle(color: Colors.grey, fontSize: 10, fontWeight: FontWeight.bold)),
                        const SizedBox(height: 8),
                        ...(donation['form_data'] as Map).entries.map((e) => Padding(
                           padding: const EdgeInsets.only(bottom: 8),
                           child: Row(
                              mainAxisAlignment: MainAxisAlignment.spaceBetween,
                              children: [
                                 Text(e.key.toString().replaceAll('_', ' ').toUpperCase(), style: const TextStyle(color: Colors.grey, fontSize: 12)),
                                 Text(e.value.toString(), style: const TextStyle(color: Colors.white, fontWeight: FontWeight.bold)),
                              ],
                           ),
                        )).toList(),
                        const SizedBox(height: 16),
                     ],
                   ),

                if (_isCreator && donation['status'] == 'paid')
                    TextButton.icon(
                        icon: const Icon(Icons.undo, color: Colors.red),
                        label: const Text('Rembourser ce don', style: TextStyle(color: Colors.red)),
                        onPressed: () {
                           Navigator.pop(context);
                           Future.delayed(const Duration(milliseconds: 200), () {
                               _refundDonation(donation['id']);
                           });
                        },
                    )
             ],
           ),
         ),
       ),
     );
  }

  Widget _buildDonorRow(dynamic donation) {
    final isAnon = donation['is_anonymous'] == true;
    final name = isAnon ? 'Donateur Anonyme' : (donation['donor_name'] ?? 'Bienfaiteur');
    final status = donation['status'];

    return Container(
      margin: const EdgeInsets.only(bottom: 12),
      child: Material(
        color: Colors.transparent,
        child: InkWell(
          onTap: () => _showDonationDetails(donation),
          borderRadius: BorderRadius.circular(12),
          child: Container(
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
          ),
        ),
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
