import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';

class SettingsPage extends StatelessWidget {
  const SettingsPage({super.key});

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: const Text('Paramètres'),
      ),
      body: ListView(
        children: [
          // Profile Section
          Container(
            padding: const EdgeInsets.all(16),
            child: Row(
              children: [
                CircleAvatar(
                  radius: 35,
                  backgroundColor: Theme.of(context).primaryColor,
                  child: const Text('JD', style: TextStyle(fontSize: 24, color: Colors.white)),
                ),
                const SizedBox(width: 16),
                Expanded(
                  child: Column(
                    crossAxisAlignment: CrossAxisAlignment.start,
                    children: [
                      const Text('John Doe', style: TextStyle(fontSize: 20, fontWeight: FontWeight.bold)),
                      Text('john.doe@email.com', style: TextStyle(color: Colors.grey[600])),
                    ],
                  ),
                ),
                IconButton(
                  icon: const Icon(Icons.edit),
                  onPressed: () => context.push('/more/profile'),
                ),
              ],
            ),
          ),
          const Divider(),
          
          // Settings Groups
          _SettingsGroup(
            title: 'Compte',
            items: [
              _SettingsItem(icon: Icons.person, title: 'Profil', onTap: () => context.push('/more/profile')),
              _SettingsItem(icon: Icons.security, title: 'Sécurité', onTap: () => context.push('/more/security')),
              _SettingsItem(icon: Icons.notifications, title: 'Notifications', onTap: () {}),
              _SettingsItem(icon: Icons.language, title: 'Langue', subtitle: 'Français', onTap: () {}),
            ],
          ),
          
          _SettingsGroup(
            title: 'Préférences',
            items: [
              _SettingsItem(icon: Icons.dark_mode, title: 'Thème sombre', trailing: Switch(value: false, onChanged: (_) {})),
              _SettingsItem(icon: Icons.currency_exchange, title: 'Devise par défaut', subtitle: 'USD', onTap: () {}),
              _SettingsItem(icon: Icons.fingerprint, title: 'Biométrie', trailing: Switch(value: true, onChanged: (_) {})),
            ],
          ),
          
          _SettingsGroup(
            title: 'Services',
            items: [
              _SettingsItem(icon: Icons.credit_card, title: 'Mes cartes', onTap: () => context.push('/more/cards')),
              _SettingsItem(icon: Icons.send, title: 'Transferts', onTap: () => context.push('/more/transfer')),
              _SettingsItem(icon: Icons.storefront, title: 'Espace Marchand', onTap: () => context.push('/more/merchant')),
            ],
          ),
          
          _SettingsGroup(
            title: 'Support',
            items: [
              _SettingsItem(icon: Icons.help, title: 'Centre d\'aide', onTap: () => context.push('/more/support')),
              _SettingsItem(icon: Icons.chat, title: 'Contacter le support', onTap: () => context.push('/more/support')),
              _SettingsItem(icon: Icons.policy, title: 'Politique de confidentialité', onTap: () {}),
              _SettingsItem(icon: Icons.description, title: 'Conditions d\'utilisation', onTap: () {}),
            ],
          ),
          
          _SettingsGroup(
            title: 'Application',
            items: [
              _SettingsItem(icon: Icons.info, title: 'Version', subtitle: '1.0.0', onTap: () {}),
              _SettingsItem(icon: Icons.star, title: 'Noter l\'application', onTap: () {}),
            ],
          ),
          
          const SizedBox(height: 16),
          Padding(
            padding: const EdgeInsets.all(16),
            child: OutlinedButton.icon(
              onPressed: () {
                // Logout logic
                context.go('/auth/login');
              },
              icon: const Icon(Icons.logout, color: Colors.red),
              label: const Text('Déconnexion', style: TextStyle(color: Colors.red)),
              style: OutlinedButton.styleFrom(
                side: const BorderSide(color: Colors.red),
                padding: const EdgeInsets.symmetric(vertical: 12),
              ),
            ),
          ),
          const SizedBox(height: 32),
        ],
      ),
    );
  }
}

class _SettingsGroup extends StatelessWidget {
  final String title;
  final List<Widget> items;

  const _SettingsGroup({required this.title, required this.items});

  @override
  Widget build(BuildContext context) {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Padding(
          padding: const EdgeInsets.fromLTRB(16, 16, 16, 8),
          child: Text(title, style: TextStyle(fontSize: 14, fontWeight: FontWeight.bold, color: Colors.grey[600])),
        ),
        ...items,
      ],
    );
  }
}

class _SettingsItem extends StatelessWidget {
  final IconData icon;
  final String title;
  final String? subtitle;
  final Widget? trailing;
  final VoidCallback? onTap;

  const _SettingsItem({
    required this.icon,
    required this.title,
    this.subtitle,
    this.trailing,
    this.onTap,
  });

  @override
  Widget build(BuildContext context) {
    return ListTile(
      leading: Icon(icon),
      title: Text(title),
      subtitle: subtitle != null ? Text(subtitle!) : null,
      trailing: trailing ?? const Icon(Icons.chevron_right),
      onTap: onTap,
    );
  }
}
