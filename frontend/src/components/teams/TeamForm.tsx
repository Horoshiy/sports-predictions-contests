import React, { useEffect } from 'react'
import { Modal, Form, Input, Button, InputNumber } from 'antd'
import { useForm, Controller } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { teamSchema, type TeamSchemaType } from '../../utils/team-validation'
import { useCreateTeam, useUpdateTeam } from '../../hooks/use-teams'
import type { Team } from '../../types/team.types'

interface TeamFormProps {
  open: boolean
  onClose: () => void
  team?: Team | null
}

export const TeamForm: React.FC<TeamFormProps> = ({ open, onClose, team }) => {
  const isEdit = !!team
  const createTeamMutation = useCreateTeam()
  const updateTeamMutation = useUpdateTeam()

  const { control, handleSubmit, reset, formState: { errors } } = useForm<TeamSchemaType>({
    resolver: zodResolver(teamSchema),
    defaultValues: { name: '', description: '', maxMembers: 10 },
  })

  useEffect(() => {
    if (team) {
      reset({ name: team.name, description: team.description, maxMembers: team.maxMembers })
    } else {
      reset({ name: '', description: '', maxMembers: 10 })
    }
  }, [team, reset])

  const onSubmit = (data: TeamSchemaType) => {
    if (isEdit && team) {
      updateTeamMutation.mutate({ id: team.id, name: data.name, description: data.description, maxMembers: data.maxMembers }, { onSuccess: onClose })
    } else {
      createTeamMutation.mutate({ name: data.name, description: data.description, maxMembers: data.maxMembers }, { onSuccess: onClose })
    }
  }

  const loading = createTeamMutation.isPending || updateTeamMutation.isPending

  return (
    <Modal
      open={open}
      title={isEdit ? 'Edit Team' : 'Create Team'}
      onCancel={onClose}
      footer={[
        <Button key="cancel" onClick={onClose} disabled={loading}>Cancel</Button>,
        <Button key="submit" type="primary" onClick={handleSubmit(onSubmit)} loading={loading}>
          {isEdit ? 'Update' : 'Create'}
        </Button>,
      ]}
    >
      <Form layout="vertical">
        <Controller name="name" control={control} render={({ field }) => (
          <Form.Item label="Team Name" required validateStatus={errors.name ? 'error' : ''} help={errors.name?.message}>
            <Input {...field} disabled={loading} />
          </Form.Item>
        )} />
        <Controller name="description" control={control} render={({ field }) => (
          <Form.Item label="Description" validateStatus={errors.description ? 'error' : ''} help={errors.description?.message}>
            <Input.TextArea {...field} rows={3} disabled={loading} />
          </Form.Item>
        )} />
        <Controller name="maxMembers" control={control} render={({ field }) => (
          <Form.Item label="Max Members" validateStatus={errors.maxMembers ? 'error' : ''} help={errors.maxMembers?.message}>
            <InputNumber {...field} min={1} max={100} disabled={loading} style={{ width: '100%' }} />
          </Form.Item>
        )} />
      </Form>
    </Modal>
  )
}

export default TeamForm
