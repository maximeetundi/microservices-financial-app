import 'package:flutter/material.dart';
import '../../data/models/shop_models.dart';
import '../../data/repositories/shop_repository.dart';

class CartPage extends StatefulWidget {
  final Shop shop;
  final List<CartItem> cart;
  final VoidCallback onClearCart;

  const CartPage({
    super.key,
    required this.shop,
    required this.cart,
    required this.onClearCart,
  });

  @override
  State<CartPage> createState() => _CartPageState();
}

class _CartPageState extends State<CartPage> {
  final ShopRepository _repository = ShopRepository(
    baseUrl: const String.fromEnvironment('API_URL', defaultValue: 'https://api.tech-afm.com'),
  );

  String? _selectedWalletId;
  String _deliveryType = 'pickup';
  bool _processing = false;
  List<dynamic> _wallets = [];

  @override
  void initState() {
    super.initState();
    _loadWallets();
  }

  Future<void> _loadWallets() async {
    // TODO: Load wallets from wallet service
    // For now, we'll leave this empty as it requires wallet API integration
  }

  double get _subtotal => widget.cart.fold(0, (sum, item) => sum + (item.product.price * item.quantity));

  String _formatPrice(double amount, String currency) {
    return '${amount.toStringAsFixed(0)} $currency';
  }

  Future<void> _checkout() async {
    if (_selectedWalletId == null) {
      ScaffoldMessenger.of(context).showSnackBar(
        const SnackBar(content: Text('Veuillez sélectionner un portefeuille')),
      );
      return;
    }

    setState(() => _processing = true);
    try {
      final order = await _repository.createOrder(
        shopId: widget.shop.id,
        items: widget.cart,
        walletId: _selectedWalletId!,
        deliveryType: _deliveryType,
      );
      widget.onClearCart();
      if (mounted) {
        Navigator.pop(context);
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(content: Text('Commande ${order.orderNumber} créée !')),
        );
      }
    } catch (e) {
      ScaffoldMessenger.of(context).showSnackBar(
        SnackBar(content: Text('Erreur: $e'), backgroundColor: Colors.red),
      );
    } finally {
      setState(() => _processing = false);
    }
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: const Color(0xFFF8FAFC),
      appBar: AppBar(
        title: const Text('🛒 Mon Panier', style: TextStyle(fontWeight: FontWeight.bold)),
        backgroundColor: Colors.white,
        foregroundColor: Colors.black87,
        elevation: 0,
      ),
      body: widget.cart.isEmpty
          ? _buildEmptyCart()
          : Column(
              children: [
                // Shop Info
                Container(
                  padding: const EdgeInsets.all(16),
                  color: Colors.white,
                  child: Row(
                    children: [
                      Container(
                        width: 40,
                        height: 40,
                        decoration: BoxDecoration(
                          color: Colors.indigo,
                          borderRadius: BorderRadius.circular(8),
                        ),
                        child: Center(
                          child: Text(
                            widget.shop.name[0],
                            style: const TextStyle(color: Colors.white, fontWeight: FontWeight.bold, fontSize: 18),
                          ),
                        ),
                      ),
                      const SizedBox(width: 12),
                      Expanded(
                        child: Column(
                          crossAxisAlignment: CrossAxisAlignment.start,
                          children: [
                            Text(widget.shop.name, style: const TextStyle(fontWeight: FontWeight.bold)),
                            Text('${widget.cart.length} article(s)', style: TextStyle(color: Colors.grey[600], fontSize: 13)),
                          ],
                        ),
                      ),
                    ],
                  ),
                ),

                // Cart Items
                Expanded(
                  child: ListView.separated(
                    padding: const EdgeInsets.all(16),
                    itemCount: widget.cart.length,
                    separatorBuilder: (_, __) => const SizedBox(height: 12),
                    itemBuilder: (context, index) => _buildCartItem(widget.cart[index]),
                  ),
                ),

                // Checkout Section
                Container(
                  padding: const EdgeInsets.all(20),
                  decoration: BoxDecoration(
                    color: Colors.white,
                    boxShadow: [BoxShadow(color: Colors.black.withOpacity(0.05), blurRadius: 10, offset: const Offset(0, -4))],
                  ),
                  child: SafeArea(
                    child: Column(
                      children: [
                        // Delivery Type
                        Row(
                          children: [
                            Expanded(
                              child: GestureDetector(
                                onTap: () => setState(() => _deliveryType = 'pickup'),
                                child: Container(
                                  padding: const EdgeInsets.all(12),
                                  decoration: BoxDecoration(
                                    border: Border.all(
                                      color: _deliveryType == 'pickup' ? Colors.indigo : Colors.grey[300]!,
                                      width: _deliveryType == 'pickup' ? 2 : 1,
                                    ),
                                    borderRadius: BorderRadius.circular(12),
                                  ),
                                  child: Column(
                                    children: [
                                      const Text('🏃', style: TextStyle(fontSize: 24)),
                                      const SizedBox(height: 4),
                                      Text('Retrait', style: TextStyle(fontWeight: _deliveryType == 'pickup' ? FontWeight.bold : FontWeight.normal)),
                                    ],
                                  ),
                                ),
                              ),
                            ),
                            const SizedBox(width: 12),
                            Expanded(
                              child: GestureDetector(
                                onTap: () => setState(() => _deliveryType = 'delivery'),
                                child: Container(
                                  padding: const EdgeInsets.all(12),
                                  decoration: BoxDecoration(
                                    border: Border.all(
                                      color: _deliveryType == 'delivery' ? Colors.indigo : Colors.grey[300]!,
                                      width: _deliveryType == 'delivery' ? 2 : 1,
                                    ),
                                    borderRadius: BorderRadius.circular(12),
                                  ),
                                  child: Column(
                                    children: [
                                      const Text('🚚', style: TextStyle(fontSize: 24)),
                                      const SizedBox(height: 4),
                                      Text('Livraison', style: TextStyle(fontWeight: _deliveryType == 'delivery' ? FontWeight.bold : FontWeight.normal)),
                                    ],
                                  ),
                                ),
                              ),
                            ),
                          ],
                        ),

                        const SizedBox(height: 16),

                        // Total Row
                        Row(
                          mainAxisAlignment: MainAxisAlignment.spaceBetween,
                          children: [
                            const Text('Total', style: TextStyle(fontSize: 16)),
                            Text(
                              _formatPrice(_subtotal, widget.shop.currency),
                              style: TextStyle(fontSize: 24, fontWeight: FontWeight.bold, color: Colors.indigo.shade600),
                            ),
                          ],
                        ),

                        const SizedBox(height: 16),

                        // Checkout Button
                        SizedBox(
                          width: double.infinity,
                          child: ElevatedButton(
                            onPressed: _processing ? null : _checkout,
                            style: ElevatedButton.styleFrom(
                              backgroundColor: Colors.indigo,
                              padding: const EdgeInsets.symmetric(vertical: 16),
                              shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(12)),
                            ),
                            child: _processing
                                ? const SizedBox(width: 20, height: 20, child: CircularProgressIndicator(strokeWidth: 2, color: Colors.white))
                                : const Text('Payer maintenant', style: TextStyle(fontSize: 16, fontWeight: FontWeight.bold)),
                          ),
                        ),
                      ],
                    ),
                  ),
                ),
              ],
            ),
    );
  }

  Widget _buildEmptyCart() {
    return Center(
      child: Column(
        mainAxisAlignment: MainAxisAlignment.center,
        children: [
          const Text('🛒', style: TextStyle(fontSize: 64)),
          const SizedBox(height: 16),
          Text('Votre panier est vide', style: TextStyle(fontSize: 18, fontWeight: FontWeight.bold, color: Colors.grey[800])),
          const SizedBox(height: 24),
          ElevatedButton(
            onPressed: () => Navigator.pop(context),
            style: ElevatedButton.styleFrom(backgroundColor: Colors.indigo),
            child: const Text('Continuer mes achats'),
          ),
        ],
      ),
    );
  }

  Widget _buildCartItem(CartItem item) {
    return Container(
      padding: const EdgeInsets.all(12),
      decoration: BoxDecoration(
        color: Colors.white,
        borderRadius: BorderRadius.circular(12),
        boxShadow: [BoxShadow(color: Colors.black.withOpacity(0.05), blurRadius: 8)],
      ),
      child: Row(
        children: [
          // Image
          Container(
            width: 64,
            height: 64,
            decoration: BoxDecoration(
              color: Colors.grey[100],
              borderRadius: BorderRadius.circular(8),
            ),
            child: item.product.images.isNotEmpty
                ? ClipRRect(
                    borderRadius: BorderRadius.circular(8),
                    child: Image.network(item.product.images[0], fit: BoxFit.cover),
                  )
                : const Center(child: Text('📦', style: TextStyle(fontSize: 28))),
          ),
          const SizedBox(width: 12),

          // Info
          Expanded(
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Text(item.product.name, style: const TextStyle(fontWeight: FontWeight.w600)),
                const SizedBox(height: 4),
                Text(
                  _formatPrice(item.product.price, item.product.currency),
                  style: TextStyle(color: Colors.indigo.shade600, fontWeight: FontWeight.bold),
                ),
              ],
            ),
          ),

          // Quantity
          Container(
            decoration: BoxDecoration(
              color: Colors.grey[100],
              borderRadius: BorderRadius.circular(8),
            ),
            child: Row(
              children: [
                IconButton(
                  icon: const Icon(Icons.remove, size: 18),
                  onPressed: () {
                    if (item.quantity > 1) {
                      setState(() {
                        final index = widget.cart.indexOf(item);
                        widget.cart[index] = CartItem(
                          productId: item.productId,
                          quantity: item.quantity - 1,
                          product: item.product,
                        );
                      });
                    } else {
                      setState(() => widget.cart.remove(item));
                    }
                  },
                ),
                Text('${item.quantity}', style: const TextStyle(fontWeight: FontWeight.bold)),
                IconButton(
                  icon: const Icon(Icons.add, size: 18),
                  onPressed: () {
                    setState(() {
                      final index = widget.cart.indexOf(item);
                      widget.cart[index] = CartItem(
                        productId: item.productId,
                        quantity: item.quantity + 1,
                        product: item.product,
                      );
                    });
                  },
                ),
              ],
            ),
          ),
        ],
      ),
    );
  }
}
