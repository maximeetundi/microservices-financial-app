import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';
import 'package:go_router/go_router.dart';
import 'package:google_fonts/google_fonts.dart';

import '../../../../core/theme/app_theme.dart';
import '../../../../core/widgets/glass_container.dart';
import '../../../../core/widgets/custom_text_field.dart';
import '../../../../core/widgets/custom_button.dart';
import '../../../../core/widgets/security_confirmation.dart';
import '../../../wallet/presentation/bloc/wallet_bloc.dart';
import '../bloc/transfer_bloc.dart';

/// Modern Transfer Page matching frontend design
class TransferPage extends StatefulWidget {
  const TransferPage({super.key});

  @override
  State<TransferPage> createState() => _TransferPageState();
}

class _TransferPageState extends State<TransferPage> {
  String _selectedType = 'p2p';
  String? _selectedWalletId;
  
  final _amountController = TextEditingController();
  final _descriptionController = TextEditingController();
  
  // P2P fields
  final _recipientController = TextEditingController();
  bool _isLookingUp = false;
  Map<String, dynamic>? _lookupResult;
  String? _lookupError;
  
  // Mobile Money fields
  String _selectedCountry = 'CI';
  String _selectedProvider = 'orange';
  final _phoneController = TextEditingController();
  
  // Wire fields  
  final _bankNameController = TextEditingController();
  final _ibanController = TextEditingController();
  final _recipientNameController = TextEditingController();

  bool _isLoading = false;

  final List<Map<String, dynamic>> _transferTypes = [
    {'id': 'p2p', 'name': 'Interne (P2P)', 'icon': '👤'},
    {'id': 'mobile', 'name': 'Mobile Money', 'icon': '📱'},
    {'id': 'wire', 'name': 'Virement', 'icon': '🏦'},
    {'id': 'crypto', 'name': 'Crypto', 'icon': '₿'},
  ];

  @override
  void initState() {
    super.initState();
    context.read<WalletBloc>().add(LoadWalletsEvent());
  }

  @override
  void dispose() {
    _amountController.dispose();
    _descriptionController.dispose();
    _recipientController.dispose();
    _phoneController.dispose();
    _bankNameController.dispose();
    _ibanController.dispose();
    _recipientNameController.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    final isDark = Theme.of(context).brightness == Brightness.dark;

    return Scaffold(
      backgroundColor: Colors.transparent, // Allow gradient to show
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
               // Custom App Bar Area
               Padding(
                 padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 8),
                 child: Row(
                   mainAxisAlignment: MainAxisAlignment.spaceBetween,
                   children: [
                     GlassContainer(
                       padding: EdgeInsets.zero,
                       width: 40, 
                       height: 40,
                       borderRadius: 12,
                       child: IconButton(
                        icon: Icon(Icons.arrow_back_ios_new, size: 20, color: isDark ? Colors.white : AppTheme.textPrimaryColor),
                        onPressed: () => context.go('/dashboard'),
                      ),
                     ),
                     Text(
                        'Envoyer de l\'argent 💸',
                         style: GoogleFonts.inter(
                           fontSize: 20,
                           fontWeight: FontWeight.bold,
                           color: isDark ? Colors.white : AppTheme.textPrimaryColor,
                         ),
                      ),
                      GlassContainer(
                        padding: EdgeInsets.zero,
                        width: 40, 
                        height: 40,
                        borderRadius: 12,
                        child: IconButton(
                         icon: Icon(Icons.home_rounded, size: 20, color: isDark ? Colors.white : AppTheme.textPrimaryColor),
                         onPressed: () => context.go('/dashboard'),
                       ),
                      ),
                   ],
                 ),
               ),
               Expanded(
                 child: SingleChildScrollView(
                  padding: const EdgeInsets.all(20),
                  child: Column(
                    crossAxisAlignment: CrossAxisAlignment.start,
                    children: [
                      // Subtitle
                      Text(
                        'Transferts P2P, Mobile Money, et virements bancaires',
                        style: GoogleFonts.inter(
                          color: isDark ? Colors.white70 : AppTheme.textSecondaryColor,
                          fontSize: 14
                        ),
                      ),
                      const SizedBox(height: 24),
                      
                      // Transfer Type Selector
                      _buildTypeSelector(),
                      const SizedBox(height: 24),
                      
                      // Transfer Form Card
                      GlassContainer(
                        padding: const EdgeInsets.all(24),
                        borderRadius: 24,
                        child: Column(
                          crossAxisAlignment: CrossAxisAlignment.start,
                          children: [
                            // From Wallet
                            _buildWalletSelector(),
                            const SizedBox(height: 24),
                            
                            // Amount
                            _buildAmountField(),
                            const SizedBox(height: 24),
                            
                            // Type-specific fields
                            _buildTypeSpecificFields(),
                            const SizedBox(height: 24),
                            
                            // Description
                            _buildDescriptionField(),
                            const SizedBox(height: 24),
                            
                            // Summary
                            _buildSummary(),
                            const SizedBox(height: 24),
                            
                            // Submit Button
                            _buildSubmitButton(),
                          ],
                        ),
                      ),
                    ],
                  ),
                ),
               ),
            ],
          ),
        ),
      ),
    );
  }

  Widget _buildTypeSelector() {
    final isDark = Theme.of(context).brightness == Brightness.dark;
    
    return GridView.count(
      crossAxisCount: 2,
      shrinkWrap: true,
      physics: const NeverScrollableScrollPhysics(),
      crossAxisSpacing: 12,
      mainAxisSpacing: 12,
      childAspectRatio: 2.2,
      children: _transferTypes.map((type) {
        final isSelected = _selectedType == type['id'];
        return GestureDetector(
          onTap: () => setState(() {
            _selectedType = type['id'];
            _lookupResult = null;
            _lookupError = null;
          }),
          child: AnimatedContainer(
            duration: const Duration(milliseconds: 200),
            padding: const EdgeInsets.all(16),
            decoration: BoxDecoration(
              color: isSelected 
                  ? const Color(0xFF6366F1).withOpacity(isDark ? 0.3 : 0.1)
                  : isDark ? const Color(0xFF1E293B) : Colors.white,
              borderRadius: BorderRadius.circular(16),
              border: Border.all(
                color: isSelected 
                    ? const Color(0xFF6366F1)
                    : isDark ? const Color(0xFF334155) : const Color(0xFFE2E8F0),
                width: isSelected ? 2 : 1,
              ),
              boxShadow: isSelected ? [
                BoxShadow(
                  color: const Color(0xFF6366F1).withOpacity(0.2),
                  blurRadius: 12,
                  offset: const Offset(0, 4),
                ),
              ] : null,
            ),
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              mainAxisAlignment: MainAxisAlignment.center,
              children: [
                Text(type['icon'], style: const TextStyle(fontSize: 24)),
                const SizedBox(height: 4),
                Text(
                  type['name'],
                  style: TextStyle(
                    fontSize: 12,
                    fontWeight: FontWeight.w600,
                    color: isSelected 
                        ? const Color(0xFF818CF8)
                        : isDark ? const Color(0xFF94A3B8) : const Color(0xFF64748B),
                  ),
                ),
              ],
            ),
          ),
        );
      }).toList(),
    );
  }

  Widget _buildWalletSelector() {
    final isDark = Theme.of(context).brightness == Brightness.dark;
    
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Text(
          'Depuis le portefeuille',
          style: TextStyle(
            fontSize: 14,
            fontWeight: FontWeight.w500,
            color: isDark ? const Color(0xFF94A3B8) : const Color(0xFF64748B),
          ),
        ),
        const SizedBox(height: 8),
        BlocBuilder<WalletBloc, WalletState>(
          builder: (context, state) {
            if (state is WalletLoadedState) {
              final wallets = state.wallets;
              if (_selectedWalletId == null && wallets.isNotEmpty) {
                _selectedWalletId = wallets.first.id;
              }
              
              return Container(
                padding: const EdgeInsets.symmetric(horizontal: 16),
                decoration: BoxDecoration(
                  color: isDark ? const Color(0xFF1E293B) : const Color(0xFFF8FAFC),
                  borderRadius: BorderRadius.circular(12),
                  border: Border.all(color: isDark ? const Color(0xFF334155) : const Color(0xFFE2E8F0)),
                ),
                child: DropdownButton<String>(
                  value: _selectedWalletId,
                  isExpanded: true,
                  underline: const SizedBox(),
                  dropdownColor: isDark ? const Color(0xFF1E293B) : Colors.white,
                  icon: Icon(Icons.keyboard_arrow_down, color: isDark ? const Color(0xFF64748B) : const Color(0xFF64748B)),
                  items: wallets.map((wallet) {
                    return DropdownMenuItem<String>(
                      value: wallet.id,
                      child: Text(
                        '${wallet.name ?? wallet.currency} - ${wallet.balance.toStringAsFixed(2)} ${wallet.currency}',
                        style: TextStyle(
                          fontSize: 15,
                          color: isDark ? Colors.white : const Color(0xFF1E293B),
                        ),
                      ),
                    );
                  }).toList(),
                  onChanged: (value) => setState(() => _selectedWalletId = value),
                ),
              );
            }
            return const Center(child: CircularProgressIndicator());
          },
        ),
      ],
    );
  }

  Widget _buildAmountField() {
    final isDark = Theme.of(context).brightness == Brightness.dark;
    
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Text(
          'Montant à envoyer',
          style: TextStyle(
            fontSize: 14,
            fontWeight: FontWeight.w500,
            color: isDark ? const Color(0xFF94A3B8) : const Color(0xFF64748B),
          ),
        ),
        const SizedBox(height: 8),
        Container(
          padding: const EdgeInsets.all(16),
          decoration: BoxDecoration(
            color: isDark ? const Color(0xFF1E293B) : const Color(0xFFF8FAFC),
            borderRadius: BorderRadius.circular(12),
            border: Border.all(color: isDark ? const Color(0xFF334155) : const Color(0xFFE2E8F0)),
          ),
          child: Row(
            children: [
              Expanded(
                child: TextField(
                  controller: _amountController,
                  keyboardType: const TextInputType.numberWithOptions(decimal: true),
                  style: TextStyle(
                    fontSize: 28,
                    fontWeight: FontWeight.bold,
                    color: isDark ? Colors.white : const Color(0xFF1a1a2e),
                  ),
                  decoration: InputDecoration(
                    border: InputBorder.none,
                    hintText: '0.00',
                    hintStyle: TextStyle(color: isDark ? const Color(0xFF475569) : const Color(0xFFCBD5E1)),
                  ),
                  onChanged: (_) => setState(() {}),
                ),
              ),
              Container(
                padding: const EdgeInsets.symmetric(horizontal: 12, vertical: 6),
                decoration: BoxDecoration(
                  color: const Color(0xFF6366F1).withOpacity(isDark ? 0.3 : 0.1),
                  borderRadius: BorderRadius.circular(8),
                ),
                child: Text(
                  _getSelectedCurrency(),
                  style: const TextStyle(
                    color: Color(0xFF818CF8),
                    fontWeight: FontWeight.bold,
                  ),
                ),
              ),
            ],
          ),
        ),
        const SizedBox(height: 8),
        BlocBuilder<WalletBloc, WalletState>(
          builder: (context, state) {
            if (state is WalletLoadedState && _selectedWalletId != null) {
              final wallet = state.wallets.firstWhere(
                (w) => w.id == _selectedWalletId,
                orElse: () => state.wallets.first,
              );
              final amount = double.tryParse(_amountController.text) ?? 0;
              final isInsufficient = amount > wallet.balance;
              
              return Text(
                isInsufficient 
                    ? 'Solde insuffisant'
                    : 'Disponible: ${wallet.balance.toStringAsFixed(2)} ${wallet.currency}',
                style: TextStyle(
                  fontSize: 12,
                  color: isInsufficient ? Colors.red : (isDark ? const Color(0xFF64748B) : const Color(0xFF64748B)),
                  fontWeight: isInsufficient ? FontWeight.w600 : FontWeight.normal,
                ),
              );
            }
            return const SizedBox();
          },
        ),
      ],
    );
  }

  Widget _buildTypeSpecificFields() {
    switch (_selectedType) {
      case 'p2p':
        return _buildP2PFields();
      case 'mobile':
        return _buildMobileMoneyFields();
      case 'wire':
        return _buildWireFields();
      case 'crypto':
        return _buildCryptoFields();
      default:
        return const SizedBox();
    }
  }

  Widget _buildP2PFields() {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        const Text(
          'Destinataire (Email ou Téléphone)',
          style: TextStyle(
            fontSize: 14,
            fontWeight: FontWeight.w500,
            color: Color(0xFF64748B),
          ),
        ),
        const SizedBox(height: 8),
        Row(
          children: [
            Expanded(
              child: TextField(
                controller: _recipientController,
                decoration: InputDecoration(
                  hintText: 'ex: ami@email.com ou +225...',
                  filled: true,
                  fillColor: const Color(0xFFF8FAFC),
                  border: OutlineInputBorder(
                    borderRadius: BorderRadius.circular(12),
                    borderSide: const BorderSide(color: Color(0xFFE2E8F0)),
                  ),
                  enabledBorder: OutlineInputBorder(
                    borderRadius: BorderRadius.circular(12),
                    borderSide: const BorderSide(color: Color(0xFFE2E8F0)),
                  ),
                ),
              ),
            ),
            const SizedBox(width: 8),
            GestureDetector(
              onTap: _lookupUser,
              child: Container(
                padding: const EdgeInsets.all(14),
                decoration: BoxDecoration(
                  color: const Color(0xFFF8FAFC),
                  borderRadius: BorderRadius.circular(12),
                  border: Border.all(color: const Color(0xFFE2E8F0)),
                ),
                child: _isLookingUp
                    ? const SizedBox(
                        width: 20,
                        height: 20,
                        child: CircularProgressIndicator(strokeWidth: 2),
                      )
                    : const Text('🔍', style: TextStyle(fontSize: 20)),
              ),
            ),
          ],
        ),
        if (_lookupError != null) ...[
          const SizedBox(height: 8),
          Text(
            _lookupError!,
            style: const TextStyle(color: Colors.red, fontSize: 12),
          ),
        ],
        if (_lookupResult != null) ...[
          const SizedBox(height: 12),
          Container(
            padding: const EdgeInsets.all(12),
            decoration: BoxDecoration(
              color: const Color(0xFF10B981).withOpacity(0.1),
              borderRadius: BorderRadius.circular(12),
              border: Border.all(color: const Color(0xFF10B981).withOpacity(0.2)),
            ),
            child: Row(
              children: [
                Container(
                  width: 40,
                  height: 40,
                  decoration: BoxDecoration(
                    color: const Color(0xFF10B981),
                    borderRadius: BorderRadius.circular(20),
                  ),
                  child: Center(
                    child: Text(
                      (_lookupResult!['first_name'] as String?)?.substring(0, 1).toUpperCase() ?? 'U',
                      style: const TextStyle(
                        color: Colors.white,
                        fontWeight: FontWeight.bold,
                        fontSize: 18,
                      ),
                    ),
                  ),
                ),
                const SizedBox(width: 12),
                Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    Text(
                      '${_lookupResult!['first_name']} ${_lookupResult!['last_name']}',
                      style: const TextStyle(
                        fontWeight: FontWeight.bold,
                        color: Color(0xFF1a1a2e),
                      ),
                    ),
                    const Row(
                      children: [
                        Icon(Icons.check_circle, color: Color(0xFF10B981), size: 14),
                        SizedBox(width: 4),
                        Text(
                          'Utilisateur vérifié',
                          style: TextStyle(
                            color: Color(0xFF10B981),
                            fontSize: 12,
                            fontWeight: FontWeight.w500,
                          ),
                        ),
                      ],
                    ),
                  ],
                ),
              ],
            ),
          ),
        ],
      ],
    );
  }

  Widget _buildMobileMoneyFields() {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Row(
          children: [
            Expanded(
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  const Text('Pays', style: TextStyle(fontSize: 14, fontWeight: FontWeight.w500, color: Color(0xFF64748B))),
                  const SizedBox(height: 8),
                  Container(
                    padding: const EdgeInsets.symmetric(horizontal: 12),
                    decoration: BoxDecoration(
                      color: const Color(0xFFF8FAFC),
                      borderRadius: BorderRadius.circular(12),
                      border: Border.all(color: const Color(0xFFE2E8F0)),
                    ),
                    child: DropdownButton<String>(
                      value: _selectedCountry,
                      isExpanded: true,
                      underline: const SizedBox(),
                      items: const [
                        DropdownMenuItem(value: 'CI', child: Text('🇨🇮 Côte d\'Ivoire')),
                        DropdownMenuItem(value: 'SN', child: Text('🇸🇳 Sénégal')),
                        DropdownMenuItem(value: 'CM', child: Text('🇨🇲 Cameroun')),
                      ],
                      onChanged: (v) => setState(() => _selectedCountry = v!),
                    ),
                  ),
                ],
              ),
            ),
            const SizedBox(width: 12),
            Expanded(
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  const Text('Opérateur', style: TextStyle(fontSize: 14, fontWeight: FontWeight.w500, color: Color(0xFF64748B))),
                  const SizedBox(height: 8),
                  Container(
                    padding: const EdgeInsets.symmetric(horizontal: 12),
                    decoration: BoxDecoration(
                      color: const Color(0xFFF8FAFC),
                      borderRadius: BorderRadius.circular(12),
                      border: Border.all(color: const Color(0xFFE2E8F0)),
                    ),
                    child: DropdownButton<String>(
                      value: _selectedProvider,
                      isExpanded: true,
                      underline: const SizedBox(),
                      items: const [
                        DropdownMenuItem(value: 'orange', child: Text('Orange Money')),
                        DropdownMenuItem(value: 'mtn', child: Text('MTN MoMo')),
                        DropdownMenuItem(value: 'wave', child: Text('Wave')),
                      ],
                      onChanged: (v) => setState(() => _selectedProvider = v!),
                    ),
                  ),
                ],
              ),
            ),
          ],
        ),
        const SizedBox(height: 16),
        const Text('Numéro de téléphone', style: TextStyle(fontSize: 14, fontWeight: FontWeight.w500, color: Color(0xFF64748B))),
        const SizedBox(height: 8),
        TextField(
          controller: _phoneController,
          keyboardType: TextInputType.phone,
          decoration: InputDecoration(
            hintText: '+225 07...',
            filled: true,
            fillColor: const Color(0xFFF8FAFC),
            border: OutlineInputBorder(
              borderRadius: BorderRadius.circular(12),
              borderSide: const BorderSide(color: Color(0xFFE2E8F0)),
            ),
            enabledBorder: OutlineInputBorder(
              borderRadius: BorderRadius.circular(12),
              borderSide: const BorderSide(color: Color(0xFFE2E8F0)),
            ),
          ),
        ),
      ],
    );
  }

  Widget _buildWireFields() {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        const Text('Nom de la banque', style: TextStyle(fontSize: 14, fontWeight: FontWeight.w500, color: Color(0xFF64748B))),
        const SizedBox(height: 8),
        TextField(
          controller: _bankNameController,
          decoration: InputDecoration(
            hintText: 'ex: Ecobank',
            filled: true,
            fillColor: const Color(0xFFF8FAFC),
            border: OutlineInputBorder(borderRadius: BorderRadius.circular(12), borderSide: const BorderSide(color: Color(0xFFE2E8F0))),
            enabledBorder: OutlineInputBorder(borderRadius: BorderRadius.circular(12), borderSide: const BorderSide(color: Color(0xFFE2E8F0))),
          ),
        ),
        const SizedBox(height: 16),
        const Text('IBAN / Numéro de compte', style: TextStyle(fontSize: 14, fontWeight: FontWeight.w500, color: Color(0xFF64748B))),
        const SizedBox(height: 8),
        TextField(
          controller: _ibanController,
          decoration: InputDecoration(
            hintText: 'FR76...',
            filled: true,
            fillColor: const Color(0xFFF8FAFC),
            border: OutlineInputBorder(borderRadius: BorderRadius.circular(12), borderSide: const BorderSide(color: Color(0xFFE2E8F0))),
            enabledBorder: OutlineInputBorder(borderRadius: BorderRadius.circular(12), borderSide: const BorderSide(color: Color(0xFFE2E8F0))),
          ),
        ),
        const SizedBox(height: 16),
        const Text('Nom du bénéficiaire', style: TextStyle(fontSize: 14, fontWeight: FontWeight.w500, color: Color(0xFF64748B))),
        const SizedBox(height: 8),
        TextField(
          controller: _recipientNameController,
          decoration: InputDecoration(
            hintText: 'Jean Dupont',
            filled: true,
            fillColor: const Color(0xFFF8FAFC),
            border: OutlineInputBorder(borderRadius: BorderRadius.circular(12), borderSide: const BorderSide(color: Color(0xFFE2E8F0))),
            enabledBorder: OutlineInputBorder(borderRadius: BorderRadius.circular(12), borderSide: const BorderSide(color: Color(0xFFE2E8F0))),
          ),
        ),
      ],
    );
  }

  Widget _buildCryptoFields() {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        const Text('Adresse du portefeuille', style: TextStyle(fontSize: 14, fontWeight: FontWeight.w500, color: Color(0xFF64748B))),
        const SizedBox(height: 8),
        TextField(
          controller: _recipientController,
          decoration: InputDecoration(
            hintText: '0x... ou bc1...',
            filled: true,
            fillColor: const Color(0xFFF8FAFC),
            border: OutlineInputBorder(borderRadius: BorderRadius.circular(12), borderSide: const BorderSide(color: Color(0xFFE2E8F0))),
            enabledBorder: OutlineInputBorder(borderRadius: BorderRadius.circular(12), borderSide: const BorderSide(color: Color(0xFFE2E8F0))),
            suffixIcon: IconButton(
              icon: const Icon(Icons.qr_code_scanner),
              onPressed: () {/* Scan QR */},
            ),
          ),
        ),
      ],
    );
  }

  Widget _buildDescriptionField() {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        const Text('Note (Optionnel)', style: TextStyle(fontSize: 14, fontWeight: FontWeight.w500, color: Color(0xFF64748B))),
        const SizedBox(height: 8),
        TextField(
          controller: _descriptionController,
          decoration: InputDecoration(
            hintText: 'Ex: Loyer',
            filled: true,
            fillColor: const Color(0xFFF8FAFC),
            border: OutlineInputBorder(borderRadius: BorderRadius.circular(12), borderSide: const BorderSide(color: Color(0xFFE2E8F0))),
            enabledBorder: OutlineInputBorder(borderRadius: BorderRadius.circular(12), borderSide: const BorderSide(color: Color(0xFFE2E8F0))),
          ),
        ),
      ],
    );
  }

  Widget _buildSummary() {
    final amount = double.tryParse(_amountController.text) ?? 0;
    final fee = _getEstimatedFee();
    final total = amount + fee;
    final currency = _getSelectedCurrency();

    return Container(
      padding: const EdgeInsets.all(16),
      decoration: BoxDecoration(
        color: const Color(0xFFF8FAFC),
        borderRadius: BorderRadius.circular(16),
        border: Border.all(color: const Color(0xFFE2E8F0)),
      ),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          const Row(
            children: [
              Text('📝', style: TextStyle(fontSize: 18)),
              SizedBox(width: 8),
              Text('Résumé', style: TextStyle(fontWeight: FontWeight.bold, fontSize: 16)),
            ],
          ),
          const SizedBox(height: 16),
          Row(
            mainAxisAlignment: MainAxisAlignment.spaceBetween,
            children: [
              const Text('Montant', style: TextStyle(color: Color(0xFF64748B))),
              Text('${amount.toStringAsFixed(2)} $currency', style: const TextStyle(fontWeight: FontWeight.w500)),
            ],
          ),
          if (fee > 0) ...[
            const SizedBox(height: 8),
            Row(
              mainAxisAlignment: MainAxisAlignment.spaceBetween,
              children: [
                const Text('Frais estimés', style: TextStyle(color: Color(0xFF667eea))),
                Text('+ ${fee.toStringAsFixed(2)} $currency', style: const TextStyle(color: Color(0xFF667eea), fontWeight: FontWeight.w500)),
              ],
            ),
          ],
          const Divider(height: 24),
          Row(
            mainAxisAlignment: MainAxisAlignment.spaceBetween,
            children: [
              const Text('Total à débiter', style: TextStyle(fontWeight: FontWeight.bold)),
              Text('${total.toStringAsFixed(2)} $currency', style: const TextStyle(fontWeight: FontWeight.bold, fontSize: 18)),
            ],
          ),
        ],
      ),
    );
  }

  Widget _buildSubmitButton() {
    final isP2PValid = _selectedType != 'p2p' || _lookupResult != null;
    final hasWallet = _selectedWalletId != null;
    final hasAmount = (double.tryParse(_amountController.text) ?? 0) > 0;
    final isEnabled = isP2PValid && hasWallet && hasAmount && !_isLoading;

    return GestureDetector(
      onTap: isEnabled ? _handleSubmit : null,
      child: AnimatedContainer(
        duration: const Duration(milliseconds: 200),
        width: double.infinity,
        padding: const EdgeInsets.symmetric(vertical: 18),
        decoration: BoxDecoration(
          gradient: isEnabled
              ? const LinearGradient(colors: [Color(0xFF667eea), Color(0xFF764ba2)])
              : null,
          color: isEnabled ? null : const Color(0xFFCBD5E1),
          borderRadius: BorderRadius.circular(16),
          boxShadow: isEnabled ? [
            BoxShadow(
              color: const Color(0xFF667eea).withOpacity(0.3),
              blurRadius: 12,
              offset: const Offset(0, 4),
            ),
          ] : null,
        ),
        child: Row(
          mainAxisAlignment: MainAxisAlignment.center,
          children: [
            if (_isLoading)
              const SizedBox(
                width: 20,
                height: 20,
                child: CircularProgressIndicator(strokeWidth: 2, color: Colors.white),
              )
            else ...[
              Text(
                'Confirmer le transfert',
                style: TextStyle(
                  color: isEnabled ? Colors.white : const Color(0xFF94A3B8),
                  fontWeight: FontWeight.bold,
                  fontSize: 16,
                ),
              ),
              const SizedBox(width: 8),
              Icon(
                Icons.arrow_forward,
                color: isEnabled ? Colors.white : const Color(0xFF94A3B8),
                size: 20,
              ),
            ],
          ],
        ),
      ),
    );
  }

  String _getSelectedCurrency() {
    final walletState = context.read<WalletBloc>().state;
    if (walletState is WalletLoadedState && _selectedWalletId != null) {
      final wallet = walletState.wallets.firstWhere(
        (w) => w.id == _selectedWalletId,
        orElse: () => walletState.wallets.first,
      );
      return wallet.currency;
    }
    return 'USD';
  }

  double _getEstimatedFee() {
    final amount = double.tryParse(_amountController.text) ?? 0;
    if (amount == 0) return 0;
    
    switch (_selectedType) {
      case 'p2p':
        return 0; // Free internal transfers
      case 'mobile':
        return amount * 0.01; // 1%
      case 'wire':
        return amount * 0.02 < 5 ? 5 : amount * 0.02; // Min 5 or 2%
      case 'crypto':
        return amount * 0.005; // 0.5%
      default:
        return 0;
    }
  }

  void _lookupUser() async {
    final query = _recipientController.text.trim();
    if (query.length < 3) return;

    setState(() {
      _isLookingUp = true;
      _lookupError = null;
      _lookupResult = null;
    });

    // Simulate API call - in real app, call userApi.lookup
    await Future.delayed(const Duration(seconds: 1));

    setState(() {
      _isLookingUp = false;
      // Simulated result
      if (query.contains('@') || query.startsWith('+')) {
        _lookupResult = {
          'first_name': 'Jean',
          'last_name': 'Dupont',
          'id': 'user-123',
        };
      } else {
        _lookupError = 'Utilisateur introuvable';
      }
    });
  }

  void _handleSubmit() async {
    // Require security confirmation before transfer
    final confirmed = await SecurityConfirmation.require(
      context,
      title: 'Confirmer le transfert',
      message: 'Authentifiez-vous pour valider ce transfert',
    );
    if (!confirmed) return;

    setState(() => _isLoading = true);

    try {
      final amount = double.tryParse(_amountController.text) ?? 0;
      
      context.read<TransferBloc>().add(SendTransferEvent(
        type: _selectedType,
        fromWallet: _selectedWalletId!,
        recipient: _selectedType == 'p2p' 
            ? _recipientController.text
            : _selectedType == 'mobile'
                ? _phoneController.text
                : _selectedType == 'wire'
                    ? _recipientNameController.text
                    : _recipientController.text,
        amount: amount,
        memo: _descriptionController.text,
      ));

      // Wait for result
      await Future.delayed(const Duration(seconds: 2));
      
      if (mounted) {
        ScaffoldMessenger.of(context).showSnackBar(
          const SnackBar(
            content: Text('Transfert initié avec succès!'),
            backgroundColor: Color(0xFF10B981),
          ),
        );
        context.go('/wallet');
      }
    } catch (e) {
      if (mounted) {
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(
            content: Text('Erreur: ${e.toString()}'),
            backgroundColor: Colors.red,
          ),
        );
      }
    } finally {
      if (mounted) {
        setState(() => _isLoading = false);
      }
    }
  }
}