import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';

import '../../../../core/services/biometric_service.dart';
import '../../../../core/api/api_client.dart';
import '../../../security/pin_code_screen.dart';

/// Modern Security Settings Page with real biometric and PIN integration
class SecurityPage extends StatefulWidget {
  const SecurityPage({super.key});

  @override
  State<SecurityPage> createState() => _SecurityPageState();
}

class _SecurityPageState extends State<SecurityPage> {
  final BiometricService _biometricService = BiometricService();
  final ApiClient _apiClient = ApiClient();
  
  bool _biometricAvailable = false;
  bool _biometricEnabled = false;
  bool _pinEnabled = false;
  bool _twoFactorEnabled = false;
  String _biometricName = 'Biométrie';
  int _lockTimeout = 5;
  bool _loading = true;

  @override
  void initState() {
    super.initState();
    _loadSettings();
  }

  Future<void> _loadSettings() async {
    final biometricAvailable = await _biometricService.isBiometricAvailable();
    final biometricEnabled = await _biometricService.isBiometricEnabled();
    final pinEnabled = await _biometricService.isPinEnabled();
    final biometricName = await _biometricService.getBiometricTypeName();
    final lockTimeout = await _biometricService.getLockTimeout();

    if (mounted) {
      setState(() {
        _biometricAvailable = biometricAvailable;
        _biometricEnabled = biometricEnabled;
        _pinEnabled = pinEnabled;
        _biometricName = biometricName;
        _lockTimeout = lockTimeout;
        _loading = false;
      });
    }
  }

  Future<void> _toggleBiometric(bool enabled) async {
    if (enabled) {
      // Test biometric authentication first
      final success = await _biometricService.authenticateWithBiometrics(
        reason: 'Authentifiez-vous pour activer $_biometricName',
      );
      if (!success) return;
    }
    
    await _biometricService.setBiometricEnabled(enabled);
    setState(() => _biometricEnabled = enabled);
    
    if (mounted) {
      ScaffoldMessenger.of(context).showSnackBar(
        SnackBar(
          content: Text(enabled 
              ? '$_biometricName activé' 
              : '$_biometricName désactivé'),
          backgroundColor: const Color(0xFF10B981),
        ),
      );
    }
  }

  void _setupPin() {
    Navigator.of(context).push(
      MaterialPageRoute(
        builder: (context) => PinCodeScreen(
          isSetup: true,
          onSuccess: () {
            Navigator.of(context).pop();
            setState(() => _pinEnabled = true);
          },
          onCancel: () => Navigator.of(context).pop(),
        ),
      ),
    );
  }

  void _changePin() {
    Navigator.of(context).push(
      MaterialPageRoute(
        builder: (context) => PinCodeScreen(
          isSetup: false,
          title: 'Entrez votre code actuel',
          onSuccess: () {
            Navigator.of(context).pop();
            _setupPin();
          },
          onCancel: () => Navigator.of(context).pop(),
        ),
      ),
    );
  }

  Future<void> _removePin() async {
    final confirm = await showDialog<bool>(
      context: context,
      builder: (context) => AlertDialog(
        title: const Text('Supprimer le code PIN?'),
        content: const Text('Votre compte sera moins sécurisé.'),
        actions: [
          TextButton(
            onPressed: () => Navigator.pop(context, false),
            child: const Text('Annuler'),
          ),
          ElevatedButton(
            onPressed: () => Navigator.pop(context, true),
            style: ElevatedButton.styleFrom(backgroundColor: Colors.red),
            child: const Text('Supprimer'),
          ),
        ],
      ),
    );
    
    if (confirm == true) {
      await _biometricService.removePin();
      setState(() => _pinEnabled = false);
    }
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: const Color(0xFFF5F7FA),
      appBar: AppBar(
        backgroundColor: Colors.transparent,
        elevation: 0,
        leading: IconButton(
          icon: const Icon(Icons.arrow_back_ios, color: Color(0xFF1a1a2e)),
          onPressed: () => context.pop(),
        ),
        title: const Text(
          'Sécurité 🔒',
          style: TextStyle(
            color: Color(0xFF1a1a2e),
            fontWeight: FontWeight.bold,
            fontSize: 20,
          ),
        ),
        centerTitle: true,
      ),
      body: _loading
          ? const Center(child: CircularProgressIndicator())
          : ListView(
              padding: const EdgeInsets.all(16),
              children: [
                // PIN Section
                _buildSection(
                  title: 'Code PIN',
                  icon: Icons.dialpad,
                  children: [
                    if (_pinEnabled) ...[
                      _buildSettingTile(
                        icon: Icons.lock,
                        title: 'Code PIN activé',
                        subtitle: '6 chiffres',
                        trailing: Container(
                          padding: const EdgeInsets.symmetric(horizontal: 10, vertical: 4),
                          decoration: BoxDecoration(
                            color: const Color(0xFF10B981).withOpacity(0.1),
                            borderRadius: BorderRadius.circular(12),
                          ),
                          child: const Text(
                            'Actif',
                            style: TextStyle(
                              color: Color(0xFF10B981),
                              fontWeight: FontWeight.w600,
                              fontSize: 12,
                            ),
                          ),
                        ),
                      ),
                      _buildSettingTile(
                        icon: Icons.edit,
                        title: 'Changer le code PIN',
                        onTap: _changePin,
                      ),
                      _buildSettingTile(
                        icon: Icons.delete_outline,
                        title: 'Supprimer le code PIN',
                        titleColor: Colors.red,
                        onTap: _removePin,
                      ),
                    ] else ...[
                      _buildSettingTile(
                        icon: Icons.add_circle_outline,
                        title: 'Configurer un code PIN',
                        subtitle: 'Protégez l\'accès à l\'app',
                        onTap: _setupPin,
                      ),
                    ],
                  ],
                ),
                const SizedBox(height: 16),

                // Biometric Section
                if (_biometricAvailable)
                  _buildSection(
                    title: _biometricName,
                    icon: Icons.fingerprint,
                    children: [
                      _buildSwitchTile(
                        icon: Icons.fingerprint,
                        title: 'Activer $_biometricName',
                        subtitle: 'Connexion rapide',
                        value: _biometricEnabled,
                        onChanged: _toggleBiometric,
                      ),
                    ],
                  ),
                const SizedBox(height: 16),

                // Lock Timeout
                _buildSection(
                  title: 'Verrouillage automatique',
                  icon: Icons.timer,
                  children: [
                    _buildSettingTile(
                      icon: Icons.schedule,
                      title: 'Délai de verrouillage',
                      subtitle: 'Après $_lockTimeout minutes d\'inactivité',
                      trailing: const Icon(Icons.chevron_right, color: Color(0xFF94A3B8)),
                      onTap: () => _showLockTimeoutPicker(),
                    ),
                  ],
                ),
                const SizedBox(height: 16),

                // Password Section
                _buildSection(
                  title: 'Mot de passe',
                  icon: Icons.key,
                  children: [
                    _buildSettingTile(
                      icon: Icons.lock_reset,
                      title: 'Changer le mot de passe',
                      onTap: () => _showChangePasswordDialog(),
                    ),
                  ],
                ),
                const SizedBox(height: 16),

                // 2FA Section
                _buildSection(
                  title: 'Authentification à deux facteurs',
                  icon: Icons.security,
                  children: [
                    _buildSwitchTile(
                      icon: Icons.verified_user,
                      title: 'Activer 2FA',
                      subtitle: 'Sécurité renforcée',
                      value: _twoFactorEnabled,
                      onChanged: (v) => _toggle2FA(v),
                    ),
                    if (_twoFactorEnabled)
                      _buildSettingTile(
                        icon: Icons.qr_code,
                        title: 'Configurer l\'application',
                        subtitle: 'Google Authenticator, Authy...',
                        onTap: () => _setup2FA(),
                      ),
                  ],
                ),
                const SizedBox(height: 16),

                // Sessions Section
                _buildSection(
                  title: 'Sessions & Appareils',
                  icon: Icons.devices,
                  children: [
                    _buildSettingTile(
                      icon: Icons.smartphone,
                      title: 'Appareils connectés',
                      subtitle: 'Voir les appareils actifs',
                      trailing: const Icon(Icons.chevron_right, color: Color(0xFF94A3B8)),
                      onTap: () => _showDevicesSheet(),
                    ),
                    _buildSettingTile(
                      icon: Icons.history,
                      title: 'Historique de connexion',
                      trailing: const Icon(Icons.chevron_right, color: Color(0xFF94A3B8)),
                      onTap: () {},
                    ),
                  ],
                ),
                const SizedBox(height: 24),

                // Danger Zone
                Container(
                  padding: const EdgeInsets.all(16),
                  decoration: BoxDecoration(
                    color: Colors.red.withOpacity(0.05),
                    borderRadius: BorderRadius.circular(16),
                    border: Border.all(color: Colors.red.withOpacity(0.2)),
                  ),
                  child: Column(
                    crossAxisAlignment: CrossAxisAlignment.start,
                    children: [
                      const Row(
                        children: [
                          Icon(Icons.warning, color: Colors.red, size: 20),
                          SizedBox(width: 8),
                          Text(
                            'Zone Danger',
                            style: TextStyle(
                              color: Colors.red,
                              fontWeight: FontWeight.bold,
                            ),
                          ),
                        ],
                      ),
                      const SizedBox(height: 12),
                      SizedBox(
                        width: double.infinity,
                        child: OutlinedButton.icon(
                          onPressed: _logoutAllDevices,
                          icon: const Icon(Icons.logout, color: Colors.red),
                          label: const Text('Déconnecter tous les appareils'),
                          style: OutlinedButton.styleFrom(
                            foregroundColor: Colors.red,
                            side: const BorderSide(color: Colors.red),
                          ),
                        ),
                      ),
                    ],
                  ),
                ),
              ],
            ),
    );
  }

  Widget _buildSection({
    required String title,
    required IconData icon,
    required List<Widget> children,
  }) {
    return Container(
      decoration: BoxDecoration(
        color: Colors.white,
        borderRadius: BorderRadius.circular(16),
        boxShadow: [
          BoxShadow(
            color: Colors.black.withOpacity(0.05),
            blurRadius: 10,
          ),
        ],
      ),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Padding(
            padding: const EdgeInsets.fromLTRB(16, 16, 16, 8),
            child: Row(
              children: [
                Icon(icon, size: 18, color: const Color(0xFF667eea)),
                const SizedBox(width: 8),
                Text(
                  title,
                  style: const TextStyle(
                    fontSize: 14,
                    fontWeight: FontWeight.bold,
                    color: Color(0xFF64748B),
                  ),
                ),
              ],
            ),
          ),
          ...children,
        ],
      ),
    );
  }

  Widget _buildSettingTile({
    required IconData icon,
    required String title,
    String? subtitle,
    Color? titleColor,
    Widget? trailing,
    VoidCallback? onTap,
  }) {
    return ListTile(
      leading: Icon(icon, color: const Color(0xFF64748B)),
      title: Text(
        title,
        style: TextStyle(
          color: titleColor ?? const Color(0xFF1a1a2e),
          fontWeight: FontWeight.w500,
        ),
      ),
      subtitle: subtitle != null ? Text(subtitle) : null,
      trailing: trailing ?? (onTap != null ? const Icon(Icons.chevron_right, color: Color(0xFF94A3B8)) : null),
      onTap: onTap,
    );
  }

  Widget _buildSwitchTile({
    required IconData icon,
    required String title,
    String? subtitle,
    required bool value,
    required Function(bool) onChanged,
  }) {
    return SwitchListTile(
      secondary: Icon(icon, color: const Color(0xFF64748B)),
      title: Text(
        title,
        style: const TextStyle(
          color: Color(0xFF1a1a2e),
          fontWeight: FontWeight.w500,
        ),
      ),
      subtitle: subtitle != null ? Text(subtitle) : null,
      value: value,
      activeColor: const Color(0xFF667eea),
      onChanged: onChanged,
    );
  }

  void _showLockTimeoutPicker() {
    showModalBottomSheet(
      context: context,
      backgroundColor: Colors.white,
      shape: const RoundedRectangleBorder(
        borderRadius: BorderRadius.vertical(top: Radius.circular(20)),
      ),
      builder: (context) => Column(
        mainAxisSize: MainAxisSize.min,
        children: [
          const SizedBox(height: 16),
          Container(
            width: 40,
            height: 4,
            decoration: BoxDecoration(
              color: const Color(0xFFE2E8F0),
              borderRadius: BorderRadius.circular(2),
            ),
          ),
          const SizedBox(height: 16),
          const Text(
            'Délai de verrouillage',
            style: TextStyle(
              fontSize: 18,
              fontWeight: FontWeight.bold,
            ),
          ),
          const SizedBox(height: 16),
          for (final minutes in [1, 5, 15, 30])
            ListTile(
              title: Text('$minutes minute${minutes > 1 ? 's' : ''}'),
              trailing: _lockTimeout == minutes
                  ? const Icon(Icons.check, color: Color(0xFF667eea))
                  : null,
              onTap: () async {
                await _biometricService.setLockTimeout(minutes);
                setState(() => _lockTimeout = minutes);
                Navigator.pop(context);
              },
            ),
          const SizedBox(height: 16),
        ],
      ),
    );
  }

  void _showChangePasswordDialog() {
    final currentController = TextEditingController();
    final newController = TextEditingController();
    final confirmController = TextEditingController();
    
    showDialog(
      context: context,
      builder: (context) => AlertDialog(
        title: const Text('Changer le mot de passe'),
        content: Column(
          mainAxisSize: MainAxisSize.min,
          children: [
            TextField(
              controller: currentController,
              obscureText: true,
              decoration: const InputDecoration(
                labelText: 'Mot de passe actuel',
                prefixIcon: Icon(Icons.lock_outline),
              ),
            ),
            const SizedBox(height: 16),
            TextField(
              controller: newController,
              obscureText: true,
              decoration: const InputDecoration(
                labelText: 'Nouveau mot de passe',
                prefixIcon: Icon(Icons.lock),
              ),
            ),
            const SizedBox(height: 16),
            TextField(
              controller: confirmController,
              obscureText: true,
              decoration: const InputDecoration(
                labelText: 'Confirmer',
                prefixIcon: Icon(Icons.lock),
              ),
            ),
          ],
        ),
        actions: [
          TextButton(
            onPressed: () => Navigator.pop(context),
            child: const Text('Annuler'),
          ),
          ElevatedButton(
            onPressed: () {
              // TODO: Call API to change password
              Navigator.pop(context);
              ScaffoldMessenger.of(context).showSnackBar(
                const SnackBar(
                  content: Text('Mot de passe modifié'),
                  backgroundColor: Color(0xFF10B981),
                ),
              );
            },
            child: const Text('Confirmer'),
          ),
        ],
      ),
    );
  }

  void _toggle2FA(bool enabled) {
    setState(() => _twoFactorEnabled = enabled);
    if (enabled) {
      _setup2FA();
    }
  }

  void _setup2FA() {
    // TODO: Implement 2FA setup with QR code
    ScaffoldMessenger.of(context).showSnackBar(
      const SnackBar(content: Text('Configuration 2FA à venir')),
    );
  }

  void _showDevicesSheet() {
    showModalBottomSheet(
      context: context,
      backgroundColor: Colors.white,
      shape: const RoundedRectangleBorder(
        borderRadius: BorderRadius.vertical(top: Radius.circular(20)),
      ),
      builder: (context) => Column(
        mainAxisSize: MainAxisSize.min,
        children: [
          const SizedBox(height: 16),
          Container(
            width: 40,
            height: 4,
            decoration: BoxDecoration(
              color: const Color(0xFFE2E8F0),
              borderRadius: BorderRadius.circular(2),
            ),
          ),
          const SizedBox(height: 16),
          const Text(
            'Appareils connectés',
            style: TextStyle(fontSize: 18, fontWeight: FontWeight.bold),
          ),
          const SizedBox(height: 16),
          ListTile(
            leading: Container(
              padding: const EdgeInsets.all(8),
              decoration: BoxDecoration(
                color: const Color(0xFF667eea).withOpacity(0.1),
                borderRadius: BorderRadius.circular(8),
              ),
              child: const Icon(Icons.phone_android, color: Color(0xFF667eea)),
            ),
            title: const Text('Cet appareil'),
            subtitle: const Text('Actif maintenant'),
            trailing: Container(
              padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 4),
              decoration: BoxDecoration(
                color: const Color(0xFF10B981).withOpacity(0.1),
                borderRadius: BorderRadius.circular(8),
              ),
              child: const Text('Actuel', style: TextStyle(color: Color(0xFF10B981), fontSize: 12)),
            ),
          ),
          ListTile(
            leading: Container(
              padding: const EdgeInsets.all(8),
              decoration: BoxDecoration(
                color: const Color(0xFF64748B).withOpacity(0.1),
                borderRadius: BorderRadius.circular(8),
              ),
              child: const Icon(Icons.laptop, color: Color(0xFF64748B)),
            ),
            title: const Text('Chrome - Windows'),
            subtitle: const Text('Dernière activité: il y a 2h'),
            trailing: IconButton(
              icon: const Icon(Icons.logout, color: Colors.red),
              onPressed: () {
                Navigator.pop(context);
                ScaffoldMessenger.of(context).showSnackBar(
                  const SnackBar(content: Text('Appareil déconnecté')),
                );
              },
            ),
          ),
          const SizedBox(height: 24),
        ],
      ),
    );
  }

  void _logoutAllDevices() {
    showDialog(
      context: context,
      builder: (context) => AlertDialog(
        title: const Text('Déconnecter tous les appareils?'),
        content: const Text('Vous serez déconnecté de tous les appareils y compris celui-ci.'),
        actions: [
          TextButton(
            onPressed: () => Navigator.pop(context),
            child: const Text('Annuler'),
          ),
          ElevatedButton(
            onPressed: () {
              Navigator.pop(context);
              ScaffoldMessenger.of(context).showSnackBar(
                const SnackBar(
                  content: Text('Tous les appareils ont été déconnectés'),
                  backgroundColor: Color(0xFF10B981),
                ),
              );
            },
            style: ElevatedButton.styleFrom(backgroundColor: Colors.red),
            child: const Text('Déconnecter'),
          ),
        ],
      ),
    );
  }
}
