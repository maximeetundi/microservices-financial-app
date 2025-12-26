import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';
import 'package:google_fonts/google_fonts.dart';
import 'package:provider/provider.dart';

import '../../../../core/widgets/glass_container.dart';
import '../../../../core/providers/theme_provider.dart';
import '../../../../core/services/auth_api_service.dart';

/// Preferences Page matching web design
class PreferencesPage extends StatefulWidget {
  const PreferencesPage({super.key});

  @override
  State<PreferencesPage> createState() => _PreferencesPageState();
}

class _PreferencesPageState extends State<PreferencesPage> {
  final AuthApiService _authApi = AuthApiService();
  
  bool _isLoading = true;
  bool _isSaving = false;
  
  // Appearance
  String _selectedTheme = 'system';
  String _selectedLanguage = 'fr';
  
  // Notifications
  bool _emailNotifications = true;
  bool _pushNotifications = true;
  bool _smsNotifications = false;
  
  // Transaction alerts
  bool _transactionAlerts = true;
  bool _securityAlerts = true;
  bool _marketingEmails = false;
  bool _priceAlerts = true;
  
  // Display preferences
  String _defaultCurrency = 'EUR';
  String _numberFormat = 'fr-FR';
  bool _showBalances = true;
  
  @override
  void initState() {
    super.initState();
    _loadPreferences();
  }

  Future<void> _loadPreferences() async {
    try {
      final prefs = await _authApi.getPreferences();
      final notifPrefs = await _authApi.getNotificationPrefs();
      
      setState(() {
        // Validate theme
        final savedTheme = prefs['theme'] ?? 'system';
        _selectedTheme = ['light', 'dark', 'system'].contains(savedTheme) ? savedTheme : 'system';

        // Validate language
        final savedLanguage = prefs['language'] ?? 'fr';
        _selectedLanguage = ['fr', 'en', 'es'].contains(savedLanguage) ? savedLanguage : 'fr';
        
        // Validate currency
        final savedCurrency = prefs['default_currency'] ?? 'EUR';
        _defaultCurrency = ['EUR', 'USD', 'GBP', 'XAF'].contains(savedCurrency) ? savedCurrency : 'EUR';

        // Validate number format
        final savedFormat = prefs['number_format'] ?? 'fr-FR';
        _numberFormat = ['fr-FR', 'en-US', 'de-DE'].contains(savedFormat) ? savedFormat : 'fr-FR';

        _showBalances = prefs['show_balances'] ?? true;
        
        _emailNotifications = notifPrefs['email_enabled'] ?? true;
        _pushNotifications = notifPrefs['push_enabled'] ?? true;
        _smsNotifications = notifPrefs['sms_enabled'] ?? false;
        _transactionAlerts = notifPrefs['transaction_alerts'] ?? true;
        _securityAlerts = notifPrefs['security_alerts'] ?? true;
        _marketingEmails = notifPrefs['marketing_emails'] ?? false;
        _priceAlerts = notifPrefs['price_alerts'] ?? true;
        
        _isLoading = false;
      });
    } catch (e) {
      setState(() => _isLoading = false);
      debugPrint('Failed to load preferences: $e');
    }
  }

  Future<void> _savePreferences() async {
    setState(() => _isSaving = true);
    
    try {
      await _authApi.updatePreferences({
        'theme': _selectedTheme,
        'language': _selectedLanguage,
        'default_currency': _defaultCurrency,
        'number_format': _numberFormat,
        'show_balances': _showBalances,
      });
      
      await _authApi.updateNotificationPrefs({
        'email_enabled': _emailNotifications,
        'push_enabled': _pushNotifications,
        'sms_enabled': _smsNotifications,
        'transaction_alerts': _transactionAlerts,
        'security_alerts': _securityAlerts,
        'marketing_emails': _marketingEmails,
        'price_alerts': _priceAlerts,
      });
      
      ScaffoldMessenger.of(context).showSnackBar(
        const SnackBar(
          content: Text('Préférences sauvegardées!'),
          backgroundColor: Colors.green,
        ),
      );
    } catch (e) {
      ScaffoldMessenger.of(context).showSnackBar(
        SnackBar(content: Text('Erreur: $e'), backgroundColor: Colors.red),
      );
    } finally {
      setState(() => _isSaving = false);
    }
  }

  @override
  Widget build(BuildContext context) {
    final isDark = Theme.of(context).brightness == Brightness.dark;

    return Scaffold(
      backgroundColor: isDark ? const Color(0xFF0F172A) : const Color(0xFFF8FAFC),
      body: SafeArea(
        child: _isLoading
            ? const Center(child: CircularProgressIndicator())
            : Column(
                children: [
                  Expanded(
                    child: ListView(
                      padding: const EdgeInsets.all(16),
                      children: [
                        _buildHeader(isDark),
                        const SizedBox(height: 24),
                        _buildAppearanceSection(isDark),
                        const SizedBox(height: 24),
                        _buildNotificationsSection(isDark),
                        const SizedBox(height: 24),
                        _buildDisplaySection(isDark),
                        const SizedBox(height: 100),
                      ],
                    ),
                  ),
                  _buildSaveButton(isDark),
                ],
              ),
      ),
    );
  }

  Widget _buildHeader(bool isDark) {
    return Row(
      children: [
        GlassContainer(
          padding: EdgeInsets.zero,
          width: 40,
          height: 40,
          borderRadius: 12,
          child: IconButton(
            icon: Icon(Icons.arrow_back_ios_new, size: 20,
                color: isDark ? Colors.white : const Color(0xFF1E293B)),
            onPressed: () => context.go('/more/settings'),
          ),
        ),
        const SizedBox(width: 16),
        Expanded(
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              Text(
                '🎨 Préférences',
                style: GoogleFonts.inter(
                  fontSize: 22,
                  fontWeight: FontWeight.bold,
                  color: isDark ? Colors.white : const Color(0xFF1E293B),
                ),
              ),
              Text(
                'Personnalisez votre expérience',
                style: GoogleFonts.inter(
                  fontSize: 12,
                  color: isDark ? const Color(0xFF94A3B8) : const Color(0xFF64748B),
                ),
              ),
            ],
          ),
        ),
      ],
    );
  }

  Widget _buildAppearanceSection(bool isDark) {
    return _buildSection(
      title: 'APPARENCE',
      icon: '🎨',
      isDark: isDark,
      children: [
        _buildDropdownRow(
          label: 'Thème',
          value: _selectedTheme,
          items: const [
            DropdownMenuItem(value: 'light', child: Text('☀️ Clair')),
            DropdownMenuItem(value: 'dark', child: Text('🌙 Sombre')),
            DropdownMenuItem(value: 'system', child: Text('💻 Système')),
          ],
          onChanged: (value) {
            setState(() => _selectedTheme = value!);
            // Apply theme immediately
            final themeProvider = Provider.of<ThemeProvider>(context, listen: false);
            themeProvider.setThemeMode(
              value == 'dark' ? ThemeMode.dark : 
              value == 'light' ? ThemeMode.light : ThemeMode.system
            );
          },
          isDark: isDark,
        ),
        const Divider(height: 24),
        _buildDropdownRow(
          label: 'Langue',
          value: _selectedLanguage,
          items: const [
            DropdownMenuItem(value: 'fr', child: Text('🇫🇷 Français')),
            DropdownMenuItem(value: 'en', child: Text('🇬🇧 English')),
            DropdownMenuItem(value: 'es', child: Text('🇪🇸 Español')),
          ],
          onChanged: (value) => setState(() => _selectedLanguage = value!),
          isDark: isDark,
        ),
      ],
    );
  }

  Widget _buildNotificationsSection(bool isDark) {
    return _buildSection(
      title: 'NOTIFICATIONS',
      icon: '🔔',
      isDark: isDark,
      children: [
        _buildSwitchRow('Notifications email', _emailNotifications, 
            (v) => setState(() => _emailNotifications = v), isDark),
        const Divider(height: 24),
        _buildSwitchRow('Notifications push', _pushNotifications, 
            (v) => setState(() => _pushNotifications = v), isDark),
        const Divider(height: 24),
        _buildSwitchRow('Notifications SMS', _smsNotifications, 
            (v) => setState(() => _smsNotifications = v), isDark),
        const Divider(height: 32),
        Text(
          'TYPES D\'ALERTES',
          style: GoogleFonts.inter(
            fontSize: 11,
            fontWeight: FontWeight.w600,
            color: isDark ? const Color(0xFF64748B) : const Color(0xFF94A3B8),
            letterSpacing: 1,
          ),
        ),
        const SizedBox(height: 16),
        _buildSwitchRow('Alertes transactions', _transactionAlerts, 
            (v) => setState(() => _transactionAlerts = v), isDark),
        const Divider(height: 24),
        _buildSwitchRow('Alertes sécurité', _securityAlerts, 
            (v) => setState(() => _securityAlerts = v), isDark),
        const Divider(height: 24),
        _buildSwitchRow('Alertes prix', _priceAlerts, 
            (v) => setState(() => _priceAlerts = v), isDark),
        const Divider(height: 24),
        _buildSwitchRow('Emails marketing', _marketingEmails, 
            (v) => setState(() => _marketingEmails = v), isDark),
      ],
    );
  }

  Widget _buildDisplaySection(bool isDark) {
    return _buildSection(
      title: 'AFFICHAGE',
      icon: '📊',
      isDark: isDark,
      children: [
        _buildDropdownRow(
          label: 'Devise par défaut',
          value: _defaultCurrency,
          items: const [
            DropdownMenuItem(value: 'EUR', child: Text('€ EUR')),
            DropdownMenuItem(value: 'USD', child: Text('\$ USD')),
            DropdownMenuItem(value: 'GBP', child: Text('£ GBP')),
            DropdownMenuItem(value: 'XAF', child: Text('FCFA XAF')),
          ],
          onChanged: (value) => setState(() => _defaultCurrency = value!),
          isDark: isDark,
        ),
        const Divider(height: 24),
        _buildDropdownRow(
          label: 'Format des nombres',
          value: _numberFormat,
          items: const [
            DropdownMenuItem(value: 'fr-FR', child: Text('1 234,56 (FR)')),
            DropdownMenuItem(value: 'en-US', child: Text('1,234.56 (US)')),
            DropdownMenuItem(value: 'de-DE', child: Text('1.234,56 (DE)')),
          ],
          onChanged: (value) => setState(() => _numberFormat = value!),
          isDark: isDark,
        ),
        const Divider(height: 24),
        _buildSwitchRow('Afficher les soldes', _showBalances, 
            (v) => setState(() => _showBalances = v), isDark),
      ],
    );
  }

  Widget _buildSection({
    required String title,
    required String icon,
    required bool isDark,
    required List<Widget> children,
  }) {
    return Container(
      padding: const EdgeInsets.all(20),
      decoration: BoxDecoration(
        color: isDark ? Colors.white.withOpacity(0.03) : Colors.white,
        borderRadius: BorderRadius.circular(16),
        border: Border.all(
          color: isDark ? Colors.white.withOpacity(0.08) : const Color(0xFFE2E8F0),
        ),
      ),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Row(
            children: [
              Text(icon, style: const TextStyle(fontSize: 20)),
              const SizedBox(width: 10),
              Text(
                title,
                style: GoogleFonts.inter(
                  fontSize: 12,
                  fontWeight: FontWeight.w600,
                  letterSpacing: 1.2,
                  color: isDark ? const Color(0xFF94A3B8) : const Color(0xFF64748B),
                ),
              ),
            ],
          ),
          const SizedBox(height: 20),
          ...children,
        ],
      ),
    );
  }

  Widget _buildSwitchRow(String label, bool value, Function(bool) onChanged, bool isDark) {
    return Row(
      mainAxisAlignment: MainAxisAlignment.spaceBetween,
      children: [
        Text(
          label,
          style: GoogleFonts.inter(
            fontSize: 14,
            color: isDark ? Colors.white : const Color(0xFF1E293B),
          ),
        ),
        Switch.adaptive(
          value: value,
          onChanged: onChanged,
          activeColor: const Color(0xFF6366F1),
        ),
      ],
    );
  }

  Widget _buildDropdownRow<T>({
    required String label,
    required T value,
    required List<DropdownMenuItem<T>> items,
    required Function(T?) onChanged,
    required bool isDark,
  }) {
    return Row(
      mainAxisAlignment: MainAxisAlignment.spaceBetween,
      children: [
        Text(
          label,
          style: GoogleFonts.inter(
            fontSize: 14,
            color: isDark ? Colors.white : const Color(0xFF1E293B),
          ),
        ),
        Container(
          padding: const EdgeInsets.symmetric(horizontal: 12),
          decoration: BoxDecoration(
            color: isDark ? Colors.white.withOpacity(0.05) : const Color(0xFFF1F5F9),
            borderRadius: BorderRadius.circular(8),
          ),
          child: DropdownButtonHideUnderline(
            child: DropdownButton<T>(
              value: value,
              items: items,
              onChanged: onChanged,
              dropdownColor: isDark ? const Color(0xFF1E293B) : Colors.white,
              style: GoogleFonts.inter(
                fontSize: 13,
                color: isDark ? Colors.white : const Color(0xFF1E293B),
              ),
            ),
          ),
        ),
      ],
    );
  }

  Widget _buildSaveButton(bool isDark) {
    return Container(
      padding: const EdgeInsets.all(16),
      decoration: BoxDecoration(
        color: isDark ? const Color(0xFF0F172A) : const Color(0xFFF8FAFC),
        boxShadow: [
          BoxShadow(
            color: Colors.black.withOpacity(0.1),
            blurRadius: 10,
            offset: const Offset(0, -5),
          ),
        ],
      ),
      child: SizedBox(
        width: double.infinity,
        child: ElevatedButton(
          onPressed: _isSaving ? null : _savePreferences,
          style: ElevatedButton.styleFrom(
            backgroundColor: const Color(0xFF6366F1),
            padding: const EdgeInsets.symmetric(vertical: 16),
            shape: RoundedRectangleBorder(
              borderRadius: BorderRadius.circular(12),
            ),
          ),
          child: _isSaving
              ? const SizedBox(
                  height: 20,
                  width: 20,
                  child: CircularProgressIndicator(color: Colors.white, strokeWidth: 2),
                )
              : Text(
                  '💾 Sauvegarder les préférences',
                  style: GoogleFonts.inter(
                    fontSize: 14,
                    fontWeight: FontWeight.w600,
                    color: Colors.white,
                  ),
                ),
        ),
      ),
    );
  }
}
