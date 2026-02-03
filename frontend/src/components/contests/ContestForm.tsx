import React, { useState, useEffect } from 'react'
import { Modal, Form, Input, Select, Button, DatePicker, InputNumber, Collapse, message } from 'antd'
import { useForm, Controller } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { contestSchema, type ContestFormData } from '../../utils/validation'
import type { Contest } from '../../types/contest.types'
import { ScoringRulesEditor, type ContestRules } from './ScoringRulesEditor'
import { EventSelector } from './EventSelector'
import dayjs from 'dayjs'

interface ContestFormProps {
  open: boolean
  onClose: () => void
  onSubmit: (data: ContestFormData, eventIds?: number[]) => void
  contest?: Contest | null
  loading?: boolean
}

const sportTypes = [
  'Football',
  'Basketball',
  'Baseball',
  'Soccer',
  'Tennis',
  'Hockey',
  'Golf',
  'Boxing',
  'MMA',
  'Other',
]

export const ContestForm: React.FC<ContestFormProps> = ({
  open,
  onClose,
  onSubmit,
  contest,
  loading = false,
}) => {
  const isEditing = !!contest
  const [selectedEventIds, setSelectedEventIds] = useState<number[]>([])
  const [contestType, setContestType] = useState<string>('standard')
  
  // Track totalizator event_count from rules
  const [totalizatorEventCount, setTotalizatorEventCount] = useState<number>(15)

  const defaultValues = React.useMemo(() => {
    const now = new Date()
    const tomorrow = new Date(now.getTime() + 24 * 60 * 60 * 1000) // +1 day
    const nextWeek = new Date(now.getTime() + 7 * 24 * 60 * 60 * 1000) // +7 days
    
    return {
      title: contest?.title || '',
      description: contest?.description || '',
      sportType: contest?.sportType || '',
      rules: contest?.rules || '',
      startDate: contest?.startDate ? new Date(contest.startDate) : tomorrow,
      endDate: contest?.endDate ? new Date(contest.endDate) : nextWeek,
      maxParticipants: contest?.maxParticipants || 0,
    }
  }, [contest])

  const {
    control,
    handleSubmit,
    reset,
    formState: { errors, isValid },
  } = useForm<ContestFormData>({
    resolver: zodResolver(contestSchema),
    defaultValues,
    mode: 'onChange',
  })

  const handleClose = () => {
    reset()
    onClose()
  }

  React.useEffect(() => {
    reset(defaultValues)
  }, [defaultValues, reset])

  return (
    <Modal
      open={open}
      title={isEditing ? 'Edit Contest' : 'Create Contest'}
      onCancel={handleClose}
      footer={[
        <Button key="cancel" onClick={handleClose} disabled={loading}>Cancel</Button>,
        <Button 
          key="submit" 
          type="primary" 
          onClick={handleSubmit((data) => {
            // Pass event IDs for totalizator
            if (contestType === 'totalizator') {
              onSubmit(data, selectedEventIds)
            } else {
              onSubmit(data)
            }
          })} 
          disabled={loading || !isValid || (contestType === 'totalizator' && selectedEventIds.length < 5)} 
          loading={loading}
        >
          {isEditing ? 'Update' : 'Create'}
        </Button>,
      ]}
      width={700}
    >
      <Form layout="vertical">
        <Controller name="title" control={control} render={({ field }) => (
          <Form.Item label="Title" required validateStatus={errors.title ? 'error' : ''} help={errors.title?.message}>
            <Input {...field} disabled={loading} />
          </Form.Item>
        )} />
        <Controller name="description" control={control} render={({ field }) => (
          <Form.Item label="Description" validateStatus={errors.description ? 'error' : ''} help={errors.description?.message}>
            <Input.TextArea {...field} rows={3} disabled={loading} />
          </Form.Item>
        )} />
        <Controller name="sportType" control={control} render={({ field }) => (
          <Form.Item label="Sport Type" required validateStatus={errors.sportType ? 'error' : ''} help={errors.sportType?.message}>
            <Select {...field} disabled={loading}>
              {sportTypes.map(sport => <Select.Option key={sport} value={sport}>{sport}</Select.Option>)}
            </Select>
          </Form.Item>
        )} />
        <Controller name="rules" control={control} render={({ field }) => {
          // Parse rules JSON or use default
          let rulesValue: ContestRules | undefined
          try {
            if (field.value && typeof field.value === 'string') {
              rulesValue = JSON.parse(field.value)
            }
          } catch {
            // Invalid JSON, will use default
          }
          
          // Track contest type for showing EventSelector
          const currentType = rulesValue?.type || 'standard'
          if (currentType !== contestType) {
            setContestType(currentType)
          }
          
          // Track totalizator event_count
          if (rulesValue?.type === 'totalizator' && rulesValue.totalizator?.event_count !== totalizatorEventCount) {
            setTotalizatorEventCount(rulesValue.totalizator?.event_count || 15)
          }
          
          return (
            <Form.Item validateStatus={errors.rules ? 'error' : ''} help={errors.rules?.message}>
              <Collapse
                items={[{
                  key: 'scoring',
                  label: 'Правила подсчёта очков',
                  children: (
                    <ScoringRulesEditor
                      value={rulesValue}
                      onChange={(newRules) => {
                        field.onChange(JSON.stringify(newRules))
                        setContestType(newRules.type)
                        if (newRules.type === 'totalizator' && newRules.totalizator) {
                          setTotalizatorEventCount(newRules.totalizator.event_count)
                        }
                      }}
                    />
                  ),
                }]}
                defaultActiveKey={isEditing ? [] : ['scoring']}
              />
            </Form.Item>
          )
        }} />
        
        {/* Event Selector for Totalizator */}
        {contestType === 'totalizator' && (
          <Form.Item label="Выбор матчей для тотализатора" required>
            <EventSelector
              value={selectedEventIds}
              onChange={setSelectedEventIds}
              maxEvents={totalizatorEventCount}
              minEvents={Math.min(5, totalizatorEventCount)}
            />
          </Form.Item>
        )}
        <Form.Item label="Contest Period" required>
          <Input.Group compact>
            <Controller name="startDate" control={control} render={({ field }) => (
              <Form.Item validateStatus={errors.startDate ? 'error' : ''} help={errors.startDate?.message} style={{ display: 'inline-block', width: 'calc(50% - 4px)' }}>
                <DatePicker
                  {...field}
                  value={field.value ? dayjs(field.value) : null}
                  onChange={(date) => field.onChange(date?.toDate())}
                  showTime
                  format="YYYY-MM-DD HH:mm"
                  placeholder="Start Date"
                  disabled={loading}
                  style={{ width: '100%' }}
                />
              </Form.Item>
            )} />
            <Controller name="endDate" control={control} render={({ field }) => (
              <Form.Item validateStatus={errors.endDate ? 'error' : ''} help={errors.endDate?.message} style={{ display: 'inline-block', width: 'calc(50% - 4px)', marginLeft: 8 }}>
                <DatePicker
                  {...field}
                  value={field.value ? dayjs(field.value) : null}
                  onChange={(date) => field.onChange(date?.toDate())}
                  showTime
                  format="YYYY-MM-DD HH:mm"
                  placeholder="End Date"
                  disabled={loading}
                  style={{ width: '100%' }}
                />
              </Form.Item>
            )} />
          </Input.Group>
        </Form.Item>
        <Controller name="maxParticipants" control={control} render={({ field }) => (
          <Form.Item label="Max Participants" validateStatus={errors.maxParticipants ? 'error' : ''} help={errors.maxParticipants?.message}>
            <InputNumber {...field} min={0} disabled={loading} style={{ width: '100%' }} />
          </Form.Item>
        )} />
      </Form>
    </Modal>
  )
}

export default ContestForm
