import React, { useState } from 'react'
import { Table, Button, Space, Tag, Switch, Popconfirm, Typography, Card, Input } from 'antd'
import { PlusOutlined, EditOutlined, DeleteOutlined, SearchOutlined } from '@ant-design/icons'
import type { ColumnsType } from 'antd/es/table'
import { 
  useRiskyEventTypes, 
  useToggleRiskyEventType, 
  useDeleteRiskyEventType 
} from '../../hooks/use-risky-events'
import { RISKY_EVENT_CATEGORIES } from '../../types/risky-events.types'
import type { RiskyEventType } from '../../types/risky-events.types'
import RiskyEventTypeForm from './RiskyEventTypeForm'

const { Text } = Typography

const categoryColors: Record<string, string> = {
  goals: 'green',
  cards: 'red',
  defense: 'blue',
  totals: 'purple',
  halves: 'orange',
  timing: 'cyan',
  special: 'magenta',
  general: 'default',
}

const RiskyEventTypesManager: React.FC = () => {
  const [formOpen, setFormOpen] = useState(false)
  const [selectedEvent, setSelectedEvent] = useState<RiskyEventType | null>(null)
  const [searchText, setSearchText] = useState('')

  const { data: eventTypes, isLoading } = useRiskyEventTypes()
  const toggleMutation = useToggleRiskyEventType()
  const deleteMutation = useDeleteRiskyEventType()

  const handleCreate = () => {
    setSelectedEvent(null)
    setFormOpen(true)
  }

  const handleEdit = (record: RiskyEventType) => {
    setSelectedEvent(record)
    setFormOpen(true)
  }

  const handleDelete = async (id: number) => {
    await deleteMutation.mutateAsync(id)
  }

  const handleToggle = async (record: RiskyEventType) => {
    await toggleMutation.mutateAsync({ 
      id: record.id, 
      isActive: !record.isActive 
    })
  }

  const filteredData = eventTypes?.filter(e => 
    e.name.toLowerCase().includes(searchText.toLowerCase()) ||
    e.slug.toLowerCase().includes(searchText.toLowerCase()) ||
    e.nameEn?.toLowerCase().includes(searchText.toLowerCase())
  )

  const columns: ColumnsType<RiskyEventType> = [
    {
      title: '',
      dataIndex: 'icon',
      key: 'icon',
      width: 50,
      render: (icon: string) => (
        <span style={{ fontSize: '1.2em' }}>{icon || '❓'}</span>
      ),
    },
    {
      title: 'Название',
      dataIndex: 'name',
      key: 'name',
      render: (name: string, record) => (
        <Space direction="vertical" size={0}>
          <Text strong>{name}</Text>
          {record.nameEn && (
            <Text type="secondary" style={{ fontSize: '0.85em' }}>
              {record.nameEn}
            </Text>
          )}
        </Space>
      ),
    },
    {
      title: 'Slug',
      dataIndex: 'slug',
      key: 'slug',
      width: 150,
      render: (slug: string) => (
        <Text code style={{ fontSize: '0.85em' }}>{slug}</Text>
      ),
    },
    {
      title: 'Категория',
      dataIndex: 'category',
      key: 'category',
      width: 120,
      filters: RISKY_EVENT_CATEGORIES.map(c => ({ text: c.label, value: c.value })),
      onFilter: (value, record) => record.category === value,
      render: (category: string) => {
        const cat = RISKY_EVENT_CATEGORIES.find(c => c.value === category)
        return (
          <Tag color={categoryColors[category] || 'default'}>
            {cat?.icon} {cat?.label || category}
          </Tag>
        )
      },
    },
    {
      title: 'Очки',
      dataIndex: 'defaultPoints',
      key: 'defaultPoints',
      width: 80,
      align: 'center',
      sorter: (a, b) => a.defaultPoints - b.defaultPoints,
      render: (points: number) => (
        <Text strong style={{ color: points >= 5 ? '#f5222d' : points >= 3 ? '#fa8c16' : '#52c41a' }}>
          {points}
        </Text>
      ),
    },
    {
      title: 'Порядок',
      dataIndex: 'sortOrder',
      key: 'sortOrder',
      width: 80,
      align: 'center',
      sorter: (a, b) => a.sortOrder - b.sortOrder,
      defaultSortOrder: 'ascend',
    },
    {
      title: 'Активно',
      dataIndex: 'isActive',
      key: 'isActive',
      width: 90,
      align: 'center',
      filters: [
        { text: 'Активные', value: true },
        { text: 'Неактивные', value: false },
      ],
      onFilter: (value, record) => record.isActive === value,
      render: (isActive: boolean, record) => (
        <Switch
          checked={isActive}
          onChange={() => handleToggle(record)}
          loading={toggleMutation.isPending}
          size="small"
        />
      ),
    },
    {
      title: 'Действия',
      key: 'actions',
      width: 100,
      render: (_, record) => (
        <Space size="small">
          <Button
            type="text"
            icon={<EditOutlined />}
            onClick={() => handleEdit(record)}
            size="small"
          />
          <Popconfirm
            title="Удалить событие?"
            description="Это действие нельзя отменить"
            onConfirm={() => handleDelete(record.id)}
            okText="Удалить"
            cancelText="Отмена"
            okButtonProps={{ danger: true }}
          >
            <Button
              type="text"
              danger
              icon={<DeleteOutlined />}
              size="small"
            />
          </Popconfirm>
        </Space>
      ),
    },
  ]

  return (
    <>
      <Card
        title="Типы рисковых событий"
        extra={
          <Space>
            <Input
              placeholder="Поиск..."
              prefix={<SearchOutlined />}
              value={searchText}
              onChange={(e) => setSearchText(e.target.value)}
              style={{ width: 200 }}
              allowClear
            />
            <Button
              type="primary"
              icon={<PlusOutlined />}
              onClick={handleCreate}
            >
              Добавить
            </Button>
          </Space>
        }
      >
        <Table
          columns={columns}
          dataSource={filteredData}
          rowKey="id"
          loading={isLoading}
          pagination={{ pageSize: 20, showSizeChanger: false }}
          size="middle"
          rowClassName={(record) => !record.isActive ? 'ant-table-row-disabled' : ''}
        />
      </Card>

      <RiskyEventTypeForm
        open={formOpen}
        onClose={() => setFormOpen(false)}
        eventType={selectedEvent}
      />

      <style>{`
        .ant-table-row-disabled {
          opacity: 0.5;
        }
      `}</style>
    </>
  )
}

export default RiskyEventTypesManager
