import React from 'react'
import { Modal, Form, Input, Select, Button, Avatar } from 'antd'
import { useForm, Controller } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { teamSchema, type TeamFormData, generateSlug } from '../../utils/sports-validation'
import { useSports } from '../../hooks/use-sports'
import type { Team } from '../../types/sports.types'

interface TeamFormProps {
  open: boolean
  onClose: () => void
  onSubmit: (data: TeamFormData) => void
  team?: Team | null
  loading?: boolean
}

export const TeamForm: React.FC<TeamFormProps> = ({ open, onClose, onSubmit, team, loading = false }) => {
  const isEditing = !!team
  const { data: sportsData } = useSports({ pagination: { page: 1, limit: 100 }, activeOnly: true })

  const defaultValues = React.useMemo(() => ({
    sportId: team?.sportId || undefined,
    name: team?.name || '',
    slug: team?.slug || '',
    shortName: team?.shortName || '',
    logoUrl: team?.logoUrl || '',
    country: team?.country || '',
  }), [team])

  const { control, handleSubmit, reset, watch, setValue, formState: { errors, isValid } } = useForm<TeamFormData>({
    resolver: zodResolver(teamSchema),
    defaultValues,
    mode: 'onChange',
  })

  const name = watch('name')
  const slug = watch('slug')
  const logoUrl = watch('logoUrl')

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
      title={isEditing ? 'Edit Team' : 'Create Team'}
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
            <Select {...field} disabled={loading} placeholder="Select a sport">
              {sportsData?.sports?.map(s => <Select.Option key={s.id} value={s.id}>{s.name}</Select.Option>)}
            </Select>
          </Form.Item>
        )} />
        <Controller name="name" control={control} render={({ field }) => (
          <Form.Item label="Name" required validateStatus={errors.name ? 'error' : ''} help={errors.name?.message}>
            <Input {...field} disabled={loading} />
          </Form.Item>
        )} />
        <Controller name="shortName" control={control} render={({ field }) => (
          <Form.Item label="Short Name" validateStatus={errors.shortName ? 'error' : ''} help={errors.shortName?.message}>
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
        <Controller name="logoUrl" control={control} render={({ field }) => (
          <Form.Item label="Logo URL" validateStatus={errors.logoUrl ? 'error' : ''} help={errors.logoUrl?.message}>
            <Input {...field} disabled={loading} addonAfter={logoUrl && <Avatar src={logoUrl} size="small" />} />
          </Form.Item>
        )} />
      </Form>
    </Modal>
  )
}

export default TeamForm
