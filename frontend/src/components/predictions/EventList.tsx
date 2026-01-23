import React, { useState } from 'react'
import { Row, Col, Select, Space, Spin, Pagination, Alert, Typography } from 'antd'
import { useEvents } from '../../hooks/use-predictions'
import { useSports } from '../../hooks/use-sports'
import EventCard from './EventCard'
import type { Event, Prediction } from '../../types/prediction.types'
import { DEFAULT_EVENT_PAGE_SIZE, MAX_PARTICIPANTS_DISPLAY } from '../../utils/constants'

const { Title } = Typography

interface EventListProps {
  contestId: number
  onPredict: (event: Event) => void
  userPredictions: Prediction[]
}

const statusOptions = ['', 'scheduled', 'live', 'completed', 'cancelled']

export const EventList: React.FC<EventListProps> = ({
  contestId,
  onPredict,
  userPredictions,
}) => {
  const [sportType, setSportType] = useState('')
  const [status, setStatus] = useState('scheduled')
  const [page, setPage] = useState(1)
  const pageSize = DEFAULT_EVENT_PAGE_SIZE

  const { data: sportsData } = useSports({ pagination: { page: 1, limit: MAX_PARTICIPANTS_DISPLAY } })
  const sports = sportsData?.sports || []

  const { data, isLoading, isError, error } = useEvents({
    sportType: sportType || undefined,
    status: status || undefined,
    pagination: { page, limit: pageSize },
  })

  const predictionMap = new Map(userPredictions.map(p => [p.eventId, p]))

  if (isError) {
    return <Alert message="Error" description={error?.message} type="error" showIcon />
  }

  return (
    <Space direction="vertical" size="large" style={{ width: '100%' }}>
      <Space wrap>
        <Select
          style={{ width: 150 }}
          placeholder="Sport Type"
          value={sportType}
          onChange={setSportType}
          allowClear
        >
          <Select.Option value="">All Sports</Select.Option>
          {sports.map(s => <Select.Option key={s.id} value={s.slug}>{s.name}</Select.Option>)}
        </Select>
        <Select
          style={{ width: 150 }}
          placeholder="Status"
          value={status}
          onChange={setStatus}
        >
          {statusOptions.map(s => <Select.Option key={s} value={s}>{s || 'All'}</Select.Option>)}
        </Select>
      </Space>

      {isLoading ? (
        <div style={{ textAlign: 'center', padding: '32px 0' }}>
          <Spin size="large" />
        </div>
      ) : (
        <>
          <Row gutter={[16, 16]}>
            {data?.events?.map(event => (
              <Col key={event.id} xs={24} sm={12} md={8} lg={6}>
                <EventCard
                  event={event}
                  onPredict={onPredict}
                  existingPrediction={predictionMap.get(event.id)}
                />
              </Col>
            ))}
          </Row>
          {data?.pagination && (
            <Pagination
              current={page}
              pageSize={pageSize}
              total={data.pagination.total}
              onChange={setPage}
              showSizeChanger={false}
            />
          )}
        </>
      )}
    </Space>
  )
}

export default EventList
