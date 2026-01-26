import 'dart:io';
import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';
import 'package:google_fonts/google_fonts.dart';
import 'package:image_picker/image_picker.dart';
import '../../../../core/theme/app_theme.dart';
import '../../../../core/services/api_service.dart';
import '../../../../core/widgets/glass_container.dart';
import '../../../../core/widgets/custom_text_field.dart';
import '../../../../core/widgets/custom_button.dart';

class AddProductPage extends StatefulWidget {
  const AddProductPage({super.key});

  @override
  State<AddProductPage> createState() => _AddProductPageState();
}

class _AddProductPageState extends State<AddProductPage> {
  final _formKey = GlobalKey<FormState>();
  final ApiService _api = ApiService();
  
  final _nameController = TextEditingController();
  final _descController = TextEditingController();
  final _priceController = TextEditingController();
  final _qtyController = TextEditingController(text: '1');
  
  String _currency = 'USD';
  File? _image;
  bool _isLoading = false;

  Future<void> _pickImage() async {
    final picker = ImagePicker();
    final pickedFile = await picker.pickImage(source: ImageSource.gallery);
    if (pickedFile != null) {
      setState(() => _image = File(pickedFile.path));
    }
  }

  Future<void> _submit() async {
    if (!_formKey.currentState!.validate()) return;
    
    setState(() => _isLoading = true);
    
    try {
      final data = {
        'name': _nameController.text,
        'description': _descController.text,
        'price': double.parse(_priceController.text),
        'currency': _currency,
        'quantity': int.parse(_qtyController.text),
        if (_image != null) 'image_path': _image!.path,
      };
      
      await _api.shop.createProduct(data);
      
      if (mounted) {
        ScaffoldMessenger.of(context).showSnackBar(
          const SnackBar(content: Text('Produit créé avec succès !')),
        );
        context.pop(); // Return to MyProductsPage
        // ideally trigger refresh there
      }
    } catch (e) {
      if (mounted) {
         ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(content: Text('Erreur: $e'), backgroundColor: Colors.red),
        );
      }
    } finally {
      setState(() => _isLoading = false);
    }
  }

  @override
  Widget build(BuildContext context) {
    final isDark = Theme.of(context).brightness == Brightness.dark;

    return Scaffold(
      backgroundColor: Colors.transparent,
      body: Container(
        decoration: BoxDecoration(
          gradient: LinearGradient(
            begin: Alignment.topLeft,
            end: Alignment.bottomRight,
            colors: isDark 
                ? [const Color(0xFF020617), const Color(0xFF0F172A)]
                : [const Color(0xFFFAFBFC), const Color(0xFFEFF6FF)],
          ),
        ),
        child: SafeArea(
          child: Column(
            children: [
               // Header
              Padding(
                padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 8),
                child: Row(
                  children: [
                    GlassContainer(
                      padding: EdgeInsets.zero,
                      width: 40,
                      height: 40,
                      borderRadius: 12,
                      child: IconButton(
                        icon: Icon(Icons.arrow_back_ios_new, size: 20, 
                            color: isDark ? Colors.white : AppTheme.textPrimaryColor),
                        onPressed: () => context.pop(),
                      ),
                    ),
                    const SizedBox(width: 16),
                    Text(
                      'Ajouter un produit',
                      style: GoogleFonts.inter(
                        fontSize: 20,
                        fontWeight: FontWeight.bold,
                        color: isDark ? Colors.white : AppTheme.textPrimaryColor,
                      ),
                    ),
                  ],
                ),
              ),
              
              Expanded(
                child: SingleChildScrollView(
                  padding: const EdgeInsets.all(20),
                  child: GlassContainer(
                    padding: const EdgeInsets.all(24),
                    borderRadius: 24,
                    child: Form(
                      key: _formKey,
                      child: Column(
                        crossAxisAlignment: CrossAxisAlignment.start,
                        children: [
                          // Image Picker
                          Center(
                            child: GestureDetector(
                              onTap: _pickImage,
                              child: Container(
                                width: 120,
                                height: 120,
                                decoration: BoxDecoration(
                                  color: isDark ? Colors.white10 : Colors.grey.shade100,
                                  borderRadius: BorderRadius.circular(20),
                                  border: Border.all(color: Colors.grey.shade300),
                                  image: _image != null 
                                      ? DecorationImage(image: FileImage(_image!), fit: BoxFit.cover)
                                      : null,
                                ),
                                child: _image == null 
                                    ? const Icon(Icons.add_a_photo, size: 40, color: Colors.grey)
                                    : null,
                              ),
                            ),
                          ),
                          const SizedBox(height: 24),
                          
                          CustomTextField(
                            controller: _nameController,
                            label: 'Nom du produit',
                            hint: 'Ex: T-Shirt Vintage',
                            validator: (v) => v!.isEmpty ? 'Requis' : null,
                          ),
                          const SizedBox(height: 16),
                          
                          Row(
                            children: [
                              Expanded(
                                flex: 2,
                                child: CustomTextField(
                                  controller: _priceController,
                                  label: 'Prix',
                                  hint: '0.00',
                                  keyboardType: const TextInputType.numberWithOptions(decimal: true),
                                  validator: (v) => v!.isEmpty ? 'Requis' : null,
                                ),
                              ),
                              const SizedBox(width: 16),
                              Expanded(
                                child: DropdownButtonFormField<String>(
                                  value: _currency,
                                  decoration: InputDecoration(
                                    labelText: 'Devise',
                                     labelStyle: TextStyle(
                                      color: isDark ? Colors.white70 : Colors.black54,
                                    ),
                                    border: OutlineInputBorder(borderRadius: BorderRadius.circular(12)),
                                  ),
                                  dropdownColor: isDark ? const Color(0xFF1E293B) : Colors.white,
                                  items: ['USD', 'EUR', 'XOF'].map((c) => DropdownMenuItem(
                                    value: c, 
                                    child: Text(c, style: TextStyle(color: isDark ? Colors.white : Colors.black))
                                  )).toList(),
                                  onChanged: (v) => setState(() => _currency = v!),
                                ),
                              ),
                            ],
                          ),
                          const SizedBox(height: 16),
                          
                          CustomTextField(
                            controller: _qtyController,
                            label: 'Quantité en stock',
                            keyboardType: TextInputType.number,
                          ),
                          const SizedBox(height: 16),
                          
                          CustomTextField(
                            controller: _descController,
                            label: 'Description',
                            maxLines: 3,
                          ),
                          const SizedBox(height: 32),
                          
                          SizedBox(
                            width: double.infinity,
                            child: CustomButton(
                              text: 'Créer le produit',
                              onPressed: _isLoading ? null : _submit,
                              isLoading: _isLoading,
                              backgroundColor: AppTheme.primaryColor,
                              textColor: Colors.white,
                            ),
                          ),
                        ],
                      ),
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
}
