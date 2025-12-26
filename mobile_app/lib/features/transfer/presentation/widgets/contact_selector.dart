import 'package:flutter/material.dart';
import 'package:google_fonts/google_fonts.dart';

import '../../../../core/theme/app_theme.dart';
import '../../../../core/widgets/custom_text_field.dart';
import '../bloc/transfer_bloc.dart';

/// Contact Selector Widget
/// Note: This displays an empty state since contacts are stored locally
/// and should be managed by the user. The UI is ready for future contact API integration.
class ContactSelector extends StatefulWidget {
  final Function(Contact) onContactSelected;

  const ContactSelector({
    Key? key,
    required this.onContactSelected,
  }) : super(key: key);

  @override
  State<ContactSelector> createState() => _ContactSelectorState();
}

class _ContactSelectorState extends State<ContactSelector> {
  final _searchController = TextEditingController();
  final _newContactNameController = TextEditingController();
  final _newContactEmailController = TextEditingController();
  final _newContactAddressController = TextEditingController();
  
  List<Contact> _filteredContacts = [];
  List<Contact> _allContacts = [];
  bool _isLoading = false;

  @override
  void initState() {
    super.initState();
    _loadContacts();
  }

  Future<void> _loadContacts() async {
    setState(() => _isLoading = true);
    
    // Contacts are stored locally - in a real app, this would fetch from SharedPreferences
    // or a local database. For now, we start with an empty list.
    // Users can add contacts manually which would be persisted locally.
    
    setState(() {
      _allContacts = [];
      _filteredContacts = [];
      _isLoading = false;
    });
  }

  void _filterContacts(String query) {
    setState(() {
      if (query.isEmpty) {
        _filteredContacts = _allContacts;
      } else {
        _filteredContacts = _allContacts
            .where((contact) =>
                contact.name.toLowerCase().contains(query.toLowerCase()) ||
                contact.email.toLowerCase().contains(query.toLowerCase()) ||
                contact.address.toLowerCase().contains(query.toLowerCase()))
            .toList();
      }
    });
  }

  void _addContact(Contact contact) {
    setState(() {
      _allContacts.add(contact);
      _filteredContacts = List.from(_allContacts);
    });
    // TODO: Persist to local storage
  }

  @override
  Widget build(BuildContext context) {
    final isDark = Theme.of(context).brightness == Brightness.dark;
    
    return Container(
      height: MediaQuery.of(context).size.height * 0.8,
      decoration: BoxDecoration(
        color: isDark ? const Color(0xFF1E293B) : Colors.white,
        borderRadius: const BorderRadius.only(
          topLeft: Radius.circular(20),
          topRight: Radius.circular(20),
        ),
      ),
      child: Column(
        children: [
          // Handle bar
          Container(
            width: 40,
            height: 4,
            margin: const EdgeInsets.only(top: 12),
            decoration: BoxDecoration(
              color: isDark ? Colors.white24 : Colors.grey.shade300,
              borderRadius: BorderRadius.circular(2),
            ),
          ),
          
          // Header
          Padding(
            padding: const EdgeInsets.all(20),
            child: Row(
              mainAxisAlignment: MainAxisAlignment.spaceBetween,
              children: [
                Text(
                  'Sélectionner un contact',
                  style: GoogleFonts.inter(
                    fontSize: 20,
                    fontWeight: FontWeight.bold,
                    color: isDark ? Colors.white : const Color(0xFF1E293B),
                  ),
                ),
                IconButton(
                  icon: Icon(
                    Icons.close,
                    color: isDark ? Colors.white70 : Colors.grey,
                  ),
                  onPressed: () => Navigator.pop(context),
                ),
              ],
            ),
          ),

          // Search Bar
          Padding(
            padding: const EdgeInsets.symmetric(horizontal: 20),
            child: SearchTextField(
              controller: _searchController,
              hint: 'Rechercher un contact...',
              onChanged: _filterContacts,
            ),
          ),

          const SizedBox(height: 16),

          // Quick Actions
          Padding(
            padding: const EdgeInsets.symmetric(horizontal: 20),
            child: Row(
              children: [
                Expanded(
                  child: _buildQuickAction(
                    icon: Icons.person_add,
                    label: 'Ajouter',
                    onTap: _showAddContactDialog,
                    isDark: isDark,
                  ),
                ),
                const SizedBox(width: 12),
                Expanded(
                  child: _buildQuickAction(
                    icon: Icons.qr_code_scanner,
                    label: 'Scanner QR',
                    onTap: () {
                      Navigator.pop(context);
                      ScaffoldMessenger.of(context).showSnackBar(
                        const SnackBar(content: Text('Scanner QR à venir')),
                      );
                    },
                    isDark: isDark,
                  ),
                ),
              ],
            ),
          ),

          const SizedBox(height: 24),

          // Contacts List
          Expanded(
            child: _isLoading
                ? const Center(child: CircularProgressIndicator())
                : _filteredContacts.isEmpty
                    ? _buildEmptyState(isDark)
                    : ListView.separated(
                        padding: const EdgeInsets.symmetric(horizontal: 20),
                        itemCount: _filteredContacts.length,
                        separatorBuilder: (context, index) => const SizedBox(height: 12),
                        itemBuilder: (context, index) {
                          final contact = _filteredContacts[index];
                          return _buildContactItem(contact, isDark);
                        },
                      ),
          ),
        ],
      ),
    );
  }

  Widget _buildQuickAction({
    required IconData icon,
    required String label,
    required VoidCallback onTap,
    required bool isDark,
  }) {
    return GestureDetector(
      onTap: onTap,
      child: Container(
        padding: const EdgeInsets.all(16),
        decoration: BoxDecoration(
          color: const Color(0xFF6366F1).withOpacity(0.1),
          borderRadius: BorderRadius.circular(12),
          border: Border.all(
            color: const Color(0xFF6366F1).withOpacity(0.2),
          ),
        ),
        child: Column(
          children: [
            Icon(
              icon,
              color: const Color(0xFF6366F1),
              size: 24,
            ),
            const SizedBox(height: 8),
            Text(
              label,
              style: GoogleFonts.inter(
                fontSize: 12,
                fontWeight: FontWeight.w500,
                color: isDark ? Colors.white : const Color(0xFF1E293B),
              ),
            ),
          ],
        ),
      ),
    );
  }

  Widget _buildContactItem(Contact contact, bool isDark) {
    return GestureDetector(
      onTap: () => widget.onContactSelected(contact),
      child: Container(
        padding: const EdgeInsets.all(16),
        decoration: BoxDecoration(
          color: isDark ? Colors.white.withOpacity(0.05) : Colors.white,
          borderRadius: BorderRadius.circular(12),
          border: Border.all(
            color: isDark ? Colors.white.withOpacity(0.1) : const Color(0xFFE2E8F0),
          ),
        ),
        child: Row(
          children: [
            CircleAvatar(
              radius: 24,
              backgroundColor: _getContactColor(contact.type),
              child: Text(
                contact.name.isNotEmpty ? contact.name[0].toUpperCase() : '?',
                style: GoogleFonts.inter(
                  color: Colors.white,
                  fontWeight: FontWeight.bold,
                ),
              ),
            ),
            const SizedBox(width: 16),
            Expanded(
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  Text(
                    contact.name,
                    style: GoogleFonts.inter(
                      fontSize: 16,
                      fontWeight: FontWeight.w600,
                      color: isDark ? Colors.white : const Color(0xFF1E293B),
                    ),
                  ),
                  const SizedBox(height: 4),
                  Text(
                    contact.email,
                    style: GoogleFonts.inter(
                      fontSize: 14,
                      color: isDark ? Colors.white54 : const Color(0xFF64748B),
                    ),
                  ),
                  const SizedBox(height: 4),
                  Container(
                    padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 4),
                    decoration: BoxDecoration(
                      color: _getContactColor(contact.type).withOpacity(0.1),
                      borderRadius: BorderRadius.circular(6),
                    ),
                    child: Text(
                      '${contact.currency} • ${_getContactTypeText(contact.type)}',
                      style: GoogleFonts.inter(
                        fontSize: 12,
                        color: _getContactColor(contact.type),
                        fontWeight: FontWeight.w500,
                      ),
                    ),
                  ),
                ],
              ),
            ),
            Icon(
              Icons.chevron_right,
              color: isDark ? Colors.white30 : Colors.grey,
            ),
          ],
        ),
      ),
    );
  }

  Widget _buildEmptyState(bool isDark) {
    return Center(
      child: Column(
        mainAxisAlignment: MainAxisAlignment.center,
        children: [
          Container(
            width: 80,
            height: 80,
            decoration: BoxDecoration(
              color: const Color(0xFF6366F1).withOpacity(0.1),
              shape: BoxShape.circle,
            ),
            child: const Icon(
              Icons.person_search,
              size: 40,
              color: Color(0xFF6366F1),
            ),
          ),
          const SizedBox(height: 16),
          Text(
            _searchController.text.isEmpty
                ? 'Aucun contact enregistré'
                : 'Aucun résultat',
            style: GoogleFonts.inter(
              fontSize: 18,
              fontWeight: FontWeight.w600,
              color: isDark ? Colors.white : const Color(0xFF1E293B),
            ),
          ),
          const SizedBox(height: 8),
          Text(
            'Ajoutez des contacts pour des transferts rapides',
            style: GoogleFonts.inter(
              fontSize: 14,
              color: isDark ? Colors.white54 : const Color(0xFF64748B),
            ),
            textAlign: TextAlign.center,
          ),
          const SizedBox(height: 24),
          ElevatedButton.icon(
            onPressed: _showAddContactDialog,
            icon: const Icon(Icons.person_add),
            label: const Text('Ajouter un contact'),
            style: ElevatedButton.styleFrom(
              backgroundColor: const Color(0xFF6366F1),
              foregroundColor: Colors.white,
              padding: const EdgeInsets.symmetric(horizontal: 24, vertical: 12),
              shape: RoundedRectangleBorder(
                borderRadius: BorderRadius.circular(12),
              ),
            ),
          ),
        ],
      ),
    );
  }

  Color _getContactColor(ContactType type) {
    switch (type) {
      case ContactType.crypto:
        return const Color(0xFFF7931A); // BTC orange
      case ContactType.fiat:
        return const Color(0xFF22C55E); // Green
      case ContactType.instant:
        return const Color(0xFF6366F1); // Indigo
    }
  }

  String _getContactTypeText(ContactType type) {
    switch (type) {
      case ContactType.crypto:
        return 'Crypto';
      case ContactType.fiat:
        return 'Banque';
      case ContactType.instant:
        return 'Instantané';
    }
  }

  void _showAddContactDialog() {
    final isDark = Theme.of(context).brightness == Brightness.dark;
    ContactType selectedType = ContactType.instant;
    String selectedCurrency = 'XOF';
    
    showDialog(
      context: context,
      builder: (dialogContext) => StatefulBuilder(
        builder: (context, setDialogState) => AlertDialog(
          backgroundColor: isDark ? const Color(0xFF1E293B) : Colors.white,
          shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(16)),
          title: Text(
            'Nouveau contact',
            style: GoogleFonts.inter(
              fontWeight: FontWeight.bold,
              color: isDark ? Colors.white : const Color(0xFF1E293B),
            ),
          ),
          content: SingleChildScrollView(
            child: Column(
              mainAxisSize: MainAxisSize.min,
              children: [
                CustomTextField(
                  controller: _newContactNameController,
                  label: 'Nom',
                  hint: 'Nom du contact',
                ),
                const SizedBox(height: 16),
                CustomTextField(
                  controller: _newContactEmailController,
                  label: 'Email ou téléphone',
                  hint: 'email@example.com',
                  keyboardType: TextInputType.emailAddress,
                ),
                const SizedBox(height: 16),
                CustomTextField(
                  controller: _newContactAddressController,
                  label: 'Adresse/Compte',
                  hint: 'Adresse wallet ou IBAN',
                ),
                const SizedBox(height: 16),
                // Type selector
                Row(
                  children: [
                    Text(
                      'Type:',
                      style: GoogleFonts.inter(
                        color: isDark ? Colors.white70 : const Color(0xFF64748B),
                      ),
                    ),
                    const SizedBox(width: 12),
                    ChoiceChip(
                      label: const Text('Instant'),
                      selected: selectedType == ContactType.instant,
                      onSelected: (v) => setDialogState(() => selectedType = ContactType.instant),
                    ),
                    const SizedBox(width: 8),
                    ChoiceChip(
                      label: const Text('Crypto'),
                      selected: selectedType == ContactType.crypto,
                      onSelected: (v) => setDialogState(() => selectedType = ContactType.crypto),
                    ),
                  ],
                ),
              ],
            ),
          ),
          actions: [
            TextButton(
              onPressed: () {
                _newContactNameController.clear();
                _newContactEmailController.clear();
                _newContactAddressController.clear();
                Navigator.pop(dialogContext);
              },
              child: Text(
                'Annuler',
                style: GoogleFonts.inter(color: Colors.grey),
              ),
            ),
            ElevatedButton(
              onPressed: () {
                if (_newContactNameController.text.isNotEmpty) {
                  final newContact = Contact(
                    id: DateTime.now().millisecondsSinceEpoch.toString(),
                    name: _newContactNameController.text,
                    email: _newContactEmailController.text,
                    address: _newContactAddressController.text,
                    currency: selectedCurrency,
                    type: selectedType,
                  );
                  _addContact(newContact);
                  _newContactNameController.clear();
                  _newContactEmailController.clear();
                  _newContactAddressController.clear();
                  Navigator.pop(dialogContext);
                }
              },
              style: ElevatedButton.styleFrom(
                backgroundColor: const Color(0xFF6366F1),
                foregroundColor: Colors.white,
              ),
              child: const Text('Ajouter'),
            ),
          ],
        ),
      ),
    );
  }

  @override
  void dispose() {
    _searchController.dispose();
    _newContactNameController.dispose();
    _newContactEmailController.dispose();
    _newContactAddressController.dispose();
    super.dispose();
  }
}