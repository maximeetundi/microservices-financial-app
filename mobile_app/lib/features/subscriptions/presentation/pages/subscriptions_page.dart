
import 'package:flutter/material.dart';
import '../../data/repositories/subscription_repository.dart';
import '../../data/models/subscription.dart';
import '../widgets/link_subscription_modal.dart';
import '../../../../core/di/injection_container.dart';

class SubscriptionsPage extends StatefulWidget {
  const SubscriptionsPage({Key? key}) : super(key: key);

  @override
  _SubscriptionsPageState createState() => _SubscriptionsPageState();
}

class _SubscriptionsPageState extends State<SubscriptionsPage> {
  final _repository = sl<SubscriptionRepository>();
  
  List<Subscription> _subscriptions = [];
  List<Map<String, dynamic>> _bills = [];
  bool _isLoading = true;

  @override
  void initState() {
    super.initState();
    _loadData();
  }

  Future<void> _loadData() async {
    setState(() => _isLoading = true);
    try {
      final subs = await _repository.getSubscriptions();
      final bills = await _repository.getPendingBills();
      if (mounted) {
        setState(() {
          _subscriptions = subs;
          _bills = bills;
        });
      }
    } catch (e) {
      debugPrint("Error loading subscriptions: $e");
    } finally {
      if (mounted) setState(() => _isLoading = false);
    }
  }

  void _showLinkModal() {
    showModalBottomSheet(
      context: context,
      isScrollControlled: true,
      builder: (_) => Padding(
        padding: EdgeInsets.only(bottom: MediaQuery.of(context).viewInsets.bottom),
        child: LinkSubscriptionModal(onSuccess: _loadData),
      ),
    );
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: Text("Mes Abonnements"),
        actions: [
          IconButton(
            icon: Icon(Icons.add),
            tooltip: "Lier un compte",
            onPressed: _showLinkModal,
          )
        ],
      ),
      body: _isLoading 
        ? Center(child: CircularProgressIndicator()) 
        : RefreshIndicator(
            onRefresh: _loadData,
            child: SingleChildScrollView(
              padding: EdgeInsets.all(16),
              physics: AlwaysScrollableScrollPhysics(),
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  // PENDING BILLS SECTION
                  if (_bills.isNotEmpty) ...[
                    Text("Factures en attente", style: TextStyle(fontSize: 18, fontWeight: FontWeight.bold, color: Colors.orange)),
                    SizedBox(height: 10),
                    ListView.builder(
                      shrinkWrap: true,
                      physics: NeverScrollableScrollPhysics(),
                      itemCount: _bills.length,
                      itemBuilder: (ctx, idx) {
                        final bill = _bills[idx];
                        return Card(
                          color: Colors.orange.shade50,
                          child: ListTile(
                            leading: Icon(Icons.receipt, color: Colors.orange),
                            title: Text(bill['service_name'] ?? 'Facture'),
                            subtitle: Text("Échéance: ${bill['due_date'] ?? 'N/A'}"),
                            trailing: ElevatedButton(
                              onPressed: () {}, // TODO: Implement Pay
                              child: Text("Payer ${bill['amount']} FCFA"),
                              style: ElevatedButton.styleFrom(backgroundColor: Colors.green),
                            ),
                          ),
                        );
                      }
                    ),
                    SizedBox(height: 24),
                  ],

                  // SUBSCRIPTIONS SECTION
                  Text("Mes Services & Écoles", style: TextStyle(fontSize: 18, fontWeight: FontWeight.bold)),
                  SizedBox(height: 10),
                  if (_subscriptions.isEmpty)
                    Center(child: Padding(
                      padding: const EdgeInsets.all(32.0),
                      child: Text("Aucun abonnement actif.\nCliquez sur + pour lier un compte.", textAlign: TextAlign.center, style: TextStyle(color: Colors.grey)),
                    ))
                  else
                    ListView.builder(
                      shrinkWrap: true,
                      physics: NeverScrollableScrollPhysics(),
                      itemCount: _subscriptions.length,
                      itemBuilder: (ctx, idx) {
                        final sub = _subscriptions[idx];
                        return Card(
                          elevation: 2,
                          child: ListTile(
                            leading: CircleAvatar(
                              child: Text(sub.enterpriseName.isNotEmpty ? sub.enterpriseName[0] : '?'),
                              backgroundColor: Colors.blue.shade100,
                            ),
                            title: Text(sub.enterpriseName),
                            subtitle: Column(
                              crossAxisAlignment: CrossAxisAlignment.start,
                              children: [
                                Text("${sub.serviceName} • ${sub.externalId}"),
                                if (sub.studentName != null)
                                  Text("Élève: ${sub.studentName} (${sub.className})", style: TextStyle(fontSize: 12, color: Colors.blueGrey)),
                              ],
                            ),
                            isThreeLine: true,
                            trailing: Icon(Icons.chevron_right),
                            onTap: () {
                              // Go to details
                            },
                          ),
                        );
                      },
                    ),
                ],
              ),
            ),
          ),
    );
  }
}
