import React from 'react'
import { Card, Form, InputNumber, Radio, Space, Typography, Divider, List, Switch } from 'antd'
import { TrophyOutlined, ThunderboltOutlined } from '@ant-design/icons'

const { Title, Text } = Typography

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

export interface ContestRules {
  type: 'standard' | 'risky'
  scoring?: StandardScoringRules
  risky?: RiskyScoringRules
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

const defaultRiskyEvents: RiskyEvent[] = [
  { slug: 'penalty', name: 'Будет пенальти', name_en: 'Penalty awarded', points: 3, enabled: true },
  { slug: 'red_card', name: 'Будет удаление', name_en: 'Red card shown', points: 4, enabled: true },
  { slug: 'own_goal', name: 'Будет автогол', name_en: 'Own goal scored', points: 5, enabled: true },
  { slug: 'hat_trick', name: 'Будет хет-трик', name_en: 'Hat-trick scored', points: 6, enabled: true },
  { slug: 'clean_sheet_home', name: 'Хозяева на ноль', name_en: 'Home clean sheet', points: 2, enabled: true },
  { slug: 'clean_sheet_away', name: 'Гости на ноль', name_en: 'Away clean sheet', points: 3, enabled: true },
  { slug: 'both_teams_score', name: 'Обе забьют', name_en: 'Both teams score', points: 2, enabled: true },
  { slug: 'over_3_goals', name: 'Больше 3 голов', name_en: 'Over 3.5 goals', points: 2, enabled: true },
  { slug: 'first_half_draw', name: 'Ничья в 1-м тайме', name_en: 'First half draw', points: 2, enabled: true },
  { slug: 'comeback', name: 'Камбэк', name_en: 'Comeback from 0:2+', points: 7, enabled: true },
]

export const ScoringRulesEditor: React.FC<ScoringRulesEditorProps> = ({
  value,
  onChange,
}) => {
  // Initialize with defaults
  const rules: ContestRules = value || {
    type: 'standard',
    scoring: { ...defaultStandardRules },
  }

  const handleTypeChange = (type: 'standard' | 'risky') => {
    const newRules: ContestRules = { type }
    if (type === 'standard') {
      newRules.scoring = rules.scoring || { ...defaultStandardRules }
    } else {
      newRules.risky = rules.risky || {
        max_selections: 5,
        events: defaultRiskyEvents.map(e => ({ ...e })),
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

            <List
              size="small"
              dataSource={rules.risky.events}
              renderItem={(event) => (
                <List.Item
                  actions={[
                    <InputNumber
                      key="points"
                      min={1}
                      max={20}
                      value={event.points}
                      onChange={(v) => handleRiskyEventChange(event.slug, 'points', v || 1)}
                      size="small"
                      style={{ width: 70 }}
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
      </Form>
    </Card>
  )
}

export default ScoringRulesEditor
