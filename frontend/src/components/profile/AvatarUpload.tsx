import React, { useState } from 'react'
import { Upload, Avatar, Button, Progress, Alert, Space, Typography } from 'antd'
import { CloudUploadOutlined, CameraOutlined, DeleteOutlined } from '@ant-design/icons'
import { profileService } from '../../services/profile-service'
import { showSuccess, showError } from '../../utils/notification'
import { MAX_AVATAR_SIZE_MB, ALLOWED_IMAGE_FORMATS } from '../../utils/constants'
import type { AvatarUploadState } from '../../types/profile.types'

const { Text } = Typography

interface AvatarUploadProps {
  currentAvatarUrl?: string
  onAvatarUpdate: (avatarUrl: string) => void
  size?: number
}

export const AvatarUpload: React.FC<AvatarUploadProps> = ({
  currentAvatarUrl,
  onAvatarUpdate,
  size = 120,
}) => {
  const [uploadState, setUploadState] = useState<AvatarUploadState>({
    isUploading: false,
    progress: 0,
    error: null,
  })
  const [previewUrl, setPreviewUrl] = useState<string | null>(null)

  const handleUpload = async (file: File) => {
    // Validate file size
    const maxSizeBytes = MAX_AVATAR_SIZE_MB * 1024 * 1024
    if (file.size > maxSizeBytes) {
      showError(`File size must be less than ${MAX_AVATAR_SIZE_MB}MB`)
      return false
    }
    
    // Validate file type
    if (!ALLOWED_IMAGE_FORMATS.includes(file.type)) {
      const allowedFormats = ALLOWED_IMAGE_FORMATS.map(f => f.split('/')[1].toUpperCase()).join(', ')
      showError(`Only ${allowedFormats} formats are allowed`)
      return false
    }
    
    setUploadState({ isUploading: true, progress: 0, error: null })
    
    try {
      const reader = new FileReader()
      reader.onloadend = () => setPreviewUrl(reader.result as string)
      reader.readAsDataURL(file)

      const avatarUrl = await profileService.uploadAvatar(file)

      onAvatarUpdate(avatarUrl)
      setUploadState({ isUploading: false, progress: 100, error: null })
      showSuccess('Avatar uploaded successfully')
    } catch (error: any) {
      const errorMessage = error?.message || 'Failed to upload avatar'
      setUploadState({ isUploading: false, progress: 0, error: errorMessage })
      showError(errorMessage)
    }
    
    return false
  }

  const handleDelete = async () => {
    try {
      onAvatarUpdate('')
      setPreviewUrl(null)
      showSuccess('Avatar removed successfully')
    } catch (error: any) {
      showError(error?.message || 'Failed to remove avatar')
    }
  }

  return (
    <Space direction="vertical" align="center" size="middle">
      <Avatar size={size} src={previewUrl || currentAvatarUrl} icon={<CameraOutlined />} />
      
      {uploadState.isUploading && <Progress percent={uploadState.progress} />}
      {uploadState.error && <Alert message={uploadState.error} type="error" showIcon closable />}

      <Space>
        <Upload
          beforeUpload={handleUpload}
          showUploadList={false}
          accept="image/*"
        >
          <Button icon={<CloudUploadOutlined />} loading={uploadState.isUploading}>
            Upload Avatar
          </Button>
        </Upload>
        {(currentAvatarUrl || previewUrl) && (
          <Button danger icon={<DeleteOutlined />} onClick={handleDelete}>
            Remove
          </Button>
        )}
      </Space>
      <Text type="secondary" style={{ fontSize: 12 }}>
        Max size: {MAX_AVATAR_SIZE_MB}MB. Formats: {ALLOWED_IMAGE_FORMATS.map(f => f.split('/')[1].toUpperCase()).join(', ')}
      </Text>
    </Space>
  )
}

export default AvatarUpload
