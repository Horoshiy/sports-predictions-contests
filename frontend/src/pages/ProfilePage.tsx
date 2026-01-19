import React, { useState, useEffect } from 'react'
import {
  Box,
  Typography,
  Container,
  Tabs,
  Tab,
  CircularProgress,
  Alert,
} from '@mui/material'
import {
  Person as PersonIcon,
  Security as SecurityIcon,
  TrendingUp as ProgressIcon,
} from '@mui/icons-material'
import { useAuth } from '../contexts/AuthContext'
import { useToast } from '../contexts/ToastContext'
import { ProfileForm } from '../components/profile/ProfileForm'
import { AvatarUpload } from '../components/profile/AvatarUpload'
import { PrivacySettings } from '../components/profile/PrivacySettings'
import { ProfileCompletion } from '../components/profile/ProfileCompletion'
import { profileService } from '../services/profile-service'
import type {
  Profile,
  UserPreferences,
  ProfileFormData,
  PreferencesFormData,
} from '../types/profile.types'

// Use the renamed interface
interface ProfileCompletionData {
  percentage: number
  missingFields: string[]
  suggestions: string[]
}

interface TabPanelProps {
  children?: React.ReactNode
  index: number
  value: number
}

const TabPanel: React.FC<TabPanelProps> = ({ children, value, index }) => (
  <div
    role="tabpanel"
    hidden={value !== index}
    id={`profile-tabpanel-${index}`}
    aria-labelledby={`profile-tab-${index}`}
  >
    {value === index && <Box sx={{ py: 3 }}>{children}</Box>}
  </div>
)

const ProfilePage: React.FC = () => {
  const { user } = useAuth()
  const { showToast } = useToast()
  
  const [activeTab, setActiveTab] = useState(0)
  const [loading, setLoading] = useState(true)
  const [saving, setSaving] = useState(false)
  
  const [profile, setProfile] = useState<Profile | null>(null)
  const [preferences, setPreferences] = useState<UserPreferences | null>(null)
  const [completion, setCompletion] = useState<ProfileCompletionData | null>(null)
  
  const [error, setError] = useState<string | null>(null)

  // Load profile data on component mount
  useEffect(() => {
    loadProfileData()
  }, [])

  const loadProfileData = async () => {
    if (!user) return

    try {
      setLoading(true)
      setError(null)

      const [profileData, preferencesData, completionData] = await Promise.all([
        profileService.getProfile().catch(() => profileService.getDefaultProfile()),
        profileService.getPreferences().catch(() => profileService.getDefaultPreferences()),
        profileService.getProfileCompletion().catch(() => ({ percentage: 0, missingFields: [], suggestions: [] })),
      ])

      setProfile(profileData as Profile)
      setPreferences(preferencesData as UserPreferences)
      setCompletion(completionData)
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : 'Failed to load profile data'
      setError(errorMessage)
      showToast(errorMessage, 'error')
    } finally {
      setLoading(false)
    }
  }

  const handleTabChange = (_event: React.SyntheticEvent, newValue: number) => {
    setActiveTab(newValue)
  }

  const handleProfileUpdate = async (profileData: ProfileFormData) => {
    try {
      setSaving(true)
      const updatedProfile = await profileService.updateProfile(profileData)
      setProfile(updatedProfile)
      
      // Refresh completion data
      const completionData = await profileService.getProfileCompletion()
      setCompletion(completionData)
      
      showToast('Profile updated successfully', 'success')
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : 'Failed to update profile'
      showToast(errorMessage, 'error')
    } finally {
      setSaving(false)
    }
  }

  const handleAvatarUpdate = async (avatarUrl: string) => {
    try {
      // Update profile with new avatar URL
      if (profile) {
        const updatedProfileData: ProfileFormData = {
          bio: profile.bio,
          location: profile.location,
          website: profile.website,
          twitterUrl: profile.twitterUrl,
          linkedinUrl: profile.linkedinUrl,
          githubUrl: profile.githubUrl,
          profileVisibility: profile.profileVisibility,
        }
        
        const updatedProfile = await profileService.updateProfile(updatedProfileData)
        setProfile({ ...updatedProfile, avatarUrl })
        
        // Refresh completion data
        const completionData = await profileService.getProfileCompletion()
        setCompletion(completionData)
        
        showToast('Avatar updated successfully', 'success')
      }
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : 'Failed to update avatar'
      showToast(errorMessage, 'error')
    }
  }

  const handlePreferencesUpdate = async (preferencesData: Partial<PreferencesFormData>) => {
    try {
      const fullPreferencesData: PreferencesFormData = {
        emailNotifications: preferencesData.emailNotifications ?? true,
        pushNotifications: preferencesData.pushNotifications ?? true,
        contestNotifications: preferencesData.contestNotifications ?? true,
        predictionReminders: preferencesData.predictionReminders ?? true,
        weeklyDigest: preferencesData.weeklyDigest ?? true,
        theme: preferencesData.theme ?? 'light',
        language: preferencesData.language ?? 'en',
        timezone: preferencesData.timezone ?? 'UTC',
      }
      
      const updatedPreferences = await profileService.updatePreferences(fullPreferencesData)
      setPreferences(updatedPreferences)
      
      showToast('Preferences updated successfully', 'success')
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : 'Failed to update preferences'
      showToast(errorMessage, 'error')
    }
  }

  const handleFieldClick = (field: string) => {
    // Switch to appropriate tab based on field
    switch (field) {
      case 'avatar':
        setActiveTab(0) // Profile tab
        break
      case 'bio':
      case 'location':
      case 'website':
      case 'twitter':
      case 'linkedin':
      case 'github':
        setActiveTab(0) // Profile tab
        break
      default:
        setActiveTab(0) // Default to profile tab
    }
  }

  if (loading) {
    return (
      <Container maxWidth="lg" sx={{ py: 4 }}>
        <Box display="flex" justifyContent="center" alignItems="center" minHeight="400px">
          <CircularProgress size={60} />
        </Box>
      </Container>
    )
  }

  if (error) {
    return (
      <Container maxWidth="lg" sx={{ py: 4 }}>
        <Alert severity="error" sx={{ mb: 3 }}>
          {error}
        </Alert>
      </Container>
    )
  }

  return (
    <Container maxWidth="lg" sx={{ py: 4 }}>
      <Typography variant="h4" component="h1" gutterBottom>
        Profile Management
      </Typography>
      
      <Typography variant="body1" color="text.secondary" paragraph>
        Manage your profile information, privacy settings, and preferences
      </Typography>

      <Box sx={{ borderBottom: 1, borderColor: 'divider', mb: 3 }}>
        <Tabs value={activeTab} onChange={handleTabChange} aria-label="profile tabs">
          <Tab
            icon={<PersonIcon />}
            label="Profile"
            id="profile-tab-0"
            aria-controls="profile-tabpanel-0"
          />
          <Tab
            icon={<SecurityIcon />}
            label="Privacy & Notifications"
            id="profile-tab-1"
            aria-controls="profile-tabpanel-1"
          />
          <Tab
            icon={<ProgressIcon />}
            label="Completion"
            id="profile-tab-2"
            aria-controls="profile-tabpanel-2"
          />
        </Tabs>
      </Box>

      <TabPanel value={activeTab} index={0}>
        <Box sx={{ display: 'flex', flexDirection: 'column', gap: 4 }}>
          {/* Avatar Upload */}
          <AvatarUpload
            currentAvatarUrl={profile?.avatarUrl}
            onAvatarUpdate={handleAvatarUpdate}
          />
          
          {/* Profile Form */}
          <ProfileForm
            initialData={profile ? {
              bio: profile.bio,
              location: profile.location,
              website: profile.website,
              twitterUrl: profile.twitterUrl,
              linkedinUrl: profile.linkedinUrl,
              githubUrl: profile.githubUrl,
              profileVisibility: profile.profileVisibility,
            } : undefined}
            onSubmit={handleProfileUpdate}
            loading={saving}
          />
        </Box>
      </TabPanel>

      <TabPanel value={activeTab} index={1}>
        <PrivacySettings
          initialData={preferences ? {
            emailNotifications: preferences.emailNotifications,
            pushNotifications: preferences.pushNotifications,
            contestNotifications: preferences.contestNotifications,
            predictionReminders: preferences.predictionReminders,
            weeklyDigest: preferences.weeklyDigest,
            theme: preferences.theme,
            language: preferences.language,
            timezone: preferences.timezone,
          } : undefined}
          onUpdate={handlePreferencesUpdate}
          loading={saving}
        />
      </TabPanel>

      <TabPanel value={activeTab} index={2}>
        {completion && (
          <ProfileCompletion
            completion={completion}
            onFieldClick={handleFieldClick}
          />
        )}
      </TabPanel>
    </Container>
  )
}

export default ProfilePage
