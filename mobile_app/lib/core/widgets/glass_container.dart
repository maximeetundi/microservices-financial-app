import 'dart:ui';
import 'package:flutter/material.dart';

/// Glass container matching web frontend glassmorphism design exactly
/// Features: backdrop blur, border with opacity, optional top edge highlight
class GlassContainer extends StatelessWidget {
  final Widget child;
  final double? width;
  final double? height;
  final EdgeInsetsGeometry? padding;
  final EdgeInsetsGeometry? margin;
  final double borderRadius;
  final double blur;
  final Color? color;
  final Color? borderColor;
  final Gradient? gradient;
  final bool showTopHighlight;
  final List<BoxShadow>? boxShadow;

  const GlassContainer({
    super.key,
    required this.child,
    this.width,
    this.height,
    this.padding,
    this.margin,
    this.borderRadius = 24.0,
    this.blur = 20.0,
    this.color,
    this.borderColor,
    this.gradient,
    this.showTopHighlight = true,
    this.boxShadow,
  });

  @override
  Widget build(BuildContext context) {
    final isDark = Theme.of(context).brightness == Brightness.dark;
    
    // Default styles matching web design exactly
    // Light: rgba(255,255,255,0.6) - Dark: rgba(255,255,255,0.05)
    final defaultColor = isDark 
        ? Colors.white.withOpacity(0.05) 
        : Colors.white.withOpacity(0.6);
    
    // Light: rgba(255,255,255,0.4) - Dark: rgba(255,255,255,0.1)
    final defaultBorderColor = isDark
        ? Colors.white.withOpacity(0.1)
        : Colors.white.withOpacity(0.4);
    
    // Default shadow matching web box-shadow
    final defaultShadow = [
      BoxShadow(
        color: isDark 
            ? Colors.black.withOpacity(0.3) 
            : Colors.black.withOpacity(0.03),
        blurRadius: 20,
        offset: const Offset(0, 4),
      ),
    ];

    return Container(
      width: width,
      height: height,
      margin: margin,
      decoration: BoxDecoration(
        boxShadow: boxShadow ?? defaultShadow,
        borderRadius: BorderRadius.circular(borderRadius),
      ),
      child: ClipRRect(
        borderRadius: BorderRadius.circular(borderRadius),
        child: BackdropFilter(
          filter: ImageFilter.blur(sigmaX: blur, sigmaY: blur),
          child: Stack(
            children: [
              // Main container
              Container(
                padding: padding,
                decoration: BoxDecoration(
                  color: gradient == null ? (color ?? defaultColor) : null,
                  borderRadius: BorderRadius.circular(borderRadius),
                  border: Border.all(
                    color: borderColor ?? defaultBorderColor,
                    width: 1,
                  ),
                  gradient: gradient,
                ),
                child: child,
              ),
              // Top edge highlight (matching web ::before pseudo-element)
              if (showTopHighlight)
                Positioned(
                  top: 0,
                  left: 0,
                  right: 0,
                  child: Container(
                    height: 1,
                    decoration: BoxDecoration(
                      gradient: LinearGradient(
                        colors: [
                          Colors.transparent,
                          Colors.white.withOpacity(0.2),
                          Colors.transparent,
                        ],
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
