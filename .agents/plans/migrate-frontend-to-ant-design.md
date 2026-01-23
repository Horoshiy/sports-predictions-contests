# Feature: Migrate Frontend from Material-UI to Ant Design

## Feature Description

Complete migration of the Sports Prediction Contests frontend from Material-UI (MUI) to Ant Design. This involves replacing all MUI components with their Ant Design equivalents across 48 component files and 8 page files, updating styling approaches, and ensuring consistent design language throughout the application.

## User Story

As a developer
I want to migrate the frontend from Material-UI to Ant Design
So that the application uses a more enterprise-focused UI library with better out-of-the-box components for data-heavy applications and admin interfaces

## Problem Statement

The current frontend uses Material-UI (MUI) v5 with extensive component usage across the application. While MUI is excellent, Ant Design provides:
- Better default components for admin/dashboard interfaces
- More comprehensive data display components (Table, Descriptions, etc.)
- Stronger design system for enterprise applications
- Better internationalization support out of the box
- More opinionated design that reduces custom styling needs

The migration requires systematic replacement of all MUI components while maintaining existing functionality and improving the overall user experience.

## Solution Statement

Systematically replace all Material-UI components with Ant Design equivalents by:
1. Installing Ant Design and removing MUI dependencies
2. Creating a comprehensive component mapping guide
3. Migrating components in logical groups (auth → layout → data display → forms)
4. Updating theme configuration and global styles
5. Testing each migrated section thoroughly
6. Removing all MUI dependencies after complete migration

## Feature Metadata

**Feature Type**: Refactor
**Estimated Complexity**: High
**Primary Systems Affected**: 
- Frontend (all 48 components + 8 pages)
- Theme configuration
- Build system
- Testing infrastructure

**Dependencies**: 
- antd@^5.22.0
- @ant-design/icons@^5.5.0
- dayjs@^1.11.13 (replaces date-fns for DatePicker)

---

## CONTEXT REFERENCES

### Relevant Codebase Files - IMPORTANT: YOU MUST READ THESE FILES BEFORE IMPLEMENTING!

**Current MUI Implementation Examples:**
- `frontend/src/App.tsx` (lines 1-100) - Main app structure with MUI ThemeProvider, AppBar, navigation
- `frontend/src/components/auth/LoginForm.tsx` - Form patterns with MUI TextField, Button, Paper
- `frontend/src/components/contests/ContestCard.tsx` - Card layout with MUI Card, CardContent, Chip
- `frontend/src/components/leaderboard/LeaderboardTable.tsx` - Table implementation with material-react-table
- `frontend/src/pages/ContestsPage.tsx` - Page layout with MUI Container, Tabs
- `frontend/src/pages/SportsPage.tsx` - Complex page with multiple MUI components
- `frontend/src/contexts/ToastContext.tsx` - Snackbar/Alert usage

**Configuration Files:**
- `frontend/package.json` - Current dependencies to update
- `frontend/vite.config.ts` - Build configuration
- `frontend/tsconfig.json` - TypeScript configuration

### New Files to Create

- `frontend/src/theme/antd-theme.ts` - Ant Design theme configuration
- `frontend/src/utils/antd-helpers.ts` - Helper functions for Ant Design patterns
- `.agents/plans/migration-checklist.md` - Detailed migration tracking document

### Relevant Documentation - YOU SHOULD READ THESE BEFORE IMPLEMENTING!

- [Ant Design Components Overview](https://ant.design/components/overview/)
  - Specific section: All component categories
  - Why: Complete reference for available components and their APIs

- [Ant Design Form Component](https://ant.design/components/form/)
  - Specific section: Form, Form.Item, Form.List API
  - Why: Critical for form migration - different pattern from MUI

- [Ant Design Table Component](https://ant.design/components/table/)
  - Specific section: Table API, columns configuration
  - Why: Replace material-react-table with native Ant Design Table

- [Ant Design ConfigProvider](https://ant.design/components/config-provider/)
  - Specific section: Theme configuration, locale
  - Why: Global configuration and theming

- [Ant Design Migration Guide](https://ant.design/docs/react/migration-v6/)
  - Specific section: Breaking changes, component updates
  - Why: Understand latest version features

### Patterns to Follow

**Component Mapping (MUI → Ant Design):**

```typescript
// MUI Pattern
import { Button, TextField, Box } from '@mui/material';

<Box sx={{ p: 2 }}>
  <TextField label="Email" fullWidth />
  <Button variant="contained">Submit</Button>
</Box>

// Ant Design Pattern
import { Button, Input, Space } from 'antd';

<Space direction="vertical" style={{ padding: 16, width: '100%' }}>
  <Input placeholder="Email" />
  <Button type="primary">Submit</Button>
</Space>
```

**Form Pattern Migration:**

```typescript
// MUI + react-hook-form
import { useForm, Controller } from 'react-hook-form';
import { TextField } from '@mui/material';

const { control } = useForm();
<Controller
  name="email"
  control={control}
  render={({ field }) => <TextField {...field} />}
/>

// Ant Design Form
import { Form, Input } from 'antd';

<Form>
  <Form.Item name="email" label="Email" rules={[{ required: true }]}>
    <Input />
  </Form.Item>
</Form>
```

**Theme Configuration:**

```typescript
// MUI Theme
import { createTheme, ThemeProvider } from '@mui/material/styles';

const theme = createTheme({
  palette: {
    primary: { main: '#1976d2' },
  },
});

// Ant Design Theme
import { ConfigProvider, theme } from 'antd';

<ConfigProvider
  theme={{
    token: {
      colorPrimary: '#1976d2',
    },
  }}
>
  {children}
</ConfigProvider>
```

**Notification/Toast Pattern:**

```typescript
// MUI Snackbar
import { Snackbar, Alert } from '@mui/material';
const [open, setOpen] = useState(false);
<Snackbar open={open}>
  <Alert severity="success">Success!</Alert>
</Snackbar>

// Ant Design Message/Notification
import { message, notification } from 'antd';
message.success('Success!');
// or
notification.success({
  message: 'Success',
  description: 'Operation completed',
});
```

---

## IMPLEMENTATION PLAN

### Phase 1: Foundation & Setup

Prepare the project for migration by installing Ant Design, creating theme configuration, and setting up the migration infrastructure.

**Tasks:**
- Install Ant Design dependencies and remove MUI
- Create Ant Design theme configuration
- Set up component mapping documentation
- Create migration tracking system

### Phase 2: Core Layout & Navigation

Migrate the main application structure, navigation, and layout components.

**Tasks:**
- Migrate App.tsx main structure
- Replace AppBar with Ant Design Layout.Header
- Update navigation patterns
- Migrate routing and protected routes

### Phase 3: Authentication Components

Migrate all authentication-related components (login, register, forms).

**Tasks:**
- Migrate LoginForm component
- Migrate RegisterForm component
- Update authentication pages
- Migrate ProtectedRoute component

### Phase 4: Data Display Components

Migrate components that display data (cards, tables, lists).

**Tasks:**
- Migrate all Card components
- Replace material-react-table with Ant Design Table
- Migrate List components
- Update data display patterns

### Phase 5: Form Components

Migrate all form-related components and pages.

**Tasks:**
- Migrate form components (ContestForm, TeamForm, etc.)
- Update form validation patterns
- Migrate date pickers and selectors
- Update form submission handlers

### Phase 6: Pages & Complex Components

Migrate page-level components and complex multi-component pages.

**Tasks:**
- Migrate all page components
- Update complex layouts
- Migrate dashboard/analytics components
- Update sports management pages

### Phase 7: Utilities & Context

Migrate utility components, contexts, and global state.

**Tasks:**
- Migrate ToastContext to Ant Design message/notification
- Update utility functions
- Migrate hooks to Ant Design patterns
- Update global styles

### Phase 8: Testing & Cleanup

Test the migrated application and remove all MUI dependencies.

**Tasks:**
- Test all migrated components
- Fix styling inconsistencies
- Remove all MUI imports and dependencies
- Update documentation

---

## STEP-BY-STEP TASKS

### Phase 1: Foundation & Setup

#### Task 1.1: UPDATE package.json dependencies

- **IMPLEMENT**: Remove all MUI dependencies and add Ant Design
- **PATTERN**: Standard npm dependency management
- **IMPORTS**: None (package.json modification)
- **GOTCHA**: Keep react-hook-form and zod for form validation
- **VALIDATE**: `cd frontend && npm install`

```json
// Remove these dependencies:
"@emotion/react": "^11.11.1",
"@emotion/styled": "^11.11.0",
"@mui/icons-material": "^5.14.18",
"@mui/material": "^5.14.18",
"@mui/x-date-pickers": "^6.18.1",
"material-react-table": "^2.0.4",
"date-fns": "^2.30.0",

// Add these dependencies:
"antd": "^5.22.0",
"@ant-design/icons": "^5.5.0",
"dayjs": "^1.11.13",
```

#### Task 1.2: CREATE frontend/src/theme/antd-theme.ts

- **IMPLEMENT**: Ant Design theme configuration matching current MUI theme
- **PATTERN**: ConfigProvider theme token configuration
- **IMPORTS**: `import type { ThemeConfig } from 'antd';`
- **GOTCHA**: Use token-based theming, not CSS variables
- **VALIDATE**: TypeScript compilation `npm run build`

#### Task 1.3: CREATE frontend/src/utils/antd-helpers.ts

- **IMPLEMENT**: Helper functions for common Ant Design patterns
- **PATTERN**: Utility functions for message, notification, form helpers
- **IMPORTS**: `import { message, notification } from 'antd';`
- **GOTCHA**: Ant Design uses imperative API for messages
- **VALIDATE**: TypeScript compilation

#### Task 1.4: CREATE .agents/plans/migration-checklist.md

- **IMPLEMENT**: Detailed checklist of all files to migrate
- **PATTERN**: Markdown checklist with file paths and status
- **IMPORTS**: None (documentation file)
- **GOTCHA**: Track both component files and their imports
- **VALIDATE**: Manual review

### Phase 2: Core Layout & Navigation

#### Task 2.1: UPDATE frontend/src/App.tsx

- **IMPLEMENT**: Replace MUI ThemeProvider with Ant Design ConfigProvider
- **PATTERN**: Wrap app with ConfigProvider, use Layout components
- **IMPORTS**: 
```typescript
import { ConfigProvider, Layout, Menu, Dropdown, Avatar, Button } from 'antd';
import { UserOutlined, LogoutOutlined } from '@ant-design/icons';
import { antdTheme } from './theme/antd-theme';
```
- **GOTCHA**: Layout.Header replaces AppBar, Menu replaces navigation buttons
- **VALIDATE**: `npm run dev` - app should render without errors

#### Task 2.2: UPDATE frontend/src/components/auth/ProtectedRoute.tsx

- **IMPLEMENT**: Replace MUI Box, CircularProgress, Typography, Button with Ant Design equivalents
- **PATTERN**: Use Spin for loading, Result for error states
- **IMPORTS**: `import { Spin, Result, Button } from 'antd';`
- **GOTCHA**: Result component provides better error/empty states
- **VALIDATE**: Navigate to protected route while logged out

### Phase 3: Authentication Components

#### Task 3.1: UPDATE frontend/src/components/auth/LoginForm.tsx

- **IMPLEMENT**: Replace MUI form components with Ant Design Form
- **PATTERN**: Use Form, Form.Item, Input, Button from Ant Design
- **IMPORTS**:
```typescript
import { Form, Input, Button, Checkbox } from 'antd';
import { MailOutlined, LockOutlined, EyeInvisibleOutlined, EyeTwoTone } from '@ant-design/icons';
```
- **GOTCHA**: Ant Design Form has built-in validation, may need to adjust react-hook-form integration
- **VALIDATE**: Test login form submission and validation

#### Task 3.2: UPDATE frontend/src/components/auth/RegisterForm.tsx

- **IMPLEMENT**: Replace MUI form components with Ant Design Form
- **PATTERN**: Similar to LoginForm, use Form.Item for each field
- **IMPORTS**: Same as LoginForm plus additional icons
- **GOTCHA**: Password confirmation validation pattern differs
- **VALIDATE**: Test registration form with all validations

#### Task 3.3: UPDATE frontend/src/pages/LoginPage.tsx

- **IMPLEMENT**: Replace MUI Container, Box, Typography with Ant Design Layout, Space, Typography
- **PATTERN**: Use Layout for page structure, Space for spacing
- **IMPORTS**: `import { Layout, Space, Typography, Card } from 'antd';`
- **GOTCHA**: Card component provides better visual container
- **VALIDATE**: Visual inspection of login page

#### Task 3.4: UPDATE frontend/src/pages/RegisterPage.tsx

- **IMPLEMENT**: Same pattern as LoginPage
- **PATTERN**: Layout + Card structure
- **IMPORTS**: Same as LoginPage
- **GOTCHA**: Maintain consistent styling with LoginPage
- **VALIDATE**: Visual inspection of register page

### Phase 4: Data Display Components

#### Task 4.1: UPDATE frontend/src/components/contests/ContestCard.tsx

- **IMPLEMENT**: Replace MUI Card with Ant Design Card
- **PATTERN**: Use Card with actions prop for buttons
- **IMPORTS**:
```typescript
import { Card, Tag, Button, Space, Tooltip } from 'antd';
import { EditOutlined, DeleteOutlined, TeamOutlined, TrophyOutlined } from '@ant-design/icons';
```
- **GOTCHA**: Card actions are passed as array, not CardActions component
- **VALIDATE**: Render contest list and verify card interactions

#### Task 4.2: UPDATE frontend/src/components/contests/ContestList.tsx

- **IMPLEMENT**: Replace MUI Grid with Ant Design Row/Col or List
- **PATTERN**: Use List component with grid prop for card layout
- **IMPORTS**: `import { List, Row, Col, Empty, Spin } from 'antd';`
- **GOTCHA**: List.grid provides responsive grid layout
- **VALIDATE**: Test responsive behavior and loading states

#### Task 4.3: UPDATE frontend/src/components/leaderboard/LeaderboardTable.tsx

- **IMPLEMENT**: Replace material-react-table with Ant Design Table
- **PATTERN**: Define columns array, use Table component
- **IMPORTS**: `import { Table, Tag, Avatar, Typography } from 'antd';`
- **GOTCHA**: Column configuration is different, need to map existing columns
- **VALIDATE**: Test sorting, pagination, and data display

#### Task 4.4: UPDATE frontend/src/components/leaderboard/UserScore.tsx

- **IMPLEMENT**: Replace MUI components with Ant Design Descriptions or Card
- **PATTERN**: Use Descriptions for key-value pairs
- **IMPORTS**: `import { Descriptions, Card, Statistic, Row, Col } from 'antd';`
- **GOTCHA**: Statistic component great for displaying scores
- **VALIDATE**: Visual inspection of score display

#### Task 4.5: UPDATE frontend/src/components/teams/TeamList.tsx

- **IMPLEMENT**: Replace MUI components with Ant Design Table or List
- **PATTERN**: Use Table for structured data display
- **IMPORTS**: `import { Table, Button, Tag, Space, Tooltip } from 'antd';`
- **GOTCHA**: Action buttons in table columns
- **VALIDATE**: Test team list interactions

#### Task 4.6: UPDATE frontend/src/components/teams/TeamMembers.tsx

- **IMPLEMENT**: Replace MUI List with Ant Design List
- **PATTERN**: Use List with List.Item for members
- **IMPORTS**: `import { List, Avatar, Tag, Button, Popconfirm } from 'antd';`
- **GOTCHA**: Popconfirm for delete confirmation
- **VALIDATE**: Test member list and actions

### Phase 5: Form Components

#### Task 5.1: UPDATE frontend/src/components/contests/ContestForm.tsx

- **IMPLEMENT**: Replace MUI Dialog and form components with Ant Design Modal and Form
- **PATTERN**: Use Modal + Form combination
- **IMPORTS**:
```typescript
import { Modal, Form, Input, DatePicker, Select, InputNumber, Switch } from 'antd';
import dayjs from 'dayjs';
```
- **GOTCHA**: DatePicker uses dayjs instead of date-fns
- **VALIDATE**: Test form submission and validation

#### Task 5.2: UPDATE frontend/src/components/teams/TeamForm.tsx

- **IMPLEMENT**: Replace MUI Dialog with Ant Design Modal
- **PATTERN**: Modal + Form with explicit field passing (already done)
- **IMPORTS**: `import { Modal, Form, Input, InputNumber, Button } from 'antd';`
- **GOTCHA**: Form.Item name prop for field binding
- **VALIDATE**: Test team creation and editing

#### Task 5.3: UPDATE frontend/src/components/sports/SportForm.tsx

- **IMPLEMENT**: Replace MUI Dialog with Ant Design Modal
- **PATTERN**: Modal + Form for sport management
- **IMPORTS**: `import { Modal, Form, Input, Button } from 'antd';`
- **GOTCHA**: Keep existing validation logic
- **VALIDATE**: Test sport CRUD operations

#### Task 5.4: UPDATE frontend/src/components/sports/LeagueForm.tsx

- **IMPLEMENT**: Replace MUI components with Ant Design
- **PATTERN**: Modal + Form with Select for sport selection
- **IMPORTS**: `import { Modal, Form, Input, Select, Button } from 'antd';`
- **GOTCHA**: Select component for sport dropdown
- **VALIDATE**: Test league creation with sport selection

#### Task 5.5: UPDATE frontend/src/components/sports/TeamForm.tsx (sports)

- **IMPLEMENT**: Replace MUI components with Ant Design
- **PATTERN**: Modal + Form with league selection
- **IMPORTS**: `import { Modal, Form, Input, Select, Upload, Button } from 'antd';`
- **GOTCHA**: Upload component for logo
- **VALIDATE**: Test team creation in sports context

#### Task 5.6: UPDATE frontend/src/components/sports/MatchForm.tsx

- **IMPLEMENT**: Replace MUI components with Ant Design
- **PATTERN**: Modal + Form with DatePicker and team selects
- **IMPORTS**: `import { Modal, Form, Select, DatePicker, InputNumber, Button } from 'antd';`
- **GOTCHA**: DatePicker with showTime for match scheduling
- **VALIDATE**: Test match creation and editing

#### Task 5.7: UPDATE frontend/src/components/predictions/PredictionForm.tsx

- **IMPLEMENT**: Replace MUI components with Ant Design
- **PATTERN**: Form with dynamic fields based on prediction type
- **IMPORTS**: `import { Form, Input, InputNumber, Select, Radio, Button } from 'antd';`
- **GOTCHA**: Form.List for dynamic prediction fields
- **VALIDATE**: Test different prediction types

### Phase 6: Pages & Complex Components

#### Task 6.1: UPDATE frontend/src/pages/ContestsPage.tsx

- **IMPLEMENT**: Replace MUI Container, Tabs with Ant Design Layout, Tabs
- **PATTERN**: Use Layout.Content + Tabs for page structure
- **IMPORTS**: `import { Layout, Tabs, Button, Space } from 'antd';`
- **GOTCHA**: Tabs items prop instead of Tab children
- **VALIDATE**: Test tab navigation and contest display

#### Task 6.2: UPDATE frontend/src/pages/SportsPage.tsx

- **IMPLEMENT**: Replace MUI components with Ant Design
- **PATTERN**: Layout + Tabs for sports management sections
- **IMPORTS**: `import { Layout, Tabs, Typography, Space } from 'antd';`
- **GOTCHA**: Complex page with multiple forms and lists
- **VALIDATE**: Test all sports management operations

#### Task 6.3: UPDATE frontend/src/pages/TeamsPage.tsx

- **IMPLEMENT**: Replace MUI components with Ant Design
- **PATTERN**: Layout + Tabs for team sections
- **IMPORTS**: `import { Layout, Tabs, Form, Input, Button, Modal } from 'antd';`
- **GOTCHA**: Join team form in separate tab
- **VALIDATE**: Test team operations and join flow

#### Task 6.4: UPDATE frontend/src/pages/PredictionsPage.tsx

- **IMPLEMENT**: Replace MUI components with Ant Design
- **PATTERN**: Layout + List/Table for predictions
- **IMPORTS**: `import { Layout, List, Card, Button, Space, Empty } from 'antd';`
- **GOTCHA**: Event list with prediction forms
- **VALIDATE**: Test prediction submission

#### Task 6.5: UPDATE frontend/src/pages/AnalyticsPage.tsx

- **IMPLEMENT**: Replace MUI components with Ant Design
- **PATTERN**: Layout + Card grid for analytics widgets
- **IMPORTS**: `import { Layout, Card, Row, Col, Statistic, Typography } from 'antd';`
- **GOTCHA**: Statistic component for metrics display
- **VALIDATE**: Visual inspection of analytics dashboard

#### Task 6.6: UPDATE frontend/src/pages/ProfilePage.tsx

- **IMPLEMENT**: Replace MUI components with Ant Design
- **PATTERN**: Layout + Tabs for profile sections
- **IMPORTS**: `import { Layout, Tabs, Card, Avatar, Upload, Button } from 'antd';`
- **GOTCHA**: Upload component for avatar
- **VALIDATE**: Test profile editing and avatar upload

### Phase 7: Utilities & Context

#### Task 7.1: UPDATE frontend/src/contexts/ToastContext.tsx

- **IMPLEMENT**: Replace MUI Snackbar with Ant Design message/notification
- **PATTERN**: Use message API for simple toasts, notification for complex ones
- **IMPORTS**: `import { message, notification } from 'antd';`
- **GOTCHA**: Imperative API instead of component-based
- **VALIDATE**: Test toast notifications throughout app

#### Task 7.2: UPDATE frontend/src/components/analytics/ExportButton.tsx

- **IMPLEMENT**: Replace MUI Button and CircularProgress with Ant Design
- **PATTERN**: Button with loading prop
- **IMPORTS**: `import { Button } from 'antd'; import { DownloadOutlined } from '@ant-design/icons';`
- **GOTCHA**: loading prop handles spinner automatically
- **VALIDATE**: Test export functionality

#### Task 7.3: UPDATE frontend/src/components/predictions/CoefficientIndicator.tsx

- **IMPLEMENT**: Replace MUI components with Ant Design
- **PATTERN**: Use Tag or Badge for coefficient display
- **IMPORTS**: `import { Tag, Tooltip, Typography } from 'antd';`
- **GOTCHA**: Tag color prop for visual indication
- **VALIDATE**: Visual inspection of coefficient display

#### Task 7.4: UPDATE frontend/src/components/challenges/ChallengeCard.tsx

- **IMPLEMENT**: Replace MUI Card with Ant Design Card
- **PATTERN**: Card with actions and meta
- **IMPORTS**: `import { Card, Tag, Button, Avatar, Space } from 'antd';`
- **GOTCHA**: Card.Meta for user info
- **VALIDATE**: Test challenge card display

#### Task 7.5: UPDATE frontend/src/components/challenges/ChallengeList.tsx

- **IMPLEMENT**: Replace MUI Grid with Ant Design List
- **PATTERN**: List with grid layout
- **IMPORTS**: `import { List, Empty, Spin } from 'antd';`
- **GOTCHA**: List.grid for responsive layout
- **VALIDATE**: Test challenge list display

### Phase 8: Testing & Cleanup

#### Task 8.1: REMOVE all MUI imports

- **IMPLEMENT**: Search and remove all @mui imports across codebase
- **PATTERN**: Find and replace operation
- **IMPORTS**: None (removal task)
- **GOTCHA**: Check for any remaining MUI references
- **VALIDATE**: `grep -r "@mui" frontend/src` should return nothing

#### Task 8.2: REMOVE MUI dependencies from package.json

- **IMPLEMENT**: Remove all MUI packages from dependencies
- **PATTERN**: Edit package.json and run npm install
- **IMPORTS**: None
- **GOTCHA**: Also remove @emotion packages if not used elsewhere
- **VALIDATE**: `npm install && npm run build`

#### Task 8.3: UPDATE global styles

- **IMPLEMENT**: Remove MUI-specific global styles, add Ant Design customizations
- **PATTERN**: Update CSS files to work with Ant Design
- **IMPORTS**: None (CSS modification)
- **GOTCHA**: Ant Design has its own reset styles
- **VALIDATE**: Visual inspection across all pages

#### Task 8.4: TEST all pages and components

- **IMPLEMENT**: Manual testing of all functionality
- **PATTERN**: Systematic page-by-page testing
- **IMPORTS**: None (testing task)
- **GOTCHA**: Test responsive behavior on different screen sizes
- **VALIDATE**: All features work as before migration

#### Task 8.5: UPDATE documentation

- **IMPLEMENT**: Update README and component documentation
- **PATTERN**: Document Ant Design usage patterns
- **IMPORTS**: None (documentation task)
- **GOTCHA**: Include theme customization guide
- **VALIDATE**: Documentation review

---

## TESTING STRATEGY

### Unit Tests

- Test form validation with Ant Design Form
- Test component rendering with Ant Design components
- Test utility functions for Ant Design helpers
- Mock Ant Design message/notification APIs

### Integration Tests

- Test complete user flows (login, create contest, make prediction)
- Test form submission and validation
- Test navigation and routing
- Test responsive behavior

### Visual Testing

- Compare before/after screenshots of all pages
- Test on different screen sizes (mobile, tablet, desktop)
- Verify consistent styling across components
- Check accessibility (ARIA labels, keyboard navigation)

### Edge Cases

- Test empty states with Ant Design Empty component
- Test loading states with Spin component
- Test error states with Result component
- Test form validation edge cases
- Test date picker with different locales

---

## VALIDATION COMMANDS

### Level 1: Syntax & Style

```bash
cd frontend
npm run lint
npm run build
```

### Level 2: Type Checking

```bash
cd frontend
npx tsc --noEmit
```

### Level 3: Build Verification

```bash
cd frontend
npm run build
npm run preview
```

### Level 4: Manual Validation

1. Start development server: `cd frontend && npm run dev`
2. Test authentication flow (login/register)
3. Test contest creation and management
4. Test predictions submission
5. Test team management
6. Test sports management
7. Test analytics dashboard
8. Test profile management
9. Verify responsive design on mobile/tablet
10. Test all form validations

### Level 5: Cross-browser Testing

- Test on Chrome, Firefox, Safari
- Test on mobile browsers (iOS Safari, Chrome Mobile)
- Verify consistent behavior across browsers

---

## ACCEPTANCE CRITERIA

- [ ] All MUI dependencies removed from package.json
- [ ] All components migrated to Ant Design
- [ ] No MUI imports remaining in codebase
- [ ] Application builds without errors
- [ ] All existing functionality works as before
- [ ] Forms validate correctly with Ant Design Form
- [ ] Tables display and sort data correctly
- [ ] Navigation and routing work properly
- [ ] Responsive design works on all screen sizes
- [ ] Theme configuration applied consistently
- [ ] Loading and error states display correctly
- [ ] Notifications/messages work properly
- [ ] Date pickers work with dayjs
- [ ] All icons display correctly
- [ ] Accessibility maintained (keyboard navigation, ARIA labels)
- [ ] Performance is equal or better than before
- [ ] No console errors or warnings
- [ ] Documentation updated

---

## COMPLETION CHECKLIST

- [ ] Phase 1: Foundation & Setup completed
- [ ] Phase 2: Core Layout & Navigation completed
- [ ] Phase 3: Authentication Components completed
- [ ] Phase 4: Data Display Components completed
- [ ] Phase 5: Form Components completed
- [ ] Phase 6: Pages & Complex Components completed
- [ ] Phase 7: Utilities & Context completed
- [ ] Phase 8: Testing & Cleanup completed
- [ ] All validation commands pass
- [ ] Manual testing completed
- [ ] Cross-browser testing completed
- [ ] Documentation updated
- [ ] Code review completed

---

## NOTES

### Component Mapping Reference

| MUI Component | Ant Design Equivalent | Notes |
|---------------|----------------------|-------|
| Box | Space, Flex | Use Space for spacing, Flex for flexbox |
| Container | Layout.Content | Wrap with Layout |
| Paper | Card | Card provides better styling |
| Typography | Typography | Similar API |
| Button | Button | type prop instead of variant |
| TextField | Input | Separate Input.TextArea for multiline |
| Select | Select | Similar API |
| Checkbox | Checkbox | Similar API |
| Radio | Radio | Use Radio.Group |
| Switch | Switch | checked prop instead of value |
| Slider | Slider | Similar API |
| Dialog | Modal | Similar API |
| Snackbar | message, notification | Imperative API |
| Alert | Alert | Similar API |
| Chip | Tag | Similar concept |
| Avatar | Avatar | Similar API |
| Badge | Badge | Similar API |
| Tooltip | Tooltip | Similar API |
| Menu | Menu | Different structure |
| Tabs | Tabs | items prop instead of children |
| Table | Table | columns configuration |
| Pagination | Pagination | Similar API |
| CircularProgress | Spin | Similar concept |
| LinearProgress | Progress | type="line" |
| Skeleton | Skeleton | Similar API |
| Drawer | Drawer | Similar API |
| AppBar | Layout.Header | Part of Layout |
| Grid | Row, Col | 24-column grid system |
| Stack | Space | direction prop |
| Divider | Divider | Similar API |
| List | List | Similar API |
| Card | Card | Similar API |
| Accordion | Collapse | Similar concept |
| DatePicker | DatePicker | Uses dayjs |
| TimePicker | TimePicker | Uses dayjs |

### Design Decisions

1. **Form Management**: Keep react-hook-form for complex validation logic, integrate with Ant Design Form for UI
2. **Date Handling**: Migrate from date-fns to dayjs (required by Ant Design DatePicker)
3. **Icons**: Use @ant-design/icons instead of @mui/icons-material
4. **Theme**: Create custom theme matching current color scheme
5. **Table**: Replace material-react-table with native Ant Design Table (simpler, better performance)
6. **Notifications**: Use message for simple toasts, notification for complex messages
7. **Layout**: Use Ant Design Layout system for consistent page structure

### Migration Strategy

- Migrate in phases to maintain working application
- Test each phase before moving to next
- Keep git commits small and focused
- Document any issues or deviations from plan
- Maintain feature parity throughout migration

### Performance Considerations

- Ant Design bundle size is similar to MUI
- Tree-shaking works well with both libraries
- Consider lazy loading for large components
- Monitor bundle size after migration

### Accessibility Considerations

- Ant Design has good accessibility support
- Verify ARIA labels after migration
- Test keyboard navigation
- Ensure color contrast meets WCAG standards

---

**Estimated Time**: 16-20 hours
**Confidence Score**: 8/10 - Straightforward migration with well-documented patterns, main risk is ensuring all edge cases are covered
