import 'package:flutter/material.dart';
import '../../core/services/ticket_api_service.dart';
import 'my_tickets_screen.dart';
import 'create_event_screen.dart';
import 'event_details_screen.dart';
import 'verify_ticket_screen.dart';

class EventsScreen extends StatefulWidget {
  const EventsScreen({super.key});

  @override
  State<EventsScreen> createState() => _EventsScreenState();
}

class _EventsScreenState extends State<EventsScreen> with SingleTickerProviderStateMixin {
  final TicketApiService _ticketApi = TicketApiService();
  late TabController _tabController;
  
  List<dynamic> _myEvents = [];
  List<dynamic> _activeEvents = [];
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
      final results = await Future.wait([
        _ticketApi.getMyEvents(),
        _ticketApi.getActiveEvents(),
      ]);
      setState(() {
        _myEvents = results[0];
        _activeEvents = results[1];
      });
    } catch (e) {
      debugPrint('Error loading events: $e');
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
              _buildTabBar(),
              Expanded(
                child: _loading
                    ? const Center(child: CircularProgressIndicator())
                    : TabBarView(
                        controller: _tabController,
                        children: [
                          _buildMyEventsTab(),
                          _buildDiscoverTab(),
                        ],
                      ),
              ),
            ],
          ),
        ),
      ),
      floatingActionButton: FloatingActionButton.extended(
        onPressed: () => Navigator.push(
          context,
          MaterialPageRoute(builder: (_) => CreateEventFormScreen()),
        ).then((_) => _loadData()),
        backgroundColor: const Color(0xFF6366f1),
        icon: const Icon(Icons.add),
        label: const Text('Créer'),
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
                '🎫 Événements',
                style: TextStyle(
                  fontSize: 28,
                  fontWeight: FontWeight.bold,
                  color: Colors.white,
                ),
              ),
              Text(
                'Créez et gérez vos événements',
                style: TextStyle(
                  fontSize: 14,
                  color: Colors.white.withOpacity(0.6),
                ),
              ),
            ],
          ),
          Row(
            children: [
              IconButton(
                onPressed: () => Navigator.push(
                  context,
                  MaterialPageRoute(builder: (_) => const MyTicketsScreen()),
                ),
                icon: const Icon(Icons.confirmation_number, color: Colors.white),
                tooltip: 'Mes tickets',
              ),
              IconButton(
                onPressed: () => Navigator.push(
                  context,
                  MaterialPageRoute(builder: (_) => const VerifyTicketScreen()),
                ),
                icon: const Icon(Icons.qr_code_scanner, color: Colors.white),
                tooltip: 'Vérifier ticket',
              ),
            ],
          ),
        ],
      ),
    );
  }

  Widget _buildTabBar() {
    return Container(
      margin: const EdgeInsets.symmetric(horizontal: 20),
      decoration: BoxDecoration(
        color: Colors.white.withOpacity(0.1),
        borderRadius: BorderRadius.circular(12),
      ),
      child: TabBar(
        controller: _tabController,
        indicator: BoxDecoration(
          color: const Color(0xFF6366f1),
          borderRadius: BorderRadius.circular(12),
        ),
        labelColor: Colors.white,
        unselectedLabelColor: Colors.white.withOpacity(0.6),
        tabs: const [
          Tab(text: 'Mes événements'),
          Tab(text: 'Découvrir'),
        ],
      ),
    );
  }

  Widget _buildMyEventsTab() {
    if (_myEvents.isEmpty) {
      return _buildEmptyState(
        icon: '🎪',
        title: 'Aucun événement',
        subtitle: 'Créez votre premier événement',
      );
    }

    return RefreshIndicator(
      onRefresh: _loadData,
      child: ListView.builder(
        padding: const EdgeInsets.all(20),
        itemCount: _myEvents.length,
        itemBuilder: (context, index) => _buildEventCard(_myEvents[index], isOwner: true),
      ),
    );
  }

  Widget _buildDiscoverTab() {
    if (_activeEvents.isEmpty) {
      return _buildEmptyState(
        icon: '🔍',
        title: 'Aucun événement actif',
        subtitle: 'Revenez bientôt',
      );
    }

    return RefreshIndicator(
      onRefresh: _loadData,
      child: ListView.builder(
        padding: const EdgeInsets.all(20),
        itemCount: _activeEvents.length,
        itemBuilder: (context, index) => _buildEventCard(_activeEvents[index], isOwner: false),
      ),
    );
  }

  Widget _buildEventCard(dynamic event, {required bool isOwner}) {
    final tiers = event['ticket_tiers'] as List? ?? [];
    final status = event['status'] ?? 'draft';
    
    return GestureDetector(
      onTap: () => Navigator.push(
        context,
        MaterialPageRoute(
          builder: (_) => EventDetailsScreen(
            eventId: event['id'],
            isOwner: isOwner,
          ),
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
            // Cover Image
            ClipRRect(
              borderRadius: const BorderRadius.vertical(top: Radius.circular(16)),
              child: Container(
                height: 120,
                width: double.infinity,
                color: Colors.white.withOpacity(0.1),
                child: event['cover_image'] != null && event['cover_image'].isNotEmpty
                    ? Image.network(event['cover_image'], fit: BoxFit.cover)
                    : const Center(
                        child: Text('🎪', style: TextStyle(fontSize: 48)),
                      ),
              ),
            ),
            
            // Status Badge
            if (isOwner)
              Padding(
                padding: const EdgeInsets.only(left: 12, top: 12),
                child: Container(
                  padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 4),
                  decoration: BoxDecoration(
                    color: _getStatusColor(status),
                    borderRadius: BorderRadius.circular(12),
                  ),
                  child: Text(
                    _getStatusLabel(status),
                    style: const TextStyle(
                      color: Colors.white,
                      fontSize: 12,
                      fontWeight: FontWeight.w600,
                    ),
                  ),
                ),
              ),
            
            // Content
            Padding(
              padding: const EdgeInsets.all(16),
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  Text(
                    event['title'] ?? 'Sans titre',
                    style: const TextStyle(
                      color: Colors.white,
                      fontSize: 18,
                      fontWeight: FontWeight.bold,
                    ),
                  ),
                  const SizedBox(height: 8),
                  Row(
                    children: [
                      Icon(Icons.location_on, size: 14, color: Colors.white.withOpacity(0.6)),
                      const SizedBox(width: 4),
                      Text(
                        event['location'] ?? 'Non défini',
                        style: TextStyle(
                          color: Colors.white.withOpacity(0.6),
                          fontSize: 13,
                        ),
                      ),
                    ],
                  ),
                  const SizedBox(height: 4),
                  Row(
                    children: [
                      Icon(Icons.calendar_today, size: 14, color: Colors.white.withOpacity(0.6)),
                      const SizedBox(width: 4),
                      Text(
                        _formatDate(event['start_date']),
                        style: TextStyle(
                          color: Colors.white.withOpacity(0.6),
                          fontSize: 13,
                        ),
                      ),
                    ],
                  ),
                  const SizedBox(height: 12),
                  
                  // Stats or Price
                  if (isOwner)
                    Row(
                      children: [
                        _buildStat('${event['total_sold'] ?? 0}', 'vendus'),
                        const SizedBox(width: 24),
                        _buildStat(_formatAmount(event['total_revenue'] ?? 0), 'revenus'),
                      ],
                    )
                  else if (tiers.isNotEmpty)
                    Wrap(
                      spacing: 8,
                      children: tiers.take(3).map<Widget>((tier) {
                        return Container(
                          padding: const EdgeInsets.symmetric(horizontal: 10, vertical: 4),
                          decoration: BoxDecoration(
                            color: _hexToColor(tier['color'] ?? '#6366f1'),
                            borderRadius: BorderRadius.circular(12),
                          ),
                          child: Text(
                            '${tier['icon'] ?? '🎫'} ${tier['name']} - ${_formatAmount(tier['price'] ?? 0)}',
                            style: const TextStyle(
                              color: Colors.white,
                              fontSize: 12,
                              fontWeight: FontWeight.w500,
                            ),
                          ),
                        );
                      }).toList(),
                    ),
                ],
              ),
            ),
          ],
        ),
      ),
    );
  }

  Widget _buildStat(String value, String label) {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Text(
          value,
          style: const TextStyle(
            color: Colors.white,
            fontSize: 16,
            fontWeight: FontWeight.bold,
          ),
        ),
        Text(
          label,
          style: TextStyle(
            color: Colors.white.withOpacity(0.5),
            fontSize: 12,
          ),
        ),
      ],
    );
  }

  Widget _buildEmptyState({required String icon, required String title, required String subtitle}) {
    return Center(
      child: Column(
        mainAxisAlignment: MainAxisAlignment.center,
        children: [
          Text(icon, style: const TextStyle(fontSize: 64)),
          const SizedBox(height: 16),
          Text(
            title,
            style: const TextStyle(
              color: Colors.white,
              fontSize: 20,
              fontWeight: FontWeight.bold,
            ),
          ),
          const SizedBox(height: 8),
          Text(
            subtitle,
            style: TextStyle(
              color: Colors.white.withOpacity(0.6),
            ),
          ),
        ],
      ),
    );
  }

  String _formatDate(String? dateStr) {
    if (dateStr == null) return 'Non défini';
    try {
      final date = DateTime.parse(dateStr);
      return '${date.day}/${date.month}/${date.year} ${date.hour}:${date.minute.toString().padLeft(2, '0')}';
    } catch (e) {
      return dateStr;
    }
  }

  String _formatAmount(dynamic amount) {
    final num = (amount is int) ? amount : (amount as double).toInt();
    return '${num.toString().replaceAllMapped(RegExp(r'(\d{1,3})(?=(\d{3})+(?!\d))'), (m) => '${m[1]} ')} XOF';
  }

  String _getStatusLabel(String status) {
    switch (status) {
      case 'draft': return 'Brouillon';
      case 'active': return 'Actif';
      case 'ended': return 'Terminé';
      case 'cancelled': return 'Annulé';
      default: return status;
    }
  }

  Color _getStatusColor(String status) {
    switch (status) {
      case 'draft': return Colors.grey;
      case 'active': return Colors.green;
      case 'ended': return Colors.orange;
      case 'cancelled': return Colors.red;
      default: return Colors.grey;
    }
  }

  Color _hexToColor(String hex) {
    hex = hex.replaceFirst('#', '');
    if (hex.length == 6) {
      hex = 'FF$hex';
    }
    return Color(int.parse(hex, radix: 16));
  }
}
