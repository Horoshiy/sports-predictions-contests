import React from 'react'
import {
  Box,
  Typography,
  Paper,
  LinearProgress,
  List,
  ListItem,
  ListItemIcon,
  ListItemText,
  Chip,
  Alert,
  Button,
} from '@mui/material'
import {
  CheckCircle,
  RadioButtonUnchecked,
  TrendingUp,
  Person,
  PhotoCamera,
  LocationOn,
  Language,
  Email,
} from '@mui/icons-material'

// Rename the interface to avoid conflict
interface ProfileCompletionData {
  percentage: number
  missingFields: string[]
  suggestions: string[]
}

interface ProfileCompletionProps {
  completion: ProfileCompletionData
  onFieldClick?: (field: string) => void
}

export const ProfileCompletion: React.FC<ProfileCompletionProps> = ({
  completion,
  onFieldClick,
}) => {
  const { percentage, missingFields, suggestions } = completion

  const getFieldIcon = (field: string) => {
    switch (field) {
      case 'name':
        return <Person />
      case 'email':
        return <Email />
      case 'bio':
        return <Person />
      case 'location':
        return <LocationOn />
      case 'avatar':
        return <PhotoCamera />
      case 'website':
        return <Language />
      default:
        return <RadioButtonUnchecked />
    }
  }

  const getFieldLabel = (field: string) => {
    switch (field) {
      case 'name':
        return 'Full Name'
      case 'email':
        return 'Email Address'
      case 'bio':
        return 'Bio/Description'
      case 'location':
        return 'Location'
      case 'avatar':
        return 'Profile Picture'
      case 'website':
        return 'Website'
      case 'twitter':
        return 'Twitter Profile'
      case 'linkedin':
        return 'LinkedIn Profile'
      case 'github':
        return 'GitHub Profile'
      default:
        return field.charAt(0).toUpperCase() + field.slice(1)
    }
  }

  const getCompletionColor = (percentage: number) => {
    if (percentage >= 80) return 'success'
    if (percentage >= 50) return 'warning'
    return 'error'
  }

  const getCompletionMessage = (percentage: number) => {
    if (percentage >= 90) return 'Excellent! Your profile is almost complete.'
    if (percentage >= 70) return 'Great progress! Just a few more details to go.'
    if (percentage >= 50) return 'Good start! Complete more fields to unlock all features.'
    return 'Get started by filling out your basic profile information.'
  }

  return (
    <Paper elevation={3} sx={{ p: 4, maxWidth: 600, mx: 'auto' }}>
      <Box sx={{ mb: 3 }}>
        <Typography variant="h5" component="h2" gutterBottom>
          <TrendingUp sx={{ mr: 1, verticalAlign: 'middle' }} />
          Profile Completion
        </Typography>
        <Typography variant="body2" color="text.secondary">
          Complete your profile to unlock all features and improve your experience
        </Typography>
      </Box>

      {/* Progress Bar */}
      <Box sx={{ mb: 3 }}>
        <Box sx={{ display: 'flex', alignItems: 'center', mb: 1 }}>
          <Typography variant="h4" sx={{ mr: 2, fontWeight: 'bold' }}>
            {percentage}%
          </Typography>
          <Chip
            label={percentage >= 80 ? 'Complete' : percentage >= 50 ? 'In Progress' : 'Getting Started'}
            color={getCompletionColor(percentage)}
            size="small"
          />
        </Box>
        
        <LinearProgress
          variant="determinate"
          value={percentage}
          color={getCompletionColor(percentage)}
          sx={{
            height: 8,
            borderRadius: 4,
            mb: 1,
          }}
        />
        
        <Typography variant="body2" color="text.secondary">
          {getCompletionMessage(percentage)}
        </Typography>
      </Box>

      {/* Missing Fields */}
      {missingFields.length > 0 && (
        <Box sx={{ mb: 3 }}>
          <Typography variant="h6" gutterBottom>
            Missing Information
          </Typography>
          <List dense>
            {missingFields.map((field) => (
              <ListItem
                key={field}
                sx={{
                  cursor: onFieldClick ? 'pointer' : 'default',
                  '&:hover': onFieldClick ? {
                    bgcolor: 'action.hover',
                    borderRadius: 1,
                  } : {},
                }}
                onClick={() => onFieldClick?.(field)}
              >
                <ListItemIcon>
                  {getFieldIcon(field)}
                </ListItemIcon>
                <ListItemText
                  primary={getFieldLabel(field)}
                  secondary={`Add your ${field} to improve your profile`}
                />
                {onFieldClick && (
                  <Button size="small" variant="outlined">
                    Add
                  </Button>
                )}
              </ListItem>
            ))}
          </List>
        </Box>
      )}

      {/* Suggestions */}
      {suggestions.length > 0 && (
        <Box sx={{ mb: 3 }}>
          <Typography variant="h6" gutterBottom>
            Suggestions
          </Typography>
          {suggestions.map((suggestion, index) => (
            <Alert
              key={index}
              severity="info"
              sx={{ mb: 1 }}
              icon={<CheckCircle />}
            >
              {suggestion}
            </Alert>
          ))}
        </Box>
      )}

      {/* Completion Benefits */}
      {percentage < 100 && (
        <Alert severity="success" sx={{ mt: 3 }}>
          <Typography variant="subtitle2" gutterBottom>
            Complete your profile to:
          </Typography>
          <Typography variant="body2" component="div">
            • Increase your credibility in contests<br />
            • Connect with other users<br />
            • Unlock advanced features<br />
            • Get personalized recommendations
          </Typography>
        </Alert>
      )}

      {percentage === 100 && (
        <Alert severity="success" sx={{ mt: 3 }}>
          <Typography variant="subtitle2" gutterBottom>
            🎉 Congratulations!
          </Typography>
          <Typography variant="body2">
            Your profile is complete! You now have access to all platform features.
          </Typography>
        </Alert>
      )}
    </Paper>
  )
}
