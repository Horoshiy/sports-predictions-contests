import React from 'react'
import { Modal, Form, Input, Select, Button } from 'antd'
import { useForm, Controller } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { leagueSchema, type LeagueFormData, generateSlug } from '../../utils/sports-validation'
import { useSports } from '../../hooks/use-sports'
import type { League } from '../../types/sports.types'

interface LeagueFormProps {
  open: boolean
  onClose: () => void
  onSubmit: (data: LeagueFormData) => void
  league?: League | null
  loading?: boolean
}

export const LeagueForm: React.FC<LeagueFormProps> = ({ open, onClose, onSubmit, league, loading = false }) => {
  const isEditing = !!league
  const { data: sportsData } = useSports({ pagination: { page: 1, limit: 100 }, activeOnly: true })

  const defaultValues = React.useMemo(() => ({
    sportId: league?.sportId || 0,
    name: league?.name || '',
    slug: league?.slug || '',
    country: league?.country || '',
    season: league?.season || '',
  }), [league])

  const { control, handleSubmit, reset, watch, setValue, formState: { errors, isValid } } = useForm<LeagueFormData>({
    resolver: zodResolver(leagueSchema),
    defaultValues,
    mode: 'onChange',
  })

  const name = watch('name')
  const slug = watch('slug')

  React.useEffect(() => {
    if (!isEditing && name && !slug) {
      setValue('slug', generateSlug(name))
    }
  }, [name, slug, isEditing, setValue])

  React.useEffect(() => {
    reset(defaultValues)
  }, [defaultValues, reset])

  const handleClose = () => {
    reset()
    onClose()
  }

  return (
    <Modal
      open={open}
      title={isEditing ? 'Edit League' : 'Create League'}
      onCancel={handleClose}
      footer={[
        <Button key="cancel" onClick={handleClose} disabled={loading}>Cancel</Button>,
        <Button key="submit" type="primary" onClick={handleSubmit(onSubmit)} disabled={loading || !isValid} loading={loading}>
          {isEditing ? 'Update' : 'Create'}
        </Button>,
      ]}
    >
      <Form layout="vertical">
        <Controller name="sportId" control={control} render={({ field }) => (
          <Form.Item label="Sport" required validateStatus={errors.sportId ? 'error' : ''} help={errors.sportId?.message}>
            <Select {...field} disabled={loading}>
              {sportsData?.sports?.map(s => <Select.Option key={s.id} value={s.id}>{s.name}</Select.Option>)}
            </Select>
          </Form.Item>
        )} />
        <Controller name="name" control={control} render={({ field }) => (
          <Form.Item label="Name" required validateStatus={errors.name ? 'error' : ''} help={errors.name?.message}>
            <Input {...field} disabled={loading} />
          </Form.Item>
        )} />
        <Controller name="slug" control={control} render={({ field }) => (
          <Form.Item label="Slug" validateStatus={errors.slug ? 'error' : ''} help={errors.slug?.message || 'Auto-generated'}>
            <Input {...field} disabled={loading} />
          </Form.Item>
        )} />
        <Controller name="country" control={control} render={({ field }) => (
          <Form.Item label="Country" validateStatus={errors.country ? 'error' : ''} help={errors.country?.message}>
            <Input {...field} disabled={loading} />
          </Form.Item>
        )} />
        <Controller name="season" control={control} render={({ field }) => (
          <Form.Item label="Season" validateStatus={errors.season ? 'error' : ''} help={errors.season?.message}>
            <Input {...field} disabled={loading} />
          </Form.Item>
        )} />
      </Form>
    </Modal>
  )
}

export default LeagueForm
