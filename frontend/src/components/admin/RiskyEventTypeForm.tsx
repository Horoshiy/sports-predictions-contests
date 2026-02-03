import React, { useEffect } from 'react'
import { Modal, Form, Input, InputNumber, Select, Space, Row, Col } from 'antd'
import { useCreateRiskyEventType, useUpdateRiskyEventType } from '../../hooks/use-risky-events'
import { RISKY_EVENT_CATEGORIES } from '../../types/risky-events.types'
import type { RiskyEventType, CreateRiskyEventTypeRequest } from '../../types/risky-events.types'

interface RiskyEventTypeFormProps {
  open: boolean
  onClose: () => void
  eventType?: RiskyEventType | null
}

interface FormValues {
  slug: string
  name: string
  nameEn?: string
  description?: string
  defaultPoints: number
  category: string
  icon?: string
  sortOrder: number
}

// Common emoji icons for quick selection
const COMMON_ICONS = ['⚽', '🟥', '🟨', '🎩', '🏠', '✈️', '📈', '📉', '⏰', '🤝', '🔄', '📺', '0️⃣', '1️⃣', '🛡️', '✨']

const RiskyEventTypeForm: React.FC<RiskyEventTypeFormProps> = ({
  open,
  onClose,
  eventType,
}) => {
  const [form] = Form.useForm<FormValues>()
  const createMutation = useCreateRiskyEventType()
  const updateMutation = useUpdateRiskyEventType()

  const isEditing = !!eventType

  useEffect(() => {
    if (open && eventType) {
      form.setFieldsValue({
        slug: eventType.slug,
        name: eventType.name,
        nameEn: eventType.nameEn,
        description: eventType.description,
        defaultPoints: eventType.defaultPoints,
        category: eventType.category,
        icon: eventType.icon,
        sortOrder: eventType.sortOrder,
      })
    } else if (open) {
      form.resetFields()
      form.setFieldsValue({
        defaultPoints: 2,
        category: 'general',
        sortOrder: 0,
      })
    }
  }, [open, eventType, form])

  const handleSubmit = async () => {
    try {
      const values = await form.validateFields()

      if (isEditing && eventType) {
        await updateMutation.mutateAsync({
          id: eventType.id,
          ...values,
        })
      } else {
        await createMutation.mutateAsync(values as CreateRiskyEventTypeRequest)
      }

      onClose()
    } catch (error) {
      // Form validation failed
      console.error('Form validation failed:', error)
    }
  }

  const generateSlug = () => {
    const name = form.getFieldValue('name')
    if (name) {
      const slug = name
        .toLowerCase()
        .replace(/[а-яё]/g, (char: string) => {
          const map: Record<string, string> = {
            'а': 'a', 'б': 'b', 'в': 'v', 'г': 'g', 'д': 'd', 'е': 'e', 'ё': 'e',
            'ж': 'zh', 'з': 'z', 'и': 'i', 'й': 'y', 'к': 'k', 'л': 'l', 'м': 'm',
            'н': 'n', 'о': 'o', 'п': 'p', 'р': 'r', 'с': 's', 'т': 't', 'у': 'u',
            'ф': 'f', 'х': 'h', 'ц': 'ts', 'ч': 'ch', 'ш': 'sh', 'щ': 'sch',
            'ъ': '', 'ы': 'y', 'ь': '', 'э': 'e', 'ю': 'yu', 'я': 'ya',
          }
          return map[char] || char
        })
        .replace(/[^a-z0-9]+/g, '_')
        .replace(/^_|_$/g, '')
      form.setFieldValue('slug', slug)
    }
  }

  return (
    <Modal
      title={isEditing ? 'Редактировать событие' : 'Новое событие'}
      open={open}
      onOk={handleSubmit}
      onCancel={onClose}
      okText={isEditing ? 'Сохранить' : 'Создать'}
      cancelText="Отмена"
      confirmLoading={createMutation.isPending || updateMutation.isPending}
      width={600}
    >
      <Form
        form={form}
        layout="vertical"
        requiredMark="optional"
      >
        <Row gutter={16}>
          <Col span={16}>
            <Form.Item
              name="name"
              label="Название"
              rules={[{ required: true, message: 'Введите название' }]}
            >
              <Input 
                placeholder="Будет пенальти" 
                onBlur={generateSlug}
              />
            </Form.Item>
          </Col>
          <Col span={8}>
            <Form.Item
              name="icon"
              label="Иконка"
            >
              <Select
                placeholder="⚽"
                allowClear
                showSearch
                options={COMMON_ICONS.map(icon => ({ value: icon, label: icon }))}
              />
            </Form.Item>
          </Col>
        </Row>

        <Row gutter={16}>
          <Col span={12}>
            <Form.Item
              name="slug"
              label="Slug (ID)"
              rules={[
                { required: true, message: 'Введите slug' },
                { pattern: /^[a-z0-9_]+$/, message: 'Только a-z, 0-9, _' },
              ]}
            >
              <Input placeholder="penalty" disabled={isEditing} />
            </Form.Item>
          </Col>
          <Col span={12}>
            <Form.Item
              name="nameEn"
              label="Название (EN)"
            >
              <Input placeholder="Penalty awarded" />
            </Form.Item>
          </Col>
        </Row>

        <Row gutter={16}>
          <Col span={8}>
            <Form.Item
              name="defaultPoints"
              label="Очки по умолчанию"
              rules={[{ required: true, message: 'Введите очки' }]}
            >
              <InputNumber
                min={0.5}
                max={20}
                step={0.5}
                style={{ width: '100%' }}
              />
            </Form.Item>
          </Col>
          <Col span={8}>
            <Form.Item
              name="category"
              label="Категория"
              rules={[{ required: true, message: 'Выберите категорию' }]}
            >
              <Select
                placeholder="Выберите"
                options={RISKY_EVENT_CATEGORIES.map(c => ({
                  value: c.value,
                  label: `${c.icon} ${c.label}`,
                }))}
              />
            </Form.Item>
          </Col>
          <Col span={8}>
            <Form.Item
              name="sortOrder"
              label="Порядок"
            >
              <InputNumber min={0} max={100} style={{ width: '100%' }} />
            </Form.Item>
          </Col>
        </Row>

        <Form.Item
          name="description"
          label="Описание"
        >
          <Input.TextArea
            rows={2}
            placeholder="Дополнительное описание события..."
          />
        </Form.Item>
      </Form>
    </Modal>
  )
}

export default RiskyEventTypeForm
