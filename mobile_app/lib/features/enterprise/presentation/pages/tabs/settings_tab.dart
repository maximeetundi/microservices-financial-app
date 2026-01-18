import 'package:flutter/material.dart';
import '../../../../../core/services/api_service.dart';
import '../../../data/models/enterprise_model.dart';

class SettingsTab extends StatefulWidget {
  final Enterprise enterprise;
  final VoidCallback onRefresh;

  const SettingsTab({Key? key, required this.enterprise, required this.onRefresh}) : super(key: key);

  @override
  State<SettingsTab> createState() => _SettingsTabState();
}

class _SettingsTabState extends State<SettingsTab> {
  final ApiService _api = ApiService();
  bool _isSaving = false;
  
  late TextEditingController _nameController;
  late TextEditingController _descriptionController;

  @override
  void initState() {
    super.initState();
    _nameController = TextEditingController(text: widget.enterprise.name);
    _descriptionController = TextEditingController(text: widget.enterprise.description ?? '');
  }

  @override
  void dispose() {
    _nameController.dispose();
    _descriptionController.dispose();
    super.dispose();
  }

  Future<void> _saveSettings() async {
    setState(() => _isSaving = true);
    try {
      await _api.enterprise.updateEnterprise(widget.enterprise.id, {
        'name': _nameController.text,
        'description': _descriptionController.text,
      });
      widget.onRefresh();
      ScaffoldMessenger.of(context).showSnackBar(
        const SnackBar(content: Text('Paramètres enregistrés')),
      );
    } catch (e) {
      ScaffoldMessenger.of(context).showSnackBar(
        SnackBar(content: Text('Erreur: ${e.toString()}')),
      );
    } finally {
      setState(() => _isSaving = false);
    }
  }

  @override
  Widget build(BuildContext context) {
    return SingleChildScrollView(
      padding: const EdgeInsets.all(16),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          const Text(
            'Paramètres de l\'entreprise',
            style: TextStyle(fontSize: 20, fontWeight: FontWeight.bold),
          ),
          const SizedBox(height: 24),
          
          // Basic Info Card
          _SettingsCard(
            title: 'Informations générales',
            icon: Icons.business,
            children: [
              _InputField(
                label: 'Nom de l\'entreprise',
                controller: _nameController,
              ),
              const SizedBox(height: 16),
              _InputField(
                label: 'Description',
                controller: _descriptionController,
                maxLines: 3,
              ),
              const SizedBox(height: 16),
              _InfoRow(label: 'Type', value: widget.enterprise.type),
              _InfoRow(label: 'Statut', value: widget.enterprise.status),
            ],
          ),
          
          const SizedBox(height: 16),
          
          // Security Card
          _SettingsCard(
            title: 'Sécurité',
            icon: Icons.security,
            children: [
              _SettingsTile(
                icon: Icons.admin_panel_settings,
                title: 'Approbation multi-admin',
                subtitle: 'Requiert l\'approbation de plusieurs admins',
                trailing: Switch(value: true, onChanged: (v) {}),
              ),
              _SettingsTile(
                icon: Icons.lock,
                title: 'Vérification PIN',
                subtitle: 'PIN requis pour les actions sensibles',
                trailing: Switch(value: true, onChanged: (v) {}),
              ),
            ],
          ),
          
          const SizedBox(height: 16),
          
          // QR Code Card
          _SettingsCard(
            title: 'Codes QR',
            icon: Icons.qr_code,
            children: [
              _SettingsTile(
                icon: Icons.qr_code_2,
                title: 'QR Code Entreprise',
                subtitle: 'Scanner pour accéder aux services',
                onTap: _showQRCode,
              ),
            ],
          ),
          
          const SizedBox(height: 24),
          
          // Save Button
          SizedBox(
            width: double.infinity,
            child: ElevatedButton(
              onPressed: _isSaving ? null : _saveSettings,
              style: ElevatedButton.styleFrom(
                backgroundColor: Colors.blue,
                foregroundColor: Colors.white,
                padding: const EdgeInsets.symmetric(vertical: 16),
                shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(12)),
              ),
              child: _isSaving
                  ? const SizedBox(
                      width: 20,
                      height: 20,
                      child: CircularProgressIndicator(color: Colors.white, strokeWidth: 2),
                    )
                  : const Text('Enregistrer les modifications'),
            ),
          ),
          
          const SizedBox(height: 32),
          
          // Danger Zone
          Container(
            padding: const EdgeInsets.all(16),
            decoration: BoxDecoration(
              color: Colors.red.shade50,
              borderRadius: BorderRadius.circular(12),
              border: Border.all(color: Colors.red.shade100),
            ),
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                const Text(
                  'Zone de danger',
                  style: TextStyle(fontWeight: FontWeight.bold, color: Colors.red),
                ),
                const SizedBox(height: 8),
                TextButton.icon(
                  onPressed: () {},
                  icon: const Icon(Icons.delete, color: Colors.red, size: 18),
                  label: const Text('Supprimer l\'entreprise', style: TextStyle(color: Colors.red)),
                ),
              ],
            ),
          ),
        ],
      ),
    );
  }

  void _showQRCode() async {
    try {
      final qrData = await _api.enterprise.getEnterpriseQR(widget.enterprise.id);
      // Show QR code dialog
      showDialog(
        context: context,
        builder: (context) => AlertDialog(
          title: const Text('QR Code'),
          content: Column(
            mainAxisSize: MainAxisSize.min,
            children: [
              Container(
                width: 200,
                height: 200,
                color: Colors.grey[200],
                child: const Center(child: Icon(Icons.qr_code_2, size: 100)),
              ),
              const SizedBox(height: 16),
              Text('Scannez pour accéder à ${widget.enterprise.name}'),
            ],
          ),
          actions: [
            TextButton(
              onPressed: () => Navigator.pop(context),
              child: const Text('Fermer'),
            ),
          ],
        ),
      );
    } catch (e) {
      ScaffoldMessenger.of(context).showSnackBar(
        SnackBar(content: Text('Erreur: ${e.toString()}')),
      );
    }
  }
}

class _SettingsCard extends StatelessWidget {
  final String title;
  final IconData icon;
  final List<Widget> children;

  const _SettingsCard({
    required this.title,
    required this.icon,
    required this.children,
  });

  @override
  Widget build(BuildContext context) {
    return Container(
      decoration: BoxDecoration(
        color: Colors.white,
        borderRadius: BorderRadius.circular(16),
        boxShadow: [
          BoxShadow(
            color: Colors.grey.withOpacity(0.1),
            blurRadius: 10,
            offset: const Offset(0, 4),
          ),
        ],
      ),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Padding(
            padding: const EdgeInsets.all(16),
            child: Row(
              children: [
                Icon(icon, color: Colors.blue),
                const SizedBox(width: 12),
                Text(title, style: const TextStyle(fontWeight: FontWeight.bold, fontSize: 16)),
              ],
            ),
          ),
          const Divider(height: 1),
          Padding(
            padding: const EdgeInsets.all(16),
            child: Column(children: children),
          ),
        ],
      ),
    );
  }
}

class _InputField extends StatelessWidget {
  final String label;
  final TextEditingController controller;
  final int maxLines;

  const _InputField({
    required this.label,
    required this.controller,
    this.maxLines = 1,
  });

  @override
  Widget build(BuildContext context) {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Text(label, style: TextStyle(color: Colors.grey[600], fontSize: 13)),
        const SizedBox(height: 8),
        TextField(
          controller: controller,
          maxLines: maxLines,
          decoration: InputDecoration(
            border: OutlineInputBorder(borderRadius: BorderRadius.circular(12)),
            contentPadding: const EdgeInsets.symmetric(horizontal: 16, vertical: 12),
          ),
        ),
      ],
    );
  }
}

class _InfoRow extends StatelessWidget {
  final String label;
  final String value;

  const _InfoRow({required this.label, required this.value});

  @override
  Widget build(BuildContext context) {
    return Padding(
      padding: const EdgeInsets.symmetric(vertical: 4),
      child: Row(
        mainAxisAlignment: MainAxisAlignment.spaceBetween,
        children: [
          Text(label, style: TextStyle(color: Colors.grey[600])),
          Text(value, style: const TextStyle(fontWeight: FontWeight.w500)),
        ],
      ),
    );
  }
}

class _SettingsTile extends StatelessWidget {
  final IconData icon;
  final String title;
  final String subtitle;
  final Widget? trailing;
  final VoidCallback? onTap;

  const _SettingsTile({
    required this.icon,
    required this.title,
    required this.subtitle,
    this.trailing,
    this.onTap,
  });

  @override
  Widget build(BuildContext context) {
    return InkWell(
      onTap: onTap,
      child: Padding(
        padding: const EdgeInsets.symmetric(vertical: 8),
        child: Row(
          children: [
            Icon(icon, color: Colors.grey[600], size: 22),
            const SizedBox(width: 12),
            Expanded(
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  Text(title, style: const TextStyle(fontWeight: FontWeight.w500)),
                  Text(subtitle, style: TextStyle(color: Colors.grey[600], fontSize: 12)),
                ],
              ),
            ),
            if (trailing != null) trailing!,
            if (onTap != null && trailing == null)
              Icon(Icons.chevron_right, color: Colors.grey[400]),
          ],
        ),
      ),
    );
  }
}
