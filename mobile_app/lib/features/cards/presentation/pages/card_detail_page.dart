import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';

class CardDetailPage extends StatelessWidget {
  final String cardId;

  const CardDetailPage({super.key, required this.cardId});

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: const Text('Détails de la carte'),
        leading: IconButton(
          icon: const Icon(Icons.arrow_back),
          onPressed: () => context.pop(),
        ),
        actions: [
          IconButton(icon: const Icon(Icons.settings), onPressed: () {}),
        ],
      ),
      body: SingleChildScrollView(
        padding: const EdgeInsets.all(16),
        child: Column(
          children: [
            // Virtual Card
            Container(
              width: double.infinity,
              height: 200,
              decoration: BoxDecoration(
                gradient: LinearGradient(
                  colors: [Colors.purple[700]!, Colors.purple[400]!],
                  begin: Alignment.topLeft,
                  end: Alignment.bottomRight,
                ),
                borderRadius: BorderRadius.circular(20),
              ),
              padding: const EdgeInsets.all(24),
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                mainAxisAlignment: MainAxisAlignment.spaceBetween,
                children: [
                  Row(
                    mainAxisAlignment: MainAxisAlignment.spaceBetween,
                    children: [
                      const Text('CryptoBank', style: TextStyle(color: Colors.white, fontSize: 18, fontWeight: FontWeight.bold)),
                      const Icon(Icons.credit_card, color: Colors.white, size: 32),
                    ],
                  ),
                  const Text(
                    '**** **** **** 4532',
                    style: TextStyle(color: Colors.white, fontSize: 24, letterSpacing: 2),
                  ),
                  Row(
                    mainAxisAlignment: MainAxisAlignment.spaceBetween,
                    children: [
                      Column(
                        crossAxisAlignment: CrossAxisAlignment.start,
                        children: [
                          Text('TITULAIRE', style: TextStyle(color: Colors.white70, fontSize: 10)),
                          const Text('JOHN DOE', style: TextStyle(color: Colors.white, fontSize: 14)),
                        ],
                      ),
                      Column(
                        crossAxisAlignment: CrossAxisAlignment.start,
                        children: [
                          Text('EXPIRE', style: TextStyle(color: Colors.white70, fontSize: 10)),
                          const Text('12/28', style: TextStyle(color: Colors.white, fontSize: 14)),
                        ],
                      ),
                      const Icon(Icons.contactless, color: Colors.white),
                    ],
                  ),
                ],
              ),
            ),
            const SizedBox(height: 24),
            
            // Balance
            Container(
              width: double.infinity,
              padding: const EdgeInsets.all(20),
              decoration: BoxDecoration(
                color: Colors.grey[100],
                borderRadius: BorderRadius.circular(12),
              ),
              child: Column(
                children: [
                  const Text('Solde disponible', style: TextStyle(color: Colors.grey)),
                  const SizedBox(height: 8),
                  const Text('\$2,450.00', style: TextStyle(fontSize: 32, fontWeight: FontWeight.bold)),
                ],
              ),
            ),
            const SizedBox(height: 24),
            
            // Actions
            Row(
              children: [
                Expanded(
                  child: _ActionCard(icon: Icons.add, label: 'Recharger', onTap: () {}),
                ),
                const SizedBox(width: 12),
                Expanded(
                  child: _ActionCard(icon: Icons.lock, label: 'Bloquer', onTap: () {}),
                ),
                const SizedBox(width: 12),
                Expanded(
                  child: _ActionCard(icon: Icons.visibility, label: 'Voir PIN', onTap: () {}),
                ),
              ],
            ),
            const SizedBox(height: 24),
            
            // Settings
            const Align(
              alignment: Alignment.centerLeft,
              child: Text('Paramètres', style: TextStyle(fontSize: 18, fontWeight: FontWeight.bold)),
            ),
            const SizedBox(height: 16),
            
            _SettingsTile(icon: Icons.shopping_bag, title: 'Limites de dépenses', subtitle: '\$5,000/jour'),
            _SettingsTile(icon: Icons.language, title: 'Paiements en ligne', subtitle: 'Activé'),
            _SettingsTile(icon: Icons.atm, title: 'Retraits DAB', subtitle: 'Activé'),
            _SettingsTile(icon: Icons.notifications, title: 'Notifications', subtitle: 'Activées'),
          ],
        ),
      ),
    );
  }
}

class _ActionCard extends StatelessWidget {
  final IconData icon;
  final String label;
  final VoidCallback onTap;

  const _ActionCard({required this.icon, required this.label, required this.onTap});

  @override
  Widget build(BuildContext context) {
    return InkWell(
      onTap: onTap,
      borderRadius: BorderRadius.circular(12),
      child: Container(
        padding: const EdgeInsets.symmetric(vertical: 16),
        decoration: BoxDecoration(
          border: Border.all(color: Colors.grey[300]!),
          borderRadius: BorderRadius.circular(12),
        ),
        child: Column(
          children: [
            Icon(icon, color: Theme.of(context).primaryColor),
            const SizedBox(height: 8),
            Text(label, style: const TextStyle(fontSize: 12)),
          ],
        ),
      ),
    );
  }
}

class _SettingsTile extends StatelessWidget {
  final IconData icon;
  final String title;
  final String subtitle;

  const _SettingsTile({required this.icon, required this.title, required this.subtitle});

  @override
  Widget build(BuildContext context) {
    return ListTile(
      leading: Container(
        padding: const EdgeInsets.all(8),
        decoration: BoxDecoration(
          color: Colors.grey[100],
          borderRadius: BorderRadius.circular(8),
        ),
        child: Icon(icon, size: 20),
      ),
      title: Text(title),
      subtitle: Text(subtitle, style: TextStyle(color: Colors.grey[600])),
      trailing: const Icon(Icons.chevron_right),
      onTap: () {},
    );
  }
}
