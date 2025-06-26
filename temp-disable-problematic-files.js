const fs = require('fs');
const path = require('path');

// Files to disable by renaming with .disabled extension
const filesToDisable = [
  // Missing composables
  'composables/useXrpAmmHeatmap.ts',
  'composables/useXrpAmmSwap.ts',
  'composables/useXrpAmmTransactions.ts',
  'composables/useXrpHeatmap.ts',
  'composables/useXrpTokenHeatmap.ts',
  'composables/useXrpGridRenderers.ts',
  'composables/useXrpAmmLiveData.ts',
  'composables/useXrpAmmSwap.ts',
  'composables/useXrpAmmTransactions.ts',
  'composables/useXrpConfigs.ts',
  'composables/useXrpHeatmap.ts',
  'composables/useXrpTokenHeatmap.ts',
  'composables/useXrpAmmHeatmap.ts',
  'composables/useXrpAmmLiveData.ts',
  'composables/useXrpAmmSwap.ts',
  'composables/useXrpAmmTransactions.ts',
  'composables/useXrpConfigs.ts',
  'composables/useXrpGridRenderers.ts',
  'composables/useXrpHeatmap.ts',
  'composables/useXrpTokenHeatmap.ts',
  
  // Missing plugins
  'plugins/web3/enhanced-xrp.client.ts',
  'plugins/web3/metamask-xrp-snap.connector.ts',
  'plugins/web3/xaman.connector.ts',
  
  // Missing components
  'components/xrp/XrpTokenHeatmap.vue',
  'components/xrp/XrpTerminal.vue',
  'components/xrp/XrpAmmHeatmap.vue',
  'components/xrp/XrpAmmSwapDialog.vue',
  'components/xrp/XrpAmmActionDialog.vue',
  
  // Missing pages
  'pages/xrp-amm-heatmap.vue',
  'pages/xrp-heatmap.vue',
  'pages/xrp-amm-pools/_id.vue',
  
  // Problematic mixins
  'mixins/configBarMixin.ts',
  
  // Problematic pages with route issues
  'pages/xrp-balances.vue',
  'pages/xrp-transactions.vue',
  'pages/xrp-terminal.vue',
  'pages/terminal.vue',
  
  // Problematic components
  'components/HeatmapChart.vue',
  'components/HeatmapConfigMenu.vue',
  'components/xrp/XrpHeatmapChart.vue',
  'components/common/GlobalSearch.vue',
  'components/common/EnhancedWalletSelectDialog.vue',
  'components/common/EnhancedXrpWalletConnector.vue'
];

console.log('Disabling problematic files...');

filesToDisable.forEach(filePath => {
  const fullPath = path.join(__dirname, filePath);
  const disabledPath = fullPath + '.disabled';
  
  if (fs.existsSync(fullPath)) {
    try {
      fs.renameSync(fullPath, disabledPath);
      console.log(`✓ Disabled: ${filePath}`);
    } catch (error) {
      console.log(`✗ Failed to disable ${filePath}:`, error.message);
    }
  } else {
    console.log(`- Not found: ${filePath}`);
  }
});

console.log('Done disabling files.'); 