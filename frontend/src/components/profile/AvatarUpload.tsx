import React, { useState, useRef } from 'react'
import {
  Box,
  Button,
  Avatar,
  Typography,
  Paper,
  LinearProgress,
  Alert,
  IconButton,
} from '@mui/material'
import { CloudUpload, PhotoCamera, Delete } from '@mui/icons-material'
import { profileService } from '../../services/profile-service'
import type { AvatarUploadState } from '../../types/profile.types'

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
  const fileInputRef = useRef<HTMLInputElement>(null)

  const handleFileSelect = (event: React.ChangeEvent<HTMLInputElement>) => {
    const file = event.target.files?.[0]
    if (!file) return

    // Validate file
    const validation = profileService.validateAvatarFile(file)
    if (!validation.isValid) {
      setUploadState({
        isUploading: false,
        progress: 0,
        error: validation.error || 'Invalid file',
      })
      return
    }

    // Clear any previous errors
    setUploadState({
      isUploading: false,
      progress: 0,
      error: null,
    })

    // Generate preview
    profileService.generateAvatarPreview(file)
      .then(setPreviewUrl)
      .catch(() => {
        setUploadState(prev => ({
          ...prev,
          error: 'Failed to generate preview',
        }))
      })

    // Upload file
    handleUpload(file)
  }

  const handleUpload = async (file: File) => {
    setUploadState({
      isUploading: true,
      progress: 0,
      error: null,
    })

    try {
      const avatarUrl = await profileService.uploadAvatar(file)
      
      setUploadState({
        isUploading: false,
        progress: 0,
        error: null,
      })

      setPreviewUrl(null)
      onAvatarUpdate(avatarUrl)

    } catch (error) {
      setUploadState({
        isUploading: false,
        progress: 0,
        error: error instanceof Error ? error.message : 'Upload failed',
      })
      setPreviewUrl(null)
    }
  }

  const handleRemoveAvatar = () => {
    setPreviewUrl(null)
    onAvatarUpdate('')
    setUploadState({
      isUploading: false,
      progress: 0,
      error: null,
    })
  }

  const handleButtonClick = () => {
    fileInputRef.current?.click()
  }

  const displayUrl = previewUrl || currentAvatarUrl
  const hasAvatar = Boolean(displayUrl)

  return (
    <Paper elevation={2} sx={{ p: 3, textAlign: 'center' }}>
      <Typography variant="h6" gutterBottom>
        Profile Picture
      </Typography>

      <Box sx={{ mb: 3, position: 'relative', display: 'inline-block' }}>
        <Avatar
          src={displayUrl}
          sx={{
            width: size,
            height: size,
            mx: 'auto',
            fontSize: size / 3,
            border: '4px solid',
            borderColor: 'divider',
          }}
        >
          {!hasAvatar && <PhotoCamera sx={{ fontSize: size / 3 }} />}
        </Avatar>

        {hasAvatar && (
          <IconButton
            size="small"
            sx={{
              position: 'absolute',
              top: -8,
              right: -8,
              bgcolor: 'error.main',
              color: 'white',
              '&:hover': {
                bgcolor: 'error.dark',
              },
            }}
            onClick={handleRemoveAvatar}
            disabled={uploadState.isUploading}
          >
            <Delete fontSize="small" />
          </IconButton>
        )}
      </Box>

      {uploadState.isUploading && (
        <Box sx={{ mb: 2 }}>
          <LinearProgress sx={{ mb: 1 }} />
          <Typography variant="body2" color="text.secondary">
            Uploading...
          </Typography>
        </Box>
      )}

      {uploadState.error && (
        <Alert severity="error" sx={{ mb: 2 }}>
          {uploadState.error}
        </Alert>
      )}

      <Box sx={{ display: 'flex', gap: 1, justifyContent: 'center' }}>
        <Button
          variant="contained"
          startIcon={<CloudUpload />}
          onClick={handleButtonClick}
          disabled={uploadState.isUploading}
        >
          {hasAvatar ? 'Change Photo' : 'Upload Photo'}
        </Button>
      </Box>

      <input
        ref={fileInputRef}
        type="file"
        accept="image/jpeg,image/png,image/gif"
        onChange={handleFileSelect}
        style={{ display: 'none' }}
      />

      <Typography variant="caption" display="block" sx={{ mt: 2, color: 'text.secondary' }}>
        Supported formats: JPEG, PNG, GIF (max 5MB)
      </Typography>
    </Paper>
  )
}
