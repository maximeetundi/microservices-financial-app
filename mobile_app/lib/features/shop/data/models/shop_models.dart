class Shop {
  final String id;
  final String ownerId;
  final String ownerType;
  final String name;
  final String slug;
  final String description;
  final String? logoUrl;
  final String? bannerUrl;
  final bool isPublic;
  final String walletId;
  final String currency;
  final List<String> tags;
  final String? qrCode;
  final String status;
  final ShopStats stats;
  final DateTime createdAt;

  Shop({
    required this.id,
    required this.ownerId,
    required this.ownerType,
    required this.name,
    required this.slug,
    required this.description,
    this.logoUrl,
    this.bannerUrl,
    required this.isPublic,
    required this.walletId,
    required this.currency,
    required this.tags,
    this.qrCode,
    required this.status,
    required this.stats,
    required this.createdAt,
  });

  factory Shop.fromJson(Map<String, dynamic> json) {
    return Shop(
      id: json['id'] ?? '',
      ownerId: json['owner_id'] ?? '',
      ownerType: json['owner_type'] ?? 'user',
      name: json['name'] ?? '',
      slug: json['slug'] ?? '',
      description: json['description'] ?? '',
      logoUrl: json['logo_url'],
      bannerUrl: json['banner_url'],
      isPublic: json['is_public'] ?? true,
      walletId: json['wallet_id'] ?? '',
      currency: json['currency'] ?? 'XOF',
      tags: List<String>.from(json['tags'] ?? []),
      qrCode: json['qr_code'],
      status: json['status'] ?? 'active',
      stats: ShopStats.fromJson(json['stats'] ?? {}),
      createdAt: DateTime.tryParse(json['created_at'] ?? '') ?? DateTime.now(),
    );
  }
}

class ShopStats {
  final int totalProducts;
  final int totalOrders;
  final double totalRevenue;
  final double averageRating;

  ShopStats({
    required this.totalProducts,
    required this.totalOrders,
    required this.totalRevenue,
    required this.averageRating,
  });

  factory ShopStats.fromJson(Map<String, dynamic> json) {
    return ShopStats(
      totalProducts: json['total_products'] ?? 0,
      totalOrders: json['total_orders'] ?? 0,
      totalRevenue: (json['total_revenue'] ?? 0).toDouble(),
      averageRating: (json['average_rating'] ?? 0).toDouble(),
    );
  }
}

class Product {
  final String id;
  final String shopId;
  final String? categoryId;
  final String name;
  final String slug;
  final String description;
  final String? shortDesc;
  final double price;
  final double? compareAtPrice;
  final String currency;
  final List<String> images;
  final int stock;
  final bool isCustomizable;
  final List<CustomField> customFields;
  final List<String> tags;
  final String? qrCode;
  final String status;
  final bool isFeatured;

  Product({
    required this.id,
    required this.shopId,
    this.categoryId,
    required this.name,
    required this.slug,
    required this.description,
    this.shortDesc,
    required this.price,
    this.compareAtPrice,
    required this.currency,
    required this.images,
    required this.stock,
    required this.isCustomizable,
    required this.customFields,
    required this.tags,
    this.qrCode,
    required this.status,
    required this.isFeatured,
  });

  factory Product.fromJson(Map<String, dynamic> json) {
    return Product(
      id: json['id'] ?? '',
      shopId: json['shop_id'] ?? '',
      categoryId: json['category_id'],
      name: json['name'] ?? '',
      slug: json['slug'] ?? '',
      description: json['description'] ?? '',
      shortDesc: json['short_desc'],
      price: (json['price'] ?? 0).toDouble(),
      compareAtPrice: json['compare_at_price']?.toDouble(),
      currency: json['currency'] ?? 'XOF',
      images: List<String>.from(json['images'] ?? []),
      stock: json['stock'] ?? 0,
      isCustomizable: json['is_customizable'] ?? false,
      customFields: (json['custom_fields'] as List<dynamic>?)
              ?.map((e) => CustomField.fromJson(e))
              .toList() ??
          [],
      tags: List<String>.from(json['tags'] ?? []),
      qrCode: json['qr_code'],
      status: json['status'] ?? 'active',
      isFeatured: json['is_featured'] ?? false,
    );
  }
}

class CustomField {
  final String name;
  final String type;
  final bool required;
  final List<String>? options;

  CustomField({
    required this.name,
    required this.type,
    required this.required,
    this.options,
  });

  factory CustomField.fromJson(Map<String, dynamic> json) {
    return CustomField(
      name: json['name'] ?? '',
      type: json['type'] ?? 'text',
      required: json['required'] ?? false,
      options: json['options'] != null ? List<String>.from(json['options']) : null,
    );
  }
}

class Category {
  final String id;
  final String shopId;
  final String? parentId;
  final String name;
  final String slug;
  final String? imageUrl;
  final int productCount;

  Category({
    required this.id,
    required this.shopId,
    this.parentId,
    required this.name,
    required this.slug,
    this.imageUrl,
    required this.productCount,
  });

  factory Category.fromJson(Map<String, dynamic> json) {
    return Category(
      id: json['id'] ?? '',
      shopId: json['shop_id'] ?? '',
      parentId: json['parent_id'],
      name: json['name'] ?? '',
      slug: json['slug'] ?? '',
      imageUrl: json['image_url'],
      productCount: json['product_count'] ?? 0,
    );
  }
}

class Order {
  final String id;
  final String orderNumber;
  final String shopId;
  final String shopName;
  final List<OrderItem> items;
  final double totalAmount;
  final String currency;
  final String paymentStatus;
  final String orderStatus;
  final DateTime createdAt;

  Order({
    required this.id,
    required this.orderNumber,
    required this.shopId,
    required this.shopName,
    required this.items,
    required this.totalAmount,
    required this.currency,
    required this.paymentStatus,
    required this.orderStatus,
    required this.createdAt,
  });

  factory Order.fromJson(Map<String, dynamic> json) {
    return Order(
      id: json['id'] ?? '',
      orderNumber: json['order_number'] ?? '',
      shopId: json['shop_id'] ?? '',
      shopName: json['shop_name'] ?? '',
      items: (json['items'] as List<dynamic>?)
              ?.map((e) => OrderItem.fromJson(e))
              .toList() ??
          [],
      totalAmount: (json['total_amount'] ?? 0).toDouble(),
      currency: json['currency'] ?? 'XOF',
      paymentStatus: json['payment_status'] ?? 'pending',
      orderStatus: json['order_status'] ?? 'pending',
      createdAt: DateTime.tryParse(json['created_at'] ?? '') ?? DateTime.now(),
    );
  }
}

class OrderItem {
  final String productId;
  final String productName;
  final String? productImage;
  final int quantity;
  final double unitPrice;
  final double totalPrice;

  OrderItem({
    required this.productId,
    required this.productName,
    this.productImage,
    required this.quantity,
    required this.unitPrice,
    required this.totalPrice,
  });

  factory OrderItem.fromJson(Map<String, dynamic> json) {
    return OrderItem(
      productId: json['product_id'] ?? '',
      productName: json['product_name'] ?? '',
      productImage: json['product_image'],
      quantity: json['quantity'] ?? 1,
      unitPrice: (json['unit_price'] ?? 0).toDouble(),
      totalPrice: (json['total_price'] ?? 0).toDouble(),
    );
  }
}

class CartItem {
  final String productId;
  final int quantity;
  final Map<String, String>? customValues;
  final Product product;

  CartItem({
    required this.productId,
    required this.quantity,
    this.customValues,
    required this.product,
  });

  Map<String, dynamic> toJson() {
    return {
      'product_id': productId,
      'quantity': quantity,
      'custom_values': customValues,
    };
  }
}
