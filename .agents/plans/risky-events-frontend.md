# Feature: Risky Events Frontend

## Feature Description

Создание фронтенд-компонентов для управления рисковыми событиями:
1. **RiskyEventTypesManager** — CRUD глобальных типов событий (админка)
2. **ContestRiskyEventsEditor** — выбор событий из API при создании конкурса
3. **MatchRiskyEventsEditor** — переопределение событий для конкретного матча

## User Story

**As an** administrator  
**I want** to manage risky event types from a web interface  
**So that** I can add/edit/disable events without deployment

**As an** administrator  
**I want** to select which risky events are available in a contest  
**So that** I can customize each contest's event list

**As an** administrator  
**I want** to override event points for specific matches  
**So that** I can adjust for derbies, cup finals, etc.

---

## CONTEXT REFERENCES

### Existing Patterns to Follow

**Services:** `frontend/src/services/sports-service.ts`
- gRPC client wrapper
- CRUD methods pattern

**Hooks:** `frontend/src/hooks/use-sports.ts`
- React Query hooks
- Query keys pattern
- useMutation with toast

**Types:** `frontend/src/types/sports.types.ts`
- Request/Response types
- Entity types

**Pages:** `frontend/src/pages/SportsPage.tsx`
- Tab-based admin interface
- Form modal pattern
- List + CRUD pattern

**Components:** `frontend/src/components/sports/SportList.tsx`
- Ant Design Table
- Edit/Delete actions

### Backend API (already implemented)

Endpoints (prediction-service):
- `GET /v1/risky-event-types` — list all event types
- `POST /v1/risky-event-types` — create event type (admin)
- `PUT /v1/risky-event-types/{id}` — update event type (admin)
- `DELETE /v1/risky-event-types/{id}` — delete event type (admin)
- `GET /v1/events/{id}/risky-events?contest_id=X` — get match risky events
- `PUT /v1/events/{id}/risky-events` — set match overrides

---

## IMPLEMENTATION PLAN

### Phase 1: Types & Service (Foundation) ✅

**Task 1:** ✅ Create TypeScript types for risky events
**Task 2:** ✅ Create risky-events-service.ts

### Phase 2: React Query Hooks ✅

**Task 3:** ✅ Create use-risky-events.ts hooks

### Phase 3: Admin Component (RiskyEventTypesManager) ✅

**Task 4:** ✅ Create RiskyEventTypesManager.tsx
**Task 5:** ✅ Create RiskyEventTypeForm.tsx
**Task 6:** ✅ Add "Risky Events" tab to SportsPage.tsx

### Phase 4: Contest Integration ✅

**Task 7:** ✅ Update ScoringRulesEditor.tsx to load events from API
**Task 8:** ✅ Create ContestRiskyEventsSelector.tsx component

### Phase 5: Match Override Component ✅

**Task 9:** ✅ Create MatchRiskyEventsEditor.tsx
**Task 10:** ✅ Integrate via ContestEventsManager (+ Events tab in ContestsPage)

---

## STEP-BY-STEP TASKS

### Task 1: CREATE Types

**File:** `frontend/src/types/risky-events.types.ts`

```typescript
// Risky Event Type (global)
export interface RiskyEventType {
  id: number
  slug: string
  name: string
  nameEn?: string
  description?: string
  defaultPoints: number
  sportType: string
  category: string
  icon?: string
  isActive: boolean
  sortOrder: number
  createdAt?: string
  updatedAt?: string
}

// Match Risky Event (with overrides)
export interface MatchRiskyEvent {
  riskyEventTypeId: number
  slug: string
  name: string
  nameEn?: string
  points: number  // final points (with overrides)
  icon?: string
  isEnabled: boolean
  outcome?: boolean | null  // null=pending, true=happened, false=didn't
}

// Request types
export interface ListRiskyEventTypesRequest {
  sportType?: string
  activeOnly?: boolean
}

export interface CreateRiskyEventTypeRequest {
  slug: string
  name: string
  nameEn?: string
  description?: string
  defaultPoints: number
  sportType?: string
  category?: string
  icon?: string
  sortOrder?: number
}

export interface UpdateRiskyEventTypeRequest extends Partial<CreateRiskyEventTypeRequest> {
  id: number
  isActive?: boolean
}

export interface GetMatchRiskyEventsRequest {
  eventId: number
  contestId: number
}

export interface SetMatchRiskyEventsRequest {
  eventId: number
  events: {
    riskyEventTypeId: number
    points: number
    isEnabled: boolean
  }[]
}

export interface SetMatchRiskyEventOutcomeRequest {
  eventId: number
  riskyEventTypeId: number
  outcome: boolean
}

// Response types
export interface RiskyEventTypeResponse {
  eventType: RiskyEventType
}

export interface ListRiskyEventTypesResponse {
  eventTypes: RiskyEventType[]
}

export interface GetMatchRiskyEventsResponse {
  events: MatchRiskyEvent[]
  maxSelections: number
}
```

**VALIDATE:** `npm run typecheck`

---

### Task 2: CREATE Service

**File:** `frontend/src/services/risky-events-service.ts`

```typescript
import grpcClient from './grpc-client'
import type {
  RiskyEventType,
  MatchRiskyEvent,
  ListRiskyEventTypesRequest,
  CreateRiskyEventTypeRequest,
  UpdateRiskyEventTypeRequest,
  GetMatchRiskyEventsRequest,
  SetMatchRiskyEventsRequest,
  SetMatchRiskyEventOutcomeRequest,
  RiskyEventTypeResponse,
  ListRiskyEventTypesResponse,
  GetMatchRiskyEventsResponse,
} from '../types/risky-events.types'

class RiskyEventsService {
  // Global Event Types
  async listEventTypes(request: ListRiskyEventTypesRequest = {}): Promise<RiskyEventType[]> {
    const params = new URLSearchParams()
    if (request.sportType) params.append('sport_type', request.sportType)
    if (request.activeOnly) params.append('active_only', 'true')
    const url = params.toString() ? `/v1/risky-event-types?${params}` : '/v1/risky-event-types'
    const response = await grpcClient.get<ListRiskyEventTypesResponse>(url)
    return response.eventTypes || []
  }

  async createEventType(request: CreateRiskyEventTypeRequest): Promise<RiskyEventType> {
    const response = await grpcClient.post<RiskyEventTypeResponse>('/v1/risky-event-types', request)
    return response.eventType
  }

  async updateEventType(request: UpdateRiskyEventTypeRequest): Promise<RiskyEventType> {
    const { id, ...data } = request
    const response = await grpcClient.put<RiskyEventTypeResponse>(`/v1/risky-event-types/${id}`, data)
    return response.eventType
  }

  async deleteEventType(id: number): Promise<void> {
    await grpcClient.delete(`/v1/risky-event-types/${id}`)
  }

  // Match Events
  async getMatchEvents(request: GetMatchRiskyEventsRequest): Promise<GetMatchRiskyEventsResponse> {
    const params = new URLSearchParams()
    params.append('contest_id', request.contestId.toString())
    const response = await grpcClient.get<GetMatchRiskyEventsResponse>(
      `/v1/events/${request.eventId}/risky-events?${params}`
    )
    return response
  }

  async setMatchEvents(request: SetMatchRiskyEventsRequest): Promise<void> {
    const { eventId, events } = request
    await grpcClient.put(`/v1/events/${eventId}/risky-events`, { events })
  }

  async setMatchEventOutcome(request: SetMatchRiskyEventOutcomeRequest): Promise<void> {
    const { eventId, riskyEventTypeId, outcome } = request
    await grpcClient.put(`/v1/events/${eventId}/risky-events/${riskyEventTypeId}/outcome`, { outcome })
  }
}

export const riskyEventsService = new RiskyEventsService()
```

**VALIDATE:** `npm run typecheck`

---

### Task 3: CREATE Hooks

**File:** `frontend/src/hooks/use-risky-events.ts`

```typescript
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { riskyEventsService } from '../services/risky-events-service'
import { useToast } from '../contexts/ToastContext'
import type {
  ListRiskyEventTypesRequest,
  CreateRiskyEventTypeRequest,
  UpdateRiskyEventTypeRequest,
  GetMatchRiskyEventsRequest,
  SetMatchRiskyEventsRequest,
  SetMatchRiskyEventOutcomeRequest,
} from '../types/risky-events.types'

// Query keys
export const riskyEventTypesKeys = {
  all: ['riskyEventTypes'] as const,
  lists: () => [...riskyEventTypesKeys.all, 'list'] as const,
  list: (filters: ListRiskyEventTypesRequest) => [...riskyEventTypesKeys.lists(), filters] as const,
}

export const matchRiskyEventsKeys = {
  all: ['matchRiskyEvents'] as const,
  match: (eventId: number, contestId: number) => [...matchRiskyEventsKeys.all, eventId, contestId] as const,
}

// Hooks for global event types
export const useRiskyEventTypes = (request: ListRiskyEventTypesRequest = {}) => {
  return useQuery({
    queryKey: riskyEventTypesKeys.list(request),
    queryFn: () => riskyEventsService.listEventTypes(request),
    staleTime: 5 * 60 * 1000,
  })
}

export const useCreateRiskyEventType = () => {
  const queryClient = useQueryClient()
  const { showToast } = useToast()
  return useMutation({
    mutationFn: (request: CreateRiskyEventTypeRequest) => riskyEventsService.createEventType(request),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: riskyEventTypesKeys.lists() })
      showToast('Событие создано!', 'success')
    },
    onError: (error: Error) => showToast(`Ошибка: ${error.message}`, 'error'),
  })
}

export const useUpdateRiskyEventType = () => {
  const queryClient = useQueryClient()
  const { showToast } = useToast()
  return useMutation({
    mutationFn: (request: UpdateRiskyEventTypeRequest) => riskyEventsService.updateEventType(request),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: riskyEventTypesKeys.lists() })
      showToast('Событие обновлено!', 'success')
    },
    onError: (error: Error) => showToast(`Ошибка: ${error.message}`, 'error'),
  })
}

export const useDeleteRiskyEventType = () => {
  const queryClient = useQueryClient()
  const { showToast } = useToast()
  return useMutation({
    mutationFn: (id: number) => riskyEventsService.deleteEventType(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: riskyEventTypesKeys.lists() })
      showToast('Событие удалено!', 'success')
    },
    onError: (error: Error) => showToast(`Ошибка: ${error.message}`, 'error'),
  })
}

// Hooks for match events
export const useMatchRiskyEvents = (eventId: number, contestId: number, enabled = true) => {
  return useQuery({
    queryKey: matchRiskyEventsKeys.match(eventId, contestId),
    queryFn: () => riskyEventsService.getMatchEvents({ eventId, contestId }),
    enabled: enabled && !!eventId && !!contestId,
  })
}

export const useSetMatchRiskyEvents = () => {
  const queryClient = useQueryClient()
  const { showToast } = useToast()
  return useMutation({
    mutationFn: (request: SetMatchRiskyEventsRequest) => riskyEventsService.setMatchEvents(request),
    onSuccess: (_, variables) => {
      queryClient.invalidateQueries({ queryKey: matchRiskyEventsKeys.all })
      showToast('События матча обновлены!', 'success')
    },
    onError: (error: Error) => showToast(`Ошибка: ${error.message}`, 'error'),
  })
}

export const useSetMatchEventOutcome = () => {
  const queryClient = useQueryClient()
  const { showToast } = useToast()
  return useMutation({
    mutationFn: (request: SetMatchRiskyEventOutcomeRequest) => riskyEventsService.setMatchEventOutcome(request),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: matchRiskyEventsKeys.all })
      showToast('Результат события сохранён!', 'success')
    },
    onError: (error: Error) => showToast(`Ошибка: ${error.message}`, 'error'),
  })
}
```

**VALIDATE:** `npm run typecheck`

---

### Task 4: CREATE RiskyEventTypesManager

**File:** `frontend/src/components/admin/RiskyEventTypesManager.tsx`

CRUD компонент для глобальных событий:
- Ant Design Table со списком событий
- Колонки: Icon, Name, Category, Points, Active, Actions
- Кнопка "Добавить событие"
- Edit/Delete actions в каждой строке
- Drag-n-drop сортировка (опционально)

```typescript
import React, { useState } from 'react'
import { Table, Button, Space, Tag, Switch, Popconfirm, Typography } from 'antd'
import { PlusOutlined, EditOutlined, DeleteOutlined } from '@ant-design/icons'
import { useRiskyEventTypes, useUpdateRiskyEventType, useDeleteRiskyEventType } from '../../hooks/use-risky-events'
import RiskyEventTypeForm from './RiskyEventTypeForm'
import type { RiskyEventType } from '../../types/risky-events.types'

const { Text } = Typography

// ... implementation
```

**VALIDATE:** `npm run build`

---

### Task 5: CREATE RiskyEventTypeForm

**File:** `frontend/src/components/admin/RiskyEventTypeForm.tsx`

Модальная форма создания/редактирования:
- Fields: slug, name, nameEn, description, defaultPoints, category, icon, sortOrder
- Category select: goals, cards, defense, totals, halves, timing, special
- Icon picker (emoji selector или input)
- Валидация: required name, slug, points >= 0

```typescript
import React, { useEffect } from 'react'
import { Modal, Form, Input, InputNumber, Select, Space } from 'antd'
import { useCreateRiskyEventType, useUpdateRiskyEventType } from '../../hooks/use-risky-events'
import type { RiskyEventType, CreateRiskyEventTypeRequest } from '../../types/risky-events.types'

// ... implementation
```

**VALIDATE:** `npm run build`

---

### Task 6: ADD Tab to SportsPage

**File:** `frontend/src/pages/SportsPage.tsx`

Добавить пятую вкладку "Risky Events":

```typescript
import RiskyEventTypesManager from '../components/admin/RiskyEventTypesManager'

// В Tabs.items добавить:
{
  key: 'risky-events',
  label: 'Risky Events',
  children: <RiskyEventTypesManager />,
}
```

**VALIDATE:** `npm run build`

---

### Task 7: UPDATE ScoringRulesEditor

**File:** `frontend/src/components/contests/ScoringRulesEditor.tsx`

Изменить для risky типа:
1. Загружать события из API вместо `defaultRiskyEvents`
2. Показать чекбоксы для выбора событий
3. InputNumber для очков каждого выбранного

```typescript
// Добавить хук
const { data: eventTypes, isLoading: loadingEvents } = useRiskyEventTypes({ activeOnly: true })

// Заменить defaultRiskyEvents на eventTypes
// При создании нового risky конкурса - использовать eventTypes с дефолтными очками
```

**VALIDATE:** `npm run build`

---

### Task 8: CREATE ContestRiskyEventsSelector

**File:** `frontend/src/components/contests/ContestRiskyEventsSelector.tsx`

Компонент выбора событий для конкурса:
- Checkbox.Group со всеми доступными событиями
- Для каждого выбранного - InputNumber для очков
- Лимит max_selections (показывать warning если превышен)
- Категории событий как Collapse или Tags

```typescript
import React from 'react'
import { Card, Checkbox, InputNumber, Space, Typography, Tag, Alert, Spin } from 'antd'
import { useRiskyEventTypes } from '../../hooks/use-risky-events'
import type { RiskyEvent } from './ScoringRulesEditor'

interface ContestRiskyEventsSelectorProps {
  value?: RiskyEvent[]
  onChange?: (events: RiskyEvent[]) => void
  maxSelections?: number
}

// ... implementation
```

**VALIDATE:** `npm run build`

---

### Task 9: CREATE MatchRiskyEventsEditor

**File:** `frontend/src/components/events/MatchRiskyEventsEditor.tsx`

Редактор переопределений для матча:
- Показать события конкурса
- Для каждого: InputNumber очков, Switch включено
- После матча: Switch "произошло" для записи outcome

```typescript
import React from 'react'
import { Card, List, InputNumber, Switch, Space, Typography, Tag, Button } from 'antd'
import { useMatchRiskyEvents, useSetMatchRiskyEvents, useSetMatchEventOutcome } from '../../hooks/use-risky-events'
import type { MatchRiskyEvent } from '../../types/risky-events.types'

interface MatchRiskyEventsEditorProps {
  eventId: number
  contestId: number
  isMatchFinished?: boolean  // показывать outcome editors
}

// ... implementation
```

**VALIDATE:** `npm run build`

---

### Task 10: INTEGRATE into MatchForm

**File:** `frontend/src/components/sports/MatchForm.tsx`

Добавить секцию Risky Events:
- Показать только если матч в risky конкурсе
- Collapse секция с MatchRiskyEventsEditor

```typescript
// Добавить проверку типа конкурса
// Если risky - показать MatchRiskyEventsEditor
```

**VALIDATE:** `npm run build`

---

## TESTING STRATEGY

### Manual Testing Checklist

- [ ] Создать новый тип события в админке
- [ ] Отредактировать существующее событие
- [ ] Деактивировать событие (оно не должно показываться при создании конкурса)
- [ ] Создать risky конкурс, выбрав 5 событий из списка
- [ ] Изменить очки события для конкурса
- [ ] Переопределить очки для конкретного матча
- [ ] После матча отметить outcomes

### E2E Tests (Playwright)

- `tests/e2e/risky-events-admin.spec.ts` — CRUD в админке
- `tests/e2e/risky-events-contest.spec.ts` — создание risky конкурса

---

## ACCEPTANCE CRITERIA

- [x] Админ может создавать/редактировать/удалять типы рисковых событий
- [x] При создании risky конкурса события загружаются из API
- [x] Можно выбрать до 10 событий и настроить очки каждого
- [x] Для матча можно переопределить очки выбранных событий
- [x] После матча можно отметить какие события произошли

**Note:** Requires gateway routing fix for `/v1/risky-event-types` endpoints

---

## FILES TO CREATE

1. `frontend/src/types/risky-events.types.ts`
2. `frontend/src/services/risky-events-service.ts`
3. `frontend/src/hooks/use-risky-events.ts`
4. `frontend/src/components/admin/RiskyEventTypesManager.tsx`
5. `frontend/src/components/admin/RiskyEventTypeForm.tsx`
6. `frontend/src/components/contests/ContestRiskyEventsSelector.tsx`
7. `frontend/src/components/events/MatchRiskyEventsEditor.tsx`

## FILES TO MODIFY

1. `frontend/src/pages/SportsPage.tsx` — добавить таб
2. `frontend/src/components/contests/ScoringRulesEditor.tsx` — загрузка из API
3. `frontend/src/components/sports/MatchForm.tsx` — интеграция редактора

---

## ESTIMATED TIME

- Phase 1 (Types & Service): 30 min
- Phase 2 (Hooks): 20 min
- Phase 3 (Admin): 1.5 hours
- Phase 4 (Contest): 1 hour
- Phase 5 (Match): 1 hour

**Total: ~4-5 hours**

---

## PRIORITY ORDER

1. **Task 1-3**: Foundation (без них ничего не работает)
2. **Task 7**: ScoringRulesEditor (критично для создания конкурсов)
3. **Task 4-6**: Admin UI (управление событиями)
4. **Task 8-10**: Match overrides (опционально, можно отложить)
