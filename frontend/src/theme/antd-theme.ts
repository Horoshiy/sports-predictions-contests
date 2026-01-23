import type { ThemeConfig } from 'antd';

export const antdTheme: ThemeConfig = {
  token: {
    // Primary color matching MUI theme
    colorPrimary: '#1976d2',
    colorSuccess: '#2e7d32',
    colorWarning: '#ed6c02',
    colorError: '#d32f2f',
    colorInfo: '#0288d1',
    
    // Typography
    fontFamily: 'Roboto, Arial, sans-serif',
    fontSize: 14,
    
    // Border radius
    borderRadius: 4,
    
    // Layout
    colorBgContainer: '#ffffff',
    colorBgLayout: '#f5f5f5',
  },
  components: {
    Button: {
      controlHeight: 36,
      borderRadius: 4,
    },
    Input: {
      controlHeight: 40,
      borderRadius: 4,
    },
    Select: {
      controlHeight: 40,
    },
    Card: {
      borderRadiusLG: 4,
    },
    Table: {
      borderRadius: 4,
    },
  },
};
