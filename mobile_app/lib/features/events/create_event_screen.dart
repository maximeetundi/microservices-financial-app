import 'package:flutter/material.dart';
import '../../core/services/ticket_api_service.dart';

class CreateEventFormScreen extends StatefulWidget {
  const CreateEventFormScreen({super.key});

  @override
  State<CreateEventFormScreen> createState() => _CreateEventFormScreenState();
}

class _CreateEventFormScreenState extends State<CreateEventFormScreen> {
  final TicketApiService _ticketApi = TicketApiService();
  final _formKey = GlobalKey<FormState>();
  
  // Event data
  String _title = '';
  String _description = '';
  String _location = '';
  String _coverImage = '';
  DateTime? _startDate;
  DateTime? _endDate;
  DateTime? _saleStartDate;
  DateTime? _saleEndDate;
  
  // Form fields
  List<Map<String, dynamic>> _formFields = [
    {'name': 'full_name', 'label': 'Nom complet', 'type': 'text', 'required': true},
    {'name': 'email', 'label': 'Email', 'type': 'email', 'required': true},
  ];
  
  // Ticket tiers
  List<Map<String, dynamic>> _tiers = [
    {
      'name': 'Standard',
      'icon': '🎫',
      'price': 5000,
      'quantity': -1,
      'description': 'Accès standard',
      'color': '#6366f1'
    },
  ];
  
  List<String> _availableIcons = [];
  bool _loading = false;
  bool _saving = false;

  @override
  void initState() {
    super.initState();
    _loadIcons();
  }

  Future<void> _loadIcons() async {
    setState(() => _loading = true);
    try {
      _availableIcons = await _ticketApi.getAvailableIcons();
    } catch (e) {
      _availableIcons = ['⭐', '🌟', '✨', '💎', '👑', '🏆', '🎫', '🎟️', '🔥'];
    } finally {
      setState(() => _loading = false);
    }
  }

  Future<void> _createEvent() async {
    if (!_formKey.currentState!.validate()) return;
    
    if (_startDate == null || _endDate == null || _saleStartDate == null || _saleEndDate == null) {
      ScaffoldMessenger.of(context).showSnackBar(
        const SnackBar(content: Text('Veuillez remplir toutes les dates')),
      );
      return;
    }

    setState(() => _saving = true);
    try {
      await _ticketApi.createEvent({
        'title': _title,
        'description': _description,
        'location': _location,
        'cover_image': _coverImage.isNotEmpty ? _coverImage : null,
        'start_date': _startDate!.toIso8601String(),
        'end_date': _endDate!.toIso8601String(),
        'sale_start_date': _saleStartDate!.toIso8601String(),
        'sale_end_date': _saleEndDate!.toIso8601String(),
        'currency': 'XOF',
        'form_fields': _formFields,
        'ticket_tiers': _tiers.asMap().entries.map((e) => {
          ...e.value,
          'sort_order': e.key,
        }).toList(),
      });

      if (!mounted) return;
      ScaffoldMessenger.of(context).showSnackBar(
        const SnackBar(content: Text('✅ Événement créé avec succès!'), backgroundColor: Colors.green),
      );
      Navigator.pop(context);
    } catch (e) {
      if (!mounted) return;
      ScaffoldMessenger.of(context).showSnackBar(
        SnackBar(content: Text('Erreur: ${e.toString().replaceAll("Exception: ", "")}')),
      );
    } finally {
      if (mounted) setState(() => _saving = false);
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
                child: SingleChildScrollView(
                  padding: const EdgeInsets.all(20),
                  child: Form(
                    key: _formKey,
                    child: Column(
                      crossAxisAlignment: CrossAxisAlignment.start,
                      children: [
                        _buildBasicInfo(),
                        const SizedBox(height: 24),
                        _buildDates(),
                        const SizedBox(height: 24),
                        _buildFormFields(),
                        const SizedBox(height: 24),
                        _buildTiers(),
                        const SizedBox(height: 32),
                        _buildCreateButton(),
                      ],
                    ),
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
            '🎪 Créer un événement',
            style: TextStyle(
              fontSize: 22,
              fontWeight: FontWeight.bold,
              color: Colors.white,
            ),
          ),
        ],
      ),
    );
  }

  Widget _buildBasicInfo() {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        const Text(
          '📋 Informations générales',
          style: TextStyle(color: Colors.white, fontSize: 18, fontWeight: FontWeight.bold),
        ),
        const SizedBox(height: 12),
        _buildTextField('Titre de l\'événement', (v) => _title = v, required: true),
        _buildTextField('Description', (v) => _description = v, maxLines: 3),
        _buildTextField('Lieu', (v) => _location = v),
        _buildTextField('URL de l\'image de couverture', (v) => _coverImage = v),
      ],
    );
  }

  Widget _buildDates() {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        const Text(
          '📅 Dates',
          style: TextStyle(color: Colors.white, fontSize: 18, fontWeight: FontWeight.bold),
        ),
        const SizedBox(height: 12),
        _buildDatePicker('Début de l\'événement', _startDate, (d) => setState(() => _startDate = d)),
        _buildDatePicker('Fin de l\'événement', _endDate, (d) => setState(() => _endDate = d)),
        _buildDatePicker('Début des ventes', _saleStartDate, (d) => setState(() => _saleStartDate = d)),
        _buildDatePicker('Fin des ventes', _saleEndDate, (d) => setState(() => _saleEndDate = d)),
      ],
    );
  }

  Widget _buildFormFields() {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        const Text(
          '📝 Formulaire d\'inscription',
          style: TextStyle(color: Colors.white, fontSize: 18, fontWeight: FontWeight.bold),
        ),
        const SizedBox(height: 12),
        ..._formFields.asMap().entries.map((entry) {
          final idx = entry.key;
          final field = entry.value;
          return Container(
            margin: const EdgeInsets.only(bottom: 12),
            padding: const EdgeInsets.all(12),
            decoration: BoxDecoration(
              color: Colors.white.withOpacity(0.05),
              borderRadius: BorderRadius.circular(12),
            ),
            child: Row(
              children: [
                Expanded(
                  child: Text(
                    field['label'],
                    style: const TextStyle(color: Colors.white),
                  ),
                ),
                Text(
                  field['required'] ? 'Requis' : 'Optionnel',
                  style: TextStyle(color: Colors.white.withOpacity(0.6), fontSize: 12),
                ),
                IconButton(
                  onPressed: () => setState(() => _formFields.removeAt(idx)),
                  icon: const Icon(Icons.close, color: Colors.red, size: 20),
                ),
              ],
            ),
          );
        }),
        OutlinedButton.icon(
          onPressed: _addFormField,
          icon: const Icon(Icons.add),
          label: const Text('Ajouter un champ'),
          style: OutlinedButton.styleFrom(foregroundColor: Colors.white),
        ),
      ],
    );
  }

  Widget _buildTiers() {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        const Text(
          '🎫 Niveaux de tickets',
          style: TextStyle(color: Colors.white, fontSize: 18, fontWeight: FontWeight.bold),
        ),
        const SizedBox(height: 12),
        ..._tiers.asMap().entries.map((entry) {
          final idx = entry.key;
          final tier = entry.value;
          return Container(
            margin: const EdgeInsets.only(bottom: 12),
            padding: const EdgeInsets.all(16),
            decoration: BoxDecoration(
              color: Colors.white.withOpacity(0.05),
              borderRadius: BorderRadius.circular(12),
              border: Border.all(color: _hexToColor(tier['color']).withOpacity(0.5)),
            ),
            child: Column(
              children: [
                Row(
                  children: [
                    GestureDetector(
                      onTap: () => _selectIcon(idx),
                      child: Container(
                        width: 50,
                        height: 50,
                        decoration: BoxDecoration(
                          color: Colors.white.withOpacity(0.1),
                          borderRadius: BorderRadius.circular(10),
                        ),
                        child: Center(
                          child: Text(tier['icon'], style: const TextStyle(fontSize: 28)),
                        ),
                      ),
                    ),
                    const SizedBox(width: 12),
                    Expanded(
                      child: Text(
                        tier['name'],
                        style: const TextStyle(color: Colors.white, fontSize: 16, fontWeight: FontWeight.bold),
                      ),
                    ),
                    GestureDetector(
                      onTap: () => _selectColor(idx),
                      child: Container(
                        width: 30,
                        height: 30,
                        decoration: BoxDecoration(
                          color: _hexToColor(tier['color']),
                          shape: BoxShape.circle,
                        ),
                      ),
                    ),
                    IconButton(
                      onPressed: () => setState(() => _tiers.removeAt(idx)),
                      icon: const Icon(Icons.close, color: Colors.red, size: 20),
                    ),
                  ],
                ),
                const SizedBox(height: 12),
                Row(
                  children: [
                    Expanded(
                      child: Text(
                        '${tier['price']} XOF',
                        style: const TextStyle(color: Colors.white70),
                      ),
                    ),
                    Text(
                      tier['quantity'] == -1 ? 'Illimité' : '${tier['quantity']} tickets',
                      style: const TextStyle(color: Colors.white70),
                    ),
                  ],
                ),
              ],
            ),
          );
        }),
        OutlinedButton.icon(
          onPressed: _addTier,
          icon: const Icon(Icons.add),
          label: const Text('Ajouter un niveau'),
          style: OutlinedButton.styleFrom(foregroundColor: Colors.white),
        ),
      ],
    );
  }

  Widget _buildCreateButton() {
    return SizedBox(
      width: double.infinity,
      child: ElevatedButton(
        onPressed: _saving ? null : _createEvent,
        style: ElevatedButton.styleFrom(
          backgroundColor: const Color(0xFF6366f1),
          padding: const EdgeInsets.all(18),
          shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(12)),
        ),
        child: _saving
            ? const SizedBox(
                height: 20,
                width: 20,
                child: CircularProgressIndicator(strokeWidth: 2, color: Colors.white),
              )
            : const Text(
                '🎉 Créer l\'événement',
                style: TextStyle(fontSize: 18, fontWeight: FontWeight.bold),
              ),
      ),
    );
  }

  Widget _buildTextField(String label, Function(String) onChanged, {bool required = false, int maxLines = 1}) {
    return Padding(
      padding: const EdgeInsets.only(bottom: 16),
      child: TextFormField(
        style: const TextStyle(color: Colors.white),
        maxLines: maxLines,
        decoration: InputDecoration(
          labelText: label + (required ? ' *' : ''),
          labelStyle: TextStyle(color: Colors.white.withOpacity(0.7)),
          filled: true,
          fillColor: Colors.white.withOpacity(0.1),
          border: OutlineInputBorder(
            borderRadius: BorderRadius.circular(12),
            borderSide: BorderSide.none,
          ),
        ),
        validator: (v) => required && (v == null || v.isEmpty) ? 'Champ requis' : null,
        onChanged: onChanged,
      ),
    );
  }

  Widget _buildDatePicker(String label, DateTime? value, Function(DateTime) onChanged) {
    return Padding(
      padding: const EdgeInsets.only(bottom: 16),
      child: InkWell(
        onTap: () async {
          final date = await showDatePicker(
            context: context,
            initialDate: value ?? DateTime.now(),
            firstDate: DateTime.now(),
            lastDate: DateTime.now().add(const Duration(days: 365)),
          );
          if (date != null && mounted) {
            final time = await showTimePicker(context: context, initialTime: TimeOfDay.now());
            if (time != null) {
              onChanged(DateTime(date.year, date.month, date.day, time.hour, time.minute));
            }
          }
        },
        child: Container(
          padding: const EdgeInsets.all(16),
          decoration: BoxDecoration(
            color: Colors.white.withOpacity(0.1),
            borderRadius: BorderRadius.circular(12),
          ),
          child: Row(
            children: [
              const Icon(Icons.calendar_today, color: Colors.white70),
              const SizedBox(width: 12),
              Expanded(
                child: Text(
                  value == null ? label : '$label: ${_formatDate(value)}',
                  style: TextStyle(color: value == null ? Colors.white54 : Colors.white),
                ),
              ),
            ],
          ),
        ),
      ),
    );
  }

  void _addFormField() {
    showDialog(
      context: context,
      builder: (context) {
        String label = '';
        String type = 'text';
        bool required = false;
        return AlertDialog(
          title: const Text('Nouveau champ'),
          content: Column(
            mainAxisSize: MainAxisSize.min,
            children: [
              TextField(
                decoration: const InputDecoration(labelText: 'Nom du champ'),
                onChanged: (v) => label = v,
              ),
              DropdownButton<String>(
                value: type,
                items: ['text', 'email', 'phone', 'number'].map((t) =>
                  DropdownMenuItem(value: t, child: Text(t)),
                ).toList(),
                onChanged: (v) => type = v!,
              ),
              CheckboxListTile(
                title: const Text('Obligatoire'),
                value: required,
                onChanged: (v) => required = v!,
              ),
            ],
          ),
          actions: [
            TextButton(onPressed: () => Navigator.pop(context), child: const Text('Annuler')),
            ElevatedButton(
              onPressed: () {
                if (label.isNotEmpty) {
                  setState(() => _formFields.add({
                    'name': label.toLowerCase().replaceAll(' ', '_'),
                    'label': label,
                    'type': type,
                    'required': required,
                  }));
                  Navigator.pop(context);
                }
              },
              child: const Text('Ajouter'),
            ),
          ],
        );
      },
    );
  }

  void _addTier() {
    setState(() => _tiers.add({
      'name': 'Nouveau niveau',
      'icon': '🎟️',
      'price': 10000,
      'quantity': -1,
      'description': '',
      'color': '#8b5cf6',
    }));
  }

  void _selectIcon(int tierIndex) {
    showDialog(
      context: context,
      builder: (context) => AlertDialog(
        title: const Text('Choisir une icône'),
        content: SizedBox(
          width: 300,
          child: GridView.count(
            crossAxisCount: 5,
            shrinkWrap: true,
            children: _availableIcons.map((icon) => GestureDetector(
              onTap: () {
                setState(() => _tiers[tierIndex]['icon'] = icon);
                Navigator.pop(context);
              },
              child: Center(child: Text(icon, style: const TextStyle(fontSize: 32))),
            )).toList(),
          ),
        ),
      ),
    );
  }

  void _selectColor(int tierIndex) {
    final colors = ['#6366f1', '#8b5cf6', '#ec4899', '#f59e0b', '#10b981', '#3b82f6'];
    showDialog(
      context: context,
      builder: (context) => AlertDialog(
        title: const Text('Choisir une couleur'),
        content: Wrap(
          spacing: 12,
          children: colors.map((color) => GestureDetector(
            onTap: () {
              setState(() => _tiers[tierIndex]['color'] = color);
              Navigator.pop(context);
            },
            child: Container(
              width: 50,
              height: 50,
              decoration: BoxDecoration(
                color: _hexToColor(color),
                shape: BoxShape.circle,
              ),
            ),
          )).toList(),
        ),
      ),
    );
  }

  String _formatDate(DateTime date) {
    return '${date.day}/${date.month}/${date.year} ${date.hour}:${date.minute.toString().padLeft(2, '0')}';
  }

  Color _hexToColor(String hex) {
    hex = hex.replaceFirst('#', '');
    if (hex.length == 6) hex = 'FF$hex';
    return Color(int.parse(hex, radix: 16));
  }
}
