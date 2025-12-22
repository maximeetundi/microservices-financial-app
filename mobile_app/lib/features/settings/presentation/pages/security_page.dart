import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';

class SecurityPage extends StatefulWidget {
  const SecurityPage({super.key});

  @override
  State<SecurityPage> createState() => _SecurityPageState();
}

class _SecurityPageState extends State<SecurityPage> {
  bool _biometricEnabled = true;
  bool _twoFactorEnabled = false;

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: const Text('Sécurité'),
        leading: IconButton(
          icon: const Icon(Icons.arrow_back),
          onPressed: () => context.pop(),
        ),
      ),
      body: ListView(
        children: [
          const SizedBox(height: 16),
          
          // Password Section
          _SecuritySection(
            title: 'Mot de passe',
            children: [
              ListTile(
                leading: const Icon(Icons.lock),
                title: const Text('Changer le mot de passe'),
                trailing: const Icon(Icons.chevron_right),
                onTap: () => _showChangePasswordDialog(context),
              ),
            ],
          ),
          
          // Biometric Section
          _SecuritySection(
            title: 'Biométrie',
            children: [
              SwitchListTile(
                secondary: const Icon(Icons.fingerprint),
                title: const Text('Empreinte / Face ID'),
                subtitle: const Text('Connexion rapide avec biométrie'),
                value: _biometricEnabled,
                onChanged: (v) => setState(() => _biometricEnabled = v),
              ),
            ],
          ),
          
          // 2FA Section
          _SecuritySection(
            title: 'Double authentification',
            children: [
              SwitchListTile(
                secondary: const Icon(Icons.security),
                title: const Text('Authentification à deux facteurs'),
                subtitle: const Text('Ajouter une couche de sécurité'),
                value: _twoFactorEnabled,
                onChanged: (v) => setState(() => _twoFactorEnabled = v),
              ),
              if (_twoFactorEnabled)
                ListTile(
                  leading: const Icon(Icons.qr_code),
                  title: const Text('Configurer l\'application'),
                  trailing: const Icon(Icons.chevron_right),
                  onTap: () {},
                ),
            ],
          ),
          
          // Session Section
          _SecuritySection(
            title: 'Sessions',
            children: [
              ListTile(
                leading: const Icon(Icons.devices),
                title: const Text('Appareils connectés'),
                subtitle: const Text('2 appareils actifs'),
                trailing: const Icon(Icons.chevron_right),
                onTap: () => _showDevicesDialog(context),
              ),
              ListTile(
                leading: const Icon(Icons.history),
                title: const Text('Historique de connexion'),
                trailing: const Icon(Icons.chevron_right),
                onTap: () {},
              ),
            ],
          ),
          
          const SizedBox(height: 16),
          Padding(
            padding: const EdgeInsets.all(16),
            child: OutlinedButton.icon(
              onPressed: () {
                // Disconnect all sessions
                ScaffoldMessenger.of(context).showSnackBar(
                  const SnackBar(content: Text('Toutes les sessions ont été déconnectées')),
                );
              },
              icon: const Icon(Icons.logout, color: Colors.orange),
              label: const Text('Déconnecter tous les appareils', style: TextStyle(color: Colors.orange)),
              style: OutlinedButton.styleFrom(
                side: const BorderSide(color: Colors.orange),
                padding: const EdgeInsets.symmetric(vertical: 12),
              ),
            ),
          ),
        ],
      ),
    );
  }

  void _showChangePasswordDialog(BuildContext context) {
    showDialog(
      context: context,
      builder: (context) => AlertDialog(
        title: const Text('Changer le mot de passe'),
        content: Column(
          mainAxisSize: MainAxisSize.min,
          children: [
            TextField(
              obscureText: true,
              decoration: const InputDecoration(labelText: 'Mot de passe actuel'),
            ),
            const SizedBox(height: 16),
            TextField(
              obscureText: true,
              decoration: const InputDecoration(labelText: 'Nouveau mot de passe'),
            ),
            const SizedBox(height: 16),
            TextField(
              obscureText: true,
              decoration: const InputDecoration(labelText: 'Confirmer'),
            ),
          ],
        ),
        actions: [
          TextButton(onPressed: () => Navigator.pop(context), child: const Text('Annuler')),
          ElevatedButton(onPressed: () => Navigator.pop(context), child: const Text('Confirmer')),
        ],
      ),
    );
  }

  void _showDevicesDialog(BuildContext context) {
    showDialog(
      context: context,
      builder: (context) => AlertDialog(
        title: const Text('Appareils connectés'),
        content: Column(
          mainAxisSize: MainAxisSize.min,
          children: [
            ListTile(
              leading: const Icon(Icons.phone_android),
              title: const Text('iPhone 15 Pro'),
              subtitle: const Text('Appareil actuel • Paris, France'),
              trailing: Container(
                padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 4),
                decoration: BoxDecoration(
                  color: Colors.green[100],
                  borderRadius: BorderRadius.circular(12),
                ),
                child: const Text('Actif', style: TextStyle(color: Colors.green, fontSize: 12)),
              ),
            ),
            ListTile(
              leading: const Icon(Icons.laptop),
              title: const Text('Chrome - Windows'),
              subtitle: const Text('Dernière activité: il y a 2h'),
              trailing: IconButton(
                icon: const Icon(Icons.close, color: Colors.red),
                onPressed: () {},
              ),
            ),
          ],
        ),
        actions: [
          TextButton(onPressed: () => Navigator.pop(context), child: const Text('Fermer')),
        ],
      ),
    );
  }
}

class _SecuritySection extends StatelessWidget {
  final String title;
  final List<Widget> children;

  const _SecuritySection({required this.title, required this.children});

  @override
  Widget build(BuildContext context) {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Padding(
          padding: const EdgeInsets.fromLTRB(16, 16, 16, 8),
          child: Text(title, style: TextStyle(fontSize: 14, fontWeight: FontWeight.bold, color: Colors.grey[600])),
        ),
        ...children,
      ],
    );
  }
}
