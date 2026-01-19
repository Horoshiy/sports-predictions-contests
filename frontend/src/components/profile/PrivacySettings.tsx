import React from 'react'
import {
  Box,
  Typography,
  Paper,
  FormControlLabel,
  Switch,
  Divider,
  Alert,
  Grid,
} from '@mui/material'
import { Security, Notifications } from '@mui/icons-material'
import { useForm, Controller } from 'react-hook-form'
import type { PreferencesFormData } from '../../types/profile.types'

interface PrivacySettingsProps {
  initialData?: Partial<PreferencesFormData>
  onUpdate: (data: Partial<PreferencesFormData>) => void
  loading?: boolean
}

export const PrivacySettings: React.FC<PrivacySettingsProps> = ({
  initialData,
  onUpdate,
  loading = false,
}) => {
  const { control, watch } = useForm<PreferencesFormData>({
    defaultValues: {
      emailNotifications: initialData?.emailNotifications ?? true,
      pushNotifications: initialData?.pushNotifications ?? true,
      contestNotifications: initialData?.contestNotifications ?? true,
      predictionReminders: initialData?.predictionReminders ?? true,
      weeklyDigest: initialData?.weeklyDigest ?? true,
      theme: initialData?.theme || 'light',
      language: initialData?.language || 'en',
      timezone: initialData?.timezone || 'UTC',
    },
    mode: 'onChange',
  })

  // Watch all form values to trigger updates immediately
  const formValues = watch()

  // Handle immediate updates when switches are toggled
  const handleSwitchChange = (field: keyof PreferencesFormData) => (
    event: React.ChangeEvent<HTMLInputElement>
  ) => {
    const updatedData = {
      ...formValues,
      [field]: event.target.checked,
    }
    onUpdate(updatedData)
  }

  return (
    <Paper elevation={3} sx={{ p: 4, maxWidth: 600, mx: 'auto' }}>
      <Box sx={{ mb: 3 }}>
        <Typography variant="h5" component="h2" gutterBottom>
          <Security sx={{ mr: 1, verticalAlign: 'middle' }} />
          Privacy & Notifications
        </Typography>
        <Typography variant="body2" color="text.secondary">
          Control your privacy settings and notification preferences
        </Typography>
      </Box>

      <Alert severity="info" sx={{ mb: 3 }}>
        Changes are saved automatically when you toggle settings
      </Alert>

      <Grid container spacing={3}>
        {/* Notification Settings */}
        <Grid item xs={12}>
          <Typography variant="h6" gutterBottom sx={{ display: 'flex', alignItems: 'center' }}>
            <Notifications sx={{ mr: 1 }} />
            Notification Preferences
          </Typography>
          
          <Box sx={{ ml: 4 }}>
            <Controller
              name="emailNotifications"
              control={control}
              render={({ field }) => (
                <FormControlLabel
                  control={
                    <Switch
                      {...field}
                      checked={field.value}
                      onChange={(e) => {
                        field.onChange(e)
                        handleSwitchChange('emailNotifications')(e)
                      }}
                      disabled={loading}
                    />
                  }
                  label={
                    <Box>
                      <Typography variant="body1">Email Notifications</Typography>
                      <Typography variant="caption" color="text.secondary">
                        Receive important updates via email
                      </Typography>
                    </Box>
                  }
                />
              )}
            />

            <Controller
              name="pushNotifications"
              control={control}
              render={({ field }) => (
                <FormControlLabel
                  control={
                    <Switch
                      {...field}
                      checked={field.value}
                      onChange={(e) => {
                        field.onChange(e)
                        handleSwitchChange('pushNotifications')(e)
                      }}
                      disabled={loading}
                    />
                  }
                  label={
                    <Box>
                      <Typography variant="body1">Push Notifications</Typography>
                      <Typography variant="caption" color="text.secondary">
                        Receive real-time notifications in your browser
                      </Typography>
                    </Box>
                  }
                />
              )}
            />

            <Controller
              name="contestNotifications"
              control={control}
              render={({ field }) => (
                <FormControlLabel
                  control={
                    <Switch
                      {...field}
                      checked={field.value}
                      onChange={(e) => {
                        field.onChange(e)
                        handleSwitchChange('contestNotifications')(e)
                      }}
                      disabled={loading}
                    />
                  }
                  label={
                    <Box>
                      <Typography variant="body1">Contest Updates</Typography>
                      <Typography variant="caption" color="text.secondary">
                        Get notified about contest results and new contests
                      </Typography>
                    </Box>
                  }
                />
              )}
            />

            <Controller
              name="predictionReminders"
              control={control}
              render={({ field }) => (
                <FormControlLabel
                  control={
                    <Switch
                      {...field}
                      checked={field.value}
                      onChange={(e) => {
                        field.onChange(e)
                        handleSwitchChange('predictionReminders')(e)
                      }}
                      disabled={loading}
                    />
                  }
                  label={
                    <Box>
                      <Typography variant="body1">Prediction Reminders</Typography>
                      <Typography variant="caption" color="text.secondary">
                        Remind me to make predictions before deadlines
                      </Typography>
                    </Box>
                  }
                />
              )}
            />

            <Controller
              name="weeklyDigest"
              control={control}
              render={({ field }) => (
                <FormControlLabel
                  control={
                    <Switch
                      {...field}
                      checked={field.value}
                      onChange={(e) => {
                        field.onChange(e)
                        handleSwitchChange('weeklyDigest')(e)
                      }}
                      disabled={loading}
                    />
                  }
                  label={
                    <Box>
                      <Typography variant="body1">Weekly Digest</Typography>
                      <Typography variant="caption" color="text.secondary">
                        Receive a weekly summary of your predictions and results
                      </Typography>
                    </Box>
                  }
                />
              )}
            />
          </Box>
        </Grid>

        <Grid item xs={12}>
          <Divider sx={{ my: 2 }} />
        </Grid>

        {/* Privacy Information */}
        <Grid item xs={12}>
          <Typography variant="h6" gutterBottom sx={{ display: 'flex', alignItems: 'center' }}>
            <Security sx={{ mr: 1 }} />
            Privacy Information
          </Typography>
          
          <Alert severity="info" sx={{ mb: 2 }}>
            Your profile visibility is controlled in the Profile tab. 
            These settings only affect notifications and communications.
          </Alert>

          <Box sx={{ ml: 2 }}>
            <Typography variant="body2" color="text.secondary" paragraph>
              • We never share your personal information with third parties
            </Typography>
            <Typography variant="body2" color="text.secondary" paragraph>
              • You can export or delete your data at any time
            </Typography>
            <Typography variant="body2" color="text.secondary" paragraph>
              • All communications are encrypted and secure
            </Typography>
            <Typography variant="body2" color="text.secondary">
              • You can unsubscribe from any notifications at any time
            </Typography>
          </Box>
        </Grid>
      </Grid>
    </Paper>
  )
}
