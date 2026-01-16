import React, { useMemo, useState } from 'react'
import {
  MaterialReactTable,
  useMaterialReactTable,
  type MRT_ColumnDef,
} from 'material-react-table'
import {
  Box,
  IconButton,
  Tooltip,
  Chip,
  Typography,
} from '@mui/material'
import { Edit as EditIcon, Delete as DeleteIcon } from '@mui/icons-material'
import { useUserPredictions, useDeletePrediction } from '../../hooks/use-predictions'
import type { Prediction } from '../../types/prediction.types'
import { formatDate, formatRelativeTime } from '../../utils/date-utils'

interface PredictionListProps {
  contestId: number
  onEdit: (prediction: Prediction) => void
}

const getStatusColor = (status: string) => {
  switch (status) {
    case 'pending': return 'warning'
    case 'scored': return 'success'
    case 'cancelled': return 'error'
    default: return 'default'
  }
}

const formatPredictionData = (data: string): string => {
  try {
    const parsed = JSON.parse(data)
    const parts: string[] = []
    if (parsed.winner) {
      parts.push(`Winner: ${parsed.winner}`)
    }
    if (parsed.homeScore !== undefined && parsed.awayScore !== undefined) {
      parts.push(`Score: ${parsed.homeScore}-${parsed.awayScore}`)
    }
    return parts.join(', ') || data
  } catch {
    return data
  }
}

export const PredictionList: React.FC<PredictionListProps> = ({
  contestId,
  onEdit,
}) => {
  const [pagination, setPagination] = useState({
    pageIndex: 0,
    pageSize: 10,
  })

  const { data, isLoading, isError, error } = useUserPredictions({
    contestId,
    pagination: {
      page: pagination.pageIndex + 1,
      limit: pagination.pageSize,
    },
  })

  const deletePredictionMutation = useDeletePrediction()

  const handleDelete = (prediction: Prediction) => {
    if (window.confirm('Are you sure you want to delete this prediction?')) {
      deletePredictionMutation.mutate(prediction.id)
    }
  }

  const columns = useMemo<MRT_ColumnDef<Prediction>[]>(
    () => [
      {
        accessorKey: 'eventId',
        header: 'Event ID',
        size: 80,
      },
      {
        accessorKey: 'predictionData',
        header: 'Prediction',
        size: 200,
        Cell: ({ cell }) => (
          <Typography variant="body2">
            {formatPredictionData(cell.getValue<string>())}
          </Typography>
        ),
      },
      {
        accessorKey: 'status',
        header: 'Status',
        size: 100,
        Cell: ({ cell }) => (
          <Chip
            label={cell.getValue<string>()}
            color={getStatusColor(cell.getValue<string>())}
            size="small"
          />
        ),
      },
      {
        accessorKey: 'submittedAt',
        header: 'Submitted',
        size: 150,
        Cell: ({ cell }) => formatRelativeTime(cell.getValue<string>()),
      },
      {
        accessorKey: 'createdAt',
        header: 'Created',
        size: 150,
        Cell: ({ cell }) => formatDate(cell.getValue<string>()),
      },
    ],
    []
  )

  const table = useMaterialReactTable({
    columns,
    data: data?.predictions ?? [],
    enableRowSelection: false,
    enableColumnOrdering: true,
    enableGlobalFilter: true,
    enableSorting: true,
    enablePagination: true,
    manualPagination: true,
    rowCount: data?.pagination?.total ?? 0,
    onPaginationChange: setPagination,
    state: {
      isLoading,
      pagination,
      showAlertBanner: isError,
      showProgressBars: isLoading,
    },
    muiToolbarAlertBannerProps: isError
      ? {
          color: 'error',
          children: `Error loading predictions: ${error?.message}`,
        }
      : undefined,
    renderRowActions: ({ row }) => (
      <Box sx={{ display: 'flex', gap: '0.5rem' }}>
        {row.original.status === 'pending' && (
          <>
            <Tooltip title="Edit Prediction">
              <IconButton
                color="primary"
                onClick={() => onEdit(row.original)}
              >
                <EditIcon />
              </IconButton>
            </Tooltip>
            <Tooltip title="Delete Prediction">
              <IconButton
                color="error"
                onClick={() => handleDelete(row.original)}
                disabled={deletePredictionMutation.isPending}
              >
                <DeleteIcon />
              </IconButton>
            </Tooltip>
          </>
        )}
      </Box>
    ),
    enableRowActions: true,
    positionActionsColumn: 'last',
    muiTableContainerProps: {
      sx: { minHeight: '400px' },
    },
  })

  return <MaterialReactTable table={table} />
}

export default PredictionList
