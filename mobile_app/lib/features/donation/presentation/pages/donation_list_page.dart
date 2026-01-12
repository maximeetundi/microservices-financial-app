import 'package:flutter/material.dart';
import '../../../../core/services/donation_api_service.dart';
import 'campaign_detail_page.dart';

class DonationListPage extends StatefulWidget {
  const DonationListPage({super.key});

  @override
  State<DonationListPage> createState() => _DonationListPageState();
}

class _DonationListPageState extends State<DonationListPage> with SingleTickerProviderStateMixin {
  final DonationApiService _api = DonationApiService();
  late TabController _tabController;
  
  List<dynamic> _campaigns = [];
  bool _loading = true;

  @override
  void initState() {
    super.initState();
    _tabController = TabController(length: 2, vsync: this);
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
      final campaigns = await _api.getCampaigns(limit: 50);
      setState(() {
        _campaigns = campaigns;
      });
    } catch (e) {
      debugPrint('Error loading campaigns: $e');
    } finally {
      setState(() => _loading = false);
    }
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
          child: Column(
            children: [
              _buildHeader(),
              Expanded(
                child: _loading
                    ? const Center(child: CircularProgressIndicator())
                    : RefreshIndicator(
                        onRefresh: _loadData,
                        child: _campaigns.isEmpty 
                          ? _buildEmptyState()
                          : ListView.builder(
                              padding: const EdgeInsets.all(20),
                              itemCount: _campaigns.length,
                              itemBuilder: (context, index) => _buildCampaignCard(_campaigns[index]),
                            ),
                      ),
              ),
            ],
          ),
        ),
      ),
      floatingActionButton: FloatingActionButton.extended(
        onPressed: () {
          // Navigate to Create Campaign (Not implemented yet for mobile MVP or reuse event form logic?)
          // For now just show snackbar or basic placeholder
          ScaffoldMessenger.of(context).showSnackBar(const SnackBar(content: Text('Création bientôt disponible sur mobile')));
        },
        backgroundColor: const Color(0xFF6366f1),
        icon: const Icon(Icons.add),
        label: const Text('Lancer'),
      ),
    );
  }

  Widget _buildHeader() {
    return Padding(
      padding: const EdgeInsets.all(20),
      child: Row(
        mainAxisAlignment: MainAxisAlignment.spaceBetween,
        children: [
          Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              const Text(
                '🤲 Dons & Solidarité',
                style: TextStyle(
                  fontSize: 28,
                  fontWeight: FontWeight.bold,
                  color: Colors.white,
                ),
              ),
              Text(
                'Soutenez des causes qui comptent',
                style: TextStyle(
                  fontSize: 14,
                  color: Colors.white.withOpacity(0.6),
                ),
              ),
            ],
          ),
        ],
      ),
    );
  }

  Widget _buildEmptyState() {
    return Center(
      child: Column(
        mainAxisAlignment: MainAxisAlignment.center,
        children: [
          const Text('🌱', style: TextStyle(fontSize: 64)),
          const SizedBox(height: 16),
          const Text(
            'Aucune campagne',
            style: TextStyle(
              color: Colors.white,
              fontSize: 20,
              fontWeight: FontWeight.bold,
            ),
          ),
          const SizedBox(height: 8),
          Text(
            'Soyez le premier à en lancer une !',
            style: TextStyle(
              color: Colors.white.withOpacity(0.6),
            ),
          ),
        ],
      ),
    );
  }

  Widget _buildCampaignCard(dynamic campaign) {
    final collected = (campaign['collected_amount'] ?? 0).toDouble();
    final target = (campaign['target_amount'] ?? 0).toDouble();
    final progress = target > 0 ? (collected / target) : 0.0;
    final percent = (progress * 100).clamp(0, 100).toInt();

    return GestureDetector(
      onTap: () => Navigator.push(
        context,
        MaterialPageRoute(
          builder: (_) => CampaignDetailPage(campaignId: campaign['id']),
        ),
      ).then((_) => _loadData()),
      child: Container(
        margin: const EdgeInsets.only(bottom: 16),
        decoration: BoxDecoration(
          color: Colors.white.withOpacity(0.05),
          borderRadius: BorderRadius.circular(16),
          border: Border.all(color: Colors.white.withOpacity(0.1)),
        ),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            // Image
            ClipRRect(
              borderRadius: const BorderRadius.vertical(top: Radius.circular(16)),
              child: Container(
                height: 150,
                width: double.infinity,
                color: Colors.white.withOpacity(0.1),
                child: campaign['image_url'] != null && campaign['image_url'].isNotEmpty
                    ? Image.network(campaign['image_url'], fit: BoxFit.cover)
                    : const Center(child: Text('🤲', style: TextStyle(fontSize: 48))),
              ),
            ),
            
            Padding(
              padding: const EdgeInsets.all(16),
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  Text(
                    campaign['title'] ?? 'Sans titre',
                    style: const TextStyle(
                      color: Colors.white,
                      fontSize: 18,
                      fontWeight: FontWeight.bold,
                    ),
                  ),
                  const SizedBox(height: 8),
                  Text(
                    campaign['description'] ?? '',
                    maxLines: 2,
                    overflow: TextOverflow.ellipsis,
                    style: TextStyle(
                      color: Colors.white.withOpacity(0.7),
                      fontSize: 14,
                    ),
                  ),
                  const SizedBox(height: 16),
                  
                  // Progress Bar
                  ClipRRect(
                    borderRadius: BorderRadius.circular(4),
                    child: LinearProgressIndicator(
                      value: progress > 1 ? 1 : progress,
                      backgroundColor: Colors.white.withOpacity(0.1),
                      valueColor: const AlwaysStoppedAnimation<Color>(Color(0xFF6366f1)),
                      minHeight: 6,
                    ),
                  ),
                  const SizedBox(height: 8),
                  
                  Row(
                    mainAxisAlignment: MainAxisAlignment.spaceBetween,
                    children: [
                      Text(
                        _formatAmount(collected, campaign['currency']),
                        style: const TextStyle(
                          color: Color(0xFF6366f1),
                          fontWeight: FontWeight.bold,
                        ),
                      ),
                      if (target > 0)
                        Text(
                          '$percent%',
                          style: TextStyle(color: Colors.white.withOpacity(0.5)),
                        )
                      else
                         const Text(
                          'Illimité',
                          style: TextStyle(color: Colors.green),
                        ),
                    ],
                  ),
                ],
              ),
            ),
          ],
        ),
      ),
    );
  }

  String _formatAmount(double amount, String? currency) {
    return '${amount.toStringAsFixed(0)} ${currency ?? "XOF"}';
  }
}
