import 'package:flutter/material.dart';
import '../../../../core/services/donation_api_service.dart';

class DonorsListPage extends StatefulWidget {
  final String campaignId;

  const DonorsListPage({super.key, required this.campaignId});

  @override
  State<DonorsListPage> createState() => _DonorsListPageState();
}

class _DonorsListPageState extends State<DonorsListPage> {
  final DonationApiService _api = DonationApiService();
  List<dynamic> _donations = [];
  bool _loading = true;

  @override
  void initState() {
    super.initState();
    _loadDonations();
  }

  Future<void> _loadDonations() async {
    setState(() => _loading = true);
    try {
      final donations = await _api.getDonations(widget.campaignId, limit: 100);
      setState(() {
        _donations = donations;
      });
    } catch (e) {
      debugPrint('Error: $e');
    } finally {
      setState(() => _loading = false);
    }
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: const Color(0xFF1a1a2e),
      appBar: AppBar(
        title: const Text('Donateurs'),
        backgroundColor: Colors.transparent,
        elevation: 0,
      ),
      body: _loading 
          ? const Center(child: CircularProgressIndicator())
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
    
    return Container(
      margin: const EdgeInsets.only(bottom: 12),
      padding: const EdgeInsets.all(12),
      decoration: BoxDecoration(
        color: Colors.white.withOpacity(0.05),
        borderRadius: BorderRadius.circular(12),
      ),
      child: Row(
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
              ],
            ),
          ),
          Text(
            _formatAmount((donation['amount'] ?? 0).toDouble(), donation['currency']),
            style: const TextStyle(color: Color(0xFF6366f1), fontWeight: FontWeight.bold),
          ),
        ],
      ),
    );
  }

  String _formatAmount(double amount, String? currency) {
    return '${amount.toStringAsFixed(0)} ${currency ?? "XOF"}';
  }
}
