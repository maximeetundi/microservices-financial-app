import 'package:flutter/material.dart';
import '../../core/services/ticket_api_service.dart';

class MyTicketsScreen extends StatefulWidget {
  const MyTicketsScreen({super.key});

  @override
  State<MyTicketsScreen> createState() => _MyTicketsScreenState();
}

class _MyTicketsScreenState extends State<MyTicketsScreen> {
  final TicketApiService _ticketApi = TicketApiService();
  List<dynamic> _tickets = [];
  bool _loading = true;
  dynamic _selectedTicket;

  @override
  void initState() {
    super.initState();
    _loadTickets();
  }

  Future<void> _loadTickets() async {
    setState(() => _loading = true);
    try {
      _tickets = await _ticketApi.getMyTickets();
    } catch (e) {
      debugPrint('Error loading tickets: $e');
    } finally {
      setState(() => _loading = false);
    }
  }

  void _showTicketDetails(dynamic ticket) {
    setState(() => _selectedTicket = ticket);
    showModalBottomSheet(
      context: context,
      isScrollControlled: true,
      backgroundColor: Colors.transparent,
      builder: (_) => _buildTicketModal(),
    );
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
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              _buildHeader(),
              Expanded(
                child: _loading
                    ? const Center(child: CircularProgressIndicator())
                    : _tickets.isEmpty
                        ? _buildEmptyState()
                        : RefreshIndicator(
                            onRefresh: _loadTickets,
                            child: ListView.builder(
                              padding: const EdgeInsets.all(20),
                              itemCount: _tickets.length,
                              itemBuilder: (context, index) => _buildTicketCard(_tickets[index]),
                            ),
                          ),
              ),
            ],
          ),
        ),
      ),
    );
  }

  Widget _buildHeader() {
    return Padding(
      padding: const EdgeInsets.all(20),
      child: Row(
        children: [
          IconButton(
            onPressed: () => Navigator.pop(context),
            icon: const Icon(Icons.arrow_back, color: Colors.white),
          ),
          const SizedBox(width: 8),
          const Text(
            '🎟️ Mes Tickets',
            style: TextStyle(
              fontSize: 24,
              fontWeight: FontWeight.bold,
              color: Colors.white,
            ),
          ),
        ],
      ),
    );
  }

  Widget _buildTicketCard(dynamic ticket) {
    final status = ticket['status'] ?? 'pending';
    final isValid = status == 'paid';
    final isUsed = status == 'used';

    return GestureDetector(
      onTap: () => _showTicketDetails(ticket),
      child: Container(
        margin: const EdgeInsets.only(bottom: 16),
        decoration: BoxDecoration(
          color: Colors.white.withOpacity(0.05),
          borderRadius: BorderRadius.circular(16),
          border: Border.all(
            color: isUsed ? Colors.grey.withOpacity(0.3) : Colors.white.withOpacity(0.1),
          ),
        ),
        child: Row(
          children: [
            // Left side - Icon
            Container(
              width: 80,
              height: 100,
              decoration: BoxDecoration(
                color: isUsed 
                    ? Colors.grey.withOpacity(0.2) 
                    : const Color(0xFF6366f1).withOpacity(0.2),
                borderRadius: const BorderRadius.horizontal(left: Radius.circular(16)),
              ),
              child: Center(
                child: Text(
                  ticket['tier_icon'] ?? '🎫',
                  style: TextStyle(
                    fontSize: 36,
                    color: isUsed ? Colors.grey : null,
                  ),
                ),
              ),
            ),
            
            // Middle - Info
            Expanded(
              child: Padding(
                padding: const EdgeInsets.all(16),
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    Text(
                      ticket['event_title'] ?? 'Événement',
                      style: TextStyle(
                        color: isUsed ? Colors.grey : Colors.white,
                        fontSize: 16,
                        fontWeight: FontWeight.bold,
                      ),
                      maxLines: 1,
                      overflow: TextOverflow.ellipsis,
                    ),
                    const SizedBox(height: 4),
                    Text(
                      ticket['tier_name'] ?? 'Standard',
                      style: TextStyle(
                        color: isUsed ? Colors.grey : Colors.white.withOpacity(0.7),
                        fontSize: 14,
                      ),
                    ),
                    const SizedBox(height: 4),
                    Row(
                      children: [
                        Icon(
                          Icons.calendar_today,
                          size: 12,
                          color: Colors.white.withOpacity(0.5),
                        ),
                        const SizedBox(width: 4),
                        Text(
                          _formatDate(ticket['event_date']),
                          style: TextStyle(
                            color: Colors.white.withOpacity(0.5),
                            fontSize: 12,
                          ),
                        ),
                      ],
                    ),
                  ],
                ),
              ),
            ),
            
            // Right - Status
            Padding(
              padding: const EdgeInsets.all(16),
              child: Column(
                children: [
                  Container(
                    padding: const EdgeInsets.symmetric(horizontal: 12, vertical: 6),
                    decoration: BoxDecoration(
                      color: isValid 
                          ? Colors.green 
                          : isUsed 
                              ? Colors.grey 
                              : Colors.orange,
                      borderRadius: BorderRadius.circular(20),
                    ),
                    child: Text(
                      isValid ? 'Valide' : isUsed ? 'Utilisé' : status,
                      style: const TextStyle(
                        color: Colors.white,
                        fontSize: 12,
                        fontWeight: FontWeight.w600,
                      ),
                    ),
                  ),
                  const SizedBox(height: 8),
                  const Icon(Icons.qr_code, color: Colors.white54, size: 24),
                ],
              ),
            ),
          ],
        ),
      ),
    );
  }

  Widget _buildTicketModal() {
    if (_selectedTicket == null) return const SizedBox();
    
    final ticket = _selectedTicket!;
    final status = ticket['status'] ?? 'pending';
    final isValid = status == 'paid';

    return Container(
      height: MediaQuery.of(context).size.height * 0.75,
      decoration: const BoxDecoration(
        color: Color(0xFF1a1a2e),
        borderRadius: BorderRadius.vertical(top: Radius.circular(24)),
      ),
      child: Column(
        children: [
          // Handle
          Container(
            margin: const EdgeInsets.only(top: 12),
            width: 40,
            height: 4,
            decoration: BoxDecoration(
              color: Colors.white.withOpacity(0.3),
              borderRadius: BorderRadius.circular(2),
            ),
          ),
          
          Expanded(
            child: SingleChildScrollView(
              padding: const EdgeInsets.all(24),
              child: Column(
                children: [
                  // Event title
                  Text(
                    ticket['event_title'] ?? 'Événement',
                    style: const TextStyle(
                      color: Colors.white,
                      fontSize: 22,
                      fontWeight: FontWeight.bold,
                    ),
                    textAlign: TextAlign.center,
                  ),
                  const SizedBox(height: 8),
                  
                  // Tier info
                  Row(
                    mainAxisAlignment: MainAxisAlignment.center,
                    children: [
                      Text(
                        ticket['tier_icon'] ?? '🎫',
                        style: const TextStyle(fontSize: 24),
                      ),
                      const SizedBox(width: 8),
                      Text(
                        ticket['tier_name'] ?? 'Standard',
                        style: TextStyle(
                          color: Colors.white.withOpacity(0.8),
                          fontSize: 18,
                        ),
                      ),
                    ],
                  ),
                  const SizedBox(height: 24),
                  
                  // QR Code
                  if (ticket['qr_code'] != null)
                    Container(
                      padding: const EdgeInsets.all(16),
                      decoration: BoxDecoration(
                        color: Colors.white,
                        borderRadius: BorderRadius.circular(16),
                      ),
                      child: Image.network(
                        ticket['qr_code'],
                        width: 200,
                        height: 200,
                        fit: BoxFit.contain,
                        errorBuilder: (_, __, ___) => Container(
                          width: 200,
                          height: 200,
                          color: Colors.grey[200],
                          child: const Icon(Icons.qr_code, size: 100, color: Colors.grey),
                        ),
                      ),
                    ),
                  const SizedBox(height: 16),
                  
                  // Ticket code
                  Container(
                    padding: const EdgeInsets.symmetric(horizontal: 20, vertical: 12),
                    decoration: BoxDecoration(
                      color: Colors.white.withOpacity(0.1),
                      borderRadius: BorderRadius.circular(12),
                    ),
                    child: Text(
                      ticket['ticket_code'] ?? '',
                      style: const TextStyle(
                        color: Colors.white,
                        fontSize: 18,
                        fontWeight: FontWeight.bold,
                        fontFamily: 'monospace',
                        letterSpacing: 2,
                      ),
                    ),
                  ),
                  const SizedBox(height: 24),
                  
                  // Status
                  Container(
                    padding: const EdgeInsets.symmetric(horizontal: 24, vertical: 12),
                    decoration: BoxDecoration(
                      color: isValid ? Colors.green : Colors.grey,
                      borderRadius: BorderRadius.circular(30),
                    ),
                    child: Row(
                      mainAxisSize: MainAxisSize.min,
                      children: [
                        Icon(
                          isValid ? Icons.check_circle : Icons.cancel,
                          color: Colors.white,
                        ),
                        const SizedBox(width: 8),
                        Text(
                          isValid ? 'Ticket valide' : 'Ticket ${status}',
                          style: const TextStyle(
                            color: Colors.white,
                            fontSize: 16,
                            fontWeight: FontWeight.w600,
                          ),
                        ),
                      ],
                    ),
                  ),
                  const SizedBox(height: 24),
                  
                  // Event details
                  _buildDetailRow(Icons.location_on, ticket['event_location'] ?? 'Non défini'),
                  _buildDetailRow(Icons.calendar_today, _formatDate(ticket['event_date'])),
                  _buildDetailRow(Icons.confirmation_number, '${_formatAmount(ticket['price'] ?? 0)} ${ticket['currency'] ?? 'XOF'}'),
                ],
              ),
            ),
          ),
        ],
      ),
    );
  }

  Widget _buildDetailRow(IconData icon, String text) {
    return Padding(
      padding: const EdgeInsets.symmetric(vertical: 8),
      child: Row(
        children: [
          Icon(icon, color: Colors.white54, size: 20),
          const SizedBox(width: 12),
          Expanded(
            child: Text(
              text,
              style: TextStyle(
                color: Colors.white.withOpacity(0.8),
                fontSize: 15,
              ),
            ),
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
          const Text('🎟️', style: TextStyle(fontSize: 64)),
          const SizedBox(height: 16),
          const Text(
            'Aucun ticket',
            style: TextStyle(
              color: Colors.white,
              fontSize: 20,
              fontWeight: FontWeight.bold,
            ),
          ),
          const SizedBox(height: 8),
          Text(
            'Vos tickets achetés apparaîtront ici',
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
}
