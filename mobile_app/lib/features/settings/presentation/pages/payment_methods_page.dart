import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';
import 'package:go_router/go_router.dart';
import 'package:google_fonts/google_fonts.dart';

import '../../../../core/theme/app_theme.dart';
import '../../../../core/widgets/glass_container.dart';
import '../../../cards/presentation/bloc/cards_bloc.dart';

/// Payment Methods Page - Match web payment-methods.vue
class PaymentMethodsPage extends StatefulWidget {
  const PaymentMethodsPage({super.key});

  @override
  State<PaymentMethodsPage> createState() => _PaymentMethodsPageState();
}

class _PaymentMethodsPageState extends State<PaymentMethodsPage> {
  List<Map<String, dynamic>> _bankAccounts = [];
  List<Map<String, dynamic>> _mobileAccounts = [];

  @override
  void initState() {
    super.initState();
    context.read<CardsBloc>().add(LoadCardsEvent());
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
              _buildAppBar(isDark),
              Expanded(
                child: SingleChildScrollView(
                  padding: const EdgeInsets.all(16),
                  child: Column(
                    crossAxisAlignment: CrossAxisAlignment.start,
                    children: [
                      _buildCardsSection(isDark),
                      const SizedBox(height: 24),
                      _buildBankAccountsSection(isDark),
                      const SizedBox(height: 24),
                      _buildMobileMoneySection(isDark),
                      const SizedBox(height: 32),
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

  Widget _buildAppBar(bool isDark) {
    return Padding(
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
              onPressed: () => context.go('/more/settings'),
            ),
          ),
          const SizedBox(width: 16),
          Expanded(
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Text(
                  '💳 Moyens de paiement',
                  style: GoogleFonts.inter(
                    fontSize: 20,
                    fontWeight: FontWeight.bold,
                    color: isDark ? Colors.white : const Color(0xFF1E293B),
                  ),
                ),
                Text(
                  'Gérez vos cartes et comptes bancaires',
                  style: GoogleFonts.inter(
                    fontSize: 12,
                    color: isDark ? const Color(0xFF94A3B8) : const Color(0xFF64748B),
                  ),
                ),
              ],
            ),
          ),
        ],
      ),
    );
  }

  Widget _buildCardsSection(bool isDark) {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        _buildSectionHeader('Mes cartes', () => context.push('/more/cards'), isDark),
        const SizedBox(height: 12),
        BlocBuilder<CardsBloc, CardsState>(
          builder: (context, state) {
            if (state is CardsLoadingState) {
              return const Center(child: CircularProgressIndicator());
            }
            
            if (state is CardsLoadedState && state.cards.isNotEmpty) {
              return Column(
                children: state.cards.map((card) => _buildCardItem(card, isDark)).toList(),
              );
            }
            
            return _buildEmptyState('💳', 'Aucune carte enregistrée', 
                'Créer une carte', () => context.push('/more/cards'), isDark);
          },
        ),
      ],
    );
  }

  Widget _buildCardItem(Map<String, dynamic> card, bool isDark) {
    final isVirtual = card['is_virtual'] == true;
    final status = card['status']?.toString() ?? 'active';
    final cardNumber = card['card_number']?.toString() ?? '0000';
    final lastFour = cardNumber.length >= 4 ? cardNumber.substring(cardNumber.length - 4) : '0000';
    
    return Container(
      margin: const EdgeInsets.only(bottom: 12),
      padding: const EdgeInsets.all(16),
      decoration: BoxDecoration(
        color: isDark 
            ? Colors.white.withOpacity(0.05)
            : Colors.white,
        borderRadius: BorderRadius.circular(16),
        border: Border.all(
          color: isDark 
              ? Colors.white.withOpacity(0.1)
              : const Color(0xFFE2E8F0),
        ),
      ),
      child: Row(
        children: [
          // Card Preview
          Container(
            width: 56,
            height: 36,
            decoration: BoxDecoration(
              gradient: isVirtual 
                  ? const LinearGradient(colors: [Color(0xFF6366F1), Color(0xFF8B5CF6)])
                  : const LinearGradient(colors: [Color(0xFF1E1E2F), Color(0xFF3D3D5C)]),
              borderRadius: BorderRadius.circular(6),
            ),
            child: Column(
              mainAxisAlignment: MainAxisAlignment.center,
              children: [
                Text(
                  isVirtual ? 'Virtuelle' : 'Physique',
                  style: GoogleFonts.inter(
                    fontSize: 6,
                    color: Colors.white70,
                  ),
                ),
                Text(
                  '•••• $lastFour',
                  style: GoogleFonts.sourceCodePro(
                    fontSize: 8,
                    color: Colors.white,
                  ),
                ),
              ],
            ),
          ),
          const SizedBox(width: 16),
          
          // Card Info
          Expanded(
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Text(
                  isVirtual ? 'Carte Virtuelle' : 'Carte Physique',
                  style: GoogleFonts.inter(
                    fontSize: 14,
                    fontWeight: FontWeight.w500,
                    color: isDark ? Colors.white : const Color(0xFF1E293B),
                  ),
                ),
                Text(
                  'Expire ${card['expiry_month'] ?? '12'}/${card['expiry_year'] ?? '28'}',
                  style: GoogleFonts.inter(
                    fontSize: 12,
                    color: isDark ? const Color(0xFF64748B) : const Color(0xFF94A3B8),
                  ),
                ),
              ],
            ),
          ),
          
          // Status Badge
          Container(
            padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 4),
            decoration: BoxDecoration(
              color: status == 'active' 
                  ? const Color(0xFF22C55E).withOpacity(0.15)
                  : const Color(0xFF6B7280).withOpacity(0.15),
              borderRadius: BorderRadius.circular(6),
            ),
            child: Text(
              status == 'active' ? 'Active' : 'Inactive',
              style: GoogleFonts.inter(
                fontSize: 10,
                fontWeight: FontWeight.bold,
                color: status == 'active' 
                    ? const Color(0xFF22C55E)
                    : const Color(0xFF9CA3AF),
              ),
            ),
          ),
        ],
      ),
    );
  }

  Widget _buildBankAccountsSection(bool isDark) {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        _buildSectionHeader('Comptes bancaires', () => _showAddBankModal(isDark), isDark),
        const SizedBox(height: 12),
        if (_bankAccounts.isEmpty)
          _buildEmptyState('🏦', 'Aucun compte bancaire lié', 
              'Ajouter un compte', () => _showAddBankModal(isDark), isDark)
        else
          Column(
            children: _bankAccounts.map((bank) => _buildBankItem(bank, isDark)).toList(),
          ),
      ],
    );
  }

  Widget _buildBankItem(Map<String, dynamic> bank, bool isDark) {
    return Container(
      margin: const EdgeInsets.only(bottom: 12),
      padding: const EdgeInsets.all(16),
      decoration: BoxDecoration(
        color: isDark 
            ? Colors.white.withOpacity(0.05)
            : Colors.white,
        borderRadius: BorderRadius.circular(14),
        border: Border.all(
          color: isDark 
              ? Colors.white.withOpacity(0.1)
              : const Color(0xFFE2E8F0),
        ),
      ),
      child: Row(
        children: [
          const Text('🏦', style: TextStyle(fontSize: 24)),
          const SizedBox(width: 12),
          Expanded(
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Text(
                  bank['bank_name'] ?? 'Banque',
                  style: GoogleFonts.inter(
                    fontSize: 14,
                    fontWeight: FontWeight.w500,
                    color: isDark ? Colors.white : const Color(0xFF1E293B),
                  ),
                ),
                Text(
                  bank['account_number'] ?? '••••',
                  style: GoogleFonts.inter(
                    fontSize: 12,
                    color: isDark ? const Color(0xFF64748B) : const Color(0xFF94A3B8),
                  ),
                ),
              ],
            ),
          ),
          GestureDetector(
            onTap: () => _removeBank(bank['id']),
            child: Container(
              width: 28,
              height: 28,
              decoration: BoxDecoration(
                color: const Color(0xFFEF4444).withOpacity(0.15),
                borderRadius: BorderRadius.circular(8),
              ),
              child: const Center(
                child: Text('✕', style: TextStyle(color: Color(0xFFEF4444))),
              ),
            ),
          ),
        ],
      ),
    );
  }

  Widget _buildMobileMoneySection(bool isDark) {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        _buildSectionHeader('Mobile Money', () => _showAddMobileModal(isDark), isDark),
        const SizedBox(height: 12),
        if (_mobileAccounts.isEmpty)
          _buildEmptyState('📱', 'Aucun compte Mobile Money', 
              'Ajouter', () => _showAddMobileModal(isDark), isDark)
        else
          Column(
            children: _mobileAccounts.map((mobile) => _buildMobileItem(mobile, isDark)).toList(),
          ),
      ],
    );
  }

  Widget _buildMobileItem(Map<String, dynamic> mobile, bool isDark) {
    final icon = _getOperatorIcon(mobile['operator'] ?? '');
    
    return Container(
      margin: const EdgeInsets.only(bottom: 12),
      padding: const EdgeInsets.all(16),
      decoration: BoxDecoration(
        color: isDark 
            ? Colors.white.withOpacity(0.05)
            : Colors.white,
        borderRadius: BorderRadius.circular(14),
        border: Border.all(
          color: isDark 
              ? Colors.white.withOpacity(0.1)
              : const Color(0xFFE2E8F0),
        ),
      ),
      child: Row(
        children: [
          Text(icon, style: const TextStyle(fontSize: 24)),
          const SizedBox(width: 12),
          Expanded(
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Text(
                  mobile['operator'] ?? 'Mobile Money',
                  style: GoogleFonts.inter(
                    fontSize: 14,
                    fontWeight: FontWeight.w500,
                    color: isDark ? Colors.white : const Color(0xFF1E293B),
                  ),
                ),
                Text(
                  mobile['phone_number'] ?? '',
                  style: GoogleFonts.inter(
                    fontSize: 12,
                    color: isDark ? const Color(0xFF64748B) : const Color(0xFF94A3B8),
                  ),
                ),
              ],
            ),
          ),
          GestureDetector(
            onTap: () => _removeMobile(mobile['id']),
            child: Container(
              width: 28,
              height: 28,
              decoration: BoxDecoration(
                color: const Color(0xFFEF4444).withOpacity(0.15),
                borderRadius: BorderRadius.circular(8),
              ),
              child: const Center(
                child: Text('✕', style: TextStyle(color: Color(0xFFEF4444))),
              ),
            ),
          ),
        ],
      ),
    );
  }

  Widget _buildSectionHeader(String title, VoidCallback onAdd, bool isDark) {
    return Row(
      mainAxisAlignment: MainAxisAlignment.spaceBetween,
      children: [
        Text(
          title.toUpperCase(),
          style: GoogleFonts.inter(
            fontSize: 12,
            fontWeight: FontWeight.w600,
            color: isDark ? const Color(0xFF64748B) : const Color(0xFF94A3B8),
            letterSpacing: 0.5,
          ),
        ),
        GestureDetector(
          onTap: onAdd,
          child: Container(
            padding: const EdgeInsets.symmetric(horizontal: 12, vertical: 6),
            decoration: BoxDecoration(
              border: Border.all(color: const Color(0xFF6366F1).withOpacity(0.3)),
              borderRadius: BorderRadius.circular(8),
            ),
            child: Text(
              '+ Ajouter',
              style: GoogleFonts.inter(
                fontSize: 12,
                fontWeight: FontWeight.w600,
                color: const Color(0xFF6366F1),
              ),
            ),
          ),
        ),
      ],
    );
  }

  Widget _buildEmptyState(String emoji, String message, String action, VoidCallback onTap, bool isDark) {
    return Container(
      width: double.infinity,
      padding: const EdgeInsets.all(32),
      decoration: BoxDecoration(
        color: isDark 
            ? Colors.white.withOpacity(0.03)
            : Colors.white,
        borderRadius: BorderRadius.circular(16),
      ),
      child: Column(
        children: [
          Text(emoji, style: TextStyle(fontSize: 40, color: Colors.grey.withOpacity(0.5))),
          const SizedBox(height: 8),
          Text(
            message,
            style: GoogleFonts.inter(
              color: isDark ? const Color(0xFF64748B) : const Color(0xFF94A3B8),
            ),
          ),
          const SizedBox(height: 12),
          GestureDetector(
            onTap: onTap,
            child: Text(
              action,
              style: GoogleFonts.inter(
                color: const Color(0xFF6366F1),
                fontWeight: FontWeight.w500,
              ),
            ),
          ),
        ],
      ),
    );
  }

  String _getOperatorIcon(String operator) {
    switch (operator) {
      case 'Orange Money':
        return '🟠';
      case 'Wave':
        return '🔵';
      case 'MTN MoMo':
        return '🟡';
      case 'Free Money':
        return '🟢';
      default:
        return '📱';
    }
  }

  void _showAddBankModal(bool isDark) {
    final nameController = TextEditingController();
    final ibanController = TextEditingController();
    
    showModalBottomSheet(
      context: context,
      isScrollControlled: true,
      backgroundColor: Colors.transparent,
      builder: (context) => Container(
        margin: const EdgeInsets.all(16),
        padding: EdgeInsets.only(
          bottom: MediaQuery.of(context).viewInsets.bottom,
        ),
        decoration: BoxDecoration(
          color: isDark ? const Color(0xFF1A1A2E) : Colors.white,
          borderRadius: BorderRadius.circular(24),
        ),
        child: Padding(
          padding: const EdgeInsets.all(24),
          child: Column(
            mainAxisSize: MainAxisSize.min,
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              Text(
                '🏦 Ajouter un compte bancaire',
                style: GoogleFonts.inter(
                  fontSize: 18,
                  fontWeight: FontWeight.bold,
                  color: isDark ? Colors.white : const Color(0xFF1E293B),
                ),
              ),
              const SizedBox(height: 24),
              _buildTextField('Nom de la banque', 'Ex: Société Générale', nameController, isDark),
              const SizedBox(height: 16),
              _buildTextField('IBAN / Numéro de compte', 'Ex: SN01 1234 5678...', ibanController, isDark),
              const SizedBox(height: 24),
              Row(
                children: [
                  Expanded(
                    child: GestureDetector(
                      onTap: () => Navigator.pop(context),
                      child: Container(
                        padding: const EdgeInsets.symmetric(vertical: 14),
                        decoration: BoxDecoration(
                          color: isDark ? Colors.white.withOpacity(0.1) : const Color(0xFFF1F5F9),
                          borderRadius: BorderRadius.circular(12),
                        ),
                        child: Center(
                          child: Text(
                            'Annuler',
                            style: GoogleFonts.inter(
                              fontWeight: FontWeight.w600,
                              color: isDark ? Colors.white : const Color(0xFF64748B),
                            ),
                          ),
                        ),
                      ),
                    ),
                  ),
                  const SizedBox(width: 12),
                  Expanded(
                    child: GestureDetector(
                      onTap: () {
                        if (nameController.text.isNotEmpty && ibanController.text.isNotEmpty) {
                          setState(() {
                            _bankAccounts.add({
                              'id': DateTime.now().millisecondsSinceEpoch,
                              'bank_name': nameController.text,
                              'account_number': '••••${ibanController.text.substring(ibanController.text.length.clamp(0, 4) > 4 ? ibanController.text.length - 4 : 0)}',
                            });
                          });
                          Navigator.pop(context);
                        }
                      },
                      child: Container(
                        padding: const EdgeInsets.symmetric(vertical: 14),
                        decoration: BoxDecoration(
                          color: const Color(0xFF6366F1),
                          borderRadius: BorderRadius.circular(12),
                        ),
                        child: Center(
                          child: Text(
                            'Ajouter',
                            style: GoogleFonts.inter(
                              fontWeight: FontWeight.w600,
                              color: Colors.white,
                            ),
                          ),
                        ),
                      ),
                    ),
                  ),
                ],
              ),
            ],
          ),
        ),
      ),
    );
  }

  void _showAddMobileModal(bool isDark) {
    String selectedOperator = '';
    final phoneController = TextEditingController();
    
    showModalBottomSheet(
      context: context,
      isScrollControlled: true,
      backgroundColor: Colors.transparent,
      builder: (context) => StatefulBuilder(
        builder: (context, setModalState) => Container(
          margin: const EdgeInsets.all(16),
          padding: EdgeInsets.only(
            bottom: MediaQuery.of(context).viewInsets.bottom,
          ),
          decoration: BoxDecoration(
            color: isDark ? const Color(0xFF1A1A2E) : Colors.white,
            borderRadius: BorderRadius.circular(24),
          ),
          child: Padding(
            padding: const EdgeInsets.all(24),
            child: Column(
              mainAxisSize: MainAxisSize.min,
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Text(
                  '📱 Ajouter Mobile Money',
                  style: GoogleFonts.inter(
                    fontSize: 18,
                    fontWeight: FontWeight.bold,
                    color: isDark ? Colors.white : const Color(0xFF1E293B),
                  ),
                ),
                const SizedBox(height: 24),
                
                // Operator Selection
                Text(
                  'Opérateur',
                  style: GoogleFonts.inter(
                    fontSize: 12,
                    color: isDark ? const Color(0xFF64748B) : const Color(0xFF94A3B8),
                  ),
                ),
                const SizedBox(height: 8),
                Wrap(
                  spacing: 8,
                  runSpacing: 8,
                  children: ['Orange Money', 'Wave', 'MTN MoMo', 'Free Money'].map((op) {
                    final isSelected = selectedOperator == op;
                    return GestureDetector(
                      onTap: () => setModalState(() => selectedOperator = op),
                      child: Container(
                        padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 10),
                        decoration: BoxDecoration(
                          color: isSelected 
                              ? const Color(0xFF6366F1).withOpacity(0.15)
                              : (isDark ? Colors.white.withOpacity(0.05) : const Color(0xFFF8FAFC)),
                          border: Border.all(
                            color: isSelected 
                                ? const Color(0xFF6366F1)
                                : (isDark ? Colors.white.withOpacity(0.1) : const Color(0xFFE2E8F0)),
                          ),
                          borderRadius: BorderRadius.circular(8),
                        ),
                        child: Row(
                          mainAxisSize: MainAxisSize.min,
                          children: [
                            Text(_getOperatorIcon(op)),
                            const SizedBox(width: 8),
                            Text(
                              op,
                              style: GoogleFonts.inter(
                                fontSize: 13,
                                color: isSelected 
                                    ? const Color(0xFF6366F1)
                                    : (isDark ? Colors.white : const Color(0xFF1E293B)),
                              ),
                            ),
                          ],
                        ),
                      ),
                    );
                  }).toList(),
                ),
                const SizedBox(height: 16),
                _buildTextField('Numéro de téléphone', '+221 77 123 45 67', phoneController, isDark),
                const SizedBox(height: 24),
                Row(
                  children: [
                    Expanded(
                      child: GestureDetector(
                        onTap: () => Navigator.pop(context),
                        child: Container(
                          padding: const EdgeInsets.symmetric(vertical: 14),
                          decoration: BoxDecoration(
                            color: isDark ? Colors.white.withOpacity(0.1) : const Color(0xFFF1F5F9),
                            borderRadius: BorderRadius.circular(12),
                          ),
                          child: Center(
                            child: Text(
                              'Annuler',
                              style: GoogleFonts.inter(
                                fontWeight: FontWeight.w600,
                                color: isDark ? Colors.white : const Color(0xFF64748B),
                              ),
                            ),
                          ),
                        ),
                      ),
                    ),
                    const SizedBox(width: 12),
                    Expanded(
                      child: GestureDetector(
                        onTap: () {
                          if (selectedOperator.isNotEmpty && phoneController.text.isNotEmpty) {
                            setState(() {
                              _mobileAccounts.add({
                                'id': DateTime.now().millisecondsSinceEpoch,
                                'operator': selectedOperator,
                                'phone_number': phoneController.text,
                              });
                            });
                            Navigator.pop(context);
                          }
                        },
                        child: Container(
                          padding: const EdgeInsets.symmetric(vertical: 14),
                          decoration: BoxDecoration(
                            color: const Color(0xFF6366F1),
                            borderRadius: BorderRadius.circular(12),
                          ),
                          child: Center(
                            child: Text(
                              'Ajouter',
                              style: GoogleFonts.inter(
                                fontWeight: FontWeight.w600,
                                color: Colors.white,
                              ),
                            ),
                          ),
                        ),
                      ),
                    ),
                  ],
                ),
              ],
            ),
          ),
        ),
      ),
    );
  }

  Widget _buildTextField(String label, String hint, TextEditingController controller, bool isDark) {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Text(
          label,
          style: GoogleFonts.inter(
            fontSize: 12,
            color: isDark ? const Color(0xFF64748B) : const Color(0xFF94A3B8),
          ),
        ),
        const SizedBox(height: 8),
        TextField(
          controller: controller,
          style: GoogleFonts.inter(
            color: isDark ? Colors.white : const Color(0xFF1E293B),
          ),
          decoration: InputDecoration(
            hintText: hint,
            hintStyle: GoogleFonts.inter(
              color: isDark ? const Color(0xFF475569) : const Color(0xFFCBD5E1),
            ),
            filled: true,
            fillColor: isDark ? Colors.white.withOpacity(0.05) : const Color(0xFFF8FAFC),
            border: OutlineInputBorder(
              borderRadius: BorderRadius.circular(12),
              borderSide: BorderSide(
                color: isDark ? Colors.white.withOpacity(0.1) : const Color(0xFFE2E8F0),
              ),
            ),
            enabledBorder: OutlineInputBorder(
              borderRadius: BorderRadius.circular(12),
              borderSide: BorderSide(
                color: isDark ? Colors.white.withOpacity(0.1) : const Color(0xFFE2E8F0),
              ),
            ),
            focusedBorder: OutlineInputBorder(
              borderRadius: BorderRadius.circular(12),
              borderSide: const BorderSide(color: Color(0xFF6366F1)),
            ),
          ),
        ),
      ],
    );
  }

  void _removeBank(int id) {
    setState(() {
      _bankAccounts.removeWhere((b) => b['id'] == id);
    });
  }

  void _removeMobile(int id) {
    setState(() {
      _mobileAccounts.removeWhere((m) => m['id'] == id);
    });
  }
}
