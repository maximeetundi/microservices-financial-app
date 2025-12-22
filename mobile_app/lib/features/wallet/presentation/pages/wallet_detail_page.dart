import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';
import 'package:flutter_bloc/flutter_bloc.dart';
import '../bloc/wallet_bloc.dart';

class WalletDetailPage extends StatefulWidget {
  final String walletId;
  
  const WalletDetailPage({super.key, required this.walletId});

  @override
  State<WalletDetailPage> createState() => _WalletDetailPageState();
}

class _WalletDetailPageState extends State<WalletDetailPage> {
  @override
  void initState() {
    super.initState();
    context.read<WalletBloc>().add(LoadWalletTransactionsEvent(walletId: widget.walletId));
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: const Text('Détails du portefeuille'),
        leading: IconButton(
          icon: const Icon(Icons.arrow_back),
          onPressed: () => context.pop(),
        ),
        actions: [
          IconButton(icon: const Icon(Icons.more_vert), onPressed: () {}),
        ],
      ),
      body: BlocBuilder<WalletBloc, WalletState>(
        builder: (context, state) {
          if (state is WalletLoadingState) {
            return const Center(child: CircularProgressIndicator());
          }
          
          if (state is WalletErrorState) {
            return Center(child: Text('Erreur: ${state.message}'));
          }
          
          if (state is WalletLoadedState) {
            final wallet = state.wallets.firstWhere(
              (w) => w.id == widget.walletId,
              orElse: () => state.wallets.first,
            );
            
            return SingleChildScrollView(
              padding: const EdgeInsets.all(16),
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  // Balance Card
                  Container(
                    width: double.infinity,
                    padding: const EdgeInsets.all(24),
                    decoration: BoxDecoration(
                      gradient: LinearGradient(
                        colors: [Colors.blue[700]!, Colors.blue[500]!],
                        begin: Alignment.topLeft,
                        end: Alignment.bottomRight,
                      ),
                      borderRadius: BorderRadius.circular(20),
                    ),
                    child: Column(
                      crossAxisAlignment: CrossAxisAlignment.start,
                      children: [
                        Text(
                          wallet.currency,
                          style: const TextStyle(color: Colors.white70, fontSize: 16),
                        ),
                        const SizedBox(height: 8),
                        Text(
                          '${wallet.balance.toStringAsFixed(2)} ${wallet.currency}',
                          style: const TextStyle(
                            color: Colors.white,
                            fontSize: 32,
                            fontWeight: FontWeight.bold,
                          ),
                        ),
                      ],
                    ),
                  ),
                  const SizedBox(height: 24),
                  
                  // Actions
                  Row(
                    mainAxisAlignment: MainAxisAlignment.spaceEvenly,
                    children: [
                      _ActionButton(
                        icon: Icons.arrow_upward,
                        label: 'Envoyer',
                        onTap: () => context.push('/transfer'),
                      ),
                      _ActionButton(
                        icon: Icons.arrow_downward,
                        label: 'Recevoir',
                        onTap: () => _showReceiveDialog(context),
                      ),
                      _ActionButton(
                        icon: Icons.swap_horiz,
                        label: 'Échanger',
                        onTap: () => context.push('/exchange'),
                      ),
                    ],
                  ),
                  const SizedBox(height: 32),
                  
                  // Transactions
                  const Text(
                    'Transactions récentes',
                    style: TextStyle(fontSize: 18, fontWeight: FontWeight.bold),
                  ),
                  const SizedBox(height: 16),
                  
                  if (state.recentTransactions.isEmpty)
                    const Center(
                      child: Padding(
                        padding: EdgeInsets.all(32),
                        child: Text('Aucune transaction', style: TextStyle(color: Colors.grey)),
                      ),
                    )
                  else
                    ListView.builder(
                      shrinkWrap: true,
                      physics: const NeverScrollableScrollPhysics(),
                      itemCount: state.recentTransactions.length,
                      itemBuilder: (context, index) {
                        final tx = state.recentTransactions[index];
                        return ListTile(
                          leading: CircleAvatar(
                            backgroundColor: tx.isIncoming ? Colors.green[100] : Colors.red[100],
                            child: Icon(
                              tx.isIncoming ? Icons.arrow_downward : Icons.arrow_upward,
                              color: tx.isIncoming ? Colors.green : Colors.red,
                            ),
                          ),
                          title: Text(tx.memo ?? tx.type.name),
                          subtitle: Text(tx.createdAt.toString().substring(0, 16)),
                          trailing: Text(
                            '${tx.isIncoming ? '+' : '-'}${tx.amount.toStringAsFixed(2)}',
                            style: TextStyle(
                              color: tx.isIncoming ? Colors.green : Colors.red,
                              fontWeight: FontWeight.bold,
                            ),
                          ),
                        );
                      },
                    ),
                ],
              ),
            );
          }
          
          return const SizedBox();
        },
      ),
    );
  }

  void _showReceiveDialog(BuildContext context) {
    showDialog(
      context: context,
      builder: (context) => AlertDialog(
        title: const Text('Recevoir des fonds'),
        content: Column(
          mainAxisSize: MainAxisSize.min,
          children: [
            const Icon(Icons.qr_code, size: 150),
            const SizedBox(height: 16),
            Text('ID: ${widget.walletId}'),
          ],
        ),
        actions: [
          TextButton(
            onPressed: () => Navigator.pop(context),
            child: const Text('Fermer'),
          ),
        ],
      ),
    );
  }
}

class _ActionButton extends StatelessWidget {
  final IconData icon;
  final String label;
  final VoidCallback onTap;

  const _ActionButton({required this.icon, required this.label, required this.onTap});

  @override
  Widget build(BuildContext context) {
    return InkWell(
      onTap: onTap,
      borderRadius: BorderRadius.circular(12),
      child: Container(
        padding: const EdgeInsets.symmetric(horizontal: 24, vertical: 16),
        child: Column(
          children: [
            Container(
              padding: const EdgeInsets.all(12),
              decoration: BoxDecoration(
                color: Theme.of(context).primaryColor.withOpacity(0.1),
                borderRadius: BorderRadius.circular(12),
              ),
              child: Icon(icon, color: Theme.of(context).primaryColor),
            ),
            const SizedBox(height: 8),
            Text(label, style: const TextStyle(fontSize: 12)),
          ],
        ),
      ),
    );
  }
}
