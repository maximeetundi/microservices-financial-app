import 'package:flutter/material.dart';
import '../../core/services/ticket_api_service.dart';

class SoldTicketsScreen extends StatefulWidget {
  final String eventId;
  final String eventTitle;
  final String currency;

  const SoldTicketsScreen({
    super.key,
    required this.eventId,
    required this.eventTitle,
    required this.currency,
  });

  @override
  State<SoldTicketsScreen> createState() => _SoldTicketsScreenState();
}

class _SoldTicketsScreenState extends State<SoldTicketsScreen> {
  final TicketApiService _ticketApi = TicketApiService();
  List<dynamic> _tickets = [];
  bool _loading = true;
  String? _error;

  @override
  void initState() {
    super.initState();
    _loadTickets();
  }

  Future<void> _loadTickets() async {
    setState(() {
      _loading = true;
      _error = null;
    });
    try {
      final tickets = await _ticketApi.getEventTickets(widget.eventId, limit: 100);
      setState(() {
        _tickets = tickets;
        _loading = false;
      });
    } catch (e) {
      setState(() {
        _error = e.toString();
        _loading = false;
      });
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
              // Header
              Padding(
                padding: const EdgeInsets.all(16.0),
                child: Row(
                  children: [
                    IconButton(
                      onPressed: () => Navigator.pop(context),
                      icon: const Icon(Icons.arrow_back, color: Colors.white),
                    ),
                    Expanded(
                      child: Column(
                        crossAxisAlignment: CrossAxisAlignment.start,
                        children: [
                          const Text(
                            'Tickets Vendus',
                            style: TextStyle(
                              color: Colors.white,
                              fontSize: 20,
                              fontWeight: FontWeight.bold,
                            ),
                          ),
                          Text(
                            widget.eventTitle,
                            style: TextStyle(
                              color: Colors.white.withOpacity(0.7),
                              fontSize: 14,
                            ),
                          ),
                        ],
                      ),
                    ),
                      ),
                    ),
                    IconButton(
                        onPressed: _confirmCancelEvent,
                        icon: const Icon(Icons.cancel_outlined, color: Colors.red),
                        tooltip: 'Annuler l\'événement',
                    ),
                    Container(
                      padding: const EdgeInsets.symmetric(horizontal: 12, vertical: 6),
                      decoration: BoxDecoration(
                        color: const Color(0xFF6366f1).withOpacity(0.2),
                        borderRadius: BorderRadius.circular(20),
                        border: Border.all(color: const Color(0xFF6366f1).withOpacity(0.5)),
                      ),
                      child: Text(
                        '${_tickets.length}',
                        style: const TextStyle(
                          color: Color(0xFF6366f1),
                          fontWeight: FontWeight.bold,
                        ),
                      ),
                    ),
                  ],
                ),
              ),

              // Content
              Expanded(
                child: _loading
                    ? const Center(child: CircularProgressIndicator())
                    : _error != null
                        ? _buildError()
                        : _tickets.isEmpty
                            ? _buildEmptyState()
                            : _buildTicketsList(),
              ),
            ],
          ),
        ),
      ),
    );
  }

  Future<void> _confirmCancelEvent() async {
    final confirmed = await showDialog<bool>(
      context: context,
      builder: (context) => AlertDialog(
        backgroundColor: const Color(0xFF1e293b),
        title: const Text('Annuler l\'événement ?', style: TextStyle(color: Colors.white)),
        content: const Text(
          'ATTENTION : Cette action est IRRÉVERSIBLE. Tous les tickets vendus seront REMBOURSÉS automatiquement et l\'événement sera annulé.',
          style: TextStyle(color: Colors.white70),
        ),
        actions: [
          TextButton(
            onPressed: () => Navigator.pop(context, false),
            child: const Text('Retour', style: TextStyle(color: Colors.white)),
          ),
          TextButton(
            onPressed: () => Navigator.pop(context, true),
            child: const Text('Confirmer l\'annulation', style: TextStyle(color: Colors.red)),
          ),
        ],
      ),
    );

    if (confirmed == true) {
      _cancelEvent();
    }
  }

  Future<void> _cancelEvent() async {
    setState(() => _loading = true);
    try {
      await _ticketApi.cancelEvent(widget.eventId);
      if (mounted) {
        ScaffoldMessenger.of(context).showSnackBar(
          const SnackBar(content: Text('Événement annulé et remboursements initiés')),
        );
        _loadTickets(); // Reload to see status changes
      }
    } catch (e) {
      if (mounted) {
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(content: Text('Erreur: ${e.toString()}')),
        );
        setState(() => _loading = false);
      }
    }
  }

  void _showTicketDetails(dynamic ticket) {
    showModalBottomSheet(
      context: context,
      isScrollControlled: true,
      backgroundColor: Colors.transparent,
      builder: (context) => _TicketDetailsSheet(
        ticket: ticket,
        currency: widget.currency,
        onRefund: () => _refundTicket(ticket['id']),
      ),
    );
  }

  Future<void> _refundTicket(String ticketId) async {
    // Close sheet first? Or keep open and show loading?
    // Let's close sheet to assume fresh state or show loading dialog.
    Navigator.pop(context); // Close details sheet

    setState(() => _loading = true);
    try {
      await _ticketApi.refundTicket(ticketId);
      if (mounted) {
        ScaffoldMessenger.of(context).showSnackBar(
          const SnackBar(content: Text('Ticket remboursé avec succès')),
        );
        _loadTickets();
      }
    } catch (e) {
      if (mounted) {
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(content: Text('Erreur: ${e.toString()}')),
        );
        setState(() => _loading = false);
      }
    }
  }

  // ... (rest of the class)
  
  Widget _buildError() {
    return Center(
      child: Column(
        mainAxisAlignment: MainAxisAlignment.center,
        children: [
          const Icon(Icons.error_outline, size: 48, color: Colors.red),
          const SizedBox(height: 16),
          Text(
            'Erreur de chargement',
            style: TextStyle(color: Colors.white.withOpacity(0.8)),
          ),
          TextButton(
            onPressed: _loadTickets,
            child: const Text('Réessayer'),
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
          const Text('🎫', style: TextStyle(fontSize: 64)),
          const SizedBox(height: 16),
          Text(
            'Aucun ticket vendu',
            style: TextStyle(
              color: Colors.white,
              fontSize: 18,
              fontWeight: FontWeight.bold,
            ),
          ),
          const SizedBox(height: 8),
          Text(
            'Les tickets vendus apparaîtront ici',
            style: TextStyle(color: Colors.white.withOpacity(0.5)),
          ),
        ],
      ),
    );
  }

  Widget _buildTicketsList() {
    return ListView.builder(
      padding: const EdgeInsets.all(16),
      itemCount: _tickets.length,
      itemBuilder: (context, index) {
        final ticket = _tickets[index];
        return _buildTicketCard(ticket);
      },
    );
  }

  Widget _buildTicketCard(dynamic ticket) {
    final formData = ticket['form_data'] as Map<String, dynamic>? ?? {};
    final buyerName = formData['name'] ?? formData['nom'] ?? formData['full_name'] ?? 'Anonyme';
    final buyerEmail = formData['email'] ?? formData['Email'] ?? '';
    final status = ticket['status'] ?? 'pending';

    return GestureDetector(
      onTap: () => _showTicketDetails(ticket),
      child: Container(
        margin: const EdgeInsets.only(bottom: 12),
        padding: const EdgeInsets.all(16),
        decoration: BoxDecoration(
          color: Colors.white.withOpacity(0.05),
          borderRadius: BorderRadius.circular(16),
          border: Border.all(color: Colors.white.withOpacity(0.1)),
        ),
        child: Row(
          children: [
            // Avatar
            CircleAvatar(
              radius: 20,
              backgroundColor: const Color(0xFF6366f1).withOpacity(0.2),
              child: Text(
                buyerName.isNotEmpty ? buyerName[0].toUpperCase() : '?',
                style: const TextStyle(
                  color: Color(0xFF6366f1),
                  fontWeight: FontWeight.bold,
                ),
              ),
            ),
            const SizedBox(width: 12),
            
            // Info
            Expanded(
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  Text(
                    buyerName,
                    style: const TextStyle(
                      color: Colors.white,
                      fontWeight: FontWeight.w600,
                    ),
                  ),
                  if (buyerEmail.isNotEmpty)
                    Text(
                      buyerEmail,
                      style: TextStyle(
                        color: Colors.white.withOpacity(0.5),
                        fontSize: 12,
                      ),
                    ),
                  const SizedBox(height: 4),
                  Row(
                    children: [
                      Container(
                        padding: const EdgeInsets.symmetric(horizontal: 6, vertical: 2),
                        decoration: BoxDecoration(
                          color: _hexToColor(ticket['tier_color'] ?? '#6366f1').withOpacity(0.2),
                          borderRadius: BorderRadius.circular(4),
                        ),
                        child: Text(
                          '${ticket['tier_icon'] ?? ''} ${ticket['tier_name'] ?? 'Ticket'}',
                          style: TextStyle(
                            color: _hexToColor(ticket['tier_color'] ?? '#6366f1'),
                            fontSize: 10,
                            fontWeight: FontWeight.bold,
                          ),
                        ),
                      ),
                      const SizedBox(width: 8),
                      // Code
                      Text(
                        ticket['ticket_code'] ?? '',
                        style: TextStyle(
                          color: Colors.white.withOpacity(0.4),
                          fontSize: 10,
                          fontFamily: 'monospace',
                        ),
                      ),
                    ],
                  ),
                ],
              ),
            ),

            // Amount & Status
            Column(
              crossAxisAlignment: CrossAxisAlignment.end,
              children: [
                Text(
                  '${_formatAmount(ticket['price'])} ${widget.currency}',
                  style: const TextStyle(
                    color: Colors.white,
                    fontWeight: FontWeight.bold,
                    fontSize: 14,
                  ),
                ),
                const SizedBox(height: 4),
                _buildStatusBadge(status),
              ],
            ),
          ],
        ),
      ),
    );
  }

  Widget _buildStatusBadge(String status) {
    Color color;
    String label;

    switch (status) {
      case 'paid':
        color = Colors.green;
        label = 'Confirmé';
        break;
      case 'used':
        color = Colors.blue;
        label = 'Utilisé';
        break;
      case 'cancelled':
        color = Colors.red;
        label = 'Annulé';
        break;
      case 'refunded':
        color = Colors.grey;
        label = 'Remboursé';
        break;
      default:
        color = Colors.amber;
        label = status;
    }

    return Container(
      padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 2),
      decoration: BoxDecoration(
        color: color.withOpacity(0.2),
        borderRadius: BorderRadius.circular(12),
        border: Border.all(color: color.withOpacity(0.5), width: 0.5),
      ),
      child: Text(
        label,
        style: TextStyle(color: color, fontSize: 10, fontWeight: FontWeight.w500),
      ),
    );
  }

  void _showTicketDetails(dynamic ticket) {
    showModalBottomSheet(
      context: context,
      isScrollControlled: true,
      backgroundColor: Colors.transparent,
      builder: (context) => _TicketDetailsSheet(
        ticket: ticket, 
        currency: widget.currency,
        onRefund: () => _refundTicket(ticket['id']),
      ),
    );
  }

  Future<void> _refundTicket(String ticketId) async {
    Navigator.pop(context); // Close details sheet
    setState(() => _loading = true);
    try {
      await _ticketApi.refundTicket(ticketId);
      if (mounted) {
        ScaffoldMessenger.of(context).showSnackBar(
          const SnackBar(content: Text('Ticket remboursé avec succès')),
        );
        _loadTickets();
      }
    } catch (e) {
      if (mounted) {
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(content: Text('Erreur: ${e.toString()}')),
        );
        setState(() => _loading = false);
      }
    }
  }

  Color _hexToColor(String hex) {
    try {
      hex = hex.replaceFirst('#', '');
      if (hex.length == 6) hex = 'FF$hex';
      return Color(int.parse(hex, radix: 16));
    } catch (e) {
      return const Color(0xFF6366f1);
    }
  }

  String _formatAmount(dynamic amount) {
    final num = (amount is int) ? amount : (amount as double).toInt();
    return num.toString().replaceAllMapped(
      RegExp(r'(\d{1,3})(?=(\d{3})+(?!\d))'),
      (m) => '${m[1]} ',
    );
  }
}

class _TicketDetailsSheet extends StatelessWidget {
  final dynamic ticket;
  final String currency;
  final VoidCallback onRefund;

  const _TicketDetailsSheet({required this.ticket, required this.currency, required this.onRefund});

  @override
  Widget build(BuildContext context) {
    final formData = ticket['form_data'] as Map<String, dynamic>? ?? {};
    final status = ticket['status'] ?? 'pending';

    return Container(
      height: MediaQuery.of(context).size.height * 0.7,
      decoration: const BoxDecoration(
        color: Color(0xFF1a1a2e),
        borderRadius: BorderRadius.vertical(top: Radius.circular(24)),
      ),
      padding: const EdgeInsets.all(24),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          // Handle
          Center(
            child: Container(
              width: 40,
              height: 4,
              decoration: BoxDecoration(
                color: Colors.white.withOpacity(0.3),
                borderRadius: BorderRadius.circular(2),
              ),
            ),
          ),
          const SizedBox(height: 24),

          Row(
            mainAxisAlignment: MainAxisAlignment.spaceBetween,
            children: [
              const Text(
                'Détails du Ticket',
                style: TextStyle(
                  color: Colors.white,
                  fontSize: 20,
                  fontWeight: FontWeight.bold,
                ),
              ),
              IconButton(
                onPressed: () => Navigator.pop(context),
                icon: const Icon(Icons.close, color: Colors.white54),
              ),
            ],
          ),
          
          Expanded(
            child: SingleChildScrollView(
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  const SizedBox(height: 16),
                  
                  // Main Info Card
                  Container(
                    padding: const EdgeInsets.all(16),
                    decoration: BoxDecoration(
                      color: Colors.white.withOpacity(0.05),
                      borderRadius: BorderRadius.circular(16),
                    ),
                    child: Column(
                      children: [
                        _buildDetailRow('Code', ticket['ticket_code'] ?? '', isMono: true),
                        const Divider(color: Colors.white10),
                        _buildDetailRow('Prix', '${ticket['price']} $currency'),
                        const Divider(color: Colors.white10),
                        _buildDetailRow('Date', _formatDate(ticket['created_at'])),
                         const Divider(color: Colors.white10),
                        _buildDetailRow('Statut', status),
                      ],
                    ),
                  ),

                  // Refund Button
                  if (status == 'paid') ...[
                    const SizedBox(height: 16),
                    SizedBox(
                      width: double.infinity,
                      child: ElevatedButton.icon(
                        onPressed: () {
                          // Confirmation dialog
                          showDialog(
                            context: context,
                            builder: (context) => AlertDialog(
                               backgroundColor: const Color(0xFF1e293b),
                               title: const Text('Rembourser le ticket ?', style: TextStyle(color: Colors.white)),
                               content: const Text('Cette action est irréversible.', style: TextStyle(color: Colors.white70)),
                               actions: [
                                  TextButton(
                                    onPressed: () => Navigator.pop(context),
                                    child: const Text('Annuler', style: TextStyle(color: Colors.white)),
                                  ),
                                  TextButton(
                                    onPressed: () {
                                      Navigator.pop(context); // Close dialog
                                      onRefund(); // Call refund
                                    },
                                    child: const Text('Rembourser', style: TextStyle(color: Colors.red)),
                                  ),
                               ],
                            )
                          );
                        },
                        icon: const Icon(Icons.refresh_outlined, color: Colors.red),
                        label: const Text('Rembourser ce ticket', style: TextStyle(color: Colors.red)),
                        style: ElevatedButton.styleFrom(
                          backgroundColor: Colors.red.withOpacity(0.1),
                          foregroundColor: Colors.red,
                          shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(12)),
                          padding: const EdgeInsets.symmetric(vertical: 12),
                          side: BorderSide(color: Colors.red.withOpacity(0.3)),
                        ),
                      ),
                    ),
                  ],

                  const SizedBox(height: 24),
                  const Text(
                    'INFOS PARTICIPANT',
                    style: TextStyle(
                      color: Colors.white54,
                      fontSize: 12,
                      fontWeight: FontWeight.bold,
                      letterSpacing: 1,
                    ),
                  ),
                  const SizedBox(height: 12),

                  Container(
                    padding: const EdgeInsets.all(16),
                    decoration: BoxDecoration(
                      color: Colors.white.withOpacity(0.05),
                      borderRadius: BorderRadius.circular(16),
                    ),
                    child: Column(
                      children: [
                        if (formData.isEmpty)
                          const Text('Aucune donnée', style: TextStyle(color: Colors.white54)),
                        ...formData.entries.map((e) => Column(
                          children: [
                            _buildDetailRow(e.key, e.value.toString()),
                            if (e.key != formData.entries.last.key)
                              const Divider(color: Colors.white10),
                          ],
                        )).toList(),
                      ],
                    ),
                  ),

                   if (ticket['transaction_id'] != null) ...[
                    const SizedBox(height: 24),
                    Center(
                      child: Text(
                        'Ref: ${ticket['transaction_id']}',
                        style: TextStyle(
                          color: Colors.white.withOpacity(0.3),
                          fontSize: 10,
                          fontFamily: 'monospace',
                        ),
                      ),
                    ),
                  ],
                ],
              ),
            ),
          ),
        ],
      ),
    );
  }

  // ... (rest of _TicketDetailsSheet methods)
  Widget _buildDetailRow(String label, String value, {bool isMono = false}) {
    return Padding(
      padding: const EdgeInsets.symmetric(vertical: 8),
      child: Row(
        mainAxisAlignment: MainAxisAlignment.spaceBetween,
        children: [
          Text(
            label.toUpperCase(),
            style: const TextStyle(color: Colors.white54, fontSize: 12),
          ),
          Text(
            value,
            style: TextStyle(
              color: Colors.white,
              fontWeight: FontWeight.w600,
              fontFamily: isMono ? 'monospace' : null,
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
}
