import React from 'react'
import { Card, Progress, List, Tag, Button, Space, Typography } from 'antd'
import { CheckCircleOutlined, CloseCircleOutlined, TrophyOutlined } from '@ant-design/icons'

const { Title, Text } = Typography

interface ProfileCompletionData {
  percentage: number
  missingFields: string[]
  suggestions: string[]
}

interface ProfileCompletionProps {
  completion: ProfileCompletionData
  onFieldClick?: (field: string) => void
}

const fieldIcons: Record<string, React.ReactNode> = {
  avatar: '📷',
  bio: '✍️',
  location: '📍',
  website: '🌐',
  email: '📧',
}

export const ProfileCompletion: React.FC<ProfileCompletionProps> = ({
  completion,
  onFieldClick,
}) => {
  const getProgressColor = (percentage: number) => {
    if (percentage >= 80) return '#52c41a'
    if (percentage >= 50) return '#faad14'
    return '#ff4d4f'
  }

  return (
    <Card>
      <Space direction="vertical" size="large" style={{ width: '100%' }}>
        <div>
          <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: 8 }}>
            <Title level={5} style={{ margin: 0 }}>Profile Completion</Title>
            <Tag color={getProgressColor(completion.percentage)}>
              {completion.percentage}%
            </Tag>
          </div>
          <Progress
            percent={completion.percentage}
            strokeColor={getProgressColor(completion.percentage)}
            showInfo={false}
          />
        </div>

        {completion.missingFields.length > 0 && (
          <div>
            <Text strong>Missing Fields:</Text>
            <List
              size="small"
              dataSource={completion.missingFields}
              renderItem={(field) => (
                <List.Item
                  actions={[
                    onFieldClick && (
                      <Button size="small" type="link" onClick={() => onFieldClick(field)}>
                        Add
                      </Button>
                    ),
                  ]}
                >
                  <Space>
                    <CloseCircleOutlined style={{ color: '#ff4d4f' }} />
                    <Text>{fieldIcons[field]} {field}</Text>
                  </Space>
                </List.Item>
              )}
            />
          </div>
        )}

        {completion.suggestions.length > 0 && (
          <div>
            <Text strong>Suggestions:</Text>
            <List
              size="small"
              dataSource={completion.suggestions}
              renderItem={(suggestion) => (
                <List.Item>
                  <Space>
                    <TrophyOutlined style={{ color: '#1890ff' }} />
                    <Text type="secondary">{suggestion}</Text>
                  </Space>
                </List.Item>
              )}
            />
          </div>
        )}

        {completion.percentage === 100 && (
          <div style={{ textAlign: 'center', padding: '16px 0' }}>
            <CheckCircleOutlined style={{ fontSize: 48, color: '#52c41a', marginBottom: 8 }} />
            <br />
            <Text strong>Your profile is complete!</Text>
          </div>
        )}
      </Space>
    </Card>
  )
}

export default ProfileCompletion
