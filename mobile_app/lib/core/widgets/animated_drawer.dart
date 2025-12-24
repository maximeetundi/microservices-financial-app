import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';
import 'dart:ui';
import '../theme/app_theme.dart';

/// Modern animated drawer with glassmorphism effect
class AnimatedDrawer extends StatefulWidget {
  final Widget child;
  final VoidCallback? onMenuTap;

  const AnimatedDrawer({
    super.key,
    required this.child,
    this.onMenuTap,
  });

  @override
  State<AnimatedDrawer> createState() => AnimatedDrawerState();
}

class AnimatedDrawerState extends State<AnimatedDrawer>
    with SingleTickerProviderStateMixin {
  late AnimationController _animationController;
  late Animation<double> _scaleAnimation;
  late Animation<double> _slideAnimation;
  late Animation<double> _rotateAnimation;
  late Animation<double> _menuScaleAnimation;

  bool _isDrawerOpen = false;

  @override
  void initState() {
    super.initState();
    _animationController = AnimationController(
      vsync: this,
      duration: const Duration(milliseconds: 400),
    );

    _scaleAnimation = Tween<double>(begin: 1.0, end: 0.75).animate(
      CurvedAnimation(parent: _animationController, curve: Curves.easeOutCubic),
    );

    _slideAnimation = Tween<double>(begin: 0.0, end: 280.0).animate(
      CurvedAnimation(parent: _animationController, curve: Curves.easeOutCubic),
    );

    _rotateAnimation = Tween<double>(begin: 0.0, end: -0.1).animate(
      CurvedAnimation(parent: _animationController, curve: Curves.easeOutCubic),
    );

    _menuScaleAnimation = Tween<double>(begin: 0.5, end: 1.0).animate(
      CurvedAnimation(parent: _animationController, curve: Curves.easeOutBack),
    );
  }

  @override
  void dispose() {
    _animationController.dispose();
    super.dispose();
  }

  void toggleDrawer() {
    if (_isDrawerOpen) {
      _animationController.reverse();
    } else {
      _animationController.forward();
    }
    setState(() => _isDrawerOpen = !_isDrawerOpen);
  }

  void closeDrawer() {
    if (_isDrawerOpen) {
      _animationController.reverse();
      setState(() => _isDrawerOpen = false);
    }
  }

  @override
  Widget build(BuildContext context) {
    return Stack(
      children: [
        // Animated Drawer Background
        AnimatedBuilder(
          animation: _animationController,
          builder: (context, child) {
            return Container(
              decoration: const BoxDecoration(
                gradient: LinearGradient(
                  begin: Alignment.topLeft,
                  end: Alignment.bottomRight,
                  colors: [
                    AppTheme.darkTheme.colorScheme.background,
                    AppTheme.darkTheme.colorScheme.surface,
                  ],
                ),
              ),
              child: SafeArea(
                child: Padding(
                  padding: const EdgeInsets.only(left: 24, top: 80, bottom: 40),
                  child: Transform.scale(
                    scale: _menuScaleAnimation.value,
                    alignment: Alignment.centerLeft,
                    child: Opacity(
                      opacity: _animationController.value,
                      child: _buildDrawerContent(context),
                    ),
                  ),
                ),
              ),
            );
          },
        ),

        // Main Content
        AnimatedBuilder(
          animation: _animationController,
          builder: (context, child) {
            return Transform(
              transform: Matrix4.identity()
                ..translate(_slideAnimation.value, 0.0)
                ..scale(_scaleAnimation.value)
                ..rotateZ(_rotateAnimation.value),
              alignment: Alignment.centerLeft,
              child: GestureDetector(
                onTap: _isDrawerOpen ? closeDrawer : null,
                onHorizontalDragUpdate: (details) {
                  if (details.delta.dx > 8 && !_isDrawerOpen) {
                    toggleDrawer();
                  } else if (details.delta.dx < -8 && _isDrawerOpen) {
                    closeDrawer();
                  }
                },
                child: AbsorbPointer(
                  absorbing: _isDrawerOpen,
                  child: ClipRRect(
                    borderRadius: BorderRadius.circular(
                      _isDrawerOpen ? 30 : 0,
                    ),
                    child: Container(
                      decoration: BoxDecoration(
                        boxShadow: [
                          BoxShadow(
                            color: Colors.black.withOpacity(0.3),
                            blurRadius: 30,
                            offset: const Offset(-10, 0),
                          ),
                        ],
                      ),
                      child: widget.child,
                    ),
                  ),
                ),
              ),
            );
          },
        ),
      ],
    );
  }

  Widget _buildDrawerContent(BuildContext context) {
    return SingleChildScrollView(
      child: ConstrainedBox(
        constraints: BoxConstraints(
          minHeight: MediaQuery.of(context).size.height - 140,
        ),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          mainAxisSize: MainAxisSize.min,
          children: [
            // User Profile Section
            _buildUserProfile(),
            const SizedBox(height: 24),

            // Main Menu Items - TRANSFER FOCUSED
            _buildMenuItem(
              icon: Icons.home_rounded,
              title: 'Accueil',
              onTap: () => _navigateTo(context, '/dashboard'),
              isActive: true,
            ),
            _buildMenuItem(
              icon: Icons.send_rounded,
              title: 'Envoyer',
              onTap: () => _navigateTo(context, '/more/transfer'),
              badge: 'Rapide',
              badgeColor: const Color(0xFF10B981),
            ),
            _buildMenuItem(
              icon: Icons.qr_code_scanner_rounded,
              title: 'Scanner & Payer',
              onTap: () => _navigateTo(context, '/more/merchant/scan'),
            ),
            _buildMenuItem(
              icon: Icons.account_balance_wallet_rounded,
              title: 'Portefeuilles',
              onTap: () => _navigateTo(context, '/wallet'),
            ),
            _buildMenuItem(
              icon: Icons.credit_card_rounded,
              title: 'Cartes',
              onTap: () => _navigateTo(context, '/more/cards'),
            ),
            
            const Padding(
              padding: EdgeInsets.symmetric(vertical: 12),
              child: Divider(color: Colors.white24, indent: 0, endIndent: 40),
            ),
            
            // Secondary Items
            _buildMenuItem(
              icon: Icons.swap_horiz_rounded,
              title: 'Exchange',
              onTap: () => _navigateTo(context, '/exchange'),
              subtitle: 'Crypto & Fiat',
            ),
            _buildMenuItem(
              icon: Icons.pie_chart_rounded,
              title: 'Portfolio',
              onTap: () => _navigateTo(context, '/portfolio'),
            ),
            _buildMenuItem(
              icon: Icons.storefront_rounded,
              title: 'Marchand',
              onTap: () => _navigateTo(context, '/more/merchant'),
            ),

            const SizedBox(height: 20),
            const Padding(
              padding: EdgeInsets.symmetric(vertical: 12),
              child: Divider(color: Colors.white24, indent: 0, endIndent: 40),
            ),

            // Bottom Items
            _buildMenuItem(
              icon: Icons.support_agent_rounded,
              title: 'Support 24/7',
              onTap: () => _navigateTo(context, '/more/support'),
            ),
            _buildMenuItem(
              icon: Icons.settings_rounded,
              title: 'Paramètres',
              onTap: () => _navigateTo(context, '/more'),
            ),
          ],
        ),
      ),
    );
  }

  Widget _buildUserProfile() {
    return Row(
      children: [
        Container(
          width: 60,
          height: 60,
          decoration: BoxDecoration(
            gradient: AppTheme.primaryGradient,
            borderRadius: BorderRadius.circular(20),
            boxShadow: [
              BoxShadow(
                color: const Color(0xFF667eea).withOpacity(0.4),
                blurRadius: 15,
                offset: const Offset(0, 5),
              ),
            ],
          ),
          child: const Center(
            child: Text(
              '👤',
              style: TextStyle(fontSize: 28),
            ),
          ),
        ),
        const SizedBox(width: 16),
        Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            const Text(
              'Bonjour! 👋',
              style: TextStyle(
                color: Colors.white60,
                fontSize: 14,
              ),
            ),
            const SizedBox(height: 4),
            const Text(
              'Utilisateur',
              style: TextStyle(
                color: Colors.white,
                fontSize: 20,
                fontWeight: FontWeight.bold,
              ),
            ),
          ],
        ),
      ],
    );
  }

  Widget _buildMenuItem({
    required IconData icon,
    required String title,
    required VoidCallback onTap,
    String? subtitle,
    String? badge,
    Color? badgeColor,
    bool isActive = false,
  }) {
    return Padding(
      padding: const EdgeInsets.only(bottom: 8),
      child: Material(
        color: Colors.transparent,
        child: InkWell(
          onTap: onTap,
          borderRadius: BorderRadius.circular(12),
          child: Container(
            padding: const EdgeInsets.symmetric(vertical: 12),
            child: Row(
              children: [
                Container(
                  width: 42,
                  height: 42,
                  decoration: BoxDecoration(
                    color: isActive 
                        ? Colors.white.withOpacity(0.2)
                        : Colors.white.withOpacity(0.08),
                    borderRadius: BorderRadius.circular(12),
                  ),
                  child: Icon(
                    icon,
                    color: isActive ? Colors.white : Colors.white70,
                    size: 22,
                  ),
                ),
                const SizedBox(width: 16),
                Expanded(
                  child: Column(
                    crossAxisAlignment: CrossAxisAlignment.start,
                    children: [
                      Row(
                        children: [
                          Text(
                            title,
                            style: TextStyle(
                              color: isActive ? Colors.white : Colors.white70,
                              fontSize: 16,
                              fontWeight: isActive ? FontWeight.w600 : FontWeight.w500,
                            ),
                          ),
                          if (badge != null) ...[
                            const SizedBox(width: 8),
                            Container(
                              padding: const EdgeInsets.symmetric(
                                horizontal: 8,
                                vertical: 2,
                              ),
                              decoration: BoxDecoration(
                                color: badgeColor ?? const Color(0xFF667eea),
                                borderRadius: BorderRadius.circular(10),
                              ),
                              child: Text(
                                badge,
                                style: const TextStyle(
                                  color: Colors.white,
                                  fontSize: 10,
                                  fontWeight: FontWeight.bold,
                                ),
                              ),
                            ),
                          ],
                        ],
                      ),
                      if (subtitle != null) ...[
                        const SizedBox(height: 2),
                        Text(
                          subtitle,
                          style: const TextStyle(
                            color: Colors.white38,
                            fontSize: 12,
                          ),
                        ),
                      ],
                    ],
                  ),
                ),
              ],
            ),
          ),
        ),
      ),
    );
  }

  void _navigateTo(BuildContext context, String route) {
    closeDrawer();
    Future.delayed(const Duration(milliseconds: 200), () {
      if (context.mounted) {
        context.go(route);
      }
    });
  }
}
