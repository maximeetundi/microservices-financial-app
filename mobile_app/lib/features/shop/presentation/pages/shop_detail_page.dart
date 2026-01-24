import 'package:flutter/material.dart';
import '../../data/models/shop_models.dart';
import '../../data/repositories/shop_repository.dart';
import 'cart_page.dart';

class ShopDetailPage extends StatefulWidget {
  final String shopSlug;

  const ShopDetailPage({super.key, required this.shopSlug});

  @override
  State<ShopDetailPage> createState() => _ShopDetailPageState();
}

class _ShopDetailPageState extends State<ShopDetailPage> {
  final ShopRepository _repository = ShopRepository(
    baseUrl: const String.fromEnvironment('API_URL', defaultValue: 'https://api.tech-afm.com'),
  );

  Shop? _shop;
  List<Product> _products = [];
  List<Category> _categories = [];
  bool _loading = true;
  String? _selectedCategory;

  // Cart
  final List<CartItem> _cart = [];

  @override
  void initState() {
    super.initState();
    _loadShop();
  }

  Future<void> _loadShop() async {
    setState(() => _loading = true);
    try {
      final shop = await _repository.getShop(widget.shopSlug);
      final products = await _repository.getProducts(widget.shopSlug);
      final categories = await _repository.getCategories(widget.shopSlug);
      setState(() {
        _shop = shop;
        _products = products;
        _categories = categories;
      });
    } catch (e) {
      debugPrint('Failed to load shop: $e');
    } finally {
      setState(() => _loading = false);
    }
  }

  void _addToCart(Product product) {
    final existingIndex = _cart.indexWhere((item) => item.productId == product.id);
    if (existingIndex >= 0) {
      setState(() {
        _cart[existingIndex] = CartItem(
          productId: product.id,
          quantity: _cart[existingIndex].quantity + 1,
          product: product,
        );
      });
    } else {
      setState(() {
        _cart.add(CartItem(productId: product.id, quantity: 1, product: product));
      });
    }
    ScaffoldMessenger.of(context).showSnackBar(
      SnackBar(content: Text('${product.name} ajouté au panier'), duration: const Duration(seconds: 1)),
    );
  }

  int get _cartItemCount => _cart.fold(0, (sum, item) => sum + item.quantity);
  double get _cartTotal => _cart.fold(0, (sum, item) => sum + (item.product.price * item.quantity));

  String _formatPrice(double amount, String currency) {
    return '${amount.toStringAsFixed(0)} $currency';
  }

  @override
  Widget build(BuildContext context) {
    if (_loading) {
      return const Scaffold(body: Center(child: CircularProgressIndicator()));
    }

    if (_shop == null) {
      return Scaffold(
        appBar: AppBar(title: const Text('Erreur')),
        body: const Center(child: Text('Boutique non trouvée')),
      );
    }

    return Scaffold(
      backgroundColor: const Color(0xFFF8FAFC),
      body: CustomScrollView(
        slivers: [
          // App Bar with Banner
          SliverAppBar(
            expandedHeight: 180,
            pinned: true,
            backgroundColor: Colors.indigo,
            flexibleSpace: FlexibleSpaceBar(
              background: Stack(
                fit: StackFit.expand,
                children: [
                  _shop!.bannerUrl != null
                      ? Image.network(_shop!.bannerUrl!, fit: BoxFit.cover)
                      : Container(
                          decoration: BoxDecoration(
                            gradient: LinearGradient(
                              colors: [Colors.indigo.shade500, Colors.purple.shade500],
                            ),
                          ),
                        ),
                  Container(
                    decoration: BoxDecoration(
                      gradient: LinearGradient(
                        begin: Alignment.topCenter,
                        end: Alignment.bottomCenter,
                        colors: [Colors.transparent, Colors.black.withOpacity(0.6)],
                      ),
                    ),
                  ),
                  Positioned(
                    bottom: 16,
                    left: 16,
                    right: 16,
                    child: Row(
                      children: [
                        Container(
                          width: 60,
                          height: 60,
                          decoration: BoxDecoration(
                            color: Colors.white,
                            borderRadius: BorderRadius.circular(12),
                            border: Border.all(color: Colors.white, width: 3),
                          ),
                          child: ClipRRect(
                            borderRadius: BorderRadius.circular(9),
                            child: _shop!.logoUrl != null
                                ? Image.network(_shop!.logoUrl!, fit: BoxFit.cover)
                                : Container(
                                    color: Colors.indigo,
                                    child: Center(
                                      child: Text(
                                        _shop!.name[0],
                                        style: const TextStyle(color: Colors.white, fontSize: 28, fontWeight: FontWeight.bold),
                                      ),
                                    ),
                                  ),
                          ),
                        ),
                        const SizedBox(width: 12),
                        Expanded(
                          child: Column(
                            crossAxisAlignment: CrossAxisAlignment.start,
                            mainAxisSize: MainAxisSize.min,
                            children: [
                              Text(_shop!.name, style: const TextStyle(color: Colors.white, fontSize: 20, fontWeight: FontWeight.bold)),
                              if (_shop!.description.isNotEmpty)
                                Text(
                                  _shop!.description,
                                  maxLines: 1,
                                  overflow: TextOverflow.ellipsis,
                                  style: const TextStyle(color: Colors.white70, fontSize: 13),
                                ),
                            ],
                          ),
                        ),
                      ],
                    ),
                  ),
                ],
              ),
            ),
          ),

          // Categories
          if (_categories.isNotEmpty)
            SliverToBoxAdapter(
              child: Container(
                height: 48,
                margin: const EdgeInsets.only(top: 8),
                child: ListView.builder(
                  scrollDirection: Axis.horizontal,
                  padding: const EdgeInsets.symmetric(horizontal: 16),
                  itemCount: _categories.length + 1,
                  itemBuilder: (context, index) {
                    if (index == 0) {
                      return _buildCategoryChip('Tous', null);
                    }
                    final cat = _categories[index - 1];
                    return _buildCategoryChip('${cat.name} (${cat.productCount})', cat.slug);
                  },
                ),
              ),
            ),

          // Products Grid
          SliverPadding(
            padding: const EdgeInsets.all(16),
            sliver: _products.isEmpty
                ? const SliverToBoxAdapter(
                    child: Center(child: Text('Aucun produit', style: TextStyle(color: Colors.grey))),
                  )
                : SliverGrid(
                    gridDelegate: const SliverGridDelegateWithFixedCrossAxisCount(
                      crossAxisCount: 2,
                      childAspectRatio: 0.7,
                      crossAxisSpacing: 12,
                      mainAxisSpacing: 12,
                    ),
                    delegate: SliverChildBuilderDelegate(
                      (context, index) => _buildProductCard(_products[index]),
                      childCount: _products.length,
                    ),
                  ),
          ),
        ],
      ),
      floatingActionButton: _cartItemCount > 0
          ? FloatingActionButton.extended(
              onPressed: () {
                Navigator.push(
                  context,
                  MaterialPageRoute(
                    builder: (_) => CartPage(
                      shop: _shop!,
                      cart: _cart,
                      onClearCart: () => setState(() => _cart.clear()),
                    ),
                  ),
                );
              },
              backgroundColor: Colors.indigo,
              icon: const Icon(Icons.shopping_cart),
              label: Text('$_cartItemCount • ${_formatPrice(_cartTotal, _shop!.currency)}'),
            )
          : null,
    );
  }

  Widget _buildCategoryChip(String label, String? slug) {
    final isSelected = _selectedCategory == slug;
    return Padding(
      padding: const EdgeInsets.only(right: 8),
      child: FilterChip(
        label: Text(label),
        selected: isSelected,
        onSelected: (selected) async {
          setState(() => _selectedCategory = selected ? slug : null);
          final products = await _repository.getProducts(widget.shopSlug, category: _selectedCategory);
          setState(() => _products = products);
        },
        selectedColor: Colors.indigo.shade100,
        checkmarkColor: Colors.indigo,
        backgroundColor: Colors.white,
        labelStyle: TextStyle(
          color: isSelected ? Colors.indigo : Colors.grey[700],
          fontWeight: isSelected ? FontWeight.bold : FontWeight.normal,
        ),
      ),
    );
  }

  Widget _buildProductCard(Product product) {
    return GestureDetector(
      onTap: () => _showProductModal(product),
      child: Container(
        decoration: BoxDecoration(
          color: Colors.white,
          borderRadius: BorderRadius.circular(12),
          boxShadow: [BoxShadow(color: Colors.black.withOpacity(0.05), blurRadius: 8, offset: const Offset(0, 2))],
        ),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            // Image
            Expanded(
              child: ClipRRect(
                borderRadius: const BorderRadius.vertical(top: Radius.circular(12)),
                child: Stack(
                  children: [
                    Container(
                      width: double.infinity,
                      color: Colors.grey[100],
                      child: product.images.isNotEmpty
                          ? Image.network(product.images[0], fit: BoxFit.cover)
                          : const Center(child: Text('📦', style: TextStyle(fontSize: 40))),
                    ),
                    if (product.isFeatured)
                      Positioned(
                        top: 8,
                        left: 8,
                        child: Container(
                          padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 4),
                          decoration: BoxDecoration(color: Colors.amber, borderRadius: BorderRadius.circular(8)),
                          child: const Text('⭐ Featured', style: TextStyle(color: Colors.white, fontSize: 10, fontWeight: FontWeight.bold)),
                        ),
                      ),
                    if (product.stock == 0)
                      Positioned(
                        top: 8,
                        right: 8,
                        child: Container(
                          padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 4),
                          decoration: BoxDecoration(color: Colors.red, borderRadius: BorderRadius.circular(8)),
                          child: const Text('Épuisé', style: TextStyle(color: Colors.white, fontSize: 10, fontWeight: FontWeight.bold)),
                        ),
                      ),
                  ],
                ),
              ),
            ),

            // Info
            Padding(
              padding: const EdgeInsets.all(10),
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  Text(
                    product.name,
                    maxLines: 1,
                    overflow: TextOverflow.ellipsis,
                    style: const TextStyle(fontWeight: FontWeight.w600, fontSize: 13),
                  ),
                  const SizedBox(height: 4),
                  Row(
                    children: [
                      Text(
                        _formatPrice(product.price, product.currency),
                        style: TextStyle(fontWeight: FontWeight.bold, color: Colors.indigo.shade600),
                      ),
                      if (product.compareAtPrice != null) ...[
                        const SizedBox(width: 6),
                        Text(
                          _formatPrice(product.compareAtPrice!, product.currency),
                          style: TextStyle(fontSize: 11, color: Colors.grey[400], decoration: TextDecoration.lineThrough),
                        ),
                      ],
                    ],
                  ),
                ],
              ),
            ),
          ],
        ),
      ),
    );
  }

  void _showProductModal(Product product) {
    showModalBottomSheet(
      context: context,
      isScrollControlled: true,
      backgroundColor: Colors.transparent,
      builder: (context) => Container(
        height: MediaQuery.of(context).size.height * 0.7,
        decoration: const BoxDecoration(
          color: Colors.white,
          borderRadius: BorderRadius.vertical(top: Radius.circular(20)),
        ),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.stretch,
          children: [
            // Image
            Expanded(
              flex: 2,
              child: ClipRRect(
                borderRadius: const BorderRadius.vertical(top: Radius.circular(20)),
                child: product.images.isNotEmpty
                    ? Image.network(product.images[0], fit: BoxFit.cover)
                    : Container(color: Colors.grey[100], child: const Center(child: Text('📦', style: TextStyle(fontSize: 60)))),
              ),
            ),

            // Content
            Expanded(
              flex: 2,
              child: Padding(
                padding: const EdgeInsets.all(20),
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    Text(product.name, style: const TextStyle(fontSize: 22, fontWeight: FontWeight.bold)),
                    const SizedBox(height: 8),
                    if (product.description.isNotEmpty)
                      Expanded(
                        child: Text(product.description, style: TextStyle(color: Colors.grey[600])),
                      ),
                    Row(
                      children: [
                        Text(
                          _formatPrice(product.price, product.currency),
                          style: TextStyle(fontSize: 28, fontWeight: FontWeight.bold, color: Colors.indigo.shade600),
                        ),
                        if (product.compareAtPrice != null) ...[
                          const SizedBox(width: 12),
                          Text(
                            _formatPrice(product.compareAtPrice!, product.currency),
                            style: TextStyle(fontSize: 18, color: Colors.grey[400], decoration: TextDecoration.lineThrough),
                          ),
                        ],
                      ],
                    ),
                    const SizedBox(height: 16),
                    SizedBox(
                      width: double.infinity,
                      child: ElevatedButton(
                        onPressed: product.stock == 0 ? null : () {
                          _addToCart(product);
                          Navigator.pop(context);
                        },
                        style: ElevatedButton.styleFrom(
                          backgroundColor: Colors.indigo,
                          padding: const EdgeInsets.symmetric(vertical: 16),
                          shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(12)),
                        ),
                        child: Text(
                          product.stock == 0 ? 'Épuisé' : 'Ajouter au panier',
                          style: const TextStyle(fontSize: 16, fontWeight: FontWeight.bold),
                        ),
                      ),
                    ),
                  ],
                ),
              ),
            ),
          ],
        ),
      ),
    );
  }
}
