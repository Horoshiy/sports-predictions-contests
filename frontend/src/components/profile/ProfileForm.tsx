import React from 'react'
import {
  Box,
  TextField,
  Button,
  Typography,
  Paper,
  FormControl,
  InputLabel,
  Select,
  MenuItem,
  Grid,
} from '@mui/material'
import { Save, Person, LocationOn, Language, Public } from '@mui/icons-material'
import { useForm, Controller } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { z } from 'zod'
import type { ProfileFormData } from '../../types/profile.types'
import { PROFILE_VISIBILITY_OPTIONS } from '../../types/profile.types'

// Validation schema
const profileSchema = z.object({
  bio: z.string().max(500, 'Bio must be less than 500 characters').optional(),
  location: z.string().max(100, 'Location must be less than 100 characters').optional(),
  website: z.string().url('Invalid website URL').or(z.literal('')).optional(),
  twitterUrl: z.string().url('Invalid Twitter URL').or(z.literal('')).optional(),
  linkedinUrl: z.string().url('Invalid LinkedIn URL').or(z.literal('')).optional(),
  githubUrl: z.string().url('Invalid GitHub URL').or(z.literal('')).optional(),
  profileVisibility: z.enum(['public', 'friends', 'private']),
})

interface ProfileFormProps {
  initialData?: Partial<ProfileFormData>
  onSubmit: (data: ProfileFormData) => void
  loading?: boolean
}

export const ProfileForm: React.FC<ProfileFormProps> = ({
  initialData,
  onSubmit,
  loading = false,
}) => {
  const {
    control,
    handleSubmit,
    formState: { errors, isDirty },
  } = useForm<ProfileFormData>({
    resolver: zodResolver(profileSchema),
    defaultValues: {
      bio: initialData?.bio || '',
      location: initialData?.location || '',
      website: initialData?.website || '',
      twitterUrl: initialData?.twitterUrl || '',
      linkedinUrl: initialData?.linkedinUrl || '',
      githubUrl: initialData?.githubUrl || '',
      profileVisibility: initialData?.profileVisibility || 'public',
    },
    mode: 'onBlur',
  })

  const handleFormSubmit = (data: ProfileFormData) => {
    onSubmit(data)
  }

  return (
    <Paper elevation={3} sx={{ p: 4, maxWidth: 600, mx: 'auto' }}>
      <Box sx={{ mb: 3 }}>
        <Typography variant="h5" component="h2" gutterBottom>
          <Person sx={{ mr: 1, verticalAlign: 'middle' }} />
          Profile Information
        </Typography>
        <Typography variant="body2" color="text.secondary">
          Update your profile information to personalize your account
        </Typography>
      </Box>

      <Box component="form" onSubmit={handleSubmit(handleFormSubmit)} noValidate>
        <Grid container spacing={3}>
          {/* Bio */}
          <Grid item xs={12}>
            <Controller
              name="bio"
              control={control}
              render={({ field }) => (
                <TextField
                  {...field}
                  fullWidth
                  label="Bio"
                  placeholder="Tell others about yourself..."
                  multiline
                  rows={4}
                  error={!!errors.bio}
                  helperText={errors.bio?.message || `${field.value?.length || 0}/500 characters`}
                  disabled={loading}
                  InputProps={{
                    startAdornment: (
                      <Person sx={{ color: 'action.active', mr: 1, my: 0.5 }} />
                    ),
                  }}
                />
              )}
            />
          </Grid>

          {/* Location */}
          <Grid item xs={12} sm={6}>
            <Controller
              name="location"
              control={control}
              render={({ field }) => (
                <TextField
                  {...field}
                  fullWidth
                  label="Location"
                  placeholder="City, Country"
                  error={!!errors.location}
                  helperText={errors.location?.message}
                  disabled={loading}
                  InputProps={{
                    startAdornment: (
                      <LocationOn sx={{ color: 'action.active', mr: 1, my: 0.5 }} />
                    ),
                  }}
                />
              )}
            />
          </Grid>

          {/* Profile Visibility */}
          <Grid item xs={12} sm={6}>
            <Controller
              name="profileVisibility"
              control={control}
              render={({ field }) => (
                <FormControl fullWidth error={!!errors.profileVisibility}>
                  <InputLabel>Profile Visibility</InputLabel>
                  <Select
                    {...field}
                    label="Profile Visibility"
                    disabled={loading}
                    startAdornment={
                      <Public sx={{ color: 'action.active', mr: 1 }} />
                    }
                  >
                    {PROFILE_VISIBILITY_OPTIONS.map((option) => (
                      <MenuItem key={option.value} value={option.value}>
                        {option.label}
                      </MenuItem>
                    ))}
                  </Select>
                </FormControl>
              )}
            />
          </Grid>

          {/* Website */}
          <Grid item xs={12} sm={6}>
            <Controller
              name="website"
              control={control}
              render={({ field }) => (
                <TextField
                  {...field}
                  fullWidth
                  label="Website"
                  placeholder="https://yourwebsite.com"
                  error={!!errors.website}
                  helperText={errors.website?.message}
                  disabled={loading}
                  InputProps={{
                    startAdornment: (
                      <Language sx={{ color: 'action.active', mr: 1, my: 0.5 }} />
                    ),
                  }}
                />
              )}
            />
          </Grid>

          {/* Twitter URL */}
          <Grid item xs={12} sm={6}>
            <Controller
              name="twitterUrl"
              control={control}
              render={({ field }) => (
                <TextField
                  {...field}
                  fullWidth
                  label="Twitter"
                  placeholder="https://twitter.com/username"
                  error={!!errors.twitterUrl}
                  helperText={errors.twitterUrl?.message}
                  disabled={loading}
                />
              )}
            />
          </Grid>

          {/* LinkedIn URL */}
          <Grid item xs={12} sm={6}>
            <Controller
              name="linkedinUrl"
              control={control}
              render={({ field }) => (
                <TextField
                  {...field}
                  fullWidth
                  label="LinkedIn"
                  placeholder="https://linkedin.com/in/username"
                  error={!!errors.linkedinUrl}
                  helperText={errors.linkedinUrl?.message}
                  disabled={loading}
                />
              )}
            />
          </Grid>

          {/* GitHub URL */}
          <Grid item xs={12} sm={6}>
            <Controller
              name="githubUrl"
              control={control}
              render={({ field }) => (
                <TextField
                  {...field}
                  fullWidth
                  label="GitHub"
                  placeholder="https://github.com/username"
                  error={!!errors.githubUrl}
                  helperText={errors.githubUrl?.message}
                  disabled={loading}
                />
              )}
            />
          </Grid>
        </Grid>

        {/* Submit Button */}
        <Box sx={{ mt: 4, display: 'flex', justifyContent: 'flex-end' }}>
          <Button
            type="submit"
            variant="contained"
            size="large"
            disabled={loading || !isDirty}
            startIcon={<Save />}
            sx={{ minWidth: 120 }}
          >
            {loading ? 'Saving...' : 'Save Profile'}
          </Button>
        </Box>
      </Box>
    </Paper>
  )
}
