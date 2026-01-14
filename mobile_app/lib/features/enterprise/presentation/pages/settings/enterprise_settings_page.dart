import 'package:flutter/material.dart';
import '../../../../core/services/api_service.dart';

class EnterpriseSettingsPage extends StatefulWidget {
  final Map<String, dynamic> enterprise;

  const EnterpriseSettingsPage({Key? key, required this.enterprise}) : super(key: key);

  @override
  State<EnterpriseSettingsPage> createState() => _EnterpriseSettingsPageState();
}

class _EnterpriseSettingsPageState extends State<EnterpriseSettingsPage> {
  final ApiService _api = ApiService();
  final _formKey = GlobalKey<FormState>();
  
  late TextEditingController _nameController;
  late String _currency;
  bool _isSaving = false;

  // School Config
  List<Map<String, dynamic>> _classes = [];

  // Transport Config
  List<Map<String, dynamic>> _routes = [];

  // Custom Services (Plans)
  List<Map<String, dynamic>> _customServices = [];

  @override
  void initState() {
    super.initState();
    _nameController = TextEditingController(text: widget.enterprise['name']);
    
    final settings = widget.enterprise['settings'] ?? {};
    _currency = settings['currency'] ?? 'XOF';

    if (widget.enterprise['type'] == 'SCHOOL') {
      final config = widget.enterprise['school_config'] ?? {};
      if (config['classes'] != null) {
        _classes = List<Map<String, dynamic>>.from(config['classes']);
      }
    }

    if (widget.enterprise['type'] == 'TRANSPORT') {
      final config = widget.enterprise['transport_config'] ?? {};
      if (config['routes'] != null) {
         _routes = List<Map<String, dynamic>>.from(config['routes']);
      }
    }

    if (widget.enterprise['custom_services'] != null) {
      _customServices = List<Map<String, dynamic>>.from(widget.enterprise['custom_services']);
    }
  }

  @override
  void dispose() {
    _nameController.dispose();
    super.dispose();
  }

  Future<void> _saveSettings() async {
    if (!_formKey.currentState!.validate()) return;

    setState(() => _isSaving = true);

    try {
      final Map<String, dynamic> updates = {
        'name': _nameController.text,
        'settings': {
          ...(widget.enterprise['settings'] ?? {}),
          'currency': _currency,
        }
      };

      if (widget.enterprise['type'] == 'SCHOOL') {
        updates['school_config'] = {'classes': _classes};
      }

      if (widget.enterprise['type'] == 'TRANSPORT') {
        updates['transport_config'] = {'routes': _routes};
      }

      // Always update custom services
      updates['custom_services'] = _customServices;

      await _api.enterprise.updateEnterprise(widget.enterprise['id'] ?? widget.enterprise['_id'], updates);

      if (mounted) {
        ScaffoldMessenger.of(context).showSnackBar(
          const SnackBar(content: Text('Settings saved successfully')),
        );
        Navigator.pop(context, true); // Return true to trigger refresh
      }
    } catch (e) {
      if (mounted) {
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(content: Text('Error saving settings: $e')),
        );
      }
    } finally {
      if (mounted) {
        setState(() => _isSaving = false);
      }
    }
  }

  void _addClass() {
    setState(() {
      _classes.add({'name': '', 'total_fees': 0});
    });
  }

  void _addRoute() {
    setState(() {
      _routes.add({'name': '', 'base_price': 0});
    });
  }

  void _addCustomService() {
    setState(() {
      _customServices.add({
        'id': 'svc_${DateTime.now().millisecondsSinceEpoch}',
        'name': '',
        'billing_type': 'FIXED', 
        'billing_frequency': 'MONTHLY',
        'base_price': 0,
        'unit': ''
      });
    });
  }

  @override
  Widget build(BuildContext context) {
    final type = widget.enterprise['type'];

    return Scaffold(
      appBar: AppBar(
        title: const Text('Enterprise Settings'),
        actions: [
          IconButton(
            icon: _isSaving 
              ? const SizedBox(width: 20, height: 20, child: CircularProgressIndicator(color: Colors.white, strokeWidth: 2)) 
              : const Icon(Icons.check),
            onPressed: _isSaving ? null : _saveSettings,
          ),
        ],
      ),
      body: Form(
        key: _formKey,
        child: ListView(
          padding: const EdgeInsets.all(16),
          children: [
            _buildSectionHeader('General Settings'),
            const SizedBox(height: 8),
            TextFormField(
              controller: _nameController,
              decoration: const InputDecoration(labelText: 'Name', border: OutlineInputBorder()),
              validator: (v) => v!.isEmpty ? 'Required' : null,
            ),
            const SizedBox(height: 16),
            DropdownButtonFormField<String>(
              value: _currency,
              decoration: const InputDecoration(labelText: 'Currency', border: OutlineInputBorder()),
              items: ['XOF', 'EUR', 'USD'].map((c) => DropdownMenuItem(value: c, child: Text(c))).toList(),
              onChanged: (v) => setState(() => _currency = v!),
            ),

            if (type == 'SCHOOL') ...[
              const SizedBox(height: 24),
              _buildSectionHeader('🎓 School Configuration', action: TextButton(onPressed: _addClass, child: const Text('+ Add Class'))),
               ..._classes.asMap().entries.map((entry) {
                 final index = entry.key;
                 final cls = entry.value;
                 return Card(
                   margin: const EdgeInsets.only(bottom: 8),
                   child: Padding(
                     padding: const EdgeInsets.all(8.0),
                     child: Row(
                       children: [
                         Expanded(
                           child: TextFormField(
                             initialValue: cls['name'],
                             decoration: const InputDecoration(labelText: 'Class Name (e.g. CP)'),
                             onChanged: (v) => cls['name'] = v,
                           ),
                         ),
                         const SizedBox(width: 8),
                         SizedBox(
                           width: 100,
                           child: TextFormField(
                             initialValue: cls['total_fees']?.toString(),
                             decoration: const InputDecoration(labelText: 'Fees'),
                             keyboardType: TextInputType.number,
                             onChanged: (v) => cls['total_fees'] = double.tryParse(v) ?? 0,
                           ),
                         ),
                         IconButton(
                           icon: const Icon(Icons.delete, color: Colors.red),
                           onPressed: () => setState(() => _classes.removeAt(index)),
                         )
                       ],
                     ),
                   ),
                 );
               }).toList(),
               if (_classes.isEmpty) const Text('No classes configured', style: TextStyle(fontStyle: FontStyle.italic, color: Colors.grey)),
            ],

            if (type == 'TRANSPORT') ...[
              const SizedBox(height: 24),
              _buildSectionHeader('🚌 Transport Configuration', action: TextButton(onPressed: _addRoute, child: const Text('+ Add Route'))),
              ..._routes.asMap().entries.map((entry) {
                 final index = entry.key;
                 final route = entry.value;
                 return Card(
                   margin: const EdgeInsets.only(bottom: 8),
                   child: Padding(
                     padding: const EdgeInsets.all(8.0),
                     child: Row(
                       children: [
                         Expanded(
                           child: TextFormField(
                             initialValue: route['name'],
                             decoration: const InputDecoration(labelText: 'Route Name'),
                             onChanged: (v) => route['name'] = v,
                           ),
                         ),
                         const SizedBox(width: 8),
                         SizedBox(
                           width: 100,
                           child: TextFormField(
                             initialValue: route['base_price']?.toString(),
                             decoration: const InputDecoration(labelText: 'Price'),
                             keyboardType: TextInputType.number,
                             onChanged: (v) => route['base_price'] = double.tryParse(v) ?? 0,
                           ),
                         ),
                         IconButton(
                           icon: const Icon(Icons.delete, color: Colors.red),
                           onPressed: () => setState(() => _routes.removeAt(index)),
                         )
                       ],
                     ),
                   ),
                 );
               }).toList(),
               if (_routes.isEmpty) const Text('No routes configured', style: TextStyle(fontStyle: FontStyle.italic, color: Colors.grey)),
            ],

            const SizedBox(height: 24),
            _buildSectionHeader('⚡ Plans & Services', action: TextButton(onPressed: _addCustomService, child: const Text('+ Add Plan'))),
            if (_customServices.isEmpty) const Text('No plans configured', style: TextStyle(fontStyle: FontStyle.italic, color: Colors.grey)),
            ..._customServices.asMap().entries.map((entry) {
               final index = entry.key;
               final svc = entry.value;
               return Card(
                 margin: const EdgeInsets.only(bottom: 12),
                 elevation: 2,
                 child: Padding(
                   padding: const EdgeInsets.all(12.0),
                   child: Column(
                     children: [
                       Row(
                         children: [
                           Expanded(
                             child: TextFormField(
                               initialValue: svc['name'],
                               decoration: const InputDecoration(labelText: 'Plan Name (e.g. Gold Subscription)'),
                               onChanged: (v) => svc['name'] = v,
                             ),
                           ),
                           IconButton(
                             icon: const Icon(Icons.delete, color: Colors.red),
                             onPressed: () => setState(() => _customServices.removeAt(index)),
                           )
                         ],
                       ),
                       const SizedBox(height: 8),
                       Row(
                         children: [
                           Expanded(
                             child: DropdownButtonFormField<String>(
                               value: svc['billing_type'] ?? 'FIXED',
                               decoration: const InputDecoration(labelText: 'Type', contentPadding: EdgeInsets.symmetric(horizontal: 8, vertical: 0)),
                               items: ['FIXED', 'USAGE'].map((t) => DropdownMenuItem(value: t, child: Text(t))).toList(),
                               onChanged: (v) => setState(() => svc['billing_type'] = v),
                             ),
                           ),
                           const SizedBox(width: 8),
                           Expanded(
                             child: DropdownButtonFormField<String>(
                               value: svc['billing_frequency'] ?? 'MONTHLY',
                               decoration: const InputDecoration(labelText: 'Frequency', contentPadding: EdgeInsets.symmetric(horizontal: 8, vertical: 0)),
                               items: ['DAILY', 'WEEKLY', 'MONTHLY', 'QUARTERLY', 'ANNUALLY', 'CUSTOM', 'ONETIME']
                                   .map((f) => DropdownMenuItem(value: f, child: Text(f))).toList(),
                               onChanged: (v) => setState(() => svc['billing_frequency'] = v),
                             ),
                           ),
                         ],
                       ),
                       const SizedBox(height: 8),
                       if (svc['billing_frequency'] == 'CUSTOM') ...[
                         Row(
                           children: [
                             Checkbox(
                               value: svc['use_schedule'] ?? false,
                               onChanged: (v) {
                                 setState(() {
                                   svc['use_schedule'] = v;
                                   if (v == true && svc['payment_schedule'] == null) {
                                     svc['payment_schedule'] = [];
                                   }
                                 });
                               },
                             ),
                             const Text('Use Schedule (Dates)'),
                           ],
                         ),
                         
                         if (svc['use_schedule'] == true) ...[
                             ...(svc['payment_schedule'] as List? ?? []).asMap().entries.map((entry) {
                               final sIdx = entry.key;
                               final item = entry.value;
                               return Padding(
                                 padding: const EdgeInsets.only(left: 16.0, bottom: 8.0),
                                 child: Row(
                                   children: [
                                     Expanded(
                                       flex: 2,
                                       child: TextFormField(
                                         initialValue: item['name'],
                                         decoration: const InputDecoration(labelText: 'Name', isDense: true),
                                         onChanged: (v) => item['name'] = v,
                                       ),
                                     ),
                                     const SizedBox(width: 8),
                                     Expanded(
                                       child: TextFormField(
                                         initialValue: item['amount']?.toString(),
                                         decoration: const InputDecoration(labelText: 'Amt', isDense: true),
                                         keyboardType: TextInputType.number,
                                         onChanged: (v) => item['amount'] = double.tryParse(v) ?? 0,
                                       ),
                                     ),
                                     IconButton(
                                       icon: const Icon(Icons.close, size: 16, color: Colors.red),
                                       onPressed: () => setState(() => (svc['payment_schedule'] as List).removeAt(sIdx)),
                                     ),
                                   ],
                                 ),
                               );
                             }).toList(),
                             TextButton(
                               onPressed: () => setState(() {
                                 if (svc['payment_schedule'] == null) svc['payment_schedule'] = [];
                                 (svc['payment_schedule'] as List).add({'name': '', 'amount': 0, 'start_date': '', 'end_date': ''});
                               }),
                               child: const Text('+ Add Period'),
                             ),
                         ] else ...[
                           Row(
                              children: [
                                 Expanded(
                                   child: TextFormField(
                                     initialValue: svc['custom_interval']?.toString(),
                                     decoration: const InputDecoration(labelText: 'Interval (Days)'),
                                     keyboardType: TextInputType.number,
                                     onChanged: (v) => svc['custom_interval'] = int.tryParse(v) ?? 30,
                                   ),
                                 ),
                              ],
                           ),
                         ],
                         const SizedBox(height: 8),
                       ],
                       Row(
                         children: [
                           if (svc['billing_type'] == 'USAGE') ...[
                             Expanded(
                               child: TextFormField(
                                 initialValue: svc['unit'],
                                 decoration: const InputDecoration(labelText: 'Unit (e.g. kWh)'),
                                 onChanged: (v) => svc['unit'] = v,
                               ),
                             ),
                             const SizedBox(width: 8),
                           ],
                           Expanded(
                             child: TextFormField(
                               initialValue: svc['base_price']?.toString(),
                               decoration: const InputDecoration(labelText: 'Price'),
                               keyboardType: TextInputType.number,
                               onChanged: (v) => svc['base_price'] = double.tryParse(v) ?? 0,
                             ),
                           ),
                         ],
                       ),
                     ],
                   ),
                 ),
               );
             }).toList(),
          ],
        ),
      ),
    );
  }

  Widget _buildSectionHeader(String title, {Widget? action}) {
    return Row(
      mainAxisAlignment: MainAxisAlignment.spaceBetween,
      children: [
        Text(title, style: const TextStyle(fontSize: 18, fontWeight: FontWeight.bold, color: Colors.blueGrey)),
        if (action != null) action,
      ],
    );
  }
}
