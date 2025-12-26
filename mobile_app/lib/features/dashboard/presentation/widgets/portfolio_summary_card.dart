import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';
import 'package:google_fonts/google_fonts.dart';
import 'package:intl/intl.dart';

import '../../../wallet/presentation/bloc/wallet_bloc.dart';

/// Portfolio Summary Card - displays total balance from wallets
class PortfolioSummaryCard extends StatelessWidget {
  const PortfolioSummaryCard({super.key});

  String _formatCurrency(double amount) {
    final formatter = NumberFormat.currency(
      locale: 'fr_FR',
      symbol: '\$',
      decimalDigits: 2,
    );
    return formatter.format(amount);
  }

  @override
  Widget build(BuildContext context) {
    return BlocBuilder<WalletBloc, WalletState>(
      builder: (context, state) {
        double totalBalance = 0.0;
        double changePercent = 0.0;

        if (state is WalletsLoadedState) {
          // Calculate total balance from all wallets
          for (final wallet in state.wallets) {
            final balance = wallet['balance'];
            if (balance is num) {
              totalBalance += balance.toDouble();
            } else if (balance is String) {
              totalBalance += double.tryParse(balance) ?? 0.0;
            }
          }
          
          // Calculate change from transaction history if available
          // For now, we show 0% since we don't have historical data
          changePercent = 0.0;
        }

        return Container(
          width: double.infinity,
          padding: const EdgeInsets.all(24),
          decoration: BoxDecoration(
            gradient: LinearGradient(
              colors: [const Color(0xFF6366F1), const Color(0xFF8B5CF6)],
              begin: Alignment.topLeft,
              end: Alignment.bottomRight,
            ),
            borderRadius: BorderRadius.circular(20),
            boxShadow: [
              BoxShadow(
                color: const Color(0xFF6366F1).withOpacity(0.3),
                blurRadius: 20,
                offset: const Offset(0, 10),
              ),
            ],
          ),
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              Row(
                mainAxisAlignment: MainAxisAlignment.spaceBetween,
                children: [
                  Text(
                    'Solde total',
                    style: GoogleFonts.inter(
                      color: Colors.white70,
                      fontSize: 14,
                    ),
                  ),
                  if (state is WalletsLoadingState)
                    const SizedBox(
                      width: 16,
                      height: 16,
                      child: CircularProgressIndicator(
                        strokeWidth: 2,
                        color: Colors.white,
                      ),
                    ),
                ],
              ),
              const SizedBox(height: 8),
              Text(
                state is WalletsLoadedState
                    ? _formatCurrency(totalBalance)
                    : '---',
                style: GoogleFonts.inter(
                  color: Colors.white,
                  fontSize: 36,
                  fontWeight: FontWeight.bold,
                ),
              ),
              const SizedBox(height: 16),
              Row(
                children: [
                  Container(
                    padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 4),
                    decoration: BoxDecoration(
                      color: changePercent >= 0 
                          ? Colors.green.withOpacity(0.2)
                          : Colors.red.withOpacity(0.2),
                      borderRadius: BorderRadius.circular(8),
                    ),
                    child: Row(
                      children: [
                        Icon(
                          changePercent >= 0 ? Icons.arrow_upward : Icons.arrow_downward,
                          color: changePercent >= 0 ? Colors.green : Colors.red,
                          size: 16,
                        ),
                        Text(
                          '${changePercent >= 0 ? '+' : ''}${changePercent.toStringAsFixed(1)}%',
                          style: GoogleFonts.inter(
                            color: changePercent >= 0 ? Colors.green : Colors.red,
                            fontWeight: FontWeight.bold,
                          ),
                        ),
                      ],
                    ),
                  ),
                  const SizedBox(width: 8),
                  Text(
                    'Ce mois',
                    style: GoogleFonts.inter(
                      color: Colors.white70,
                      fontSize: 12,
                    ),
                  ),
                ],
              ),
            ],
          ),
        );
      },
    );
  }
}
