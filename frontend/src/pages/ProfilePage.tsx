import React, { useState, useEffect } from 'react'
import { Space, Typography, Tabs, Spin, Alert } from 'antd'
import { UserOutlined, SafetyOutlined, RiseOutlined } from '@ant-design/icons'
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
  UpdateProfileRequest,
  UpdatePreferencesRequest,
} from '../types/profile.types'

const { Title } = Typography

interface ProfileCompletionData {
  percentage: number
  missingFields: string[]
  suggestions: string[]
}

const ProfilePage: React.FC = () => {
  const { user } = useAuth()
  const { showToast } = useToast()
  
  const [activeTab, setActiveTab] = useState('profile')
  const [loading, setLoading] = useState(true)
  const [saving, setSaving] = useState(false)
  
  const [profile, setProfile] = useState<Profile | null>(null)
  const [preferences, setPreferences] = useState<UserPreferences | null>(null)
  const [completion, setCompletion] = useState<ProfileCompletionData | null>(null)
  
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    loadProfileData()
  }, [])

  const loadProfileData = async () => {
    if (!user) return

    try {
      setLoading(true)
      setError(null)

      const [profileData, preferencesData, completionData] = await Promise.all([
        profileService.getProfile(),
        profileService.getPreferences(),
        profileService.getProfileCompletion(),
      ])

      setProfile(profileData)
      setPreferences(preferencesData)
      setCompletion(completionData)
    } catch (err: any) {
      setError(err.message || 'Failed to load profile data')
    } finally {
      setLoading(false)
    }
  }

  const handleProfileUpdate = async (data: ProfileFormData) => {
    if (!user) return

    try {
      setSaving(true)
      await profileService.updateProfile(data)
      showToast('Profile updated successfully', 'success')
      await loadProfileData()
    } catch (err: any) {
      showToast(err.message || 'Failed to update profile', 'error')
    } finally {
      setSaving(false)
    }
  }

  const handleAvatarUpdate = async (avatarUrl: string) => {
    if (!user) return

    try {
      // Avatar is already uploaded via uploadAvatar service call
      // Just reload profile data to get the updated avatar
      await loadProfileData()
      showToast('Avatar updated successfully', 'success')
    } catch (err: any) {
      showToast(err.message || 'Failed to update avatar', 'error')
    }
  }

  const handlePreferencesUpdate = async (data: Partial<PreferencesFormData>) => {
    if (!user) return

    try {
      // Merge with existing preferences to ensure all required fields are present
      const fullData: PreferencesFormData = {
        ...preferences,
        ...data,
      } as PreferencesFormData

      await profileService.updatePreferences(fullData)
      setPreferences(prev => ({ ...prev, ...data } as UserPreferences))
    } catch (err: any) {
      showToast(err.message || 'Failed to update preferences', 'error')
    }
  }

  if (loading) {
    return (
      <div style={{ textAlign: 'center', padding: '64px 0' }}>
        <Spin size="large" />
      </div>
    )
  }

  if (error) {
    return <Alert message="Error" description={error} type="error" showIcon />
  }

  return (
    <Space direction="vertical" size="large" style={{ width: '100%', padding: '24px' }}>
      <Title level={2}>Profile Settings</Title>

      <Tabs
        activeKey={activeTab}
        onChange={setActiveTab}
        items={[
          {
            key: 'profile',
            label: <span><UserOutlined /> Profile</span>,
            children: (
              <Space direction="vertical" size="large" style={{ width: '100%' }}>
                <AvatarUpload
                  currentAvatarUrl={profile?.avatarUrl}
                  onAvatarUpdate={handleAvatarUpdate}
                />
                <ProfileForm
                  initialData={profile || undefined}
                  onSubmit={handleProfileUpdate}
                  loading={saving}
                />
              </Space>
            ),
          },
          {
            key: 'privacy',
            label: <span><SafetyOutlined /> Privacy & Notifications</span>,
            children: (
              <PrivacySettings
                initialData={preferences || undefined}
                onUpdate={handlePreferencesUpdate}
                loading={saving}
              />
            ),
          },
          {
            key: 'completion',
            label: <span><RiseOutlined /> Profile Completion</span>,
            children: completion && (
              <ProfileCompletion
                completion={completion}
                onFieldClick={(field) => {
                  setActiveTab('profile')
                  showToast(`Please complete the ${field} field`, 'info')
                }}
              />
            ),
          },
        ]}
      />
    </Space>
  )
}

export default ProfilePage
