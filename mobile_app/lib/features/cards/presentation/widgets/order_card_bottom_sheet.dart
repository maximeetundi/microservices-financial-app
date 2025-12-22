import 'package:flutter/material.dart';

class OrderCardBottomSheet extends StatefulWidget {
  final VoidCallback? onOrder;

  const OrderCardBottomSheet({super.key, this.onOrder});

  @override
  State<OrderCardBottomSheet> createState() => _OrderCardBottomSheetState();
}

class _OrderCardBottomSheetState extends State<OrderCardBottomSheet> {
  String _selectedType = 'virtual';

  @override
  Widget build(BuildContext context) {
    return Container(
      padding: const EdgeInsets.all(24),
      child: Column(
        mainAxisSize: MainAxisSize.min,
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Center(
            child: Container(
              width: 40,
              height: 4,
              decoration: BoxDecoration(
                color: Colors.grey[300],
                borderRadius: BorderRadius.circular(2),
              ),
            ),
          ),
          const SizedBox(height: 24),
          
          const Text('Commander une carte', style: TextStyle(fontSize: 20, fontWeight: FontWeight.bold)),
          const SizedBox(height: 24),
          
          // Card type selection
          _CardTypeOption(
            type: 'virtual',
            title: 'Carte Virtuelle',
            description: 'Disponible instantanément pour les achats en ligne',
            price: 'Gratuite',
            isSelected: _selectedType == 'virtual',
            onTap: () => setState(() => _selectedType = 'virtual'),
          ),
          const SizedBox(height: 12),
          _CardTypeOption(
            type: 'physical',
            title: 'Carte Physique',
            description: 'Livrée sous 5-7 jours ouvrables',
            price: '\$9.99',
            isSelected: _selectedType == 'physical',
            onTap: () => setState(() => _selectedType = 'physical'),
          ),
          const SizedBox(height: 24),
          
          // Features
          const Text('Fonctionnalités incluses:', style: TextStyle(fontWeight: FontWeight.w500)),
          const SizedBox(height: 8),
          _FeatureItem(icon: Icons.check, text: 'Paiements sans contact'),
          _FeatureItem(icon: Icons.check, text: 'Notifications en temps réel'),
          _FeatureItem(icon: Icons.check, text: 'Contrôles de sécurité avancés'),
          _FeatureItem(icon: Icons.check, text: 'Limites personnalisables'),
          const SizedBox(height: 24),
          
          SizedBox(
            width: double.infinity,
            child: ElevatedButton(
              onPressed: () {
                widget.onOrder?.call();
                Navigator.pop(context);
                ScaffoldMessenger.of(context).showSnackBar(
                  const SnackBar(content: Text('Carte commandée avec succès!'), backgroundColor: Colors.green),
                );
              },
              style: ElevatedButton.styleFrom(padding: const EdgeInsets.symmetric(vertical: 16)),
              child: Text('Commander la carte ${_selectedType == 'physical' ? '(\$9.99)' : '(Gratuite)'}'),
            ),
          ),
        ],
      ),
    );
  }
}

class _CardTypeOption extends StatelessWidget {
  final String type;
  final String title;
  final String description;
  final String price;
  final bool isSelected;
  final VoidCallback onTap;

  const _CardTypeOption({
    required this.type,
    required this.title,
    required this.description,
    required this.price,
    required this.isSelected,
    required this.onTap,
  });

  @override
  Widget build(BuildContext context) {
    return InkWell(
      onTap: onTap,
      borderRadius: BorderRadius.circular(12),
      child: Container(
        padding: const EdgeInsets.all(16),
        decoration: BoxDecoration(
          border: Border.all(
            color: isSelected ? Theme.of(context).primaryColor : Colors.grey[300]!,
            width: isSelected ? 2 : 1,
          ),
          borderRadius: BorderRadius.circular(12),
          color: isSelected ? Theme.of(context).primaryColor.withOpacity(0.05) : null,
        ),
        child: Row(
          children: [
            Icon(
              type == 'virtual' ? Icons.credit_card : Icons.credit_card_outlined,
              color: isSelected ? Theme.of(context).primaryColor : Colors.grey,
              size: 32,
            ),
            const SizedBox(width: 16),
            Expanded(
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  Text(title, style: const TextStyle(fontWeight: FontWeight.bold)),
                  Text(description, style: TextStyle(fontSize: 12, color: Colors.grey[600])),
                ],
              ),
            ),
            Text(price, style: TextStyle(fontWeight: FontWeight.bold, color: Theme.of(context).primaryColor)),
          ],
        ),
      ),
    );
  }
}

class _FeatureItem extends StatelessWidget {
  final IconData icon;
  final String text;

  const _FeatureItem({required this.icon, required this.text});

  @override
  Widget build(BuildContext context) {
    return Padding(
      padding: const EdgeInsets.symmetric(vertical: 4),
      child: Row(
        children: [
          Icon(icon, color: Colors.green, size: 18),
          const SizedBox(width: 8),
          Text(text, style: const TextStyle(fontSize: 14)),
        ],
      ),
    );
  }
}
