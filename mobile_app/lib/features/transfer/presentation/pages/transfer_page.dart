import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';
import 'package:go_router/go_router.dart';
import 'package:qr_code_scanner/qr_code_scanner.dart';

import '../../../../core/theme/app_theme.dart';
import '../../../../core/widgets/custom_button.dart';
import '../../../../core/widgets/custom_text_field.dart';
import '../../../../core/widgets/loading_widget.dart';
import '../bloc/transfer_bloc.dart';
import '../widgets/transfer_type_card.dart';
import '../widgets/recent_transfers_list.dart';
import '../widgets/contact_selector.dart';
import '../widgets/transfer_confirmation_sheet.dart';

class TransferPage extends StatefulWidget {
  const TransferPage({Key? key}) : super(key: key);

  @override
  State<TransferPage> createState() => _TransferPageState();
}

class _TransferPageState extends State<TransferPage>
    with SingleTickerProviderStateMixin {
  late TabController _tabController;
  
  final _formKey = GlobalKey<FormState>();
  final _amountController = TextEditingController();
  final _recipientController = TextEditingController();
  final _memoController = TextEditingController();
  
  String _selectedTransferType = 'crypto';
  String? _selectedWallet;
  String? _selectedRecipient;

  @override
  void initState() {
    super.initState();
    _tabController = TabController(length: 2, vsync: this);
    _loadTransferData();
  }

  void _loadTransferData() {
    context.read<TransferBloc>().add(LoadTransferDataEvent());
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: const Text('Send & Transfer'),
        centerTitle: true,
        bottom: TabBar(
          controller: _tabController,
          tabs: const [
            Tab(text: 'Send Money'),
            Tab(text: 'Recent'),
          ],
        ),
      ),
      body: BlocListener<TransferBloc, TransferState>(
        listener: (context, state) {
          if (state is TransferSuccessState) {
            _showSuccessDialog(state.transfer);
          } else if (state is TransferErrorState) {
            _showErrorSnackBar(state.message);
          }
        },
        child: TabBarView(
          controller: _tabController,
          children: [
            _buildSendMoneyTab(),
            _buildRecentTransfersTab(),
          ],
        ),
      ),
    );
  }

  Widget _buildSendMoneyTab() {
    return SingleChildScrollView(
      padding: const EdgeInsets.all(16),
      child: Form(
        key: _formKey,
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.stretch,
          children: [
            // Transfer Type Selection
            const Text(
              'Transfer Type',
              style: TextStyle(
                fontSize: 18,
                fontWeight: FontWeight.w600,
              ),
            ),
            const SizedBox(height: 12),
            
            Row(
              children: [
                Expanded(
                  child: TransferTypeCard(
                    icon: Icons.currency_bitcoin,
                    title: 'Crypto',
                    description: 'Send cryptocurrencies\nto any wallet address',
                    color: AppTheme.bitcoinColor,
                    isSelected: _selectedTransferType == 'crypto',
                    onTap: () => _setTransferType('crypto'),
                  ),
                ),
                const SizedBox(width: 12),
                Expanded(
                  child: TransferTypeCard(
                    icon: Icons.account_balance,
                    title: 'Bank',
                    description: 'SEPA/SWIFT transfers\nto bank accounts',
                    color: AppTheme.secondaryColor,
                    isSelected: _selectedTransferType == 'fiat',
                    onTap: () => _setTransferType('fiat'),
                  ),
                ),
                const SizedBox(width: 12),
                Expanded(
                  child: TransferTypeCard(
                    icon: Icons.flash_on,
                    title: 'Instant',
                    description: 'Free transfers between\nCrypto Bank users',
                    color: AppTheme.primaryColor,
                    isSelected: _selectedTransferType == 'instant',
                    onTap: () => _setTransferType('instant'),
                  ),
                ),
              ],
            ),

            const SizedBox(height: 32),

            // From Wallet Selection
            BlocBuilder<TransferBloc, TransferState>(
              builder: (context, state) {
                if (state is TransferLoadedState) {
                  return _buildFromWalletSection(state.wallets);
                }
                return const SizedBox.shrink();
              },
            ),

            const SizedBox(height: 24),

            // Recipient Selection
            _buildRecipientSection(),

            const SizedBox(height: 24),

            // Amount Input
            CustomTextField(
              controller: _amountController,
              label: 'Amount',
              hint: '0.00',
              keyboardType: const TextInputType.numberWithOptions(decimal: true),
              prefixIcon: Icons.attach_money,
              validator: _validateAmount,
            ),

            const SizedBox(height: 16),

            // Memo/Notes (optional)
            CustomTextField(
              controller: _memoController,
              label: 'Notes (Optional)',
              hint: 'Payment reference or notes',
              maxLines: 3,
              prefixIcon: Icons.note_outlined,
            ),

            const SizedBox(height: 32),

            // Fee Information
            _buildFeeInformation(),

            const SizedBox(height: 32),

            // Send Button
            BlocBuilder<TransferBloc, TransferState>(
              builder: (context, state) {
                final isLoading = state is TransferLoadingState;
                
                return CustomButton(
                  text: isLoading ? 'Processing...' : 'Send Transfer',
                  onPressed: isLoading ? null : _handleSendTransfer,
                  isLoading: isLoading,
                  icon: Icons.send,
                );
              },
            ),
          ],
        ),
      ),
    );
  }

  Widget _buildRecentTransfersTab() {
    return BlocBuilder<TransferBloc, TransferState>(
      builder: (context, state) {
        if (state is TransferLoadingState) {
          return const LoadingWidget();
        } else if (state is TransferLoadedState) {
          return RecentTransfersList(
            transfers: state.recentTransfers,
            onTransferTap: (transfer) {
              context.push('/more/transfer/${transfer.id}');
            },
          );
        } else if (state is TransferErrorState) {
          return Center(
            child: Column(
              mainAxisAlignment: MainAxisAlignment.center,
              children: [
                Icon(
                  Icons.error_outline,
                  size: 64,
                  color: AppTheme.errorColor,
                ),
                const SizedBox(height: 16),
                Text(
                  'Failed to load transfers',
                  style: Theme.of(context).textTheme.headlineSmall,
                ),
                const SizedBox(height: 8),
                Text(
                  state.message,
                  textAlign: TextAlign.center,
                ),
                const SizedBox(height: 24),
                CustomButton(
                  text: 'Try Again',
                  onPressed: _loadTransferData,
                ),
              ],
            ),
          );
        }
        return const SizedBox.shrink();
      },
    );
  }

  Widget _buildFromWalletSection(List<dynamic> wallets) {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        const Text(
          'From Wallet',
          style: TextStyle(
            fontSize: 16,
            fontWeight: FontWeight.w600,
          ),
        ),
        const SizedBox(height: 12),
        Container(
          decoration: BoxDecoration(
            border: Border.all(color: Colors.grey.shade300),
            borderRadius: BorderRadius.circular(12),
          ),
          child: DropdownButtonFormField<String>(
            value: _selectedWallet,
            decoration: const InputDecoration(
              contentPadding: EdgeInsets.symmetric(horizontal: 16, vertical: 12),
              border: InputBorder.none,
              hintText: 'Select wallet',
            ),
            items: wallets.map<DropdownMenuItem<String>>((wallet) {
              return DropdownMenuItem<String>(
                value: wallet.id,
                child: Row(
                  children: [
                    Icon(
                      _getCurrencyIcon(wallet.currency),
                      size: 20,
                      color: AppTheme.primaryColor,
                    ),
                    const SizedBox(width: 12),
                    Expanded(
                      child: Column(
                        crossAxisAlignment: CrossAxisAlignment.start,
                        children: [
                          Text(
                            '${wallet.currency} Wallet',
                            style: const TextStyle(fontWeight: FontWeight.w600),
                          ),
                          Text(
                            '${wallet.balance.toStringAsFixed(8)} ${wallet.currency}',
                            style: TextStyle(
                              fontSize: 12,
                              color: Colors.grey.shade600,
                            ),
                          ),
                        ],
                      ),
                    ),
                  ],
                ),
              );
            }).toList(),
            onChanged: (value) {
              setState(() {
                _selectedWallet = value;
              });
            },
            validator: (value) {
              if (value == null || value.isEmpty) {
                return 'Please select a wallet';
              }
              return null;
            },
          ),
        ),
      ],
    );
  }

  Widget _buildRecipientSection() {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Row(
          mainAxisAlignment: MainAxisAlignment.spaceBetween,
          children: [
            const Text(
              'Send To',
              style: TextStyle(
                fontSize: 16,
                fontWeight: FontWeight.w600,
              ),
            ),
            Row(
              children: [
                IconButton(
                  icon: const Icon(Icons.qr_code_scanner),
                  onPressed: _scanQRCode,
                  tooltip: 'Scan QR Code',
                ),
                IconButton(
                  icon: const Icon(Icons.contacts),
                  onPressed: _selectFromContacts,
                  tooltip: 'Select Contact',
                ),
              ],
            ),
          ],
        ),
        const SizedBox(height: 12),
        
        if (_selectedTransferType == 'crypto')
          CustomTextField(
            controller: _recipientController,
            label: 'Wallet Address',
            hint: 'Enter recipient wallet address',
            prefixIcon: Icons.account_balance_wallet,
            validator: _validateRecipient,
          )
        else if (_selectedTransferType == 'fiat')
          Column(
            children: [
              CustomTextField(
                controller: _recipientController,
                label: 'IBAN / Account Number',
                hint: 'Enter bank account details',
                prefixIcon: Icons.account_balance,
                validator: _validateRecipient,
              ),
              const SizedBox(height: 16),
              Row(
                children: [
                  Expanded(
                    child: CustomTextField(
                      label: 'Bank Code / SWIFT',
                      hint: 'SWIFT/BIC code',
                      prefixIcon: Icons.business,
                    ),
                  ),
                  const SizedBox(width: 12),
                  Expanded(
                    child: CustomTextField(
                      label: 'Recipient Name',
                      hint: 'Full name',
                      prefixIcon: Icons.person,
                    ),
                  ),
                ],
              ),
            ],
          )
        else // instant transfer
          CustomTextField(
            controller: _recipientController,
            label: 'Email Address',
            hint: 'user@email.com',
            keyboardType: TextInputType.emailAddress,
            prefixIcon: Icons.email,
            validator: _validateRecipient,
          ),
      ],
    );
  }

  Widget _buildFeeInformation() {
    return Container(
      padding: const EdgeInsets.all(16),
      decoration: BoxDecoration(
        color: Colors.grey.shade50,
        borderRadius: BorderRadius.circular(12),
        border: Border.all(color: Colors.grey.shade200),
      ),
      child: Column(
        children: [
          Row(
            mainAxisAlignment: MainAxisAlignment.spaceBetween,
            children: [
              const Text('Transfer Amount:'),
              Text(
                '${_amountController.text.isEmpty ? '0.00' : _amountController.text} ${_getSelectedCurrency()}',
                style: const TextStyle(fontWeight: FontWeight.w600),
              ),
            ],
          ),
          const SizedBox(height: 8),
          Row(
            mainAxisAlignment: MainAxisAlignment.spaceBetween,
            children: [
              const Text('Network Fee:'),
              Text(
                _getEstimatedFee(),
                style: const TextStyle(fontWeight: FontWeight.w600),
              ),
            ],
          ),
          const Divider(height: 24),
          Row(
            mainAxisAlignment: MainAxisAlignment.spaceBetween,
            children: [
              const Text(
                'Total:',
                style: TextStyle(
                  fontSize: 16,
                  fontWeight: FontWeight.bold,
                ),
              ),
              Text(
                _getTotalAmount(),
                style: const TextStyle(
                  fontSize: 16,
                  fontWeight: FontWeight.bold,
                ),
              ),
            ],
          ),
        ],
      ),
    );
  }

  void _setTransferType(String type) {
    setState(() {
      _selectedTransferType = type;
      _recipientController.clear();
    });
  }

  String? _validateAmount(String? value) {
    if (value == null || value.isEmpty) {
      return 'Please enter an amount';
    }
    final amount = double.tryParse(value);
    if (amount == null || amount <= 0) {
      return 'Please enter a valid amount';
    }
    return null;
  }

  String? _validateRecipient(String? value) {
    if (value == null || value.isEmpty) {
      return 'This field is required';
    }
    
    if (_selectedTransferType == 'crypto') {
      // Basic crypto address validation
      if (value.length < 26) {
        return 'Invalid wallet address';
      }
    } else if (_selectedTransferType == 'instant') {
      // Email validation
      if (!RegExp(r'^[\w-\.]+@([\w-]+\.)+[\w-]{2,4}$').hasMatch(value)) {
        return 'Please enter a valid email address';
      }
    }
    
    return null;
  }

  IconData _getCurrencyIcon(String currency) {
    switch (currency.toUpperCase()) {
      case 'BTC': return Icons.currency_bitcoin;
      case 'ETH': return Icons.diamond;
      case 'USD': return Icons.attach_money;
      case 'EUR': return Icons.euro;
      default: return Icons.monetization_on;
    }
  }

  String _getSelectedCurrency() {
    // This should come from the selected wallet
    return 'USD'; // Placeholder
  }

  String _getEstimatedFee() {
    final amount = double.tryParse(_amountController.text) ?? 0;
    switch (_selectedTransferType) {
      case 'crypto':
        return '${(amount * 0.0025).toStringAsFixed(8)} ${_getSelectedCurrency()}';
      case 'fiat':
        return '\$${(amount * 0.001).clamp(2.5, 50.0).toStringAsFixed(2)}';
      case 'instant':
        return 'Free';
      default:
        return '0.00';
    }
  }

  String _getTotalAmount() {
    final amount = double.tryParse(_amountController.text) ?? 0;
    final fee = _selectedTransferType == 'instant' ? 0.0 : 
                _selectedTransferType == 'crypto' ? amount * 0.0025 :
                (amount * 0.001).clamp(2.5, 50.0);
    
    return '${(amount + fee).toStringAsFixed(8)} ${_getSelectedCurrency()}';
  }

  Future<void> _scanQRCode() async {
    // Implement QR code scanning
    showModalBottomSheet(
      context: context,
      builder: (context) => Container(
        height: 400,
        child: Column(
          children: [
            AppBar(
              title: const Text('Scan QR Code'),
              automaticallyImplyLeading: false,
              actions: [
                IconButton(
                  icon: const Icon(Icons.close),
                  onPressed: () => Navigator.pop(context),
                ),
              ],
            ),
            Expanded(
              child: QRView(
                key: GlobalKey(debugLabel: 'QR'),
                onQRViewCreated: (controller) {
                  controller.scannedDataStream.listen((scanData) {
                    Navigator.pop(context);
                    setState(() {
                      _recipientController.text = scanData.code ?? '';
                    });
                  });
                },
              ),
            ),
          ],
        ),
      ),
    );
  }

  void _selectFromContacts() {
    showModalBottomSheet(
      context: context,
      isScrollControlled: true,
      builder: (context) => ContactSelector(
        onContactSelected: (contact) {
          Navigator.pop(context);
          setState(() {
            _selectedRecipient = contact.id;
            _recipientController.text = contact.address;
          });
        },
      ),
    );
  }

  void _handleSendTransfer() {
    if (_formKey.currentState!.validate()) {
      final amount = double.parse(_amountController.text);
      
      showModalBottomSheet(
        context: context,
        isScrollControlled: true,
        builder: (context) => TransferConfirmationSheet(
          transferType: _selectedTransferType,
          amount: amount,
          recipient: _recipientController.text,
          currency: _getSelectedCurrency(),
          fee: _getEstimatedFee(),
          onConfirm: () {
            Navigator.pop(context);
            context.read<TransferBloc>().add(
              SendTransferEvent(
                type: _selectedTransferType,
                fromWallet: _selectedWallet!,
                recipient: _recipientController.text,
                amount: amount,
                memo: _memoController.text.isEmpty ? null : _memoController.text,
              ),
            );
          },
        ),
      );
    }
  }

  void _showSuccessDialog(dynamic transfer) {
    showDialog(
      context: context,
      builder: (context) => AlertDialog(
        icon: const Icon(
          Icons.check_circle,
          color: AppTheme.successColor,
          size: 48,
        ),
        title: const Text('Transfer Sent'),
        content: Column(
          mainAxisSize: MainAxisSize.min,
          children: [
            const Text('Your transfer has been sent successfully!'),
            const SizedBox(height: 16),
            Container(
              padding: const EdgeInsets.all(12),
              decoration: BoxDecoration(
                color: Colors.grey.shade100,
                borderRadius: BorderRadius.circular(8),
              ),
              child: Column(
                children: [
                  Row(
                    mainAxisAlignment: MainAxisAlignment.spaceBetween,
                    children: [
                      const Text('Transaction ID:'),
                      Text(
                        transfer.id.substring(0, 8),
                        style: const TextStyle(fontFamily: 'monospace'),
                      ),
                    ],
                  ),
                  const SizedBox(height: 8),
                  Row(
                    mainAxisAlignment: MainAxisAlignment.spaceBetween,
                    children: [
                      const Text('Status:'),
                      Text(transfer.status),
                    ],
                  ),
                ],
              ),
            ),
          ],
        ),
        actions: [
          TextButton(
            onPressed: () => Navigator.pop(context),
            child: const Text('View Details'),
          ),
          ElevatedButton(
            onPressed: () => Navigator.pop(context),
            child: const Text('Done'),
          ),
        ],
      ),
    );
  }

  void _showErrorSnackBar(String message) {
    ScaffoldMessenger.of(context).showSnackBar(
      SnackBar(
        content: Text(message),
        backgroundColor: AppTheme.errorColor,
      ),
    );
  }

  @override
  void dispose() {
    _tabController.dispose();
    _amountController.dispose();
    _recipientController.dispose();
    _memoController.dispose();
    super.dispose();
  }
}