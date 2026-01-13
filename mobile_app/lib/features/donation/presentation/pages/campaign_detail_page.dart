import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import 'dart:convert';
import 'package:qr_flutter/qr_flutter.dart';
import '../../../../core/services/donation_api_service.dart';
import '../widgets/donate_modal.dart';
import 'donors_list_page.dart';

class CampaignDetailPage extends StatefulWidget {
  final String campaignId;

  const CampaignDetailPage({super.key, required this.campaignId});

  @override
  State<CampaignDetailPage> createState() => _CampaignDetailPageState();
}

class _CampaignDetailPageState extends State<CampaignDetailPage> {
  final DonationApiService _api = DonationApiService();
  
  Map<String, dynamic>? _campaign;
  List<dynamic> _donations = [];
  bool _loading = true;

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
        _api.getDonations(widget.campaignId, limit: 1), // Only fetch 1 to check emptiness or just skip
      ]);
      setState(() {
        _campaign = results[0] as Map<String, dynamic>;
        // We might not need _donations anymore here if we just link to page
        // But keeping it for now if we want to show count etc.
        _donations = results[1] as List<dynamic>;
      });
    } catch (e) {
      debugPrint('Error loading campaign: $e');
    } finally {
      setState(() => _loading = false);
    }
  }

  void _openDonateModal() {
    if (_campaign == null) return;
    
    showModalBottomSheet(
      context: context,
      isScrollControlled: true,
      backgroundColor: Colors.transparent,
      builder: (context) => DonateModal(
        campaignId: widget.campaignId,
        currency: _campaign!['currency'] ?? 'XOF',
        formSchema: _campaign!['form_schema'] != null ? List<dynamic>.from(_campaign!['form_schema']) : null,
        donationType: _campaign!['donation_type'] ?? 'free',
        fixedAmount: double.tryParse(_campaign!['fixed_amount']?.toString() ?? ''),
        minAmount: double.tryParse(_campaign!['min_amount']?.toString() ?? ''),
        maxAmount: double.tryParse(_campaign!['max_amount']?.toString() ?? ''),
        donationTiers: _campaign!['donation_tiers'] != null ? List<dynamic>.from(_campaign!['donation_tiers']) : null,
      ),
    ).then((result) {
      if (result == true) {
        _loadData(); // Refresh on success
      }
    });
  }

  @override
  Widget build(BuildContext context) {
    if (_loading) {
      return const Scaffold(
        backgroundColor: Color(0xFF1a1a2e),
        body: Center(child: CircularProgressIndicator()),
      );
    }

    if (_campaign == null) {
      return const Scaffold(
        backgroundColor: Color(0xFF1a1a2e),
        body: Center(child: Text('Erreur lors du chargement', style: TextStyle(color: Colors.white))),
      );
    }

    final campaign = _campaign!;
    final collected = (campaign['collected_amount'] ?? 0).toDouble();
    final target = (campaign['target_amount'] ?? 0).toDouble();
    final progress = target > 0 ? (collected / target) : 0.0;
    final percent = (progress * 100).clamp(0, 100).toInt();

    return Scaffold(
      backgroundColor: const Color(0xFF1a1a2e),
      body: CustomScrollView(
        slivers: [
          // App Bar with Image
          SliverAppBar(
            expandedHeight: 250,
            pinned: true,
            backgroundColor: const Color(0xFF1a1a2e),
            flexibleSpace: FlexibleSpaceBar(
              background: Stack(
                fit: StackFit.expand,
                children: [
                   campaign['image_url'] != null && campaign['image_url'].isNotEmpty
                      ? Image.network(campaign['image_url'], fit: BoxFit.cover)
                      : Container(color: Colors.indigo, child: const Center(child: Text('🤲', style: TextStyle(fontSize: 64)))),
                   Container(
                     decoration: BoxDecoration(
                       gradient: LinearGradient(
                         begin: Alignment.topCenter,
                         end: Alignment.bottomCenter,
                         colors: [Colors.transparent, const Color(0xFF1a1a2e)],
                       ),
                     ),
                   ),
                ],
              ),
            ),
          ),

          // Content
          SliverToBoxAdapter(
            child: Padding(
              padding: const EdgeInsets.all(20),
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  Text(
                    campaign['title'] ?? 'Sans titre',
                    style: const TextStyle(
                      color: Colors.white,
                      fontSize: 24,
                      fontWeight: FontWeight.bold,
                    ),
                  ),
                  const SizedBox(height: 8),
                  Text(
                    'Organisé par ${campaign['creator_id']}', 
                    style: TextStyle(color: Colors.white.withOpacity(0.5), fontSize: 12),
                  ),
                  const SizedBox(height: 24),

                  // QR & Code
                  if (campaign['campaign_code'] != null)
                    Container(
                      margin: const EdgeInsets.only(bottom: 24),
                      padding: const EdgeInsets.all(16),
                      decoration: BoxDecoration(
                        color: Colors.white.withOpacity(0.05),
                        borderRadius: BorderRadius.circular(16),
                        border: Border.all(color: Colors.white24),
                      ),
                      child: Column(
                        children: [
                          Row(
                            children: [
                                 // QR Image
                                 Container(
                                   height: 100, 
                                   width: 100,
                                   decoration: BoxDecoration(
                                     color: Colors.white,
                                     borderRadius: BorderRadius.circular(8),
                                   ),
                                   padding: const EdgeInsets.all(4),
                                   child: QrImageView(
                                     data: 'https://app.maximeetundi.store/donations/${widget.campaignId}',
                                     version: QrVersions.auto,
                                     size: 92.0,
                                     backgroundColor: Colors.white,
                                   ),
                                 ),
                                 const SizedBox(width: 16),
                                 Expanded(
                                   child: Column(
                                     crossAxisAlignment: CrossAxisAlignment.start,
                                     children: [
                                       const Text('Code Campagne', style: TextStyle(color: Colors.grey, fontSize: 12)),
                                       Row(
                                         children: [
                                           Text(
                                             campaign['campaign_code'], 
                                             style: const TextStyle(color: Colors.white, fontSize: 20, fontWeight: FontWeight.bold, letterSpacing: 1.5),
                                           ),
                                           IconButton(
                                             icon: const Icon(Icons.copy, color: Color(0xFF6366f1), size: 20),
                                             onPressed: () {
                                               Clipboard.setData(ClipboardData(text: campaign['campaign_code']));
                                               ScaffoldMessenger.of(context).showSnackBar(
                                                 const SnackBar(content: Text('Code copié !')),
                                               );
                                             },
                                           ),
                                         ],
                                       ),
                                       const SizedBox(height: 4),
                                       const Text('Partagez ce code ou faites scanner le QR pour recevoir des dons.', style: TextStyle(color: Colors.grey, fontSize: 12)),
                                     ],
                                   )
                                 )
                            ],
                          ),
                        ],
                      ),
                    ),

                  // Progress
                  ClipRRect(
                    borderRadius: BorderRadius.circular(8),
                    child: LinearProgressIndicator(
                      value: progress > 1 ? 1 : progress,
                      backgroundColor: Colors.white.withOpacity(0.1),
                      valueColor: const AlwaysStoppedAnimation<Color>(Color(0xFF6366f1)),
                      minHeight: 12,
                    ),
                  ),
                  const SizedBox(height: 12),
                  Row(
                    mainAxisAlignment: MainAxisAlignment.spaceBetween,
                    children: [
                      Column(
                        crossAxisAlignment: CrossAxisAlignment.start,
                        children: [
                          Text(
                            _formatAmount(collected, campaign['currency']),
                            style: const TextStyle(color: Color(0xFF6366f1), fontSize: 20, fontWeight: FontWeight.bold),
                          ),
                          const Text('récoltés', style: TextStyle(color: Colors.grey, fontSize: 12)),
                        ],
                      ),
                       Column(
                        crossAxisAlignment: CrossAxisAlignment.end,
                        children: [
                          if (target > 0) ...[
                            Text(
                            _formatAmount(target, campaign['currency']),
                             style: const TextStyle(color: Colors.white, fontSize: 16, fontWeight: FontWeight.w500),
                            ),
                            Text('objectif ($percent%)', style: const TextStyle(color: Colors.grey, fontSize: 12)),
                          ] else 
                             const Text('Objectif illimité', style: TextStyle(color: Colors.green, fontWeight: FontWeight.bold)),
                        ],
                      ),
                    ],
                  ),
                  const SizedBox(height: 32),

                  // Description
                  const Text(
                    'À propos',
                    style: TextStyle(color: Colors.white, fontSize: 18, fontWeight: FontWeight.bold),
                  ),
                  const SizedBox(height: 12),
                  Text(
                    campaign['description'] ?? '',
                    style: TextStyle(color: Colors.white.withOpacity(0.8), height: 1.5),
                  ),
                  const SizedBox(height: 32),

                  // Donors
                  Container(
                    width: double.infinity,
                    padding: const EdgeInsets.all(20),
                    decoration: BoxDecoration(
                      gradient: LinearGradient(
                        colors: [Colors.indigo.shade900.withOpacity(0.5), Colors.indigo.shade800.withOpacity(0.3)],
                      ),
                      borderRadius: BorderRadius.circular(16),
                      border: Border.all(color: Colors.indigo.withOpacity(0.3)),
                    ),
                    child: Column(
                      children: [
                        const Text(
                           '🏆 Mur des donateurs',
                           style: TextStyle(color: Colors.white, fontSize: 18, fontWeight: FontWeight.bold),
                        ),
                        const SizedBox(height: 8),
                        Text(
                          'Voir la liste complète des donateurs',
                          style: TextStyle(color: Colors.indigo.shade200, fontSize: 14),
                        ),
                        const SizedBox(height: 16),
                        SizedBox(
                          width: double.infinity,
                          child: ElevatedButton(
                            onPressed: () {
                              Navigator.push(context, MaterialPageRoute(builder: (_) => DonorsListPage(campaignId: widget.campaignId)));
                            },
                            style: ElevatedButton.styleFrom(
                              backgroundColor: Colors.indigo,
                              padding: const EdgeInsets.symmetric(vertical: 12),
                              shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(12)),
                            ),
                            child: const Text('Voir les donateurs', style: TextStyle(fontWeight: FontWeight.bold)),
                          ),
                        ),
                      ],
                    ),
                  ),
                  
                  const SizedBox(height: 100), // Spacing for fab
                ],
              ),
            ),
          ),
        ],
      ),
      bottomNavigationBar: Container(
        padding: const EdgeInsets.all(20),
        decoration: const BoxDecoration(
          color: Color(0xFF1a1a2e),
          border: Border(top: BorderSide(color: Colors.white12)),
        ),
        child: SafeArea(
          child: ElevatedButton(
            onPressed: _openDonateModal,
            style: ElevatedButton.styleFrom(
              backgroundColor: const Color(0xFF6366f1),
              padding: const EdgeInsets.symmetric(vertical: 16),
              shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(16)),
            ),
            child: const Text('Faire un don ❤️', style: TextStyle(fontSize: 18, fontWeight: FontWeight.bold)),
          ),
        ),
      ),
    );
  }

  String _formatAmount(double amount, String? currency) {
    return '${amount.toStringAsFixed(0)} ${currency ?? "XOF"}';
  }
}
