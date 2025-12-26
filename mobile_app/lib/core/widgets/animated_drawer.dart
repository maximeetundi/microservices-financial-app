import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';
import 'package:go_router/go_router.dart';
import 'package:google_fonts/google_fonts.dart';

import '../theme/app_theme.dart';
import '../../features/auth/presentation/bloc/auth_bloc.dart';

/// Sidebar drawer matching web frontend design exactly
/// Simple slide drawer without excessive animations
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
  late Animation<double> _slideAnimation;

  bool _isDrawerOpen = false;

  @override
  void initState() {
    super.initState();
    _animationController = AnimationController(
      vsync: this,
      duration: const Duration(milliseconds: 250),
    );

    _slideAnimation = Tween<double>(begin: 0.0, end: 280.0).animate(
      CurvedAnimation(parent: _animationController, curve: Curves.easeOut),
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
    final isDark = Theme.of(context).brightness == Brightness.dark;
    
    return Stack(
      children: [
        // Drawer Background - matching web sidebar
        Container(
          width: 280,
          decoration: BoxDecoration(
            color: isDark ? const Color(0xFF0F172A) : const Color(0xFFF1F5F9),
            border: Border(
              right: BorderSide(
                color: isDark ? const Color(0xFF1E293B) : const Color(0xFFE2E8F0),
              ),
            ),
          ),
          child: SafeArea(
            child: Column(
              children: [
                // Logo Header - matching web
                _buildLogoHeader(isDark),
                
                // Navigation - matching web exactly
                Expanded(
                  child: SingleChildScrollView(
                    padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 24),
                    child: Column(
                      crossAxisAlignment: CrossAxisAlignment.start,
                      children: [
                        _buildNavItem('📊', 'Tableau de bord', '/dashboard', isDark),
                        _buildNavItem('👛', 'Portefeuilles', '/wallet', isDark),
                        _buildNavItem('💳', 'Mes Cartes', '/more/cards', isDark),
                        
                        _buildSectionTitle('Échange', isDark),
                        _buildNavItem('₿', 'Crypto', '/exchange', isDark),
                        _buildNavItem('💱', 'Fiat', '/exchange', isDark),
                        
                        _buildSectionTitle('Opérations', isDark),
                        _buildNavItem('💸', 'Virements', '/more/transfer', isDark),
                        _buildNavItem('🏪', 'Espace Marchand', '/more/merchant', isDark),
                        _buildNavItem('📷', 'Scanner / Payer', '/more/merchant/scan', isDark),
                        _buildNavItem('🔔', 'Notifications', '/dashboard/notifications', isDark),
                        _buildNavItem('⚙️', 'Paramètres', '/more', isDark),
                      ],
                    ),
                  ),
                ),
                
                // User Section - matching web
                _buildUserSection(isDark),
              ],
            ),
          ),
        ),

        // Main Content with slide
        AnimatedBuilder(
          animation: _animationController,
          builder: (context, child) {
            return Transform.translate(
              offset: Offset(_slideAnimation.value, 0),
              child: GestureDetector(
                onTap: _isDrawerOpen ? closeDrawer : null,
                onHorizontalDragUpdate: (details) {
                  if (details.delta.dx > 8 && !_isDrawerOpen) {
                    toggleDrawer();
                  } else if (details.delta.dx < -8 && _isDrawerOpen) {
                    closeDrawer();
                  }
                },
                child: Container(
                  decoration: BoxDecoration(
                    boxShadow: _isDrawerOpen ? [
                      BoxShadow(
                        color: Colors.black.withOpacity(0.1),
                        blurRadius: 20,
                        offset: const Offset(-5, 0),
                      ),
                    ] : null,
                  ),
                  child: AbsorbPointer(
                    absorbing: _isDrawerOpen,
                    child: widget.child,
                  ),
                ),
              ),
            );
          },
        ),

        // Overlay when drawer is open
        if (_isDrawerOpen)
          Positioned.fill(
            left: 280,
            child: GestureDetector(
              onTap: closeDrawer,
              child: Container(
                color: Colors.black.withOpacity(0.3),
              ),
            ),
          ),
      ],
    );
  }

  Widget _buildLogoHeader(bool isDark) {
    return Container(
      padding: const EdgeInsets.all(24),
      decoration: BoxDecoration(
        border: Border(
          bottom: BorderSide(
            color: isDark ? const Color(0xFF1E293B) : const Color(0xFFE2E8F0),
          ),
        ),
      ),
      child: Row(
        children: [
          Container(
            width: 40,
            height: 40,
            decoration: BoxDecoration(
              gradient: const LinearGradient(
                colors: [Color(0xFF6366F1), Color(0xFF8B5CF6)],
              ),
              borderRadius: BorderRadius.circular(12),
              boxShadow: [
                BoxShadow(
                  color: const Color(0xFF6366F1).withOpacity(0.3),
                  blurRadius: 12,
                  offset: const Offset(0, 4),
                ),
              ],
            ),
            child: const Center(
              child: Text('🏦', style: TextStyle(fontSize: 20)),
            ),
          ),
          const SizedBox(width: 12),
          Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              SelectionContainer.disabled(
                child: Text(
                  'Zekora',
                  style: GoogleFonts.inter(
                    fontSize: 18,
                    fontWeight: FontWeight.bold,
                    color: isDark ? Colors.white : const Color(0xFF1E293B),
                    decoration: TextDecoration.none,
                    decorationColor: Colors.transparent,
                  ),
                ),
              ),
              SelectionContainer.disabled(
                child: Text(
                  'Premium Banking',
                  style: GoogleFonts.inter(
                    fontSize: 12,
                    color: isDark ? const Color(0xFF94A3B8) : const Color(0xFF64748B),
                    decoration: TextDecoration.none,
                    decorationColor: Colors.transparent,
                  ),
                ),
              ),
            ],
          ),
        ],
      ),
    );
  }

  Widget _buildSectionTitle(String title, bool isDark) {
    return Padding(
      padding: const EdgeInsets.only(top: 24, bottom: 12, left: 16),
      child: SelectionContainer.disabled(
        child: Text(
          title.toUpperCase(),
          style: GoogleFonts.inter(
            fontSize: 11,
            fontWeight: FontWeight.w600,
            letterSpacing: 1.2,
            color: isDark ? const Color(0xFF64748B) : const Color(0xFF94A3B8),
            decoration: TextDecoration.none,
            decorationColor: Colors.transparent,
          ),
        ),
      ),
    );
  }

  Widget _buildNavItem(String emoji, String title, String route, bool isDark) {
    final currentRoute = GoRouterState.of(context).uri.toString();
    final isActive = currentRoute == route || currentRoute.startsWith(route);
    
    return Padding(
      padding: const EdgeInsets.only(bottom: 4),
      child: Material(
        color: Colors.transparent,
        child: InkWell(
          onTap: () => _navigateTo(context, route),
          borderRadius: BorderRadius.circular(12),
          child: Container(
            padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 12),
            decoration: BoxDecoration(
              color: isActive 
                  ? (isDark ? const Color(0xFF1E293B) : Colors.white)
                  : Colors.transparent,
              borderRadius: BorderRadius.circular(12),
            ),
            child: Row(
              children: [
                Text(emoji, style: const TextStyle(fontSize: 20)),
                const SizedBox(width: 12),
                Text(
                  title,
                  style: GoogleFonts.inter(
                    fontSize: 15,
                    fontWeight: isActive ? FontWeight.w600 : FontWeight.w500,
                    color: isActive 
                        ? (isDark ? Colors.white : const Color(0xFF1E293B))
                        : (isDark ? Colors.white70 : const Color(0xFF334155)),
                    decoration: TextDecoration.none,
                  ),
                ),
              ],
            ),
          ),
        ),
      ),
    );
  }

  Widget _buildUserSection(bool isDark) {
    return BlocBuilder<AuthBloc, AuthState>(
      builder: (context, state) {
        String userName = 'Utilisateur';
        String userEmail = 'user@zekora.com';
        String initials = 'U';
        
        if (state is AuthenticatedState) {
          userName = state.user.fullName;
          userEmail = state.user.email;
          initials = state.user.initials;
        }
        
        return Container(
          padding: const EdgeInsets.all(16),
          decoration: BoxDecoration(
            border: Border(
              top: BorderSide(
                color: isDark ? const Color(0xFF1E293B) : const Color(0xFFE2E8F0),
              ),
            ),
            color: isDark 
                ? const Color(0xFF0F172A).withOpacity(0.5)
                : const Color(0xFFF8FAFC),
          ),
          child: Row(
            children: [
              Container(
                width: 40,
                height: 40,
                decoration: BoxDecoration(
                  gradient: const LinearGradient(
                    colors: [Color(0xFF6366F1), Color(0xFF8B5CF6)],
                  ),
                  borderRadius: BorderRadius.circular(20),
                ),
                child: Center(
                  child: Text(
                    initials,
                    style: GoogleFonts.inter(
                      color: Colors.white,
                      fontWeight: FontWeight.bold,
                      fontSize: 14,
                    ),
                  ),
                ),
              ),
              const SizedBox(width: 12),
              Expanded(
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    SelectionContainer.disabled(
                      child: Text(
                        userName,
                        style: GoogleFonts.inter(
                          fontSize: 14,
                          fontWeight: FontWeight.w600,
                          color: isDark ? Colors.white : const Color(0xFF1E293B),
                          decoration: TextDecoration.none,
                          decorationColor: Colors.transparent,
                        ),
                        maxLines: 1,
                        overflow: TextOverflow.ellipsis,
                      ),
                    ),
                    SelectionContainer.disabled(
                      child: Text(
                        userEmail,
                        style: GoogleFonts.inter(
                          fontSize: 12,
                          color: isDark ? const Color(0xFF94A3B8) : const Color(0xFF64748B),
                          decoration: TextDecoration.none,
                          decorationColor: Colors.transparent,
                        ),
                        maxLines: 1,
                        overflow: TextOverflow.ellipsis,
                      ),
                    ),
                  ],
                ),
              ),
              IconButton(
                icon: Icon(
                  Icons.logout_rounded,
                  size: 20,
                  color: isDark ? const Color(0xFF64748B) : const Color(0xFF94A3B8),
                ),
                onPressed: () {
                  context.read<AuthBloc>().add(SignOutEvent());
                  context.go('/auth/login');
                },
              ),
            ],
          ),
        );
      },
    );
  }

  void _navigateTo(BuildContext context, String route) {
    closeDrawer();
    Future.delayed(const Duration(milliseconds: 150), () {
      if (context.mounted) {
        context.go(route);
      }
    });
  }
}
