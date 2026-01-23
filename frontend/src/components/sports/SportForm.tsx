import React from 'react'
import { Modal, Form, Input, Button } from 'antd'
import { useForm, Controller } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { sportSchema, type SportFormData, generateSlug } from '../../utils/sports-validation'
import type { Sport } from '../../types/sports.types'

interface SportFormProps {
  open: boolean
  onClose: () => void
  onSubmit: (data: SportFormData) => void
  sport?: Sport | null
  loading?: boolean
}

export const SportForm: React.FC<SportFormProps> = ({ open, onClose, onSubmit, sport, loading = false }) => {
  const isEditing = !!sport

  const defaultValues = React.useMemo(() => ({
    name: sport?.name || '',
    slug: sport?.slug || '',
    description: sport?.description || '',
    iconUrl: sport?.iconUrl || '',
  }), [sport])

  const { control, handleSubmit, reset, watch, setValue, formState: { errors, isValid } } = useForm<SportFormData>({
    resolver: zodResolver(sportSchema),
    defaultValues,
    mode: 'onChange',
  })

  const name = watch('name')
  const slugTouched = React.useRef(false)

  React.useEffect(() => {
    if (!isEditing && name && !slugTouched.current) {
      setValue('slug', generateSlug(name))
    }
  }, [name, isEditing, setValue])

  React.useEffect(() => {
    reset(defaultValues)
  }, [defaultValues, reset])

  const handleClose = () => {
    slugTouched.current = false
    reset()
    onClose()
  }

  return (
    <Modal
      open={open}
      title={isEditing ? 'Edit Sport' : 'Create Sport'}
      onCancel={handleClose}
      footer={[
        <Button key="cancel" onClick={handleClose} disabled={loading}>
          Cancel
        </Button>,
        <Button key="submit" type="primary" onClick={handleSubmit(onSubmit)} disabled={loading || !isValid} loading={loading}>
          {isEditing ? 'Update' : 'Create'}
        </Button>,
      ]}
    >
      <Form layout="vertical">
        <Controller
          name="name"
          control={control}
          render={({ field }) => (
            <Form.Item label="Name" required validateStatus={errors.name ? 'error' : ''} help={errors.name?.message}>
              <Input {...field} disabled={loading} />
            </Form.Item>
          )}
        />
        <Controller
          name="slug"
          control={control}
          render={({ field }) => (
            <Form.Item label="Slug" validateStatus={errors.slug ? 'error' : ''} help={errors.slug?.message || 'Auto-generated from name'}>
              <Input {...field} disabled={loading} onChange={(e) => { slugTouched.current = true; field.onChange(e) }} />
            </Form.Item>
          )}
        />
        <Controller
          name="description"
          control={control}
          render={({ field }) => (
            <Form.Item label="Description" validateStatus={errors.description ? 'error' : ''} help={errors.description?.message}>
              <Input.TextArea {...field} rows={3} disabled={loading} />
            </Form.Item>
          )}
        />
        <Controller
          name="iconUrl"
          control={control}
          render={({ field }) => (
            <Form.Item label="Icon URL" validateStatus={errors.iconUrl ? 'error' : ''} help={errors.iconUrl?.message}>
              <Input {...field} disabled={loading} />
            </Form.Item>
          )}
        />
      </Form>
    </Modal>
  )
}

export default SportForm
