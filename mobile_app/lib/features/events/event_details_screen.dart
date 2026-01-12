import 'package:flutter/material.dart';
import '../../core/services/ticket_api_service.dart';
import 'purchase_ticket_screen.dart';
import 'ticket_scanner_screen.dart';
import 'sold_tickets_screen.dart';
import 'edit_event_screen.dart';

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

  Future<void> _editEvent(Map<String, dynamic> event) async {
    final result = await Navigator.push(
      context,
      MaterialPageRoute(
        builder: (context) => EditEventScreen(event: event),
      ),
    );

    if (result == true) {
      _loadEvent();
    }
  }

  Future<void> _confirmDelete(Map<String, dynamic> event) async {
    final confirmed = await showDialog<bool>(
      context: context,
      builder: (context) => AlertDialog(
        title: const Text('Supprimer l\'événement ?'),
        content: const Text('Cette action est irréversible. Êtes-vous sûr ?'),
        actions: [
          TextButton(
            onPressed: () => Navigator.pop(context, false),
            child: const Text('Annuler'),
          ),
          TextButton(
            onPressed: () => Navigator.pop(context, true),
            child: const Text('Supprimer', style: TextStyle(color: Colors.red)),
          ),
        ],
      ),
    );

    if (confirmed == true) {
      setState(() => _loading = true);
      try {
        await _ticketApi.deleteEvent(widget.eventId);
        if (mounted) {
          Navigator.pop(context); // Go back to list
          ScaffoldMessenger.of(context).showSnackBar(
            const SnackBar(content: Text('Événement supprimé'), backgroundColor: Colors.green),
          );
        }
      } catch (e) {
        if (mounted) {
          setState(() => _loading = false);
          ScaffoldMessenger.of(context).showSnackBar(
            SnackBar(content: Text('Erreur: ${e.toString()}')),
          );
        }
      }
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
                  IconButton(
                    onPressed: () => _editEvent(event),
                    icon: const Icon(Icons.edit),
                    tooltip: 'Modifier',
                  ),
                  IconButton(
                    onPressed: () => _confirmDelete(event),
                    icon: const Icon(Icons.delete, color: Colors.white),
                    tooltip: 'Supprimer',
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
                if (widget.isOwner) ...[
                  const SizedBox(height: 24),
                  _buildEventCodeSection(event),
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
          if (!widget.isOwner) ...[
            const SizedBox(height: 12),
            SizedBox(
              width: double.infinity,
              child: ElevatedButton(
                onPressed: _event!['status'] == 'active' && !isSoldOut
                    ? () {
                        Navigator.push(
                          context,
                          MaterialPageRoute(
                            builder: (_) => PurchaseTicketScreen(
                              event: _event!,
                              tier: tier,
                            ),
                          ),
                        ).then((_) => _loadEvent()); // Refresh after purchase
                      }
                    : null,
                style: ElevatedButton.styleFrom(
                  backgroundColor: _hexToColor(tier['color'] ?? '#6366f1'),
                  disabledBackgroundColor: Colors.white.withOpacity(0.1),
                  padding: const EdgeInsets.symmetric(vertical: 12),
                  shape: RoundedRectangleBorder(
                    borderRadius: BorderRadius.circular(10),
                  ),
                ),
                child: Text(
                  _event!['status'] == 'active'
                      ? (isSoldOut ? 'Épuisé' : 'Acheter ce ticket')
                      : 'Ventes fermées',
                  style: TextStyle(
                    fontWeight: FontWeight.w600,
                    color: _event!['status'] == 'active' && !isSoldOut ? Colors.white : Colors.white54,
                  ),
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

  Widget _buildEventCodeSection(Map<String, dynamic> event) {
    return Container(
      padding: const EdgeInsets.all(16),
      decoration: BoxDecoration(
        color: Colors.white.withOpacity(0.05),
        borderRadius: BorderRadius.circular(16),
        border: Border.all(color: Colors.white.withOpacity(0.1)),
      ),
      child: Column(
        children: [
          const Text(
            '🔲 Code de l\'événement',
            style: TextStyle(
              color: Colors.white,
              fontSize: 16,
              fontWeight: FontWeight.w600,
            ),
          ),
          const SizedBox(height: 12),
          GestureDetector(
            onTap: () => _copyToClipboard(event['event_code'] ?? ''),
            child: Container(
              padding: const EdgeInsets.symmetric(horizontal: 20, vertical: 12),
              decoration: BoxDecoration(
                color: Colors.white.withOpacity(0.1),
                borderRadius: BorderRadius.circular(10),
              ),
              child: Row(
                mainAxisSize: MainAxisSize.min,
                children: [
                  Text(
                    event['event_code'] ?? 'N/A',
                    style: const TextStyle(
                      color: Colors.white,
                      fontSize: 20,
                      fontWeight: FontWeight.bold,
                      fontFamily: 'monospace',
                      letterSpacing: 2,
                    ),
                  ),
                  const SizedBox(width: 12),
                  const Icon(Icons.copy, color: Colors.white70, size: 20),
                ],
              ),
            ),
          ),
          const SizedBox(height: 12),
          SizedBox(
            width: double.infinity,
            child: ElevatedButton.icon(
              onPressed: () => _showQRCodeModal(event),
              icon: const Icon(Icons.qr_code),
              label: const Text('Voir QR Code & Télécharger'),
              style: ElevatedButton.styleFrom(
                backgroundColor: const Color(0xFF6366f1),
                foregroundColor: Colors.white,
                padding: const EdgeInsets.symmetric(vertical: 14),
                shape: RoundedRectangleBorder(
                  borderRadius: BorderRadius.circular(12),
                ),
              ),
            ),
          ),
          const SizedBox(height: 8),
          SizedBox(
            width: double.infinity,
            child: ElevatedButton.icon(
              onPressed: () => _openSoldTickets(event),
              icon: const Icon(Icons.confirmation_number_outlined),
              label: const Text('📊 Voir les tickets vendus'),
              style: ElevatedButton.styleFrom(
                backgroundColor: const Color(0xFF10b981),
                foregroundColor: Colors.white,
                padding: const EdgeInsets.symmetric(vertical: 14),
                shape: RoundedRectangleBorder(
                  borderRadius: BorderRadius.circular(12),
                ),
              ),
            ),
          ),
          const SizedBox(height: 8),
          SizedBox(
            width: double.infinity,
            child: ElevatedButton.icon(
              onPressed: () => _openTicketScanner(event),
              icon: const Icon(Icons.qr_code_scanner),
              label: const Text('📷 Scanner les tickets'),
              style: ElevatedButton.styleFrom(
                backgroundColor: const Color(0xFFf59e0b),
                foregroundColor: Colors.white,
                padding: const EdgeInsets.symmetric(vertical: 14),
                shape: RoundedRectangleBorder(
                  borderRadius: BorderRadius.circular(12),
                ),
              ),
            ),
          ),
        ],
      ),
    );
  }

        ],
      ),
    );
  }

  void _openTicketScanner(Map<String, dynamic> event) {
    Navigator.push(
      context,
      MaterialPageRoute(
        builder: (context) => TicketScannerScreen(
          eventId: widget.eventId,
          eventTitle: event['title'] ?? 'Événement',
        ),
      ),
    );
  }

  void _openSoldTickets(Map<String, dynamic> event) {
    Navigator.push(
      context,
      MaterialPageRoute(
        builder: (context) => SoldTicketsScreen(
          eventId: widget.eventId,
          eventTitle: event['title'] ?? 'Événement',
          currency: event['currency'] ?? 'XOF',
        ),
      ),
    );
  }

  void _showQRCodeModal(Map<String, dynamic> event) {
    showModalBottomSheet(
      context: context,
      isScrollControlled: true,
      backgroundColor: Colors.transparent,
      builder: (context) => Container(
        decoration: const BoxDecoration(
          color: Color(0xFF1a1a2e),
          borderRadius: BorderRadius.vertical(top: Radius.circular(24)),
        ),
        padding: const EdgeInsets.all(24),
        child: Column(
          mainAxisSize: MainAxisSize.min,
          children: [
            // Handle bar
            Container(
              width: 40,
              height: 4,
              decoration: BoxDecoration(
                color: Colors.white.withOpacity(0.3),
                borderRadius: BorderRadius.circular(2),
              ),
            ),
            const SizedBox(height: 24),
            
            // QR Code
            if (event['qr_code'] != null)
              Container(
                padding: const EdgeInsets.all(16),
                decoration: BoxDecoration(
                  color: Colors.white,
                  borderRadius: BorderRadius.circular(16),
                ),
                child: Image.network(
                  event['qr_code'],
                  width: 200,
                  height: 200,
                  fit: BoxFit.contain,
                ),
              )
            else
              Container(
                width: 200,
                height: 200,
                decoration: BoxDecoration(
                  color: Colors.white.withOpacity(0.1),
                  borderRadius: BorderRadius.circular(16),
                ),
                child: const Column(
                  mainAxisAlignment: MainAxisAlignment.center,
                  children: [
                    Icon(Icons.qr_code_2, size: 64, color: Colors.white54),
                    SizedBox(height: 8),
                    Text('QR non disponible', style: TextStyle(color: Colors.white54)),
                  ],
                ),
              ),
            
            const SizedBox(height: 20),
            
            // Event Info
            Container(
              padding: const EdgeInsets.all(16),
              decoration: BoxDecoration(
                color: Colors.white.withOpacity(0.05),
                borderRadius: BorderRadius.circular(12),
              ),
              child: Column(
                children: [
                  Row(
                    mainAxisAlignment: MainAxisAlignment.spaceBetween,
                    children: [
                      const Text('Événement', style: TextStyle(color: Colors.white54)),
                      Text(
                        event['title'] ?? '',
                        style: const TextStyle(color: Colors.white, fontWeight: FontWeight.w600),
                      ),
                    ],
                  ),
                  const Divider(color: Colors.white12, height: 20),
                  Row(
                    mainAxisAlignment: MainAxisAlignment.spaceBetween,
                    children: [
                      const Text('Statut', style: TextStyle(color: Colors.white54)),
                      Container(
                        padding: const EdgeInsets.symmetric(horizontal: 12, vertical: 4),
                        decoration: BoxDecoration(
                          color: _getStatusColor(event['status'] ?? 'draft'),
                          borderRadius: BorderRadius.circular(12),
                        ),
                        child: Text(
                          _getStatusLabel(event['status'] ?? 'draft'),
                          style: const TextStyle(color: Colors.white, fontSize: 12),
                        ),
                      ),
                    ],
                  ),
                ],
              ),
            ),
            
            const SizedBox(height: 20),
            
            // Action Buttons
            Row(
              children: [
                Expanded(
                  child: _buildActionButton(
                    icon: Icons.copy,
                    label: 'Code',
                    onTap: () {
                      _copyToClipboard(event['event_code'] ?? '');
                      Navigator.pop(context);
                    },
                  ),
                ),
                const SizedBox(width: 12),
                Expanded(
                  child: _buildActionButton(
                    icon: Icons.download,
                    label: 'DL PNG',
                    onTap: () {
                      ScaffoldMessenger.of(context).showSnackBar(
                        const SnackBar(content: Text('Téléchargement bientôt disponible')),
                      );
                    },
                  ),
                ),
                const SizedBox(width: 12),
                Expanded(
                  child: _buildActionButton(
                    icon: Icons.share,
                    label: 'Partager',
                    onTap: () => _shareEvent(event),
                  ),
                ),
              ],
            ),
            
            const SizedBox(height: 20),
            
            // Event Code
            Container(
              padding: const EdgeInsets.all(16),
              decoration: BoxDecoration(
                color: Colors.white.withOpacity(0.05),
                borderRadius: BorderRadius.circular(12),
              ),
              child: Column(
                children: [
                  const Text(
                    'CODE DE L\'ÉVÉNEMENT',
                    style: TextStyle(color: Colors.white54, fontSize: 11, letterSpacing: 1),
                  ),
                  const SizedBox(height: 8),
                  GestureDetector(
                    onTap: () => _copyToClipboard(event['event_code'] ?? ''),
                    child: Container(
                      padding: const EdgeInsets.symmetric(horizontal: 20, vertical: 12),
                      decoration: BoxDecoration(
                        border: Border.all(color: const Color(0xFF6366f1)),
                        borderRadius: BorderRadius.circular(10),
                      ),
                      child: Row(
                        mainAxisAlignment: MainAxisAlignment.center,
                        children: [
                          Text(
                            event['event_code'] ?? '',
                            style: const TextStyle(
                              color: Colors.white,
                              fontSize: 18,
                              fontWeight: FontWeight.bold,
                              fontFamily: 'monospace',
                              letterSpacing: 2,
                            ),
                          ),
                          const SizedBox(width: 12),
                          const Icon(Icons.copy, color: Colors.white70, size: 18),
                        ],
                      ),
                    ),
                  ),
                ],
              ),
            ),
            
            const SizedBox(height: 16),
          ],
        ),
      ),
    );
  }

  Widget _buildActionButton({
    required IconData icon,
    required String label,
    required VoidCallback onTap,
  }) {
    return GestureDetector(
      onTap: onTap,
      child: Container(
        padding: const EdgeInsets.symmetric(vertical: 16),
        decoration: BoxDecoration(
          color: Colors.white.withOpacity(0.05),
          borderRadius: BorderRadius.circular(12),
          border: Border.all(color: Colors.white.withOpacity(0.1)),
        ),
        child: Column(
          children: [
            Icon(icon, color: Colors.white, size: 28),
            const SizedBox(height: 4),
            Text(
              label,
              style: const TextStyle(color: Colors.white70, fontSize: 12),
            ),
          ],
        ),
      ),
    );
  }

  void _copyToClipboard(String text) {
    // Using services import for clipboard
    ScaffoldMessenger.of(context).showSnackBar(
      SnackBar(content: Text('Code copié: $text')),
    );
  }

  void _shareEvent(Map<String, dynamic> event) {
    final text = 'Rejoignez l\'événement "${event['title']}" - Code: ${event['event_code']}';
    ScaffoldMessenger.of(context).showSnackBar(
      SnackBar(content: Text('Partage: $text')),
    );
  }
}
