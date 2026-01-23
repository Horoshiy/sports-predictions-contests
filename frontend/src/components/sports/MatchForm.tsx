import React, { useMemo } from 'react'
import { Modal, Form, Input, Select, Button, DatePicker, InputNumber } from 'antd'
import { useForm, Controller } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { matchSchema, type MatchFormData } from '../../utils/sports-validation'
import { useLeagues, useTeams } from '../../hooks/use-sports'
import type { Match } from '../../types/sports.types'
import dayjs from 'dayjs'

interface MatchFormProps {
  open: boolean
  onClose: () => void
  onSubmit: (data: MatchFormData) => void
  match?: Match | null
  loading?: boolean
}

export const MatchForm: React.FC<MatchFormProps> = ({ open, onClose, onSubmit, match, loading = false }) => {
  const isEditing = !!match
  const { data: leaguesData } = useLeagues({ pagination: { page: 1, limit: 100 }, activeOnly: true })
  const { data: teamsData } = useTeams({ pagination: { page: 1, limit: 200 }, activeOnly: true })

  const defaultValues = React.useMemo(() => ({
    leagueId: match?.leagueId || 0,
    homeTeamId: match?.homeTeamId || 0,
    awayTeamId: match?.awayTeamId || 0,
    scheduledAt: match?.scheduledAt ? new Date(match.scheduledAt) : new Date(),
    status: match?.status || 'scheduled',
    homeScore: match?.homeScore || 0,
    awayScore: match?.awayScore || 0,
    resultData: match?.resultData || '',
  }), [match])

  const { control, handleSubmit, reset, watch, setValue, formState: { errors, isValid } } = useForm<MatchFormData>({
    resolver: zodResolver(matchSchema),
    defaultValues,
    mode: 'onChange',
  })

  const selectedLeagueId = watch('leagueId')
  const selectedLeague = leaguesData?.leagues?.find(l => l.id === selectedLeagueId)

  const filteredTeams = useMemo(() => {
    if (!selectedLeague || !teamsData?.teams) return []
    return teamsData.teams.filter(t => t.sportId === selectedLeague.sportId)
  }, [selectedLeague, teamsData])

  const prevLeagueId = React.useRef(selectedLeagueId)
  React.useEffect(() => {
    if (!isEditing && prevLeagueId.current !== selectedLeagueId && prevLeagueId.current !== 0) {
      setValue('homeTeamId', 0)
      setValue('awayTeamId', 0)
    }
    prevLeagueId.current = selectedLeagueId
  }, [selectedLeagueId, isEditing, setValue])

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
      title={isEditing ? 'Edit Match' : 'Create Match'}
      onCancel={handleClose}
      footer={[
        <Button key="cancel" onClick={handleClose} disabled={loading}>Cancel</Button>,
        <Button key="submit" type="primary" onClick={handleSubmit(onSubmit)} disabled={loading || !isValid} loading={loading}>
          {isEditing ? 'Update' : 'Create'}
        </Button>,
      ]}
      width={600}
    >
      <Form layout="vertical">
        <Controller name="leagueId" control={control} render={({ field }) => (
          <Form.Item label="League" required validateStatus={errors.leagueId ? 'error' : ''} help={errors.leagueId?.message}>
            <Select {...field} disabled={loading}>
              {leaguesData?.leagues?.map(l => <Select.Option key={l.id} value={l.id}>{l.name}</Select.Option>)}
            </Select>
          </Form.Item>
        )} />
        <Form.Item label="Teams" required>
          <Input.Group compact>
            <Controller name="homeTeamId" control={control} render={({ field }) => (
              <Form.Item validateStatus={errors.homeTeamId ? 'error' : ''} help={errors.homeTeamId?.message} style={{ display: 'inline-block', width: 'calc(50% - 4px)' }}>
                <Select {...field} placeholder="Home Team" disabled={loading || !selectedLeagueId}>
                  {filteredTeams.map(t => <Select.Option key={t.id} value={t.id}>{t.name}</Select.Option>)}
                </Select>
              </Form.Item>
            )} />
            <Controller name="awayTeamId" control={control} render={({ field }) => (
              <Form.Item validateStatus={errors.awayTeamId ? 'error' : ''} help={errors.awayTeamId?.message} style={{ display: 'inline-block', width: 'calc(50% - 4px)', marginLeft: 8 }}>
                <Select {...field} placeholder="Away Team" disabled={loading || !selectedLeagueId}>
                  {filteredTeams.map(t => <Select.Option key={t.id} value={t.id}>{t.name}</Select.Option>)}
                </Select>
              </Form.Item>
            )} />
          </Input.Group>
        </Form.Item>
        <Controller name="scheduledAt" control={control} render={({ field }) => (
          <Form.Item label="Scheduled At" required validateStatus={errors.scheduledAt ? 'error' : ''} help={errors.scheduledAt?.message}>
            <DatePicker
              {...field}
              value={field.value ? dayjs(field.value) : null}
              onChange={(date) => field.onChange(date?.toDate())}
              showTime
              format="YYYY-MM-DD HH:mm"
              disabled={loading}
              style={{ width: '100%' }}
            />
          </Form.Item>
        )} />
        {isEditing && (
          <>
            <Controller name="status" control={control} render={({ field }) => (
              <Form.Item label="Status">
                <Select {...field} disabled={loading}>
                  <Select.Option value="scheduled">Scheduled</Select.Option>
                  <Select.Option value="live">Live</Select.Option>
                  <Select.Option value="finished">Finished</Select.Option>
                  <Select.Option value="cancelled">Cancelled</Select.Option>
                  <Select.Option value="postponed">Postponed</Select.Option>
                </Select>
              </Form.Item>
            )} />
            <Form.Item label="Score">
              <Input.Group compact>
                <Controller name="homeScore" control={control} render={({ field }) => (
                  <Form.Item style={{ display: 'inline-block', width: 'calc(50% - 4px)' }}>
                    <InputNumber {...field} placeholder="Home" disabled={loading} style={{ width: '100%' }} />
                  </Form.Item>
                )} />
                <Controller name="awayScore" control={control} render={({ field }) => (
                  <Form.Item style={{ display: 'inline-block', width: 'calc(50% - 4px)', marginLeft: 8 }}>
                    <InputNumber {...field} placeholder="Away" disabled={loading} style={{ width: '100%' }} />
                  </Form.Item>
                )} />
              </Input.Group>
            </Form.Item>
          </>
        )}
      </Form>
    </Modal>
  )
}

export default MatchForm
