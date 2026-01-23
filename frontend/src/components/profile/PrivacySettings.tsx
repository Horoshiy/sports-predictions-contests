import React from 'react'
import { Card, Switch, Space, Typography, Divider } from 'antd'
import { SecurityScanOutlined, BellOutlined } from '@ant-design/icons'
import { useForm, Controller } from 'react-hook-form'
import { debounce } from '../../utils/debounce'
import type { PreferencesFormData } from '../../types/profile.types'

const { Title, Text } = Typography

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
    },
  })

  React.useEffect(() => {
    let isMounted = true
    
    const debouncedUpdate = debounce(async (value: Partial<PreferencesFormData>) => {
      if (isMounted) {
        await onUpdate(value)
      }
    }, 500) // Wait 500ms after last change
    
    const subscription = watch((value) => {
      debouncedUpdate(value as Partial<PreferencesFormData>)
    })
    
    return () => {
      isMounted = false
      debouncedUpdate.cancel()
      subscription.unsubscribe()
    }
  }, [watch, onUpdate])

  return (
    <Space direction="vertical" size="large" style={{ width: '100%' }}>
      <Card>
        <Space direction="vertical" size="middle" style={{ width: '100%' }}>
          <Space>
            <BellOutlined style={{ fontSize: 20 }} />
            <Title level={5} style={{ margin: 0 }}>Notifications</Title>
          </Space>
          <Divider style={{ margin: '8px 0' }} />
          
          <Controller name="emailNotifications" control={control} render={({ field }) => (
            <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
              <div>
                <Text strong>Email Notifications</Text>
                <br />
                <Text type="secondary" style={{ fontSize: 12 }}>Receive notifications via email</Text>
              </div>
              <Switch {...field} checked={field.value} disabled={loading} />
            </div>
          )} />

          <Controller name="pushNotifications" control={control} render={({ field }) => (
            <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
              <div>
                <Text strong>Push Notifications</Text>
                <br />
                <Text type="secondary" style={{ fontSize: 12 }}>Receive push notifications</Text>
              </div>
              <Switch {...field} checked={field.value} disabled={loading} />
            </div>
          )} />

          <Controller name="contestNotifications" control={control} render={({ field }) => (
            <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
              <div>
                <Text strong>Contest Updates</Text>
                <br />
                <Text type="secondary" style={{ fontSize: 12 }}>Get notified about contest updates</Text>
              </div>
              <Switch {...field} checked={field.value} disabled={loading} />
            </div>
          )} />

          <Controller name="predictionReminders" control={control} render={({ field }) => (
            <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
              <div>
                <Text strong>Prediction Reminders</Text>
                <br />
                <Text type="secondary" style={{ fontSize: 12 }}>Remind me to make predictions</Text>
              </div>
              <Switch {...field} checked={field.value} disabled={loading} />
            </div>
          )} />

          <Controller name="weeklyDigest" control={control} render={({ field }) => (
            <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
              <div>
                <Text strong>Weekly Digest</Text>
                <br />
                <Text type="secondary" style={{ fontSize: 12 }}>Receive weekly summary emails</Text>
              </div>
              <Switch {...field} checked={field.value} disabled={loading} />
            </div>
          )} />
        </Space>
      </Card>
    </Space>
  )
}

export default PrivacySettings
