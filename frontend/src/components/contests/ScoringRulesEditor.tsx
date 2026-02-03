import React, { useEffect } from 'react'
import { Card, Form, InputNumber, Radio, Space, Typography, Divider, List, Switch, Spin, Alert } from 'antd'
import { TrophyOutlined, ThunderboltOutlined, DollarOutlined, TeamOutlined } from '@ant-design/icons'
import { useRiskyEventTypes } from '../../hooks/use-risky-events'

const { Text } = Typography

// Types for scoring rules
export interface StandardScoringRules {
  exact_score: number
  goal_difference: number
  correct_outcome: number
  outcome_plus_team_goals: number
  any_other: number
}

export interface RiskyEvent {
  slug: string
  name: string
  name_en?: string
  points: number
  enabled?: boolean
}

export interface RiskyScoringRules {
  max_selections: number
  events: RiskyEvent[]
}

export interface TotalizatorRules {
  event_count: number
  scoring: StandardScoringRules
}

export interface RelayRules {
  team_size: number
  event_count: number
  scoring: StandardScoringRules
  allow_reassign: boolean
}

export interface ContestRules {
  type: 'standard' | 'risky' | 'totalizator' | 'relay'
  scoring?: StandardScoringRules
  risky?: RiskyScoringRules
  totalizator?: TotalizatorRules
  relay?: RelayRules
}

interface ScoringRulesEditorProps {
  value?: ContestRules
  onChange?: (rules: ContestRules) => void
}

// Default values
const defaultStandardRules: StandardScoringRules = {
  exact_score: 5,
  goal_difference: 3,
  correct_outcome: 1,
  outcome_plus_team_goals: 1,
  any_other: 4,
}

// Fallback events if API fails
const fallbackRiskyEvents: RiskyEvent[] = [
  { slug: 'penalty', name: 'Будет пенальти', name_en: 'Penalty awarded', points: 3, enabled: true },
  { slug: 'red_card', name: 'Будет удаление', name_en: 'Red card shown', points: 4, enabled: true },
  { slug: 'own_goal', name: 'Будет автогол', name_en: 'Own goal scored', points: 5, enabled: true },
  { slug: 'hat_trick', name: 'Будет хет-трик', name_en: 'Hat-trick scored', points: 6, enabled: true },
  { slug: 'clean_sheet_home', name: 'Хозяева на ноль', name_en: 'Home clean sheet', points: 2, enabled: true },
]

export const ScoringRulesEditor: React.FC<ScoringRulesEditorProps> = ({
  value,
  onChange,
}) => {
  // Load risky event types from API
  const { data: apiEventTypes, isLoading: loadingEvents } = useRiskyEventTypes({ activeOnly: true })

  // Convert API event types to RiskyEvent format
  const availableRiskyEvents: RiskyEvent[] = apiEventTypes?.map(et => ({
    slug: et.slug,
    name: et.name,
    name_en: et.nameEn,
    points: et.defaultPoints,
    enabled: true,
  })) || fallbackRiskyEvents

  // Initialize with defaults
  const rules: ContestRules = value || {
    type: 'standard',
    scoring: { ...defaultStandardRules },
  }

  const handleTypeChange = (type: 'standard' | 'risky' | 'totalizator' | 'relay') => {
    const newRules: ContestRules = { type }
    if (type === 'standard') {
      newRules.scoring = rules.scoring || { ...defaultStandardRules }
    } else if (type === 'risky') {
      newRules.risky = rules.risky || {
        max_selections: 5,
        events: availableRiskyEvents.map(e => ({ ...e })),
      }
    } else if (type === 'totalizator') {
      newRules.totalizator = rules.totalizator || {
        event_count: 15,
        scoring: { ...defaultStandardRules },
      }
    } else if (type === 'relay') {
      newRules.relay = rules.relay || {
        team_size: 5,
        event_count: 15,
        scoring: { ...defaultStandardRules },
        allow_reassign: true,
      }
    }
    onChange?.(newRules)
  }

  const handleStandardChange = (field: keyof StandardScoringRules, val: number | null) => {
    if (rules.type !== 'standard') return
    const newScoring = { ...(rules.scoring || defaultStandardRules), [field]: val || 0 }
    onChange?.({ ...rules, scoring: newScoring })
  }

  const handleRiskyMaxChange = (val: number | null) => {
    if (rules.type !== 'risky' || !rules.risky) return
    onChange?.({ ...rules, risky: { ...rules.risky, max_selections: val || 5 } })
  }

  const handleRiskyEventChange = (slug: string, field: 'points' | 'enabled', val: number | boolean) => {
    if (rules.type !== 'risky' || !rules.risky) return
    const newEvents = rules.risky.events.map(e =>
      e.slug === slug ? { ...e, [field]: val } : e
    )
    onChange?.({ ...rules, risky: { ...rules.risky, events: newEvents } })
  }

  const handleTotalizatorEventCountChange = (val: number | null) => {
    if (rules.type !== 'totalizator' || !rules.totalizator) return
    onChange?.({ ...rules, totalizator: { ...rules.totalizator, event_count: val || 15 } })
  }

  const handleTotalizatorScoringChange = (field: keyof StandardScoringRules, val: number | null) => {
    if (rules.type !== 'totalizator' || !rules.totalizator) return
    const newScoring = { ...rules.totalizator.scoring, [field]: val || 0 }
    onChange?.({ ...rules, totalizator: { ...rules.totalizator, scoring: newScoring } })
  }

  const handleRelayTeamSizeChange = (val: number | null) => {
    if (rules.type !== 'relay' || !rules.relay) return
    onChange?.({ ...rules, relay: { ...rules.relay, team_size: val || 5 } })
  }

  const handleRelayEventCountChange = (val: number | null) => {
    if (rules.type !== 'relay' || !rules.relay) return
    onChange?.({ ...rules, relay: { ...rules.relay, event_count: val || 15 } })
  }

  const handleRelayAllowReassignChange = (val: boolean) => {
    if (rules.type !== 'relay' || !rules.relay) return
    onChange?.({ ...rules, relay: { ...rules.relay, allow_reassign: val } })
  }

  const handleRelayScoringChange = (field: keyof StandardScoringRules, val: number | null) => {
    if (rules.type !== 'relay' || !rules.relay) return
    const newScoring = { ...rules.relay.scoring, [field]: val || 0 }
    onChange?.({ ...rules, relay: { ...rules.relay, scoring: newScoring } })
  }

  return (
    <Card title="Правила подсчёта очков" size="small">
      <Form layout="vertical">
        <Form.Item label="Тип конкурса">
          <Radio.Group
            value={rules.type}
            onChange={(e) => handleTypeChange(e.target.value)}
            optionType="button"
            buttonStyle="solid"
          >
            <Radio.Button value="standard">
              <TrophyOutlined /> Обычный
            </Radio.Button>
            <Radio.Button value="risky">
              <ThunderboltOutlined /> Рисковый
            </Radio.Button>
            <Radio.Button value="totalizator">
              <DollarOutlined /> Тотализатор
            </Radio.Button>
            <Radio.Button value="relay">
              <TeamOutlined /> Эстафета
            </Radio.Button>
          </Radio.Group>
        </Form.Item>

        {rules.type === 'standard' && (
          <>
            <Divider orientation="left">Очки за прогноз</Divider>
            <Space direction="vertical" style={{ width: '100%' }}>
              <Form.Item label="Точный счёт" style={{ marginBottom: 8 }}>
                <InputNumber
                  min={0}
                  max={100}
                  value={rules.scoring?.exact_score}
                  onChange={(v) => handleStandardChange('exact_score', v)}
                  addonAfter="очков"
                />
              </Form.Item>
              <Form.Item label="Разница мячей" style={{ marginBottom: 8 }}>
                <InputNumber
                  min={0}
                  max={100}
                  value={rules.scoring?.goal_difference}
                  onChange={(v) => handleStandardChange('goal_difference', v)}
                  addonAfter="очков"
                />
              </Form.Item>
              <Form.Item label="Верный исход" style={{ marginBottom: 8 }}>
                <InputNumber
                  min={0}
                  max={100}
                  value={rules.scoring?.correct_outcome}
                  onChange={(v) => handleStandardChange('correct_outcome', v)}
                  addonAfter="очков"
                />
              </Form.Item>
              <Form.Item label="Исход + голы команды (бонус)" style={{ marginBottom: 8 }}>
                <InputNumber
                  min={0}
                  max={100}
                  value={rules.scoring?.outcome_plus_team_goals}
                  onChange={(v) => handleStandardChange('outcome_plus_team_goals', v)}
                  addonAfter="очков"
                />
              </Form.Item>
              <Form.Item label='Прогноз "Другой счёт"' style={{ marginBottom: 8 }}>
                <InputNumber
                  min={0}
                  max={100}
                  value={rules.scoring?.any_other}
                  onChange={(v) => handleStandardChange('any_other', v)}
                  addonAfter="очков"
                />
              </Form.Item>
            </Space>
          </>
        )}

        {rules.type === 'risky' && rules.risky && (
          <>
            <Divider orientation="left">Настройки рискового конкурса</Divider>
            <Form.Item label="Макс. выборов на матч">
              <InputNumber
                min={1}
                max={10}
                value={rules.risky.max_selections}
                onChange={handleRiskyMaxChange}
              />
            </Form.Item>
            
            <Text type="secondary" style={{ display: 'block', marginBottom: 12 }}>
              Угадал событие → +очки, не угадал → −очки
            </Text>

            {loadingEvents ? (
              <div style={{ textAlign: 'center', padding: 20 }}>
                <Spin tip="Загрузка событий..." />
              </div>
            ) : (
              <>
                {apiEventTypes && apiEventTypes.length > 0 && (
                  <Alert 
                    message={`Доступно ${apiEventTypes.length} типов событий`}
                    type="info"
                    showIcon
                    style={{ marginBottom: 12 }}
                  />
                )}
                <List
                  size="small"
                  dataSource={rules.risky.events}
                  renderItem={(event) => (
                    <List.Item
                      actions={[
                        <InputNumber
                          key="points"
                          min={0.5}
                          max={20}
                          step={0.5}
                          value={event.points}
                          onChange={(v) => handleRiskyEventChange(event.slug, 'points', v || 1)}
                          size="small"
                          style={{ width: 80 }}
                        />,
                        <Switch
                          key="enabled"
                          size="small"
                          checked={event.enabled !== false}
                          onChange={(v) => handleRiskyEventChange(event.slug, 'enabled', v)}
                        />,
                      ]}
                    >
                      <List.Item.Meta
                        title={event.name}
                        description={event.name_en}
                      />
                    </List.Item>
                  )}
                />
              </>
            )}
          </>
        )}

        {rules.type === 'totalizator' && rules.totalizator && (
          <>
            <Divider orientation="left">Настройки тотализатора</Divider>
            <Text type="secondary" style={{ display: 'block', marginBottom: 12 }}>
              Выбери матчи из разных лиг. Очки считаются как в обычном конкурсе.
            </Text>
            
            <Form.Item label="Количество матчей">
              <InputNumber
                min={5}
                max={30}
                value={rules.totalizator.event_count}
                onChange={handleTotalizatorEventCountChange}
                addonAfter="матчей"
              />
            </Form.Item>

            <Divider orientation="left">Очки за прогноз</Divider>
            <Space direction="vertical" style={{ width: '100%' }}>
              <Form.Item label="Точный счёт" style={{ marginBottom: 8 }}>
                <InputNumber
                  min={0}
                  max={100}
                  value={rules.totalizator.scoring.exact_score}
                  onChange={(v) => handleTotalizatorScoringChange('exact_score', v)}
                  addonAfter="очков"
                />
              </Form.Item>
              <Form.Item label="Разница мячей" style={{ marginBottom: 8 }}>
                <InputNumber
                  min={0}
                  max={100}
                  value={rules.totalizator.scoring.goal_difference}
                  onChange={(v) => handleTotalizatorScoringChange('goal_difference', v)}
                  addonAfter="очков"
                />
              </Form.Item>
              <Form.Item label="Верный исход" style={{ marginBottom: 8 }}>
                <InputNumber
                  min={0}
                  max={100}
                  value={rules.totalizator.scoring.correct_outcome}
                  onChange={(v) => handleTotalizatorScoringChange('correct_outcome', v)}
                  addonAfter="очков"
                />
              </Form.Item>
              <Form.Item label="Исход + голы команды (бонус)" style={{ marginBottom: 8 }}>
                <InputNumber
                  min={0}
                  max={100}
                  value={rules.totalizator.scoring.outcome_plus_team_goals}
                  onChange={(v) => handleTotalizatorScoringChange('outcome_plus_team_goals', v)}
                  addonAfter="очков"
                />
              </Form.Item>
              <Form.Item label='Прогноз "Другой счёт"' style={{ marginBottom: 8 }}>
                <InputNumber
                  min={0}
                  max={100}
                  value={rules.totalizator.scoring.any_other}
                  onChange={(v) => handleTotalizatorScoringChange('any_other', v)}
                  addonAfter="очков"
                />
              </Form.Item>
            </Space>
          </>
        )}

        {rules.type === 'relay' && rules.relay && (
          <>
            <Divider orientation="left">Настройки эстафеты</Divider>
            <Text type="secondary" style={{ display: 'block', marginBottom: 12 }}>
              Командный конкурс. Капитан распределяет матчи между участниками команды.
            </Text>
            
            <Space direction="vertical" style={{ width: '100%' }}>
              <Form.Item label="Участников в команде">
                <InputNumber
                  min={2}
                  max={10}
                  value={rules.relay.team_size}
                  onChange={handleRelayTeamSizeChange}
                  addonAfter="человек"
                />
              </Form.Item>
              
              <Form.Item label="Количество матчей">
                <InputNumber
                  min={5}
                  max={50}
                  value={rules.relay.event_count}
                  onChange={handleRelayEventCountChange}
                  addonAfter="матчей"
                />
              </Form.Item>

              <Form.Item label="Переназначение матчей">
                <Switch
                  checked={rules.relay.allow_reassign}
                  onChange={handleRelayAllowReassignChange}
                  checkedChildren="Разрешено"
                  unCheckedChildren="Запрещено"
                />
                <Text type="secondary" style={{ marginLeft: 8 }}>
                  Капитан может менять распределение до начала
                </Text>
              </Form.Item>
            </Space>

            <Divider orientation="left">Очки за прогноз</Divider>
            <Space direction="vertical" style={{ width: '100%' }}>
              <Form.Item label="Точный счёт" style={{ marginBottom: 8 }}>
                <InputNumber
                  min={0}
                  max={100}
                  value={rules.relay.scoring.exact_score}
                  onChange={(v) => handleRelayScoringChange('exact_score', v)}
                  addonAfter="очков"
                />
              </Form.Item>
              <Form.Item label="Разница мячей" style={{ marginBottom: 8 }}>
                <InputNumber
                  min={0}
                  max={100}
                  value={rules.relay.scoring.goal_difference}
                  onChange={(v) => handleRelayScoringChange('goal_difference', v)}
                  addonAfter="очков"
                />
              </Form.Item>
              <Form.Item label="Верный исход" style={{ marginBottom: 8 }}>
                <InputNumber
                  min={0}
                  max={100}
                  value={rules.relay.scoring.correct_outcome}
                  onChange={(v) => handleRelayScoringChange('correct_outcome', v)}
                  addonAfter="очков"
                />
              </Form.Item>
              <Form.Item label="Исход + голы команды (бонус)" style={{ marginBottom: 8 }}>
                <InputNumber
                  min={0}
                  max={100}
                  value={rules.relay.scoring.outcome_plus_team_goals}
                  onChange={(v) => handleRelayScoringChange('outcome_plus_team_goals', v)}
                  addonAfter="очков"
                />
              </Form.Item>
              <Form.Item label='Прогноз "Другой счёт"' style={{ marginBottom: 8 }}>
                <InputNumber
                  min={0}
                  max={100}
                  value={rules.relay.scoring.any_other}
                  onChange={(v) => handleRelayScoringChange('any_other', v)}
                  addonAfter="очков"
                />
              </Form.Item>
            </Space>
          </>
        )}
      </Form>
    </Card>
  )
}

export default ScoringRulesEditor
