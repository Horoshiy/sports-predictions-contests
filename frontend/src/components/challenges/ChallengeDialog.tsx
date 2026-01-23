import React from 'react'
import { Modal, Form, Input, Select, Button, Alert } from 'antd'
import { useForm, Controller } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { z } from 'zod'
import type { ChallengeFormData } from '../../types/challenge.types'

const challengeSchema = z.object({
  opponentId: z.number().min(1, 'Please select an opponent'),
  eventId: z.number().min(1, 'Please select an event'),
  message: z.string().max(500, 'Message cannot exceed 500 characters').optional(),
})

interface ChallengeDialogProps {
  open: boolean
  onClose: () => void
  onSubmit: (data: ChallengeFormData) => void
  availableOpponents?: Array<{ id: number; name: string }>
  availableEvents?: Array<{ id: number; title: string; startDate: string }>
  loading?: boolean
  error?: string | null
}

export const ChallengeDialog: React.FC<ChallengeDialogProps> = ({
  open,
  onClose,
  onSubmit,
  availableOpponents = [],
  availableEvents = [],
  loading = false,
  error = null,
}) => {
  const {
    control,
    handleSubmit,
    reset,
    formState: { errors, isValid },
  } = useForm<ChallengeFormData>({
    resolver: zodResolver(challengeSchema),
    defaultValues: {
      opponentId: 0,
      eventId: 0,
      message: '',
    },
  })

  const handleClose = () => {
    reset()
    onClose()
  }

  return (
    <Modal
      open={open}
      title="Create Challenge"
      onCancel={handleClose}
      footer={[
        <Button key="cancel" onClick={handleClose} disabled={loading}>Cancel</Button>,
        <Button key="submit" type="primary" onClick={handleSubmit(onSubmit)} disabled={loading || !isValid} loading={loading}>
          Create Challenge
        </Button>,
      ]}
    >
      {error && <Alert message={error} type="error" showIcon style={{ marginBottom: 16 }} />}
      <Form layout="vertical">
        <Controller name="opponentId" control={control} render={({ field }) => (
          <Form.Item label="Opponent" required validateStatus={errors.opponentId ? 'error' : ''} help={errors.opponentId?.message}>
            <Select {...field} placeholder="Select opponent" disabled={loading}>
              {availableOpponents.map(o => <Select.Option key={o.id} value={o.id}>{o.name}</Select.Option>)}
            </Select>
          </Form.Item>
        )} />
        <Controller name="eventId" control={control} render={({ field }) => (
          <Form.Item label="Event" required validateStatus={errors.eventId ? 'error' : ''} help={errors.eventId?.message}>
            <Select {...field} placeholder="Select event" disabled={loading}>
              {availableEvents.map(e => <Select.Option key={e.id} value={e.id}>{e.title}</Select.Option>)}
            </Select>
          </Form.Item>
        )} />
        <Controller name="message" control={control} render={({ field }) => (
          <Form.Item label="Message (Optional)" validateStatus={errors.message ? 'error' : ''} help={errors.message?.message}>
            <Input.TextArea {...field} rows={3} placeholder="Add a message to your challenge..." disabled={loading} maxLength={500} showCount />
          </Form.Item>
        )} />
      </Form>
    </Modal>
  )
}

export default ChallengeDialog
