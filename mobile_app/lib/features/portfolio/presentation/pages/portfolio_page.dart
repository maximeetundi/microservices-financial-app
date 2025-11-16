import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';
import 'package:fl_chart/fl_chart.dart';
import 'package:go_router/go_router.dart';

import '../../../../core/theme/app_theme.dart';
import '../../../../core/widgets/loading_widget.dart';
import '../bloc/portfolio_bloc.dart';
import '../widgets/portfolio_chart.dart';
import '../widgets/performance_metrics.dart';
import '../widgets/asset_allocation.dart';
import '../widgets/holdings_list.dart';
import '../widgets/portfolio_summary.dart';

class PortfolioPage extends StatefulWidget {
  const PortfolioPage({Key? key}) : super(key: key);

  @override
  State<PortfolioPage> createState() => _PortfolioPageState();
}

class _PortfolioPageState extends State<PortfolioPage>
    with SingleTickerProviderStateMixin {
  late TabController _tabController;
  String _selectedTimeframe = '1D';
  
  final List<String> _timeframes = ['1D', '1W', '1M', '3M', '1Y', 'ALL'];

  @override
  void initState() {
    super.initState();
    _tabController = TabController(length: 3, vsync: this);
    _loadPortfolioData();
  }

  void _loadPortfolioData() {
    context.read<PortfolioBloc>().add(LoadPortfolioEvent());
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: const Text('Portfolio'),
        centerTitle: true,
        actions: [
          IconButton(
            icon: const Icon(Icons.analytics_outlined),
            onPressed: () {
              // Navigate to detailed analytics
            },
          ),
        ],
      ),
      body: BlocBuilder<PortfolioBloc, PortfolioState>(
        builder: (context, state) {
          if (state is PortfolioLoadingState) {
            return const LoadingWidget();
          } else if (state is PortfolioLoadedState) {
            return _buildPortfolioContent(state);
          } else if (state is PortfolioErrorState) {
            return _buildErrorState(state.message);
          }
          return const SizedBox.shrink();
        },
      ),
    );
  }

  Widget _buildPortfolioContent(PortfolioLoadedState state) {
    return RefreshIndicator(
      onRefresh: () async {
        _loadPortfolioData();
      },
      child: CustomScrollView(
        slivers: [
          // Portfolio Summary
          SliverToBoxAdapter(
            child: PortfolioSummary(
              portfolio: state.portfolio,
            ),
          ),

          // Time Frame Selector
          SliverToBoxAdapter(
            child: Container(
              height: 60,
              margin: const EdgeInsets.symmetric(vertical: 16),
              child: ListView.builder(
                scrollDirection: Axis.horizontal,
                padding: const EdgeInsets.symmetric(horizontal: 16),
                itemCount: _timeframes.length,
                itemBuilder: (context, index) {
                  final timeframe = _timeframes[index];
                  final isSelected = timeframe == _selectedTimeframe;
                  
                  return GestureDetector(
                    onTap: () {
                      setState(() {
                        _selectedTimeframe = timeframe;
                      });
                      context.read<PortfolioBloc>().add(
                        ChangeTimeframeEvent(timeframe: timeframe),
                      );
                    },
                    child: Container(
                      margin: const EdgeInsets.only(right: 8),
                      padding: const EdgeInsets.symmetric(horizontal: 20, vertical: 8),
                      decoration: BoxDecoration(
                        color: isSelected
                            ? AppTheme.primaryColor
                            : Colors.grey.shade100,
                        borderRadius: BorderRadius.circular(20),
                      ),
                      child: Center(
                        child: Text(
                          timeframe,
                          style: TextStyle(
                            color: isSelected ? Colors.white : Colors.grey.shade700,
                            fontWeight: isSelected ? FontWeight.w600 : FontWeight.normal,
                          ),
                        ),
                      ),
                    ),
                  );
                },
              ),
            ),
          ),

          // Portfolio Chart
          SliverToBoxAdapter(
            child: Container(
              height: 300,
              margin: const EdgeInsets.symmetric(horizontal: 16),
              child: PortfolioChart(
                data: state.chartData,
                timeframe: _selectedTimeframe,
              ),
            ),
          ),

          // Tab Bar
          SliverToBoxAdapter(
            child: Container(
              margin: const EdgeInsets.only(top: 32, bottom: 16),
              child: TabBar(
                controller: _tabController,
                labelColor: AppTheme.primaryColor,
                unselectedLabelColor: Colors.grey,
                indicatorColor: AppTheme.primaryColor,
                tabs: const [
                  Tab(text: 'Assets'),
                  Tab(text: 'Performance'),
                  Tab(text: 'Allocation'),
                ],
              ),
            ),
          ),

          // Tab Content
          SliverFillRemaining(
            child: TabBarView(
              controller: _tabController,
              children: [
                // Assets Tab
                HoldingsList(
                  holdings: state.portfolio.holdings,
                  onAssetTap: (asset) {
                    // Navigate to asset details
                  },
                ),

                // Performance Tab
                PerformanceMetrics(
                  performance: state.portfolio.performance,
                ),

                // Allocation Tab
                AssetAllocation(
                  holdings: state.portfolio.holdings,
                  totalValue: state.portfolio.totalValue,
                ),
              ],
            ),
          ),
        ],
      ),
    );
  }

  Widget _buildErrorState(String message) {
    return Center(
      child: Column(
        mainAxisAlignment: MainAxisAlignment.center,
        children: [
          Icon(
            Icons.error_outline,
            size: 80,
            color: AppTheme.errorColor,
          ),
          const SizedBox(height: 20),
          Text(
            'Failed to load portfolio',
            style: Theme.of(context).textTheme.headlineMedium,
          ),
          const SizedBox(height: 8),
          Text(
            message,
            textAlign: TextAlign.center,
            style: Theme.of(context).textTheme.bodyMedium,
          ),
          const SizedBox(height: 32),
          ElevatedButton(
            onPressed: _loadPortfolioData,
            child: const Text('Try Again'),
          ),
        ],
      ),
    );
  }

  @override
  void dispose() {
    _tabController.dispose();
    super.dispose();
  }
}