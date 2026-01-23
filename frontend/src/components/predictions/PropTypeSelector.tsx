import React from 'react'
import { Card, Radio, Input, Button, Space, Tag, Typography } from 'antd'
import { PlusOutlined, DeleteOutlined } from '@ant-design/icons'
import type { PropType } from '../../types/props.types'
import type { PropPredictionFormData } from '../../utils/prediction-validation'

const { Text } = Typography

interface PropTypeSelectorProps {
  propTypes: PropType[]
  selectedProps: PropPredictionFormData[]
  onPropsChange: (props: PropPredictionFormData[]) => void
  homeTeam: string
  awayTeam: string
  disabled?: boolean
}

export const PropTypeSelector: React.FC<PropTypeSelectorProps> = ({
  propTypes,
  selectedProps,
  onPropsChange,
  homeTeam,
  awayTeam,
  disabled = false,
}) => {
  const addProp = (propType: PropType) => {
    const newProp: PropPredictionFormData = {
      propTypeId: propType.id,
      propSlug: propType.slug,
      selection: '',
      pointsValue: 2,
    }
    onPropsChange([...selectedProps, newProp])
  }

  const removeProp = (index: number) => {
    onPropsChange(selectedProps.filter((_, i) => i !== index))
  }

  const updateProp = (index: number, selection: string) => {
    const updated = [...selectedProps]
    updated[index] = { ...updated[index], selection }
    onPropsChange(updated)
  }

  const availableProps = propTypes.filter(
    pt => !selectedProps.some(sp => sp.propTypeId === pt.id)
  )

  return (
    <Space direction="vertical" size="middle" style={{ width: '100%' }}>
      <Text strong>Props Predictions (Optional)</Text>
      
      {selectedProps.map((prop, index) => {
        const propType = propTypes.find(pt => pt.id === prop.propTypeId)
        if (!propType) return null

        return (
          <Card key={index} size="small">
            <Space direction="vertical" style={{ width: '100%' }}>
              <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
                <Tag color="blue">{propType.name}</Tag>
                <Button
                  danger
                  size="small"
                  icon={<DeleteOutlined />}
                  onClick={() => removeProp(index)}
                  disabled={disabled}
                />
              </div>
              
              {propType.valueType === 'exact_value' ? (
                <Input
                  type="number"
                  placeholder="Enter value"
                  value={prop.selection}
                  onChange={e => updateProp(index, e.target.value)}
                  disabled={disabled}
                />
              ) : propType.valueType === 'yes_no' ? (
                <Radio.Group
                  value={prop.selection}
                  onChange={e => updateProp(index, e.target.value)}
                  disabled={disabled}
                >
                  <Radio value="yes">Yes</Radio>
                  <Radio value="no">No</Radio>
                </Radio.Group>
              ) : propType.valueType === 'team_select' ? (
                <Radio.Group
                  value={prop.selection}
                  onChange={e => updateProp(index, e.target.value)}
                  disabled={disabled}
                >
                  <Radio value="home">{homeTeam}</Radio>
                  <Radio value="away">{awayTeam}</Radio>
                </Radio.Group>
              ) : (
                <Input
                  placeholder="Enter value"
                  value={prop.selection}
                  onChange={e => updateProp(index, e.target.value)}
                  disabled={disabled}
                />
              )}
            </Space>
          </Card>
        )
      })}

      {availableProps.length > 0 && (
        <Space wrap>
          {availableProps.map(propType => (
            <Button
              key={propType.id}
              icon={<PlusOutlined />}
              onClick={() => addProp(propType)}
              disabled={disabled}
              size="small"
            >
              Add {propType.name}
            </Button>
          ))}
        </Space>
      )}
    </Space>
  )
}

export default PropTypeSelector
