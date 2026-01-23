import React from 'react'
import { Form, Input, Button, Select, Space } from 'antd'
import { SaveOutlined, UserOutlined, EnvironmentOutlined, GlobalOutlined } from '@ant-design/icons'
import { useForm, Controller } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { z } from 'zod'
import type { ProfileFormData } from '../../types/profile.types'
import { PROFILE_VISIBILITY_OPTIONS } from '../../types/profile.types'

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
  })

  return (
    <Form layout="vertical" onFinish={handleSubmit(onSubmit)}>
      <Controller name="bio" control={control} render={({ field }) => (
        <Form.Item label="Bio" validateStatus={errors.bio ? 'error' : ''} help={errors.bio?.message}>
          <Input.TextArea {...field} rows={4} placeholder="Tell us about yourself..." disabled={loading} maxLength={500} showCount />
        </Form.Item>
      )} />

      <Controller name="location" control={control} render={({ field }) => (
        <Form.Item label="Location" validateStatus={errors.location ? 'error' : ''} help={errors.location?.message}>
          <Input {...field} prefix={<EnvironmentOutlined />} placeholder="City, Country" disabled={loading} />
        </Form.Item>
      )} />

      <Controller name="website" control={control} render={({ field }) => (
        <Form.Item label="Website" validateStatus={errors.website ? 'error' : ''} help={errors.website?.message}>
          <Input {...field} prefix={<GlobalOutlined />} placeholder="https://example.com" disabled={loading} />
        </Form.Item>
      )} />

      <Controller name="twitterUrl" control={control} render={({ field }) => (
        <Form.Item label="Twitter" validateStatus={errors.twitterUrl ? 'error' : ''} help={errors.twitterUrl?.message}>
          <Input {...field} placeholder="https://twitter.com/username" disabled={loading} />
        </Form.Item>
      )} />

      <Controller name="linkedinUrl" control={control} render={({ field }) => (
        <Form.Item label="LinkedIn" validateStatus={errors.linkedinUrl ? 'error' : ''} help={errors.linkedinUrl?.message}>
          <Input {...field} placeholder="https://linkedin.com/in/username" disabled={loading} />
        </Form.Item>
      )} />

      <Controller name="githubUrl" control={control} render={({ field }) => (
        <Form.Item label="GitHub" validateStatus={errors.githubUrl ? 'error' : ''} help={errors.githubUrl?.message}>
          <Input {...field} placeholder="https://github.com/username" disabled={loading} />
        </Form.Item>
      )} />

      <Controller name="profileVisibility" control={control} render={({ field }) => (
        <Form.Item label="Profile Visibility">
          <Select {...field} disabled={loading}>
            {PROFILE_VISIBILITY_OPTIONS.map(option => (
              <Select.Option key={option.value} value={option.value}>
                {option.label}
              </Select.Option>
            ))}
          </Select>
        </Form.Item>
      )} />

      <Form.Item>
        <Button
          type="primary"
          htmlType="submit"
          icon={<SaveOutlined />}
          loading={loading}
          disabled={!isDirty}
        >
          Save Changes
        </Button>
      </Form.Item>
    </Form>
  )
}

export default ProfileForm
