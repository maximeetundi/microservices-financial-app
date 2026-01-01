import 'package:flutter/material.dart';
import '../../core/services/ticket_api_service.dart';
import 'purchase_ticket_screen.dart';

class EventDetailsScreen extends StatefulWidget {
  final String eventId;
  final bool isOwner;

  const EventDetailsScreen({
    super.key,
    required this.eventId,
    required this.isOwner,
  });

  @override
  State<EventDetailsScreen> createState() => _EventDetailsScreenState();
}

class _EventDetailsScreenState extends State<EventDetailsScreen> {
  final TicketApiService _ticketApi = TicketApiService();
  Map<String, dynamic>? _event;
  bool _loading = true;

  @override
  void initState() {
    super.initState();
    _loadEvent();
  }

  Future<void> _loadEvent() async {
    setState(() => _loading = true);
    try {
      _event = await _ticketApi.getEvent(widget.eventId);
    } catch (e) {
      debugPrint('Error loading event: $e');
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
          child: _loading
              ? const Center(child: CircularProgressIndicator())
              : _event == null
                  ? _buildError()
                  : _buildContent(),
        ),
      ),
    );
  }

  Widget _buildContent() {
    final event = _event!;
    final tiers = event['ticket_tiers'] as List? ?? [];

    return CustomScrollView(
      slivers: [
        // App Bar
        SliverAppBar(
          expandedHeight: 200,
          pinned: true,
          backgroundColor: const Color(0xFF1a1a2e),
          leading: IconButton(
            onPressed: () => Navigator.pop(context),
            icon: const Icon(Icons.arrow_back),
          ),
          actions: widget.isOwner
              ? [
                  IconButton(
                    onPressed: () {
                      // TODO: Edit event
                    },
                    icon: const Icon(Icons.edit),
                  ),
                ]
              : null,
          flexibleSpace: FlexibleSpaceBar(
            background: event['cover_image'] != null
                ? Image.network(
                    event['cover_image'],
                    fit: BoxFit.cover,
                    color: Colors.black.withOpacity(0.3),
                    colorBlendMode: BlendMode.darken,
                  )
                : Container(
                    color: const Color(0xFF6366f1).withOpacity(0.3),
                    child: const Center(
                      child: Text('🎪', style: TextStyle(fontSize: 64)),
                    ),
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
                // Title and status
                Row(
                  children: [
                    Expanded(
                      child: Text(
                        event['title'] ?? 'Sans titre',
                        style: const TextStyle(
                          color: Colors.white,
                          fontSize: 24,
                          fontWeight: FontWeight.bold,
                        ),
                      ),
                    ),
                    if (widget.isOwner)
                      Container(
                        padding: const EdgeInsets.symmetric(horizontal: 12, vertical: 6),
                        decoration: BoxDecoration(
                          color: _getStatusColor(event['status'] ?? 'draft'),
                          borderRadius: BorderRadius.circular(20),
                        ),
                        child: Text(
                          _getStatusLabel(event['status'] ?? 'draft'),
                          style: const TextStyle(
                            color: Colors.white,
                            fontSize: 12,
                            fontWeight: FontWeight.w600,
                          ),
                        ),
                      ),
                  ],
                ),
                const SizedBox(height: 16),

                // Event info
                _buildInfoRow(Icons.location_on, event['location'] ?? 'Non défini'),
                _buildInfoRow(Icons.calendar_today, _formatDate(event['start_date'])),
                _buildInfoRow(Icons.schedule, 'Ventes: ${_formatDate(event['sale_start_date'])} - ${_formatDate(event['sale_end_date'])}'),
                
                if (event['description']?.isNotEmpty ?? false) ...[
                  const SizedBox(height: 16),
                  Text(
                    event['description'],
                    style: TextStyle(
                      color: Colors.white.withOpacity(0.8),
                      fontSize: 14,
                      height: 1.5,
                    ),
                  ),
                ],

                // QR Code for event (organizer)
                if (widget.isOwner && event['qr_code'] != null) ...[
                  const SizedBox(height: 24),
                  Center(
                    child: Column(
                      children: [
                        const Text(
                          'QR Code de l\'événement',
                          style: TextStyle(
                            color: Colors.white,
                            fontSize: 16,
                            fontWeight: FontWeight.w600,
                          ),
                        ),
                        const SizedBox(height: 12),
                        Container(
                          padding: const EdgeInsets.all(16),
                          decoration: BoxDecoration(
                            color: Colors.white,
                            borderRadius: BorderRadius.circular(16),
                          ),
                          child: Image.network(
                            event['qr_code'],
                            width: 150,
                            height: 150,
                          ),
                        ),
                        const SizedBox(height: 8),
                        Text(
                          event['event_code'] ?? '',
                          style: const TextStyle(
                            color: Colors.white70,
                            fontSize: 16,
                            fontFamily: 'monospace',
                          ),
                        ),
                      ],
                    ),
                  ),
                ],

                // Ticket Tiers
                const SizedBox(height: 24),
                const Text(
                  'Niveaux de tickets',
                  style: TextStyle(
                    color: Colors.white,
                    fontSize: 18,
                    fontWeight: FontWeight.bold,
                  ),
                ),
                const SizedBox(height: 12),
                ...tiers.map((tier) => _buildTierCard(tier)),

                // Action buttons
                const SizedBox(height: 24),
                if (widget.isOwner && event['status'] == 'draft')
                  SizedBox(
                    width: double.infinity,
                    child: ElevatedButton(
                      onPressed: () async {
                        try {
                          await _ticketApi.publishEvent(widget.eventId);
                          _loadEvent();
                          ScaffoldMessenger.of(context).showSnackBar(
                            const SnackBar(content: Text('Événement publié !')),
                          );
                        } catch (e) {
                          ScaffoldMessenger.of(context).showSnackBar(
                            SnackBar(content: Text('Erreur: $e')),
                          );
                        }
                      },
                      style: ElevatedButton.styleFrom(
                        backgroundColor: const Color(0xFF6366f1),
                        padding: const EdgeInsets.all(16),
                        shape: RoundedRectangleBorder(
                          borderRadius: BorderRadius.circular(12),
                        ),
                      ),
                      child: const Text(
                        '🚀 Publier l\'événement',
                        style: TextStyle(fontSize: 16, fontWeight: FontWeight.w600),
                      ),
                    ),
                  ),
              ],
            ),
          ),
        ),
      ],
    );
  }

  Widget _buildInfoRow(IconData icon, String text) {
    return Padding(
      padding: const EdgeInsets.symmetric(vertical: 6),
      child: Row(
        children: [
          Icon(icon, color: Colors.white54, size: 18),
          const SizedBox(width: 12),
          Expanded(
            child: Text(
              text,
              style: TextStyle(
                color: Colors.white.withOpacity(0.8),
                fontSize: 14,
              ),
            ),
          ),
        ],
      ),
    );
  }

  Widget _buildTierCard(dynamic tier) {
    final available = tier['quantity'] == -1 
        ? 'Illimité' 
        : '${tier['quantity'] - (tier['sold'] ?? 0)} restants';
    final isSoldOut = tier['quantity'] != -1 && (tier['sold'] ?? 0) >= tier['quantity'];

    return Container(
      margin: const EdgeInsets.only(bottom: 12),
      padding: const EdgeInsets.all(16),
      decoration: BoxDecoration(
        color: Colors.white.withOpacity(0.05),
        borderRadius: BorderRadius.circular(12),
        border: Border.all(color: _hexToColor(tier['color'] ?? '#6366f1').withOpacity(0.5)),
      ),
      child: Column(
        children: [
          Row(
            children: [
              Text(tier['icon'] ?? '🎫', style: const TextStyle(fontSize: 32)),
              const SizedBox(width: 16),
              Expanded(
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    Text(
                      tier['name'] ?? 'Standard',
                      style: const TextStyle(
                        color: Colors.white,
                        fontSize: 16,
                        fontWeight: FontWeight.bold,
                      ),
                    ),
                    if (tier['description']?.isNotEmpty ?? false)
                      Text(
                        tier['description'],
                        style: TextStyle(
                          color: Colors.white.withOpacity(0.6),
                          fontSize: 12,
                        ),
                      ),
                    Text(
                      available,
                      style: TextStyle(
                        color: isSoldOut ? Colors.red : Colors.white.withOpacity(0.5),
                        fontSize: 12,
                        fontWeight: isSoldOut ? FontWeight.bold : FontWeight.normal,
                      ),
                    ),
                  ],
                ),
              ),
              Column(
                crossAxisAlignment: CrossAxisAlignment.end,
                children: [
                  Text(
                    '${_formatAmount(tier['price'] ?? 0)}',
                    style: const TextStyle(
                      color: Colors.white,
                      fontSize: 18,
                      fontWeight: FontWeight.bold,
                    ),
                  ),
                  Text(
                    'XOF',
                    style: TextStyle(
                      color: Colors.white.withOpacity(0.6),
                      fontSize: 12,
                    ),
                  ),
                ],
              ),
            ],
          ),
          if (!widget.isOwner && !isSoldOut) ...[
            const SizedBox(height: 12),
            SizedBox(
              width: double.infinity,
              child: ElevatedButton(
                onPressed: () {
                  Navigator.push(
                    context,
                    MaterialPageRoute(
                      builder: (_) => PurchaseTicketScreen(
                        event: _event!,
                        tier: tier,
                      ),
                    ),
                  ).then((_) => _loadEvent()); // Refresh after purchase
                },
                style: ElevatedButton.styleFrom(
                  backgroundColor: _hexToColor(tier['color'] ?? '#6366f1'),
                  padding: const EdgeInsets.symmetric(vertical: 12),
                  shape: RoundedRectangleBorder(
                    borderRadius: BorderRadius.circular(10),
                  ),
                ),
                child: const Text(
                  'Acheter ce ticket',
                  style: TextStyle(fontWeight: FontWeight.w600),
                ),
              ),
            ),
          ],
        ],
      ),
    );
  }

  Widget _buildError() {
    return Center(
      child: Column(
        mainAxisAlignment: MainAxisAlignment.center,
        children: [
          const Text('❌', style: TextStyle(fontSize: 64)),
          const SizedBox(height: 16),
          const Text(
            'Événement non trouvé',
            style: TextStyle(color: Colors.white, fontSize: 18),
          ),
          const SizedBox(height: 16),
          ElevatedButton(
            onPressed: () => Navigator.pop(context),
            child: const Text('Retour'),
          ),
        ],
      ),
    );
  }

  String _formatDate(String? dateStr) {
    if (dateStr == null) return 'Non défini';
    try {
      final date = DateTime.parse(dateStr);
      return '${date.day}/${date.month}/${date.year}';
    } catch (e) {
      return dateStr;
    }
  }

  String _formatAmount(dynamic amount) {
    final num = (amount is int) ? amount : (amount as double).toInt();
    return num.toString().replaceAllMapped(
      RegExp(r'(\d{1,3})(?=(\d{3})+(?!\d))'),
      (m) => '${m[1]} ',
    );
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
    if (hex.length == 6) hex = 'FF$hex';
    return Color(int.parse(hex, radix: 16));
  }
}
